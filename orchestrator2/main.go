package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
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

func sleep() {
	time.Sleep(pollInterval)
}

func mustMkdir(path string) {
	if err := os.MkdirAll(path, 0755); err != nil {
		panic(fmt.Sprintf("failed to create directory %s: %v", path, err))
	}
}
