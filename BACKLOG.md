# BACKLOG

Pending work only. Active work lives in TASKS.md. Completed work lives in ARCHIVE.md.

## Backlog

### Exhaustive Mapper Pattern Cleanup

Owner: Unassigned
Branch:
Status: Backlog

Outcome:
Make app-owned error and status mapping functions consistently use exhaustive `ts-pattern` matches.

Scope:
- Convert service error mappers such as team and epic mutation error mapping from `switch` statements to `match(...).with(...).exhaustive()`.
- Check existing app-owned union or enum mapping functions for the same cleanup opportunity.
- Preserve all existing messages, return values, and behavior.
- Keep simple guard clauses, null checks, and untrusted external input handling unchanged.
- Add or update focused tests only if an affected mapper lacks coverage.

Coordination:
- Keep this as a narrow consistency cleanup in service and mapper code only.
- Use `.exhaustive()` only for typed unions or enums owned by the application, following `TECH.md`.
- Do not change route behavior, UI styling, component structure, database schema, or feature workflows.
- Avoid broad refactors while Dev Agent 2's `Epic Management UI` task is active.

Follow-up:
- Apply this convention as new app-owned error/status mappers are added.

---


### Component Folder and CSS Module Cleanup

Owner: Unassigned
Branch:
Status: Backlog

Outcome:
Move shared UI components into per-component folders with colocated tests and CSS modules.

Scope:
- Create a dedicated subfolder for each existing shared component under `app/components`.
- Move each component implementation into its own folder.
- Move each component's focused test into the same folder as the component.
- Move component-specific CSS out of global stylesheets into colocated `*.module.css` files.
- Update imports, route usage, and test references after files move.
- Keep `app/styles.css` limited to global resets, document defaults, and intentionally shared application-level primitives.
- Preserve existing component behavior and visual simplicity.
- Add or update automated tests only where moves expose missing component coverage.

Coordination:
- Build after Dev Agent 2 completes `Epic Management UI`, or coordinate carefully if that branch changes shared component imports or global styles.
- Keep this as a structural cleanup only; do not redesign components, change route behavior, or add visual polish.
- Use the component organization and CSS module conventions in `TECH.md`.
- Avoid ticket, comment, board, team, and epic feature behavior changes.

Follow-up:
- Apply the same folder-and-module convention to future route-specific components as they are extracted.

---


---

### Ticket Data Model and Services

Owner: Unassigned
Branch:
Status: Backlog

Outcome:
Implement backend ticket persistence, validation, timestamps, and team/epic relationship rules.

Scope:
- Store each ticket with a stable unique system-generated identifier.
- Require each ticket to reference an existing team.
- Support exactly these ticket type values: bug, feature, fix.
- Treat ticket type as a classification label only with no workflow differences.
- Support exactly these ticket state values: new, ready_for_implementation, in_progress, ready_for_acceptance, done.
- Keep the workflow fixed with no custom states.
- Allow each ticket to optionally reference one epic.
- Ensure a ticket's epic is null or references an epic from the same team as the ticket.
- Require ticket title to be non-empty after trimming.
- Require ticket body to be non-empty.
- Support plain text or Markdown ticket body content without requiring rich-text editing.
- Set created-at on the server in UTC when the ticket is created.
- Set modified-at on the server in UTC whenever ticket fields or state change.
- Set created-by automatically from the authenticated user.
- Do not advance modified-at when saving unchanged ticket values.
- Reject in the backend any ticket whose epic belongs to a different team.
- Validate all submitted enum values and references in the backend.
- Add automated tests for required fields, enum validation, reference validation, same-team epic constraints, timestamp behavior, created-by assignment, and unchanged saves.

Coordination:
- Build on merged API, persistence, authentication, teams, and epics behavior.
- Keep backend business rules in reusable ticket services and keep route modules focused on request validation, authentication boundaries, loading, invoking services, and rendering.
- Do not implement ticket comments, board drag-and-drop UI, activity history, or broad styling.

Follow-up:
- Add custom workflows or ticket-type-specific behavior only if later required.

---


---

### Ticket CRUD UI and Routes

Owner: Unassigned
Branch:
Status: Backlog

Outcome:
Allow authenticated users to create, view, edit, and delete tickets through the existing ticket route placeholders.

Scope:
- Allow authenticated users to create tickets.
- Allow users to open tickets and view all fields, including created-by, created-at, and modified-at.
- Allow editing ticket type, team, epic, title, body, and state.
- Display ticket states in the UI with human-readable labels using spaces.
- Ensure the ticket epic drop-down lists only epics for the ticket's team.
- Clear or replace any selected epic in the UI when a ticket's team changes.
- Delete tickets only after explicit confirmation.
- Delete a ticket's comments when the ticket is deleted if comments already exist.
- Add automated tests for create/view/edit/delete behavior and same-team epic selection behavior.

Coordination:
- Build after `Ticket Data Model and Services`.
- Keep UI work scoped to the existing ticket create, edit, and details placeholder routes with minimal styling only.
- Avoid shared dialog, header, button internals, broad styling, comments UI, activity history, and the dedicated Kanban board UI.
- Do not implement comment creation or display beyond preserving deletion behavior for existing ticket comments if the comments schema exists.

Follow-up:
- Add rich-text editing only if later required.

---


---

### Ticket State Update API

Owner: Unassigned
Branch:
Status: Backlog

Outcome:
Provide a focused backend path for persistent ticket state updates used by ticket editing and the Kanban board.

Scope:
- Persist ticket state changes immediately in the database.
- Allow cards or forms to move tickets directly between any two states without enforcing sequential transitions.
- Reuse the ticket update API when practical.
- Return meaningful validation, missing-record, authentication, and conflict responses.
- Do not persist a custom manual order.
- Add automated tests for drag-and-drop-oriented state persistence and enum validation.

Coordination:
- Build after `Ticket Data Model and Services` and before the Kanban drag-and-drop task.
- Avoid creating a parallel ticket state API unless the merged ticket workflow lacks a usable endpoint.
- Keep work backend-focused and do not implement board UI.

Follow-up:
- Add sequential transition enforcement only if later required.

---


---

### Immutable Ticket Comments

Owner: Unassigned
Branch:
Status: Backlog

Outcome:
Implement append-only ticket comments for authenticated users.

Scope:
- Allow authenticated users to add comments to a ticket.
- Store each comment with an identifier, ticket reference, author, body, and created timestamp.
- Require comment bodies to be non-empty.
- Display comments chronologically with the oldest comment first.
- Do not update the ticket modified timestamp when adding a comment.
- Ensure adding comments does not change ticket board ordering.
- Keep comments immutable after creation for the mandatory scope.
- Add automated tests for authenticated access, required body validation, author assignment, chronological ordering, and ticket modified timestamp behavior.

Coordination:
- Build after ticket CRUD is merged.
- Keep backend business rules in reusable comment services and keep route modules focused on request validation, authentication boundaries, loading, invoking services, and rendering.
- Keep UI work scoped to the existing ticket details route with minimal styling only.
- Avoid shared dialog, header, button internals, broad styling, and the dedicated Kanban board UI while Dev Agent 1 UI tasks and the Kanban Board backlog task are separate.
- Do not implement comment editing, deletion, moderation, or ticket activity history beyond append-only comment creation and chronological display.

Follow-up:
- Add comment editing and deletion only as stretch features if later required.

---


---

### Kanban Board Shell and Ticket Cards

Owner: Unassigned
Branch:
Status: Backlog

Outcome:
Implement the primary team Kanban board layout, ticket card rendering, ordering, and ticket navigation.

Scope:
- Make the primary screen a Kanban board for one selected team.
- Render exactly five columns, one for each ticket state, in workflow order.
- Display each ticket as a card showing at least title and type.
- Show the ticket's epic on the card where practical.
- Order cards within each column by most recently modified first.
- Provide a clear way to create a ticket from the board.
- Provide a clear way to open an existing ticket from the board.
- Keep the interface usable with at least 100 tickets on one team board where practical without virtualization.
- Add automated tests for column ordering, card rendering, create/open affordances, and 100-ticket usability where practical.

Coordination:
- Build on merged authentication, teams, epics, tickets, and ticket state API behavior.
- Keep board work focused on the Kanban route, card rendering, ticket navigation affordances, and focused tests.
- Avoid drag-and-drop behavior, filters, shared header internals, ticket service rewrites, and broad styling.
- Keep styling minimal and limited to simple flexbox layouts, `padding: 10px`, and `border: 1px solid grey`.

Follow-up:
- Add manual ordering only if later required.

---


---

### Kanban Drag and Drop Persistence

Owner: Unassigned
Branch:
Status: Backlog

Outcome:
Allow users to move tickets between Kanban columns with immediate persisted state updates and rollback on failure.

Scope:
- Allow users to drag a ticket card from one column to another.
- Persist dropped card state changes immediately through the backend API.
- Return a card to its previous column when a drag-and-drop update fails.
- Display a clear UI error when a drag-and-drop update fails.
- Allow cards to move directly between any two states without enforcing sequential transitions.
- Do not persist a custom manual order.
- Add automated tests for drag-and-drop persistence and failed drag rollback.

Coordination:
- Build after `Kanban Board Shell and Ticket Cards` and `Ticket State Update API`.
- Use existing ticket update APIs for drag-and-drop persistence when available.
- Keep work focused on board interactions, rollback behavior, and focused tests.

Follow-up:
- Add keyboard drag alternatives only if accessibility testing shows the selected library requires extra work.

---


---

### Kanban Board Filters

Owner: Unassigned
Branch:
Status: Backlog

Outcome:
Add client-side or server-side filtering to the Kanban board.

Scope:
- Provide filtering by ticket type.
- Provide filtering by epic.
- Provide case-insensitive substring search over ticket title.
- Combine active filters using AND logic.
- Keep filtered column ordering by most recently modified first.
- Add automated tests for ticket type, epic, title search, and combined filters.

Coordination:
- Build after `Kanban Board Shell and Ticket Cards`.
- Keep work focused on board filtering controls and filtering behavior.
- Avoid drag-and-drop internals, ticket CRUD routes, backend pagination, and broad styling.

Follow-up:
- Add saved filters only if later required.

---


---

### Ticket Activity History

Owner: Unassigned
Branch:
Status: Backlog

Outcome:
Record and display meaningful ticket changes over time.

Scope:
- Record ticket creation, status changes, title changes, description changes, team changes, epic changes, and assignment changes where supported.
- Store activity entries with actor, timestamp, action type, and relevant changed values.
- Display ticket activity history on the ticket details screen.
- Keep activity entries append-only.
- Avoid recording sensitive authentication data in activity history.
- Add automated tests for activity creation and display ordering.

Coordination:
- Build on the ticket schema, ticket services, authentication, teams, and epics behavior from `Ticket Data Model and Services`; if that task is not merged yet, wait for it to land.
- Keep activity writes in backend services that already handle ticket creation and ticket updates so route modules remain thin.
- Keep activity entries append-only; do not add edit or delete workflows for activity history.
- Keep UI work scoped to the existing ticket details route with minimal styling only.
- Avoid shared dialog, header, button internals, broad styling, and the dedicated Kanban board UI while Dev Agent 1 UI tasks are active.
- Do not record sensitive authentication data, session identifiers, verification tokens, password-reset tokens, or password-related values in activity history.

Follow-up:
- Add filtering or export only if activity volume makes it necessary.

---


---

### Edit or Delete Own Comments

Owner: Unassigned
Branch:
Status: Backlog

Outcome:
Allow authenticated users to modify or remove comments they created.

Scope:
- Allow users to edit their own ticket comments.
- Allow users to delete their own ticket comments.
- Prevent users from editing or deleting comments created by other users.
- Preserve created and modified timestamps.
- Show clear UI affordances only where the current user owns the comment.
- Add automated tests for ownership checks, edit behavior, delete behavior, and unauthorized attempts.

Coordination:
- Build on the comment schema, comment services, authentication, and ticket details route behavior from `Immutable Ticket Comments`; if that task is not merged yet, wait for it to land.
- Keep backend ownership checks in reusable comment services and keep route modules focused on request validation, authentication boundaries, loading, invoking services, and rendering.
- Keep UI work scoped to the existing ticket details route with minimal styling only.
- Avoid shared dialog, header, button internals, broad styling, moderation controls, and ticket activity history while those concerns remain separate.
- Preserve ticket modified timestamp behavior unless the existing comments model explicitly defines comment edits as ticket changes.

Follow-up:
- Add moderator or admin comment controls only if later required.

---


---

### Password Reset Flow

Owner: Unassigned
Branch:
Status: Backlog

Outcome:
Allow users to reset forgotten passwords securely.

Scope:
- Provide a public forgot-password screen where users can request a reset email.
- Issue single-use password reset tokens.
- Expire password reset tokens after a short configurable window.
- Allow users with valid reset tokens to set a new password.
- Enforce the same password rules used during sign-up.
- Invalidate prior unused reset tokens when a new token is issued.
- Avoid revealing whether an email address is registered.
- Add automated tests for token lifecycle, expiration, invalidation, password validation, and successful reset behavior.

Coordination:
- Build on the authentication schema, password validation, password hashing, SMTP/email conventions, and public auth route boundaries from `agent/2/user-accounts-authentication`; if that branch is not merged yet, keep this branch rebased onto that work or wait for it to land.
- Keep changes focused on password-reset services, token persistence, reset email delivery, public forgot/reset route actions/loaders, and focused tests.
- Do not reveal account existence in UI copy, API responses, logs, or email-request behavior.
- Avoid shared UI component internals owned by Dev Agent 1. Use existing auth placeholders and minimal styling only.
- Keep styling minimal and limited to simple flexbox layouts, `padding: 10px`, and `border: 1px solid grey`.

Follow-up:
- Add account recovery hardening such as throttling only if abuse becomes a concern.

---


---

### Virtualized Large Board Rendering

Owner: Unassigned
Branch:
Status: Backlog

Outcome:
Keep Kanban board rendering responsive when a board contains many tickets.

Scope:
- Virtualize ticket rendering within board columns.
- Preserve keyboard and pointer interactions for visible ticket cards.
- Keep column headers and empty states stable while lists virtualize.
- Avoid changing backend pagination unless needed for frontend performance.
- Add focused tests or profiling notes that demonstrate large-board rendering remains usable.

Coordination:
- Build on the completed Kanban board route and component structure; if `Kanban Board Shell and Ticket Cards` has not landed in product `main`, wait for it before changing board internals.
- Keep work focused on board column/card rendering performance, visible-card interactions, and focused tests or profiling notes.
- Do not change backend pagination or ticket APIs unless profiling shows client-side virtualization is insufficient; record that as a follow-up before expanding backend scope.
- Avoid shared header, dialog, team, epic, ticket service, and comment internals owned by active Dev Agent 1 and Dev Agent 2 tasks.
- Keep styling minimal and limited to simple flexbox layouts, `padding: 10px`, and `border: 1px solid grey`.

Follow-up:
- Add server-side pagination or infinite loading only if client-side virtualization is insufficient.

---


---

### Definition of Done

Owner: Unassigned
Branch:
Status: Backlog

Outcome:
Verify the complete product against the mandatory acceptance criteria before final release.

Scope:
- [ ] A user can sign up, receive a verification email through the configured SMTP service, verify the account, and log in.
- [ ] Teams and epics can be managed through the UI and persist in the database.
- [ ] A verified user can create, view, edit, and delete tickets.
- [ ] A user can add comments and see their author and timestamp.
- [ ] The Kanban board shows tickets in the correct state columns for the selected team.
- [ ] Dragging a ticket to another column updates the server and remains correct after refreshing the page.
- [ ] The application can be started from a clean checkout with `docker compose up --build` from the repository root.
- [ ] The solution contains no hard-coded user password or committed secret.
- [ ] A fresh database starts with schema and migration metadata only; no application data is preloaded.
- [ ] QA can create all required test or demo data through the application UI or API without manually changing database records.

Coordination:
- Start final acceptance verification only after the mandatory authentication, teams, epics, tickets, comments, and Kanban board branches have landed in product `main`.
- Keep this task focused on verification, release-readiness notes, and focused test or documentation updates if required by the verification process.
- Do not implement feature fixes in this branch; record failed acceptance criteria as new backlog tasks or rejection notes against the responsible implementation task.
- Run full Team Lead Agent verification before hand-off: `npm test`, `npm run typecheck`, and `npm run build`.
- Verify default styling remains limited to simple flexbox layouts, `padding: 10px`, and `border: 1px solid grey`.

Follow-up:
- Record any failed acceptance criterion as a new backlog task or rejection against the responsible implementation task.

---


---

