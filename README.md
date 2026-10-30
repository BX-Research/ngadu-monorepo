# Ngadu

A badminton matchmaking app for friends. The leaderboard reads fictional global development rankings from PostgreSQL through the Go API and profile photos from private local S3 storage. Authentication and matchmaking are still pending; dashboard and match pages use browser demo data.

## Structure

```text
ngadu-monorepo/
  AGENTS.md
  package.json          # Root commands delegating to the frontend
  frontend/           # Next.js + React + TypeScript, using Bun
  backend/           # Go/Fiber backend and database migrations
  docs/
    functional-spec.md
```

Services share one Git repository and keep independent toolchains. The frontend manifest and lockfile remain in `frontend`. No extra monorepo build system is required initially.

## Prerequisites

Install these tools before running the project:

```text
Git                    Version control and shared Git hooks
GNU Make               Local development commands
Bun 1.4.3              Frontend runtime and package manager (package.json)
Go 1.26.0 or compatible Backend toolchain (backend/go.mod)
Docker Engine/Desktop  Local PostgreSQL and SeaweedFS containers; daemon must be running
Docker Compose v2      Docker Compose plugin (docker compose)
Air                    Go backend live reload
```

Install [Air](https://github.com/air-verse/air) and ensure its installation directory is on your `PATH`:

```bash
go install github.com/air-verse/air@latest
# For the default Go installation directory on Linux/macOS/WSL:
export PATH="$(go env GOPATH)/bin:$PATH"
air -v
```

If you set `GOBIN`, add that directory to `PATH` instead. The Makefile does not install these prerequisite tools.

## Common commands

From the repository root:

```bash
make install     # Download dependencies for both services and enable Git hooks
make dev         # Seed PostgreSQL and private S3 photos, then run both apps
make seed        # Seed local PostgreSQL and SeaweedFS photo fixtures
make db-seed     # Seed only PostgreSQL
make storage-seed # Seed only local private S3 photo fixtures
make db-migrate  # Apply schema migrations with the Goose Go library
make db-status   # Show applied and pending migrations
make db-revision msg=add_matches # Create a SQL migration
make check       # Frontend lint/types/build, Go tests/vet
make help        # List the available commands
```

`make dev` first copies `.env.example` to `.env` if `.env` is missing, preserving any existing `.env`. It then starts PostgreSQL, applies migrations, seeds development fixtures, starts SeaweedFS and seeds profile photos, then runs `make dev-backend` and `make dev-frontend` concurrently. Open <http://127.0.0.1:3000>; backend liveness is at <http://127.0.0.1:8080/healthz> and database readiness at <http://127.0.0.1:8080/readyz>. The backend also provides a local development `GET /api/users/{id}` lookup; see [the backend guide](backend/README.md). Open `/leaderboard` to view database standings with signed photo URLs and a placeholder for missing images. Protected product APIs remain pending. Air rebuilds and restarts the backend after Go edits; Next.js reloads frontend edits automatically.

Ctrl+C stops both app processes. PostgreSQL and SeaweedFS stay running; `make storage-stop` stops storage while preserving its data. `make db-stop` stops PostgreSQL without deleting its named data volume. `make db-up` starts it again and `make db-logs` shows its logs. The recipe waits for both apps; if one exits, stop the remaining app with Ctrl+C.

Required database settings are in [.env.example](.env.example); copy it to root `.env`, or set the listed variables in your shell. Database settings have no defaults. Next.js uses its default port 3000; the Go backend listens on localhost:8080 and PostgreSQL on localhost:5432. Use `DEV_DB_PORT` and `DEV_S3_PORT` for local infrastructure ports, `PORT` for Go, `FRONTEND_PORT` for Next.js through Make, and `NGADU_BACKEND_ORIGIN` in the Next.js shell when Go uses a different port. The checked-in database password is a local development placeholder. Set a local password before first use. Changing database credentials after the volume is initialized requires updating the database role too; it does not reset existing data.

`make seed` starts the same Docker database and writes the development fixtures described in the [seeder guide](backend/scripts/README.md). Migrations run Go directly, load root `.env`, and construct a local database URL from `DEV_DB_*`. Seeding executes `backend/seeds/dev.sql` using `psql` inside the Compose database container. Both ignore ambient `DATABASE_URL` and host PostgreSQL connection variables. All Make service/database commands target local development. Goose is pinned in `backend/go.mod`, so no separate Goose installation is needed. `make db-seed` applies migrations and seeds PostgreSQL alone; `make seed` also starts and seeds SeaweedFS. Fresh setups receive 24 fictional users and the same rating snapshots. Reruns preserve profile/rating edits and existing photo objects. Go database tests require `TEST_DATABASE_URL`; otherwise they are skipped. Existing Bun commands remain available. Use `make install-frontend` or `make install-backend` to install dependencies for one service.

Use `make dev-frontend` to start only Next.js. `make dev-backend` runs only the Go app through Air using `backend/.air.toml`, which builds `./app`. You can also run `air` directly from `backend` with an explicit `DATABASE_URL` in root `.env` or your shell. `make dev` and `make dev-backend` override `DATABASE_URL` for the backend process with a local URL built from required `DEV_DB_*` settings; they never rewrite `.env`. Override the executable with `make dev-backend AIR=/path/to/air` if needed. Start PostgreSQL separately with `make db-up` (or `make db-seed` to migrate and seed it).

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
bun run --cwd frontend test
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
- Backend: separate Go service using Fiber v3; remaining API tooling is undecided.
- Application database: PostgreSQL.
- Development file storage: SeaweedFS.
- Local development: Docker Compose runs PostgreSQL; `make dev` runs Next.js and Go locally. SeaweedFS exposes private S3 storage on loopback port 8333.

Firebase manages login identities; PostgreSQL will store application data linked to Firebase UIDs. Firestore is not required. Production hosting and storage are undecided.

## Documentation

Read [the functional specification](docs/functional-spec.md). Confirmed choices are separated from proposed features and open decisions.

Read [AGENTS.md](AGENTS.md) and nested instructions before making changes. Keep credentials, environment files, dependencies, build files, and local service data out of Git.
