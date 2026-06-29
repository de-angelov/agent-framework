# TASKS

Live execution lanes only. Pending backlog lives in BACKLOG.md. Completed work lives in ARCHIVE.md.

## Dev Agent 1 In Progress

### Wireframe Team Management UI Alignment

Owner: Dev Agent 1
Branch: agent/1/wireframe-team-management-ui-alignment
Status: In Progress

Outcome:
Align the team management screen with the team wireframe hierarchy.

Scope:
- Shape team management around a table-style list with name, ticket count, epic count, modified timestamp, row edit/delete actions, disabled delete affordance when referenced, explanatory blocked-delete copy, and create/edit form placement consistent with the wireframe.
- Preserve existing team create, rename, delete, unique-name validation, and blocked-delete behavior.
- Preserve required loading, empty, success, and error states while applying the wireframe hierarchy.
- Add or update focused tests for team counts, modified timestamp display, row edit/delete actions, disabled delete affordance, blocked-delete explanation, create form, and edit form.

Coordination:
- `Teams` and `Component Folder and CSS Module Cleanup` are already recorded as Done, so this can proceed without overlapping the completed shared-component cleanup lane.
- Reuse the shared authenticated shell if `Wireframe Shared App Shell` has landed first.
- Keep this task scoped to team route presentation and focused tests.
- Avoid changing backend team business rules, database schema, ticket behavior, epic behavior, or route permissions.

Follow-up:
- If ticket or epic counts are unavailable from existing loaders/services, add a separate backend count-query task before expanding this UI branch.

Progress:
- Assigned by Team Lead on 2026-06-29.

---

## Dev Agent 2 In Progress
