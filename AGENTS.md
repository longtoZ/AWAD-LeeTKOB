# Project rules

## PA2 handoff gate

- Before beginning any remaining PA#2 work, read [`docs/PA2-HANDOFF.md`](docs/PA2-HANDOFF.md). It is the current status, owner, prerequisite, and acceptance-check list for the team.
- Do not mark an item done without its listed evidence. When status changes, update the handoff and the relevant PA2 evidence in the same change.
- The official rubric in `docs/PA#2 rubric — Harness and quality gate.md` remains the scoring authority.

## Stack

- Frontend: React, TypeScript, Vite, and npm in `frontend/`.
- Backend: Go in `backend/`; PostgreSQL runs through Docker Compose.
- Planned judging: a queue and worker with an isolated sandbox. See `docs/PA2-PLAN.md` for planned components.

## Commands

- Install frontend dependencies: `cd frontend` then `npm ci`.
- Run the frontend: `npm run dev` from `frontend/`; use the README command for lint, format, type, and build checks.
- Start the local database: `docker compose up -d` from the root; host port is 5433. See `AI-LOG.md`, entry 1.
- Run the documented Go format check from the root: `node scripts/check-go-format.mjs`.
- Run Go vet and tests from the root: `go -C backend vet ./...` and `go -C backend test ./... -count=1`.
- Run frontend lint, formatting, type-check, and build with `npm run check --prefix frontend` after `npm ci`.
- Run those same commands in CI; do not let a CI step hide failures.

## Conventions

- Keep API behavior consistent with the endpoint contract; enforce role and ownership checks on the server.
- Use `.env.example` for variable names without values; keep local configuration in ignored `.env` files. See `AI-LOG.md`, entry 2.
- Keep model verdicts, evidence, and TA-selected deductions separate; a TA must review before publishing.

## Never

- Never expose secrets, the model API key, the Docker socket, or host mounts to a student-code sandbox.
- Never make real model calls in unit tests; use a fake or stub.
- Never change tests, eval data, or labels merely to make a check pass.
- Never publish marks while required TA reviews are pending.
- Never put environment values in `.env.example`; setup found values there and moved them to ignored `.env` (see `AI-LOG.md`, entry 2).
- Never assume local PostgreSQL can bind host port 5432; a port conflict led to Compose using 5433 (see `AI-LOG.md`, entry 1).

## Permissions

- **Allowed:** Read project files, edit code and documentation for the requested task, and run relevant local checks.
- **Ask first:** Deploy or publish externally, change paid services or spending limits, delete persistent data, or perform destructive Git operations, unless the user already authorized the action.
- **Never:** Commit secrets, bypass required checks or branch protection, or send messages to others without explicit authorization.
