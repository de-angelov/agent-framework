# AI Development Workflow

These instructions define how AI agents collaborate within this repository.

## Goals

- Prevent duplicate work.
- Maintain a single source of truth for project planning.
- Keep implementation isolated.
- Ensure completed work is verified before acceptance.

---

# Repository Files

The repository uses three coordination documents.

## Repository Boundaries

There are two separate git repositories involved in this workflow:

- The top-level `agent-framework` repository contains the orchestrator, coordination logic, logs, and shared workflow instructions.
- The actual product repository is checked out separately inside `workspaces/repo-tl`, `workspaces/repo-agent-1`, and `workspaces/repo-agent-2`.

When answering whether product work is merged, verified, or complete, inspect the product repository workspace. Do not infer product merge status from the top-level `agent-framework` git history.

## AGENTS.md

Defines:

- Roles
- Workflow
- Git process
- Task lifecycle

## TECH.md

Defines:

- Technology stack
- Architecture
- Coding standards
- Testing
- Verification
- Project-specific engineering decisions

## TASKS.md

The single source of truth for project work.

Contains:

- Backlog
- Agent 1 In Progress
- Agent 2 In Progress
- Done

---

# Roles

Each Codex session has exactly one active role.

Available roles:

- Team Lead
- Agent 1
- Agent 2

The orchestrator selects the role. Do not change roles unless explicitly instructed.

---

# Team Lead

Responsible for:

- Backlog grooming
- Breaking milestones into tasks
- Assigning work
- Assigning new work

The Team Lead does **not** implement feature code unless explicitly instructed.

---

# Agent 1 / Agent 2

Implementation agents.

Responsibilities:

- Implement assigned tasks
- Update implementation progress
- Write tests
- Run verification
- Commit focused changes
- Push assigned branch
- Update task status
- Move completed work into **Done**

Implementation agents must **never**:

- Reprioritize work
- Approve their own work
- Edit another agent's assigned branch

---

# Git Workflow

Each implementation task uses its own branch.

Example:

```text
agent/1/login-page
agent/2/dashboard
```

Implementation agents squash-merge their completed work and move finished tasks into **Done**. The Team Lead only assigns new work.

Completed task branches must be merged into product `main` with a squash merge so `main` receives one final commit per task.

---

# Workspace Isolation

Each role has its own workspace directory.
The top-level repository owns `TASKS.md`; workspace clones are execution-only and do not keep separate boards.

Active workspaces:

- `workspaces/repo-tl` (Team Lead workspace)
- `workspaces/repo-agent-1` (Agent 1 workspace)
- `workspaces/repo-agent-2` (Agent 2 workspace)

Implementation and merges occur only inside implementation workspaces. The Team Lead uses the Team Lead workspace for planning and coordination checks only.

---

# Project Tracking & Syntax

`TASKS.md` is the single source of truth.

> **CRITICAL**
>
> When reading or writing tasks, metadata fields must use strict plaintext line prefixes. Do **not** use Markdown bolding (such as `**Owner:**`) for metadata keys.

Every task must follow this exact format:

```markdown
### Task Title Here

Owner: Agent 1
Branch: agent/1/branch-name
Status: In Progress

[Task Scope, Body, and Progress Notes Go Here]
```

Track only meaningful work. Do **not** track formatting changes, temporary debugging, exploratory edits, or routine commands.

## Task Lifecycle

```text
Backlog → In Progress → Done
```

Backlog work becomes active by moving the task into the appropriate agent lane and setting `Status: In Progress`.
Rejected or blocked work stays in its implementation lane until fixed or reassigned by the Team Lead.

## Board Updates

Implementation agents may update:

- Implementation progress
- Blocked notes
- Task status
- Move their completed task into **Done**

Implementation agents may **not**:

- Reprioritize the backlog
- Assign work

Only implementation agents move their own completed work into **Done** after verification and merge.

Completed tasks must include an explicit timestamp:

```text
Completed: YYYY-MM-DD
```

## Conflict Prevention

Assign non-overlapping work whenever practical.

Separate work by:

- Features
- Routes
- Services
- Directories

Implementation agents modify only files required for their assigned task.

If another active task owns required files:

1. Stop immediately.
2. Document the dependency.
3. Request Team Lead coordination.

## Completion Workflow

Before marking work complete, implementation agents:

1. Run verification defined in `TECH.md`.
2. Commit focused changes.
3. Squash-merge their completed branch into product `main`.
4. Push the completed work.
5. Update the task with verification and merge notes.

After the implementation agent verifies, merges, and pushes the work:

- The implementation agent moves the task to **Done**.
- The implementation agent sets `Status: Done`.
- The implementation agent adds the completion date.
- The Team Lead may assign the next task.

If work cannot be completed or merged:

- Append a `[REJECTED]` section to the task body containing:
  - Failing command
  - Exact output
  - Short explanation of what must be fixed
- Keep the task in its original implementation lane.

The Team Lead does **not** modify implementation code during completion handling.

---

# Orchestrator

The orchestrator coordinates AI execution.

Responsibilities:

- Monitor `TASKS.md`
- Launch implementation agents
- Terminate obsolete sessions when task assignments change
- Use `workspaces/repo-tl` as the Team Lead execution workspace; the top-level repo owns `TASKS.md`

Agents never launch one another.

---

# Instruction Priority

Instructions are applied in the following order:

1. User instructions
2. Repository-specific instructions
3. `AGENTS.md`
4. `TECH.md`
5. General engineering best practices
