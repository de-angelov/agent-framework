# Technical Standards

This document defines the engineering standards for this repository.

---

# Technology Stack

## Runtime

- Node.js (current LTS)
- TypeScript (strict mode)

---

## Frontend

- React
- React Router

Keep route modules focused on request boundaries, data loading, and rendering.

Shared UI components must live in their own subfolders under the component area.
Each component folder should contain the component implementation, its focused
test file, and its CSS module when styling is needed.

Example:

```text
app/components/button/button.tsx
app/components/button/button.test.tsx
app/components/button/button.module.css
```

Prefer CSS modules for component styling. Avoid adding new component-specific
styles to global stylesheets unless the style is truly global application
infrastructure.

---

## Backend

- React Router server actions/loaders
- Service layer for business logic

Business logic belongs in reusable server-side service modules rather than route handlers.

---

## Database

- SQLite (local MVP)
- Drizzle ORM

Rules:

- Never modify an applied migration.
- Always create a new migration for schema changes.
- Run database migrations after schema changes:

```bash
npm run db:migrate
```

---

# Architecture

Business rules belong in reusable services or pure functions.

Separate business logic from:

- database access
- networking
- framework code
- session handling
- UI rendering

Keep route modules responsible only for:

- request validation
- authentication/session boundaries
- loading data
- invoking services
- rendering responses

Prefer dependency injection over direct infrastructure imports where practical.

Remove unused starter or template code.

---

# Coding Standards

Prefer:

- strict TypeScript
- small composable functions
- explicit interfaces and types
- reusable service modules
- pure business logic
- function composition
- deterministic transformations

Avoid:

- `any`
- dead code
- speculative abstractions
- duplicated business rules

Validate inputs at application boundaries.

Use named functions when logic is:

- reused
- domain-specific
- directly tested

Short inline lambdas are preferred for simple one-off transformations.

---

# Error Handling

Use `neverthrow` for expected failures.

Expected domain failures should be returned as `Result` values.

Unexpected failures may throw.

Do not use exceptions for normal control flow.

Services should communicate recoverable failures explicitly.

---

# Branching

Use `ts-pattern` when branching over:

- domain states
- action intents
- enums
- status values
- error mappings
- complex conditional logic

Prefer exhaustive `match(...)` expressions.

Use `.exhaustive()` when matching a typed union or enum owned by the
application, such as domain states, service errors, and known status values.
Use `.otherwise(...)` only when the input can contain unknown external values,
such as raw form intents, query parameters, or untrusted request payloads.

Do not replace simple guard clauses or nullish defaults with pattern matching when it reduces readability.

---

# Utilities

Prefer Remeda whenever an equivalent helper exists.

Especially prefer:

- `pipe`
- collection transformations
- immutable utility helpers

Avoid writing custom helpers that duplicate existing Remeda functionality.

---

# Testing

Every feature and business rule requires automated tests.

Testing stack:

- Vitest for unit tests
- Playwright for end-to-end and workflow testing
- React Router `createRoutesStub` only for isolated component tests requiring router context

Place tests close to the code they validate whenever practical.

---

# Verification

Dev agents must run relevant verification before hand-off.

Typical commands:

```bash
npm test
npm run typecheck
npm run build
```

Run migrations whenever schema changes:

```bash
npm run db:migrate
```

Before hand-off remove:

- console logging
- temporary debugging
- commented-out code
- temporary test files

Team Lead Agent verification:

- full typecheck
- full test suite
- production build

---

# Styling

Until core workflows are complete:

Use only simple layouts.

Required default styling:

- Flexbox layouts
- `padding: 10px`
- `border: 1px solid grey`

Avoid unnecessary visual polish until core functionality is complete.
Do not add custom fonts, shadows, gradients, rounded corners, or decorative spacing unless explicitly required.
Keep component styles colocated in `*.module.css` files inside the component's
own folder. Route-level styles may use route-specific CSS modules when the
styles are not shared. Keep `app/styles.css` limited to global resets, document
defaults, and intentionally shared application-level primitives.

---

# Git

Create focused commits.

Do not:

- mix unrelated work
- rewrite history without instruction
- revert user changes without instruction

Each completed task should be committed separately.

When merging a completed task branch into `main`, use a squash merge. `main` should receive one final commit for the task, even if the implementation branch contains multiple working commits.

Commit messages must include:

- the ticket/task number
- a concise description of the completed work

If no task number exists, use the next sequential task number from the previous commit.

---

# Project Workflow

Use the coordination board files as the project source of truth:

- `BACKLOG.md` for pending work
- `TASKS.md` for active dev-agent lanes
- `ARCHIVE.md` for completed work

When starting a new implementation session:

- determine whether the session is Dev Agent 1 or Dev Agent 2 before modifying the board
- move active work into the appropriate In Progress column
- move completed work into `ARCHIVE.md` with the completion date
- record meaningful follow-up work in `BACKLOG.md`

Track project features and milestones, not individual commands or minor edits.

---

# Project Notes

Current project conventions:

- React Router application
- SQLite for local development
- Drizzle ORM for persistence
- `neverthrow` for expected errors
- `ts-pattern` for structured branching
- Remeda for functional utilities
- Vitest for unit tests
- Playwright for end-to-end tests
- Keep business rules in reusable services and pure functions
- Keep routes thin and focused on request boundaries
- MVP-first development with intentionally simple UI

Update this section as architecture, deployment, conventions, or preferred libraries evolve.
