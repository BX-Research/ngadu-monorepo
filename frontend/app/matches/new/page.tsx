"use client";

import Link from "next/link";
import { useState } from "react";
import { useRouter } from "next/navigation";
import AppFrame from "../../components/AppFrame";
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
    return <label className="field">
      <span className="field-label">{label}</span>
      <select value={value} required onChange={(event) => onChange(event.target.value)}>
        <option value="" disabled>Choose a player</option>
        {choices(excluded).map((player) => <option value={player.id} key={player.id}>{player.name}</option>)}
      </select>
      <small>{selected ? `Rating ${selected.elo.toLocaleString("en-US")}` : "Registered players only"}</small>
    </label>;
  }

  function submit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (!selectedOpponent1 || (format === "Doubles" && (!selectedAlly || !selectedOpponent2))) return;
    const opponents = format === "Doubles" ? [selectedOpponent1, selectedOpponent2!] : [selectedOpponent1];
    const match = createDemoMatch({ format, sport: "Badminton", players: { you: signedInPlayer, ...(format === "Doubles" && selectedAlly ? { ally: selectedAlly } : {}), opponents } });
    router.push(`/matches/${match.id}`);
  }

  return (
    <AppFrame>
      <h1 className="page-title">Record a result</h1>
      <p className="page-lede">Pick who played. You’ll enter the score on the next screen.</p>

      <form className="sheet" onSubmit={submit}>
        <div className="format-switch" role="group" aria-label="Match format">
          <button type="button" aria-pressed={format === "Singles"} onClick={() => { setFormat("Singles"); setAllyId(""); setOpponent2(""); }}>Singles <small>1 v 1</small></button>
          <button type="button" aria-pressed={format === "Doubles"} onClick={() => setFormat("Doubles")}>Doubles <small>2 v 2</small></button>
        </div>

        <fieldset className="team">
          <legend>Your side</legend>
          <div className="locked-player">
            <strong>{signedInPlayer.name}</strong>
            <small>You, rating {signedInPlayer.elo.toLocaleString("en-US")}</small>
          </div>
          {format === "Doubles" && playerPicker("Partner", allyId, setAllyId, [opponent1, opponent2], selectedAlly)}
        </fieldset>

        <fieldset className="team">
          <legend>Opponents</legend>
          {playerPicker(format === "Doubles" ? "Opponent 1" : "Opponent", opponent1, setOpponent1, [allyId, ...(format === "Doubles" ? [opponent2] : [])], selectedOpponent1)}
          {format === "Doubles" && playerPicker("Opponent 2", opponent2, setOpponent2, [allyId, opponent1], selectedOpponent2)}
        </fieldset>

        <button className="primary-button" type="submit">Continue to score</button>
      </form>
      <p className="demo-note">Demo only. The match is saved in this browser. <Link className="text-link" href="/matches">Back to results</Link></p>
    </AppFrame>
  );
}
