package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	root, err := os.Getwd()
	if err != nil {
		log.Fatalf("Fatal: failed to get working directory: %v", err)
	}

	fmt.Println("Setting up AI development workflow...")
	fmt.Println("Root:", root)

	mustMkdir(filepath.Join(root, "orchestrator"))

	createFile(filepath.Join(root, "AGENTS.md"), agentsTemplate())
	createFile(filepath.Join(root, "TECH.md"), techTemplate())
	createFile(filepath.Join(root, "TASKS.md"), tasksTemplate())

	ensureGitRepo(root)
	ensureMainBranch(root)

	// Set up isolated workspaces via Git worktrees
	createWorktree(root, filepath.Join(".worktrees", "tl"), "main", false)
	createWorktree(root, filepath.Join(".worktrees", "agent-1"), "agent/1/current-task", true)
	createWorktree(root, filepath.Join(".worktrees", "agent-2"), "agent/2/current-task", true)

	fmt.Println()
	fmt.Println("Setup complete.")
	fmt.Println()
	fmt.Println("Run orchestrator:")
	fmt.Println("  cd orchestrator")
	fmt.Println("  go run orchestrator.go")
}

func ensureGitRepo(root string) {
	// If a .git directory exists, make sure it actually has a commit history before skipping.
	// Otherwise, worktrees cannot be created.
	if _, err := os.Stat(filepath.Join(root, ".git")); err == nil {
		if hasCommits(root) {
			return
		}
		fmt.Println("• Git repo exists but has no commits. Initializing base history...")
	} else {
		run(root, "git", "init")
	}

	run(root, "git", "add", "AGENTS.md", "TECH.md", "TASKS.md")
	run(root, "git", "commit", "-m", "chore: initialize multi-agent workflow")
}

func ensureMainBranch(root string) {
	run(root, "git", "branch", "-M", "main")
}

func createWorktree(root, dir, branch string, createBranch bool) {
	path := filepath.Join(root, dir)

	if _, err := os.Stat(path); err == nil {
		fmt.Println("• worktree directory already exists:", dir)
		return
	}

	// Use cross-platform slashes for Git CLI predictability
	gitPath := filepath.ToSlash(path)
	args := []string{"worktree", "add"}

	if createBranch {
		// If the branch already exists in Git history, don't pass the "-b" flag
		// which would cause a Git crash. Just check out the existing branch.
		if branchExists(root, branch) {
			fmt.Printf("• branch %s already exists, linking worktree to it...\n", branch)
		} else {
			args = append(args, "-b", branch)
		}
	}

	args = append(args, gitPath, branch)

	// Run worktree command safely
	cmd := exec.Command("git", args...)
	cmd.Dir = root
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	fmt.Println("$ git", strings.Join(args, " "))
	if err := cmd.Run(); err != nil {
		log.Fatalf("Fatal: failed to create worktree %s: %v", dir, err)
	}

	fmt.Println("✓ worktree:", dir, "→", branch)
}

func hasCommits(dir string) bool {
	cmd := exec.Command("git", "rev-parse", "--verify", "HEAD")
	cmd.Dir = dir
	return cmd.Run() == nil
}

func branchExists(dir string, branch string) bool {
	cmd := exec.Command("git", "show-ref", "--verify", "refs/heads/"+branch)
	cmd.Dir = dir
	return cmd.Run() == nil
}

func mustMkdir(path string) {
	if err := os.MkdirAll(path, 0755); err != nil {
		log.Fatalf("Fatal: failed to create directory %s: %v", path, err)
	}
	fmt.Println("✓ dir:", filepath.Base(path))
}

func createFile(path string, contents string) {
	if _, err := os.Stat(path); err == nil {
		fmt.Println("• exists:", filepath.Base(path))
		return
	}

	if err := os.WriteFile(path, []byte(contents), 0644); err != nil {
		log.Fatalf("Fatal: failed to write file %s: %v", path, err)
	}
	fmt.Println("✓ file:", filepath.Base(path))
}

func run(dir string, name string, args ...string) {
	fmt.Println("$", name, strings.Join(args, " "))

	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		log.Fatalf("Fatal: command %s failed: %v", name, err)
	}
}

func agentsTemplate() string {
	return `# AGENTS

Defines how AI agents collaborate in this repository.

## Roles

- Team Lead
- Agent 1
- Agent 2

## Team Lead

Responsible for:

- planning
- task assignment
- review
- verification
- merging approved work
- moving tasks to Done

The Team Lead must not implement feature code during review.

## Agent 1 / Agent 2

Responsible for:

- implementing assigned tasks
- writing tests
- running focused verification
- committing focused changes
- pushing assigned branches
- moving completed work to Ready For Review

Implementation agents must not:

- merge branches
- move tasks to Done
- approve their own work
- reprioritize backlog

## Coordination Files

- AGENTS.md defines workflow.
- TECH.md defines technical standards.
- TASKS.md is the live project board.

## Task Flow

Backlog → Agent In Progress → Ready For Review → Done

## Conflict Prevention

Agents must only edit files required for their assigned task.

The Team Lead should assign non-overlapping tasks.

## Instruction Priority

1. User instructions
2. Repository instructions
3. AGENTS.md
4. TECH.md
5. General engineering best practices
`
}

func techTemplate() string {
	return `# TECH

Defines project-specific technical standards.

## Stack

Fill this in for the project.

Example:

- TypeScript
- React
- React Router
- Drizzle ORM
- SQLite
- Vitest
- Playwright

## Architecture

Keep business logic separate from framework, database, networking, and session code.

Prefer pure functions and dependency injection.

## Coding Standards

- Prefer strict typing.
- Avoid any.
- Use small single-purpose functions.
- Delete dead code.
- Validate inputs at boundaries.

## Testing

Every new feature or business rule should include appropriate tests.

## Verification

Implementation agents should run:

- typecheck
- relevant tests

Team Lead should run:

- full typecheck
- full test suite
- build
`
}

func tasksTemplate() string {
	return `# TASKS

## Backlog

---

## Agent 1 In Progress

---

## Agent 2 In Progress

---

## Ready For Review

---

## Done
`
}
