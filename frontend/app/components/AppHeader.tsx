import Link from "next/link";

export default function AppHeader() {
  return (
    <header className="app-header">
      <Link className="brand" href="/" aria-label="Ngadu home">
        <span className="brand-mark" aria-hidden="true">n</span>
        <span>ngadu<span className="brand-dot">.</span></span>
      </Link>
      <span className="avatar" role="img" aria-label="Alex Morgan, demo player">A</span>
    </header>
  );
}
