# Repository instructions

## Purpose and scope

This monorepo contains Ngadu, a badminton matchmaking app for a community of friends. Read `docs/functional-spec.md` before changing product behavior. Distinguish confirmed decisions from proposed features and open choices.

## Layout

- `frontend/`: Next.js, React, and TypeScript frontend using Bun. Read its nested `AGENTS.md` before frontend edits.
- `backend/`: reserved for the backend. Language, framework, database library, and build tools are undecided. Do not scaffold a backend or add language-specific dependencies without a user decision. Rust is not a confirmed choice.
- `docs/`: functional specification and PDF export. Markdown is the source of truth; keep PDF exports consistent.

## Confirmed architecture

Next.js runs on Bun. A separate backend will verify Firebase Authentication ID tokens, enforce community permissions, query PostgreSQL, and handle SeaweedFS uploads if uploads are included. Firebase provides Google login; it does not require Firestore. SeaweedFS is selected for development. Docker Compose is planned, not implemented. Production hosting and storage remain undecided.

## Development commands

From the repository root:

```bash
bun run setup:frontend
bun run setup:hooks
bun run dev
bun run lint
bun run typecheck
bun run build
bun run start
```

Commands delegate to `frontend`; its dependencies and Bun lockfile remain there. Each service keeps its own toolchain. A JavaScript workspace or additional build orchestrator is not required.

## Validation

- Run `bun run setup:hooks` once per clone to enable the shared Git hooks. Install the Bun version declared in `package.json` and use `bun run setup:frontend` for a frozen dependency install.
- For frontend code, dependencies, or build configuration changes, run `bun run lint`, `bun run typecheck`, and `bun run build` before reporting completion, committing, or opening a PR. If behavior changes, exercise the affected flow and add focused automated tests when a test runner is available; these commands do not replace behavior tests.
- For Git hook changes, run `sh -n .githooks/pre-commit`, verify the hook skips documentation-only commits and blocks failing checks in a temporary repository, and run the frontend checks above.
- For documentation-only changes, run `git diff --check`, verify links and commands, and check any PDF export against its Markdown source. Application lint, typecheck, and build are unnecessary unless executable code or configuration also changed.
- Before committing, inspect `git diff --cached` for scope and sensitive content and run `git diff --cached --check`. Report which checks passed, failed, or could not run; do not claim unrun checks passed.
- The pre-commit hook selects checks from staged paths but runs lint and typecheck against the working tree. It does not build or test behavior. Partially staged changes require reviewing the staged result separately. Do not bypass failing hooks without the user's explicit instruction.
- Backend tooling and tests remain undecided; add checks when its stack is selected.

## Working conventions

- Preserve existing work and history. Git tracking belongs at the monorepo root; do not initialize nested repositories.
- Keep rendering and browser interaction in Next.js and authoritative application rules in the future backend. Agree on the API contract before connecting services.
- Follow installed Next.js documentation as directed by the frontend `AGENTS.md`.
- Never commit credentials, real environment files, user data, dependencies, build output, or database/storage volumes. Use placeholders in configuration examples.
- Add tooling only when required for agreed work. Do not assume an auth provider, database, backend language, or product feature beyond the spec.
- Run checks appropriate to changed files and report actual results. Documentation-only changes do not need application tests. Verify PDF exports against their source.
- Before publishing, inspect staged files and diffs for unintended content or sensitive data. Public GitHub publishing requires explicit user authorization.
