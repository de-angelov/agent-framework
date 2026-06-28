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

## AGENTS.md

Defines:

- Roles
- Workflow
- Git process
- Review process
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
- Ready For Review
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
- Architecture review
- Implementation review
- Verification
- Merging approved work
- Moving tasks into **Done**
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

Implementation agents must **never**:

- Reprioritize work
- Merge branches
- Approve their own work
- Move tasks into **Done**
- Edit another agent's assigned branch

---

# Git Workflow

Each implementation task uses its own branch.

Example:

```text
agent/1/login-page
agent/2/dashboard
```

Only the Team Lead merges into `main`. Implementation agents never work directly on `main`.

---

# Workspace Isolation

Each role has its own workspace directory.
The top-level repository owns `TASKS.md`; workspace clones are execution-only and do not keep separate boards.

Active workspaces:

- `workspaces/repo-tl` (Team Lead workspace)
- `workspaces/repo-agent-1` (Agent 1 workspace)
- `workspaces/repo-agent-2` (Agent 2 workspace)

Implementation occurs only inside implementation workspaces. The Team Lead performs reviews and merges inside the Team Lead workspace.

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
Status: Assigned

[Task Scope, Body, and Progress Notes Go Here]
```

Track only meaningful work. Do **not** track formatting changes, temporary debugging, exploratory edits, or routine commands.

## Task Lifecycle

```text
Backlog → Assigned → In Progress → Ready For Review → Done
```

Rejected work returns to the assigned implementation lane.

## Board Updates

Implementation agents may update:

- Implementation progress
- Blocked notes
- Task status

Implementation agents may **not**:

- Reprioritize the backlog
- Assign work
- Move tasks into **Done**

Only the Team Lead moves completed work into **Done**.

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

## Review Workflow

When work reaches **Ready For Review**, the Team Lead:

1. Reviews the implementation.
2. Runs verification defined in `TECH.md`.
3. Accepts or rejects the work.

If accepted:

- Merge the branch into `main`.
- Update `TASKS.md` (move the task to **Done** and add the completion date).
- Assign the next task.

If rejected:

- Return the task to its original implementation lane.
- Append a `[REJECTED]` section to the task body containing:
  - Failing command
  - Exact output
  - Short explanation of what must be fixed

The Team Lead does **not** modify implementation code during a review.

---

# Orchestrator

The orchestrator coordinates AI execution.

Responsibilities:

- Monitor `TASKS.md`
- Launch implementation agents
- Launch Team Lead reviews
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
