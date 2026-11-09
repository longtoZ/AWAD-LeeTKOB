# PA2 quality gate verification

Fresh-clone verification was run on 9 November 2026 from commit `3e4728c` in a checkout whose path does not contain `#` (Vite treats `#` in the requested `PA#2` workspace path as a URL fragment on Windows). Verification commit `90b68ef8ca1c667dc70adddd4b278bb607acaac5` passed `gates`, `secret-history`, and `e2e` in [run 42661351986](https://github.com/longtoZ/AWAD-LeeTKOB/actions/runs/42661351986). Its Go test log shows 5 named unit tests and 46 contract subtests.

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

The deliberate enrollment-guard break at commit `3d030a5` made the `gates` job fail on the unenrolled and missing-lab contract cases in [red run 42661114861](https://github.com/longtoZ/AWAD-LeeTKOB/actions/runs/42661114861); restoring the guard at `90b68ef` made the checks green. This proves the contract test catches that defect. It does not prove merge blocking: branch protection is not configured, so the red check was not a required merge check. The provider spending cap is also not established by these checks.
