import { existsSync, mkdtempSync, rmSync } from "node:fs";
import { tmpdir } from "node:os";
import { join } from "node:path";
import { fileURLToPath } from "node:url";

const root = fileURLToPath(new URL("../", import.meta.url));

export function developmentEnvironment(input = process.env) {
  const env = { ...input };
  for (const [key, fallback] of Object.entries({
    DEV_DB_USER: "ngadu", DEV_DB_PASSWORD: "change-me-local-only",
    DEV_DB_NAME: "ngadu_dev", DEV_DB_PORT: "5432", API_PORT: "8080", FRONTEND_PORT: "3000",
  })) env[key] ||= fallback;
  for (const key of ["DEV_DB_PORT", "API_PORT", "FRONTEND_PORT"]) {
    if (!/^\d+$/.test(env[key]) || Number(env[key]) < 1 || Number(env[key]) > 65535) {
      throw new Error(`${key} must be a port between 1 and 65535`);
    }
  }
  return env;
}

export function backendEnvironment(env) {
  const url = new URL(`postgresql://127.0.0.1:${env.DEV_DB_PORT}/`);
  url.username = env.DEV_DB_USER;
  url.password = env.DEV_DB_PASSWORD;
  url.pathname = "/" + encodeURIComponent(env.DEV_DB_NAME);
  url.searchParams.set("sslmode", "disable");
  return { ...env, DATABASE_URL: url.href, API_ADDR: `127.0.0.1:${env.API_PORT}` };
}

function signalTree(child, signal) {
  if (process.platform === "win32") {
    Bun.spawnSync(["taskkill", "/pid", String(child.pid), "/T", "/F"], { stdout: "ignore", stderr: "ignore" });
    child.kill(signal);
  } else {
    try { process.kill(-child.pid, signal); }
    catch (error) { if (error.code !== "ESRCH") throw error; }
  }
}

export async function stopChildren(children) {
  for (const child of children) signalTree(child, "SIGTERM");
  const exited = Promise.all([...children].map(child => child.exited));
  let timer;
  try {
    const timeout = new Promise(resolve => { timer = setTimeout(() => resolve(false), 5000); });
    if (!await Promise.race([exited.then(() => true), timeout])) {
      for (const child of children) signalTree(child, "SIGKILL");
      await exited;
    }
  } finally {
    clearTimeout(timer);
  }
}

async function main(mode) {
  const env = developmentEnvironment();
  const docker = env.DOCKER_BIN || "docker";
  const go = env.GO_BIN || "go";
  const bun = env.BUN_BIN || process.execPath;
  const children = new Set();
  let directory;
  let interrupted = 0;
  function start(command, environment = env) {
    if (interrupted) throw new Error("Development startup interrupted");
    const child = Bun.spawn(command, {
      cwd: root, env: environment, stdin: "inherit", stdout: "inherit", stderr: "inherit",
      detached: process.platform !== "win32",
    });
    children.add(child);
    return child;
  }
  async function run(command, environment = env) {
    const child = start(command, environment);
    const code = await child.exited;
    children.delete(child);
    if (code !== 0 || interrupted) throw new Error(`${command[0]} command failed (exit ${code})`);
  }
  const onInterrupt = () => {
    interrupted = 130;
    for (const child of children) signalTree(child, "SIGTERM");
  };
  const onTerminate = () => {
    interrupted = 143;
    for (const child of children) signalTree(child, "SIGTERM");
  };
  process.on("SIGINT", onInterrupt);
  process.on("SIGTERM", onTerminate);
  try {
    if (["db-stop", "db-logs"].includes(mode)) {
      await run([docker, "compose", ...(mode === "db-stop" ? ["stop", "db"] : ["logs", "--follow", "db"])]);
      return;
    }
    if (!["dev", "seed", "db-up"].includes(mode)) throw new Error("Use dev, seed, db-up, db-stop, or db-logs");
    if (mode !== "db-up" && !Bun.which(go)) throw new Error("Go was not found. Install Go and add it to PATH.");
    const nextCLI = join(root, "frontend/node_modules/next/dist/bin/next");
    if (mode === "dev") {
      if (!existsSync(nextCLI)) throw new Error("Frontend dependencies are missing. Run make install first.");
      directory = mkdtempSync(join(tmpdir(), "ngadu-api-"));
      await run([go, "-C", "backend", "build", "-o", join(directory, process.platform === "win32" ? "api.exe" : "api"), "./cmd/api"]);
    }
    await run([docker, "compose", "up", "--detach", "--wait", "--wait-timeout", "60", "db"]);
    if (mode === "db-up") return;
    if (mode === "seed") {
      await run([go, "-C", "backend", "run", "./scripts/seed"], backendEnvironment(env));
      return;
    }
    const backend = start([join(directory, process.platform === "win32" ? "api.exe" : "api")], backendEnvironment(env));
    const readyURL = `http://127.0.0.1:${env.API_PORT}/readyz`;
    const deadline = Date.now() + 20000;
    let ready = false;
    while (!ready && Date.now() < deadline && !interrupted && backend.exitCode === null) {
      try { ready = (await fetch(readyURL, { signal: AbortSignal.timeout(1000) })).ok; }
      catch { /* The server may still be starting. */ }
      if (!ready) await Bun.sleep(200);
    }
    if (!ready) throw new Error("Backend did not become ready; check the database settings and API port.");
    const frontend = start([bun, nextCLI, "dev", "frontend", "--hostname", "127.0.0.1", "--port", env.FRONTEND_PORT]);
    console.log(`Frontend: http://127.0.0.1:${env.FRONTEND_PORT} | Backend: http://127.0.0.1:${env.API_PORT}`);
    console.log("Ctrl+C stops both app services. PostgreSQL stays running; use make db-stop to stop it.");
    const result = await Promise.race([
      backend.exited.then(code => ({ name: "Backend", code })),
      frontend.exited.then(code => ({ name: "Frontend", code })),
    ]);
    if (!interrupted) throw new Error(`${result.name} exited (code ${result.code}); stopping the other service.`);
  } catch (error) {
    if (!interrupted) console.error(error.message);
    process.exitCode = interrupted || 1;
  } finally {
    if (interrupted) process.exitCode = interrupted;
    await stopChildren(children);
    if (directory) rmSync(directory, { recursive: true, force: true });
    process.off("SIGINT", onInterrupt);
    process.off("SIGTERM", onTerminate);
  }
}

if (import.meta.main) await main(process.argv[2] || "dev");
