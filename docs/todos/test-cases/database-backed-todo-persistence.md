# Test Cases — Database-backed Todo Persistence

Module: `todos`
Function: Database-backed todo persistence
Requirements covered: TODOS-005, TODOS-006
Risk level: Medium — this function controls durable storage for create, read, update, and delete operations. The app has no roles or payments, but persistence failures directly affect user trust.

## Automated test cases

**Scenario**: Load saved active and completed todos
**Given**: The database contains two non-deleted todos for the deployment: `Buy milk` saved as active and `Pay bills` saved as completed
**When**: The User opens the app
**Then**: The todo list shows `Buy milk` as active and `Pay bills` as completed

Traceability: TODOS-005 AC-1

**Scenario**: Added todo remains after refresh
**Given**: The app is open and the database is reachable
**When**: The User adds `Book dentist` and refreshes the page after the add action succeeds
**Then**: The todo list shows `Book dentist` as an active todo after refresh

Traceability: TODOS-005 AC-2; TODOS-006 AC-1

**Scenario**: Completed todo remains completed after refresh
**Given**: The database contains an active todo labelled `Water plants` and the app shows it as active
**When**: The User marks `Water plants` complete and refreshes the page after the status change succeeds
**Then**: The todo list shows `Water plants` as completed after refresh

Traceability: TODOS-005 AC-3; TODOS-006 AC-2

**Scenario**: Uncompleted todo remains active after refresh
**Given**: The database contains a completed todo labelled `Submit report` and the app shows it as completed
**When**: The User marks `Submit report` incomplete and refreshes the page after the status change succeeds
**Then**: The todo list shows `Submit report` as active after refresh

Traceability: TODOS-006 AC-3

**Scenario**: Deleted todo remains absent after refresh
**Given**: The database contains a visible todo labelled `Archive notes`
**When**: The User deletes `Archive notes` and refreshes the page after the delete action succeeds
**Then**: The todo list does not show `Archive notes` after refresh

Traceability: TODOS-005 AC-4; TODOS-006 AC-4

**Scenario**: Empty saved list shows empty state
**Given**: The database contains no non-deleted todos for the deployment
**When**: The User opens the app
**Then**: The todo list area shows the friendly empty state explaining that the User can add a first task

Traceability: TODOS-005 AC-5

**Scenario**: Load failure shows retryable non-technical error
**Given**: Saved todos cannot be loaded from the database
**When**: The User opens the app
**Then**: The todo list area shows a clear non-technical error state with a retry path

Traceability: TODOS-005 AC-6

**Scenario**: Failed add preserves current input text
**Given**: The app is open, the todo input contains `Call Sam`, and saving a new todo cannot complete
**When**: The User adds the todo and the save error is shown
**Then**: The todo input still contains `Call Sam`

Traceability: TODOS-006 AC-5

**Scenario**: Failed status change returns to last saved state
**Given**: The database contains an active todo labelled `Feed cat`, the app shows it as active, and saving a status change cannot complete
**When**: The User marks `Feed cat` complete and the save error is shown
**Then**: The visible todo `Feed cat` is shown as active, matching its last saved state

Traceability: TODOS-006 AC-6

**Scenario**: Failed delete restores last saved list state
**Given**: The database contains a visible todo labelled `Renew license`, the app shows it, and saving a delete cannot complete
**When**: The User deletes `Renew license` and the save error is shown
**Then**: The visible list shows `Renew license`, matching its last saved state

Traceability: TODOS-006 AC-6

## Manual test cases

None. These behaviours are observable through UI or integration tests by controlling database contents and simulated service responses.
