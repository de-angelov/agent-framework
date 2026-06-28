package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const defaultRepoURL = "https://github.com/de-angelov/agent-task-test"
const defaultRepoSSHURL = "git@github.com:de-angelov/agent-task-test.git"

func main() {
	root, err := os.Getwd()
	if err != nil {
		log.Fatalf("Fatal: failed to get working directory: %v", err)
	}

	fmt.Println("Setting up AI development workflow...")
	fmt.Println("Root:", root)

	workspacesRoot := filepath.Join(root, "workspaces")
	mustMkdir(workspacesRoot)

	createClone(workspacesRoot, "repo-tl")
	createClone(workspacesRoot, "repo-agent-1")
	createClone(workspacesRoot, "repo-agent-2")

	fmt.Println()
	fmt.Println("Setup complete.")
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println("  cd workspaces/repo-tl && go run orchestrator.go")
	fmt.Println("  cd workspaces/repo-agent-1 && go run orchestrator.go")
	fmt.Println("  cd workspaces/repo-agent-2 && go run orchestrator.go")
}

func createClone(root, dir string) {
	path := filepath.Join(root, dir)
	if _, err := os.Stat(path); err == nil {
		fmt.Println("• exists:", dir)
		ensureSSHRemote(path)
		ensureMainBranch(path)
		removeWorkspaceTaskBoard(path)
		return
	}

	run(root, "git", "clone", defaultRepoURL, dir)
	fmt.Println("✓ clone:", dir)
	ensureSSHRemote(path)
	ensureMainBranch(path)
	removeWorkspaceTaskBoard(path)
}

func ensureSSHRemote(repoPath string) {
	run(repoPath, "git", "remote", "set-url", "origin", defaultRepoSSHURL)
	fmt.Println("✓ remote:", defaultRepoSSHURL)
}

func ensureMainBranch(repoPath string) {
	run(repoPath, "git", "fetch", "origin", "main")
	run(repoPath, "git", "remote", "set-head", "origin", "main")
	run(repoPath, "git", "checkout", "main")
	run(repoPath, "git", "branch", "--set-upstream-to", "origin/main", "main")
	fmt.Println("✓ default branch: main")
}

func removeWorkspaceTaskBoard(repoPath string) {
	taskBoard := filepath.Join(repoPath, "TASKS.md")
	if err := os.Remove(taskBoard); err == nil {
		fmt.Println("✓ removed:", filepath.Join(filepath.Base(repoPath), "TASKS.md"))
	}
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

func mustMkdir(path string) {
	if err := os.MkdirAll(path, 0755); err != nil {
		log.Fatalf("Fatal: failed to create directory %s: %v", path, err)
	}
	fmt.Println("✓ dir:", filepath.Base(path))
}
