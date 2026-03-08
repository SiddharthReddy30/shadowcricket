package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSearch_HappyPath(t *testing.T) {
	router := newTestRouter(t, "production")
	req := httptest.NewRequest("GET", "/api/players/search?q=Virat", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	results := decodeJSONArray(t, rec.Body)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0]["name"] != "Virat Kohli" {
		t.Errorf("expected name Virat Kohli, got %v", results[0]["name"])
	}
}

func TestSearch_EmptyQuery(t *testing.T) {
	router := newTestRouter(t, "production")
	req := httptest.NewRequest("GET", "/api/players/search?q=", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	results := decodeJSONArray(t, rec.Body)
	if len(results) != 0 {
		t.Errorf("expected 0 results for empty query, got %d", len(results))
	}
}

func TestSearch_NoQueryParam(t *testing.T) {
	router := newTestRouter(t, "production")
	req := httptest.NewRequest("GET", "/api/players/search", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	results := decodeJSONArray(t, rec.Body)
	if len(results) != 0 {
		t.Errorf("expected 0 results when q is missing, got %d", len(results))
	}
}

func TestSearch_NoResults(t *testing.T) {
	router := newTestRouter(t, "production")
	req := httptest.NewRequest("GET", "/api/players/search?q=xyz", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	results := decodeJSONArray(t, rec.Body)
	if len(results) != 0 {
		t.Errorf("expected 0 results for unknown query, got %d", len(results))
	}
}

func TestSearch_PartialMatch(t *testing.T) {
	router := newTestRouter(t, "production")
	req := httptest.NewRequest("GET", "/api/players/search?q=vi", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	results := decodeJSONArray(t, rec.Body)
	if len(results) < 1 {
		t.Fatal("expected at least 1 result for partial match 'vi'")
	}

	found := false
	for _, r := range results {
		if r["name"] == "Virat Kohli" {
			found = true
		}
	}
	if !found {
		t.Error("expected Virat Kohli in partial match results")
	}
}

func TestSearch_CaseInsensitive(t *testing.T) {
	router := newTestRouter(t, "production")
	req := httptest.NewRequest("GET", "/api/players/search?q=virat", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	results := decodeJSONArray(t, rec.Body)
	if len(results) != 1 {
		t.Fatalf("expected 1 result for case-insensitive search, got %d", len(results))
	}
	if results[0]["name"] != "Virat Kohli" {
		t.Errorf("expected Virat Kohli, got %v", results[0]["name"])
	}
}

func TestSearch_ResultStructure(t *testing.T) {
	router := newTestRouter(t, "production")
	req := httptest.NewRequest("GET", "/api/players/search?q=Dhoni", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	results := decodeJSONArray(t, rec.Body)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}

	r := results[0]
	if _, ok := r["id"]; !ok {
		t.Error("expected result to have 'id' field")
	}
	if _, ok := r["name"]; !ok {
		t.Error("expected result to have 'name' field")
	}
	// should NOT leak other player fields
	if _, ok := r["country"]; ok {
		t.Error("result should not contain 'country' field")
	}
	if _, ok := r["jersey_number"]; ok {
		t.Error("result should not contain 'jersey_number' field")
	}
	if _, ok := r["role"]; ok {
		t.Error("result should not contain 'role' field")
	}
}

func TestSearch_ContentTypeJSON(t *testing.T) {
	router := newTestRouter(t, "production")
	req := httptest.NewRequest("GET", "/api/players/search?q=Virat", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	ct := rec.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %q", ct)
	}
}
