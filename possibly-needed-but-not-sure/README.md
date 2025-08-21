# Possibly Needed But Not Sure - Complete Collection

This directory contains additional protocol and instruction files from the production TMC implementation that may be useful depending on your project's specific needs. These were not included in the main template because they are either optional, examples, or need customization.

## Categories of Files

### 📋 Quick Reference Guides
**When to Use**: If you want quick lookup guides for agents

- **CODE-REVIEWER-QUICK-REFERENCE.md** - Quick checks and decision tree for code reviewers
- **ORCHESTRATOR-QUICK-REFERENCE.md** - Quick commands and state transitions
- **ORCHESTRATOR-WORKFLOW-SUMMARY.md** - Visual workflow summary with examples

### 🏗️ Architect Protocols
**When to Use**: For projects requiring strict architectural governance

- **ARCHITECT-REVIEWER-WAVE-INSTRUCTIONS.md** - Detailed instructions for architect reviews
- **PHASE-START-ARCHITECT-REVIEW-PROTOCOL.md** - How to assess phase readiness
- **PHASE-COMPLETION-FUNCTIONAL-TESTING.md** - End-of-phase testing requirements

### 📝 Review and Integration Guides
**When to Use**: For complex projects with multiple reviewers

- **CODE-REVIEWER-COMPREHENSIVE-GUIDE-EXAMPLE.md** - Complete review process (KCP-specific example)
- **ORCHESTRATOR-CODE-REVIEW-INTEGRATION.md** - How orchestrator integrates reviews
- **CODE-REVIEW-ENFORCEMENT-SUMMARY.md** - Summary of review enforcement points
- **CODE-REVIEW-EXAMPLES.md** - Real examples of code reviews

### 🔄 State Management
**When to Use**: For advanced TODO and state management

- **TODO-STATE-MANAGEMENT-PROTOCOL.md** - Detailed TODO persistence and recovery
- **ORCHESTRATOR-NEVER-WRITES-CODE-RULE.md** - Standalone rule document (already in CLAUDE.md)

### 📊 Split Examples and Diagrams
**When to Use**: To understand how splits work in practice

- **SPLIT-EXAMPLE-E3.1.1-SYNC-ENGINE.md** - Real example of splitting a 2000+ line effort
- **SPLIT-REVIEW-LOOP-DIAGRAM.md** - Visual diagram of the split review process

### 📄 Files Already Moved to Active Use
These were identified as critical and moved to main directories:
- ✅ TEST-DRIVEN-VALIDATION-REQUIREMENTS.md → /protocols/
- ✅ WORK-LOG-TEMPLATE.md → /protocols/

## Quick Decision Guide

### You SHOULD Move These to Active Use If:

**Your Project Has Multiple Phases/Waves**:
- Move all Architect Protocol files
- Move Phase Completion Testing

**You Want Better Agent Guidance**:
- Move Quick Reference guides
- Move Comprehensive Review Guide

**You're Managing Complex State**:
- Move TODO-STATE-MANAGEMENT-PROTOCOL.md

**You're Learning the System**:
- Keep Split Examples handy for reference
- Review CODE-REVIEW-EXAMPLES.md

### You Can IGNORE These If:

**Simple Single-Phase Project**:
- Don't need phase protocols
- Don't need architect reviews

**Small Team**:
- Quick references may be overkill
- Comprehensive guides not needed

**Already Familiar with System**:
- Examples not needed
- Diagrams not needed

## How to Activate Files

1. **Review the file** to understand its purpose
2. **Customize** for your project (remove KCP/TMC references)
3. **Move to appropriate directory**:
   ```bash
   # Move to protocols
   mv ARCHITECT-REVIEWER-WAVE-INSTRUCTIONS.md ../protocols/
   
   # Move to core
   mv ORCHESTRATOR-QUICK-REFERENCE.md ../core/
   ```
4. **Update agent configurations** to reference the file
5. **Add to CLAUDE.md** reading requirements if critical

## File Descriptions

### CODE-REVIEWER-QUICK-REFERENCE.md
- Decision trees for review outcomes
- Common issues checklist
- Size violation handling
- Quick commands

### ORCHESTRATOR-QUICK-REFERENCE.md
- State transition table
- Agent spawning templates
- Common commands
- Troubleshooting guide

### ARCHITECT-REVIEWER-WAVE-INSTRUCTIONS.md
- Complete architect review process
- Assessment criteria
- Decision framework
- Integration requirements

### PHASE-START-ARCHITECT-REVIEW-PROTOCOL.md
- Phase readiness assessment
- Feature completeness checks
- Go/no-go decisions
- Course correction procedures

### TODO-STATE-MANAGEMENT-PROTOCOL.md
- TODO file formats
- Save/load procedures
- Cleanup rules
- Recovery procedures

### SPLIT-EXAMPLE-E3.1.1-SYNC-ENGINE.md
- Real-world split of 2400 lines → 3 parts
- Shows decision process
- Demonstrates file organization
- Includes review outcomes

## Important Notes

1. **These files contain battle-tested wisdom** from production use
2. **Many have KCP/TMC specific references** that need updating
3. **They provide valuable patterns** even if not directly used
4. **Keep as reference** even if not actively using

## Recommendation

Start with the minimal set in the main template. As your project grows and you encounter specific needs, selectively move files from here to active use. The examples and diagrams are particularly valuable for training new team members on the system.