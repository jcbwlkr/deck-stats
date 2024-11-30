package moxfield

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/jcbwlkr/deck-stats/internal/domains/magic"
	"github.com/jcbwlkr/deck-stats/internal/services"
)

type Client struct {
	url       string
	userAgent string
	client    *http.Client

	gate chan struct{}
}

func NewClient(userAgent string, sleep time.Duration) *Client {

	// This semaphore goroutine acts as a traffic cop. Expensive operations
	// receive from this channel to get permission to do work. This inifinite
	// loop will push a value, wait for someone to take it, sleep a bit, then
	// start over. This is a simple way to ensure we don't abuse the api.
	gate := make(chan struct{})
	go func() {
		for {
			gate <- struct{}{}
			time.Sleep(sleep)
		}
	}()

	return &Client{
		url:       "https://api2.moxfield.com",
		userAgent: userAgent,
		client:    &http.Client{},
		gate:      gate,
	}
}

func (c *Client) ListMyDecks(ctx context.Context, token string) ([]magic.Deck, error) {
	url := fmt.Sprintf("%s/v3/decks", c.url)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("User-Agent", c.userAgent)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		//b, _ := io.ReadAll(resp.Body)
		//fmt.Println(string(b[0:min(512, len(b)-1)]))
		return nil, fmt.Errorf("moxfield: api status %s", resp.Status)
	}

	var data struct {
		Decks []deckListDeck `json:"decks"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	decks := make([]magic.Deck, 0, len(data.Decks))
	for _, deck := range data.Decks {
		d := magic.Deck{
			Name:      deck.Name,
			Service:   services.Moxfield,
			ServiceID: deck.PublicID,
			URL:       deck.PublicURL,
			Folder: magic.Folder{
				ID:   deck.Folder.ID,
				Name: deck.Folder.Name,
			},
			UpdatedAt: deck.LastUpdatedAtUtc,
		}
		for _, c := range deck.ColorIdentity {
			d.ColorIdentity = append(d.ColorIdentity, magic.Color(strings.ToLower(c)))
		}

		decks = append(decks, d)
	}

	return decks, nil
}

func (c *Client) AddDeckDetails(ctx context.Context, token string, d *magic.Deck) error {

	// Block until we have permission to call or our context is canceled.
	select {
	case <-c.gate:
	case <-ctx.Done():
	}

	url := fmt.Sprintf("%s/v3/decks/all/%s", c.url, d.ServiceID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("User-Agent", c.userAgent)

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("moxfield: api status %s", resp.Status)
	}

	var data deck
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}

	d.Archetypes = []magic.Archetype{}
	for _, hub := range data.Hubs {
		d.Archetypes = append(d.Archetypes, magic.Archetype{
			Name:        hub.Name,
			Description: hub.Description,
		})
	}

	if data.Format == "oathbreaker" {
		for _, card := range data.Boards.Commanders.Cards {
			d.Leaders.Oathbreakers = append(d.Leaders.Oathbreakers, magic.Card{
				ID:   card.Card.ScryfallID,
				Name: card.Card.Name,
			})
		}
		for _, card := range data.Boards.SignatureSpells.Cards {
			d.Leaders.SignatureSpells = append(d.Leaders.SignatureSpells, magic.Card{
				ID:   card.Card.ScryfallID,
				Name: card.Card.Name,
			})
		}
	} else {
		for _, card := range data.Boards.Commanders.Cards {
			d.Leaders.Commanders = append(d.Leaders.Commanders, magic.Card{
				ID:   card.Card.ScryfallID,
				Name: card.Card.Name,
			})
		}
	}

	for _, card := range data.Boards.Companions.Cards {
		d.Leaders.Companion = &magic.Card{
			ID:   card.Card.ScryfallID,
			Name: card.Card.Name,
		}
	}

	return nil
}

type deckListDeck struct {
	ID                   string `json:"id"`
	Name                 string `json:"name"`
	HasPrimer            bool   `json:"hasPrimer"`
	Format               string `json:"format"`
	AreCommentsEnabled   bool   `json:"areCommentsEnabled"`
	Visibility           string `json:"visibility"`
	PublicURL            string `json:"publicUrl"`
	PublicID             string `json:"publicId"`
	LikeCount            int    `json:"likeCount"`
	ViewCount            int    `json:"viewCount"`
	CommentCount         int    `json:"commentCount"`
	SfwCommentCount      int    `json:"sfwCommentCount"`
	IsLegal              bool   `json:"isLegal"`
	AuthorsCanEdit       bool   `json:"authorsCanEdit"`
	IsShared             bool   `json:"isShared"`
	MainCardID           string `json:"mainCardId"`
	MainCardIDIsCardFace bool   `json:"mainCardIdIsCardFace"`
	MainCardIDIsBackFace bool   `json:"mainCardIdIsBackFace"`
	CreatedByUser        struct {
		UserName    string `json:"userName"`
		DisplayName string `json:"displayName"`
		Badges      []any  `json:"badges"`
	} `json:"createdByUser"`
	Authors []struct {
		UserName    string `json:"userName"`
		DisplayName string `json:"displayName"`
		Badges      []any  `json:"badges"`
	} `json:"authors"`
	CreatedAtUtc     time.Time `json:"createdAtUtc"`
	LastUpdatedAtUtc time.Time `json:"lastUpdatedAtUtc"`
	MainboardCount   int       `json:"mainboardCount"`
	SideboardCount   int       `json:"sideboardCount"`
	MaybeboardCount  int       `json:"maybeboardCount"`
	HubNames         []any     `json:"hubNames"`
	Folder           struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"folder"`
	Colors           []string `json:"colors"`
	ColorPercentages struct {
		White float64 `json:"white"`
		Blue  float64 `json:"blue"`
		Black float64 `json:"black"`
		Red   float64 `json:"red"`
		Green float64 `json:"green"`
	} `json:"colorPercentages"`
	ColorIdentity            []string `json:"colorIdentity"`
	ColorIdentityPercentages struct {
		White float64 `json:"white"`
		Blue  float64 `json:"blue"`
		Black float64 `json:"black"`
		Red   float64 `json:"red"`
		Green float64 `json:"green"`
	} `json:"colorIdentityPercentages"`
	IsPinned       bool `json:"isPinned"`
	DeckTier       int  `json:"deckTier"`
	CommanderTier  int  `json:"commanderTier"`
	DeckTier1Count int  `json:"deckTier1Count"`
	DeckTier2Count int  `json:"deckTier2Count"`
	DeckTier3Count int  `json:"deckTier3Count"`
	DeckTier4Count int  `json:"deckTier4Count"`
	Commanders     []struct {
		ID                    string `json:"id"`
		UniqueCardID          string `json:"uniqueCardId"`
		Name                  string `json:"name"`
		ImageCardID           string `json:"imageCardId"`
		ImageCardIDIsCardFace bool   `json:"imageCardIdIsCardFace"`
	} `json:"commanders,omitempty"`
	SignatureSpells []struct {
		ID                    string `json:"id"`
		UniqueCardID          string `json:"uniqueCardId"`
		Name                  string `json:"name"`
		ImageCardID           string `json:"imageCardId"`
		ImageCardIDIsCardFace bool   `json:"imageCardIdIsCardFace"`
	} `json:"signatureSpells,omitempty"`
	IsMature     bool `json:"isMature,omitempty"`
	IsMatureAuto bool `json:"isMatureAuto,omitempty"`
}

type card struct {
	Quantity int `json:"quantity"`
	Card     struct {
		ID             string   `json:"id"`
		UniqueCardID   string   `json:"uniqueCardId"`
		ScryfallID     string   `json:"scryfall_id"`
		Set            string   `json:"set"`
		SetName        string   `json:"set_name"`
		Name           string   `json:"name"`
		Cn             string   `json:"cn"`
		Layout         string   `json:"layout"`
		Cmc            float64  `json:"cmc"`
		Type           string   `json:"type"`
		TypeLine       string   `json:"type_line"`
		OracleText     string   `json:"oracle_text"`
		ManaCost       string   `json:"mana_cost"`
		Colors         []string `json:"colors"`
		ColorIndicator []string `json:"color_indicator"`
		ColorIdentity  []string `json:"color_identity"`
	} `json:"card"`
}

type deck struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Format      string `json:"format"`
	Visibility  string `json:"visibility"`
	PublicURL   string `json:"publicUrl"`
	PublicID    string `json:"publicId"`
	Boards      struct {
		Commanders struct {
			Count int             `json:"count"`
			Cards map[string]card `json:"cards"`
		} `json:"commanders"`
		Companions struct {
			Count int             `json:"count"`
			Cards map[string]card `json:"cards"`
		} `json:"companions"`
		SignatureSpells struct {
			Count int             `json:"count"`
			Cards map[string]card `json:"cards"`
		} `json:"signatureSpells"`
	} `json:"boards"`
	Version int `json:"version"`
	Hubs    []struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	} `json:"hubs"`
	CreatedAtUtc     time.Time `json:"createdAtUtc"`
	LastUpdatedAtUtc time.Time `json:"lastUpdatedAtUtc"`
	Colors           []string  `json:"colors"`
	ColorPercentages struct {
		White float64 `json:"white"`
		Blue  float64 `json:"blue"`
		Black float64 `json:"black"`
		Red   float64 `json:"red"`
		Green float64 `json:"green"`
	} `json:"colorPercentages"`
	ColorIdentity            []string `json:"colorIdentity"`
	ColorIdentityPercentages struct {
		White float64 `json:"white"`
		Blue  float64 `json:"blue"`
		Black float64 `json:"black"`
		Red   float64 `json:"red"`
		Green float64 `json:"green"`
	} `json:"colorIdentityPercentages"`
	OwnerUserID    string `json:"ownerUserId"`
	DeckTier       int    `json:"deckTier"`
	CommanderTier  int    `json:"commanderTier"`
	DeckTier1Count int    `json:"deckTier1Count"`
	DeckTier2Count int    `json:"deckTier2Count"`
	DeckTier3Count int    `json:"deckTier3Count"`
	DeckTier4Count int    `json:"deckTier4Count"`
}
