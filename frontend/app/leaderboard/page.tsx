"use client";

import { useSyncExternalStore } from "react";
import AppNavigation from "../components/AppNavigation";
import AppHeader from "../components/AppHeader";
import ProfilePortrait from "../components/ProfilePortrait";
import SportsIllustration from "../components/SportsIllustration";
import { demoMatchesSnapshot, getDemoStandings, parseDemoMatchesSnapshot, subscribeDemoMatches } from "../lib/demoMatches";

const medals = ["🥇", "🥈", "🥉"];

export default function LeaderboardRoute() {
  const snapshot = useSyncExternalStore(subscribeDemoMatches, demoMatchesSnapshot, () => null);
  const players = getDemoStandings(parseDemoMatchesSnapshot(snapshot ?? "[]"));
  const podium = [players[1], players[0], players[2]];
  return <main className="app-shell"><AppNavigation/><section className="main-area"><AppHeader/><div className="content leaderboard-route">
    <section className="leader-hero"><div className="leader-hero-copy"><span className="hero-kicker">🏸 &nbsp; NGADU CLUB RANKS</span><h1>Who’s on a<br/>winning streak?</h1><p>Every rally gets you closer to the top. Keep playing, keep climbing!</p><span className="hero-season">🏆 &nbsp; Season 1 · Singles + Doubles</span></div><div className="hero-doodle"><SportsIllustration/><span className="doodle-spark spark-a">✦</span><span className="doodle-spark spark-b">✧</span></div><span className="hero-ball ball-one">🏸</span><span className="hero-ball ball-two">🎾</span></section>
    <div className="podium-intro"><div><p className="eyebrow">THE BIG THREE</p><h2>Top players</h2></div><p>Current club standings</p></div>
    <section className="photo-podium" aria-label="Leaderboard top three">{podium.map((player,index)=>{const rank=index===1?1:index===0?2:3;return <article className={`photo-podium-card photo-rank-${rank}`} key={player.name}><div className="photo-medal">{medals[rank-1]}</div><ProfilePortrait className="profile-photo" variant={player.avatar} name={player.name}/><span className="photo-rank-label">RANK #{rank}</span><h3>{player.name}</h3><strong className="photo-elo">{player.elo.toLocaleString()} <small>ELO</small></strong><span className="photo-wins">✨ &nbsp;{player.wins} wins</span></article>})}</section>
    <section className="full-rankings"><div className="full-rankings-head"><div><h2>Club leaderboard</h2><p>Friendly competition, excellent rallies.</p></div><span className="ranking-count">6 PLAYERS</span></div>{players.map((player,index)=><article className="ranking-row" key={player.name}><span className={`rank-number ${index<3?`rank-medal rank-medal-${index+1}`:""}`}>{index<3?medals[index]:String(index+1).padStart(2,"0")}</span><div className="ranking-player"><ProfilePortrait className="ranking-photo" variant={player.avatar} name={player.name}/><div><strong>{player.name}{player.name==="Alex Morgan"&&<em>YOU</em>}</strong><small>{player.wins} wins · Badminton</small></div></div><strong className="ranking-elo">{player.elo.toLocaleString()} <small>ELO</small></strong></article>)}</section>
    <p className="demo-note">Illustrated player avatars and rankings · demo data</p>
  </div></section></main>;
}
