#!/bin/bash
# Phase 2 Wave 2 Integration Demo Script

echo "🎬 Running Phase 2 Wave 2 Wave Demo..."
echo "========================="
echo

echo "📦 Wave 2 Components:"
echo "  - Auth Flow Implementation (effort 2.2.2)"
echo "  - Flow Tests (effort 2.2.1)"
echo
echo "Note: Tests cannot compile due to interface mismatch"
echo "      This is documented as an upstream bug"
echo

# Run auth-flow demo
if [[ -x "./demo-auth-flow.sh" ]]; then
    echo "Running auth-flow demo..."
    ./demo-auth-flow.sh --with-flags
    echo
    ./demo-auth-flow.sh --with-secrets  
    echo
fi

echo "========================="
echo "✅ Wave demo completed (with noted test compilation issues)"
