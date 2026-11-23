import type { ReactNode } from "react";
import AppHeader from "./AppHeader";
import AppNavigation from "./AppNavigation";

export default function AppFrame({ children, neutral = false }: { children: ReactNode; neutral?: boolean }) {
  return (
    <main className="app-shell">
      <AppHeader neutral={neutral} />
      <div className="content">{children}</div>
      <AppNavigation />
    </main>
  );
}
