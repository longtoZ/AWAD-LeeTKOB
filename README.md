# LeeTKOB

LeeTKOB supports C++ lab judging and constraint review before a TA publishes marks. See [the project brief](docs/LeeTKOB%20%E2%80%94%20class%20project%20brief.md) and [PA2 plan](docs/PA2-PLAN.md).

## Stack

- Frontend: React, TypeScript, Vite, and npm.
- Backend: Go, Chi, and SCS sessions stored in PostgreSQL.
- Database: PostgreSQL, started locally with Docker Compose.

## Local setup

Requirements: Node.js and npm, Go 1.26 or newer, Docker with Compose, and Git. On Windows PowerShell, use `npm.cmd` if script execution prevents `npm` from running.

1. Clone `https://github.com/longtoZ/AWAD-LeeTKOB.git`.
2. Copy `.env.example` to `.env` and set `POSTGRES_DB`, `POSTGRES_USER`, `POSTGRES_PASSWORD`, `DATABASE_URL`, and `API_PORT` (normally `8080`). Use host port `5433` in `DATABASE_URL`. Set `APP_ORIGIN=http://localhost:5173` and `COOKIE_SECURE=false` for local HTTP. Use HTTPS and `COOKIE_SECURE=true` in deployment. Never commit `.env`.
3. Run `docker compose up -d` from the repository root to start PostgreSQL.
4. In `frontend`, run `npm ci`, then `npm run dev`. Vite prints the local URL.
5. Run `npm run build` in `frontend` to build the frontend.

## Backend

```text
backend/
├── cmd/
│   ├── server/          # API entry point
│   ├── migrate/         # Database migration command
│   └── create-user/     # Initial account command
├── internal/api/
│   ├── server.go        # Router and shared server dependencies
│   ├── auth.go          # Login and logout handlers
│   ├── enrollment.go    # Enrollment handler
│   ├── submission.go   # Submission and idempotency handling
│   ├── problem.go      # Problem details responses
│   ├── password.go     # Password hashing
│   └── api_test.go     # Unit and HTTP contract checks
├── migrations/          # Versioned PostgreSQL schema changes
├── go.mod
└── go.sum
```

Files in `internal/api` share the `api` package. The SQL used by each handler stays alongside that handler for now.

The Go programs read shell environment variables; they do not load `.env` automatically. In PowerShell, load the local file from the repository root:

```powershell
Get-Content .env | Where-Object { $_ -match '^[A-Za-z_][A-Za-z0-9_]*=' } | ForEach-Object {
  $name, $value = $_ -split '=', 2
  [Environment]::SetEnvironmentVariable($name, $value, 'Process')
}
cd backend
go run ./cmd/migrate
go run ./cmd/server
```

From `backend`, run `go test ./...` for unit checks. Set `TEST_DATABASE_URL` to a PostgreSQL test database to include integration checks; these create and remove an isolated temporary schema. CI must set this variable so database contract tests are not skipped.

For initial development accounts, set `USER_EMAIL`, `USER_PASSWORD` (at least 12 bytes), and `USER_ROLE` (`admin`, `ta`, or `student`) in the shell, then run `go run ./cmd/create-user`. Remove the password environment variable afterwards. Labs currently need to be created in PostgreSQL; an authenticated admin can enroll a student using `POST /labs/{labId}/enrollments` with JSON `{ "studentId": "<uuid>" }`.

`POST /auth/login` accepts `email` and `password`; `POST /auth/logout` ends the session. These endpoints and enrollment require the configured `Origin` header. Submission details are in [the contract](docs/contract/submissions.md). No registration UI, judging worker, or sandbox is implemented yet.

For a walkthrough with example requests, see [manual backend testing](docs/MANUAL-TESTING.md).

## Model and cost

Gemini 3.8 Flash is the planned evaluation candidate. The proposed project budget is US$10; it has **not yet been configured or verified as an enforced spending cap**. No model integration is implemented yet.

At Gemini 3.8 Flash introductory standard rates of US$0.75 per million input tokens and US$3.75 per million output tokens, 80 calls with 2,000 input and 500 output tokens each would cost approximately **US$0.27 per lab**. Output pricing includes thinking tokens. This is an estimate, not a measured cost or spending limit; retries, usage, other charges, and pricing changes can change it. The introductory rates apply through 31 December 2026; check [current pricing and data terms](https://ai.google.dev/gemini-api/docs/pricing) before use.

The team reports AI Studio limits of 1,000 requests per minute and 10,000 requests per day for its account. These limit request frequency, not total spending; active limits can change.

## PA2 checks

Use Node.js 22.13.1 and Go 1.26 or newer. From the repository root, run the same checks used by the `gates` workflow:

| Check | Command |
|---|---|
| Frontend lint, format, type-check, and build | `npm ci --prefix frontend && npm run check --prefix frontend` |
| Go format check | `node scripts/check-go-format.mjs` |
| Go vet | `go -C backend vet ./...` |
| Backend unit and PostgreSQL contract tests | `go -C backend test ./... -count=1 -v` |
| Dependency audit | `npm audit --prefix frontend` |
| Browser smoke test | `npm ci --prefix frontend && npx --prefix frontend playwright install chromium && npm run e2e --prefix frontend` |

The database-backed Go tests require `TEST_DATABASE_URL`; CI starts PostgreSQL and sets it so the integration suite cannot silently skip. The browser smoke test checks that the app shell renders. It does not replace API contract tests.

CI runs the `gates`, `secret-history`, and `e2e` jobs on pull requests, pushes to `main`, and the `pa2` tag. Configure `gates`, `secret-history`, and `e2e` as required checks in GitHub branch protection. A passing workflow run alone does not configure branch protection.

`npm audit` checks the frontend dependency tree. The evaluated lint-toolchain versions in `package.json` use a targeted `minimatch` override to avoid the advisory in the compatible upstream dependency range; review the override when upgrading that toolchain.

## Spending cap status

The model integration is not implemented. The proposed US$10 budget is not an enforced provider limit. Before connecting a model, configure a provider-side cap or an application-side hard limit, verify it, and update this section with the measured cost boundary. Do not treat the request limits or illustrative US$0.27 estimate above as a cap.
