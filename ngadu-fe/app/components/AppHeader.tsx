import type { ReactNode } from "react";
import Link from "next/link";

type AppHeaderProps = {
  title: string;
  children?: ReactNode;
};

export default function AppHeader({ title, children }: AppHeaderProps) {
  return (
    <header className="topbar">
      <Link className="topbar-brand" href="/" aria-label="Ngadu home">
        <span className="topbar-brand-mark">n</span>
        <span>ngadu<span className="topbar-brand-dot">.</span></span>
      </Link>
      <div className="crumb">
        <span className="crumb-workspace">Workspace</span>
        <span>/</span>
        <strong>{title}</strong>
      </div>
      <div className="top-actions">{children}</div>
    </header>
  );
}
