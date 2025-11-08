# 🔴🔴🔴 RULE R381 - Library Version Consistency Protocol

**Criticality:** SUPREME LAW - Project Stability Foundation
**Grading Impact:** -100% for unauthorized version updates
**Enforcement:** CONTINUOUS - Every dependency change monitored

## Rule Statement

ALL library versions chosen in previous efforts/waves/phases are IMMUTABLE. NO agent may update pinned versions without explicit user approval and mandatory fix cascade.

## Core Principle

**LIBRARY VERSIONS ARE SACRED**
- Previous work was tested with SPECIFIC versions
- Changing versions can break existing functionality
- Version consistency across entire project is CRITICAL
- Convenience updates are FORBIDDEN

## Version Immutability Requirements

### 1. Metadata Files Are READ-ONLY for Versions
```yaml
PROTECTED_FILES:
  - go.mod            # Go module versions
  - go.sum            # Go checksums
  - package.json      # Node.js dependencies
  - package-lock.json # Locked versions
  - requirements.txt  # Python dependencies
  - Pipfile.lock     # Python locked versions
  - Cargo.toml       # Rust dependencies
  - Cargo.lock       # Rust locked versions
  - pom.xml          # Java/Maven
  - build.gradle     # Java/Gradle
  - composer.json    # PHP dependencies
  - composer.lock    # PHP locked versions
```

### 2. ALLOWED Operations
```bash
# ✅ ALLOWED: Adding NEW dependencies
echo 'github.com/new/library v1.0.0' >> go.mod

# ✅ ALLOWED: Using existing versions
import "github.com/spf13/cobra"  # Use v1.7.0 as locked

# ✅ ALLOWED: Reading version information
grep "cobra" go.mod  # Check what version is in use
```

### 3. FORBIDDEN Operations
```bash
# ❌ FORBIDDEN: Updating existing versions
sed -i 's/v1.7.0/v1.8.0/' go.mod  # NEVER DO THIS

# ❌ FORBIDDEN: Using "latest" without checking
go get github.com/spf13/cobra@latest  # NO! Check current first

# ❌ FORBIDDEN: Changing version ranges
npm install package@^2.0.0  # NO! Keep exact version
```

## Version Tracking Requirements

### 1. State File Tracking
```json
{
  "version_locks": {
    "established_by": "phase1-wave1-effort1",
    "locked_at": "2025-01-20T10:00:00Z",
    "libraries": {
      "github.com/spf13/cobra": {
        "version": "v1.7.0",
        "locked_by": "phase1-wave1-effort1",
        "reason": "CLI framework",
        "modification_requires": "fix_cascade"
      },
      "github.com/google/go-containerregistry": {
        "version": "v0.16.1",
        "locked_by": "phase1-wave1-effort2",
        "reason": "Container registry client",
        "modification_requires": "fix_cascade"
      }
    }
  }
}
```

### 2. Planning Documentation
```markdown
## Dependency Requirements
- MUST use cobra v1.7.0 (locked by phase1)
- MUST use go-containerregistry v0.16.1 (locked by phase1)
- New dependencies allowed: prometheus/client_golang (suggest v1.17.0)
```

## Exception Handling Protocol

### ONLY Valid Reasons for Updates:
1. **Security Vulnerability** (CVE with HIGH/CRITICAL severity)
2. **Critical Bug** (blocks core functionality)
3. **Upstream Deprecation** (with removal deadline)
4. **Explicit User Approval** (documented in plan)

### NOT Valid Reasons:
- New features available ❌
- Performance improvements ❌
- Developer convenience ❌
- "Latest is better" ❌
- Test framework preferences ❌

## Enforcement Mechanisms

### 1. Code Reviewer Checks
```bash
check_version_consistency() {
    local pr_branch="$1"
    local base_branch="$2"

    # Get version changes
    git diff "$base_branch...$pr_branch" -- go.mod package.json requirements.txt

    # ANY version change triggers alert
    if version_changed; then
        echo "🔴🔴🔴 CRITICAL: Version change detected!"
        echo "Requires: 1) User approval 2) Fix cascade (R382)"
        exit 381
    fi
}
```

### 2. SW Engineer Guards
```bash
before_adding_dependency() {
    local package="$1"

    # Check if already exists
    if grep -q "$package" go.mod; then
        echo "⚠️ Package exists - MUST use locked version"
        show_locked_version "$package"
        exit 381
    fi
}
```

### 3. Integration Validation
```bash
validate_version_consistency() {
    local effort_branches=("$@")

    for branch in "${effort_branches[@]}"; do
        # All must have identical versions
        extract_versions "$branch" > "/tmp/$branch.versions"
    done

    # Compare all version files
    if ! all_identical /tmp/*.versions; then
        echo "❌ Version drift detected between efforts!"
        exit 381
    fi
}
```

## Grading Impact

### Automatic Failures (-100%)
- Updating any locked version without approval
- Changing version in metadata files
- Using different versions across efforts
- Ignoring version lock warnings

### Major Penalties (-50%)
- Not checking existing versions before adding
- Suggesting version updates in reviews
- Planning with incompatible versions

### Minor Penalties (-20%)
- Not documenting version choices
- Missing version tracking in state
- Incomplete version verification

## Examples

### ❌ BAD: Casual Version Update
```go
// Developer thinks: "Let's use the latest cobra"
// go get github.com/spf13/cobra@latest
// This BREAKS all previous work tested with v1.7.0!
```

### ❌ BAD: Version Drift
```yaml
# Effort 1: Uses cobra v1.7.0
# Effort 2: Updates to cobra v1.8.0
# Integration: CONFLICT! Different versions!
```

### ✅ GOOD: Respecting Locks
```go
// Check current version first
// $ grep cobra go.mod
// github.com/spf13/cobra v1.7.0

// Use the locked version
import "github.com/spf13/cobra"
// Works perfectly with all other efforts
```

### ✅ GOOD: Adding New Dependency
```go
// Check it doesn't exist
// $ grep prometheus go.mod
// (no output - safe to add)

// Add with specific version
// go get github.com/prometheus/client_golang@v1.17.0
```

## Integration with Other Rules

- **R382**: Triggers fix cascade if update required
- **R300**: Fix cascade includes version validation
- **R362**: Changing versions = architectural violation
- **R363**: Version consistency enables sequential merging

## Recovery Protocol

If version update was made:
1. **IMMEDIATE STOP** - No further work
2. **Assess Impact** - Which efforts affected?
3. **User Decision** - Revert or cascade?
4. **If Cascade** - Trigger R382 protocol
5. **If Revert** - Reset to locked version

## Self-Monitoring Commands

```bash
# For SW Engineers
check_my_versions() {
    echo "📋 Checking version compliance..."
    git diff HEAD -- go.mod package.json
    if [ $? -ne 0 ]; then
        echo "⚠️ Version changes detected - VERIFY AUTHORIZED"
    fi
}

# For Code Reviewers
validate_pr_versions() {
    echo "🔍 Validating version consistency..."
    git diff main...HEAD -- '**/go.mod' '**/package.json'
    # ANY output = potential violation
}
```

---

**REMEMBER**: Version stability is project stability. Every version change cascades through the entire codebase. RESPECT THE LOCKS!