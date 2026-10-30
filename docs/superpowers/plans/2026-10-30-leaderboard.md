# Database-backed leaderboard plan

Status: Implemented and verified.

## Goal and scope

Replace the leaderboard's hardcoded players and browser-local standings with
fictional PostgreSQL records and profile images read from S3-compatible storage.
Preserve the current page layout and display a local placeholder whenever an
image is absent or cannot load.

This plan covers `/leaderboard`. Dashboard and match creation/history remain
separate demo flows. Their browser-local results will not change database
ratings. The shared header must stop identifying the leaderboard visitor as
the fictional Alex Morgan; use a neutral header on this page until real
authentication exists.

## Proposed design

### Approaches

1. **Recommended: Go leaderboard API, PostgreSQL ratings, private S3 bucket.**
   Next.js requests the Go API through a same-origin route. Go queries current
   ratings and provides short-lived signed photo URLs. This follows the existing
   service boundary and keeps storage credentials on the backend.
2. **Go API with public photo URLs.** Fewer signing requirements, but requires
   public bucket access. Avoid making public profile photos a default decision.
3. **Next.js queries PostgreSQL directly.** Fewer hops, but moves application
   queries into the rendering service, contrary to the repository architecture.

The private S3 approach was approved in chat.

### Data

- Add nullable `users.profile_photo_key`; store an object key, never credentials
  or a signed URL.
- Add `player_ratings` with `community_id`, `user_id`, integer `rating`,
  nonnegative `wins` and `losses`, and `updated_at`.
- Use `(community_id, user_id)` as the primary key and a composite foreign key
  to memberships, so nonmembers cannot acquire a community ranking.
- Index `(community_id, rating DESC, user_id ASC)` for ordered reads.
- Seed at least 20 fictional players in the demo club, including tied ratings,
  null photo keys, and missing-object references. Preserve existing fixtures,
  unrelated users, and edited profiles; reruns must not reset existing ratings.
- Ratings and wins are fixture snapshots. Implementing match completion,
  rating formulas, singles/doubles separation, and seasons is separate work.

### API contract

Proposed route: `GET /api/communities/:communityId/leaderboard`.

Return `{ community: { id, name }, total_players, players }`. Each player contains
`rank`, `user_id`, `display_name`, `rating`, `wins`, `losses`, and nullable
`photo_url`. Return the first 10 players, ordered by rating descending, then
user UUID ascending; ranks are sequential positions, including for tied ratings.
Use a single database statement or consistent snapshot for count and list.

Validate UUIDs, use parameterized SQL and a three-second query timeout. Return
400 for malformed IDs, 404 for unknown communities, 200 with an empty list for
an empty community, and 503 for unavailable PostgreSQL. Errors must not include
connection details. Storage failure must not fail the standings response.

Each page visit fetches current database standings. Use `Cache-Control: no-store`
on the API/proxy and disable fetch caching; no Redis, polling, or leaderboard
cache in this iteration. Photo objects may use their own browser caching.

Authentication is currently unimplemented. Initially expose this route only in
an explicit local-development mode with loopback binding and fictional records;
default to disabling it elsewhere. Do not treat a community ID as authorization.
Real deployment requires Firebase verification and membership enforcement.

### Storage

Use the already selected SeaweedFS for local development, exposing its S3 API.
Add local-only orchestration, health checks, persistent volumes, and an idempotent
storage seed command. Use an isolated development bucket and upload a few small
repository-owned fixture images; leave some users without photos.

Configure endpoint, region, bucket, access key, secret key, path-style access,
and signing lifetime on the backend. Use a browser-reachable endpoint when
signing URLs. Proposed lifetime: 15 minutes. Never send credentials to Next.js
client code or commit a real environment file.

Use an S3 client supporting SeaweedFS and private-bucket signed GET URLs. Before
implementation, verify current SeaweedFS startup/configuration and the chosen
SDK against official documentation; this plan does not pin unverified commands
or versions. Do not list the bucket or check every object on leaderboard reads:
sign configured keys, then let image loading handle missing objects. Signing
configuration failure yields null photo URLs and leaves standings available.

### Frontend

Add a Next.js route handler forwarding to the configured Go endpoint. Keep the
backend origin server-only, construct the upstream path from validated IDs, and
allow only the fixed leaderboard operation. Read installed Next.js docs before
choosing route-handler/fetch APIs, particularly with `cacheComponents` enabled.

Fetch standings on page mount and provide loading, empty, retryable error, and
success states. Remove this page's `getDemoStandings` and localStorage
subscription. Render up to 10 rows, dynamic community/count labels, and a podium
that supports zero, one, or two players. Show both the displayed count and total
when the club has more than 10 ranked members. Remove unsupported season labels
and the name-based `YOU` badge.

Add a reusable profile image component with a local default image. Use it for
the podium and rows. Fall back for null/blank URLs, missing objects, denied or
expired URLs, invalid image bytes, and network failures. Prevent fallback loops;
reset error state if the source changes. Keep accessible names and fixed image
dimensions. No image upload UI is included.

## Implementation tasks

- [x] **1. Schema and fixtures.** Add a new Goose migration without modifying
  the initial migration. Extend `backend/seeds/dev.sql` and its tests. Update the
  integration-test setup to apply all migrations in order. Verify unique
  memberships, constraints, tied ordering fixtures, and rerun preservation on a
  disposable PostgreSQL database.
- [x] **2. Local S3 storage.** Update `docker-compose.yaml`, `.env.example`, and
  Make targets; add a storage seeding script and fixture assets. Verify an
  uploaded object can be read through a signed URL and a private object cannot
  be fetched anonymously. Verify repeat seeding and persistent restart behavior.
  Ensure ordinary database seed operations remain local-only.
- [x] **3. Backend read API.** Add focused leaderboard files under
  `backend/app/types`, `dal`, `services`, and `routes`, plus an S3 signing adapter.
  Wire them in `backend/app/main.go`. Write failing tests first for top-10
  selection, stable ties, community isolation, counts, invalid/unknown IDs,
  empty clubs, database errors, storage errors, and development-mode gating.
  Keep the existing Fiber/pgx architecture.
- [x] **4. Frontend integration.** Add the same-origin route handler and typed
  leaderboard client, update `frontend/app/leaderboard/page.tsx`, add the profile
  image component/default asset, and allow a neutral shared header. Exercise
  the full database-to-page flow plus loading/error/retry, fewer than three
  players, more than ten players, and all photo fallback cases. Use an existing
  test runner if available; otherwise document browser verification without
  introducing a broad test-tooling project.
- [x] **5. Documentation and verification.** Update the functional spec to
  distinguish this local leaderboard implementation from still-proposed
  production behavior; document API/storage configuration, seed commands, and
  the separation from demo match results in the relevant READMEs.

## Acceptance and validation

- Editing a seeded rating in PostgreSQL and reopening the leaderboard changes
  the displayed order without frontend source edits or a rebuild.
- The page displays only ranked members of its configured community, sorted
  deterministically, with a maximum of 10 rows and an accurate total.
- At least one displayed avatar loads from private S3-compatible storage;
  missing/unavailable images show the local default without hiding rankings.
- Database failures show a retryable error, never hardcoded sample standings.
- Frontend: `bun run lint`, `bun run typecheck`, `bun run build`.
- Backend: `make test-backend`, `make vet-backend`; run integration tests using
  `TEST_DATABASE_URL` against a disposable local database and report any skips.
- Verify `/healthz` and `/readyz`; the current source only registers `/healthz`
  despite documentation claiming readiness. Add/test database readiness if it
  remains absent when implementation starts.
- Orchestration: `make help`, target dry runs, startup/shutdown with isolated
  ports and a separate Docker project; seed only disposable local resources.
  Parameterize database ports if needed for isolated validation.
- Review browser/network output for exposed credentials and request paths;
  run `git diff --check`. Inspect staged changes and run
  `git diff --cached --check` before any authorized commit.

## Verification results

- Frontend lint, typecheck, three Bun tests, and production build passed.
- The first build stopped for an existing Next.js security advisory; the retry
  passed while retaining the upgrade recommendation. An existing
  `/matches/[id]` client-hook/Suspense warning also remains outside this scope.
- All Go tests and vet passed with PostgreSQL and S3 integration enabled; no
  integration tests were skipped in the final run.
- Browser checks passed against the built frontend and real Go/PostgreSQL/S3:
  top ten, ties/counts, database rating refresh, signed photos and all fallbacks,
  empty/one/two-player podiums, loading, failure/retry, phone/desktop layout,
  and no browser runtime errors.
- `make help`, target dry runs, isolated Compose startup/shutdown/restart,
  Goose migrations, repeated seeding, health/readiness, and real database outage
  behavior passed. API and built frontend startup/shutdown were checked
  separately; the full Air plus Next.js `make dev` process was not launched
  alongside the checkout's already-running development servers.
- A first-start SeaweedFS lock-ring race was found on the normal local stack.
  Bounded conditional-write retries were added and tested; a new disposable
  storage volume seeded successfully. Both independent reviews found no
  actionable issues after the fix.
- The configured local development database and private photo bucket were
  seeded successfully. Disposable test services were stopped/removed afterward.
- Database fixtures include 24 fictional users, 25 memberships, and 25 rating
  records (23 in the main club, two in the second). Fresh setups share this
  baseline; local edits survive reruns.
