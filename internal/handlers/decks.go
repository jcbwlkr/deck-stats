package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/jcbwlkr/deck-stats/internal/auth"
	"github.com/jcbwlkr/deck-stats/internal/domains/magic"
	"github.com/jcbwlkr/deck-stats/internal/domains/users"
)

type DeckHandlers struct {
	svc     *magic.Service
	userSvc *users.Service
}

func (h *DeckHandlers) GetDecks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, ok := auth.User(ctx)
	if !ok {
		http.Error(w, "system error", http.StatusInternalServerError)
		return
	}

	decks, err := h.svc.GetDecksForUser(ctx, user)
	if err != nil {
		slog.ErrorContext(ctx, "could not list decks", "error", err)
		http.Error(w, "system error", http.StatusInternalServerError)
		return
	}

	response := struct {
		Decks []magic.Deck `json:"decks"`
	}{
		Decks: decks,
	}
	json.NewEncoder(w).Encode(response)
}

func (h *DeckHandlers) CreateAccount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, ok := auth.User(ctx)
	if !ok {
		http.Error(w, "system error", http.StatusInternalServerError)
		return
	}

	var input users.NewAccount
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		slog.WarnContext(ctx, "could not decode input", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	account, err := h.userSvc.CreateAccount(ctx, user.ID, input)
	if err != nil {
		slog.ErrorContext(ctx, "could not create account for user", "error", err)
		http.Error(w, "system error", http.StatusInternalServerError)
		return
	}

	response := struct {
		Account users.Account `json:"account"`
	}{
		Account: account,
	}
	json.NewEncoder(w).Encode(response)
}

func (h *DeckHandlers) RefreshAccount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, ok := auth.User(ctx)
	if !ok {
		http.Error(w, "system error", http.StatusInternalServerError)
		return
	}

	accountID := r.PathValue("id")

	account, err := h.userSvc.GetAccount(ctx, accountID, user.ID)
	if err != nil {
		slog.ErrorContext(ctx, "could not find account for user", "error", err)
		http.Error(w, "system error", http.StatusInternalServerError)
		return
	}

	h.svc.RefreshDecks(user, account)
}
