package decks

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"slices"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/jcbwlkr/deck-stats/internal/magic"
	"github.com/jcbwlkr/deck-stats/internal/moxfield"
	"github.com/jcbwlkr/deck-stats/internal/services/users"
)

type Service struct {
	db *sqlx.DB
	us *users.Service
	mc *moxfield.Client
	wg sync.WaitGroup
}

func NewService(db *sqlx.DB, us *users.Service, mc *moxfield.Client) *Service {
	return &Service{db: db, us: us, mc: mc, wg: sync.WaitGroup{}}
}

func (s *Service) Wait() {
	s.wg.Wait()
}

func (s *Service) RefreshDecks(ctx context.Context, user users.User) error {

	for _, account := range user.Accounts {
		// Make a new context for these goroutines.
		ctx := context.Background()

		log := slog.With("user", user.ID, "service", account.Service)

		// Channel for these two goroutines to talk to each other.
		result := make(chan error)

		// Goroutine that updates decks for this service.
		s.wg.Add(1)
		go func(account users.Account) {
			defer s.wg.Done()
			defer close(result) // In case of panic, let other goroutine terminate

			switch account.Service {
			case moxfield.ServiceName:
				result <- s.refreshMoxfield(ctx, log, user, account)
			default:
				result <- errors.New("unknown service for refresh")
			}
		}(account)

		// Goroutine that periodically marks the account as either still refreshing
		// or completed/failed.
		s.wg.Add(1)
		go func(account users.Account) {
			defer s.wg.Done()
			now := time.Now()
			account.RefreshStartedAt = &now
			account.RefreshActiveAt = &now
			account.RefreshCompletedAt = nil
			account.RefreshStatus = "pending"
			if err := s.us.UpdateAccount(ctx, account); err != nil {
				log.ErrorContext(ctx, "account status goroutine died", "error", err)
				return
			}

			tick := time.NewTicker(10 * time.Second)
			defer tick.Stop()

			for {
				select {
				case now := <-tick.C:
					log.InfoContext(ctx, "service still refreshing")
					account.RefreshActiveAt = &now
					if err := s.us.UpdateAccount(ctx, account); err != nil {
						log.ErrorContext(ctx, "account status goroutine died", "error", err)
						return
					}

				case err := <-result:
					if err != nil {
						log.ErrorContext(ctx, "failed to refresh", "error", err)
						account.RefreshStatus = "failed: " + err.Error()
					} else {
						log.InfoContext(ctx, "refresh complete")
						account.RefreshStatus = "completed"
					}
					now := time.Now()
					account.RefreshCompletedAt = &now
					if err := s.us.UpdateAccount(ctx, account); err != nil {
						log.ErrorContext(ctx, "account status goroutine died", "error", err)
					}
					break
				}
			}
		}(account)
	}
	return nil
}

func (s *Service) refreshMoxfield(ctx context.Context, log *slog.Logger, user users.User, account users.Account) error {
	start := time.Now()

	existingDecks, err := s.GetDecksForUserAndService(ctx, user, account.Service)
	if err != nil {
		return fmt.Errorf("could not list existing decks: %w", err)
	}

	moxfieldDecks, err := s.mc.ListMyDecks(ctx, account.Token)
	if err != nil {
		return fmt.Errorf("could not list moxfield decks: %w", err)
	}

	for _, moxDeck := range moxfieldDecks {

		// Look for this deck in our db results
		i := slices.IndexFunc(existingDecks, func(d magic.Deck) bool {
			return d.ServiceID == moxDeck.ServiceID
		})

		var deck *magic.Deck
		if i >= 0 {
			deck = &existingDecks[i]
			log = log.With("id", deck.ID)
		}

		switch {
		case deck == nil:
			log.Info("new deck found")
			moxDeck.UserID = user.ID
			moxDeck.RefreshedAt = time.Now()
			if err := s.mc.AddDeckDetails(ctx, account.Token, &moxDeck); err != nil {
				return err
			}
			if err := s.InsertDeck(ctx, moxDeck); err != nil {
				return err
			}

		case deck.RefreshedAt.Before(moxDeck.UpdatedAt):
			log.Info("stale deck found")
			moxDeck.ID = deck.ID
			moxDeck.UserID = deck.UserID
			moxDeck.RefreshedAt = time.Now()
			if err := s.mc.AddDeckDetails(ctx, account.Token, &moxDeck); err != nil {
				return err
			}
			if err := s.UpdateDeck(ctx, moxDeck); err != nil {
				return err
			}

		default:
			log.Info("deck is up to date")
			deck.RefreshedAt = time.Now()
			if err := s.UpdateDeck(ctx, *deck); err != nil {
				return err
			}
		}
	}

	for _, eDeck := range existingDecks {
		if eDeck.UpdatedAt.Before(start) {
			log.Info("deleting deck that wasn't on moxfield", "id", eDeck.ID)
		}
	}

	return nil
}

func (s *Service) GetDecksForUserAndService(ctx context.Context, user users.User, service string) ([]magic.Deck, error) {

	const q = `
	SELECT
		id,
		user_id,
		service,
		service_id,
		name,
		format,
		url,
		color_identity,
		folder,
		archetypes,
		updated_at
	FROM decks
	WHERE user_id = $1
		AND service = $2`

	decks := []magic.Deck{}
	err := s.db.SelectContext(ctx, &decks, q, user.ID, service)
	return decks, err
}

func (s *Service) InsertDeck(ctx context.Context, deck magic.Deck) error {

	const q = `
	INSERT INTO decks
		id,
		user_id,
		service,
		service_id,
		name,
		format,
		url,
		color_identity,
		folder,
		leaders,
		archetypes,
		updated_at
	VALUES (
		:id,
		:user_id,
		:service,
		:service_id,
		:name,
		:format,
		:url,
		:color_identity,
		:folder,
		:leaders,
		:archetypes,
		:updated_at
	)`

	deck.ID = uuid.New().String()

	_, err := s.db.NamedExecContext(ctx, q, deck)
	return err
}

func (s *Service) UpdateDeck(ctx context.Context, deck magic.Deck) error {
	const q = `
	UPDATE decks SET
		name = :name,
		format = :format,
		url = :url,
		color_identity = :color_identity,
		folder = :folder,
		leaders = :leaders,
		archetypes = :archetypes,
		updated_a = :updated_at
	WHERE id = :id
	`

	_, err := s.db.NamedExecContext(ctx, q, deck)
	return err
}
