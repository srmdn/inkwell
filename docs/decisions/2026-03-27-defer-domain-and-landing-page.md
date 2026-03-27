# Decision Record: Defer Domain and Landing Page

- Date: 2026-03-27
- Status: Accepted
- Owners: @srmdn

## Context

The primary `.com` TLD for "inkwell" is taken. Considered buying an
alternative TLD and building a Ghost-style marketing site before development
begins.

## Decision

Defer domain purchase and landing page until after v0.1 ships.

## Rationale

- No product exists yet — a marketing site before v0.1 is premature
- GitHub serves as the project home during early development
- Buying a domain now locks in a name before the project has proven itself
- The right time to invest in marketing infrastructure is when there is
  something worth marketing

## Consequences

- GitHub repo (`github.com/srmdn/inkwell`) is the canonical project URL
  until a domain is purchased
- README must be clear enough to serve as the project landing page for now

## Post-v0.1 Actions

When v0.1 ships, revisit:

- [ ] Check domain availability: `getinkwell.com`, `inkwellcms.com`,
      `inkwell.dev`, `inkwell.app`, `tryinkwell.com`
- [ ] Build a static landing page (Astro, hosted on existing VPS or
      GitHub Pages)
- [ ] Update README to point to new domain once acquired
- [ ] Consider a separate GitHub org (`inkwell-cms`) if the project
      grows beyond a single maintainer
