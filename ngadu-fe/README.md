# Frontend

Next.js, React, and TypeScript frontend for Ngadu, using Bun. This directory contains the starter application; badminton features and Firebase integration have not been implemented.

Prefer commands in the [root README](../README.md). From this directory:

```bash
bun install --frozen-lockfile
bun --bun run dev
bun run lint
bun run --bun tsc --noEmit
bun --bun run build
bun --bun run start
```

Open <http://localhost:3000>. Production start requires a successful build.

Read [AGENTS.md](AGENTS.md) before changing frontend code. Application scope is in [the functional specification](../docs/functional-spec.md).
