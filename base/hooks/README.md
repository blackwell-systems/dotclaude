# Custom Hooks

This directory contains custom hooks that run alongside dotclaude's built-in hooks.

## Hook Types

| Directory | Trigger | Use Case |
|-----------|---------|----------|
| `session-start/` | When Claude Code session starts | Custom greeting, environment checks |
| `post-tool-bash/` | After Bash tool execution | Post-command notifications |
| `post-tool-edit/` | After Edit tool execution | Linting reminders, format checks |
| `pre-tool-bash/` | Before Bash tool execution | Safety checks |
| `pre-tool-edit/` | Before Edit tool execution | Pre-edit validation |

## Adding Custom Hooks

1. Create an executable script in the appropriate directory
2. Use numeric prefix for ordering (e.g., `20-myhook.sh`)
3. Lower numbers run first (00 before 50)

### Example: Custom Session Start Hook

```bash
# base/hooks/session-start/20-custom-greeting.sh
#!/bin/bash
echo "Welcome to $(basename $(pwd))!"
```

Make it executable:
```bash
chmod +x base/hooks/session-start/20-custom-greeting.sh
```

## Supported Script Types

- **Unix**: `.sh`, `.bash`, or any executable with shebang
- **Windows**: `.ps1` (PowerShell), `.cmd`, `.bat`, `.exe`
- **Cross-platform**: Use the appropriate extension for your OS

## Environment Variables

Hooks receive these environment variables:
- `DOTCLAUDE_REPO_DIR`: Path to dotclaude repository
- `CLAUDE_DIR`: Path to ~/.claude directory
- `TOOL_USE_ARGS`: (post-tool hooks) Arguments passed to the tool

## Built-in Hooks

dotclaude includes these built-in hooks (run via Go):

**session-start:**
- `session-info` (priority 00): Displays session info, git branch
- `check-dotclaude` (priority 10): Checks for profile mismatch

**post-tool-bash:**
- `git-tips` (priority 10): Suggests sync after git operations

View all hooks: `dotclaude hook list`
