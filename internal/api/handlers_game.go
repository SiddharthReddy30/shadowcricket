package api

import (
	"encoding/json"
	"net/http"

	"github.com/siddharthreddy/shadowcricket/internal/game"
	"github.com/siddharthreddy/shadowcricket/internal/player"
)

type Handler struct {
	store  *player.Store
	secret string
}

func (h *Handler) randomGame(w http.ResponseWriter, r *http.Request) {
	video, target := h.store.RandomVideo()
	token, err := game.CreateToken(target.ID, h.secret)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create token")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"token":    token,
		"video_id": video.ID,
	})
}

type guessRequest struct {
	Token    string `json:"token"`
	PlayerID int    `json:"player_id"`
}

func (h *Handler) guess(w http.ResponseWriter, r *http.Request) {
	var req guessRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	targetID, err := game.DecryptToken(req.Token, h.secret)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid token")
		return
	}

	target, ok := h.store.Lookup(targetID)
	if !ok {
		writeError(w, http.StatusInternalServerError, "target player not found")
		return
	}

	guessedPlayer, ok := h.store.Lookup(req.PlayerID)
	if !ok {
		writeError(w, http.StatusBadRequest, "player not found")
		return
	}

	result := game.EvaluateGuess(target, guessedPlayer)
	writeJSON(w, http.StatusOK, map[string]any{
		"correct":  result.Correct,
		"feedback": result.Feedback,
		"token":    req.Token,
	})
}
