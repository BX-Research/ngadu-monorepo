# Ngadu

A badminton community matchmaking app. This monorepo currently contains the Next.js starter and functional specification. Authentication, backend endpoints, database/storage integration, and matchmaking have not been implemented.

## Structure

```text
ngadu-monorepo/
  AGENTS.md
  package.json          # Root commands delegating to the frontend
  ngadu-fe/           # Next.js + React + TypeScript, using Bun
  ngadu-be/           # Backend placeholder; language undecided
  docs/
    functional-spec.md
    functional-spec.pdf
```

Services share one Git repository and keep independent toolchains. The frontend manifest and lockfile remain in `ngadu-fe`. No extra monorepo build system is required initially.

## Frontend development

Use Bun 1.4.3. From the repository root:

```bash
bun run setup:frontend
bun run dev
```

Open <http://localhost:3000>. Check or build the frontend with:

```bash
bun run lint
bun run typecheck
bun run build
bun run start
```

Production start requires a successful build. Development, build, and start commands explicitly select Bun as the Next.js runtime.

## Architecture

- Frontend: Next.js, React, TypeScript, and Bun.
- Authentication: Firebase Authentication with Google sign-in.
- Backend: separate service; language and framework are undecided.
- Application database: PostgreSQL.
- Development file storage: SeaweedFS.
- Local orchestration: Docker Compose, to be configured after service requirements are settled.

Firebase manages login identities; PostgreSQL will store application data linked to Firebase UIDs. Firestore is not required. Production hosting and storage are undecided.

## Documentation

Read [the functional specification](docs/functional-spec.md) or [its PDF export](docs/functional-spec.pdf). Confirmed choices are separated from proposed features and open decisions.

Read [AGENTS.md](AGENTS.md) and nested instructions before making changes. Keep credentials, environment files, dependencies, build files, and local service data out of Git.
