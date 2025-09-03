# Agent Rule Exposure Map

## Where Agents Encounter Rules in Software Factory 2.0

### 1. 🚨 Direct Agent Configuration Files
**Location:** `.claude/agents/*.md`
**When Read:** On agent startup/spawn
**Rules Present:**
- `orchestrator.md`: 3 rules (R001, R010, R006)
- `sw-engineer.md`: 3 rules (R001, R010, R060)  
- `code-reviewer.md`: 3 rules (R001, R010, R055)
- `architect.md`: 3 rules (R001, R010, R057)

**Key Rules:**
- R001 (Pre-flight checks) - BLOCKING
- R010 (Wrong location handling) - MANDATORY

### 2. 🚨 CRITICAL Folder (Always Loaded)
**Location:** `🚨-CRITICAL/*.md`
**When Read:** Automatically on every agent startup
**Rules Present:**
- `000-PRE-FLIGHT-CHECKS.md`: 2 rules
- `001-AGENT-ACKNOWLEDGMENT.md`: 5 rules
- `002-GRADING-SYSTEM.md`: 5 rules
- `003-STATE-MACHINE-NAV.md`: 4 rules
- `004-CONTEXT-RECOVERY.md`: 9 rules
- `005-TEMPLATE-USAGE.md`: 1 rule

**Total:** 26 critical rules always visible

### 3. 📋 State-Specific Rule Files
**Location:** `agent-states/{agent}/{STATE}/rules.md`
**When Read:** When agent enters specific state
**Examples:**
- `orchestrator/SPAWN_AGENTS/rules.md`: R052, R053, R151
- `sw-engineer/IMPLEMENTATION/rules.md`: R011, R012, etc.
- `code-reviewer/CODE_REVIEW/rules.md`: R055, R056, etc.
- `architect/WAVE_REVIEW/rules.md`: R070-R076

**Pattern:** Rules are contextual to current work state

### 4. 🎯 Slash Command Files
**Location:** `.claude/commands/*.md`
**When Read:** When user invokes command
**Rule References:**
- `/continue-orchestrating`: References state machine rules
- `/continue-implementing`: References implementation rules
- `/continue-reviewing`: References review rules
- `/check-status`: References recovery rules

### 5. 📚 State Machine Definitions
**Location:** `state-machines/*.md`
**When Read:** During state transitions
**Rule Integration:**
- Define when specific rules apply
- Link states to rule files
- Specify rule enforcement points

### 6. 📖 Quick Reference Guides
**Location:** `quick-reference/*.md`
**When Read:** As needed for guidance
**Notable Files:**
- `critical-rules-cheatsheet.md`: 10 rule references
- `grading-metrics.md`: Grading rule details
- Agent-specific quick refs: 1-2 rules each

### 7. 🏗️ Rule Library (Source of Truth)
**Location:** `rule-library/`
**When Read:** For rule details/clarification
- `RULE-REGISTRY.md`: Master list of all rules
- Individual rule files: `R###-*.md`
- Formatting guides and meta-documentation

## Agent Rule Loading Sequence

### Orchestrator Startup:
1. **Loads:** `.claude/agents/orchestrator.md` (R001, R010, R006)
2. **Loads:** `🚨-CRITICAL/*.md` (26 rules)
3. **Checks:** Current state in `orchestrator-state.yaml`
4. **Loads:** `agent-states/orchestrator/{STATE}/rules.md`
5. **References:** State machine for transitions

### SW Engineer Spawn:
1. **Loads:** `.claude/agents/sw-engineer.md` (R001, R010, R060)
2. **Loads:** `🚨-CRITICAL/*.md` (if first startup)
3. **Receives:** Context from orchestrator (including rules)
4. **Loads:** `agent-states/sw-engineer/IMPLEMENTATION/rules.md`
5. **Monitors:** Size rules (R002) continuously

### Code Reviewer Spawn:
1. **Loads:** `.claude/agents/code-reviewer.md` (R001, R010, R055)
2. **Loads:** `🚨-CRITICAL/*.md` (if first startup)
3. **Loads:** Either `EFFORT_PLANNING` or `CODE_REVIEW` rules
4. **Enforces:** Size limits, quality standards

### Architect Review:
1. **Loads:** `.claude/agents/architect.md` (R001, R010, R057)
2. **Loads:** `🚨-CRITICAL/*.md` (if first startup)
3. **Loads:** `WAVE_REVIEW` or `PHASE_ASSESSMENT` rules
4. **Evaluates:** Against architectural rules (R070-R080)

## Rule Visibility by Criticality

### Always Visible (Can't Miss):
- 🚨🚨🚨 **BLOCKING** rules in agent configs (R001)
- 🚨🚨 **MANDATORY** rules in CRITICAL folder
- 🚨 **CRITICAL** grading rules

### Context-Specific:
- State-specific rules when in that state
- Task-specific rules when performing that task
- Recovery rules when context lost

### Reference Only:
- ℹ️ **INFO** rules in quick references
- Details in rule library
- Examples in documentation

## Key Observations

1. **Multi-Layer Exposure**: Rules appear in multiple places to ensure agents see them
2. **Criticality-Based Placement**: Most critical rules are in startup files
3. **Context-Aware Loading**: State-specific rules load only when needed
4. **Redundancy by Design**: Important rules appear in multiple locations
5. **Clear Hierarchy**: Visual formatting ensures critical rules stand out

## Files Agents MUST Read

### Every Agent, Every Time:
1. Their own config: `.claude/agents/{agent-type}.md`
2. Pre-flight checks: `🚨-CRITICAL/000-PRE-FLIGHT-CHECKS.md`
3. Agent acknowledgment: `🚨-CRITICAL/001-AGENT-ACKNOWLEDGMENT.md`

### Based on Role:
- **Orchestrator**: State machine, grading system
- **SW Engineer**: Implementation rules, size limits
- **Code Reviewer**: Review criteria, split protocols
- **Architect**: Assessment criteria, patterns

### Based on State:
- Current state rules from `agent-states/{agent}/{STATE}/rules.md`
- Checkpoint files for state recovery
- Grading files for performance metrics

This multi-layered approach ensures agents encounter the right rules at the right time, with critical rules impossible to miss!