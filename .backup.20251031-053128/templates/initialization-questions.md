# Software Factory 2.0 Initialization - Question Templates

## Overview
This document provides structured questions for the architect agent to use during the INIT_REQUIREMENTS_GATHERING state. Questions are organized by category and include example answers to guide users.

## Question Flow Strategy
1. Start with project type to determine path
2. Gather technical requirements
3. Collect quality and deployment needs
4. Confirm and summarize

---

## CATEGORY 1: Project Foundation

### Q1.1: Project Type
**Question**: "Is this project enhancing an existing codebase or creating something new?"

**Options**:
- `existing` - Adding features to an existing project
- `new` - Creating a new project from scratch

**Example Answers**:
- "existing" → Follow up with repository questions
- "new" → Follow up with project structure questions

### Q1.2: For Existing Codebases
**Question**: "What is the GitHub repository URL for the existing codebase?"

**Format**: `https://github.com/[org]/[repo]`

**Example Answers**:
- The URL configured in `$PROJECT-TARGET-REPO` (from setup-config.yaml)
- "https://github.com/your-org/your-project"
- "https://github.com/kubernetes/kubectl"

**Follow-up**: "Do you have a fork of this repository?"
- If yes: "What is your fork's URL?"
- If no: "Should I help you create one?"

### Q1.3: For New Projects
**Question**: "What type of application are you building?"

**Options**:
- `cli` - Command-line tool
- `api` - REST/GraphQL API service
- `library` - Shared library/package
- `webapp` - Web application
- `service` - Background service/daemon
- `plugin` - Plugin for another system

**Example Answers**:
- "cli" → Focus on command structure
- "api" → Focus on endpoints and data models
- "library" → Focus on public interfaces

---

## CATEGORY 2: Technology Stack

### Q2.1: Primary Language
**Question**: "What programming language will be the primary language for this project?"

**Common Options**:
- Go
- Python
- TypeScript/JavaScript
- Java
- Rust
- C++
- Ruby

**Example Answer with Context**:
- "Go" (if enhancing a Go project, this should match)
- "Python with type hints"
- "TypeScript (strict mode)"

### Q2.2: Key Frameworks
**Question**: "What are the main frameworks or libraries you'll use? (comma-separated)"

**Language-Specific Examples**:

**Go**:
- "cobra, viper, gin"
- "echo, gorm, testify"
- "chi, sqlx, zap"

**Python**:
- "FastAPI, SQLAlchemy, pydantic"
- "Django, celery, pytest"
- "Flask, requests, pandas"

**TypeScript**:
- "React, Next.js, Prisma"
- "Express, TypeORM, Jest"
- "NestJS, GraphQL, class-validator"

### Q2.3: Build System
**Question**: "What build system/tool will you use?"

**Language-Specific Defaults**:
- Go: `make`, `go build`, `goreleaser`
- Python: `poetry`, `pip`, `setuptools`
- JavaScript: `npm`, `yarn`, `pnpm`
- Java: `gradle`, `maven`
- Rust: `cargo`
- C++: `cmake`, `make`, `bazel`

### Q2.4: Testing Framework
**Question**: "What testing framework will you use?"

**Language-Specific Options**:
- Go: `go test`, `testify`, `ginkgo`
- Python: `pytest`, `unittest`, `nose2`
- JavaScript: `jest`, `mocha`, `vitest`
- Java: `junit`, `testng`

---

## CATEGORY 3: Architecture & Design

### Q3.1: Architecture Pattern
**Question**: "What architecture pattern best describes your project?"

**Options with Descriptions**:
- `monolith` - Single deployable unit
- `microservices` - Multiple small services
- `serverless` - Function-based, event-driven
- `plugin` - Extends existing system
- `library` - Reusable code package
- `modular-monolith` - Monolith with clear module boundaries

### Q3.2: Data Storage
**Question**: "Will this project need data persistence? If yes, what kind?"

**Options**:
- `none` - No persistence needed
- `file` - File-based storage
- `sqlite` - Embedded SQL database
- `postgres` - PostgreSQL database
- `mysql` - MySQL/MariaDB database
- `mongodb` - Document database
- `redis` - Cache/key-value store
- `cloud` - Cloud storage (S3, GCS, etc.)

### Q3.3: External Integration
**Question**: "What external systems or APIs will this integrate with?"

**Example Answers**:
- "Kubernetes API, Docker API"
- "GitHub API, Slack webhooks"
- "AWS S3, Lambda"
- "None - standalone application"

---

## CATEGORY 4: Deployment & Operations

### Q4.1: Deployment Target
**Question**: "Where will this be deployed?"

**Options**:
- `local` - Developer machines only
- `kubernetes` - Kubernetes cluster
- `docker` - Docker/Docker Compose
- `vm` - Virtual machines
- `serverless` - Lambda/Cloud Functions
- `binary` - Distributed as binary
- `package` - NPM, PyPI, etc.

### Q4.2: Container Strategy
**Question**: "Will you use containers? If yes, what's your strategy?"

**Options**:
- `none` - No containers
- `docker` - Standard Docker images
- `multi-stage` - Multi-stage Docker builds
- `distroless` - Minimal distroless images
- `scratch` - From scratch images

### Q4.3: CI/CD Platform
**Question**: "What CI/CD platform will you use?"

**Options**:
- `github-actions` - GitHub Actions
- `gitlab-ci` - GitLab CI
- `jenkins` - Jenkins
- `circleci` - CircleCI
- `none` - No CI/CD initially
- `custom` - Custom solution

---

## CATEGORY 5: Quality Requirements

### Q5.1: Code Coverage Target
**Question**: "What's your target code coverage percentage?"

**Common Answers**:
- "80" - Standard target
- "90" - High coverage
- "70" - Pragmatic coverage
- "100" - Full coverage

### Q5.2: Performance Requirements
**Question**: "Do you have specific performance requirements?"

**Example Answers**:
- "Sub-100ms API response time"
- "Handle 1000 requests/second"
- "Process 1GB files efficiently"
- "No specific requirements"

### Q5.3: Security Requirements
**Question**: "What security requirements do you have?"

**Example Answers**:
- "OAuth2/OIDC authentication"
- "TLS everywhere, secrets management"
- "OWASP Top 10 compliance"
- "Basic authentication sufficient"
- "No authentication needed"

---

## CATEGORY 6: Development Preferences

### Q6.1: Development Environment
**Question**: "What's your preferred development environment?"

**Options**:
- `vscode` - VS Code with extensions
- `intellij` - IntelliJ IDEA
- `vim/neovim` - Terminal-based
- `codespaces` - GitHub Codespaces
- `devcontainer` - VS Code Dev Containers

### Q6.2: Documentation Style
**Question**: "How will you document this project?"

**Options**:
- `markdown` - Markdown files in repo
- `godoc` - Go documentation
- `sphinx` - Sphinx (Python)
- `typedoc` - TypeDoc (TypeScript)
- `swagger` - OpenAPI/Swagger
- `minimal` - Just README

---

## SMART FOLLOW-UP PATTERNS

### For Unclear Answers
"I need a bit more detail. For [topic], could you specify [specific aspect]?"

### For Validation
"Just to confirm: You're building a [type] using [language] with [frameworks], correct?"

### For Missing Information
"One more thing I need for the configuration: [specific requirement]"

### For Options
"I can suggest [option1] or [option2] based on your requirements. Which would you prefer?"

---

## EXAMPLE CONVERSATION FLOW

```
Architect: "Is this project enhancing an existing codebase or creating something new?"
User: "existing"

Architect: "What is the GitHub repository URL for the existing codebase?"
User: "$PROJECT-TARGET-REPO" (This should be configured in setup-config.yaml)

Architect: "Do you have a fork of this repository?"
User: "Yes, https://github.com/myuser/idpbuilder"

Architect: "What programming language will be the primary language for this project?"
User: "Go"

Architect: "What are the main frameworks or libraries you'll use? (comma-separated)"
User: "cobra for CLI, go-containerregistry for OCI, gin for HTTP"

[continues through all categories...]

Architect: "Perfect! I have all the information I need. Let me summarize what I've gathered..."
```

---

## COMPLETION CHECKLIST

Before proceeding from requirements gathering, ensure:

### Required Information
- [ ] Project type (existing/new) determined
- [ ] Repository information collected (if applicable)
- [ ] Primary language specified
- [ ] Key frameworks listed
- [ ] Build system identified
- [ ] Test framework selected
- [ ] Architecture pattern chosen
- [ ] Deployment target defined

### Optional but Helpful
- [ ] Performance requirements noted
- [ ] Security requirements understood
- [ ] Documentation approach decided
- [ ] Development environment preferences

### Validation
- [ ] No placeholder answers
- [ ] Technology choices are compatible
- [ ] All required config fields can be populated