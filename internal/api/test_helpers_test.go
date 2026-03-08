package api

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/siddharthreddy/shadowcricket/internal/player"
)

const testSecret = "test-secret-key-32-bytes-long!x!"

const testPlayersJSON = `[
  {"id":1,"name":"Virat Kohli","country":"India","jersey_number":18,"role":"Middle-Order Batsman","ipl_team":"Royal Challengers Bengaluru","is_wicket_keeper":false},
  {"id":2,"name":"MS Dhoni","country":"India","jersey_number":7,"role":"Finisher","ipl_team":"Chennai Super Kings","is_wicket_keeper":true},
  {"id":3,"name":"AB de Villiers","country":"South Africa","jersey_number":17,"role":"Middle-Order Batsman","ipl_team":"Royal Challengers Bengaluru","is_wicket_keeper":true}
]`

const testVideosJSON = `[
  {"id":1,"player_id":1,"raw_video":"virat_1.mp4","silhouette_video":"virat_1_sil.mp4"},
  {"id":2,"player_id":2,"raw_video":"dhoni_1.mp4","silhouette_video":"dhoni_1_sil.mp4"}
]`

func newTestRouter(t *testing.T, env string) http.Handler {
	t.Helper()

	dir := t.TempDir()
	playersFile := filepath.Join(dir, "players.json")
	videosFile := filepath.Join(dir, "videos.json")

	if err := os.WriteFile(playersFile, []byte(testPlayersJSON), 0644); err != nil {
		t.Fatalf("failed to write test players file: %v", err)
	}
	if err := os.WriteFile(videosFile, []byte(testVideosJSON), 0644); err != nil {
		t.Fatalf("failed to write test videos file: %v", err)
	}

	store, err := player.LoadStore(playersFile, videosFile)
	if err != nil {
		t.Fatalf("failed to load test store: %v", err)
	}

	return NewRouter(store, testSecret, env)
}

func decodeJSON(t *testing.T, body io.Reader) map[string]any {
	t.Helper()
	var result map[string]any
	if err := json.NewDecoder(body).Decode(&result); err != nil {
		t.Fatalf("failed to decode JSON response: %v", err)
	}
	return result
}

func decodeJSONArray(t *testing.T, body io.Reader) []map[string]any {
	t.Helper()
	var result []map[string]any
	if err := json.NewDecoder(body).Decode(&result); err != nil {
		t.Fatalf("failed to decode JSON array response: %v", err)
	}
	return result
}
