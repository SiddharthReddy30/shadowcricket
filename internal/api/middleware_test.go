package api

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// --- statusRecorder tests ---

func TestStatusRecorder_DefaultStatus(t *testing.T) {
	rec := &statusRecorder{
		ResponseWriter: httptest.NewRecorder(),
		status:         http.StatusOK,
	}
	if rec.status != http.StatusOK {
		t.Errorf("expected default status 200, got %d", rec.status)
	}
}

func TestStatusRecorder_CapturesWriteHeader(t *testing.T) {
	inner := httptest.NewRecorder()
	rec := &statusRecorder{ResponseWriter: inner, status: http.StatusOK}

	rec.WriteHeader(http.StatusNotFound)

	if rec.status != http.StatusNotFound {
		t.Errorf("expected captured status 404, got %d", rec.status)
	}
	if inner.Code != http.StatusNotFound {
		t.Errorf("expected underlying writer status 404, got %d", inner.Code)
	}
}

// --- CORS tests ---

func TestCORS_SetsHeaders(t *testing.T) {
	handler := CORS("http://localhost:5173")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/api/health", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	origin := rec.Header().Get("Access-Control-Allow-Origin")
	if origin != "http://localhost:5173" {
		t.Errorf("expected origin http://localhost:5173, got %q", origin)
	}

	methods := rec.Header().Get("Access-Control-Allow-Methods")
	if methods != "GET, POST" {
		t.Errorf("expected methods GET, POST, got %q", methods)
	}

	headers := rec.Header().Get("Access-Control-Allow-Headers")
	if headers != "Content-Type" {
		t.Errorf("expected headers Content-Type, got %q", headers)
	}
}

func TestCORS_PreflightReturns204(t *testing.T) {
	nextCalled := false
	handler := CORS("http://localhost:5173")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
	}))

	req := httptest.NewRequest("OPTIONS", "/api/game/random", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", rec.Code)
	}
	if nextCalled {
		t.Error("expected next handler to NOT be called on preflight")
	}
}

func TestCORS_NonPreflightCallsNext(t *testing.T) {
	nextCalled := false
	handler := CORS("http://localhost:5173")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/api/health", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if !nextCalled {
		t.Error("expected next handler to be called on GET request")
	}
}

// --- Logging tests ---

func TestLogging_CallsNextHandler(t *testing.T) {
	nextCalled := false
	handler := Logging(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/api/health", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if !nextCalled {
		t.Error("expected next handler to be called")
	}
}

func TestLogging_LogsRequestDetails(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	t.Cleanup(func() { log.SetOutput(os.Stderr) })

	handler := Logging(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))

	req := httptest.NewRequest("GET", "/api/game/random", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	logLine := buf.String()
	if !strings.Contains(logLine, "GET") {
		t.Errorf("expected log to contain method GET, got %q", logLine)
	}
	if !strings.Contains(logLine, "/api/game/random") {
		t.Errorf("expected log to contain path, got %q", logLine)
	}
	if !strings.Contains(logLine, "404") {
		t.Errorf("expected log to contain status 404, got %q", logLine)
	}
}
