"use client";

import { useEffect, useState } from "react";
import AppFrame from "../components/AppFrame";
import ProfilePhoto from "../components/ProfilePhoto";
import { loadLeaderboard, type Leaderboard } from "../lib/leaderboard";

type State = { status: "loading" } | { status: "error" } | { status: "ready"; data: Leaderboard };

export default function LeaderboardRoute() {
  const [state, setState] = useState<State>({ status: "loading" });
  const [attempt, setAttempt] = useState(0);
  useEffect(() => {
    const controller = new AbortController();
    loadLeaderboard(controller.signal).then((data) => {
      if (!controller.signal.aborted) setState({ status: "ready", data });
    }).catch(() => {
      if (!controller.signal.aborted) setState({ status: "error" });
    });
    return () => controller.abort();
  }, [attempt]);

  const data = state.status === "ready" ? state.data : null;
  const players = data?.players ?? [];
  // Podium reads 2 · 1 · 3 so the leader stands in the middle.
  const podium = [players[1], players[0], players[2]].filter((player) => player !== undefined);

  return (
    <AppFrame neutral>
      <div className="page-head">
        <h1 className="page-title">Leaderboard</h1>
        {data && players.length > 0 && (
          <span className="page-meta">{players.length < data.total_players ? `Top ${players.length} of ${data.total_players}` : `${data.total_players} players`}</span>
        )}
      </div>

      {state.status === "loading" && <section className="state-panel" role="status">Loading standings…</section>}
      {state.status === "error" && (
        <section className="state-panel" role="alert">
          <h2>Standings didn’t load</h2>
          <p>The leaderboard server didn’t respond. Check that the backend is running, then try again.</p>
          <button className="primary-button" onClick={() => { setState({ status: "loading" }); setAttempt((value) => value + 1); }}>Try again</button>
        </section>
      )}
      {data && players.length === 0 && (
        <section className="state-panel">
          <h2>No ranked players yet</h2>
          <p>Standings appear here once players have a rating.</p>
        </section>
      )}

      {data && players.length > 0 && <>
        <section className="podium" aria-label="Top three players">
          {podium.map((player) => (
            <article className={`podium-step podium-rank-${player.rank}`} key={player.user_id}>
              <span className="podium-rank numeral">{player.rank}</span>
              <ProfilePhoto className="podium-photo" photoUrl={player.photo_url} name={player.display_name} />
              <h2>{player.display_name}</h2>
              <strong className="numeral">{player.rating.toLocaleString("en-US")}</strong>
            </article>
          ))}
        </section>

        <ol className="ladder">
          {players.map((player) => (
            <li className="ladder-row" key={player.user_id}>
              <span className="ladder-rank numeral">{player.rank}</span>
              <ProfilePhoto className="ladder-photo" photoUrl={player.photo_url} name={player.display_name} />
              <span className="ladder-name">
                <strong>{player.display_name}</strong>
                <small>{player.wins} won, {player.losses} lost</small>
              </span>
              <strong className="ladder-rating numeral">{player.rating.toLocaleString("en-US")}</strong>
            </li>
          ))}
        </ol>
      </>}

      <p className="demo-note">Fictional development standings. Results recorded in this browser don’t change them.</p>
    </AppFrame>
  );
}
