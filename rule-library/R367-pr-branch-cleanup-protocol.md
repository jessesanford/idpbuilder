# 🚨🚨🚨 BLOCKING RULE R367 - PR Branch Cleanup Protocol

**Criticality:** BLOCKING - Clean branches for production PRs
**Enforcement:** MANDATORY - All PR-Ready operations
**Created:** 2025-01-21

## PURPOSE
Remove all Software Factory artifacts, temporary files, and development remnants from branches before PR creation to ensure only production-ready code is submitted.

## CLEANUP CATEGORIES

### Category 1: Software Factory Artifacts (MUST REMOVE)
```bash
# These MUST be removed from ALL branches
.software-factory/              # All SF metadata
orchestrator-state.json         # State tracking
pr-ready-state.json            # PR state tracking
*.todo                         # TODO files
effort-*/                      # Effort directories
phase-plans/                   # Planning documents
wave-plans/                    # Wave documents
SPLIT-*/                       # Split tracking
*-PLAN.md                      # Plan documents
*-REPORT.md                    # Report documents
.state_rules_read_*            # Rule markers
```

### Category 2: Development Artifacts (MUST REMOVE)
```bash
# Development and debugging artifacts
*.log                          # Log files (unless needed)
*.tmp                          # Temporary files
*.bak                          # Backup files
*.orig                         # Merge conflict originals
*.rej                          # Patch rejects
.DS_Store                      # macOS artifacts
Thumbs.db                      # Windows artifacts
*.swp                          # Vim swap files
*~                            # Editor backups
```

### Category 3: Documentation (EVALUATE)
```bash
# Evaluate case-by-case
INTEGRATION-*.md               # May be useful for PR description
CODE-REVIEW-*.md              # Document review history
ARCHITECTURE-*.md             # Keep if adds value
IMPLEMENTATION-*.md           # Usually remove
FIX-*.md                      # Usually remove
```

## CLEANUP PROTOCOL

### Step 1: Create Cleanup Branch
```bash
# Work on cleanup branch to preserve original
ORIGINAL_BRANCH=$(git branch --show-current)
CLEANUP_BRANCH="${ORIGINAL_BRANCH}-cleanup"

git checkout -b $CLEANUP_BRANCH
echo "✅ Created cleanup branch: $CLEANUP_BRANCH"
```

### Step 2: Scan for Artifacts
```bash
# Run artifact detection (R365)
echo "🔍 Scanning for SF artifacts..."

# Find all SF artifacts
SF_ARTIFACTS=$(find . -type f \( \
    -name "orchestrator-state.json" -o \
    -name "pr-ready-state.json" -o \
    -name "*.todo" -o \
    -name "*-PLAN.md" -o \
    -name "*-REPORT.md" -o \
    -name ".state_rules_read_*" \
\) -o -type d \( \
    -name ".software-factory" -o \
    -name "effort-*" -o \
    -name "phase-plans" -o \
    -name "wave-plans" -o \
    -name "SPLIT-*" \
\))

echo "Found $(echo "$SF_ARTIFACTS" | wc -l) SF artifacts"
```

### Step 3: Remove SF Artifacts
```bash
# Remove each artifact with confirmation
for artifact in $SF_ARTIFACTS; do
    echo "Removing: $artifact"

    if [ -d "$artifact" ]; then
        # Remove directory
        git rm -rf "$artifact" 2>/dev/null || rm -rf "$artifact"
    else
        # Remove file
        git rm -f "$artifact" 2>/dev/null || rm -f "$artifact"
    fi
done

# Verify removal
if [ -n "$(find . -name 'orchestrator-state.json' 2>/dev/null)" ]; then
    echo "❌ ERROR: SF artifacts still present!"
    exit 1
fi
```

### Step 4: Remove Development Artifacts
```bash
# Clean development artifacts
echo "🧹 Cleaning development artifacts..."

# Remove common dev files
find . -type f \( \
    -name "*.log" -o \
    -name "*.tmp" -o \
    -name "*.bak" -o \
    -name "*.orig" -o \
    -name "*.rej" -o \
    -name ".DS_Store" -o \
    -name "Thumbs.db" -o \
    -name "*.swp" -o \
    -name "*~" \
\) -exec rm -f {} \;

echo "✅ Development artifacts removed"
```

### Step 5: Evaluate Documentation
```bash
# Review documentation files
echo "📄 Evaluating documentation..."

DOCS=$(find . -name "*.md" | grep -E "(INTEGRATION|CODE-REVIEW|ARCHITECTURE|IMPLEMENTATION|FIX)")

for doc in $DOCS; do
    echo "Found: $doc"

    # Check if referenced in code
    REFS=$(grep -r "$(basename $doc)" --include="*.md" --exclude="$doc" . | wc -l)

    if [ $REFS -eq 0 ]; then
        echo "  → No references, removing"
        git rm -f "$doc" 2>/dev/null || rm -f "$doc"
    else
        echo "  → Referenced $REFS times, keeping"
    fi
done
```

### Step 6: Update .gitignore
```bash
# Ensure SF artifacts are ignored going forward
cat >> .gitignore << 'EOF'

# Software Factory artifacts (auto-added by cleanup)
.software-factory/
orchestrator-state.json
pr-ready-state.json
*.todo
effort-*/
phase-plans/
wave-plans/
SPLIT-*/
*-PLAN.md
*-REPORT.md
.state_rules_read_*
EOF

git add .gitignore
echo "✅ Updated .gitignore"
```

### Step 7: Commit Cleanup
```bash
# Commit all cleanup changes
git add -A
git commit -m "chore: remove Software Factory artifacts for PR preparation

- Removed all SF metadata and state files
- Cleaned development artifacts
- Updated .gitignore to prevent re-addition
- Preserved valuable documentation

This cleanup ensures the branch contains only production code."
```

### Step 8: Verify Cleanliness
```bash
# Final verification
echo "🔍 Final verification..."

# Check for any remaining SF artifacts
REMAINING=$(find . -type f -name "*orchestrator*" -o -name "*.todo" -o -name "*-PLAN.md")

if [ -n "$REMAINING" ]; then
    echo "❌ ERROR: Artifacts still present:"
    echo "$REMAINING"
    exit 1
fi

echo "✅ Branch is clean and PR-ready"
```

### Step 9: Update Original Branch
```bash
# Apply cleanup to original branch
git checkout $ORIGINAL_BRANCH
git reset --hard $CLEANUP_BRANCH

# Push cleaned branch
git push origin $ORIGINAL_BRANCH --force-with-lease

# Clean up cleanup branch
git branch -D $CLEANUP_BRANCH
```

## VERIFICATION CHECKLIST

### Pre-Cleanup
- [ ] Backup branch created
- [ ] Artifact inventory completed
- [ ] Important docs identified

### During Cleanup
- [ ] All SF artifacts removed
- [ ] Dev artifacts cleaned
- [ ] Documentation evaluated
- [ ] .gitignore updated

### Post-Cleanup
- [ ] No SF artifacts remain
- [ ] Branch builds successfully
- [ ] Tests still pass
- [ ] Pushed to origin

## OUTPUT ARTIFACTS

### Required Reports
```bash
cat > PR-CLEANUP-REPORT.md << 'EOF'
# PR Branch Cleanup Report
Branch: ${ORIGINAL_BRANCH}
Date: $(date)

## Artifacts Removed
- SF artifacts: [count]
- Dev artifacts: [count]
- Documentation: [count]

## Documentation Preserved
- [List kept documents]

## Verification
- [ ] No SF artifacts remain
- [ ] Build successful
- [ ] Tests passing
- [ ] .gitignore updated

## Branch Status
✅ Clean and PR-ready
EOF
```

## ERROR RECOVERY

### If Important File Deleted
```bash
# Recover from git history
git checkout HEAD~1 -- path/to/important/file
git add path/to/important/file
git commit -m "restore: recover important file"
```

### If Cleanup Breaks Build
```bash
# Revert cleanup commit
git revert HEAD
# Selectively re-apply cleanup
```

## SUCCESS CRITERIA
✅ All SF artifacts removed
✅ No development artifacts
✅ Documentation evaluated
✅ .gitignore updated
✅ Branch builds and tests pass

## GRADING PENALTIES
- SF artifacts in PR: -50%
- No cleanup verification: -25%
- Breaking build during cleanup: -30%
- No .gitignore update: -15%

## INTEGRATION WITH OTHER RULES
- Depends on **R365** (Artifact Detection)
- Enables **R369** (Validation Protocol)
- Required before **R370** (PR Plan Creation)

---

*This rule ensures branches are production-clean before PR submission.*