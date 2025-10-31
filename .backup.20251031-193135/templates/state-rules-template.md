# [Agent] - [STATE_NAME] State Rules

## 📋 PRIMARY DIRECTIVES FOR [STATE_NAME]

**YOU MUST READ EACH RULE LISTED HERE. YOUR READ TOOL CALLS ARE BEING MONITORED.**

### State-Specific Rules (NOT in [agent].md):
1. **R###** - [Rule Name]
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R###-[rule-file].md`
   - Criticality: [BLOCKING/CRITICAL/MANDATORY/INFO] - [Brief description]

2. **R###** - [Rule Name]
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R###-[rule-file].md`
   - Criticality: [Level] - [Brief description]

**Note**: List rules already in the agent's main config that don't need re-reading here.

## 📋 RULE SUMMARY FOR [STATE_NAME] STATE

### Rules Enforced in This State:
- R###: [Rule description] [Criticality level]
- R###: [Rule description] [Criticality level]

### Critical Requirements:
1. [Requirement] - Penalty: -##%
2. [Requirement] - Penalty: -##%

### Success Criteria:
- ✅ [Success condition]
- ✅ [Success condition]

### Failure Triggers:
- ❌ [Failure condition] = [Consequence]
- ❌ [Failure condition] = [Consequence]

## 🚨 [STATE_NAME] IS A VERB - START [ACTION] IMMEDIATELY! 🚨

### IMMEDIATE ACTIONS UPON ENTERING [STATE_NAME]

**THE MOMENT YOU ENTER THIS STATE, YOU MUST:**
1. [Action 1] NOW
2. [Action 2] immediately
3. Check TodoWrite for pending items and process them

**FORBIDDEN - AUTOMATIC FAILURE:**
- ❌ "STATE TRANSITION COMPLETE: Now in [STATE_NAME]" [stops]
- ❌ "Successfully entered [STATE_NAME] state" [waits]
- ❌ "Ready to start [action]" [pauses]
- ❌ "I'm in [STATE_NAME] state" [does nothing]

**REQUIRED - IMMEDIATE ACTION:**
- ✅ "Entering [STATE_NAME], [Starting action] NOW..."
- ✅ "START [ACTION], [doing specific work]..."

## State Context
[Brief description of what this state is for and what should be happening]

## [Additional sections as needed for state-specific guidance]