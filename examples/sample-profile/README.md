# Sample Profile Example

This directory contains a complete example profile showing how to customize Claude Code for a specific project or context.

## What's Included

### CLAUDE.md
The main instruction file that gets merged with your base CLAUDE.md when this profile is active.

**What to customize:**
- **Project Overview**: Describe your project's purpose and context
- **Tech Stack**: List your specific frameworks, libraries, and tools
- **Coding Standards**: Define your team's style guide and conventions
- **API Design**: Document your REST/GraphQL patterns
- **Security Requirements**: Specify your auth, validation, and compliance needs
- **Git Workflow**: Describe your branching strategy and commit conventions
- **Testing Guidelines**: Set expectations for test coverage and practices
- **Deployment**: Document your environments and deployment process

**When Claude Code runs with this profile active**, it will have all these guidelines in context and follow them automatically.

### settings.json
Profile-specific Claude Code settings that override your global settings.

**Common settings to customize:**

```json
{
  "provider": "anthropic-claude-max",  // or "aws-bedrock"
  "model": "sonnet",                   // or "opus", "haiku"

  "hooks": {
    // Run shell commands at specific events
    "sessionStart": "echo 'Starting work on Project X'",
    "userPromptSubmit": "",
    "preToolUse": "",
    "postToolUse": "git status"
  },

  "workingDirectories": [
    // Additional directories Claude can access
    "/path/to/project",
    "/path/to/related/repo"
  ],

  "editor": "code",  // "vim", "nvim", "cursor", etc.

  "preferences": {
    "autoApproveTools": ["Read", "Grep", "Glob"],
    "confirmBeforeGitPush": true,
    "showToolApprovalDetails": true
  }
}
```

## How to Use This Example

### Option 1: Copy and Customize
```bash
# Copy to create your own profile
cp -r examples/sample-profile profiles/my-project

# Edit the files
vim profiles/my-project/CLAUDE.md
vim profiles/my-project/settings.json

# Activate it
activate-profile my-project
```

### Option 2: Use as Reference
Keep this example as documentation and create your profiles from scratch using these files as a guide.

## Profile Structure

```
profiles/my-project/
├── CLAUDE.md          # Required: Instructions for Claude
├── settings.json      # Optional: Claude Code settings
└── README.md          # Optional: Notes for yourself
```

## Real-World Profile Examples

### For an Open Source Project
```markdown
# Profile: My OSS Library

## Project Overview
A TypeScript library for data validation used by 10k+ developers.

## Communication Style
- Friendly and welcoming to contributors
- Detailed explanations in code comments
- Public API must be thoroughly documented

## Code Standards
- 100% TypeScript strict mode
- Comprehensive unit tests (>95% coverage)
- Semantic versioning for releases
- Changelog must be updated with every PR

## Documentation
- JSDoc for all public APIs
- README examples must be tested
- TypeScript types are the source of truth
```

### For Client Work
```markdown
# Profile: Acme Corp Website

## Project Overview
E-commerce site for Acme Corp. Billing project #12345.

## Client Preferences
- jQuery (legacy codebase, no React)
- IE11 compatibility required
- WCAG 2.1 AA compliance mandatory
- No external dependencies without approval

## Code Standards
- Follow existing jQuery patterns
- Prefix all IDs with 'acme-'
- Use client's CSS framework (Bootstrap 3)
- Comment extensively for future maintainers
```

### For Personal Projects
```markdown
# Profile: My Side Project

## Project Overview
Experimenting with new tech and building for fun.

## Preferences
- Move fast, be pragmatic
- Try latest language features
- Ok to refactor extensively
- Focus on learning over polish

## Tech Stack
- Whatever seems interesting today
- Experiment with bleeding-edge features
- No tests unless needed
```

## Tips for Creating Profiles

1. **Start Simple**: Begin with just tech stack and coding style preferences
2. **Evolve Over Time**: Add more detail as you discover what helps
3. **Be Specific**: "Use TypeScript strict mode" is better than "Write good code"
4. **Include Examples**: Show desired patterns, not just rules
5. **Keep It Current**: Update profiles as your practices evolve

## When to Create a New Profile

- Working on a different project with different tech/standards
- Switching between client work and personal projects
- Different compliance requirements (HIPAA, SOC2, etc.)
- Different communication styles (internal vs public)
- Different languages or frameworks

## Benefits of Profile-Based Configuration

- **Context Switching**: Claude adapts to each project automatically
- **Team Alignment**: Share profiles with teammates for consistency
- **Compliance**: Enforce security/legal requirements per project
- **Learning**: Help Claude learn your project's specific patterns
- **Efficiency**: No need to repeat preferences each session

---

**Need Help?**
- [Full Documentation](../../docs/README.md)
- [Usage Guide](../../docs/USAGE.md)
- [Architecture Details](../../docs/ARCHITECTURE.md)
