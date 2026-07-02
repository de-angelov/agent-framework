package main

import (
	"bufio"
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

type ExperimentConfig struct {
	Name                  string              `json:"name"`
	SourceWorkspace       string              `json:"sourceWorkspace"`
	BaseBranch            string              `json:"baseBranch"`
	TicketFile            string              `json:"ticketFile"`
	TaskSourceFile        string              `json:"taskSourceFile"`
	TaskTitle             string              `json:"taskTitle"`
	PromptMode            string              `json:"promptMode"`
	PromptRole            string              `json:"promptRole"`
	PromptBranch          string              `json:"promptBranch"`
	OutputDir             string              `json:"outputDir"`
	TimeoutMinutes        int                 `json:"timeoutMinutes"`
	PrepareTimeoutMinutes int                 `json:"prepareTimeoutMinutes"`
	PrepareCommands       []string            `json:"prepareCommands"`
	SkipPrepare           bool                `json:"skipPrepare"`
	Variants              []ExperimentVariant `json:"variants"`
}

type ExperimentVariant struct {
	Name       string            `json:"name"`
	Model      string            `json:"model"`
	Profile    string            `json:"profile"`
	PromptFile string            `json:"promptFile"`
	Config     map[string]string `json:"config"`
}

type ExperimentRun struct {
	Name       string                    `json:"name"`
	StartedAt  time.Time                 `json:"startedAt"`
	FinishedAt time.Time                 `json:"finishedAt"`
	BaseBranch string                    `json:"baseBranch"`
	BaseCommit string                    `json:"baseCommit"`
	TicketFile string                    `json:"ticketFile"`
	Results    []ExperimentVariantResult `json:"results"`
}

type ExperimentVariantResult struct {
	Name                 string            `json:"name"`
	Model                string            `json:"model,omitempty"`
	Profile              string            `json:"profile,omitempty"`
	Config               map[string]string `json:"config,omitempty"`
	Branch               string            `json:"branch"`
	Worktree             string            `json:"worktree"`
	LogFile              string            `json:"logFile"`
	PrepareLogFile       string            `json:"prepareLogFile,omitempty"`
	LastMessageFile      string            `json:"lastMessageFile"`
	PatchFile            string            `json:"patchFile"`
	Status               string            `json:"status"`
	PrepareStatus        string            `json:"prepareStatus,omitempty"`
	PrepareError         string            `json:"prepareError,omitempty"`
	PrepareCommands      []string          `json:"prepareCommands,omitempty"`
	PrepareMilliseconds  int64             `json:"prepareMilliseconds,omitempty"`
	ExitError            string            `json:"exitError,omitempty"`
	StartedAt            time.Time         `json:"startedAt"`
	FinishedAt           time.Time         `json:"finishedAt"`
	DurationMilliseconds int64             `json:"durationMilliseconds"`
	PromptBytes          int               `json:"promptBytes"`
	ApproxPromptTokens   int               `json:"approxPromptTokens"`
	DetectedInputTokens  int               `json:"detectedInputTokens,omitempty"`
	DetectedOutputTokens int               `json:"detectedOutputTokens,omitempty"`
	DetectedTotalTokens  int               `json:"detectedTotalTokens,omitempty"`
	BaseCommit           string            `json:"baseCommit"`
	HeadCommit           string            `json:"headCommit"`
	CommitCount          int               `json:"commitCount"`
	ChangedFiles         int               `json:"changedFiles"`
	Insertions           int               `json:"insertions"`
	Deletions            int               `json:"deletions"`
	DirtyAfterRun        bool              `json:"dirtyAfterRun"`
	AutoCommittedChanges bool              `json:"autoCommittedChanges"`
	FinalResponseSummary string            `json:"finalResponseSummary,omitempty"`
}

type tokenUsage struct {
	Input  int
	Output int
	Total  int
}

func runExperimentCLI(args []string) int {
	flags := flag.NewFlagSet("experiment", flag.ContinueOnError)
	flags.SetOutput(os.Stderr)
	configPath := flags.String("config", "", "path to experiment JSON config")
	dryRun := flags.Bool("dry-run", false, "prepare prompts and worktrees without running Codex")

	if err := flags.Parse(args); err != nil {
		return 2
	}
	if *configPath == "" {
		fmt.Fprintln(os.Stderr, "usage: cd orchestrator && go run . experiment --config ../experiments/example-agent-loop.json")
		return 2
	}

	config, err := readExperimentConfig(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "experiment config error: %v\n", err)
		return 1
	}

	run, err := runExperiment(context.Background(), config, *dryRun)
	if err != nil {
		fmt.Fprintf(os.Stderr, "experiment failed: %v\n", err)
		return 1
	}

	fmt.Printf("Experiment complete: %s\n", filepath.Join(config.resolvedOutputDir(), run.Name, "report.md"))
	return 0
}

func readExperimentConfig(path string) (ExperimentConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return ExperimentConfig{}, err
	}

	var config ExperimentConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return ExperimentConfig{}, err
	}
	if config.Name == "" {
		return ExperimentConfig{}, fmt.Errorf("name is required")
	}
	if config.SourceWorkspace == "" {
		config.SourceWorkspace = agent1Path
	}
	if !filepath.IsAbs(config.SourceWorkspace) {
		config.SourceWorkspace = filepath.Join(repoRoot, config.SourceWorkspace)
	}
	if config.BaseBranch == "" {
		config.BaseBranch = currentBranchName(config.SourceWorkspace)
		if config.BaseBranch == "" {
			config.BaseBranch = "main"
		}
	}
	if config.TicketFile != "" && !filepath.IsAbs(config.TicketFile) {
		config.TicketFile = filepath.Join(repoRoot, config.TicketFile)
	}
	if config.TaskSourceFile != "" && !filepath.IsAbs(config.TaskSourceFile) {
		config.TaskSourceFile = filepath.Join(repoRoot, config.TaskSourceFile)
	}
	if config.TicketFile == "" && config.TaskTitle == "" {
		return ExperimentConfig{}, fmt.Errorf("ticketFile or taskTitle is required")
	}
	if config.TicketFile != "" && config.TaskTitle != "" {
		return ExperimentConfig{}, fmt.Errorf("use either ticketFile or taskTitle, not both")
	}
	if config.PromptMode == "" {
		config.PromptMode = "bounded"
	}
	if config.PromptMode != "bounded" && config.PromptMode != "orchestrator-dev" {
		return ExperimentConfig{}, fmt.Errorf("unsupported promptMode %q", config.PromptMode)
	}
	if config.PromptRole == "" {
		config.PromptRole = devAgent1Role
	}
	if config.PromptMode == "orchestrator-dev" && config.PromptRole != devAgent1Role && config.PromptRole != devAgent2Role {
		return ExperimentConfig{}, fmt.Errorf("orchestrator-dev promptRole must be %q or %q", devAgent1Role, devAgent2Role)
	}
	if config.OutputDir == "" {
		config.OutputDir = filepath.Join(repoRoot, "experiments")
	} else if !filepath.IsAbs(config.OutputDir) {
		config.OutputDir = filepath.Join(repoRoot, config.OutputDir)
	}
	if config.TimeoutMinutes == 0 {
		config.TimeoutMinutes = 90
	}
	if config.PrepareTimeoutMinutes == 0 {
		config.PrepareTimeoutMinutes = 20
	}
	if len(config.Variants) == 0 {
		return ExperimentConfig{}, fmt.Errorf("at least one variant is required")
	}
	seen := map[string]bool{}
	for i := range config.Variants {
		if config.Variants[i].Name == "" {
			return ExperimentConfig{}, fmt.Errorf("variants[%d].name is required", i)
		}
		if seen[config.Variants[i].Name] {
			return ExperimentConfig{}, fmt.Errorf("duplicate variant name %q", config.Variants[i].Name)
		}
		seen[config.Variants[i].Name] = true
		if config.Variants[i].PromptFile != "" && !filepath.IsAbs(config.Variants[i].PromptFile) {
			config.Variants[i].PromptFile = filepath.Join(repoRoot, config.Variants[i].PromptFile)
		}
	}

	return config, nil
}

func (config ExperimentConfig) resolvedOutputDir() string {
	if config.OutputDir == "" {
		return filepath.Join(repoRoot, "experiments")
	}
	return config.OutputDir
}

func runExperiment(ctx context.Context, config ExperimentConfig, dryRun bool) (ExperimentRun, error) {
	runName := experimentRunName(config.Name, time.Now())
	runDir := filepath.Join(config.resolvedOutputDir(), runName)
	worktreeRoot := filepath.Join(runDir, "worktrees")
	if err := os.MkdirAll(worktreeRoot, 0755); err != nil {
		return ExperimentRun{}, err
	}

	baseCommit, err := gitOutput(config.SourceWorkspace, "rev-parse", config.BaseBranch)
	if err != nil {
		return ExperimentRun{}, fmt.Errorf("resolve base branch %s: %w", config.BaseBranch, err)
	}
	baseCommit = strings.TrimSpace(baseCommit)

	task, ticketSource, err := resolveExperimentTask(config)
	if err != nil {
		return ExperimentRun{}, err
	}

	run := ExperimentRun{
		Name:       runName,
		StartedAt:  time.Now(),
		BaseBranch: config.BaseBranch,
		BaseCommit: baseCommit,
		TicketFile: ticketSource,
	}

	for _, variant := range config.Variants {
		result, err := runExperimentVariant(ctx, config, runName, worktreeRoot, baseCommit, task, variant, dryRun)
		if err != nil {
			return run, err
		}
		run.Results = append(run.Results, result)
	}

	run.FinishedAt = time.Now()
	if err := writeExperimentReports(filepath.Join(config.resolvedOutputDir(), runName), run); err != nil {
		return run, err
	}
	return run, nil
}

func resolveExperimentTask(config ExperimentConfig) (Task, string, error) {
	if config.TicketFile != "" {
		ticket, err := os.ReadFile(config.TicketFile)
		if err != nil {
			return Task{}, "", err
		}
		return Task{
			Section: "Experiment",
			Title:   strings.TrimSuffix(filepath.Base(config.TicketFile), filepath.Ext(config.TicketFile)),
			Owner:   emptyAs(config.PromptRole, devAgent1Role),
			Branch:  config.PromptBranch,
			Status:  "In Progress",
			Body:    strings.TrimSpace(string(ticket)),
		}, config.TicketFile, nil
	}

	sourceFile := config.TaskSourceFile
	if sourceFile == "" {
		sourceFile = backlogFile
	}
	tasks, err := readTasks(sourceFile)
	if err != nil {
		return Task{}, "", err
	}
	for _, task := range tasks {
		if task.Title == config.TaskTitle {
			if config.PromptMode == "orchestrator-dev" {
				task.Section = config.PromptRole + " In Progress"
				task.Owner = config.PromptRole
				task.Status = "In Progress"
				if config.PromptBranch != "" {
					task.Branch = config.PromptBranch
				}
			}
			return task, fmt.Sprintf("%s:%s", sourceFile, config.TaskTitle), nil
		}
	}
	return Task{}, "", fmt.Errorf("task %q not found in %s", config.TaskTitle, sourceFile)
}

func runExperimentVariant(ctx context.Context, config ExperimentConfig, runName string, worktreeRoot string, baseCommit string, task Task, variant ExperimentVariant, dryRun bool) (ExperimentVariantResult, error) {
	start := time.Now()
	variantID := sanitizeBranchPart(variant.Name)
	branch := fmt.Sprintf("experiment/%s/%s", sanitizeBranchPart(runName), variantID)
	worktree := filepath.Join(worktreeRoot, variantID)
	logFile := filepath.Join(worktreeRoot, variantID+".jsonl")
	prepareLogFile := filepath.Join(worktreeRoot, variantID+"-prepare.log")
	lastMessageFile := filepath.Join(worktreeRoot, variantID+"-last-message.txt")
	patchFile := filepath.Join(worktreeRoot, variantID+".patch")

	if err := runGitChecked(config.SourceWorkspace, "worktree", "add", "-B", branch, worktree, baseCommit); err != nil {
		return ExperimentVariantResult{}, err
	}

	prompt, err := buildExperimentPrompt(config, task, variant)
	if err != nil {
		return ExperimentVariantResult{}, err
	}
	if err := os.WriteFile(filepath.Join(worktreeRoot, variantID+"-prompt.txt"), []byte(prompt), 0644); err != nil {
		return ExperimentVariantResult{}, err
	}

	result := ExperimentVariantResult{
		Name:               variant.Name,
		Model:              variant.Model,
		Profile:            variant.Profile,
		Config:             variant.Config,
		Branch:             branch,
		Worktree:           worktree,
		LogFile:            logFile,
		PrepareLogFile:     prepareLogFile,
		LastMessageFile:    lastMessageFile,
		PatchFile:          patchFile,
		Status:             "prepared",
		StartedAt:          start,
		PromptBytes:        len(prompt),
		ApproxPromptTokens: approximateTokens(prompt),
		BaseCommit:         baseCommit,
	}

	if dryRun {
		result.Status = "dry-run"
		result.FinishedAt = time.Now()
		result.DurationMilliseconds = result.FinishedAt.Sub(result.StartedAt).Milliseconds()
		return finalizeExperimentResult(worktree, baseCommit, patchFile, result)
	}

	result = prepareExperimentWorktree(ctx, config, worktree, prepareLogFile, result)
	if result.PrepareStatus == "failed" {
		result.Status = "prepare-failed"
		result.FinishedAt = time.Now()
		result.DurationMilliseconds = result.FinishedAt.Sub(result.StartedAt).Milliseconds()
		return finalizeExperimentResult(worktree, baseCommit, patchFile, result)
	}

	timeout := time.Duration(config.TimeoutMinutes) * time.Minute
	runCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	err = runCodexExperiment(runCtx, worktree, prompt, variant, logFile, lastMessageFile)
	result.FinishedAt = time.Now()
	result.DurationMilliseconds = result.FinishedAt.Sub(result.StartedAt).Milliseconds()
	if err != nil {
		result.Status = "failed"
		result.ExitError = err.Error()
	} else {
		result.Status = "completed"
	}
	if runCtx.Err() == context.DeadlineExceeded {
		result.Status = "timeout"
		result.ExitError = runCtx.Err().Error()
	}

	if usage, err := parseCodexTokenUsage(logFile); err == nil {
		result.DetectedInputTokens = usage.Input
		result.DetectedOutputTokens = usage.Output
		result.DetectedTotalTokens = usage.Total
	}
	if data, err := os.ReadFile(lastMessageFile); err == nil {
		result.FinalResponseSummary = oneLine(string(data))
	}

	return finalizeExperimentResult(worktree, baseCommit, patchFile, result)
}

func prepareExperimentWorktree(ctx context.Context, config ExperimentConfig, worktree string, logFile string, result ExperimentVariantResult) ExperimentVariantResult {
	start := time.Now()
	commands := experimentPrepareCommands(config, worktree)
	result.PrepareCommands = commands

	if len(commands) == 0 {
		result.PrepareStatus = "skipped"
		return result
	}

	timeout := time.Duration(config.PrepareTimeoutMinutes) * time.Minute
	prepareCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	logOutput, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		result.PrepareStatus = "failed"
		result.PrepareError = err.Error()
		return result
	}
	defer logOutput.Close()

	for _, command := range commands {
		if _, err := fmt.Fprintf(logOutput, "$ %s\n", command); err != nil {
			result.PrepareStatus = "failed"
			result.PrepareError = err.Error()
			return result
		}

		cmd := exec.CommandContext(prepareCtx, "bash", "-lc", command)
		cmd.Dir = worktree
		cmd.Stdout = logOutput
		cmd.Stderr = logOutput
		if err := cmd.Run(); err != nil {
			result.PrepareStatus = "failed"
			result.PrepareError = err.Error()
			if prepareCtx.Err() == context.DeadlineExceeded {
				result.PrepareError = prepareCtx.Err().Error()
			}
			result.PrepareMilliseconds = time.Since(start).Milliseconds()
			return result
		}
	}

	result.PrepareStatus = "completed"
	result.PrepareMilliseconds = time.Since(start).Milliseconds()
	return result
}

func experimentPrepareCommands(config ExperimentConfig, worktree string) []string {
	if config.SkipPrepare {
		return nil
	}
	if len(config.PrepareCommands) > 0 {
		return config.PrepareCommands
	}
	if _, err := os.Stat(filepath.Join(worktree, "package-lock.json")); err != nil {
		return nil
	}
	if _, err := os.Stat(filepath.Join(worktree, "node_modules")); err == nil {
		return nil
	}

	sourceNodeModules := filepath.Join(config.SourceWorkspace, "node_modules")
	if _, err := os.Stat(sourceNodeModules); err == nil {
		return []string{fmt.Sprintf("ln -s %s node_modules", shellQuote(sourceNodeModules))}
	}

	return []string{"npm ci"}
}

func buildExperimentPrompt(config ExperimentConfig, task Task, variant ExperimentVariant) (string, error) {
	variantPrompt := ""
	if variant.PromptFile != "" {
		data, err := os.ReadFile(variant.PromptFile)
		if err != nil {
			return "", err
		}
		variantPrompt = string(data)
	}

	if config.PromptMode == "orchestrator-dev" {
		tasks, err := readBoardTasks()
		if err != nil {
			return "", err
		}
		prompt := buildPrompt(config.PromptRole, task, tasks, `
Role: Dev Agent
Runtime Rules:
- Work only on the assigned task and keep changes focused.
`)
		return prompt + fmt.Sprintf(`

================ EXPERIMENT SAFETY OVERRIDES ================

- This is an isolated experiment branch, not a live orchestrator session.
- The current working directory is already the assigned product repository. Do not cd into "workspaces/repo-agent-*" unless that directory exists from the current working directory.
- Do not push, merge, update BACKLOG.md, TASKS.md, or ARCHIVE.md.
- Do not install dependencies or change package manager files unless the ticket explicitly requires it.
- Dependencies are prepared before you start.
- If verification reveals unrelated failures, stop and report them. Do not fix unrelated files.
- Do not commit; the experiment harness records and commits the final diff after the run.
- Testing discipline: prefer existing test style and user-visible/rendered assertions over inspecting React component internals.
- Do not add custom test traversal helpers, mock frameworks, DOM harnesses, or new test utilities unless the ticket explicitly requires interaction behavior that cannot be covered otherwise.
- Test contract, not implementation: assert rendered text, form fields, button labels, form ids/actions, validation messages, and preserved existing behavior.
- Avoid asserting component prop wiring directly when the same behavior can be observed in rendered output.
- If the repo lacks jsdom/testing-library, do not invent an interaction harness. Use the existing SSR/render test pattern and state the interaction limit if relevant.
- Scope budget: for small route/UI tickets, keep test changes close to the changed file and avoid adding more test code than implementation code unless a failing test proves it is necessary.

================ VARIANT PROMPT ================

%s
`, variantPrompt), nil
	}

	return fmt.Sprintf(`You are running an isolated implementation experiment.

================ AGENTS.md COMMON RULES ================

%s

================ TECH.md ================

%s

================ BOUNDED EXPERIMENT RULES ================

- Implement the ticket in this isolated worktree only.
- The current working directory is already the assigned product repository. Do not cd into "workspaces/repo-agent-*" unless that directory exists from the current working directory.
- Do not push, merge, update BACKLOG.md, TASKS.md, or ARCHIVE.md.
- Do not install dependencies or change package manager files unless the ticket explicitly requires it; dependencies are prepared before you start.
- Use the existing repository test stack. If a missing tool or unrelated type error blocks verification, report it instead of expanding scope.
- Exploration budget: inspect at most 8 relevant files before editing unless you have a concrete failing command that requires more context.
- Implementation budget: make one focused implementation pass for this ticket.
- Verification budget: run only the ticket's listed verification command or the narrowest equivalent focused command.
- Testing discipline: prefer existing test style and user-visible/rendered assertions over inspecting React component internals.
- Do not add custom test traversal helpers, mock frameworks, DOM harnesses, or new test utilities unless the ticket explicitly requires interaction behavior that cannot be covered otherwise.
- Test contract, not implementation: assert rendered text, form fields, button labels, form ids/actions, validation messages, and preserved existing behavior.
- Avoid asserting component prop wiring directly when the same behavior can be observed in rendered output.
- If the repo lacks jsdom/testing-library, do not invent an interaction harness. Use the existing SSR/render test pattern and state the interaction limit if relevant.
- Scope budget: for small route/UI tickets, keep test changes close to the changed file and avoid adding more test code than implementation code unless a failing test proves it is necessary.
- If verification reveals unrelated failures, stop and report them. Do not fix unrelated files.
- Do not commit; the experiment harness records and commits the final diff after the run.
- Keep the implementation focused so this branch can be compared against other variants.

================ VARIANT PROMPT ================

%s

================ TICKET ================

%s
`, mustRead(agentsFile), mustRead(techFile), variantPrompt, task.Body), nil
}

func runCodexExperiment(ctx context.Context, workspace string, prompt string, variant ExperimentVariant, logFile string, lastMessageFile string) error {
	args := []string{"exec", "--sandbox", "danger-full-access", "--json", "--output-last-message", lastMessageFile}
	if variant.Model != "" {
		args = append(args, "--model", variant.Model)
	}
	if variant.Profile != "" {
		args = append(args, "--profile", variant.Profile)
	}
	for _, key := range sortedConfigKeys(variant.Config) {
		args = append(args, "--config", fmt.Sprintf("%s=%s", key, variant.Config[key]))
	}
	args = append(args, "-")

	cmd := exec.CommandContext(ctx, "codex", args...)
	cmd.Dir = workspace
	cmd.Stdin = strings.NewReader(prompt)

	logOutput, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer logOutput.Close()

	var stderr bytes.Buffer
	cmd.Stdout = logOutput
	cmd.Stderr = io.MultiWriter(logOutput, &stderr)

	if err := cmd.Run(); err != nil {
		if stderr.Len() > 0 {
			return fmt.Errorf("%w: %s", err, oneLine(stderr.String()))
		}
		return err
	}
	return nil
}

func finalizeExperimentResult(worktree string, baseCommit string, patchFile string, result ExperimentVariantResult) (ExperimentVariantResult, error) {
	if workspaceHasChanges(worktree) {
		result.DirtyAfterRun = true
		if err := stageExperimentChanges(worktree); err != nil {
			return result, err
		}
		if workspaceHasStagedChanges(worktree) {
			if err := runGitChecked(worktree, "commit", "-m", "Experiment result: "+result.Name); err != nil {
				return result, err
			}
			result.AutoCommittedChanges = true
			result.DirtyAfterRun = false
		}
	}

	head, err := gitOutput(worktree, "rev-parse", "HEAD")
	if err != nil {
		return result, err
	}
	result.HeadCommit = strings.TrimSpace(head)

	count, err := gitOutput(worktree, "rev-list", "--count", baseCommit+"..HEAD")
	if err == nil {
		result.CommitCount, _ = strconv.Atoi(strings.TrimSpace(count))
	}

	stat, _ := gitOutput(worktree, "diff", "--shortstat", baseCommit+"..HEAD")
	result.ChangedFiles, result.Insertions, result.Deletions = parseDiffShortstat(stat)

	patch, err := gitOutput(worktree, "diff", "--patch", baseCommit+"..HEAD")
	if err != nil {
		return result, err
	}
	if err := os.WriteFile(patchFile, []byte(patch), 0644); err != nil {
		return result, err
	}

	result.DirtyAfterRun = workspaceHasChanges(worktree)
	return result, nil
}

func stageExperimentChanges(worktree string) error {
	if err := runGitChecked(worktree, "add", "-u", "--", "."); err != nil {
		return err
	}

	output, err := gitOutput(worktree, "ls-files", "--others", "--exclude-standard", "-z")
	if err != nil {
		return err
	}

	var paths []string
	for _, path := range strings.Split(output, "\x00") {
		if path == "" || path == "node_modules" || strings.HasPrefix(path, "node_modules/") {
			continue
		}
		paths = append(paths, path)
	}
	if len(paths) == 0 {
		return nil
	}

	args := append([]string{"add", "--"}, paths...)
	return runGitChecked(worktree, args...)
}

func writeExperimentReports(runDir string, run ExperimentRun) error {
	data, err := json.MarshalIndent(run, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(runDir, "report.json"), append(data, '\n'), 0644); err != nil {
		return err
	}

	var b strings.Builder
	fmt.Fprintf(&b, "# Experiment: %s\n\n", run.Name)
	fmt.Fprintf(&b, "- Base branch: `%s`\n", run.BaseBranch)
	fmt.Fprintf(&b, "- Base commit: `%s`\n", run.BaseCommit)
	fmt.Fprintf(&b, "- Ticket: `%s`\n", run.TicketFile)
	fmt.Fprintf(&b, "- Started: `%s`\n", run.StartedAt.Format(time.RFC3339))
	fmt.Fprintf(&b, "- Finished: `%s`\n\n", run.FinishedAt.Format(time.RFC3339))

	b.WriteString("| Variant | Status | Prepare | Duration | Approx prompt tokens | Detected total tokens | Commits | Files | +/- | Branch |\n")
	b.WriteString("| --- | --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | --- |\n")
	for _, result := range run.Results {
		fmt.Fprintf(&b, "| %s | %s | %s | %s | %d | %s | %d | %d | +%d/-%d | `%s` |\n",
			escapeTable(result.Name),
			escapeTable(result.Status),
			escapeTable(formatPrepareStatus(result)),
			(time.Duration(result.DurationMilliseconds) * time.Millisecond).Round(time.Second),
			result.ApproxPromptTokens,
			emptyInt(result.DetectedTotalTokens),
			result.CommitCount,
			result.ChangedFiles,
			result.Insertions,
			result.Deletions,
			result.Branch,
		)
	}
	b.WriteString("\n## Artifacts\n\n")
	for _, result := range run.Results {
		fmt.Fprintf(&b, "### %s\n\n", result.Name)
		fmt.Fprintf(&b, "- Worktree: `%s`\n", result.Worktree)
		if result.PrepareLogFile != "" {
			fmt.Fprintf(&b, "- Prepare log: `%s`\n", result.PrepareLogFile)
		}
		fmt.Fprintf(&b, "- Log: `%s`\n", result.LogFile)
		fmt.Fprintf(&b, "- Patch: `%s`\n", result.PatchFile)
		fmt.Fprintf(&b, "- Head: `%s`\n", result.HeadCommit)
		if len(result.PrepareCommands) > 0 {
			fmt.Fprintf(&b, "- Prepare commands: `%s`\n", strings.Join(result.PrepareCommands, " && "))
		}
		if result.PrepareError != "" {
			fmt.Fprintf(&b, "- Prepare error: `%s`\n", result.PrepareError)
		}
		if result.ExitError != "" {
			fmt.Fprintf(&b, "- Error: `%s`\n", result.ExitError)
		}
		if result.FinalResponseSummary != "" {
			fmt.Fprintf(&b, "- Final response: %s\n", result.FinalResponseSummary)
		}
		b.WriteString("\n")
	}

	return os.WriteFile(filepath.Join(runDir, "report.md"), []byte(b.String()), 0644)
}

func formatPrepareStatus(result ExperimentVariantResult) string {
	if result.PrepareStatus == "" {
		return ""
	}
	if result.PrepareMilliseconds == 0 {
		return result.PrepareStatus
	}
	return fmt.Sprintf("%s %s", result.PrepareStatus, (time.Duration(result.PrepareMilliseconds) * time.Millisecond).Round(time.Second))
}

func parseCodexTokenUsage(path string) (tokenUsage, error) {
	f, err := os.Open(path)
	if err != nil {
		return tokenUsage{}, err
	}
	defer f.Close()

	var usage tokenUsage
	nextLineIsTotal := false
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if nextLineIsTotal {
			if value := parseFlexibleInt(line); value > 0 {
				usage.Total += value
			}
			nextLineIsTotal = false
		}
		if strings.EqualFold(strings.TrimSpace(line), "tokens used") {
			nextLineIsTotal = true
			continue
		}
		for _, pair := range tokenRegex.FindAllStringSubmatch(line, -1) {
			value := parseFlexibleInt(pair[2])
			key := strings.ToLower(pair[1])
			switch {
			case strings.Contains(key, "input"):
				usage.Input += value
			case strings.Contains(key, "output"):
				usage.Output += value
			case strings.Contains(key, "total"):
				usage.Total += value
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return tokenUsage{}, err
	}
	if usage.Total == 0 {
		usage.Total = usage.Input + usage.Output
	}
	return usage, nil
}

var tokenRegex = regexp.MustCompile(`(?i)"?([a-z_]*tokens?)"?\s*[:=]\s*([0-9,]+)`)

func parseFlexibleInt(value string) int {
	value = strings.ReplaceAll(strings.TrimSpace(value), ",", "")
	parsed, _ := strconv.Atoi(value)
	return parsed
}

func parseDiffShortstat(stat string) (int, int, int) {
	files := extractShortstatNumber(stat, `([0-9]+) files? changed`)
	insertions := extractShortstatNumber(stat, `([0-9]+) insertions?\(\+\)`)
	deletions := extractShortstatNumber(stat, `([0-9]+) deletions?\(-\)`)
	return files, insertions, deletions
}

func extractShortstatNumber(stat string, pattern string) int {
	re := regexp.MustCompile(pattern)
	match := re.FindStringSubmatch(stat)
	if len(match) < 2 {
		return 0
	}
	value, _ := strconv.Atoi(match[1])
	return value
}

func approximateTokens(text string) int {
	if text == "" {
		return 0
	}
	return (len([]rune(text)) + 3) / 4
}

func experimentRunName(name string, now time.Time) string {
	return fmt.Sprintf("%s-%s-%s", sanitizeBranchPart(name), now.Format("20060102-150405"), shortHash(name+now.String()))
}

func shortHash(value string) string {
	sum := sha1.Sum([]byte(value))
	return hex.EncodeToString(sum[:])[:8]
}

func sortedConfigKeys(config map[string]string) []string {
	keys := make([]string, 0, len(config))
	for key := range config {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func escapeTable(value string) string {
	return strings.ReplaceAll(value, "|", "\\|")
}

func emptyInt(value int) string {
	if value == 0 {
		return ""
	}
	return strconv.Itoa(value)
}

func shellQuote(value string) string {
	return "'" + strings.ReplaceAll(value, "'", "'\\''") + "'"
}
