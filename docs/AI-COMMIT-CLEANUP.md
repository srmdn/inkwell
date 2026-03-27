# AI Commit Attribution Cleanup

Use this runbook to remove banned AI attribution lines from commit messages.

## Warning Before Rewrite

Rewriting history changes commit hashes.

Before continuing:
- notify collaborators
- pause merges to the branch
- confirm nobody is building new work on old commits

If history was already pushed, force-push with lease only:

```sh
git push --force-with-lease
```

## Case A: Latest Commit Only

```sh
git commit --amend
scripts/check-commit-attribution.sh
```

Then push normally, or with lease if already pushed.

## Case B: Older Commits

1. Pick base commit before the first bad commit.

```sh
git rebase -i <base-commit>
```

2. Mark offending commits as `reword`.
3. Remove banned lines from each commit message.
4. Validate:

```sh
scripts/check-commit-attribution.sh
```

5. Update remote safely:

```sh
git push --force-with-lease
```

## Team Recovery After Rewrite

Share these commands with collaborators:

```sh
git fetch origin
git checkout <branch>
git reset --hard origin/<branch>
```

Only use reset after collaborators confirm local work is backed up.

## Portable vs Project-Specific

Portable:
- amend/rebase/force-with-lease sequence
- pre-rewrite collaboration warning

Project-specific:
- who must approve rewrites
- branch protection exceptions
- notification channel
