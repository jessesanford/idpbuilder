# R530: Cascade Terminology Disambiguation Guide

**Rule ID**: R530
**Type**: REFERENCE
**Severity**: INFORMATIONAL
**Applies to**: All Agents
**Introduced**: 2025-10-06

## Purpose

The Software Factory uses the term "cascade" in multiple contexts. This guide disambiguates the three distinct meanings to prevent agent confusion.

## The Three Types of Cascade

### 1. Integration Cascade (R406, R348, R350, R351)

**Definition**: The process of recreating dependent integrations when an upstream bug is fixed.

**Example**:
- Bug found in Wave 1 integration
- Bug fixed in effort branch
- Wave 1 integration MUST be recreated (rebase)
- Wave 2 integration depends on Wave 1 → MUST be recreated
- Phase 1 integration depends on Wave 2 → MUST be recreated
- This chain is the "integration cascade"

**Metadata**:
- `cascade_source` - Bug that triggered cascade
- `cascade_level` - Wave/phase/project level
- `cascade_dependencies` - What depends on this

**Rules**: R406 (tracking), R348 (state), R350 (calculation), R351 (execution)

**State Machine**: `CASCADE_REINTEGRATION` state

**When to use this term**: When discussing the rebuild/recreation of dependent integrations

---

### 2. Bug Status Propagation (R524)

**Definition**: The process of updating duplicate bug statuses when the canonical bug is resolved.

**Example**:
- BUG-001 is canonical
- BUG-006, BUG-007 are duplicates of BUG-001
- BUG-001 gets fixed
- BUG-006 and BUG-007 statuses automatically update to "FIXED_AS_DUPLICATE"
- This status update flow is "bug status propagation"

**Metadata**:
- `is_duplicate` - Boolean flag
- `duplicate_of` - Canonical bug reference
- `duplicates` - Array of duplicate bugs

**Rules**: R524 (propagation protocol)

**When to use this term**: When discussing status updates flowing from canonical to duplicate bugs

**Historical Note**: Previously called "bug resolution cascade" but renamed to avoid confusion with integration cascade

---

### 3. Layered Cascade (R410)

**Definition**: Multiple integration cascades occurring in succession due to fixes at different levels.

**Example**:
- Cascade 1: Fix BUG-001 → rebuild Wave 1, Wave 2
- During Wave 2 rebuild: BUG-008 discovered
- Cascade 2: Fix BUG-008 → rebuild Wave 2, Phase 1
- During Phase 1 rebuild: BUG-015 discovered
- Cascade 3: Fix BUG-015 → rebuild Phase 1, Project
- This succession of cascades is a "layered cascade"

**Metadata**:
- `cascade_iteration` - Which cascade iteration (1st, 2nd, 3rd)
- `cascade_depth` - How many layers deep

**Rules**: R410 (layered cascade management)

**When to use this term**: When discussing multiple integration cascades happening sequentially

---

## Terminology Quick Reference

| Concept | Term to Use | Rules | What It Affects |
|---------|-------------|-------|-----------------|
| Rebuilding dependent integrations | **Integration Cascade** | R406, R348, R350, R351 | Branches, merges, rebases |
| Updating duplicate bug statuses | **Bug Status Propagation** | R524 | Bug registry metadata |
| Multiple cascades in succession | **Layered Cascade** | R410 | Cascade iteration count |

## Usage Guidelines

### ✅ DO

- Use "integration cascade" when discussing R406/R348/R350/R351
- Use "bug status propagation" when discussing R524
- Use "layered cascade" when discussing R410
- Use "cascade" alone only when context is crystal clear

### ❌ DON'T

- Use "cascade" without qualifier when ambiguity exists
- Mix cascade types in the same discussion without clarification
- Use "bug resolution cascade" (deprecated, use "bug status propagation")

## Agent Responsibilities

**All agents MUST**:
1. Read this guide before working with cascade-related rules
2. Use precise terminology in bug reports, commit messages, and agent communications
3. Ask for clarification if "cascade" is used ambiguously
4. Reference this guide (R530) when documenting cascade-related work

## Cross-References

- **Integration Cascade**: R406, R348, R350, R351, R410
- **Bug Status Propagation**: R521, R522, R523, R524, R525
- **Orchestrator State**: `bug_registry.cascade_*` fields (integration), `bug_registry.is_duplicate` fields (propagation)

## Examples

### Example 1: Clear Usage ✅

> "After fixing BUG-001, we need to execute the **integration cascade** (R350) to rebuild Wave 2 and Phase 1. Once the integration cascade completes successfully, **bug status propagation** (R524) will automatically mark the duplicate bugs as fixed."

**Why it's good**: Uses precise terms, distinguishes the two systems clearly.

### Example 2: Ambiguous Usage ❌

> "After fixing BUG-001, we need to cascade the changes and cascade the bug status."

**Why it's bad**: Uses "cascade" twice for different meanings without qualification.

### Example 3: Corrected Usage ✅

> "After fixing BUG-001, we need to execute the **integration cascade** (R350/R351) to propagate the code changes to dependent integrations, and then **bug status propagation** (R524) will update the duplicate bugs."

**Why it's good**: Clarifies both types, references specific rules.

---

## Revision History

- **2025-10-06**: Initial version (R530) - Disambiguates three cascade types
