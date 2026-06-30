package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestResolveRepoRootFromCurrentRoot(t *testing.T) {
	root := makeRepoRoot(t)

	resolved, ok := resolveRepoRootFrom(root)
	if !ok {
		t.Fatal("expected repo root to be resolved")
	}
	if resolved != root {
		t.Fatalf("resolved root = %q, want %q", resolved, root)
	}
}

func TestResolveRepoRootFromOrchestratorSubdirectory(t *testing.T) {
	root := makeRepoRoot(t)
	orchestratorDir := filepath.Join(root, "orchestrator")
	if err := os.Mkdir(orchestratorDir, 0755); err != nil {
		t.Fatal(err)
	}

	resolved, ok := resolveRepoRootFrom(orchestratorDir)
	if !ok {
		t.Fatal("expected repo root to be resolved")
	}
	if resolved != root {
		t.Fatalf("resolved root = %q, want %q", resolved, root)
	}
}

func TestResolveRepoRootFromDirectoryWithoutMarkers(t *testing.T) {
	dir := t.TempDir()

	resolved, ok := resolveRepoRootFrom(dir)
	if ok {
		t.Fatalf("resolved root = %q, want no root", resolved)
	}
}

func TestNewRunLogFilePathUsesTimestampedLogFile(t *testing.T) {
	oldLogsRoot := logsRoot
	logsRoot = t.TempDir()

	t.Cleanup(func() {
		logsRoot = oldLogsRoot
	})

	got := newRunLogFilePath(time.Date(2026, 6, 30, 13, 45, 12, 123456789, time.UTC))
	want := filepath.Join(logsRoot, "orchestrator-20260630-134512.123456789.log")

	if got != want {
		t.Fatalf("newRunLogFilePath() = %q, want %q", got, want)
	}
}

func makeRepoRoot(t *testing.T) string {
	t.Helper()

	root := t.TempDir()
	for _, marker := range repoRootMarkers {
		path := filepath.Join(root, marker)
		if err := os.WriteFile(path, []byte(marker+"\n"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	return root
}
