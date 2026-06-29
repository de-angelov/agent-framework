# BACKLOG

Pending work only. Active work lives in TASKS.md. Completed work lives in ARCHIVE.md.

## Backlog

### Ticket Create Service

Task ID: TICKET-02
Category: AFK
Owner: Unassigned
Branch:
Status: Backlog
Dependencies: TICKET-01
Blocking tasks: TICKET-03, TICKET-05

#### Objective

Create tickets through a reusable backend service with required validation and timestamps.

#### Scope

- Implement `createTicket` in a ticket service module.
- Require existing team, non-empty trimmed title, non-empty body, valid ticket type, and valid ticket state.
- Allow a null epic or an epic from the same team as the ticket.
- Reject epics from other teams.
- Set created-by from the authenticated user id supplied by the caller.
- Set created-at and modified-at on the server in UTC.
- Add focused service tests for required fields, enum validation, reference validation, same-team epic validation, timestamp assignment, and created-by assignment.

#### Out of Scope

- Do not implement update, delete, comments, route actions, or UI.
- Do not add workflow transition restrictions.

#### Acceptance Criteria

- Ticket creation service tests pass.
- Existing team and epic tests continue to pass.
- Verification passes.
- No placeholder TODOs remain.

#### Verification

Run from `workspaces/repo-agent-*`: `npm test -- app/services/tickets.server.test.ts app/services/teams.server.test.ts app/services/epics.server.test.ts && npm run typecheck`

---

### Ticket Read and List Services

Task ID: TICKET-03
Category: AFK
Owner: Unassigned
Branch:
Status: Backlog
Dependencies: TICKET-02
Blocking tasks: TICKET-05, TICKET-06, BOARD-01

#### Objective

Provide reusable ticket read models for ticket screens and the team board.

#### Scope

- Implement service functions to get one ticket by id.
- Implement service functions to list tickets for a team.
- Include enough joined data for UI display: team name, optional epic title, created-by email, created-at, and modified-at.
- Order listed tickets by most recently modified first.
- Add focused tests for missing tickets, joined display data, null epic handling, team filtering, and ordering.

#### Out of Scope

- Do not implement create, update, delete, comments, route UI, board UI, or filters.
- Do not add pagination or virtualization.

#### Acceptance Criteria

- Ticket read/list service tests pass.
- Existing service tests continue to pass.
- Verification passes.
- No placeholder TODOs remain.

#### Verification

Run from `workspaces/repo-agent-*`: `npm test -- app/services/tickets.server.test.ts && npm run typecheck`

---

### Ticket Update Service

Task ID: TICKET-04
Category: AFK
Owner: Unassigned
Branch:
Status: Backlog
Dependencies: TICKET-03
Blocking tasks: TICKET-04B, TICKET-07, TICKET-08

#### Objective

Update ticket fields through a reusable backend service.

#### Scope

- Implement ticket field updates for type, team, epic, title, body, and state.
- Validate submitted enum values and references in the backend.
- Ensure the selected epic is null or belongs to the ticket's team after the update.
- Advance modified-at only when persisted ticket values actually change.
- Allow direct state changes between any two valid states.
- Add focused service tests for updates, unchanged saves, same-team epic rules, and direct state changes.

#### Out of Scope

- Do not implement ticket deletion, route actions, forms, comments creation, activity history, or board drag-and-drop UI.
- Do not persist custom manual ordering.

#### Acceptance Criteria

- Ticket update service tests pass.
- Existing ticket create/read tests continue to pass.
- Verification passes.
- No placeholder TODOs remain.

#### Verification

Run from `workspaces/repo-agent-*`: `npm test -- app/services/tickets.server.test.ts && npm run typecheck`

---

### Ticket Delete Service

Task ID: TICKET-04B
Category: AFK
Owner: Unassigned
Branch:
Status: Backlog
Dependencies: TICKET-04
Blocking tasks: TICKET-09, COMMENT-01

#### Objective

Delete tickets through a reusable backend service.

#### Scope

- Implement ticket deletion by id.
- Return a missing-record error when the ticket does not exist.
- Preserve team and epic blocked-delete behavior for existing referenced tickets.
- If a comments table already exists by implementation time, delete ticket comments when the ticket is deleted.
- Add focused service tests for successful deletion, missing-ticket deletion, and continued team/epic blocked-delete behavior.

#### Out of Scope

- Do not implement route actions, delete confirmation UI, comment creation, activity history, or board behavior.
- Do not change ticket update rules.

#### Acceptance Criteria

- Ticket delete service tests pass.
- Existing ticket create/read/update tests continue to pass.
- Verification passes.
- No placeholder TODOs remain.

#### Verification

Run from `workspaces/repo-agent-*`: `npm test -- app/services/tickets.server.test.ts app/services/teams.server.test.ts app/services/epics.server.test.ts && npm run typecheck`

---

### Ticket Create Route

Task ID: TICKET-05
Category: AFK
Owner: Unassigned
Branch:
Status: Backlog
Dependencies: TICKET-03
Blocking tasks: None

#### Objective

Replace the create-ticket placeholder with a working authenticated create screen.

#### Scope

- Load teams and team-scoped epics needed by the create form.
- Render fields for team, epic, type, state, title, and body.
- Submit through the ticket create service.
- Show validation errors from the service.
- Redirect to the created ticket details route on success.
- Ensure the epic options shown for a selected team are same-team only.
- Add focused route tests for form rendering, successful create, validation errors, and same-team epic options.

#### Out of Scope

- Do not implement edit, delete, comments, board behavior, rich-text editing, or visual redesign.
- Do not change ticket service business rules.

#### Acceptance Criteria

- Create-ticket route tests pass.
- Ticket service tests continue to pass.
- Verification passes.
- No placeholder create-ticket copy remains.

#### Verification

Run from `workspaces/repo-agent-*`: `npm test -- app/routes/tickets.new.test.tsx app/services/tickets.server.test.ts && npm run typecheck && npm run build`

---

### Ticket Details Route

Task ID: TICKET-06
Category: AFK
Owner: Unassigned
Branch:
Status: Backlog
Dependencies: TICKET-03
Blocking tasks: COMMENT-02, TICKET-WF-01

#### Objective

Replace the ticket details placeholder with a working authenticated read screen.

#### Scope

- Load one ticket by id through the ticket read service.
- Display type, team, optional epic, state with human-readable label, title, body, created-by, created-at, and modified-at.
- Show a missing-record response for unknown ticket ids.
- Provide navigation to the edit route.
- Add focused route tests for field display, human-readable state labels, null epic display, missing ticket handling, and edit navigation.

#### Out of Scope

- Do not implement edit form, delete action, comments, board behavior, or wireframe presentation polish.
- Do not change ticket service business rules.

#### Acceptance Criteria

- Ticket details route tests pass.
- Ticket service tests continue to pass.
- Verification passes.
- No placeholder details copy remains.

#### Verification

Run from `workspaces/repo-agent-*`: `npm test -- 'app/routes/tickets.$ticketId.test.tsx' app/services/tickets.server.test.ts && npm run typecheck && npm run build`

---

### Ticket Edit Route

Task ID: TICKET-07
Category: AFK
Owner: Unassigned
Branch:
Status: Backlog
Dependencies: TICKET-04, TICKET-06
Blocking tasks: TICKET-09, TICKET-WF-01

#### Objective

Replace the ticket edit placeholder with working edit behavior.

#### Scope

- Load the ticket, teams, and team-scoped epics needed by the edit form.
- Allow editing type, team, epic, title, body, and state.
- Clear or require replacement of an invalid selected epic when the ticket team changes.
- Save through the ticket update service.
- Redirect to the ticket details route or remain on the edit route with a success message after successful save, following existing route conventions.
- Add focused route tests for edit rendering, successful save, validation errors, and same-team epic selection.

#### Out of Scope

- Do not implement delete behavior, comments UI, board drag-and-drop, rich-text editing, or wireframe presentation polish.
- Do not change ticket service business rules.

#### Acceptance Criteria

- Ticket edit route tests pass.
- Ticket service tests continue to pass.
- Verification passes.
- No placeholder edit-ticket copy remains.

#### Verification

Run from `workspaces/repo-agent-*`: `npm test -- 'app/routes/tickets.$ticketId.edit.test.tsx' app/services/tickets.server.test.ts && npm run typecheck && npm run build`

---

### Ticket Delete Route

Task ID: TICKET-09
Category: AFK
Owner: Unassigned
Branch:
Status: Backlog
Dependencies: TICKET-04B, TICKET-06
Blocking tasks: TICKET-WF-01

#### Objective

Allow users to delete a ticket only after explicit confirmation.

#### Scope

- Add delete action handling to the ticket details or edit route, following existing route conventions.
- Require explicit confirmation in the submitted form data before deletion.
- Delete through the ticket delete service.
- Show a validation error when confirmation is missing.
- Redirect to the board or another existing safe route after successful delete.
- Add focused route tests for missing confirmation, successful delete, missing ticket handling, and redirect behavior.

#### Out of Scope

- Do not implement edit fields, comments UI, board drag-and-drop, or wireframe presentation polish.
- Do not change ticket delete service business rules.

#### Acceptance Criteria

- Ticket delete route tests pass.
- Ticket delete service tests continue to pass.
- Verification passes.
- No placeholder delete behavior remains.

#### Verification

Run from `workspaces/repo-agent-*`: `npm test -- 'app/routes/tickets.$ticketId.test.tsx' app/services/tickets.server.test.ts && npm run typecheck && npm run build`

---

### Ticket State Update API

Task ID: TICKET-08
Category: AFK
Owner: Unassigned
Branch:
Status: Backlog
Dependencies: TICKET-04
Blocking tasks: BOARD-03

#### Objective

Provide a focused persistent state-update path for ticket editing and Kanban movement.

#### Scope

- Add or reuse a route action/API path that updates only ticket state where practical.
- Persist state changes immediately in the database.
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

Run from `workspaces/repo-agent-*`: `npm test -- app/services/tickets.server.test.ts app/routes/board.test.tsx && npm run typecheck`

---

### Immutable Ticket Comment Schema and Services

Task ID: COMMENT-01
Category: AFK
Owner: Unassigned
Branch:
Status: Backlog
Dependencies: TICKET-04B
Blocking tasks: COMMENT-02, COMMENT-03

#### Objective

Implement append-only ticket comment persistence and service behavior.

#### Scope

- Add a new migration for ticket comments with identifier, ticket reference, author reference, body, and created timestamp.
- Implement service functions to add a comment and list comments for a ticket.
- Require authenticated author id supplied by the caller.
- Require non-empty comment bodies.
- Return comments chronologically with oldest first.
- Ensure adding a comment does not update the ticket modified timestamp.
- Add focused service tests for required body validation, author assignment, chronological ordering, missing ticket handling, and ticket modified timestamp behavior.

#### Out of Scope

- Do not implement comment UI.
- Do not implement comment editing, deletion, moderation, or activity history.
- Do not change ticket board ordering.

#### Acceptance Criteria

- Comment service tests pass.
- Ticket service tests continue to pass.
- `npm run db:migrate` succeeds.
- No placeholder TODOs remain.

#### Verification

Run from `workspaces/repo-agent-*`: `npm test -- app/services/comments.server.test.ts app/services/tickets.server.test.ts && npm run db:migrate && npm run typecheck`

---

### Ticket Comments UI

Task ID: COMMENT-02
Category: AFK
Owner: Unassigned
Branch:
Status: Backlog
Dependencies: COMMENT-01, TICKET-06
Blocking tasks: TICKET-WF-01, COMMENT-03

#### Objective

Allow authenticated users to add and view immutable comments on the ticket details screen.

#### Scope

- Load comments on the ticket details route.
- Display comment author, body, and created timestamp in chronological order.
- Add a comment form to the ticket details route.
- Submit new comments through the comment service.
- Show validation errors for empty comment bodies.
- Add focused route tests for authenticated access, comment display order, successful creation, author display, timestamp display, and validation errors.

#### Out of Scope

- Do not implement comment editing, deletion, moderation, activity history, or board behavior.
- Do not change ticket modified timestamp behavior.

#### Acceptance Criteria

- Ticket details comment tests pass.
- Comment service tests continue to pass.
- Verification passes.
- No placeholder comment copy remains.

#### Verification

Run from `workspaces/repo-agent-*`: `npm test -- 'app/routes/tickets.$ticketId.test.tsx' app/services/comments.server.test.ts && npm run typecheck && npm run build`

---

### Edit or Delete Own Comments

Task ID: COMMENT-03
Category: HITL
Owner: Unassigned
Branch:
Status: Backlog
Dependencies: COMMENT-02
Blocking tasks: None

#### Objective

Allow authenticated users to modify or remove comments they created.

#### Scope

- Allow users to edit their own ticket comments.
- Allow users to delete their own ticket comments.
- Prevent users from editing or deleting comments created by other users.
- Preserve created timestamps and add or preserve modified timestamps according to the existing comment model.
- Show edit/delete affordances only where the current user owns the comment.
- Add automated tests for ownership checks, edit behavior, delete behavior, and unauthorized attempts.

#### Out of Scope

- Do not add moderator or admin comment controls.
- Do not add ticket activity history.
- Do not change ticket modified timestamp behavior unless the existing comments model defines comment edits as ticket changes.

#### Acceptance Criteria

- Ownership checks are enforced in backend services.
- Focused route and service tests pass.
- Verification passes.
- No placeholder TODOs remain.

#### Verification

Run from `workspaces/repo-agent-*`: `npm test -- app/services/comments.server.test.ts 'app/routes/tickets.$ticketId.test.tsx' && npm run typecheck && npm run build`

---

### Kanban Board Data Loader

Task ID: BOARD-01
Category: AFK
Owner: Unassigned
Branch:
Status: Backlog
Dependencies: TICKET-03
Blocking tasks: BOARD-02, BOARD-04

#### Objective

Load the real team board data needed by the Kanban screen.

#### Scope

- Replace placeholder board loader data with real teams, selected team, team epics, and selected-team tickets.
- Select a deterministic default team when no team is requested.
- Preserve authenticated access.
- Keep tickets ordered by most recently modified first from the ticket list service.
- Add focused loader tests for default team selection, explicit team selection, empty team state, ticket loading, and unauthenticated redirects.

#### Out of Scope

- Do not change board presentation beyond removing placeholder data dependencies.
- Do not add filters, drag-and-drop, virtualization, or shared shell changes.

#### Acceptance Criteria

- Board loader tests pass.
- Ticket service tests continue to pass.
- Verification passes.
- No placeholder board data remains in the loader.

#### Verification

Run from `workspaces/repo-agent-*`: `npm test -- app/routes/board.test.tsx app/services/tickets.server.test.ts && npm run typecheck`

---

### Kanban Columns and Ticket Cards

Task ID: BOARD-02
Category: AFK
Owner: Unassigned
Branch:
Status: Backlog
Dependencies: BOARD-01
Blocking tasks: BOARD-03, BOARD-04, BOARD-WF-01

#### Objective

Render the primary team Kanban board with real workflow columns and ticket cards.

#### Scope

- Render exactly five columns in workflow order.
- Display each selected-team ticket as a card in its current state column.
- Show at least ticket title and type on each card.
- Show the ticket epic on the card when present.
- Provide a clear create-ticket link from the board.
- Provide a clear open-ticket link for each card.
- Keep the interface usable with at least 100 tickets on one team board where practical without virtualization.
- Add focused route/component tests for column ordering, card rendering, create/open affordances, and 100-ticket rendering where practical.

#### Out of Scope

- Do not add drag-and-drop, filters, manual ordering, virtualization, shared header internals, or broad styling.
- Do not change ticket services.

#### Acceptance Criteria

- Board rendering tests pass.
- Tickets are grouped by state and ordered by modified timestamp within each column.
- Verification passes.
- No placeholder card data remains.

#### Verification

Run from `workspaces/repo-agent-*`: `npm test -- app/routes/board.test.tsx && npm run typecheck && npm run build`

---

### Kanban Drag and Drop Persistence

Task ID: BOARD-03
Category: HITL
Owner: Unassigned
Branch:
Status: Backlog
Dependencies: BOARD-02, TICKET-08
Blocking tasks: BOARD-WF-01

#### Objective

Allow users to move tickets between Kanban columns with immediate persisted state updates and rollback on failure.

#### Scope

- Allow users to drag a ticket card from one column to another.
- Persist dropped card state changes immediately through the backend state-update path.
- Return a card to its previous column when a drag-and-drop update fails.
- Display a clear UI error when a drag-and-drop update fails.
- Allow cards to move directly between any two states.
- Add automated tests for drag-and-drop persistence and failed drag rollback.

#### Out of Scope

- Do not persist custom manual order.
- Do not add filters, virtualization, or keyboard drag alternatives unless required by the chosen drag library.

#### Acceptance Criteria

- Drag-and-drop happy path and rollback tests pass.
- State changes remain correct after refresh.
- Verification passes.
- No placeholder TODOs remain.

#### Verification

Run from `workspaces/repo-agent-*`: `npm test -- app/routes/board.test.tsx && npm run typecheck && npm run build`

---

### Kanban Board Filters

Task ID: BOARD-04
Category: AFK
Owner: Unassigned
Branch:
Status: Backlog
Dependencies: BOARD-02
Blocking tasks: BOARD-WF-01

#### Objective

Add filtering to the Kanban board.

#### Scope

- Provide filtering by ticket type.
- Provide filtering by epic.
- Provide case-insensitive substring search over ticket title.
- Combine active filters using AND logic.
- Keep filtered column ordering by most recently modified first.
- Add focused tests for ticket type, epic, title search, and combined filters.

#### Out of Scope

- Do not change drag-and-drop internals, ticket CRUD routes, backend pagination, or broad styling.
- Do not add saved filters.

#### Acceptance Criteria

- Board filter tests pass.
- Existing board rendering tests continue to pass.
- Verification passes.
- No placeholder TODOs remain.

#### Verification

Run from `workspaces/repo-agent-*`: `npm test -- app/routes/board.test.tsx && npm run typecheck && npm run build`

---

### Wireframe Shared App Shell

Task ID: SHELL-01
Category: AFK
Owner: Unassigned
Branch:
Status: Backlog
Dependencies: BOARD-02, TICKET-05, TICKET-06, TICKET-07, TICKET-09, Epic Management UI merged, Teams UI merged
Blocking tasks: Wireframe Board UI Alignment, Wireframe Ticket Details and Comments UI, Wireframe Epic Management UI Alignment

#### Objective

Create the shared authenticated application shell shown across the wireframes.

#### Scope

- Add or refine a consistent authenticated header with `TICKET TRACKER` branding.
- Include Board, Teams, and Epics navigation.
- Show active-route indication, current user email, and a logout affordance.
- Keep public authentication screens outside the authenticated shell.
- Ensure the header works consistently on board, ticket, team, and epic screens.
- Add focused tests for visible navigation, active-route indication, current user display, and logout affordance.

#### Out of Scope

- Do not change backend business rules, database schema, authentication behavior, route permissions, or session behavior.
- Do not redesign public auth screens.

#### Acceptance Criteria

- Authenticated shell tests pass.
- Public auth screens remain outside the shell.
- Verification passes.
- No placeholder TODOs remain.

#### Verification

Run from `workspaces/repo-agent-*`: `npm test -- app/components/authenticated-header/authenticated-header.test.tsx app/routes/board.test.tsx app/routes/teams.test.tsx app/routes/epics.test.tsx && npm run typecheck && npm run build`

---

### Wireframe Ticket Details and Comments UI

Task ID: TICKET-WF-01
Category: HITL
Owner: Unassigned
Branch:
Status: Backlog
Dependencies: TICKET-07, TICKET-09, COMMENT-02; SHELL-01 if it lands first
Blocking tasks: None

#### Objective

Align ticket create, edit, details, and comments presentation with the ticket wireframe hierarchy.

#### Scope

- Shape ticket details/editing around back-to-board navigation, ticket metadata, editable team/type/state/epic/title/body fields, clear save and delete actions, and a comments panel.
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

Run from `workspaces/repo-agent-*`: `npm test -- 'app/routes/tickets.$ticketId.test.tsx' 'app/routes/tickets.$ticketId.edit.test.tsx' app/routes/tickets.new.test.tsx && npm run typecheck && npm run build`

---

### Wireframe Board UI Alignment

Task ID: BOARD-WF-01
Category: HITL
Owner: Unassigned
Branch:
Status: Backlog
Dependencies: BOARD-03, BOARD-04; SHELL-01 if it lands first
Blocking tasks: None

#### Objective

Align the primary Kanban board screen with the board wireframe hierarchy.

#### Scope

- Show team selector, prominent new-ticket action, search/type/epic filters, clear-filter action, total ticket count, five workflow columns, and per-column counts.
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

Run from `workspaces/repo-agent-*`: `npm test -- app/routes/board.test.tsx && npm run typecheck && npm run build`

---

### Ticket Activity History

Task ID: ACTIVITY-01
Category: HITL
Owner: Unassigned
Branch:
Status: Backlog
Dependencies: TICKET-07
Blocking tasks: None

#### Objective

Record and display meaningful ticket changes over time.

#### Scope

- Record ticket creation, status changes, title changes, description/body changes, team changes, epic changes, and assignment changes where supported.
- Store activity entries with actor, timestamp, action type, and relevant changed values.
- Display ticket activity history on the ticket details screen.
- Keep activity entries append-only.
- Avoid recording sensitive authentication data.
- Add automated tests for activity creation and display ordering.

#### Out of Scope

- Do not add filtering, export, moderation, or editable activity entries.
- Do not record session identifiers, verification tokens, password-reset tokens, or password-related values.

#### Acceptance Criteria

- Activity is written from backend ticket services, not route-only code.
- Activity display tests pass.
- Verification passes.
- No placeholder TODOs remain.

#### Verification

Run from `workspaces/repo-agent-*`: `npm test -- app/services/ticket-activity.server.test.ts 'app/routes/tickets.$ticketId.test.tsx' && npm run typecheck && npm run build`

---

### Password Reset Token Service

Task ID: AUTH-01
Category: AFK
Owner: Unassigned
Branch:
Status: Backlog
Dependencies: User Accounts and Authentication merged
Blocking tasks: AUTH-02

#### Objective

Implement secure password reset token persistence and email issuance.

#### Scope

- Add password reset token schema through a new migration.
- Issue single-use password reset tokens.
- Expire password reset tokens after a short configurable window.
- Invalidate prior unused reset tokens when a new token is issued.
- Send reset emails using existing SMTP/email conventions.
- Avoid revealing whether an email address is registered from the service API.
- Add focused service tests for token lifecycle, expiration, invalidation, email issuance, and unknown-email behavior.

#### Out of Scope

- Do not add public forgot/reset screens.
- Do not change login, signup, or email verification behavior.
- Do not add throttling unless already present in auth conventions.

#### Acceptance Criteria

- Password reset service tests pass.
- Existing auth tests continue to pass.
- `npm run db:migrate` succeeds.
- No placeholder TODOs remain.

#### Verification

Run from `workspaces/repo-agent-*`: `npm test -- app/services/auth.server.test.ts app/services/password-reset.server.test.ts && npm run db:migrate && npm run typecheck`

---

### Password Reset Public Routes

Task ID: AUTH-02
Category: AFK
Owner: Unassigned
Branch:
Status: Backlog
Dependencies: AUTH-01
Blocking tasks: None

#### Objective

Allow users to request a password reset email and set a new password with a valid token.

#### Scope

- Provide a public forgot-password screen where users can request a reset email.
- Provide a public reset-password screen for valid tokens.
- Enforce the same password rules used during sign-up.
- Consume valid reset tokens after successful password change.
- Keep responses and UI copy from revealing whether an email address is registered.
- Add focused route tests for request flow, invalid/expired token states, password validation, and successful reset.

#### Out of Scope

- Do not change authentication layout beyond minimal public auth screen integration.
- Do not add account recovery hardening such as throttling.

#### Acceptance Criteria

- Password reset route tests pass.
- Existing auth route tests continue to pass.
- Verification passes.
- No placeholder TODOs remain.

#### Verification

Run from `workspaces/repo-agent-*`: `npm test -- app/routes/auth-screens.test.tsx app/services/password-reset.server.test.ts && npm run typecheck && npm run build`

---

### Virtualized Large Board Rendering

Task ID: BOARD-PERF-01
Category: HITL
Owner: Unassigned
Branch:
Status: Backlog
Dependencies: BOARD-02
Blocking tasks: None

#### Objective

Keep Kanban board rendering responsive when a board contains many tickets.

#### Scope

- Virtualize ticket rendering within board columns if profiling shows the non-virtualized board is insufficient.
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

Run from `workspaces/repo-agent-*`: `npm test -- app/routes/board.test.tsx && npm run typecheck && npm run build`

---

### Remove Placeholder Scaffolding and Starter Copy

Task ID: CLEANUP-01
Category: AFK
Owner: Unassigned
Branch:
Status: Backlog
Dependencies: AUTH-02, EPIC-UI-01, TICKET-WF-01, BOARD-WF-01, COMMENT-02
Blocking tasks: DONE-01

#### Objective

Remove remaining placeholder-only screens, services, tests, and README copy after real product workflows are implemented.

#### Scope

- Remove the placeholder service module and focused tests once no production route depends on it.
- Replace the home screen placeholder with a real authenticated entry point, such as redirecting to the primary board.
- Remove placeholder-only route copy, sample ticket cards, sample team/epic options, and placeholder status payloads from ticket and board routes.
- Remove or rewrite placeholder-focused tests so coverage validates real route behavior.
- Update README introduction text so it describes the product.
- Preserve public authentication route behavior and shared helpers still used by implemented routes.

#### Out of Scope

- Do not add new feature behavior while removing placeholders.
- Do not change business rules.

#### Acceptance Criteria

- No placeholder-only production copy remains.
- README describes the product accurately.
- Full verification passes.
- No placeholder TODOs remain.

#### Verification

Run from `workspaces/repo-agent-*`: `npm test && npm run typecheck && npm run build`

---

### Definition of Done

Task ID: DONE-01
Category: HITL
Owner: Unassigned
Branch:
Status: Backlog
Dependencies: CLEANUP-01
Blocking tasks: None

#### Objective

Verify the complete product against the mandatory acceptance criteria before final release.

#### Scope

- Verify that a user can sign up, receive a verification email through the configured SMTP service, verify the account, and log in.
- Verify that teams and epics can be managed through the UI and persist in the database.
- Verify that a verified user can create, view, edit, and delete tickets.
- Verify that a user can add comments and see their author and timestamp.
- Verify that the Kanban board shows tickets in the correct state columns for the selected team.
- Verify that dragging a ticket to another column updates the server and remains correct after refreshing the page.
- Verify that the application can be started from a clean checkout with `docker compose up --build`.
- Verify that the solution contains no hard-coded user password or committed secret.
- Verify that a fresh database starts with schema and migration metadata only.
- Verify that QA can create required test or demo data through the application UI or API.

#### Out of Scope

- Do not implement feature fixes in this branch.
- Record failed acceptance criteria as new backlog tasks or rejection notes against the responsible implementation task.

#### Acceptance Criteria

- All mandatory acceptance checks are recorded.
- Full Team Lead verification passes: `npm test`, `npm run typecheck`, and `npm run build`.
- Any failed check has a linked backlog task or rejection note.

#### Verification

Run from `workspaces/repo-tl`: `npm test && npm run typecheck && npm run build`

---
