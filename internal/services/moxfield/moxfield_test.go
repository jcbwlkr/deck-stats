package moxfield

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"github.com/jcbwlkr/deck-stats/internal/domains/magic"
	"github.com/jcbwlkr/deck-stats/internal/services"
)

func TestClientListDecks(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		filename := ""
		switch r.URL.Path {
		case "/v2/users/jcbwlkr/decks":
			if r.URL.Query().Get("pageNumber") == "1" {
				filename = "testdata/decks-page-01.json"
			} else {
				filename = "testdata/decks-page-02.json"
			}
		case "/v3/decks/all/CZ0EOEAF0Ue1wZQRUNARZw":
			filename = "testdata/deck-vial-smasher.json"
		case "/v3/decks/all/H7Es1cn9y0qJVkT8_2A-Zw":
			filename = "testdata/deck-winota.json"
		case "/v3/decks/all/ryBvdOgMCEqkKKpbNVd7tQ":
			filename = "testdata/deck-cats.json"
		case "/v3/decks/all/NM5x9FZ6cEGEEfx8RSU5mQ":
			filename = "testdata/deck-poison.json"
		case "/v3/decks/all/hJ0iLLhZ1ECZlQpyX8czKA":
			filename = "testdata/deck-blanka.json"
		}

		if filename == "" {
			w.Write([]byte("{}"))
			return
		}

		f, err := os.Open(filename)
		if err != nil {
			t.Fatal(err)
		}
		io.Copy(w, f)
	}))
	defer srv.Close()

	client := NewClient("", 0)
	client.url = srv.URL

	ctx := context.Background()

	deckList, err := client.ListDecks(ctx, "jcbwlkr")
	if err != nil {
		t.Fatal(err)
	}

	if got, want := len(deckList), 30; got != want {
		t.Fatalf("should have %d decks but got %d", want, got)
	}

	for i := range deckList {
		err := client.AddDeckDetails(ctx, &deckList[i])
		if err != nil {
			t.Fatal(err)
		}
	}

	vialSmasher := magic.Deck{
		Name:          "1,000 Smashed Vials",
		Format:        "commander",
		Service:       services.Moxfield,
		ServiceID:     "CZ0EOEAF0Ue1wZQRUNARZw",
		URL:           "https://www.moxfield.com/decks/CZ0EOEAF0Ue1wZQRUNARZw",
		ColorIdentity: magic.Grixis,
		Leaders: magic.Leaders{
			Commanders: []magic.Card{
				{
					ID:   "4e439cd0-5ba1-45af-a868-f408e0a50465",
					Name: "Vial Smasher the Fierce",
				},
				{
					ID:   "714c3a1f-7b30-4ed8-8f38-6176758741fb",
					Name: "Sakashima of a Thousand Faces",
				},
			},
		},
		Archetypes: []magic.Archetype{
			{Name: "Burn", Description: "A deck that focuses on dealing direct damage to an opponent."},
			{Name: "Clones", Description: "A deck focusing on clone effects, or the ability to copy another creature."},
		},
		UpdatedAt: parseUTC("2024-11-30T21:41:45.337Z"),
	}
	if diff := cmp.Diff(vialSmasher, deckList[1]); diff != "" {
		t.Errorf("second response should be vial smasher but was:\n%s", diff)
	}

	cats := magic.Deck{
		Name:          "Cats!",
		Format:        "commander",
		Service:       services.Moxfield,
		ServiceID:     "ryBvdOgMCEqkKKpbNVd7tQ",
		URL:           "https://www.moxfield.com/decks/ryBvdOgMCEqkKKpbNVd7tQ",
		ColorIdentity: magic.Selesnya,
		Leaders: magic.Leaders{
			Commanders: []magic.Card{
				{
					ID:   "81dc3d00-97cd-4549-b5a4-15a1e08767f5",
					Name: "Arahbo, Roar of the World",
				},
			},
			Companion: &magic.Card{
				ID:   "d4ebed0b-8060-4a7b-a060-5cfcd2172b16",
				Name: "Kaheera, the Orphanguard",
			},
		},
		Archetypes: []magic.Archetype{},
		UpdatedAt:  parseUTC("2023-03-31T17:33:58.21Z"),
	}
	if diff := cmp.Diff(cats, deckList[22]); diff != "" {
		t.Errorf("second to last response should be Arahbo but was:\n%s", diff)
	}

	poison := magic.Deck{
		Name:          "Drink your Poison",
		Format:        "oathbreaker",
		Service:       services.Moxfield,
		ServiceID:     "NM5x9FZ6cEGEEfx8RSU5mQ",
		URL:           "https://www.moxfield.com/decks/NM5x9FZ6cEGEEfx8RSU5mQ",
		ColorIdentity: magic.Simic,
		Leaders: magic.Leaders{
			Oathbreakers: []magic.Card{
				{
					ID:   "222a736e-d819-452d-aeda-eb848c4b2302",
					Name: "Tamiyo, Compleated Sage",
				},
			},
			SignatureSpells: []magic.Card{
				{
					ID:   "ac625f30-ed91-4b21-ada8-aaa5b2ad79b8",
					Name: "Prologue to Phyresis",
				},
			},
		},
		Archetypes: []magic.Archetype{},
		UpdatedAt:  parseUTC("2024-11-30T21:42:01.61Z"),
	}
	if diff := cmp.Diff(poison, deckList[0]); diff != "" {
		t.Errorf("first response should be poison but was:\n%s", diff)
	}

	blanka := magic.Deck{
		Name:          "Blanka never loses!",
		Format:        "commander",
		Service:       services.Moxfield,
		ServiceID:     "hJ0iLLhZ1ECZlQpyX8czKA",
		URL:           "https://www.moxfield.com/decks/hJ0iLLhZ1ECZlQpyX8czKA",
		ColorIdentity: magic.Gruul,
		Leaders: magic.Leaders{
			Commanders: []magic.Card{
				{
					ID:   "26fecfd4-40c1-4ed3-bf58-22c7b98699a9",
					Name: "Blanka, Ferocious Friend",
				},
			},
		},
		Archetypes: []magic.Archetype{},
		UpdatedAt:  parseUTC("2024-11-30T21:41:34.123Z"),
	}
	if diff := cmp.Diff(blanka, deckList[2]); diff != "" {
		t.Errorf("third response should be blanka but was:\n%s", diff)
	}

	winota := magic.Deck{
		Name:          "Winota Ryder",
		Format:        "commander",
		Service:       services.Moxfield,
		ServiceID:     "H7Es1cn9y0qJVkT8_2A-Zw",
		URL:           "https://www.moxfield.com/decks/H7Es1cn9y0qJVkT8_2A-Zw",
		ColorIdentity: magic.Boros,
		Leaders: magic.Leaders{
			Commanders: []magic.Card{
				{
					ID:   "5dd13a6c-23d3-44ce-a628-cb1c19d777c4",
					Name: "Winota, Joiner of Forces",
				},
			},
		},
		Archetypes: []magic.Archetype{},
		UpdatedAt:  parseUTC("2024-10-20T13:45:59.85Z"),
	}
	if diff := cmp.Diff(winota, deckList[9]); diff != "" {
		t.Errorf("tenth response should be winota but was:\n%s", diff)
	}
}

func parseUTC(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}

	return t
}
