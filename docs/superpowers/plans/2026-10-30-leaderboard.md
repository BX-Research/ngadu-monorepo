# Global leaderboard implementation plan

Status: Implemented and verified against the global-ladder decision.

## Goal

Provide one application-wide ladder backed by PostgreSQL and private local S3
photos. The leaderboard API is available by default and returns the top ten.
Dashboard and match pages remain separate browser demo flows.

## Data and migration

- One `player_ratings` row per user, keyed by `user_id`, with rating, wins,
  losses, and updated time. The user reference cascades on deletion.
- Index ratings descending and user UUID ascending for deterministic ordering.
- Apply a new migration; retain historical migration files unchanged.
- Reset the old ratings as explicitly selected by the user. Preserve users,
  display names, and stored avatar objects. Remove the obsolete grouping tables.
- The reset is forward-only: discarded rows cannot be reconstructed by a down
  migration. `make seed` restores the shared global baseline after migration.
- Seed all 24 fictional players, including Drew, with 24 global ratings. Include
  ties and missing avatar objects. Later seed runs preserve
  locally edited profiles and ratings and resolve stored IDs by Firebase UID.

## API and frontend

- Go exposes `GET /api/leaderboard`. Return `{ total_players, players }`; player
  fields are rank, user_id, display_name, rating, wins, losses, photo_url.
- Return up to ten sequentially ranked rows, ordered by rating descending then
  user UUID ascending. A single SQL statement produces the count and list.
- An empty ladder returns 200 with `players: []`; database errors return 503
  without internal details. Queries have a three-second timeout.
- The Next.js `/api/leaderboard` route proxies only this fixed read operation
  to the server-only loopback `NGADU_BACKEND_ORIGIN`. No ladder-selection setting
  or runtime enable flag remains. Responses/fetches use `no-store`.
- Use global ladder labels throughout the UI, with loading, empty, error/retry,
  and success states. Support zero, one, two, or three podium players.
- Preserve private signed photo reads, a 15-minute signing lifetime, local
  placeholder images, and fallback for missing or failed images.
- Backend and development frontend bind to loopback. Authentication and future
  application authorization remain separate work; this local read API is unauthenticated.

## Storage and development setup

- Keep SeaweedFS 4.51, its private development bucket, and committed JPEG portraits.
- `make seed` prepares PostgreSQL and local S3; `make db-seed` runs database-only
  setup; `make storage-seed` preserves existing objects while adding fixtures.
- Storage seeding retries conditional-write readiness failures within 60 seconds.
- Keep credentials server-side; ignore ambient AWS credentials and remote storage
  endpoints in development seed commands.

## Validation

- Migration test: old ratings are removed, grouping tables no longer exist,
  profiles survive, and seeding restores exactly 24 ratings.
- Seed tests: rerun preservation, fixed baseline, alternate stored user IDs,
  atomic rollback, unique user ratings, valid foreign keys, nonnegative wins.
- DAL/API tests: global top ten and count, stable ties, rating refresh, empty
  results, default route availability, photo signing failure, database failure.
- Frontend tests and browser checks: global route and response, no selection
  parameter, loading, empty/small podiums, error/retry, and image fallbacks.
- Run frontend lint/typecheck/build, Bun tests, Go tests/vet with disposable DB
  integration enabled, and health/readiness checks.
- Verify fresh Goose migrations plus repeated seeds in an isolated Compose
  project before applying the reset to the configured local development database.
- Review diffs and run `git diff --check`; do not publish cloud resources.

## Verification results

- Frontend lint, typecheck, three Bun tests, and production build passed. The
  build still reports the existing Next.js security advisory and match-detail
  Suspense warning.
- Go tests and vet passed with disposable PostgreSQL integration and private S3
  reads enabled; no integration tests were skipped in the final run.
- Fresh Goose migrations and repeated SQL seeds passed in an isolated Compose
  project. Upgrade tests prove reset behavior, profile preservation, removal of
  obsolete tables, and restoration of the 24-player shared baseline.
- Browser checks passed for the global top ten/count, database rating updates,
  phone/desktop layout, loading, empty/one/two-player podiums, errors/retry,
  signed images, and all photo fallbacks; no runtime errors were observed.
- Health/readiness and global response checks passed; the old scoped URL is 404.
- The configured local development database was migrated and seeded: 24 global
  ratings, zero obsolete grouping tables. Private photo fixtures were preserved.
- Independent review found no actionable issues. Disposable test services were
  stopped and removed after verification; the normal local stack stays running.

## Configuration and dependency injection revision

- `main.go` delegates startup to `server.Run`; the composition root loads
  `config.Config`, opens concrete clients, and injects them into `server.New`.
- `server.Dependencies` supplies user lookup, leaderboard read, photo signing,
  and readiness interfaces. The HTTP constructor never loads environment files.
- Generic `S3_*` configuration replaces `DevConfig`. Application clients support
  configured remote endpoints; only local seeding enforces loopback storage.
- `make dev` prepares `S3_ENDPOINT_HOST=localhost` and
  `S3_ENDPOINT_URL=http://localhost:9000` in root `.env`. `S3_PORT` supports isolated
  setups. Existing settings and quoted/multiline credentials are preserved.
- Go tests/vet passed with PostgreSQL, private S3 reads on port 9000, and Compose
  dotenv compatibility checks enabled, with no integration skips. Regression
  tests cover dependency injection, optional photos, generic configuration,
  repeated environment setup, credential punctuation, and multiline values.
- An isolated `make dev` check passed startup, migrations/seeding, API and Next.js
  proxy requests, signed photo reads, readiness, and Ctrl+C shutdown. The temporary
  checkout used the installed dependencies with Webpack and disabled Go VCS
  stamping because it had no Git metadata. The existing Next.js advisory appeared.
- Both review findings about dotenv value preservation were fixed and re-tested;
  the final focused review found no actionable issues.

## Derived JPEG avatars revision

- Use the canonical `avatars/{userId}.jpg` key for signing, writes and seeding.
- Remove `users.profile_photo_key` through a new migration; preserve users/ratings.
- Seed 14 per-user avatars from three committed JPEG illustrations; resolve
  actual UUIDs by Firebase UID and preserve existing objects on repeated runs.
- Keep the frontend SVG placeholder for missing/failed avatar reads.
- Add a backend-only upload service with an injected writer. Normalize JPEG,
  PNG, GIF and WebP to static 512×512 quality-85 JPEGs, white background,
  JPEG EXIF orientation correction and stripped source metadata. Limit input
  to 10 MiB and 25 million pixels before full decode.
- Set private/no-cache object headers for deterministic-key replacement.
- No HTTP upload endpoint is included, as selected by the user. The future
  caller must authorize the user ID before invoking the pipeline.

Avatar validation results: all Go tests and vet passed with PostgreSQL, S3 and
Compose dotenv integration enabled, with no skips. JPEG upload/replacement and
private/missing object behavior passed against S3. Browser checks passed for
JPEG display, SVG fallbacks, errors/retry and small podiums. Cropping before
resizing bounds extreme-aspect-ratio allocations; regressions for 100000×1 and
1×100000 inputs pass. The focused review found no remaining actionable issues.

The configured local database now has no `profile_photo_key` column, and local
S3 on port 9000 contains the per-user JPEG fixtures. The real local storage
passed private/signed/missing-avatar and normalization/replacement smoke checks.
Disposable avatar test services and volumes were removed afterward.
