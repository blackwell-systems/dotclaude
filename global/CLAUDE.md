# Global Claude Code Instructions

These instructions apply to ALL projects unless overridden by project-specific CLAUDE.md.

## Development Standards

### Code Quality
- Write clean, maintainable code
- Follow existing patterns and conventions in each project
- Use meaningful variable and function names
- Keep functions focused and single-purpose

### File Operations
- Always use absolute file paths
- Read files before editing to understand context
- Prefer editing existing files over creating new ones
- Only create documentation when explicitly requested

### Security
- Never commit sensitive data (API keys, credentials, tokens)
- Use environment variables for configuration
- Validate user input at system boundaries
- Follow OWASP guidelines for web applications

### Git Practices
- Write clear, descriptive commit messages
- Focus commit messages on "why" rather than "what"
- Review changes before committing
- Never force push to main/master without explicit approval

### Long-Lived Feature Branches
When working repeatedly on the same feature branch:

**After PR is merged:**
1. `git checkout main && git pull`
2. `git checkout feature-branch`
3. Run `sync-feature-branch` (choose rebase or merge)
4. Continue development on the synced branch

**Automated reminders:**
- SessionStart hook checks if current branch is behind main
- PostToolUse hook reminds after git operations
- Use `check-branches` to see all branch statuses
- Use `pr-merged` for guided post-merge workflow

**Why this matters:**
- Prevents merge conflicts from accumulating
- Keeps feature branches deployable
- Avoids confusion about which code is current
- Enables continuous iteration on the same branch

## Tool Usage

### Prefer Specialized Tools
- Use Read instead of `cat`
- Use Edit instead of `sed`
- Use Write instead of `echo >` or heredoc
- Use Grep instead of bash `grep`
- Use Glob instead of `find`

### Task Management
- Use TodoWrite for complex multi-step tasks
- Mark todos in_progress before starting work
- Mark todos completed immediately when done
- Keep only ONE task in_progress at a time

### Background Operations
- Run long-running commands with run_in_background
- Monitor output with BashOutput
- Clean up with KillShell when done

## Project Context

### Before Making Changes
1. Read relevant files to understand current implementation
2. Search for similar patterns in the codebase
3. Understand the project's architecture and conventions
4. Ask clarifying questions if requirements are unclear

### When Implementing Features
- Start with the simplest solution that works
- Avoid over-engineering
- Don't add features beyond what was requested
- Test changes in the project's environment

## Communication

### With Users
- Be concise and direct
- Avoid unnecessary emojis unless requested
- Ask questions when requirements are unclear
- Provide file:line references for code locations

### In Code
- Add comments only where logic isn't self-evident
- Write self-documenting code when possible
- Document public APIs and complex algorithms
- Keep comments up-to-date with code changes

## Common Patterns

### Error Handling
- Validate at system boundaries (user input, external APIs)
- Trust internal code and framework guarantees
- Don't add error handling for impossible scenarios
- Fail fast with clear error messages

### Performance
- Optimize only when necessary
- Profile before optimizing
- Use appropriate data structures
- Consider time/space tradeoffs

## Notes

These are general guidelines. Project-specific requirements always take precedence.
