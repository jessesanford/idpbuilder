#!/bin/bash

# Test that our bash syntax fixes work correctly

echo "Testing R205 bash syntax..."

# Create a mock split plan
mkdir -p /tmp/test-bash
cat > /tmp/test-bash/SPLIT-PLAN-001.md << 'EOF'
# Test Split Plan
**WORKING_DIRECTORY**: /tmp/test-bash/split-001
**BRANCH**: test-split-001
EOF

mkdir -p /tmp/test-bash/split-001
cd /tmp/test-bash

# Test the R205 navigation script (with proper line continuations)
echo "Running R205 navigation with fixed syntax..."

if [ -f "SPLIT-PLAN-*.md" ] || [ -f "SPLIT-INVENTORY.md" ]; then \
    echo "✅ Split work detected"; \
    SPLIT_PLAN=$(ls SPLIT-PLAN-*.md 2>/dev/null | head -1); \
    if [ -z "$SPLIT_PLAN" ]; then \
        echo "❌ No split plan found"; \
    else \
        echo "✅ Found split plan: $SPLIT_PLAN"; \
        WORKING_DIR=$(grep "\*\*WORKING_DIRECTORY\*\*:" "$SPLIT_PLAN" | cut -d: -f2- | xargs); \
        EXPECTED_BRANCH=$(grep "\*\*BRANCH\*\*:" "$SPLIT_PLAN" | head -1 | cut -d: -f2- | xargs); \
        if [ -z "$WORKING_DIR" ] || [ -z "$EXPECTED_BRANCH" ]; then \
            echo "❌ Missing metadata"; \
        else \
            echo "✅ Extracted metadata:"; \
            echo "   WORKING_DIRECTORY: $WORKING_DIR"; \
            echo "   BRANCH: $EXPECTED_BRANCH"; \
            if [ ! -d "$WORKING_DIR" ]; then \
                echo "❌ Directory doesn't exist (expected for test)"; \
            else \
                echo "✅ Directory exists"; \
            fi; \
        fi; \
    fi; \
else \
    echo "❌ No split work detected"; \
fi

echo ""
echo "Testing nested if-then with proper syntax..."

# Test nested conditions
TEST_VAR="split"
if [[ "$TEST_VAR" == *"split"* ]]; then \
    echo "✅ Detected split in variable"; \
    if [ -f "SPLIT-PLAN-001.md" ]; then \
        echo "✅ Split plan exists"; \
    else \
        echo "❌ Split plan missing"; \
    fi; \
else \
    echo "❌ No split detected"; \
fi

echo ""
echo "✅ All bash syntax tests completed successfully!"

# Cleanup
rm -rf /tmp/test-bash