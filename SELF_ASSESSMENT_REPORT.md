# Self-assessment — PA#2 (draft)

**Draft for team review; do not submit yet.** Confirm student IDs and the claimed marks as a team after the missing branch-protection, deliberate-break, public-repository, and tag evidence is resolved. The assignment ZIP filename total must match the final table.

Submitted by: Tran Anh Khoa, Pham Bao Huy, To Thanh Long — student IDs to be added in ascending order.

Proposed total to review: **80 / 100**

| Criterion | Max | Proposed claim | Evidence and remaining gap |
|---|---:|---:|---|
| Rules file | 15 | 13 | `AGENTS.md` covers stack, commands, conventions, never rules, and permission tiers. Two setup-specific rules link to incidents in `AI-LOG.md`. Team should confirm each line reflects a real incident. |
| Permissions and secrets | 10 | 8 | `.env.example`, `.gitignore`, `AGENTS.md` permissions, and green full-history secret scan in [run 42661351986](https://github.com/longtoZ/AWAD-LeeTKOB/actions/runs/42661351986). No provider/application spending cap is enforced. |
| One command per check | 10 | 9 | Root-level README commands match the `gates` workflow and passed on fresh GitHub Actions checkout in [run 42661351986](https://github.com/longtoZ/AWAD-LeeTKOB/actions/runs/42661351986). |
| CI that blocks a merge | 20 | 11 | CI runs on every PR and [run 42661114861](https://github.com/longtoZ/AWAD-LeeTKOB/actions/runs/42661114861) went red on a deliberate contract defect; [run 42661351986](https://github.com/longtoZ/AWAD-LeeTKOB/actions/runs/42661351986) passed after the fix. [Ruleset and rejected-push screenshots](PA2.md) are recorded; a failed check blocking a merge remains to be demonstrated. |
| Endpoint contract | 25 | 21 | `docs/contract/submissions.md` defines caller/ownership, request and response, RFC 9457 errors, idempotency, and persistence behavior. Add the freeze SHA and product-owner FAQ/no-answer record to `PA2.md`. |
| Contract tests that bite | 15 | 13 | `backend/internal/api/api_test.go` logs 46 contract subtests in [green run 42661351986](https://github.com/longtoZ/AWAD-LeeTKOB/actions/runs/42661351986). The code-only enrollment break at `3d030a5` made the contract cases fail in [red run 42661114861](https://github.com/longtoZ/AWAD-LeeTKOB/actions/runs/42661114861); restoring the guard in `90b68ef` made them green. The earlier red run does not demonstrate merge blocking under the ruleset. |
| AI-LOG.md | 5 | 5 | `AI-LOG.md` identifies the tool and work, records concrete changes, and includes rejected choices and verification evidence. Keep the final CI run reference current. |
| **Total** | **100** | **80** | **Provisional only. Re-score after completing the listed evidence and team review.** |

## What we did not manage yet

The project record says the connected GitHub account cannot make the private repository public. The repository's existing history contains a personal email address, so the owner should decide how to handle that history before publication. No-bypass settings, the spending cap, proof that a failing check blocked a PR, final `pa2` tag/run, and final public-repository links are not complete. This is not the final score or ZIP filename total.

- **Contract tests were not written red-first.** The submission endpoint (`2c7e6ef`) was implemented before its tests were run failing, and no output from that time was kept. Our test-first evidence comes later: the duplicate-Origin red run in `docs/PA2-H1-H2-VERIFICATION.md`, and the tests-first PR #9 before the deliberate break in PR #10. We do not present either as the original run.

