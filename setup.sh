#!/bin/bash

# Software Factory Template Setup Script
# This script helps you set up the software factory for your project

set -e

echo "========================================="
echo "Software Factory Template Setup"
echo "========================================="
echo ""

# Get project name
read -p "Enter your project name (e.g., my-awesome-app): " PROJECT_NAME
if [ -z "$PROJECT_NAME" ]; then
    echo "Error: Project name is required"
    exit 1
fi

# Get target directory
read -p "Enter target directory (default: /workspaces/$PROJECT_NAME): " TARGET_DIR
if [ -z "$TARGET_DIR" ]; then
    TARGET_DIR="/workspaces/$PROJECT_NAME"
fi

# Get primary language
echo ""
echo "Select your primary programming language:"
echo "1) Go"
echo "2) Python"
echo "3) JavaScript/TypeScript"
echo "4) Java"
echo "5) Rust"
echo "6) C/C++"
echo "7) Other"
read -p "Enter choice (1-7): " LANG_CHOICE

case $LANG_CHOICE in
    1) PRIMARY_LANG="go"
       FILE_PATTERNS='**/*.go'
       EXCLUDE_PATTERNS='**/vendor/**,**/*_test.go,**/zz_generated*.go,**/*.pb.go'
       ;;
    2) PRIMARY_LANG="python"
       FILE_PATTERNS='**/*.py'
       EXCLUDE_PATTERNS='**/venv/**,**/__pycache__/**,**/tests/**,**/*_test.py'
       ;;
    3) PRIMARY_LANG="javascript"
       FILE_PATTERNS='**/*.js,**/*.jsx,**/*.ts,**/*.tsx'
       EXCLUDE_PATTERNS='**/node_modules/**,**/dist/**,**/build/**,**/*.test.*,**/*.spec.*'
       ;;
    4) PRIMARY_LANG="java"
       FILE_PATTERNS='**/*.java'
       EXCLUDE_PATTERNS='**/target/**,**/build/**,**/test/**,**/*Test.java'
       ;;
    5) PRIMARY_LANG="rust"
       FILE_PATTERNS='**/*.rs'
       EXCLUDE_PATTERNS='**/target/**,**/tests/**'
       ;;
    6) PRIMARY_LANG="c"
       FILE_PATTERNS='**/*.c,**/*.cpp,**/*.h,**/*.hpp'
       EXCLUDE_PATTERNS='**/build/**,**/tests/**,**/*_test.*'
       ;;
    *) PRIMARY_LANG="other"
       FILE_PATTERNS='**/*'
       EXCLUDE_PATTERNS=''
       ;;
esac

# Get complexity level
echo ""
echo "Select project complexity:"
echo "1) Simple (1-2 developers, <10K lines)"
echo "2) Medium (3-10 developers, 10K-50K lines)"
echo "3) Complex (10+ developers, 50K+ lines)"
read -p "Enter choice (1-3): " COMPLEXITY

# Create target directory
echo ""
echo "Creating project structure in $TARGET_DIR..."
mkdir -p "$TARGET_DIR"

# Copy template files
echo "Copying template files..."
cp -r . "$TARGET_DIR/" 2>/dev/null || true

# Clean up the copied setup script
rm -f "$TARGET_DIR/setup.sh"

# Update paths in all files
echo "Updating paths for $PROJECT_NAME..."
find "$TARGET_DIR" -type f -name "*.md" -o -name "*.yaml" -o -name "*.sh" | while read file; do
    # Update project-specific paths
    sed -i "s|/workspaces/software-factory-template|$TARGET_DIR|g" "$file"
    sed -i "s|/workspaces/\[project\]|$TARGET_DIR|g" "$file"
    sed -i "s|\[project\]|$PROJECT_NAME|g" "$file"
done

# Configure critical settings.json file
echo "Configuring critical settings.json for compaction recovery..."
SETTINGS_FILE="$TARGET_DIR/.claude/settings.json"
if [ -f "$SETTINGS_FILE" ]; then
    # Update TODO directory path in settings.json
    sed -i "s|TODO_DIR='./todos'|TODO_DIR='$TARGET_DIR/todos'|g" "$SETTINGS_FILE"
    echo "✅ settings.json configured for TODO preservation"
else
    echo "⚠️ WARNING: settings.json not found - compaction recovery will not work!"
fi

# Configure line counter
echo "Configuring line counter for $PRIMARY_LANG..."
LINE_COUNTER="$TARGET_DIR/tools/line-counter.sh"
if [ -f "$LINE_COUNTER" ]; then
    # Update file patterns
    sed -i "s|FILE_PATTERNS=.*|FILE_PATTERNS='$FILE_PATTERNS'|g" "$LINE_COUNTER"
    sed -i "s|EXCLUDE_PATTERNS=.*|EXCLUDE_PATTERNS='$EXCLUDE_PATTERNS'|g" "$LINE_COUNTER"
    chmod +x "$LINE_COUNTER"
fi

# Customize SW engineer agent for language
SW_ENG_FILE="$TARGET_DIR/.claude/agents/sw-engineer-example-go.md"
NEW_SW_ENG_FILE="$TARGET_DIR/.claude/agents/sw-engineer-$PRIMARY_LANG.md"
if [ -f "$SW_ENG_FILE" ]; then
    mv "$SW_ENG_FILE" "$NEW_SW_ENG_FILE"
    
    # Update agent file with language-specific guidance
    case $PRIMARY_LANG in
        go)
            echo "# Already configured for Go"
            ;;
        python)
            sed -i 's/Go/Python/g' "$NEW_SW_ENG_FILE"
            sed -i 's/goroutines/async\/await/g' "$NEW_SW_ENG_FILE"
            sed -i 's/channels/queues/g' "$NEW_SW_ENG_FILE"
            ;;
        javascript)
            sed -i 's/Go/JavaScript\/TypeScript/g' "$NEW_SW_ENG_FILE"
            sed -i 's/goroutines/async\/await/g' "$NEW_SW_ENG_FILE"
            sed -i 's/channels/promises/g' "$NEW_SW_ENG_FILE"
            ;;
        *)
            sed -i "s/Go/$PRIMARY_LANG/g" "$NEW_SW_ENG_FILE"
            ;;
    esac
fi

# Activate optional files based on complexity
case $COMPLEXITY in
    2)
        echo "Activating medium complexity protocols..."
        if [ -d "$TARGET_DIR/possibly-needed-but-not-sure" ]; then
            cp "$TARGET_DIR/possibly-needed-but-not-sure/CODE-REVIEWER-QUICK-REFERENCE.md" "$TARGET_DIR/protocols/" 2>/dev/null || true
            cp "$TARGET_DIR/possibly-needed-but-not-sure/ORCHESTRATOR-QUICK-REFERENCE.md" "$TARGET_DIR/protocols/" 2>/dev/null || true
            cp "$TARGET_DIR/possibly-needed-but-not-sure/ORCHESTRATOR-WORKFLOW-SUMMARY.md" "$TARGET_DIR/protocols/" 2>/dev/null || true
        fi
        ;;
    3)
        echo "Activating complex project protocols..."
        if [ -d "$TARGET_DIR/possibly-needed-but-not-sure" ]; then
            # Copy all optional files to active locations
            cp "$TARGET_DIR/possibly-needed-but-not-sure/"*.md "$TARGET_DIR/protocols/" 2>/dev/null || true
        fi
        ;;
esac

# Create initial directories
echo "Creating working directories..."
mkdir -p "$TARGET_DIR/efforts"
mkdir -p "$TARGET_DIR/todos"
mkdir -p "$TARGET_DIR/orchestrator"

# Create initial project implementation plan template
cat > "$TARGET_DIR/orchestrator/PROJECT-IMPLEMENTATION-PLAN.md" << 'EOF'
# Project Implementation Plan

## Project: PROJECT_NAME_PLACEHOLDER

### Phase 1: Foundation
**Goal**: Establish core architecture and basic functionality

#### Wave 1: Core Types and Interfaces
- E1.1.1: Define basic data models
- E1.1.2: Create API interfaces
- E1.1.3: Setup configuration structures

#### Wave 2: Core Implementation
- E1.2.1: Implement business logic
- E1.2.2: Add data persistence
- E1.2.3: Create service layer

#### Wave 3: Testing and Validation
- E1.3.1: Unit test coverage
- E1.3.2: Integration tests
- E1.3.3: Validation framework

### Phase 2: Features
**Goal**: Implement primary features

#### Wave 1: Feature Set A
- E2.1.1: [Define your features]

### Phase 3: Optimization
**Goal**: Performance and reliability

#### Wave 1: Performance
- E3.1.1: [Define optimization goals]

## Success Criteria
- [ ] All efforts under 800 lines
- [ ] 100% code review coverage
- [ ] Test coverage meets requirements
- [ ] Architecture reviews pass
- [ ] Integration successful
EOF

sed -i "s/PROJECT_NAME_PLACEHOLDER/$PROJECT_NAME/g" "$TARGET_DIR/orchestrator/PROJECT-IMPLEMENTATION-PLAN.md"

# Create initial state file
cat > "$TARGET_DIR/orchestrator/orchestrator-state.yaml" << 'EOF'
# Orchestrator State File
# This file tracks the current state of the orchestration

current_phase: 1
current_wave: 1
current_state: "INIT"

efforts_completed: []
efforts_in_progress: []
efforts_pending:
  - "E1.1.1"
  - "E1.1.2"
  - "E1.1.3"

integration_branches: []

metadata:
  project_name: "PROJECT_NAME_PLACEHOLDER"
  created_at: "DATE_PLACEHOLDER"
  last_updated: "DATE_PLACEHOLDER"
EOF

CURRENT_DATE=$(date -Iseconds)
sed -i "s/PROJECT_NAME_PLACEHOLDER/$PROJECT_NAME/g" "$TARGET_DIR/orchestrator/orchestrator-state.yaml"
sed -i "s/DATE_PLACEHOLDER/$CURRENT_DATE/g" "$TARGET_DIR/orchestrator/orchestrator-state.yaml"

# Create a quick start guide
cat > "$TARGET_DIR/QUICK-START.md" << EOF
# Quick Start Guide for $PROJECT_NAME

## Setup Complete! 🎉

Your software factory has been configured with:
- **Language**: $PRIMARY_LANG
- **Complexity**: Level $COMPLEXITY
- **Location**: $TARGET_DIR

## Next Steps

### 1. Review and Customize Your Plan
Edit \`orchestrator/PROJECT-IMPLEMENTATION-PLAN.md\` to define your specific:
- Phases and goals
- Waves and groupings
- Individual efforts (keep each under 800 lines!)

### 2. Configure Your Agents
The following agents are ready:
- Orchestrator: \`.claude/agents/orchestrator-task-master.md\`
- Code Reviewer: \`.claude/agents/code-reviewer.md\`
- SW Engineer: \`.claude/agents/sw-engineer-$PRIMARY_LANG.md\`
- Architect: \`.claude/agents/architect-reviewer.md\`

### 3. Start Orchestration
\`\`\`bash
cd $TARGET_DIR
# In Claude Code, run:
/continue-orchestrating
\`\`\`

### 4. Monitor Progress
- State tracking: \`orchestrator/orchestrator-state.yaml\`
- TODOs: \`todos/\` directory
- Work logs: Each effort directory

## Key Commands

### Measure Lines (Always Use This!)
\`\`\`bash
$TARGET_DIR/tools/line-counter.sh -c branch-name
\`\`\`

### Check State
\`\`\`bash
cat orchestrator/orchestrator-state.yaml
\`\`\`

## Critical Configuration

### 🔴 settings.json is ESSENTIAL
The `.claude/settings.json` file enables:
- Compaction recovery
- TODO state preservation  
- Context maintenance

**Without it, you WILL lose work during compaction!**

## Important Rules

1. **NEVER exceed 800 lines per effort** - Split if needed
2. **ALWAYS review before proceeding** - No skipping reviews
3. **SEQUENTIAL splits only** - Never parallel
4. **Orchestrator NEVER writes code** - Only coordinates

## Troubleshooting

If context is lost:
1. Check \`orchestrator/orchestrator-state.yaml\`
2. Look for TODO files in \`todos/\`
3. Read \`.claude/CLAUDE.md\` sections 7-9

## Support

Refer to:
- Main documentation: \`README.md\`
- State machine: \`core/SOFTWARE-FACTORY-STATE-MACHINE.md\`
- Operations guide: \`core/ORCHESTRATOR-MASTER-OPERATIONS-GUIDE.md\`

Happy coding with your new software factory! 🏭
EOF

# Final summary
echo ""
echo "========================================="
echo "✅ Setup Complete!"
echo "========================================="
echo ""
echo "Project: $PROJECT_NAME"
echo "Location: $TARGET_DIR"
echo "Language: $PRIMARY_LANG"
echo "Complexity: Level $COMPLEXITY"
echo ""
echo "Next steps:"
echo "1. cd $TARGET_DIR"
echo "2. Review QUICK-START.md"
echo "3. Edit orchestrator/PROJECT-IMPLEMENTATION-PLAN.md"
echo "4. Run: /continue-orchestrating in Claude Code"
echo ""
echo "Happy orchestrating! 🎼"