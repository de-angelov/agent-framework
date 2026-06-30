package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

var repoRootMarkers = []string{
	"BACKLOG.md",
	"TASKS.md",
	"ARCHIVE.md",
	"AGENTS.md",
	"DEV_AGENT.md",
	"TEAM_LEAD_AGENT.md",
	"TECH.md",
}

func main() {
	mustMkdir(logsRoot)
	logFilePath = newRunLogFilePath(time.Now())
	logEvent("orchestrator started")
	logEvent("repo root: %s", repoRoot)

	for {
		tasks, err := readBoardTasks()
		if err != nil {
			logEvent("failed to read board files: %v", err)
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
	cwd, err := os.Getwd()
	if err != nil {
		return "."
	}

	if root, ok := resolveRepoRootFrom(cwd); ok {
		return root
	}

	if executable, err := os.Executable(); err == nil {
		if root, ok := resolveRepoRootFrom(filepath.Dir(executable)); ok {
			return root
		}
	}

	abs, err := filepath.Abs(cwd)
	if err != nil {
		return cwd
	}
	return abs
}

func resolveRepoRootFrom(start string) (string, bool) {
	current, err := filepath.Abs(start)
	if err != nil {
		return "", false
	}

	for {
		if hasRepoRootMarkers(current) {
			return current, true
		}

		parent := filepath.Dir(current)
		if parent == current {
			return "", false
		}
		current = parent
	}
}

func hasRepoRootMarkers(dir string) bool {
	for _, marker := range repoRootMarkers {
		if _, err := os.Stat(filepath.Join(dir, marker)); err != nil {
			return false
		}
	}
	return true
}

func sleep() {
	time.Sleep(pollInterval)
}

func mustMkdir(path string) {
	if err := os.MkdirAll(path, 0755); err != nil {
		panic(fmt.Sprintf("failed to create directory %s: %v", path, err))
	}
}

func defaultLogFilePath() string {
	return filepath.Join(logsRoot, "orchestrator.log")
}

func newRunLogFilePath(now time.Time) string {
	return filepath.Join(logsRoot, fmt.Sprintf(
		"orchestrator-%s.log",
		now.Format("20060102-150405.000000000"),
	))
}
