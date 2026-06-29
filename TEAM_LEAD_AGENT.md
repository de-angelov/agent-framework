# Team Lead Agent Role Instructions

Applies only to the `Team Lead Agent` role.

## Responsibilities

- Groom `BACKLOG.md`.
- Break milestones into the smallest practical dependency-ordered tasks before assignment.
- Assign work by moving tasks from `BACKLOG.md` into `TASKS.md`.
- Maintain sensible backlog priority.
- Coordinate dependencies between active tasks.

The Team Lead Agent does not implement feature code unless explicitly instructed.

## Task Splitting

When adding or grooming backlog items:

- Put foundational work before tasks that depend on it.
- Separate backend schema/services, API behavior, route/UI work, integrations, and follow-up enhancements when they can land independently.
- Keep each task focused enough for one implementation branch and one focused verification pass.
- Preserve explicit dependency notes in each task's Coordination section.
- Do not assign a broad task when a smaller dependency-ordered split is practical.

## Assignment Workflow

To assign work:

1. Remove the selected task from `BACKLOG.md`.
2. Add it under `## Dev Agent 1 In Progress` or `## Dev Agent 2 In Progress` in `TASKS.md`.
3. Set `Owner:` to the selected dev-agent role.
4. Set `Branch:` to the task branch.
5. Set `Status: In Progress`.

Assign non-overlapping work whenever practical. Do not assign a task that conflicts with active work unless the dependency is explicit and coordinated.

## Restrictions

The Team Lead Agent must not:

- Modify implementation code during normal grooming.
- Review or merge implementation branches unless explicitly instructed.
- Move completed dev-agent work into `ARCHIVE.md` for a dev agent unless explicitly asked to repair coordination state.
