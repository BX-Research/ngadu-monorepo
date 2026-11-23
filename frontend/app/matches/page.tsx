"use client";

import Link from "next/link";
import { useSyncExternalStore } from "react";
import AppFrame from "../components/AppFrame";
import { demoMatchesSnapshot, parseDemoMatchesSnapshot, subscribeDemoMatches } from "../lib/demoMatches";
import { opponentNames, shortDate, signedPoints, teamNames } from "../lib/display";

type Row = { id: string; title: string; format: string; date: string; score?: [number, number]; won?: boolean; change?: number; href?: string };

const sampleResults: Row[] = [
  { id: "m-101", title: "You vs Dimas Putra", format: "Singles", score: [21, 17], won: true, change: 18, date: "Oct 8" },
  { id: "m-098", title: "You vs Kevin Hartono", format: "Singles", score: [21, 19], won: true, change: 14, date: "Oct 5" },
  { id: "m-094", title: "You & Sarah Wijaya vs Nadia Putri & Rizky Pratama", format: "Doubles", score: [16, 21], won: false, change: -11, date: "Oct 2" },
  { id: "m-090", title: "You vs Sarah Wijaya", format: "Singles", score: [21, 12], won: true, change: 16, date: "Sep 28" },
];

export default function MatchesRoute() {
  const snapshot = useSyncExternalStore(subscribeDemoMatches, demoMatchesSnapshot, () => null);
  const created = parseDemoMatchesSnapshot(snapshot ?? "[]");
  const rows: Row[] = [
    ...created.map((match) => ({
      id: match.id,
      title: `${teamNames(match)} vs ${opponentNames(match)}`,
      format: match.format,
      date: shortDate(match.createdAt),
      score: match.status === "completed" ? match.score : undefined,
      won: match.won,
      change: match.ratingChange,
      href: `/matches/${match.id}`,
    })),
    ...sampleResults,
  ];
  const played = rows.filter((row) => row.score);
  const won = played.filter((row) => row.won).length;
  const winRate = played.length ? Math.round((won / played.length) * 100) : 0;

  return (
    <AppFrame>
      <h1 className="page-title">Results</h1>
      <dl className="tally">
        <div><dt>Played</dt><dd className="numeral">{played.length}</dd></div>
        <div><dt>Won</dt><dd className="numeral">{won}</dd></div>
        <div><dt>Win rate</dt><dd className="numeral">{winRate}%</dd></div>
      </dl>

      <ul className="result-list history-list">
        {rows.map((row) => {
          const body = <>
            <span className={`result-mark ${row.score ? row.won ? "mark-won" : "mark-lost" : "mark-open"}`} aria-label={row.score ? row.won ? "Won" : "Lost" : "Score needed"}>{row.score ? row.won ? "W" : "L" : "?"}</span>
            <span className="result-players">
              {row.title}
              <small>{row.format}, {row.date}</small>
            </span>
            {row.score
              ? <span className="result-outcome"><span className="result-score numeral">{row.score[0]}–{row.score[1]}</span><small className={(row.change ?? 0) < 0 ? "down" : "up"}>{signedPoints(row.change ?? 0)}</small></span>
              : <span className="result-action">Enter score</span>}
          </>;
          return <li key={row.id}>{row.href ? <Link className="result-row" href={row.href}>{body}</Link> : <div className="result-row">{body}</div>}</li>;
        })}
      </ul>

      <Link className="primary-button" href="/matches/new">Record a result</Link>
      <p className="demo-note">Names, scores, and ratings are sample data. Matches you record are saved in this browser only.</p>
    </AppFrame>
  );
}
