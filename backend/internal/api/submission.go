package api

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"io"
	"mime"
	"net/http"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"
)

var keyPattern = regexp.MustCompile(`^[A-Za-z0-9_-]{1,128}$`)

type submission struct {
	ID        string `json:"id"`
	LabID     string `json:"labId"`
	Status    string `json:"status"`
	CreatedAt string `json:"createdAt"`
}

func checkJSONValue(d *json.Decoder) error {
	token, err := d.Token()
	if err != nil {
		return err
	}
	delim, ok := token.(json.Delim)
	if !ok {
		return nil
	}
	switch delim {
	case '{':
		seen := map[string]bool{}
		for d.More() {
			key, err := d.Token()
			if err != nil {
				return err
			}
			name, ok := key.(string)
			if !ok || seen[name] {
				return errors.New("duplicate or invalid key")
			}
			seen[name] = true
			if err := checkJSONValue(d); err != nil {
				return err
			}
		}
		end, err := d.Token()
		if err != nil || end != json.Delim('}') {
			return errors.New("invalid object")
		}
	case '[':
		for d.More() {
			if err := checkJSONValue(d); err != nil {
				return err
			}
		}
		end, err := d.Token()
		if err != nil || end != json.Delim(']') {
			return errors.New("invalid array")
		}
	default:
		return errors.New("invalid JSON")
	}
	return nil
}

// decodeObject rejects duplicate object keys while preserving JSON values for fingerprinting.
func decodeObject(b []byte) (map[string]json.RawMessage, error) {
	if !utf8.Valid(b) {
		return nil, errors.New("invalid UTF-8")
	}
	validator := json.NewDecoder(bytes.NewReader(b))
	if err := checkJSONValue(validator); err != nil {
		return nil, err
	}
	if _, err := validator.Token(); err != io.EOF {
		return nil, errors.New("trailing JSON")
	}
	d := json.NewDecoder(bytes.NewReader(b))
	tok, err := d.Token()
	if err != nil {
		return nil, err
	}
	if tok != json.Delim('{') {
		var v any
		if json.Unmarshal(b, &v) != nil {
			return nil, errors.New("invalid JSON")
		}
		return map[string]json.RawMessage{}, nil
	}
	result := map[string]json.RawMessage{}
	for d.More() {
		t, err := d.Token()
		if err != nil {
			return nil, err
		}
		k, ok := t.(string)
		if !ok {
			return nil, errors.New("invalid key")
		}
		if _, exists := result[k]; exists {
			return nil, errors.New("duplicate key")
		}
		var v json.RawMessage
		if err := d.Decode(&v); err != nil {
			return nil, err
		}
		result[k] = v
	}
	if _, err := d.Token(); err != nil {
		return nil, err
	}
	if d.Decode(&struct{}{}) != io.EOF {
		return nil, errors.New("trailing JSON")
	}
	return result, nil
}
func (s *Server) submit(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := s.Sessions.GetString(ctx, "user_id")
	if user == "" {
		fail(w, 401, "unauthenticated", "Authentication required", "A valid session is required.")
		return
	}
	var role string
	err := s.DB.QueryRow(ctx, "SELECT role FROM users WHERE id=$1", user).Scan(&role)
	if errors.Is(err, pgx.ErrNoRows) {
		fail(w, 401, "unauthenticated", "Authentication required", "A valid session is required.")
		return
	}
	if err != nil {
		internal(w)
		return
	}
	if role != "student" {
		fail(w, 403, "forbidden", "Forbidden", "Only students may submit.")
		return
	}
	if !s.validOrigin(r) {
		fail(w, 403, "invalid-origin", "Invalid origin", "The request origin is not allowed.")
		return
	}
	labID, err := uuid.Parse(chi.URLParam(r, "labId"))
	if err != nil {
		fail(w, 404, "not-found", "Lab not found", "The lab was not found.")
		return
	}
	lab := labID.String()
	var enrolled bool
	if err := s.DB.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM enrollments WHERE lab_id=$1 AND student_id=$2)", lab, user).Scan(&enrolled); err != nil {
		internal(w)
		return
	}
	if !enrolled {
		fail(w, 404, "not-found", "Lab not found", "The lab was not found.")
		return
	}
	key := r.Header.Get("Idempotency-Key")
	if !keyPattern.MatchString(key) {
		fail(w, 400, "invalid-idempotency-key", "Invalid idempotency key", "A valid idempotency key is required.")
		return
	}
	media, params, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || media != "application/json" || (params["charset"] != "" && !strings.EqualFold(params["charset"], "utf-8")) || (len(params) > 0 && (len(params) != 1 || params["charset"] == "")) {
		fail(w, 415, "unsupported-media-type", "Unsupported media type", "Use application/json with UTF-8.")
		return
	}
	body, err := io.ReadAll(io.LimitReader(r.Body, 256*1024+1))
	if err != nil {
		internal(w)
		return
	}
	if len(body) > 256*1024 {
		fail(w, 413, "payload-too-large", "Payload too large", "The request body is too large.")
		return
	}
	obj, err := decodeObject(body)
	if err != nil {
		fail(w, 400, "invalid-json", "Malformed JSON", "The request body is not valid JSON.")
		return
	}
	var decoded any
	if err := json.Unmarshal(body, &decoded); err != nil {
		fail(w, 400, "invalid-json", "Malformed JSON", "The request body is not valid JSON.")
		return
	}
	canonical, err := json.Marshal(decoded)
	if err != nil {
		internal(w)
		return
	}
	fingerprint := sha256.Sum256(canonical)
	tx, err := s.DB.Begin(ctx)
	if err != nil {
		fail(w, 503, "unavailable", "Service unavailable", "The submission could not be recorded.")
		return
	}
	defer tx.Rollback(context.Background())
	var locked bool
	err = tx.QueryRow(ctx, "SELECT pg_try_advisory_xact_lock(hashtextextended($1, 0))", user+":"+lab+":"+key).Scan(&locked)
	if err != nil {
		internal(w)
		return
	}
	if !locked {
		fail(w, 409, "idempotency-key-in-flight", "Request in progress", "Retry the request with the same key.")
		return
	}
	var storedFingerprint []byte
	var stored []byte
	err = tx.QueryRow(ctx, "SELECT fingerprint, response FROM idempotency_keys WHERE student_id=$1 AND lab_id=$2 AND key=$3", user, lab, key).Scan(&storedFingerprint, &stored)
	if err == nil {
		if !bytes.Equal(storedFingerprint, fingerprint[:]) {
			fail(w, 422, "idempotency-key-reused", "Idempotency key reused", "This key was used for another request.")
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Idempotent-Replayed", "true")
		w.WriteHeader(201)
		_, _ = w.Write(stored)
		return
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		internal(w)
		return
	}
	var source string
	if len(obj) != 1 || obj["source"] == nil || json.Unmarshal(obj["source"], &source) != nil || source == "" || len([]byte(source)) > 128*1024 {
		fail(w, 400, "invalid-request", "Invalid request", "Provide a nonempty source string within the size limit.")
		return
	}
	id := uuid.NewString()
	now := time.Now().UTC().Truncate(time.Microsecond)
	out := submission{id, lab, "queued", now.Format(time.RFC3339Nano)}
	response, _ := json.Marshal(out)
	_, err = tx.Exec(ctx, "INSERT INTO submissions (id,lab_id,student_id,source,created_at) VALUES ($1,$2,$3,$4,$5)", id, lab, user, []byte(source), now)
	if err == nil {
		_, err = tx.Exec(ctx, "INSERT INTO jobs (id,submission_id) VALUES ($1,$2)", uuid.NewString(), id)
	}
	if err == nil {
		_, err = tx.Exec(ctx, "INSERT INTO idempotency_keys (student_id,lab_id,key,fingerprint,response) VALUES ($1,$2,$3,$4,$5)", user, lab, key, fingerprint[:], response)
	}
	if err != nil {
		fail(w, 503, "unavailable", "Service unavailable", "The submission could not be recorded.")
		return
	}
	if err = tx.Commit(ctx); err != nil {
		fail(w, 503, "unavailable", "Service unavailable", "The request outcome is uncertain; retry with the same key.")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	_, _ = w.Write(response)
}
