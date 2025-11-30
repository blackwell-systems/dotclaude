# dotclaude Usage Guide

Complete reference for using dotclaude, the definitive profile management system for Claude Code.

## Table of Contents

- [Multi-Profile System](#multi-profile-system)
- [Multi-Provider Strategy](#multi-provider-strategy)
- [Auto-Detection](#auto-detection)
- [Command Reference](#command-reference)
- [Long-Lived Feature Branch Management](#long-lived-feature-branch-management)
- [Profile Management](#profile-management)
- [Shell Integration](#shell-integration)
- [Advanced Usage](#advanced-usage)

---

## Multi-Profile System

dotclaude supports multiple work profiles, allowing different configurations for different contexts (OSS projects, proprietary work, employer work, etc.).

### How Profiles Work

**Architecture:**
- **Base** - Shared configuration (git workflow, security, tool usage) that applies to ALL profiles
- **Profiles** - Context-specific additions (coding standards, compliance, tech stacks)
- **Activation** - Merge base + profile → `~/.claude/` when activated

**Getting Started:**
- See `examples/sample-profile/` for a complete example profile
- Copy and customize it to create your own profiles
- The `profiles/` directory starts empty - you create your own based on the examples

### What Happens When You Activate a Profile

1. **Merges CLAUDE.md**: Base guidelines + profile-specific additions
2. **Applies settings.json**: Profile settings override base (if profile has custom settings)
3. **Marks active profile**: Creates `~/.claude/.current-profile` marker
4. **Backs up existing**: Previous config backed up with timestamp

**Example merged CLAUDE.md:**

When you run `dotclaude activate my-project`, the final `~/.claude/CLAUDE.md` contains:

```markdown
# Global Claude Code Instructions

These instructions apply to ALL projects unless overridden...

## Development Standards

### Code Quality
- Write clean, maintainable code
- Follow existing patterns...

### File Operations
- Always use absolute file paths
- Read files before editing...

### Security
- Never commit sensitive data
- Use environment variables...

### Git Practices
- Write clear, descriptive commit messages
- Focus on "why" rather than "what"...

### Tool Usage
- Use Read instead of cat
- Use Edit instead of sed...

### Task Management
- Use TodoWrite for complex tasks
- Mark todos completed immediately...

[... 100+ more lines from base/CLAUDE.md ...]

# =========================================
# Profile-Specific Additions: my-project
# =========================================

## Tech Stack Preferences

### Backend
- Language: Node.js (TypeScript preferred)
- Framework: Express.js or Fastify
- Database: PostgreSQL with Prisma ORM

### Frontend
- Framework: React with TypeScript
- State Management: React Query + Context API
- Styling: Tailwind CSS

## API Design Principles
- Use HTTP methods appropriately
- Return appropriate status codes
- Version APIs in the URL (/api/v1/users)

[... project-specific content continues ...]
```

**Key point:** Claude Code reads the ENTIRE merged file, so it knows both your universal practices (from base) AND your project-specific requirements (from profile).

### Profile Use Cases

**Example Profiles You Might Create:**

**For open source projects:**
- Open source best practices
- Public documentation emphasis
- MIT/Apache licensing guidance
- Community contribution guidelines

**For client/proprietary work:**
- Proprietary code handling
- Internal documentation standards
- Business-specific tech stack
- Private repo security

**For employer work:**
- Corporate compliance policies
- Company coding standards
- Specific frameworks/tools
- Security/audit requirements

### Creating Custom Profiles

Create your own profiles based on the example:

```bash
# Start from the example
cp -r examples/sample-profile profiles/my-project

# Edit to customize
dotclaude edit my-project

# Or create from scratch
mkdir -p profiles/my-new-profile

cat > profiles/my-new-profile/CLAUDE.md << 'EOF'
# Profile: My New Profile

Profile-specific guidelines here...

## Context
Describe the work context this profile is for.

## Standards
Context-specific coding standards.

## Tech Stack
Preferred tools and frameworks for this context.
EOF

# (Optional) Create profile-specific settings.json
cat > profiles/my-new-profile/settings.json << 'EOF'
{
  "hooks": {
    "SessionStart": [...]
  }
}
EOF

dotclaude activate my-new-profile
```

---

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

---

## Auto-Detection

dotclaude can automatically detect when you're working on a project that requires a specific profile.

### The `.dotclaude` File

Place a `.dotclaude` file in your project root to specify which profile should be used:

```bash
cd ~/code/my-my-project
echo "profile: my-project" > .dotclaude
```

**Supported formats:**

```yaml
# YAML-style (recommended)
profile: my-project
```

```bash
# Shell-style
profile=my-project
```

### How It Works

When Claude Code starts a session, it automatically:
1. Checks for `.dotclaude` file in current directory
2. Reads the specified profile name
3. Compares with currently active profile
4. If they differ, displays a reminder to switch

**Example session output with profile mismatch:**

```
=== Claude Code Session Started ===
Fri Nov 29 14:30:00 PST 2024
Working directory: /home/user/code/my-my-project
Git branch: main

╭─────────────────────────────────────────────────────────────╮
│  🍃 Profile Mismatch Detected                               │
╰─────────────────────────────────────────────────────────────╯

  This project uses:    my-project
  Currently active:     work-project

  To activate the project profile:
    dotclaude activate my-project
```

### Use Cases

**1. Team Collaboration**

Commit `.dotclaude` to your repository:

```bash
# In your OSS project
echo "profile: my-project" > .dotclaude
git add .dotclaude
git commit -m "Add dotclaude profile configuration"
```

All team members using dotclaude will be reminded to use the correct profile.

**2. Personal Organization**

Use `.dotclaude` files across your projects:

```
~/code/
├── my-my-project/
│   └── .dotclaude          # profile: my-project
├── proprietary-business/
│   └── .dotclaude          # profile: client-work
└── work-project/
    └── .dotclaude          # profile: work-project
```

Never forget which profile to use for each project.

**3. Context Switching Safety**

When jumping between projects with different contexts, auto-detection prevents mistakes:

```bash
# Working on employer project
cd ~/work-project
# Reminded to use: work-project profile

# Switch to OSS project
cd ~/my-my-project
# Reminded to use: my-project profile
```

### Security

The `.dotclaude` file is validated for security:
- Profile names must be alphanumeric + hyphens/underscores only
- Path traversal attempts are blocked
- Profile existence is verified before displaying reminder
- **Detection only** - never auto-activates without your confirmation

### Git Integration

**Should you commit `.dotclaude`?**

✅ **Commit when:**
- Team project with shared profile
- Want consistent setup across machines
- OSS project with documented standards

❌ **Don't commit when:**
- Personal project with unique profile
- Profile is machine-specific
- Team members use different dotclaude setups

Add to `.gitignore` if needed:
```bash
echo ".dotclaude" >> .gitignore
```

### Complete Documentation

For full details, see **[docs/DOTCLAUDE-FILE.md](DOTCLAUDE-FILE.md)**

---

## Command Reference

### Profile Management Commands

#### `dotclaude show`
Show current active profile and configuration status.

```bash
dotclaude show
```

**Output:**
```
╭─────────────────────────────────────────────────────────────╮
│  🌲 dotclaude                                               │
╰─────────────────────────────────────────────────────────────╯

  Active Profile: my-project

  Configuration:
    • CLAUDE.md: 245 lines
    • settings.json: configured

  Location:
    • ~/.claude/

╭─────────────────────────────────────────────────────────────╮
│  🍃 Tip: Run 'dotclaude switch' to change profiles         │
╰─────────────────────────────────────────────────────────────╯
```

#### `dotclaude list`
List all available profiles.

```bash
dotclaude list
# Aliases: dotclaude ls
```

**Output:**
```
╭─────────────────────────────────────────────────────────────╮
│  🌲 dotclaude                                               │
╰─────────────────────────────────────────────────────────────╯

  Available Profiles:

    ● my-project (active)
    ○ client-work
    ○ work-project

  Total: 3 profiles

╭─────────────────────────────────────────────────────────────╮
│  🍃 Tip: Use 'dotclaude activate <name>' or 'dotclaude switch'│
╰─────────────────────────────────────────────────────────────╯
```

#### `dotclaude activate <profile-name>`
Activate a specific profile.

```bash
dotclaude activate my-project
# Aliases: dotclaude use <profile-name>
```

**Output:**
```
╭─────────────────────────────────────────────────────────────╮
│  🌲 dotclaude                                               │
╰─────────────────────────────────────────────────────────────╯

  Activating profile: my-project

  [1/3] Backed up existing CLAUDE.md
  [2/3] Merged base + profile configuration
  [3/3] Applied profile settings

╭─────────────────────────────────────────────────────────────╮
│  ✓ Profile 'my-project' activated               │
╰─────────────────────────────────────────────────────────────╯

  Configuration deployed to: /home/user/.claude

  Verify with:
    • dotclaude show
    • cat ~/.claude/CLAUDE.md

╭─────────────────────────────────────────────────────────────╮
│  🍃 Tip: Happy coding!                                      │
╰─────────────────────────────────────────────────────────────╯
```

**Dry-Run Mode:**

Preview changes before activating:

```bash
dotclaude activate my-project --dry-run
# Or: dotclaude activate my-project --preview
```

**Output:**
```
╭─────────────────────────────────────────────────────────────╮
│  🌲 dotclaude                                               │
╰─────────────────────────────────────────────────────────────╯

  DRY-RUN MODE - Preview changes for: my-project

  Changes that would be made:

    [BACKUP] Existing CLAUDE.md would be backed up
    [MERGE] CLAUDE.md: base + my-project
            Base: 150 lines
            Profile: 45 lines
            Result: ~200 lines

    [APPLY] settings.json: profile-specific
            [BACKUP] Existing settings.json would be backed up

    [SET] Active profile: my-project

  Files that would be modified:
    • ~/.claude/CLAUDE.md
    • ~/.claude/settings.json
    • ~/.claude/.current-profile
    • ~/.claude/CLAUDE.md.backup.20241129-HHMMSS

╭─────────────────────────────────────────────────────────────╮
│  🍃 Tip: Run without --dry-run to apply changes            │
╰─────────────────────────────────────────────────────────────╯
```

**Use case:** Preview what will change before committing to a profile switch.

#### `dotclaude switch`
Interactive profile switcher.

```bash
dotclaude switch
# Aliases: dotclaude select
```

**Output:**
```
╭─────────────────────────────────────────────────────────────╮
│  🌲 dotclaude                                               │
╰─────────────────────────────────────────────────────────────╯

  Select a profile to activate:

    [1] my-project (active)
    [2] client-work
    [3] work-project

  Enter number (or 'q' to quit): 2

  Activating profile: client-work
  ...
```

#### `dotclaude create <profile-name>`
Create a new profile.

```bash
dotclaude create my-new-profile
# Aliases: dotclaude new <profile-name>
```

**Output:**
```
╭─────────────────────────────────────────────────────────────╮
│  🌲 dotclaude                                               │
╰─────────────────────────────────────────────────────────────╯

  Creating new profile: my-new-profile

╭─────────────────────────────────────────────────────────────╮
│  ✓ Profile 'my-new-profile' created                        │
╰─────────────────────────────────────────────────────────────╯

  Location: /home/user/code/dotclaude/profiles/my-new-profile

  Next steps:
    • dotclaude edit my-new-profile
    • dotclaude activate my-new-profile

╭─────────────────────────────────────────────────────────────╮
│  🍃 Tip: Edit the profile to add your guidelines           │
╰─────────────────────────────────────────────────────────────╯
```

#### `dotclaude edit [profile-name]`
Edit a profile's CLAUDE.md in $EDITOR.

```bash
# Edit current active profile
dotclaude edit

# Edit specific profile
dotclaude edit my-project
```

Opens the profile's `CLAUDE.md` in your configured editor (`$EDITOR` or `nano` by default).

#### `dotclaude diff <profile1> [profile2]`
Compare two profiles or compare current profile with another.

```bash
# Compare two profiles
dotclaude diff my-project client-work

# Compare current profile with another
dotclaude diff work-project
```

**Output:**
```
╭─────────────────────────────────────────────────────────────╮
│  🌲 dotclaude                                               │
╰─────────────────────────────────────────────────────────────╯

  Comparing profiles: my-project vs client-work

  CLAUDE.md differences:

    Differences found:

      @@ -1,5 +1,5 @@
      -# Profile: my-project
      +# Profile: client-work

      -Open source best practices
      +Proprietary code handling

      -## Licensing
      -All projects use MIT or Apache 2.0 licenses
      +## Confidentiality
      +Proprietary code, no public sharing

      ... (150 more lines)

    Tip: See full diff with:
      diff -u profiles/my-project/CLAUDE.md profiles/client-work/CLAUDE.md

  settings.json:

    ✓ No differences

╭─────────────────────────────────────────────────────────────╮
│  🍃 Tip: Use 'dotclaude activate <profile> --dry-run' to preview activation│
╰─────────────────────────────────────────────────────────────╯
```

**Use cases:**
- See differences before switching profiles
- Compare OSS vs proprietary guidelines
- Verify profile-specific settings

#### `dotclaude restore`
Restore from backup interactively.

```bash
dotclaude restore
```

**Output:**
```
╭─────────────────────────────────────────────────────────────╮
│  🌲 dotclaude                                               │
╰─────────────────────────────────────────────────────────────╯

  Backup Restoration

  Available backups:

  CLAUDE.md backups:
    [1] 20241129-143022 (24K)
    [2] 20241129-120145 (22K)
    [3] 20241128-183045 (23K)

  settings.json backups:
    [4] 20241129-143022 (4.2K)
    [5] 20241129-120145 (3.8K)

  Select backup to restore (or 'q' to quit): 1

  ⚠  This will overwrite:
    /home/user/.claude/CLAUDE.md

  Continue? (y/N): y

  [BACKUP] Current file backed up to:
    CLAUDE.md.backup.20241129-150322
  [RESTORE] Restored from: CLAUDE.md.backup.20241129-143022
  [UPDATE] Active profile: my-project

╭─────────────────────────────────────────────────────────────╮
│  ✓ Backup restored successfully                            │
╰─────────────────────────────────────────────────────────────╯

  Restored: /home/user/.claude/CLAUDE.md

╭─────────────────────────────────────────────────────────────╮
│  🍃 Tip: Verify with 'dotclaude show'                      │
╰─────────────────────────────────────────────────────────────╯
```

**Features:**
- Interactive selection from all available backups
- Shows timestamp and file size
- Separate lists for CLAUDE.md and settings.json backups
- Creates backup of current file before restoring
- Auto-detects and updates active profile marker

**Use cases:**
- Undo a profile switch
- Recover from accidental changes
- Roll back to previous configuration

### Git Workflow Commands

#### `dotclaude sync`
Run the interactive feature branch sync tool.

```bash
dotclaude sync
```

See [Long-Lived Feature Branch Management](#long-lived-feature-branch-management) for details.

#### `dotclaude branches`
Check status of all branches.

```bash
dotclaude branches
# Aliases: dotclaude br
```

**Output:**
```
Checking branches against main...

  feature-add-auth              10 ahead, 5 behind
  feature-refactor              2 ahead, 12 behind
  feature-experimental          0 ahead, 3 behind
```

### System Commands

#### `dotclaude version`
Show dotclaude version.

```bash
dotclaude version
# Aliases: dotclaude -v, dotclaude --version
```

**Output:**
```
╭─────────────────────────────────────────────────────────────╮
│  🌲 dotclaude                                               │
╰─────────────────────────────────────────────────────────────╯

  dotclaude version 1.0.0

  The definitive profile management system for Claude Code

  Repository: /home/user/code/dotclaude
  Configuration: /home/user/.claude
```

#### `dotclaude help [command]`
Show help information.

```bash
# General help
dotclaude help
# Aliases: dotclaude -h, dotclaude --help

# Command-specific help (coming soon)
dotclaude help activate
```

---

## Long-Lived Feature Branch Management

Automated workflow for keeping feature branches in sync with main when working iteratively.

### The Problem

**Scenario:**
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
$ dotclaude sync
Currently on main

Feature branches that are behind:
  - feature-add-auth (behind by 5 commits)
  - feature-refactor (behind by 12 commits)

To sync a branch, run:
  git checkout <branch-name>
  dotclaude sync
```

**When on feature branch:**

```bash
$ dotclaude sync
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

These are sourced shell functions available after adding to your `~/.bashrc` or `~/.zshrc`.

**`sync-feature-branch` - Direct sync command**

```bash
# Wrapper that calls dotclaude sync
sync-feature-branch
```

**`check-branches` - Quick Status Check**

```bash
$ check-branches
Checking branches against main...

  feature-add-auth              10 ahead, 5 behind
  feature-refactor              2 ahead, 12 behind
  feature-experimental          0 ahead, 3 behind
```

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
$ dotclaude sync  # Choose rebase or merge

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

### Key Benefits

1. **Never forget to sync** - Automated hooks remind you every session
2. **No mistakes** - Interactive tool validates safety before executing
3. **Flexible** - Choose rebase or merge based on your needs
4. **Status visibility** - Always know which branches need attention
5. **Guided workflows** - `pr-merged` walks you through the entire process
6. **Conflict prevention** - Keep branches synced = fewer merge conflicts
7. **Team friendly** - Works across multiple developers and projects

---

## Profile Management

### Deployed Structure

After activation, your `~/.claude/` directory contains:

```
~/.claude/
├── .current-profile               # Active profile marker
├── CLAUDE.md                      # Base + Profile merged
├── settings.json                  # Base or Profile settings
├── scripts/                       # Management scripts
│   ├── dotclaude                 # Main CLI
│   ├── sync-feature-branch.sh
│   ├── shell-functions.sh
│   ├── activate-profile.sh
│   ├── profile-management.sh
│   └── lib/
│       └── validation.sh
└── agents/                        # Shared agents
    └── best-in-class-gap-analysis/
```

### Backups

When switching profiles, dotclaude automatically backs up your existing configuration:

```
~/.claude/
├── CLAUDE.md                      # Current active
├── CLAUDE.md.backup.20241129-143022
├── CLAUDE.md.backup.20241129-120145
├── settings.json                  # Current active
└── settings.json.backup.20241129-143022
```

**Backup behavior:**
- Only created when switching between different profiles
- Re-activating the same profile updates in place (no backup)
- Limited to 5 most recent backups (oldest deleted automatically)
- Secure permissions: `chmod 600` (only you can read)

### Global vs Project Configuration

**Global (`~/.claude/`):**
- Shared standards across all projects
- Applied to every Claude Code session
- Good for: organization policies, tool preferences, git workflows

**Project (`.claude/` in project root):**
- Project-specific overrides
- Team-shared via git
- Good for: tech stack choices, project standards

**Settings Precedence (highest → lowest):**
1. Enterprise policies
2. CLI arguments
3. Project `.claude/settings.local.json` (gitignored)
4. Project `.claude/settings.json` (team-shared)
5. Global `~/.claude/settings.json` (from dotclaude)

### Global Instructions (CLAUDE.md)

Loaded at session start for ALL projects. Use for:
- Standard coding conventions
- Common development guidelines
- Organization-wide policies
- Git workflow practices
- Security standards

### Global Hooks (settings.json)

Execute automatically across all projects:
- `SessionStart` - Session initialization
- `PreToolUse` - Before tool execution
- `PostToolUse` - After tool completion
- `UserPromptSubmit` - When user submits a message

**Example: SessionStart hook**
```json
{
  "hooks": {
    "SessionStart": [
      {
        "matcher": "*",
        "hooks": [
          {
            "type": "command",
            "command": "echo '=== Session Started ===' && date"
          }
        ]
      }
    ]
  }
}
```

### Shared Agents

Available in all projects without per-project configuration.

**Included agent:**
- `best-in-class-gap-analysis` - Competitive analysis agent that researches industry standards

---

## Shell Integration

### Adding to Shell

Add sourced functions to your shell for convenience:

**For bash (`~/.bashrc`):**
```bash
export DOTCLAUDE_REPO_DIR="$HOME/code/dotclaude"
export PATH="$HOME/.local/bin:$PATH"

# Optional: Source convenience functions
if [ -f "$HOME/.claude/scripts/shell-functions.sh" ]; then
    source "$HOME/.claude/scripts/shell-functions.sh"
fi
```

**For zsh (`~/.zshrc`):**
```bash
export DOTCLAUDE_REPO_DIR="$HOME/code/dotclaude"
export PATH="$HOME/.local/bin:$PATH"

# Optional: Source convenience functions
if [ -f "$HOME/.claude/scripts/shell-functions.sh" ]; then
    source "$HOME/.claude/scripts/shell-functions.sh"
fi
```

### Shell Compatibility

**Main CLI (`dotclaude`):**
- Uses `#!/bin/bash` shebang
- Runs in bash regardless of your shell
- Works in any shell environment (bash, zsh, fish, etc.)

**Sourced functions:**
- POSIX-compatible syntax
- Tested with bash 5.x and zsh 5.x
- Use `[[ ]]` conditionals (bash/zsh)
- `export -f` fails gracefully in zsh (not needed, functions already available)

**Functions available after sourcing:**
- `sync-feature-branch` - Sync current branch
- `check-branches` - Check all branches
- `pr-merged` - Post-PR workflow
- `list-feature-branches` - List branches with status

---

## Advanced Usage

### Environment Variables

**`DOTCLAUDE_REPO_DIR`**
- Location of your dotclaude repository
- Default: `$HOME/code/dotclaude`
- Used by CLI to find profiles and base configuration

**`EDITOR`**
- Your preferred text editor
- Default: `nano`
- Used by `dotclaude edit` command

**Example:**
```bash
export DOTCLAUDE_REPO_DIR="$HOME/repos/dotclaude"
export EDITOR="vim"
```

### Installation Options

```bash
# Basic installation (interactive)
./install.sh

# Force overwrite existing files
./install.sh --force

# Non-interactive mode (for automation/CI)
./install.sh --non-interactive

# Show help
./install.sh --help
```

### Updating dotclaude

```bash
# In your dotclaude repo
git pull origin main

# Redeploy (updates CLI and scripts)
./install.sh --force

# Or just update CLI
cp base/scripts/dotclaude ~/.local/bin/dotclaude
chmod +x ~/.local/bin/dotclaude
```

### Sharing Profiles Across Machines

```bash
# Machine 1: Push changes
cd ~/code/dotclaude
git add profiles/my-profile/
git commit -m "Update my-profile guidelines"
git push

# Machine 2: Pull and activate
cd ~/code/dotclaude
git pull
dotclaude activate my-profile
```

### Troubleshooting

**"Profile not found"**
```bash
# Check DOTCLAUDE_REPO_DIR is set correctly
echo $DOTCLAUDE_REPO_DIR

# List available profiles
dotclaude list

# Check profile directory exists
ls -la $DOTCLAUDE_REPO_DIR/profiles/
```

**"Command not found: dotclaude"**
```bash
# Check if ~/.local/bin is in PATH
echo $PATH | grep -o "$HOME/.local/bin"

# If not, add to shell config
echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc

# Or run installer again
./install.sh
```

**"Another dotclaude operation in progress"**
```bash
# Lock file exists from previous operation
# Wait for other operation to complete, or remove stale lock:
rm ~/.claude/.lock
```

### Security Notes

**File Permissions:**
- Backups: `chmod 600` (readable only by you)
- Lock files: Prevent concurrent modifications
- Input validation: Prevents path traversal attacks

**Hook Security:**
- Hooks run with your user credentials
- Review hook commands carefully
- Keep sensitive data in `.local.json` (gitignored)
- Don't commit API keys or credentials

**See also:** [docs/SECURITY-AUDIT.md](docs/SECURITY-AUDIT.md)

---

## Notes

- Global configs in `~/.claude/` apply to ALL projects
- Project `.claude/settings.json` overrides global settings
- Hooks run with your environment credentials - review carefully
- Keep sensitive data (API keys) in project-level `.claude/settings.local.json`
- Backup files are kept in `~/.claude/` for safety
- Only 5 most recent backups retained to save space

---

**Back to:** [README.md](../README.md) | **See also:** [ARCHITECTURE.md](ARCHITECTURE.md)
