# Multi-Agent Orchestrator Refactoring

## Goal

Refactor the Go multi-agent orchestrator to minimize LLM token usage while preserving agent quality and autonomy.

The current architecture launches long-running `codex exec` sessions. Investigation of Codex session logs showed that token spikes are primarily caused by accumulated conversation history (hundreds of turns and repeated context compaction), not by build/test output.

The new architecture should resemble Cursor/Aider style operation:

* short-lived AI sessions
* external persistent state
* retrieval of only relevant context
* deterministic prompts
* minimal history

---

## High-Level Design

### 1. Stateless Agent Sessions

Every `codex exec` invocation should be treated as an isolated transaction.

A session should:

* read current task
* read persistent state files
* perform exactly one milestone of work
* update persistent state
* exit

Never rely on Codex conversation history for memory.

---

### 2. Persistent Agent State

Create a `.agent/` directory inside every workspace.

Example:

.agent/
session.json
repo-state.json
decisions.md
summary.md
touched-files.txt
rich-doc-summary.md

Purpose:

session.json

* current phase
* current objective
* next milestone
* blockers

repo-state.json

* changed files
* tests status
* build status
* last successful commit
* last validation timestamp

decisions.md

Permanent architectural decisions.

summary.md

Compact running summary of completed work.

touched-files.txt

Only files modified during the current task.

rich-doc-summary.md

Extracted information from Word/PDF/images.

Never store screenshots, Base64, or full documents.

---

### 3. Session Lifecycle

Implementation should become:

Start Session

↓

Read:

* AGENTS.md
* TECH.md
* role instructions
* current task
* .agent/* state

↓

Perform ONE milestone

↓

Update .agent/*

↓

Exit Codex

Examples of milestones:

* implement feature
* validation pass
* fix failing tests
* review
* merge
* backlog grooming

---

### 4. Team Lead Rich Document Pipeline

If the Team Lead receives:

* Word
* PDF
* Screenshots
* Images

It must

1. extract information

2. generate compact Markdown

3. write

.agent/rich-doc-summary.md

4. update backlog/tasks

5. terminate session

Developer agents must never receive original documents or Base64 payloads.

---

### 5. Prompt Construction

Every new Codex session should read only:

AGENTS.md

TECH.md

Role instructions

Current task

.agent/session.json

.agent/summary.md

.agent/decisions.md

.agent/repo-state.json

Never include previous Codex chat history.

---

### 6. Context Limits

Enforce limits.

Examples:

Task body:
max 12 KB

Board summaries:
2 lines per task

Persistent summary:
max 8 KB

decisions.md:
max 4 KB

Never include

* build logs
* test logs
* compiler output
* Base64
* screenshots
* full diffs
* full source files

---

### 7. Runtime Rules

Update prompts.

Developer agents:

* Complete only one milestone.
* Prefer concise summaries.
* Do not inspect unrelated files.
* Do not continue after milestone completion.
* Update persistent state before exiting.

Team Lead:

* Extract requirements.
* Produce backlog.
* Stop.

---

### 8. Validation Wrappers

Create lightweight wrappers for:

* build
* test
* typecheck

Wrappers should output concise summaries:

PASS

or

FAILED

plus:

* failing files
* failing tests
* compiler errors

Do not emit successful build artifact lists.

---

### 9. Go Engine Improvements

Refactor runCodex()

* configurable timeout
* bounded output buffer
* Codex CLI config overrides
* configurable profiles

Example:

* model_reasoning_effort="medium"
* model_reasoning_summary="none"
* web_search="disabled"

---

### 10. Future Extensions

Design the orchestrator so that future improvements can include:

* semantic repository index
* automatic relevant-file retrieval
* vector search
* incremental summaries
* task dependency graph
* per-agent persistent memory

without changing the core orchestration architecture.

---

## Success Criteria

The orchestrator should:

* minimize cumulative token usage
* eliminate extremely long Codex sessions
* preserve implementation quality
* allow frequent fresh sessions
* keep deterministic external memory
* never rely on long conversation history as persistent state
* be resilient to Codex session resets and context compaction
