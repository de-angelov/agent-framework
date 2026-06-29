package main

import (
	"context"
	"path/filepath"
	"regexp"
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

	mu       sync.Mutex
	running  = map[string]RunningSession{}
	finished = map[string]FinishedSession{}
	logMu    sync.Mutex

	fileCache = FileCache{
		items: map[string]cachedFile{},
	}

	ansiColorRe = regexp.MustCompile(`\x1b\[[0-9;]*m`)
)

type Task struct {
	Section string
	Title   string
	Owner   string
	Branch  string
	Status  string
	Body    string
	Key     string
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
