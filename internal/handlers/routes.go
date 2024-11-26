package handlers

import (
	"net/http"

	"github.com/jcbwlkr/deck-stats/internal/services/decks"
)

func App(svc *decks.Service) http.Handler {
	mux := http.NewServeMux()

	/*
		GET /api/decks
		POST /api/refresh
	*/

	return mux
}
