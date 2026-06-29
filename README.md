# One-Time Setup

## Project Structure

Create the following repository layout:

```text
my-project/
│
├── .git/
│
├── AGENTS.md          ← Common workflow rules
├── DEV_AGENT.md             ← Dev agent role rules
├── TEAM_LEAD_AGENT.md       ← Team Lead Agent role rules
├── TECH.md            ← Stack, architecture, coding standards
├── BACKLOG.md         ← Pending work
├── TASKS.md           ← Active dev-agent lanes
├── ARCHIVE.md         ← Completed work history
│
├── setup.go           ← One-time bootstrap
│
├── orchestrator.go
│
└── workspaces/
    ├── repo-tl/       ← Team Lead Agent clone
    ├── repo-agent-1/  ← Dev Agent 1 clone
    └── repo-agent-2/  ← Dev Agent 2 clone
```

The top-level repo owns `BACKLOG.md`, `TASKS.md`, `ARCHIVE.md`, `AGENTS.md`, `DEV_AGENT.md`, `TEAM_LEAD_AGENT.md`, and `TECH.md`.
Each `repo-*` directory is a separate execution clone used by the orchestrator.

---

# Step 1 — Create the Project

The user creates:

- `AGENTS.md`
- `DEV_AGENT.md`
- `TEAM_LEAD_AGENT.md`
- `TECH.md`

Then starts with empty coordination files.

`BACKLOG.md`:

```md
## Backlog
```

`TASKS.md`:

```md
## Dev Agent 1 In Progress

## Dev Agent 2 In Progress
```

`ARCHIVE.md`:

```md
## Done
```

---

# Step 2 — Start the Orchestrator

Run the bootstrap first:

```bash
cd setup
go run setup.go
```

Then start the orchestrator from the top-level repo:

```bash
go run orchestrator.go
```

Output:

```text
Setting up AI development workflow...
Root: /path/to/my-project
✓ dir: workspaces
...

Setup complete.

Next steps:
  go run orchestrator.go
  cd workspaces/repo-tl && go run orchestrator.go
  cd workspaces/repo-agent-1 && go run orchestrator.go
  cd workspaces/repo-agent-2 && go run orchestrator.go

Orchestrator started...

repo root: /path/to/my-project
```

Nothing happens because no work has been assigned.

---

# Step 3 — User Gives a Goal

Instead of editing `TASKS.md` manually, the user should ask the Team Lead Agent:

> Build a todo application.

The Team Lead Agent updates the top-level `BACKLOG.md`:

```md
## Backlog

- Database
- API
- UI
```

Or immediately assigns work.

---

# Step 4 — Team Lead Agent Assigns Work

### Dev Agent 1

```md
## Dev Agent 1 In Progress

### Create Todo Database

Owner: Dev Agent 1
Branch: agent/1/todo-db
Status: In Progress
```

### Dev Agent 2

```md
## Dev Agent 2 In Progress

### Build Home Page

Owner: Dev Agent 2
Branch: agent/2/home-page
Status: In Progress
```

---

# Step 5 — Orchestrator Notices

On the next poll it sees:

- Dev Agent 1
- `Status = In Progress`

It starts Codex in:

```text
workspaces/repo-agent-1/
```

With the role:

```text
Dev Agent 1
```

Likewise for Dev Agent 2.

---

# Step 6 — Dev Agents Work

Each agent:

1. Reads `AGENTS.md`
2. Reads `DEV_AGENT.md`
3. Reads `TECH.md`
4. Reads board context from `BACKLOG.md` and `TASKS.md`
5. Checks out its branch
6. Writes code
7. Runs tests
8. Commits
9. Pushes
10. Squash-merges its completed branch into product `main`
11. Moves its completed task to `ARCHIVE.md`

When finished, the task moves from:

```md
## Dev Agent 1 In Progress
```

to `ARCHIVE.md`:

```md
## Done
```

and its status becomes:

```text
Status: Done
```

---

# Step 7 — Orchestrator Notices Again

When a lane is free and backlog work exists, the orchestrator launches the Team Lead Agent in:

```text
workspaces/repo-tl/
```

Role:

```text
Team Lead Agent
```

---

# Step 8 — Team Lead Agent Assigns

The Team Lead Agent grooms and assigns backlog work. Dev agents own verification, pushing, squash-merging, and archiving their completed tasks.

Example:

```md
## Done

### Create Todo Database

Completed: 2026-06-28
```

---

# Step 9 — Team Lead Agent Assigns More Work

Example:

```md
## Dev Agent 1 In Progress

### Implement Authentication

Owner: Dev Agent 1
Branch: agent/1/auth
Status: In Progress
```

The orchestrator notices the new assignment and starts Dev Agent 1 again.

---

# Continuous Loop

```text
User
  │
  ▼
Team Lead Agent
  │
  ▼
BACKLOG.md / TASKS.md
  │
  ▼
Orchestrator
  │
  ▼
Agent
  │
  ▼
Git Push + Squash Merge
  │
  ▼
ARCHIVE.md
  │
  ▼
Assign Next Task
  │
  └───────────────► Repeat
```
