package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestPassword(t *testing.T) {
	encoded, err := HashPassword("correct password")
	if err != nil {
		t.Fatal(err)
	}
	if !verifyPassword(encoded, "correct password") || verifyPassword(encoded, "wrong") || verifyPassword("invalid", "wrong") {
		t.Fatal("password verification failed")
	}
}
func TestDecodeObject(t *testing.T) {
	for _, tc := range []struct {
		body    string
		invalid bool
	}{{`{"source":"x"}`, false}, {`{"source":"x","source":"y"}`, true}, {`{"source":`, true}, {`{} {}`, true}, {`[]`, false}, {`{"extra":{"a":1,"a":2}}`, true}} {
		_, err := decodeObject([]byte(tc.body))
		if (err != nil) != tc.invalid {
			t.Errorf("%q: %v", tc.body, err)
		}
	}
}

func TestIdempotencyKeyFormat(t *testing.T) {
	for _, test := range []struct {
		key   string
		valid bool
	}{{"abc-123_", true}, {"", false}, {"a b", false}, {strings.Repeat("x", 128), true}, {strings.Repeat("x", 129), false}} {
		if keyPattern.MatchString(test.key) != test.valid {
			t.Errorf("unexpected validation for key length %d", len(test.key))
		}
	}
}

func TestOriginMustMatchConfiguredValue(t *testing.T) {
	s := &Server{Origin: "https://leetkob.test"}
	r := httptest.NewRequest("POST", "/", nil)
	r.Header.Set("Origin", "https://leetkob.test")
	if !s.validOrigin(r) {
		t.Fatal("expected configured origin")
	}
	r.Header.Set("Origin", "https://other.test")
	if s.validOrigin(r) {
		t.Fatal("unexpected origin accepted")
	}
}

func TestPasswordHashUsesFreshSalt(t *testing.T) {
	one, err := HashPassword("same password")
	if err != nil {
		t.Fatal(err)
	}
	two, err := HashPassword("same password")
	if err != nil {
		t.Fatal(err)
	}
	if one == two || !verifyPassword(one, "same password") || !verifyPassword(two, "same password") {
		t.Fatal("password salts or verification failed")
	}
}

func TestSubmissionContract(t *testing.T) {
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("set TEST_DATABASE_URL to run PostgreSQL integration tests")
	}
	ctx := context.Background()
	base, err := pgxpool.New(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	defer base.Close()
	schema := "test_" + strings.ReplaceAll(uuid.NewString(), "-", "")
	if _, err = base.Exec(ctx, "CREATE SCHEMA "+schema); err != nil {
		t.Fatal(err)
	}
	defer base.Exec(ctx, "DROP SCHEMA "+schema+" CASCADE")
	config, err := pgxpool.ParseConfig(url)
	if err != nil {
		t.Fatal(err)
	}
	config.ConnConfig.RuntimeParams["search_path"] = schema
	db, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	sql, err := os.ReadFile(filepath.Join("..", "..", "migrations", "00001_initial.sql"))
	if err != nil {
		t.Fatal(err)
	}
	up := strings.Split(string(sql), "-- +goose Down")[0]
	if _, err = db.Exec(ctx, up); err != nil {
		t.Fatal(err)
	}
	second, err := os.ReadFile(filepath.Join("..", "..", "migrations", "00002_source_bytes.sql"))
	if err != nil {
		t.Fatal(err)
	}
	if _, err = db.Exec(ctx, strings.Split(string(second), "-- +goose Down")[0]); err != nil {
		t.Fatal(err)
	}
	student, other, ta, admin, lab := uuid.NewString(), uuid.NewString(), uuid.NewString(), uuid.NewString(), uuid.NewString()
	password, err := HashPassword("test-password")
	if err != nil {
		t.Fatal(err)
	}
	if _, err = db.Exec(ctx, "INSERT INTO users (id,email,password_hash,role) VALUES ($1,'student@example.test',$2,'student'),($3,'other@example.test',$2,'student'),($4,'ta@example.test',$2,'ta'),($5,'admin@example.test',$2,'admin')", student, password, other, ta, admin); err != nil {
		t.Fatal(err)
	}
	if _, err = db.Exec(ctx, "INSERT INTO labs VALUES ($1,'Lab')", lab); err != nil {
		t.Fatal(err)
	}
	if _, err = db.Exec(ctx, "INSERT INTO enrollments VALUES ($1,$2)", lab, student); err != nil {
		t.Fatal(err)
	}
	store := pgxstore.NewWithCleanupInterval(db, 0)
	defer store.StopCleanup()
	sessions := scs.New()
	sessions.Store = store
	h := (&Server{DB: db, Sessions: sessions, Origin: "http://localhost:5173"}).Routes()
	request := func(method, path, body string, cookie *http.Cookie, key, origin string) *httptest.ResponseRecorder {
		req := httptest.NewRequest(method, path, strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		if key != "" {
			req.Header.Set("Idempotency-Key", key)
		}
		if origin != "" {
			req.Header.Set("Origin", origin)
		}
		if cookie != nil {
			req.AddCookie(cookie)
		}
		out := httptest.NewRecorder()
		h.ServeHTTP(out, req)
		if out.Code >= 400 {
			var detail problem
			if out.Header().Get("Content-Type") != "application/problem+json" || json.Unmarshal(out.Body.Bytes(), &detail) != nil || detail.Status != out.Code || detail.Title == "" || detail.Detail == "" || !strings.HasPrefix(detail.Type, "https://leetkob.example/problems/") {
				t.Fatalf("invalid problem response: %d %s", out.Code, out.Body.String())
			}
		}
		return out
	}
	login := func(email string) *http.Cookie {
		out := request("POST", "/auth/login", fmt.Sprintf(`{"email":%q,"password":"test-password"}`, email), nil, "", "http://localhost:5173")
		if out.Code != 204 {
			t.Fatalf("login: %d %s", out.Code, out.Body.String())
		}
		return out.Result().Cookies()[0]
	}
	cookie := login("student@example.test")
	otherCookie := login("other@example.test")
	taCookie := login("ta@example.test")
	adminCookie := login("admin@example.test")
	path := "/labs/" + lab + "/submissions"
	origin := "http://localhost:5173"
	t.Run("unauthenticated", func(t *testing.T) {
		if out := request("POST", path, `{"source":"x"}`, nil, "key", origin); out.Code != 401 {
			t.Fatal(out.Code)
		}
	})
	t.Run("enrollment", func(t *testing.T) {
		if out := request("POST", path, `{"source":"x"}`, otherCookie, "key", origin); out.Code != 404 {
			t.Fatal(out.Code)
		}
	})
	t.Run("non-student", func(t *testing.T) {
		out := request("POST", path, `{"source":"x"}`, taCookie, "role", origin)
		if out.Code != 403 || !strings.Contains(out.Body.String(), `"type":"https://leetkob.example/problems/forbidden"`) {
			t.Fatal(out.Code, out.Body.String())
		}
	})
	t.Run("body size limit", func(t *testing.T) {
		out := request("POST", path, strings.Repeat("x", 256*1024+1), cookie, "large", origin)
		if out.Code != 413 {
			t.Fatal(out.Code)
		}
	})
	t.Run("media type", func(t *testing.T) {
		req := httptest.NewRequest("POST", path, strings.NewReader(`{"source":"x"}`))
		req.AddCookie(cookie)
		req.Header.Set("Origin", origin)
		req.Header.Set("Idempotency-Key", "media")
		req.Header.Set("Content-Type", "text/plain")
		out := httptest.NewRecorder()
		h.ServeHTTP(out, req)
		if out.Code != 415 {
			t.Fatal(out.Code)
		}
	})
	t.Run("in-flight key", func(t *testing.T) {
		tx, err := db.Begin(ctx)
		if err != nil {
			t.Fatal(err)
		}
		defer tx.Rollback(ctx)
		if _, err = tx.Exec(ctx, "SELECT pg_advisory_xact_lock(hashtextextended($1,0))", student+":"+lab+":busy"); err != nil {
			t.Fatal(err)
		}
		out := request("POST", path, `{"source":"x"}`, cookie, "busy", origin)
		if out.Code != 409 {
			t.Fatal(out.Code, out.Body.String())
		}
	})
	t.Run("job insert rollback", func(t *testing.T) {
		if _, err := db.Exec(ctx, "ALTER TABLE jobs ADD CONSTRAINT reject_test_job CHECK (false) NOT VALID"); err != nil {
			t.Fatal(err)
		}
		defer db.Exec(ctx, "ALTER TABLE jobs DROP CONSTRAINT reject_test_job")
		out := request("POST", path, `{"source":"rollback-test"}`, cookie, "rollback", origin)
		if out.Code != 503 {
			t.Fatal(out.Code, out.Body.String())
		}
		var count int
		if err := db.QueryRow(ctx, "SELECT count(*) FROM submissions WHERE source=$1", []byte("rollback-test")).Scan(&count); err != nil || count != 0 {
			t.Fatal(count, err)
		}
	})
	t.Run("unexpected database error", func(t *testing.T) {
		if _, err := db.Exec(ctx, "ALTER TABLE users RENAME TO users_unavailable"); err != nil {
			t.Fatal(err)
		}
		defer db.Exec(ctx, "ALTER TABLE users_unavailable RENAME TO users")
		out := request("POST", path, `{"source":"x"}`, cookie, "error", origin)
		if out.Code != 500 || !strings.Contains(out.Body.String(), "/internal") {
			t.Fatal(out.Code, out.Body.String())
		}
	})
	t.Run("origin", func(t *testing.T) {
		if out := request("POST", path, `{"source":"x"}`, cookie, "key", ""); out.Code != 403 {
			t.Fatal(out.Code)
		}
	})
	t.Run("missing key", func(t *testing.T) {
		if out := request("POST", path, `{"source":"x"}`, cookie, "", origin); out.Code != 400 {
			t.Fatal(out.Code)
		}
	})
	t.Run("malformed JSON", func(t *testing.T) {
		if out := request("POST", path, `{`, cookie, "malformed", origin); out.Code != 400 {
			t.Fatal(out.Code)
		}
	})
	t.Run("invalid field", func(t *testing.T) {
		if out := request("POST", path, `{"source":""}`, cookie, "invalid", origin); out.Code != 400 {
			t.Fatal(out.Code)
		}
	})
	t.Run("creation and replay", func(t *testing.T) {
		first := request("POST", path, `{"source":"hello"}`, cookie, "same", origin)
		if first.Code != 201 {
			t.Fatalf("%d %s", first.Code, first.Body.String())
		}
		var resource submission
		if err := json.Unmarshal(first.Body.Bytes(), &resource); err != nil || resource.Status != "queued" || resource.LabID != lab {
			t.Fatal("bad response")
		}
		retry := request("POST", path, `{ "source": "hello" }`, cookie, "same", origin)
		if retry.Code != 201 || retry.Header().Get("Idempotent-Replayed") != "true" || !bytes.Equal(first.Body.Bytes(), retry.Body.Bytes()) {
			t.Fatalf("replay: %d %s", retry.Code, retry.Body.String())
		}
		var count int
		if err := db.QueryRow(ctx, "SELECT count(*) FROM jobs WHERE submission_id=$1", resource.ID).Scan(&count); err != nil || count != 1 {
			t.Fatalf("job count: %d, %v", count, err)
		}
	})
	t.Run("key reuse", func(t *testing.T) {
		if out := request("POST", path, `{"source":"different"}`, cookie, "same", origin); out.Code != 422 {
			t.Fatal(out.Code)
		}
	})
	t.Run("source bytes preserved", func(t *testing.T) {
		out := request("POST", path, `{"source":"a\u0000b\n"}`, cookie, "source-bytes", origin)
		if out.Code != 201 {
			t.Fatal(out.Code, out.Body.String())
		}
		var result submission
		if err := json.Unmarshal(out.Body.Bytes(), &result); err != nil {
			t.Fatal(err)
		}
		var source []byte
		if err := db.QueryRow(ctx, "SELECT source FROM submissions WHERE id=$1", result.ID).Scan(&source); err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(source, []byte("a\x00b\n")) {
			t.Fatal("source bytes changed")
		}
	})
	t.Run("admin enrolls student", func(t *testing.T) {
		enrollPath := "/labs/" + lab + "/enrollments"
		body := fmt.Sprintf(`{"studentId":%q}`, other)
		if out := request("POST", enrollPath, body, otherCookie, "", origin); out.Code != 403 {
			t.Fatal(out.Code)
		}
		if out := request("POST", enrollPath, body, adminCookie, "", origin); out.Code != 204 {
			t.Fatal(out.Code, out.Body.String())
		}
		if out := request("POST", path, `{"source":"new"}`, otherCookie, "enrolled", origin); out.Code != 201 {
			t.Fatal(out.Code, out.Body.String())
		}
	})
	t.Run("logout", func(t *testing.T) {
		out := request("POST", "/auth/logout", "", cookie, "", origin)
		if out.Code != 204 {
			t.Fatal(out.Code)
		}
		out = request("POST", path, `{"source":"x"}`, cookie, "after-logout", origin)
		if out.Code != 401 {
			t.Fatal(out.Code)
		}
	})
}
