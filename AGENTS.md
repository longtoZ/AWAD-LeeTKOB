# Project rules

## Stack

- Frontend: React, TypeScript, Vite, and npm in `frontend/`.
- Backend: Go in `backend/`; PostgreSQL runs through Docker Compose.
- Planned judging: a queue and worker with an isolated sandbox. See `docs/PA2-PLAN.md` for planned components.

## Commands

- Install frontend dependencies: `cd frontend` then `npm ci`.
- Run or build the frontend: `npm run dev` or `npm run build` from `frontend/`.
- Start the local database: `docker compose up -d` from the root; host port is 5433. See `AI-LOG.md`, entry 1.
- Run Go tests from `backend/`: `go test ./...` once Go packages exist.
- Document the final lint, formatting, type-check, and test commands in the README; run those same commands in CI.

## Conventions

- Keep API behavior consistent with the endpoint contract; enforce role and ownership checks on the server.
- Use `.env.example` for variable names without values; keep local configuration in ignored `.env` files. See `AI-LOG.md`, entry 2.
- Keep model verdicts, evidence, and TA-selected deductions separate; a TA must review before publishing.

## Never

- Never expose secrets, the model API key, the Docker socket, or host mounts to a student-code sandbox.
- Never make real model calls in unit tests; use a fake or stub.
- Never change tests, eval data, or labels merely to make a check pass.
- Never publish marks while required TA reviews are pending.

## Permissions

- **Allowed:** Read project files, edit code and documentation for the requested task, and run relevant local checks.
- **Ask first:** Deploy or publish externally, change paid services or spending limits, delete persistent data, or perform destructive Git operations, unless the user already authorized the action.
- **Never:** Commit secrets, bypass required checks or branch protection, or send messages to others without explicit authorization.
