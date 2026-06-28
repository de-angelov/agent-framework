# One-Time Setup

## Project Structure

Create the following repository layout:

```text
my-project/
│
├── .git/
│
├── AGENTS.md          ← Workflow (roles, TL, agents, review)
├── TECH.md            ← Stack, architecture, coding standards
├── TASKS.md           ← Live project board
│
├── setup.go           ← One-time bootstrap
│
├── orchestrator.go
│
└── workspaces/
    ├── repo-tl/       ← Team Lead clone
    ├── repo-agent-1/  ← Agent 1 clone
    └── repo-agent-2/  ← Agent 2 clone
```

The top-level repo owns `TASKS.md`, `AGENTS.md`, and `TECH.md`.
Each `repo-*` directory is a separate execution clone used by the orchestrator.

---

# Step 1 — Create the Project

The user creates:

- `AGENTS.md`
- `TECH.md`

Then starts with an empty `TASKS.md`.

Example:

```md
## Backlog

## Agent 1 In Progress

## Agent 2 In Progress

## Ready For Review

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

Instead of editing `TASKS.md` manually, the user should ask the Team Lead:

> Build a todo application.

The Team Lead updates the top-level `TASKS.md`:

```md
## Backlog

- Database
- API
- UI
```

Or immediately assigns work.

---

# Step 4 — Team Lead Assigns Work

### Agent 1

```md
## Agent 1 In Progress

### Create Todo Database

Owner: Agent 1
Branch: agent/1/todo-db
Status: Assigned
```

### Agent 2

```md
## Agent 2 In Progress

### Build Home Page

Owner: Agent 2
Branch: agent/2/home-page
Status: Assigned
```

---

# Step 5 — Orchestrator Notices

On the next poll it sees:

- Agent 1
- `Status = Assigned`

It starts Codex in:

```text
workspaces/repo-agent-1/
```

With the role:

```text
Agent 1
```

Likewise for Agent 2.

---

# Step 6 — Agents Work

Each agent:

1. Reads `AGENTS.md`
2. Reads `TECH.md`
3. Reads `TASKS.md`
4. Checks out its branch
5. Writes code
6. Runs tests
7. Commits
8. Pushes

When finished, the task moves from:

```md
## Agent 1 In Progress
```

to:

```md
## Ready For Review
```

and its status becomes:

```text
Status: Ready For Review
```

---

# Step 7 — Orchestrator Notices Again

The orchestrator launches the Team Lead in:

```text
workspaces/repo-tl/
```

Role:

```text
Team Lead
```

---

# Step 8 — Team Lead Reviews

The Team Lead:

1. Fetches the branch
2. Reviews the diff
3. Runs tests
4. Merges
5. Updates `TASKS.md`

Example:

```md
## Done

### Create Todo Database

Completed: 2026-06-28
```

---

# Step 9 — Team Lead Assigns More Work

Example:

```md
## Agent 1 In Progress

### Implement Authentication

Owner: Agent 1
Branch: agent/1/auth
Status: Assigned
```

The orchestrator notices the new assignment and starts Agent 1 again.

---

# Continuous Loop

```text
User
  │
  ▼
Team Lead
  │
  ▼
TASKS.md
  │
  ▼
Orchestrator
  │
  ▼
Agent
  │
  ▼
Git Push
  │
  ▼
Ready For Review
  │
  ▼
Orchestrator
  │
  ▼
Team Lead
  │
  ▼
Merge
  │
  ▼
Assign Next Task
  │
  └───────────────► Repeat
```
