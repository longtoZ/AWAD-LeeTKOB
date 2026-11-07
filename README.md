# LeeTKOB

LeeTKOB supports C++ lab judging and constraint review before a TA publishes marks. See [the project brief](docs/LeeTKOB%20%E2%80%94%20class%20project%20brief.md) and [PA2 plan](docs/PA2-PLAN.md).

## Stack

- Frontend: React, TypeScript, Vite, and npm.
- Backend: Go; the API is planned but not implemented yet.
- Database: PostgreSQL, started locally with Docker Compose.

## Local setup

Requirements: Node.js and npm, Go 1.23 or newer, Docker with Compose, and Git. On Windows PowerShell, use `npm.cmd` if script execution prevents `npm` from running.

1. Clone `https://github.com/longtoZ/AWAD-LeeTKOB.git`.
2. Copy `.env.example` to `.env`. Set `POSTGRES_DB`, `POSTGRES_USER`, and a local `POSTGRES_PASSWORD`; use those values in `DATABASE_URL` (for example, `postgres://user:password@localhost:5433/leetkob?sslmode=disable`). Set `API_PORT` to `8080`. These variables are planned for the API; the API does not load them yet. Never commit `.env`.
3. Run `docker compose up -d` from the repository root to start PostgreSQL.
4. In `frontend`, run `npm ci`, then `npm run dev`. Vite prints the local URL.
5. Run `npm run build` in `frontend` to build the frontend. The Go API and its commands will be added with the endpoint.

## Model and cost

Gemini 3.8 Flash is the planned evaluation candidate. The proposed project budget is US$10; it has **not yet been configured or verified as an enforced spending cap**. No model integration is implemented yet.

An illustrative estimate for 80 model calls, each using 2,000 input and 500 output tokens, is **US$0.10 per lab** at published standard rates of $0.25 per million input tokens and $1.50 per million output tokens. Output includes any billed thinking tokens. This is an assumption, not a measured cost or an upper bound: constraints, retries, token usage, and pricing can change the total. Measure actual use on the dev set and check [provider pricing and data terms](https://ai.google.dev/gemini-api/docs/pricing) before sending student code.

## PA2 checks

Lint, format, type-check, and test commands and CI are part of the next steps in [PA2-PLAN.md](docs/PA2-PLAN.md). The frontend build above is available once dependencies are installed.
