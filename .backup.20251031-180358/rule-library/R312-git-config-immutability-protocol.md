# 🔴🔴🔴 R312: Git Config Immutability Protocol (SUPREME LAW) 🔴🔴🔴

**Rule ID:** R312.0.0  
**Criticality:** 🔴🔴🔴 SUPREME LAW  
**Category:** Effort Isolation, Repository Security  
**Agents:** orchestrator, sw-engineer, code-reviewer  
**States:** CREATE_NEXT_INFRASTRUCTURE, CREATE_NEXT_INFRASTRUCTURE, IMPLEMENTATION  

## 🚨🚨🚨 ABSOLUTE REQUIREMENT: .git/config MUST BE READONLY 🚨🚨🚨

### RULE SUMMARY

After the orchestrator creates effort/split infrastructure, the `.git/config` file MUST be:
1. **Changed to root ownership** (chown root:root) - prevents permission changes
2. **Made READONLY** (chmod 444) - prevents modifications

This DOUBLE PROTECTION prevents SW engineers from:
- Changing branches
- Modifying remotes
- Pulling from other branches
- Contaminating their isolated work environment
- Bypassing protection with chmod (they're not the owner)

**VIOLATION = -100% AUTOMATIC FAILURE**

## 🔴 WHY THIS MATTERS - THE ISOLATION PRINCIPLE 🔴

### The Software Factory Isolation Model

1. **Each effort is an ISOLATED unit of work**
   - ONLY contains its own changes
   - NO contamination from other efforts
   - NO pulling from main or integration branches
   - NO switching to other branches

2. **SW Engineers are IMPLEMENTERS, not INTEGRATORS**
   - They implement ONLY their assigned scope
   - They work ONLY on their assigned branch
   - They push ONLY to their assigned remote
   - Integration happens ONLY during integration states

3. **Without This Protection, Everything Breaks:**
   ```bash
   # ❌ SW engineer pulls from main → includes changes not in scope
   # ❌ SW engineer switches branches → corrupts both efforts
   # ❌ SW engineer changes remote → pushes to wrong repository
   # ❌ SW engineer pulls from integration → gets everyone else's work
   # Result: Size measurements wrong, integration impossible, efforts contaminated
   ```

## 🔴🔴🔴 ORCHESTRATOR IMPLEMENTATION REQUIREMENTS 🔴🔴🔴

### In CREATE_NEXT_INFRASTRUCTURE State

After creating each effort's infrastructure, the orchestrator MUST:

```bash
# MANDATORY: After branch creation and remote setup
lock_git_config() {
    local EFFORT_DIR="$1"
    
    echo "🔒 R312: Locking git config for effort isolation..."
    
    # Ensure we're in the effort directory
    cd "$EFFORT_DIR"
    
    # Verify .git/config exists
    if [ ! -f .git/config ]; then
        echo "❌ FATAL: No .git/config found in $EFFORT_DIR"
        exit 312
    fi
    
    # Store current permissions and ownership for verification
    BEFORE_PERMS=$(stat -c %a .git/config 2>/dev/null || stat -f %A .git/config)
    BEFORE_OWNER=$(stat -c %U:%G .git/config 2>/dev/null || stat -f %Su:%Sg .git/config)
    
    # DOUBLE PROTECTION: Change ownership AND permissions
    if command -v sudo >/dev/null 2>&1; then
        # Full protection with root ownership
        echo "🔐 Applying full protection (root ownership + readonly)..."
        sudo chown root:root .git/config
        sudo chmod 444 .git/config
        PROTECTION_LEVEL="FULL"
    else
        # Fallback to permission-only protection
        echo "⚠️ sudo not available - applying permission-only protection"
        chmod 444 .git/config
        PROTECTION_LEVEL="PARTIAL"
    fi
    
    # Verify protection
    if [ -w .git/config ]; then
        echo "❌ R312 VIOLATION: Failed to make .git/config readonly!"
        echo "Config is still writable - effort isolation compromised!"
        exit 312
    fi
    
    # Verify ownership if sudo was available
    if [ "$PROTECTION_LEVEL" = "FULL" ]; then
        CURRENT_OWNER=$(stat -c %U:%G .git/config 2>/dev/null || stat -f %Su:%Sg .git/config)
        if [ "$CURRENT_OWNER" != "root:root" ]; then
            echo "⚠️ WARNING: Ownership not changed to root:root (got $CURRENT_OWNER)"
            echo "Protection may be weaker than intended"
        fi
    fi
    
    # Log the protection
    echo "✅ R312: Git config locked"
    echo "   Protection level: $PROTECTION_LEVEL"
    echo "   Ownership: $BEFORE_OWNER → $(stat -c %U:%G .git/config 2>/dev/null || echo 'unknown')"
    echo "   Permissions: $BEFORE_PERMS → 444"
    echo "📝 Protected operations now prevented:"
    echo "   ❌ git checkout [other-branch]"
    echo "   ❌ git pull origin main"
    echo "   ❌ git remote add/remove"
    echo "   ❌ git branch --set-upstream-to"
    
    # Create protection marker
    touch .git/R312_CONFIG_LOCKED
    echo "$(date): Config locked by orchestrator per R312" > .git/R312_CONFIG_LOCKED
}

# Call after each effort setup
prepare_effort_for_agent() {
    # ... existing infrastructure setup ...
    
    # CRITICAL NEW STEP: Lock the config
    lock_git_config "$EFFORT_DIR"
    
    echo "✅ R312 enforced: Git config is IMMUTABLE for SW engineer"
}
```

### In CREATE_NEXT_INFRASTRUCTURE State

Similarly for split infrastructure:

```bash
create_single_split_infrastructure() {
    # ... existing split setup ...
    
    # CRITICAL: Lock config before SW engineer spawns
    lock_git_config "$SPLIT_DIR"
    
    echo "✅ R312 enforced: Split git config is IMMUTABLE"
}
```

## 🚨🚨🚨 SW ENGINEER VALIDATION REQUIREMENTS 🚨🚨🚨

### Mandatory Pre-Flight Check

SW Engineers MUST verify config is readonly as part of R001 pre-flight checks:

```bash
# R312 VALIDATION: Config must be READONLY
validate_r312_config_lock() {
    echo "🔍 R312: Validating git config immutability..."
    
    if [ ! -f .git/config ]; then
        echo "❌ FATAL: No .git/config found!"
        exit 1
    fi
    
    # Check if config is writable
    if [ -w .git/config ]; then
        echo "❌❌❌ R312 SECURITY VIOLATION DETECTED!"
        echo "Git config is WRITABLE - this should be READONLY!"
        echo "This violates effort isolation protocol!"
        echo ""
        echo "Expected: .git/config with 444 permissions (readonly)"
        echo "Found: .git/config is writable"
        echo ""
        echo "STOPPING: Cannot proceed with writable config"
        exit 312
    fi
    
    # Check ownership for full protection validation
    CURRENT_OWNER=$(stat -c %U:%G .git/config 2>/dev/null || stat -f %Su:%Sg .git/config)
    if [ "$CURRENT_OWNER" = "root:root" ]; then
        echo "   🔐 FULL protection: root-owned + readonly"
    else
        echo "   🔒 PARTIAL protection: readonly only (owner: $CURRENT_OWNER)"
    fi
    
    # Verify lock marker exists
    if [ ! -f .git/R312_CONFIG_LOCKED ]; then
        echo "⚠️ WARNING: R312 lock marker missing"
        echo "Config appears readonly but marker absent"
    fi
    
    echo "✅ R312 VALIDATED: Git config is properly locked (readonly)"
    echo "📋 Protected from:"
    echo "   ✅ Cannot switch branches"
    echo "   ✅ Cannot change remotes"
    echo "   ✅ Cannot pull from other sources"
    echo "   ✅ Effort isolation guaranteed"
}

# Add to SW Engineer startup
perform_preflight_checks() {
    # ... existing checks ...
    
    # R312: Validate config immutability
    validate_r312_config_lock
}
```

## 🔴 PREVENTED OPERATIONS (MOSTLY BLOCKED) 🔴

With .git/config protected, most dangerous operations are prevented:

```bash
# ❌ THESE ARE BLOCKED OR SIGNIFICANTLY HARDER:

git checkout main                           # Blocked (config changes required)
git checkout -b new-branch                  # Blocked (config changes required)
git pull origin main                        # Blocked (cannot update branch config)
git pull origin phase1-integration          # Blocked (cannot update branch config)
git fetch origin other-branch               # Can fetch but cannot checkout/merge
git remote add upstream https://...         # Blocked (cannot modify config)
git remote remove origin                    # Blocked (cannot modify config)
git remote set-url origin https://...       # Blocked (cannot modify config)
git branch --set-upstream-to=origin/main    # Blocked (cannot change tracking)

# ⚠️ NOTE: Some git operations create new config files
# Git may work around readonly by creating .git/config.lock and renaming
# This is why BOTH technical protection AND SW engineer discipline are required
# The protection makes violations DELIBERATE rather than accidental
chmod 644 .git/config                       # Blocked if root-owned (not file owner)
sudo chmod 644 .git/config                  # Requires deliberate root override
```

## ✅ ALLOWED OPERATIONS (STILL WORK) ✅

These operations remain functional:

```bash
# ✅ THESE STILL WORK AS EXPECTED:

git add file.go                    # Can stage files
git commit -m "message"            # Can commit changes
git push                           # Can push to tracked remote
git push origin HEAD               # Can push current branch
git status                         # Can check status
git diff                           # Can view differences
git log                            # Can view history
git stash                          # Can stash changes
git fetch                          # Can fetch (but not checkout)
```

## 🔴🔴🔴 SPECIAL HANDLING FOR INTEGRATE_WAVE_EFFORTS AGENTS 🔴🔴🔴

### Integration Agents NEED Writable Configs

Integration agents (spawned in INTEGRATE_WAVE_EFFORTS states) require the ability to:
- Pull from multiple effort branches
- Create integration branches
- Merge efforts together

```bash
# For integration agents ONLY:
unlock_git_config_for_integration() {
    local INTEGRATE_WAVE_EFFORTS_DIR="$1"
    
    echo "🔓 R312 EXCEPTION: Unlocking config for INTEGRATE_WAVE_EFFORTS agent"
    
    cd "$INTEGRATE_WAVE_EFFORTS_DIR"
    
    # Check current ownership
    CURRENT_OWNER=$(stat -c %U:%G .git/config 2>/dev/null || stat -f %Su:%Sg .git/config)
    
    # Restore ownership and permissions for integration work
    if [ "$CURRENT_OWNER" = "root:root" ]; then
        # Need sudo to change from root ownership
        if command -v sudo >/dev/null 2>&1; then
            echo "🔓 Restoring user ownership from root..."
            sudo chown $(id -u):$(id -g) .git/config
            sudo chmod 644 .git/config
        else
            echo "❌ ERROR: Config is root-owned but sudo not available!"
            echo "Cannot unlock config for integration"
            exit 312
        fi
    else
        # Simple permission change
        chmod 644 .git/config
    fi
    
    # Verify it's now writable
    if [ ! -w .git/config ]; then
        echo "❌ Failed to unlock config for integration!"
        exit 312
    fi
    
    # Create exception marker
    echo "$(date): Config unlocked for integration per R312 exception" > .git/R312_INTEGRATE_WAVE_EFFORTS_EXCEPTION
    
    echo "✅ Config unlocked for integration operations"
    echo "   Owner: $(stat -c %U:%G .git/config 2>/dev/null || echo 'unknown')"
    echo "   Permissions: $(stat -c %a .git/config 2>/dev/null || echo 'unknown')"
}
```

### Integration State Detection

```bash
# Check if current state is integration-related
is_integration_state() {
    local STATE="$1"
    case "$STATE" in
        INTEGRATE_WAVE_EFFORTS|WAVE_INTEGRATE_WAVE_EFFORTS|INTEGRATE_PHASE_WAVES|SPAWN_INTEGRATION_AGENT)
            return 0  # True - is integration state
            ;;
        *)
            return 1  # False - not integration state
            ;;
    esac
}
```

## 🔐 PROTECTION LEVELS

### Full Protection (Preferred)
- **Ownership**: root:root
- **Permissions**: 444 (r--r--r--)
- **Result**: Even with chmod attempts, SW engineer cannot modify (not owner)

### Partial Protection (Fallback)
- **Ownership**: Original user
- **Permissions**: 444 (r--r--r--)
- **Result**: Prevents accidental changes but can be bypassed with chmod

### Docker/Container Considerations
- Containers running as root: Protection still applies
- Containers without sudo: Falls back to permission-only
- Rootless containers: Permission-only protection

## 📊 GRADING IMPACT

### Violations and Penalties

1. **Orchestrator fails to lock config**: -100% (Effort isolation broken)
2. **SW Engineer modifies config permissions**: -100% (Security violation)
3. **SW Engineer bypasses protection**: -100% (Deliberate violation)
4. **Missing R312 validation in pre-flight**: -50% (Incomplete checks)
5. **Integration agent with locked config**: -30% (Prevents integration)
6. **No protection applied at all**: -100% (Complete failure)

### Success Metrics

- ✅ All effort configs locked after setup
- ✅ All SW engineers validate lock on startup
- ✅ Zero branch contamination incidents
- ✅ Zero cross-effort pollution
- ✅ 100% effort isolation maintained

## 🚨 ENFORCEMENT AND MONITORING_SWE_PROGRESS 🚨

### Orchestrator Monitoring

```bash
# Periodic validation during MONITOR state
monitor_r312_compliance() {
    echo "🔍 R312 Compliance Check..."
    
    for effort_dir in efforts/phase*/wave*/E*; do
        if [ -d "$effort_dir/.git" ]; then
            if [ -w "$effort_dir/.git/config" ]; then
                echo "❌ R312 VIOLATION: $effort_dir has writable config!"
                return 312
            fi
        fi
    done
    
    echo "✅ R312: All effort configs properly locked"
}
```

### Audit Trail

Each locked config creates an audit entry:
- `.git/R312_CONFIG_LOCKED` - Lock timestamp and actor
- `.git/R312_INTEGRATE_WAVE_EFFORTS_EXCEPTION` - Exception for integration

## 🔴 CRITICAL REMINDERS 🔴

1. **This is a SECURITY measure, not a suggestion**
2. **Effort isolation is FUNDAMENTAL to Software Factory**
3. **Without this, line counting becomes WRONG**
4. **Without this, integration becomes IMPOSSIBLE**
5. **Without this, the entire model BREAKS**

## Example Implementation Flow

```bash
# 1. Orchestrator creates effort infrastructure
prepare_effort_for_agent 2 1 "E2.1.1"

# 2. Config is automatically locked
# .git/config now has 444 permissions

# 3. SW Engineer is spawned
/spawn software-engineer --effort E2.1.1

# 4. SW Engineer validates lock in pre-flight
validate_r312_config_lock  # Must pass or exit

# 5. SW Engineer works in isolation
git add src/feature.go     # ✅ Works
git commit -m "feat: add"  # ✅ Works
git push                   # ✅ Works
git checkout main          # ❌ BLOCKED! Permission denied
chmod 644 .git/config      # ❌ BLOCKED! Not owner (if root-owned)
sudo chmod 644 .git/config # ❌ Requires deliberate root override

# 6. Integration agent later gets exception
unlock_git_config_for_integration "/efforts/integration"
# Now can pull from multiple branches
```

---

**REMEMBER**: Effort isolation is not optional. It's the foundation of the entire Software Factory model. This rule ensures that foundation remains solid.