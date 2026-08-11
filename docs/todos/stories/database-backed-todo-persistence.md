# Story — Database-backed todo persistence

Module: `todos`
Plan item: Database-backed todo persistence

## User story

As a User, I want todo changes to be saved in the deployment database, so that my tasks remain available after refreshes and later sessions on the same deployment.

## Requirements covered

- TODOS-005 — Load persisted todos
- TODOS-006 — Save todo changes durably

## In scope

- Load all non-deleted todos from the database when the app opens or refreshes.
- Preserve each todo title, completion status, stable identifier, created time, and updated time in PostgreSQL.
- Save create, complete, uncomplete, and delete operations through the backend API before treating them as durable.
- Keep list order stable across reloads, with older created todos shown before newer created todos.
- Show the approved loading, empty, validation, and non-technical error states when persistence operations are pending or fail.
- Preserve the User's current add-todo input text when creating a todo fails.
- Restore or retain the last saved visible todo state when completion or deletion fails.
- Support the shared deployment-level todo list; all visitors see the same saved list because there is no login.

## Out of scope

- Login, accounts, per-user private todo lists, roles, and permissions beyond public access.
- Due dates, priorities, labels, reminders, recurring todos, search, manual sorting, and bulk actions.
- Offline-first storage, local-only persistence, background sync, conflict-resolution UI, and undo after deletion.
- Pagination or virtualized rendering beyond the SRS requirement that 100 todos render without broken controls.
- Showing created or updated timestamps in the UI.
- Changing the approved visual design, palette, page structure, or component inventory.

## UI scope

This story touches the approved single-page todo app only. The hero and interactive todo card must remain visually aligned with `design/design-system.md` and use the existing blue-and-white palette, todo rows, completion controls, delete controls, input field, loading state, empty state, validation message, and sync/error notice patterns.

The UI must add real persisted data behavior to the same screen: initial load shows a loading state, successful load shows saved todos or the empty state, failed load shows a non-technical error with a retry path, and save failures show a clear non-technical message without discarding the current input or falsely displaying unsaved state as durable.

## Acceptance criteria

1. Given saved active and completed todos exist in the database, when the User opens the app, then the list displays those todos with their saved titles and completion statuses.
2. Given no non-deleted todos exist in the database, when the User opens the app, then the approved empty state is shown.
3. Given todos are being fetched, when the app is waiting for the backend response, then the todo list area shows the approved loading state rather than a blank or broken screen.
4. Given loading saved todos fails, when the app receives the failure, then a clear non-technical error state is shown with a retry path and the page does not crash.
5. Given the User adds a valid todo and the database save succeeds, when the User refreshes the page, then the added todo still appears as active.
6. Given the User enters a valid todo and the database save fails, when the error is shown, then the input still contains the entered text and no durable saved state is implied.
7. Given the User completes an active todo and the database save succeeds, when the User refreshes the page, then the todo appears completed.
8. Given the User uncompletes a completed todo and the database save succeeds, when the User refreshes the page, then the todo appears active.
9. Given changing a todo completion status fails to save, when the error is shown, then the visible todo returns to its last saved completion status.
10. Given the User deletes a todo and the database save succeeds, when the User refreshes the page, then the deleted todo does not reappear.
11. Given deleting a todo fails to save, when the error is shown, then the todo remains visible or is restored in the list.
12. Given multiple saved todos exist, when the User refreshes the page, then older created todos appear before newer created todos.
13. Given the database contains 100 non-deleted todos, when the app loads them, then the list renders without horizontal page scroll or broken todo controls.
14. Given any User opens or changes the list, when the backend handles the request, then the action is allowed without login and affects the shared deployment-level list.
15. Given a save targets a todo that no longer exists, when the backend reports the missing todo, then the UI refreshes or removes the missing todo and shows a non-technical message without crashing.
16. Given two visitors save changes to the same todo close together, when the list is refreshed, then the most recent successful save is shown as the source of truth.

## Dependencies

- The fullstack scaffold from `docs/architecture/overview.md` must exist: Next.js frontend, Go backend, PostgreSQL, Docker Compose, and `DATABASE_URL`.
- Technical design must define the todo table, indexes, migrations, API routes, request and response payloads, and validation rules before backend implementation.
- The single-page todo experience story should provide the approved UI structure and component behavior that this story wires to persisted data.
- No external accounts or credentials are required beyond the deployment PostgreSQL connection supplied through `DATABASE_URL`.

## Notes for downstream teams

- Backend validation must trim todo titles and enforce the 1 to 200 character rule from the SRS.
- Backend queries must be parameterized.
- User-facing persistence errors must be non-technical; detailed operational context belongs in logs.
- Frontend must call the backend through `NEXT_PUBLIC_API_URL` as defined in the architecture overview.
- Do not add stakeholder questions for this story; the SRS and architecture provide enough direction to proceed.
