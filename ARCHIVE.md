# ARCHIVE

Completed work history. Normal orchestrator prompts do not load this file.

## Done

### API and Persistence Foundation

Owner: Dev Agent 1
Branch: agent/1/api-persistence-foundation
Status: Done
Completed: 2026-06-28

Outcome:
Established the database, migration, identifier, timestamp, and no-seed-data conventions shared by all workflows.

Scope:
- Added shared UUID and ISO-8601 UTC timestamp helpers for server-side persistence code.
- Updated team and epic services to use the shared identifier and timestamp helpers.
- Enabled SQLite foreign key enforcement for application and migration database connections.
- Added automated coverage for timestamp serialization and fresh database initialization through Drizzle migrations.
- Verified fresh migrations create schema and migration metadata without preloading application users, teams, epics, tickets, or comments.

Coordination:
- Rebased onto merged authentication work before merge.
- Kept changes focused on persistence infrastructure, schema/migration initialization, timestamp helpers, identifier conventions, and focused tests.

Verification:
- `npm test` passed.
- `npm run typecheck` passed.
- `npm run build` passed.
- `npm run db:migrate` passed.

Merge:
- Branch `agent/1/api-persistence-foundation` pushed.
- Squash-merged into product `main` as commit `c5cf894` (`Task 11: add API persistence foundation`) and pushed.

Follow-up:
- Add seed data or demo fixtures only behind an explicit non-default development command if later required.

---


### User Accounts and Authentication

Owner: Dev Agent 2
Branch: agent/2/user-accounts-authentication
Status: Done
Completed: 2026-06-28

Outcome:
Implemented local user accounts, email verification, login, logout, and authentication protection for business application surfaces.

Scope:
- Added signup, login, logout, email verification, and verification-email resend routes.
- Trimmed and lowercased email addresses for case-insensitive uniqueness.
- Enforced minimum 8-character passwords and Argon2id password hashing.
- Added SMTP-backed verification email delivery with `relay1.dataart.com` as the default host.
- Added single-use, 24-hour verification tokens and invalidated earlier unused tokens on resend.
- Prevented unverified accounts from logging in or accessing business screens.
- Guarded business loaders/actions with cookie-backed sessions while leaving public auth routes available.
- Added automated coverage for normalization, password rules, hashing, token lifecycle, resend invalidation, login sessions, logout, and route guards.

Coordination:
- Rebased onto merged teams and epic data-service work.
- Replaced the temporary teams auth helper with the real session-backed auth boundary.
- Added auth persistence as `drizzle/0003_user-auth.sql` after the merged teams and epics migrations.

Verification:
- `npm test` passed.
- `npm run typecheck` passed.
- `npm run build` passed.
- `npm run db:migrate` passed with `DATABASE_URL` pointed at a fresh temporary SQLite database.

Merge:
- Branch `agent/2/user-accounts-authentication` pushed.
- Squash-merged into product `main` as commit `b602bbd` (`Task 10: implement user accounts and authentication`) and pushed.

Follow-up:
- Add SSO or external identity providers only if later required.

---


### Epic Data Model and Services

Owner: Dev Agent 1
Branch: agent/1/epics
Status: Done
Completed: 2026-06-28

Outcome:
Implemented backend epic persistence and business rules for team-scoped epics.

Scope:
- Made each epic belong to exactly one team.
- Stored each epic with an identifier, team reference, title, optional description, created timestamp, and modified timestamp.
- Required epic titles to be non-empty after trimming.
- Prevented changing an epic's team after creation.
- Kept moving epics between teams out of scope.
- Prevented deleting an epic while tickets reference it.
- Added automated tests for epic creation, listing, editing, deletion, title validation, immutable team assignment, and blocked deletion.

Coordination:
- Kept backend business rules in `app/services/epics.server.ts`.
- Added nullable `tickets.epic_id` only as the minimal persistence reference required for blocked epic deletion; ticket-side same-team selection remains out of scope for ticket tasks.
- Added `drizzle/0002_material_famine.sql` and migration metadata for epic fields and ticket epic references.

Verification:
- `npm run db:migrate` passed.
- `npm test` passed.
- `npm run typecheck` passed.
- `npm run build` passed.

Merge:
- Branch `agent/1/epics` pushed.
- Squash-merged into product `main` as commit `ad4f2e5` and pushed.

Follow-up:
- Add moving epics between teams only if later required.

---


### Teams

Owner: Dev Agent 1
Branch: agent/1/teams-management
Status: Done
Completed: 2026-06-28

Outcome:
Implemented team management and enforced team rules for grouping tickets.

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
- Built the route auth boundary as a narrow helper because `agent/2/user-accounts-authentication` had not landed in product `main`; future auth work can replace that helper without changing team business rules.
- Kept backend business rules in `app/services/teams.server.ts` and kept the route focused on authentication, form parsing, service calls, and rendering.
- Kept UI work scoped to `app/routes/teams.tsx` with existing simple layout styles.
- Added only minimal `epics` and `tickets` references required for team deletion guards.

Verification:
- Passed `npm test -- teams`.
- Passed `npm run typecheck`.
- Passed `npm run db:migrate`.
- Passed `npm test`.
- Passed `npm run build`.

Merge:
- Pushed branch `agent/1/teams-management`.
- Squash-merged into product `main` as commit `d2de4ec` (`Task 8: add teams management`).
- Pushed product `main`.

Follow-up:
- Add team ownership and membership only if later required.

---


---

### Docker Compose Runtime

Owner: Dev Agent 1
Branch: agent/1/docker-compose-runtime
Status: Done
Completed: 2026-06-28

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
- Coordinate before changing package scripts, migration commands, or database defaults that overlap with Dev Agent 2's API and persistence expectations work.

Progress:
- Added Docker runtime files for a production build, local SQLite volume storage, and startup migrations.
- Added README instructions for `docker compose up --build`, exposed `http://localhost:3000`, stopping, and resetting the Docker volume.
- Resolved README merge overlap with the merged branch-cleanup documentation by preserving both sections.

Verification:
- Passed `npm run typecheck`.
- Passed `npm test`.
- Passed `npm run build`.
- Passed `DATABASE_URL=local.db node scripts/migrate.mjs`.
- `HOST=127.0.0.1 PORT=3000 DATABASE_URL=local.db timeout 10s npm run start` reached `EADDRINUSE` because port `3000` was already occupied in this environment.
- Passed local production startup smoke test with migrated SQLite database on an alternate port: `HOST=127.0.0.1 PORT=3001 DATABASE_URL=local.db timeout 10s npm run start` printed `[react-router-serve] http://localhost:3001 (http://127.0.0.1:3001)`.
- Could not run `docker compose up --build` in this environment because the `docker` CLI is not installed: `timeout: failed to run command ‘docker’: No such file or directory`.

Merge:
- Commit pushed to `origin/agent/1/docker-compose-runtime`: `8259f7c Task 6: add Docker Compose runtime`.
- Squash-merged into product `main` and pushed commit `c5aafe7 Task 7: merge Docker Compose runtime`.

Follow-up:
- Add production deployment manifests only if a deployment target is later defined.

---


---

### GitHub Action Merged Branch Cleanup

Owner: Dev Agent 2
Branch: agent/2/github-action-merged-branch-cleanup
Status: Done
Completed: 2026-06-28

Outcome:
Add a scheduled GitHub Actions workflow that removes remote branches already merged into `main`.

Scope:
- Add a GitHub Actions workflow under `.github/workflows`.
- Run the workflow daily at 00:00 UTC using a cron schedule.
- Fetch remote branches and identify branches already merged into `main`.
- Delete merged remote branches while preserving `main` and any protected or configured keep branches.
- Avoid deleting active dev-agent branches that are not merged.
- Allow manual workflow dispatch for explicit cleanup runs.
- Document the branch cleanup behavior and safety exclusions where practical.

Coordination:
- Keep the workflow scoped to repository branch maintenance.
- Do not change product build, test, deploy, or application runtime behavior.
- Use the built-in GitHub token and least-permission workflow settings where practical.

Progress:
- Added `.github/workflows/cleanup-merged-branches.yml` with a daily 00:00 UTC schedule and manual dispatch.
- Fetches remote branches, checks whether each remote tip is an ancestor of `origin/main`, and deletes only merged branches.
- Skips `main`, default configured keep branches, manual keep patterns, protected branches reported by GitHub, and branches not merged into `main`.
- Added manual `dry_run` support and README documentation for cleanup behavior and safety exclusions.
- Verified with `actionlint`, `npm test`, `npm run typecheck`, and `npm run build`.
- Pushed branch `agent/2/github-action-merged-branch-cleanup` at commit `49b83ec`.
- Merged into product `main` and pushed commit `d24501d`.

Follow-up:
- Add notifications only if cleanup failures need visibility later.

---


---

### Authenticated Header Navigation

Owner: Dev Agent 2
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

Owner: Dev Agent 1
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

Owner: Dev Agent 1
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

Owner: Dev Agent 2
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

Owner: Dev Agent 1
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
- Keep replacements narrow to avoid colliding with Dev Agent 2 placeholder route work.

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

Owner: Dev Agent 1
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
