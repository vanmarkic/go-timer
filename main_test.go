package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestServerServesIndex(t *testing.T) {
	h, err := handler()
	if err != nil {
		t.Fatalf("handler: %v", err)
	}
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d", rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "Countdown Timer") {
		t.Fatalf("index did not contain expected content; got %q", rr.Body.String())
	}
}

func TestServerServesStaticAssets(t *testing.T) {
	h, err := handler()
	if err != nil {
		t.Fatalf("handler: %v", err)
	}
	for _, path := range []string{"/style.css", "/app.js"} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		rr := httptest.NewRecorder()
		h.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Fatalf("%s status = %d", path, rr.Code)
		}
		if rr.Body.Len() == 0 {
			t.Fatalf("%s returned empty body", path)
		}
	}
}

func TestServerReturns404ForMissing(t *testing.T) {
	h, err := handler()
	if err != nil {
		t.Fatalf("handler: %v", err)
	}
	req := httptest.NewRequest(http.MethodGet, "/does-not-exist", nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rr.Code)
	}
}
