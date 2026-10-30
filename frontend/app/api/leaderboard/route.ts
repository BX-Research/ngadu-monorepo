import { connection } from "next/server";
import { serveLeaderboard } from "../../lib/server/leaderboard";

export async function GET() {
  await connection();
  return serveLeaderboard({
    origin: process.env.NGADU_BACKEND_ORIGIN,
  });
}
