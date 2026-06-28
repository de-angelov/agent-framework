package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
)

const pollInterval = 10 * time.Second

var (
	repoRoot       = mustResolveRepoRoot()
	workspacesRoot = filepath.Join(repoRoot, "workspaces")
	logsRoot       = filepath.Join(repoRoot, "logs")
	logFilePath    = filepath.Join(logsRoot, "orchestrator.log")

	tasksFile  = filepath.Join(repoRoot, "TASKS.md")
	agentsFile = filepath.Join(repoRoot, "AGENTS.md")
	techFile   = filepath.Join(repoRoot, "TECH.md")

	teamLeadPath = filepath.Join(workspacesRoot, "repo-tl")
	agent1Path   = filepath.Join(workspacesRoot, "repo-agent-1")
	agent2Path   = filepath.Join(workspacesRoot, "repo-agent-2")

	// Protects the running sessions map from concurrent map read/write panics
	mu       sync.Mutex
	running  = map[string]RunningSession{}
	finished = map[string]FinishedSession{}
	logMu    sync.Mutex

	ansiColorRe = regexp.MustCompile(`\x1b\[[0-9;]*m`)
)

type Task struct {
	Section string
	Title   string
	Owner   string
	Branch  string
	Status  string
	Body    string
}

type RunningSession struct {
	Role    string
	Task    string
	TaskKey string
	Branch  string
	Cancel  context.CancelFunc
}

type FinishedSession struct {
	Role       string
	Task       string
	TaskKey    string
	Branch     string
	Outcome    SessionOutcome
	FinishedAt time.Time
}

type SessionOutcome string

const (
	sessionCompleted SessionOutcome = "completed"
	sessionFailed    SessionOutcome = "failed"
	sessionCancelled SessionOutcome = "cancelled"
	sessionUnknown   SessionOutcome = "unknown"
)

func main() {
	mustMkdir(logsRoot)
	logEvent("orchestrator started")
	logEvent("repo root: %s", repoRoot)

	for {
		tasks, err := readTasks(tasksFile)
		if err != nil {
			logEvent("failed to read TASKS.md: %v", err)
			renderUI(nil, err)
			sleep()
			continue
		}

		reconcile(tasks)
		renderUI(tasks, nil)
		sleep()
	}
}

func mustResolveRepoRoot() string {
	cwd, err := filepath.Abs(".")
	if err != nil {
		return "."
	}
	return cwd
}

func reconcile(tasks []Task) {
	desired := map[string]Task{}
	var backlogTask *Task

	for _, task := range tasks {
		switch {
		case task.Section == "Agent 1 In Progress" &&
			task.Owner == "Agent 1" &&
			task.Status == "In Progress":
			desired["Agent 1"] = task

		case task.Section == "Agent 2 In Progress" &&
			task.Owner == "Agent 2" &&
			task.Status == "In Progress":
			desired["Agent 2"] = task

		case task.Section == "Backlog" && backlogTask == nil:
			if task.Status == "Backlog" || task.Status == "" {
				copy := task
				backlogTask = &copy
			}
		}
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
		taskKey := taskFingerprint(task)

		if !stillDesired || taskKey != session.TaskKey {
			logEvent("stopping %s because TASKS.md changed", role)
			session.Cancel()
			delete(running, role)
		}
	}
	mu.Unlock()

	for role, task := range desired {
		taskKey := taskFingerprint(task)

		if role != "Team Lead" && !workspaceExists(role) {
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

		if exists {
			continue
		}

		if alreadyFinished {
			continue
		}

		startSession(role, task)
	}
}

func workspaceExists(role string) bool {
	switch role {
	case "Agent 1":
		_, err := os.Stat(agent1Path)
		return err == nil
	case "Agent 2":
		_, err := os.Stat(agent2Path)
		return err == nil
	case "Team Lead":
		_, err := os.Stat(teamLeadPath)
		return err == nil
	default:
		return false
	}
}

func startSession(role string, task Task) {
	ctx, cancel := context.WithCancel(context.Background())
	taskKey := taskFingerprint(task)

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
			outcome = runAgent(ctx, "Agent 1", agent1Path, task)

		case "Agent 2":
			outcome = runAgent(ctx, "Agent 2", agent2Path, task)

		case "Team Lead":
			outcome = runTeamLead(ctx, task)
		}
	}()
}

func runAgent(ctx context.Context, role string, workspace string, task Task) SessionOutcome {
	logEvent("starting %s on %s in %s", role, task.Title, workspace)

	if task.Branch != "" {
		runGit(workspace, "fetch", "--all", "--prune")
		if branchExists(workspace, task.Branch) {
			runGit(workspace, "checkout", task.Branch)
			runGit(workspace, "pull", "--rebase", "origin", task.Branch)
		} else {
			logEvent("branch does not exist yet; letting agent create: %s", task.Branch)
		}
	}

	prompt := buildPrompt(role, task, `
You are an implementation agent.

Rules:
- Follow AGENTS.md.
- Follow TECH.md.
- Keep UI styling minimal until core workflows are complete.
- Use only flexbox layouts, padding: 10px, and border: 1px solid grey for default styling.
- Do not add custom fonts, shadows, gradients, rounded corners, or decorative spacing unless explicitly required.
- Work only on your assigned task.
- Do not modify backlog priority.
- Do not approve your own work.
- Keep changes focused.
- Write or update tests where appropriate.
- Run focused verification.
- Commit your completed work.
- Push your branch.
- Squash-merge your branch into product main when done.
- Move this task to Done and set Status: Done when complete.
`)

	return runCodex(ctx, workspace, prompt)
}

func runTeamLead(ctx context.Context, task Task) SessionOutcome {
	runGit(teamLeadPath, "fetch", "--all", "--prune")
	currentBranch := currentBranchName(teamLeadPath)
	if currentBranch != "" {
		logEvent("team lead branch: %s", currentBranch)
	}

	prompt := buildPrompt("Team Lead", task, `
You are the Team Lead.

Rules:
- Follow AGENTS.md.
- Follow TECH.md.
- Keep UI styling minimal until core workflows are complete.
- Verify that default styling stays limited to flexbox layouts, padding: 10px, and border: 1px solid grey.
- If the active task is in Backlog, groom it into the correct Agent 1 or Agent 2 lane, set Owner, Branch, and Status: In Progress, and keep backlog priority sensible.
- Do not review implementation branches.
- Do not merge agent branches.
- Do not implement feature code during grooming.
- Assign a new non-overlapping task from Backlog if appropriate.
`)

	return runCodex(ctx, teamLeadPath, prompt)
}

func currentBranchName(workspace string) string {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = workspace
	out, err := cmd.Output()
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(out))
}

func buildPrompt(role string, task Task, roleInstructions string) string {
	agents := mustRead(agentsFile)
	tech := mustRead(techFile)
	taskContext := buildTaskContext(role, task)

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

func buildTaskContext(role string, activeTask Task) string {
	tasks, err := readTasks(tasksFile)
	if err != nil {
		return fmt.Sprintf("[failed to read task context: %v]", err)
	}

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
			if task.Title == activeTask.Title && task.Owner == activeTask.Owner && task.Branch == activeTask.Branch {
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
	for i, line := range lines {
		if strings.TrimSpace(line) != "Outcome:" {
			continue
		}

		for _, candidate := range lines[i+1:] {
			trimmed := strings.TrimSpace(candidate)
			if trimmed != "" {
				return trimmed
			}
		}
	}

	for _, line := range strings.Split(body, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" ||
			strings.HasPrefix(trimmed, "### ") ||
			strings.HasPrefix(trimmed, "Owner:") ||
			strings.HasPrefix(trimmed, "Branch:") ||
			strings.HasPrefix(trimmed, "Status:") ||
			strings.HasSuffix(trimmed, ":") {
			continue
		}
		return trimmed
	}
	return "(no summary)"
}

func emptyAs(value string, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}

func runCodex(ctx context.Context, workspace string, prompt string) SessionOutcome {
	logEvent("running codex in %s", workspace)

	cmd := exec.CommandContext(
		ctx,
		"codex",
		"exec",
		"--sandbox", "danger-full-access",
		prompt,
	)

	cmd.Dir = workspace
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err := cmd.Run()

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
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")

	var tasks []Task
	var current *Task
	var body []string
	currentSection := ""

	flush := func() {
		if current == nil {
			return
		}

		current.Body = strings.TrimSpace(strings.Join(body, "\n"))
		tasks = append(tasks, *current)

		current = nil
		body = nil
	}

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "## ") && !strings.HasPrefix(trimmed, "### ") {
			flush()
			currentSection = strings.TrimSpace(strings.TrimPrefix(trimmed, "## "))
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

func mustRead(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Sprintf("[failed to read %s: %v]", path, err)
	}

	return string(data)
}

func runGit(workspace string, args ...string) {
	logEvent("git %s [%s]", strings.Join(args, " "), workspace)

	cmd := exec.Command("git", args...)
	cmd.Dir = workspace
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		logEvent("git failed in %s: %v", workspace, err)
	}
}

func branchExists(dir string, branch string) bool {
	cmd := exec.Command("git", "show-ref", "--verify", "refs/heads/"+branch)
	cmd.Dir = dir
	return cmd.Run() == nil
}

func remoteBranchExists(dir string, ref string) bool {
	cmd := exec.Command("git", "show-ref", "--verify", "refs/remotes/"+ref)
	cmd.Dir = dir
	return cmd.Run() == nil
}

func sleep() {
	time.Sleep(pollInterval)
}

func mustMkdir(path string) {
	if err := os.MkdirAll(path, 0755); err != nil {
		panic(fmt.Sprintf("failed to create directory %s: %v", path, err))
	}
}

func logEvent(format string, args ...any) {
	logMu.Lock()
	defer logMu.Unlock()

	line := fmt.Sprintf("%s %s\n", time.Now().Format(time.RFC3339), fmt.Sprintf(format, args...))
	fmt.Print(line)

	f, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("log write failed:", err)
		return
	}
	defer f.Close()

	_, _ = f.WriteString(line)
}

func renderUI(tasks []Task, readErr error) {
	mu.Lock()
	rows := buildRows(tasks, running, finished)
	mu.Unlock()
	latest := latestLogLine()

	var b strings.Builder
	b.WriteString("\x1b[2J\x1b[H")
	b.WriteString("Orchestrator\n")
	b.WriteString("Board: " + filepath.Base(repoRoot) + "\n")
	b.WriteString("Time: " + time.Now().Format(time.RFC3339) + "\n")
	if readErr != nil {
		b.WriteString("TASKS error: " + readErr.Error() + "\n")
	}
	b.WriteString("\n")
	b.WriteString("ROLE       STATUS            TASK                           BRANCH\n")
	b.WriteString("---------  ----------------  -----------------------------  ------------------------------\n")
	for _, row := range rows {
		b.WriteString(fmt.Sprintf("%s %-9s %-16s %-29s %s\n",
			row.marker,
			row.role,
			colorStatus(row.status),
			truncate(row.task, 29),
			truncate(row.branch, 30)))
	}
	b.WriteString("\n")
	b.WriteString("Latest: " + latest + "\n")
	b.WriteString("Logs: " + logFilePath + "\n")
	fmt.Print(b.String())
}

func buildRows(tasks []Task, sessions map[string]RunningSession, finishedSessions map[string]FinishedSession) []struct {
	marker string
	role   string
	status string
	task   string
	branch string
} {
	rows := []struct {
		marker string
		role   string
		status string
		task   string
		branch string
	}{
		{role: "Team Lead"},
		{role: "Agent 1"},
		{role: "Agent 2"},
	}

	for i := range rows {
		if session, ok := sessions[rows[i].role]; ok {
			rows[i].marker = ">"
			rows[i].status = "running"
			rows[i].task = session.Task
			rows[i].branch = session.Branch
			continue
		}

		task := findDesiredTaskForRole(tasks, rows[i].role)
		if task.Title != "" {
			if finishedSession, ok := finishedSessions[rows[i].role]; ok && finishedSession.TaskKey == taskFingerprint(task) {
				rows[i].status = string(finishedSession.Outcome)
			} else {
				rows[i].status = task.Status
			}
			rows[i].task = task.Title
			rows[i].branch = task.Branch
			continue
		}

		if rows[i].role == "Team Lead" && hasBacklog(tasks) {
			rows[i].status = "backlog pending"
		} else {
			rows[i].status = "idle"
		}
	}

	return rows
}

func findDesiredTaskForRole(tasks []Task, role string) Task {
	for _, task := range tasks {
		switch role {
		case "Agent 1":
			if task.Section == "Agent 1 In Progress" && task.Owner == "Agent 1" && task.Status == "In Progress" {
				return task
			}
		case "Agent 2":
			if task.Section == "Agent 2 In Progress" && task.Owner == "Agent 2" && task.Status == "In Progress" {
				return task
			}
		case "Team Lead":
			if task.Section == "Backlog" && (task.Status == "Backlog" || task.Status == "") {
				return task
			}
		}
	}
	return Task{}
}

func hasBacklog(tasks []Task) bool {
	for _, task := range tasks {
		if task.Section == "Backlog" && (task.Status == "Backlog" || task.Status == "") {
			return true
		}
	}
	return false
}

func truncate(s string, limit int) string {
	if len(s) <= limit {
		return s
	}
	if limit <= 3 {
		return s[:limit]
	}
	return s[:limit-3] + "..."
}

func colorStatus(status string) string {
	switch {
	case status == "running":
		return "\x1b[32mrunning\x1b[0m"
	case status == "idle":
		return "\x1b[90midle\x1b[0m"
	case status == "backlog pending":
		return "\x1b[33mbacklog pending\x1b[0m"
	case status == "In Progress":
		return "\x1b[36mIn Progress\x1b[0m"
	case status == "completed":
		return "\x1b[32mcompleted\x1b[0m"
	case status == "failed":
		return "\x1b[31mfailed\x1b[0m"
	case status == "Backlog":
		return "\x1b[90mBacklog\x1b[0m"
	default:
		return status
	}
}

func latestLogLine() string {
	data, err := os.ReadFile(logFilePath)
	if err != nil {
		return "[no log]"
	}

	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(lines) == 0 {
		return "[no log]"
	}

	last := ansiColorRe.ReplaceAllString(lines[len(lines)-1], "")
	return last
}
