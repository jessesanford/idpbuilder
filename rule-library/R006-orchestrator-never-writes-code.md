# 🚨🚨🚨 RULE R006 - Orchestrator NEVER Writes, Measures, or Reviews Code

**Criticality:** BLOCKING - Automatic termination  
**Grading Impact:** IMMEDIATE FAILURE (0% grade)  
**Enforcement:** ZERO TOLERANCE - Single violation = termination
**Detection:** Automatic via file extension, operation monitoring, and tool usage

## Rule Statement

The Orchestrator is a COORDINATOR ONLY. The Orchestrator MUST NEVER write, modify, copy, move, manipulate, measure, or review ANY code files or implementation-related resources. This includes ALL technical assessments of code, not just writing.

## 🔴🔴🔴 ABSOLUTE PROHIBITIONS 🔴🔴🔴

The Orchestrator MUST NEVER:
- ❌ Write ANY code files (.go, .py, .js, .ts, .java, .cpp, .c, .rs, .rb, .php, etc.)
- ❌ Modify existing source code
- ❌ Create implementation files
- ❌ Write test files
- ❌ Generate code snippets (except in instructions to agents)
- ❌ Use Edit/Write/MultiEdit tools on code files
- ❌ Implement any functionality directly
- ❌ **COPY code files between directories** (even if "just infrastructure")
- ❌ **MOVE code files to new locations** (agents must do this)
- ❌ **CREATE code files as templates** (agents create all code)
- ❌ **SYMLINK or hardlink code files**
- ❌ **USE `cp`, `mv`, `ln` commands on code files**
- ❌ **MANIPULATE `.gitignore` to hide code operations**
- ❌ **MEASURE code size using line-counter.sh** (Code Reviewers do this)
- ❌ **ESTIMATE effort sizes or line counts** (Code Reviewers assess this)
- ❌ **DETERMINE if splits are needed** (Code Reviewers decide this)
- ❌ **REVIEW code quality or compliance** (Code Reviewers handle this)
- ❌ **MAKE technical assessments about code** (delegate ALL to specialists)
- ❌ **BACKPORT fixes to effort branches** (SW Engineers must do this)
- ❌ **CHERRY-PICK commits between branches** (SW Engineers handle this)
- ❌ **APPLY patches or fixes directly** (SW Engineers apply all fixes)

## ⚠️⚠️⚠️ CRITICAL CLARIFICATION ⚠️⚠️⚠️

**COPYING FILES IS NOT "CREATING INFRASTRUCTURE"**
- Creating directories = ✅ Infrastructure
- Copying code files = ❌ IMPLEMENTATION (FORBIDDEN)
- Moving code files = ❌ IMPLEMENTATION (FORBIDDEN)
- Creating empty folders = ✅ Infrastructure
- Populating folders with code = ❌ IMPLEMENTATION (FORBIDDEN)

**MEASURING CODE IS NOT "MONITORING"**
- Checking agent status = ✅ Monitoring
- Measuring code size = ❌ TECHNICAL ASSESSMENT (FORBIDDEN)
- Running line-counter.sh = ❌ CODE REVIEW WORK (FORBIDDEN)
- Estimating effort size = ❌ TECHNICAL ASSESSMENT (FORBIDDEN)
- Reading review reports = ✅ Monitoring

## What Orchestrator CAN Do

The Orchestrator MAY:
- ✅ Update orchestrator-state.json
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

When size measurement is needed, Orchestrator MUST:

```bash
# CORRECT: Delegate to Code Reviewer
echo "📏 Size measurement needed for effort"
echo "🚀 Spawning Code Reviewer to assess..."

# Spawn Code Reviewer for measurement
Task: Measure implementation size and compliance
Agent: code-reviewer
Directory: efforts/phase1/wave2/auth-module
Instructions:
- Use line-counter.sh to measure size
- Check if within 800 line limit
- Determine if splits are needed
- Create review report
```

## Detection Mechanisms

```bash
# Enhanced detection of ALL forbidden orchestrator actions
detect_orchestrator_violation() {
    local action="$1"
    local file="$2"
    
    # Expanded list of forbidden actions
    local forbidden_actions="write|edit|create|copy|cp|move|mv|link|ln"
    
    # Check if orchestrator is trying ANY operation on code
    if [[ "$action" =~ ^($forbidden_actions)$ ]]; then
        if [[ "$file" =~ \.(go|py|js|ts|java|cpp|c|rs|rb|php|jsx|tsx|vue|swift|kt|scala|lua|dart|r|m|h|hpp)$ ]]; then
            echo "🚨🚨🚨 CRITICAL VIOLATION: R006 🚨🚨🚨"
            echo "ORCHESTRATOR ATTEMPTED CODE OPERATION!"
            echo "File: $file"
            echo "Action: $action"
            echo "CONSEQUENCE: IMMEDIATE TERMINATION"
            echo "CORRECT: Spawn SW Engineer to handle this"
            exit 1
        fi
    fi
}

# Detection for measurement violations
detect_orchestrator_measurement_violation() {
    local command="$1"
    
    # Check if orchestrator is trying to measure code
    if [[ "$command" =~ line-counter\.sh ]]; then
        echo "🚨🚨🚨 CRITICAL VIOLATION: R006 🚨🚨🚨"
        echo "ORCHESTRATOR ATTEMPTED CODE MEASUREMENT!"
        echo "Command: $command"
        echo "CONSEQUENCE: IMMEDIATE TERMINATION"
        echo "CORRECT: Spawn Code Reviewer to measure"
        exit 1
    fi
    
    # Check for manual counting attempts
    if [[ "$command" =~ (wc[[:space:]]-l|find.*xargs.*wc|cloc|sloccount) ]]; then
        echo "🚨🚨🚨 CRITICAL VIOLATION: R006 🚨🚨🚨"
        echo "ORCHESTRATOR ATTEMPTED MANUAL LINE COUNTING!"
        echo "Command: $command"
        echo "CONSEQUENCE: IMMEDIATE TERMINATION"
        echo "CORRECT: Spawn Code Reviewer for size assessment"
        exit 1
    fi
}

# Monitor shell commands for violations
monitor_shell_commands() {
    local cmd="$1"
    
    # Detect cp/mv/ln operations on code files
    if [[ "$cmd" =~ (cp|mv|ln)[[:space:]].*\.(go|py|js|ts|java|cpp|c|rs|rb|php) ]]; then
        echo "🚨🚨🚨 SHELL VIOLATION DETECTED: R006 🚨🚨🚨"
        echo "Command attempted: $cmd"
        echo "BLOCKED: Orchestrator cannot manipulate code files"
        return 1
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

### 🔴 REAL VIOLATION FROM TRANSCRIPT: Copying Code Files
```bash
# ❌ CATASTROPHIC VIOLATION - Orchestrator copying implementation files
# This ACTUALLY HAPPENED and caused immediate failure:
cp api_test.go efforts/phase1/wave1/split2/api_test.go
cp models_test.go efforts/phase1/wave1/split2/models_test.go
cp middleware_test.go efforts/phase1/wave1/split2/middleware_test.go

# CONSEQUENCE: Immediate 0% grade, orchestrator terminated
# WHY: Copying code is IMPLEMENTATION work, not infrastructure
# CORRECT: SW Engineer agent should have been spawned to handle file operations
```

### 🔴 VIOLATION: Measuring Code Size
```bash
# ❌ CATASTROPHIC VIOLATION - Orchestrator measuring code
cd efforts/phase1/wave1/effort1
../../tools/line-counter.sh  # FORBIDDEN! Only Code Reviewers can do this

# CONSEQUENCE: Immediate 0% grade, orchestrator terminated
# WHY: Measuring code is TECHNICAL ASSESSMENT, not coordination
# CORRECT: Code Reviewer agent should be spawned to measure and assess
```

### VIOLATION: Moving Code Between Splits
```bash
# ❌ WRONG - Orchestrator moving code files
mv efforts/phase1/wave1/split1/*.go efforts/phase1/wave1/split2/

# Even for "infrastructure setup" - this is FORBIDDEN
```

### 🔴 VIOLATION: Backporting Fixes Directly
```bash
# ❌ CATASTROPHIC VIOLATION - Orchestrator backporting fixes
cd /efforts/effort1
git cherry-pick abc123  # FORBIDDEN! Only SW Engineers can do this

# ❌ WRONG - Orchestrator applying patches
git apply /efforts/integration-testing/fixes.patch

# ❌ WRONG - Orchestrator editing code for fixes  
Edit main.go to apply integration fixes

# CONSEQUENCE: Immediate 0% grade, orchestrator terminated
# WHY: Backporting is IMPLEMENTATION work requiring code changes
# CORRECT: SW Engineer agents must be spawned for EACH effort needing backports
```

### 🔴 VIOLATION: Making Excuses to Avoid Delegation
```bash
# ❌ INVALID EXCUSE: "Everything is already integrated, so I'll just apply the fixes"
# ❌ INVALID EXCUSE: "The backport process is simpler since it's all in one branch"  
# ❌ INVALID EXCUSE: "We only need to commit the fixes there"
# ❌ INVALID EXCUSE: "It would be faster if I just do it"

# REALITY: Even if all efforts are integrated, backporting is MANDATORY
# REALITY: Each effort branch MUST be independently mergeable
# REALITY: Skipping backports = ALL PRs WILL FAIL
# CORRECT: ALWAYS spawn SW Engineers for backporting, NO EXCEPTIONS
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
I SPAWN, never measure
I READ REPORTS, never assess
Code files are FORBIDDEN to me
Technical judgments are NOT mine to make
```

---
**REMEMBER:** One line of code = Career over. One measurement = Career over. Delegate EVERYTHING technical.