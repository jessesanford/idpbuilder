# Files Not Included from Original TMC Config

These files were in the original TMC orchestrator config but were NOT included in the generic template because they are either:
1. TMC/KCP specific
2. Example artifacts from actual runs
3. Redundant with other files
4. Too specific to generalize

## Files to Copy Here

### Quick Reference Guides (Useful but Optional)
- **CODE-REVIEWER-QUICK-REFERENCE.md** - Quick lookup for reviewers
- **ORCHESTRATOR-QUICK-REFERENCE.md** - Quick lookup for orchestrator
- **ORCHESTRATOR-WORKFLOW-SUMMARY.md** - Workflow summary

### Additional Protocols (May be Useful)
- **ARCHITECT-REVIEWER-WAVE-INSTRUCTIONS.md** - Detailed architect instructions
- **PHASE-START-ARCHITECT-REVIEW-PROTOCOL.md** - Phase assessment protocol
- **PHASE-COMPLETION-FUNCTIONAL-TESTING.md** - End-of-phase testing
- **TODO-STATE-MANAGEMENT-PROTOCOL.md** - Detailed TODO management
- **KCP-CODE-REVIEWER-COMPREHENSIVE-GUIDE.md** - Needs to be generalized
- **ORCHESTRATOR-CODE-REVIEW-INTEGRATION.md** - Review integration details
- **ORCHESTRATOR-NEVER-WRITES-CODE-RULE.md** - Already in CLAUDE.md but could be standalone

### Example Artifacts (For Reference)
- **SPLIT-EXAMPLE-E3.1.1-SYNC-ENGINE.md** - Example of an actual split
- **SPLIT-REVIEW-LOOP-DIAGRAM.md** - Visual diagram of split process
- **CODE-REVIEW-EXAMPLES.md** - Example reviews
- **CODE-REVIEW-ENFORCEMENT-SUMMARY.md** - Summary of enforcement

### TMC-Specific (Not Generic)
- **TMC-ORCHESTRATOR-IMPLEMENTATION-PLAN-8-20-2025.md** - TMC specific
- **TMC-FEATURE-GAP-ANALYSIS-SYNTHESIS-PLAN-8-20-2025.md** - TMC specific
- **PHASE{1-5}-SPECIFIC-IMPL-PLAN-8-20-25.md** - TMC specific phases
- **PHASE2-CORRECTION-NOTE.md** - TMC specific correction
- **SOFTWARE-ENG-AGENT-EXPLICIT-INSTRUCTIONS.md** - Has TMC-specific Git commands
- **BRANCH-STRUCTURE-SUMMARY*.md** - TMC branch structure
- **CURRENT-TODO-STATE.md** - Active TMC state

### Already Included or Moved
✅ SOFTWARE-FACTORY-STATE-MACHINE.md - In core/
✅ ORCHESTRATOR-MASTER-OPERATIONS-GUIDE.md - In core/
✅ EFFORT-SPLIT-CONTINUOUS-EXECUTION-PROTOCOL.md - In protocols/
✅ IMPERATIVE-LINE-COUNT-RULE.md - In protocols/
✅ ORCHESTRATOR-EFFORT-PLANNING-PROTOCOL.md - In protocols/
✅ ORCHESTRATOR-TASKMASTER-EXECUTION-PLAN.md - In protocols/
✅ CODE-REVIEWER-EFFORT-PLANNING-INSTRUCTIONS.md - In protocols/
✅ WAVE-COMPLETION-ARCHITECT-REVIEW-PROTOCOL.md - In protocols/
✅ SW-ENGINEER-STARTUP-REQUIREMENTS.md - In protocols/ (as SOFTWARE-ENG-AGENT-STARTUP-REQUIREMENTS.md)
✅ TEST-DRIVEN-VALIDATION-REQUIREMENTS.md - In protocols/
✅ WORK-LOG-TEMPLATE.md - In protocols/