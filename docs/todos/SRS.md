# SRS — Todos

Module: `todos`
Last updated: 2026-08-11
Design: [View Design](http://localhost:8080/design/0e8b79be-50d7-46ff-9a1a-33ebb7dc3349)
Design system: `design/design-system.md`

> One file per module, at `docs/{module}/SRS.md`. It covers only the functions
> that belong to this module. Never write `docs/SRS.md`.

## 1. Purpose

The `todos` module provides the complete single-page experience for “Todo List App v2”. A visitor can add tasks, see the current list, mark items complete or incomplete, delete items, and return later to the same saved list on the same deployment. Without this module, the product has no usable task management flow and no persistence of work.

## 2. Actors

| Actor | Who they are | What they may do in this module |
|---|---|---|
| User | Any visitor using the app; there is no login or account identity | View the shared todo list for the deployment, add todos, complete or uncomplete todos, delete todos, and see loading, empty, validation, and error states |

## 3. Scope

**In scope** — the functions specified below, by their plan titles:

- Single-page todo experience
- Database-backed todo persistence

**Out of scope** — related functionality that is deliberately not part of this module:

- Login, accounts, roles, and per-user private lists — deliberately not built; the stakeholder requested no login.
- Due dates, priorities, labels, comments, reminders, recurring tasks, search, and sorting — deliberately not built; this version is a simple add/complete/delete todo app.
- Offline-first editing and cross-device conflict resolution — deliberately not built; the app depends on the deployment being reachable.
- Bulk actions and undo/restore after deletion — deliberately not built; each todo is managed individually.

## 4. Functional requirements

### 4.1 Single-page todo experience

**Requirement TODOS-001 — View todo page**

*As a* User, *I want to* open a single todo page, *so that* I can manage my tasks without navigating between pages.

Behaviour:

1. When the User opens the app, the system displays one page containing the app heading, add-todo input, add action, todo list area, and supporting state area.
2. The page uses the approved blue-and-white visual direction: #2563EB primary blue, #F8FAFC soft background, #FFFFFF surface, #10B981 completed accent, and #EF4444 delete/error.
3. The layout is clean, lightly polished, and loosely inspired by Todoist while remaining much simpler.
4. The page uses minimal motion only for small interaction feedback such as button hover, item completion, or item removal.
5. The page remains usable on desktop and mobile viewport widths.

**Acceptance criteria** — each maps one-to-one onto a test case in `docs/todos/test-cases/single-page-todo-experience.md`.

| # | Given | When | Then |
|---|---|---|---|
| AC-1 | The User has access to the deployed app | The User opens the app URL | The page shows the app heading, input, add action, and todo list area without requiring navigation |
| AC-2 | The User opens the page | The page renders | The visible styling uses the approved blue-and-white palette with white card surfaces |
| AC-3 | The User uses a desktop viewport at least 1024px wide | The page renders | The todo card is readable, centered or clearly contained, and has no horizontal page scroll |
| AC-4 | The User uses a mobile viewport 320px wide | The page renders | The input, add action, and todo controls remain visible and usable without horizontal page scroll |
| AC-5 | The User interacts with buttons or todo status controls | The interaction occurs | Feedback is immediate and uses minimal motion without blocking the next action |

**Failure, boundary and permission behaviour**

| Case | Condition | Expected behaviour |
|---|---|---|
| Empty data | No todos exist | The page shows a friendly empty state that explains the User can add a first task |
| Loading | Todos are being loaded | The page shows a loading state in the todo list area rather than a blank or broken screen |
| Visual boundary | The viewport width is 320px | The layout remains usable and does not require horizontal page scrolling |
| Not permitted | Any User accesses the page | Access is allowed because the app has no login or restricted role |
| Motion sensitivity | The User's environment requests reduced motion | Essential interactions remain usable without relying on animation |

**Data touched**

| Field | Type | Required | Rule |
|---|---|---|---|
| Todo title | text | yes | Displayed as the task label; cannot be blank after trimming whitespace |
| Todo completion status | boolean | yes | Displayed as active or completed |
| Todo created time | timestamp | yes | Used to keep a stable list order; not shown unless the design later adds it |

**Requirement TODOS-002 — Add todo**

*As a* User, *I want to* add a todo from the page, *so that* I can record a task immediately.

Behaviour:

1. The User enters text in the add-todo input and activates the add action.
2. The system trims leading and trailing whitespace before validating the todo title.
3. If the trimmed title is valid, the system adds the todo to the list as active and clears the input.
4. If the trimmed title is blank, the system keeps the input available and shows a clear inline validation message.
5. The User's entered text is not discarded when a save error occurs.

**Acceptance criteria** — each maps one-to-one onto a test case in `docs/todos/test-cases/single-page-todo-experience.md`.

| # | Given | When | Then |
|---|---|---|---|
| AC-1 | The todo input is empty | The User types `Buy milk` and adds it | A todo labelled `Buy milk` appears in the list as active |
| AC-2 | The todo input contains `  Buy milk  ` | The User adds it | A todo labelled `Buy milk` appears without leading or trailing spaces |
| AC-3 | The todo input contains only spaces | The User adds it | No todo is added and an inline message says the task cannot be blank |
| AC-4 | A valid todo is added successfully | The add action completes | The input is cleared |
| AC-5 | Saving a valid todo fails | The add action completes with an error | The input still contains the User's current text |

**Failure, boundary and permission behaviour**

| Case | Condition | Expected behaviour |
|---|---|---|
| Invalid input | Trimmed todo title is blank | Inline validation appears; nothing is saved; the input remains focused or easy to continue editing |
| Boundary | Todo title is 1 character after trimming | The todo is accepted |
| Boundary | Todo title exceeds 200 characters after trimming | The todo is rejected with a message naming the 200 character limit |
| Duplicate text | A todo already has the same title | The new todo is accepted; duplicate titles are allowed |
| Upstream failure | The save cannot complete | A clear non-technical error is shown and the User's current input is preserved |
| Not permitted | Any User adds a todo | The action is allowed because the app has no login or restricted role |

**Data touched**

| Field | Type | Required | Rule |
|---|---|---|---|
| Todo title | text | yes | Trimmed before validation; 1 to 200 characters after trimming |
| Todo completion status | boolean | yes | New todos start as active |
| Todo created time | timestamp | yes | Set when the todo is created |

**Requirement TODOS-003 — Complete and uncomplete todo**

*As a* User, *I want to* mark a todo complete or incomplete, *so that* the list reflects the current state of my tasks.

Behaviour:

1. The User activates the completion control for an active todo.
2. The system marks the todo completed and updates its visual treatment immediately.
3. The User activates the completion control for a completed todo.
4. The system marks the todo active and updates its visual treatment immediately.
5. If the status change cannot be saved, the system shows a clear non-technical error and returns the displayed todo to its last saved status.

**Acceptance criteria** — each maps one-to-one onto a test case in `docs/todos/test-cases/single-page-todo-experience.md`.

| # | Given | When | Then |
|---|---|---|---|
| AC-1 | An active todo is visible | The User marks it complete | The todo appears completed |
| AC-2 | A completed todo is visible | The User marks it incomplete | The todo appears active |
| AC-3 | A status change is saved successfully | The User refreshes the page | The todo appears with the saved status |
| AC-4 | A status change fails to save | The action completes with an error | The todo returns to its last saved status and a non-technical error is shown |

**Failure, boundary and permission behaviour**

| Case | Condition | Expected behaviour |
|---|---|---|
| Not found | The target todo no longer exists | The todo is removed from the visible list or the list is refreshed; a non-technical message explains it is no longer available |
| Conflict | Two Users change the same todo close together | The last saved change is the source of truth after refresh |
| Upstream failure | The status update cannot complete | The UI returns to the last saved status and shows a clear non-technical error |
| Not permitted | Any User changes status | The action is allowed because the app has no login or restricted role |

**Data touched**

| Field | Type | Required | Rule |
|---|---|---|---|
| Todo identifier | identifier | yes | Identifies the todo being changed |
| Todo completion status | boolean | yes | Can switch between active and completed |
| Todo updated time | timestamp | yes | Changes when completion status changes |

**Requirement TODOS-004 — Delete todo**

*As a* User, *I want to* delete a todo, *so that* I can remove tasks I no longer need.

Behaviour:

1. The User activates the delete action for a visible todo.
2. The system removes the todo from the list after the delete action succeeds.
3. If the deleted todo was the last todo, the system shows the empty state.
4. If the delete action cannot be saved, the system keeps or restores the todo in the visible list and shows a clear non-technical error.

**Acceptance criteria** — each maps one-to-one onto a test case in `docs/todos/test-cases/single-page-todo-experience.md`.

| # | Given | When | Then |
|---|---|---|---|
| AC-1 | A todo is visible | The User deletes it | The todo is removed from the list |
| AC-2 | The last remaining todo is visible | The User deletes it | The empty state is shown |
| AC-3 | A delete action is saved successfully | The User refreshes the page | The deleted todo does not reappear |
| AC-4 | A delete action fails to save | The action completes with an error | The todo remains visible or is restored, and a non-technical error is shown |

**Failure, boundary and permission behaviour**

| Case | Condition | Expected behaviour |
|---|---|---|
| Not found | The target todo no longer exists | The visible list refreshes or removes the missing todo; the page does not crash |
| Conflict | Two Users delete or update the same todo close together | A deleted todo stays deleted after refresh if deletion was the last saved change |
| Upstream failure | The delete cannot complete | The todo remains visible or is restored and a clear non-technical error is shown |
| Not permitted | Any User deletes a todo | The action is allowed because the app has no login or restricted role |

**Data touched**

| Field | Type | Required | Rule |
|---|---|---|---|
| Todo identifier | identifier | yes | Identifies the todo being deleted |
| Todo deleted state | boolean or absence | yes | Deleted todos are no longer shown in the active list |

### 4.2 Database-backed todo persistence

**Requirement TODOS-005 — Load persisted todos**

*As a* User, *I want to* see saved todos when I open or refresh the app, *so that* my tasks remain available across sessions on the same deployment.

Behaviour:

1. When the User opens or refreshes the app, the system loads the saved todos for the deployment.
2. The system displays all non-deleted todos with their saved titles and completion statuses.
3. The list order remains stable across refreshes, with older created todos appearing before newer created todos unless the approved design later specifies a different order.
4. If no saved todos exist, the system shows the empty state.
5. If saved todos cannot be loaded, the system shows a clear non-technical error state and offers a retry path.

**Acceptance criteria** — each maps one-to-one onto a test case in `docs/todos/test-cases/database-backed-todo-persistence.md`.

| # | Given | When | Then |
|---|---|---|---|
| AC-1 | Saved active and completed todos exist | The User opens the app | The list shows the saved todos with their saved statuses |
| AC-2 | A todo was added successfully | The User refreshes the page | The todo still appears |
| AC-3 | A todo status was changed successfully | The User refreshes the page | The todo shows the changed status |
| AC-4 | A todo was deleted successfully | The User refreshes the page | The todo does not appear |
| AC-5 | No saved todos exist | The User opens the app | The empty state is shown |
| AC-6 | Saved todos cannot be loaded | The User opens the app | A non-technical error state is shown with a retry path |

**Failure, boundary and permission behaviour**

| Case | Condition | Expected behaviour |
|---|---|---|
| Empty data | The database contains no non-deleted todos | The empty state appears |
| Loading | The list is being fetched | A loading state appears until the list is available or an error occurs |
| Upstream failure | The list cannot be fetched | A clear non-technical error appears and the page remains usable enough to retry |
| Boundary | The list contains 100 todos | The page renders the todos without horizontal scroll or broken controls |
| Not permitted | Any User loads the list | The action is allowed because the app has no login or restricted role |

**Data touched**

| Field | Type | Required | Rule |
|---|---|---|---|
| Todo identifier | identifier | yes | Stable identifier for each saved todo |
| Todo title | text | yes | Displayed as the task label |
| Todo completion status | boolean | yes | Displayed as active or completed |
| Todo created time | timestamp | yes | Used for stable ordering |
| Todo updated time | timestamp | yes | Used to determine the latest saved status if needed |

**Requirement TODOS-006 — Save todo changes durably**

*As a* User, *I want to* have add, status change, and delete actions saved durably, *so that* the list remains correct after reloads and later sessions.

Behaviour:

1. When the User adds a valid todo, the system saves the new todo before treating it as durable.
2. When the User changes a todo status, the system saves the new status and uses that saved status after reload.
3. When the User deletes a todo, the system saves the deletion and excludes the todo from later list loads.
4. If any save operation fails, the system shows a clear non-technical error and does not claim the failed change was saved.
5. The User does not lose unsaved input text when an add operation fails.

**Acceptance criteria** — each maps one-to-one onto a test case in `docs/todos/test-cases/database-backed-todo-persistence.md`.

| # | Given | When | Then |
|---|---|---|---|
| AC-1 | The User adds a valid todo and the save succeeds | The User refreshes the page | The added todo appears after refresh |
| AC-2 | The User completes a todo and the save succeeds | The User refreshes the page | The todo appears completed after refresh |
| AC-3 | The User uncompletes a todo and the save succeeds | The User refreshes the page | The todo appears active after refresh |
| AC-4 | The User deletes a todo and the save succeeds | The User refreshes the page | The todo remains absent after refresh |
| AC-5 | The User adds a valid todo and the save fails | The error is shown | The input still contains the entered todo text |
| AC-6 | The User changes or deletes a todo and the save fails | The error is shown | The visible list reflects the last saved state |

**Failure, boundary and permission behaviour**

| Case | Condition | Expected behaviour |
|---|---|---|
| Upstream failure | Creating a todo cannot be saved | A non-technical error appears; the input text remains available; no durable todo is implied |
| Upstream failure | Updating completion status cannot be saved | A non-technical error appears; the todo returns to the last saved status |
| Upstream failure | Deleting a todo cannot be saved | A non-technical error appears; the todo remains visible or is restored |
| Conflict | Two Users save changes to the same todo close together | The most recent successful save is shown after reload |
| Not found | A save targets a todo that no longer exists | The list refreshes or removes the missing todo and shows a non-technical message |
| Not permitted | Any User saves changes | The action is allowed because the app has no login or restricted role |

**Data touched**

| Field | Type | Required | Rule |
|---|---|---|---|
| Todo identifier | identifier | yes | Required for status changes and deletion |
| Todo title | text | yes for create | 1 to 200 characters after trimming |
| Todo completion status | boolean | yes | Persisted with each todo |
| Todo created time | timestamp | yes | Persists across reloads |
| Todo updated time | timestamp | yes | Changes when a todo is updated |

## 5. Screens

The design is the source of truth for appearance; this section maps functions onto it so nothing in the design is unaccounted for and nothing specified here is missing from the design.

Approved design preview: [View Design](http://localhost:8080/design/0e8b79be-50d7-46ff-9a1a-33ebb7dc3349)

Approved color palette:

- #2563EB — primary blue for key actions and emphasis.
- #F8FAFC — soft page background.
- #FFFFFF — card and input surfaces.
- #10B981 — completed state accent.
- #EF4444 — delete and error accent.

Approved screens and sections:

- Single-page Todo App: hero plus interactive todo card for adding, completing/uncompleting, deleting, filtering, validation, loading, empty, and error states.
- States & Details: supporting sections showing loading/empty/error behavior, responsive behavior, database persistence intent, and minimal motion guidance.

| Screen | Section in the design | Functions it serves | States that must exist |
|---|---|---|---|
| Single-page Todo App | Hero and interactive todo card | TODOS-001, TODOS-002, TODOS-003, TODOS-004, TODOS-005, TODOS-006 | default, loading, empty, validation error, save error, active todo, completed todo |
| States & Details | Loading/empty/error, responsive behavior, persistence, minimal motion guidance | TODOS-001, TODOS-005, TODOS-006 | loading, empty, error, mobile, desktop, reduced motion |

## 6. Non-functional requirements

| Area | Requirement |
|---|---|
| Performance | Initial page render shows the app shell within 2 seconds on a typical broadband connection, excluding a cold deployment start |
| Performance | Add, complete/uncomplete, and delete actions show visible feedback within 200 milliseconds of the User action |
| Performance | Loading 100 todos completes within 2 seconds on a typical broadband connection, excluding a cold deployment start |
| Accessibility | Input, add action, completion controls, and delete controls are keyboard reachable, have visible focus states, and expose accessible labels |
| Accessibility | Text and interactive controls meet contrast ratio of at least 4.5:1 against their backgrounds |
| Responsive | The page works from 320px viewport width upward with no horizontal page scroll |
| Privacy | The app stores todo titles and completion state only; it stores no login identity or profile data |
| Reliability | A failed load, create, update, or delete operation shows a non-technical error message and does not silently discard the User's current input |

## 7. Dependencies and assumptions

- **Depends on:** a database, for persisted todos across refresh and later sessions on the same deployment.
- **Depends on:** the approved design and `design/design-system.md`, for the blue-and-white visual direction, interaction states, and component styling.
- **Assumption:** The todo list is shared at the deployment level because the stakeholder requested no login. If private per-user lists are later required, authentication and account-scoped data become new scope.
- **Assumption:** Duplicate todo titles are allowed because the app is intentionally simple and has no uniqueness requirement.
- **Assumption:** The latest successful save is the source of truth for concurrent edits because the app has no user identity or collaboration model.

| Open question | Proposed default | Who decides |
|---|---|---|
| None | Use the assumptions above unless scope changes | Stakeholder |

## 8. Traceability

Every plan item in this module appears exactly once, and every requirement id traces to a test case. A gap in this table is a gap in the build.

| Plan item | Requirement ids | Test cases |
|---|---|---|
| Single-page todo experience | TODOS-001, TODOS-002, TODOS-003, TODOS-004 | `docs/todos/test-cases/single-page-todo-experience.md` |
| Database-backed todo persistence | TODOS-005, TODOS-006 | `docs/todos/test-cases/database-backed-todo-persistence.md` |
