# Development seed

The Go command in `seed/` creates a minimal PostgreSQL schema and inserts four
fictional users, two communities, and five memberships. The fixture structs live
in `seed/main.go`. This is a development bootstrap, not a production migration.
The second community is a fixture for testing permissions and does not decide
the application's community scope.

With PostgreSQL running and an empty development database already created, run
from `backend/` in a shell that supports the syntax below:

```bash
export PGHOST=localhost PGPORT=5432 PGUSER=postgres PGDATABASE=ngadu_dev
go run ./scripts/seed
```

Alternatively, run `make seed` from the repository root. It starts the Docker PostgreSQL service and constructs its local connection settings from root `.env` or `DEV_DB_*` shell variables (see [the example](../../.env.example)). It uses that local database rather than an ambient `DATABASE_URL`. Run `make install-backend` first to download Go dependencies. Direct `go run` commands use the environment described above and do not start Docker.

This requires Go 1.26 or newer, as declared in `backend/go.mod`. The `pgx` driver reads the PostgreSQL environment
variables above and a local PostgreSQL password file if needed. Alternatively,
set `DATABASE_URL` to your development connection string. Never commit
credentials. The command does not load `.env` files automatically.

The script runs in one transaction and reports results after commit. It looks up
users by Firebase UID, preserves existing display names, and reuses their stored
application IDs for memberships. Communities use fixed IDs; existing community
names are preserved. Existing membership roles are updated to match the fixture
values, while join dates stay unchanged. No other memberships are deleted.
Reruns do not duplicate fixtures. Conflicting IDs or other constraint violations
roll back the entire run. Existing tables must match `seed/schema.sql`;
`CREATE TABLE IF NOT EXISTS` does not upgrade an older schema.

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
