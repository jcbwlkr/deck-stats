package handlers

import (
	"net/http"

	"github.com/jcbwlkr/deck-stats/internal/domains/magic"
)

func App(svc *magic.Service) http.Handler {
	mux := http.NewServeMux()

	/*
		GET /api/decks
		POST /api/refresh
	*/

	return mux
}
