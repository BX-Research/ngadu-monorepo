# Repository instructions

## Purpose and scope

This monorepo contains Ngadu, a badminton matchmaking app for a community of friends. Read `docs/functional-spec.md` before changing product behavior. Distinguish confirmed decisions from proposed features and open choices.

## Layout

- `frontend/`: Next.js, React, and TypeScript frontend using Bun. Read its nested `AGENTS.md` before frontend edits.
- `backend/`: Go backend. The HTTP framework, API database library, and production build/deployment configuration remain undecided. Add dependencies only for agreed work; existing development utilities do not settle all API tooling choices.
- `docs/`: functional specification maintained in Markdown.

## Confirmed architecture

Next.js runs on Bun. A separate Go backend will verify Firebase Authentication ID tokens, enforce community permissions, query PostgreSQL, and handle SeaweedFS uploads if uploads are included. Firebase provides Google login; it does not require Firestore. SeaweedFS is selected for development. Docker Compose runs local PostgreSQL; `make dev` starts the database container and runs Next.js and Go locally. SeaweedFS orchestration is pending. Production hosting and storage remain undecided.

## Development commands

The root Makefile provides shared shortcuts. `make install` downloads dependencies for both services and enables hooks; `make dev` starts Docker PostgreSQL, waits for database health, migrates and seeds it, then runs the Go server and frontend locally; `make seed` starts, migrates, and seeds that same local Docker database. Goose is a Go library invoked by `backend/scripts/migrate`; `make db-migrate`, `make db-status`, and `make db-revision msg=NAME` manage SQL migrations. All Make service/database targets are local development only. Root `.env` or shell variables configure the local stack; service, migration, and seed commands construct their own local database URL from `DEV_DB_*`, ignoring ambient `DATABASE_URL`. `make check` runs frontend lint/typecheck/build, then Go tests and vet. `make help` lists targets. Install GNU Make, Bun, Go, and Docker with Compose separately. Never run seed commands against production databases.

Existing Bun commands, from the repository root:

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
- For Makefile changes, verify `make help`, dry-run the targets, and test startup/shutdown with isolated ports and a separate Docker project. Seed only a disposable development database when validating orchestration.
- For Go server changes, run `make test-backend` and `make vet-backend`; verify `/healthz` and `/readyz`. Database integration tests require `TEST_DATABASE_URL`; report skips explicitly.
- For Git hook changes, run `sh -n .githooks/pre-commit`, verify the hook skips documentation-only commits and blocks failing checks in a temporary repository, and run the frontend checks above.
- For documentation-only changes, run `git diff --check`, verify links and commands. Application lint, typecheck, and build are unnecessary unless executable code or configuration also changed.
- Before committing, inspect `git diff --cached` for scope and sensitive content and run `git diff --cached --check`. Report which checks passed, failed, or could not run; do not claim unrun checks passed.
- The pre-commit hook selects checks from staged paths but runs lint and typecheck against the working tree. It does not build or test behavior. Partially staged changes require reviewing the staged result separately. Do not bypass failing hooks without the user's explicit instruction.
- Go is selected for the backend. Follow checks documented for existing backend utilities; define API service checks as its remaining tooling is selected.

## Pull requests

Use the repository [PR skill](.agents/skills/pr/SKILL.md) when drafting, creating, or updating a pull request. It defines the shared title, description, validation, and issue-link conventions.

## Working conventions

- Preserve existing work and history. Git tracking belongs at the monorepo root; do not initialize nested repositories.
- Keep rendering and browser interaction in Next.js and authoritative application rules in the future backend. Agree on the API contract before connecting services.
- Follow installed Next.js documentation as directed by the frontend `AGENTS.md`.
- Never commit credentials, real environment files, user data, dependencies, build output, or database/storage volumes. Use placeholders in configuration examples.
- Add tooling only when required for agreed work. Do not assume an auth provider, database, backend language, or product feature beyond the spec.
- Run checks appropriate to changed files and report actual results. Documentation-only changes do not need application tests.
- Before publishing, inspect staged files and diffs for unintended content or sensitive data. Public GitHub publishing requires explicit user authorization.
