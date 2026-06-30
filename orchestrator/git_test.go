package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestPrepareBranchPreservesDirtyWorkspaceOnAssignedBranch(t *testing.T) {
	withRetryTestState(t)

	workspace := initGitWorkspace(t)
	runTestGit(t, workspace, "checkout", "-b", "agent/1/current-task")

	dirtyPath := filepath.Join(workspace, "work.txt")
	if err := os.WriteFile(dirtyPath, []byte("dirty work\n"), 0644); err != nil {
		t.Fatal(err)
	}

	if err := prepareBranch(workspace, "agent/1/current-task"); err != nil {
		t.Fatal(err)
	}

	if got := currentBranchName(workspace); got != "agent/1/current-task" {
		t.Fatalf("current branch = %q, want assigned branch", got)
	}

	status := runTestGitOutput(t, workspace, "status", "--porcelain")
	if !strings.Contains(status, "work.txt") {
		t.Fatalf("expected dirty work to remain in workspace, status = %q", status)
	}

	branches := runTestGitOutput(t, workspace, "branch", "--list", "wip/orchestrator/*")
	if strings.TrimSpace(branches) != "" {
		t.Fatalf("expected no WIP branch on restart, got %q", branches)
	}
}

func initGitWorkspace(t *testing.T) string {
	t.Helper()

	workspace := t.TempDir()
	runTestGit(t, workspace, "init", "-b", "main")
	runTestGit(t, workspace, "config", "user.email", "test@example.com")
	runTestGit(t, workspace, "config", "user.name", "Test User")

	readmePath := filepath.Join(workspace, "README.md")
	if err := os.WriteFile(readmePath, []byte("# Test\n"), 0644); err != nil {
		t.Fatal(err)
	}
	runTestGit(t, workspace, "add", "README.md")
	runTestGit(t, workspace, "commit", "-m", "initial")

	return workspace
}

func runTestGit(t *testing.T, workspace string, args ...string) {
	t.Helper()

	cmd := exec.Command("git", args...)
	cmd.Dir = workspace
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %s failed: %v\n%s", strings.Join(args, " "), err, output)
	}
}

func runTestGitOutput(t *testing.T, workspace string, args ...string) string {
	t.Helper()

	cmd := exec.Command("git", args...)
	cmd.Dir = workspace
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("git %s failed: %v", strings.Join(args, " "), err)
	}
	return string(output)
}
