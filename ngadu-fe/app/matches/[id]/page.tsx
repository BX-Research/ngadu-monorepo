"use client";

import Link from "next/link";
import { Suspense, useState, useSyncExternalStore } from "react";
import { useParams } from "next/navigation";
import AppSidebar from "../../components/AppSidebar";
import AppSocialRail from "../../components/AppSocialRail";
import AppHeader from "../../components/AppHeader";
import ProfilePortrait from "../../components/ProfilePortrait";
import { demoMatchesSnapshot, parseDemoMatchesSnapshot, subscribeDemoMatches, submitDemoMatchResult } from "../../lib/demoMatches";

export default function MatchDetailRoute() {
  return <Suspense fallback={<main className="app-shell"><section className="main-area"><div className="content match-detail-route">Loading match…</div></section></main>}><MatchDetailContent/></Suspense>;
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
      setError(ownScore === opponentScore ? "Skor seri belum bisa dikirim. Tentukan pemenang dengan skor lebih tinggi." : "Masukkan dua skor bulat yang valid.");
      return;
    }
    setError("");
  }

  const ownPlayers = match ? [match.players.you, ...(match.players.ally ? [match.players.ally] : [])] : [];
  const opponents = match?.players.opponents ?? [];
  const createdAt = match ? new Date(match.createdAt).toLocaleString("id-ID", { dateStyle: "medium", timeStyle: "short" }) : "";

  return <main className="app-shell"><AppSidebar/><AppSocialRail/><section className="main-area"><AppHeader title="Match detail"><span className="leader-season">🏸 {match?.format ?? "Match"}</span></AppHeader><div className="content match-detail-route">
    {!loaded ? <section className="new-match-card"><p>Loading match…</p></section> : !match ? <section className="new-match-card"><h1>Match tidak ditemukan</h1><p>Match demo ini mungkin dibuat di browser atau perangkat lain.</p><Link className="primary-button" href="/matches">Kembali ke matches</Link></section> : <>
      <section className="new-match-hero"><div><span className="hero-kicker">🏸 &nbsp;{match.status === "completed" ? "MATCH COMPLETED" : "MATCH IN PROGRESS"}</span><h1>{match.status === "completed" ? "Good game!" : "Enter your score."}</h1><p>{match.sport} · {match.format} · dibuat {createdAt}</p></div><span className="match-detail-art">🏸<b>✨</b></span></section>
      <section className="match-detail-card"><div className="new-match-card-head"><div><span className="eyebrow">MATCH RESULT</span><h2>Who played?</h2><p>Pastikan tim dan skor yang dimasukkan sesuai hasil permainan.</p></div><span className="setup-icon">🏆</span></div>
        <div className="match-teams"><article className="match-team-card own-team-card"><span className="new-section-label">YOUR TEAM</span>{ownPlayers.map((player) => <div className="match-player-line" key={player.id}><ProfilePortrait variant={player.avatar} name={player.name}/><span><strong>{player.name}{player.id === match.players.you.id ? <small>YOU</small> : null}</strong><small>{player.elo.toLocaleString()} ELO</small></span></div>)}</article><div className="match-versus">VS</div><article className="match-team-card opponent-team-card"><span className="new-section-label">OPPONENT TEAM</span>{opponents.map((player) => <div className="match-player-line" key={player.id}><ProfilePortrait variant={player.avatar} name={player.name}/><span><strong>{player.name}</strong><small>{player.elo.toLocaleString()} ELO</small></span></div>)}</article></div>
        {match.status === "awaiting_result" ? <form className="score-submit-form" onSubmit={submit}><h2>Final score</h2><p>Skor lebih tinggi menentukan pemenang. Seri belum bisa disubmit.</p><div className="score-submit-row"><label><span>YOUR TEAM</span><input aria-label="Your team score" type="number" min="0" step="1" required value={scores[0]} onChange={(event) => setScores([event.target.value, scores[1]])} placeholder="0"/></label><b>:</b><label><span>OPPONENT TEAM</span><input aria-label="Opponent team score" type="number" min="0" step="1" required value={scores[1]} onChange={(event) => setScores([scores[0], event.target.value])} placeholder="0"/></label></div>{error && <div className="notice" role="alert">{error}</div>}<button className="primary-button full" type="submit">Submit score <span>→</span></button><small className="score-required-note">Score submission is required to complete this match.</small></form> : <section className={`match-result-banner ${match.won ? "result-won" : "result-lost"}`}><span className="result-sparkle">{match.won ? "🎉" : "💪"}</span><div><h2>{match.resultMessage}</h2><p>Skor akhir: {match.score?.[0]} – {match.score?.[1]} · {match.won ? "Menang" : "Kalah"} · demo rating {match.ratingChange && match.ratingChange > 0 ? "+" : ""}{match.ratingChange} ELO</p><small>Hasil tersimpan di browser ini. Rating masih simulasi, belum diperbarui di server.</small></div></section>}
      </section><div className="match-detail-actions"><Link className="text-button" href="/matches">Lihat riwayat matches →</Link><Link className="text-button" href="/matches/new">Buat match lagi ＋</Link></div>
    </>}
  </div></section></main>;
}
