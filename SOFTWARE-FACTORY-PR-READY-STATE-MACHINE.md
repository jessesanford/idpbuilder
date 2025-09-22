# Software Factory 2.0 - PR-Ready State Machine
## Automated Transformation of Effort Branches to Production-Ready Pull Request Branches

---

## 🔴🔴🔴 SUPREME DIRECTIVE: CLEAN UPSTREAM PR BRANCHES 🔴🔴🔴

This state machine transforms Software Factory effort branches into production-ready branches suitable for upstream pull requests by:
- **REMOVING** all Software Factory artifacts and metadata
- **PRESERVING** all core application functionality
- **CONSOLIDATING** commits into clean, atomic changes
- **VALIDATING** merge compatibility with upstream
- **ENSURING** no destructive changes reach production

---

## 🚨🚨🚨 CRITICAL CONSTRAINTS 🚨🚨🚨

### MANDATORY AGENT ASSIGNMENTS
- **Orchestrator**: Coordination ONLY - NEVER writes code or performs operations
- **Integration Agent**: ALL branch operations, rebasing, merging, conflict resolution
- **SW Engineer**: ALL file modifications, cleanups, artifact removal
- **Code Reviewer**: Validation, verification, assessment of cleanliness

### QUALITY GATES (BLOCKING)
1. **ZERO SF ARTIFACTS**: Not a single Software Factory file may remain
2. **CORE INTEGRITY**: NO core application files may be deleted
3. **CLEAN HISTORY**: Commits must be properly consolidated
4. **MERGE COMPATIBILITY**: Branches must merge cleanly to upstream
5. **BUILD VERIFICATION**: Code must compile and pass tests

---

## 📊 STATE MACHINE OVERVIEW

### Primary Flow
```
PR_READY_INIT
  → PR_DISCOVERY_ASSESSMENT
  → PR_SPAWN_DISCOVERY_AGENTS
  → PR_MONITOR_DISCOVERY
  → PR_CLEANUP_PLANNING
  → PR_SPAWN_CLEANUP_AGENTS
  → PR_MONITOR_CLEANUP
  → PR_CONSOLIDATION_PLANNING
  → PR_SPAWN_CONSOLIDATION_AGENTS
  → PR_MONITOR_CONSOLIDATION
  → PR_INTEGRITY_VERIFICATION
  → PR_SPAWN_VERIFICATION_AGENTS
  → PR_MONITOR_VERIFICATION
  → PR_SEQUENTIAL_REBASE_PLANNING
  → PR_SPAWN_REBASE_AGENTS
  → PR_MONITOR_REBASE
  → PR_VALIDATION_TESTING
  → PR_SPAWN_VALIDATION_AGENTS
  → PR_MONITOR_VALIDATION
  → PR_FINAL_PREPARATION
  → PR_READY_SUCCESS
```

### Error Recovery Flow
```
Any State → PR_ERROR_DETECTED
  → PR_ANALYZE_ERROR
  → PR_SPAWN_RECOVERY_AGENTS
  → PR_MONITOR_RECOVERY
  → [Return to failed state or PR_READY_ABORT]
```

---

## 🎯 PHASE 1: DISCOVERY AND ASSESSMENT

### STATE: PR_READY_INIT
**Owner**: Orchestrator
**Purpose**: Initialize PR-ready transformation process
**Actions**:
- Load PR-ready transformation configuration
- Read branch transformation requirements
- Initialize PR-ready state file
- Set up transformation workspace

**Transitions**:
- Success → PR_DISCOVERY_ASSESSMENT
- Error → PR_ERROR_DETECTED

### STATE: PR_DISCOVERY_ASSESSMENT
**Owner**: Orchestrator
**Purpose**: Plan discovery of SF artifacts and branch analysis
**Actions**:
- Create discovery task list
- Identify branches to transform
- Plan contamination assessment
- Prepare discovery agent instructions

**Transitions**:
- Success → PR_SPAWN_DISCOVERY_AGENTS
- Error → PR_ERROR_DETECTED

### STATE: PR_SPAWN_DISCOVERY_AGENTS
**Owner**: Orchestrator
**Purpose**: Spawn agents to perform discovery
**Actions**:
- Spawn Integration Agent for branch inventory
- Spawn SW Engineer for artifact scanning
- Spawn Code Reviewer for contamination assessment
- Record spawned agents in state file

**Transitions**:
- Success → PR_MONITOR_DISCOVERY
- Error → PR_ERROR_DETECTED

### STATE: PR_MONITOR_DISCOVERY
**Owner**: Orchestrator
**Purpose**: Monitor discovery progress
**Actions**:
- Track agent discovery progress
- Collect artifact inventory reports
- Compile branch dependency analysis
- Aggregate contamination findings

**Transitions**:
- All complete → PR_CLEANUP_PLANNING
- Error → PR_ERROR_DETECTED
- Timeout → PR_ERROR_DETECTED

---

## 🧹 PHASE 2: CLEANUP AND SANITIZATION

### STATE: PR_CLEANUP_PLANNING
**Owner**: Orchestrator
**Purpose**: Plan artifact removal strategy
**Actions**:
- Review discovered artifacts
- Create cleanup task assignments
- Prioritize branch cleanup order
- Define rollback checkpoints

**Transitions**:
- Success → PR_SPAWN_CLEANUP_AGENTS
- No artifacts found → PR_CONSOLIDATION_PLANNING
- Error → PR_ERROR_DETECTED

### STATE: PR_SPAWN_CLEANUP_AGENTS
**Owner**: Orchestrator
**Purpose**: Spawn agents for artifact removal
**Actions**:
- Spawn SW Engineer(s) for artifact removal
- Spawn Integration Agent for main reset (if needed)
- Assign specific branches to each agent
- Record cleanup assignments

**Transitions**:
- Success → PR_MONITOR_CLEANUP
- Error → PR_ERROR_DETECTED

### STATE: PR_MONITOR_CLEANUP
**Owner**: Orchestrator
**Purpose**: Monitor cleanup progress
**Actions**:
- Track artifact removal progress
- Verify no core files deleted
- Collect cleanup completion reports
- Update transformation status

**Transitions**:
- All complete → PR_CONSOLIDATION_PLANNING
- Core files deleted → PR_ERROR_DETECTED
- Error → PR_ERROR_DETECTED

---

## 📝 PHASE 3: COMMIT CONSOLIDATION

### STATE: PR_CONSOLIDATION_PLANNING
**Owner**: Orchestrator
**Purpose**: Plan commit squashing strategy
**Actions**:
- Analyze commit history per branch
- Create consolidation instructions
- Define commit message templates
- Plan preservation of history

**Transitions**:
- Success → PR_SPAWN_CONSOLIDATION_AGENTS
- Already consolidated → PR_INTEGRITY_VERIFICATION
- Error → PR_ERROR_DETECTED

### STATE: PR_SPAWN_CONSOLIDATION_AGENTS
**Owner**: Orchestrator
**Purpose**: Spawn agents for commit consolidation
**Actions**:
- Spawn Integration Agent(s) for commit operations
- Assign branches for consolidation
- Provide commit message templates
- Set consolidation parameters

**Transitions**:
- Success → PR_MONITOR_CONSOLIDATION
- Error → PR_ERROR_DETECTED

### STATE: PR_MONITOR_CONSOLIDATION
**Owner**: Orchestrator
**Purpose**: Monitor consolidation progress
**Actions**:
- Track squashing progress
- Verify commit message quality
- Ensure history preservation
- Collect consolidation reports

**Transitions**:
- All complete → PR_INTEGRITY_VERIFICATION
- Error → PR_ERROR_DETECTED

---

## ✅ PHASE 4: INTEGRITY VERIFICATION

### STATE: PR_INTEGRITY_VERIFICATION
**Owner**: Orchestrator
**Purpose**: Plan integrity checks
**Actions**:
- Create verification checklist
- Define core file preservation rules
- Plan destructive change detection
- Prepare verification instructions

**Transitions**:
- Success → PR_SPAWN_VERIFICATION_AGENTS
- Error → PR_ERROR_DETECTED

### STATE: PR_SPAWN_VERIFICATION_AGENTS
**Owner**: Orchestrator
**Purpose**: Spawn agents for integrity verification
**Actions**:
- Spawn Code Reviewer for core file checks
- Spawn SW Engineer for change analysis
- Assign verification tasks
- Set verification thresholds

**Transitions**:
- Success → PR_MONITOR_VERIFICATION
- Error → PR_ERROR_DETECTED

### STATE: PR_MONITOR_VERIFICATION
**Owner**: Orchestrator
**Purpose**: Monitor verification results
**Actions**:
- Collect integrity reports
- Identify branches with issues
- Document destructive changes
- Determine remediation needs

**Transitions**:
- All pass → PR_SEQUENTIAL_REBASE_PLANNING
- Destructive changes found → PR_FIX_DESTRUCTIVE_CHANGES
- Error → PR_ERROR_DETECTED

### STATE: PR_FIX_DESTRUCTIVE_CHANGES
**Owner**: Orchestrator
**Purpose**: Coordinate fixes for destructive branches
**Actions**:
- Create fix plans for problematic branches
- Spawn SW Engineer to rebuild branches
- Monitor fix progress
- Re-verify after fixes

**Transitions**:
- Fixed → PR_SEQUENTIAL_REBASE_PLANNING
- Cannot fix → PR_READY_ABORT
- Error → PR_ERROR_DETECTED

---

## 🔄 PHASE 5: SEQUENTIAL REBASE

### STATE: PR_SEQUENTIAL_REBASE_PLANNING
**Owner**: Orchestrator
**Purpose**: Plan rebase sequence
**Actions**:
- Determine branch dependencies
- Create rebase order
- Identify potential conflicts
- Prepare conflict resolution strategies

**Transitions**:
- Success → PR_SPAWN_REBASE_AGENTS
- No rebase needed → PR_VALIDATION_TESTING
- Error → PR_ERROR_DETECTED

### STATE: PR_SPAWN_REBASE_AGENTS
**Owner**: Orchestrator
**Purpose**: Spawn agents for rebasing
**Actions**:
- Spawn Integration Agent for rebase operations
- Provide rebase sequence
- Set conflict resolution rules
- Define incremental change strategy

**Transitions**:
- Success → PR_MONITOR_REBASE
- Error → PR_ERROR_DETECTED

### STATE: PR_MONITOR_REBASE
**Owner**: Orchestrator
**Purpose**: Monitor rebase progress
**Actions**:
- Track rebase operations
- Monitor conflict resolution
- Verify branch relationships
- Collect rebase reports

**Transitions**:
- All complete → PR_VALIDATION_TESTING
- Unresolvable conflicts → PR_MANUAL_CONFLICT_RESOLUTION
- Error → PR_ERROR_DETECTED

### STATE: PR_MANUAL_CONFLICT_RESOLUTION
**Owner**: Orchestrator
**Purpose**: Coordinate manual conflict resolution
**Actions**:
- Document unresolvable conflicts
- Create manual resolution instructions
- Pause for human intervention
- Verify resolution completion

**Transitions**:
- Resolved → PR_VALIDATION_TESTING
- Cannot resolve → PR_READY_ABORT
- Error → PR_ERROR_DETECTED

---

## 🧪 PHASE 6: VALIDATION AND TESTING

### STATE: PR_VALIDATION_TESTING
**Owner**: Orchestrator
**Purpose**: Plan validation tests
**Actions**:
- Create merge test scenarios
- Define build verification steps
- Plan test execution sequence
- Prepare validation criteria

**Transitions**:
- Success → PR_SPAWN_VALIDATION_AGENTS
- Error → PR_ERROR_DETECTED

### STATE: PR_SPAWN_VALIDATION_AGENTS
**Owner**: Orchestrator
**Purpose**: Spawn agents for validation
**Actions**:
- Spawn Integration Agent for merge testing
- Spawn SW Engineer for build verification
- Spawn Code Reviewer for final assessment
- Set validation parameters

**Transitions**:
- Success → PR_MONITOR_VALIDATION
- Error → PR_ERROR_DETECTED

### STATE: PR_MONITOR_VALIDATION
**Owner**: Orchestrator
**Purpose**: Monitor validation results
**Actions**:
- Track merge test results
- Monitor build status
- Collect validation reports
- Assess overall readiness

**Transitions**:
- All pass → PR_FINAL_PREPARATION
- Tests fail → PR_FIX_VALIDATION_ISSUES
- Error → PR_ERROR_DETECTED

### STATE: PR_FIX_VALIDATION_ISSUES
**Owner**: Orchestrator
**Purpose**: Coordinate validation issue fixes
**Actions**:
- Analyze test failures
- Create fix assignments
- Spawn appropriate agents
- Re-run validation after fixes

**Transitions**:
- Fixed → PR_VALIDATION_TESTING
- Cannot fix → PR_READY_ABORT
- Error → PR_ERROR_DETECTED

---

## 🎁 PHASE 7: FINAL PREPARATION

### STATE: PR_FINAL_PREPARATION
**Owner**: Orchestrator
**Purpose**: Prepare final deliverables
**Actions**:
- Create PR documentation
- Generate merge order instructions
- Document known conflicts and resolutions
- Push all branches to origin
- Create transformation report

**Transitions**:
- Success → PR_READY_SUCCESS
- Error → PR_ERROR_DETECTED

### STATE: PR_READY_SUCCESS
**Owner**: Orchestrator
**Purpose**: Successfully complete PR transformation
**Actions**:
- Mark transformation complete
- Archive transformation logs
- Generate success report
- Update branch status tracking

**Terminal State**: Process complete

---

## 🚨 ERROR RECOVERY STATES

### STATE: PR_ERROR_DETECTED
**Owner**: Orchestrator
**Purpose**: Handle detected errors
**Actions**:
- Capture error details
- Determine error severity
- Identify affected branches
- Create error report

**Transitions**:
- Recoverable → PR_ANALYZE_ERROR
- Non-recoverable → PR_READY_ABORT

### STATE: PR_ANALYZE_ERROR
**Owner**: Orchestrator
**Purpose**: Analyze error and plan recovery
**Actions**:
- Determine root cause
- Identify recovery strategy
- Create recovery plan
- Prepare recovery instructions

**Transitions**:
- Recovery planned → PR_SPAWN_RECOVERY_AGENTS
- Cannot recover → PR_READY_ABORT

### STATE: PR_SPAWN_RECOVERY_AGENTS
**Owner**: Orchestrator
**Purpose**: Spawn agents for error recovery
**Actions**:
- Spawn appropriate recovery agents
- Provide recovery instructions
- Set recovery parameters
- Monitor recovery initiation

**Transitions**:
- Success → PR_MONITOR_RECOVERY
- Error → PR_READY_ABORT

### STATE: PR_MONITOR_RECOVERY
**Owner**: Orchestrator
**Purpose**: Monitor recovery progress
**Actions**:
- Track recovery operations
- Verify issue resolution
- Assess recovery success
- Update state for retry

**Transitions**:
- Recovered → [Return to appropriate state]
- Failed → PR_READY_ABORT

### STATE: PR_READY_ABORT
**Owner**: Orchestrator
**Purpose**: Abort transformation with rollback
**Actions**:
- Document abort reason
- Rollback any partial changes
- Create abort report
- Preserve diagnostic information

**Terminal State**: Process aborted

---

## 📋 VALIDATION GATES

### Gate 1: Post-Discovery
**Required Before**: PR_CLEANUP_PLANNING
**Validations**:
- All branches inventoried
- All SF artifacts identified
- Dependency tree complete
- Contamination assessment done

### Gate 2: Post-Cleanup
**Required Before**: PR_CONSOLIDATION_PLANNING
**Validations**:
- Zero SF artifacts remain
- No core files deleted
- All branches cleaned
- Cleanup verified

### Gate 3: Post-Consolidation
**Required Before**: PR_INTEGRITY_VERIFICATION
**Validations**:
- All commits consolidated
- History preserved in messages
- Commit messages follow standards
- No orphaned commits

### Gate 4: Post-Verification
**Required Before**: PR_SEQUENTIAL_REBASE_PLANNING
**Validations**:
- Core files intact
- No destructive changes
- Change counts reasonable
- All branches verified

### Gate 5: Post-Rebase
**Required Before**: PR_VALIDATION_TESTING
**Validations**:
- Branches properly sequenced
- Dependencies resolved
- Conflicts documented
- Incremental changes clean

### Gate 6: Post-Validation
**Required Before**: PR_FINAL_PREPARATION
**Validations**:
- All merges clean
- Build successful
- Tests passing
- No regressions

---

## 🔄 STATE TRANSITION RULES

### Fundamental Laws
1. **R206 Compliance**: All states must be validated against this file
2. **R233 Compliance**: Every state requires immediate action upon entry
3. **R313 Compliance**: Orchestrator must stop after spawning agents
4. **Single Atomic Operation**: Each state performs ONE operation then stops

### Transition Requirements
- **State File Update**: MANDATORY before every transition
- **TODO Persistence**: Save TODOs per R287 before transitions
- **Error Capture**: All errors trigger PR_ERROR_DETECTED
- **Progress Tracking**: Document completion percentage

### Rollback Capability
Each phase maintains rollback points:
- Pre-cleanup branch snapshots
- Pre-consolidation commit SHAs
- Pre-rebase branch states
- Pre-merge test checkpoints

---

## 📊 METRICS AND REPORTING

### Success Metrics
- **Artifact Removal Rate**: 100% required
- **Core File Preservation**: 100% required
- **Commit Consolidation**: Target 1 commit per feature
- **Merge Success Rate**: >95% expected
- **Build Success**: 100% required

### Progress Tracking
- Branches processed / total
- Artifacts removed / discovered
- Commits consolidated / original
- Conflicts resolved / encountered
- Tests passed / total

### Reporting Requirements
- Real-time progress updates
- Phase completion summaries
- Error and recovery logs
- Final transformation report
- Merge order documentation

---

## 🚀 PARALLELIZATION OPPORTUNITIES

### Parallelizable Operations
- **Discovery**: Multiple branches can be scanned simultaneously
- **Cleanup**: Independent branches cleaned in parallel
- **Consolidation**: Non-dependent branches processed together
- **Verification**: Multiple branches verified concurrently

### Sequential Requirements
- **Rebasing**: Must follow dependency order
- **Merge Testing**: Sequential to detect conflicts
- **Conflict Resolution**: One at a time for clarity
- **Final Push**: Ordered to maintain relationships

---

## 🛡️ SAFETY MECHANISMS

### Critical Protections
1. **Core File Guard**: Never delete main.*, Makefile, README, LICENSE
2. **Mass Deletion Detection**: Alert if >10,000 lines deleted
3. **Commit Preservation**: Always preserve history in squashed commits
4. **Rollback Points**: Checkpoint before each destructive operation
5. **Dry Run Mode**: Optional validation without changes

### Abort Triggers
- Core application files deleted
- Unrecoverable merge conflicts
- Build failures after fixes
- Timeout on critical operations
- Manual abort request

---

## 📝 AGENT-SPECIFIC RESPONSIBILITIES

### Orchestrator Agent
- Coordinate all phases
- Spawn appropriate agents
- Monitor progress
- Handle state transitions
- Generate reports

### Integration Agent
- Perform all git operations
- Handle rebasing
- Resolve conflicts
- Test merges
- Push branches

### SW Engineer Agent
- Remove SF artifacts
- Fix file issues
- Rebuild problematic branches
- Verify builds
- Apply fixes

### Code Reviewer Agent
- Assess contamination
- Verify cleanup completeness
- Check core file integrity
- Validate transformations
- Final quality assessment

---

## State File Locations

Each PR-Ready state has corresponding rule files:

```
agent-states/
└── pr-ready/
    ├── orchestrator/
    │   ├── PR_READY_INIT/rules.md
    │   ├── PR_DISCOVERY_ASSESSMENT/rules.md
    │   ├── PR_CLEANUP_PLANNING/rules.md
    │   ├── PR_CONSOLIDATION_PLANNING/rules.md
    │   ├── PR_SEQUENTIAL_REBASE_PLANNING/rules.md
    │   ├── PR_VALIDATION_TESTING/rules.md
    │   ├── PR_FINAL_PREPARATION/rules.md
    │   ├── PR_READY_SUCCESS/rules.md
    │   └── ... (30 states total)
    ├── integration/
    │   ├── PR_BRANCH_INVENTORY/rules.md
    │   ├── PR_BRANCH_PUSH/rules.md
    │   ├── PR_COMMIT_SQUASH/rules.md
    │   ├── PR_CONFLICT_RESOLUTION/rules.md
    │   ├── PR_MAIN_RESET/rules.md
    │   ├── PR_MERGE_TEST/rules.md
    │   └── PR_REBASE_SEQUENCE/rules.md
    └── sw-engineer/
        ├── PR_ARTIFACT_SCAN/rules.md
        ├── PR_ARTIFACT_REMOVAL/rules.md
        ├── PR_BRANCH_REBUILD/rules.md
        ├── PR_BUILD_VERIFICATION/rules.md
        └── PR_CHANGE_ANALYSIS/rules.md
```

---

## Document Version
- **Version**: 1.0
- **Created**: 2025-09-21
- **Purpose**: Automate Software Factory effort branch to PR-ready transformation
- **Compatibility**: Software Factory 2.0 State Machine System