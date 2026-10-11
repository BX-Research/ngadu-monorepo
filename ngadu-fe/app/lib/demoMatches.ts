export type DemoPlayer = { id: string; name: string; elo: number; avatar: number };
export type MatchFormat = "Singles" | "Doubles";
export type DemoMatch = {
  id: string;
  format: MatchFormat;
  sport: "Badminton";
  createdAt: string;
  players: { you: DemoPlayer; ally?: DemoPlayer; opponents: DemoPlayer[] };
  status: "awaiting_result" | "completed";
  score?: [number, number];
  won?: boolean;
  ratingChange?: number;
  resultMessage?: string;
};

export const DEMO_PLAYERS: DemoPlayer[] = [
  { id: "alex", name: "Alex Morgan", elo: 1248, avatar: 0 },
  { id: "rizky", name: "Rizky Pratama", elo: 1320, avatar: 1 },
  { id: "kevin", name: "Kevin Hartono", elo: 1284, avatar: 2 },
  { id: "dimas", name: "Dimas Putra", elo: 1212, avatar: 3 },
  { id: "nadia", name: "Nadia Putri", elo: 1368, avatar: 4 },
  { id: "sarah", name: "Sarah Wijaya", elo: 1196, avatar: 5 },
];

const STORAGE_KEY = "ngadu.demo.matches.v1";
const CHANGE_EVENT = "ngadu:demo-matches";

export function demoMatchesSnapshot(): string {
  return typeof window === "undefined" ? "[]" : window.localStorage.getItem(STORAGE_KEY) ?? "[]";
}

export function parseDemoMatchesSnapshot(snapshot: string): DemoMatch[] {
  try {
    const value: unknown = JSON.parse(snapshot);
    return Array.isArray(value) ? (value as DemoMatch[]) : [];
  } catch {
    return [];
  }
}

export function subscribeDemoMatches(callback: () => void) {
  if (typeof window === "undefined") return () => {};
  window.addEventListener(CHANGE_EVENT, callback);
  window.addEventListener("storage", callback);
  return () => {
    window.removeEventListener(CHANGE_EVENT, callback);
    window.removeEventListener("storage", callback);
  };
}

export function readDemoMatches(): DemoMatch[] {
  return parseDemoMatchesSnapshot(demoMatchesSnapshot());
}

function writeDemoMatches(matches: DemoMatch[]) {
  window.localStorage.setItem(STORAGE_KEY, JSON.stringify(matches));
  window.dispatchEvent(new Event(CHANGE_EVENT));
}

export function createDemoMatch(input: Omit<DemoMatch, "id" | "createdAt" | "status">): DemoMatch {
  const match: DemoMatch = {
    ...input,
    id: `demo-${Date.now()}-${Math.random().toString(36).slice(2, 7)}`,
    createdAt: new Date().toISOString(),
    status: "awaiting_result",
  };
  writeDemoMatches([match, ...readDemoMatches()]);
  return match;
}

export function findDemoMatch(id: string): DemoMatch | undefined {
  return readDemoMatches().find((match) => match.id === id);
}

export function submitDemoMatchResult(id: string, ownScore: number, opponentScore: number): DemoMatch | undefined {
  if (!Number.isInteger(ownScore) || !Number.isInteger(opponentScore) || ownScore < 0 || opponentScore < 0 || ownScore === opponentScore) return undefined;
  let updated: DemoMatch | undefined;
  const messages = ownScore > opponentScore
    ? ["Ciyeee, menang! Rally kamu lagi panas!", "Mantap, kamu menang! Gas pertahankan streak-nya!", "Wih, lawan tumbang. Good game!"]
    : ["Yah, cupu dikit—rematch, yuk!", "Belum rezeki menang. Latihan, lalu balas!", "Hari ini lawan unggul. Next match giliranmu!"];
  const matches = readDemoMatches().map((match) => {
    if (match.id !== id || match.status !== "awaiting_result") return match;
    const won = ownScore > opponentScore;
    updated = {
      ...match,
      status: "completed",
      score: [ownScore, opponentScore],
      won,
      ratingChange: won ? 18 : -14,
      resultMessage: messages[Math.floor(Math.random() * messages.length)],
    };
    return updated;
  });
  if (updated) writeDemoMatches(matches);
  return updated;
}
