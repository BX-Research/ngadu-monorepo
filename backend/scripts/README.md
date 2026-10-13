# Development seed

The Go command in `seed/` inserts four fictional users, two communities, and
five memberships into the migrated local development database. Goose is used as
a Go library by `migrate/`; no separate Goose executable is required. SQL schema
migrations live in `backend/migrations/`. The fixture structs live in
`seed/main.go`. The second community is a permission-testing fixture and does
not decide the application's community scope.

From the repository root:

```bash
make db-migrate                         # Start local PostgreSQL and apply schema
make db-status                          # Show migration status
make db-revision msg=add_match_tables    # Generate a SQL migration template
make db-seed                            # Start, migrate, and seed local PostgreSQL
```

`make seed` is an alias for the same seeding workflow. Install Go dependencies
with `make install-backend`. These commands require Go compatible with
`backend/go.mod`; database startup requires Docker Compose and a running daemon.
Creating a migration does not require PostgreSQL or Docker.

Migration and seed commands load root `.env` and use `DEV_DB_*` settings (see
[the example](../../.env.example)); shell settings take precedence. They always
connect to `127.0.0.1` and ignore `DATABASE_URL` and PostgreSQL connection
variables. There are no production Make targets. Directly running
`go run ./scripts/migrate up` or `go run ./scripts/seed` from `backend/` uses
these same local settings but does not start Docker. Migrate before seeding.

The initial migration preserves tables created by the former bootstrap seeder
with `CREATE TABLE IF NOT EXISTS`. Existing tables must match that original
schema; future schema changes belong in new migrations. Write and review the
SQL in each generated migration; creation does not infer schema changes.

The shared development supervisor and its tests are also under `backend/scripts/`.
It uses Bun to coordinate the local Go server and Next.js and handle shutdown.

The script runs in one transaction and reports results after commit. It looks up
users by Firebase UID, preserves existing display names, and reuses their stored
application IDs for memberships. Communities use fixed IDs; existing community
names are preserved. Existing membership roles are updated to match the fixture
values, while join dates stay unchanged. No other memberships are deleted.
Reruns do not duplicate fixtures. Conflicting IDs or other constraint violations
roll back the entire run. Schema changes are managed separately through Goose migrations.

The `seed-user-*` Firebase UIDs are placeholders. For authenticated API testing,
set the user's `FirebaseUID` in the local fixture definition to the UID of a
Firebase development or Auth emulator user before the first run. If already
seeded, update the matching local database user's UID as well. Seeding PostgreSQL
does not create Firebase accounts, passwords, or login tokens. Roles are fixture
data and must not be granted automatically on login.

To run the integration tests against a disposable PostgreSQL database:

```bash
export TEST_DATABASE_URL='postgres://localhost/ngadu_seed_test?sslmode=disable'
go test -v ./...
go vet ./...
```

Tests create and remove isolated schemas inside that database. Without
`TEST_DATABASE_URL`, database tests are skipped. Always point seed and test
commands at development databases.
