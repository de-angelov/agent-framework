# BACKLOG

Pending work only. Active work lives in TASKS.md. Completed work lives in ARCHIVE.md.

## Backlog

### Team Create Dialog

Task ID: UI-DIALOG-01
Category: AFK
Owner: Unassigned
Branch:
Status: Backlog
Dependencies: TEAM-UI-01
Blocking tasks: CLEANUP-02

#### Objective

Move team creation into the shared dialog component.

#### Scope

- Replace the always-visible create-team form on the teams screen with a shared `Dialog` opened by a shared `Button`.
- Keep the existing `create` action intent and team service behavior unchanged.
- Keep team rename and delete row actions unchanged.
- Use the shared `Button` component for dialog trigger, cancel, and submit actions.
- Show existing create validation or success messages without changing backend copy.
- Add focused route/component tests for opening the dialog, submitting create, cancel/close behavior, validation error display, and preserving the teams table.

#### Out of Scope

- Do not change team service rules, team edit/delete behavior, authenticated shell behavior, table internals, or broad styling.
- Do not change epic, ticket, or comment screens.

#### Acceptance Criteria

- Team dialog tests pass.
- Existing team management behavior still works.
- Verification passes.
- No always-visible create-team form remains.

#### Verification

Run from `workspaces/repo-agent-*`: `npm test -- app/routes/teams/teams.test.tsx app/components/dialog/dialog.test.tsx && npm run typecheck && npm run build`

---

### Epic Create Dialog

Task ID: UI-DIALOG-02
Category: AFK
Owner: Unassigned
Branch:
Status: Backlog
Dependencies: EPIC-UI-01
Blocking tasks: CLEANUP-02

#### Objective

Move epic creation into the shared dialog component.

#### Scope

- Replace the always-visible create-epic form on the epics screen with a shared `Dialog` opened by a shared `Button`.
- Keep the existing `create` action intent and epic service behavior unchanged.
- Keep epic edit and delete row actions unchanged.
- Use the shared `Button` component for dialog trigger, cancel, and submit actions.
- Preserve team selection, title, and description fields inside the dialog.
- Show existing create validation or success messages without changing backend copy.
- Add focused route/component tests for opening the dialog, submitting create, cancel/close behavior, validation error display, team selection, and preserving the epics table.

#### Out of Scope

- Do not change epic service rules, epic edit/delete behavior, authenticated shell behavior, table internals, or broad styling.
- Do not change team, ticket, or comment screens.

#### Acceptance Criteria

- Epic dialog tests pass.
- Existing epic management behavior still works.
- Verification passes.
- No always-visible create-epic form remains.

#### Verification

Run from `workspaces/repo-agent-*`: `npm test -- app/routes/epics/epics.test.tsx app/components/dialog/dialog.test.tsx && npm run typecheck && npm run build`

---

### Ticket Create Dialog Entry

Task ID: UI-DIALOG-03
Category: HITL
Owner: Unassigned
Branch:
Status: Backlog
Dependencies: TICKET-05, BOARD-02B
Blocking tasks: TICKET-WF-01, BOARD-WF-01, CLEANUP-02

#### Objective

Use the shared dialog and button components for the primary new-ticket workflow.

#### Scope

- Replace the board's plain create-ticket affordance with a shared `Button` that opens a shared `Dialog`.
- Render the existing create-ticket fields in the dialog: team, epic, type, state, title, and body.
- Submit through the existing ticket create action/service path without duplicating ticket business rules.
- Preserve same-team epic options for the selected team.
- Keep `/tickets/new` available as a direct route if existing navigation or redirects depend on it, but make the board dialog the primary new-ticket entry point.
- Add focused tests for opening the dialog, rendering the create fields, successful create, validation error display, same-team epic options, and preserving board content behind the dialog.

#### Out of Scope

- Do not change ticket create service rules, ticket edit/delete behavior, board columns, filters, drag-and-drop, comments, or broad styling.
- Do not add rich-text editing.

#### Acceptance Criteria

- Ticket create dialog tests pass.
- Existing create-ticket route/action tests continue to pass.
- Verification passes.
- Board no longer uses a plain link as the primary new-ticket affordance.

#### Verification

Run from `workspaces/repo-agent-*`: `npm test -- app/routes/board/board.test.tsx app/routes/tickets/new.test.tsx app/services/tickets/tickets.server.test.ts && npm run typecheck && npm run build`

---

### Comment Add Dialog

Task ID: UI-DIALOG-04
Category: AFK
Owner: Unassigned
Branch:
Status: Backlog
Dependencies: COMMENT-02C
Blocking tasks: TICKET-WF-01, CLEANUP-02

#### Objective

Use the shared dialog and button components for adding ticket comments.

#### Scope

- Replace the always-visible add-comment form on ticket details with a shared `Button` that opens a shared `Dialog`.
- Render the comment body field and submit action inside the dialog.
- Submit through the existing add-comment action/service path.
- Preserve chronological comment display outside the dialog.
- Show validation errors for empty comment bodies without changing comment service rules.
- Add focused tests for opening the dialog, submitting a comment, cancel/close behavior, validation error display, and preserving the comments list.

#### Out of Scope

- Do not change comment service rules, comment ordering, ticket modified timestamp behavior, comment editing/deletion, or board behavior.
- Do not change ticket create/edit forms.

#### Acceptance Criteria

- Add-comment dialog tests pass.
- Existing comment route/action/service tests continue to pass.
- Verification passes.
- No always-visible add-comment form remains.

#### Verification

Run from `workspaces/repo-agent-*`: `npm test -- app/routes/tickets/details.test.tsx app/services/comments/comments.server.test.ts app/components/dialog/dialog.test.tsx && npm run typecheck && npm run build`

---

### Ticket Edit Action

Task ID: TICKET-07B
Category: AFK
Owner: Unassigned
Branch:
Status: Backlog
Dependencies: TICKET-07A
Blocking tasks: TICKET-07C, TICKET-WF-01

#### Objective

Persist edits from the ticket edit route through the reusable update service.

#### Scope

- Add edit-route action handling for type, team, epic, title, body, and state.
- Invoke `updateTicket`.
- Map service validation errors to route action data.
- Redirect to the ticket details route after a successful save.
- Add focused action tests for successful save, validation errors, missing tickets, and same-team epic enforcement.

#### Out of Scope

- Do not implement delete behavior, comments, board drag-and-drop, rich-text editing, or presentation polish.
- Do not change ticket update service rules.

#### Acceptance Criteria

- Edit action tests pass.
- Ticket update service tests continue to pass.
- Verification passes.
- No placeholder edit action response remains.

#### Verification

Run from `workspaces/repo-agent-*`: `npm test -- app/routes/tickets/edit.test.tsx app/services/tickets/tickets.server.test.ts && npm run typecheck && npm run build`

---

### Ticket Edit Form UI

Task ID: TICKET-07C
Category: AFK
Owner: Unassigned
Branch:
Status: Backlog
Dependencies: TICKET-07B
Blocking tasks: TICKET-09A, TICKET-WF-01

#### Objective

Render the real ticket edit form using loader and action data.

#### Scope

- Render editable fields for team, epic, type, state, title, and body.
- Populate initial values from the loaded ticket.
- Show same-team epic options for the selected ticket team.
- Show action validation errors.
- Preserve navigation back to ticket details.
- Add focused render tests for initial values, same-team epic options, validation messages, and details navigation.

#### Out of Scope

- Do not implement delete behavior, comments, board drag-and-drop, rich-text editing, or broad styling.
- Do not change backend ticket rules.

#### Acceptance Criteria

- Edit form tests pass.
- Edit action and ticket service tests continue to pass.
- Verification passes.
- No placeholder edit-ticket copy remains.

#### Verification

Run from `workspaces/repo-agent-*`: `npm test -- app/routes/tickets/edit.test.tsx app/services/tickets/tickets.server.test.ts && npm run typecheck && npm run build`

---

### Ticket Delete Action

Task ID: TICKET-09A
Category: AFK
Owner: Unassigned
Branch:
Status: Backlog
Dependencies: TICKET-06, TICKET-04B
Blocking tasks: TICKET-09B

#### Objective

Delete tickets through a confirmed authenticated route action.

#### Scope

- Add delete action handling to the ticket details route or route-specific action helper.
- Require explicit confirmation in submitted form data before deletion.
- Invoke `deleteTicket`.
- Return validation errors for missing confirmation and missing tickets.
- Redirect to the board after successful deletion.
- Add focused action tests for missing confirmation, successful delete, missing ticket handling, and redirect behavior.

#### Out of Scope

- Do not implement edit fields, comments UI, board drag-and-drop, or presentation polish.
- Do not change ticket delete service rules.

#### Acceptance Criteria

- Delete action tests pass.
- Ticket delete service tests continue to pass.
- Verification passes.
- No placeholder delete behavior remains in route action code.

#### Verification

Run from `workspaces/repo-agent-*`: `npm test -- app/routes/tickets/details.test.tsx app/services/tickets/tickets.server.test.ts && npm run typecheck && npm run build`

---

### Ticket Delete Confirmation UI

Task ID: TICKET-09B
Category: AFK
Owner: Unassigned
Branch:
Status: Backlog
Dependencies: TICKET-09A
Blocking tasks: TICKET-WF-01

#### Objective

Expose confirmed ticket deletion from the ticket details screen.

#### Scope

- Render a delete form on ticket details.
- Include an explicit confirmation control or confirmation value required by `TICKET-09A`.
- Show delete validation errors returned by the action.
- Keep successful delete redirect behavior from `TICKET-09A`.
- Add focused render tests for the confirmation affordance and error display.

#### Out of Scope

- Do not change ticket service rules, edit form behavior, comments, or board behavior.

#### Acceptance Criteria

- Delete confirmation UI tests pass.
- Delete action tests continue to pass.
- Verification passes.
- No placeholder delete UI remains.

#### Verification

Run from `workspaces/repo-agent-*`: `npm test -- app/routes/tickets/details.test.tsx && npm run typecheck && npm run build`

---

### Ticket State Update Action

Task ID: TICKET-08
Category: AFK
Owner: Unassigned
Branch:
Status: Backlog
Dependencies: TICKET-04
Blocking tasks: BOARD-03B

#### Objective

Provide a focused authenticated route action for updating only ticket state.

#### Scope

- Add or reuse a route action/API path that updates only ticket state.
- Persist state changes through `updateTicket`.
- Allow direct moves between any two valid states.
- Return meaningful validation, missing-record, authentication, and conflict responses.
- Add focused tests for state persistence, enum validation, missing tickets, and unauthenticated access.

#### Out of Scope

- Do not implement board drag-and-drop UI.
- Do not persist custom manual ordering.
- Do not add sequential transition rules.

#### Acceptance Criteria

- Focused state-update tests pass.
- Ticket update service tests continue to pass.
- Verification passes.
- No placeholder TODOs remain.

#### Verification

Run from `workspaces/repo-agent-*`: `npm test -- app/routes/board/board.test.tsx app/services/tickets/tickets.server.test.ts && npm run typecheck`

---

### Comment Storage Schema

Task ID: COMMENT-01A
Category: AFK
Owner: Unassigned
Branch:
Status: Backlog
Dependencies: TICKET-04B
Blocking tasks: COMMENT-01B

#### Objective

Add append-only ticket comment persistence.

#### Scope

- Add a new migration for ticket comments.
- Add schema fields for id, ticket id, author id, body, and created timestamp.
- Add foreign keys for ticket and author references.
- Ensure ticket delete cascades to comments or is explicitly compatible with `deleteTicket`.
- Add fresh-database coverage for the comments table and references.

#### Out of Scope

- Do not implement comment services, routes, UI, editing, deletion, moderation, or activity history.
- Do not change ticket board ordering.

#### Acceptance Criteria

- Fresh database migration coverage passes.
- `npm run db:migrate` succeeds.
- Existing ticket delete tests continue to pass.
- No placeholder TODOs remain.

#### Verification

Run from `workspaces/repo-agent-*`: `npm test -- app/db/fresh-database.test.ts app/services/tickets/tickets.server.test.ts && npm run db:migrate && npm run typecheck`

---

### Comment Add Service

Task ID: COMMENT-01B
Category: AFK
Owner: Unassigned
Branch:
Status: Backlog
Dependencies: COMMENT-01A
Blocking tasks: COMMENT-01C

#### Objective

Create immutable ticket comments through a reusable backend service.

#### Scope

- Implement `addTicketComment` in a comment service module.
- Require an existing ticket.
- Require authenticated author id supplied by the caller.
- Require non-empty trimmed comment bodies.
- Set created timestamp on the server in UTC.
- Ensure adding a comment does not update the ticket modified timestamp.
- Add focused service tests for required body validation, missing ticket, missing author, author assignment, timestamp assignment, and ticket modified timestamp behavior.

#### Out of Scope

- Do not implement comment listing, route actions, UI, editing, deletion, moderation, or activity history.

#### Acceptance Criteria

- Comment add service tests pass.
- Ticket service tests continue to pass.
- Verification passes.
- No placeholder TODOs remain.

#### Verification

Run from `workspaces/repo-agent-*`: `npm test -- app/services/comments/comments.server.test.ts app/services/tickets/tickets.server.test.ts && npm run typecheck`

---

### Comment List Service

Task ID: COMMENT-01C
Category: AFK
Owner: Unassigned
Branch:
Status: Backlog
Dependencies: COMMENT-01B
Blocking tasks: COMMENT-02A, COMMENT-03A

#### Objective

List immutable comments for a ticket in chronological display order.

#### Scope

- Implement a service function to list comments for a ticket.
- Include author email, body, and created timestamp.
- Return comments oldest first.
- Return an empty list for tickets without comments.
- Add focused tests for chronological ordering, joined author data, empty list behavior, and team-ticket isolation where relevant.

#### Out of Scope

- Do not implement comment routes, UI, editing, deletion, moderation, or activity history.
- Do not change ticket modified timestamp behavior.

#### Acceptance Criteria

- Comment list service tests pass.
- Comment add and ticket service tests continue to pass.
- Verification passes.
- No placeholder TODOs remain.

#### Verification

Run from `workspaces/repo-agent-*`: `npm test -- app/services/comments/comments.server.test.ts app/services/tickets/tickets.server.test.ts && npm run typecheck`

---

### Ticket Details Comments Loader

Task ID: COMMENT-02A
Category: AFK
Owner: Unassigned
Branch:
Status: Backlog
Dependencies: COMMENT-01C, TICKET-06
Blocking tasks: COMMENT-02B, COMMENT-02C

#### Objective

Load comments for the ticket details screen.

#### Scope

- Extend ticket details loader data to include comments for found tickets.
- Use the comment list service.
- Preserve missing-ticket behavior.
- Add focused loader tests for comment ordering, author display data, empty comments, and missing-ticket behavior.

#### Out of Scope

- Do not implement add-comment action, comment form UI, edit/delete comments, moderation, or activity history.

#### Acceptance Criteria

- Details comments loader tests pass.
- Comment service tests continue to pass.
- Verification passes.
- No placeholder comment data remains.

#### Verification

Run from `workspaces/repo-agent-*`: `npm test -- app/routes/tickets/details.test.tsx app/services/comments/comments.server.test.ts && npm run typecheck`

---

### Ticket Details Add Comment Action

Task ID: COMMENT-02B
Category: AFK
Owner: Unassigned
Branch:
Status: Backlog
Dependencies: COMMENT-02A
Blocking tasks: COMMENT-02C

#### Objective

Allow authenticated users to add immutable comments from ticket details.

#### Scope

- Add ticket-details action handling for add-comment submissions.
- Invoke `addTicketComment`.
- Show validation errors for empty comment bodies and missing tickets.
- Redirect back to the same ticket details route after successful creation.
- Add focused action tests for authenticated access, successful creation, validation errors, and missing ticket handling.

#### Out of Scope

- Do not implement comment editing, deletion, moderation, activity history, or board behavior.
- Do not change ticket modified timestamp behavior.

#### Acceptance Criteria

- Add-comment action tests pass.
- Comment service tests continue to pass.
- Verification passes.
- No placeholder comment action remains.

#### Verification

Run from `workspaces/repo-agent-*`: `npm test -- app/routes/tickets/details.test.tsx app/services/comments/comments.server.test.ts && npm run typecheck && npm run build`

---

### Ticket Details Comments UI

Task ID: COMMENT-02C
Category: AFK
Owner: Unassigned
Branch:
Status: Backlog
Dependencies: COMMENT-02B
Blocking tasks: TICKET-WF-01, COMMENT-03B

#### Objective

Display and submit immutable comments on ticket details.

#### Scope

- Display comment author, body, and created timestamp in chronological order.
- Render an add-comment form.
- Show validation errors from the add-comment action.
- Add focused render tests for comment display order, author display, timestamp display, form rendering, and validation error display.

#### Out of Scope

- Do not implement comment editing, deletion, moderation, activity history, or board behavior.
- Do not change ticket modified timestamp behavior.

#### Acceptance Criteria

- Comments UI tests pass.
- Comment action and service tests continue to pass.
- Verification passes.
- No placeholder comment copy remains.

#### Verification

Run from `workspaces/repo-agent-*`: `npm test -- app/routes/tickets/details.test.tsx app/services/comments/comments.server.test.ts && npm run typecheck && npm run build`

---

### Own Comment Mutation Services

Task ID: COMMENT-03A
Category: HITL
Owner: Unassigned
Branch:
Status: Backlog
Dependencies: COMMENT-01C
Blocking tasks: COMMENT-03B

#### Objective

Allow comment authors to edit or delete only their own comments through backend services.

#### Scope

- Implement service functions to edit a comment owned by the caller.
- Implement service functions to delete a comment owned by the caller.
- Prevent users from editing or deleting comments created by other users.
- Preserve created timestamps.
- Add a modified timestamp only if the existing comment model supports it by implementation time; otherwise keep comments immutable except for deletion.
- Add service tests for ownership checks, edit behavior, delete behavior, missing comments, and unauthorized attempts.

#### Out of Scope

- Do not add moderator/admin controls.
- Do not add activity history.
- Do not change ticket modified timestamp behavior unless the comment model explicitly defines comment edits as ticket changes.

#### Acceptance Criteria

- Ownership checks are enforced in backend services.
- Focused service tests pass.
- Verification passes.
- No placeholder TODOs remain.

#### Verification

Run from `workspaces/repo-agent-*`: `npm test -- app/services/comments/comments.server.test.ts && npm run typecheck`

---

### Own Comment Edit/Delete UI

Task ID: COMMENT-03B
Category: HITL
Owner: Unassigned
Branch:
Status: Backlog
Dependencies: COMMENT-03A, COMMENT-02C
Blocking tasks: None

#### Objective

Expose edit and delete controls only for comments owned by the current user.

#### Scope

- Add route action handling for own-comment edit and delete submissions.
- Show edit/delete affordances only for comments owned by the current user.
- Surface unauthorized and validation errors.
- Add route tests for visible owner controls, hidden non-owner controls, successful edit, successful delete, and unauthorized attempts.

#### Out of Scope

- Do not add moderator/admin controls.
- Do not add activity history.
- Do not change ticket modified timestamp behavior unless the service model requires it.

#### Acceptance Criteria

- Route and service tests pass.
- Verification passes.
- No placeholder TODOs remain.

#### Verification

Run from `workspaces/repo-agent-*`: `npm test -- app/routes/tickets/details.test.tsx app/services/comments/comments.server.test.ts && npm run typecheck && npm run build`

---

### Board Column Grouping Presenter

Task ID: BOARD-02A
Category: AFK
Owner: Unassigned
Branch:
Status: Backlog
Dependencies: BOARD-01
Blocking tasks: BOARD-02B

#### Objective

Group loaded board tickets into the four workflow columns in a deterministic order.

#### Scope

- Add a route-local pure grouping helper or presenter.
- Render exactly four columns in workflow order: `backlog`, `todo`, `in-progress`, and `done`.
- Preserve ticket ordering by modified timestamp within each column using the loader-provided order.
- Add focused tests for column order, empty columns, grouping by state, and per-column ordering.

#### Out of Scope

- Do not implement card styling, filters, drag-and-drop, manual ordering, virtualization, or ticket mutations.
- Do not change ticket services.

#### Acceptance Criteria

- Board grouping tests pass.
- Board loader tests continue to pass.
- Verification passes.
- No placeholder column state list remains.

#### Verification

Run from `workspaces/repo-agent-*`: `npm test -- app/routes/board/board.test.tsx && npm run typecheck`

---

### Board Ticket Cards

Task ID: BOARD-02B
Category: AFK
Owner: Unassigned
Branch:
Status: Backlog
Dependencies: BOARD-02A
Blocking tasks: BOARD-03A, BOARD-04A, BOARD-WF-01

#### Objective

Render selected-team tickets as navigable cards in their workflow columns.

#### Scope

- Display each selected-team ticket as a card in its current state column.
- Show ticket title and type.
- Show ticket epic on the card when present.
- Provide a clear create-ticket link from the board.
- Provide a clear open-ticket link for each card.
- Add focused tests for card rendering, missing epic display, create-ticket affordance, and open-ticket affordance.

#### Out of Scope

- Do not add drag-and-drop, filters, manual ordering, virtualization, shared header internals, or broad styling.
- Do not change ticket services.

#### Acceptance Criteria

- Board rendering tests pass.
- Tickets are grouped by state and ordered by modified timestamp within each column.
- Verification passes.
- No placeholder card data remains.

#### Verification

Run from `workspaces/repo-agent-*`: `npm test -- app/routes/board/board.test.tsx && npm run typecheck && npm run build`

---

### Board Filter Query Model

Task ID: BOARD-04A
Category: AFK
Owner: Unassigned
Branch:
Status: Backlog
Dependencies: BOARD-01
Blocking tasks: BOARD-04B

#### Objective

Parse and apply board filters to loaded selected-team tickets.

#### Scope

- Support filtering by ticket type.
- Support filtering by epic.
- Support case-insensitive substring search over ticket title.
- Combine active filters using AND logic.
- Preserve filtered ticket ordering from the loader.
- Add focused tests for ticket type, epic, title search, combined filters, invalid query values, and clear/default query state.

#### Out of Scope

- Do not render filter controls.
- Do not change drag-and-drop, ticket CRUD routes, backend pagination, or broad styling.
- Do not add saved filters.

#### Acceptance Criteria

- Filter query model tests pass.
- Existing board loader tests continue to pass.
- Verification passes.
- No placeholder filter logic remains.

#### Verification

Run from `workspaces/repo-agent-*`: `npm test -- app/routes/board/board.test.tsx && npm run typecheck`

---

### Board Filter Controls

Task ID: BOARD-04B
Category: AFK
Owner: Unassigned
Branch:
Status: Backlog
Dependencies: BOARD-04A, BOARD-02B
Blocking tasks: BOARD-WF-01

#### Objective

Render filter controls that drive the board filter query model.

#### Scope

- Render ticket type, epic, and search controls.
- Reflect active filter values from the current URL query.
- Provide a clear-filters action or link.
- Show filtered tickets in the existing columns.
- Add focused tests for rendered controls, active values, clear action, and filtered card output.

#### Out of Scope

- Do not change drag-and-drop internals, ticket CRUD routes, backend pagination, or broad styling.
- Do not add saved filters.

#### Acceptance Criteria

- Board filter control tests pass.
- Existing board rendering tests continue to pass.
- Verification passes.
- No placeholder filter controls remain.

#### Verification

Run from `workspaces/repo-agent-*`: `npm test -- app/routes/board/board.test.tsx && npm run typecheck && npm run build`

---

### Board Drag Interaction

Task ID: BOARD-03A
Category: HITL
Owner: Unassigned
Branch:
Status: Backlog
Dependencies: BOARD-02B
Blocking tasks: BOARD-03B

#### Objective

Allow users to drag ticket cards between Kanban columns in the client UI.

#### Scope

- Add a drag-and-drop implementation for moving a card from one column to another.
- Update the visible column optimistically on drop.
- Keep cards movable directly between any two states.
- Add automated tests for the client-side drag interaction where practical.

#### Out of Scope

- Do not persist state changes.
- Do not implement rollback on server failure.
- Do not persist custom manual order.
- Do not add filters, virtualization, or keyboard drag alternatives unless required by the chosen drag library.

#### Acceptance Criteria

- Drag interaction tests or documented manual verification pass.
- Existing board render tests continue to pass.
- Verification passes.
- No placeholder TODOs remain.

#### Verification

Run from `workspaces/repo-agent-*`: `npm test -- app/routes/board/board.test.tsx && npm run typecheck && npm run build`

---

### Board Drag Persistence

Task ID: BOARD-03B
Category: HITL
Owner: Unassigned
Branch:
Status: Backlog
Dependencies: BOARD-03A, TICKET-08
Blocking tasks: BOARD-WF-01

#### Objective

Persist Kanban drag-and-drop state changes and roll back failed moves.

#### Scope

- Persist dropped card state changes immediately through the state-update action.
- Return a card to its previous column when persistence fails.
- Display a clear UI error when drag-and-drop persistence fails.
- Add automated tests for successful persistence, failed drag rollback, and state correctness after refresh.

#### Out of Scope

- Do not persist custom manual order.
- Do not add filters, virtualization, or unrelated board styling.

#### Acceptance Criteria

- Drag persistence happy path and rollback tests pass.
- State changes remain correct after refresh.
- Verification passes.
- No placeholder TODOs remain.

#### Verification

Run from `workspaces/repo-agent-*`: `npm test -- app/routes/board/board.test.tsx && npm run typecheck && npm run build`

---

### Authenticated Shell Navigation

Task ID: SHELL-01A
Category: AFK
Owner: Unassigned
Branch:
Status: Backlog
Dependencies: BOARD-02B, TICKET-07C, TICKET-09B, EPIC-UI-01, TEAM-UI-01
Blocking tasks: SHELL-01B, TICKET-WF-01, BOARD-WF-01

#### Objective

Provide a consistent authenticated header component for product screens.

#### Scope

- Refine the existing authenticated header or add the missing shell wrapper.
- Include `TICKET TRACKER` branding.
- Include Board, Teams, and Epics navigation.
- Show active-route indication, current user email, and logout affordance.
- Add focused component tests for navigation, active-route indication, current user display, and logout affordance.

#### Out of Scope

- Do not change backend business rules, database schema, authentication behavior, route permissions, or session behavior.
- Do not redesign public auth screens.
- Do not wire every route into the shell in this task unless the existing component API already does so.

#### Acceptance Criteria

- Authenticated header component tests pass.
- Verification passes.
- No placeholder TODOs remain.

#### Verification

Run from `workspaces/repo-agent-*`: `npm test -- app/components/authenticated-header/authenticated-header.test.tsx && npm run typecheck && npm run build`

---

### Authenticated Shell Route Adoption

Task ID: SHELL-01B
Category: AFK
Owner: Unassigned
Branch:
Status: Backlog
Dependencies: SHELL-01A
Blocking tasks: TICKET-WF-01, BOARD-WF-01

#### Objective

Use the authenticated shell consistently on board, ticket, team, and epic screens.

#### Scope

- Apply the authenticated shell/header to board, ticket create, ticket details, ticket edit, teams, and epics screens.
- Keep public authentication screens outside the authenticated shell.
- Preserve existing route behavior and actions.
- Add focused route tests for visible shell navigation on authenticated product screens and absence on public auth screens.

#### Out of Scope

- Do not change backend business rules, database schema, authentication behavior, route permissions, or public auth screen design.

#### Acceptance Criteria

- Authenticated shell route tests pass.
- Public auth screens remain outside the shell.
- Verification passes.
- No placeholder shell copy remains.

#### Verification

Run from `workspaces/repo-agent-*`: `npm test -- app/routes/board/board.test.tsx app/routes/teams/teams.test.tsx app/routes/epics/epics.test.tsx app/routes/auth/auth-screens.test.tsx && npm run typecheck && npm run build`

---

### Ticket Wireframe UI Alignment

Task ID: TICKET-WF-01
Category: HITL
Owner: Unassigned
Branch:
Status: Backlog
Dependencies: TICKET-07C, TICKET-09B, COMMENT-02C, UI-DIALOG-03, UI-DIALOG-04
Blocking tasks: CLEANUP-02

#### Objective

Align ticket create, edit, details, and comments presentation with the ticket wireframe hierarchy.

#### Scope

- Shape ticket details/editing around back-to-board navigation, ticket metadata, editable fields, clear save/delete actions, and comments panel.
- Show comments count, chronological comments, and add-comment form.
- Preserve existing ticket create, view, edit, delete, same-team epic selection, comment creation, and comment ordering behavior.
- Keep ticket deletion behind explicit confirmation.
- Preserve required loading, empty, success, and error states.
- Add or update focused tests for back navigation, metadata display, editable fields, save/delete actions, confirmation, comments count, chronological comments, and add-comment form.

#### Out of Scope

- Do not change backend ticket services, comment services, ticket workflow rules, comment mutability, database schema, or board drag-and-drop behavior.
- Do not add comment editing or deletion.

#### Acceptance Criteria

- Existing ticket and comment behavior still works.
- Focused wireframe UI tests pass.
- Verification passes.
- No placeholder TODOs remain.

#### Verification

Run from `workspaces/repo-agent-*`: `npm test -- app/routes/tickets/details.test.tsx app/routes/tickets/edit.test.tsx app/routes/tickets/new.test.tsx && npm run typecheck && npm run build`

---

### Board Wireframe UI Alignment

Task ID: BOARD-WF-01
Category: HITL
Owner: Unassigned
Branch:
Status: Backlog
Dependencies: BOARD-03B, BOARD-04B, UI-DIALOG-03
Blocking tasks: CLEANUP-02

#### Objective

Align the primary Kanban board screen with the board wireframe hierarchy.

#### Scope

- Show team selector, prominent new-ticket action, search/type/epic filters, clear-filter action, total ticket count, four workflow columns, and per-column counts.
- Show compact cards with type, title, epic when present, and modified recency where practical.
- Preserve ticket state ordering, filtering behavior, create-ticket navigation, open-ticket navigation, drag-and-drop persistence, rollback, and error handling.
- Preserve required loading, empty, success, and error states.
- Add or update focused tests for team selector, new-ticket action, filters, clear action, ticket count, column counts, card content, and drag/drop error visibility.

#### Out of Scope

- Do not change backend ticket services, ticket state API behavior, database schema, or drag-and-drop persistence semantics.
- Do not add virtualization.

#### Acceptance Criteria

- Existing board behavior still works.
- Focused board UI tests pass.
- Verification passes.
- No placeholder TODOs remain.

#### Verification

Run from `workspaces/repo-agent-*`: `npm test -- app/routes/board/board.test.tsx && npm run typecheck && npm run build`

---

### Ticket Activity Schema

Task ID: ACTIVITY-01A
Category: HITL
Owner: Unassigned
Branch:
Status: Backlog
Dependencies: TICKET-07B
Blocking tasks: ACTIVITY-01B

#### Objective

Persist append-only ticket activity entries.

#### Scope

- Add a new migration and schema for ticket activity entries.
- Store ticket id, actor id, timestamp, action type, and relevant changed values.
- Add database references for ticket and actor where practical.
- Add fresh-database tests for the table and references.

#### Out of Scope

- Do not write activity from ticket services.
- Do not display activity in the UI.
- Do not record session identifiers, verification tokens, password-reset tokens, or password-related values.

#### Acceptance Criteria

- Fresh database migration coverage passes.
- `npm run db:migrate` succeeds.
- Verification passes.
- No placeholder TODOs remain.

#### Verification

Run from `workspaces/repo-agent-*`: `npm test -- app/db/fresh-database.test.ts && npm run db:migrate && npm run typecheck`

---

### Ticket Activity Service

Task ID: ACTIVITY-01B
Category: HITL
Owner: Unassigned
Branch:
Status: Backlog
Dependencies: ACTIVITY-01A
Blocking tasks: ACTIVITY-01C

#### Objective

Create and list append-only ticket activity entries through a backend service.

#### Scope

- Implement service functions for recording and listing ticket activity.
- Keep entries append-only.
- Validate known action types.
- Return entries oldest first.
- Add service tests for creation, ordering, invalid action types, missing ticket, and actor assignment.

#### Out of Scope

- Do not wire activity into ticket create/update/delete services.
- Do not display activity in the UI.

#### Acceptance Criteria

- Activity service tests pass.
- Verification passes.
- No placeholder TODOs remain.

#### Verification

Run from `workspaces/repo-agent-*`: `npm test -- app/services/ticket-activity/ticket-activity.server.test.ts && npm run typecheck`

---

### Ticket Activity Write Integration

Task ID: ACTIVITY-01C
Category: HITL
Owner: Unassigned
Branch:
Status: Backlog
Dependencies: ACTIVITY-01B
Blocking tasks: ACTIVITY-01D

#### Objective

Record ticket activity from backend ticket mutations.

#### Scope

- Record ticket creation, status changes, title changes, body changes, team changes, epic changes, and deletion where supported.
- Write entries from backend services or action helpers, not display-only route code.
- Avoid recording sensitive authentication data.
- Add focused tests for each recorded change type and unchanged-save behavior.

#### Out of Scope

- Do not display activity in the UI.
- Do not add filtering, export, moderation, or editable activity entries.

#### Acceptance Criteria

- Activity write tests pass.
- Existing ticket service tests continue to pass.
- Verification passes.
- No placeholder TODOs remain.

#### Verification

Run from `workspaces/repo-agent-*`: `npm test -- app/services/ticket-activity/ticket-activity.server.test.ts app/services/tickets/tickets.server.test.ts && npm run typecheck`

---

### Ticket Activity UI

Task ID: ACTIVITY-01D
Category: HITL
Owner: Unassigned
Branch:
Status: Backlog
Dependencies: ACTIVITY-01C, TICKET-06
Blocking tasks: None

#### Objective

Display ticket activity history on the ticket details screen.

#### Scope

- Load activity entries on ticket details.
- Display activity entries in chronological order.
- Add focused route tests for activity display and ordering.

#### Out of Scope

- Do not add filtering, export, moderation, or editable activity entries.
- Do not change ticket mutation rules.

#### Acceptance Criteria

- Activity display tests pass.
- Verification passes.
- No placeholder TODOs remain.

#### Verification

Run from `workspaces/repo-agent-*`: `npm test -- app/routes/tickets/details.test.tsx app/services/ticket-activity/ticket-activity.server.test.ts && npm run typecheck && npm run build`

---

### Large Board Profiling

Task ID: BOARD-PERF-01A
Category: HITL
Owner: Unassigned
Branch:
Status: Backlog
Dependencies: BOARD-02B
Blocking tasks: BOARD-PERF-01B

#### Objective

Determine whether the non-virtualized board remains usable with large ticket counts.

#### Scope

- Profile or manually measure board rendering with a large selected-team ticket set.
- Record profiling notes in the task progress.
- Recommend whether virtualization is necessary.
- Add a focused test or fixture only if it can be kept stable and fast.

#### Out of Scope

- Do not implement virtualization.
- Do not change backend pagination, filters, drag-and-drop, or ticket services.

#### Acceptance Criteria

- Profiling notes justify whether implementation is needed.
- Existing board tests continue to pass.
- Verification passes.

#### Verification

Run from `workspaces/repo-agent-*`: `npm test -- app/routes/board/board.test.tsx && npm run typecheck && npm run build`

---

### Large Board Virtualization

Task ID: BOARD-PERF-01B
Category: HITL
Owner: Unassigned
Branch:
Status: Backlog
Dependencies: BOARD-PERF-01A
Blocking tasks: None

#### Objective

Virtualize ticket rendering only if profiling shows it is necessary.

#### Scope

- Virtualize ticket rendering within board columns if justified by `BOARD-PERF-01A`.
- Preserve keyboard and pointer interactions for visible ticket cards.
- Keep column headers and empty states stable while lists virtualize.
- Add focused tests or profiling notes that demonstrate large-board rendering remains usable.

#### Out of Scope

- Do not change backend pagination or ticket APIs unless profiling shows client-side virtualization is insufficient.
- Do not change filters, drag-and-drop persistence semantics, or ticket services.

#### Acceptance Criteria

- Profiling notes justify the implementation approach.
- Large-board rendering remains usable.
- Verification passes.
- No placeholder TODOs remain.

#### Verification

Run from `workspaces/repo-agent-*`: `npm test -- app/routes/board/board.test.tsx && npm run typecheck && npm run build`

---

### Home Entry Cleanup

Task ID: CLEANUP-01
Category: AFK
Owner: Unassigned
Branch:
Status: Backlog
Dependencies: BOARD-01
Blocking tasks: CLEANUP-02

#### Objective

Replace the home screen placeholder with a real authenticated entry point.

#### Scope

- Redirect authenticated users from home to the primary board or render a minimal real product entry point following existing route conventions.
- Preserve public authentication behavior.
- Update focused home route tests.

#### Out of Scope

- Do not remove placeholder services or broad README copy.
- Do not add new feature behavior.
- Do not change business rules.

#### Acceptance Criteria

- Home route tests pass.
- Verification passes.
- No placeholder-only home copy remains.

#### Verification

Run from `workspaces/repo-agent-*`: `npm test -- app/routes/home/home.test.tsx && npm run typecheck && npm run build`

---

### Placeholder Scaffolding Removal

Task ID: CLEANUP-02
Category: AFK
Owner: Unassigned
Branch:
Status: Backlog
Dependencies: CLEANUP-01, TICKET-WF-01, BOARD-WF-01, COMMENT-02C
Blocking tasks: CLEANUP-03

#### Objective

Remove placeholder-only routes, services, sample data, and tests after real workflows replace them.

#### Scope

- Remove the placeholder service module and focused tests once no production route depends on it.
- Remove placeholder-only route copy, sample ticket cards, sample team/epic options, and placeholder status payloads from ticket and board routes.
- Remove or rewrite placeholder-focused tests so coverage validates real route behavior.
- Preserve public authentication route behavior and shared helpers still used by implemented routes.

#### Out of Scope

- Do not add new feature behavior while removing placeholders.
- Do not change business rules.
- Do not update README in this task.

#### Acceptance Criteria

- Placeholder-only code is gone where no longer used.
- Relevant route and service tests pass.
- Verification passes.
- No placeholder TODOs remain.

#### Verification

Run from `workspaces/repo-agent-*`: `npm test && npm run typecheck && npm run build`

---

### Product README Cleanup

Task ID: CLEANUP-03
Category: AFK
Owner: Unassigned
Branch:
Status: Backlog
Dependencies: CLEANUP-02
Blocking tasks: DONE-01

#### Objective

Update product README copy so it describes the implemented ticket tracker.

#### Scope

- Replace starter or placeholder introduction text with product-specific setup and usage copy.
- Keep commands accurate for the current package scripts.
- Preserve useful existing development setup instructions.

#### Out of Scope

- Do not change product code.
- Do not add new feature behavior.

#### Acceptance Criteria

- README no longer describes starter scaffolding as the product.
- Setup/test commands are accurate.
- Verification passes.

#### Verification

Run from `workspaces/repo-agent-*`: `npm run typecheck`

---

### Final Completion Check

Task ID: DONE-01
Category: HITL
Owner: Unassigned
Branch:
Status: Backlog
Dependencies: CLEANUP-03
Blocking tasks: None

#### Objective

Verify the product is ready for final handoff after feature and cleanup work.

#### Scope

- Run the full product verification suite.
- Check that `BACKLOG.md` and `TASKS.md` contain no remaining required product work.
- Check that completed work has been archived.
- Record final verification notes.

#### Out of Scope

- Do not add feature behavior.
- Do not perform unrelated refactors.

#### Acceptance Criteria

- Full test suite passes.
- Typecheck passes.
- Production build passes.
- No active task remains unarchived.

#### Verification

Run from `workspaces/repo-agent-*`: `npm test && npm run typecheck && npm run build`
