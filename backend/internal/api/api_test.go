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
	"time"

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
	r.Header.Set("Origin", "https://leetkob.test")
	r.Header.Add("Origin", "https://other.test")
	if s.validOrigin(r) {
		t.Fatal("multiple Origin values accepted")
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
	expectProblem := func(t *testing.T, out *httptest.ResponseRecorder, status int, slug, title string) {
		t.Helper()
		var got problem
		if err := json.Unmarshal(out.Body.Bytes(), &got); err != nil {
			t.Fatal(err)
		}
		if out.Code != status || out.Header().Get("Content-Type") != "application/problem+json" || got.Type != "https://leetkob.example/problems/"+slug || got.Title != title || got.Status != status || got.Detail == "" {
			t.Fatalf("unexpected problem: HTTP %d, body %s", out.Code, out.Body.String())
		}
	}
	// Contract §4: a refusal creates nothing, and a replay creates no other
	// submission or job. A status check alone cannot see a write that happened
	// before the refusal, so these tests also read the database.
	type rowCounts struct{ submissions, jobs, keys int }
	countRows := func(t *testing.T) rowCounts {
		t.Helper()
		var c rowCounts
		if err := db.QueryRow(ctx, "SELECT (SELECT count(*) FROM submissions), (SELECT count(*) FROM jobs), (SELECT count(*) FROM idempotency_keys)").Scan(&c.submissions, &c.jobs, &c.keys); err != nil {
			t.Fatal(err)
		}
		return c
	}
	expectNothingRecorded := func(t *testing.T, before rowCounts) {
		t.Helper()
		if after := countRows(t); after != before {
			t.Fatalf("stored rows changed: before %+v, after %+v", before, after)
		}
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
		expectProblem(t, out, 415, "unsupported-media-type", "Unsupported media type")
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
		before := countRows(t)
		out := request("POST", path, `{"source":"x"}`, cookie, "busy", origin)
		expectProblem(t, out, 409, "idempotency-key-in-flight", "Request in progress")
		expectNothingRecorded(t, before)
	})
	t.Run("job insert rollback", func(t *testing.T) {
		if _, err := db.Exec(ctx, "ALTER TABLE jobs ADD CONSTRAINT reject_test_job CHECK (false) NOT VALID"); err != nil {
			t.Fatal(err)
		}
		defer db.Exec(ctx, "ALTER TABLE jobs DROP CONSTRAINT reject_test_job")
		out := request("POST", path, `{"source":"rollback-test"}`, cookie, "rollback", origin)
		expectProblem(t, out, 503, "unavailable", "Service unavailable")
		for _, query := range []string{
			"SELECT count(*) FROM submissions WHERE source=$1",
			"SELECT count(*) FROM jobs WHERE submission_id IN (SELECT id FROM submissions WHERE source=$1)",
		} {
			var count int
			if err := db.QueryRow(ctx, query, []byte("rollback-test")).Scan(&count); err != nil || count != 0 {
				t.Fatal(query, count, err)
			}
		}
		var keyCount int
		if err := db.QueryRow(ctx, "SELECT count(*) FROM idempotency_keys WHERE student_id=$1 AND lab_id=$2 AND key='rollback'", student, lab).Scan(&keyCount); err != nil || keyCount != 0 {
			t.Fatal("rolled-back key retained", keyCount, err)
		}
	})
	t.Run("retry after rollback", func(t *testing.T) {
		out := request("POST", path, `{"source":"rollback-test"}`, cookie, "rollback", origin)
		if out.Code != 201 {
			t.Fatal(out.Code, out.Body.String())
		}
	})
	t.Run("unexpected database error", func(t *testing.T) {
		if _, err := db.Exec(ctx, "ALTER TABLE users RENAME TO users_unavailable"); err != nil {
			t.Fatal(err)
		}
		defer db.Exec(ctx, "ALTER TABLE users_unavailable RENAME TO users")
		out := request("POST", path, `{"source":"x"}`, cookie, "error", origin)
		expectProblem(t, out, 500, "internal", "Internal error")
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
	t.Run("contract error matrix", func(t *testing.T) {
		for _, tc := range []struct {
			name, path, body, key, origin, slug, title string
			cookie                                     *http.Cookie
			status                                     int
		}{
			{"unauthenticated", path, `{"source":"x"}`, "case-auth", origin, "unauthenticated", "Authentication required", nil, 401},
			{"non-student admin", path, `{"source":"x"}`, "case-role", origin, "forbidden", "Forbidden", adminCookie, 403},
			{"unenrolled", path, `{"source":"x"}`, "case-enrollment", origin, "not-found", "Lab not found", otherCookie, 404},
			{"missing lab", "/labs/" + uuid.NewString() + "/submissions", `{"source":"x"}`, "case-lab", origin, "not-found", "Lab not found", cookie, 404},
			{"invalid lab id", "/labs/not-a-uuid/submissions", `{"source":"x"}`, "case-lab-id", origin, "not-found", "Lab not found", cookie, 404},
			{"missing origin", path, `{"source":"x"}`, "case-origin", "", "invalid-origin", "Invalid origin", cookie, 403},
			{"wrong origin", path, `{"source":"x"}`, "case-origin", "https://other.test", "invalid-origin", "Invalid origin", cookie, 403},
			{"missing key", path, `{"source":"x"}`, "", origin, "invalid-idempotency-key", "Invalid idempotency key", cookie, 400},
			{"invalid key", path, `{"source":"x"}`, "invalid key", origin, "invalid-idempotency-key", "Invalid idempotency key", cookie, 400},
			{"overlong key", path, `{"source":"x"}`, strings.Repeat("x", 129), origin, "invalid-idempotency-key", "Invalid idempotency key", cookie, 400},
			{"oversized body", path, strings.Repeat("x", 256*1024+1), "case-large", origin, "payload-too-large", "Payload too large", cookie, 413},
			{"malformed json", path, `{`, "case-json", origin, "invalid-json", "Malformed JSON", cookie, 400},
			{"duplicate key", path, `{"source":"x","source":"y"}`, "case-duplicate", origin, "invalid-json", "Malformed JSON", cookie, 400},
			{"multiple values", path, `{"source":"x"} {}`, "case-trailing", origin, "invalid-json", "Malformed JSON", cookie, 400},
			{"invalid UTF-8", path, string([]byte{'{', '"', 's', 'o', 'u', 'r', 'c', 'e', '"', ':', '"', 0xff, '"', '}'}), "case-utf8", origin, "invalid-json", "Malformed JSON", cookie, 400},
			{"non-object", path, `[]`, "case-array", origin, "invalid-request", "Invalid request", cookie, 400},
			{"missing source", path, `{}`, "case-missing-source", origin, "invalid-request", "Invalid request", cookie, 400},
			{"null source", path, `{"source":null}`, "case-null", origin, "invalid-request", "Invalid request", cookie, 400},
			{"non-string source", path, `{"source":1}`, "case-number", origin, "invalid-request", "Invalid request", cookie, 400},
			{"empty source", path, `{"source":""}`, "case-empty", origin, "invalid-request", "Invalid request", cookie, 400},
			{"unknown field", path, `{"source":"x","owner":"other"}`, "case-unknown", origin, "invalid-request", "Invalid request", cookie, 400},
			{"oversized source", path, fmt.Sprintf(`{"source":%q}`, strings.Repeat("x", 128*1024+1)), "case-source-large", origin, "invalid-request", "Invalid request", cookie, 400},
		} {
			t.Run(tc.name, func(t *testing.T) {
				before := countRows(t)
				expectProblem(t, request("POST", tc.path, tc.body, tc.cookie, tc.key, tc.origin), tc.status, tc.slug, tc.title)
				expectNothingRecorded(t, before)
			})
		}
	})
	t.Run("check order stops at the first failure", func(t *testing.T) {
		before := countRows(t)
		defer expectNothingRecorded(t, before)
		expectProblem(t, request("POST", "/labs/not-a-uuid/submissions", `{`, nil, "", ""), 401, "unauthenticated", "Authentication required")
		expectProblem(t, request("POST", "/labs/not-a-uuid/submissions", `{`, taCookie, "", ""), 403, "forbidden", "Forbidden")
		expectProblem(t, request("POST", "/labs/not-a-uuid/submissions", `{`, cookie, "", ""), 403, "invalid-origin", "Invalid origin")
		expectProblem(t, request("POST", "/labs/not-a-uuid/submissions", `{`, cookie, "", origin), 404, "not-found", "Lab not found")
		expectProblem(t, request("POST", path, `{`, cookie, "", origin), 400, "invalid-idempotency-key", "Invalid idempotency key")
		expectProblem(t, request("POST", path, `{`, cookie, "order-json", origin), 400, "invalid-json", "Malformed JSON")
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
		if first.Code != 201 || first.Header().Get("Content-Type") != "application/json" || first.Header().Get("Idempotent-Replayed") != "" {
			t.Fatalf("%d %s", first.Code, first.Body.String())
		}
		var resource submission
		if err := json.Unmarshal(first.Body.Bytes(), &resource); err != nil || resource.Status != "queued" || resource.LabID != lab {
			t.Fatal("bad response", err)
		}
		if _, err := uuid.Parse(resource.ID); err != nil {
			t.Fatal("submission ID is not a UUID", err)
		}
		if parsed, err := time.Parse(time.RFC3339Nano, resource.CreatedAt); err != nil || parsed.Location() != time.UTC {
			t.Fatal("createdAt is not RFC 3339 UTC", resource.CreatedAt, err)
		}
		var owner, status string
		var source []byte
		if err := db.QueryRow(ctx, "SELECT student_id, source FROM submissions WHERE id=$1", resource.ID).Scan(&owner, &source); err != nil || owner != student || !bytes.Equal(source, []byte("hello")) {
			t.Fatal("submission ownership or source changed", owner, err)
		}
		if err := db.QueryRow(ctx, "SELECT status FROM jobs WHERE submission_id=$1", resource.ID).Scan(&status); err != nil || status != "pending" {
			t.Fatal("durable pending job missing", status, err)
		}
		beforeRetry := countRows(t)
		retry := request("POST", path, `{ "source": "hello" }`, cookie, "same", origin)
		if retry.Code != 201 || retry.Header().Get("Idempotent-Replayed") != "true" || !bytes.Equal(first.Body.Bytes(), retry.Body.Bytes()) {
			t.Fatalf("replay: %d %s", retry.Code, retry.Body.String())
		}
		expectNothingRecorded(t, beforeRetry)
		var count int
		if err := db.QueryRow(ctx, "SELECT count(*) FROM jobs WHERE submission_id=$1", resource.ID).Scan(&count); err != nil || count != 1 {
			t.Fatalf("job count: %d, %v", count, err)
		}
	})
	t.Run("key reuse", func(t *testing.T) {
		before := countRows(t)
		expectProblem(t, request("POST", path, `{"source":"different"}`, cookie, "same", origin), 422, "idempotency-key-reused", "Idempotency key reused")
		expectNothingRecorded(t, before)
	})
	t.Run("key retained after 23 hours", func(t *testing.T) {
		if _, err := db.Exec(ctx, "UPDATE idempotency_keys SET created_at=now()-interval '23 hours' WHERE student_id=$1 AND lab_id=$2 AND key='same'", student, lab); err != nil {
			t.Fatal(err)
		}
		before := countRows(t)
		out := request("POST", path, `{"source":"hello"}`, cookie, "same", origin)
		if out.Code != 201 || out.Header().Get("Idempotent-Replayed") != "true" {
			t.Fatal(out.Code, out.Body.String())
		}
		expectNothingRecorded(t, before)
	})
	t.Run("same key in another lab", func(t *testing.T) {
		secondLab := uuid.NewString()
		if _, err := db.Exec(ctx, "INSERT INTO labs VALUES ($1,'Second lab')", secondLab); err != nil {
			t.Fatal(err)
		}
		if _, err := db.Exec(ctx, "INSERT INTO enrollments VALUES ($1,$2)", secondLab, student); err != nil {
			t.Fatal(err)
		}
		out := request("POST", "/labs/"+secondLab+"/submissions", `{"source":"hello"}`, cookie, "same", origin)
		if out.Code != 201 || out.Header().Get("Idempotent-Replayed") != "" {
			t.Fatal(out.Code, out.Body.String())
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
	t.Run("whitespace source and 128 KiB boundary", func(t *testing.T) {
		for _, tc := range []struct{ key, source string }{
			{"whitespace", " \n\t"},
			{"source-boundary", strings.Repeat("x", 128*1024)},
		} {
			out := request("POST", path, fmt.Sprintf(`{"source":%q}`, tc.source), cookie, tc.key, origin)
			if out.Code != 201 {
				t.Fatal(tc.key, out.Code, out.Body.String())
			}
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
		out := request("POST", path, `{"source":"hello"}`, otherCookie, "same", origin)
		if out.Code != 201 || out.Header().Get("Idempotent-Replayed") != "" {
			t.Fatal("key incorrectly shared across students", out.Code, out.Body.String())
		}
	})
	t.Run("invalid and expired sessions", func(t *testing.T) {
		invalid := &http.Cookie{Name: cookie.Name, Value: "invalid-session-token"}
		expectProblem(t, request("POST", path, `{"source":"x"}`, invalid, "invalid-session", origin), 401, "unauthenticated", "Authentication required")
		expired := login("student@example.test")
		if _, err := db.Exec(ctx, "UPDATE sessions SET expiry=now()-interval '1 second' WHERE token=$1", expired.Value); err != nil {
			t.Fatal(err)
		}
		expectProblem(t, request("POST", path, `{"source":"x"}`, expired, "expired-session", origin), 401, "unauthenticated", "Authentication required")
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
