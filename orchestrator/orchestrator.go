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
	repoRoot, _ = filepath.Abs(".")

	tasksFile  = filepath.Join(repoRoot, "TASKS.md")
	agentsFile = filepath.Join(repoRoot, "AGENTS.md")
	techFile   = filepath.Join(repoRoot, "TECH.md")

	teamLeadPath = filepath.Join(repoRoot, ".worktrees", "tl")
	agent1Path   = filepath.Join(repoRoot, ".worktrees", "agent-1")
	agent2Path   = filepath.Join(repoRoot, ".worktrees", "agent-2")

	// Protects the running sessions map from concurrent map read/write panics
	mu      sync.Mutex
	running = map[string]RunningSession{}
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
	fmt.Println("orchestrator started")
	fmt.Println("repo root:", repoRoot)

	for {
		tasks, err := readTasks(tasksFile)
		if err != nil {
			fmt.Println("failed to read TASKS.md:", err)
			sleep()
			continue
		}

		reconcile(tasks)
		sleep()
	}
}

func reconcile(tasks []Task) {
	desired := map[string]Task{}

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
		}
	}

	mu.Lock()
	for role, session := range running {
		task, stillDesired := desired[role]

		if !stillDesired || task.Title != session.Task {
			fmt.Println("stopping", role, "because TASKS.md changed")
			session.Cancel()
			delete(running, role)
		}
	}
	mu.Unlock()

	for role, task := range desired {
		mu.Lock()
		_, exists := running[role]
		mu.Unlock()

		if exists {
			continue
		}

		startSession(role, task)
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

	if task.Branch != "" {
		runGit(workspace, "fetch", "--all", "--prune")
		runGit(workspace, "checkout", task.Branch)
		runGit(workspace, "pull", "--rebase", "origin", task.Branch)
	}

	prompt := buildPrompt(role, task, `
You are an implementation agent.

Rules:
- Follow AGENTS.md.
- Follow TECH.md.
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
- Review the task in the Ready For Review section.
- Fetch and inspect the implementation branch.
- Run full verification.
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
		"--approval-mode", "never",
		"--sandbox", "workspace-write",
		prompt,
	)

	cmd.Dir = workspace
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err := cmd.Run()

	if ctx.Err() == context.Canceled {
		fmt.Println("codex session cancelled in", workspace)
		return
	}

	if err != nil {
		fmt.Println("codex failed:", err)
	}
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

	cmd := exec.Command("git", args...)
	cmd.Dir = workspace
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Println("git failed:", err)
	}
}

func sleep() {
	time.Sleep(pollInterval)
}
