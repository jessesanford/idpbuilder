# 📋 SOFTWARE FACTORY MANAGER - SPARSE CLONE INVESTIGATION REPORT

**Date**: 2025-08-28  
**Requested By**: User  
**Investigator**: @agent-software-factory-manager  

## 🔍 Executive Summary

After comprehensive search and analysis of the Software Factory 2.0 template, I've found that **sparse clones have been EXPLICITLY FORBIDDEN and SUPERSEDED**. The system now mandates **FULL single-branch checkouts** for all effort directories.

## 🔴🔴🔴 KEY FINDING: R271 SUPREME LAW 🔴🔴🔴

### R193 (Original Rule) - NOW SUPERSEDED
- **Title**: Effort Clone Protocol
- **Status**: SUPERSEDED by R271 on 2025-08-28
- **Original requirement**: Orchestrator creates sparse git clones for efforts
- **Location**: `/rule-library/R193-effort-clone-protocol.md`

### R271 (Current Supreme Law) - ACTIVE
- **Title**: Single-Branch Full Checkout Protocol (SUPREME LAW)
- **Criticality**: SUPREME - Supersedes R193
- **Location**: `/rule-library/R271-single-branch-full-checkout.md`

## 📊 Critical Requirements Summary

### MANDATORY Protocol (R271):

1. **THINK About Base Branch First**
   - Analyze effort dependencies
   - Determine correct base branch
   - Document selection rationale

2. **SINGLE-BRANCH Clone Only**
   ```bash
   git clone \
       --single-branch \
       --branch "$BASE_BRANCH" \
       "$TARGET_REPO_URL" \
       "$EFFORT_DIR"
   ```

3. **FULL Code Checkout**
   - NO sparse checkout allowed
   - Agents need complete codebase context
   - All files must be available

4. **FORBIDDEN Practices**:
   - ❌ `git sparse-checkout init`
   - ❌ `git sparse-checkout set`
   - ❌ Partial directory checkouts
   - ❌ Limiting agent visibility

## 🏭 Orchestrator Setup Infrastructure Protocol

### From SETUP_EFFORT_INFRASTRUCTURE State Rules:

```bash
prepare_effort_for_agent() {
    local PHASE=$1 WAVE=$2 EFFORT=$3
    
    # 1. THINK - Determine base branch
    echo "🧠 THINKING: What base branch should $EFFORT use?"
    
    # 2. Create effort directory
    EFFORT_DIR="/efforts/phase${PHASE}/wave${WAVE}/${EFFORT}"
    mkdir -p "$(dirname "$EFFORT_DIR")"
    
    # 3. SINGLE-BRANCH FULL CLONE (R271 Supreme Law)
    git clone \
        --single-branch \
        --branch "$BASE_BRANCH" \
        "$TARGET_REPO_URL" \
        "$EFFORT_DIR"
    
    # 4. Create effort branch
    cd "$EFFORT_DIR"
    git checkout -b "$BRANCH_NAME"
    git push -u origin "$BRANCH_NAME"
    
    # 5. Verify FULL checkout (NO sparse)
    if [ -f ".git/info/sparse-checkout" ]; then
        echo "🔴🔴🔴 SUPREME LAW VIOLATION: Sparse checkout detected!"
        exit 1
    fi
}
```

## 🔧 Supporting Infrastructure Rules

### R208: Orchestrator Working Directory Spawn Protocol
- Orchestrator MUST change to effort directory before spawning agents
- Agents inherit working directory from spawner
- Critical for workspace isolation

### R209: Effort Directory Isolation Protocol
- SW Engineers work EXCLUSIVELY in assigned effort directory
- Metadata injection into implementation plans
- Directory boundary enforcement

### R251: Universal Repository Separation Law
- Software Factory instance for planning only
- Target repository clones under `/efforts/`
- Complete separation of concerns

## 📁 Directory Structure

```
software-factory-instance/              # Planning & Orchestration
└── efforts/                           # Target repo clones (FULL)
    └── phase1/
        └── wave1/
            ├── api-types/              # Full clone 1
            │   ├── .git/
            │   ├── pkg/                # ALL code visible
            │   ├── go.mod
            │   └── Makefile
            ├── controllers/            # Full clone 2
            │   ├── .git/
            │   ├── pkg/                # ALL code visible
            │   └── [complete codebase]
            └── webhooks/               # Full clone 3
                └── [complete codebase]
```

## 🛠️ Utility Support

### `/utilities/setup-effort-workspace.sh`
- Implements R271 compliance
- Creates FULL single-branch clones
- Verifies no sparse checkout
- Documents base branch decisions

## 📋 Grading Impact

### R271 Violations:
- Using sparse checkout: **-100% (SUPREME LAW VIOLATION)**
- Not documenting base branch: -30%
- Wrong base branch selected: -50%
- Partial code given to agents: -80%

### R193 (Superseded):
- Original penalties no longer apply
- Must follow R271 instead

## ✅ Compliance Checklist

For EVERY effort setup, orchestrator MUST:

1. ☐ THINK about appropriate base branch
2. ☐ Document base branch decision
3. ☐ Create effort directory under `/efforts/`
4. ☐ Clone with `--single-branch` flag ONLY
5. ☐ Use FULL checkout (no sparse)
6. ☐ Verify no `.git/info/sparse-checkout` file
7. ☐ Create and push effort branch
8. ☐ Inject R209 metadata into plans
9. ☐ Change to effort directory before spawning (R208)
10. ☐ Return to orchestrator directory after spawn

## 🎯 Why This Change Matters

### Problems with Sparse Checkout:
1. **Limited Context**: Agents can't see related code
2. **Import Failures**: Missing dependencies
3. **Test Failures**: Tests may touch any part of codebase
4. **Build Issues**: Build systems expect complete tree

### Benefits of Full Checkout:
1. **Complete Context**: Agents see entire codebase
2. **Proper Dependencies**: All imports available
3. **Testing Works**: Full test suite can run
4. **Integration Ready**: Can validate entire application

## 📝 Summary for Orchestrator

**NEVER USE SPARSE CHECKOUT** - This is now a SUPREME LAW violation per R271.

Every effort gets:
- A FULL clone of the target repository
- From an appropriate base branch (with documented reasoning)
- In an isolated directory under `/efforts/`
- With complete code visibility for agents

The transition from R193 (sparse) to R271 (full) represents a fundamental shift in how the Software Factory manages working copies. This ensures agents have the complete context they need for successful implementation.

---

**Report Generated**: 2025-08-28 06:09:00 UTC  
**Factory Manager**: @agent-software-factory-manager  
**Status**: Investigation Complete