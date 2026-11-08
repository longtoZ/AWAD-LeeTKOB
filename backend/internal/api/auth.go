package api

import (
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v5"
	"io"
	"net/http"
	"strings"
)

func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	if !s.validOrigin(r) {
		fail(w, 403, "invalid-origin", "Invalid origin", "The request origin is not allowed.")
		return
	}
	var in struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, 4096))
	dec.DisallowUnknownFields()
	if dec.Decode(&in) != nil || in.Email == "" || in.Password == "" || dec.Decode(&struct{}{}) != io.EOF {
		fail(w, 400, "invalid-request", "Invalid request", "Email and password are required.")
		return
	}
	var id, hash string
	err := s.DB.QueryRow(r.Context(), "SELECT id::text, password_hash FROM users WHERE email=$1", strings.ToLower(strings.TrimSpace(in.Email))).Scan(&id, &hash)
	if errors.Is(err, pgx.ErrNoRows) || (err == nil && !verifyPassword(hash, in.Password)) {
		fail(w, 401, "unauthenticated", "Authentication required", "Invalid credentials.")
		return
	}
	if err != nil {
		internal(w)
		return
	}
	if err := s.Sessions.RenewToken(r.Context()); err != nil {
		internal(w)
		return
	}
	s.Sessions.Put(r.Context(), "user_id", id)
	w.WriteHeader(http.StatusNoContent)
}
func (s *Server) logout(w http.ResponseWriter, r *http.Request) {
	if !s.validOrigin(r) {
		fail(w, 403, "invalid-origin", "Invalid origin", "The request origin is not allowed.")
		return
	}
	if err := s.Sessions.Destroy(r.Context()); err != nil {
		internal(w)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
