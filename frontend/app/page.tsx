"use client";

import Link from "next/link";
import { useSyncExternalStore } from "react";
import AppFrame from "./components/AppFrame";
import Scoreboard from "./components/Scoreboard";
import { demoMatchesSnapshot, DEMO_PLAYERS, getDemoStandings, parseDemoMatchesSnapshot, subscribeDemoMatches } from "./lib/demoMatches";
import { opponentNames, shortDate, signedPoints, teamNames } from "./lib/display";

type FeedItem = { id: string; players: string; score: [number, number]; date: string; href?: string };

// Sample results from around the group, shown until real match history exists.
const sampleLast = { format: "Singles", you: "You", opponent: "Dimas Putra", score: [21, 17] as [number, number], won: true, change: 18, date: "Oct 8" };
const sampleFeed: FeedItem[] = [
  { id: "s-1", players: "Nadia Putri vs Sarah Wijaya", score: [21, 15], date: "Oct 8" },
  { id: "s-2", players: "Rizky Pratama vs Kevin Hartono", score: [21, 19], date: "Oct 7" },
  { id: "s-3", players: "Kevin Hartono vs Sarah Wijaya", score: [22, 20], date: "Oct 5" },
];

export default function Home() {
  const snapshot = useSyncExternalStore(subscribeDemoMatches, demoMatchesSnapshot, () => null);
  const created = parseDemoMatchesSnapshot(snapshot ?? "[]");
  const completed = created.filter((match) => match.status === "completed" && match.score);
  const pending = created.filter((match) => match.status === "awaiting_result");
  const standings = getDemoStandings(created);
  const me = standings.find((player) => player.id === DEMO_PLAYERS[0].id) ?? DEMO_PLAYERS[0];

  const latest = completed[0];
  const last = latest
    ? { format: latest.format, you: teamNames(latest), opponent: opponentNames(latest), score: latest.score!, won: latest.won === true, change: latest.ratingChange ?? 0, date: shortDate(latest.createdAt) }
    : sampleLast;

  const feed: FeedItem[] = [
    ...completed.slice(latest ? 1 : 0).map((match) => ({
      id: match.id,
      players: `${teamNames(match)} vs ${opponentNames(match)}`,
      score: match.score!,
      date: shortDate(match.createdAt),
      href: `/matches/${match.id}`,
    })),
    ...sampleFeed,
  ].slice(0, 3);

  return (
    <AppFrame>
      <h1 className="page-title">Your last result</h1>
      <Scoreboard
        title={`${last.format} vs ${last.opponent}`}
        date={last.date}
        left={{ score: last.score[0], label: last.you }}
        right={{ score: last.score[1], label: last.opponent }}
        outcome={last.won ? "won" : "lost"}
      />

      <div className="rating-line">
        <span>Rating</span>
        <strong className="numeral">{me.elo.toLocaleString("en-US")}</strong>
        <span className={`points-chip${last.change < 0 ? " points-down" : ""}`}>{signedPoints(last.change)}</span>
      </div>

      {pending.length > 0 && (
        <section className="list-section">
          <h2 className="section-title">Waiting for a score</h2>
          <ul className="result-list">
            {pending.map((match) => (
              <li key={match.id}>
                <Link className="result-row" href={`/matches/${match.id}`}>
                  <span className="result-players">{teamNames(match)} vs {opponentNames(match)}</span>
                  <span className="result-action">Enter score</span>
                </Link>
              </li>
            ))}
          </ul>
        </section>
      )}

      <section className="list-section">
        <div className="section-head">
          <h2 className="section-title">Recent matches</h2>
          <Link className="text-link" href="/matches">All results</Link>
        </div>
        <ul className="result-list">
          {feed.map((item) => {
            const body = <>
              <span className="result-players">{item.players}</span>
              <span className="result-score numeral">{item.score[0]}–{item.score[1]}</span>
              <span className="result-date">{item.date}</span>
            </>;
            return <li key={item.id}>{item.href ? <Link className="result-row" href={item.href}>{body}</Link> : <div className="result-row">{body}</div>}</li>;
          })}
        </ul>
      </section>

      <Link className="primary-button" href="/matches/new">Record a result</Link>
      <p className="demo-note">Demo data. Matches you record are saved in this browser only.</p>
    </AppFrame>
  );
}
