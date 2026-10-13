import { describe, expect, test } from "bun:test";
import { developmentEnvironment, backendEnvironment, stopChildren } from "./dev.mjs";

describe("local database settings", () => {
  test("defaults use the same local port and database as Compose", () => {
    const env = backendEnvironment(developmentEnvironment({ DATABASE_URL: "postgres://remote.invalid/production" }));
    const url = new URL(env.DATABASE_URL);
    expect(url.hostname).toBe("127.0.0.1");
    expect(url.port).toBe("5432");
    expect(url.pathname).toBe("/ngadu_dev");
    expect(env.API_ADDR).toBe("127.0.0.1:8080");
  });
  test("overrides and special characters survive connection URL construction", () => {
    const env = backendEnvironment(developmentEnvironment({
      DEV_DB_PORT: "55432", DEV_DB_USER: "user@example", DEV_DB_PASSWORD: "p@ss:/?#",
      DEV_DB_NAME: "dev database", API_PORT: "18080",
    }));
    const url = new URL(env.DATABASE_URL);
    expect(url.port).toBe("55432");
    expect(decodeURIComponent(url.username)).toBe("user@example");
    expect(decodeURIComponent(url.password)).toBe("p@ss:/?#");
    expect(decodeURIComponent(url.pathname)).toBe("/dev database");
    expect(env.API_ADDR).toBe("127.0.0.1:18080");
  });
  test("invalid port settings fail before services start", () => {
    for (const port of ["0", "65536", "abc", "5432; echo bad"]) {
      expect(() => developmentEnvironment({ DEV_DB_PORT: port })).toThrow("DEV_DB_PORT");
    }
  });
});

test("shutdown stops managed processes", async () => {
  const child = Bun.spawn([process.execPath, "-e", "setInterval(() => {}, 1000)"], {
    detached: process.platform !== "win32", stdout: "ignore", stderr: "ignore",
  });
  await stopChildren(new Set([child]));
  expect(await child.exited).not.toBe(0);
});
