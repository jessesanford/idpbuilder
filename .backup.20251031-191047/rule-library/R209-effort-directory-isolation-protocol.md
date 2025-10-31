# R209: Effort Directory Isolation Protocol

**Category:** Critical Rules  
**Agents:** Orchestrator, SW Engineer, Code Reviewer  
**Criticality:** MISSION CRITICAL - Directory violations = Implementation failures  
**Priority:** HIGHEST - Enforce at every stage

## 🚨 MISSION CRITICAL: ALL WORK MUST STAY IN EFFORT DIRECTORY 🚨

SW Engineers MUST work EXCLUSIVELY in their assigned effort directory. No exceptions. No wandering. No confusion.

## The Problem This Solves

SW Engineers frequently:
- Start in the correct directory but navigate away
- Create code in the wrong location
- Lose track of their isolation boundary
- Don't have clear metadata about their workspace

This causes:
- Code in wrong locations
- Merge conflicts
- Lost work
- Grading failures

## The Solution: Three-Part Protocol

### Part 1: Orchestrator MUST Inject Metadata (Like Split Plans!)

Just as we inject metadata into split plans, the orchestrator MUST inject workspace metadata into EVERY implementation plan:

```bash
# ORCHESTRATOR RESPONSIBILITY: Add metadata to implementation plans
inject_effort_metadata() {
    local EFFORT_NAME="$1"
    local PHASE="$2"
    local WAVE="$3"
    local BRANCH="$4"
    local IMPL_PLAN="efforts/phase${PHASE}/wave${WAVE}/${EFFORT_NAME}/IMPLEMENTATION-PLAN.md"
    
    # Check if metadata already exists
    if grep -q "EFFORT INFRASTRUCTURE METADATA" "$IMPL_PLAN"; then
        echo "✅ Metadata already present"
        return 0
    fi
    
    # Inject metadata at the TOP of the plan
    cat > /tmp/metadata_header.md << EOF
<!-- ⚠️ EFFORT INFRASTRUCTURE METADATA (Added by Orchestrator per R209) ⚠️ -->
<!-- SW ENGINEERS: YOU MUST STAY IN THIS DIRECTORY FOR ALL WORK! -->
**WORKING_DIRECTORY**: $(pwd)/efforts/phase${PHASE}/wave${WAVE}/${EFFORT_NAME}
**BRANCH**: ${BRANCH}
**REMOTE**: origin/${BRANCH}
**BASE_BRANCH**: main
**EFFORT_NAME**: ${EFFORT_NAME}
**PHASE**: ${PHASE}
**WAVE**: ${WAVE}
**ISOLATION_BOUNDARY**: efforts/phase${PHASE}/wave${WAVE}/${EFFORT_NAME}

⚠️ **CRITICAL ISOLATION RULES** ⚠️
1. ALL work MUST happen in: efforts/phase${PHASE}/wave${WAVE}/${EFFORT_NAME}
2. ALL code MUST go in: efforts/phase${PHASE}/wave${WAVE}/${EFFORT_NAME}/pkg/
3. NEVER cd out of this directory during implementation
4. NEVER create files outside this boundary
5. This is YOUR isolated workspace - stay in it!
<!-- END METADATA -->

EOF
    
    # Prepend metadata to implementation plan
    if [ -f "$IMPL_PLAN" ]; then
        cat /tmp/metadata_header.md "$IMPL_PLAN" > /tmp/updated_plan.md
        mv /tmp/updated_plan.md "$IMPL_PLAN"
        echo "✅ Metadata injected into $IMPL_PLAN"
    else
        echo "❌ ERROR: $IMPL_PLAN not found!"
        return 1
    fi
}
```

### Part 2: SW Engineer Enhanced Preflight Checks

SW Engineers MUST verify and enforce directory isolation:

```bash
# SW ENGINEER PREFLIGHT: Directory Isolation Check
enforce_directory_isolation() {
    echo "═══════════════════════════════════════════════════════"
    echo "🚨 R209: DIRECTORY ISOLATION ENFORCEMENT 🚨"
    echo "═══════════════════════════════════════════════════════"
    
    # 1. Extract metadata from implementation plan
    if [ -f "IMPLEMENTATION-PLAN.md" ]; then
        WORKING_DIR=$(grep "**WORKING_DIRECTORY**:" IMPLEMENTATION-PLAN.md | cut -d: -f2- | xargs)
        BRANCH=$(grep "**BRANCH**:" IMPLEMENTATION-PLAN.md | cut -d: -f2- | xargs)
        ISOLATION=$(grep "**ISOLATION_BOUNDARY**:" IMPLEMENTATION-PLAN.md | cut -d: -f2- | xargs)
    else
        echo "❌ FATAL: No IMPLEMENTATION-PLAN.md found!"
        exit 1
    fi
    
    # 2. Verify we're in the correct directory
    CURRENT_DIR=$(pwd)
    if [ "$CURRENT_DIR" != "$WORKING_DIR" ]; then
        echo "❌ FATAL: Wrong directory!"
        echo "   Expected: $WORKING_DIR"
        echo "   Actual: $CURRENT_DIR"
        
        # Attempt to navigate to correct directory
        if [ -d "$WORKING_DIR" ]; then
            echo "🔄 Navigating to correct directory..."
            cd "$WORKING_DIR" || exit 1
            echo "✅ Now in: $(pwd)"
        else
            echo "❌ Directory doesn't exist: $WORKING_DIR"
            exit 1
        fi
    else
        echo "✅ Correct directory: $CURRENT_DIR"
    fi
    
    # 3. Verify git branch
    CURRENT_BRANCH=$(git branch --show-current)
    if [ "$CURRENT_BRANCH" != "$BRANCH" ]; then
        echo "⚠️ WARNING: Wrong branch!"
        echo "   Expected: $BRANCH"
        echo "   Actual: $CURRENT_BRANCH"
    fi
    
    # 4. Set up isolation enforcement
    echo "📍 ISOLATION BOUNDARY SET: $ISOLATION"
    echo "⚠️ YOU MUST NOT LEAVE THIS DIRECTORY!"
    echo ""
    echo "Your workspace structure:"
    echo "  $(pwd)/"
    echo "  ├── IMPLEMENTATION-PLAN.md  (your guide)"
    echo "  ├── work-log.md            (track progress)"
    echo "  └── pkg/                   (ALL CODE HERE)"
    echo ""
    echo "❌ FORBIDDEN ACTIONS:"
    echo "  - cd .. (leaving effort directory)"
    echo "  - cd / (going to root)"
    echo "  - Creating files outside pkg/"
    echo "  - Modifying files outside this effort"
    echo ""
    echo "✅ ALLOWED ACTIONS:"
    echo "  - cd pkg (entering code directory)"
    echo "  - Creating subdirectories under pkg/"
    echo "  - All work within $(pwd)"
    echo "═══════════════════════════════════════════════════════"
}

# Run at startup
enforce_directory_isolation

# Periodic check during work
verify_still_isolated() {
    CURRENT=$(pwd)
    if [[ "$CURRENT" != *"efforts/phase"*/wave*/* ]]; then
        echo "❌❌❌ FATAL: You've left your effort directory!"
        echo "Current: $CURRENT"
        echo "YOU MUST RETURN TO YOUR EFFORT DIRECTORY!"
        exit 1
    fi
}
```

### Part 3: Continuous Isolation Verification

```bash
# Add to SW Engineer work cycle
work_cycle_with_isolation() {
    # Before EVERY file operation
    verify_still_isolated
    
    # Before creating files
    create_file_with_check() {
        local FILE="$1"
        verify_still_isolated
        
        # Ensure file is under pkg/
        if [[ "$FILE" != pkg/* ]] && [[ "$FILE" != ./pkg/* ]]; then
            echo "❌ ERROR: Files must be created under pkg/"
            echo "   Attempted: $FILE"
            echo "   Use: pkg/$FILE instead"
            return 1
        fi
        
        # Create file
        mkdir -p $(dirname "$FILE")
        touch "$FILE"
        echo "✅ Created: $FILE"
    }
    
    # Before git operations
    git_with_isolation_check() {
        verify_still_isolated
        git "$@"
    }
}
```

## Orchestrator Implementation

The orchestrator MUST:

1. **Create effort directory structure**
2. **Create implementation plan**
3. **Inject metadata into plan** (NEW!)
4. **Verify metadata before spawn**

```bash
# COMPLETE ORCHESTRATOR WORKFLOW WITH R209
prepare_effort_with_metadata() {
    local PHASE="$1"
    local WAVE="$2"
    local EFFORT="$3"
    local BRANCH="phase${PHASE}/wave${WAVE}/${EFFORT}"
    local EFFORT_DIR="efforts/phase${PHASE}/wave${WAVE}/${EFFORT}"
    
    echo "═══════════════════════════════════════════════════════"
    echo "R209: Preparing Effort with Metadata Injection"
    echo "═══════════════════════════════════════════════════════"
    
    # 1. Create directory structure
    mkdir -p "$EFFORT_DIR"/{pkg,tests,docs}
    cd "$EFFORT_DIR" || exit 1
    
    # 2. Initialize git
    git init
    git checkout -b "$BRANCH"
    
    # 3. Create initial files
    touch IMPLEMENTATION-PLAN.md
    touch work-log.md
    
    # 4. Wait for Code Reviewer to create plan
    echo "Spawning Code Reviewer to create plan..."
    # [spawn code reviewer]
    
    # 5. CRITICAL: Inject metadata after plan exists
    echo "🚨 R209: Injecting effort metadata..."
    inject_effort_metadata "$EFFORT" "$PHASE" "$WAVE" "$BRANCH"
    
    # 6. Verify metadata
    if grep -q "EFFORT INFRASTRUCTURE METADATA" IMPLEMENTATION-PLAN.md; then
        echo "✅ Metadata successfully injected"
        echo "✅ SW Engineer will be locked to: $EFFORT_DIR"
    else
        echo "❌ FATAL: Metadata injection failed!"
        exit 1
    fi
    
    # 7. Return to orchestrator directory
    cd - > /dev/null
    
    echo "✅ Effort prepared with isolation metadata"
    echo "═══════════════════════════════════════════════════════"
}
```

## SW Engineer Startup Sequence - MANDATORY ACKNOWLEDGMENT

```bash
# MANDATORY SW ENGINEER STARTUP WITH EXPLICIT ACKNOWLEDGMENT
sw_engineer_startup() {
    echo "═══════════════════════════════════════════════════════"
    echo "🚨 R209: MANDATORY DIRECTORY ISOLATION CHECK 🚨"
    echo "═══════════════════════════════════════════════════════"
    
    # 1. Read implementation plan
    if [ ! -f "IMPLEMENTATION-PLAN.md" ]; then
        echo "❌ FATAL: No implementation plan!"
        exit 1
    fi
    
    # 2. Extract and verify metadata
    WORKING_DIR=$(grep "**WORKING_DIRECTORY**:" IMPLEMENTATION-PLAN.md | cut -d: -f2- | xargs)
    BRANCH=$(grep "**BRANCH**:" IMPLEMENTATION-PLAN.md | cut -d: -f2- | xargs)
    ISOLATION=$(grep "**ISOLATION_BOUNDARY**:" IMPLEMENTATION-PLAN.md | cut -d: -f2- | xargs)
    EFFORT_NAME=$(grep "**EFFORT_NAME**:" IMPLEMENTATION-PLAN.md | cut -d: -f2- | xargs)
    PHASE=$(grep "**PHASE**:" IMPLEMENTATION-PLAN.md | cut -d: -f2- | xargs)
    WAVE=$(grep "**WAVE**:" IMPLEMENTATION-PLAN.md | cut -d: -f2- | xargs)
    
    if [ -z "$WORKING_DIR" ] || [ -z "$ISOLATION" ]; then
        echo "❌ FATAL: No R209 metadata in plan!"
        echo "Orchestrator failed to inject metadata!"
        exit 1
    fi
    
    # 3. Navigate to correct directory if needed
    CURRENT_DIR=$(pwd)
    if [ "$CURRENT_DIR" != "$WORKING_DIR" ]; then
        echo "⚠️ Not in correct directory yet"
        echo "   Expected: $WORKING_DIR"
        echo "   Current: $CURRENT_DIR"
        
        # Attempt to navigate to correct directory
        if [ -d "$WORKING_DIR" ]; then
            echo "📂 Navigating to correct directory..."
            cd "$WORKING_DIR"
            CURRENT_DIR=$(pwd)
            echo "✅ Successfully navigated to: $CURRENT_DIR"
        else
            echo "❌ FATAL: Target directory doesn't exist: $WORKING_DIR"
            echo "Orchestrator failed to create infrastructure!"
            exit 1
        fi
    else
        echo "✅ Already in correct directory: $CURRENT_DIR"
    fi
    
    # 4. Verify git branch
    CURRENT_BRANCH=$(git branch --show-current)
    
    # 5. MANDATORY EXPLICIT ACKNOWLEDGMENT
    echo ""
    echo "═══════════════════════════════════════════════════════"
    echo "📋 R209 DIRECTORY ISOLATION ACKNOWLEDGMENT"
    echo "═══════════════════════════════════════════════════════"
    echo "I, SW ENGINEER, EXPLICITLY ACKNOWLEDGE:"
    echo ""
    echo "✅ WORKING DIRECTORY CONFIRMED:"
    echo "   Required: $WORKING_DIR"
    echo "   Current:  $CURRENT_DIR"
    echo "   Status:   CORRECT ✓"
    echo ""
    echo "✅ GIT BRANCH CONFIRMED:"
    echo "   Required: $BRANCH"
    echo "   Current:  $CURRENT_BRANCH"
    if [ "$CURRENT_BRANCH" = "$BRANCH" ]; then
        echo "   Status:   CORRECT ✓"
    else
        echo "   Status:   MISMATCH ⚠️ (may need to checkout)"
    fi
    echo ""
    echo "✅ ISOLATION BOUNDARY CONFIRMED:"
    echo "   Boundary: $ISOLATION"
    echo "   Effort:   $EFFORT_NAME (Phase $PHASE, Wave $WAVE)"
    echo ""
    echo "📍 I UNDERSTAND AND AGREE TO:"
    echo "   1. ✓ ALL work happens in: $(pwd)"
    echo "   2. ✓ ALL code goes in: $(pwd)/pkg/"
    echo "   3. ✓ NEVER cd out of this directory"
    echo "   4. ✓ NEVER create files outside this boundary"
    echo "   5. ✓ This is MY isolated workspace"
    echo ""
    echo "🔒 ISOLATION LOCK ENGAGED: $(pwd)"
    echo "═══════════════════════════════════════════════════════"
    
    # 6. Set up UNREMOVABLE guard environment variables
    export EFFORT_ISOLATION_DIR="$(pwd)"
    export EFFORT_NAME="$EFFORT_NAME"
    export EFFORT_PHASE="$PHASE"
    export EFFORT_WAVE="$WAVE"
    
    # CRITICAL: Make environment variable unmodifiable (readonly)
    readonly EFFORT_ISOLATION_DIR
    readonly EFFORT_NAME
    readonly EFFORT_PHASE
    readonly EFFORT_WAVE
    
    # Set up cd guard function to prevent leaving directory
    cd() {
        local target="${1:-$HOME}"
        local abs_target=$(realpath "$target" 2>/dev/null || echo "$target")
        
        # Check if target is within isolation boundary
        if [[ "$abs_target" != "$EFFORT_ISOLATION_DIR"* ]]; then
            echo "❌❌❌ R209 VIOLATION: Cannot leave effort directory!"
            echo "   Attempted: $abs_target"
            echo "   Boundary:  $EFFORT_ISOLATION_DIR"
            echo "   BLOCKED - Staying in: $(pwd)"
            return 1
        fi
        
        # Allow cd within the boundary
        builtin cd "$@"
    }
    
    # Export the cd function
    export -f cd
    
    # 7. Create acknowledgment file for audit
    echo "R209 Acknowledged at $(date '+%Y-%m-%d %H:%M:%S')" >> .r209-acknowledged
    echo "Directory: $(pwd)" >> .r209-acknowledged
    echo "Branch: $CURRENT_BRANCH" >> .r209-acknowledged
    echo "Readonly Environment Set: YES" >> .r209-acknowledged
    echo "CD Guard Active: YES" >> .r209-acknowledged
    echo "" >> .r209-acknowledged
    
    # 8. Final verification message
    echo ""
    echo "🔐 ENVIRONMENT LOCKED:"
    echo "   EFFORT_ISOLATION_DIR=$EFFORT_ISOLATION_DIR (READONLY)"
    echo "   cd() function overridden - cannot leave boundary"
    echo "   Acknowledgment file created: .r209-acknowledged"
}
```

## Validation Script

```bash
#!/bin/bash
# validate-effort-isolation.sh

validate_implementation_plan() {
    local PLAN="$1"
    
    echo "Validating R209 metadata in $PLAN..."
    
    # Check for required metadata fields
    REQUIRED_FIELDS=(
        "WORKING_DIRECTORY"
        "BRANCH"
        "ISOLATION_BOUNDARY"
        "EFFORT_NAME"
        "CRITICAL ISOLATION RULES"
    )
    
    for field in "${REQUIRED_FIELDS[@]}"; do
        if grep -q "$field" "$PLAN"; then
            echo "✅ Found: $field"
        else
            echo "❌ Missing: $field"
            return 1
        fi
    done
    
    echo "✅ All R209 metadata present"
    return 0
}
```

## Examples

### ✅ Correct Implementation Plan with Metadata

```markdown
<!-- ⚠️ EFFORT INFRASTRUCTURE METADATA (Added by Orchestrator per R209) ⚠️ -->
<!-- SW ENGINEERS: YOU MUST STAY IN THIS DIRECTORY FOR ALL WORK! -->
**WORKING_DIRECTORY**: /workspaces/project/efforts/phase1/wave1/api-types
**BRANCH**: phase1/wave1/api-types
**REMOTE**: origin/phase1/wave1/api-types
**BASE_BRANCH**: main
**EFFORT_NAME**: api-types
**PHASE**: 1
**WAVE**: 1
**ISOLATION_BOUNDARY**: efforts/phase1/wave1/api-types

⚠️ **CRITICAL ISOLATION RULES** ⚠️
1. ALL work MUST happen in: efforts/phase1/wave1/api-types
2. ALL code MUST go in: efforts/phase1/wave1/api-types/pkg/
<!-- END METADATA -->

# API Types Implementation Plan
[rest of plan...]
```

### ❌ Wrong: Plan without metadata

```markdown
# API Types Implementation Plan
[plan content without metadata]
```

## Enforcement Points

1. **Orchestrator**: Must inject metadata after plan creation
2. **SW Engineer**: Must verify metadata on startup
3. **SW Engineer**: Must check isolation before every operation
4. **Code Reviewer**: Must verify metadata present in plans

## Summary

- **Orchestrator** injects metadata into implementation plans (like split plans!)
- **SW Engineers** read metadata and enforce isolation
- **All work** stays in the effort directory
- **No exceptions** to the isolation boundary
- **This is MISSION CRITICAL** for implementation success