# H1/H2 submission endpoint verification — 9 November 2026

The existing handler and integration tests were present on `main` at `9f25fab`. The first implementation's pre-code failing test output was not retained. This follow-up cannot recreate that historical evidence.

## Observed follow-up red/green

The new `TestOriginMustMatchConfiguredValue` case sent two `Origin` header values, the first allowed and the second unapproved. Before changing the handler, the test failed:

```text
$ go test ./internal/api -run '^TestOriginMustMatchConfiguredValue$' -count=1
--- FAIL: TestOriginMustMatchConfiguredValue (0.00s)
    api_test.go:67: multiple Origin values accepted
FAIL
FAIL    leetkob/backend/internal/api    1.360s
FAIL
```

`validOrigin` now requires exactly one `Origin` value matching the configured origin. The same test then passed. The handler's session, role, enrollment, and idempotency logic remained unchanged.

## PostgreSQL contract run

Docker Desktop was unavailable on this machine. An isolated local PostgreSQL 18 cluster on `127.0.0.1:55439` supplied a disposable UTF-8 test database. `TestSubmissionContract` creates and removes its own schema, so it does not use application data.

```text
$env:TEST_DATABASE_URL='postgres://pa2test@127.0.0.1:55439/pa2test?sslmode=disable'
go test ./... -count=1
?   leetkob/backend/cmd/create-user  [no test files]
?   leetkob/backend/cmd/migrate      [no test files]
?   leetkob/backend/cmd/server       [no test files]
ok  leetkob/backend/internal/api     6.489s
```

The expanded tests assert every error type's status, stable type, and title; authentication and enrollment; check order; source and body boundaries; successful ownership and one durable pending job; replay and key scope; 23-hour retention; rollback of submission, job, and key; and retry after rollback. The deliberate red/green above is local test evidence. Required CI and a blocked break PR are separate L3/H3 tasks and are not established by this run.

`go vet ./...` also exited successfully with no findings.
