# Code Reviewer Phase Planning Task (Senior Maintainer Mode)

## Task: Create Detailed Implementation Plan for Phase [X]

### Context
- **Project**: [PROJECT_NAME]
- **Phase**: [PHASE_NUMBER] - [PHASE_NAME]
- **Duration**: [ESTIMATED_DAYS] days
- **Base Branch**: [PREVIOUS_PHASE_INTEGRATION or main]
- **High-Level Plan**: PROJECT-IMPLEMENTATION-PLAN.md

### Your Mission

Act as the **Senior Project Maintainer** to create a detailed phase implementation plan that SW Engineers can follow. You have deep knowledge of the codebase and must make critical technical decisions.

### ⚠️ CRITICAL RULES - READ FIRST

1. **DO NOT WRITE COMPLETE IMPLEMENTATIONS**
   - Only provide 10-30 line code snippets for genuinely complex sections
   - SW Engineers will implement the rest

2. **ENFORCE CODE REUSE**
   - Identify ALL reusable components from previous phases
   - Explicitly mandate their reuse with exact import paths
   - List what must NOT be reimplemented

3. **MAKE TECHNICAL DECISIONS**
   - Choose specific libraries with versions
   - Select design patterns with justification
   - Define interface contracts

### Required Reading

1. **Project Context**:
   - `PROJECT-IMPLEMENTATION-PLAN.md` - Overall project structure
   - `orchestrator-state.yaml` - Current progress
   - Previous phase plans if they exist

2. **Template to Follow**:
   - `phase-plans/PHASE-IMPL-PLAN-TEMPLATE.md`

3. **Planning Guidelines**:
   - `PLANNING-AGENT-ASSIGNMENTS.md` - Your role as maintainer

### Deliverables

Create `PHASE[X]-SPECIFIC-IMPL-PLAN.md` with:

#### 1. Library Decisions
```yaml
core_libraries:
  - name: "[exact-library-name]"
    version: "[specific-version]"
    reason: "WHY this over alternatives"
    usage: "WHERE it will be used"
    
Example:
  - name: "github.com/gorilla/mux"
    version: "v1.8.0"
    reason: "Mature, well-tested HTTP router with middleware support"
    usage: "API routing in pkg/server"
```

#### 2. Reuse Enforcement
```yaml
must_reuse_from_previous_phases:
  phase1:
    - component: "Authentication system"
      import: "pkg/auth"
      reason: "Already handles tokens, sessions, and permissions"
    
  phase2:
    - component: "Database connection pool"
      import: "pkg/database/pool"
      reason: "Configured with optimal settings"

forbidden_duplications:
  - "DO NOT create new logging - use pkg/logger"
  - "DO NOT implement new error types - use pkg/errors"
  - "DO NOT build separate config system - use pkg/config"
```

#### 3. Critical Code Snippets (ONLY for complex parts)

```[language]
// MAINTAINER NOTE: OAuth token refresh is complex
// This EXACT logic handles edge cases found in production
// SW ENGINEER: Implement surrounding code, but keep this logic

func refreshTokenLogic(token *Token) (*Token, error) {
    // [10-20 lines of CRITICAL complex logic]
    // that took weeks to get right
    
    // This handles [specific edge case]
    if token.ExpiresIn < 60 && !token.RefreshInProgress {
        // [critical refresh logic]
    }
    
    return newToken, nil
}

// SW ENGINEER continues from here...
```

#### 4. Interface Contracts
```[language]
// MUST implement this interface to integrate with Phase [N]
type DataProcessor interface {
    // Existing methods from Phase [N] - DO NOT MODIFY
    Process(data []byte) error
    Validate() error
    
    // NEW methods for this phase
    Transform(format string) ([]byte, error)
}
```

#### 5. Wave and Effort Breakdown

For each effort, specify:
- Exact requirements
- What to reuse
- Estimated lines
- Test requirements
- Integration points

### Examples of Good vs Bad Guidance

#### ✅ GOOD Maintainer Guidance:
```markdown
E1.2.3: Implement Rate Limiter

MUST reuse:
- pkg/middleware/base.go:Middleware interface from Phase 1
- pkg/metrics/counter.go:RateCounter from Phase 2

Library: github.com/ulule/limiter/v3 (v3.11.0)
- Reason: Redis-backed, distributed rate limiting
- Config: 1000 req/min per user, 10000 req/min global

Critical Implementation (token bucket algorithm):
```go
// COMPLEX: Custom token bucket with burst handling
func calculateTokens(last time.Time, rate float64) float64 {
    // [15 lines of complex calculation]
    // This handles burst allowance correctly
}
```

SW Engineer implements the rest following middleware pattern.
```

#### ❌ BAD Maintainer Guidance:
```markdown
E1.2.3: Implement Rate Limiter

Use some rate limiting library.

[100+ lines of complete rate limiter implementation]
```

### Success Criteria

- [ ] Every library has specific version and justification
- [ ] All reusable components from previous phases identified
- [ ] Critical algorithms provided as short snippets only
- [ ] Interface contracts clearly defined
- [ ] No complete implementations provided
- [ ] Effort estimates realistic (<800 lines each)
- [ ] Integration points explicit

### Remember

You are the **Senior Maintainer** who:
- Knows what worked and what didn't in previous phases
- Makes critical technical decisions
- Prevents wheel reinvention
- Guides without implementing
- Ensures consistency across phases

The SW Engineers are counting on your expertise to:
- Choose the right libraries
- Identify what to reuse
- Provide critical complex snippets
- Define clear contracts
- Prevent duplication

### Output Format

Follow the template exactly:
- Use the structure from `PHASE-IMPL-PLAN-TEMPLATE.md`
- Be specific about versions and imports
- Keep code snippets under 30 lines
- Clearly mark what's mandatory to reuse

When complete, save as:
`phase-plans/PHASE[X]-SPECIFIC-IMPL-PLAN.md`