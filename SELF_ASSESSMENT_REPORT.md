# Self-assessment — PA#2 (draft)

**Draft for team review; do not submit yet.** Confirm student IDs and the claimed marks as a team after the missing branch-protection, deliberate-break, public-repository, and tag evidence is resolved. The assignment ZIP filename total must match the final table.

Submitted by: Tran Anh Khoa, Pham Bao Huy, To Thanh Long — student IDs to be added in ascending order.

Proposed total to review: **75 / 100**

| Criterion | Max | Proposed claim | Evidence and remaining gap |
|---|---:|---:|---|
| Rules file | 15 | 13 | `AGENTS.md` covers stack, commands, conventions, never rules, and permission tiers. Two setup-specific rules link to incidents in `AI-LOG.md`. Team should confirm each line reflects a real incident. |
| Permissions and secrets | 10 | 8 | `.env.example`, `.gitignore`, `AGENTS.md` permissions, and green full-history secret scan in [run 42660477736](https://github.com/longtoZ/AWAD-LeeTKOB/actions/runs/42660477736). No provider/application spending cap is enforced. |
| One command per check | 10 | 9 | Root-level README commands match the `gates` workflow and passed on fresh GitHub Actions checkout in [run 42660477736](https://github.com/longtoZ/AWAD-LeeTKOB/actions/runs/42660477736). |
| CI that blocks a merge | 20 | 11 | CI runs on every PR and has green runs; [run 42616143908](https://github.com/longtoZ/AWAD-LeeTKOB/actions/runs/42616143908) went red while the Gitleaks permission was being corrected. Branch protection and the deliberate blocked-break/fix evidence are still missing. |
| Endpoint contract | 25 | 21 | `docs/contract/submissions.md` defines caller/ownership, request and response, RFC 9457 errors, idempotency, and persistence behavior. Add the freeze SHA and product-owner FAQ/no-answer record to `PA2.md`. |
| Contract tests that bite | 15 | 8 | `backend/internal/api/api_test.go` exercises success, errors, ownership/enrollment, replay, and persistence; [run 42660477736](https://github.com/longtoZ/AWAD-LeeTKOB/actions/runs/42660477736) logs 46 named contract subtests. A deliberate contract break with red CI and a later green fix is still missing. |
| AI-LOG.md | 5 | 5 | `AI-LOG.md` identifies the tool and work, records concrete changes, and includes rejected choices and verification evidence. Keep the final CI run reference current. |
| **Total** | **100** | **75** | **Provisional only. Re-score after completing the listed evidence and team review.** |

## What we did not manage yet

The connected GitHub account cannot make the private repository public or configure branch protection. The repository's existing history contains a personal email address, so the owner should decide how to handle that history before publication. The spending cap, blocked deliberate-break PR, final `pa2` tag/run, and final public-repository links are not complete. This is not the final score or ZIP filename total.

