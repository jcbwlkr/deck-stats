package handlers

import (
	"net/http"

	"github.com/jcbwlkr/deck-stats/internal/auth"
	"github.com/jcbwlkr/deck-stats/internal/domains/magic"
	"github.com/jcbwlkr/deck-stats/internal/domains/users"
)

func App(
	magicService *magic.Service,
	userService *users.Service,
	authenticator *auth.Authenticator,
) http.Handler {
	mux := http.NewServeMux()

	authMW := authenticator.Middleware()

	deckHandlers := DeckHandlers{
		svc:     magicService,
		userSvc: userService,
	}
	mux.HandleFunc("GET /api/decks", authMW(deckHandlers.GetDecks))

	mux.HandleFunc("POST /api/accounts", authMW(deckHandlers.CreateAccount))
	mux.HandleFunc("POST /api/accounts/{id}/refresh", authMW(deckHandlers.RefreshAccount))

	userHandlers := UserHandlers{
		a:   authenticator,
		svc: userService,
	}
	mux.HandleFunc("POST /api/auth/register", userHandlers.Register)
	mux.HandleFunc("POST /api/auth/login", userHandlers.Login)

	return mux
}
