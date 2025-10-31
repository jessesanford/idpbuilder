#!/bin/bash

# Test the cascade calculation logic directly
# This bypasses the parse script to verify our cascade implementation

# Copy the calculate_cascade_base function for testing
calculate_cascade_base() {
    local phase="$1"
    local wave="$2"
    local effort_name="$3"
    local effort_index="$4"
    local previous_effort_branch="$5"
    local previous_wave_last_branch="$6"
    local previous_phase_last_branch="$7"

    local base_branch=""
    local cascade_reason=""

    # R308/R501 CASCADE ALGORITHM
    if [[ "$phase" == "phase1" && "$wave" == "wave1" && "$effort_index" -eq 1 ]]; then
        base_branch="main"
        cascade_reason="First effort of Phase 1 Wave 1 starts from main (R308/R501)"
    elif [[ "$effort_index" -eq 1 ]]; then
        if [[ "$wave" == "wave1" ]]; then
            if [[ -n "$previous_phase_last_branch" && "$previous_phase_last_branch" != "null" ]]; then
                base_branch="$previous_phase_last_branch"
                cascade_reason="First effort of $phase $wave cascades from previous phase's last effort: $previous_phase_last_branch (R308/R501)"
            else
                base_branch="main"
                cascade_reason="No previous phase found, defaulting to main (R308/R501)"
            fi
        else
            if [[ -n "$previous_wave_last_branch" && "$previous_wave_last_branch" != "null" ]]; then
                base_branch="$previous_wave_last_branch"
                cascade_reason="First effort of $phase $wave cascades from previous wave's last effort: $previous_wave_last_branch (R308/R501)"
            else
                echo "{\"base\": \"ERROR\", \"reason\": \"CASCADE VIOLATION: Cannot find previous wave's last effort for $phase $wave\"}"
                return
            fi
        fi
    else
        if [[ -n "$previous_effort_branch" && "$previous_effort_branch" != "null" ]]; then
            base_branch="$previous_effort_branch"
            cascade_reason="Effort $effort_index cascades from previous effort: $previous_effort_branch (R308/R501)"
        else
            echo "{\"base\": \"ERROR\", \"reason\": \"CASCADE VIOLATION: Cannot find previous effort for $effort_name (index $effort_index)\"}"
            return
        fi
    fi

    echo "{\"base\": \"$base_branch\", \"reason\": \"$cascade_reason\"}"
}

echo "Testing CASCADE calculation per R308/R501"
echo "========================================="
echo ""

# Test Case 1: First effort of P1W1 should be main
result=$(calculate_cascade_base "phase1" "wave1" "auth" 1 "" "" "")
base=$(echo "$result" | jq -r '.base')
reason=$(echo "$result" | jq -r '.reason')
echo "Test 1 - P1W1 first effort:"
echo "  Expected: main"
echo "  Got: $base"
echo "  Reason: $reason"
[[ "$base" == "main" ]] && echo "  ✅ PASS" || echo "  ❌ FAIL"
echo ""

# Test Case 2: Second effort of P1W1 should cascade from first
result=$(calculate_cascade_base "phase1" "wave1" "users" 2 "project/phase1/wave1/auth" "" "")
base=$(echo "$result" | jq -r '.base')
reason=$(echo "$result" | jq -r '.reason')
echo "Test 2 - P1W1 second effort:"
echo "  Expected: project/phase1/wave1/auth"
echo "  Got: $base"
echo "  Reason: $reason"
[[ "$base" == "project/phase1/wave1/auth" ]] && echo "  ✅ PASS" || echo "  ❌ FAIL"
echo ""

# Test Case 3: Third effort of P1W1 should cascade from second
result=$(calculate_cascade_base "phase1" "wave1" "perms" 3 "project/phase1/wave1/users" "" "")
base=$(echo "$result" | jq -r '.base')
reason=$(echo "$result" | jq -r '.reason')
echo "Test 3 - P1W1 third effort:"
echo "  Expected: project/phase1/wave1/users"
echo "  Got: $base"
echo "  Reason: $reason"
[[ "$base" == "project/phase1/wave1/users" ]] && echo "  ✅ PASS" || echo "  ❌ FAIL"
echo ""

# Test Case 4: First effort of P1W2 should cascade from last of P1W1
result=$(calculate_cascade_base "phase1" "wave2" "api" 1 "" "project/phase1/wave1/perms" "")
base=$(echo "$result" | jq -r '.base')
reason=$(echo "$result" | jq -r '.reason')
echo "Test 4 - P1W2 first effort (cross-wave):"
echo "  Expected: project/phase1/wave1/perms"
echo "  Got: $base"
echo "  Reason: $reason"
[[ "$base" == "project/phase1/wave1/perms" ]] && echo "  ✅ PASS" || echo "  ❌ FAIL"
echo ""

# Test Case 5: Second effort of P1W2 should cascade from first of P1W2
result=$(calculate_cascade_base "phase1" "wave2" "docs" 2 "project/phase1/wave2/api" "project/phase1/wave1/perms" "")
base=$(echo "$result" | jq -r '.base')
reason=$(echo "$result" | jq -r '.reason')
echo "Test 5 - P1W2 second effort:"
echo "  Expected: project/phase1/wave2/api"
echo "  Got: $base"
echo "  Reason: $reason"
[[ "$base" == "project/phase1/wave2/api" ]] && echo "  ✅ PASS" || echo "  ❌ FAIL"
echo ""

# Test Case 6: First effort of P2W1 should cascade from last of P1
result=$(calculate_cascade_base "phase2" "wave1" "ui" 1 "" "" "project/phase1/wave2/docs")
base=$(echo "$result" | jq -r '.base')
reason=$(echo "$result" | jq -r '.reason')
echo "Test 6 - P2W1 first effort (cross-phase):"
echo "  Expected: project/phase1/wave2/docs"
echo "  Got: $base"
echo "  Reason: $reason"
[[ "$base" == "project/phase1/wave2/docs" ]] && echo "  ✅ PASS" || echo "  ❌ FAIL"
echo ""

echo "========================================="
echo "CASCADE CHAIN VISUALIZATION:"
echo ""
echo "main"
echo "  └─→ project/phase1/wave1/auth (P1W1 effort 1)"
echo "        └─→ project/phase1/wave1/users (P1W1 effort 2)"
echo "              └─→ project/phase1/wave1/perms (P1W1 effort 3)"
echo "                    └─→ project/phase1/wave2/api (P1W2 effort 1)"
echo "                          └─→ project/phase1/wave2/docs (P1W2 effort 2)"
echo "                                └─→ project/phase2/wave1/ui (P2W1 effort 1)"
echo ""
echo "✅ This is PROGRESSIVE trunk-based development!"
echo "========================================="