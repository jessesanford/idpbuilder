#!/bin/bash
# API Compatibility Verification Script for Phase 2 Wave 2
# This script checks for old API calls and ensures new API is present

echo "========================================="
echo "API Compatibility Verification"
echo "========================================="

ERRORS=0
WARNINGS=0

# Check for FORBIDDEN old API calls
echo "Checking for forbidden old API calls..."

# Certificate API checks
if grep -r "ExtractCertificate" pkg/ 2>/dev/null | grep -v "// OLD" | grep -v "DEPRECATED"; then
    echo "❌ ERROR: Old ExtractCertificate API found!"
    ERRORS=$((ERRORS + 1))
fi

# Gitea client checks
if grep -r "NewClientAutoDetect" pkg/ 2>/dev/null; then
    echo "❌ ERROR: Non-existent NewClientAutoDetect found!"
    ERRORS=$((ERRORS + 1))
fi

if grep -r "PushImage" pkg/ 2>/dev/null | grep -v "// OLD" | grep -v "DEPRECATED"; then
    echo "❌ ERROR: Old PushImage method found!"
    ERRORS=$((ERRORS + 1))
fi

if grep -r "ValidateCredentials" pkg/ 2>/dev/null; then
    echo "❌ ERROR: Non-existent ValidateCredentials method found!"
    ERRORS=$((ERRORS + 1))
fi

# Check for 4-parameter NewClient (old signature)
if grep -r "NewClient.*,.*,.*,.*)" pkg/gitea/ 2>/dev/null; then
    echo "❌ ERROR: Old 4-parameter NewClient signature found!"
    ERRORS=$((ERRORS + 1))
fi

echo ""
echo "Checking for required new API..."

# Check for REQUIRED new API calls
if ! grep -r "NewTrustStore" pkg/ 2>/dev/null > /dev/null; then
    echo "❌ ERROR: NewTrustStore API missing!"
    ERRORS=$((ERRORS + 1))
fi

if ! grep -r "func.*Push(" pkg/gitea/ 2>/dev/null > /dev/null; then
    echo "❌ ERROR: Push method missing in gitea client!"
    ERRORS=$((ERRORS + 1))
fi

# Check specific file signatures
echo ""
echo "Checking specific file signatures..."

# Check push.go
if [ -f "pkg/cmd/push.go" ]; then
    if grep -q "ExtractCertificate" pkg/cmd/push.go; then
        echo "❌ ERROR: push.go contains old ExtractCertificate call"
        ERRORS=$((ERRORS + 1))
    fi
    if ! grep -q "NewTrustStore" pkg/cmd/push.go; then
        echo "⚠️  WARNING: push.go doesn't use NewTrustStore"
        WARNINGS=$((WARNINGS + 1))
    fi
fi

# Check gitea client
if [ -f "pkg/gitea/client.go" ]; then
    # Check for correct NewClient signature (2 params)
    if ! grep -q "func NewClient(registryURL string, certManager" pkg/gitea/client.go; then
        echo "❌ ERROR: gitea client has wrong NewClient signature"
        ERRORS=$((ERRORS + 1))
    fi
    # Check for Push method
    if ! grep -q "func.*Push(" pkg/gitea/client.go; then
        echo "❌ ERROR: gitea client missing Push method"
        ERRORS=$((ERRORS + 1))
    fi
fi

echo ""
echo "========================================="
echo "VERIFICATION RESULTS"
echo "========================================="

if [ $ERRORS -eq 0 ] && [ $WARNINGS -eq 0 ]; then
    echo "✅ ALL CHECKS PASSED - API is compatible!"
    echo "Integration is safe to proceed."
    exit 0
elif [ $ERRORS -eq 0 ]; then
    echo "⚠️  Verification passed with $WARNINGS warnings"
    echo "Review warnings but integration can proceed."
    exit 0
else
    echo "❌ VERIFICATION FAILED"
    echo "Found $ERRORS critical errors and $WARNINGS warnings"
    echo "DO NOT PROCEED with integration until fixed!"
    exit 1
fi