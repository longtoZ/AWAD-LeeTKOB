# PA#2 — Harness and quality gate evidence

**Status: draft; not ready to submit.** This page records evidence that exists and names the remaining acceptance items. Do not treat the pending entries as completed.

## Project and sources

- Repository: [AWAD-LeeTKOB](https://github.com/longtoZ/AWAD-LeeTKOB)
- [PA2 plan](docs/PA2-PLAN.md)
- [Submission endpoint contract](docs/contract/submissions.md)
- [Quality-gate verification record](docs/PA2-QUALITY-GATE-VERIFICATION.md)
- [PA#2 assignment and rubric](https://classroom.google.com/c/ODg2MTU1OTQ4NDc0/a/ODkzNzE1NjU1NDc5)
- [LeeTKOB project brief](https://drive.google.com/file/d/1VK58AnV3_LWNWVX6jKxj1baPMDIA8QFm/view)

## Current implementation evidence

- Quality-gate PR: [#6](https://github.com/longtoZ/AWAD-LeeTKOB/pull/6)
- Latest code-verification run: [run 42661351986](https://github.com/longtoZ/AWAD-LeeTKOB/actions/runs/42661351986) (`gates`, `secret-history`, `e2e`) on fix commit `90b68ef`. The Go log shows 5 named unit tests and 46 contract subtests.
- Deliberate contract break: code-only commit `3d030a5` temporarily bypassed the enrollment guard; [red run 42661114861](https://github.com/longtoZ/AWAD-LeeTKOB/actions/runs/42661114861) failed the unenrolled and missing-lab contract cases. Fix commit `90b68ef` restored the guard and passed in run 42661351986. No tests or contract text changed.
- The workflow runs the Go format check, `go vet`, PostgreSQL-backed API tests, frontend lint/format/type/build, full-history secret scan, `npm audit`, and a Playwright shell check.
- `AI-LOG.md` records the assistant's work and setup incidents.
- Branch ruleset screenshots show the branch protection for `main` with required checks and no bypass: [branch-ruleset-1](img/evidence/branch-ruleset-1.png), [branch-ruleset-2](img/evidence/branch-ruleset-2.png), [branch-ruleset-3](img/evidence/branch-ruleset-3.png)
- [Direct push rejection](img/evidence/main-branch-block.png): GitHub rejected a push to `main`, citing the pull request requirement and three required checks.

## Submission contract evidence

- [Submission contract](docs/contract/submissions.md): last changed in commit [`abf949c`](https://github.com/longtoZ/AWAD-LeeTKOB/commit/abf949cd718088de321cd188295878ebb7b763ae). The plan records team review and freeze as confirmed by Tran Anh Khoa; the exact freeze commit still needs his confirmation.
- [Contract test checklist](docs/contract/submissions-test-scenarios.md): review scenarios and remaining fault-injection coverage. The checklist itself is not evidence that every scenario was tested; see the [H1/H2 verification](docs/PA2-H1-H2-VERIFICATION.md) and CI run above.
- **Product-owner FAQ or no-answer assumption — pending Khoa:** add the actual FAQ entry, or document the question, whether an answer was received, and the team's specific assumption. No answer or assumption is established by the available records.

## Required evidence still pending

- **Merge block — Huy:** placeholder for the PR and evidence showing a failing required check blocks merging, followed by the passing fix. The previous red check was not required.
- **Submission contract evidence:** Khoa to confirm the exact freeze commit and supply the product-owner FAQ link or the team's specific documented no-answer assumption.
- **Final release evidence:** create and push the `pa2` tag after the assessed commit; link its green CI run here.
- **Hand-in:** add the final `SELF_ASSESSMENT_REPORT.md`, package the repository including `.git` with team IDs in ascending order, and submit it from a member's student account by the Classroom deadline.

