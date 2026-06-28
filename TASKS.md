# TASKS

This document is the single source of truth for project work.

---

## Backlog

### Docker Compose Runtime

Owner: Unassigned
Branch:
Status: Backlog

Outcome:
Allow the project to run from a clean checkout with `docker compose up --build`.

Scope:
- Add a production-oriented application `Dockerfile`.
- Add a `docker-compose.yml` at the product repository root.
- Ensure `docker compose up --build` installs dependencies, builds the app, and starts the application.
- Expose the application on a documented local port.
- Provide required runtime environment defaults suitable for local development without committing secrets.
- Ensure database initialization or migrations run in a repeatable way for the containerized app.
- Keep generated artifacts, dependency directories, and local database files out of git.
- Add or update documentation for running and stopping the Docker Compose environment.
- Add a focused verification note or script where practical.

Coordination:
- Build on the existing project scripts and database conventions.
- Avoid changing application behavior beyond runtime configuration needed for Docker.
- Keep Docker setup compatible with the final Definition of Done clean-checkout requirement.
- Keep work scoped to runtime configuration, documentation, and verification notes; do not change product behavior beyond container startup requirements.
- Coordinate before changing package scripts, migration commands, or database defaults that overlap with Agent 2's API and persistence expectations work.

Follow-up:
- Add production deployment manifests only if a deployment target is later defined.

---


---

### GitHub Action Merged Branch Cleanup

Owner: Unassigned
Branch:
Status: Backlog

Outcome:
Add a scheduled GitHub Actions workflow that removes remote branches already merged into `main`.

Scope:
- Add a GitHub Actions workflow under `.github/workflows`.
- Run the workflow daily at 00:00 UTC using a cron schedule.
- Fetch remote branches and identify branches already merged into `main`.
- Delete merged remote branches while preserving `main` and any protected or configured keep branches.
- Avoid deleting active agent branches that are not merged.
- Allow manual workflow dispatch for explicit cleanup runs.
- Document the branch cleanup behavior and safety exclusions where practical.

Coordination:
- Keep the workflow scoped to repository branch maintenance.
- Do not change product build, test, deploy, or application runtime behavior.
- Use the built-in GitHub token and least-permission workflow settings where practical.

Follow-up:
- Add notifications only if cleanup failures need visibility later.

---


---

### User Accounts and Authentication

Owner: Unassigned
Branch:
Status: Backlog

Outcome:
Implement local user accounts, email verification, login, logout, and authentication protection for business application surfaces.

Scope:
- Allow users to sign up with email address and password.
- Trim email addresses, compare them case-insensitively, and enforce uniqueness.
- Allow users to log in and log out with local credentials.
- Require passwords to contain at least 8 characters.
- Hash passwords with an established password-hashing algorithm such as Argon2id and never store plain text passwords.
- Send email-verification messages after sign-up through a configurable SMTP service.
- Ensure the SMTP implementation supports relay1.dataart.com.
- Prevent newly registered accounts from using the main application until their email address is verified.
- Make verification links or tokens expire after 24 hours.
- Make verification links or tokens single-use.
- Send successfully verified users to the login screen without automatic login.
- Allow unverified users to request a new verification email from the login or verification-result screen.
- Invalidate earlier unused verification tokens whenever a new token is issued.
- Require authentication for all business application screens and API endpoints.
- Keep sign-up, login, email verification, verification-email resend, static frontend assets, and optional health/readiness endpoints public.
- Add automated tests for account normalization, password rules, password hashing behavior, verification-token lifecycle, resend invalidation, login/logout, and authentication guards.

Coordination:
- Build on the API and persistence conventions from `agent/2/api-persistence-expectations`; if those conventions are not merged yet, keep this branch rebased onto that work or wait for it to land.
- Keep changes focused on authentication services, auth schema/migrations, session handling, email-verification flow, auth route actions/loaders, guards, and focused tests.
- Avoid shared UI component internals owned by Agent 1. Use existing placeholders and minimal styling only.
- Coordinate with `agent/2/minimum-screen-placeholders` for public auth route names and form boundaries before replacing placeholders with real behavior.

Follow-up:
- Add SSO or external identity providers only if later required.

---


---

### Teams

Owner: Unassigned
Branch:
Status: Backlog

Outcome:
Implement team management and enforce team rules for grouping tickets.

Scope:
- Group tickets by team.
- Allow authenticated users to view the list of teams.
- Allow authenticated users to create, rename, and delete teams.
- Store each team with at least an identifier, name, created timestamp, and modified timestamp.
- Require team names to be non-empty after trimming.
- Enforce team-name uniqueness case-insensitively.
- Prevent deleting a team while it contains tickets or epics.
- Show a clear UI validation message when team deletion is blocked.
- Do not cascade-delete tickets or epics when deleting a team.
- Allow all verified users to view and manage all teams.
- Exclude team ownership and membership from the mandatory scope.
- Add automated tests for name normalization, uniqueness, create/rename/delete behavior, blocked deletion, and authenticated access.

Coordination:
- Build on the API, persistence, and authentication conventions from `agent/2/api-persistence-expectations` and `agent/2/user-accounts-authentication`; if those are not merged yet, keep this branch rebased onto that work or wait for them to land.
- Keep backend business rules in services and route modules focused on request validation, authentication boundaries, loading, invoking services, and rendering.
- Keep UI work scoped to the existing team management placeholder route and minimal styling only.
- Avoid shared dialog, header, button internals, and broad styling while Agent 1 UI tasks are active.
- Do not implement epic or ticket CRUD beyond the references and deletion guards required to enforce team deletion rules.

Follow-up:
- Add team ownership and membership only if later required.

---


---

### Epics

Owner: Unassigned
Branch:
Status: Backlog

Outcome:
Implement epic CRUD and enforce team-scoped epic relationships.

Scope:
- Make each epic belong to exactly one team.
- Select the team when an epic is created.
- Prevent changing an epic's team after creation.
- Keep moving epics between teams out of scope.
- Provide a separate epic management screen for creating, listing, editing, and deleting epics.
- Store each epic with at least an identifier, team reference, title, optional description, created timestamp, and modified timestamp.
- Require epic titles to be non-empty after trimming.
- Allow tickets to optionally reference one epic selected from a drop-down list.
- Ensure the ticket epic drop-down lists only epics for the ticket's team.
- Enforce in the backend that a ticket may reference only an epic belonging to the same team as the ticket.
- Prevent deleting an epic while tickets reference it.
- Show a clear UI validation message when epic deletion is blocked.
- Add automated tests for epic creation, listing, editing, deletion, title validation, immutable team assignment, same-team ticket references, and blocked deletion.

Coordination:
- Build on the API, persistence, authentication, and teams conventions from `agent/2/api-persistence-expectations`, `agent/2/user-accounts-authentication`, and `agent/2/teams-management`; if those are not merged yet, keep this branch rebased onto that work or wait for them to land.
- Keep backend business rules in services and route modules focused on request validation, authentication boundaries, loading, invoking services, and rendering.
- Keep UI work scoped to the existing epic management placeholder route and minimal styling only.
- Avoid shared dialog, header, button internals, broad styling, and unrelated ticket CRUD while Agent 1 UI tasks and the ticket backlog task are separate.
- Implement only the ticket-side epic reference validation and drop-down behavior needed to enforce same-team epic selection; leave full ticket CRUD workflow scope to the Tickets task.

Follow-up:
- Add moving epics between teams only if later required.

---


---

### Tickets

Owner: Unassigned
Branch:
Status: Backlog

Outcome:
Implement ticket CRUD, validation, workflow state persistence, and team/epic relationships.

Scope:
- Store each ticket with a stable unique system-generated identifier.
- Require each ticket to reference an existing team.
- Support exactly these ticket type values: bug, feature, fix.
- Treat ticket type as a classification label only with no workflow differences.
- Support exactly these ticket state values: new, ready_for_implementation, in_progress, ready_for_acceptance, done.
- Display ticket states in the UI with human-readable labels using spaces.
- Keep the workflow fixed with no custom states.
- Allow each ticket to optionally reference one epic.
- Ensure a ticket's epic is null or references an epic from the same team as the ticket.
- Require ticket title to be non-empty after trimming.
- Require ticket body to be non-empty.
- Support plain text or Markdown ticket body content without requiring rich-text editing.
- Set created-at on the server in UTC when the ticket is created.
- Set modified-at on the server in UTC whenever ticket fields or state change.
- Do not update modified-at when adding a comment.
- Set created-by automatically from the authenticated user.
- Allow authenticated users to create tickets.
- Allow users to open tickets and view all fields, including created-by, created-at, and modified-at.
- Allow editing ticket type, team, epic, title, body, and state.
- Do not advance modified-at when saving unchanged ticket values.
- Clear or replace any selected epic in the UI when a ticket's team changes.
- Reject in the backend any ticket whose epic belongs to a different team.
- Delete tickets only after explicit confirmation.
- Delete a ticket's comments when the ticket is deleted.
- Persist drag-and-drop state changes immediately in the database.
- Validate all submitted enum values and references in the backend.
- Add automated tests for required fields, enum validation, reference validation, same-team epic constraints, create/view/edit/delete behavior, timestamp behavior, created-by assignment, unchanged saves, and drag-and-drop state persistence.

Coordination:
- Build on the API, persistence, authentication, teams, and epics conventions from `agent/2/api-persistence-expectations`, `agent/2/user-accounts-authentication`, `agent/2/teams-management`, and `agent/2/epics`; if those are not merged yet, keep this branch rebased onto the needed work or wait for them to land.
- Keep backend business rules in reusable ticket services and keep route modules focused on request validation, authentication boundaries, loading, invoking services, and rendering.
- Keep UI work scoped to the existing ticket create, edit, and details placeholder routes with minimal styling only.
- Avoid shared dialog, header, button internals, broad styling, and the dedicated Kanban board UI while Agent 1 UI tasks and the Kanban Board backlog task are separate.
- Implement only the ticket state persistence endpoint behavior needed for ticket CRUD and later board drag-and-drop integration; leave full board filtering, card rendering, rollback UI, and large-board behavior to the Kanban Board task.
- Do not implement comment creation or display beyond preserving deletion behavior for existing ticket comments if the comments schema exists.

Follow-up:
- Add custom workflows or ticket-type-specific behavior only if later required.

---


---

### Comments

Owner: Unassigned
Branch:
Status: Backlog

Outcome:
Implement immutable ticket comments for authenticated users.

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
- Build on the API, persistence, authentication, teams, epics, and tickets conventions from `agent/2/api-persistence-expectations`, `agent/2/user-accounts-authentication`, `agent/2/teams-management`, `agent/2/epics`, and `agent/2/tickets`; if those are not merged yet, keep this branch rebased onto the needed work or wait for them to land.
- Keep backend business rules in reusable comment services and keep route modules focused on request validation, authentication boundaries, loading, invoking services, and rendering.
- Keep UI work scoped to the existing ticket details route with minimal styling only.
- Avoid shared dialog, header, button internals, broad styling, and the dedicated Kanban board UI while Agent 1 UI tasks and the Kanban Board backlog task are separate.
- Do not implement comment editing, deletion, moderation, or ticket activity history beyond append-only comment creation and chronological display.

Follow-up:
- Add comment editing and deletion only as stretch features if later required.

---


---

### Kanban Board

Owner: Unassigned
Branch:
Status: Backlog

Outcome:
Implement the primary team Kanban board with drag-and-drop state changes, filtering, and ticket navigation.

Scope:
- Make the primary screen a Kanban board for one selected team.
- Render exactly five columns, one for each ticket state, in workflow order.
- Display each ticket as a card showing at least title and type.
- Show the ticket's epic on the card where practical.
- Allow users to drag a ticket card from one column to another.
- Persist dropped card state changes immediately through the backend API.
- Return a card to its previous column when a drag-and-drop update fails.
- Display a clear UI error when a drag-and-drop update fails.
- Allow cards to move directly between any two states without enforcing sequential transitions.
- Order cards within each column by most recently modified first.
- Do not persist a custom manual order.
- Provide a clear way to create a ticket from the board.
- Provide a clear way to open an existing ticket from the board.
- Provide filtering by ticket type.
- Provide filtering by epic.
- Provide case-insensitive substring search over ticket title.
- Combine active filters using AND logic.
- Implement filters client-side or server-side.
- Keep the interface usable with at least 100 tickets on one team board.
- Add automated tests for column ordering, card rendering, drag-and-drop persistence, failed drag rollback, filters, create/open affordances, and 100-ticket usability where practical.

Coordination:
- Build on merged authentication, teams, epics, and tickets service/API behavior from Agent 2; if those branches are not merged yet, rebase onto the needed work or pause before changing shared backend contracts.
- Keep board work focused on the Kanban route, board interactions, filtering, ticket navigation affordances, and focused tests.
- Use existing ticket update APIs for drag-and-drop persistence when available; avoid creating a parallel ticket state API unless the merged ticket workflow lacks a usable endpoint.
- Avoid shared header internals owned by `agent/1/authenticated-header-navigation`; integrate with the header only after that branch has landed or through a narrow route-level hook point.
- Keep styling minimal and limited to simple flexbox layouts, `padding: 10px`, and `border: 1px solid grey`.

Follow-up:
- Add custom manual ordering or sequential transition enforcement only if later required.

---


---

### API and Persistence Expectations

Owner: Unassigned
Branch:
Status: Backlog

Outcome:
Establish API and persistence behavior shared by all application workflows.

Scope:
- Perform all create, update, and delete operations through the backend API.
- Persist all create, update, and delete operations in the RDBMS.
- Do not rely on browser local storage as the system of record.
- Use database constraints and/or server-side validation to maintain referential integrity.
- Return meaningful HTTP status codes and error messages for validation failures.
- Return meaningful HTTP status codes and error messages for authentication failures.
- Return meaningful HTTP status codes and error messages for missing records.
- Return meaningful HTTP status codes and error messages for conflicts.
- Return HTTP 409 Conflict when deleting a team that contains tickets or epics.
- Return HTTP 409 Conflict when deleting an epic referenced by tickets.
- Use UUIDs or database-generated numeric values for identifiers.
- Represent API timestamps as ISO-8601 UTC values.
- Use cookie-based sessions or bearer-token authentication.
- Never place session identifiers, access tokens, or bearer tokens in URLs.
- Allow single-use email-verification tokens in verification URLs.
- Use last successful write wins; concurrent-edit conflict detection is not required.
- Automate database schema creation through migrations or an equivalent repeatable initialization mechanism.
- Ensure a fresh database contains no application users, teams, epics, tickets, or comments after migrations or initialization.
- Allow migration metadata in a fresh database.
- Do not load sample or seed data in the default startup path.
- Ensure QA can create test data through the application UI or API.
- Add automated tests for API error mapping, referential integrity behavior, timestamp serialization, auth-token URL handling where practical, and fresh database initialization.

Coordination:
- Keep changes focused on backend API conventions, persistence infrastructure, schema/migration initialization, service-level error mapping, and focused tests.
- Avoid `app/styles.css`, broad placeholder route shells, and shared UI component internals while Agent 1 UI tasks are active.
- Build on the completed `agent/2/minimum-screen-placeholders` route boundaries without reworking placeholder screens outside API and persistence needs.

Follow-up:
- Add seed data or demo fixtures only behind an explicit non-default development command if later required.

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
- Build on the ticket schema, ticket services, authentication, teams, and epics behavior from `agent/2/tickets`; if that branch is not merged yet, keep this branch rebased onto that work or wait for it to land.
- Keep activity writes in backend services that already handle ticket creation and ticket updates so route modules remain thin.
- Keep activity entries append-only; do not add edit or delete workflows for activity history.
- Keep UI work scoped to the existing ticket details route with minimal styling only.
- Avoid shared dialog, header, button internals, broad styling, and the dedicated Kanban board UI while Agent 1 UI tasks are active.
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
- Build on the comment schema, comment services, authentication, and ticket details route behavior from `agent/2/comments`; if that branch is not merged yet, keep this branch rebased onto that work or wait for it to land.
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
- Avoid shared UI component internals owned by Agent 1. Use existing auth placeholders and minimal styling only.
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
- Build on the completed Kanban board route and component structure; if `agent/1/kanban-board` has not landed in product `main`, wait for it or rebase onto it before changing board internals.
- Keep work focused on board column/card rendering performance, visible-card interactions, and focused tests or profiling notes.
- Do not change backend pagination or ticket APIs unless profiling shows client-side virtualization is insufficient; record that as a follow-up before expanding backend scope.
- Avoid shared header, dialog, team, epic, ticket service, and comment internals owned by active Agent 1 and Agent 2 tasks.
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
- Run full Team Lead verification before hand-off: `npm test`, `npm run typecheck`, and `npm run build`.
- Verify default styling remains limited to simple flexbox layouts, `padding: 10px`, and `border: 1px solid grey`.

Follow-up:
- Record any failed acceptance criterion as a new backlog task or rejection against the responsible implementation task.

---


---

## Agent 1 In Progress

## Agent 2 In Progress

## Done

### Authenticated Header Navigation

Owner: Agent 2
Branch: agent/2/authenticated-header-navigation
Status: Done
Completed: 2026-06-28

Outcome:
Add a reusable header for authenticated users that shows the logged-in user's email and links to main pages.

Scope:
- Create a shared authenticated header component.
- Display the logged-in user's email.
- Provide navigation links to the Kanban board, team management, epic management, and other implemented business pages.
- Include a clear logout affordance when authentication is implemented.
- Hide or replace the authenticated header on public authentication screens where appropriate.
- Ensure the header works in light and dark modes.
- Add focused component or route coverage where practical.

Coordination:
- Keep implementation focused on the reusable header component, narrow layout integration, and focused tests.
- Avoid broad route placeholder changes; `agent/2/minimum-screen-placeholders` is complete and route placeholder follow-up work should remain separate.
- Avoid overlapping dialog component internals from `agent/1/reusable-dialog-component`; complete or rebase against that work before merging if shared UI exports or styles overlap.

Progress:
- Added a reusable authenticated header component with business navigation, logged-in email display, and a logout form affordance.
- Integrated the header through the authenticated placeholder shell for board, team, epic, and ticket surfaces.
- Added a public screen shell for sign-up, login, and email verification screens so authenticated navigation stays off public auth routes.
- Kept styling to existing minimal flexbox, padding, and grey border conventions with inherited light and dark mode colors.
- Added focused component coverage and route coverage for authenticated and public shell behavior.
- Verified on product `main` with `npm test`, `npm run typecheck`, and `npm run build`.
- Pushed branch `agent/2/authenticated-header-navigation` at commit `f8d64c6`.
- Merged into product `main` and pushed commit `c3bb3a3`.

Follow-up:
- Add role-aware or team-aware navigation only if later required.

---


---

### Reusable Table Component

Owner: Agent 1
Branch: agent/1/reusable-table-component
Status: Done
Completed: 2026-06-28

Outcome:
Add a reusable table component for list and management screens.

Scope:
- Create a shared table component with configurable column headers.
- Render rows from structured data or row render callbacks.
- Allow callers to inject custom components into row cells for actions, status, links, or inline controls.
- Support empty, loading, and error states where practical.
- Ensure table markup and focus behavior are accessible.
- Ensure styling works in light and dark modes using minimal existing conventions.
- Add focused component coverage where practical.

Coordination:
- Keep changes scoped to shared table component files and focused component tests.
- Avoid broad route rewrites; integrate the table into one narrow placeholder or management screen only if it stays low risk.
- Build on the completed reusable dialog component without coupling the two components.
- Coordinate any shared UI barrel exports, test setup, or minimal style updates with `agent/2/authenticated-header-navigation` before merging either task.
- Keep styling minimal and limited to simple flexbox layouts, `padding: 10px`, and `border: 1px solid grey`.

Progress:
- Added a generic reusable table component with column headers, structured row cell renderers, and optional custom row rendering.
- Added support for caller-injected cell components, plus empty, loading, and error state rows.
- Used semantic table markup with column header scopes, optional captions, busy state, and alert markup for errors.
- Added minimal table styling that inherits current light and dark mode colors.
- Added focused component coverage for structured rows, custom rows, injected controls, and state rendering.
- Verified with `npm test`, `npm run typecheck`, and `npm run build`.
- Pushed branch `agent/1/reusable-table-component` at commit `6324807`.
- Merged into product `main` and pushed commit `d1cdd24`.

Follow-up:
- Add sorting, filtering, pagination, or virtualization only when real screens require them.

---


---

### Reusable Dialog Component

Owner: Agent 1
Branch: agent/1/reusable-dialog-component
Status: Done
Completed: 2026-06-28

Outcome:
Add a reusable dialog component for confirmations and modal workflows.

Scope:
- Create a shared dialog component with title, body, and action slots.
- Support confirmation and cancellation actions.
- Ensure keyboard and focus behavior is accessible.
- Ensure dialog styling works in light and dark modes.
- Use the dialog for at least one confirmation flow where practical.
- Add focused component coverage where practical.

Coordination:
- Keep changes scoped to shared UI component files, focused component tests, and one narrow confirmation-flow integration where practical.
- Avoid broad route styling changes; `agent/2/minimum-screen-placeholders` is complete and route placeholder follow-up work should remain separate.

Progress:
- Added a reusable dialog component using the native dialog API with title, body, cancel action, and confirm action slots.
- Added focus entry and restoration behavior for modal open and close transitions.
- Added light and dark mode dialog styling using the existing minimal style conventions.
- Integrated the dialog into the home placeholder as a narrow confirmation flow while preserving placeholder navigation on `main`.
- Added focused component coverage for accessible dialog markup and controlled modal state.
- Verified on product `main` with `npm test`, `npm run typecheck`, and `npm run build`.
- Pushed branch `agent/1/reusable-dialog-component` at commit `1b36e34`.
- Merged into product `main` and pushed commit `b4f9f54`.

Follow-up:
- Add specialized dialog variants only when real workflows require them.

---


---

### Minimum Screen Placeholders

Owner: Agent 2
Branch: agent/2/minimum-screen-placeholders
Status: Done
Completed: 2026-06-28

Outcome:
Add placeholder routes and screen shells for the minimum product surface.

Scope:
- Add a sign-up screen placeholder with form fields and submit affordance.
- Add an email verification result screen placeholder covering success, invalid-token, and expired-token states.
- Add a verification-email resend action placeholder for unverified accounts and expired-token cases.
- Add a login screen placeholder with form fields and submit affordance.
- Add a Kanban board placeholder with team selector, columns, and placeholder ticket cards.
- Add ticket create, edit, and details placeholders.
- Add a team management screen placeholder with list and create/edit affordances.
- Add an epic management screen placeholder with list and create/edit affordances.
- Keep route loaders/actions thin and ready for later service integration.
- Add focused smoke coverage for the placeholder routes where practical.

Progress:
- Added placeholder routes for authentication, email verification, verification resend, board, ticket, team, and epic screens.
- Kept loaders and actions as thin placeholder boundaries for later service integration.
- Added focused smoke coverage for placeholder route rendering and boundary return values.
- Resolved the `app/styles.css` merge conflict against current product `main` while preserving the reusable button styles.
- Verified on the feature branch and merged `main` with `npm test`, `npm run typecheck`, and `npm run build`.
- Pushed branch `agent/2/minimum-screen-placeholders` at commit `c6d30e3`.
- Merged into product `main` and pushed commit `0883c8e`.

Follow-up:
- Implement real account creation, verification, authentication, team, epic, ticket, and board workflows.

---


---

### Reusable Button and Theme Modes

Owner: Agent 1
Branch: agent/1/reusable-button-theme-modes
Status: Done
Completed: 2026-06-28

Outcome:
Add a reusable button component that works cleanly in light and dark modes.

Scope:
- Create a shared button component for common actions.
- Support primary, secondary, and destructive visual variants where practical.
- Support disabled and loading states.
- Ensure the component has accessible focus states.
- Add light mode and dark mode styling support.
- Use the shared button in placeholder or existing screens where practical without broad unrelated refactors.
- Add focused component coverage where practical.

Coordination:
- Start after Initial Project Setup establishes the React Router app structure.
- Keep replacements narrow to avoid colliding with Agent 2 placeholder route work.

Progress:
- Added a reusable button component with primary, secondary, and destructive variants.
- Added disabled and loading state support with accessible busy and focus behavior.
- Added light and dark mode button styling.
- Used the shared button on the placeholder home route.
- Added focused component coverage.
- Verified with `npm test`, `npm run typecheck`, and `npm run build`.
- Pushed branch `agent/1/reusable-button-theme-modes` at commit `cff8265`.
- Merged into product `main` and pushed commit `ec4feb7`.

Follow-up:
- Expand the component system only as repeated UI needs emerge.

---


---

### Initial Project Setup

Owner: Agent 1
Branch: agent/1/initial-project-setup
Status: Done
Completed: 2026-06-28

Outcome:
Bootstrap the application as a minimal React Router framework-mode repo ready for frontend, backend, and persistence work.

Scope:
- Initialize the app with React Router framework mode using the standard React Router init flow.
- Keep the generated application minimal by removing unnecessary starter/demo files, assets, routes, and sample content.
- Add a simple placeholder frontend route that confirms the app renders.
- Add a minimal backend/service layer placeholder for server-side business logic.
- Add a minimal database layer placeholder using SQLite and Drizzle conventions.
- Configure TypeScript strict mode and expected project scripts for development, typecheck, tests, build, and database migration.
- Add focused smoke coverage for the placeholder app structure where practical.

Progress:
- Added a minimal React Router framework-mode application structure.
- Added placeholder service and SQLite/Drizzle database scaffolding.
- Added strict TypeScript, expected npm scripts, and smoke tests.
- Verified with `npm test`, `npm run typecheck`, and `npm run build`.
- Merged into product `main` and pushed commit `db59f62`.

Follow-up:
- Replace placeholders with the first real product workflow.
