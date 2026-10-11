"use client";

import { useMemo, useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import AppSidebar from "./components/AppSidebar";
import { createDemoMatch, DEMO_PLAYERS, type MatchFormat } from "./lib/demoMatches";

type Match = { id: number; sport: string; format: string; opponent: string; ally?: string; opponentRating: number; rating: number; date: string; status: "upcoming" | "completed"; score?: [number, number]; change?: number };

const playerAccounts = [{ id: "p1", name: "Rizky Pratama", elo: 1320 }, { id: "p2", name: "Kevin Hartono", elo: 1284 }, { id: "p3", name: "Dimas Putra", elo: 1212 }, { id: "p4", name: "Nadia Putri", elo: 1368 }, { id: "p5", name: "Sarah Wijaya", elo: 1196 }];

const initialMatches: Match[] = [
  { id: 1, sport: "Badminton", format: "Singles", opponent: "Kevin Hartono", opponentRating: 1284, rating: 1248, date: "Today · 18:30", status: "upcoming" },
  { id: 2, sport: "Badminton", format: "Doubles", opponent: "Rizky Pratama & Dimas Putra", opponentRating: 1266, rating: 1248, date: "Tomorrow · 19:00", status: "upcoming" },
  { id: 3, sport: "Badminton", format: "Singles", opponent: "Dimas Putra", opponentRating: 1212, rating: 1248, date: "Oct 08 · 20:15", status: "completed", score: [21, 17], change: 18 },
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
  const router = useRouter();
  const [tab, setTab] = useState<"Overview" | "Leaderboard" | "Matches">("Overview");
  const [matches] = useState(initialMatches);
  const [format, setFormat] = useState<MatchFormat>("Singles");
  const [allyId, setAllyId] = useState("");
  const [opponent1, setOpponent1] = useState("");
  const [opponent2, setOpponent2] = useState("");
  const [notice, setNotice] = useState("");
  const rating = useMemo(() => 1248 + matches.reduce((n, m) => n + (m.change ?? 0), 0), [matches]);
  const selectedOpponent1 = playerAccounts.find((player) => player.id === opponent1);
  const selectedOpponent2 = playerAccounts.find((player) => player.id === opponent2);
  const selectedAlly = playerAccounts.find((player) => player.id === allyId);



  function createMatch(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (!selectedOpponent1 || (format === "Doubles" && (!selectedOpponent2 || !selectedAlly))) return;
    if (format === "Doubles" && selectedOpponent1.id === selectedOpponent2?.id) { setNotice("Pilih dua akun pemain yang berbeda."); return; }
    const resolvePlayer = (name: string) => DEMO_PLAYERS.find((player) => player.name === name)!;
    const opponents = format === "Doubles" ? [resolvePlayer(selectedOpponent1.name), resolvePlayer(selectedOpponent2!.name)] : [resolvePlayer(selectedOpponent1.name)];
    const match = createDemoMatch({ format, sport: "Badminton", players: { you: DEMO_PLAYERS[0], ...(format === "Doubles" && selectedAlly ? { ally: resolvePlayer(selectedAlly.name) } : {}), opponents } });
    setTab("Matches");
    router.push(`/matches/${match.id}`);
  }

  return <main className="app-shell">
    <AppSidebar />
    <section className="main-area"><header className="topbar"><div className="crumb">Workspace <span>/</span> {tab}</div><div className="top-actions"><button className="icon-button" aria-label="Notifications">♧</button><div className="avatar small">A</div></div></header>
      <div className={`content ${tab === "Leaderboard" ? "leaderboard-view" : tab === "Matches" ? "matches-view" : ""}`}><div className="welcome"><div><p className="eyebrow">SUNDAY, OCTOBER 11</p><h1>{tab === "Overview" ? "Your game, your ranking." : tab}</h1><p className="subhead">Find your level. Play your game. Climb the ranks.</p></div><button className="primary-button" onClick={() => document.getElementById("new-match")?.scrollIntoView({ behavior: "smooth" })}><span>＋</span> Record a match</button></div>
        {notice && <div className="notice" role="status">{notice}<button onClick={() => setNotice("")} aria-label="Dismiss">×</button></div>}
        {tab === "Leaderboard" && <LeaderboardPage />}
        {tab === "Matches" && <MatchesHistory matches={matches} />}
        <div className="stat-grid"><article className="stat-card rating-card"><div className="stat-top"><span>YOUR ELO RATING</span><span className="live-dot">● LIVE</span></div><div className="rating-value">{rating.toLocaleString()}<span>pts</span></div><div className="stat-foot"><span className="up">↗ +18</span><span>from last match</span><span className="rating-badge">Intermediate</span></div><div className="sparkline"><i/><i/><i/><i/><i/><i/><i/><i/><i/><i/><i/><i/></div></article>
          <article className="stat-card"><div className="stat-top">MATCHES PLAYED <span className="stat-icon">⌁</span></div><div className="stat-number">{matches.filter(m => m.status === "completed").length}<span className="muted"> matches</span></div><div className="stat-foot"><span className="neutral">This season</span></div><div className="mini-bars">{[24,38,31,54,43,66,50,72,59,84,68,95].map((h,i)=><i key={i} style={{height:`${h}%`}}/>)}</div></article>
          <article className="stat-card"><div className="stat-top">WIN RATE <span className="stat-icon">◎</span></div><div className="stat-number">67<span className="percent">%</span></div><div className="stat-foot"><span className="up">↗ +4.2%</span><span>vs last month</span></div><div className="progress"><i/></div></article></div>
        <div className="section-heading"><div><h2>{tab === "Leaderboard" ? "Top players" : "Recent matches"}</h2><p>Your latest games and rating changes</p></div><Link className="text-button" href="/matches">View all <span>→</span></Link></div>
        <div className="matches-table"><div className="table-head"><span>OPPONENT</span><span>GAME TYPE</span><span>RESULT</span><span>RATING CHANGE</span><span>DATE</span></div>{matches.map(m=><div className="match-row" key={m.id}><div className="player-cell"><div className="opponent-avatar">{m.opponent.slice(0,1)}</div><div><strong>{m.opponent}</strong><small>{m.status === "upcoming" ? `ELO ${m.opponentRating}` : `ELO ${m.opponentRating}`}</small></div></div><span className="format-pill">{m.sport} · {m.format}</span><span className={`result ${m.status === "upcoming" ? "pending" : "win"}`}>{m.status === "upcoming" ? "Awaiting score" : `Won ${m.score?.[0]}–${m.score?.[1]}`}</span><span className={m.change && m.change < 0 ? "down" : "up"}>{m.status === "upcoming" ? "—" : `${(m.change ?? 0) > 0 ? "+" : ""}${m.change} pts`}</span><span className="date-cell">{m.date}</span></div>)}</div>
        <div className="lower-grid"><section className="panel" id="new-match"><div className="panel-title"><div><h2>Set up a match</h2><p>Choose how you want to play</p></div><span className="panel-icon">＋</span></div><div className="sport-label">SPORT</div><div className="sport-choice"><span className="sport-emoji">🏸</span><span><strong>Badminton</strong><small>Sport selected</small></span><span className="check">✓</span></div><form onSubmit={createMatch}><div className="sport-label">GAME FORMAT</div><div className="format-switch"><button type="button" className={format === "Singles" ? "selected" : ""} onClick={() => setFormat("Singles")}>Singles <small>1 vs 1</small></button><button type="button" className={format === "Doubles" ? "selected" : ""} onClick={() => setFormat("Doubles")}>Doubles <small>2 vs 2</small></button></div><div className="doubles-team">{format === "Doubles" && <div className="our-team"><div className="locked-player"><span>YOUR PLAYER 1</span><strong>Alex Morgan <small>Signed-in account</small></strong><b>1,248 ELO</b></div><label className="player-picker"><span className="sport-label">YOUR ALLY</span><select required value={allyId} onChange={(event) => setAllyId(event.target.value)}><option value="" disabled>Select your teammate…</option>{playerAccounts.filter((player) => player.id !== opponent1 && player.id !== opponent2).map((player) => <option key={player.id} value={player.id}>{player.name} · ELO {player.elo}</option>)}</select><small>{selectedAlly ? `Current rating: ${selectedAlly.elo} ELO` : "Choose your teammate's account"}</small></label></div>}<div className="player-pickers">{[1, ...(format === "Doubles" ? [2] : [])].map((slot) => { const chosen = slot === 1 ? opponent1 : opponent2; const setChosen = slot === 1 ? setOpponent1 : setOpponent2; const chosenPlayer = playerAccounts.find((player) => player.id === chosen); return <label className="player-picker" key={slot}><span className="sport-label">OPPONENT PLAYER {slot}</span><select required value={chosen} onChange={(event) => setChosen(event.target.value)}><option value="" disabled>Select registered player…</option>{playerAccounts.filter((player) => player.id !== allyId && (slot === 1 || player.id !== opponent1)).map((player) => <option key={player.id} value={player.id}>{player.name} · ELO {player.elo}</option>)}</select><small>{chosenPlayer ? `Current rating: ${chosenPlayer.elo} ELO` : "Choose an account to see their ELO"}</small></label>; })}</div></div><button className="secondary-button" type="submit">Create match <span>→</span></button></form></section>
          <section className="panel score-panel"><div className="panel-title"><div><h2>Every match needs a score</h2><p>Submit the final result after you play</p></div><span className="panel-icon">↗</span></div><div className="empty-state">Create a match and Ngadu will take you to its score page straight away.</div><Link className="primary-button full" href="/matches/new">Create a match <span>→</span></Link></section></div>
        <footer>© 2026 Ngadu <span>Play your game. Own your rank.</span></footer>
      </div>
    </section>
  </main>;
}
