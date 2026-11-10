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

## 2026-11-08 — 4. Submission setup: login, enrollment and `POST /labs/{labId}/submissions` (H2, commit `2c7e6ef`, PR #3)
Tool: Codex
Asked for: implement login, enrollment and the submission endpoint with a simple PostgreSQL queue, following `docs/contract/submissions.md`.
Kept: SCS sessions stored in PostgreSQL, versioned migrations, and the submission and its pending job inserted in one transaction; PostgreSQL contract checks run in an isolated temporary schema. The second migration (`00002_source_bytes.sql`) stores source as bytes, because PostgreSQL `text` cannot hold a NUL byte while the contract allows one in `source`; a contract check verifies the bytes are preserved.
Changed: No separate human rewrite is recorded for this task. The assistant's NUL-byte storage correction is included under Kept.
Rejected: None recorded.
By hand: The team selected the endpoint and contract; the available record does not attribute further hand edits to a member.

## 2026-11-09 — 5. H1/H2 contract follow-up (commit `2b18a48`, PR #5)
Tool: Codex
Asked for: recheck H1/H2 against the frozen submission contract at `9f25fab`, and add the missing contract tests.
Kept: the duplicate-Origin finding — the handler accepted two `Origin` values when the first matched, although contract §4 says "exactly the configured application origin". The new `TestOriginMustMatchConfiguredValue` failed first (`api_test.go:67: multiple Origin values accepted`), then `validOrigin` was changed to require exactly one matching value. Also kept the expanded PostgreSQL HTTP tests: every error type's status, type and title, check order, ownership/enrollment, idempotency, rollback, source boundaries and expired sessions. Red and green runs: `docs/PA2-H1-H2-VERIFICATION.md`. The tests and `server.go` changed in the same commit; the tests got stricter (no assertion weakened) and the contract text did not change, because the code was the part that disagreed with it. The original red run from before the endpoint was implemented was never kept; this follow-up does not replace it.
Changed: No separate human rewrite is recorded. The assistant made the Origin fix and added assertions against the existing contract.
Rejected: None recorded.
By hand: The team requested the contract comparison; the available record does not attribute further hand edits to a member.

## 6. PA2 checks and CI

- **Task:** Add the highest-priority reproducible checks and a required-ready quality gate, following the course harness Session 7 guidance.
- **Changed:** Codex added frontend lint/format/type/build commands, a Go format check, CI jobs for Go/PostgreSQL and frontend checks, a full-history secret scan, and a separate Playwright shell check. README and `AGENTS.md` now list the commands CI runs. The lockfile audit exposed a moderate `minimatch` advisory in the lint toolchain; a scoped override removed the reported findings.
- **Changed after CI review:** The first Gitleaks PR run returned HTTP 403 because the security job lacked pull-request read permission. Codex added only `pull-requests: read` to that job; the subsequent CI run passed. This was a workflow-permission failure, not a deliberate contract-break demonstration.
- **Observed:** Running Vite from the requested local `PA#2` path fails because `#` in that filesystem path is parsed as a URL fragment. Fresh-clone checks passed from a path without `#`; the Windows path limitation is recorded in `docs/PA2-QUALITY-GATE-VERIFICATION.md`.
- **Evidence:** Verification commit `90b68ef` passed `gates`, `secret-history`, and `e2e` in [run 42661351986](https://github.com/longtoZ/AWAD-LeeTKOB/actions/runs/42661351986); the Go log shows 5 unit tests and 46 contract subtests. Codex then temporarily bypassed the enrollment check in code-only commit `3d030a5`; `gates` failed on the unenrolled/missing-lab cases in [run 42661114861](https://github.com/longtoZ/AWAD-LeeTKOB/actions/runs/42661114861). Restoring the guard in `90b68ef` produced the green run. This demonstrates that the tests catch this defect; it does not prove that GitHub blocked a merge.
- **Status at that time:** Branch protection and public visibility were still unresolved. They were established later on 10 November; see `PA2.md` and PR #10. The proposed US$10 provider budget was not an enforced spending cap.

## 2026-11-10 — H3: protected merge block and deliberate contract break (Pham Bao Huy)
Tool: Claude Code (Claude Opus 5.5).
Asked for: check my PA#2 tasks (H1–H4) against the rubric and finish H3: a red required check that blocks a merge, then the commit that turns it green.
Kept: the finding that every refusal test checked only the problem response, and that the replay test counted jobs for the first submission ID only, so a handler that wrote a second submission and job on a retry and still replayed the stored `201` stayed green. Kept the before/after row-count checks (PR #9, tests only, +30 lines, nothing removed) and the code-only replay break `22030bc` (PR #10). Before opening PR #10 the assistant checked the claim both ways in throwaway postgres:17 / Go 1.26 containers: the old test file passed against the break, the new one failed on `creation_and_replay` and `key_retained_after_23_hours`.
Changed: nothing in the code or tests. I read the diff of PR #9 (30 added lines, no assertion removed) and of `22030bc` before pushing, and checked that PR #10 touched only `submission.go`.
Rejected: repeating the enrollment-guard break from the team handoff (PR #7). That break was already shown red in run 42661114861, so a second one would prove nothing new; the replay break showed a hole the old tests could not catch.
By hand: opened and merged PR #9 and PR #10; checked that the required `gates` check blocked the merge and took the screenshot; chose which break to make. No additional hand work is recorded.

## 2026-11-10 — PA#2 evidence and hand-in preparation
Tool: Codex (this task).
Asked for: resolve the remaining PA#2 gaps against the live rubric and prepare a submission ZIP.
Kept: the existing contract, tests, and red/green PR evidence. A fresh second reader identified the duplicate write at `submission.go:230–235` in the deliberate-break commit against contract §4; it did not run CI. No product implementation was generated in this task.
Changed: corrected stale status in `PA2.md` and the self-assessment, replaced unfilled historical log prompts with explicit "not recorded" statements, and aligned the report with confirmed team IDs and live evidence.
Rejected: treating the draft 80/100, an unverified spending cap, or an uncited product-owner answer as established facts.
By hand: Tran Anh Khoa supplied the team's three student IDs. Further team review of the final score and archive is not recorded here.

## 2026-11-11 — Gemini project spend-cap evidence
Tool: Codex (this task), using screenshots supplied by the user.
Asked for: update the PA#2 permissions-and-secrets evidence based on Google AI Studio screenshots.
Changed: saved screenshots showing the `AWAD-LeeTKOB` project identity, linked billing/paid-tier setup, and a ₫100,000 monthly project spend cap. Updated the README, PA2 evidence page, and self-assessment. The score claims the rubric's top band while disclosing Google's experimental cap and approximate 10-minute billing delay, which can permit overages.
By hand: the screenshots were supplied by the user. No API key value or identifier was added to the repository; model integration is not implemented.

## 2026-11-11 — Remaining PA#2 rubric evidence
Tool: Codex (this task).
Asked for: close the remaining self-assessment gaps for commands, protected CI, and the endpoint contract.
Kept: the README commands and CI workflow, existing red/green PR #10 evidence, active main ruleset, and the frozen submission behavior. No application code, tests, test data, or labels were changed.
Changed: recorded a clean local run of every documented check from a path without `#`; documented the post-protection main history and linked PRs #9–#13; added contract source and cookie-jar notes distinguishing the class brief from team choices and stating that no product-owner FAQ answer is relied on. Updated the draft claims for criteria 3–5 to 10/10, 20/20, and 25/25 respectively. The draft total is 98/100 pending merge and final tag/package.
Verified: frontend install/check, Go format/vet/PostgreSQL tests, dependency audit, and Playwright browser smoke test passed locally. The main branch history and active ruleset were read from GitHub; no external changes were made.

## 2026-11-11 — Self-assessment rubric calibration
Tool: Codex (this task), comparing the report against the supplied PA#2 rubric.
Changed: revised the draft total from 98/100 to 100/100. Criterion 6's full-credit band requires contract-error/ownership/idempotency coverage in required CI and a deliberate red/green break; uncertain-commit fault injection is not listed. Criterion 7's five-point description requires named tools/tasks and specific Changed/Rejected entries, which the log contains. Kept the absent fault-injection test and missing historical H1/H2 evidence disclosed as limits, without using them to deduct points the rubric does not specify.
Rejected: treating the rubric's ±10 self-assessment tolerance as permission to inflate a score; the updated score is based on the written full-credit descriptions. If the instructor awards 98, the 2-point difference is within the rubric's no-adjustment range.

## 2026-11-11 — Clarify browser test setup
Tool: Codex, following read-only review of the PA#2 evidence changes.
Changed: separated the Playwright test command from browser installation in README and clarified that Windows/macOS and Linux CI need different browser setup flags, while both run the same `npm run e2e --prefix frontend` check.
Kept: the existing local test result and the CI workflow command; no application or test files changed.
