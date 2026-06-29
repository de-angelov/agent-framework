# AI Development Workflow

Common rules for all roles in this repository.

## Repository Boundaries

- The top-level `agent-framework` repository owns orchestration, logs, workflow instructions, and coordination files.
- The product repository is checked out separately in `workspaces/repo-tl`, `workspaces/repo-agent-1`, and `workspaces/repo-agent-2`.
- Inspect product work in the product workspace. Do not infer product merge, verification, or completion status from top-level `agent-framework` git history.

## Coordination Files

- `BACKLOG.md`: pending unassigned work.
- `TASKS.md`: live execution lanes only.
- `ARCHIVE.md`: completed work history.
- `AGENTS.md`: common workflow rules.
- `DEV_AGENT.md`: dev agent role rules.
- `TEAM_LEAD_AGENT.md`: team lead agent role rules.
- `TECH.md`: product technical standards and verification.

## Roles

Each Codex session has exactly one active role selected by the orchestrator:

- Team Lead Agent
- Dev Agent 1
- Dev Agent 2

Do not change roles unless explicitly instructed.

## Task Format

Task metadata must use strict plaintext line prefixes. Do not use Markdown bolding for metadata keys.

```markdown
### Task Title Here

Owner: Dev Agent 1
Branch: agent/1/branch-name
Status: In Progress

[Task Scope, Body, and Progress Notes Go Here]
```

Track meaningful project work only. Do not track formatting changes, temporary debugging, exploratory edits, or routine commands.

## Board Lifecycle

```text
BACKLOG.md → TASKS.md active lane → ARCHIVE.md
```

- Backlog work becomes active when the Team Lead Agent moves it from `BACKLOG.md` into a dev-agent lane in `TASKS.md`.
- Rejected or blocked work stays in its dev-agent lane until fixed or reassigned by the Team Lead Agent.
- Completed tasks move to `ARCHIVE.md` with `Status: Done` and `Completed: YYYY-MM-DD`.

## Workspace Isolation

- Team Lead Agent workspace: `workspaces/repo-tl`
- Dev Agent 1 workspace: `workspaces/repo-agent-1`
- Dev Agent 2 workspace: `workspaces/repo-agent-2`

Development and merges occur only inside implementation workspaces. The Team Lead Agent uses the Team Lead Agent workspace for planning and coordination checks only.

## Conflict Prevention

Assign non-overlapping work whenever practical. Separate work by foundations, features, routes, services, and directories.

If another active task owns required files:

1. Stop immediately.
2. Document the dependency.
3. Request Team Lead Agent coordination.

## Instruction Priority

Instructions apply in this order:

1. User instructions
2. Repository-specific instructions
3. `AGENTS.md`
4. Role-specific instructions
5. `TECH.md`
6. General engineering best practices
