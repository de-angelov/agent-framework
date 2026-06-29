# BACKLOG

Pending work only. Active work lives in TASKS.md. Completed work lives in ARCHIVE.md.

## Backlog

### Wireframe Epic Management UI Alignment

Owner: Unassigned
Branch:
Status: Backlog

Outcome:
Align the epic management screen with the epic wireframe hierarchy.

Scope:
- Shape epic management around team selection, table-style epic list with title/description, ticket count, modified timestamp, row edit/delete actions, disabled delete affordance when referenced, explanatory blocked-delete copy, and adjacent create/edit form placement where practical.
- Preserve existing epic create, edit, delete, title validation, immutable team assignment, and blocked-delete behavior.
- Preserve required loading, empty, success, and error states while applying the wireframe hierarchy.
- Add or update focused tests for team selection, title/description display, ticket counts, modified timestamp display, row edit/delete actions, disabled delete affordance, blocked-delete explanation, create form, and edit form.

Coordination:
- Build after Dev Agent 2's active `Epic Management UI` task has landed in product `main`.
- Reuse the shared authenticated shell if `Wireframe Shared App Shell` has landed first.
- Keep this task scoped to epic route presentation and focused tests.
- Avoid changing backend epic business rules, database schema, team behavior, ticket behavior, or route permissions.
- Do not overlap Dev Agent 2's active `Epic Management UI` branch unless the Team Lead Agent explicitly coordinates the work.

Follow-up:
- If ticket counts are unavailable from existing loaders/services, add a separate backend count-query task before expanding this UI branch.

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

### Wireframe Ticket Details and Comments UI

Owner: Unassigned
Branch:
Status: Backlog

Outcome:
Align ticket create/edit/details and comments presentation with the ticket wireframe hierarchy.

Scope:
- Shape ticket details/editing around the wireframe hierarchy: back-to-team-board navigation, ticket metadata, editable team/type/state/epic/title/body fields, clear save and delete actions, and a comments panel with count, chronological comments, and add-comment form.
- Preserve existing ticket create, view, edit, delete, same-team epic selection, comment creation, and comment ordering behavior.
- Keep ticket deletion behind explicit confirmation.
- Preserve required loading, empty, success, and error states while applying the wireframe hierarchy.
- Add or update focused tests for back navigation, metadata display, editable fields, save/delete actions, confirmation, comments count, chronological comments, and add-comment form.

Coordination:
- Build after `Ticket CRUD UI and Routes` and `Immutable Ticket Comments` have landed in product `main`.
- Reuse the shared authenticated shell if `Wireframe Shared App Shell` has landed first.
- Keep this task scoped to ticket routes and ticket/comment route-level components.
- Avoid changing backend ticket services, comment services, ticket workflow rules, comment mutability, database schema, or board drag-and-drop behavior.

Follow-up:
- Add comment editing/deletion only through the separate `Edit or Delete Own Comments` task.

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

### Wireframe Shared App Shell

Owner: Unassigned
Branch:
Status: Backlog

Outcome:
Create the shared authenticated application shell shown across the wireframes.

Scope:
- Add or refine a consistent authenticated header with `TICKET TRACKER` branding, Board/Teams/Epics navigation, active-route indication, current user email, and a user menu or equivalent logout affordance.
- Keep public authentication screens outside the authenticated shell.
- Ensure the header works consistently on board, ticket, team, and epic screens.
- Add focused tests for visible navigation, active-route indication, current user display, and logout affordance.

Coordination:
- Build after authentication, teams, epics, ticket routes, and the board shell have landed in product `main`.
- Treat the wireframes as low-fidelity hierarchy guidance, not a visual redesign mandate.
- Keep styling restrained and compatible with `TECH.md`; use simple layouts, readable spacing, `padding: 10px`, and `border: 1px solid grey` unless an existing component convention has replaced those defaults.
- Coordinate with `Component Folder and CSS Module Cleanup` if shared component structure or CSS module ownership changes first.
- Avoid changing backend business rules, database schema, authentication behavior, route permissions, or session behavior.
- Do not work in files owned by active implementation branches until those branches land or the Team Lead Agent explicitly coordinates the overlap.

Follow-up:
- Extend the shared shell only if later screens need additional global navigation.

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

### Wireframe Board UI Alignment

Owner: Unassigned
Branch:
Status: Backlog

Outcome:
Align the primary Kanban board screen with the board wireframe hierarchy.

Scope:
- Make the board the primary authenticated work surface with a team selector, prominent new-ticket action, combined search/type/epic filters, clear-filter action, total ticket count, five workflow columns, per-column counts, and compact cards showing type, title, epic when present, and modified recency where practical.
- Preserve existing ticket state ordering, filtering behavior, create-ticket navigation, open-ticket navigation, drag-and-drop persistence, rollback, and error handling.
- Preserve required loading, empty, success, and error states while applying the wireframe hierarchy.
- Add or update focused tests for team selector, new-ticket action, filters, clear action, ticket count, column counts, card content, and drag/drop error visibility where existing drag/drop coverage exposes it.

Coordination:
- Build after `Kanban Board Shell and Ticket Cards`, `Kanban Board Filters`, and `Kanban Drag and Drop Persistence` have landed in product `main`.
- Reuse the shared authenticated shell if `Wireframe Shared App Shell` has landed first.
- Treat the wireframe as low-fidelity hierarchy guidance, not a visual redesign mandate.
- Keep styling restrained and compatible with `TECH.md`.
- Avoid changing backend ticket services, ticket state API behavior, database schema, or drag-and-drop persistence semantics.

Follow-up:
- Record any board usability issue with 100 tickets as a focused follow-up unless it belongs to `Virtualized Large Board Rendering`.

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

### Remove Placeholder Scaffolding and Starter Copy

Owner: Unassigned
Branch:
Status: Backlog

Outcome:
Remove remaining placeholder-only screens, services, tests, and README copy after the real product workflows are implemented.

Scope:
- Remove the placeholder service module and its focused tests once no production route depends on it.
- Replace the home screen placeholder with a real authenticated entry point, such as redirecting to the primary board or rendering a minimal product home that links to live workflows.
- Remove placeholder-only route copy, sample ticket cards, sample team/epic options, and placeholder status payloads from ticket and board routes after those routes connect to real services.
- Remove or rewrite placeholder-focused tests so coverage validates real route behavior instead of scaffold text.
- Update README introduction text so it describes the product rather than a bootstrapped placeholder application.
- Preserve public authentication route behavior and any shared helpers still used by implemented routes.
- Run full verification after cleanup.

Coordination:
- Build after the mandatory authentication, teams, epics, tickets, comments, and Kanban board workflows have landed in product `main`.
- Keep this as a cleanup task only; do not add new feature behavior while removing placeholders.
- Coordinate with `Component Folder and CSS Module Cleanup` if placeholder UI helpers or shared component locations have changed.
- Avoid touching active epic, ticket, comment, and board implementation branches before their real route behavior lands.

Follow-up:
- Treat any remaining placeholder text found during final acceptance as a release blocker or a small follow-up cleanup task.

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
