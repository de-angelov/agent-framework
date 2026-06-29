package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func openLogOutput() (*os.File, error) {
	return os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
}

type lockedLogWriter struct{}

func (w lockedLogWriter) Write(p []byte) (n int, err error) {
	logMu.Lock()
	defer logMu.Unlock()

	f, err := openLogOutput()
	if err != nil {
		return 0, err
	}
	defer f.Close()

	return f.Write(p)
}

func logEvent(format string, args ...any) {
	logMu.Lock()
	defer logMu.Unlock()

	line := fmt.Sprintf("%s %s\n", time.Now().Format(time.RFC3339), fmt.Sprintf(format, args...))

	f, err := openLogOutput()
	if err != nil {
		return
	}
	defer f.Close()

	_, _ = f.WriteString(line)
}

func renderUI(tasks []Task, readErr error) {
	mu.Lock()
	rows := buildRows(tasks, running, finished)
	resourceStatus := codexResourceStatusLine()
	mu.Unlock()

	latest := latestLogLine()

	var b strings.Builder
	b.WriteString("\x1b[2J\x1b[H")
	b.WriteString("Orchestrator\n")
	b.WriteString("Board: " + filepath.Base(repoRoot) + "\n")
	b.WriteString("Time: " + time.Now().Format(time.RFC3339) + "\n")
	if resourceStatus != "" {
		b.WriteString("Codex: " + resourceStatus + "\n")
	}

	if readErr != nil {
		b.WriteString("Board error: " + readErr.Error() + "\n")
	}

	b.WriteString("\n")
	b.WriteString("ROLE             STATUS           TASK                          BRANCH\n")
	b.WriteString("--------------   ----------------  -----------------------------  ------------------------------\n")

	for _, row := range rows {
		b.WriteString(fmt.Sprintf("%s %-14s %-16s %-29s %s\n",
			row.marker,
			row.role,
			colorStatus(row.status),
			truncate(row.task, 29),
			truncate(row.branch, 30)))
	}

	b.WriteString("\n")
	b.WriteString("Latest: " + latest + "\n")
	b.WriteString("Logs: " + logFilePath + "\n")

	fmt.Print(b.String())
}

func buildRows(tasks []Task, sessions map[string]RunningSession, finishedSessions map[string]FinishedSession) []struct {
	marker string
	role   string
	status string
	task   string
	branch string
} {
	rows := []struct {
		marker string
		role   string
		status string
		task   string
		branch string
	}{
		{role: teamLeadRole},
		{role: devAgent1Role},
		{role: devAgent2Role},
	}

	for i := range rows {
		if session, ok := sessions[rows[i].role]; ok {
			rows[i].marker = ">"
			rows[i].status = "running"
			rows[i].task = session.Task
			rows[i].branch = session.Branch
			continue
		}

		if rows[i].role == devAgent1Role || rows[i].role == devAgent2Role {
			activeTasks := activeTasksForRole(tasks, rows[i].role)
			if len(activeTasks) > 1 {
				rows[i].status = "board error"
				rows[i].task = fmt.Sprintf("%d active tasks", len(activeTasks))
				continue
			}
		}

		task := findDesiredTaskForRole(tasks, rows[i].role)
		if task.Title != "" {
			if finishedSession, ok := finishedSessions[rows[i].role]; ok && finishedSession.TaskKey == task.Key {
				rows[i].status = string(finishedSession.Outcome)
			} else {
				rows[i].status = task.Status
			}
			rows[i].task = task.Title
			rows[i].branch = task.Branch
			continue
		}

		if rows[i].role == teamLeadRole && hasBoardError(tasks) {
			rows[i].status = "board error"
		} else if rows[i].role == teamLeadRole && hasBacklog(tasks) {
			if lanesHaveCapacity(tasks) {
				rows[i].status = "backlog pending"
			} else {
				rows[i].status = "waiting for lane"
			}
		} else {
			rows[i].status = "idle"
		}
	}

	return rows
}

func findDesiredTaskForRole(tasks []Task, role string) Task {
	switch role {
	case devAgent1Role, devAgent2Role:
		activeTasks := activeTasksForRole(tasks, role)
		if len(activeTasks) == 1 {
			return activeTasks[0]
		}
	case teamLeadRole:
		if hasBoardError(tasks) || !lanesHaveCapacity(tasks) {
			return Task{}
		}
		return firstBacklogTask(tasks)
	}
	return Task{}
}

func activeTasksForRole(tasks []Task, role string) []Task {
	var active []Task
	for _, task := range tasks {
		switch role {
		case devAgent1Role:
			if task.Section == "Dev Agent 1 In Progress" && task.Owner == devAgent1Role && task.Status == "In Progress" {
				active = append(active, task)
			}
		case devAgent2Role:
			if task.Section == "Dev Agent 2 In Progress" && task.Owner == devAgent2Role && task.Status == "In Progress" {
				active = append(active, task)
			}
		}
	}
	return active
}

func hasBoardError(tasks []Task) bool {
	return len(activeTasksForRole(tasks, devAgent1Role)) > 1 || len(activeTasksForRole(tasks, devAgent2Role)) > 1
}

func lanesHaveCapacity(tasks []Task) bool {
	return len(activeTasksForRole(tasks, devAgent1Role)) == 0 || len(activeTasksForRole(tasks, devAgent2Role)) == 0
}

func hasBacklog(tasks []Task) bool {
	return firstBacklogTask(tasks).Title != ""
}

func firstBacklogTask(tasks []Task) Task {
	for _, task := range tasks {
		if task.Section == "Backlog" && (task.Status == "Backlog" || task.Status == "") {
			return task
		}
	}
	return Task{}
}

func truncate(s string, limit int) string {
	if len(s) <= limit {
		return s
	}
	if limit <= 3 {
		return s[:limit]
	}
	return s[:limit-3] + "..."
}

func colorStatus(status string) string {
	switch status {
	case "running":
		return "\x1b[32mrunning\x1b[0m"
	case "idle":
		return "\x1b[90midle\x1b[0m"
	case "backlog pending":
		return "\x1b[33mbacklog pending\x1b[0m"
	case "waiting for lane":
		return "\x1b[90mwaiting for lane\x1b[0m"
	case "board error":
		return "\x1b[31mboard error\x1b[0m"
	case "In Progress":
		return "\x1b[36mIn Progress\x1b[0m"
	case "completed":
		return "\x1b[32mcompleted\x1b[0m"
	case "failed":
		return "\x1b[31mfailed\x1b[0m"
	case "rate limited":
		return "\x1b[33mrate limited\x1b[0m"
	case "Backlog":
		return "\x1b[90mBacklog\x1b[0m"
	default:
		return status
	}
}

func codexResourceStatusLine() string {
	now := time.Now()
	if !codexResourcePausedUntil.IsZero() && now.Before(codexResourcePausedUntil) {
		return "paused until " + codexResourcePausedUntil.Format(time.RFC3339) + " (" + lastCodexStatusMessage + ")"
	}
	if lastCodexStatusMessage != "" {
		return lastCodexStatusMessage
	}
	return ""
}

func latestLogLine() string {
	f, err := os.Open(logFilePath)
	if err != nil {
		return "[no log]"
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil || info.Size() == 0 {
		return "[no log]"
	}

	const maxRead int64 = 4096

	size := info.Size()
	start := size - maxRead
	if start < 0 {
		start = 0
	}

	buf := make([]byte, size-start)
	_, err = f.ReadAt(buf, start)
	if err != nil && err != io.EOF {
		return "[no log]"
	}

	text := strings.TrimSpace(string(buf))
	if text == "" {
		return "[no log]"
	}

	lines := strings.Split(text, "\n")
	last := ansiColorRe.ReplaceAllString(lines[len(lines)-1], "")

	return last
}

func emptyAs(value string, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}
