"use client";

import Link from "next/link";
import { useEffect, useState, useSyncExternalStore } from "react";
import AppNavigation from "./components/AppNavigation";
import AppHeader from "./components/AppHeader";
import ProfilePortrait from "./components/ProfilePortrait";
import ProfilePhoto from "./components/ProfilePhoto";
import { demoMatchesSnapshot, DEMO_PLAYERS, parseDemoMatchesSnapshot, subscribeDemoMatches, type DemoMatch } from "./lib/demoMatches";
import { loadLeaderboard, type Leaderboard } from "./lib/leaderboard";

type LeaderboardState = { status: "loading" } | { status: "error" } | { status: "ready"; data: Leaderboard };

type Match = { id: string | number; sport: string; format: string; opponent: string; ally?: string; opponentRating: number | string; date: string; status: "upcoming" | "completed"; score?: [number, number]; won?: boolean; change?: number };

const initialMatches: Match[] = [
  { id: 1, sport: "Badminton", format: "Singles", opponent: "Kevin Hartono", opponentRating: 1476, date: "Today · 18:30", status: "upcoming" },
  { id: 2, sport: "Badminton", format: "Doubles", opponent: "Rizky Pratama & Dimas Putra", opponentRating: "1,512 / 1,328", date: "Tomorrow · 19:00", status: "upcoming" },
  { id: 3, sport: "Badminton", format: "Singles", opponent: "Dimas Putra", opponentRating: 1328, date: "Oct 08 · 20:15", status: "completed", score: [21, 17], won: true, change: 18 },
];

export default function Home() {
  const snapshot = useSyncExternalStore(subscribeDemoMatches, demoMatchesSnapshot, () => null);
  const [leaderboardState, setLeaderboardState] = useState<LeaderboardState>({ status: "loading" });
  const [attempt, setAttempt] = useState(0);
  useEffect(() => {
    const controller = new AbortController();
    loadLeaderboard(controller.signal).then((data) => {
      if (!controller.signal.aborted) setLeaderboardState({ status: "ready", data });
    }).catch(() => {
      if (!controller.signal.aborted) setLeaderboardState({ status: "error" });
    });
    return () => controller.abort();
  }, [attempt]);
  const createdMatches = parseDemoMatchesSnapshot(snapshot ?? "[]");
  const liveMatches: Match[] = createdMatches.map((match: DemoMatch) => ({
    id: match.id,
    sport: match.sport,
    format: match.format,
    opponent: match.players.opponents.map((player) => player.name).join(" & "),
    opponentRating: match.players.opponents.map((player) => player.elo.toLocaleString()).join(" / "),
    date: new Date(match.createdAt).toLocaleDateString("en-US", { month: "short", day: "2-digit" }),
    status: match.status === "completed" ? "completed" : "upcoming",
    score: match.score,
    won: match.won,
    change: match.ratingChange,
  }));
  const matches = [...liveMatches, ...initialMatches];
  const leaderboard = leaderboardState.status === "ready" ? leaderboardState.data : null;
  const topPlayer = leaderboard?.players[0];
  const highlights = leaderboard?.players.slice(0, 5) ?? [];

  return (
    <main className="app-shell">
      <AppNavigation />
      <section className="main-area">
        <AppHeader neutral />
        <div className="content overview-route">
          <div className="welcome">
            <p className="eyebrow">BADMINTON · GLOBAL LADDER</p>
            <h1>Your game,<br />your community.</h1>
            <p className="subhead">Explore the current local-development standings.</p>
          </div>

          <div className="stat-grid" aria-label="Global leaderboard summary">
            <article className="stat-card rating-card">
              <div className="stat-top">TOP GLOBAL RATING <span className="rating-badge">{topPlayer ? `#1 · ${topPlayer.display_name}` : "Leaderboard"}</span></div>
              <div className="rating-value">{topPlayer ? topPlayer.rating.toLocaleString() : "—"}<span>ELO</span></div>
              <p className="rating-change">{topPlayer ? `${topPlayer.wins} wins · ${topPlayer.losses} losses` : "Current leader"}</p>
            </article>
            <article className="stat-card">
              <h2>Ranked players</h2>
              <div className="stat-number">{leaderboard ? leaderboard.total_players : "—"}</div>
            </article>
            <article className="stat-card">
              <h2>Players shown</h2>
              <div className="stat-number">{leaderboard ? leaderboard.players.length : "—"}<span>/10</span></div>
            </article>
          </div>

          <Link className="primary-button create-match-button" href="/matches/new"><span aria-hidden="true">+</span> Create a match</Link>
          <section className="full-rankings dashboard-ladder" aria-labelledby="dashboard-ladder-title">
            <div className="full-rankings-head">
              <div><h2 id="dashboard-ladder-title">Global leaderboard</h2><p>Players ranked by their current rating.</p></div>
              <Link className="text-button" href="/leaderboard">View all <span aria-hidden="true">→</span></Link>
            </div>
            {leaderboardState.status === "loading" && <p className="leaderboard-state" role="status">Loading standings…</p>}
            {leaderboardState.status === "error" && <div className="leaderboard-state" role="alert"><p>Standings are unavailable. Try again.</p><button className="primary-button" onClick={() => { setLeaderboardState({ status: "loading" }); setAttempt((value) => value + 1); }}>Retry</button></div>}
            {leaderboard && leaderboard.players.length === 0 && <p className="leaderboard-state">No ranked players yet.</p>}
            {highlights.map((player) => (
              <Link className="ranking-row dashboard-ranking-row" href={`/players/${encodeURIComponent(player.user_id)}`} key={player.user_id} aria-label={`View ${player.display_name}'s read-only profile`}>
                <span className={`rank-number ${player.rank <= 3 ? `rank-medal rank-medal-${player.rank}` : ""}`}>{player.rank <= 3 ? ["🥇", "🥈", "🥉"][player.rank - 1] : String(player.rank).padStart(2, "0")}</span>
                <span className="ranking-player"><ProfilePhoto className="ranking-photo" photoUrl={player.photo_url} name={player.display_name} /><span><strong>{player.display_name}</strong><small>{player.wins} wins · {player.losses} losses</small></span></span>
                <strong className="ranking-elo">{player.rating.toLocaleString()} <small>ELO</small></strong>
              </Link>
            ))}
          </section>

          <div className="section-heading">
            <h2>Recent demo matches</h2>
            <Link className="text-button" href="/matches">View matches <span aria-hidden="true">→</span></Link>
          </div>
          <div className="recent-matches">
            {matches.map((match) => (
              <article className="overview-match-card" key={match.id}>
                <div className="player-cell">
                  <div className="overview-avatar-stack">
                    {match.opponent.split(" & ").map((name) => {
                      const player = DEMO_PLAYERS.find((candidate) => candidate.name === name);
                      return <ProfilePortrait key={name} className="overview-profile-avatar" variant={player?.avatar ?? 0} name={name}/>;
                    })}
                  </div>
                  <div className="overview-match-info">
                    <strong>{match.opponent}</strong>
                    <p>{match.format} · {match.date}</p>
                    <small>ELO {match.opponentRating}</small>
                  </div>
                </div>
                <div className="overview-match-result">
                  <span className={`result ${match.status === "upcoming" ? "pending" : match.won ? "win" : "loss"}`}>
                    {match.status === "upcoming" ? "Awaiting score" : `${match.won ? "Won" : "Lost"} ${match.score?.[0]}–${match.score?.[1]}`}
                  </span>
                  {match.status === "completed" && <span className={(match.change ?? 0) < 0 ? "down" : "up"}>{(match.change ?? 0) > 0 ? "+" : ""}{match.change} pts</span>}
                  {typeof match.id === "string" && <Link href={`/matches/${match.id}`} className="match-detail-link">{match.status === "completed" ? "View match" : "Enter score"}<span aria-hidden="true"> →</span></Link>}
                </div>
              </article>
            ))}
          </div>
          <p className="demo-note">Match activity above is sample/local-browser data and is not connected to the backend.</p>
          <footer>© 2026 Ngadu <span>Play your game. Own your rank.</span></footer>
        </div>
      </section>
    </main>
  );
}
