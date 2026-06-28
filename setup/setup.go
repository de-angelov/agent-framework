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
		return
	}

	run(root, "git", "clone", defaultRepoURL, dir)
	fmt.Println("✓ clone:", dir)
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
