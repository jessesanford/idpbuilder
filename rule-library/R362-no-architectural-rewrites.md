# 🔴🔴🔴 SUPREME LAW: RULE R362 - ABSOLUTE PROHIBITION OF ARCHITECTURAL REWRITES 🔴🔴🔴

## Rule Type: BLOCKING | SUPREME LAW | ARCHITECTURAL INTEGRITY

## Overview
This rule ABSOLUTELY FORBIDS any changes to approved architectural decisions, user-recommended libraries, or established patterns without EXPLICIT written approval from the user during the planning phase.

## 🚨🚨🚨 CRITICAL VIOLATIONS (IMMEDIATE PROJECT FAILURE)

### ABSOLUTELY FORBIDDEN ACTIONS:
1. **Removing User-Recommended Libraries** - NEVER remove libraries the user explicitly recommended
2. **Replacing Established Patterns** - NEVER replace approved architectural patterns
3. **Rewriting Core Architecture** - NEVER rewrite fundamental architectural decisions
4. **Changing Technology Stack** - NEVER substitute technologies without approval
5. **Deviating from Approved Plans** - NEVER implement differently than planned

### SPECIFIC EXAMPLE VIOLATIONS:
- ❌ Replacing go-containerregistry with custom HTTP implementation
- ❌ Removing approved ORM for raw SQL
- ❌ Replacing approved framework with custom implementation
- ❌ Changing from microservices to monolith (or vice versa)
- ❌ Substituting approved messaging systems

## 🔴 ENFORCEMENT MECHANISMS

### 1. PLANNING PHASE LOCK-IN
```yaml
architectural_decisions:
  - decision: "Use go-containerregistry for registry operations"
    approved: true
    locked: true
    change_requires: "EXPLICIT_USER_APPROVAL"
  - decision: "Use GORM for database operations"
    approved: true
    locked: true
    change_requires: "EXPLICIT_USER_APPROVAL"
```

### 2. IMPLEMENTATION PHASE VALIDATION
Every Software Engineer MUST:
```bash
# Before starting implementation
validate_architecture() {
    echo "🔍 Validating architectural compliance..."

    # Check for required libraries
    if grep -q "go-containerregistry" go.mod; then
        echo "✅ Required library: go-containerregistry present"
    else
        echo "🚨 VIOLATION: Required library go-containerregistry missing!"
        exit 1
    fi

    # Verify no unauthorized replacements
    if grep -q "net/http.*registry" *.go; then
        echo "🚨 VIOLATION: Unauthorized HTTP implementation for registry!"
        exit 1
    fi
}
```

### 3. CODE REVIEW ARCHITECTURAL CHECKLIST
Code Reviewer MUST verify:
```markdown
## 🏗️ ARCHITECTURAL COMPLIANCE CHECKLIST
- [ ] All user-recommended libraries still in use
- [ ] No unauthorized pattern replacements
- [ ] Implementation matches approved plan EXACTLY
- [ ] No custom implementations replacing standard libraries
- [ ] No technology stack deviations

**VIOLATIONS FOUND**: [List any architectural violations]
**ACTION**: REJECT if any violations found
```

### 4. INTEGRATION PHASE REJECTION
Orchestrator MUST:
```bash
# Before integration
check_architectural_integrity() {
    local branch="$1"

    # Compare against approved architecture
    git diff main..."$branch" -- go.mod | grep "^-.*go-containerregistry" && {
        echo "🚨🚨🚨 CRITICAL: Approved library removed!"
        echo "INTEGRATION REJECTED: Architectural violation"
        return 1
    }
}
```

## 📋 REQUIRED ARCHITECTURE DOCUMENTATION

### Every Effort MUST Include:
```markdown
## ARCHITECTURAL DECISIONS
Library/Framework: [name]
Purpose: [why chosen]
User Approved: [YES/NO]
Can Be Changed: NO (requires explicit approval)
```

## 🚨 VIOLATION PENALTIES

### Severity Levels:
1. **Removing User-Recommended Library**: IMMEDIATE FAILURE (-100%)
2. **Replacing Core Architecture**: IMMEDIATE FAILURE (-100%)
3. **Unauthorized Pattern Change**: CRITICAL FAILURE (-75%)
4. **Technology Stack Deviation**: CRITICAL FAILURE (-75%)
5. **Minor Architectural Change**: MAJOR PENALTY (-50%)

## 🔍 DETECTION MECHANISMS

### Automated Checks:
```bash
# Run at every phase transition
architecture_compliance_check() {
    # Check go.mod for required libraries
    REQUIRED_LIBS="go-containerregistry gorm gin"
    for lib in $REQUIRED_LIBS; do
        grep -q "$lib" go.mod || {
            echo "🚨 MISSING REQUIRED: $lib"
            exit 1
        }
    done

    # Check for unauthorized implementations
    FORBIDDEN_PATTERNS="net/http.*registry custom.*ORM raw.*sql"
    for pattern in $FORBIDDEN_PATTERNS; do
        grep -r "$pattern" pkg/ && {
            echo "🚨 FORBIDDEN PATTERN: $pattern"
            exit 1
        }
    done
}
```

## 🛡️ PREVENTION STRATEGIES

### 1. Lock Architecture in Planning
- Document EVERY architectural decision
- Get explicit approval for each
- Mark as IMMUTABLE

### 2. Validate Before Implementation
- Check approved libraries are available
- Verify patterns are understood
- Confirm no alternatives will be used

### 3. Monitor During Development
- Regular architecture compliance checks
- Immediate stop if deviation detected
- Escalate to user for any changes

### 4. Reject Non-Compliant Code
- No merge if architecture violated
- No exceptions without user approval
- Document all rejections

## 📝 CHANGE REQUEST PROCESS

If architecture change is genuinely needed:
1. STOP all implementation
2. Document why change is needed
3. Present alternatives to user
4. Get EXPLICIT written approval
5. Update all planning documents
6. Only then proceed with change

## 🔴 INTEGRATION WITH OTHER RULES

This rule supersedes and enhances:
- R220/R221: Size limits (don't change to avoid limits)
- R251: Planning compliance
- R304: Line counting standards
- R309: Integration standards

## EXAMPLES OF PROPER COMPLIANCE

### ✅ CORRECT: Following Approved Architecture
```go
// Using approved library as planned
import "github.com/google/go-containerregistry/pkg/v1/remote"

func PushImage(ref string, img v1.Image) error {
    return remote.Write(ref, img, remote.WithAuthFromKeychain(keychain))
}
```

### ❌ WRONG: Custom Implementation
```go
// VIOLATION: Replacing approved library with custom HTTP
func PushImage(ref string, img []byte) error {
    req, _ := http.NewRequest("POST", registryURL, bytes.NewReader(img))
    // Custom implementation...
}
```

## MANDATORY ACKNOWLEDGMENT

All agents MUST acknowledge:
```
I understand Rule R362 absolutely forbids:
- Removing user-recommended libraries
- Changing approved architecture
- Replacing established patterns
- Any deviation without explicit approval
Violation = IMMEDIATE PROJECT FAILURE
```

---
**ENFORCEMENT**: Automatic via state transitions
**VERIFICATION**: Architecture compliance checks
**PENALTY**: IMMEDIATE FAILURE for any violation