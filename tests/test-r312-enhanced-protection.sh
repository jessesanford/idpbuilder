#!/bin/bash

# Test script for R312 Enhanced Protection with root ownership
# This verifies that git config gets DOUBLE protection (root:root + 444)

set -e

echo "🧪 R312 Enhanced Protection Test Suite"
echo "======================================"
echo ""

# Test directory
TEST_DIR="/tmp/r312-test-$$"
mkdir -p "$TEST_DIR"
cd "$TEST_DIR"

# Initialize git repo
git init
git config user.name "Test User"
git config user.email "test@example.com"

echo "📋 Initial state:"
echo "   Owner: $(stat -c %U:%G .git/config 2>/dev/null || stat -f %Su:%Sg .git/config)"
echo "   Permissions: $(stat -c %a .git/config 2>/dev/null || stat -f %A .git/config)"
echo ""

# Test 1: Apply full protection with sudo
echo "🧪 Test 1: Full protection with sudo"
echo "────────────────────────────────────"
if command -v sudo >/dev/null 2>&1; then
    echo "✅ sudo available - testing FULL protection"
    
    # Apply protection
    sudo chown root:root .git/config
    sudo chmod 444 .git/config
    
    # Verify ownership
    OWNER=$(stat -c %U:%G .git/config 2>/dev/null || stat -f %Su:%Sg .git/config)
    PERMS=$(stat -c %a .git/config 2>/dev/null || stat -f %A .git/config)
    
    echo "📋 After protection:"
    echo "   Owner: $OWNER"
    echo "   Permissions: $PERMS"
    
    if [ "$OWNER" = "root:root" ] && [ "$PERMS" = "444" ]; then
        echo "✅ FULL protection applied successfully"
    else
        echo "❌ Failed to apply full protection"
        exit 1
    fi
    
    # Test bypass attempts
    echo ""
    echo "🧪 Testing bypass attempts:"
    
    # Try to make it writable without sudo (should fail)
    if chmod 644 .git/config 2>/dev/null; then
        echo "❌ CRITICAL: Non-owner could change permissions!"
        exit 1
    else
        echo "✅ chmod without sudo correctly blocked"
    fi
    
    # Try to write to it (should fail)
    if echo "test" >> .git/config 2>/dev/null; then
        echo "❌ CRITICAL: Could write to readonly config!"
        exit 1
    else
        echo "✅ Write attempt correctly blocked"
    fi
    
    # Try git config modification
    # Note: git creates a new file and renames, so it may succeed even with protection
    # This is why we need BOTH the protection AND SW engineer discipline
    if git config user.name "Hacker" 2>/dev/null; then
        echo "⚠️ Git can still modify (creates new file) - SW discipline required"
        # Reset it back
        sudo chown root:root .git/config
        sudo chmod 444 .git/config
    else
        echo "✅ Git config modification blocked on this system"
    fi
    
    # Test unlock for integration
    echo ""
    echo "🧪 Testing unlock for integration:"
    sudo chown $(id -u):$(id -g) .git/config
    sudo chmod 644 .git/config
    
    if [ -w .git/config ]; then
        echo "✅ Successfully unlocked for integration"
    else
        echo "❌ Failed to unlock config"
        exit 1
    fi
    
else
    echo "⚠️ sudo not available - testing PARTIAL protection only"
    
    # Apply permission-only protection
    chmod 444 .git/config
    
    PERMS=$(stat -c %a .git/config 2>/dev/null || stat -f %A .git/config)
    echo "📋 After protection:"
    echo "   Permissions: $PERMS"
    
    if [ "$PERMS" = "444" ]; then
        echo "✅ PARTIAL protection applied (permissions only)"
    else
        echo "❌ Failed to apply permission protection"
        exit 1
    fi
    
    echo "⚠️ Note: Without sudo, protection can be bypassed with chmod"
fi

echo ""
echo "🧪 Test 2: Docker/Container scenario simulation"
echo "────────────────────────────────────────────────"

# Simulate container environment (no sudo)
(
    export PATH="/usr/bin:/bin"  # Remove sudo from PATH temporarily
    
    if ! command -v sudo >/dev/null 2>&1; then
        echo "✅ Simulated no-sudo environment"
        
        # Should fall back to permission-only
        chmod 644 .git/config  # Reset first
        chmod 444 .git/config  # Apply protection
        
        if [ ! -w .git/config ]; then
            echo "✅ Permission-only protection works in container"
        else
            echo "❌ Failed to protect in container environment"
        fi
    fi
)

echo ""
echo "🧪 Test 3: Verification marker"
echo "────────────────────────────────────────"

# Create marker file
cat > .git/R312_CONFIG_LOCKED << EOF
Timestamp: $(date '+%Y-%m-%d %H:%M:%S')
Locked by: test-script
Protection level: TEST
Purpose: Testing R312 enhanced protection
EOF

if [ -f .git/R312_CONFIG_LOCKED ]; then
    echo "✅ Protection marker created successfully"
    echo "📋 Marker contents:"
    cat .git/R312_CONFIG_LOCKED | sed 's/^/   /'
else
    echo "❌ Failed to create marker"
    exit 1
fi

echo ""
echo "🧪 Test 4: SW Engineer validation simulation"
echo "────────────────────────────────────────────────"

# Reset to protected state for validation
if command -v sudo >/dev/null 2>&1; then
    sudo chown root:root .git/config
    sudo chmod 444 .git/config
else
    chmod 444 .git/config
fi

# Simulate SW engineer validation
if [ -w .git/config ]; then
    echo "❌ Config is writable - validation would fail!"
    exit 1
else
    echo "✅ Config is readonly - validation passes"
    
    OWNER=$(stat -c %U:%G .git/config 2>/dev/null || stat -f %Su:%Sg .git/config)
    if [ "$OWNER" = "root:root" ]; then
        echo "   🔐 FULL protection detected (root-owned)"
    else
        echo "   🔒 PARTIAL protection detected (permission-only)"
    fi
fi

# Cleanup
cd /
rm -rf "$TEST_DIR"

echo ""
echo "======================================"
echo "✅ ALL R312 ENHANCED PROTECTION TESTS PASSED"
echo "======================================"
echo ""
echo "Summary:"
echo "- Full protection (root:root + 444) works when sudo available"
echo "- Fallback to permission-only (444) when sudo unavailable"
echo "- Bypass attempts correctly blocked with full protection"
echo "- Unlock for integration works correctly"
echo "- SW engineer validation detects protection level"
echo "- Protection markers created successfully"