import { test } from "node:test";
import assert from "node:assert/strict";
import { serveLeaderboard } from "./server/leaderboard";
import { loadLeaderboard } from "./leaderboard";

const data = { total_players: 0, players: [] };

test("proxy rejects invalid configuration", async () => {
  const neverFetch = (async () => { throw new Error("must not fetch"); }) as typeof fetch;
  assert.equal((await serveLeaderboard({ origin: "https://example.com" }, neverFetch)).status, 503);
});

test("proxy fetches the global ladder without caching and hides upstream errors", async () => {
  const fetcher = (async (input, options) => {
    assert.equal(String(input), "http://127.0.0.1:8080/api/leaderboard");
    assert.equal(options?.cache, "no-store");
    return Response.json(data);
  }) as typeof fetch;
  const response = await serveLeaderboard({}, fetcher);
  assert.equal(response.headers.get("Cache-Control"), "no-store");
  assert.deepEqual(await response.json(), data);
  const failing = (async () => new Response("secret database details", { status: 500 })) as typeof fetch;
  const error = await serveLeaderboard({}, failing);
  assert.equal(error.status, 503);
  assert.equal((await error.text()).includes("secret"), false);
});

test("client returns empty standings and reports unavailable API", async () => {
  const fetcher = (async (_input, options) => {
    assert.equal(options?.cache, "no-store");
    return Response.json(data);
  }) as typeof fetch;
  assert.deepEqual(await loadLeaderboard(undefined, fetcher), data);
  const failing = (async () => new Response("unavailable", { status: 503 })) as typeof fetch;
  await assert.rejects(loadLeaderboard(undefined, failing), /unavailable/i);
});
