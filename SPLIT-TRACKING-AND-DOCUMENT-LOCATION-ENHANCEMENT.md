# Split Tracking and Document Location Enhancement Report

**Date**: 2025-09-02
**Manager**: software-factory-manager
**Mission**: Investigate and enhance split branch tracking and phase/wave document organization

## Executive Summary

This report documents the investigation and enhancement of split branch tracking and phase/wave document location protocols in Software Factory 2.0. Two new BLOCKING rules (R302 and R303) were created to address critical gaps in tracking and organization.

## Investigation Findings

### Current State Analysis

#### Split Tracking (Existing Rules)
- **R296**: Deprecated Branch Marking Protocol - Handles marking original branches as deprecated
- **R297**: Architect Split Detection Protocol - Ensures architects check split_count before measuring
- **R014**: Branch Naming Convention - Defines naming for split branches

#### Gaps Identified in Split Tracking
1. **No comprehensive tracking structure** for split relationships
2. **Missing split lifecycle states** (PLANNED, IN_PROGRESS, COMPLETED)
3. **No centralized split_tracking section** in orchestrator-state.json
4. **Incomplete tracking of split branch details** (lines, descriptions, statuses)
5. **No integration planning guidance** for using split branches

#### Document Locations (Existing Rules)
- **R210**: Architect creates architecture plans in `phase-plans/`
- **R218**: References to phase-plans for wave implementation plans
- **R213**: Wave metadata references phase-plans directory

#### Gaps Identified in Document Locations
1. **No comprehensive protocol** for ALL document types and locations
2. **Missing branch requirements** (which branch to commit to)
3. **Incomplete naming conventions** for all document types
4. **No clear separation** between phase/wave docs and effort docs

## Solutions Implemented

### R302: Comprehensive Split Tracking Protocol

**Purpose**: Provide meticulous tracking of all split operations

**Key Features**:
1. **Comprehensive Data Structure**:
   ```yaml
   split_tracking:
     effort-001:
       original_branch: "..."
       status: "SPLIT_DEPRECATED"
       split_count: 3
       split_branches:
         - branch: "...-split-001"
           status: "COMPLETED"
           lines: 450
           description: "..."
   ```

2. **Split Lifecycle Management**:
   - SPLIT_PLANNED → SPLIT_IN_PROGRESS → SPLIT_DEPRECATED
   - Individual split statuses: ACTIVE → COMPLETED → REVIEWED → INTEGRATED

3. **Integration Planning Support**:
   - Automatic branch selection for merges
   - Validation before integration
   - Clear replacement branch identification

4. **Architect Verification Integration**:
   - Direct support for R297 compliance
   - Simplified split detection queries

### R303: Phase/Wave Document Location Protocol

**Purpose**: Centralize phase/wave documents for accessibility and persistence

**Key Features**:
1. **Mandatory Directory Structure**:
   ```
   ${SF_INSTANCE_DIR}/
   ├── phase-plans/           # ALL phase/wave documents
   │   ├── PHASE-1-ARCHITECTURE-PLAN.md
   │   ├── PHASE-1-WAVE-1-IMPLEMENTATION-PLAN.md
   │   └── ...
   └── efforts/               # Only effort-specific documents
   ```

2. **Document Type Classification**:
   - **Phase-Level**: Architecture, Implementation, Assessment, Integration plans
   - **Wave-Level**: Architecture, Implementation, Merge, Review reports
   - **Effort-Level**: Implementation plans, Split plans, Review reports

3. **Git Branch Requirements**:
   - Phase/wave documents → main branch of SF instance
   - Effort documents → effort branches in target repo

4. **Strict Naming Conventions**:
   - Format: `PHASE-${N}-WAVE-${M}-${TYPE}-${SUFFIX}.md`
   - Enforced consistency across all agents

## Agent State Updates

### Updated States
1. **sw-engineer/SPLIT_IMPLEMENTATION**:
   - Added R302 split tracking requirements
   - Mandatory reporting of split details

2. **architect/PHASE_ARCHITECTURE_PLANNING**:
   - Added R303 document location requirements
   - Explicit phase-plans directory usage

3. **orchestrator/MONITOR**:
   - Added R302 split tracking monitoring
   - Update split_tracking section requirement

## Implementation Impact

### Benefits
1. **Meticulous Split Tracking**:
   - Never lose track of split relationships
   - Clear identification of current branches for integration
   - Prevents duplicate work or missed branches

2. **Organized Document Management**:
   - All phase/wave documents in one location
   - Easy discovery by all agents
   - Persistent across effort implementations

3. **Reduced Integration Errors**:
   - Automatic prevention of deprecated branch integration
   - Clear branch selection for merges
   - Validation before integration

### Grading Impact
- **R302 Violations**: -30% for missing tracking, -40% for wrong branches
- **R303 Violations**: -30% for wrong location, -20% for wrong branch

## Verification Methods

### Split Tracking Verification
```bash
# Check split tracking completeness
yq '.split_tracking | to_entries | length' orchestrator-state.json

# Find current splits for integration
yq '.split_tracking | to_entries | .[] | 
    select(.value.status == "SPLIT_DEPRECATED") | 
    .value.split_branches[].branch' orchestrator-state.json
```

### Document Location Verification
```bash
# Verify phase-plans structure
find ${SF_INSTANCE_DIR}/phase-plans -name "*.md" | 
    grep -E "^PHASE-[0-9]+-" | wc -l

# Check document naming compliance
find phase-plans -name "*.md" | 
    grep -vE "^PHASE-[0-9]+-" || echo "All compliant"
```

## Recommendations

### Immediate Actions
1. **Update orchestrator-state.json template** to include split_tracking section
2. **Create phase-plans directory** in all SF instances
3. **Train agents** on new tracking requirements

### Future Enhancements
1. **Automated split tracking updates** via helper scripts
2. **Document discovery service** for cross-agent access
3. **Split metrics dashboard** for monitoring

## Files Modified

### New Rules Created
- `/rule-library/R302-comprehensive-split-tracking-protocol.md`
- `/rule-library/R303-phase-wave-document-location-protocol.md`

### States Updated
- `/agent-states/sw-engineer/SPLIT_IMPLEMENTATION/rules.md`
- `/agent-states/architect/PHASE_ARCHITECTURE_PLANNING/rules.md`
- `/agent-states/orchestrator/MONITOR/rules.md`

### Registry Updated
- `/rule-library/RULE-REGISTRY.md` - Added R302 and R303

## Conclusion

The implementation of R302 and R303 addresses critical gaps in split tracking and document organization. These BLOCKING rules ensure meticulous tracking of all split operations and proper centralization of phase/wave documents. The changes support the Software Factory's goal of maintaining perfect consistency and traceability across all operations.

---
**Status**: COMPLETE
**Next Steps**: Commit and push all changes