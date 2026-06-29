package main

import (
	"os"
	"path/filepath"
	"testing"
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
