# Story — Single-page todo experience

Module: `todos`
Plan item: Single-page todo experience

## User story

As a User, I want to manage todos from one responsive page, so that I can quickly add tasks, mark them complete or incomplete, and delete tasks without navigating or signing in.

## In scope

- Render the single-page todo app shell using the approved blue-and-white design direction.
- Show the app heading, supporting copy, add-todo input, add action, todo list area, and status/state area on one page.
- Let the User add a todo by entering a title and submitting the add action.
- Trim leading and trailing whitespace before validating and displaying a new todo.
- Reject blank todo titles after trimming with an inline validation message.
- Reject todo titles longer than 200 characters after trimming with an inline validation message.
- Let the User mark each visible todo complete and incomplete.
- Let the User delete each visible todo.
- Show active, completed, loading, empty, validation-error, and save-error states in the UI.
- Keep interactions visually immediate, with minimal motion for hover, completion, and removal feedback.
- Keep the page responsive from 320px mobile width through desktop layouts without horizontal page scroll.
- Keep controls keyboard reachable with visible focus states and accessible labels.
- Use mock or local client state only for this story until the persistence story wires the backend API.

## Out of scope

- Database schema, migrations, API routes, backend handlers, and durable persistence; those belong to `Database-backed todo persistence`.
- Login, accounts, per-user lists, roles, sharing, and permissions.
- Due dates, priorities, labels, reminders, recurring tasks, comments, search, sorting controls, filters, bulk actions, and undo.
- Offline-first behavior, cross-device conflict resolution, and real-time multi-user synchronization.
- Editing an existing todo title after creation.
- Confirmation dialogs for deletion.
- Navigation to additional app pages.

## UI scope

This story touches the approved `Single-page Todo App` screen and implements the hero plus interactive todo card. The UI must use the approved design system tokens: primary blue `#2563EB`, soft background `#F8FAFC`, white surfaces `#FFFFFF`, completed accent `#10B981`, and delete/error accent `#EF4444`.

Required UI states:

- Default state with at least one active todo and one completed todo represented during development or mock-state review.
- Empty state explaining that the User can add a first task.
- Loading state in the todo list area.
- Inline validation state for blank or over-length todo titles.
- Non-technical save-error state that preserves the User's current input text.
- Mobile layout at 320px width with input, add action, completion controls, and delete controls usable without horizontal scrolling.
- Desktop layout at 1024px and wider with the todo card readable and clearly contained.
- Reduced-motion-compatible interactions that do not require animation to understand state changes.

## Acceptance criteria

1. Given the User opens the app URL, when the page renders, then the app heading, add-todo input, add action, todo list area, and supporting state area are visible on a single page without requiring navigation.
2. Given the page renders, then visible styling follows the approved blue-and-white palette with white card surfaces and a clean, lightly polished Todoist-inspired layout that remains much simpler than Todoist.
3. Given the User enters `Buy milk` and submits the add action, when the add succeeds in client state, then a todo labelled `Buy milk` appears as active and the input is cleared.
4. Given the User enters `  Buy milk  ` and submits the add action, when validation runs, then the todo appears as `Buy milk` without leading or trailing spaces.
5. Given the User submits an empty input or an input containing only spaces, when validation runs, then no todo is added, the input remains available, and an inline message says the task cannot be blank.
6. Given the User submits a title longer than 200 characters after trimming, when validation runs, then no todo is added and an inline message names the 200 character limit.
7. Given an active todo is visible, when the User activates its completion control, then the todo immediately appears completed with completed-state styling.
8. Given a completed todo is visible, when the User activates its completion control, then the todo immediately appears active with active-state styling.
9. Given a todo is visible, when the User activates its delete control and the action succeeds in client state, then the todo is removed from the visible list.
10. Given the last visible todo is deleted, when the list becomes empty, then the empty state is shown.
11. Given todos are loading, when the page is waiting for list data, then the todo list area shows a loading state rather than a blank or broken screen.
12. Given an add, status-change, or delete action fails in the UI layer, when the action completes with an error, then a clear non-technical error message is shown and the visible list or input reflects the last safe state.
13. Given the User uses a desktop viewport at least 1024px wide, when the page renders, then the todo card is readable, centered or clearly contained, and there is no horizontal page scroll.
14. Given the User uses a 320px mobile viewport, when the page renders, then the input, add action, completion controls, and delete controls remain visible and usable without horizontal page scroll.
15. Given the User navigates by keyboard, when focus reaches the input, add action, completion controls, or delete controls, then each control has a visible focus state and a meaningful accessible label.
16. Given the User's environment requests reduced motion, when the User adds, completes, uncompletes, or deletes todos, then state changes remain understandable without relying on animation.

## Dependencies

- Depends on the approved design and `design/design-system.md` for layout, colors, typography, components, focus states, and motion limits.
- Depends on `docs/todos/SRS.md` requirements TODOS-001 through TODOS-004 for page, add, complete/uncomplete, and delete behavior.
- Depends on `docs/architecture/overview.md` for the Next.js App Router, TypeScript, Tailwind v3, and frontend location under `code/frontend/`.
- Does not depend on the persistence story to render and exercise the UI; persistence integration will replace mock or local client state later.
- No external accounts or secrets are required.
