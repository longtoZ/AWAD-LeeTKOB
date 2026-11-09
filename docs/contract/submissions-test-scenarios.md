# Submission contract: scenarios for step 5

Use `submissions.md` as the source of truth. Most scenarios below are automated in `backend/internal/api/api_test.go`; an uncertain database commit outcome still needs fault-injection coverage. The table remains a review checklist, not test evidence by itself.

| Area | Scenarios |
|---|---|
| Success | An enrolled student submits valid source; assert `201`, JSON content type, UUIDs, timestamp, `queued`, server-derived ownership, and one submission and job. |
| Authentication | Missing, invalid, and expired sessions return `401 unauthenticated`. |
| Authorization | A TA or admin receives `403 forbidden`; a student not enrolled receives the same `404 not-found` as for a missing lab. |
| Origin | Missing or unapproved origin returns `403 invalid-origin`; approved origin proceeds. |
| Lab | Invalid UUID and absent lab return `404 not-found`. |
| Key | Missing, empty, too long, and invalid-character keys return `400 invalid-idempotency-key`. |
| Media and size | Unsupported content type returns `415`; a body above 256 KiB returns `413`. |
| JSON | Invalid UTF-8, malformed JSON, duplicate names, or multiple values return `400 invalid-json`. |
| Body rules | Non-object, missing, null, non-string, empty, oversized source, and unknown fields return `400 invalid-request`; whitespace-only source is accepted. |
| Problem details | For every error type, check HTTP status, content type, stable type, title, status member, and safe detail. |
| Replay | Same student, lab, key, and decoded JSON replays the original `201` body with the replay header; assert one submission and one job. |
| Conflict and scope | Changed source with a completed key returns `422`; the same key under another student or lab is independent. |
| Concurrent requests | A request with a key already in flight returns `409`; retries after completion replay the result. |
| Failure and recovery | Simulate rollback, unavailable storage, unexpected failure, and uncertain commit; assert atomicity and safe retries. |
| Order | Combine failures to verify the documented check order. |
| Retention | Verify successful keys remain effective for at least 24 hours. |
