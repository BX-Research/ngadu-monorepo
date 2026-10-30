"use client";

import Link from "next/link";
import { useSyncExternalStore } from "react";
import AppNavigation from "./components/AppNavigation";
import AppHeader from "./components/AppHeader";
import ProfilePortrait from "./components/ProfilePortrait";
import { demoMatchesSnapshot, DEMO_PLAYERS, getDemoStandings, parseDemoMatchesSnapshot, subscribeDemoMatches, type DemoMatch } from "./lib/demoMatches";

type Match = { id: string | number; sport: string; format: string; opponent: string; ally?: string; opponentRating: number | string; date: string; status: "upcoming" | "completed"; score?: [number, number]; won?: boolean; change?: number };

const initialMatches: Match[] = [
  { id: 1, sport: "Badminton", format: "Singles", opponent: "Kevin Hartono", opponentRating: 1476, date: "Today · 18:30", status: "upcoming" },
  { id: 2, sport: "Badminton", format: "Doubles", opponent: "Rizky Pratama & Dimas Putra", opponentRating: "1,512 / 1,328", date: "Tomorrow · 19:00", status: "upcoming" },
  { id: 3, sport: "Badminton", format: "Singles", opponent: "Dimas Putra", opponentRating: 1328, date: "Oct 08 · 20:15", status: "completed", score: [21, 17], won: true, change: 18 },
];

export default function Home() {
  const snapshot = useSyncExternalStore(subscribeDemoMatches, demoMatchesSnapshot, () => null);
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
  const standings = getDemoStandings(createdMatches);
  const currentPlayer = standings.find((player) => player.id === DEMO_PLAYERS[0].id) ?? DEMO_PLAYERS[0];
  const playedMatches = matches.filter((match) => match.status === "completed");
  const winRate = playedMatches.length ? Math.round(playedMatches.filter((match) => match.won).length / playedMatches.length * 100) : 0;
  const lastCompleted = playedMatches[0];
  const lastChange = lastCompleted?.change ?? 0;

  return (
    <main className="app-shell">
      <AppNavigation />
      <section className="main-area">
        <AppHeader />
        <div className="content overview-route">
          <div className="welcome">
            <p className="eyebrow">BADMINTON</p>
            <h1>Your game,<br />your ranking.</h1>
            <p className="subhead">Play with friends. Keep getting better.</p>
          </div>

          <div className="stat-grid">
            <article className="stat-card rating-card">
              <div className="stat-top">YOUR ELO RATING <span className="rating-badge">Demo rating</span></div>
              <div className="rating-value">{currentPlayer.elo.toLocaleString()}<span>pts</span></div>
              <p className="rating-change">{lastChange > 0 ? "+" : ""}{lastChange} last match</p>
              <svg className="rating-chart" viewBox="0 0 120 55" fill="none" aria-hidden="true"><path d="M2 50 18 38 31 40 45 26 59 29 73 16 88 20 104 8 117 4" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/><circle cx="117" cy="4" r="3" fill="currentColor"/></svg>
            </article>
            <article className="stat-card">
              <h2>Matches played</h2>
              <div className="stat-number">{playedMatches.length}</div>
            </article>
            <article className="stat-card">
              <h2>Win rate</h2>
              <div className="stat-number">{winRate}<span>%</span></div>
            </article>
          </div>

          <Link className="primary-button create-match-button" href="/matches/new"><span aria-hidden="true">+</span> Create a match</Link>
          <div className="section-heading">
            <h2>Recent matches</h2>
            <Link className="text-button" href="/matches">View all <span aria-hidden="true">→</span></Link>
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
          <p className="demo-note">Demo data · matches are saved in this browser.</p>
          <footer>© 2026 Ngadu <span>Play your game. Own your rank.</span></footer>
        </div>
      </section>
    </main>
  );
}
