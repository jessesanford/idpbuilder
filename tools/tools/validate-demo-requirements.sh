#!/bin/bash
# Validates demo requirements across all integrations per R330/R291

set -euo pipefail

echo "🔍 Validating Demo Requirements (R330/R291)"
echo "=========================================="

VIOLATIONS=0
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# Check wave integrations
echo ""
echo "📊 Checking Wave Integrations..."
for wave_int in "$PROJECT_ROOT"/efforts/phase*/wave*/integration-workspace; do
    if [ -d "$wave_int" ]; then
        wave=$(basename $(dirname "$wave_int"))
        phase=$(basename $(dirname $(dirname "$wave_int")))

        echo ""
        echo "Checking $phase/$wave integration..."

        # Check for demo script
        if ! ls "$wave_int"/demo-*.sh 1>/dev/null 2>&1; then
            echo "  ❌ Missing demo script"
            ((VIOLATIONS++))
        else
            echo "  ✅ Demo script found"

            # Check if executable
            for script in "$wave_int"/demo-*.sh; do
                if [ ! -x "$script" ]; then
                    echo "  ⚠️  Demo script not executable: $(basename "$script")"
                fi
            done
        fi

        # Check for demo documentation
        if ! ls "$wave_int"/*DEMO*.md 1>/dev/null 2>&1 && \
           ! ls "$wave_int"/DEMO.md 1>/dev/null 2>&1; then
            echo "  ❌ Missing demo documentation"
            ((VIOLATIONS++))
        else
            echo "  ✅ Demo documentation found"
        fi

        # Check for DEMO-STATUS.md (indicates R291 failure)
        if [ -f "$wave_int/DEMO-STATUS.md" ]; then
            if grep -q "FAILED" "$wave_int/DEMO-STATUS.md"; then
                echo "  🔴 DEMO-STATUS.md indicates R291 gate FAILURE!"
                ((VIOLATIONS++))
            fi
        fi
    fi
done

# Check phase integrations
echo ""
echo "📊 Checking Phase Integrations..."
for phase_int in "$PROJECT_ROOT"/efforts/phase*/integration; do
    if [ -d "$phase_int" ]; then
        phase=$(basename $(dirname "$phase_int"))

        echo ""
        echo "Checking $phase integration..."

        # Check for demo script
        if ! ls "$phase_int"/demo-*.sh 1>/dev/null 2>&1; then
            echo "  ❌ Missing demo script"
            ((VIOLATIONS++))
        else
            echo "  ✅ Demo script found"

            # Check if executable
            for script in "$phase_int"/demo-*.sh; do
                if [ ! -x "$script" ]; then
                    echo "  ⚠️  Demo script not executable: $(basename "$script")"
                fi
            done
        fi

        # Check for demo documentation
        if ! ls "$phase_int"/*DEMO*.md 1>/dev/null 2>&1 && \
           ! ls "$phase_int"/DEMO.md 1>/dev/null 2>&1; then
            echo "  ❌ Missing demo documentation"
            ((VIOLATIONS++))
        else
            echo "  ✅ Demo documentation found"
        fi

        # Check for DEMO-STATUS.md (indicates R291 failure)
        if [ -f "$phase_int/DEMO-STATUS.md" ]; then
            if grep -q "FAILED" "$phase_int/DEMO-STATUS.md"; then
                echo "  🔴 DEMO-STATUS.md indicates R291 gate FAILURE!"
                ((VIOLATIONS++))
            fi
        fi
    fi
done

# Check project integration (if exists)
echo ""
echo "📊 Checking Project Integration..."
if [ -d "$PROJECT_ROOT/integration" ]; then
    # Check for demo script
    if ! ls "$PROJECT_ROOT/integration"/demo-*.sh 1>/dev/null 2>&1; then
        echo "  ❌ Missing project demo script"
        ((VIOLATIONS++))
    else
        echo "  ✅ Project demo script found"
    fi

    # Check for demo documentation
    if ! ls "$PROJECT_ROOT/integration"/*DEMO*.md 1>/dev/null 2>&1 && \
       ! ls "$PROJECT_ROOT/integration"/DEMO.md 1>/dev/null 2>&1; then
        echo "  ❌ Missing project demo documentation"
        ((VIOLATIONS++))
    else
        echo "  ✅ Project demo documentation found"
    fi
fi

# Summary
echo ""
echo "=========================================="
if [ $VIOLATIONS -gt 0 ]; then
    echo "🔴 R291/R330 VIOLATIONS: Found $VIOLATIONS missing demo artifacts"
    echo ""
    echo "ALL integrations MUST have:"
    echo "  1. Demo script (demo-*.sh) that is executable"
    echo "  2. Demo documentation (DEMO.md, *-DEMO.md, etc.)"
    echo "  3. No DEMO-STATUS.md with FAILED status"
    echo ""
    echo "See:"
    echo "  - rule-library/R330-demo-planning-requirements.md"
    echo "  - rule-library/R291-integration-demo-requirement.md"
    exit 1
else
    echo "✅ All integrations have demo requirements fulfilled"
    echo ""
    echo "Validated:"
    echo "  - Wave integration demos"
    echo "  - Phase integration demos"
    echo "  - Project integration demos (if exists)"
fi

exit 0
