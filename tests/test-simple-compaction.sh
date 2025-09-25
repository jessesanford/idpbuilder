#!/bin/bash

# Simple test to verify the compaction check command works

echo "Testing simplified compaction detection command..."

# The command as agents should run it (single line, no backslashes)
CMD='if [ -f /tmp/compaction_marker.txt ]; then echo "Compaction detected"; cat /tmp/compaction_marker.txt; rm -f /tmp/compaction_marker.txt; echo "TODO recovery required"; else echo "No compaction detected"; fi'

# Test with marker
touch /tmp/compaction_marker.txt
echo "TEST_MARKER" > /tmp/compaction_marker.txt

echo "Test 1: With marker..."
if bash -c "$CMD" | grep -q "Compaction detected"; then
    echo "✅ PASS - Detected compaction marker"
else
    echo "❌ FAIL - Did not detect marker"
fi

# Test without marker  
echo ""
echo "Test 2: Without marker..."
if bash -c "$CMD" | grep -q "No compaction detected"; then
    echo "✅ PASS - No compaction detected"
else
    echo "❌ FAIL - Incorrect response without marker"
fi

# Clean up
rm -f /tmp/compaction_marker.txt

echo ""
echo "This is the working command agents should use:"
echo "$CMD"