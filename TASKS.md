# TASKS

This document is the single source of truth for project work.

---

## Backlog

### Reusable Button and Theme Modes

Owner: Unassigned
Branch:
Status: Backlog

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

Follow-up:
- Expand the component system only as repeated UI needs emerge.

---

### Reusable Dialog Component

Owner: Unassigned
Branch:
Status: Backlog

Outcome:
Add a reusable dialog component for confirmations and modal workflows.

Scope:
- Create a shared dialog component with title, body, and action slots.
- Support confirmation and cancellation actions.
- Ensure keyboard and focus behavior is accessible.
- Ensure dialog styling works in light and dark modes.
- Use the dialog for at least one confirmation flow where practical.
- Add focused component coverage where practical.

Follow-up:
- Add specialized dialog variants only when real workflows require them.

---

### Authenticated Header Navigation

Owner: Unassigned
Branch:
Status: Backlog

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

Follow-up:
- Add role-aware or team-aware navigation only if later required.

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

Follow-up:
- Add seed data or demo fixtures only behind an explicit non-default development command if later required.

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

Follow-up:
- Add SSO or external identity providers only if later required.

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

Follow-up:
- Add team ownership and membership only if later required.

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

Follow-up:
- Add moving epics between teams only if later required.

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

Follow-up:
- Add custom workflows or ticket-type-specific behavior only if later required.

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

Follow-up:
- Add comment editing and deletion only as stretch features if later required.

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

Follow-up:
- Add custom manual ordering or sequential transition enforcement only if later required.

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

Follow-up:
- Record any failed acceptance criterion as a new backlog task or rejection against the responsible implementation task.

---

## Agent 1 In Progress

### Initial Project Setup

Owner: Agent 1
Branch: agent/1/initial-project-setup
Status: Assigned

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

Follow-up:
- Replace placeholders with the first real product workflow.

---

## Agent 2 In Progress

### Minimum Screen Placeholders

Owner: Agent 2
Branch: agent/2/minimum-screen-placeholders
Status: Assigned

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

Follow-up:
- Implement real account creation, verification, authentication, team, epic, ticket, and board workflows.

---

## Ready For Review

---

## Done
