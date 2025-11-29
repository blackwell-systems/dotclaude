# .dotclaude File Format

The `.dotclaude` file is a project-level configuration file that specifies which dotclaude profile should be used for a project.

## Purpose

When working on multiple projects with different contexts (OSS, proprietary, employer work), it's easy to forget which profile should be active. The `.dotclaude` file automates profile detection and reminds you to switch.

## Usage

Place a `.dotclaude` file in the root of your project directory:

```bash
cd ~/code/my-project
cat > .dotclaude << EOF
profile: blackwell-systems-oss
EOF
```

## File Format

The `.dotclaude` file supports two formats:

### YAML-Style (Recommended)

```yaml
profile: blackwell-systems-oss
```

### Shell-Style

```bash
profile=blackwell-systems-oss
```

Both formats are equivalent. Use whichever your team prefers.

## How It Works

When Claude Code starts a session, the SessionStart hook:

1. Checks if `.dotclaude` file exists in current directory
2. Reads the specified profile name
3. Compares with currently active profile
4. If they differ, displays a reminder to switch

**Example output when profile mismatch detected:**

```
=== Claude Code Session Started ===
Fri Nov 29 14:30:00 PST 2024
Working directory: /home/user/code/my-oss-project
Git branch: main

╭─────────────────────────────────────────────────────────────╮
│  🍃 Profile Mismatch Detected                               │
╰─────────────────────────────────────────────────────────────╯

  This project uses:    blackwell-systems-oss
  Currently active:     best-western

  To activate the project profile:
    dotclaude activate blackwell-systems-oss
```

## Security

The `.dotclaude` file is validated for security:

- **Profile name validation**: Only alphanumeric characters, hyphens, and underscores allowed
- **No path traversal**: Prevents `profile: ../../etc/passwd`
- **Profile existence check**: Verifies profile exists in your dotclaude repository
- **No auto-execution**: Only shows reminder, never auto-activates (you control when to switch)

## Use Cases

### 1. Team Collaboration

Commit `.dotclaude` to your project repository:

```bash
# In your OSS project
echo "profile: blackwell-systems-oss" > .dotclaude
git add .dotclaude
git commit -m "Add dotclaude profile configuration"
git push
```

Now all team members using dotclaude will be reminded to use the correct profile.

### 2. Multiple Projects

Organize your projects with appropriate profiles:

```
~/code/
├── my-oss-project/
│   └── .dotclaude          # profile: blackwell-systems-oss
├── proprietary-business/
│   └── .dotclaude          # profile: blackwell-systems
└── employer-work/
    └── .dotclaude          # profile: best-western
```

### 3. Context Switching

When jumping between projects, you'll automatically be reminded:

```bash
cd ~/code/my-oss-project
# Session starts, detects blackwell-systems-oss profile
# Reminds you to switch if needed

cd ~/code/employer-work
# Session starts, detects best-western profile
# Reminds you to switch if needed
```

## Git Integration

### Should You Commit .dotclaude?

**Commit when:**
- Team project where everyone should use same profile
- OSS project with shared standards
- Want consistent setup across machines

**Don't commit when:**
- Personal project with unique profile
- Profile is machine-specific
- Different team members use different dotclaude setups

### Gitignore

If you don't want to commit `.dotclaude`:

```bash
echo ".dotclaude" >> .gitignore
```

## Advanced Usage

### Multiple Environments

If your project has different environments, use subdirectories:

```
my-project/
├── .dotclaude              # profile: blackwell-systems-oss
├── internal/
│   └── .dotclaude          # profile: blackwell-systems
└── docs/
    └── .dotclaude          # profile: documentation-profile
```

Detection works based on current working directory when Claude Code starts.

### Validation

Profile names must match these rules:
- Only letters (a-z, A-Z)
- Only numbers (0-9)
- Only hyphens (-)
- Only underscores (_)

**Valid:**
```
profile: my-oss-profile
profile: work_project
profile: client-123
```

**Invalid:**
```
profile: ../../../etc/passwd    # Path traversal
profile: my profile              # Space
profile: profile@company         # Special char
```

## Comparison with Other Tools

Similar to:
- **direnv** (`.envrc`) - Auto-loads environment variables per directory
- **rbenv** (`.ruby-version`) - Specifies Ruby version per project
- **nvm** (`.nvmrc`) - Specifies Node.js version per project
- **pyenv** (`.python-version`) - Specifies Python version per project

The `.dotclaude` file follows the same pattern: project-level configuration that helps maintain consistency.

## Troubleshooting

### Profile not found

```
⚠️  Profile 'my-profile' specified in .dotclaude not found
   Available profiles: dotclaude list
```

**Solution:** Either create the profile or fix the typo in `.dotclaude`

### Invalid profile name

```
⚠️  Invalid profile name in .dotclaude: my profile
   Profile names must contain only letters, numbers, hyphens, and underscores
```

**Solution:** Rename profile to use valid characters (e.g., `my-profile` instead of `my profile`)

### Detection not working

1. Verify `.dotclaude` file exists: `ls -la .dotclaude`
2. Check file contents: `cat .dotclaude`
3. Ensure SessionStart hook is enabled in `~/.claude/settings.json`
4. Verify check-dotclaude.sh script exists: `ls ~/.claude/scripts/check-dotclaude.sh`

## Examples

### Minimal

```yaml
profile: my-profile
```

### With Comments (YAML-style)

```yaml
# OSS project - public standards
profile: blackwell-systems-oss
```

### Shell-style with Comments

```bash
# Employer work - corporate compliance
profile=best-western
```

## Future Enhancements

Potential future additions to `.dotclaude` file:

```yaml
profile: blackwell-systems-oss

# Future: Auto-activation (opt-in)
auto_activate: false

# Future: Profile-specific overrides
overrides:
  git_workflow: disabled

# Future: Project metadata
project:
  name: "My Project"
  type: "oss"
```

Currently only `profile:` is supported and used.

---

**Back to:** [README.md](../README.md) | [USAGE.md](USAGE.md) | [ARCHITECTURE.md](ARCHITECTURE.md)
