# Architect - INIT_REQUIREMENTS_GATHERING State Rules

## Purpose
Conduct interactive Q&A session to gather all project requirements.

## Entry Criteria
- Examples and templates loaded
- Question framework prepared
- User available for interaction

## Required Actions

### 1. Begin with Context
"Based on your project idea: [idea], I need to gather some requirements to set up your Software Factory properly."

### 2. Ask Questions by Category

#### Category 1: Target Codebase
```
Q1: Is this for an existing codebase or a new project?
    If existing:
    Q1a: What is the repository name/URL?
    Q1b: Do you have a fork to use as the upstream target?
    Q1c: What is your fork's URL?
    Q1d: What is the main branch name?

    If new:
    Q1a: What should the repository be named?
    Q1b: Should I create a git repository for it?
```

#### Category 2: Technology Stack
```
Q2: What programming language(s) will this use?
Q3: What are the key libraries/frameworks? (comma-separated)
Q4: What build system? (make, npm, cargo, gradle, etc.)
Q5: What testing framework(s)?
Q6: Any code generation tools?
```

#### Category 3: Architecture & Patterns
```
Q7: What type of project? (CLI tool, library, web service, etc.)
Q8: Architecture pattern? (monolith, microservices, serverless, etc.)
Q9: Any specific design patterns to follow?
Q10: Any existing code patterns to maintain compatibility with?
```

#### Category 4: Deployment & Infrastructure
```
Q11: Target deployment environment? (cloud provider, k8s, local, etc.)
Q12: Container requirements? (Docker, specific base images)
Q13: CI/CD preferences? (GitHub Actions, Jenkins, etc.)
Q14: Any infrastructure as code? (Terraform, Helm, etc.)
```

#### Category 5: Quality Requirements
```
Q15: Performance requirements? (response time, throughput)
Q16: Security requirements? (auth method, encryption)
Q17: Compliance needs? (HIPAA, SOC2, etc.)
Q18: Code coverage target? (percentage)
```

#### Category 6: Project Specifics
```
Q19: Who are the end users?
Q20: What problem does this solve?
Q21: Any competitor products to reference?
Q22: Documentation requirements?
```

### 3. Process Answers
For each answer:
- Validate format/content
- Ask follow-ups if unclear
- Store in structured format
- Mark field as complete

### 4. Build Requirements Document
Create structured requirements in state file:
```json
"requirements": {
  "codebase": {
    "type": "existing|new",
    "upstream_url": "...",
    "fork_url": "...",
    "main_branch": "..."
  },
  "technology": {
    "languages": ["..."],
    "frameworks": ["..."],
    "build_system": "...",
    "test_framework": "..."
  },
  "architecture": {
    "type": "...",
    "patterns": ["..."]
  },
  "deployment": {
    "environment": "...",
    "containerization": "...",
    "ci_cd": "..."
  },
  "quality": {
    "performance": "...",
    "security": "...",
    "coverage_target": "..."
  }
}
```

### 5. Validate Completeness
Ensure all required fields have values:
- For setup-config.yaml: 12 required fields
- For target-repo-config.yaml: 5-8 fields
- For IMPLEMENTATION-PLAN.md: All sections

## Interactive Q&A Best Practices

### Asking Questions
- One category at a time
- Provide examples in questions
- Accept "none" or "N/A" as valid
- Offer common options

### Example Interactions
```
Architect: "What programming language(s) will this project use?"
User: "Go"

Architect: "What are the key libraries/frameworks? (e.g., gin, cobra, gorm)"
User: "cobra for CLI, go-containerregistry for OCI"

Architect: "What build system? (make, go build, bazel, etc.)"
User: "make"
```

### Handling Unclear Answers
If answer is ambiguous:
"Could you clarify X? For example, do you mean Y or Z?"

If answer missing key info:
"That's helpful. Could you also specify [missing detail]?"

## Exit Criteria
- All required fields have values
- Requirements document complete
- Ready for repository decision
- User confirmed answers

## Transition
**MANDATORY**: → INIT_REPO_DECISION (return to orchestrator)

## Time Guidance
- ~15-20 questions total
- Allow 30-60 seconds per answer
- Total time: 10-15 minutes

## Error Handling
- User unsure → Provide defaults/examples
- Contradictory answers → Clarify priority
- Missing critical info → Mark required, re-ask