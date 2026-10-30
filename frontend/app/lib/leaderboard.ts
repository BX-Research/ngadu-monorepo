export type RankedPlayer = {
  rank: number;
  user_id: string;
  display_name: string;
  rating: number;
  wins: number;
  losses: number;
  photo_url: string | null;
};

export type Leaderboard = {
  total_players: number;
  players: RankedPlayer[];
};

export async function loadLeaderboard(signal?: AbortSignal, fetcher: typeof fetch = fetch): Promise<Leaderboard> {
  const response = await fetcher("/api/leaderboard", { cache: "no-store", signal });
  if (!response.ok) throw new Error("Leaderboard unavailable. Please try again.");
  return response.json();
}
