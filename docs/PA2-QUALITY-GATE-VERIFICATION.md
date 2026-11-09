# PA2 quality gate verification

Fresh-clone verification was run on 9 November 2026 from commit `3e4728c` in a checkout whose path does not contain `#` (Vite treats `#` in the requested `PA#2` workspace path as a URL fragment on Windows). Current PR head `3c95afc64f2196cc5b78e50bd7f4e0288b1cfd8a` passed `gates`, `secret-history`, and `e2e` in [run 42660477736](https://github.com/longtoZ/AWAD-LeeTKOB/actions/runs/42660477736). The Go test log shows 5 named unit tests and 46 contract subtests.

Commands and results:

| Command | Result |
|---|---|
| `npm ci --prefix frontend` from repo root | Passed in CI; lockfile install, 0 reported vulnerabilities. |
| `npm run check --prefix frontend` from repo root | Passed in CI: ESLint, Prettier, TypeScript check, and production build. |
| `node scripts/check-go-format.mjs` from repo root | Passed. |
| `go -C backend vet ./...` | Passed. |
| `go -C backend test ./... -count=1 -v` with isolated PostgreSQL test database | Passed in CI; log shows 5 named unit tests and 46 `TestSubmissionContract` subtests. |
| `npm audit --prefix frontend` from repo root | Passed in CI; 0 vulnerabilities. |
| `npm run e2e --prefix frontend` from repo root | Passed in CI; 1 Playwright test confirmed the LeeTKOB page heading and shell text. |

Required branch protection and the provider spending cap are not established by these checks. The green jobs prove the configured checks pass for this head; a deliberate contract-breaking change has not yet been used to prove that a required check blocks merge. Branch protection is not configured, so no job is currently enforced as a required merge check.
