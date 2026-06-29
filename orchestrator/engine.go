package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

func reconcile(tasks []Task) {
	desired := map[string]Task{}
	invalidRoles := map[string]bool{}
	var backlogTask *Task

	addDesired := func(role string, task Task) {
		if invalidRoles[role] {
			return
		}
		if existing, exists := desired[role]; exists {
			logEvent("TASKS board error: multiple In Progress tasks for %s (%q and %q); refusing to start role", role, existing.Title, task.Title)
			delete(desired, role)
			invalidRoles[role] = true
			return
		}
		desired[role] = task
	}

	for _, task := range tasks {
		switch {
		case task.Section == "Agent 1 In Progress" && task.Owner == "Agent 1" && task.Status == "In Progress":
			addDesired("Agent 1", task)
		case task.Section == "Agent 2 In Progress" && task.Owner == "Agent 2" && task.Status == "In Progress":
			addDesired("Agent 2", task)
		case task.Section == "Backlog" && backlogTask == nil:
			if task.Status == "Backlog" || task.Status == "" {
				copy := task
				backlogTask = &copy
			}
		}
	}

	if backlogTask != nil && len(invalidRoles) == 0 {
		_, agent1Busy := desired["Agent 1"]
		_, agent2Busy := desired["Agent 2"]
		if agent1Busy && agent2Busy {
			backlogTask = nil
		}
	}

	if backlogTask != nil && len(invalidRoles) > 0 {
		backlogTask = nil
	}

	if backlogTask != nil {
		desired["Team Lead"] = *backlogTask
	}

	mu.Lock()
	for role := range finished {
		if _, stillDesired := desired[role]; !stillDesired {
			delete(finished, role)
		}
	}

	for role, session := range running {
		task, stillDesired := desired[role]
		taskKey := task.Key

		if !stillDesired || taskKey != session.TaskKey {
			logEvent("stopping %s because TASKS.md changed", role)
			session.Cancel()
			delete(running, role)
		}
	}
	mu.Unlock()

	for role, task := range desired {
		taskKey := task.Key

		if !workspaceExists(role) {
			logEvent("skipping %s because workspace is missing", role)
			continue
		}

		mu.Lock()
		_, exists := running[role]
		finishedSession, alreadyFinished := finished[role]
		if alreadyFinished && finishedSession.TaskKey != taskKey {
			delete(finished, role)
			alreadyFinished = false
		}
		mu.Unlock()

		if exists || alreadyFinished {
			continue
		}

		startSession(role, task, tasks)
	}
}

func startSession(role string, task Task, tasks []Task) {
	ctx, cancel := context.WithCancel(context.Background())
	taskKey := task.Key

	mu.Lock()
	running[role] = RunningSession{
		Role:    role,
		Task:    task.Title,
		TaskKey: taskKey,
		Branch:  task.Branch,
		Cancel:  cancel,
	}
	mu.Unlock()

	logEvent("starting %s on %s", role, task.Title)

	go func() {
		outcome := sessionUnknown
		defer func() {
			mu.Lock()
			delete(running, role)
			if outcome != sessionCancelled {
				finished[role] = FinishedSession{
					Role:       role,
					Task:       task.Title,
					TaskKey:    taskKey,
					Branch:     task.Branch,
					Outcome:    outcome,
					FinishedAt: time.Now(),
				}
			}
			mu.Unlock()
		}()

		switch role {
		case "Agent 1":
			outcome = runAgent(ctx, "Agent 1", agent1Path, task, tasks)
		case "Agent 2":
			outcome = runAgent(ctx, "Agent 2", agent2Path, task, tasks)
		case "Team Lead":
			outcome = runTeamLead(ctx, task, tasks)
		}
	}()
}

func runAgent(ctx context.Context, role string, workspace string, task Task, tasks []Task) SessionOutcome {
	logEvent("starting %s on %s in %s", role, task.Title, workspace)

	if task.Branch != "" {
		if err := prepareBranch(workspace, task.Branch); err != nil {
			logEvent("failed to prepare branch %s in %s: %v", task.Branch, workspace, err)
			return sessionFailed
		}
	}

	prompt := buildPrompt(role, task, tasks, `
Role: Implementation Agent
Rules:
- Follow AGENTS.md and TECH.md.
- UI-Limits: flexbox layout only; padding:10px; border:1px solid grey. No custom fonts/shadows/gradients/rounded-corners unless explicitly requested.
- Scope: Work ONLY on assigned task. Focused changes. No backlog modifications. No self-approval.
- Workflow: Write/update tests. Run verification. Commit, push branch, squash-merge into product main when done.
- Updates: Document progress/verification/merges in TASKS.md. On complete, move task to Done (Status: Done, Completed: YYYY-MM-DD). Do not alter other tasks.
`)

	return runCodex(ctx, workspace, prompt)
}

func runTeamLead(ctx context.Context, task Task, tasks []Task) SessionOutcome {
	runGit(teamLeadPath, "fetch", "--all", "--prune")

	currentBranch := currentBranchName(teamLeadPath)
	if currentBranch != "" {
		logEvent("team lead branch: %s", currentBranch)
	}

	prompt := buildPrompt("Team Lead", task, tasks, `
Role: Team Lead
Rules:
- Follow AGENTS.md and TECH.md.
- Enforcement: Verify default styling matches flexbox limits (padding:10px, border:1px solid grey).
- Board Management: TASKS.md is single source of truth. Groom Backlog items into active Agent lanes; set Owner, Branch, and Status: In Progress.
- Constraints: No code implementation during grooming. Do not review agent branches or merge them (Agents merge own completed work). Maintain sensible backlog priorities.
`)

	return runCodex(ctx, teamLeadPath, prompt)
}

func buildPrompt(role string, task Task, tasks []Task, roleInstructions string) string {
	agents := mustRead(agentsFile)
	tech := mustRead(techFile)
	taskContext := buildTaskContext(role, task, tasks)

	return fmt.Sprintf(`
You are running inside the multi-agent development workflow.

Active role: %s

================ AGENTS.md ================

%s

================ TECH.md ================

%s

================ TASKS.md CONTEXT ================

%s

================ ACTIVE TASK ================

Section: %s
Title: %s
Owner: %s
Branch: %s
Status: %s

Task body:

%s

================ ROLE INSTRUCTIONS ================

%s
`, role, agents, tech, taskContext, task.Section, task.Title, task.Owner, task.Branch, task.Status, task.Body, roleInstructions)
}

func buildTaskContext(role string, activeTask Task, tasks []Task) string {
	var b strings.Builder
	b.WriteString("Active task body is shown separately below. This context is summarized to save tokens.\n")

	switch role {
	case "Team Lead":
		b.WriteString("\nBacklog:\n")
		writeTaskSummaries(&b, tasks, func(task Task) bool {
			return task.Section == "Backlog" && (task.Status == "Backlog" || task.Status == "")
		})

		b.WriteString("\nImplementation lanes:\n")
		writeTaskSummaries(&b, tasks, func(task Task) bool {
			return strings.HasSuffix(task.Section, "In Progress")
		})

	default:
		b.WriteString("\nOther active implementation work:\n")
		writeTaskSummaries(&b, tasks, func(task Task) bool {
			if task.Title == activeTask.Title &&
				task.Owner == activeTask.Owner &&
				task.Branch == activeTask.Branch {
				return false
			}
			return strings.HasSuffix(task.Section, "In Progress")
		})

		b.WriteString("\nBacklog titles:\n")
		writeTaskSummaries(&b, tasks, func(task Task) bool {
			return task.Section == "Backlog" && (task.Status == "Backlog" || task.Status == "")
		})
	}

	return strings.TrimSpace(b.String())
}

func writeTaskSummaries(b *strings.Builder, tasks []Task, include func(Task) bool) {
	wrote := false
	for _, task := range tasks {
		if !include(task) {
			continue
		}

		wrote = true
		fmt.Fprintf(b, "- %s | Owner: %s | Branch: %s | Status: %s | %s\n",
			task.Title,
			emptyAs(task.Owner, "Unassigned"),
			emptyAs(task.Branch, "(none)"),
			emptyAs(task.Status, "(none)"),
			taskSummary(task.Body),
		)
	}

	if !wrote {
		b.WriteString("- none\n")
	}
}

func taskSummary(body string) string {
	lines := strings.Split(body, "\n")
	var pruned []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		// Strip heavy structural markdown junk that doesn't add context
		if strings.HasPrefix(trimmed, "### ") ||
			strings.HasPrefix(trimmed, "Owner:") ||
			strings.HasPrefix(trimmed, "Branch:") ||
			strings.HasPrefix(trimmed, "Status:") {
			continue
		}
		// Strip common Markdown checklist brackets [ ] or [x] to save tokens
		trimmed = strings.NewReplacer("[ ] ", "", "[x] ", "", "- ", "").Replace(trimmed)

		if trimmed != "" {
			pruned = append(pruned, trimmed)
		}
		if len(pruned) >= 2 { // Keep a strict max of 2 meaningful context lines per background task
			break
		}
	}

	if len(pruned) == 0 {
		return "no summary"
	}
	return strings.Join(pruned, " | ")
}

func runCodex(ctx context.Context, workspace string, prompt string) SessionOutcome {
	logEvent("running codex in %s", workspace)

	cmd := exec.CommandContext(ctx, "codex", "exec", "--sandbox", "danger-full-access", prompt)
	cmd.Dir = workspace
	cmd.Stdin = os.Stdin

	logOutput, err := openLogOutput()
	if err != nil {
		return sessionFailed
	}
	defer logOutput.Close()

	cmd.Stdout = logOutput
	cmd.Stderr = logOutput

	err = cmd.Run()

	if ctx.Err() == context.Canceled {
		logEvent("codex session cancelled in %s", workspace)
		return sessionCancelled
	}

	if err != nil {
		logEvent("codex failed in %s: %v", workspace, err)
		return sessionFailed
	}

	logEvent("codex completed in %s", workspace)
	return sessionCompleted
}

func taskFingerprint(task Task) string {
	return strings.Join([]string{
		task.Section,
		task.Title,
		task.Owner,
		task.Branch,
		task.Status,
		task.Body,
	}, "\x00")
}

func readTasks(path string) ([]Task, error) {
	data, err := fileCache.Read(path)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(data, "\n")

	var tasks []Task
	var current *Task
	var body []string
	currentSection := ""

	flush := func() {
		if current == nil {
			return
		}

		current.Body = strings.TrimSpace(strings.Join(body, "\n"))
		current.Key = taskFingerprint(*current)
		tasks = append(tasks, *current)

		current = nil
		body = nil
	}

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "## ") && !strings.HasPrefix(trimmed, "### ") {
			flush()
			currentSection = strings.TrimSpace(strings.TrimPrefix(trimmed, "## "))

			// TOKEN OPTIMIZATION: Stop parsing if we hit historical/completed logs
			if currentSection == "Done" || currentSection == "Archive" || currentSection == "Completed" {
				currentSection = "SKIP"
			}
			continue
		}

		// Fast-forward past skipped sections entirely
		if currentSection == "SKIP" {
			continue
		}

		if strings.HasPrefix(trimmed, "### ") {
			flush()

			current = &Task{
				Section: currentSection,
				Title:   strings.TrimSpace(strings.TrimPrefix(trimmed, "### ")),
			}

			body = append(body, line)
			continue
		}

		if current == nil {
			continue
		}

		switch {
		case strings.HasPrefix(trimmed, "Owner:"):
			current.Owner = strings.TrimSpace(strings.TrimPrefix(trimmed, "Owner:"))
		case strings.HasPrefix(trimmed, "Branch:"):
			current.Branch = strings.TrimSpace(strings.TrimPrefix(trimmed, "Branch:"))
		case strings.HasPrefix(trimmed, "Status:"):
			current.Status = strings.TrimSpace(strings.TrimPrefix(trimmed, "Status:"))
		}

		body = append(body, line)
	}

	flush()
	return tasks, nil
}
