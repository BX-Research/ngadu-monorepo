# Player profile API contract

Status: proposed for backend implementation; the profile endpoint does not exist yet.

## Current frontend data source

The dashboard and current read-only profile pages consume the existing same-origin `GET /api/leaderboard` endpoint. The Next.js route forwards to the Go backend's `GET /api/leaderboard` and returns up to 10 ranked players. Each entry contains `rank`, `user_id`, `display_name`, `rating`, `wins`, `losses`, and nullable `photo_url`; the response also contains `total_players`.

The current profile page can therefore show only players included in that top-10 response. It must not call the local-development `GET /api/users/{id}` lookup from the browser: that endpoint is an inspection utility and exposes `firebase_uid`. The frontend does not forward that field.

## Proposed endpoint

### `GET /api/players/{user_id}`

Return the public, read-only profile for one player. The Next.js same-origin route should proxy the Go endpoint and expose only the response fields below. The backend must not return Firebase UID, email, or storage object keys.

`user_id` must be a canonical UUID. The endpoint should return `Cache-Control: no-store` because `photo_url` may be a short-lived signed URL and profile data should reflect the latest rating.

### Success response — `200 OK`

```json
{
  "user_id": "00000000-0000-4000-8000-000000000201",
  "display_name": "Adi Demo",
  "rank": 1,
  "rating": 1750,
  "wins": 30,
  "losses": 5,
  "photo_url": null
}
```

Fields:

| Field | Type | Meaning |
| --- | --- | --- |
| `user_id` | UUID string | Application user ID used by the route |
| `display_name` | string | Player's display name |
| `rank` | integer or `null` | Global rank; `null` when the player has no rating record |
| `rating` | integer or `null` | Global rating; `null` when the player has no rating record |
| `wins` | integer | Wins from the global rating record; `0` if unranked |
| `losses` | integer | Losses from the global rating record; `0` if unranked |
| `photo_url` | string or `null` | Short-lived signed private photo URL, or `null` when unavailable |

Do not include `updated_at` unless the product decides to display it. Do not accept client-supplied identity fields or allow writes through this endpoint.

### Errors

All errors use JSON `{ "error": "..." }` and omit database/storage details.

| Status | Meaning |
| --- | --- |
| `400 Bad Request` | `user_id` is not a valid UUID |
| `404 Not Found` | No player exists for that ID |
| `503 Service Unavailable` | Database/profile lookup is unavailable |

### Privacy and deployment

This contract describes read-only fictional development data. The existing leaderboard is unauthenticated in local development. Production visibility and authentication/authorization policy for player profiles are not decided; do not expose this proposed endpoint publicly until that policy is approved and implemented.

## Frontend integration boundary

Until `GET /api/players/{user_id}` is implemented and proxied, the frontend profile page uses the existing leaderboard response and shows an explanatory not-in-top-10 state for IDs absent from its returned list. It does not fabricate profile data or fall back to the local inspection endpoint.
