package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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
	mu      sync.Mutex
	running = map[string]RunningSession{}
	logMu   sync.Mutex
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
	Role   string
	Task   string
	Cancel context.CancelFunc
}

func main() {
	mustMkdir(logsRoot)
	fmt.Println("orchestrator started")
	fmt.Println("repo root:", repoRoot)
	logEvent("orchestrator started")
	logEvent("repo root: %s", repoRoot)

	for {
		tasks, err := readTasks(tasksFile)
		if err != nil {
			fmt.Println("failed to read TASKS.md:", err)
			logEvent("failed to read TASKS.md: %v", err)
			printStatusTable(nil)
			sleep()
			continue
		}

		reconcile(tasks)
		printStatusTable(tasks)
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
			task.Status == "Assigned":
			desired["Agent 1"] = task

		case task.Section == "Agent 2 In Progress" &&
			task.Owner == "Agent 2" &&
			task.Status == "Assigned":
			desired["Agent 2"] = task

		case task.Section == "Ready For Review" &&
			task.Status == "Ready For Review":
			desired["Team Lead"] = task

		case task.Section == "Backlog" && backlogTask == nil:
			if task.Status == "Backlog" || task.Status == "" {
				copy := task
				backlogTask = &copy
			}
		}
	}

	if _, hasActiveReview := desired["Team Lead"]; !hasActiveReview && backlogTask != nil {
		desired["Team Lead"] = *backlogTask
	}

	mu.Lock()
	for role, session := range running {
		task, stillDesired := desired[role]

		if !stillDesired || task.Title != session.Task {
			fmt.Println("stopping", role, "because TASKS.md changed")
			logEvent("stopping %s because TASKS.md changed", role)
			session.Cancel()
			delete(running, role)
		}
	}
	mu.Unlock()

	for role, task := range desired {
		if role != "Team Lead" && !workspaceExists(role) {
			fmt.Println("skipping", role, "because workspace is missing")
			logEvent("skipping %s because workspace is missing", role)
			continue
		}

		mu.Lock()
		_, exists := running[role]
		mu.Unlock()

		if exists {
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

	mu.Lock()
	running[role] = RunningSession{
		Role:   role,
		Task:   task.Title,
		Cancel: cancel,
	}
	mu.Unlock()
	logEvent("starting %s on %s", role, task.Title)

	go func() {
		defer func() {
			mu.Lock()
			delete(running, role)
			mu.Unlock()
		}()

		switch role {
		case "Agent 1":
			runAgent(ctx, "Agent 1", agent1Path, task)

		case "Agent 2":
			runAgent(ctx, "Agent 2", agent2Path, task)

		case "Team Lead":
			runTeamLead(ctx, task)
		}
	}()
}

func runAgent(ctx context.Context, role string, workspace string, task Task) {
	fmt.Println("starting", role, "on", task.Title)
	logEvent("starting %s on %s in %s", role, task.Title, workspace)

	if task.Branch != "" {
		runGit(workspace, "fetch", "--all", "--prune")
		if branchExists(workspace, task.Branch) {
			runGit(workspace, "checkout", task.Branch)
			runGit(workspace, "pull", "--rebase", "origin", task.Branch)
		} else {
			fmt.Println("branch does not exist yet; letting agent create:", task.Branch)
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
- Do not move tasks into Done.
- Do not merge branches.
- Do not approve your own work.
- Keep changes focused.
- Write or update tests where appropriate.
- Run focused verification.
- Commit your completed work.
- Push your branch.
- Move this task to the Ready For Review section and set Status: Ready For Review when complete.
`)

	runCodex(ctx, workspace, prompt)
}

func runTeamLead(ctx context.Context, task Task) {
	fmt.Println("starting Team Lead review for", task.Title)

	runGit(teamLeadPath, "fetch", "--all", "--prune")
	runGit(teamLeadPath, "checkout", "main")
	runGit(teamLeadPath, "pull", "--rebase", "origin", "main")

	prompt := buildPrompt("Team Lead", task, `
You are the Team Lead.

Rules:
- Follow AGENTS.md.
- Follow TECH.md.
- Keep UI styling minimal until core workflows are complete.
- Verify that default styling stays limited to flexbox layouts, padding: 10px, and border: 1px solid grey.
- If the active task is in Backlog, groom it into the correct Agent 1 or Agent 2 lane, set Owner, Branch, and Status: Assigned, and keep backlog priority sensible.
- If the active task is in Ready For Review, review the implementation branch.
- Fetch and inspect the implementation branch before review.
- Run full verification before approval.
- If approved, merge to main and move the task to Done.
- Add Completed: YYYY-MM-DD when moving to Done.
- If rejected, return the task to the assigned agent lane and add a [REJECTED] section.
- Do not implement feature code during review.
- Do not silently resolve merge conflicts.
- Assign a new non-overlapping task from Backlog if appropriate.
`)

	runCodex(ctx, teamLeadPath, prompt)
}

func buildPrompt(role string, task Task, roleInstructions string) string {
	agents := mustRead(agentsFile)
	tech := mustRead(techFile)

	// Filter out the completed tasks to save tokens
	filteredTasks := readFilteredTasks(tasksFile)

	return fmt.Sprintf(`
You are running inside the multi-agent development workflow.

Active role: %s

================ AGENTS.md ================

%s

================ TECH.md ================

%s

================ TASKS.md (Active & Backlog only) ================

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
`, role, agents, tech, filteredTasks, task.Section, task.Title, task.Owner, task.Branch, task.Status, task.Body, roleInstructions)
}

// readFilteredTasks reads TASKS.md but reconstructs it without the "Done" section.
func readFilteredTasks(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Sprintf("[failed to read %s: %v]", path, err)
	}

	lines := strings.Split(string(data), "\n")
	var keptLines []string
	skipSection := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// If we encounter a top-level heading, check if it's the "Done" section
		if strings.HasPrefix(trimmed, "## ") && !strings.HasPrefix(trimmed, "### ") {
			sectionName := strings.ToLower(strings.TrimSpace(strings.TrimPrefix(trimmed, "## ")))
			if sectionName == "done" || sectionName == "completed" {
				skipSection = true
			} else {
				skipSection = false
			}
		}

		if !skipSection {
			keptLines = append(keptLines, line)
		}
	}

	return strings.Join(keptLines, "\n")
}

func runCodex(ctx context.Context, workspace string, prompt string) {
	fmt.Println("running codex in", workspace)

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
		fmt.Println("codex session cancelled in", workspace)
		logEvent("codex session cancelled in %s", workspace)
		return
	}

	if err != nil {
		fmt.Println("codex failed:", err)
		logEvent("codex failed in %s: %v", workspace, err)
		return
	}

	logEvent("codex completed in %s", workspace)
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
	fmt.Println("git", strings.Join(args, " "))
	logEvent("git %s [%s]", strings.Join(args, " "), workspace)

	cmd := exec.Command("git", args...)
	cmd.Dir = workspace
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Println("git failed:", err)
		logEvent("git failed in %s: %v", workspace, err)
	}
}

func branchExists(dir string, branch string) bool {
	cmd := exec.Command("git", "show-ref", "--verify", "refs/heads/"+branch)
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

func printStatusTable(tasks []Task) {
	mu.Lock()
	defer mu.Unlock()

	rows := []struct {
		role      string
		status    string
		task      string
		branch    string
		workspace string
	}{
		{role: "Team Lead", workspace: teamLeadPath},
		{role: "Agent 1", workspace: agent1Path},
		{role: "Agent 2", workspace: agent2Path},
	}

	for i := range rows {
		if session, ok := running[rows[i].role]; ok {
			rows[i].status = "running"
			rows[i].task = session.Task
			continue
		}

		task := findDesiredTaskForRole(tasks, rows[i].role)
		if task.Title != "" {
			rows[i].status = task.Status
			rows[i].task = task.Title
			rows[i].branch = task.Branch
		} else if rows[i].role == "Team Lead" && hasBacklog(tasks) {
			rows[i].status = "backlog pending"
		} else {
			rows[i].status = "idle"
		}
	}

	fmt.Println()
	fmt.Println("Current Work")
	fmt.Println("ROLE      | STATUS          | TASK                          | BRANCH")
	fmt.Println("----------+-----------------+-------------------------------+------------------------------")
	for _, row := range rows {
		fmt.Printf("%-9s | %-15s | %-29s | %s\n", row.role, truncate(row.status, 15), truncate(row.task, 29), row.branch)
	}
	fmt.Println()
}

func findDesiredTaskForRole(tasks []Task, role string) Task {
	for _, task := range tasks {
		switch role {
		case "Agent 1":
			if task.Section == "Agent 1 In Progress" && task.Owner == "Agent 1" && task.Status == "Assigned" {
				return task
			}
		case "Agent 2":
			if task.Section == "Agent 2 In Progress" && task.Owner == "Agent 2" && task.Status == "Assigned" {
				return task
			}
		case "Team Lead":
			if task.Section == "Ready For Review" && task.Status == "Ready For Review" {
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
