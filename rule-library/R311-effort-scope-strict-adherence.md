# 🚨🚨🚨 BLOCKING: RULE R311 - Effort Scope Strict Adherence Protocol 🚨🚨🚨

**Category:** Implementation Control  
**Agents:** SW Engineer (PRIMARY), Code Reviewer  
**Criticality:** BLOCKING - Over-engineering causes 3-5X size violations  
**Penalty:** -100% for exceeding effort scope, -50% for unrequested features

## 🔴🔴🔴 SUPREME DIRECTIVE: IMPLEMENT EXACTLY WHAT'S SPECIFIED 🔴🔴🔴

**NO MORE, NO LESS - EVEN IF IT SEEMS INCOMPLETE!**

## The Critical Problem This Solves

**THE MINDSET PROBLEM:**
Engineers want to "complete features" but efforts need "partial implementations"

**ACTUAL VIOLATION FROM PRODUCTION:**
```
Effort Plan Said: "Implement authentication service"
SW Engineer Did: Added OAuth2, JWT, SAML, session management, 2FA, password reset
Result: 2,847 lines instead of 600 (4.7X OVERRUN)
Consequence: COMPLETE FAILURE - Effort rejected, required emergency split
```

**RECENT SPLIT-003 VIOLATION FROM TRANSCRIPT:**
```
Plan Said: "selected command files and utilities"
SW Thought: "I should implement all important commands to make it work"
Result: 2,215 lines instead of 650 (3.4X OVERRUN)
Root Cause: "Complete the feature" mindset instead of "stay within budget"
```

## 🚨 MANDATORY ADHERENCE PROTOCOL

### 0. ADOPT THE CORRECT MINDSET

**❌ WRONG MINDSET (Causes overruns):**
- "This feature needs X, Y, and Z to work properly"
- "I should make this production-ready"
- "While I'm here, I'll add this related thing"
- "This seems incomplete without authentication"

**✅ CORRECT MINDSET (Stays within budget):**
- "I will implement ONLY what's explicitly listed"
- "Other efforts will handle the missing pieces"
- "Incomplete is intentional and expected"
- "My job is to stay within budget, not complete the feature"

### 1. READ THE "DO NOT IMPLEMENT" SECTION FIRST

**Before writing ANY code:**
```bash
# MANDATORY: Extract and display what NOT to do
echo "═══════════════════════════════════════════════════════"
echo "🛑 DO NOT IMPLEMENT LIST FROM EFFORT PLAN:"
echo "═══════════════════════════════════════════════════════"
grep -A 20 "DO NOT IMPLEMENT\|SCOPE BOUNDARIES\|FORBIDDEN" EFFORT-IMPLEMENTATION-PLAN.md

# Acknowledge the boundaries
echo "✅ I will NOT implement anything in the DO NOT list"
echo "✅ I will STOP at the specified boundaries"
echo "✅ I will RESIST adding 'helpful' extras"
```

### 2. EXTRACT EXACT SPECIFICATIONS

**Parse the effort plan for explicit items:**
```bash
# Extract exact specifications
echo "📊 EXACT SCOPE FOR THIS EFFORT:"
grep -A 10 "IMPLEMENT EXACTLY\|EXACTLY What to Implement" EFFORT-IMPLEMENTATION-PLAN.md

# Count specific items
FUNCTION_COUNT=$(grep -c "^  - Function:" EFFORT-IMPLEMENTATION-PLAN.md)
ENDPOINT_COUNT=$(grep -c "^  - Endpoint:" EFFORT-IMPLEMENTATION-PLAN.md)
TYPE_COUNT=$(grep -c "^  - Type:" EFFORT-IMPLEMENTATION-PLAN.md)

echo "  Functions to implement: ${FUNCTION_COUNT:-0}"
echo "  Endpoints to implement: ${ENDPOINT_COUNT:-0}"
echo "  Types to define: ${TYPE_COUNT:-0}"
echo "  ANYTHING ELSE: FORBIDDEN WITHOUT JUSTIFICATION"
```

### 3. IMPLEMENT ONLY WHAT'S EXPLICITLY NAMED

#### ❌ WRONG - Adding Unlisted Features:
```go
// Effort plan says: "Implement basic authentication with login endpoint"
// SW Engineer writes:
func HandleLogin(w http.ResponseWriter, r *http.Request) { ... }  // ✅ Correct
func HandleLogout(w http.ResponseWriter, r *http.Request) { ... } // ❌ NOT IN PLAN!
func HandleRefreshToken(w http.ResponseWriter, r *http.Request) { ... } // ❌ NOT IN PLAN!
func Handle2FA(w http.ResponseWriter, r *http.Request) { ... }    // ❌ NOT IN PLAN!
func HandlePasswordReset(w http.ResponseWriter, r *http.Request) { ... } // ❌ NOT IN PLAN!
// Result: 5 endpoints instead of 1 = VIOLATION
```

#### ✅ RIGHT - Exactly What's Listed:
```go
// Effort plan says: "Implement basic authentication with login endpoint"
// SW Engineer writes:
func HandleLogin(w http.ResponseWriter, r *http.Request) { ... }  // ✅ Listed
// STOP HERE - Even if logout seems essential!
```

### 4. JUSTIFICATION PROTOCOL FOR NECESSARY ADDITIONS

**If you TRULY need something not in the plan:**

```bash
# STOP AND THINK
echo "⚠️ ADDITION CONSIDERATION:"
echo "Need to add: [describe addition]"
echo "Reason: [why it's ESSENTIAL, not just nice]"
echo "Impact if omitted: [what BREAKS without it]"
echo "Line count: [estimated lines]"

# Document in implementation report
cat >> IMPLEMENTATION-REPORT.md << EOF
## Justified Additions
### Addition: [Name]
- **What**: [Description]
- **Why Essential**: [Not just "nice to have"]
- **What Breaks Without It**: [Specific failure]
- **Lines Added**: [Count]
- **Approved By**: [Self-review justification]
EOF
```

### 5. VALIDATION CHECKPOINTS

#### After Each Major Component:
```bash
# Validate scope adherence
validate_effort_scope() {
    echo "🔍 SCOPE VALIDATION CHECKPOINT"
    
    # Count actual implementations
    ACTUAL_FUNCS=$(grep -c "^func " *.go 2>/dev/null || echo 0)
    ACTUAL_ENDPOINTS=$(grep -c "HandleFunc\|Handler(" *.go 2>/dev/null || echo 0)
    
    echo "📈 Current implementation:"
    echo "  Functions: $ACTUAL_FUNCS"
    echo "  Endpoints: $ACTUAL_ENDPOINTS"
    
    # Check for scope creep
    if grep -q "OAuth\|SAML\|JWT" *.go && ! grep -q "OAuth\|SAML\|JWT" EFFORT-IMPLEMENTATION-PLAN.md; then
        echo "❌ VIOLATION: Adding authentication methods not in plan!"
        exit 1
    fi
}
```

## 🛑 WHEN THE PLAN SEEMS INSUFFICIENT

### ASK BEFORE ASSUMING

**If the effort plan seems to miss something critical:**

```markdown
## SCOPE CLARIFICATION REQUEST

The effort plan specifies:
- Login endpoint only
- No logout functionality

This seems incomplete because:
- Users who log in need a way to log out
- Sessions will accumulate without cleanup

Options:
A) Implement ONLY login as specified ✅ (RECOMMENDED)
B) Add logout with justification in report
C) Request plan clarification

Proceeding with Option A - implementing ONLY what's specified.
```

## 📏 REALISTIC SIZE GUIDELINES

### Typical Component Sizes:

| Component | Simple | Medium | Complex |
|-----------|--------|--------|---------|
| REST Endpoint | 30-50 lines | 75-150 lines | 200-300 lines |
| Service Function | 20-40 lines | 50-100 lines | 150-250 lines |
| Data Model | 10-30 lines | 40-80 lines | 100-200 lines |
| Unit Test | 15-30 lines | 40-70 lines | 100-150 lines |
| Integration Test | 30-60 lines | 80-150 lines | 200-400 lines |

### Over-Engineering Red Flags:

| Plan Says | SW Does | Why It's Wrong |
|-----------|---------|----------------|
| "Basic auth" | Adds OAuth2, SAML | 10X complexity increase |
| "CRUD endpoints" | Adds search, filter, sort | 3X size increase |
| "Simple service" | Adds caching, retry logic | 5X complexity |
| "Data model" | Adds validation, serialization | 4X size increase |

## 🔴 ENFORCEMENT MECHANISMS

### Pre-Implementation Gate:
```bash
# Must extract and acknowledge scope
if ! test -f .scope-acknowledgment; then
    echo "❌ Must acknowledge effort scope first!"
    echo "Run: validate_effort_scope > .scope-acknowledgment"
    exit 1
fi
```

### During Implementation:
```bash
# Monitor line growth
LINES_ADDED=$(git diff --stat | tail -1 | awk '{print $4}')
EFFORT_LIMIT=800

if [ "$LINES_ADDED" -gt 500 ]; then
    echo "⚠️ WARNING: Approaching effort limit ($LINES_ADDED/800)"
    echo "Focus on completing specified features only!"
fi
```

### Pre-Commit Gate:
```bash
# Final scope validation
check_effort_compliance() {
    # Check for unauthorized additions
    JUSTIFICATION_COUNT=$(grep -c "## Justified Additions" IMPLEMENTATION-REPORT.md)
    ADDITION_COUNT=$(grep -c "Added.*not in plan" git diff)
    
    if [ "$ADDITION_COUNT" -gt "$JUSTIFICATION_COUNT" ]; then
        echo "❌ BLOCKING: Unjustified additions detected!"
        echo "Either remove additions or document justification!"
        exit 1
    fi
}
```

## 🎯 Success Criteria

An effort implementation is successful when:

1. ✅ Implements EXACTLY the specified components
2. ✅ Stays under 800 lines (hard limit)
3. ✅ Contains NO unrequested features without justification
4. ✅ Follows all "DO NOT" instructions
5. ✅ Documents any essential additions with clear reasoning
6. ✅ Stops at specified boundaries

## 📊 Grading Impact

- Exceeding scope without justification: **-100% (AUTOMATIC FAILURE)**
- Adding "nice to have" features: **-50% per feature**
- Missing justification for additions: **-40%**
- Ignoring "DO NOT" instructions: **-75%**
- Over-engineering simple requirements: **-60%**

## 💡 Key Mental Models

### The Restaurant Order Analogy
- Customer orders: "Hamburger with lettuce and tomato"
- ❌ WRONG: Add cheese, bacon, onions, special sauce (they seem good!)
- ✅ RIGHT: Make exactly hamburger with lettuce and tomato
- If something seems missing: ASK, don't assume

### The Building Blueprint Analogy
- Blueprint shows: 3 rooms, 2 doors, 4 windows
- ❌ WRONG: Add a bathroom because "every house needs one"
- ✅ RIGHT: Build exactly 3 rooms, 2 doors, 4 windows
- Missing bathroom? That's Phase 2's problem

## 🔑 Mantras to Remember

- **"If it's not specified, don't build it"**
- **"Incomplete today is better than oversized"**
- **"The plan is the contract, not a suggestion"**
- **"Adding features without asking is sabotage"**
- **"Document or delete additions"**
- **"Stay within budget, don't complete the feature"**
- **"Partial implementation is the goal"**
- **"Other efforts will handle what's missing"**

## Examples of Proper Adherence

### Example 1: Authentication Service
```
Plan: Implement login endpoint with username/password validation
Did: Created one endpoint, validated credentials, returned token
Did NOT: Add logout, refresh, password reset, 2FA
Size: 285 lines
Result: SUCCESS ✅
```

### Example 2: Data Service
```
Plan: Implement Create and Read operations for User model
Did: Implemented CreateUser and GetUser functions only
Did NOT: Add Update, Delete, List, Search operations
Size: 420 lines
Result: SUCCESS ✅
```

### Example 3: API Gateway
```
Plan: Implement request routing to 3 backend services
Did: Created router with exactly 3 route mappings
Did NOT: Add authentication, rate limiting, caching, logging
Size: 195 lines
Result: SUCCESS ✅
```

## Integration with R310 (Split Scope)

This rule works in harmony with R310:
- **R310**: Controls scope in split implementations
- **R311**: Controls scope in regular effort implementations
- **Same Principle**: Build only what's specified
- **Same Enforcement**: Strict boundaries and documentation

## Summary

**The #1 cause of effort failures is the "complete the feature" mindset when "partial implementation" is required.**

**KEY INSIGHT FROM TRANSCRIPT:**
Engineers aren't trying to break the system - they're trying to be helpful by making features complete. But the Software Factory REQUIRES partial implementations that fit within size budgets.

This rule ensures:
- SW Engineers adopt "stay within budget" mindset
- Partial implementations are recognized as intentional
- EXACTLY what's specified gets implemented
- Any additions require explicit justification
- Clear boundaries prevent scope creep

**Remember: Your job is NOT to complete the feature. Your job is to implement EXACTLY what's specified within the size budget. Other efforts will handle the rest.**