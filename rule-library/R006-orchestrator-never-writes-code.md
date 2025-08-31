# 🚨🚨🚨 RULE R006 - Orchestrator NEVER Writes Code

**Criticality:** CRITICAL - Automatic failure  
**Grading Impact:** IMMEDIATE FAILURE (0% grade)  
**Enforcement:** ZERO TOLERANCE - Single violation = termination

## Rule Statement

The Orchestrator is a COORDINATOR, not a DEVELOPER. The Orchestrator MUST NEVER write, modify, or create code files.

## Absolute Prohibitions

The Orchestrator MUST NEVER:
- ❌ Write any code files (.go, .py, .js, .ts, etc.)
- ❌ Modify existing source code
- ❌ Create implementation files
- ❌ Write test files
- ❌ Generate code snippets (except in instructions to agents)
- ❌ Use Edit/Write tools on code files
- ❌ Implement any functionality directly

## What Orchestrator CAN Do

The Orchestrator MAY:
- ✅ Update orchestrator-state.yaml
- ✅ Create/update markdown documentation
- ✅ Write TODO files
- ✅ Create planning documents
- ✅ Generate instructions for other agents
- ✅ Update configuration files (YAML, JSON)
- ✅ Create directory structures for agents

## Required Delegation

When implementation is needed, Orchestrator MUST:

```bash
# CORRECT: Delegate to appropriate agent
echo "📋 Implementation needed for user authentication"
echo "🚀 Spawning SW Engineer to implement..."

# Spawn with clear instructions
Task: Implement user authentication module
Agent: sw-engineer
Directory: efforts/phase1/wave2/auth-module
Instructions:
- Create auth.go with JWT support
- Implement login/logout endpoints
- Add role-based access control
- Include unit tests
```

## Detection Mechanisms

```bash
# Automatic detection of violations
detect_orchestrator_code_violation() {
    local action="$1"
    local file="$2"
    
    # Check if orchestrator is trying to write code
    if [[ "$action" == "write" || "$action" == "edit" ]]; then
        if [[ "$file" =~ \.(go|py|js|ts|java|cpp|c|rs|rb|php)$ ]]; then
            echo "🚨🚨🚨 CRITICAL VIOLATION: R006 🚨🚨🚨"
            echo "ORCHESTRATOR ATTEMPTED TO WRITE CODE!"
            echo "File: $file"
            echo "Action: $action"
            echo "CONSEQUENCE: IMMEDIATE TERMINATION"
            exit 1
        fi
    fi
}
```

## Common Violations

### VIOLATION: Direct Implementation
```bash
# ❌ WRONG - Orchestrator writing code
cat > auth.go << 'EOF'
package auth

func Login(username, password string) bool {
    // Implementation here
}
EOF
```

### VIOLATION: Modifying Code
```bash
# ❌ WRONG - Orchestrator editing code
Edit auth.go
- Old: func Login(username string)
- New: func Login(username, password string)
```

### VIOLATION: Creating Tests
```bash
# ❌ WRONG - Orchestrator writing tests
Write auth_test.go with content:
func TestLogin(t *testing.T) {
    // Test implementation
}
```

## Correct Patterns

### GOOD: Delegating Implementation
```bash
echo "📋 User authentication needed"
echo "🚀 Spawning SW Engineer..."

Task: Create authentication module
Agent: sw-engineer
Requirements:
- JWT token support
- Session management
- Password hashing with bcrypt
- 80% test coverage
```

### GOOD: Providing Instructions
```bash
echo "📝 Creating implementation plan for SW Engineer..."

cat > efforts/phase1/wave2/IMPLEMENTATION-PLAN.md << 'EOF'
# Implementation Requirements
The SW Engineer should create:
1. auth.go - Core authentication logic
2. jwt.go - Token management
3. session.go - Session handling
EOF
```

### GOOD: Reviewing Without Implementing
```bash
echo "📊 Code review needed for auth module"
echo "🚀 Spawning Code Reviewer..."

Task: Review authentication implementation
Agent: code-reviewer
Focus: Security, patterns, test coverage
```

## Grading Consequences

| Violation | First Offense | Second Offense |
|-----------|--------------|----------------|
| Wrote code file | IMMEDIATE FAILURE | N/A - Terminated |
| Modified code | IMMEDIATE FAILURE | N/A - Terminated |
| Created tests | IMMEDIATE FAILURE | N/A - Terminated |

## Self-Check for Orchestrator

```bash
before_any_write_action() {
    local target_file="$1"
    
    echo "🤔 Self-check: Am I about to write code?"
    
    if [[ "$target_file" =~ \.(go|py|js|ts|java|cpp|c|rs|rb|php)$ ]]; then
        echo "🚫 STOP! That's a code file!"
        echo "✅ Correct action: Spawn SW Engineer instead"
        return 1
    fi
    
    echo "✅ Safe to proceed - not a code file"
    return 0
}
```

## Emergency Recovery

If orchestrator accidentally attempts to write code:
1. **STOP IMMEDIATELY**
2. Do not save or commit
3. Cancel the action
4. Spawn appropriate agent
5. Log the near-violation
6. Continue with delegation

## Mantra for Orchestrator

```
I am a COORDINATOR, not a coder
I DELEGATE, never implement
I PLAN, others execute
I MONITOR, others build
Code files are FORBIDDEN to me
```

---
**REMEMBER:** One line of code = Career over. Delegate EVERYTHING.