"use client";

import Link from "next/link";
import { useSyncExternalStore } from "react";
import AppSidebar from "./components/AppSidebar";
import AppSocialRail from "./components/AppSocialRail";
import AppHeader from "./components/AppHeader";
import ProfilePortrait from "./components/ProfilePortrait";
import { demoMatchesSnapshot, DEMO_PLAYERS, getDemoStandings, parseDemoMatchesSnapshot, subscribeDemoMatches, type DemoMatch } from "./lib/demoMatches";

type Match = { id: string | number; sport: string; format: string; opponent: string; ally?: string; opponentRating: number | string; date: string; status: "upcoming" | "completed"; score?: [number, number]; won?: boolean; change?: number };



const initialMatches: Match[] = [
  { id: 1, sport: "Badminton", format: "Singles", opponent: "Kevin Hartono", opponentRating: 1476, date: "Today · 18:30", status: "upcoming" },
  { id: 2, sport: "Badminton", format: "Doubles", opponent: "Rizky Pratama & Dimas Putra", opponentRating: "1,512 / 1,328", date: "Tomorrow · 19:00", status: "upcoming" },
  { id: 3, sport: "Badminton", format: "Singles", opponent: "Dimas Putra", opponentRating: 1328, date: "Oct 08 · 20:15", status: "completed", score: [21, 17], won: true, change: 18 },
];

const leaderboardPlayers = [{name:"Nadia Putri",elo:1588,wins:28,initials:"NP"},{name:"Rizky Pratama",elo:1512,wins:24,initials:"RP"},{name:"Kevin Hartono",elo:1476,wins:21,initials:"KH"},{name:"Sarah Wijaya",elo:1394,wins:19,initials:"SW"},{name:"Dimas Putra",elo:1328,wins:17,initials:"DP"},{name:"Alex Morgan",elo:1248,wins:14,initials:"AM"}];

function LeaderboardPage() {
  const podium=[leaderboardPlayers[1],leaderboardPlayers[0],leaderboardPlayers[2]];
  return <section className="leaderboard-page"><div className="leader-intro"><div><p className="eyebrow">BADMINTON · ALL FORMATS</p><h2>Leaderboard</h2><p className="subhead">The friendly race to the top. Keep playing, keep climbing!</p></div><span className="leader-season">🏸 &nbsp;Season 1</span></div><div className="podium" aria-label="Top three players">{podium.map((player,index)=>{const rank=index===1?1:index===0?2:3;return <article className={`podium-card rank-${rank}`} key={player.name}><span className="medal">{rank===1?"🏆":rank===2?"🥈":"🥉"}</span><div className="podium-avatar">{player.initials}</div><span className="podium-rank">#{rank} PLAYER</span><strong className="podium-name">{player.name}</strong><span className="podium-elo">{player.elo.toLocaleString()} <small>ELO</small></span><span className="podium-wins">{player.wins} wins</span></article>})}</div><div className="leader-list"><div className="leader-list-head"><h3>All players</h3><span>Updated after each match</span></div>{leaderboardPlayers.map((player,index)=><div className={`leader-row ${index<3?"top-row":""}`} key={player.name}><span className="leader-position">{index<3?["🥇","🥈","🥉"][index]:String(index+1).padStart(2,"0")}</span><div className="leader-person"><span className="leader-avatar">{player.initials}</span><strong>{player.name}{player.name==="Alex Morgan"&&<em>You</em>}</strong></div><span className="leader-wins">{player.wins} wins</span><strong className="leader-rating">{player.elo.toLocaleString()} <small>ELO</small></strong></div>)}</div><p className="demo-note">Demo leaderboard · sample player accounts and ratings</p></section>;
}

function MatchesHistory({ matches }: { matches: Match[] }) {
  const history = matches.filter((match) => match.status === "completed");
  const wins = history.filter((match) => (match.change ?? 0) > 0).length;
  return <section className="history-page"><div className="history-intro"><div><p className="eyebrow">YOUR COURT JOURNAL</p><h2>Match history</h2><p className="subhead">Every game counts. Look how far you’ve come!</p></div><span className="history-total">🏸 &nbsp;{history.length} matches</span></div><div className="history-stats"><article><span>PLAYED</span><strong>{history.length}</strong></article><article><span>WON</span><strong>{wins}</strong></article><article><span>WIN RATE</span><strong>{history.length?Math.round(wins/history.length*100):0}%</strong></article></div><div className="history-list"><div className="history-list-title"><h3>Past matches</h3><span>Most recent first</span></div>{history.length?history.map(match=>{const won=(match.change??0)>0;return <article className="history-card" key={match.id}><div className={`history-outcome ${won?"history-win":"history-loss"}`}>{won?"W":"L"}</div><div className="history-detail"><div className="history-topline"><strong>{won?"Victory":"Tough match"} <span className={`history-tag ${won?"tag-win":"tag-loss"}`}>{won?"WON":"LOST"}</span></strong><time>{match.date}</time></div><p>{match.format} · vs {match.opponent}</p><small>{match.sport} · Opponent ELO {match.opponentRating}</small></div><div className="history-score"><strong>{match.score?.[0]} – {match.score?.[1]}</strong><span className={won?"up":"down"}>{(match.change??0)>0?"+":""}{match.change} ELO</span></div></article>}):<div className="empty-state">No finished matches yet. Submitted scores will appear here.</div>}</div><p className="demo-note">Match history uses local demo data for now.</p></section>;
}

export default function Home() {
  const tab: string = "Overview";
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
  const winCount = playedMatches.filter((match) => match.won).length;
  const winRate = playedMatches.length ? Math.round(winCount / playedMatches.length * 100) : 0;
  const lastCompleted = playedMatches[0];
  const lastChange = lastCompleted?.change ?? 0;

  return <main className="app-shell">
    <AppSidebar />
    <AppSocialRail />
    <section className="main-area"><AppHeader title={tab}><button className="icon-button" aria-label="Notifications">♧</button><div className="avatar small">A</div></AppHeader>
      <div className={`content ${tab === "Leaderboard" ? "leaderboard-view" : tab === "Matches" ? "matches-view" : ""}`}>
        <section className="home-hero" aria-labelledby="home-hero-title">
          <div className="home-hero-copy"><span className="home-hero-kicker"><span>🏸</span> BADMINTON · YOUR COMMUNITY</span><h1 id="home-hero-title">Own every<br/>rally.</h1><p>Get on court, log your matches, and watch your ranking grow.</p><Link className="primary-button" href="/matches/new"><span>＋</span> Create a match</Link></div>
          <span className="home-hero-art" aria-hidden="true"/><span className="home-hero-shuttle" aria-hidden="true"/>
          <div className="home-hero-stamp"><span>PLAY</span><strong>YOUR</strong><span>GAME</span></div>
        </section>

        {tab === "Leaderboard" && <LeaderboardPage />}
        {tab === "Matches" && <MatchesHistory matches={matches} />}
        <div className="vertical-section-intro"><span>YOUR SEASON</span><h2>At a glance</h2></div>
        <div className="stat-grid"><article className="stat-card rating-card"><div className="stat-top"><span>YOUR ELO RATING</span><span className="live-dot">● DEMO</span></div><div className="rating-value">{currentPlayer.elo.toLocaleString()}<span>pts</span></div><div className="stat-foot"><span className="up">{lastChange > 0 ? "↗ +" : lastChange < 0 ? "↘ " : ""}{lastChange} pts</span><span>last completed match</span><span className="rating-badge">Demo rating</span></div><div className="sparkline"><i/><i/><i/><i/><i/><i/><i/><i/><i/><i/><i/><i/></div></article>
          <article className="stat-card"><div className="stat-top">MATCHES PLAYED <span className="stat-icon">⌁</span></div><div className="stat-number">{playedMatches.length}<span className="muted">{playedMatches.length === 1 ? " match" : " matches"}</span></div><div className="stat-foot"><span className="neutral">This season</span></div><div className="mini-bars">{[24,38,31,54,43,66,50,72,59,84,68,95].map((h,i)=><i key={i} style={{height:`${h}%`}}/>)}</div></article>
          <article className="stat-card"><div className="stat-top">WIN RATE <span className="stat-icon">◎</span></div><div className="stat-number">{winRate}<span className="percent">%</span></div><div className="stat-foot"><span className="neutral">{winCount} {winCount === 1 ? "win" : "wins"}</span><span>from {playedMatches.length} finished</span></div><div className="progress"><i style={{ width: `${winRate}%` }}/></div></article></div>

        <div className="section-heading"><div><h2>Recent matches</h2><p>Your latest games and rating changes</p></div><Link className="text-button" href="/matches">View all <span>→</span></Link></div>
        <div className="matches-table"><div className="table-head"><span>OPPONENT</span><span>GAME TYPE</span><span>RESULT</span><span>RATING CHANGE</span><span>DATE</span></div>{matches.map((match) => <div className="match-row" key={match.id}><div className="player-cell"><div className="overview-avatar-stack">{match.opponent.split(" & ").map((name) => { const player = DEMO_PLAYERS.find((candidate) => candidate.name === name); return <ProfilePortrait key={name} className="overview-profile-avatar" variant={player?.avatar ?? 0} name={name}/>; })}</div><div><strong>{match.opponent}</strong><small>ELO {match.opponentRating}</small></div></div><span className="format-pill">{match.sport} · {match.format}</span><span className={`result ${match.status === "upcoming" ? "pending" : match.won ? "win" : "loss"}`}>{match.status === "upcoming" ? "Awaiting score" : `${match.won ? "Won" : "Lost"} ${match.score?.[0]}–${match.score?.[1]}`}</span><span className={match.change && match.change < 0 ? "down" : "up"}>{match.status === "upcoming" ? "—" : `${(match.change ?? 0) > 0 ? "+" : ""}${match.change} pts`}</span><span className="date-cell">{match.date}</span></div>)}</div>
        <footer>© 2026 Ngadu <span>Play your game. Own your rank.</span></footer>
      </div>
    </section>
  </main>;
}
