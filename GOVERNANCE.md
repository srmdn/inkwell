# Governance

## Purpose

Define how this project makes decisions, resolves disputes, and maintains
accountability for AI-related risk.

## Roles

- Maintainers: merge code, enforce policy, manage releases.
- Reviewers: review code/policy changes and provide risk feedback.
- Contributors: propose and implement changes.
- Security contacts: coordinate security and AI safety incidents.

## Decision Model

- Routine changes: maintainer approval in pull request.
- High-impact policy/architecture changes: at least 2 maintainer approvals.
- Breaking governance changes: maintainers vote; simple majority wins.

## Escalation

Escalate when:
- AI feature is classified as high risk,
- a policy exception is requested,
- there is unresolved disagreement after review.

Escalations require a written decision record in `docs/decisions/`.

## Transparency

Major governance outcomes should be summarized in repository discussions,
issues, or release notes.
