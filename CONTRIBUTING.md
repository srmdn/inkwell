# Contributing

Thanks for contributing.

## Quick Rules

- Be respectful and constructive.
- Keep pull requests focused and reviewable.
- Humans own final code, tests, docs, and commits.
- Disclose meaningful AI assistance in PR/MR descriptions.

## AI Contribution Policy

This project allows AI-assisted work, but requires human accountability.

Required in each PR/MR:
- AI tools/models used (if any)
- files or sections materially influenced by AI
- human validation performed (tests, review, security/license checks)

Prohibited in commit history:
- AI branding lines and AI co-author trailers
- automated attribution text injected by AI tools

Install the pre-commit hook to catch issues locally before push:

```sh
cp scripts/pre-commit .git/hooks/pre-commit
chmod +x .git/hooks/pre-commit
```

Or run the check manually:

```sh
scripts/check-commit-attribution.sh
```

## Recovery if Check Fails

- Last commit only:

```sh
git commit --amend
```

- Earlier commits:

```sh
git rebase -i <base-commit>
```

- After rewriting pushed history:

```sh
git push --force-with-lease
```

Detailed runbook: `docs/AI-COMMIT-CLEANUP.md`

## Portable vs Project-Specific

Portable:
- AI disclosure expectations
- attribution check policy
- cleanup commands

Project-specific:
- reviewer approval thresholds
- branch naming rules
- test matrix and release requirements
