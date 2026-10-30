type Config = { enabled?: string; origin?: string; communityId?: string };
const defaultCommunity = "00000000-0000-4000-8000-000000000101";
const uuid = /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i;

export async function serveLeaderboard(config: Config, fetcher: typeof fetch = fetch): Promise<Response> {
  const headers = { "Cache-Control": "no-store" };
  const fail = (status: number, error: string) => Response.json({ error }, { status, headers });
  if (config.enabled !== "true") return fail(404, "Leaderboard unavailable.");
  try {
    const origin = new URL(config.origin ?? "http://127.0.0.1:8080");
    const communityId = config.communityId ?? defaultCommunity;
    // This proxy exposes only the fixed fictional development read operation.
    if (!uuid.test(communityId) || origin.protocol !== "http:" || !["127.0.0.1", "[::1]"].includes(origin.hostname)
      || origin.username || origin.password || origin.pathname !== "/" || origin.search || origin.hash) {
      return fail(503, "Leaderboard configuration unavailable.");
    }
    const response = await fetcher(new URL(`/api/communities/${communityId}/leaderboard`, origin), {
      cache: "no-store", signal: AbortSignal.timeout(5000), redirect: "error",
    });
    if (!response.ok) return fail(response.status === 404 ? 404 : 503, "Leaderboard unavailable.");
    return Response.json(await response.json(), { headers });
  } catch {
    return fail(503, "Leaderboard unavailable.");
  }
}
