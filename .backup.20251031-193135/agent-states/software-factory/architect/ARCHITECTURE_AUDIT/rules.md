# ARCHITECTURE_AUDIT State Rules

## SF 3.0 Architecture Context

This state performs architecture audit within SF 3.0:
- Reads current system state from `state_machine.current_state` in orchestrator-state-v3.json to understand project progression
- Reviews architectural decisions and patterns from previous phases/waves stored in orchestrator-state-v3.json metadata
- Assesses integration outcomes from `integration-containers.json` to identify architectural issues or technical debt
- Documents findings that may create `bug-tracking.json` entries if architectural violations or anti-patterns are detected
- All audit results are reported back for orchestrator to record in orchestrator-state-v3.json per R288

## Core Architecture Audit Rules

---
### ℹ️ RULE R077.0.0 - System-Wide Architecture Audit
**Source:** rule-library/RULE-REGISTRY.md#R077
**Criticality:** INFO - Best practice

MANDATE: Architect must perform comprehensive system audit
across all phases and components. Evaluate architectural
consistency, technical debt, and system health.

AUDIT SCOPE:
- Cross-phase architectural consistency
- Technical debt accumulation and impact
- KCP pattern adherence across entire system
- Performance characteristics at system scale
- Security posture and vulnerability assessment
- Maintainability and evolution readiness
---

---
### ℹ️ RULE R078.0.0 - Technical Debt Assessment and Management
**Source:** rule-library/RULE-REGISTRY.md#R078
**Criticality:** INFO - Best practice

MANDATE: Identify, categorize, and prioritize technical debt
accumulated across development phases. Provide remediation
strategy and impact analysis.

DEBT ASSESSMENT AREAS:
- Code quality degradation and anti-patterns
- Architecture drift from original design
- Performance regressions and scalability issues
- Security vulnerabilities and hardening gaps
- Documentation debt and knowledge gaps
- Test coverage gaps and quality issues
---

---
### ℹ️ RULE R079.0.0 - Architecture Evolution Assessment
**Source:** rule-library/RULE-REGISTRY.md#R079
**Criticality:** INFO - Best practice

MANDATE: Evaluate system's readiness for future evolution and
expansion. Assess extensibility, modularity, and adaptation
capabilities.

EVOLUTION CRITERIA:
- API design stability and backward compatibility
- Component modularity and loose coupling
- Extension points and plugin architectures
- Configuration flexibility and customization
- Deployment and operational scalability
- Integration readiness with external systems
---

---
### ℹ️ RULE R080.0.0 - System Health and Operational Readiness
**Source:** rule-library/RULE-REGISTRY.md#R080
**Criticality:** INFO - Best practice

MANDATE: Assess operational readiness including monitoring,
observability, debugging capabilities, and production
deployment characteristics.

OPERATIONAL CRITERIA:
- Monitoring and alerting coverage
- Logging and tracing capabilities
- Debugging and troubleshooting tools
- Performance profiling and optimization hooks
- Disaster recovery and backup capabilities
- Security monitoring and incident response
---

---
### ℹ️ RULE R037.0.0 - KCP Resource Pattern Enforcement
**Source:** rule-library/RULE-REGISTRY.md#R037
**Criticality:** INFO - Best practice

MANDATE: System-wide audit of KCP multi-tenant patterns to
ensure consistency and proper implementation across all
components and phases.

SYSTEM-LEVEL VALIDATION:
- Uniform LogicalCluster usage across all resource types
- Consistent workspace isolation patterns
- Compatible multi-tenancy implementations
- Coherent RBAC models across components
- Aligned event and status reporting strategies
---

## 🔴🔴🔴 MANDATORY MEASUREMENT RULE - R304 🔴🔴🔴

### R304: Mandatory Line Counter Tool Enforcement
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R304-mandatory-line-counter-enforcement.md`
**Criticality**: BLOCKING - Manual counting = AUTOMATIC -100% FAILURE

**ABSOLUTE REQUIREMENTS:**
- ✅ MUST use `$CLAUDE_PROJECT_DIR/tools/line-counter.sh` for ALL line counting
- ❌ NEVER use `wc -l` or any manual counting method
- ❌ NEVER count lines any other way - this is a -100% automatic failure
- ✅ MUST specify both -b (base branch) and -c (current branch) parameters
- ✅ Base branch MUST be phase integration branch (NOT "main")

**Failure to use the line counter tool = AUTOMATIC -100% GRADE**

## State-Specific Context

### ARCHITECTURE_AUDIT State Purpose
This state is entered for comprehensive system-wide architecture audits. This can occur at major milestones, before production deployment, or when architectural concerns arise across multiple phases or components.

### State Transitions

**ENTRY CONDITIONS:**
- Major milestone reached (e.g., multiple phases complete)
- Architecture review requested due to concerns
- Pre-production readiness assessment required
- Technical debt remediation planning needed

**EXIT CONDITIONS:**
- **ARCHITECTURE_HEALTHY**: System meets architectural standards → Continue with current approach
- **IMPROVEMENTS_REQUIRED**: Issues identified needing attention → Provide remediation plan
- **ARCHITECTURE_DEGRADED**: Significant problems requiring major work → Halt progression for refactoring

**STATE TRANSITION FLOW:**
```
ARCHITECTURE_AUDIT → [Comprehensive Assessment] → Decision State
├─ ARCHITECTURE_HEALTHY → Continue Current Development
├─ IMPROVEMENTS_REQUIRED → Implement Remediation Plan
└─ ARCHITECTURE_DEGRADED → Major Refactoring Required
```

## Architecture Audit Scope and Criteria

### Primary Audit Dimensions

| Audit Area | Weight | Assessment Method | Health Threshold |
|------------|--------|-------------------|------------------|
| Architectural Consistency | 20% | Cross-component pattern analysis | >90% consistency |
| Technical Debt Assessment | 20% | Code quality metrics + manual review | Low to Medium debt |
| Performance Characteristics | 15% | System-wide performance analysis | Meets scalability targets |
| Security Posture | 15% | Comprehensive security audit | Zero critical vulnerabilities |
| KCP Pattern Adherence | 15% | Multi-tenancy validation | 100% pattern compliance |
| Operational Readiness | 10% | Monitoring and ops capability review | Production-ready |
| Evolution Capability | 5% | Extensibility and modularity assessment | Future-ready |

### Architectural Consistency Assessment

**CONSISTENCY EVALUATION AREAS**:
- Design pattern usage across components
- API design and interaction patterns
- Data modeling and storage approaches
- Error handling and logging strategies
- Configuration and deployment patterns

**Cross-Component Analysis**:
- Identify architectural drift between phases
- Validate design decision consistency
- Assess component coupling and cohesion
- Review interface and contract stability
- Evaluate shared library and utility usage

### Technical Debt Assessment

**DEBT IDENTIFICATION CATEGORIES**:

| Debt Type | Impact | Assessment Method | Remediation Priority |
|-----------|--------|-------------------|-------------------|
| Code Quality Debt | Medium-High | Static analysis + review | High |
| Architecture Drift | High | Design comparison | Critical |
| Performance Debt | Medium | Benchmark analysis | High |
| Security Debt | Critical | Vulnerability assessment | Critical |
| Documentation Debt | Low-Medium | Coverage analysis | Medium |
| Test Debt | Medium | Coverage + quality review | High |

**Technical Debt Calculation**:
```yaml
debt_assessment:
  code_quality:
    complexity_score: 7.2/10       # Acceptable
    maintainability_index: 78      # Good
    duplication_percentage: 12%    # Moderate
    
  architecture_drift:
    pattern_violations: 3          # Low
    design_inconsistencies: 5      # Moderate
    coupling_issues: 2             # Low
    
  performance_debt:
    known_bottlenecks: 4           # Moderate
    optimization_opportunities: 8   # High
    scalability_concerns: 2        # Low
    
  security_debt:
    known_vulnerabilities: 0       # Excellent
    hardening_gaps: 3             # Low
    audit_findings: 1             # Low
```

### KCP Pattern System-Wide Validation

**COMPREHENSIVE KCP ASSESSMENT**:
- LogicalCluster field presence in ALL CRDs across phases
- Workspace isolation consistency across all controllers
- RBAC pattern uniformity across all components
- Event scoping and handling consistency
- Multi-tenant data isolation verification

**Pattern Compliance Matrix**:

| Component Category | LogicalCluster | Workspace Isolation | RBAC Patterns | Event Scoping | Score |
|-------------------|----------------|-------------------|---------------|---------------|-------|
| API Types (Phase 1) | ✅ 100% | ✅ 100% | ✅ 100% | ✅ 100% | 100% |
| Controllers (Phase 1) | ✅ 100% | ✅ 95% | ✅ 98% | ✅ 92% | 96% |
| Webhooks (Phase 2) | ⚠️ 90% | ✅ 100% | ✅ 100% | ✅ 95% | 96% |
| Storage Layer (Phase 3) | ✅ 100% | ✅ 100% | ✅ 100% | ✅ 100% | 100% |

### System Performance Characteristics

**SYSTEM-WIDE PERFORMANCE AUDIT**:
- End-to-end performance across all phases
- Resource utilization patterns and optimization
- Scalability testing at enterprise levels
- Performance regression trend analysis
- Bottleneck identification and impact assessment

**Performance Health Metrics**:
```yaml
system_performance:
  api_performance:
    p50_response_time: "89ms"      # Excellent
    p95_response_time: "156ms"     # Good
    p99_response_time: "284ms"     # Acceptable
    error_rate: "0.01%"           # Excellent
    
  resource_utilization:
    cpu_efficiency: 78%            # Good
    memory_efficiency: 82%         # Good
    storage_efficiency: 91%        # Excellent
    network_efficiency: 85%        # Good
    
  scalability_metrics:
    max_workspaces_tested: 1500    # Exceeds requirements
    concurrent_operations: 2000    # Good
    data_volume_capacity: "10TB"   # Sufficient
    user_capacity: 5000            # Exceeds requirements
```

## Architecture Health Assessment

### System Health Categories

**ARCHITECTURE_HEALTHY** Requirements:
- Architectural consistency >90%
- Technical debt in Low-Medium range
- Performance meets all scalability targets
- Zero critical security vulnerabilities
- KCP pattern compliance 100%
- Operational readiness confirmed

**IMPROVEMENTS_REQUIRED** Triggers:
- Architectural consistency 80-90%
- Technical debt in Medium-High range
- Performance regressions but within limits
- Minor security vulnerabilities
- KCP pattern compliance 95-99%
- Some operational gaps

**ARCHITECTURE_DEGRADED** Triggers:
- Architectural consistency <80%
- Technical debt in High-Critical range
- Performance failing scalability targets
- Critical security vulnerabilities
- KCP pattern compliance <95%
- Major operational readiness gaps

### Audit Reporting Requirements

**COMPREHENSIVE AUDIT REPORT SECTIONS**:

1. **Executive Summary**
   - Overall architecture health assessment
   - Key findings and recommendations
   - Critical issues requiring immediate attention
   - System readiness for next phase/production

2. **Detailed Assessment Results**
   - Architectural consistency analysis
   - Technical debt inventory and prioritization
   - Performance characteristics and benchmarks
   - Security posture and vulnerability assessment
   - KCP pattern compliance verification
   - Operational readiness evaluation

3. **Remediation Plan** (if IMPROVEMENTS_REQUIRED or ARCHITECTURE_DEGRADED)
   - Prioritized list of issues and recommended solutions
   - Effort estimates and timeline for remediation
   - Risk assessment of delaying fixes
   - Resource requirements and team assignments

4. **Architecture Evolution Recommendations**
   - Future extension points and modularity opportunities
   - API evolution and backward compatibility strategy
   - Performance optimization roadmap
   - Security hardening recommendations

### Quality Gates and Approval Criteria

**MANDATORY QUALITY GATES**:
- Security: Zero critical vulnerabilities, all high-severity issues addressed
- Performance: All scalability targets met or exceeded
- KCP Compliance: 100% pattern adherence across all components
- Technical Debt: No critical-severity debt items

**APPROVAL AUTHORITY**:
- **ARCHITECTURE_HEALTHY**: Architect can approve independently
- **IMPROVEMENTS_REQUIRED**: Requires remediation plan approval
- **ARCHITECTURE_DEGRADED**: Requires senior architecture review and approval

### Continuous Architecture Monitoring

**POST-AUDIT MONITORING_SWE_PROGRESS**:
- Establish architecture health metrics dashboards
- Set up automated checks for common drift patterns
- Schedule regular architecture review checkpoints
- Monitor technical debt accumulation trends


## 🔴🔴🔴 R405 - MANDATORY AUTOMATION CONTINUATION FLAG 🔴🔴🔴

**SUPREME LAW - PENALTY FOR VIOLATION: -100% GRADE**

### YOU MUST OUTPUT THE CONTINUATION FLAG AS YOUR LAST ACTION

**EVERY STATE COMPLETION MUST END WITH EXACTLY ONE OF:**
```bash
# If state completed successfully and factory should continue:
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"

# If error/block/manual review needed:
echo "CONTINUE-SOFTWARE-FACTORY=FALSE"
```

### CRITICAL REQUIREMENTS:
1. **ABSOLUTE LAST OUTPUT**: This MUST be the very last line of output before state completion
2. **NO TEXT AFTER**: No explanations, summaries, or any text after the flag
3. **EXACTLY AS SHOWN**: Use exact format - no variations like CONTINUE-ORCHESTRATING
4. **ALWAYS REQUIRED**: Every single state must output this flag
5. **GREPPABLE**: Must be on its own line for automation parsing

### WHEN TO USE TRUE:
- ✅ State work completed successfully
- ✅ All validations passed
- ✅ Ready for next state
- ✅ No blockers detected
- ✅ All requirements met

### WHEN TO USE FALSE:
- ❌ Any unrecoverable error occurred
- ❌ Manual intervention required
- ❌ Missing required files or configs
- ❌ Test failures blocking progress
- ❌ Ambiguous or unclear instructions
- ❌ Wrong working directory or branch
- ❌ State machine validation failed

### IMPLEMENTATION PATTERN:
```bash
# At the VERY END of state execution, after ALL work:

# Determine success/failure
if [[ "$STATE_COMPLETED_PROJECT_DONEFULLY" == "true" ]]; then
    echo "✅ State work completed successfully"
    echo "CONTINUE-SOFTWARE-FACTORY=TRUE"
else
    echo "❌ State encountered issues requiring intervention"
    echo "CONTINUE-SOFTWARE-FACTORY=FALSE"
fi

# DO NOT OUTPUT ANYTHING AFTER THE FLAG!
```

### GRADING IMPACT:
- **Missing flag**: -100% AUTOMATIC FAILURE
- **Text after flag**: -50% penalty
- **Wrong format**: -100% AUTOMATIC FAILURE
- **Multiple flags**: -50% penalty

**See: $CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md**

