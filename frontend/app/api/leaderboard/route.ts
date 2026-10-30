import { connection } from "next/server";
import { serveLeaderboard } from "../../lib/server/leaderboard";

export async function GET() {
  await connection();
  return serveLeaderboard({
    enabled: process.env.ENABLE_DEV_LEADERBOARD,
    origin: process.env.NGADU_BACKEND_ORIGIN,
    communityId: process.env.DEV_LEADERBOARD_COMMUNITY_ID,
  });
}
