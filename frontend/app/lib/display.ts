import type { DemoMatch } from "./demoMatches";

export function shortDate(iso: string) {
  return new Date(iso).toLocaleDateString("en-US", { month: "short", day: "numeric" });
}

export function opponentNames(match: DemoMatch) {
  return match.players.opponents.map((player) => player.name).join(" & ");
}

export function teamNames(match: DemoMatch) {
  return match.players.ally ? `You & ${match.players.ally.name}` : "You";
}

export function signedPoints(change: number) {
  return `${change > 0 ? "+" : change < 0 ? "−" : ""}${Math.abs(change)} pts`;
}
