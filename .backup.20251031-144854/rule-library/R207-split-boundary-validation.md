# R207: Split Boundary Validation Protocol

**Category:** Quality Gates  
**Agents:** Code Reviewer, SW Engineer, Orchestrator  
**Criticality:** BLOCKING - Cross-effort references corrupt split integrity

## Purpose
Prevent split plans from referencing splits from different efforts, ensuring each effort's splits remain isolated and properly bounded.

## The Problem This Solves

When Split 002 of `registry-auth-types` incorrectly references "Split 001 (oci-types)" as its previous split, this creates:
- Confusion about dependencies
- Incorrect merge ordering
- Lost context about what was actually implemented before
- Broken split chains

## Validation Rules

### 1. Same Effort Only
```bash
# When creating Split 002, validate previous split reference
CURRENT_EFFORT=$(pwd | grep -oP 'effort-[^/]+')
PREVIOUS_SPLIT_REF="Split 001 of phase2/wave1/$CURRENT_EFFORT"

# NEVER reference different effort
❌ WRONG: "Previous Split: Split 001 (oci-types)"
✅ RIGHT: "Previous Split: Split 001 of phase2/wave1/registry-auth-types"
```

### 2. Full Path Requirements

Every split boundary MUST include:
```markdown
## Boundaries
- **Previous Split**: Split [N-1] of phase[X]/wave[Y]/effort-[name]
  - Path: efforts/phase[X]/wave[Y]/effort-[name]/split-[N-1]/
  - Branch: phase[X]/wave[Y]/effort-[name]-split-[N-1]
  - Summary: [What the previous split in THIS effort implemented]
```

### 3. Validation Function

```bash
validate_split_boundaries() {
    local split_plan="$1"
    local current_effort="$2"
    
    # Extract previous split reference
    PREV_REF=$(grep "Previous Split:" "$split_plan" | head -1)
    
    # Ensure it references same effort
    if ! echo "$PREV_REF" | grep -q "$current_effort"; then
        echo "❌ FATAL: Split boundary references different effort!"
        echo "   Found: $PREV_REF"
        echo "   Expected: Reference to $current_effort"
        exit 1
    fi
    
    # Ensure full path is present
    if ! grep -q "Path: efforts/phase.*/wave.*/effort-$current_effort" "$split_plan"; then
        echo "❌ FATAL: Split boundary missing full path!"
        exit 1
    fi
    
    echo "✅ Split boundaries validated - same effort maintained"
}
```

## Code Reviewer Responsibilities

When creating split plans:

1. **Identify Current Effort**
   ```bash
   CURRENT_PHASE=$(pwd | grep -oP 'phase\d+')
   CURRENT_WAVE=$(pwd | grep -oP 'wave\d+')
   CURRENT_EFFORT=$(pwd | grep -oP 'effort-[^/]+')
   ```

2. **Reference Only Same Effort Splits**
   ```markdown
   # For Split 002 of registry-auth-types:
   Previous Split: Split 001 of phase2/wave1/registry-auth-types
   NOT: Split 001 (oci-types) ← Different effort!
   ```

3. **Include Full Paths**
   ```markdown
   Path: efforts/phase2/wave1/registry-auth-types/split-001/
   Branch: phase2/wave1/registry-auth-types-split-001
   ```

## SW Engineer Responsibilities

When reading split plans:

1. **Verify Previous Split**
   ```bash
   # Check that previous split is from same effort
   if [ -f "../split-001/SPLIT-PLAN-001.md" ]; then
       echo "✅ Previous split found in same effort"
   else
       echo "❌ Previous split not in expected location!"
       exit 1
   fi
   ```

2. **Validate Branch Naming**
   ```bash
   # Branch should match effort
   EXPECTED_BRANCH="${CURRENT_EFFORT}-split-${SPLIT_NUM}"
   if [[ "$CURRENT_BRANCH" != *"$EXPECTED_BRANCH"* ]]; then
       echo "❌ Branch name doesn't match effort!"
       exit 1
   fi
   ```

## Orchestrator Responsibilities

When creating split infrastructure:

1. **Maintain Effort Isolation**
   ```bash
   # Each split directory under same effort
   efforts/phase2/wave1/registry-auth-types/
   ├── split-001/
   ├── split-002/
   └── split-003/
   ```

2. **Validate Split Chain**
   ```bash
   # Ensure all splits reference same parent
   for split in split-*/SPLIT-PLAN-*.md; do
       if ! grep -q "effort-$EFFORT_NAME" "$split"; then
           echo "❌ Split $split references wrong effort!"
           exit 1
       fi
   done
   ```

## Examples

### ✅ Correct Split 002 Header

```markdown
# SPLIT-PLAN-002.md
## Split 002 of 3: Controller Implementation
**Parent Effort**: registry-auth-types

## Boundaries
- **Previous Split**: Split 001 of phase2/wave1/registry-auth-types
  - Path: efforts/phase2/wave1/registry-auth-types/split-001/
  - Branch: phase2/wave1/registry-auth-types-split-001
  - Summary: Implemented core API types and validation
- **This Split**: Split 002 of phase2/wave1/registry-auth-types
  - Path: efforts/phase2/wave1/registry-auth-types/split-002/
  - Branch: phase2/wave1/registry-auth-types-split-002
```

### ❌ Incorrect Split 002 Header

```markdown
# SPLIT-PLAN-002.md
## Split 002 of 2: Stack Types

## Boundaries
- **Previous Split**: Split 001 (OCI types)  ← WRONG EFFORT!
- **This Split Focus**: Stack type definitions
```

## Validation Script

```bash
#!/bin/bash
# validate-split-boundaries.sh

SPLIT_PLAN="$1"
EFFORT_NAME=$(pwd | grep -oP 'effort-[^/]+')

echo "Validating split boundaries for $EFFORT_NAME..."

# Check previous split reference
if grep -q "Previous Split.*$EFFORT_NAME" "$SPLIT_PLAN"; then
    echo "✅ Previous split correctly references $EFFORT_NAME"
else
    echo "❌ FATAL: Previous split references different effort!"
    grep "Previous Split:" "$SPLIT_PLAN"
    exit 1
fi

# Check for full paths
if grep -q "Path: efforts/phase.*/wave.*/$EFFORT_NAME" "$SPLIT_PLAN"; then
    echo "✅ Full paths included"
else
    echo "❌ FATAL: Missing full path information!"
    exit 1
fi

# Check branch naming
if grep -q "Branch:.*$EFFORT_NAME-split-" "$SPLIT_PLAN"; then
    echo "✅ Branch naming consistent"
else
    echo "❌ WARNING: Branch naming may be inconsistent"
fi

echo "✅ Split boundary validation passed!"
```

## Enforcement

1. **Pre-split Check**: Validate boundaries before creating splits
2. **Review Gate**: Code reviewer must verify split boundaries
3. **SW Engineer Check**: Validate before starting implementation
4. **Merge Gate**: Orchestrator validates before integration

## Summary

- Every split MUST reference only splits from the SAME effort
- Full phase/wave/effort paths are REQUIRED
- Cross-effort references are FATAL errors
- Validate boundaries before AND after split creation