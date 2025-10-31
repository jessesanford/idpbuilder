# EXAMPLE: GOOD Effort Plan with Explicit Boundaries

## Effort Metadata
- **Effort**: E1.1.2 - Basic Authentication
- **Branch**: `phase1/wave1/effort-auth`
- **Size Estimate**: 385 lines (well under 800)
- **Created**: 2025-01-21 10:00:00

## ✅ EXPLICIT SCOPE DEFINITION (R311 COMPLIANT)

### IMPLEMENT EXACTLY (BE SPECIFIC!)

#### Functions (EXACTLY 3 - NO MORE)
```go
1. HandleLogin(w http.ResponseWriter, r *http.Request)     // ~80 lines - Basic username/password validation
2. ValidateToken(token string) (*Claims, error)            // ~45 lines - Check JWT signature only
3. GenerateToken(userID string) (string, error)            // ~40 lines - Create simple JWT
// STOP HERE - DO NOT ADD MORE ENDPOINTS OR FUNCTIONS
```

#### Types/Structs (EXACTLY 2)
```go
// Type 1: Credentials for login
type LoginRequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

// Type 2: JWT claims
type Claims struct {
    UserID string `json:"user_id"`
    jwt.StandardClaims
}
// NO additional types, NO methods on these structs
```

#### API Endpoints (EXACTLY 1)
```
POST /api/v1/auth/login - Login with username/password (~80 lines)
// NO logout endpoint
// NO refresh endpoint  
// NO password reset
// NO user registration
```

### 🛑 DO NOT IMPLEMENT (SCOPE BOUNDARIES)

**EXPLICITLY FORBIDDEN IN THIS EFFORT:**
- ❌ DO NOT implement logout (future effort E1.1.3)
- ❌ DO NOT implement token refresh (future effort E1.1.4)
- ❌ DO NOT implement password reset (phase 2)
- ❌ DO NOT implement user registration (effort E1.1.1 handles this)
- ❌ DO NOT add 2FA/MFA support (phase 3)
- ❌ DO NOT implement OAuth/SAML (not in requirements)
- ❌ DO NOT add rate limiting (effort E1.2.5)
- ❌ DO NOT implement session management (using stateless JWT)
- ❌ DO NOT add comprehensive error messages (basic only)
- ❌ DO NOT create middleware (future effort)

### 📊 REALISTIC SIZE CALCULATION

```
Production Code:
- HandleLogin function:      80 lines
- ValidateToken function:    45 lines
- GenerateToken function:    40 lines
- Type definitions:          20 lines
- Imports and setup:         15 lines
Subtotal:                   200 lines

Test Code:
- TestHandleLogin:          60 lines (2 cases: success, failure)
- TestValidateToken:        50 lines (2 cases: valid, invalid)
- TestGenerateToken:        30 lines (1 case: success)
- Test helpers:             25 lines
Subtotal:                   165 lines

Supporting Files:
- Config constants:         20 lines
Subtotal:                   20 lines

TOTAL: 385 lines (48% of limit - safe buffer)
```

## 🔴 MANDATORY ADHERENCE CHECKPOINTS

### Before Starting:
```bash
echo "EFFORT SCOPE LOCKED:"
echo "✓ Functions: EXACTLY 3 (HandleLogin, ValidateToken, GenerateToken)"
echo "✓ Types: EXACTLY 2 (LoginRequest, Claims)"
echo "✓ Endpoints: EXACTLY 1 (POST /api/v1/auth/login)"
echo "✓ Tests: EXACTLY 5 basic tests"
echo "✗ NO logout functionality"
echo "✗ NO token refresh"
echo "✗ NO user registration"
```

### During Implementation:
```bash
# After implementing each function
FUNC_COUNT=$(grep -c "^func Handle\|^func Validate\|^func Generate" *.go)
if [ "$FUNC_COUNT" -gt 3 ]; then
    echo "⚠️ SCOPE VIOLATION: Too many functions!"
    exit 1
fi
```

## Files to Create

### auth_handler.go (~165 lines MAX)
```go
package auth

// ONLY these imports
import (
    "encoding/json"
    "net/http"
    "time"
    "github.com/dgrijalva/jwt-go"
)

// ONLY these 2 types
type LoginRequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

type Claims struct {
    UserID string `json:"user_id"`
    jwt.StandardClaims
}

// ONLY these 3 functions
func HandleLogin(w http.ResponseWriter, r *http.Request) { ... }
func ValidateToken(token string) (*Claims, error) { ... }
func GenerateToken(userID string) (string, error) { ... }
```

### auth_handler_test.go (~165 lines MAX)
```go
// ONLY these test functions
func TestHandleLogin_Success(t *testing.T) { ... }     // ~30 lines
func TestHandleLogin_InvalidCreds(t *testing.T) { ... } // ~30 lines
func TestValidateToken_Valid(t *testing.T) { ... }      // ~25 lines
func TestValidateToken_Invalid(t *testing.T) { ... }    // ~25 lines
func TestGenerateToken_Success(t *testing.T) { ... }    // ~30 lines
// NO edge cases, NO benchmark tests, NO integration tests
```

## Success Criteria
- [ ] Implemented EXACTLY 3 functions (no more, no less)
- [ ] Created EXACTLY 2 types (no more, no less)
- [ ] Added EXACTLY 1 endpoint (no more, no less)
- [ ] Wrote EXACTLY 5 tests (no more, no less)
- [ ] Total lines under 400
- [ ] NO logout functionality added
- [ ] NO token refresh added
- [ ] NO additional "helpful" features

## If Something Seems Missing

**DON'T ASSUME - DOCUMENT!**

If during implementation you think "this needs logout to be useful":
1. STOP - Don't implement logout
2. Document in IMPLEMENTATION-REPORT.md:
   ```markdown
   ## Note: Logout Not Implemented
   - Per effort plan, logout is effort E1.1.3
   - Current implementation uses stateless JWT
   - Client can discard token for "logout" effect
   ```
3. Continue with ONLY what's specified

## Why This Plan Succeeds

1. **Explicit Functions**: Names each function with line estimates
2. **Clear Boundaries**: Lists what NOT to implement
3. **Exact Counts**: "EXACTLY 3 functions" leaves no ambiguity
4. **Realistic Sizes**: Based on actual code, not optimistic guesses
5. **Checkpoints**: Automated verification during implementation
6. **Buffer Space**: 385/800 lines leaves room for discoveries

## Comparison with Bad Example

| Aspect | Bad Plan | Good Plan | Result |
|--------|----------|-----------|---------|
| Scope | "Auth service" | "3 specific functions" | No ambiguity |
| Boundaries | None | 10 explicit exclusions | No scope creep |
| Size | "~600 lines" | Detailed breakdown = 385 | Accurate |
| Testing | "Add tests" | "EXACTLY 5 tests listed" | Clear target |
| Outcome | 4,847 lines | 385 lines | PROJECT_DONE |

**REMEMBER**: Specific boundaries aren't restrictions - they're protection against failure!