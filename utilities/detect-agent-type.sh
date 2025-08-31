#!/bin/bash
# Detect agent type from various context clues

# Method 1: Check recent command history for agent-specific commands
RECENT_COMMANDS=$(history 20 2>/dev/null | tail -20)

# Method 2: Check current directory structure
CURRENT_DIR=$(pwd)

# Method 3: Check for agent-specific files in current directory
AGENT_TYPE="unknown"

# Check for orchestrator patterns
if echo "$CURRENT_DIR" | grep -q "orchestrator\|agent-configs"; then
    AGENT_TYPE="orchestrator"
elif [ -f "orchestrator-state.yaml" ] || [ -f "../orchestrator-state.yaml" ]; then
    AGENT_TYPE="orchestrator"
elif echo "$RECENT_COMMANDS" | grep -q "continue-orchestrating\|/orchestrator"; then
    AGENT_TYPE="orchestrator"

# Check for SW engineer patterns
elif echo "$CURRENT_DIR" | grep -q "efforts/phase[0-9]*/wave[0-9]*/"; then
    AGENT_TYPE="sw-engineer"
elif [ -f "pkg/main.go" ] || [ -d "pkg" ]; then
    AGENT_TYPE="sw-engineer"
elif echo "$RECENT_COMMANDS" | grep -q "sw-eng\|implementation"; then
    AGENT_TYPE="sw-engineer"

# Check for code reviewer patterns
elif [ -f "IMPLEMENTATION-PLAN.md" ] || [ -f "REVIEW-FEEDBACK.md" ]; then
    AGENT_TYPE="code-reviewer"
elif echo "$RECENT_COMMANDS" | grep -q "code-review\|/reviewer"; then
    AGENT_TYPE="code-reviewer"

# Check for architect patterns
elif [ -f "ARCHITECTURE-REVIEW.md" ] || [ -f "PHASE-ASSESSMENT.md" ]; then
    AGENT_TYPE="architect"
elif echo "$RECENT_COMMANDS" | grep -q "architect\|/architecture"; then
    AGENT_TYPE="architect"
fi

# Also save process info that might help
PARENT_PID=$(ps -o ppid= -p $$ | tr -d ' ')
PROCESS_CMD=$(ps -o cmd= -p $PARENT_PID 2>/dev/null | head -1)

# Output both the detected type and additional context
echo "$AGENT_TYPE"