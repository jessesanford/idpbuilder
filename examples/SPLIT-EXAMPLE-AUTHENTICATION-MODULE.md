# Split Example: Authentication Module

## Scenario
An authentication module implementation exceeded 800 lines and needs to be split into manageable parts.

## Original Implementation (1,247 lines)

The original attempt tried to implement everything in one effort:
- User model and database schema
- Password hashing and validation
- JWT token generation and validation
- Session management
- OAuth2 integration
- Rate limiting
- Account recovery flows
- Admin user management

## Line Count Analysis

```bash
$ tools/line-counter.sh -c phase2/wave1/effort3-auth-module -d

===========================================
DETAILED BREAKDOWN BY DIRECTORY
===========================================

src/models/:
  user.js                           156 lines
  session.js                         89 lines
  oauth_provider.js                  67 lines
  Subtotal:                         312 lines

src/services/:
  auth_service.js                   287 lines
  token_service.js                  143 lines
  password_service.js                98 lines
  oauth_service.js                  176 lines
  Subtotal:                         704 lines

src/controllers/:
  auth_controller.js                156 lines
  admin_controller.js                75 lines
  Subtotal:                         231 lines

TOTAL:                           1,247 lines ❌ EXCEEDS 800 LIMIT
```

## Split Strategy

### Part 1: Core Authentication (378 lines)
**Focus**: Basic authentication flow
```
src/models/
  └── user.js                     156 lines
src/services/
  ├── password_service.js          98 lines
  └── auth_service.js (partial)   124 lines
```

### Part 2: Token & Session Management (232 lines)
**Focus**: JWT and session handling
```
src/models/
  └── session.js                   89 lines
src/services/
  └── token_service.js            143 lines
```

### Part 3: OAuth Integration (243 lines)
**Focus**: Third-party authentication
```
src/models/
  └── oauth_provider.js            67 lines
src/services/
  └── oauth_service.js            176 lines
```

### Part 4: Controllers & Admin (394 lines)
**Focus**: HTTP endpoints and admin features
```
src/controllers/
  ├── auth_controller.js          156 lines
  └── admin_controller.js          75 lines
src/services/
  └── auth_service.js (remaining) 163 lines
```

## Implementation Order

1. **Part 1 First** - Core authentication is foundation
2. **Part 2 Second** - Builds on user model
3. **Part 3 Third** - Extends authentication options
4. **Part 4 Last** - Ties everything together

## Branch Structure

```
phase2/wave1/effort3-auth-module-part1  (378 lines) ✅
phase2/wave1/effort3-auth-module-part2  (232 lines) ✅
phase2/wave1/effort3-auth-module-part3  (243 lines) ✅
phase2/wave1/effort3-auth-module-part4  (394 lines) ✅
```

## Key Principles Applied

### 1. Logical Grouping
Each part has a clear, focused purpose

### 2. Dependency Management
Parts build on each other in sequence

### 3. Independent Testing
Each part can be tested independently

### 4. Clean Interfaces
Clear boundaries between parts

### 5. Under Limit
All parts well under 800 lines

## Lessons Learned

### DO:
- ✅ Group related functionality together
- ✅ Keep interfaces minimal between parts
- ✅ Test each part independently
- ✅ Document integration points
- ✅ Plan splits during initial design

### DON'T:
- ❌ Split arbitrarily by file size
- ❌ Create circular dependencies
- ❌ Break logical units apart
- ❌ Exceed 700 lines (leave buffer)
- ❌ Mix unrelated features

## Alternative Split Strategies

### By Layer (Horizontal)
```
Part 1: All models (312 lines)
Part 2: Core services (385 lines)
Part 3: Extended services (319 lines)
Part 4: All controllers (231 lines)
```

### By Feature (Vertical)
```
Part 1: Basic auth (login/logout)
Part 2: OAuth integration
Part 3: Admin features
Part 4: Account recovery
```

### By Risk (Progressive)
```
Part 1: Minimal viable auth
Part 2: Security enhancements
Part 3: Extended features
Part 4: Nice-to-haves
```

## Template for Split Plan

```markdown
# Split Plan: [Effort Name]

## Current Status
- Total Lines: [X]
- Files: [count]
- Must split to <800 lines each

## Proposed Splits

### Part 1: [Descriptive Name] (~[Y] lines)
**Purpose**: [What this part does]
**Files**:
- [file1.ext] ([N] lines)
- [file2.ext] ([M] lines)
**Dependencies**: [What it needs]
**Provides**: [What others can use]

### Part 2: [Descriptive Name] (~[Z] lines)
[Same structure...]

## Integration Points
- Part 1 → Part 2: [How they connect]
- Part 2 → Part 3: [Interface description]

## Testing Strategy
- Part 1: [Test approach]
- Part 2: [Test approach]
- Integration: [How to test together]

## Risk Assessment
- [Potential issues]
- [Mitigation strategies]
```

## Remember

1. **Plan splits during design** - Don't wait until after implementation
2. **Keep logical units together** - Don't break cohesive functionality
3. **Test each part** - Ensure parts work independently
4. **Document interfaces** - Clear contracts between parts
5. **Stay under 700 lines** - Leave buffer for reviewer requests