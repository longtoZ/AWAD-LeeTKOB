# AI log

Tool used: Codex, for planning and initial project setup.

## 1. Local PostgreSQL port

- **Task:** Start the local PostgreSQL service.
- **Problem:** Docker could not bind host port 5432 during setup.
- **Changed:** The assistant changed the host port to 5433, updated the connection instructions, and confirmed PostgreSQL accepted connections. `AGENTS.md` records the port so the team uses the working setup.
- **Evidence:** `compose.yaml` and `README.md`; the first repository commit preceded this change.

## 2. Example environment file

- **Task:** Prepare local configuration while following the PA2 rubric.
- **Problem:** `.env.example` contained values, while the rubric requires variable names without values.
- **Changed:** The assistant copied the local settings to ignored `.env` and cleared the values from `.env.example`. `AGENTS.md` now states this convention.
- **Evidence:** `.env.example` and `.gitignore`; see the current change history for the setup update.

## 3. Framework choice

- **Task:** Discuss a frontend and backend stack.
- **Rejected:** The user rejected the assistant's recommendation to use Next.js for the backend, preferring a backend-focused design.
- **Changed:** The plan uses React with TypeScript and a separate Go backend.
- **Evidence:** `docs/PA2-PLAN.md`, step 1; repository setup commit `aa16e18`.

## 4. Submission setup

- **Task:** Implement login, enrollment, and the submission endpoint with a simple PostgreSQL queue.
- **Changed:** Codex added SCS-backed sessions, versioned migrations, and an atomic submission/job insert. PostgreSQL contract checks run in an isolated temporary schema.
- **Changed:** PostgreSQL `text` cannot store a NUL byte, although the contract allows it in a source string. A second migration stores source as bytes; a contract check verifies preservation.
- **Remaining:** Complete the team's contract review, cover the remaining contract scenarios, and run the checks in required CI.

## 5. H1/H2 contract follow-up

- **Task:** Recheck H1/H2 against the frozen submission contract after pulling `main` at `9f25fab`.
- **Problem:** The handler accepted a request with two `Origin` header values when the first matched the configured origin. Several contract scenarios lacked exact error-type or persistence assertions. The original pre-implementation red output was not retained.
- **Changed:** Codex added a failing duplicate-Origin test, changed validation to require exactly one matching value, and expanded PostgreSQL HTTP tests for contract errors, ordering, ownership, idempotency, rollback, source boundaries, and expired sessions.
- **Evidence:** `docs/PA2-H1-H2-VERIFICATION.md` records the observed local red and green runs. This is follow-up test evidence, not the original test-first run or the required CI break PR.

## 6. PA2 checks and CI

- **Task:** Add the highest-priority reproducible checks and a required-ready quality gate, following the course harness Session 7 guidance.
- **Changed:** Codex added frontend lint/format/type/build commands, a Go format check, CI jobs for Go/PostgreSQL and frontend checks, a full-history secret scan, and a separate Playwright shell check. README and `AGENTS.md` now list the commands CI runs. The lockfile audit exposed a moderate `minimatch` advisory in the lint toolchain; a scoped override removed the reported findings.
- **Observed:** Running Vite from the requested local `PA#2` path fails because `#` in that filesystem path is parsed as a URL fragment. The same checks will be verified from a fresh clone in a path without `#`; this environment-specific failure is not hidden.
- **Remaining:** GitHub branch protection and the provider spending cap are account/provider settings. No CI green run or enforced cap is claimed until each is actually verified.
