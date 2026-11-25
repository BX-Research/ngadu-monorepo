"use client";

import Link from "next/link";
import { Suspense, useState, useSyncExternalStore } from "react";
import { useParams } from "next/navigation";
import AppFrame from "../../components/AppFrame";
import AppHeader from "../../components/AppHeader";
import Scoreboard from "../../components/Scoreboard";
import { demoMatchesSnapshot, parseDemoMatchesSnapshot, subscribeDemoMatches, submitDemoMatchResult } from "../../lib/demoMatches";
import { opponentNames, shortDate, signedPoints, teamNames } from "../../lib/display";

export default function MatchDetailRoute() {
  return <Suspense fallback={<main className="app-shell"><AppHeader /><div className="content"><p className="state-panel" role="status">Loading match…</p></div></main>}><MatchDetailContent /></Suspense>;
}

function MatchDetailContent() {
  const { id } = useParams<{ id: string }>();
  const snapshot = useSyncExternalStore(subscribeDemoMatches, demoMatchesSnapshot, () => null);
  const match = parseDemoMatchesSnapshot(snapshot ?? "[]").find((item) => item.id === id);
  const loaded = snapshot !== null;
  const [scores, setScores] = useState<[string, string]>(["", ""]);
  const [error, setError] = useState("");

  function submit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const ownScore = Number(scores[0]);
    const opponentScore = Number(scores[1]);
    const saved = submitDemoMatchResult(id, ownScore, opponentScore);
    if (!saved) {
      setError(ownScore === opponentScore ? "Scores can’t be tied. The side with more points wins." : "Enter a whole number for each side.");
      return;
    }
    setError("");
  }

  if (!loaded) return <AppFrame><p className="state-panel" role="status">Loading match…</p></AppFrame>;
  if (!match) {
    return (
      <AppFrame>
        <section className="state-panel">
          <h2>Match not found</h2>
          <p>Demo matches live in the browser they were recorded in. This one may be on another device.</p>
          <Link className="primary-button" href="/matches">Back to results</Link>
        </section>
      </AppFrame>
    );
  }

  const ownPlayers = [match.players.you, ...(match.players.ally ? [match.players.ally] : [])];
  const title = `${match.format} vs ${opponentNames(match)}`;
  const left = teamNames(match);
  const right = opponentNames(match);
  const awaiting = match.status === "awaiting_result";

  return (
    <AppFrame>
      <h1 className="page-title">{awaiting ? "Enter the score" : "Final score"}</h1>

      {awaiting ? (
        <form onSubmit={submit}>
          <Scoreboard
            title={title}
            date={shortDate(match.createdAt)}
            left={{ label: left, score: <input className="score-input" aria-label="Your side’s score" type="number" inputMode="numeric" min="0" step="1" required value={scores[0]} onChange={(event) => setScores([event.target.value, scores[1]])} placeholder="0" /> }}
            right={{ label: right, score: <input className="score-input" aria-label="Opponents’ score" type="number" inputMode="numeric" min="0" step="1" required value={scores[1]} onChange={(event) => setScores([scores[0], event.target.value])} placeholder="0" /> }}
          />
          {error && <p className="notice" role="alert">{error}</p>}
          <button className="primary-button" type="submit">Save result</button>
          <p className="demo-note">The side with more points wins. Ties can’t be saved.</p>
        </form>
      ) : (
        <>
          <Scoreboard
            title={title}
            date={shortDate(match.createdAt)}
            left={{ label: left, score: match.score?.[0] }}
            right={{ label: right, score: match.score?.[1] }}
            outcome={match.won ? "won" : "lost"}
          >
            {match.resultMessage && <p className="scoreboard-quip">“{match.resultMessage}”</p>}
          </Scoreboard>
          <div className="rating-line">
            <span>Rating change</span>
            <span className={`points-chip${(match.ratingChange ?? 0) < 0 ? " points-down" : ""}`}>{signedPoints(match.ratingChange ?? 0)}</span>
          </div>
        </>
      )}

      <section className="list-section">
        <h2 className="section-title">Who played</h2>
        <ul className="result-list">
          {[...ownPlayers.map((player) => ({ player, side: "Your side" })), ...match.players.opponents.map((player) => ({ player, side: "Opponents" }))].map(({ player, side }) => (
            <li className="result-row" key={player.id}>
              <span className="result-players">{player.id === match.players.you.id ? `${player.name} (you)` : player.name}<small>{side}</small></span>
              <span className="result-score numeral">{player.elo.toLocaleString("en-US")}</span>
            </li>
          ))}
        </ul>
      </section>

      {!awaiting && <p className="demo-note">Saved in this browser. Demo ratings don’t update the leaderboard.</p>}
      <div className="link-pair">
        <Link className="text-link" href="/matches">All results</Link>
        <Link className="text-link" href="/matches/new">Record another</Link>
      </div>
    </AppFrame>
  );
}
