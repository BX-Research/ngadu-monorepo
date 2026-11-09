import { test } from "node:test";
import assert from "node:assert/strict";
import { findRankedPlayer } from "./playerProfiles";
import type { Leaderboard } from "./leaderboard";

const leaderboard: Leaderboard = {
  total_players: 12,
  players: [
    { rank: 1, user_id: "player-1", display_name: "Ari", rating: 1700, wins: 12, losses: 2, photo_url: null },
    { rank: 2, user_id: "player-2", display_name: "Bima", rating: 1600, wins: 9, losses: 4, photo_url: null },
  ],
};

test("finds a profile by the exact backend user ID", () => {
  assert.deepEqual(findRankedPlayer(leaderboard, "player-2"), leaderboard.players[1]);
});

test("does not invent profile data for users outside the returned top ten", () => {
  assert.equal(findRankedPlayer(leaderboard, "player-11"), null);
});
