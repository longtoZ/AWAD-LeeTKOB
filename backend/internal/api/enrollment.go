package api

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"io"
	"net/http"
)

func (s *Server) enroll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := s.Sessions.GetString(ctx, "user_id")
	if user == "" {
		fail(w, 401, "unauthenticated", "Authentication required", "A valid session is required.")
		return
	}
	var role string
	if err := s.DB.QueryRow(ctx, "SELECT role FROM users WHERE id=$1", user).Scan(&role); errors.Is(err, pgx.ErrNoRows) {
		fail(w, 401, "unauthenticated", "Authentication required", "A valid session is required.")
		return
	} else if err != nil {
		internal(w)
		return
	}
	if role != "admin" {
		fail(w, 403, "forbidden", "Forbidden", "Only administrators may enroll students.")
		return
	}
	if !s.validOrigin(r) {
		fail(w, 403, "invalid-origin", "Invalid origin", "The request origin is not allowed.")
		return
	}
	lab := chi.URLParam(r, "labId")
	var in struct {
		StudentID string `json:"studentId"`
	}
	dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, 4096))
	dec.DisallowUnknownFields()
	if _, err := uuid.Parse(lab); err != nil {
		fail(w, 404, "not-found", "Lab not found", "The lab was not found.")
		return
	}
	if dec.Decode(&in) != nil || dec.Decode(&struct{}{}) != io.EOF {
		fail(w, 400, "invalid-request", "Invalid request", "Provide a studentId.")
		return
	}
	if _, err := uuid.Parse(in.StudentID); err != nil {
		fail(w, 400, "invalid-request", "Invalid request", "Provide a valid studentId.")
		return
	}
	tag, err := s.DB.Exec(ctx, "INSERT INTO enrollments (lab_id,student_id) SELECT l.id,u.id FROM labs l CROSS JOIN users u WHERE l.id=$1 AND u.id=$2 AND u.role='student' ON CONFLICT DO NOTHING", lab, in.StudentID)
	if err != nil {
		internal(w)
		return
	}
	if tag.RowsAffected() == 0 {
		var exists bool
		if err := s.DB.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM enrollments WHERE lab_id=$1 AND student_id=$2)", lab, in.StudentID).Scan(&exists); err != nil {
			internal(w)
			return
		}
		if !exists {
			fail(w, 404, "not-found", "Not found", "The lab or student was not found.")
			return
		}
	}
	w.WriteHeader(http.StatusNoContent)
}
