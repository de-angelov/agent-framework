package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func workspaceExists(role string) bool {
	switch role {
	case devAgent1Role:
		_, err := os.Stat(agent1Path)
		return err == nil
	case devAgent2Role:
		_, err := os.Stat(agent2Path)
		return err == nil
	case teamLeadRole:
		_, err := os.Stat(teamLeadPath)
		return err == nil
	default:
		return false
	}
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

func runGit(workspace string, args ...string) {
	logEvent("git %s [%s]", strings.Join(args, " "), workspace)

	cmd := exec.Command("git", args...)
	cmd.Dir = workspace

	logOutput, err := openLogOutput()
	if err == nil {
		defer logOutput.Close()
		cmd.Stdout = logOutput
		cmd.Stderr = logOutput
	}

	if err := cmd.Run(); err != nil {
		logEvent("git failed in %s: %v", workspace, err)
	}
}

func prepareBranch(workspace string, branch string) error {
	if err := checkpointDirtyWorkspace(workspace, branch); err != nil {
		return err
	}

	if err := runGitChecked(workspace, "fetch", "--all", "--prune"); err != nil {
		return err
	}

	if branchExists(workspace, branch) {
		if err := runGitChecked(workspace, "checkout", branch); err != nil {
			return err
		}
		if remoteBranchExists(workspace, "origin/"+branch) {
			return runGitChecked(workspace, "pull", "--rebase", "origin", branch)
		}
		return runGitChecked(workspace, "push", "-u", "origin", branch)
	}

	remoteRef := "origin/" + branch
	if remoteBranchExists(workspace, remoteRef) {
		return runGitChecked(workspace, "checkout", "-B", branch, remoteRef)
	}

	if err := runGitChecked(workspace, "checkout", "main"); err != nil {
		return err
	}
	if err := runGitChecked(workspace, "pull", "--rebase", "origin", "main"); err != nil {
		return err
	}
	if err := runGitChecked(workspace, "checkout", "-b", branch); err != nil {
		return err
	}
	return runGitChecked(workspace, "push", "-u", "origin", branch)
}

func checkpointDirtyWorkspace(workspace string, nextBranch string) error {
	if !workspaceHasChanges(workspace) && !mergeInProgress(workspace) {
		return nil
	}

	current := currentBranchName(workspace)
	if current == "" || current == "HEAD" {
		current = "detached"
	}

	wipBranch := fmt.Sprintf(
		"wip/orchestrator/%s/%s/%d",
		sanitizeBranchPart(current),
		sanitizeBranchPart(nextBranch),
		time.Now().Unix(),
	)

	logEvent("checkpointing dirty workspace %s on %s before switching to %s", workspace, wipBranch, nextBranch)

	if mergeInProgress(workspace) {
		_ = runGitChecked(workspace, "merge", "--abort")
		if mergeInProgress(workspace) {
			return fmt.Errorf("workspace has unresolved merge state; manual cleanup required before switching tasks")
		}
	}

	if err := runGitChecked(workspace, "checkout", "-b", wipBranch); err != nil {
		return err
	}
	if err := runGitChecked(workspace, "add", "-A"); err != nil {
		return err
	}
	if workspaceHasStagedChanges(workspace) {
		if err := runGitChecked(workspace, "commit", "-m", "WIP before switching to "+nextBranch); err != nil {
			return err
		}
		if err := runGitChecked(workspace, "push", "-u", "origin", wipBranch); err != nil {
			return err
		}
		logEvent("saved dirty workspace to %s", wipBranch)
		return nil
	}

	logEvent("dirty workspace resolved without commit while switching to %s", nextBranch)
	return nil
}

func workspaceHasChanges(workspace string) bool {
	out, err := gitOutput(workspace, "status", "--porcelain")
	return err == nil && strings.TrimSpace(out) != ""
}

func workspaceHasStagedChanges(workspace string) bool {
	cmd := exec.Command("git", "diff", "--cached", "--quiet")
	cmd.Dir = workspace
	return cmd.Run() != nil
}

func mergeInProgress(workspace string) bool {
	out, err := gitOutput(workspace, "rev-parse", "-q", "--verify", "MERGE_HEAD")
	return err == nil && strings.TrimSpace(out) != ""
}

func gitOutput(workspace string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = workspace
	out, err := cmd.Output()
	return string(out), err
}

func sanitizeBranchPart(value string) string {
	var b strings.Builder
	lastDash := false

	for _, r := range strings.ToLower(value) {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
			lastDash = false
		case r == '/', r == '-', r == '_', r == '.':
			if !lastDash {
				b.WriteRune('-')
				lastDash = true
			}
		default:
			if !lastDash {
				b.WriteRune('-')
				lastDash = true
			}
		}
	}

	result := strings.Trim(b.String(), "-")
	if result == "" {
		return "unknown"
	}

	return strconv.Itoa(len(result)) + "-" + result
}

func runGitChecked(workspace string, args ...string) error {
	logEvent("git %s [%s]", strings.Join(args, " "), workspace)

	cmd := exec.Command("git", args...)
	cmd.Dir = workspace

	logOutput, err := openLogOutput()
	if err != nil {
		return err
	}
	defer logOutput.Close()

	cmd.Stdout = logOutput
	cmd.Stderr = logOutput

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git %s failed: %w", strings.Join(args, " "), err)
	}

	return nil
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
