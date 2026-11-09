# PA#2 — Harness and quality gate evidence

**Status: draft; not ready to submit.** This page records evidence that exists and names the remaining acceptance items. Do not treat the pending entries as completed.

**Team status and next actions:** Start at the [PA2 team handoff](docs/PA2-HANDOFF.md) before doing remaining work. It records the current merged-PR status, PICs, blockers, prerequisites, and acceptance evidence.

## Project and sources

- Repository: [AWAD-LeeTKOB](https://github.com/longtoZ/AWAD-LeeTKOB)
- [PA2 plan](docs/PA2-PLAN.md)
- [Submission endpoint contract](docs/contract/submissions.md)
- [Quality-gate verification record](docs/PA2-QUALITY-GATE-VERIFICATION.md)
- [PA#2 assignment and rubric](https://classroom.google.com/c/ODg2MTU1OTQ4NDc0/a/ODkzNzE1NjU1NDc5)
- [LeeTKOB project brief](https://drive.google.com/file/d/1VK58AnV3_LWNWVX6jKxj1baPMDIA8QFm/view)

## Current implementation evidence

- Quality-gate PR: [#6](https://github.com/longtoZ/AWAD-LeeTKOB/pull/6), merged into `main` at [`109267345fad067751edcae7a59a9f10dfec856b`](https://github.com/longtoZ/AWAD-LeeTKOB/commit/109267345fad067751edcae7a59a9f10dfec856b).
- PR-head verification: [run 42661937462](https://github.com/longtoZ/AWAD-LeeTKOB/actions/runs/42661937462) passed `gates`, `secret-history`, and `e2e`. Post-merge `main` verification: [run 42662815719](https://github.com/longtoZ/AWAD-LeeTKOB/actions/runs/42662815719) passed. The Go log shows 5 named unit tests and 46 contract subtests.
- Deliberate contract break: code-only commit `3d030a5` temporarily bypassed the enrollment guard; [red run 42661114861](https://github.com/longtoZ/AWAD-LeeTKOB/actions/runs/42661114861) failed the unenrolled and missing-lab contract cases. Fix commit `90b68ef` restored the guard and passed in run 42661351986. No tests or contract text changed.
- The workflow runs the Go format check, `go vet`, PostgreSQL-backed API tests, frontend lint/format/type/build, full-history secret scan, `npm audit`, and a Playwright shell check.
- `AI-LOG.md` records the assistant's work and setup incidents.

## Required evidence still pending

- **Public repository:** GitHub currently reports the repository as private. The connected account has push access but not repository-admin access. Before changing visibility, the owner should review the existing commit history because it contains a personal email address.
- **Protected `main`:** required checks, no-bypass policy, and direct-push rejection have not been verified. The connected account cannot change repository settings.
- **Merge block:** the red and green runs are preserved, but the red check was not required and branch protection is unverified. The owner must enable protection and confirm the rule blocks merge during a red run.
- **Submission contract evidence:** frozen contract blob SHA is `dd91f4dec2967c5aedd49a41d202f11bde212e13`; add the test checklist and product-owner FAQ link or the team's explicitly labeled no-answer assumption. Record the second-reader review of the deliberate-break PR.
- **Spending cap:** the proposed US$10 budget is not an enforceable provider or application cap; configure one and update the README with the actual enforcement mechanism.
- **Final release evidence:** create and push the `pa2` tag after the assessed commit; link its green CI run here.
- **Hand-in:** add the final `SELF_ASSESSMENT_REPORT.md`, package the repository including `.git` with team IDs in ascending order, and submit it from a member's student account by the Classroom deadline.

