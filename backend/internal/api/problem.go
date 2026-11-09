package api

import (
	"encoding/json"
	"net/http"
)

type problem struct {
	Type   string `json:"type"`
	Title  string `json:"title"`
	Status int    `json:"status"`
	Detail string `json:"detail"`
}

func fail(w http.ResponseWriter, status int, slug, title, detail string) {
	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(problem{"https://leetkob.example/problems/" + slug, title, status, detail})
}
func internal(w http.ResponseWriter) {
	fail(w, 500, "internal", "Internal error", "The request could not be completed.")
}
