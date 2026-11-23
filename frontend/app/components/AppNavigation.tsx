"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";

const items = [
  { label: "Overview", href: "/", icon: <path d="M12 2.6 2.5 10.2V21h7v-6.5h5V21h7V10.2L12 2.6Z" fill="currentColor" /> },
  { label: "Leaderboard", href: "/leaderboard", icon: <path d="M3 14h4.5v7H3v-7Zm6.75-5h4.5v12h-4.5V9ZM16.5 3H21v18h-4.5V3Z" fill="currentColor" /> },
  { label: "Results", href: "/matches", icon: <path d="M3 4.5h3v3H3v-3Zm5 .5h13v2H8V5Zm-5 5.5h3v3H3v-3Zm5 .5h13v2H8v-2Zm-5 5.5h3v3H3v-3Zm5 .5h13v2H8v-2Z" fill="currentColor" /> },
];

export default function AppNavigation() {
  const pathname = usePathname();

  return (
    <nav className="app-navigation" aria-label="Main navigation">
      {items.map((item) => {
        const active = pathname === item.href || (item.href !== "/" && pathname.startsWith(`${item.href}/`));
        return (
          <Link key={item.href} href={item.href} className={`nav-item${active ? " active" : ""}`} aria-current={active ? "page" : undefined}>
            <svg viewBox="0 0 24 24" aria-hidden="true">{item.icon}</svg>
            {item.label}
          </Link>
        );
      })}
    </nav>
  );
}
