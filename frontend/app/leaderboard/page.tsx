"use client";

import { useEffect, useState } from "react";
import AppNavigation from "../components/AppNavigation";
import AppHeader from "../components/AppHeader";
import ProfilePhoto from "../components/ProfilePhoto";
import SportsIllustration from "../components/SportsIllustration";
import { loadLeaderboard, type Leaderboard } from "../lib/leaderboard";

const medals = ["🥇", "🥈", "🥉"];
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
  const podium = [players[1], players[0], players[2]].filter((player) => player !== undefined);
  return <main className="app-shell"><AppNavigation /><section className="main-area"><AppHeader neutral /><div className="content leaderboard-route">
    <section className="leader-hero"><div className="leader-hero-copy"><span className="hero-kicker">🏸 &nbsp; NGADU GLOBAL LADDER</span><h1>Who’s on a<br />winning streak?</h1><p>Every rally gets you closer to the top. Keep playing, keep climbing!</p><span className="hero-season">🏆 &nbsp; Global leaderboard</span></div><div className="hero-doodle"><SportsIllustration /><span className="doodle-spark spark-a">✦</span><span className="doodle-spark spark-b">✧</span></div><span className="hero-ball ball-one">🏸</span><span className="hero-ball ball-two">🎾</span></section>
    {state.status === "loading" && <section className="leaderboard-state" role="status">Loading global standings…</section>}
    {state.status === "error" && <section className="leaderboard-state" role="alert"><h2>Standings are unavailable</h2><p>We couldn’t load the leaderboard. Please try again.</p><button className="primary-button" onClick={() => { setState({ status: "loading" }); setAttempt((value) => value + 1); }}>Try again</button></section>}
    {data && players.length === 0 && <section className="leaderboard-state"><h2>No ranked players yet</h2><p>There are no standings to show yet.</p></section>}
    {data && players.length > 0 && <>
      <div className="podium-intro"><div><p className="eyebrow">TOP PLAYERS</p><h2>Leading the ladder</h2></div><p>Current global standings</p></div>
      <section className="photo-podium" aria-label="Leaderboard top three">{podium.map((player) => <article className={`photo-podium-card photo-rank-${player.rank}`} key={player.user_id}><div className="photo-medal">{medals[player.rank - 1]}</div><ProfilePhoto className="profile-photo" photoUrl={player.photo_url} name={player.display_name} /><span className="photo-rank-label">RANK #{player.rank}</span><h3>{player.display_name}</h3><strong className="photo-elo">{player.rating.toLocaleString()} <small>ELO</small></strong><span className="photo-wins">✨ &nbsp;{player.wins} wins</span></article>)}</section>
      <section className="full-rankings"><div className="full-rankings-head"><div><h2>Global leaderboard</h2><p>Friendly competition, excellent rallies.</p></div><span className="ranking-count">{players.length < data.total_players ? `TOP ${players.length} OF ${data.total_players}` : `${data.total_players} PLAYERS`}</span></div>{players.map((player) => <article className="ranking-row" key={player.user_id}><span className={`rank-number ${player.rank <= 3 ? `rank-medal rank-medal-${player.rank}` : ""}`}>{player.rank <= 3 ? medals[player.rank - 1] : String(player.rank).padStart(2, "0")}</span><div className="ranking-player"><ProfilePhoto className="ranking-photo" photoUrl={player.photo_url} name={player.display_name} /><div><strong>{player.display_name}</strong><small>{player.wins} wins · {player.losses} losses</small></div></div><strong className="ranking-elo">{player.rating.toLocaleString()} <small>ELO</small></strong></article>)}</section>
    </>}
    <p className="demo-note">Fictional development standings · Demo match results are tracked separately.</p>
  </div></section></main>;
}
