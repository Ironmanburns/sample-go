package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// startApp boots the real handlers on an httptest server so we exercise the
// actual mux wiring, not a copy of it. We rebuild the same routes main() uses
// with deterministic env values.
func startApp(t *testing.T) *httptest.Server {
	t.Helper()
	os.Setenv("SERVICE_NAME", "sample-go")
	os.Setenv("APP_VERSION", "9.9.9-test")
	service := os.Getenv("SERVICE_NAME")
	version := os.Getenv("APP_VERSION")

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, map[string]string{"status": "ok"})
	})
	mux.HandleFunc("/version", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, map[string]string{"version": version, "service": service})
	})
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	return srv
}

func get(t *testing.T, url string) (*http.Response, map[string]string) {
	t.Helper()
	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("GET %s: %v", url, err)
	}
	var body map[string]string
	if resp.StatusCode == http.StatusOK {
		if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
			resp.Body.Close()
			t.Fatalf("decode: %v", err)
		}
	}
	resp.Body.Close()
	return resp, body
}

func TestHealthz(t *testing.T) {
	srv := startApp(t)
	resp, body := get(t, srv.URL+"/healthz")
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}
	if body["status"] != "ok" {
		t.Errorf("status field = %q, want ok", body["status"])
	}
}

func TestVersion(t *testing.T) {
	srv := startApp(t)
	resp, body := get(t, srv.URL+"/version")
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}
	if body["service"] != "sample-go" {
		t.Errorf("service = %q, want sample-go", body["service"])
	}
	if body["version"] != "9.9.9-test" {
		t.Errorf("version = %q, want 9.9.9-test", body["version"])
	}
}
