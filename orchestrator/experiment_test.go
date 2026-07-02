package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestReadExperimentConfigDefaultsAndResolvesPaths(t *testing.T) {
	root := t.TempDir()
	oldRepoRoot := repoRoot
	oldAgent1Path := agent1Path
	repoRoot = root
	agent1Path = initGitWorkspace(t)
	t.Cleanup(func() {
		repoRoot = oldRepoRoot
		agent1Path = oldAgent1Path
	})

	configPath := filepath.Join(root, "experiment.json")
	if err := os.WriteFile(configPath, []byte(`{
  "name": "Prompt trial",
  "ticketFile": "tickets/auth.md",
  "variants": [{"name": "baseline"}]
}`), 0644); err != nil {
		t.Fatal(err)
	}

	config, err := readExperimentConfig(configPath)
	if err != nil {
		t.Fatal(err)
	}

	if config.SourceWorkspace != agent1Path {
		t.Fatalf("SourceWorkspace = %q, want %q", config.SourceWorkspace, agent1Path)
	}
	if config.BaseBranch != "main" {
		t.Fatalf("BaseBranch = %q, want current branch main", config.BaseBranch)
	}
	if config.TicketFile != filepath.Join(root, "tickets", "auth.md") {
		t.Fatalf("TicketFile = %q, want resolved path", config.TicketFile)
	}
	if config.OutputDir != filepath.Join(root, "experiments") {
		t.Fatalf("OutputDir = %q, want default experiments dir", config.OutputDir)
	}
	if config.TimeoutMinutes != 90 {
		t.Fatalf("TimeoutMinutes = %d, want 90", config.TimeoutMinutes)
	}
	if config.PrepareTimeoutMinutes != 20 {
		t.Fatalf("PrepareTimeoutMinutes = %d, want 20", config.PrepareTimeoutMinutes)
	}
}

func TestBuildExperimentPromptIncludesVariantAndRules(t *testing.T) {
	oldAgentsFile := agentsFile
	oldDevAgentInstructionsFile := devAgentInstructionsFile
	oldTechFile := techFile
	dir := t.TempDir()
	agentsFile = filepath.Join(dir, "AGENTS.md")
	devAgentInstructionsFile = filepath.Join(dir, "DEV_AGENT.md")
	techFile = filepath.Join(dir, "TECH.md")
	t.Cleanup(func() {
		agentsFile = oldAgentsFile
		devAgentInstructionsFile = oldDevAgentInstructionsFile
		techFile = oldTechFile
	})

	mustWriteTestFile(t, agentsFile, "common rules")
	mustWriteTestFile(t, devAgentInstructionsFile, "dev rules")
	mustWriteTestFile(t, techFile, "tech rules")
	promptFile := filepath.Join(dir, "variant.md")
	mustWriteTestFile(t, promptFile, "use low reasoning")

	prompt, err := buildExperimentPrompt(ExperimentConfig{PromptMode: "bounded"}, Task{Body: "Build auth"}, ExperimentVariant{Name: "low", PromptFile: promptFile})
	if err != nil {
		t.Fatal(err)
	}

	for _, want := range []string{
		"common rules",
		"tech rules",
		"use low reasoning",
		"Build auth",
		"Exploration budget",
		"Do not push",
		"Do not install dependencies",
		"current working directory is already the assigned product repository",
		"Test contract, not implementation",
		"do not invent an interaction harness",
		"avoid adding more test code than implementation code",
	} {
		if !strings.Contains(prompt, want) {
			t.Fatalf("prompt missing %q\n%s", want, prompt)
		}
	}
	if strings.Contains(prompt, "dev rules") {
		t.Fatalf("prompt should not include full DEV_AGENT.md instructions\n%s", prompt)
	}
}

func TestBuildExperimentPromptCanMirrorDevOrchestratorPrompt(t *testing.T) {
	oldRepoRoot := repoRoot
	oldAgentsFile := agentsFile
	oldDevAgentInstructionsFile := devAgentInstructionsFile
	oldTechFile := techFile
	oldBacklogFile := backlogFile
	oldTasksFile := tasksFile
	oldArchiveFile := archiveFile
	dir := t.TempDir()
	repoRoot = dir
	agentsFile = filepath.Join(dir, "AGENTS.md")
	devAgentInstructionsFile = filepath.Join(dir, "DEV_AGENT.md")
	techFile = filepath.Join(dir, "TECH.md")
	backlogFile = filepath.Join(dir, "BACKLOG.md")
	tasksFile = filepath.Join(dir, "TASKS.md")
	archiveFile = filepath.Join(dir, "ARCHIVE.md")
	t.Cleanup(func() {
		repoRoot = oldRepoRoot
		agentsFile = oldAgentsFile
		devAgentInstructionsFile = oldDevAgentInstructionsFile
		techFile = oldTechFile
		backlogFile = oldBacklogFile
		tasksFile = oldTasksFile
		archiveFile = oldArchiveFile
	})

	mustWriteTestFile(t, agentsFile, "common rules")
	mustWriteTestFile(t, devAgentInstructionsFile, "dev rules")
	mustWriteTestFile(t, techFile, "tech rules")
	mustWriteTestFile(t, backlogFile, "# BACKLOG\n")
	mustWriteTestFile(t, tasksFile, "# TASKS\n\n## Dev Agent 2 In Progress\n\n### Other Task\nOwner: Dev Agent 2\nBranch: agent/2/other\nStatus: In Progress\n\nOther body\n")
	mustWriteTestFile(t, archiveFile, "# ARCHIVE\n")

	prompt, err := buildExperimentPrompt(
		ExperimentConfig{PromptMode: "orchestrator-dev", PromptRole: devAgent1Role},
		Task{
			Section: "Dev Agent 1 In Progress",
			Title:   "Build auth",
			Owner:   devAgent1Role,
			Branch:  "agent/1/auth",
			Status:  "In Progress",
			Body:    "Build auth body",
		},
		ExperimentVariant{Name: "orchestrator"},
	)
	if err != nil {
		t.Fatal(err)
	}

	for _, want := range []string{
		"Active role: Dev Agent 1",
		"dev rules",
		"Other active dev-agent work",
		"Other Task",
		"Build auth body",
		"EXPERIMENT SAFETY OVERRIDES",
		"Do not push",
		"current working directory is already the assigned product repository",
		"Test contract, not implementation",
		"do not invent an interaction harness",
		"avoid adding more test code than implementation code",
	} {
		if !strings.Contains(prompt, want) {
			t.Fatalf("prompt missing %q\n%s", want, prompt)
		}
	}
}

func TestExperimentPrepareCommandsSymlinkSourceNodeModules(t *testing.T) {
	source := t.TempDir()
	worktree := t.TempDir()
	if err := os.Mkdir(filepath.Join(source, "node_modules"), 0755); err != nil {
		t.Fatal(err)
	}
	mustWriteTestFile(t, filepath.Join(worktree, "package-lock.json"), "{}")

	commands := experimentPrepareCommands(ExperimentConfig{SourceWorkspace: source}, worktree)
	if len(commands) != 1 {
		t.Fatalf("commands = %v, want one symlink command", commands)
	}
	if !strings.Contains(commands[0], "ln -s") || !strings.Contains(commands[0], "node_modules") {
		t.Fatalf("command = %q, want node_modules symlink", commands[0])
	}
}

func TestExperimentPrepareCommandsFallsBackToNpmCi(t *testing.T) {
	worktree := t.TempDir()
	mustWriteTestFile(t, filepath.Join(worktree, "package-lock.json"), "{}")

	commands := experimentPrepareCommands(ExperimentConfig{SourceWorkspace: t.TempDir()}, worktree)
	if len(commands) != 1 || commands[0] != "npm ci" {
		t.Fatalf("commands = %v, want npm ci", commands)
	}
}

func TestExperimentPrepareCommandsCanBeSkipped(t *testing.T) {
	worktree := t.TempDir()
	mustWriteTestFile(t, filepath.Join(worktree, "package-lock.json"), "{}")

	commands := experimentPrepareCommands(ExperimentConfig{SkipPrepare: true}, worktree)
	if len(commands) != 0 {
		t.Fatalf("commands = %v, want none", commands)
	}
}

func TestResolveExperimentTaskFromBoardTask(t *testing.T) {
	oldBacklogFile := backlogFile
	dir := t.TempDir()
	backlogFile = filepath.Join(dir, "BACKLOG.md")
	t.Cleanup(func() {
		backlogFile = oldBacklogFile
	})
	mustWriteTestFile(t, backlogFile, `# BACKLOG

## Backlog

### Team Create Dialog

Task ID: UI-DIALOG-01
Status: Backlog

#### Objective

Move team creation into the shared dialog component.
`)

	task, source, err := resolveExperimentTask(ExperimentConfig{TaskTitle: "Team Create Dialog"})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(task.Body, "Move team creation") {
		t.Fatalf("task body = %q", task.Body)
	}
	if task.Title != "Team Create Dialog" {
		t.Fatalf("task title = %q", task.Title)
	}
	if !strings.Contains(source, "Team Create Dialog") {
		t.Fatalf("source = %q", source)
	}
}

func TestParseDiffShortstat(t *testing.T) {
	files, insertions, deletions := parseDiffShortstat(" 3 files changed, 24 insertions(+), 7 deletions(-)")
	if files != 3 || insertions != 24 || deletions != 7 {
		t.Fatalf("got %d/%d/%d", files, insertions, deletions)
	}
}

func TestParseCodexTokenUsageFromPlainTextSummary(t *testing.T) {
	path := filepath.Join(t.TempDir(), "codex.log")
	mustWriteTestFile(t, path, "tokens used\n360,377\n")

	usage, err := parseCodexTokenUsage(path)
	if err != nil {
		t.Fatal(err)
	}
	if usage.Total != 360377 {
		t.Fatalf("Total = %d, want 360377", usage.Total)
	}
}

func TestExperimentRunNameIsStableShape(t *testing.T) {
	name := experimentRunName("Auth Prompt Trial", time.Date(2026, 7, 2, 9, 8, 7, 0, time.UTC))
	if !strings.HasPrefix(name, "17-auth-prompt-trial-20260702-090807-") {
		t.Fatalf("unexpected run name %q", name)
	}
}

func mustWriteTestFile(t *testing.T, path string, contents string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(contents), 0644); err != nil {
		t.Fatal(err)
	}
}
