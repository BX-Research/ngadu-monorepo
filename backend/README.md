# Backend

A minimal Go server uses the standard library HTTP server for `GET /healthz` (liveness) and `GET /readyz` (PostgreSQL connectivity). A local development user lookup is available; Firebase token verification and protected product APIs are not implemented yet. The readiness check uses `pgx`; this bootstrap does not settle the full API framework or database-access design.

Run `make dev` from the repository root to start local Docker PostgreSQL, this backend, and the frontend. The server listens on localhost port 8080. Root `.env` settings are described in [.env.example](../.env.example). Standalone startup uses `go run ./cmd/api` from this directory. It loads root `.env` and shell development database credentials, ignores ambient `DATABASE_URL`, and binds to localhost:8080. PostgreSQL uses localhost:5432.

Run `make test-backend` and `make vet-backend` from the root to check Go code.

SQL fixtures in `seeds/dev.sql` populate PostgreSQL through Docker’s `psql`; see [scripts/README.md](scripts/README.md) for setup,
usage, and checks. The seed does not choose an HTTP framework or implement login.

Expected responsibilities include verifying Firebase ID tokens, enforcing community permissions, managing sessions and pairings, querying PostgreSQL, and handling uploads through SeaweedFS if uploads are selected.

See [the functional specification](../docs/functional-spec.md) and [repository instructions](../AGENTS.md).

## User lookup (local development)

`GET /api/users/{id}` looks up the application UUID in PostgreSQL using a
parameterized query and returns `id`, `firebase_uid`, `display_name`, and
`created_at` (an RFC 3339 timestamp). No membership information is included.

```bash
make db-seed
make dev-backend
# In another terminal:
curl http://127.0.0.1:8080/api/users/00000000-0000-4000-8000-000000000001
```

Responses are JSON: 200 for a user, 400 for an invalid UUID, 404 when no user
matches, and 503 when the database lookup fails. Errors use an `error` field
without connection details. Requests have a three-second database timeout.
This endpoint is for local inspection; it has no Firebase authentication or
community authorization and must gain those checks before product use.
