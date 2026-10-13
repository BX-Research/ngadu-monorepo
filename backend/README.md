# Backend

A minimal Go server uses the standard library HTTP server for `GET /healthz` (liveness) and `GET /readyz` (PostgreSQL connectivity). Application endpoints and Firebase token verification are not implemented yet. The readiness check reuses `pgx`, already used by the seeder; this bootstrap does not settle the full API framework or database-access design.

Run `make dev` from the repository root to start local Docker PostgreSQL, this backend, and the frontend. The server listens on localhost port 8080 by default. Root `.env` settings are described in [.env.example](../.env.example). Standalone startup uses `go run ./cmd/api` from this directory with `DATABASE_URL` (or PostgreSQL environment variables) and optional `API_ADDR` configured in your shell.

Run `make test-backend` and `make vet-backend` from the root to check Go code.

A Go development seeder uses
`pgx` to populate PostgreSQL; see [scripts/README.md](scripts/README.md) for setup,
usage, and checks. The seed does not choose an HTTP framework or implement login.

Expected responsibilities include verifying Firebase ID tokens, enforcing community permissions, managing sessions and pairings, querying PostgreSQL, and handling uploads through SeaweedFS if uploads are selected.

See [the functional specification](../docs/functional-spec.md) and [repository instructions](../AGENTS.md).
