---
status: idea
kind: bug
---

# Build Failure: bborbe/go-skeleton

Filed automatically by the build-fix agent for the CI episode `0db9b0dac90f4d8d7114e3698b366668f5fdb49b`.

## Summary

The default-branch build for `bborbe/go-skeleton` is failing; the build-fix diagnosis classified this as a code/test bug (verdict `file_spec`).

## Reproduction

Failing workflow(s): test

Episode SHA: `0db9b0dac90f4d8d7114e3698b366668f5fdb49b`

Log evidence:

```text
| Workflow | Job | Failed Step | Run |
|---|---|---|---|
| CI | test | Run precommit checks | [Run](https://github.com/bborbe/go-skeleton/actions/runs/32784236062) |

```
2026-08-24T22:22:52.8262219Z Temporarily overriding HOME='/home/runner/work/_temp/3bbcf548-bc41-493a-9326-40c55ded7376' before making global git config changes
2026-08-24T22:22:52.8263624Z Adding repository directory to the temporary git global config as a safe directory
2026-08-24T22:22:52.8268384Z [command]/usr/bin/git config --global --add safe.directory /home/runner/work/go-skeleton/go-skeleton
2026-08-24T22:22:52.8299459Z Removing SSH command configuration
2026-08-24T22:22:52.8307384Z [command]/usr/bin/git config --local --name-only --get-regexp core\.sshCommand
2026-08-24T22:22:52.8340901Z [command]/usr/bin/git submodule foreach --recursive sh -c "git config --local --name-only --get-regexp 'core\.sshCommand' && git config --local --unset-all 'core.sshCommand' || :"
2026-08-24T22:22:52.8539016Z Removing HTTP extra header
2026-08-24T22:22:52.8544204Z [command]/usr/bin/git config --local --name-only --get-regexp http\.https\:\/\/github\.com\/\.extraheader
2026-08-24T22:22:52.8578172Z [command]/usr/bin/git submodule foreach --recursive sh -c "git config --local --name-only --get-regexp 'http\.https\:\/\/github\.com\/\.extraheader' && git config --local --unset-all 'http.https://github.com/.extraheader' || :"
2026-08-24T22:22:52.8766741Z Removing includeIf entries pointing to credentials config files
2026-08-24T22:22:52.8772981Z [command]/usr/bin/git config --local --name-only --get-regexp ^includeIf\.gitdir:
2026-08-24T22:22:52.8795846Z includeif.gitdir:/home/runner/work/go-skeleton/go-skeleton/.git.path
2026-08-24T22:22:52.8797058Z includeif.gitdir:/home/runner/work/go-skeleton/go-skeleton/.git/worktrees/*.path
2026-08-24T22:22:52.8798143Z includeif.gitdir:/github/workspace/.git.path
2026-08-24T22:22:52.8799259Z includeif.gitdir:/github/workspace/.git/worktrees/*.path
2026-08-24T22:22:52.8806015Z [command]/usr/bin/git config --local --get-all includeif.gitdir:/home/runner/work/go-skeleton/go-skeleton/.git.path
2026-08-24T22:22:52.8830980Z /home/runner/work/_temp/git-credentials-b534e2b3-d748-44df-a647-d3a38efac081.config
2026-08-24T22:22:52.8842527Z [command]/usr/bin/git config --local --unset includeif.gitdir:/home/runner/work/go-skeleton/go-skeleton/.git.path /home/runner/work/_temp/git-credentials-b534e2b3-d748-44df-a647-d3a38efac081.config
2026-08-24T22:22:52.8879348Z [command]/usr/bin/git config --local --get-all includeif.gitdir:/home/runner/work/go-skeleton/go-skeleton/.git/worktrees/*.path
2026-08-24T22:22:52.8900026Z /home/runner/work/_temp/git-credentials-b534e2b3-d748-44df-a647-d3a38efac081.config
2026-08-24T22:22:52.8909114Z [command]/usr/bin/git config --local --unset includeif.gitdir:/home/runner/work/go-skeleton/go-skeleton/.git/worktrees/*.path /home/runner/work/_temp/git-credentials-b534e2b3-d748-44df-a647-d3a38efac081.config
2026-08-24T22:22:52.8943769Z [command]/usr/bin/git config --local --get-all includeif.gitdir:/github/workspace/.git.path
2026-08-24T22:22:52.9083151Z /github/runner_temp/git-credentials-b534e2b3-d748-44df-a647-d3a38efac081.config
2026-08-24T22:22:52.9086517Z [command]/usr/bin/git config --local --unset includeif.gitdir:/github/workspace/.git.path /github/runner_temp/git-credentials-b534e2b3-d748-44df-a647-d3a38efac081.config
2026-08-24T22:22:52.9088884Z [command]/usr/bin/git config --local --get-all includeif.gitdir:/github/workspace/.git/worktrees/*.path
2026-08-24T22:22:52.9090047Z /github/runner_temp/git-credentials-b534e2b3-d748-44df-a647-d3a38efac081.config
2026-08-24T22:22:52.9092278Z [command]/usr/bin/git config --local --unset includeif.gitdir:/github/workspace/.git/worktrees/*.path /github/runner_temp/git-credentials-b534e2b3-d748-44df-a647-d3a38efac081.config
2026-08-24T22:22:52.9094829Z [command]/usr/bin/git submodule foreach --recursive git config --local --show-origin --name-only --get-regexp remote.origin.url
2026-08-24T22:22:52.9259570Z Removing credentials config '/home/runner/work/_temp/git-credentials-b534e2b3-d748-44df-a647-d3a38efac081.config'
2026-08-24T22:22:52.9421331Z Cleaning up orphan processes
```
```

## Expected vs Actual

**Expected:** green CI on the default branch.
**Actual:** `The test step fails with a compile error in main.go:34:24 referencing undefinedSymbolForE2E1, which is a code bug in the repo itself, not a dependency issue.`

## Why this is a bug

The default-branch build is the repository's quality gate; a red build blocks merges. Diagnosis: `The test step fails with a compile error in main.go:34:24 referencing undefinedSymbolForE2E1, which is a code bug in the repo itself, not a dependency issue.`
