"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";

const items = [
  { label: "Overview", href: "/", icon: "◫" },
  { label: "Leaderboard", href: "/leaderboard", icon: "♜" },
  { label: "Matches", href: "/matches", icon: "⌁" },
];

export default function AppSidebar() {
  const pathname = usePathname();
  return <aside className="sidebar">
    <Link className="brand" href="/"><span className="brand-mark">n</span> ngadu<span className="brand-dot">.</span></Link>
    <p className="side-label">MENU</p>
    {items.map((item) => <Link key={item.href} href={item.href} className={`nav-item ${pathname === item.href || (item.href !== "/" && pathname.startsWith(item.href)) ? "active" : ""}`}><span className="nav-icon">{item.icon}</span>{item.label}</Link>)}
    <div className="sidebar-bottom"><div className="avatar">A</div><div><strong>Alex Morgan</strong><small>Demo player account</small></div><span className="more">•••</span></div>
  </aside>;
}
