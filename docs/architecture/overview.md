# Architecture Overview — Todo List App v2

## Scope and shape

Todo List App v2 is a fullstack single-page todo application. It ships:

- A Next.js frontend for the blue-and-white todo UI.
- A Go HTTP backend for health checks and future todo API handlers.
- PostgreSQL for durable todo persistence.

The app has no login. All visitors share the same deployment-level todo list. Detailed table design and API contracts are intentionally left to the ERD and service-design tasks; this document fixes the foundation that later work must follow.

## Tech stack

### Frontend

- Next.js 15 App Router.
- TypeScript with strict checking.
- Tailwind CSS v3 for styling.
- ESLint using Next.js core web vitals and TypeScript rules.
- Location: `code/frontend/`.

### Backend

- Go 1.22+ module in `code/backend/`.
- `net/http` server with one `main` package under `cmd/api`.
- PostgreSQL driver: `github.com/jackc/pgx/v5/stdlib`.
- Embedded SQL migrations applied on server startup before readiness.

### Database

- PostgreSQL 16 in local Docker Compose and managed PostgreSQL in deployed runtime.
- Backend reads a single `DATABASE_URL`; it does not assemble credentials from separate `DB_*` variables.

## Repository layout

```text
code/
  backend/
    cmd/api/main.go          # HTTP entrypoint; exactly one main package
    internal/db/migrations.go # embedded migration runner
    migrations/              # timestamped .up.sql/.down.sql pairs
    .env.example             # backend env documentation
    go.mod / go.sum          # Go module files
  frontend/
    app/                     # Next.js App Router routes and global CSS
    next.config.js           # standalone output for container runtime
    package.json             # scripts and pinned dependencies
    tailwind.config.ts       # design tokens mapped into Tailwind
    .env.example             # frontend env documentation
docs/architecture/overview.md
.github/workflows/ci.yml
```

## Runtime data flow

1. A browser loads the Next.js single page from the frontend service.
2. Interactive UI code will call the backend using `NEXT_PUBLIC_API_URL`.
3. The backend validates requests, executes parameterized PostgreSQL queries, and returns JSON.
4. PostgreSQL stores todos durably for the deployment.

For this scaffold, the frontend renders a static shell and the backend exposes `/healthz`. Product API handlers and UI data wiring will be added by later stories.

## Startup and readiness

The backend startup sequence is fixed:

1. Read `DATABASE_URL`; fail fast if it is missing.
2. Connect to PostgreSQL.
3. Apply every pending migration from `code/backend/migrations/` in filename order.
4. Start listening on `PORT`, falling back to `APP_PORT`, then `8080`.
5. Return `200` from `/healthz` only after migrations succeeded and `SELECT 1` succeeds.

This prevents a container from reporting healthy while connected to an empty or unreachable database.

## Environment variables

### Root / Docker Compose

- `POSTGRES_USER` — local PostgreSQL user for the Compose database.
- `POSTGRES_PASSWORD` — local PostgreSQL password for the Compose database.
- `POSTGRES_DB` — local PostgreSQL database name.
- `BACKEND_PORT` — optional host port for the backend, default `8080`.
- `FRONTEND_PORT` — optional host port for the frontend, default `3000`.
- `NEXT_PUBLIC_API_URL` — browser-visible backend URL, default `http://localhost:8080`.

### Backend (`code/backend/.env.example`)

- `DATABASE_URL` — PostgreSQL connection URL injected by runtime/Compose.
- `PORT` — HTTP listen port.
- `APP_PORT` — fallback HTTP listen port for platforms that use this name.

### Frontend (`code/frontend/.env.example`)

- `NEXT_PUBLIC_API_URL` — browser-visible backend URL.

## Naming conventions

- Go packages use lowercase names and stay under `internal/` unless they are the executable in `cmd/api`.
- SQL migration files use `YYYYMMDDHHMMSS_description.up.sql` and matching `.down.sql` names.
- React component files use default exported function declarations: `export default function ComponentName()`.
- App Router pages remain Server Components unless client behavior is required. Files using browser APIs, event handlers, or React client hooks must start with the literal first line `"use client"`.
- Tailwind tokens should reference the approved design colors rather than ad-hoc hex values in components.

## Security and reliability conventions

- No secrets are committed. `.env.example` documents keys only.
- Backend code must use parameterized SQL queries.
- User input is validated at the backend boundary; todo titles are trimmed and constrained by the service contract once implemented.
- External errors returned to users are non-technical; detailed operational context belongs in logs.
- Database migrations are idempotent at the runner level through `schema_migrations`.

## CI and local checks

`.github/workflows/ci.yml` runs on pull requests and pushes to `main`:

- Backend: `go build ./...`, `go vet ./...`, `go test ./...`.
- Frontend: `npm ci`, `npm run lint`, `npm run build`, `npm test --if-present`.
- Compose: `docker compose config -q`.

The existing container workflows and Docker Compose conventions are treated as compatibility constraints.

## Running locally

From the repository root:

```bash
cp .env.example .env
cp code/backend/.env.example code/backend/.env
cp code/frontend/.env.example code/frontend/.env.local
docker compose --profile local up --build
```

Then open `http://localhost:3000`. The backend health endpoint is `http://localhost:8080/healthz`.

## Decisions and tradeoffs

### Decision: fullstack shape with PostgreSQL

Tasks must persist across refreshes and later sessions, so the product needs a database-backed backend instead of static local-only state.

Rejected alternative: frontend-only local storage. It would be simpler and cheaper, but it would not meet the SRS requirement that todos persist in a deployment database.

### Decision: Go backend with self-applied migrations

The runtime provisions an empty PostgreSQL database. Applying migrations on backend boot gives every environment the schema before the first query.

Rejected alternative: manual migration command. It creates a hidden deployment step that CI and local `docker compose up` can miss, causing a healthy-looking app to fail on first use.

### Decision: Next.js App Router with standalone output

The committed frontend container expects `code/frontend/` and `.next/standalone/server.js`; aligning the scaffold keeps Docker builds predictable.

Rejected alternative: Vite or a custom Node server. Either would be lighter for a single page, but it would contradict the repository container conventions and require unrelated infrastructure changes.

### Decision: one additional CI workflow

A separate `ci.yml` covers lint, unit tests, type checks, and Compose validation while leaving the orchestrator-owned container and publish workflows untouched.

Rejected alternative: editing `container.yml`. That risks breaking the fixed workflow shared by the pipeline and violates the repository convention for this task.

## Known follow-up work

- ERD task will define the todo table and indexes.
- Service design task will define the todo API routes, request/response payloads, and validation rules.
- Feature stories will replace the scaffold UI with the approved interactive experience and connect it to the backend.
