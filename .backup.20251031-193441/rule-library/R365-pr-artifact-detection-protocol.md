# 🚨🚨🚨 BLOCKING RULE R365 - PR Artifact Detection and Inventory Protocol

**Criticality:** BLOCKING - Undetected artifacts = PR contamination
**Enforcement:** MANDATORY - All PR-Ready operations
**Created:** 2025-01-21

## PURPOSE
Ensure all Software Factory artifacts are detected and inventoried before PR creation to prevent contamination of production code with development artifacts.

## SOFTWARE FACTORY ARTIFACTS TO DETECT

### Critical Artifacts (MUST be removed)
```
.software-factory/          # All SF metadata directories
orchestrator-state-v3.json     # State tracking files
pr-ready-state.json        # PR state tracking
*.todo                     # TODO persistence files
effort-*/                  # Effort working directories
phase-plans/               # Phase planning documents
wave-plans/                # Wave planning documents
*-PLAN.md                  # Planning documents
*-REPORT.md                # Report documents
SPLIT-*/                   # Split tracking files
```

### Documentation Artifacts (MUST be evaluated)
```
INTEGRATE_WAVE_EFFORTS-*.md           # Integration reports
CODE-REVIEW-*.md          # Code review reports
ARCHITECTURE-*.md         # Architecture documents
IMPLEMENTATION-*.md       # Implementation plans
FIX-*.md                  # Fix documentation
```

### Temporary Artifacts (MUST be cleaned)
```
.state_rules_read_*       # Rule acknowledgment markers
.git_config_backup_*      # Config backups
tmp-*                     # Temporary files
*.orig                    # Merge conflict originals
*.rej                     # Patch rejects
```

## DETECTION PROTOCOL

### Step 1: Comprehensive Scan
```bash
# Create artifact patterns file
cat > sf-artifact-patterns.txt << 'EOF'
.software-factory/
orchestrator-state-v3.json
pr-ready-state.json
*.todo
effort-*/
phase-plans/
wave-plans/
*-PLAN.md
*-REPORT.md
SPLIT-*/
INTEGRATE_WAVE_EFFORTS-*.md
CODE-REVIEW-*.md
ARCHITECTURE-*.md
IMPLEMENTATION-*.md
FIX-*.md
.state_rules_read_*
.git_config_backup_*
tmp-*
*.orig
*.rej
EOF

# Scan for all artifacts
ARTIFACTS_FOUND=0
for pattern in $(cat sf-artifact-patterns.txt); do
    FILES=$(find . -name "$pattern" -o -type d -name "$pattern" 2>/dev/null)
    if [ -n "$FILES" ]; then
        echo "Found artifacts matching: $pattern"
        echo "$FILES" | while read -r file; do
            echo "  - $file"
            ARTIFACTS_FOUND=$((ARTIFACTS_FOUND + 1))
        done
    fi
done
```

### Step 2: Create Inventory
```bash
# Generate JSON inventory
cat > PR-ARTIFACT-INVENTORY.json << 'EOF'
{
  "scan_timestamp": "$(date -u +%Y-%m-%dT%H:%M:%SZ)",
  "branch": "$(git branch --show-current)",
  "artifacts": {
    "critical": [],
    "documentation": [],
    "temporary": []
  },
  "total_count": 0,
  "action_required": true
}
EOF
```

### Step 3: Categorize Artifacts
```bash
# Categorize each found artifact
for artifact in $ARTIFACTS_LIST; do
    if [[ "$artifact" =~ \.(software-factory|todo)$ ]]; then
        CATEGORY="critical"
    elif [[ "$artifact" =~ \.(md)$ ]]; then
        CATEGORY="documentation"
    else
        CATEGORY="temporary"
    fi

    # Add to inventory
    jq ".artifacts.$CATEGORY += [\"$artifact\"]" PR-ARTIFACT-INVENTORY.json > tmp.json
    mv tmp.json PR-ARTIFACT-INVENTORY.json
done
```

### Step 4: Generate Report
```bash
cat > PR-ARTIFACT-SCAN-REPORT.md << 'EOF'
# PR Artifact Scan Report
Generated: $(date)
Branch: $(git branch --show-current)

## Summary
- Total artifacts found: $ARTIFACTS_FOUND
- Critical artifacts: $(jq '.artifacts.critical | length' PR-ARTIFACT-INVENTORY.json)
- Documentation artifacts: $(jq '.artifacts.documentation | length' PR-ARTIFACT-INVENTORY.json)
- Temporary artifacts: $(jq '.artifacts.temporary | length' PR-ARTIFACT-INVENTORY.json)

## Action Required
$([ $ARTIFACTS_FOUND -gt 0 ] && echo "⚠️ Artifacts must be removed before PR creation" || echo "✅ No artifacts detected")

## Detailed Inventory
[See PR-ARTIFACT-INVENTORY.json for complete list]
EOF
```

## ENFORCEMENT REQUIREMENTS

### MANDATORY Actions
1. **SCAN EVERY BRANCH** before PR preparation
2. **INVENTORY ALL ARTIFACTS** in structured format
3. **CATEGORIZE BY CRITICALITY** for removal priority
4. **REPORT FINDINGS** to orchestrator
5. **VERIFY REMOVAL** after cleanup

### PROHIBITED Actions
❌ **NEVER** skip artifact detection
❌ **NEVER** ignore found artifacts
❌ **NEVER** delete without inventory
❌ **NEVER** proceed with artifacts present

## PROJECT_DONE CRITERIA
✅ All branches scanned
✅ Inventory JSON created
✅ Report markdown generated
✅ All artifacts documented
✅ Removal plan created

## GRADING PENALTIES
- Missing artifact scan: -30%
- Incomplete inventory: -20%
- PR with SF artifacts: -50%
- No verification after cleanup: -25%

## INTEGRATE_WAVE_EFFORTS WITH OTHER RULES
- Works with **R367** (Cleanup Protocol)
- Supports **R369** (Validation Protocol)
- Enables **R370** (PR Plan Creation)

---

*This rule ensures Software Factory artifacts never contaminate production code.*