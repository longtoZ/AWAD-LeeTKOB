package api

import (
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
)

type Server struct {
	DB       *pgxpool.Pool
	Sessions *scs.SessionManager
	Origin   string
}

func (s *Server) Routes() http.Handler {
	r := chi.NewRouter()
	r.Post("/auth/login", s.login)
	r.Post("/auth/logout", s.logout)
	r.Post("/labs/{labId}/enrollments", s.enroll)
	r.Post("/labs/{labId}/submissions", s.submit)
	return s.Sessions.LoadAndSave(r)
}

func (s *Server) validOrigin(r *http.Request) bool {
	return r.Header.Get("Origin") == s.Origin && s.Origin != ""
}
