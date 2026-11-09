# PA2 quality gate verification

Fresh-clone verification was run on 9 November 2026 from commit `3e4728c` in a checkout whose path does not contain `#` (Vite treats `#` in the requested `PA#2` workspace path as a URL fragment on Windows).

Commands and results:

| Command | Result |
|---|---|
| `npm ci` in `frontend/` | Passed; lockfile install, 0 reported vulnerabilities. |
| `npm run check` in `frontend/` | Passed: ESLint, Prettier, TypeScript check, and production build. |
| `node scripts/check-go-format.mjs` from repo root | Passed. |
| `go -C backend vet ./...` | Passed. |
| `go -C backend test ./... -count=1` with isolated PostgreSQL 18 test database | Passed, including `backend/internal/api` integration tests. |
| `npm audit` in `frontend/` | Passed; 0 vulnerabilities. |
| `npm run e2e` in `frontend/` | Passed; 1 Playwright test confirmed the LeeTKOB page heading and shell text. |

The checks were run from a fresh clone so Vite could build and serve correctly. GitHub Actions has not yet reported these jobs on the pull request. Required branch protection, the full-history secret-scan result on the GitHub refs, and the provider spending cap are not established by this local run.
