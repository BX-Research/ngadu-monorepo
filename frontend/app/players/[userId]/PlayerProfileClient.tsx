"use client";

import Link from "next/link";
import { useEffect, useState } from "react";
import AppNavigation from "../../components/AppNavigation";
import AppHeader from "../../components/AppHeader";
import ProfilePhoto from "../../components/ProfilePhoto";
import { loadLeaderboard, type Leaderboard } from "../../lib/leaderboard";
import { findRankedPlayer } from "../../lib/playerProfiles";

type State = { status: "loading" } | { status: "error" } | { status: "ready"; data: Leaderboard };

type Props = { userId: string };

export default function PlayerProfileClient({ userId }: Props) {
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

  const player = state.status === "ready" ? findRankedPlayer(state.data, userId) : null;

  return (
    <main className="app-shell">
      <AppNavigation />
      <section className="main-area">
        <AppHeader neutral />
        <div className="content player-profile-route">
          <Link className="text-button" href="/leaderboard">← Back to leaderboard</Link>
          <p className="eyebrow">PLAYER PROFILE · READ ONLY</p>
          {state.status === "loading" && <section className="leaderboard-state" role="status">Loading player profile…</section>}
          {state.status === "error" && <section className="leaderboard-state" role="alert"><h2>Profile unavailable</h2><p>Could not load the current leaderboard.</p><button className="primary-button" onClick={() => { setState({ status: "loading" }); setAttempt((value) => value + 1); }}>Try again</button></section>}
          {state.status === "ready" && !player && <section className="leaderboard-state"><h1>Profile not available</h1><p>This read-only page currently uses the leaderboard API, which only returns the top 10 players. A player-profile endpoint is needed to look up everyone else.</p></section>}
          {player && <>
            <section className="player-profile-card" aria-labelledby="player-profile-name">
              <div className="player-profile-identity">
                <ProfilePhoto className="player-profile-photo" photoUrl={player.photo_url} name={player.display_name} />
                <div><span className="hero-kicker">GLOBAL RANK #{player.rank}</span><h1 id="player-profile-name">{player.display_name}</h1><p>Read-only local-development profile</p></div>
              </div>
              <div className="player-profile-stats">
                <article><span>RATING</span><strong>{player.rating.toLocaleString()}</strong><small>ELO</small></article>
                <article><span>WINS</span><strong>{player.wins}</strong></article>
                <article><span>LOSSES</span><strong>{player.losses}</strong></article>
              </div>
            </section>
            <p className="demo-note">Values come from the global leaderboard. Profile editing and match history are not available from the current API.</p>
          </>}
        </div>
      </section>
    </main>
  );
}
