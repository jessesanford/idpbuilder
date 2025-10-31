# EXAMPLE: BAD Effort Plan with Vague Scope

## Effort Metadata
- **Effort**: E1.1.2 - Authentication Service
- **Branch**: `phase1/wave1/effort-auth`
- **Size Estimate**: 600 lines
- **Created**: 2025-01-21 10:00:00

## ❌ VAGUE SCOPE DEFINITION (CAUSES 3-5X OVERRUNS)

### Implementation Requirements
"Implement a complete authentication service with user management"

### Components to Build
- Authentication endpoints
- User management functionality
- Security features
- Database integration
- Testing

## ❌ WHY THIS IS BAD

### Problem 1: No Specific Function List
**Vague**: "Authentication endpoints"
**Result**: SW Engineer implements login, logout, refresh, 2FA, password reset, email verification, OAuth, SAML
**Actual**: 2,400 lines instead of 600

### Problem 2: No DO NOT List
**Missing**: No boundaries specified
**Result**: SW Engineer adds caching, rate limiting, session management, audit logging
**Actual**: Additional 800 lines of "helpful" features

### Problem 3: Open-Ended Requirements
**Vague**: "Security features"
**Result**: Engineer implements encryption, JWT, CORS, CSRF, XSS protection, SQL injection prevention
**Actual**: 1,200 lines of security code

### Problem 4: No Explicit Counts
**Missing**: How many endpoints? How many functions? How many types?
**Result**: Engineer keeps adding until it "feels complete"
**Actual**: 47 functions instead of planned 5

## ACTUAL OUTCOME FROM PRODUCTION

```
Plan said: "Implement authentication service"
Engineer did: Full enterprise auth system
Result: 4,847 lines (8X overrun)
Consequence: Effort failed, required emergency split into 6 parts
Grade: -100% AUTOMATIC FAILURE
```

## Common Vague Phrases That Cause Problems

| Vague Phrase | Engineer Interpretation | Actual Need |
|--------------|------------------------|-------------|
| "Complete implementation" | Everything possible | 3 specific functions |
| "Handle errors properly" | Comprehensive error system | Basic nil checks |
| "Add tests" | 100% coverage with edge cases | 3 happy path tests |
| "User management" | Full CRUD + search + filter | Just Create and Read |
| "Security features" | Enterprise security suite | Basic auth header check |

## The Cascading Failure Pattern

1. **Vague Plan**: "Implement auth service"
2. **Engineer Assumption**: "Must be production-ready"
3. **Scope Creep**: Adds 15 features not requested
4. **Size Explosion**: 600 → 4,847 lines
5. **Review Failure**: Code reviewer rejects
6. **Emergency Split**: Days of rework
7. **Schedule Impact**: 3-day effort becomes 2 weeks

## How This Could Have Been Prevented

See EFFORT-PLAN-EXAMPLE-GOOD.md for the correct way to specify scope with:
- EXACT function names and counts
- Explicit DO NOT IMPLEMENT list
- Specific line estimates per component
- Clear boundaries and limits

**REMEMBER**: Vague instructions are not flexibility - they're a recipe for failure!