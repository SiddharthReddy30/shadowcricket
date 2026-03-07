package api

import (
	"net/http"

	"github.com/siddharthreddy/shadowcricket/internal/player"
)

func NewRouter(store *player.Store, secret string) http.Handler {
	h := &Handler{store: store, secret: secret}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	mux.HandleFunc("GET /api/game/random", h.randomGame)
	mux.HandleFunc("POST /api/game/guess", h.guess)
	mux.HandleFunc("GET /api/players/search", h.search)

	return mux
}
