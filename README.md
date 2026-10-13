# Ngadu

A badminton community matchmaking app. This monorepo currently contains the Next.js starter and functional specification. Authentication, backend endpoints, database/storage integration, and matchmaking have not been implemented.

## Structure

```text
ngadu-monorepo/
  AGENTS.md
  package.json          # Root commands delegating to the frontend
  frontend/           # Next.js + React + TypeScript, using Bun
  backend/           # Go backend; API implementation pending
  docs/
    functional-spec.md
```

Services share one Git repository and keep independent toolchains. The frontend manifest and lockfile remain in `frontend`. No extra monorepo build system is required initially.

## Common commands

Install GNU Make, Bun (the version in `package.json`), Go (compatible with `backend/go.mod`), and Docker with Compose on your machine. Start the Docker daemon. The Makefile does not install these tools.

From the repository root:

```bash
make install     # Download dependencies for both services and enable Git hooks
make dev         # Start, migrate, and seed PostgreSQL, then run both apps
make seed        # Start, migrate, and seed local Docker PostgreSQL
make db-migrate  # Apply schema migrations with the Goose Go library
make db-status   # Show applied and pending migrations
make db-revision msg=add_matches # Create a SQL migration
make check       # Frontend lint/types/build, Go tests/vet
make help        # List the available commands
```

`make dev` starts PostgreSQL and waits for its health check, applies migrations, seeds development fixtures, then runs `make dev-backend` and `make dev-frontend` concurrently. Open <http://127.0.0.1:3000>; backend liveness is at <http://127.0.0.1:8080/healthz> and database readiness at <http://127.0.0.1:8080/readyz>. The backend also provides a local development `GET /api/users/{id}` lookup; see [the backend guide](backend/README.md). Protected product APIs and frontend API integration remain pending. Restart `make dev` after changing Go code; Next.js reloads frontend edits automatically.

Ctrl+C stops both app processes. PostgreSQL stays running; `make db-stop` stops it without deleting its named data volume. `make db-up` starts it again and `make db-logs` shows its logs. The recipe waits for both apps; if one exits, stop the remaining app with Ctrl+C.

Optional local overrides are in [.env.example](.env.example); copy it to `.env`, or set the listed variables in your shell. Next.js uses its default port 3000; the Go backend listens on localhost:8080 and PostgreSQL on localhost:5432. There are no custom port settings. The checked-in database password is a local development placeholder. Set a local password before first use. Changing database credentials after the volume is initialized requires updating the database role too; it does not reset existing data.

`make seed` starts the same Docker database and writes the development fixtures described in the [seeder guide](backend/scripts/README.md). Migration and seed commands run Go directly, load root `.env`, and construct a local database URL from `DEV_DB_*`; they ignore ambient `DATABASE_URL` and PostgreSQL connection variables. All Make service/database commands target local development. Goose is pinned in `backend/go.mod`, so no separate Goose installation is needed. `make db-seed` aliases `make seed`, which applies migrations before inserting fixtures. Go database tests require `TEST_DATABASE_URL`; otherwise they are skipped. Existing Bun commands remain available. Use `make install-frontend` or `make install-backend` to install dependencies for one service.

Use `make dev-frontend` to start only Next.js. `make dev-backend` runs only the Go app (`go run ./cmd/api`) with local development settings. Start PostgreSQL separately with `make db-up` (or `make db-seed` to migrate and seed it).

## Frontend development

Use Bun 1.4.3. From the repository root:

```bash
bun run setup:frontend
bun run dev
```

Open <http://localhost:3000>. Check or build the frontend with:

```bash
bun run lint
bun run typecheck
bun run build
bun run start
```

Production start requires a successful build. Development, build, and start commands explicitly select Bun as the Next.js runtime.

## Git hooks

Enable the shared hooks once per clone after installing Bun:

```bash
bun run setup:hooks
```

The pre-commit hook runs lint and typecheck when staged changes affect frontend code, dependencies, or validation tooling. Documentation-only commits skip application checks. Checks run against the working tree, so review partially staged changes separately. Hook scripts use Git's shell and keep LF line endings, including on Windows with Git for Windows.

For frontend changes, also run `bun run build` before committing or opening a PR. Follow the backend documentation for checks on existing utilities; API service tooling remains undecided.

## Pull requests

Agents use the shared [PR skill](.agents/skills/pr/SKILL.md), available as `$pr`. Titles use Conventional Commits style; descriptions include Summary and Validation, with Notes and Related issues when applicable.

## Architecture

- Frontend: Next.js, React, TypeScript, and Bun.
- Authentication: Firebase Authentication with Google sign-in.
- Backend: separate Go service; HTTP framework and remaining API tooling are undecided.
- Application database: PostgreSQL.
- Development file storage: SeaweedFS.
- Local development: Docker Compose runs PostgreSQL; `make dev` runs Next.js and Go locally. SeaweedFS orchestration is pending.

Firebase manages login identities; PostgreSQL will store application data linked to Firebase UIDs. Firestore is not required. Production hosting and storage are undecided.

## Documentation

Read [the functional specification](docs/functional-spec.md). Confirmed choices are separated from proposed features and open decisions.

Read [AGENTS.md](AGENTS.md) and nested instructions before making changes. Keep credentials, environment files, dependencies, build files, and local service data out of Git.
