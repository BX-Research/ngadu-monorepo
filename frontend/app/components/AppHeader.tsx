import Link from "next/link";

function Shuttlecock() {
  return (
    <svg className="brand-shuttle" viewBox="0 0 32 32" fill="none" aria-hidden="true">
      <g transform="rotate(40 16 16)">
        <path d="M10.5 4h11l-2.6 15h-5.8L10.5 4Z" fill="currentColor" opacity=".22" />
        <path d="M10.5 4l2.6 15M16 3.5V19M21.5 4l-2.6 15M11.4 9.5h9.2M12.3 14.5h7.4" stroke="currentColor" strokeWidth="2" strokeLinecap="round" />
        <path d="M12.6 19h6.8v2.2a3.4 3.4 0 0 1-6.8 0V19Z" fill="currentColor" />
      </g>
    </svg>
  );
}

export default function AppHeader({ neutral = false }: { neutral?: boolean }) {
  return (
    <header className="app-header">
      <Link className="brand" href="/" aria-label="Ngadu home">
        <span>ngadu.</span>
        <Shuttlecock />
      </Link>
      <span className="avatar" role="img" aria-label={neutral ? "Guest profile" : "Alex Morgan, demo player"}>{neutral ? "?" : "A"}</span>
    </header>
  );
}
