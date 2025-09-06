# FILE NAMING COLLISION ANALYSIS REPORT

**Date**: 2025-09-01
**Analyst**: Software Factory Manager
**Status**: INVESTIGATION COMPLETE

## EXECUTIVE SUMMARY

After comprehensive analysis of the Software Factory 2.0 template, I found that **file naming collisions are mostly prevented by the directory structure**, but there are **specific scenarios where collisions could occur** during integration and parallel work.

## KEY FINDINGS

### 1. CURRENT STATE: Partial Protection

#### Files Created in Effort Directories (NO COLLISION RISK)
Each effort works in its own directory structure:
- `efforts/phase1/wave1/auth-system/work-log.md`
- `efforts/phase1/wave1/user-mgmt/work-log.md`
- `efforts/phase1/wave1/api-gateway/work-log.md`

These are **separate directories**, so no collision occurs for standard effort files.

#### Files with Timestamps (PROTECTED)
Some rules already use timestamps:
- **R287**: TODO files use `${AGENT_NAME}-${STATE}-${YYYYMMDD-HHMMSS}.todo`
- **R293**: Archived files use `SPLIT-PLAN-COMPLETED-$(date +%Y%m%d-%H%M%S).md`
- **R239**: Fix plans use `INTEGRATION_FIX_PLAN_$(date +%Y%m%d-%H%M%S).md`

### 2. COLLISION RISK SCENARIOS

#### Scenario 1: Integration Workspace Files
When multiple integration agents work in the same integration workspace:
- `/efforts/phase1/integration-workspace/INTEGRATION-REPORT.md`
- `/efforts/phase1/integration-workspace/work-log.md`

**Risk**: If multiple integration attempts happen, files could be overwritten.

#### Scenario 2: Wave/Phase Level Plans
Files created at wave or phase level without effort context:
- `/efforts/phase1/wave1/WAVE-ARCHITECTURE-PLAN.md`
- `/efforts/phase1/wave1/WAVE-IMPLEMENTATION-PLAN.md`

**Risk**: Multiple architects or planners could overwrite each other's work.

#### Scenario 3: Split Work Within Efforts
When efforts are split, multiple SW engineers might create:
- `efforts/phase1/wave1/auth-system/split-1/work-log.md`
- `efforts/phase1/wave1/auth-system/split-2/work-log.md`

**Current Status**: PROTECTED by split directory structure.

#### Scenario 4: Review Reports
Code reviewers create reports that might collide:
- `CODE-REVIEW-REPORT.md` (if created in same directory)
- `SPLIT-PLAN.md` (if multiple splits planned)

**Partial Protection**: R293 archives old plans with timestamps.

## EXISTING PROTECTIONS

### Rules with Timestamp Requirements:
1. **R264** - Work Log Tracking: Requires timestamps in operation headers
2. **R287** - TODO Persistence: Requires timestamp in filename
3. **R293** - Integration Report Distribution: Archives with timestamps
4. **R239** - Fix Plan Distribution: Uses timestamps in filenames

### Directory-Based Protection:
1. **R208** - Spawn Directory Protocol: Creates separate effort directories
2. **R197** - One Agent Per Effort: Ensures isolation
3. **R271** - Single Branch Full Checkout: Separate workspaces

## GAPS IDENTIFIED

### 1. No Universal File Naming Rule
There's no single rule mandating timestamps/suffixes for ALL generated files.

### 2. Integration Workspace Collisions
Integration workspace files lack systematic naming protection.

### 3. Inconsistent Timestamp Formats
Different rules use different timestamp formats:
- `YYYYMMDD-HHMMSS`
- `%Y%m%d-%H%M%S`
- `YYYYMMDD`

### 4. No Collision Detection Mechanism
No rule requires checking for existing files before creation.

## RECOMMENDATIONS

### Option 1: Create R301 - Universal File Naming Collision Prevention
A new rule requiring ALL generated files to include:
- Effort identifier (when applicable)
- Timestamp (for uniqueness)
- Agent type (for clarity)

Example: `WORK-LOG-auth-system-swe-20250120-143000.md`

### Option 2: Enhance Existing Rules
Update specific high-risk rules:
- R263 - Integration Documentation: Add timestamp requirements
- R264 - Work Log Tracking: Add filename timestamp requirement
- R219 - Code Reviewer Planning: Add timestamp to plan files

### Option 3: Directory-Based Solution
Mandate unique directories for each operation:
- `integration-20250120-143000/INTEGRATION-REPORT.md`
- `review-20250120-145500/CODE-REVIEW-REPORT.md`

## RISK ASSESSMENT

**Current Risk Level**: MEDIUM
- Most common cases are protected by directory structure
- Some edge cases could cause overwrites
- Integration and review phases are most vulnerable

**Impact if Collision Occurs**:
- Lost work from overwrites
- Merge conflicts during integration
- Confusion about which file is authoritative
- Potential grading penalties for lost documentation

## DECISION REQUIRED

Should we:
1. **CREATE R301** - New comprehensive file naming rule
2. **ENHANCE EXISTING** - Update current rules with naming requirements
3. **ACCEPT CURRENT STATE** - Rely on directory structure protection
4. **HYBRID APPROACH** - R301 for high-risk files only

## FILES AT HIGHEST RISK

Priority for protection:
1. **INTEGRATION-REPORT.md** - Created during integration
2. **work-log.md** - Created by multiple agents
3. **CODE-REVIEW-REPORT.md** - Created by reviewers
4. **WAVE-IMPLEMENTATION-PLAN.md** - Wave-level planning
5. **PHASE-ARCHITECTURE-PLAN.md** - Phase-level planning

## CONCLUSION

While the current system provides **good protection through directory isolation**, there are **specific scenarios where file collisions could occur**. The highest risk is during integration and review phases where multiple agents might work in shared directories.

**Recommendation**: Create R301 as a targeted rule for high-risk files rather than universal application, focusing on integration, review, and planning documents that exist outside effort directories.