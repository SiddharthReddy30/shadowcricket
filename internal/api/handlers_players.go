package api

import (
	"net/http"
)

type playerResult struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (h *Handler) search(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		writeJSON(w, http.StatusOK, []playerResult{})
		return
	}

	players := h.store.Search(query)
	results := make([]playerResult, len(players))
	for i, p := range players {
		results[i] = playerResult{ID: p.ID, Name: p.Name}
	}
	writeJSON(w, http.StatusOK, results)
}
