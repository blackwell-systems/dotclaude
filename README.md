# Claude Code Global Configuration

This repository manages global Claude Code configurations that can be shared across all projects.

## Overview

This repo contains:
- Global CLAUDE.md instructions
- Global settings.json (hooks, preferences)
- Shared agents/skills
- Deployment scripts

## Multi-Provider Strategy

### Provider Support
These configs work with BOTH:
- **AWS Bedrock** (e.g., `us.anthropic.claude-sonnet-4-5-20250929-v1:0`)
- **Claude Max** (e.g., `claude-sonnet-4-5-20250929`)

### Configuration Strategy

**Global configs are provider-agnostic:**
- Don't hard-code provider-specific model IDs in global settings
- Let projects specify models via `.claude/settings.json`
- Hooks and agents work identically across providers

**Settings Precedence (highest → lowest):**
1. Enterprise policies
2. CLI arguments
3. Project `.claude/settings.local.json` (gitignored)
4. Project `.claude/settings.json` (team-shared)
5. Global `~/.claude/settings.json` (from this repo)

## Domain-Specific Configurations

### Current Approach: Project-Level Override

For separate work contexts (e.g., Blackwell Systems vs employer work):

**Option 1: Single Global + Project Override (Recommended)**
- Keep ONE global `~/.claude/` with shared tools (hooks, common agents)
- Use project-level `.claude/settings.json` to differentiate domains
- Simple, maintainable, no context switching needed

**Option 2: Symbolic Link Switching**
```bash
# Shell functions
use-blackwell() { ln -sf ~/.claude-blackwell ~/.claude; }
use-work() { ln -sf ~/.claude-work ~/.claude; }
```

**Option 3: Environment Variable (Experimental)**
```bash
export CLAUDE_CONFIG_DIR=~/.claude-blackwell
export CLAUDE_CONFIG_DIR=~/.claude-work
```

## Repository Structure

```
CLAUDE/
├── README.md                    # This file
├── install.sh                   # Deploy configs to ~/.claude/
├── global/
│   ├── CLAUDE.md               # Global instructions for all projects
│   ├── settings.json           # Global settings & hooks
│   ├── scripts/
│   │   ├── sync-feature-branch.sh    # Interactive branch sync tool
│   │   └── shell-functions.sh        # Shell helper functions
│   └── agents/
│       └── best-in-class-gap-analysis/
│           └── definition.json
└── docs/
    └── strategy.md             # Detailed configuration strategy
```

## Features

### Long-Lived Feature Branch Management

Automated workflow for keeping feature branches in sync with main when working iteratively.

**The Problem:** When you merge a PR and continue working on the same feature branch, it falls behind main, causing merge conflicts and confusion.

**The Solution:** Automated hooks and helper tools that keep branches synchronized.

**What's Included:**

1. **Interactive Sync Tool** (`sync-feature-branch`)
   - Guides you through rebasing or merging with main
   - Checks for uncommitted changes
   - Offers to push after sync
   - Prevents mistakes

2. **Automated Reminders**
   - SessionStart hook warns if current branch is behind
   - PostToolUse hook reminds after git operations
   - Never miss a sync opportunity

3. **Shell Helper Functions**
   - `sync-feature-branch` - Interactive sync current branch
   - `check-branches` - Show status of all branches
   - `pr-merged` - Guided workflow after PR merge
   - `list-feature-branches` - List all branches with status

**Example Workflow:**

```bash
# Your feature branch PR just got merged
git checkout main
git pull

# Back to your feature branch
git checkout feature-branch

# Interactive sync (choose rebase or merge)
sync-feature-branch

# Continue working on the same branch, now synced with main
```

**Or use the guided workflow:**

```bash
# On your feature branch after PR merge
pr-merged
# This automates: checkout main → pull → checkout feature → sync
```

---

## How It Works

### The Problem Scenario

```
1. You create feature-branch from main
2. You make commits and open a PR
3. PR gets merged into main
4. You go back to feature-branch to continue work
5. ❌ feature-branch is now BEHIND main (missing the merged code)
```

Without syncing, your feature branch accumulates merge conflicts and becomes harder to maintain.

### The Solution - 3 Layers of Automation

#### Layer 1: Automated Detection (Hooks)

**SessionStart Hook** - Runs every time Claude Code starts/resumes

```bash
# Automatically checks at session start:
1. Are we in a git repo?
2. What branch are we on?
3. Is this branch behind main?
4. If yes → Display warning with commit count
```

**Example output:**
```
=== Claude Code Session Started ===
Fri Nov 29 12:15:00 PST 2024
Working directory: /home/user/code/my-project
Git branch: feature-add-auth
⚠️  Branch is 5 commits behind main - consider running: sync-feature-branch
```

**PostToolUse Hook** - Runs after git operations

```bash
# Triggers after: git checkout main, git pull, etc.
1. Detect: Did we just update main?
2. If yes → Remind: "Feature branches may be behind"
```

#### Layer 2: Interactive Sync Tool

**When on main branch:**

```bash
$ sync-feature-branch
Currently on main

Feature branches that are behind:
  - feature-add-auth (behind by 5 commits)
  - feature-refactor (behind by 12 commits)

To sync a branch, run:
  git checkout <branch-name>
  sync-feature-branch
```

**When on feature branch:**

```bash
$ sync-feature-branch
Current branch: feature-add-auth
Status: 10 commits ahead, 5 commits behind main

Branch is 5 commits behind main

Choose sync method:
  1) Rebase (cleaner history, requires force push)
  2) Merge (preserves history, no force push needed)
  3) Cancel

Selection (1/2/3): 1

Rebasing feature-add-auth onto main...
✓ Rebase successful

To push changes:
  git push --force-with-lease

Push now? (y/N): y
✓ Branch synced and pushed
```

**What the tool does internally:**

1. **Safety checks:**
   - Verify git repo exists
   - Check for uncommitted changes (fails if found)
   - Calculate commits ahead/behind main

2. **User choice:**
   - **Rebase** → cleaner history, linear commits, requires force push
   - **Merge** → preserves all history, creates merge commit, regular push

3. **Execute sync:**
   - Fetches latest from origin
   - Rebases or merges depending on choice
   - Handles conflicts (pauses if needed)
   - Offers to push changes

4. **Conflict handling:**
   ```bash
   # If conflicts occur during rebase
   Rebase failed. Resolve conflicts and run:
     git rebase --continue
     git push --force-with-lease
   ```

#### Layer 3: Workflow Helpers

**`pr-merged` - Post-PR Workflow**

```bash
# You're on feature-branch, PR just merged
$ pr-merged

PR merged workflow:
  1. Switching to main and pulling latest
  2. Switching back to feature-branch and syncing

Continue? (y/N): y

# Automatically executes:
# git checkout main && git pull
# git checkout feature-branch
# sync-feature-branch  # Interactive from here
```

**`check-branches` - Quick Status Check**

```bash
$ check-branches
Checking branches against main...

  feature-add-auth              10 ahead, 5 behind
  feature-refactor              2 ahead, 12 behind
  feature-experimental          0 ahead, 3 behind
```

**`list-feature-branches` - Detailed Branch View**

```bash
$ list-feature-branches
Feature branches:

  BRANCH                         AHEAD           BEHIND          LAST COMMIT
  ------                         -----           ------          -----------
  feature-add-auth              10              5               2 hours ago
  feature-refactor              2               12              3 days ago
  feature-experimental          0               3               1 week ago
```

---

### Complete Workflow Example

**Day 1: Initial work**
```bash
git checkout -b feature-add-auth
# ... make changes ...
git commit -m "Add initial auth"
git push -u origin feature-add-auth
# Open PR
```

**Day 2: PR merged, continue work**
```bash
# Start Claude Code session
# 👆 Hook warns: "Branch is 5 commits behind main"

# Option A: Use guided helper
$ pr-merged
# Guides you through: main → pull → feature-branch → sync

# Option B: Manual approach
$ git checkout main && git pull
$ git checkout feature-add-auth
$ sync-feature-branch  # Choose rebase or merge

# Now your feature branch is synced!
# Continue working...
git commit -m "Add OAuth support"
git push
```

**Day 3+: Continued iterations**
```bash
# Same process - branch stays synced
# No accumulated merge conflicts
# Clean, maintainable feature branch
```

---

### Technical Details

#### Rebase vs Merge - What's the Difference?

**Rebase (Option 1):**
```
Before:
main:     A---B---C---D---E
               \
feature:        F---G---H

After rebase:
main:     A---B---C---D---E
                           \
feature:                    F'---G'---H'

- Commits F,G,H are "replayed" on top of E
- History is linear and clean
- Commits get new SHAs (F', G', H')
- Requires force push (rewrites history)
```

**Pros:** Clean, linear history. Easy to follow git log.
**Cons:** Rewrites history. Requires force push.
**Use when:** You want clean history and haven't shared commits with others.

**Merge (Option 2):**
```
Before:
main:     A---B---C---D---E
               \
feature:        F---G---H

After merge:
main:     A---B---C---D---E
               \           \
feature:        F---G---H---M

- Creates merge commit M
- Preserves all history
- Original commits unchanged
- Regular push (no force needed)
```

**Pros:** Preserves complete history. No force push needed.
**Cons:** More complex git log. Extra merge commits.
**Use when:** You want to preserve exact history or have shared the branch.

#### Why Force Push with Rebase?

```bash
# Before rebase (on remote):
origin/feature-branch: F---G---H

# After rebase (locally):
local/feature-branch: F'---G'---H'

# F', G', H' are NEW commits with different SHAs
# Git sees these as diverged histories
# Need --force-with-lease to tell git "I know what I'm doing"
```

**Why `--force-with-lease` instead of `--force`?**

- `--force` blindly overwrites remote branch (dangerous)
- `--force-with-lease` checks that no one else pushed to the branch
- Prevents accidentally overwriting teammate's work

#### Behind the Scenes: Hook Mechanics

**How hooks execute:**

1. Claude Code starts/resumes
2. Reads `~/.claude/settings.json`
3. Finds hooks for the current event (SessionStart, PostToolUse, etc.)
4. Executes bash commands defined in hooks
5. Displays output to user
6. Claude Code continues normal operation

**Example hook configuration:**
```json
{
  "hooks": {
    "SessionStart": [
      {
        "matcher": "*",
        "hooks": [
          {
            "type": "command",
            "command": "if git rev-parse --git-dir > /dev/null 2>&1; then ..."
          }
        ]
      }
    ]
  }
}
```

Hooks are just **bash scripts** that run at specific lifecycle events. They have access to:
- Current working directory
- Environment variables
- Git repository information
- Tool use context (for PostToolUse hooks)

---

### Key Benefits

1. **Never forget to sync** - Automated hooks remind you every session
2. **No mistakes** - Interactive tool validates safety before executing
3. **Flexible** - Choose rebase or merge based on your needs
4. **Status visibility** - Always know which branches need attention
5. **Guided workflows** - `pr-merged` walks you through the entire process
6. **Conflict prevention** - Keep branches synced = fewer merge conflicts
7. **Team friendly** - Works across multiple developers and projects

---

## Installation

### 1. Deploy Configs

Deploy these configs to your `~/.claude/` directory:

```bash
./install.sh
```

Or manually:
```bash
mkdir -p ~/.claude/{agents,scripts}
cp global/CLAUDE.md ~/.claude/
cp global/settings.json ~/.claude/
cp -r global/agents/* ~/.claude/agents/
cp -r global/scripts/* ~/.claude/scripts/
chmod +x ~/.claude/scripts/*.sh
```

### 2. Enable Shell Functions (Optional but Recommended)

To use `sync-feature-branch` and other git workflow helpers, add to your `~/.bashrc` or `~/.zshrc`:

```bash
# Claude Code git workflow functions
if [ -f "$HOME/.claude/scripts/shell-functions.sh" ]; then
    source "$HOME/.claude/scripts/shell-functions.sh"
fi
```

Then restart your shell:
```bash
source ~/.bashrc  # or source ~/.zshrc
```

## Usage

### Global Instructions (CLAUDE.md)
Loaded at session start for ALL projects. Use for:
- Standard coding conventions
- Common development guidelines
- Organization-wide policies

### Global Hooks (settings.json)
Execute automatically across all projects:
- `SessionStart` - Session initialization
- `PreToolUse` - Before tool execution
- `PostToolUse` - After tool completion

### Shared Agents
Available in all projects without per-project configuration.

## Maintenance

1. Edit configs in this repo
2. Test changes in a project
3. Commit to version control
4. Run `./install.sh` to deploy
5. Share repo across machines for consistent setup

## Notes

- Global configs in `~/.claude/` apply to ALL projects
- Project `.claude/settings.json` overrides global settings
- Hooks run with your environment credentials - review carefully
- Keep sensitive data (API keys) in project-level `.claude/settings.local.json`
