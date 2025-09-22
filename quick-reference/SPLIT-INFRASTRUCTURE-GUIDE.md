# Split Infrastructure Quick Reference Guide

## 🔴🔴🔴 CRITICAL: Each Split Needs Its Own Directory and Clone 🔴🔴🔴

### The Golden Rule
**ONE SPLIT = ONE DIRECTORY = ONE CLONE = ONE BRANCH**

## Directory Structure

### Before Splitting
```
efforts/phase1/wave1/
└── auth-system/                    # Original effort (will be too large)
    ├── .git/                       # Git repository
    ├── pkg/                        # Implementation
    └── IMPLEMENTATION-PLAN.md      # Original plan
```

### After Splitting
```
efforts/phase1/wave1/
├── auth-system/                    # DEPRECATED (too large)
│   ├── SPLIT-INVENTORY.md          # List of all splits
│   ├── SPLIT-PLAN-auth-system-split001-20250120-143000.md  # Plan for split 1 (timestamped per R301)
│   ├── SPLIT-PLAN-auth-system-split002-20250120-143000.md  # Plan for split 2 (timestamped per R301)
│   └── SPLIT-PLAN-auth-system-split003-20250120-143000.md  # Plan for split 3 (timestamped per R301)
├── auth-system-SPLIT-001/          # Split 1 (NEW DIRECTORY)
│   ├── .git/                       # SEPARATE CLONE
│   ├── SPLIT-PLAN-auth-system-split001-20250120-143000.md  # Copied timestamped plan
│   └── pkg/                        # Implementation
├── auth-system-SPLIT-002/          # Split 2 (NEW DIRECTORY)
│   ├── .git/                       # SEPARATE CLONE
│   ├── SPLIT-PLAN-auth-system-split002-20250120-143000.md  # Copied timestamped plan
│   └── pkg/                        # Implementation
└── auth-system-SPLIT-003/          # Split 3 (NEW DIRECTORY)
    ├── .git/                       # SEPARATE CLONE
    ├── SPLIT-PLAN-auth-system-split003-20250120-143000.md  # Copied timestamped plan
    └── pkg/                        # Implementation
```

## Directory Naming Convention

### Format
```
${EFFORT_NAME}-SPLIT-${SPLIT_NUMBER}
```

### Rules
- **EFFORT_NAME**: Original effort name (e.g., `auth-system`)
- **-SPLIT-**: Literal text "SPLIT" in UPPERCASE with hyphens
- **SPLIT_NUMBER**: Zero-padded 3 digits (001, 002, 003)

### Examples
✅ CORRECT:
- `auth-system-SPLIT-001`
- `user-management-SPLIT-002`
- `api-gateway-SPLIT-003`

❌ WRONG:
- `auth-system-split-1` (lowercase, no padding)
- `auth-system-001` (missing SPLIT keyword)
- `auth-system--split-001` (double hyphen)

## Sequential Branching Strategy

### Branch Creation Order
```
1. Split-001: Based on phase-integration (same as original)
   └── 2. Split-002: Based on Split-001 (NOT phase-integration!)
       └── 3. Split-003: Based on Split-002 (NOT Split-001!)
```

### Why Sequential?
- Each split measures ONLY its own additions
- No cumulative line counting
- Clean integration without conflicts
- Later splits can use earlier split code

## Who Does What?

### Code Reviewer
1. Creates SPLIT-INVENTORY.md in too-large branch
2. Creates timestamped SPLIT-PLAN-{effort}-split{num}-{timestamp}.md files (per R301)
3. Commits and pushes to too-large branch

### Orchestrator
1. Reads split plans from too-large branch
2. Creates split directories with -SPLIT-XXX suffix
3. Clones target repo into each split directory
4. Creates branches sequentially
5. Copies relevant timestamped split plan to each directory
6. Spawns SW Engineer

### SW Engineer
1. Verifies split infrastructure exists
2. Confirms in correct -SPLIT-XXX directory
3. Finds and reads timestamped split plan for their split
4. Stays under 800 lines
5. Notifies orchestrator when all splits complete

## Verification Commands

### Check Split Infrastructure
```bash
# Verify split structure for an effort
$CLAUDE_PROJECT_DIR/utilities/verify-split-infrastructure.sh auth-system

# List all splits in system
$CLAUDE_PROJECT_DIR/utilities/verify-split-infrastructure.sh --list

# Check sequential branching
$CLAUDE_PROJECT_DIR/utilities/verify-split-infrastructure.sh --branching auth-system
```

### Manual Checks
```bash
# Are we in a split directory?
pwd | grep -q "SPLIT-" && echo "YES" || echo "NO"

# What split number are we?
pwd | grep -o 'SPLIT-[0-9]*' | grep -o '[0-9]*'

# Is this a separate git repository?
[[ -d .git ]] && echo "YES" || echo "NO"

# What branch are we on?
git branch --show-current

# Does our split plan exist? (Check for timestamped format)
SPLIT_NUM=$(pwd | grep -o 'SPLIT-[0-9]*' | grep -o '[0-9]*')
SPLIT_PLAN=$(ls -t SPLIT-PLAN-*-split${SPLIT_NUM}-*.md 2>/dev/null | head -1)
[[ -n "$SPLIT_PLAN" ]] && echo "YES: $SPLIT_PLAN" || echo "NO"
```

## Common Problems and Solutions

### Problem: "Not in a split directory"
**Solution**: Orchestrator must create split directories first with -SPLIT-XXX suffix

### Problem: "No split plan found"
**Solution**: Code Reviewer must create plans first, Orchestrator must copy them

### Problem: "Not a git repository"
**Solution**: Each split needs its own clone - orchestrator must clone into each directory

### Problem: "Wrong branch"
**Solution**: Branch name must contain split number (e.g., auth-system-SPLIT-001)

### Problem: Line count includes previous splits
**Solution**: Use sequential branching - each split based on previous, not all from same base

## Integration After Splits

### Merge Order (Sequential)
```bash
# After all splits complete:
1. Merge split-003 → split-002
2. Merge split-002 → split-001  
3. Merge split-001 → phase-integration
4. Mark original as deprecated
```

### Why This Order Works
- No conflicts between splits
- Each builds on previous
- Clean final result
- Original can be abandoned

## Key Rules

- **R204**: Orchestrator creates split infrastructure
- **R199**: Single code reviewer creates all split plans
- **R202**: Single SW engineer implements all splits
- **R205**: Splits implemented sequentially, not parallel
- **R296**: Mark original as deprecated after splits complete

## Quick Checklist

Before starting split implementation:
- [ ] Split directory exists with -SPLIT-XXX suffix
- [ ] Directory contains separate git clone
- [ ] Branch created with split number in name
- [ ] SPLIT-PLAN-XXX.md exists in directory
- [ ] Previous split complete (if not split-001)

## Remember

**Every split is a completely separate working environment!**
- Own directory
- Own git clone
- Own branch
- Own plan file
- Sequential dependencies

This ensures clean, measurable, mergeable splits!