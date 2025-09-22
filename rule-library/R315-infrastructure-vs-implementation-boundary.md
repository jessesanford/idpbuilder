# 🚨🚨🚨 RULE R315 - Infrastructure vs Implementation Boundary

**Criticality:** BLOCKING  
**Grading Impact:** -50% for violations  
**Enforcement:** Strict separation of concerns
**Applies To:** Orchestrator during infrastructure creation

## Rule Statement

The Orchestrator may ONLY create empty infrastructure. ANY content, code, or implementation MUST be delegated to appropriate agents. Infrastructure is LIMITED to empty directories and configuration files.

## 🔴🔴🔴 CRITICAL DISTINCTIONS 🔴🔴🔴

### ✅ INFRASTRUCTURE (Orchestrator Allowed)
- Creating empty directories
- Creating orchestrator-state.json
- Creating markdown planning documents
- Creating empty .gitkeep files
- Writing TODO files
- Creating YAML/JSON config (non-code)

### ❌ IMPLEMENTATION (Orchestrator FORBIDDEN)
- **Copying ANY files with code**
- **Moving existing code files**
- **Creating template code files**
- **Writing starter/boilerplate code**
- **Creating test file structures**
- **Populating directories with code**
- **Creating symbolic links to code**
- **Generating code from templates**

## Violation Examples

### ❌ FORBIDDEN: Copying Code as "Infrastructure"
```bash
# VIOLATION - Even for split setup
cp *.go efforts/phase1/wave1/split2/
cp *.py efforts/phase1/wave1/split2/

# WHY: Code operations are ALWAYS implementation
```

### ❌ FORBIDDEN: Creating Starter Files
```bash
# VIOLATION - No boilerplate allowed
echo "package main" > efforts/phase1/wave1/main.go
touch efforts/phase1/wave1/test.py

# WHY: Even empty code files are implementation
```

### ❌ FORBIDDEN: Moving Code Between Efforts
```bash
# VIOLATION - File operations on code
mv split1/*.js split2/
ln -s ../shared/utils.go efforts/phase1/utils.go

# WHY: All code manipulation is implementation
```

## Correct Patterns

### ✅ CORRECT: Create Empty Infrastructure
```bash
# Orchestrator creates structure
mkdir -p efforts/phase1/wave1/split2
touch efforts/phase1/wave1/split2/.gitkeep

# Then spawn agent for content
echo "🚀 Spawning SW Engineer to populate split2"
```

### ✅ CORRECT: Delegate All Code Operations
```bash
# Orchestrator identifies need
echo "📋 Split2 needs test files from split1"

# Spawn agent with instructions
Task: Copy relevant test files to split2
Agent: sw-engineer
Instructions:
- Copy api_test.go to split2
- Copy models_test.go to split2
- Verify tests still pass
```

## 🔴 COMPREHENSIVE COMMAND WHITELIST/BLACKLIST 🔴

### ✅ ORCHESTRATOR ALLOWED COMMANDS (Infrastructure Only)
```bash
# Directory Operations
mkdir -p                    # Create empty directories
rmdir                      # Remove empty directories  
ls, find, tree             # Read-only exploration

# Navigation
cd, pwd                    # Directory navigation
pushd, popd               # Directory stack management

# Git Operations (Infrastructure)
git clone --sparse         # Initial repository setup
git checkout -b           # Create new branches
git branch                # List/manage branches
git status, git log       # Read-only git info

# File Creation (Non-Code Only)
touch .gitkeep            # Create marker files
echo "text" > file.md     # Create markdown docs
cat > file.yaml           # Create YAML configs
yq, jq                    # Process YAML/JSON

# State Management
jq orchestrator-state.json    # Update state file
git add orchestrator-state.json    # Stage state changes
git commit -m "orchestrator: ..." # Commit orchestration files
```

### ❌ ORCHESTRATOR FORBIDDEN COMMANDS (Implementation)
```bash
# Code File Operations - ABSOLUTELY FORBIDDEN
cp *.go *.py *.js *.ts *.java     # NO copying code files
mv *.cpp *.c *.rs *.rb *.php      # NO moving code files
ln -s any_code_file               # NO linking code files
rsync --include="*.go"            # NO syncing code
tar czf archive.tar.gz *.py       # NO archiving code
zip -r code.zip src/              # NO zipping code

# Code Modification - ABSOLUTELY FORBIDDEN  
sed -i 's/old/new/' *.go         # NO editing code
awk '{print}' > file.py          # NO generating code
perl -pi -e 's/x/y/' *.js       # NO modifying code
> main.go                        # NO creating code files
>> test.py                       # NO appending to code

# Package/Build Operations - FORBIDDEN
npm install                      # NO package installation
pip install                      # NO Python packages
go mod init                      # NO module initialization
cargo build                      # NO building
make                            # NO make operations
gradle, maven                    # NO build tools

# Test Operations - FORBIDDEN
go test                         # NO running tests
pytest                          # NO Python tests
npm test                        # NO JavaScript tests
cargo test                      # NO Rust tests
```

## Enforcement Mechanism

```bash
# Enhanced validation with whitelist/blacklist
validate_orchestrator_command() {
    local full_command="$1"
    local base_command=$(echo "$full_command" | awk '{print $1}')
    
    # Check against blacklist
    local blacklist="cp|mv|ln|rsync|tar|zip|sed|awk|perl|npm|pip|go|cargo|make|gradle|maven|pytest"
    
    if [[ "$base_command" =~ ^($blacklist)$ ]]; then
        # Check if operating on code files
        if [[ "$full_command" =~ \.(go|py|js|ts|java|cpp|c|rs|rb|php|jsx|tsx|vue|swift|kt|scala|lua|dart)(\s|$) ]]; then
            echo "🚨🚨🚨 R315 VIOLATION DETECTED!"
            echo "❌ FORBIDDEN: $base_command on code files"
            echo "❌ Command attempted: $full_command"
            echo "✅ REQUIRED: Spawn appropriate agent instead"
            return 1
        fi
    fi
    
    # Check for redirect operations to code files
    if [[ "$full_command" =~ (>|>>).*\.(go|py|js|ts|java|cpp|c|rs|rb|php) ]]; then
        echo "🚨🚨🚨 R315 VIOLATION DETECTED!"
        echo "❌ FORBIDDEN: Creating/modifying code files"
        echo "✅ REQUIRED: Spawn SW Engineer instead"
        return 1
    fi
    
    return 0
}

# Hook into command execution
command_not_found_handle() {
    validate_orchestrator_command "$*" || return 1
    command "$@"
}
```

## Key Principle

**"If it has code, it's not infrastructure"**

The Orchestrator builds the empty stage. The actors (agents) bring the performance (code).

## Grading Impact

- First violation: -50% with opportunity to correct
- Second violation: -75% 
- Third violation: FAILURE

## Self-Check Questions

Before ANY file operation, ask:
1. Am I creating an EMPTY structure? ✅ Proceed
2. Am I moving/copying existing files? ❌ Stop, spawn agent
3. Does this involve ANY code? ❌ Stop, spawn agent
4. Is this purely organizational? ✅ Proceed

---
**REMEMBER:** Infrastructure = empty containers. Implementation = ANY content.