# 🏗️ ARCHITECTURE COMPLIANCE CHECKLIST

## 🔴🔴🔴 MANDATORY FOR ALL IMPLEMENTATIONS 🔴🔴🔴

This checklist MUST be completed:
1. **BEFORE** starting any implementation
2. **DURING** code reviews
3. **BEFORE** integration
4. **BEFORE** creating PRs

## 📋 PLANNING PHASE - LOCK IN ARCHITECTURE

### Architectural Decisions Document
```yaml
architectural_decisions:
  - library: "[name]"
    purpose: "[why chosen]"
    user_recommended: [YES/NO]
    approved: true
    locked: true  # CANNOT BE CHANGED
    change_requires: "EXPLICIT_USER_APPROVAL"
```

### Required Documentation
- [ ] All libraries documented
- [ ] All frameworks documented
- [ ] All patterns documented
- [ ] All technology choices documented
- [ ] User recommendations highlighted
- [ ] Approval status marked

## 🔍 IMPLEMENTATION PHASE - VERIFY COMPLIANCE

### Before Starting Implementation
```bash
# Run these checks BEFORE writing any code
cd $EFFORT_DIR

# 1. Check required libraries are available
echo "=== CHECKING REQUIRED LIBRARIES ==="
grep "go-containerregistry" go.mod || echo "⚠️ Missing required library!"

# 2. Verify no plans to replace
echo "=== VERIFYING IMPLEMENTATION PLAN ==="
# Review the effort plan for any deviations
```

### During Implementation
- [ ] Using ALL specified libraries
- [ ] Following approved patterns
- [ ] NO custom replacements
- [ ] NO technology substitutions
- [ ] Implementation matches plan EXACTLY

### Red Flags to Watch For
```go
// ❌ VIOLATION: Custom implementation
func customHTTPClient() { ... }

// ❌ VIOLATION: Replacing standard library
type MyORM struct { ... }  // When GORM was specified

// ❌ VIOLATION: Different pattern
// Plan said: Repository pattern
// Implementation: Direct database access
```

## 🔎 CODE REVIEW PHASE - VALIDATE ARCHITECTURE

### Code Reviewer Checklist
```markdown
## R362 ARCHITECTURAL COMPLIANCE

### Libraries & Frameworks
- [ ] All user-recommended libraries present
- [ ] No unauthorized library removals
- [ ] No custom replacements for standard libraries
- [ ] Dependencies match approved list

### Patterns & Design
- [ ] Implementation follows approved patterns
- [ ] No architectural deviations
- [ ] Technology stack unchanged
- [ ] Design decisions preserved

### Specific Checks
- [ ] go-containerregistry still used (if specified)
- [ ] No custom HTTP implementations for registry
- [ ] Framework usage as planned
- [ ] Database access as specified

**Violations Found**: [NONE / List here]
**Action Required**: [APPROVE / REJECT / FIX]
```

### Automated Validation Script
```bash
#!/bin/bash
# architecture-check.sh

echo "🏗️ R362 ARCHITECTURE COMPLIANCE CHECK"

VIOLATIONS=0

# Check for required libraries
REQUIRED_LIBS="go-containerregistry gorm gin"
for lib in $REQUIRED_LIBS; do
    if ! grep -q "$lib" go.mod; then
        echo "❌ MISSING REQUIRED: $lib"
        ((VIOLATIONS++))
    else
        echo "✅ Found: $lib"
    fi
done

# Check for forbidden patterns
FORBIDDEN="custom.*client|MyORM|DirectDB"
if grep -r "$FORBIDDEN" pkg/; then
    echo "❌ FORBIDDEN PATTERN DETECTED"
    ((VIOLATIONS++))
fi

if [ $VIOLATIONS -gt 0 ]; then
    echo "🔴 ARCHITECTURE VIOLATIONS: $VIOLATIONS"
    exit 362
else
    echo "✅ ARCHITECTURE COMPLIANT"
fi
```

## 🚧 INTEGRATION PHASE - REJECT NON-COMPLIANT CODE

### Integration Checklist
- [ ] Architecture unchanged from plan
- [ ] All libraries still present
- [ ] No custom implementations added
- [ ] Patterns remain consistent
- [ ] Technology stack preserved

### Integration Rejection Criteria
```bash
# Automatic rejection if:
git diff main...branch -- go.mod | grep "^-.*go-containerregistry"
# Result: REJECT - Required library removed

git diff main...branch -- pkg/ | grep "custom.*registry"
# Result: REJECT - Unauthorized custom implementation
```

## 🛑 VIOLATION RESPONSE PROTOCOL

### If Architecture Violation Detected:

1. **IMMEDIATE STOP**
   ```bash
   echo "🔴🔴🔴 R362 VIOLATION DETECTED"
   exit 362
   ```

2. **DOCUMENT VIOLATION**
   - What was changed
   - Why it differs from plan
   - Impact assessment

3. **ESCALATE TO USER**
   - Request explicit approval
   - Present alternatives
   - Wait for decision

4. **NO PROCEEDING WITHOUT APPROVAL**
   - Stop ALL work
   - No workarounds
   - No "temporary" solutions

## 📝 CHANGE REQUEST PROCESS

### If Architecture Change Is Needed:

1. **STOP Implementation**
2. **Document Request**
   ```markdown
   ## ARCHITECTURE CHANGE REQUEST

   **Current Architecture**: [what was approved]
   **Proposed Change**: [what you want to change]
   **Reason**: [why change is needed]
   **Impact**: [what will be affected]
   **Alternatives**: [other options considered]
   ```

3. **Get Explicit Approval**
   - User must approve IN WRITING
   - Update planning documents
   - Update this checklist

4. **Only Then Proceed**
   - Update all documentation
   - Implement approved change
   - Note approval in commits

## 🔴 PENALTIES FOR VIOLATIONS

### Severity Matrix
| Violation | Penalty | Recovery |
|-----------|---------|----------|
| Removed user library | -100% FAIL | Full rewrite required |
| Custom replacement | -100% FAIL | Revert and reimplement |
| Pattern deviation | -75% | Fix and resubmit |
| Minor change | -50% | Document and fix |

## ✅ SIGN-OFF REQUIREMENTS

### Every Implementation Must Include:
```markdown
## ARCHITECTURE COMPLIANCE CERTIFICATION

I certify that this implementation:
- [ ] Uses ALL specified libraries
- [ ] Follows approved patterns EXACTLY
- [ ] Makes NO unauthorized changes
- [ ] Matches the plan completely

Agent: [name]
Date: [date]
Effort: [effort-name]
```

---

**REMEMBER**: Architecture decisions are IMMUTABLE once approved. ANY change requires EXPLICIT user approval. Violations result in IMMEDIATE project failure.