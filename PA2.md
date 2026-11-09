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
- Latest green CI run before this draft: [run 42616773213](https://github.com/longtoZ/AWAD-LeeTKOB/actions/runs/42616773213) (`gates`, `secret-history`, `e2e`). Update this link to the final green run after the last evidence change.
- The workflow runs the Go format check, `go vet`, PostgreSQL-backed API tests, frontend lint/format/type/build, full-history secret scan, `npm audit`, and a Playwright shell check.
- `AI-LOG.md` records the assistant's work and setup incidents.

## Required evidence still pending

- **Public repository:** GitHub currently reports the repository as private. The connected account has push access but not repository-admin access. Before changing visibility, the owner should review the existing commit history because it contains a personal email address.
- **Protected `main`:** required checks, no-bypass policy, and direct-push rejection have not been verified. The connected account cannot change repository settings.
- **Blocked red PR and green fix:** no deliberate contract-breaking commit has yet been shown to turn a required check red and block merge. Preserve the red run and later fix commit/green run when the owner can enable protection.
- **Submission contract evidence:** add the frozen contract SHA, test checklist, and the product-owner FAQ link or the team's explicitly labeled no-answer assumption.
- **Final release evidence:** create and push the `pa2` tag after the assessed commit; link its green CI run here.
- **Hand-in:** add the final `SELF_ASSESSMENT_REPORT.md`, package the repository including `.git` with team IDs in ascending order, and submit it from a member's student account by the Classroom deadline.

