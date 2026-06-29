# TASKS

Live execution lanes only. Pending backlog lives in BACKLOG.md. Completed work lives in ARCHIVE.md.

## Dev Agent 1 In Progress

### Ticket Schema Expansion

Task ID: TICKET-01
Category: AFK
Owner: Dev Agent 1
Branch: agent/1/ticket-schema-expansion
Status: In Progress
Dependencies: API and Persistence Foundation; User Accounts and Authentication; Teams and epics behavior merged
Blocking tasks: TICKET-02, TICKET-03, TICKET-04, COMMENT-01, BOARD-01

Objective:
Expand ticket persistence so tickets can store the required workflow fields.

Scope:
- Add a new migration that extends the existing `tickets` table for title, body, type, state, created-by, created-at, and modified-at.
- Keep team and optional epic references intact.
- Add database-level foreign keys where practical for team, epic, and created-by references.
- Add app-owned ticket type and ticket state constants/types in a server-safe module.
- Add focused schema or migration tests for a fresh database containing the expanded ticket fields.

Out of Scope:
- Do not implement ticket create, update, delete, or UI behavior.
- Do not add comments or activity history.
- Do not modify already-applied migrations.

Acceptance Criteria:
- Fresh database migration coverage passes.
- Existing team and epic blocked-delete behavior still works with the expanded schema.
- `npm run db:migrate` succeeds.
- No placeholder TODOs remain.

Verification:
Run from `workspaces/repo-agent-1`: `npm test -- app/db/fresh-database.test.ts app/services/teams.server.test.ts app/services/epics.server.test.ts && npm run db:migrate && npm run typecheck`

Coordination:
- Dependencies are resolved in `ARCHIVE.md`: API and Persistence Foundation, User Accounts and Authentication, Epic Data Model and Services, Team Management, and Epic Management UI.
- Keep edits scoped to database schema/migrations, server-safe ticket constants/types, and focused schema/service tests.
- Avoid route, UI, ticket service behavior, comments, board behavior, and starter cleanup changes.
- Coordinate before touching epic route presentation files because Dev Agent 2 owns `EPIC-UI-01`.

Progress:
- Assigned by Team Lead on 2026-06-29 for parallel execution with `EPIC-UI-01`.
