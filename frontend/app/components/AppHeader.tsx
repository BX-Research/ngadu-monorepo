import Link from "next/link";

export default function AppHeader({ neutral = false }: { neutral?: boolean }) {
  return (
    <header className="app-header">
      <Link className="brand" href="/" aria-label="Ngadu home">
        <span className="brand-mark" aria-hidden="true">n</span>
        <span>ngadu<span className="brand-dot">.</span></span>
      </Link>
      <span className="avatar" role="img" aria-label={neutral ? "Guest profile" : "Alex Morgan, demo player"}>{neutral ? "○" : "A"}</span>
    </header>
  );
}
