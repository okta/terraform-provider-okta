---
name: review-queue
description: Show a summary of GitHub PRs awaiting the user's review, organized by recency. Use this skill whenever the user asks about their PR review queue, pending reviews, what PRs need attention, what they should review, or anything related to seeing open pull requests requesting their review — even across multiple GitHub orgs. Also trigger when the user says things like "what's waiting on me", "show my reviews", or "PR summary".
user_invocable: true
---

# PR Review Queue

Show all GitHub PRs where the user's review has been requested, filtered to human-authored and non-draft, organized by recency (last 7 days vs older).

## How to use

Run the bundled script. It handles everything: fetching PRs via `gh search prs`, filtering out bots and drafts, and rendering a markdown summary.

```bash
bash ~/.claude/skills/review-queue/scripts/review_queue.sh
```

### Options

- `--days N` — Only show PRs updated in the last N days (default: 30)
- `--org ORG` — Limit to a specific GitHub org

### Examples

```bash
# Default: all orgs, last 30 days
bash ~/.claude/skills/review-queue/scripts/review_queue.sh

# Last 7 days only
bash ~/.claude/skills/review-queue/scripts/review_queue.sh --days 7

# Single org
bash ~/.claude/skills/review-queue/scripts/review_queue.sh --org okta
```

## Presenting the output

The script already outputs well-formatted markdown tables. Present its output directly to the user as-is — do not summarize, reformat, or add your own commentary on top. The tables are the deliverable.

The output includes:
- Markdown link to each PR (e.g. `[repo#123](https://github.com/...)`)
- PR title
- Author
- Last update date

Tables are grouped by recency (Recent: Last 7 days, Older PRs), sorted by most recently updated within each group.

## Prerequisites

- `gh` CLI authenticated
- `python3` available
