# Development seed

The SQL file [`../seeds/dev.sql`](../seeds/dev.sql) inserts 24 fictional users,
two communities, 25 memberships, and 25 community rating records into the migrated local development database.
Goose is used as a Go library by `migrate/`; no separate Goose executable is
required. SQL schema migrations live in `backend/migrations/`. The second
community is a permission-testing fixture and does not decide community scope.

From the repository root:

```bash
make db-migrate                         # Start local PostgreSQL and apply schema
make db-status                          # Show migration status
make db-revision msg=add_match_tables    # Generate a SQL migration template
make db-seed                            # Start, migrate, and seed local PostgreSQL
```

`make seed` runs the database workflow plus local SeaweedFS startup and photo
seeding. `make db-seed` operates on PostgreSQL alone. `make storage-seed`
starts storage and seeds only photos. Install Go dependencies
with `make install-backend`. These commands require Go compatible with
`backend/go.mod`; database startup requires Docker Compose and a running daemon.
Creating a migration does not require PostgreSQL or Docker.

Migrations load root `.env` and use `DEV_DB_*` settings (see
[the example](../../.env.example)); shell settings take precedence. Seeding runs
`psql` inside the local Compose database container using its configured
`POSTGRES_USER` and `POSTGRES_DB`. No host `psql` installation is needed.
Both workflows ignore ambient `DATABASE_URL` and host PostgreSQL connection
variables. There are no production Make targets. Directly running
`go run ./scripts/migrate up` from `backend/` uses the same local settings but
does not start Docker. Use `make seed` to migrate and execute the SQL fixtures.

The initial migration preserves tables created by the former bootstrap seeder
with `CREATE TABLE IF NOT EXISTS`. Existing tables must match that original
schema; future schema changes belong in new migrations. Write and review the
SQL in each generated migration; creation does not infer schema changes.

The root Makefile runs the local Go app and Next.js concurrently after database
and storage startup, migrations, and seeding. Its shell trap stops the app processes on exit.

The SQL runs in one transaction; `psql` stops immediately on errors. It looks up
users by Firebase UID, preserves existing display names, and reuses their stored
application IDs for memberships. Communities use fixed IDs; existing community
names are preserved. Existing membership roles are updated to match the fixture
values, while join dates stay unchanged. No other memberships are deleted.
Reruns do not duplicate fixtures. Conflicting IDs or other constraint violations
roll back the entire run. Schema changes are managed separately through Goose migrations.

The main club has 23 ranked members; Drew remains a nonmember. The other club
has two independent ratings. Ties and missing photo references are intentional.
Profiles and ratings use insert-only conflict handling, so edited names, photo
keys, wins/losses, and ratings survive reruns. The original five memberships
retain their existing role-sync behavior; added memberships preserve edited roles.
Every fresh database starts with the same records; locally edited databases may
differ. Seeds never reset application data automatically.

## Profile photo fixtures

`seed-storage/` uses the AWS Go S3 SDK with a loopback-only endpoint, static local
credentials, and path-style requests. It ignores ambient AWS credentials and
never selects AWS endpoints. `DEV_S3_*` settings configure the local bucket,
region, credentials, and port. `make storage-seed` uploads the three SVG portraits
in `backend/seeds/photos/` under `profiles/`, preserving existing objects.
`profiles/missing.svg` is deliberately not uploaded. These illustrations are
repository-owned fictional fixtures. No real people or external image downloads
are involved. Bucket initialization and uploads are repeatable. Conditional uploads recheck
object existence and retry within a 60-second deadline while SeaweedFS initializes
its internal write locks; existing objects are never replaced.

Only the S3 port is published on loopback; internal SeaweedFS UIs are not
published. Credentials switch the S3 gateway into authenticated mode. The
leaderboard signs reads for 15 minutes, while anonymous reads remain denied.
Stopping storage with `make storage-stop` preserves the named volume. Use a
separate Compose project and unused ports for disposable integration checks:

```bash
export COMPOSE_PROJECT_NAME=ngadu-ladder-test
export DEV_DB_PORT=55440 DEV_S3_PORT=18539
make seed
# Point TEST_DATABASE_URL at this local database using your development credentials.
export TEST_S3_ENDPOINT=http://127.0.0.1:18539
make test-backend
make db-stop storage-stop
```

SeaweedFS configuration follows its [S3 credential documentation](https://github.com/seaweedfs/seaweedfs/wiki/S3-Credentials).

The `seed-user-*` and `seed-player-*` Firebase UIDs are placeholders. For authenticated API testing,
set the user's `firebase_uid` in `backend/seeds/dev.sql` to the UID of a
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
