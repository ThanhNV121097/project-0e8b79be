# Test Cases — Single-page todo experience

Module: `todos`
Function: Single-page todo experience
Risk level: Medium — this is the primary user-facing data-entry flow for the app, covering core interactions, visual states, and responsive behavior. These cases focus on approved happy-path behavior and requirement-explicit validation/state checks.

## Automated test cases

**Scenario**: Open the single todo page
**Given**: The User has access to the deployed app
**When**: The User opens the app URL
**Then**: The page shows the app heading, add-todo input, add action, and todo list area without requiring navigation

Traceability: TODOS-001 AC-1

**Scenario**: Render approved blue-and-white visual style
**Given**: The User opens the page
**When**: The page renders
**Then**: The visible styling uses the approved blue-and-white palette with white card surfaces

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

**Scenario**: Show immediate minimal interaction feedback
**Given**: The User interacts with buttons or todo status controls
**When**: The interaction occurs
**Then**: Feedback is visible immediately and uses minimal motion without blocking the next action

Traceability: TODOS-001 AC-5

**Scenario**: Show friendly empty state
**Given**: No todos exist
**When**: The User opens the app
**Then**: The todo list area shows a friendly empty state explaining that the User can add a first task

Traceability: TODOS-001 empty data behavior

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

**Scenario**: Validate blank todo input
**Given**: The todo input contains only spaces
**When**: The User adds it
**Then**: No todo is added and an inline message says the task cannot be blank

Traceability: TODOS-002 AC-3

**Scenario**: Clear input after successful add
**Given**: A valid todo is added successfully
**When**: The add action completes
**Then**: The input is cleared

Traceability: TODOS-002 AC-4

**Scenario**: Complete an active todo
**Given**: An active todo is visible
**When**: The User marks it complete
**Then**: The todo appears completed

Traceability: TODOS-003 AC-1

**Scenario**: Uncomplete a completed todo
**Given**: A completed todo is visible
**When**: The User marks it incomplete
**Then**: The todo appears active

Traceability: TODOS-003 AC-2

**Scenario**: Persist completed status after refresh
**Given**: A status change is saved successfully
**When**: The User refreshes the page
**Then**: The todo appears with the saved status

Traceability: TODOS-003 AC-3

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

## Manual verification

No manual-only cases are required for this function. Visual styling, responsive layout, minimal motion, empty state, validation, add, complete/uncomplete, delete, and refresh outcomes are observable through automated UI tests.
