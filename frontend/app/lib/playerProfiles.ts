import type { Leaderboard, RankedPlayer } from "./leaderboard";

export function findRankedPlayer(leaderboard: Leaderboard, userId: string): RankedPlayer | null {
  return leaderboard.players.find((player) => player.user_id === userId) ?? null;
}
