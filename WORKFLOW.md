---
polling:
  interval_ms: 30000
workspace:
  root: .worktrees
agent:
  max_concurrent_agents: 1
  max_turns: 20
  continuation_turns_enabled: false
  max_retry_attempts: 3
  max_retry_backoff_ms: 300000
codex:
  command: codex app-server
  read_timeout_ms: 5000
  turn_timeout_ms: 3600000
  stall_timeout_ms: 300000
---

# Task

Issue ID: {{ issue.id }}
Title: {{ issue.title }}
Status: {{ issue.status }}
Priority: {{ issue.priority }}
Assignee: {{ issue.assignee }}
Attempt: {{ attempt }}

## Description

{{ issue.description }}

## Required Flow

1. Confirm the task scope from the issue title and description above.
2. Create or switch to an isolated task branch/worktree before editing.
3. Make focused changes that satisfy the issue.
4. Run the narrowest useful verification first, then broaden checks when shared behavior is affected.
5. Commit the change and create or update a pull request.
6. Leave a progress or handoff comment on the issue.
7. Move the issue to `review` when the pull request is ready for human review.

For repository-specific contribution rules, follow [CONTRIBUTE.md](CONTRIBUTE.md).

