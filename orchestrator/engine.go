package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
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
		case task.Section == "Dev Agent 1 In Progress" && task.Owner == devAgent1Role && task.Status == "In Progress":
			addDesired(devAgent1Role, task)
		case task.Section == "Dev Agent 2 In Progress" && task.Owner == devAgent2Role && task.Status == "In Progress":
			addDesired(devAgent2Role, task)
		case task.Section == "Backlog" && backlogTask == nil:
			if task.Status == "Backlog" || task.Status == "" {
				copy := task
				backlogTask = &copy
			}
		}
	}

	if backlogTask != nil && len(invalidRoles) == 0 {
		_, agent1Busy := desired[devAgent1Role]
		_, agent2Busy := desired[devAgent2Role]
		if agent1Busy && agent2Busy {
			backlogTask = nil
		}
	}

	if backlogTask != nil && len(invalidRoles) > 0 {
		backlogTask = nil
	}

	if backlogTask != nil {
		desired[teamLeadRole] = *backlogTask
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
			logEvent("stopping %s because board files changed", role)
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
		if alreadyFinished && finishedSession.Outcome == sessionRateLimited && !codexResourcePauseActiveLocked(time.Now()) {
			delete(finished, role)
			alreadyFinished = false
		}
		mu.Unlock()

		if exists || alreadyFinished {
			continue
		}

		if !codexResourcesAvailable() {
			logEvent("skipping %s because Codex resources are paused", role)
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
		case devAgent1Role:
			outcome = runAgent(ctx, devAgent1Role, agent1Path, task, tasks)
		case devAgent2Role:
			outcome = runAgent(ctx, devAgent2Role, agent2Path, task, tasks)
		case teamLeadRole:
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
Role: Dev Agent
Runtime Rules:
- Work only on the assigned task and keep changes focused.
`)

	return runCodex(ctx, workspace, prompt)
}

func runTeamLead(ctx context.Context, task Task, tasks []Task) SessionOutcome {
	runGit(teamLeadPath, "fetch", "--all", "--prune")

	currentBranch := currentBranchName(teamLeadPath)
	if currentBranch != "" {
		logEvent("team lead branch: %s", currentBranch)
	}

	prompt := buildPrompt(teamLeadRole, task, tasks, `
Role: Team Lead Agent
Runtime Rules:
- No code implementation during grooming. Do not review dev-agent branches or merge them. Maintain sensible backlog priorities.
`)

	return runCodex(ctx, teamLeadPath, prompt)
}

func buildPrompt(role string, task Task, tasks []Task, roleInstructions string) string {
	commonInstructions := mustRead(agentsFile)
	specificInstructions := mustRead(roleInstructionsPath(role))
	tech := mustRead(techFile)
	taskContext := buildTaskContext(role, task, tasks)

	return fmt.Sprintf(`
You are running inside the multi-agent development workflow.

Active role: %s

================ AGENTS.md COMMON RULES ================

%s

================ ROLE-SPECIFIC INSTRUCTIONS ================

%s

================ TECH.md ================

%s

================ BOARD CONTEXT ================

%s

================ ACTIVE TASK ================

Section: %s
Title: %s
Owner: %s
Branch: %s
Status: %s

Task body:

%s

================ RUNTIME INSTRUCTIONS ================

%s
`, role, commonInstructions, specificInstructions, tech, taskContext, task.Section, task.Title, task.Owner, task.Branch, task.Status, task.Body, roleInstructions)
}

func roleInstructionsPath(role string) string {
	if role == teamLeadRole {
		return tlAgentInstructionsFile
	}
	return devAgentInstructionsFile
}

func buildTaskContext(role string, activeTask Task, tasks []Task) string {
	var b strings.Builder
	b.WriteString("Active task body is shown separately below. This context is summarized to save tokens.\n")

	switch role {
	case teamLeadRole:
		b.WriteString("\nBacklog:\n")
		writeTaskSummaries(&b, tasks, func(task Task) bool {
			return task.Section == "Backlog" && (task.Status == "Backlog" || task.Status == "")
		})

		b.WriteString("\nDev-agent lanes:\n")
		writeTaskSummaries(&b, tasks, func(task Task) bool {
			return strings.HasSuffix(task.Section, "In Progress")
		})

	default:
		b.WriteString("\nOther active dev-agent work:\n")
		writeTaskSummaries(&b, tasks, func(task Task) bool {
			if task.Title == activeTask.Title &&
				task.Owner == activeTask.Owner &&
				task.Branch == activeTask.Branch {
				return false
			}
			return strings.HasSuffix(task.Section, "In Progress")
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

func readBoardTasks() ([]Task, error) {
	var all []Task

	for _, path := range []string{backlogFile, tasksFile} {
		tasks, err := readTasks(path)
		if err != nil {
			return nil, err
		}
		all = append(all, tasks...)
	}

	return all, nil
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

func codexResourcesAvailable() bool {
	now := time.Now()

	mu.Lock()
	if codexResourcePauseActiveLocked(now) {
		mu.Unlock()
		return false
	}
	if !lastCodexStatusCheck.IsZero() && now.Sub(lastCodexStatusCheck) < codexStatusCheckInterval {
		mu.Unlock()
		return true
	}
	lastCodexStatusCheck = now
	mu.Unlock()

	output, err := runCodexStatus()
	trimmed := strings.TrimSpace(output)
	if trimmed == "" {
		trimmed = "codex status returned no output"
	}

	if isCodexResourceExhausted(trimmed) {
		pauseCodexResources(codexResourceRetryDelay, "codex status: "+oneLine(trimmed))
		return false
	}

	if err != nil {
		mu.Lock()
		lastCodexStatusMessage = "codex status unavailable: " + oneLine(trimmed)
		mu.Unlock()
		logEvent("codex status unavailable; continuing: %v: %s", err, oneLine(trimmed))
		return true
	}

	mu.Lock()
	lastCodexStatusMessage = "codex status ok: " + oneLine(trimmed)
	mu.Unlock()

	logEvent("codex status ok: %s", oneLine(trimmed))
	return true
}

// codexResourcePauseActiveLocked checks if resource pause is currently active.
// Caller MUST hold mu lock.
func codexResourcePauseActiveLocked(now time.Time) bool {
	return !codexResourcePausedUntil.IsZero() && now.Before(codexResourcePausedUntil)
}

func runCodexStatus() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "codex", "status")
	cmd.Dir = repoRoot
	cmd.Stdin = os.Stdin

	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output

	err := cmd.Run()
	return output.String(), err
}

func pauseCodexResources(delay time.Duration, reason string) {
	retryAt := time.Now().Add(delay)

	mu.Lock()
	if retryAt.After(codexResourcePausedUntil) {
		codexResourcePausedUntil = retryAt
	}
	lastCodexStatusMessage = reason
	mu.Unlock()

	logEvent("pausing Codex launches until %s: %s", retryAt.Format(time.RFC3339), reason)
}

func isCodexResourceExhausted(output string) bool {
	text := strings.ToLower(output)
	resourceWords := []string{
		"rate limit",
		"rate_limit",
		"quota",
		"usage limit",
		"usage_limit",
		"limit reached",
		"too many requests",
		"insufficient credits",
		"out of credits",
		"resources exhausted",
		"resource exhausted",
		"try again later",
		"reset in",
		"resets in",
	}

	for _, word := range resourceWords {
		if strings.Contains(text, word) {
			return true
		}
	}
	return false
}

func oneLine(value string) string {
	fields := strings.Fields(value)
	if len(fields) == 0 {
		return ""
	}
	line := strings.Join(fields, " ")
	if len(line) > 240 {
		return line[:237] + "..."
	}
	return line
}

func runCodex(ctx context.Context, workspace string, prompt string) SessionOutcome {
	logEvent("running codex in %s", workspace)

	cmd := exec.CommandContext(ctx, "codex", "exec", "--sandbox", "danger-full-access", "-")
	cmd.Dir = workspace
	cmd.Stdin = strings.NewReader(prompt)

	var captured bytes.Buffer
	output := io.MultiWriter(lockedLogWriter{}, &captured)
	cmd.Stdout = output
	cmd.Stderr = output

	err := cmd.Run()

	if ctx.Err() == context.Canceled {
		logEvent("codex session cancelled in %s", workspace)
		return sessionCancelled
	}

	if err != nil {
		if isCodexResourceExhausted(captured.String()) {
			pauseCodexResources(codexResourceRetryDelay, "codex exec: "+oneLine(captured.String()))
			logEvent("codex resource limit detected in %s: %v", workspace, err)
			return sessionRateLimited
		}
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
