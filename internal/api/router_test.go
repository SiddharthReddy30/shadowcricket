package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// --- health endpoint ---

func TestHealth_Returns200(t *testing.T) {
	router := newTestRouter(t, "production")
	req := httptest.NewRequest("GET", "/api/health", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
}

func TestHealth_ReturnsStatusOk(t *testing.T) {
	router := newTestRouter(t, "production")
	req := httptest.NewRequest("GET", "/api/health", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	body := decodeJSON(t, rec.Body)
	if body["status"] != "ok" {
		t.Errorf("expected status ok, got %v", body["status"])
	}
}

func TestHealth_ContentTypeJSON(t *testing.T) {
	router := newTestRouter(t, "production")
	req := httptest.NewRequest("GET", "/api/health", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	ct := rec.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %q", ct)
	}
}

// --- 404 handling ---

func TestUnknownRoute_Returns404(t *testing.T) {
	router := newTestRouter(t, "production")
	req := httptest.NewRequest("GET", "/api/nonexistent", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", rec.Code)
	}
}

// --- wrong method handling ---

func TestWrongMethod_PostToHealth(t *testing.T) {
	router := newTestRouter(t, "production")
	req := httptest.NewRequest("POST", "/api/health", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status 405, got %d", rec.Code)
	}
}

func TestWrongMethod_GetToGuess(t *testing.T) {
	router := newTestRouter(t, "production")
	req := httptest.NewRequest("GET", "/api/game/guess", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status 405, got %d", rec.Code)
	}
}

// --- CORS in dev vs prod ---

func TestCORSInDev_SetsOriginHeader(t *testing.T) {
	router := newTestRouter(t, "development")
	req := httptest.NewRequest("GET", "/api/health", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	origin := rec.Header().Get("Access-Control-Allow-Origin")
	if origin != "http://localhost:5173" {
		t.Errorf("expected CORS origin http://localhost:5173 in dev, got %q", origin)
	}
}

func TestCORSInDev_PreflightReturns204(t *testing.T) {
	router := newTestRouter(t, "development")
	req := httptest.NewRequest("OPTIONS", "/api/game/random", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Errorf("expected status 204 for preflight in dev, got %d", rec.Code)
	}
}

func TestCORSInProd_NoOriginHeader(t *testing.T) {
	router := newTestRouter(t, "production")
	req := httptest.NewRequest("GET", "/api/health", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	origin := rec.Header().Get("Access-Control-Allow-Origin")
	if origin != "" {
		t.Errorf("expected no CORS origin in production, got %q", origin)
	}
}
