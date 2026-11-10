# PA#2 — Harness and quality gate evidence

This page identifies the repository state and evidence for the seven rubric criteria. The existing `pa2` tag predates the new evidence and score updates in this worktree; after the changes are merged, create and verify the final assessed tag and package. The draft self-assessed score is 100/100; see [`SELF_ASSESSMENT_REPORT.md`](SELF_ASSESSMENT_REPORT.md) for criterion-level evidence.

## Sources and scope

- [Classroom PA#2 post](https://classroom.google.com/c/ODg2MTU1OTQ4NDc0/a/ODkzNzE1NjU1NDc5/details), [attached rubric](https://drive.google.com/file/d/1wWH6xtNdzngJdGl7t4Qxe6VBarwl8lH9/view), and [LeeTKOB brief](https://drive.google.com/file/d/1VK58AnV3_LWNWVX6jKxj1baPMDIA8QFm/view) were checked on 10 November 2026. The post and rubric specify the 11 November, 22:00 Vietnam-time deadline.
- This checkpoint covers the harness, CI, and one endpoint contract with tests. The judging sandbox and model review are later checkpoints in the brief; they are not claimed as implemented here.
- Public repository: [AWAD-LeeTKOB](https://github.com/longtoZ/AWAD-LeeTKOB). GitHub reported `visibility: public` on 10 November 2026.

## Rules, secrets, and reproducible checks

- [`AGENTS.md`](AGENTS.md) names the stack, commands, conventions, product-specific never rules, and allowed / ask-first / never permissions. Its port and example-environment rules trace to the recorded setup incidents in [`AI-LOG.md`](AI-LOG.md).
- [`.env.example`](.env.example) lists names without values; [`.gitignore`](.gitignore) excludes local `.env` files. The `secret-history` job runs Gitleaks with full Git history and `npm audit`; it passed in the [latest main run](https://github.com/longtoZ/AWAD-LeeTKOB/actions/runs/42851830119). A passing scan is evidence of the scan, not proof that every possible secret is absent.
- The [README PA2 checks](README.md#pa2-checks) give the same frontend check, Go format, `go vet`, and PostgreSQL-backed test commands that the [workflow](.github/workflows/gates.yml) runs. On `main` at `0c5de65`, [run 42851830119](https://github.com/longtoZ/AWAD-LeeTKOB/actions/runs/42851830119) passed `gates`, `secret-history`, and `e2e` on a clean runner. CI sets `TEST_DATABASE_URL` so the contract suite runs instead of skipping.
- **Spending cap:** Google AI Studio shows a ₫100,000 monthly project cap for `AWAD-LeeTKOB` ([project details](img/evidence/gemini-project-details.png), [paid-tier setup](img/evidence/gemini-paid-tier.png), [cap display](img/evidence/gemini-monthly-cap.png)). The [README](README.md#spending-cap-status) records that Google's project cap is experimental and may overrun during billing latency. The model is not integrated yet.

## Protected merge and contract break

- The active [`main` ruleset](https://github.com/longtoZ/AWAD-LeeTKOB/rules/26947349) requires a pull request and the `gates`, `secret-history`, and `e2e` checks, with no bypass actor in the ruleset. [Settings screenshots](img/evidence/branch-ruleset-1.png) and [direct-push rejection](img/evidence/main-branch-block.png) provide additional evidence.
- Since the direct-push rejection was captured, main has advanced through merged pull requests; see the [main commit history](https://github.com/longtoZ/AWAD-LeeTKOB/commits/main/) and [merged PRs #9–#13](https://github.com/longtoZ/AWAD-LeeTKOB/pulls?q=is%3Apr+is%3Amerged).
- [PR #9](https://github.com/longtoZ/AWAD-LeeTKOB/pull/9) added database row-count assertions to refusals and replays before the deliberate break. Its [green run](https://github.com/longtoZ/AWAD-LeeTKOB/actions/runs/42840683728) passed.
- [PR #10](https://github.com/longtoZ/AWAD-LeeTKOB/pull/10) introduced a **code-only** replay defect in [`22030bc`](https://github.com/longtoZ/AWAD-LeeTKOB/commit/22030bc8253aa160328b84011ab14a02eba8157a): a duplicate submission and job were written while the replay still returned the stored `201`. The required `gates` job [failed](https://github.com/longtoZ/AWAD-LeeTKOB/actions/runs/42841435023) on `creation_and_replay` and `key_retained_after_23_hours`; [GitHub showed merging blocked](img/evidence/pr-blocked-by-red-gates.png).
- The [fix commit `473475d`](https://github.com/longtoZ/AWAD-LeeTKOB/commit/473475d5e601a314ef8b4345096b97c264497778) reverted only the defective code. Tests and contract were unchanged; the [PR checks turned green](https://github.com/longtoZ/AWAD-LeeTKOB/actions/runs/42842548890) before merge.
- A fresh second-reader review of `22030bc` found the duplicate write at [`submission.go:230–235`](https://github.com/longtoZ/AWAD-LeeTKOB/blob/22030bc8253aa160328b84011ab14a02eba8157a/backend/internal/api/submission.go#L230-L235), followed by the stored-response replay at lines 239–242. It matched these to [contract §4](docs/contract/submissions.md) and found no other independent issue supported by the diff. The reviewer did not run tests or inspect CI; the red and blocked-run evidence above was checked separately.

## Commands

- [README PA2 checks](README.md#pa2-checks) match the commands in [the CI workflow](.github/workflows/gates.yml). A clean detached worktree at `8cae248`, checked out under a Windows path without `#`, passed every listed command locally; see the [11 November verification record](docs/PA2-QUALITY-GATE-VERIFICATION.md).

## Endpoint contract and tests

- [`docs/contract/submissions.md`](docs/contract/submissions.md) specifies `POST /labs/{labId}/submissions`: enrolled-student ownership, request and response fields, all problem-detail error types and statuses, check order, 24-hour key retention, replay/conflict behavior, and atomic submission/job persistence. It follows the five sections in [`docs/IA3/CONTRACT.md`](docs/IA3/CONTRACT.md), identifies the class brief as authoritative, explains cookie-jar handling for authentication, and states that no product-owner FAQ answer is relied on. Commit [`abf949c`](https://github.com/longtoZ/AWAD-LeeTKOB/commit/abf949cd718088de321cd188295878ebb7b763ae) introduced the contract; its current blob includes the source note and cookie guidance.
- [`backend/internal/api/api_test.go`](backend/internal/api/api_test.go) tests the handler over HTTP with PostgreSQL. It checks success, all listed error types, other-student enrollment refusal, replay and conflict, and database row counts after refusals and retries. The [latest main CI run](https://github.com/longtoZ/AWAD-LeeTKOB/actions/runs/42851830119) executes this suite. [Test scenarios](docs/contract/submissions-test-scenarios.md) are a review checklist, not a claim that every scenario is automated; uncertain-commit fault injection remains incomplete.
- The brief directs teams to cite a product-owner FAQ entry **when relying on one**. This contract relies on none; it explicitly labels endpoint details left open by the brief as team choices.

## Release and hand-in

- The final self-score and its evidence are in [`SELF_ASSESSMENT_REPORT.md`](SELF_ASSESSMENT_REPORT.md). Historical AI work and this preparation task are recorded in [`AI-LOG.md`](AI-LOG.md).
- **Tag CI:** After the final assessed commit is tagged and pushed as `pa2`, use the [tag-filtered workflow runs](https://github.com/longtoZ/AWAD-LeeTKOB/actions/workflows/gates.yml?query=branch%3Apa2) to find its `gates`, `secret-history`, and `e2e` results. The ZIP must contain that tagged repository, including `.git`.
- The team will submit one ZIP from a member's student Google account. ZIP preparation does not itself submit coursework.
