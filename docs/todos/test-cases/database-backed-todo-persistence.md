# Test Cases — Database-backed Todo Persistence

Module: `todos`
Function: Database-backed todo persistence
Risk level: Medium. This function writes and reads user task data durably through a database; loss of saved state would break the product's core value. Cases focus on the approved acceptance criteria for persisted create, read, update, delete, reload, empty, and explicit error-recovery behaviour.

## Automated test cases

**Scenario**: Load saved active and completed todos
**Given**: The database contains two non-deleted todos for the deployment: `Buy milk` saved as active and `Pay rent` saved as completed, with `Buy milk` created before `Pay rent`.
**When**: The User opens the app.
**Then**: The todo list shows `Buy milk` as active before `Pay rent`, and shows `Pay rent` as completed.

Traceability: TODOS-005 AC-1; TODOS-005 behaviour 1, 2, 3.

**Scenario**: Added todo remains after page refresh
**Given**: The app is open and the database is reachable.
**When**: The User enters `Buy milk`, adds the todo, waits for the add action to complete successfully, and refreshes the page.
**Then**: A todo labelled `Buy milk` appears in the list as active after the refresh.

Traceability: TODOS-005 AC-2; TODOS-006 AC-1; TODOS-006 behaviour 1.

**Scenario**: Completed status remains after page refresh
**Given**: A saved active todo labelled `Buy milk` is visible in the list.
**When**: The User marks `Buy milk` complete, waits for the status change to complete successfully, and refreshes the page.
**Then**: The todo labelled `Buy milk` appears completed after the refresh.

Traceability: TODOS-005 AC-3; TODOS-006 AC-2; TODOS-006 behaviour 2.

**Scenario**: Uncompleted status remains after page refresh
**Given**: A saved completed todo labelled `Buy milk` is visible in the list.
**When**: The User marks `Buy milk` incomplete, waits for the status change to complete successfully, and refreshes the page.
**Then**: The todo labelled `Buy milk` appears active after the refresh.

Traceability: TODOS-006 AC-3; TODOS-006 behaviour 2.

**Scenario**: Deleted todo remains absent after page refresh
**Given**: A saved todo labelled `Buy milk` is visible in the list.
**When**: The User deletes `Buy milk`, waits for the delete action to complete successfully, and refreshes the page.
**Then**: No todo labelled `Buy milk` appears in the list after the refresh.

Traceability: TODOS-005 AC-4; TODOS-006 AC-4; TODOS-006 behaviour 3.

**Scenario**: Empty database shows empty state
**Given**: The database contains no non-deleted todos for the deployment.
**When**: The User opens the app.
**Then**: The todo list area shows the friendly empty state explaining that the User can add a first task.

Traceability: TODOS-005 AC-5; TODOS-005 behaviour 4.

**Scenario**: Failed load shows retryable non-technical error state
**Given**: Saved todos cannot be loaded from the database.
**When**: The User opens the app.
**Then**: The todo list area shows a non-technical error state and includes a retry path the User can activate.

Traceability: TODOS-005 AC-6; TODOS-005 behaviour 5; Reliability NFR.

**Scenario**: Failed add preserves current input text
**Given**: The app is open and saving a new todo cannot complete.
**When**: The User enters `Buy milk` and activates the add action.
**Then**: A non-technical error is shown, and the add-todo input still contains `Buy milk`.

Traceability: TODOS-006 AC-5; TODOS-006 behaviour 4, 5; Reliability NFR.

**Scenario**: Failed status change returns to last saved state
**Given**: A saved active todo labelled `Buy milk` is visible in the list and saving status changes cannot complete.
**When**: The User marks `Buy milk` complete.
**Then**: A non-technical error is shown, and `Buy milk` is displayed as active, matching its last saved status.

Traceability: TODOS-006 AC-6; TODOS-006 behaviour 4; Reliability NFR.

**Scenario**: Failed delete keeps last saved todo visible
**Given**: A saved todo labelled `Buy milk` is visible in the list and saving deletes cannot complete.
**When**: The User deletes `Buy milk`.
**Then**: A non-technical error is shown, and `Buy milk` remains visible or is restored in the list.

Traceability: TODOS-006 AC-6; TODOS-006 behaviour 4; Reliability NFR.

## Manual checks

None. The required persistence, reload, status, deletion, empty-state, and explicit error-recovery behaviours are observable through automated UI/integration tests with controllable database and service failure fixtures.
