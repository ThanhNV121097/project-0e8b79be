# Database Design (ERD) — Todo List App v2

Engine: PostgreSQL 16
Last updated: 2026-08-11
Source requirements: `docs/todos/SRS.md`, `docs/todos/stories/single-page-todo-experience.md`, reviewed UI mock module from PR #13 (`code/frontend/lib/mock/single-page-todo-experience.ts`)

## 1. Overview

This schema stores the deployment-level todo list for Todo List App v2. The only aggregate root is `todos`: each row represents one task with its title, completion state, and timestamps needed for stable ordering and last-saved state. Login, user identity, private lists, due dates, priorities, labels, comments, reminders, recurrence, search metadata, filters, and undo/restore history are deliberately kept out of the database because the SRS excludes them.

The reviewed Single-page todo UI mock exposes todo status as the string union `"active" | "completed"` and timestamps as camelCase `createdAt` / `updatedAt`. The database keeps the normalized boolean column `is_completed`; the Go API maps `false -> "active"` and `true -> "completed"` at the boundary so the approved UI can replace its mock adapter without reworking components.

## 2. Diagram

```mermaid
erDiagram
    TODOS {
        uuid id PK
        text title
        boolean is_completed
        timestamptz created_at
        timestamptz updated_at
    }
```

Cardinality notation: `||` exactly one, `o|` zero or one, `}o` zero or many, `}|` one or many. There are no inter-table relationships in this schema because the app has no login, accounts, lists, or child entities.

## 3. Entities

### 3.1 `todos`

**Purpose** — Stores each saved todo task shown in the shared deployment-level list. **Traces to** — TODOS-001, TODOS-002, TODOS-003, TODOS-004, TODOS-005, TODOS-006.

| Column | Type | Null | Default | Unique | Description |
|---|---|---|---|---|---|
| `id` | `uuid` | no | `gen_random_uuid()` | PK | Stable surrogate identifier for create response, status changes, and deletion. |
| `title` | `text` | no | none | no | Trimmed todo label displayed to users; duplicate titles are allowed. |
| `is_completed` | `boolean` | no | `false` | no | Completion status in storage; `false` maps to API/UI status `active`, `true` maps to `completed`. |
| `created_at` | `timestamptz` | no | `now()` | no | Creation time used for stable ordering across refreshes and mapped to API/UI `createdAt`. |
| `updated_at` | `timestamptz` | no | `now()` | no | Last successful row update time and mapped to API/UI `updatedAt`. |

**Nullable columns** — none.

**Foreign keys** — none. The SRS explicitly excludes login, accounts, private lists, labels, comments, and other parent/child entities.

**Constraints**

| Name | Definition | Rule enforced |
|---|---|---|
| `ck_todos_title_trimmed_not_blank` | `CHECK (length(btrim(title)) BETWEEN 1 AND 200)` | Todo titles must be non-blank after trimming and must not exceed the 200 character limit from TODOS-002 and TODOS-006. |
| `ck_todos_title_is_trimmed` | `CHECK (title = btrim(title))` | Stored titles are already trimmed so readers do not need to trim every display value. |

**Indexes**

| Name | Columns | Type | Query it serves |
|---|---|---|---|
| `idx_todos_created_at_id` | `created_at`, `id` | btree | List all current todos in stable oldest-first order for TODOS-005; `id` is a deterministic tie-breaker when multiple rows share the same timestamp. |

**Lifecycle** — hard delete. TODOS-004 and TODOS-005 require deleted todos to be absent after refresh, and the SRS excludes undo/restore and audit history. A deleted todo has no child rows and no reporting requirement, so physical deletion keeps reads simple and avoids `deleted_at IS NULL` filters on every list query.

## 4. API field mapping for the Single-page todo experience

The reviewed UI mock module is the frontend contract for this story:

```ts
type TodoDto = {
  id: string;
  title: string;
  status: "active" | "completed";
  createdAt: string;
  updatedAt: string;
};

type TodoListResponse = {
  data: TodoDto[];
  meta: { total: number; active: number; completed: number };
};
```

| Database field | API/UI field | Mapping |
|---|---|---|
| `id` | `id` | UUID rendered as a string. |
| `title` | `title` | Stored trimmed text returned as-is. |
| `is_completed` | `status` | `false` => `active`; `true` => `completed`. |
| `created_at` | `createdAt` | RFC 3339 UTC timestamp string. |
| `updated_at` | `updatedAt` | RFC 3339 UTC timestamp string. |

No additional columns are needed for UI-only counts. `meta.total`, `meta.active`, and `meta.completed` are computed by the API from the returned collection for this small app scope.

## 5. Enumerations

No database enumerations are needed. Completion status is a boolean in storage because the SRS defines exactly two states: active and completed. The API exposes the reviewed UI string union `active | completed` for component compatibility.

## 6. Access patterns

| # | Pattern | Frequency | Index used |
|---|---|---|---|
| 1 | Load all todos ordered by `created_at ASC, id ASC` when the page opens, refreshes, or retries after load failure. | High: every page load/refresh/retry. | `idx_todos_created_at_id` |
| 2 | Insert one todo with trimmed title, `is_completed = false`, and generated timestamps. | Medium: each add action. | Primary key only; no secondary index is needed for insert. |
| 3 | Update one todo's `is_completed` and `updated_at` by `id`. | Medium: each complete/uncomplete action. | Primary key index on `id`. |
| 4 | Delete one todo by `id`. | Medium: each delete action. | Primary key index on `id`. |

## 7. Data volume and growth

| Table | Rows at launch | Growth | Retention |
|---|---|---|---|
| `todos` | 0 | Expected to remain small; SRS boundary explicitly verifies 100 visible todos, with no requirement suggesting high-volume ingestion. | Rows are retained until the user deletes them; deletion is hard delete. |

No table is expected to exceed 10M rows within a year for this simple shared-list app. No partitioning or archival strategy is needed for the approved scope.

## 8. Integrity, privacy, and security

- Database-enforced invariants: every todo has a surrogate UUID primary key, a non-null trimmed title of 1 to 200 characters, a non-null boolean completion state, and non-null creation/update timestamps.
- Application-enforced invariants: request bodies are validated at the HTTP boundary before writes so users receive clear non-technical validation errors instead of raw database constraint errors; the database constraints remain the final protection against invalid rows.
- Personal data: only todo titles and completion state are stored. Todo titles may contain user-entered personal content, but the app stores no login identity, profile, email, IP-derived user table, or account identifier.
- Secrets: no table stores secrets, credentials, tokens, password hashes, or private keys.
- Row-level access: none. The SRS defines one shared deployment-level list accessible to any visitor without login.
- Concurrency: last successful write wins for status changes and deletes, matching TODOS-003 and TODOS-006. Updating or deleting a missing `id` returns a not-found application response and does not require extra schema state.

## 9. Migrations

| # | Change | Forward | Backward | Safe on non-empty table |
|---|---|---|---|---|
| 1 | Enable UUID generation | `CREATE EXTENSION IF NOT EXISTS pgcrypto;` | `DROP EXTENSION IF EXISTS pgcrypto;` only if no remaining database object depends on it | Safe. Idempotent extension creation is safe before application tables exist. The down migration may be skipped or guarded in shared databases if other objects use `pgcrypto`. |
| 2 | Initial `todos` table | Create `todos` with `id uuid PRIMARY KEY DEFAULT gen_random_uuid()`, `title text NOT NULL`, `is_completed boolean NOT NULL DEFAULT false`, `created_at timestamptz NOT NULL DEFAULT now()`, `updated_at timestamptz NOT NULL DEFAULT now()`, and the two title check constraints. | `DROP TABLE IF EXISTS todos;` | Safe for an empty database. On a populated table, the backward migration is destructive because it deletes todo data; only run after backup or in disposable environments. |
| 3 | Stable list-order index | `CREATE INDEX idx_todos_created_at_id ON todos (created_at, id);` | `DROP INDEX IF EXISTS idx_todos_created_at_id;` | Safe on an empty table. On a populated table, prefer `CREATE INDEX CONCURRENTLY` outside a transaction to avoid blocking writes. |

Initial migration filenames should follow the architecture convention, for example `20260811000100_create_todos.up.sql` and `20260811000100_create_todos.down.sql`. This task records the schema design only; Dev will implement migrations in the backend story.

## 10. Open questions

| Question | Owner | Blocking |
|---|---|---|
| none | n/a | no |
