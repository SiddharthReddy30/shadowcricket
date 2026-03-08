package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/siddharthreddy/shadowcricket/internal/game"
)

// --- helper to get a valid token from /api/game/random ---

func getRandomToken(t *testing.T, router http.Handler) (string, int) {
	t.Helper()
	req := httptest.NewRequest("GET", "/api/game/random", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("GET /api/game/random returned %d", rec.Code)
	}

	body := decodeJSON(t, rec.Body)
	token, ok := body["token"].(string)
	if !ok || token == "" {
		t.Fatal("expected non-empty token from /api/game/random")
	}

	targetID, err := game.DecryptToken(token, testSecret)
	if err != nil {
		t.Fatalf("failed to decrypt token: %v", err)
	}

	return token, targetID
}

// --- randomGame tests ---

func TestRandomGame_Status200(t *testing.T) {
	router := newTestRouter(t, "production")
	req := httptest.NewRequest("GET", "/api/game/random", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
}

func TestRandomGame_ReturnsTokenAndVideoID(t *testing.T) {
	router := newTestRouter(t, "production")
	req := httptest.NewRequest("GET", "/api/game/random", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	body := decodeJSON(t, rec.Body)

	token, ok := body["token"].(string)
	if !ok || token == "" {
		t.Error("expected non-empty token string")
	}

	videoID, ok := body["video_id"].(float64)
	if !ok || videoID <= 0 {
		t.Errorf("expected positive video_id, got %v", body["video_id"])
	}
}

func TestRandomGame_TokenIsDecryptable(t *testing.T) {
	router := newTestRouter(t, "production")
	req := httptest.NewRequest("GET", "/api/game/random", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	body := decodeJSON(t, rec.Body)
	token := body["token"].(string)

	playerID, err := game.DecryptToken(token, testSecret)
	if err != nil {
		t.Fatalf("token should be decryptable: %v", err)
	}
	// player ID should be one of our test players (1, 2, or 3)
	if playerID < 1 || playerID > 3 {
		t.Errorf("expected player ID between 1-3, got %d", playerID)
	}
}

func TestRandomGame_ContentTypeJSON(t *testing.T) {
	router := newTestRouter(t, "production")
	req := httptest.NewRequest("GET", "/api/game/random", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	ct := rec.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %q", ct)
	}
}

// --- guess tests: happy paths ---

func TestGuess_CorrectGuess(t *testing.T) {
	router := newTestRouter(t, "production")
	token, targetID := getRandomToken(t, router)

	guessBody := fmt.Sprintf(`{"token":"%s","player_id":%d}`, token, targetID)
	req := httptest.NewRequest("POST", "/api/game/guess", strings.NewReader(guessBody))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	body := decodeJSON(t, rec.Body)
	if body["correct"] != true {
		t.Error("expected correct to be true for matching player_id")
	}

	feedback, ok := body["feedback"].([]any)
	if !ok {
		t.Fatal("expected feedback array")
	}
	if len(feedback) != 5 {
		t.Errorf("expected 5 feedback fields, got %d", len(feedback))
	}
	for _, f := range feedback {
		field := f.(map[string]any)
		if field["color"] != "green" {
			t.Errorf("expected all green for correct guess, field %s got %s", field["field"], field["color"])
		}
	}
}

func TestGuess_WrongGuess(t *testing.T) {
	router := newTestRouter(t, "production")
	token, targetID := getRandomToken(t, router)

	// pick a different player ID
	wrongID := 1
	if targetID == 1 {
		wrongID = 2
	}

	guessBody := fmt.Sprintf(`{"token":"%s","player_id":%d}`, token, wrongID)
	req := httptest.NewRequest("POST", "/api/game/guess", strings.NewReader(guessBody))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	body := decodeJSON(t, rec.Body)
	if body["correct"] != false {
		t.Error("expected correct to be false for wrong player_id")
	}

	feedback, ok := body["feedback"].([]any)
	if !ok {
		t.Fatal("expected feedback array")
	}
	if len(feedback) != 5 {
		t.Errorf("expected 5 feedback fields, got %d", len(feedback))
	}
}

func TestGuess_EchoesToken(t *testing.T) {
	router := newTestRouter(t, "production")
	token, targetID := getRandomToken(t, router)

	guessBody := fmt.Sprintf(`{"token":"%s","player_id":%d}`, token, targetID)
	req := httptest.NewRequest("POST", "/api/game/guess", strings.NewReader(guessBody))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	body := decodeJSON(t, rec.Body)
	if body["token"] != token {
		t.Errorf("expected token to be echoed back, got %v", body["token"])
	}
}

// --- guess tests: error cases ---

func TestGuess_InvalidJSON(t *testing.T) {
	router := newTestRouter(t, "production")
	req := httptest.NewRequest("POST", "/api/game/guess", strings.NewReader("{invalid"))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}
	body := decodeJSON(t, rec.Body)
	if body["error"] != "invalid request body" {
		t.Errorf("expected error 'invalid request body', got %v", body["error"])
	}
}

func TestGuess_EmptyBody(t *testing.T) {
	router := newTestRouter(t, "production")
	req := httptest.NewRequest("POST", "/api/game/guess", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}
}

func TestGuess_EmptyToken(t *testing.T) {
	router := newTestRouter(t, "production")
	req := httptest.NewRequest("POST", "/api/game/guess", strings.NewReader(`{"token":"","player_id":1}`))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}
	body := decodeJSON(t, rec.Body)
	if body["error"] != "invalid token" {
		t.Errorf("expected error 'invalid token', got %v", body["error"])
	}
}

func TestGuess_TamperedToken(t *testing.T) {
	router := newTestRouter(t, "production")
	token, _ := getRandomToken(t, router)

	// flip a character in the middle of the token
	tampered := token[:len(token)/2] + "X" + token[len(token)/2+1:]

	guessBody := fmt.Sprintf(`{"token":"%s","player_id":1}`, tampered)
	req := httptest.NewRequest("POST", "/api/game/guess", strings.NewReader(guessBody))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 for tampered token, got %d", rec.Code)
	}
	body := decodeJSON(t, rec.Body)
	if body["error"] != "invalid token" {
		t.Errorf("expected error 'invalid token', got %v", body["error"])
	}
}

func TestGuess_CompletelyFakeToken(t *testing.T) {
	router := newTestRouter(t, "production")
	req := httptest.NewRequest("POST", "/api/game/guess", strings.NewReader(`{"token":"dGhpcyBpcyBub3QgYSB0b2tlbg==","player_id":1}`))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}
}

func TestGuess_PlayerNotFound(t *testing.T) {
	router := newTestRouter(t, "production")
	token, _ := getRandomToken(t, router)

	guessBody := fmt.Sprintf(`{"token":"%s","player_id":9999}`, token)
	req := httptest.NewRequest("POST", "/api/game/guess", strings.NewReader(guessBody))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}
	body := decodeJSON(t, rec.Body)
	if body["error"] != "player not found" {
		t.Errorf("expected error 'player not found', got %v", body["error"])
	}
}

func TestGuess_NegativePlayerID(t *testing.T) {
	router := newTestRouter(t, "production")
	token, _ := getRandomToken(t, router)

	guessBody := fmt.Sprintf(`{"token":"%s","player_id":-1}`, token)
	req := httptest.NewRequest("POST", "/api/game/guess", strings.NewReader(guessBody))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}
}
