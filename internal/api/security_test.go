package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSecurity_TamperedTokenBase64(t *testing.T) {
	router := newTestRouter(t, "production")
	token, _ := getRandomToken(t, router)

	// flip a character near the middle
	mid := len(token) / 2
	flipped := token[:mid] + string(token[mid]^0x01) + token[mid+1:]

	body := fmt.Sprintf(`{"token":"%s","player_id":1}`, flipped)
	req := httptest.NewRequest("POST", "/api/game/guess", strings.NewReader(body))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for tampered base64 token, got %d", rec.Code)
	}
}

func TestSecurity_LargeRequestBody(t *testing.T) {
	router := newTestRouter(t, "production")

	// ~1MB of garbage JSON
	largeBody := `{"token":"` + strings.Repeat("A", 1024*1024) + `","player_id":1}`
	req := httptest.NewRequest("POST", "/api/game/guess", strings.NewReader(largeBody))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	// should respond (not hang) — either 400 (invalid token) or similar
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for large body, got %d", rec.Code)
	}
}

func TestSecurity_SQLInjectionInSearch(t *testing.T) {
	router := newTestRouter(t, "production")
	req := httptest.NewRequest("GET", "/api/players/search?q='+OR+1%3D1+--", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
	results := decodeJSONArray(t, rec.Body)
	if len(results) != 0 {
		t.Errorf("SQL injection should match nothing, got %d results", len(results))
	}
}

func TestSecurity_XSSInSearch(t *testing.T) {
	router := newTestRouter(t, "production")
	req := httptest.NewRequest("GET", "/api/players/search?q=%3Cscript%3Ealert(1)%3C/script%3E", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}

	ct := rec.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("response must be JSON (not HTML), got %q", ct)
	}

	results := decodeJSONArray(t, rec.Body)
	if len(results) != 0 {
		t.Errorf("XSS payload should match nothing, got %d results", len(results))
	}
}

func TestSecurity_UnicodeInSearch(t *testing.T) {
	router := newTestRouter(t, "production")
	req := httptest.NewRequest("GET", "/api/players/search?q=%F0%9F%8F%8F", nil) // cricket emoji
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200 for unicode search, got %d", rec.Code)
	}
}

func TestSecurity_StringPlayerID(t *testing.T) {
	router := newTestRouter(t, "production")
	token, _ := getRandomToken(t, router)

	body := fmt.Sprintf(`{"token":"%s","player_id":"abc"}`, token)
	req := httptest.NewRequest("POST", "/api/game/guess", strings.NewReader(body))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for string player_id, got %d", rec.Code)
	}
}

func TestSecurity_NullFields(t *testing.T) {
	router := newTestRouter(t, "production")
	req := httptest.NewRequest("POST", "/api/game/guess", strings.NewReader(`{"token":null,"player_id":null}`))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	// null token decodes to empty string "" → invalid token
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for null fields, got %d", rec.Code)
	}
}

func TestSecurity_ExtraFieldsIgnored(t *testing.T) {
	router := newTestRouter(t, "production")
	token, targetID := getRandomToken(t, router)

	body := fmt.Sprintf(`{"token":"%s","player_id":%d,"extra":"field","admin":true}`, token, targetID)
	req := httptest.NewRequest("POST", "/api/game/guess", strings.NewReader(body))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200 (extra fields ignored), got %d", rec.Code)
	}
	result := decodeJSON(t, rec.Body)
	if result["correct"] != true {
		t.Error("expected correct guess despite extra fields")
	}
}

func TestSecurity_MissingContentTypeHeader(t *testing.T) {
	router := newTestRouter(t, "production")
	token, targetID := getRandomToken(t, router)

	body := fmt.Sprintf(`{"token":"%s","player_id":%d}`, token, targetID)
	req := httptest.NewRequest("POST", "/api/game/guess", strings.NewReader(body))
	// explicitly do NOT set Content-Type
	req.Header.Del("Content-Type")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	// json.Decoder doesn't check Content-Type, so this should succeed
	if rec.Code != http.StatusOK {
		t.Errorf("expected 200 without Content-Type header, got %d", rec.Code)
	}
}

func TestSecurity_IntegerOverflowPlayerID(t *testing.T) {
	router := newTestRouter(t, "production")
	token, _ := getRandomToken(t, router)

	body := fmt.Sprintf(`{"token":"%s","player_id":99999999999999999999}`, token)
	req := httptest.NewRequest("POST", "/api/game/guess", strings.NewReader(body))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	// should return 400 (json decode error for overflow) — should not panic
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for integer overflow, got %d", rec.Code)
	}
}
