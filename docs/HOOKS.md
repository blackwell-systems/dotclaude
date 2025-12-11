# Hooks System

dotclaude includes a cross-platform hook system that integrates with Claude Code's lifecycle events. Hooks enable automation and customization without platform-specific scripts.

## Overview

Hooks are executable scripts or built-in commands that run at specific points during Claude Code sessions. The hook system provides:

- **Cross-platform support**: Works on Linux, macOS, and Windows
- **Built-in hooks**: Core functionality implemented in Go
- **Custom hooks**: Add your own scripts in any language
- **Priority ordering**: Control execution order with numeric prefixes

## Hook Types

| Hook Type | Trigger | Use Case |
|-----------|---------|----------|
| `session-start` | When Claude Code session starts | Session info, profile checks, environment validation |
| `post-tool-bash` | After Bash tool execution | Git workflow tips, command notifications |
| `post-tool-edit` | After Edit tool execution | Linting reminders, format checks |
| `pre-tool-bash` | Before Bash tool execution | Safety checks, command validation |
| `pre-tool-edit` | Before Edit tool execution | Pre-edit validation |

## Commands

### Run Hooks

Execute all hooks of a given type:

```bash
dotclaude hook run session-start
dotclaude hook run post-tool-bash
```

### List Hooks

View all available hooks:

```bash
# List all hooks
dotclaude hook list

# List hooks of specific type
dotclaude hook list session-start
```

Example output:
```
session-start:
  [00] session-info (built-in, enabled)
  [10] check-dotclaude (built-in, enabled)
  [20] custom-greeting.sh (enabled)

post-tool-bash:
  [10] git-tips (built-in, enabled)
```

### Initialize Hooks Directory

Create the hooks directory structure:

```bash
dotclaude hook init
```

This creates:
```
~/.claude/hooks/
├── session-start/
├── post-tool-bash/
├── post-tool-edit/
├── pre-tool-bash/
└── pre-tool-edit/
```

## Built-in Hooks

dotclaude includes these built-in hooks implemented in Go:

### session-info (Priority: 00)

Displays session start information:
- Timestamp
- Working directory
- Git branch (if in a git repository)
- Warning if branch is behind main

### check-dotclaude (Priority: 10)

Checks for `.dotclaude` file in the current directory and detects profile mismatches:
- Reads desired profile from `.dotclaude` file
- Compares with currently active profile
- Shows reminder to switch if they differ

### git-tips (Priority: 10, post-tool-bash)

Provides helpful tips after git operations:
- Suggests running `dotclaude sync` after checkout/pull operations

## Custom Hooks

### Adding Custom Hooks

1. Create an executable script in the appropriate hooks directory
2. Use numeric prefix for ordering (e.g., `20-myhook.sh`)
3. Lower numbers run first (00 before 50)

### Example: Custom Greeting

```bash
# ~/.claude/hooks/session-start/20-greeting.sh
#!/bin/bash
echo "Welcome to $(basename $(pwd))!"
echo "Today is $(date +%A)"
```

Make it executable:
```bash
chmod +x ~/.claude/hooks/session-start/20-greeting.sh
```

### Example: PowerShell Hook (Windows)

```powershell
# ~/.claude/hooks/session-start/20-greeting.ps1
Write-Host "Welcome to $((Get-Item .).Name)!"
Write-Host "Today is $((Get-Date).DayOfWeek)"
```

### Supported Script Types

| Platform | Extensions | Notes |
|----------|------------|-------|
| Linux/macOS | `.sh`, `.bash`, any executable | Must have execute permission |
| Windows | `.ps1`, `.cmd`, `.bat`, `.exe` | PowerShell scripts run with `-ExecutionPolicy Bypass` |

### Environment Variables

Hooks receive these environment variables:

| Variable | Description |
|----------|-------------|
| `DOTCLAUDE_REPO_DIR` | Path to dotclaude repository |
| `CLAUDE_DIR` | Path to ~/.claude directory |
| `TOOL_USE_ARGS` | (post-tool hooks) Arguments passed to the tool |

## Priority System

Hooks are executed in alphabetical order by filename. Use numeric prefixes to control order:

| Prefix | Typical Use |
|--------|-------------|
| `00-09` | Critical/first hooks |
| `10-19` | Early hooks |
| `20-49` | Normal priority |
| `50` | Default (no prefix) |
| `50-89` | Later hooks |
| `90-99` | Final hooks |

Example execution order:
```
00-init.sh
10-check-env.sh
20-greeting.sh
50-custom.sh (no prefix = 50)
90-cleanup.sh
```

## Configuration

### Claude Code Integration

The hooks are configured in `~/.claude/settings.json`:

```json
{
  "hooks": {
    "SessionStart": [
      {
        "matcher": "*",
        "hooks": [
          {
            "type": "command",
            "command": "dotclaude hook run session-start"
          }
        ]
      }
    ],
    "PostToolUse": [
      {
        "matcher": "Bash",
        "hooks": [
          {
            "type": "command",
            "command": "dotclaude hook run post-tool-bash"
          }
        ]
      }
    ]
  }
}
```

### Profile-Specific Hooks

Profiles can include custom hooks that are copied during activation:

```
profiles/my-profile/
├── CLAUDE.md
├── settings.json
└── hooks/
    └── session-start/
        └── 20-project-check.sh
```

## Troubleshooting

### Hook Not Running

1. Check the hook is executable: `ls -la ~/.claude/hooks/session-start/`
2. Verify the hook directory exists: `dotclaude hook init`
3. List hooks to confirm registration: `dotclaude hook list session-start`

### Hook Errors

Hooks log warnings but don't fail the session:
```
Hook 20-myhook.sh warning: exit status 1
```

To debug, run the hook directly:
```bash
~/.claude/hooks/session-start/20-myhook.sh
```

### Windows Script Execution

If PowerShell scripts don't run, ensure execution policy allows them:
```powershell
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
```

## Best Practices

1. **Keep hooks fast**: Hooks run synchronously; slow hooks delay session start
2. **Don't fail on errors**: Return success even if optional checks fail
3. **Use meaningful prefixes**: Reserve 00-09 for critical hooks
4. **Test hooks manually**: Run scripts directly before adding to hooks directory
5. **Document your hooks**: Add comments explaining what each hook does

---

**Back to:** [README.md](../README.md) | [COMMANDS.md](COMMANDS.md) | [ARCHITECTURE.md](ARCHITECTURE.md)
