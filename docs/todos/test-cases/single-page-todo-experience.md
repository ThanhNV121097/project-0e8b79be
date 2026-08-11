# Test Cases — Single-page todo experience

Module: `todos`
Function: Single-page todo experience
Risk level: Medium — this is the core user-facing flow for the app and includes data-changing interactions, responsive layout, and required visual/accessibility states. This artifact covers the requested happy-path scenarios and the explicitly required blank-task validation scenario.

## Automated test cases

**Scenario**: Open the single-page todo app
**Given**: The User has access to the deployed app
**When**: The User opens the app URL
**Then**: The page shows the app heading, add-todo input, add action, and todo list area without requiring navigation
Traceability: TODOS-001 AC-1

**Scenario**: Render approved blue-and-white visual direction
**Given**: The User opens the page
**When**: The page renders
**Then**: The visible styling uses the approved blue-and-white palette, including primary blue emphasis and white card surfaces on a soft background
Traceability: TODOS-001 AC-2

**Scenario**: Render usable desktop layout
**Given**: The User uses a desktop viewport at least 1024px wide
**When**: The page renders
**Then**: The todo card is readable, centered or clearly contained, and the page has no horizontal scroll
Traceability: TODOS-001 AC-3

**Scenario**: Render usable mobile layout
**Given**: The User uses a mobile viewport 320px wide
**When**: The page renders
**Then**: The input, add action, and todo controls remain visible and usable without horizontal page scroll
Traceability: TODOS-001 AC-4

**Scenario**: Provide immediate minimal interaction feedback
**Given**: The User can see buttons or todo status controls
**When**: The User interacts with a button or todo status control
**Then**: The control feedback appears immediately and uses minimal motion without blocking the next action
Traceability: TODOS-001 AC-5

**Scenario**: Add a new active todo
**Given**: The todo input is empty
**When**: The User types `Buy milk` and adds it
**Then**: A todo labelled `Buy milk` appears in the list as active
Traceability: TODOS-002 AC-1

**Scenario**: Trim whitespace when adding a todo
**Given**: The todo input contains `  Buy milk  `
**When**: The User adds it
**Then**: A todo labelled `Buy milk` appears without leading or trailing spaces
Traceability: TODOS-002 AC-2

**Scenario**: Show inline validation for blank todo
**Given**: The todo input contains only spaces
**When**: The User adds it
**Then**: No todo is added and an inline message says the task cannot be blank
Traceability: TODOS-002 AC-3

**Scenario**: Clear input after successful add
**Given**: A valid todo is added successfully
**When**: The add action completes
**Then**: The input is cleared
Traceability: TODOS-002 AC-4

**Scenario**: Preserve input after failed add save
**Given**: Saving a valid todo fails
**When**: The add action completes with an error
**Then**: The input still contains the User's current text
Traceability: TODOS-002 AC-5

**Scenario**: Mark an active todo complete
**Given**: An active todo is visible
**When**: The User marks it complete
**Then**: The todo appears completed
Traceability: TODOS-003 AC-1

**Scenario**: Mark a completed todo active
**Given**: A completed todo is visible
**When**: The User marks it incomplete
**Then**: The todo appears active
Traceability: TODOS-003 AC-2

**Scenario**: Keep saved completed status after refresh
**Given**: A status change is saved successfully
**When**: The User refreshes the page
**Then**: The todo appears with the saved status
Traceability: TODOS-003 AC-3

**Scenario**: Restore last saved status after failed status save
**Given**: A status change fails to save
**When**: The action completes with an error
**Then**: The todo returns to its last saved status and a non-technical error is shown
Traceability: TODOS-003 AC-4

**Scenario**: Delete a visible todo
**Given**: A todo is visible
**When**: The User deletes it
**Then**: The todo is removed from the list
Traceability: TODOS-004 AC-1

**Scenario**: Show empty state after deleting the last todo
**Given**: The last remaining todo is visible
**When**: The User deletes it
**Then**: The empty state is shown
Traceability: TODOS-004 AC-2

**Scenario**: Keep deleted todo absent after refresh
**Given**: A delete action is saved successfully
**When**: The User refreshes the page
**Then**: The deleted todo does not reappear
Traceability: TODOS-004 AC-3

**Scenario**: Keep or restore todo after failed delete save
**Given**: A delete action fails to save
**When**: The action completes with an error
**Then**: The todo remains visible or is restored, and a non-technical error is shown
Traceability: TODOS-004 AC-4

## Manual checks

No manual-only checks are required for these scenarios. Visual, responsive, interaction, validation, status, delete, and refresh outcomes are observable through browser-based automated tests.
