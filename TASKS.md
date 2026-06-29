# TASKS

Live execution lanes only. Pending backlog lives in BACKLOG.md. Completed work lives in ARCHIVE.md.

## Dev Agent 1 In Progress

## Dev Agent 2 In Progress

### Epic Management UI

Owner: Dev Agent 2
Branch: agent/2/epic-management-ui
Status: In Progress

Outcome:
Provide a separate epic management screen for creating, listing, editing, and deleting epics.

Scope:
- Select the team when an epic is created.
- List epics with their team, title, optional description, created timestamp, and modified timestamp where practical.
- Allow authenticated users to create, edit, and delete epics.
- Show a clear UI validation message when epic deletion is blocked.
- Keep moving epics between teams out of scope.
- Add focused route/component coverage for listing, creation, editing, deletion, title validation, and blocked deletion messaging.

Coordination:
- Build after Dev Agent 1 completes `Epic Data Model and Services`.
- `Epic Data Model and Services`, `Teams`, and `User Accounts and Authentication` are already recorded as Done.
- Keep UI work scoped to the epic management route and minimal styling only.
- Use only simple flexbox layouts, `padding: 10px`, and `border: 1px solid grey`.
- Avoid shared dialog, header, button internals, broad styling, ticket CRUD, and board behavior.
- Do not modify Dev Agent 1's active `API and Persistence Foundation` branch or persistence-foundation internals.

Follow-up:
- Add bulk epic management only if later required.

---


---

