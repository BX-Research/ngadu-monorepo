# Repository instructions

## Purpose and scope

This monorepo contains Ngadu, a badminton matchmaking app for a community of friends. Read `docs/functional-spec.md` before changing product behavior. Distinguish confirmed decisions from proposed features and open choices.

## Layout

- `ngadu-fe/`: Next.js, React, and TypeScript frontend using Bun. Read its nested `AGENTS.md` before frontend edits.
- `ngadu-be/`: reserved for the backend. Language, framework, database library, and build tools are undecided. Do not scaffold a backend or add language-specific dependencies without a user decision. Rust is not a confirmed choice.
- `docs/`: functional specification and PDF export. Markdown is the source of truth; keep PDF exports consistent.

## Confirmed architecture

Next.js runs on Bun. A separate backend will verify Firebase Authentication ID tokens, enforce community permissions, query PostgreSQL, and handle SeaweedFS uploads if uploads are included. Firebase provides Google login; it does not require Firestore. SeaweedFS is selected for development. Docker Compose is planned, not implemented. Production hosting and storage remain undecided.

## Commands and environment

Run repository commands in WSL using `/bin/bash` with an explicit working directory. The visible terminal and command-execution tool may use separate environments.

From the repository root:

```bash
bun run setup:frontend
bun run dev
bun run lint
bun run typecheck
bun run build
bun run start
```

Commands delegate to `ngadu-fe`; its dependencies and Bun lockfile remain there. Each service keeps its own toolchain. A JavaScript workspace or additional build orchestrator is not required. If needed, add `/home/reyza/.bun/bin` to PATH for the current WSL Bun installation.

Try normal sandboxed execution first. If creation fails with `Failed to create unified exec process: No such file or directory (os error 2)`, use a read-only probe through the approval system with `sandbox_permissions: "require_escalated"`, `/bin/bash`, and the explicit working directory when policy permits. Respect the result. The launcher error alone does not establish that WSL, Bash, or the repository is missing.

## Working conventions

- Preserve existing work and history. Git tracking belongs at the monorepo root; do not initialize nested repositories.
- Keep rendering and browser interaction in Next.js and authoritative application rules in the future backend. Agree on the API contract before connecting services.
- Follow installed Next.js documentation as directed by the frontend `AGENTS.md`.
- Never commit credentials, real environment files, user data, dependencies, build output, or database/storage volumes. Use placeholders in configuration examples.
- Add tooling only when required for agreed work. Do not assume an auth provider, database, backend language, or product feature beyond the spec.
- Run checks appropriate to changed files and report actual results. Documentation-only changes do not need application tests. Verify PDF exports against their source.
- Before publishing, inspect staged files and diffs for unintended content or sensitive data. Public GitHub publishing requires explicit user authorization.
