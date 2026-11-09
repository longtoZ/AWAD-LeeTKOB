# PA2 quality gate verification

Fresh-clone verification was run on 9 November 2026 from commit `3e4728c` in a checkout whose path does not contain `#` (Vite treats `#` in the requested `PA#2` workspace path as a URL fragment on Windows). PR head `12caa1a8d47f9bb2b7dab949303eee60d48f2cde` passed `gates`, `secret-history`, and `e2e` in [run 42616773213](https://github.com/longtoZ/AWAD-LeeTKOB/actions/runs/42616773213). A later run for any updated head must be recorded here before hand-in.

Commands and results:

| Command | Result |
|---|---|
| `npm ci` in `frontend/` | Passed; lockfile install, 0 reported vulnerabilities. |
| `npm run check` in `frontend/` | Passed: ESLint, Prettier, TypeScript check, and production build. |
| `node scripts/check-go-format.mjs` from repo root | Passed. |
| `go -C backend vet ./...` | Passed. |
| `go -C backend test ./... -count=1 -v` with isolated PostgreSQL 18 test database | Pending verification in the updated PR run. The earlier local PostgreSQL suite passed without `-v`; the updated CI log must show named unit and `TestSubmissionContract` subtests. |
| `npm audit` in `frontend/` | Passed; 0 vulnerabilities. |
| `npm run e2e` in `frontend/` | Passed; 1 Playwright test confirmed the LeeTKOB page heading and shell text. |

Required branch protection and the provider spending cap are not established by these checks. The green jobs prove the configured checks pass for this head; a deliberate contract-breaking change has not yet been used to prove that a required check blocks merge. Branch protection is not configured, so no job is currently enforced as a required merge check.
