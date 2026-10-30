"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import AppNavigation from "../../components/AppNavigation";
import AppHeader from "../../components/AppHeader";
import SportsIllustration from "../../components/SportsIllustration";
import { createDemoMatch, DEMO_PLAYERS, type DemoPlayer, type MatchFormat } from "../../lib/demoMatches";

const signedInPlayer = DEMO_PLAYERS[0];
const accounts = DEMO_PLAYERS.filter((player) => player.id !== signedInPlayer.id);

export default function NewMatchRoute() {
  const router = useRouter();
  const [format, setFormat] = useState<MatchFormat>("Singles");
  const [allyId, setAllyId] = useState("");
  const [opponent1, setOpponent1] = useState("");
  const [opponent2, setOpponent2] = useState("");
  const selectedAlly = accounts.find((player) => player.id === allyId);
  const selectedOpponent1 = accounts.find((player) => player.id === opponent1);
  const selectedOpponent2 = accounts.find((player) => player.id === opponent2);
  const choices = (excluded: string[]) => accounts.filter((player) => !excluded.includes(player.id));

  function playerPicker(label: string, value: string, onChange: (id: string) => void, excluded: string[], selected?: DemoPlayer) {
    return <label className="new-match-field"><span>{label}</span><select value={value} required onChange={(event) => onChange(event.target.value)}><option value="" disabled>Choose a player…</option>{choices(excluded).map((player) => <option value={player.id} key={player.id}>{player.name} · {player.elo} ELO</option>)}</select><small>{selected ? `🏅 ${selected.elo} ELO` : "Choose a registered player account"}</small></label>;
  }

  function submit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (!selectedOpponent1 || (format === "Doubles" && (!selectedAlly || !selectedOpponent2))) return;
    const opponents = format === "Doubles" ? [selectedOpponent1, selectedOpponent2!] : [selectedOpponent1];
    const match = createDemoMatch({ format, sport: "Badminton", players: { you: signedInPlayer, ...(format === "Doubles" && selectedAlly ? { ally: selectedAlly } : {}), opponents } });
    router.push(`/matches/${match.id}`);
  }

  return <main className="app-shell"><AppNavigation/><section className="main-area"><AppHeader/><div className="content new-match-route">
    <section className="new-match-hero"><div><span className="hero-kicker">🏸 &nbsp;READY, SET, PLAY!</span><h1>Set up your<br/>next match.</h1><p>Choose your format and players. You’ll enter the score on the match page right after.</p></div><SportsIllustration/></section>
    <form className="new-match-card" onSubmit={submit}>
      <div className="new-match-card-head"><div><span className="eyebrow">MATCH SETUP</span><h2>Who’s playing?</h2><p>Pick singles or doubles, then choose the registered players.</p></div><span className="setup-icon">🏸</span></div>
      <div className="format-switch"><button type="button" className={format === "Singles" ? "selected" : ""} onClick={() => { setFormat("Singles"); setAllyId(""); setOpponent2(""); }}>Singles <small>1 vs 1</small></button><button type="button" className={format === "Doubles" ? "selected" : ""} onClick={() => setFormat("Doubles")}>Doubles <small>2 vs 2</small></button></div>
      {format === "Doubles" && <div className="new-team-card"><span className="new-section-label">YOUR TEAM</span><div className="locked-player"><span>PLAYER 1 · YOU</span><strong>{signedInPlayer.name} <small>Locked to signed-in demo account</small></strong><b>{signedInPlayer.elo.toLocaleString()} ELO</b></div>{playerPicker("YOUR ALLY", allyId, setAllyId, [opponent1, opponent2], selectedAlly)}</div>}
      <div className="new-team-card opponent-team"><span className="new-section-label">OPPONENT TEAM</span>{playerPicker("OPPONENT PLAYER 1", opponent1, setOpponent1, [allyId, ...(format === "Doubles" ? [opponent2] : [])], selectedOpponent1)}{format === "Doubles" && playerPicker("OPPONENT PLAYER 2", opponent2, setOpponent2, [allyId, opponent1], selectedOpponent2)}</div>
      <button className="primary-button new-match-submit" type="submit">Create match & enter score <span>→</span></button>
      <p className="demo-note">Demo only: match setup is saved in this browser. The next page requires a score submission to finish the match.</p>
    </form>
  </div></section></main>;
}
