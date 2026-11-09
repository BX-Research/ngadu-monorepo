import { Suspense } from "react";
import PlayerProfileClient from "./PlayerProfileClient";

type PlayerProfilePageProps = {
  params: Promise<{ userId: string }>;
};

export default function PlayerProfilePage({ params }: PlayerProfilePageProps) {
  return (
    <Suspense fallback={<main className="app-shell"><section className="main-area"><div className="content player-profile-route" role="status">Loading player profile…</div></section></main>}>
      <PlayerProfilePageContent params={params} />
    </Suspense>
  );
}

async function PlayerProfilePageContent({ params }: PlayerProfilePageProps) {
  const { userId } = await params;
  return <PlayerProfileClient userId={userId} />;
}
