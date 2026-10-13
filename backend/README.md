# Backend

The separate API service is not implemented yet. A Go development seeder uses
`pgx` to populate PostgreSQL; see [scripts/README.md](scripts/README.md) for setup,
usage, and checks. The seed does not choose an HTTP framework or implement login.

Expected responsibilities include verifying Firebase ID tokens, enforcing community permissions, managing sessions and pairings, querying PostgreSQL, and handling uploads through SeaweedFS if uploads are selected.

See [the functional specification](../docs/functional-spec.md) and [repository instructions](../AGENTS.md).
