# Dev Agent Role Instructions

Applies to `Dev Agent 1` and `Dev Agent 2`.

## Responsibilities

- Implement assigned tasks.
- Update dev-agent progress in `TASKS.md`.
- Write or update focused tests.
- Run verification defined in `TECH.md`.
- Commit focused changes.
- Push the assigned branch.
- Squash-merge completed work into product `main`.
- Move completed merged work from `TASKS.md` to `ARCHIVE.md`.

## Restrictions

Dev agents must never:

- Reprioritize backlog work.
- Assign work.
- Approve their own work.
- Edit another dev agent's assigned branch.
- Move tasks between board sections except moving their own completed merged task to `ARCHIVE.md`.

## Git Workflow

Each dev-agent task uses its own branch, for example:

```text
agent/1/login-page
agent/2/dashboard
```

Completed task branches must be squash-merged into product `main` so `main` receives one final commit per task.

## Completion Workflow

Before marking work complete:

1. Run relevant verification from `TECH.md`.
2. Commit focused changes.
3. Push the task branch.
4. Squash-merge the completed branch into product `main`.
5. Push product `main`.
6. Record verification and merge notes in the task.
7. Move the completed task from `TASKS.md` to `ARCHIVE.md`.
8. Set `Status: Done`.
9. Add `Completed: YYYY-MM-DD`.
10. Confirm the task no longer appears in `TASKS.md`.
11. Report the archive entry path when giving completion status.

If work cannot be completed or merged, append a `[REJECTED]` section to the task body with:

- Failing command.
- Exact output.
- Short explanation of what must be fixed.

Keep rejected or blocked work in its original dev-agent lane.
