# 🚨 DEPRECATED - DO NOT USE 🚨

**This file was deprecated on 2025-09-01**

## Why This File Was Deprecated

This file contained massive duplication with rules already present in:
- `.claude/agents/orchestrator.md` (supreme laws)
- State-specific rule files in `agent-states/orchestrator/[STATE]/rules.md`

The duplication was causing:
- Confusion about which rules to follow
- Unnecessary reading overhead (862 lines of mostly duplicate content)
- Risk of rules getting out of sync
- Violation of single source of truth principle

## Where The Rules Moved To

All rules have been redistributed to their proper locations:

### Supreme Laws
Now exclusively in `.claude/agents/orchestrator.md`:
- R234 - Mandatory State Traversal (HIGHEST SUPREME LAW)
- R208 - CD Before Spawn (SUPREME LAW #2)
- R290 - State Rule Reading Verification (SUPREME LAW #3)
- R221 - Bash Directory Reset Protocol (SUPREME LAW #4)
- R235 - Mandatory Pre-flight Verification (SUPREME LAW #5)
- R021 - Orchestrator Never Stops (SUPREME LAW #8)
- R231 - Continuous Operation (SUPREME LAW #9)
- R232 - TodoWrite Enforcement (SUPREME LAW #10)
- R281 - Complete State File Initialization (SUPREME LAW #7)
- R288 - State File Update and Commit Protocol (SUPREME LAW)

### State-Specific Rules
Moved to appropriate state directories:
- **R209 metadata injection code** → `agent-states/orchestrator/SETUP_EFFORT_INFRASTRUCTURE/rules.md`
- **R254 sequential split processing** → Already in `agent-states/orchestrator/MONITOR/rules.md`
- **R255 post-agent verification** → Already in `agent-states/orchestrator/MONITOR/rules.md`
- **R250 integration isolation** → Already in `agent-states/orchestrator/INTEGRATION/rules.md`
- **R212 wave completion** → `agent-states/orchestrator/WAVE_COMPLETE/rules.md`
- **R213 phase completion** → `agent-states/orchestrator/SUCCESS/rules.md`

### Mandatory Rules
In `.claude/agents/orchestrator.md`:
- R287 - Todo Persistence
- R216 - Bash Execution Syntax
- R206 - State Machine Validation
- R203 - State-Aware Startup

## DO NOT READ THIS FILE

**This file is no longer maintained and should not be read during agent startup.**

The orchestrator startup sequence has been updated to remove this file from the required reading list.

## For Historical Reference Only

This file is kept for historical reference and audit trail purposes only.
All active rules are maintained in their proper locations as listed above.