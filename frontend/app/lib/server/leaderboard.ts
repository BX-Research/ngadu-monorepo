type Config = { origin?: string };
export async function serveLeaderboard(config: Config, fetcher: typeof fetch = fetch): Promise<Response> {
  const headers = { "Cache-Control": "no-store" };
  const fail = (status: number, error: string) => Response.json({ error }, { status, headers });
  try {
    const origin = new URL(config.origin ?? "http://127.0.0.1:8080");
    // This proxy exposes only the fixed fictional development read operation.
    if (origin.protocol !== "http:" || !["127.0.0.1", "[::1]"].includes(origin.hostname)
      || origin.username || origin.password || origin.pathname !== "/" || origin.search || origin.hash) {
      return fail(503, "Leaderboard configuration unavailable.");
    }
    const response = await fetcher(new URL("/api/leaderboard", origin), {
      cache: "no-store", signal: AbortSignal.timeout(5000), redirect: "error",
    });
    if (!response.ok) return fail(response.status === 404 ? 404 : 503, "Leaderboard unavailable.");
    return Response.json(await response.json(), { headers });
  } catch {
    return fail(503, "Leaderboard unavailable.");
  }
}
