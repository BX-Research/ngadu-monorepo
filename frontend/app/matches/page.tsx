"use client";

import Link from "next/link";
import { useSyncExternalStore } from "react";
import AppNavigation from "../components/AppNavigation";
import AppHeader from "../components/AppHeader";
import ProfilePortrait from "../components/ProfilePortrait";
import SportsIllustration from "../components/SportsIllustration";
import { demoMatchesSnapshot, parseDemoMatchesSnapshot, subscribeDemoMatches, type DemoMatch } from "../lib/demoMatches";

const sampleMatches = [
  { id: "m-101", result: "WIN", opponent: "Dimas Putra", format: "Singles", score: "21 – 17", elo: "+18", date: "Oct 08 · 20:15", avatar: 3, won: true },
  { id: "m-098", result: "WIN", opponent: "Kevin Hartono", format: "Singles", score: "21 – 19", elo: "+14", date: "Oct 05 · 18:30", avatar: 2, won: true },
  { id: "m-094", result: "LOSS", opponent: "Nadia Putri", format: "Doubles", score: "16 – 21", elo: "−11", date: "Oct 02 · 19:00", avatar: 4, won: false },
  { id: "m-090", result: "WIN", opponent: "Sarah Wijaya", format: "Singles", score: "21 – 12", elo: "+16", date: "Sep 28 · 17:45", avatar: 5, won: true },
];

function describeMatch(match: DemoMatch) {
  const names = match.players.opponents.map((player) => player.name).join(" & ");
  return `${match.format} · vs ${names}`;
}

export default function MatchesRoute() {
  const snapshot = useSyncExternalStore(subscribeDemoMatches, demoMatchesSnapshot, () => null);
  const createdMatches = parseDemoMatchesSnapshot(snapshot ?? "[]");
  const completed = createdMatches.filter((match) => match.status === "completed");
  const playedCount = sampleMatches.length + completed.length;
  const wonCount = sampleMatches.filter((match) => match.won).length + completed.filter((match) => match.won).length;
  const winRate = playedCount ? Math.round((wonCount / playedCount) * 100) : 0;

  return <main className="app-shell"><AppNavigation/><section className="main-area"><AppHeader/><div className="content matches-route">
    <section className="matches-hero"><div><span className="hero-kicker">🏸 &nbsp;YOUR COURT JOURNAL</span><h1>Every match<br/>has a story.</h1><p>Your rallies, results, and little rating victories—all in one place.</p><Link className="primary-button" href="/matches/new"><span>＋</span> Record a match</Link></div><div className="matches-hero-art"><SportsIllustration/><span>Nice rally! ✨</span></div></section>
    <div className="match-summary"><article><span>PLAYED</span><strong>{playedCount.toString().padStart(2, "0")}</strong><small>completed matches</small></article><article><span>WON</span><strong>{wonCount.toString().padStart(2, "0")}</strong><small>keep up the rallies!</small></article><article><span>WIN RATE</span><strong>{winRate}%</strong><small>demo stats</small></article></div>
    <section className="history-route-list"><div className="history-route-head"><div><h2>Match history</h2><p>Latest first · demo results</p></div><span>🏆 &nbsp;Your highlights</span></div>
      {createdMatches.map((match) => { const won = match.won === true; const opponent = match.players.opponents[0]; return <Link className="route-match-card created-match-card" key={match.id} href={`/matches/${match.id}`}><span className={`route-result-icon ${match.status === "completed" ? won ? "won" : "lost" : "pending-match-icon"}`}>{match.status === "completed" ? won ? "W" : "L" : "…"}</span><ProfilePortrait className="route-opponent-photo" variant={opponent?.avatar ?? 0} name={opponent?.name ?? "Opponent"}/><div className="route-match-info"><strong>{opponent?.name ?? "Opponent"} <span className={`result-pill ${match.status === "completed" ? won ? "pill-win" : "pill-loss" : "pill-pending"}`}>{match.status === "completed" ? won ? "WIN" : "LOSS" : "SCORE REQUIRED"}</span></strong><small>{describeMatch(match)}</small><time>{new Date(match.createdAt).toLocaleString("id-ID", { dateStyle: "medium", timeStyle: "short" })}</time></div><div className="route-match-score">{match.status === "completed" ? <><strong>{match.score?.[0]} – {match.score?.[1]}</strong><span className={won ? "up" : "down"}>{match.ratingChange && match.ratingChange > 0 ? "+" : ""}{match.ratingChange} ELO</span></> : <strong>Submit score →</strong>}</div><span className="route-chevron">→</span></Link>; })}
      {sampleMatches.map((match) => <article className="route-match-card" key={match.id}><span className={`route-result-icon ${match.won ? "won" : "lost"}`}>{match.won ? "W" : "L"}</span><ProfilePortrait className="route-opponent-photo" variant={match.avatar} name={match.opponent}/><div className="route-match-info"><strong>{match.opponent} <span className={`result-pill ${match.won ? "pill-win" : "pill-loss"}`}>{match.result}</span></strong><small>{match.format} · Badminton</small><time>{match.date}</time></div><div className="route-match-score"><strong>{match.score}</strong><span className={match.won ? "up" : "down"}>{match.elo} ELO</span></div><span className="route-chevron">→</span></article>)}
    </section>
    <p className="demo-note">Player names, scores, and ratings are sample data; newly created demo matches are stored only in this browser.</p>
  </div></section></main>;
}
