package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/jcbwlkr/deck-stats/internal/auth"
	"github.com/jcbwlkr/deck-stats/internal/domains/users"
)

type UserHandlers struct {
	a   *auth.Authenticator
	svc *users.Service
}

func (h *UserHandlers) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var input users.NewUser

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		slog.WarnContext(ctx, "could not decode input", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.svc.RegisterUser(ctx, input)
	if err != nil {
		switch err {
		case users.ErrBlankPassword,
			users.ErrPasswordsDoNotMatch,
			users.ErrUsernameRegistered:
			slog.WarnContext(ctx, "could not register user", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			slog.ErrorContext(ctx, "problem registering user", "error", err)
			http.Error(w, "system error", http.StatusInternalServerError)
		}
		return
	}

	json.NewEncoder(w).Encode(user)
}

func (h *UserHandlers) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var input users.LoginUser

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		slog.WarnContext(ctx, "could not decode input", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.svc.Login(ctx, input)
	if err != nil {
		switch err {
		case users.ErrUserNotFound,
			users.ErrPasswordInvalid:
			slog.WarnContext(ctx, "could not log user in", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			slog.ErrorContext(ctx, "problem logging user in", "error", err)
			http.Error(w, "system error", http.StatusInternalServerError)
		}
		return
	}

	token, err := h.a.GenerateJWT(user)
	if err != nil {
		slog.ErrorContext(ctx, "problem generating token", "error", err)
		http.Error(w, "system error", http.StatusInternalServerError)
		return
	}

	var response struct {
		Token string      `json:"token"`
		User  *users.User `json:"user"`
	}
	response.User = user
	response.Token = token

	json.NewEncoder(w).Encode(response)
}
