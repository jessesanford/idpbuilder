# 🚨🚨🚨 BLOCKING RULE R239: Fix Plan Distribution Protocol

## Criticality: BLOCKING
**Failure to distribute fix plans properly = -75% GRADE**

## Description
Fix plans created by Code Reviewer MUST be distributed to individual effort directories for engineers to access.

## Requirements

### 1. Fix Plan Creation (Code Reviewer)
Code Reviewer in CREATE_INTEGRATION_FIX_PLAN state MUST:
- Analyze integration report
- Create specific fix plan for each failed effort
- Output to fix-plans directory
- Create FIX_PLAN_SUMMARY.yaml

### 2. Distribution Process (Orchestrator)
In DISTRIBUTE_FIX_PLANS state, orchestrator MUST:

```bash
for each effort in fix_plans:
    # Copy fix plan to effort directory
    cp fix-plans/FIX_PLAN_${effort}.md \
       efforts/phase${PHASE}/wave${WAVE}/${effort}/INTEGRATION_FIX_PLAN_${TIMESTAMP}.md
    
    # Create marker files
    echo "INTEGRATION_FIX_REQUIRED" > ${effort_dir}/FIX_REQUIRED.flag
    echo "${fix_plan_path}" > ${effort_dir}/FIX_PLAN_LOCATION.txt
    
    # Commit to effort branch
    cd ${effort_dir}
    git add INTEGRATION_FIX_PLAN_*.md FIX_REQUIRED.flag FIX_PLAN_LOCATION.txt
    git commit -m "fix-plan: Integration fixes required"
    git push
```

### 3. Marker Files (MANDATORY)
Each effort directory receiving fixes MUST have:
- `FIX_REQUIRED.flag` - Signals engineer that fixes are needed
- `FIX_PLAN_LOCATION.txt` - Contains path to fix plan
- `INTEGRATION_FIX_PLAN_*.md` - The actual fix plan

### 4. State File Updates
Orchestrator MUST record:
- Which efforts received fix plans
- Timestamp of distribution
- Location of distribution log

## Violations

### AUTOMATIC FAILURE (-75%)
- Not distributing fix plans to effort directories
- Missing marker files
- Wrong directory structure

### MAJOR VIOLATIONS (-50%)
- Not committing to effort branches
- Missing push to remote
- Incomplete distribution

## Implementation Example

```bash
distribute_fix_plan() {
    local effort="$1"
    local plan_file="$2"
    local effort_dir="efforts/phase${PHASE}/wave${WAVE}/${effort}"
    
    # Distribute
    cp "$plan_file" "${effort_dir}/INTEGRATION_FIX_PLAN_$(date +%Y%m%d-%H%M%S).md"
    
    # Mark
    echo "INTEGRATION_FIX_REQUIRED" > "${effort_dir}/FIX_REQUIRED.flag"
    
    # Commit
    cd "$effort_dir"
    git add -A
    git commit -m "fix-plan: Distributed from Code Reviewer"
    git push
}
```

## Related Rules
- R238: Integration Report Evaluation Protocol
- R209: Effort Directory Isolation Protocol
- R194: Remote Branch Tracking
- R300: Comprehensive Fix Management Protocol

## Grading Impact
- **Correct distribution**: +15% compliance bonus
- **Missing distribution**: -75% major failure
- **Partial distribution**: -30% violation