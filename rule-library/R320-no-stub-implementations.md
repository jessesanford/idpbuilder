# 🚨🚨🚨 BLOCKING RULE R320: No Stub Implementations 🚨🚨🚨

## 🔴🔴🔴 CRITICAL BLOCKER: STUB IMPLEMENTATIONS ARE FORBIDDEN 🔴🔴🔴

**THIS RULE IS A CRITICAL BLOCKER - VIOLATION = IMMEDIATE REVIEW FAILURE**

## Rule Definition

**ANY CODE THAT DOES NOT PROVIDE ACTUAL FUNCTIONALITY IS A STUB AND MUST BE REJECTED**

### What Constitutes a Stub Implementation:

1. **Explicit "Not Implemented" Returns**
   ```go
   return fmt.Errorf("not implemented")
   return errors.New("not yet implemented")
   return nil, errors.New("TODO")
   ```

2. **Panic Placeholders**
   ```go
   panic("unimplemented")
   panic("TODO")
   panic("not implemented")
   ```

3. **Empty Function Bodies (in non-interface contexts)**
   ```go
   func DoSomething() error {
       return nil  // With no actual logic
   }
   ```

4. **TODO Comments Without Implementation**
   ```go
   func ProcessData() {
       // TODO: implement this
   }
   ```

5. **Placeholder Returns**
   ```python
   raise NotImplementedError
   return "TODO"
   pass  # in non-abstract methods
   ```

6. **JavaScript/TypeScript Stubs**
   ```javascript
   throw new Error("Not implemented")
   return Promise.reject("TODO")
   console.warn("Not implemented")
   ```

## 🚨🚨🚨 SEVERITY CLASSIFICATION 🚨🚨🚨

**ALL STUB IMPLEMENTATIONS ARE CRITICAL BLOCKERS - NO EXCEPTIONS**

- **Core Functionality Stubs**: CRITICAL BLOCKER
- **Helper Function Stubs**: CRITICAL BLOCKER  
- **Command Handler Stubs**: CRITICAL BLOCKER
- **API Endpoint Stubs**: CRITICAL BLOCKER
- **ANY Stub**: CRITICAL BLOCKER

## Code Reviewer Responsibilities

### MANDATORY STUB DETECTION PROTOCOL:

1. **Search for Common Stub Patterns**
   ```bash
   # Go patterns
   grep -r "not.*implemented\|TODO\|unimplemented" --include="*.go"
   grep -r "panic.*TODO\|panic.*unimplemented" --include="*.go"
   
   # Python patterns
   grep -r "NotImplementedError\|pass.*#.*TODO" --include="*.py"
   
   # JavaScript/TypeScript patterns
   grep -r "Not implemented\|TODO.*throw\|Promise.reject" --include="*.js" --include="*.ts"
   ```

2. **Verify Core Functionality**
   - Don't just check if code compiles
   - Don't just verify structure exists
   - VERIFY ACTUAL IMPLEMENTATION EXISTS
   - Test that functions DO SOMETHING MEANINGFUL

3. **Classification Requirements**
   - Stub found = CRITICAL BLOCKER (never "minor issue")
   - Missing functionality = FAILED REVIEW
   - "Not implemented" = IMMEDIATE REJECTION

## 🔴🔴🔴 GRADING PENALTIES 🔴🔴🔴

### For Code Reviewers:
- **-50%**: Passing ANY stub implementation through review
- **-30%**: Classifying stub as "minor issue" instead of CRITICAL
- **-40%**: Marking stub code as "properly implemented"
- **-100%**: Multiple stub implementations passed in single review

### For SW Engineers:
- **-20%**: Submitting code with stub implementations
- **-30%**: Not fixing stubs after review feedback
- **-50%**: Arguing that stubs are acceptable

## Contradictory Assessment Prevention

**THESE COMBINATIONS ARE FORBIDDEN:**
- ❌ "✅ properly implemented" + "returns not implemented"
- ❌ "Minor issue" + "core functionality missing"
- ❌ "Implementation looks good" + "TODO: implement this"
- ❌ "Code structure correct" + "panic(unimplemented)"

**If code contains stubs, the ONLY valid assessment is:**
- "❌ CRITICAL BLOCKER: Stub implementation detected"
- "❌ FAILED REVIEW: Core functionality not implemented"

## Integration with Other Rules

- **Works with R007**: Size limits don't excuse stubs
- **Works with R031**: Code review MUST catch stubs
- **Works with R220**: Atomic PRs need complete functionality
- **Overrides structure checks**: Good structure + stub = FAIL

## Examples of Review Failures

### FAILURE EXAMPLE 1: Go Command Handler
```go
// ❌ THIS MUST FAIL REVIEW
func (c *Controller) HandlePush(args []string) error {
    // TODO: implement push logic
    return fmt.Errorf("push command not yet implemented")
}
```
**Review Result**: CRITICAL BLOCKER - Core command returns stub error

### FAILURE EXAMPLE 2: Python API Endpoint
```python
# ❌ THIS MUST FAIL REVIEW
def process_payment(amount, card_info):
    """Process a payment transaction"""
    raise NotImplementedError("Payment processing coming soon")
```
**Review Result**: CRITICAL BLOCKER - API endpoint not implemented

### FAILURE EXAMPLE 3: JavaScript Service
```javascript
// ❌ THIS MUST FAIL REVIEW
async function syncDatabase() {
    console.warn("Database sync not implemented");
    return Promise.resolve({ status: "TODO" });
}
```
**Review Result**: CRITICAL BLOCKER - Service method is stub

## Acceptable Patterns (NOT Stubs)

These patterns are acceptable and NOT considered stubs:

1. **Interface Definitions** (no implementation expected)
2. **Abstract Base Classes** (meant to be overridden)
3. **Deprecated Functions** with clear deprecation notices
4. **Feature Flags** that conditionally disable features
5. **Proper Error Handling** for unsupported operations

## Enforcement Checklist for Code Reviewers

**MUST CHECK EVERY REVIEW:**
- [ ] Search for "not implemented" patterns
- [ ] Search for "TODO" in function bodies
- [ ] Search for panic/throw with stub messages
- [ ] Verify each function has actual logic
- [ ] Confirm no placeholder returns
- [ ] Validate core functionality works
- [ ] Check for empty function bodies
- [ ] Ensure no contradictory assessments

## Recovery Protocol When Stubs Found

1. **IMMEDIATE**: Mark review as FAILED
2. **CLASSIFY**: All stubs as CRITICAL BLOCKERS
3. **DOCUMENT**: List every stub found with location
4. **REQUIRE**: Complete implementation before re-review
5. **VERIFY**: Re-review confirms ALL stubs replaced

## Special Note on Size Limit Distractions

**SIZE LIMITS DO NOT EXCUSE STUBS!**

Even if a PR is at 799 lines (just under limit), stub implementations are still CRITICAL BLOCKERS. The correct action is:
1. FAIL the review for stubs
2. Fix the stubs (may exceed size limit)
3. If over limit, create split plan
4. Each split must have complete functionality

## Summary

**ZERO TOLERANCE FOR STUB IMPLEMENTATIONS**

- Any "not implemented" = FAILED REVIEW
- Any TODO in code = CRITICAL BLOCKER
- Any placeholder = IMMEDIATE REJECTION
- Structure without function = UNACCEPTABLE
- Good architecture with stubs = STILL FAILS

This rule exists because stub implementations:
- Break production systems
- Create false completion signals
- Hide actual work remaining
- Mislead stakeholders
- Violate atomic PR principles

**Remember**: The Software Factory produces WORKING SOFTWARE, not templates!