# AI-LOG.md

> Updated 2026-11-10 to match the course harness template. Earlier notes were topic summaries, not contemporaneous task entries. The entries below use repository history and the existing log; any human review or work not recorded there is explicitly marked as unrecorded. Confirm those fields with the task PICs before submission rather than filling them from memory without evidence.

## 2026-11-07 — Initial stack and local repository setup
Tool: Codex (model not recorded)
Asked for: Set up the project foundation with a frontend, Go backend, and local PostgreSQL.
Kept: React with TypeScript, a separate Go backend, and PostgreSQL setup; see `docs/PA2-PLAN.md`, `README.md`, and commit `aa16e18`.
Changed: The team chose a backend-focused Go service instead of the initial Next.js backend direction.
Rejected: Next.js as the backend, because the team preferred a separate backend-focused design.
By hand: The team made the stack decision. The previous log attributes the port change to 5433 and the connection check to the assistant; it does not record a separate human command or review for those steps.

## 2026-11-08 — Example environment file
Tool: Codex (model not recorded)
Asked for: Prepare local configuration that follows the PA2 rule that `.env.example` contains variable names, not values.
Kept: Names-only `.env.example` and ignored local `.env`; see `.env.example`, `.gitignore`, and commit `354ec2a`.
Changed: No human-authored edit to the generated setup is recorded; the assistant's file edits are summarized under Kept.
Rejected: None recorded.
By hand: The task was framed around the rubric requirement. Separate human review details are not recorded in the available history.

## 2026-11-08 — Login, enrollment, and submission endpoint
Tool: Codex (model not recorded)
Asked for: Implement the agreed submission flow with server-side sessions, enrollment checks, and durable submission/job recording.
Kept: SCS-backed sessions, PostgreSQL migrations, and atomic submission/job persistence; see commit `2c7e6ef`, `backend/`, and `docs/contract/submissions.md`.
Changed: No human-authored implementation change is recorded. The assistant's NUL-byte storage change is described under Kept; see migration `00002_source_bytes.sql` and its contract test.
Rejected: None recorded.
By hand: The team selected the endpoint and contract. The prior log does not identify which team member reviewed each implementation detail; confirm with the PIC before claiming a manual review.

## 2026-11-09 — H1/H2 contract follow-up
Tool: Codex (model not recorded)
Asked for: Recheck the endpoint against the frozen contract after pulling `main` at `9f25fab`.
Kept: A duplicate-Origin regression test, stricter Origin validation, and expanded PostgreSQL contract assertions; see `docs/PA2-H1-H2-VERIFICATION.md` and commit `2b18a48`. That verification file states the original pre-implementation red output was not preserved.
Changed: No human-authored code change is recorded; the assistant's follow-up changes are summarized under Kept.
Rejected: None recorded for this follow-up.
By hand: The team requested comparison with the frozen contract. Further manual review by the PIC is not recorded here.

## 2026-11-10 — Historical failure evidence correction
Tool: Codex (model not recorded)
Asked for: Supply the original expected failing output from before endpoint implementation, if it exists; do not describe a later red run as the original.
Kept: The truthful note in `docs/PA2-H1-H2-VERIFICATION.md` that the original run was not preserved.
Changed: No file or code change was requested to recreate the missing run.
Rejected: Reconstructing or relabeling a later red run as historical evidence.
By hand: The team member explicitly clarified the evidence requirement and corrected the record boundary; no original run was available.

## 2026-11-09 — PA2 checks and CI
Tool: Codex (model not recorded)
Asked for: Add reproducible checks and a fail-closed CI gate aligned with Session 7 and the README commands.
Kept: Go format/vet/tests, frontend lint/format/type/build, full-history secret scan, dependency audit, and Playwright smoke check; see `.github/workflows/`, `README.md`, and `docs/PA2-QUALITY-GATE-VERIFICATION.md`.
Changed: No human-authored workflow edit is recorded; the assistant's Gitleaks permission fix after the HTTP 403 and scoped `minimatch` override are described under Kept.
Rejected: None recorded.
By hand: The team asked to align the work with the course harness and Session 7. The available log does not record a teammate’s line-by-line review of every workflow command; the PIC should add that review only if they performed it.

## 2026-11-10 — Deliberate contract-break evidence
Tool: Codex (model not recorded)
Asked for: Demonstrate that the unchanged contract tests catch an enrollment-guard defect and preserve truthful red/green evidence.
Kept: Code-only defect commit `3d030a5`, red run [42661114861](https://github.com/longtoZ/AWAD-LeeTKOB/actions/runs/42661114861), guard restoration in `90b68ef`, and green run [42661351986](https://github.com/longtoZ/AWAD-LeeTKOB/actions/runs/42661351986).
Changed: No human-authored code change is recorded; the assistant temporarily removed and restored the guard, without changing the contract or tests.
Rejected: Treating this red run as the original pre-implementation failure; that output was not preserved.
By hand: The team required the historical test-first evidence to be distinguished from this later demonstration. A protected merge block remains to be evidenced by the team.

## 2026-11-10 — Team handoff and status email
Tool: Codex (model not recorded); Gmail connector for the requested email.
Asked for: Create a checkable handoff with the remaining PA2 tasks, PICs, prerequisites, requirements, and outcomes; compile the priorities into an email to the named teammates.
Kept: `docs/PA2-HANDOFF.md`, links from `AGENTS.md`, `README.md`, `PA2.md`, and `docs/PA2-PLAN.md`, and the sent priority email. Handoff commit: `e143b52`; its checks passed in run [42665230607](https://github.com/longtoZ/AWAD-LeeTKOB/actions/runs/42665230607).
Changed: The user clarified that PR #7 did not need to be reopened or replaced; the priority list was corrected accordingly.
Rejected: Reopening or replacing PR #7.
By hand: The team member specified the assignees and corrected the PR #7 action. The assistant checked the repository/CI state and sent the email only after the explicit request. Team review of the handoff is not recorded.

## 2026-11-10 — AI log format alignment
Tool: Codex (model not recorded)
Asked for: Check whether `AI-LOG.md` follows the course harness and correct its format without inventing past work.
Kept: The harness fields `Tool`, `Asked for`, `Kept`, `Changed`, `Rejected`, and `By hand`; see the [course template](https://github.com/akus176/wada-hcmus-k23-skill/blob/main/harness/templates/AI-LOG.template.md).
Changed: The assistant reorganized the existing topic summaries into dated task entries and marked human review details as unrecorded where the evidence does not establish them.
Rejected: Filling missing student decisions, reviews, or manual test results from assumptions.
By hand: A team member flagged the mismatch and requested harness alignment. The team PICs still need to verify the factual entries and add any human review or work they personally performed.

## Current open evidence
- Repository visibility is still private, and branch protection has not been verified by an admin.
- The proposed US$10 model budget is not an enforced spending cap.
- The protected red-PR merge block, second-reader review, final evidence, `pa2` tag, ZIP, and Classroom receipt remain to be completed.
- The original pre-implementation H1 failing output was not preserved. Do not recreate it and label it as original.
