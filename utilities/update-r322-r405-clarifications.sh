#!/bin/bash
# Script to add R322/R405 clarifications to all state files
# This eliminates confusion between agent stopping and factory continuing

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
AGENT_STATES_DIR="$PROJECT_ROOT/agent-states"

# Counters
TOTAL_FILES=0
UPDATED_FILES=0
SKIPPED_FILES=0

# Log file
LOG_FILE="$PROJECT_ROOT/r322-r405-update-log.txt"
> "$LOG_FILE"  # Clear log file

log() {
    echo "$@" | tee -a "$LOG_FILE"
}

# R322 Clarification text
R322_CLARIFICATION='
### 🛑 R322 CHECKPOINT CLARIFICATION

**AGENT STOPS ≠ FACTORY STOPS**

At R322 checkpoints:
- ✅ Agent STOPS work (completes state, exits, preserves context)
- ✅ Factory CONTINUES (set CONTINUE-SOFTWARE-FACTORY=TRUE)
- ✅ User runs /continue-[agent] to resume next state
- ✅ This is NORMAL, EXPECTED workflow

**R322 = Agent checkpoint = CONTINUE-SOFTWARE-FACTORY=TRUE**

**DO NOT CONFUSE:**
- "MUST STOP" means agent exits (saves context)
- Does NOT mean set CONTINUE-SOFTWARE-FACTORY=FALSE
- Agent stopping is for context preservation
- Factory continuation flag is about automation readiness'

# R405 Enhanced clarification
R405_ENHANCEMENT='
### ✅ WHEN TO USE TRUE (99% of cases):
- ✅ State work completed successfully
- ✅ R322 checkpoint - agent stops BUT factory continues
- ✅ Ready for /continue-[agent] command
- ✅ Waiting for user to continue (NORMAL workflow)
- ✅ Plan ready for review (agent done, factory proceeds)
- ✅ All validations passed
- ✅ No blockers detected

### ❌ WHEN TO USE FALSE (catastrophic only):
- ❌ CATASTROPHIC error preventing ANY continuation
- ❌ Data corruption that would spread
- ❌ State file corruption that prevents recovery
- ❌ **NOT for agent checkpoints (use TRUE)**
- ❌ **NOT for user review points (use TRUE)**
- ❌ **NOT for R322 stops (use TRUE)**
- ❌ **NOT for "being cautious" (use TRUE)**

### 🔴 CRITICAL UNDERSTANDING:
Two independent concepts:
1. **Agent Completion** = Agent finishes state work, exits
2. **Factory Continuation** = Should automation proceed when user continues?

At R322 checkpoints:
- Agent COMPLETES and EXITS (preserves context)
- Factory CONTINUES (TRUE flag)
- These are BOTH normal operations!'

update_file() {
    local file="$1"
    local temp_file="${file}.tmp"

    TOTAL_FILES=$((TOTAL_FILES + 1))

    # Check if file has R322 or R405 sections
    if ! grep -q "R322\|R405\|MUST STOP" "$file"; then
        log "   SKIP: $file (no R322/R405 content)"
        SKIPPED_FILES=$((SKIPPED_FILES + 1))
        return
    fi

    local updated=false

    # Check if clarification already exists
    if grep -q "R322 CHECKPOINT CLARIFICATION" "$file"; then
        log "   SKIP: $file (clarification already exists)"
        SKIPPED_FILES=$((SKIPPED_FILES + 1))
        return
    fi

    # Create temp file with updates
    awk -v r322_text="$R322_CLARIFICATION" -v r405_text="$R405_ENHANCEMENT" '
    BEGIN {
        in_r322 = 0
        in_r405 = 0
        added_r322_clarification = 0
        added_r405_enhancement = 0
    }

    # Detect R322 section start
    /^##.*R322.*MANDATORY STOP/ {
        in_r322 = 1
        print
        next
    }

    # Add R322 clarification after R322 header
    in_r322 == 1 && !added_r322_clarification && /^$/ {
        print r322_text
        print ""
        added_r322_clarification = 1
        in_r322 = 0
        next
    }

    # Detect R405 section
    /^## .*R405.*MANDATORY AUTOMATION CONTINUATION FLAG/ {
        in_r405 = 1
        print
        next
    }

    # Look for WHEN TO USE TRUE section in R405
    in_r405 == 1 && /^### WHEN TO USE TRUE:/ {
        # Skip old TRUE section until we hit FALSE section
        while (getline && !/^### WHEN TO USE FALSE:/) {
            # Skip these lines
        }
        # Now print enhanced version
        print r405_text
        print ""
        # Print the FALSE line we just read
        print
        added_r405_enhancement = 1
        in_r405 = 0
        next
    }

    # Default: print line as-is
    { print }
    ' "$file" > "$temp_file"

    # Check if file was actually modified
    if ! diff -q "$file" "$temp_file" > /dev/null 2>&1; then
        mv "$temp_file" "$file"
        log "   ✅ UPDATED: $file"
        UPDATED_FILES=$((UPDATED_FILES + 1))
    else
        rm -f "$temp_file"
        log "   SKIP: $file (no changes needed)"
        SKIPPED_FILES=$((SKIPPED_FILES + 1))
    fi
}

log "═══════════════════════════════════════════════════════"
log "R322/R405 CLARIFICATION UPDATE"
log "═══════════════════════════════════════════════════════"
log "Starting: $(date)"
log ""

# Process all rules.md files
find "$AGENT_STATES_DIR" -name "rules.md" -type f | while read -r file; do
    update_file "$file"
done

log ""
log "═══════════════════════════════════════════════════════"
log "UPDATE COMPLETE"
log "═══════════════════════════════════════════════════════"
log "Total files processed: $TOTAL_FILES"
log "Files updated: $UPDATED_FILES"
log "Files skipped: $SKIPPED_FILES"
log "Completed: $(date)"
log ""
log "Full log saved to: $LOG_FILE"
