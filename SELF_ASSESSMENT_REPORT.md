# Self-assessment — PA#2

Submitted by (ascending student ID): Pham Bao Huy **23120131**, Tran Anh Khoa **23120135**, To Thanh Long **23120143**.

**Self-assessed total: 89 / 100.** This is an evidence-based claim, not an instructor mark. The ZIP filename is `23120131-23120135-23120143_89.zip`.

| Rubric criterion | Max | Claim | Evidence and limit |
|---|---:|---:|---|
| 1. Rules file | 15 | 13 | [`AGENTS.md`](AGENTS.md) has stack, commands, conventions and product-specific never rules. The port and names-only environment instructions trace to the setup incidents in [`AI-LOG.md`](AI-LOG.md). Historical human edits to some generated rules are not separately recorded. |
| 2. Permissions and secrets | 10 | 8 | [`AGENTS.md`](AGENTS.md#permissions) states allowed / ask-first / never; [`.env.example`](.env.example) has names without values; [`.gitignore`](.gitignore) excludes local values. The full-history `secret-history` check passed in [run 42851830119](https://github.com/longtoZ/AWAD-LeeTKOB/actions/runs/42851830119). The [README](README.md#spending-cap-status) says the proposed US$10 cap is not enforced, so we do not claim the 9–10 band. |
| 3. One command for every check | 10 | 9 | [README PA2 checks](README.md#pa2-checks) document frontend lint/format/type/build and Go format/vet/test commands matching the [workflow](.github/workflows/gates.yml). The clean-runner [main CI run](https://github.com/longtoZ/AWAD-LeeTKOB/actions/runs/42851830119) passed, including PostgreSQL-backed tests. Local Vite checks need a checkout path without `#` on Windows; see [verification record](docs/PA2-QUALITY-GATE-VERIFICATION.md). |
| 4. CI that blocks a merge | 20 | 18 | The active [main ruleset](https://github.com/longtoZ/AWAD-LeeTKOB/rules/26947349) requires a PR and `gates`, `secret-history`, `e2e`; [settings](img/evidence/branch-ruleset-1.png) and [direct-push rejection](img/evidence/main-branch-block.png) are saved. The required `gates` check [failed on PR #10](https://github.com/longtoZ/AWAD-LeeTKOB/actions/runs/42841435023) and [GitHub blocked merging](img/evidence/pr-blocked-by-red-gates.png); fix `473475d` [passed](https://github.com/longtoZ/AWAD-LeeTKOB/actions/runs/42842548890). |
| 5. One endpoint contract | 25 | 23 | [Submission contract](docs/contract/submissions.md) states enrolled-student ownership, exact request and success shape, each RFC 9457 problem type/status, idempotency, atomic state changes, and version policy. Contract blob `dd91f4dec2967c5aedd49a41d202f11bde212e13` was introduced in [`abf949c`](https://github.com/longtoZ/AWAD-LeeTKOB/commit/abf949cd718088de321cd188295878ebb7b763ae) and is unchanged. No product-owner FAQ answer is claimed. |
| 6. Contract tests that bite | 15 | 14 | [PostgreSQL HTTP tests](backend/internal/api/api_test.go) assert every error type, ownership/enrollment refusal, replay and row counts. [PR #9](https://github.com/longtoZ/AWAD-LeeTKOB/pull/9) tightened assertions before the code-only break `22030bc`; the break [failed CI](https://github.com/longtoZ/AWAD-LeeTKOB/actions/runs/42841435023) and the fix [passed](https://github.com/longtoZ/AWAD-LeeTKOB/actions/runs/42842548890), without editing tests or contract in the fix. Uncertain-commit fault injection is still absent. |
| 7. AI-LOG.md | 5 | 4 | [`AI-LOG.md`](AI-LOG.md) names the tools and tasks, explains the Gitleaks permission correction and replay-test finding, and distinguishes recorded work from unavailable human-review history. Some original H1/H2 hand edits and the first pre-implementation failing output were not recorded. |
| **Total** | **100** | **89** | One evidence line per criterion. |

## What we did not manage

- An enforced model API spending cap has not been verified; the README labels the US$10 figure as a proposal. The model is not integrated yet.
- The original failing test output from before the endpoint implementation was not kept. The later duplicate-Origin red/green run and protected replay-break PR are separate, accurately dated evidence.
- The uncertain database-commit outcome has no fault-injection test. The contract test checklist records this gap.
- The available record does not establish a product-owner FAQ answer or independent manual edits for every assistant-authored H1/H2 change. We do not claim them.
