#!/bin/bash
# Test script to verify line-counter.sh handles complex branch patterns

set -uo pipefail

echo "================================================================"
echo "Line Counter Pattern Testing Suite"
echo "================================================================"
echo ""

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

# Test patterns that should match
echo "Testing branch pattern matching..."
echo "--------------------------------"

test_patterns=(
    # Simple patterns (branch:prefix:type:phase:wave:effort)
    "phase1/wave1/api:none:effort:1:1:api"
    "phase2/wave3/feature:none:effort:2:3:feature"
    "phase10/wave99/test:none:effort:10:99:test"
    
    # With project prefixes
    "my-project/phase1/wave1/api:my-project:effort:1:1:api"
    "idpbuilder-oci-go-cr/phase2/wave1/go-containerregistry-image-builder:idpbuilder-oci-go-cr:effort:2:1:go-containerregistry-image-builder"
    "complex-project-name/phase3/wave2/complex-feature-name:complex-project-name:effort:3:2:complex-feature-name"
    "a-b-c-d-e-f/phase10/wave99/x-y-z-1-2-3:a-b-c-d-e-f:effort:10:99:x-y-z-1-2-3"
    
    # Split branches
    "phase1/wave1/api--split-001:none:split:api:001:"
    "my-project/phase2/wave1/feature--split-003:my-project:split:feature:003:"
    
    # Integration branches
    "phase1-wave1-integration:none:wave-integration:1:1:"
    "phase2-integration:none:phase-integration:2::"
    "my-project/phase3-wave2-integration:my-project:wave-integration:3:2:"
    "my-project/phase4-integration:my-project:phase-integration:4::"
)

success_count=0
failure_count=0

for test_case in "${test_patterns[@]}"; do
    # Split test case into branch and expected values
    IFS=':' read -r branch expected_prefix expected_type expected_phase expected_wave expected_effort <<< "$test_case"
    
    echo -n "Testing: $branch ... "
    
    # Test effort pattern
    if [[ "$expected_type" == "split" ]]; then
        if [[ "$branch" =~ (.*)--split-([0-9]+) ]]; then
            echo -e "${GREEN}✓ Split pattern matched${NC}"
            ((success_count++))
        else
            echo -e "${RED}✗ Split pattern failed${NC}"
            ((failure_count++))
        fi
    elif [[ "$expected_type" == "wave-integration" ]] || [[ "$expected_type" == "phase-integration" ]]; then
        # Test integration patterns
        if [[ "$branch" =~ (([^/]+)/)?phase([0-9]+)-(wave([0-9]+)-)?integration$ ]]; then
            echo -e "${GREEN}✓ Integration pattern matched${NC}"
            ((success_count++))
        else
            echo -e "${RED}✗ Integration pattern failed${NC}"
            ((failure_count++))
        fi
    else
        # Test effort pattern
        if [[ "$branch" =~ ^(([^/]+)/)?phase([0-9]+)/wave([0-9]+)/([^/]+)$ ]]; then
            actual_prefix="${BASH_REMATCH[2]:-none}"
            actual_phase="${BASH_REMATCH[3]}"
            actual_wave="${BASH_REMATCH[4]}"
            actual_effort="${BASH_REMATCH[5]}"
            
            if [[ "$actual_prefix" == "$expected_prefix" ]] && \
               [[ "$actual_phase" == "$expected_phase" ]] && \
               [[ "$actual_wave" == "$expected_wave" ]] && \
               [[ "$actual_effort" == "$expected_effort" ]]; then
                echo -e "${GREEN}✓ Pattern matched correctly${NC}"
                echo "    Prefix: $actual_prefix, Phase: $actual_phase, Wave: $actual_wave, Effort: $actual_effort"
                ((success_count++))
            else
                echo -e "${RED}✗ Pattern matched but values incorrect${NC}"
                echo "    Expected: prefix=$expected_prefix, phase=$expected_phase, wave=$expected_wave, effort=$expected_effort"
                echo "    Got:      prefix=$actual_prefix, phase=$actual_phase, wave=$actual_wave, effort=$actual_effort"
                ((failure_count++))
            fi
        else
            echo -e "${RED}✗ Pattern did not match${NC}"
            ((failure_count++))
        fi
    fi
done

echo ""
echo "================================================================"
echo "Test Results:"
echo "  Passed: $success_count"
echo "  Failed: $failure_count"
echo "================================================================"

if [ $failure_count -eq 0 ]; then
    echo -e "${GREEN}All tests passed!${NC}"
    exit 0
else
    echo -e "${RED}Some tests failed!${NC}"
    exit 1
fi