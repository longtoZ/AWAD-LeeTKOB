# PA#2 Team Handoff — Read Before Any Remaining Work

> **This is the team’s current PA#2 handoff and work gate. All three teammates must read it before taking a remaining task.** The course rubric is the scoring authority; this file is the current status and execution checklist. The older plan remains useful for history and technical detail, but its dated status must not override this handoff.

**Status checked:** 10 November 2026 (Vietnam time)

**Submission deadline:** 11 November 2026, 22:00 Asia/Ho_Chi_Minh

**Repository:** [AWAD-LeeTKOB](https://github.com/longtoZ/AWAD-LeeTKOB)

**Course skill/harness:** [wada-hcmus-k23-skill](https://github.com/akus176/wada-hcmus-k23-skill)

**Official rubric copy:** [docs/PA#2 rubric](PA%232%20rubric%20%E2%80%94%20Harness%20and%20quality%20gate.md)

**Frozen contract:** [`docs/contract/submissions.md`](contract/submissions.md), blob SHA `dd91f4dec2967c5aedd49a41d202f11bde212e13`.

## Read these first

1. [Official PA#2 rubric](PA%232%20rubric%20%E2%80%94%20Harness%20and%20quality%20gate.md) — seven scored criteria and required hand-in.
2. [Project brief](LeeTKOB%20%E2%80%94%20class%20project%20brief.md) — product constraints.
3. [Submission endpoint contract](contract/submissions.md) and [contract test scenarios](contract/submissions-test-scenarios.md).
4. [PA2 plan](PA2-PLAN.md) — design choices and task history.
5. This file — current status, owner, prerequisites, and acceptance checks for all remaining work.

If this handoff conflicts with an old checklist status, use the status below and update the stale checklist when the task is completed. Do not mark a task done from a plan, an intention, or a passing unrelated check; attach the evidence listed below.

## Catch-up: what has happened

- PR [#6](https://github.com/longtoZ/AWAD-LeeTKOB/pull/6) is **merged** into `main` at merge commit [`109267345fad067751edcae7a59a9f10dfec856b`](https://github.com/longtoZ/AWAD-LeeTKOB/commit/109267345fad067751edcae7a59a9f10dfec856b). The branch was `PA2/quality-gate-and-checks`.
- The post-merge `main` workflow passed `gates`, `secret-history`, and `e2e`: [run 42662815719](https://github.com/longtoZ/AWAD-LeeTKOB/actions/runs/42662815719). The PR-head run also passed: [run 42661937462](https://github.com/longtoZ/AWAD-LeeTKOB/actions/runs/42661937462).
- The contract tests caught the deliberate enrollment-guard defect in commit [`3d030a5`](https://github.com/longtoZ/AWAD-LeeTKOB/commit/3d030a5), with [red run 42661114861](https://github.com/longtoZ/AWAD-LeeTKOB/actions/runs/42661114861). Commit [`90b68ef`](https://github.com/longtoZ/AWAD-LeeTKOB/commit/90b68ef) restored the guard and passed. This proves the tests catch that defect; it **does not** prove a required check blocked a merge.
- K1 and K2 are marked done in the plan. H2 and H4 are verified done. The endpoint, CI workflow, test suite, README commands, and evidence draft exist.
- The `SELF_ASSESSMENT_REPORT.md` draft proposes **80/100**. It is not an official grade or a final team score; the team must re-score from evidence after the remaining work.

## Current blockers and evidence limits

The repo was fetched and checked on 10 November. `main` includes PR #6. The repository API reports it is **private**; the connected account has push access but **not admin access**. Reading `main` protection returned 404 to that non-admin account, so branch protection is **not verified**. Ask a repository admin to inspect and configure it; do not report the 404 alone as proof either way.

The repository history contains a personal email address. Before making the repository public, the owner/team must decide whether to rewrite history or knowingly publish that history. Do not rewrite shared history or change visibility until the owner has made that decision.

The README’s proposed US$10 model budget is not an enforced provider or application spending cap. No provider cap has been verified. The rubric expects an actual cap and a truthful README statement.

The original failing output from before the endpoint implementation was not preserved. `docs/PA2-H1-H2-VERIFICATION.md` says this plainly. Do not recreate a later failure and label it as the historical original. The separate H3 red/green evidence above is valid for showing that the contract tests bite.

## Remaining tasks — execute in this order

### 1. L4 — Protect `main` and capture a real merge block

**PIC:** To Thanh Long, working with a repository admin

**Prerequisite:** Admin access to repository settings; PR #6 is already merged, so its old red run cannot demonstrate that protection blocked merging.

**Requirement:** Protect `main`; require pull requests and the relevant PA2 checks (`gates`, `secret-history`, and `e2e`); disallow bypass; verify that a direct push is rejected. Do not treat the historical red run from PR #6 as proof of merge blocking.

**Outcome/evidence:** Settings link or screenshot proving the rule and required checks, plus evidence of a rejected direct push. The deliberate red PR and its fix are task 4.

### 2. L2 — Enforce and document the model spending cap

**PIC:** To Thanh Long

**Prerequisite:** Agree the cap amount with the account/provider owner; the plan’s US$10 is only a proposal.

**Requirement:** Configure an enforceable provider-side cap for the model key, or a tested application-level cap if the provider cannot enforce one. Keep the key out of Git and logs. Measure or clearly label the cost-per-lab estimate and state exactly what enforces the cap in `README.md`.

**Outcome/evidence:** Provider setting or implementation evidence, a check showing the cap behavior, and matching README text. Do not claim that a budget is enforced based only on an estimate.

### 3. Resolve public visibility and history exposure

**PIC:** To Thanh Long, with the repository owner/admin

**Prerequisite:** Owner decision on the personal email in commit history and any history rewrite. Coordinate with every teammate before rewriting shared history.

**Requirement:** Make the repository public as required by the assignment after the history decision. Recheck that the required files and PA2 evidence are visible without authentication.

**Outcome/evidence:** Public repository URL and successful anonymous access; document the history decision. Do not expose credentials or a live secret.

### 4. H3/L5 — Complete the protected red/green and second-reader evidence

**PIC:** Pham Bao Huy for the controlled defect/fix; To Thanh Long for the second-reader record

**Prerequisites:** Task 1 is complete and `main` is protected; the current code is green.

**Requirement:** On a separate PR, make only the enrollment-guard code change needed to violate the frozen contract. Confirm the required CI check fails and GitHub prevents merge. Restore the guard without modifying tests or contract, then confirm the required checks pass. In a fresh assistant session, ask for a review of the break diff against the contract; record the exact finding, cited line, and any incorrect or unsupported claims.

**Outcome/evidence:** Red PR/check evidence, fix commit and green run, plus a concise second-reader result in `PA2.md`. The prior red run 42661114861 remains useful test evidence but is not proof of protected merge blocking.

### 5. K3 — Finish root `PA2.md` and evidence links

**PIC:** Tran Anh Khoa

**Prerequisites:** Tasks 1–4; contract owner review; provider-cap and visibility decisions.

**Requirement:** Link the repository, protected-branch settings, blocked red PR, green fix, deliberate contract break, frozen contract SHA, product-owner FAQ or an explicitly labeled no-answer assumption, and the final `pa2` tag CI run. Keep `AI-LOG.md` limited to actual work and evidence; do not invent missing historical activity.

**Outcome/evidence:** Root `PA2.md` is complete and every link opens to evidence supporting its statement; `AI-LOG.md` remains truthful and current.

### 6. K4 — Finalize the self-assessment and team identifiers

**PIC:** Tran Anh Khoa

**Prerequisites:** K3 complete; all three members confirm their student IDs and score claims.

**Requirement:** Add student IDs in ascending order; provide one row per official rubric criterion, a mark and evidence for every row, a short “What we did not manage” section, and a total that matches the eventual ZIP filename. Remove or revise unsupported claims.

**Outcome/evidence:** Team-approved final `SELF_ASSESSMENT_REPORT.md`; its total is an evidence-based team claim, not a guarantee of the instructor’s grade.

### 7. K5 — Tag, package, and submit

**PIC:** Tran Anh Khoa

**Prerequisites:** K4 approved; final assessed commit is on public `main`; all checks pass; an active student Google account is available.

**Requirement:** Create and push the `pa2` tag on the assessed `main` commit; verify the tag’s CI run is green. Create a ZIP of the repository **including `.git`**, named with the three student IDs in ascending order and the final self-assessed total, following the exact rubric filename pattern. Submit one ZIP from a member’s `@student.hcmus.edu.vn` account before the deadline.

**Outcome/evidence:** Tag URL and green run linked in `PA2.md`; inspect ZIP contents and filename; save/confirm the Classroom submission receipt before 11 November 2026, 22:00 Vietnam time.

## Done / not done snapshot

| Work item | Status | Evidence or next acceptance check |
|---|---|---|
| K1, K2 contract review | Done | Checked in `docs/PA2-PLAN.md`; contract and decisions are in `docs/contract/submissions.md`. |
| H2 endpoint implementation | Done | PostgreSQL-backed contract tests; see `docs/PA2-H1-H2-VERIFICATION.md` and the green CI run. |
| L1/L3 commands and CI jobs | Implemented and green | Post-merge [run 42662815719](https://github.com/longtoZ/AWAD-LeeTKOB/actions/runs/42662815719). |
| H3 test catches an enrollment defect | Demonstrated | Red [run 42661114861](https://github.com/longtoZ/AWAD-LeeTKOB/actions/runs/42661114861), restored by `90b68ef`. Protected merge block still open. |
| H1 original pre-implementation red output | Unavailable | Original run was not preserved; keep this limitation explicit. Do not fabricate it. |
| L2 actual spending cap | Open | Configure and verify; update README. |
| L4 protected `main` / direct-push rejection | Open / admin evidence required | Settings and rejected-push evidence. |
| Protected red PR and second-reader result | Open | Complete task 4 after L4. |
| Public repository | Open | Private as of the latest check; decide on history/email first. |
| K3/K4 final evidence and self-assessment | Draft | Finish after evidence tasks and collect student IDs. |
| K5 `pa2` tag, ZIP, and Classroom receipt | Open | Complete last, before the deadline. |

## Team rules for finishing

- Read this handoff before beginning a remaining task. Each PIC posts the resulting commit, setting screenshot, PR, CI run, or receipt as evidence.
- A checkbox is complete only when its acceptance evidence exists. Update this handoff and the relevant plan/report in the same change.
- Keep one accountable PIC per task. Reviewers verify evidence; they do not mark another member’s task done without checking it.
- If an acceptance item cannot be completed, state the blocker and its likely rubric impact in `PA2.md` and the self-assessment. Never fabricate, backdate, or relabel evidence.
- The rubric—not this checklist, the proposed score, or the assistant’s assessment—determines the final mark.
