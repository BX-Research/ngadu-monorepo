# Backend

A minimal Go server uses Fiber v3 for `GET /healthz` (liveness) and `GET /readyz` (PostgreSQL connectivity). A local development user lookup is available; Firebase token verification and protected product APIs are not implemented yet. The readiness check uses `pgx`; Fiber is the selected HTTP framework; the full API database-access design remains undecided.

Run `make dev` from the repository root to start local Docker PostgreSQL, this backend, and the frontend. The server listens on localhost port 8080. Root `.env` settings are described in [.env.example](../.env.example). Install Air as described in the [root README](../README.md#common-commands). Standalone startup uses `air` from this directory; `.air.toml` builds `./app` and reloads the server after Go edits. For startup without live reload, use `go run ./app`. Direct startup loads root `.env` without overriding shell variables and requires `DATABASE_URL`. Make startup supplies a local URL from required `DEV_DB_USER`, `DEV_DB_PASSWORD`, and `DEV_DB_NAME`, overriding any other `DATABASE_URL` for that process without editing `.env`. No database credentials have fallback values. The server binds to localhost:8080. PostgreSQL uses localhost:5432.

Run `make test-backend` and `make vet-backend` from the root to check Go code.

SQL fixtures in `seeds/dev.sql` populate PostgreSQL through Docker’s `psql`; see [scripts/README.md](scripts/README.md) for setup,
usage, and checks. The seed does not implement login.

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

## Leaderboard (local development)

`make dev-backend` explicitly enables the fictional leaderboard using
`ENABLE_DEV_LEADERBOARD=true`. Direct Go startup requires that flag to opt in;
without it the route is absent (404). The backend binds to loopback. Production
use requires Firebase ID token verification and community membership checks.

`GET /api/communities/{communityId}/leaderboard` returns:

```json
{
  "community": { "id": "00000000-0000-4000-8000-000000000101", "name": "Demo Shuttle Club" },
  "total_players": 23,
  "players": [
    { "rank": 1, "user_id": "00000000-0000-4000-8000-000000000201", "display_name": "Adi Demo", "rating": 1750, "wins": 30, "losses": 5, "photo_url": null }
  ]
}
```

The example shows one row; responses contain up to 10. `photo_url` is either a
15-minute signed private S3 URL or null. Object keys and storage credentials are
not returned. No bucket listing or per-player object checks occur during reads.
Missing objects fall back in the browser; invalid storage configuration yields
null URLs without failing the leaderboard.

Ranks are sequential, ordered by rating descending then user UUID ascending.
Count and list use a single database statement. Each request queries stored
rating snapshots with a three-second timeout and `Cache-Control: no-store`.
Malformed UUIDs return 400, unknown communities 404, empty communities 200 with
an empty array, and database failure 503 without internal error details.
The dashboard and browser-local demo matches do not update these ratings.

The Next.js endpoint `/api/leaderboard` forwards only this fixed operation to a
loopback backend, using server-only `NGADU_BACKEND_ORIGIN` (default
`http://127.0.0.1:8080`) and `DEV_LEADERBOARD_COMMUNITY_ID` (default main fixture
club). It also requires `ENABLE_DEV_LEADERBOARD=true`. Set optional overrides in
the shell running Next.js, or in its local environment file; root `.env` loading
is owned by the Go utilities/Compose. Never use `NEXT_PUBLIC_` for credentials.

See the [seed guide](scripts/README.md) for deterministic fixtures and storage
commands. S3 clients use the SDK's [custom endpoint support](https://docs.aws.amazon.com/sdk-for-go/v2/developer-guide/configure-endpoints.html).
