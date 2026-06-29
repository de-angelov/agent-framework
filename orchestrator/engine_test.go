package main

import (
	"path/filepath"
	"testing"
	"time"
)

func TestShouldRetryFailedSessionAfterDelay(t *testing.T) {
	withRetryTestState(t)
	now := time.Now()
	session := FinishedSession{
		Role:       devAgent1Role,
		TaskKey:    "task-1",
		Outcome:    sessionFailed,
		FinishedAt: now.Add(-failedSessionRetryDelay),
	}

	mu.Lock()
	defer mu.Unlock()

	shouldRetry, retryCount := shouldRetryFailedSessionLocked(devAgent1Role, session, now)
	if !shouldRetry {
		t.Fatal("expected failed session to be retried after cooldown")
	}
	if retryCount != 1 {
		t.Fatalf("retry count return = %d, want 1", retryCount)
	}

	key := failedSessionRetryKey(devAgent1Role, session.TaskKey)
	if failedSessionRetryCounts[key] != 1 {
		t.Fatalf("retry count = %d, want 1", failedSessionRetryCounts[key])
	}
}

func TestShouldNotRetryFailedSessionBeforeDelay(t *testing.T) {
	withRetryTestState(t)
	now := time.Now()
	session := FinishedSession{
		Role:       devAgent1Role,
		TaskKey:    "task-1",
		Outcome:    sessionFailed,
		FinishedAt: now.Add(-failedSessionRetryDelay + time.Second),
	}

	mu.Lock()
	defer mu.Unlock()

	shouldRetry, retryCount := shouldRetryFailedSessionLocked(devAgent1Role, session, now)
	if shouldRetry {
		t.Fatal("expected failed session to wait for cooldown")
	}
	if retryCount != 0 {
		t.Fatalf("retry count return = %d, want 0", retryCount)
	}
}

func TestShouldStopRetryingFailedSessionAtLimit(t *testing.T) {
	withRetryTestState(t)
	now := time.Now()
	session := FinishedSession{
		Role:       devAgent1Role,
		TaskKey:    "task-1",
		Outcome:    sessionFailed,
		FinishedAt: now.Add(-failedSessionRetryDelay),
	}
	key := failedSessionRetryKey(devAgent1Role, session.TaskKey)
	failedSessionRetryCounts[key] = maxFailedSessionRetries

	mu.Lock()
	defer mu.Unlock()

	shouldRetry, retryCount := shouldRetryFailedSessionLocked(devAgent1Role, session, now)
	if shouldRetry {
		t.Fatal("expected failed session retry limit to be enforced")
	}
	if retryCount != maxFailedSessionRetries {
		t.Fatalf("retry count return = %d, want %d", retryCount, maxFailedSessionRetries)
	}
	if failedSessionRetryCounts[key] != maxFailedSessionRetries {
		t.Fatalf("retry count = %d, want %d", failedSessionRetryCounts[key], maxFailedSessionRetries)
	}
}

func TestShouldNotRetryCompletedSession(t *testing.T) {
	withRetryTestState(t)
	now := time.Now()
	session := FinishedSession{
		Role:       devAgent1Role,
		TaskKey:    "task-1",
		Outcome:    sessionCompleted,
		FinishedAt: now.Add(-failedSessionRetryDelay),
	}

	mu.Lock()
	defer mu.Unlock()

	shouldRetry, retryCount := shouldRetryFailedSessionLocked(devAgent1Role, session, now)
	if shouldRetry {
		t.Fatal("expected completed session not to be retried")
	}
	if retryCount != 0 {
		t.Fatalf("retry count return = %d, want 0", retryCount)
	}
}

func withRetryTestState(t *testing.T) {
	t.Helper()

	oldLogFilePath := logFilePath
	logFilePath = filepath.Join(t.TempDir(), "orchestrator.log")
	failedSessionRetryCounts = map[string]int{}

	t.Cleanup(func() {
		logFilePath = oldLogFilePath
		failedSessionRetryCounts = map[string]int{}
	})
}
