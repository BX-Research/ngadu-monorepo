import { Suspense } from "react";
import MatchDetailClient from "./MatchDetailClient";

type MatchDetailRouteProps = {
  params: Promise<{ id: string }>;
};

export default function MatchDetailRoute({ params }: MatchDetailRouteProps) {
  return (
    <Suspense fallback={<main className="app-shell"><section className="main-area"><div className="content match-detail-route">Loading match…</div></section></main>}>
      <MatchDetailPageContent params={params} />
    </Suspense>
  );
}

async function MatchDetailPageContent({ params }: MatchDetailRouteProps) {
  const { id } = await params;
  return <MatchDetailClient id={id} />;
}
