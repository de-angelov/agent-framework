# TASKS

Live execution lanes only. Pending backlog lives in BACKLOG.md. Completed work lives in ARCHIVE.md.

## Dev Agent 2 In Progress

### Kanban Board Data Loader

Task ID: BOARD-01
Category: AFK
Owner: Dev Agent 2
Branch: agent/2/kanban-board-data-loader
Status: In Progress
Dependencies: TICKET-03
Blocking tasks: BOARD-02A, BOARD-04A

Objective:
Replace only the board loader's placeholder payload with real authenticated team, epic, and ticket data.

Scope:
- Edit only `app/routes/board/board.tsx` and `app/routes/board/board.test.tsx` unless a tiny route-local loader helper is needed.
- In the loader, require authentication and load teams with `listTeams(db)`.
- Read `teamId` from the URL query string.
- Select the requested team only when it exists in the loaded teams.
- When no valid team is requested, select the first team from `listTeams` ordering.
- When there are no teams, return `selectedTeamId: ""`, `epics: []`, and `tickets: []`.
- For a selected team, load epics with `listEpics(db, { teamId: selectedTeamId })`.
- For a selected team, load tickets with `listTicketsForTeam(db, { teamId: selectedTeamId })`.
- Return `teams`, `selectedTeamId`, `epics`, `tickets`, and `userEmail` from the loader.
- Preserve the existing placeholder board presentation as much as possible; only consume the new loader shape enough for tests to prove real data is loaded.
- Add focused loader tests for unauthenticated redirect, default team selection, explicit valid team selection, invalid team fallback, empty team state, selected-team epics, and selected-team tickets.

Out of Scope:
- Do not change Kanban column/card presentation beyond removing loader placeholder data dependencies.
- Do not add filters, drag-and-drop, virtualization, ticket mutations, or shared shell changes.
- Do not change ticket services.
- Do not implement column grouping, card rendering from real tickets, or filter behavior; those are `BOARD-02A`, `BOARD-02B`, `BOARD-04A`, and `BOARD-04B`.

Acceptance Criteria:
- Board loader tests pass.
- Ticket service tests continue to pass.
- Verification passes.
- No placeholder board data remains in the loader.

Verification:
Run from `workspaces/repo-agent-2`: `npm test -- app/routes/board/board.test.tsx app/services/tickets/tickets.server.test.ts && npm run typecheck`

Coordination:
- Dependencies are resolved on product `main`: TICKET-03 is merged.
- `listTeams`, `listEpics`, and `listTicketsForTeam` already exist; do not add new service functions.
- The current product has four ticket states, not five.
- Avoid ticket details/edit route files owned by Dev Agent 1, and avoid ticket service/schema changes.

Progress:
- Reconciled by Team Lead on 2026-06-30 after finding `workspaces/repo-agent-2` on `agent/2/kanban-board-data-loader`.

---
