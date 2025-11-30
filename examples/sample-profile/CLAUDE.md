# Profile: Sample Project

<!--
═══════════════════════════════════════════════════════════════════
IMPORTANT: This file is MERGED with base/CLAUDE.md when activated
═══════════════════════════════════════════════════════════════════

When you activate this profile, Claude Code sees:
  1. base/CLAUDE.md (universal standards: git, security, tools)
  2. THIS FILE (project-specific additions)

The base already includes:
  ✓ Git workflow and commit message standards
  ✓ Security practices (no secrets, input validation)
  ✓ Tool usage (Read instead of cat, Edit instead of sed)
  ✓ Task management (TodoWrite, marking todos complete)
  ✓ File operations (absolute paths, read before edit)

This profile ADDS project-specific context on top of that base.
-->

## Project Overview

This profile demonstrates configuration for a web application project with specific tech stack and practices.

**This is an overlay** - it adds project-specific details to the universal standards in base/CLAUDE.md.

## Tech Stack Preferences

### Backend
- **Language**: Node.js (TypeScript preferred)
- **Framework**: Express.js or Fastify
- **Database**: PostgreSQL with Prisma ORM
- **API Style**: RESTful with OpenAPI documentation

### Frontend
- **Framework**: React with TypeScript
- **State Management**: React Query + Context API
- **Styling**: Tailwind CSS
- **Build Tool**: Vite

### Testing
- **Unit Tests**: Vitest
- **Integration Tests**: Supertest
- **E2E Tests**: Playwright
- **Coverage Target**: 80%+ for critical paths

## Coding Standards

### TypeScript
- Always use explicit types, avoid `any`
- Enable strict mode in tsconfig.json
- Use interfaces for object shapes
- Use type aliases for unions and complex types

### Code Style
- Use ESLint + Prettier for formatting
- 2-space indentation
- Single quotes for strings
- Trailing commas in multiline structures

### Naming Conventions
- **Files**: kebab-case (e.g., `user-service.ts`)
- **Classes**: PascalCase (e.g., `UserService`)
- **Functions**: camelCase (e.g., `getUserById`)
- **Constants**: UPPER_SNAKE_CASE (e.g., `MAX_RETRY_ATTEMPTS`)

## API Design Principles

### RESTful Conventions
- Use HTTP methods appropriately (GET, POST, PUT, DELETE)
- Return appropriate status codes (200, 201, 400, 404, 500)
- Use plural nouns for resources (`/users`, not `/user`)
- Version APIs in the URL (`/api/v1/users`)

### Response Format
```json
{
  "success": true,
  "data": {},
  "error": null,
  "meta": {
    "timestamp": "2025-01-15T10:30:00Z",
    "requestId": "abc123"
  }
}
```

## Database Practices

### Migrations
- Always create migrations for schema changes
- Name migrations descriptively
- Never edit existing migrations in production
- Test migrations in staging first

### Queries
- Use Prisma for type-safe queries
- Avoid N+1 queries - use `include` or `select` wisely
- Add indexes for frequently queried fields
- Use transactions for multi-step operations

## Security Requirements

### Authentication
- Use JWT tokens with short expiration (15 minutes)
- Implement refresh token rotation
- Hash passwords with bcrypt (10+ rounds)
- Rate limit authentication endpoints

### Input Validation
- Validate all user input at API boundaries
- Use Zod or similar for schema validation
- Sanitize inputs to prevent XSS
- Use parameterized queries to prevent SQL injection

### Environment Variables
- Never commit `.env` files
- Use `.env.example` for documentation
- Validate required env vars on startup
- Use different credentials per environment

## Git Workflow

### Branches
- `main` - production code, always deployable
- `develop` - integration branch
- `feature/*` - new features
- `fix/*` - bug fixes
- `hotfix/*` - urgent production fixes

### Commit Messages
Follow conventional commits:
```
feat: add user registration endpoint
fix: resolve race condition in payment processing
docs: update API documentation for v2
refactor: extract validation middleware
test: add unit tests for user service
```

### Pull Requests
- Link to issue/ticket number
- Include test plan in description
- Request review from at least one teammate
- Ensure CI passes before merging
- Squash commits when merging to main

## Documentation Standards

### Code Comments
- Document "why", not "what"
- Add JSDoc comments for public functions
- Include examples for complex APIs
- Keep comments up-to-date with code changes

### README Files
- Include setup instructions
- Document environment variables
- Add troubleshooting section
- Include links to additional docs

### API Documentation
- Maintain OpenAPI/Swagger specs
- Include request/response examples
- Document error codes and meanings
- Keep docs in sync with implementation

## Testing Guidelines

### Unit Tests
- Test business logic in isolation
- Mock external dependencies
- Aim for fast execution (<1s per test)
- One assertion concept per test

### Integration Tests
- Test API endpoints end-to-end
- Use test database (not production)
- Clean up test data after each test
- Test error cases, not just happy paths

### Test Organization
```
src/
  users/
    user.service.ts
    user.service.test.ts
    user.controller.ts
    user.controller.test.ts
```

## Performance Considerations

### Caching Strategy
- Cache static content at CDN level
- Use Redis for session storage
- Implement cache invalidation strategy
- Set appropriate TTLs per data type

### Optimization
- Lazy load heavy dependencies
- Paginate large result sets (max 100 items)
- Use database indexes for slow queries
- Profile before optimizing

## Deployment

### Environments
- **Development**: Local machine
- **Staging**: Mirrors production, for testing
- **Production**: Live system

### Pre-Deployment Checklist
- [ ] All tests passing
- [ ] Database migrations tested
- [ ] Environment variables configured
- [ ] Monitoring and logging enabled
- [ ] Rollback plan documented

## Team Practices

### Code Review Focus
- Security vulnerabilities
- Logic errors and edge cases
- Code clarity and maintainability
- Test coverage
- Performance implications

### Communication
- Use GitHub Issues for bug reports
- Use GitHub Discussions for questions
- Keep PRs focused and reasonably sized
- Respond to review comments within 24 hours

---

<!--
═══════════════════════════════════════════════════════════════════
MERGE BEHAVIOR
═══════════════════════════════════════════════════════════════════

When this profile is activated, the final ~/.claude/CLAUDE.md contains:

  1. ALL of base/CLAUDE.md (universal standards)
  2. THEN this file appended below it
  3. A separator comment showing where profile starts

If there are conflicts, profile instructions take precedence.

The merged file is what Claude Code actually reads.
-->
