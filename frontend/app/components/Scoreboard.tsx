import type { ReactNode } from "react";

/** The court net that sits between the two scores. */
function Net() {
  return (
    <svg className="score-net" viewBox="0 0 120 64" fill="none" aria-hidden="true" preserveAspectRatio="none">
      <defs>
        <pattern id="net-mesh" width="7" height="7" patternUnits="userSpaceOnUse">
          <path d="M7 0H0v7" stroke="currentColor" strokeWidth=".8" opacity=".55" />
        </pattern>
        <clipPath id="net-shape"><path d="M6 17Q60 24 114 17V50Q60 56 6 50Z" /></clipPath>
      </defs>
      <rect x="0" y="0" width="120" height="64" fill="url(#net-mesh)" clipPath="url(#net-shape)" />
      <path d="M6 17Q60 24 114 17" stroke="currentColor" strokeWidth="3" strokeLinecap="round" />
      <path d="M6 50Q60 56 114 50" stroke="currentColor" strokeWidth="1" opacity=".6" />
      <path d="M5 12V58M115 12V58" stroke="currentColor" strokeWidth="2.4" strokeLinecap="round" />
    </svg>
  );
}

type Side = { label: ReactNode; score: ReactNode };

type Props = {
  title: ReactNode;
  date?: ReactNode;
  left: Side;
  right: Side;
  outcome?: "won" | "lost";
  children?: ReactNode;
};

export default function Scoreboard({ title, date, left, right, outcome, children }: Props) {
  return (
    <section className="scoreboard">
      <header className="scoreboard-head">
        <h2>{title}</h2>
        {date && <span>{date}</span>}
      </header>
      <div className="scoreboard-court">
        <div className="score-figure">{left.score}</div>
        <Net />
        <div className="score-figure score-figure-right">{right.score}</div>
      </div>
      <div className="score-names">
        <div className="score-name">{left.label}</div>
        <div className="score-name score-name-right">{right.label}</div>
      </div>
      {outcome && <span className={`outcome-tag outcome-${outcome}`}>{outcome === "won" ? "Won" : "Lost"}</span>}
      {children}
    </section>
  );
}
