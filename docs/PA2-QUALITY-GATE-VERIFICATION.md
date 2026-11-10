# PA2 quality gate verification

Fresh-clone verification was run on 9 November 2026 from commit `3e4728c` in a checkout whose path does not contain `#` (Vite treats `#` in the requested `PA#2` workspace path as a URL fragment on Windows). Verification commit `90b68ef8ca1c667dc70adddd4b278bb607acaac5` passed `gates`, `secret-history`, and `e2e` in [run 42661351986](https://github.com/longtoZ/AWAD-LeeTKOB/actions/runs/42661351986). Its Go test log shows 5 named unit tests and 46 contract subtests.

## Local command verification — 11 November 2026

All README check commands passed again from a clean detached worktree at `8cae248` on Windows. Its path was `C:\Users\ADMIN\AppData\Local\Temp\pa2-command-check-20261111`, without `#`; this confirms the earlier issue was the original workspace path, not the project commands. Local versions were Node.js 24.11.0 / npm 11.21.1 and Go 1.27.3; CI remains pinned to Node.js 22.13.1 and Go 1.26.x. PostgreSQL 17 ran in a temporary Docker container on host port 55433 for the integration suite.

| Command | Result |
|---|---|
| `npm ci --prefix frontend` | Passed; lockfile install reported 0 vulnerabilities. |
| `npm run check --prefix frontend` | Passed: lint, Prettier, TypeScript, and Vite production build. |
| `node scripts/check-go-format.mjs` | Passed. |
| `go -C backend vet ./...` | Passed. |
| `go -C backend test ./... -count=1 -v` with `TEST_DATABASE_URL` | Passed; 5 named unit tests and all PostgreSQL contract tests passed. |
| `npm audit --prefix frontend` | Passed; 0 vulnerabilities. |
| `npx --prefix frontend playwright install chromium` then `npm run e2e --prefix frontend` | Passed; 1 browser smoke test. |

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

The deliberate enrollment-guard break at commit `3d030a5` made the `gates` job fail on the unenrolled and missing-lab contract cases in [red run 42661114861](https://github.com/longtoZ/AWAD-LeeTKOB/actions/runs/42661114861); restoring the guard at `90b68ef` made the checks green. This earlier run proves the test catches that defect; it predates the protected-merge demonstration.

On 10 November, tests-first [PR #9](https://github.com/longtoZ/AWAD-LeeTKOB/pull/9) added row-count assertions. Code-only break `22030bc` on [PR #10](https://github.com/longtoZ/AWAD-LeeTKOB/pull/10) made the required `gates` job [fail](https://github.com/longtoZ/AWAD-LeeTKOB/actions/runs/42841435023), and GitHub [blocked merging](../img/evidence/pr-blocked-by-red-gates.png). Fix `473475d` restored the code without changing tests or contract, and the [PR checks passed](https://github.com/longtoZ/AWAD-LeeTKOB/actions/runs/42842548890). The active [main ruleset](https://github.com/longtoZ/AWAD-LeeTKOB/rules/26947349) requires `gates`, `secret-history`, and `e2e`. The latest [main run](https://github.com/longtoZ/AWAD-LeeTKOB/actions/runs/42851830119) passed all three. That run does not establish a model spending cap; the provider project-cap screenshots were supplied later and are documented in [README](../README.md#spending-cap-status).
