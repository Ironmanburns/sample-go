package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

func main() {
	port := env("PORT", "8080")
	service := env("SERVICE_NAME", "sample-go")
	version := env("APP_VERSION", "0.1.0")

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, map[string]string{"status": "ok"})
	})
	mux.HandleFunc("/version", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, map[string]string{"version": version, "service": service})
	})

	log.Printf("%s ready on %s", service, port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

func env(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("content-type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}
