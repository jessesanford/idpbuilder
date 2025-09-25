# 🚨 LINE COUNTER USAGE FOR SPLIT EFFORTS 🚨

## The Problem with Split Efforts
Split efforts use FULL single-branch checkout (R271), creating a NEW git repository in the split directory. This means:
- `git rev-parse --show-toplevel` returns the split directory, NOT the main project
- The line counter is in the main project, NOT in your split directory
- You need to look UP the directory tree, not at git root

## The CORRECT Way for Split Efforts

```bash
# FOR SPLIT EFFORTS - Find main project root by looking UP
find_line_counter_for_splits() {
    echo "🔍 Finding line counter for split effort..."
    
    # Start from current directory and go UP
    SEARCH_DIR=$(pwd)
    
    # Keep going up until we find the main project with tools/line-counter.sh
    while [ "$SEARCH_DIR" != "/" ]; do
        # Check if tools/line-counter.sh exists at this level
        if [ -f "$SEARCH_DIR/tools/line-counter.sh" ]; then
            LINE_COUNTER="$SEARCH_DIR/tools/line-counter.sh"
            echo "✅ Found line counter at: $LINE_COUNTER"
            break
        fi
        
        # Also check one level up for the main project
        PARENT_DIR=$(dirname "$SEARCH_DIR")
        if [ -f "$PARENT_DIR/tools/line-counter.sh" ]; then
            LINE_COUNTER="$PARENT_DIR/tools/line-counter.sh"
            echo "✅ Found line counter at: $LINE_COUNTER"
            break
        fi
        
        # Go up one directory
        SEARCH_DIR=$(dirname "$SEARCH_DIR")
    done
    
    # Verify we found it
    if [ -z "$LINE_COUNTER" ] || [ ! -f "$LINE_COUNTER" ]; then
        echo "❌ FATAL: Could not find line-counter.sh"
        echo "Searched from $(pwd) up to root"
        exit 1
    fi
    
    # Now run it with current branch
    CURRENT_BRANCH=$(git branch --show-current)
    echo "📊 Measuring branch: $CURRENT_BRANCH"
    
    # Run with explicit branch parameter since we're in a separate checkout
    "$LINE_COUNTER" -c "$CURRENT_BRANCH"
}

# Execute
find_line_counter_for_splits
```

## Quick One-Liner for Split Efforts

```bash
# Find and run line counter from split directory
LC=$(d=$(pwd); while [ "$d" != "/" ]; do [ -f "$d/tools/line-counter.sh" ] && echo "$d/tools/line-counter.sh" && break; [ -f "$(dirname "$d")/tools/line-counter.sh" ] && echo "$(dirname "$d")/tools/line-counter.sh" && break; d=$(dirname "$d"); done); [ -n "$LC" ] && "$LC" -c "$(git branch --show-current)" || echo "Line counter not found"
```

## Why Split Efforts Are Different

1. **Sparse Checkout**: Each split is its own git repository
2. **Directory Structure**: 
   ```
   /home/vscode/workspaces/idpbuilder-oci-mgmt/        <- Main project (has tools/)
   └── efforts/
       └── phase2/
           └── wave1/
               └── buildah-integration-split-002/      <- Your split (separate git repo)
                   └── .git/                            <- This makes git root HERE, not main project
   ```
3. **Git Root Confusion**: `git rev-parse --show-toplevel` returns split directory, not main project

## The Solution Pattern

```bash
# STEP 1: Recognize you're in a split
if ls SPLIT-PLAN-*.md 2>/dev/null | head -1 > /dev/null; then
    echo "🔄 Split effort detected - using special line counter search"
    IS_SPLIT=true
else
    IS_SPLIT=false
fi

# STEP 2: Find line counter differently for splits
if [ "$IS_SPLIT" = true ]; then
    # For splits: search UP the directory tree
    CURRENT_DIR=$(pwd)
    while [ "$CURRENT_DIR" != "/" ]; do
        # Check parent directories for main project
        for check_dir in "$CURRENT_DIR" "$(dirname "$CURRENT_DIR")" "$(dirname "$(dirname "$CURRENT_DIR")")"; do
            if [ -f "$check_dir/tools/line-counter.sh" ]; then
                LINE_COUNTER="$check_dir/tools/line-counter.sh"
                echo "✅ Found at: $LINE_COUNTER"
                break 2
            fi
        done
        CURRENT_DIR=$(dirname "$CURRENT_DIR")
    done
else
    # For normal efforts: use standard pattern
    PROJECT_ROOT=$(pwd)
    while [ "$PROJECT_ROOT" != "/" ]; do
        if [ -f "$PROJECT_ROOT/orchestrator-state.json" ]; then
            break
        fi
        PROJECT_ROOT=$(dirname "$PROJECT_ROOT")
    done
    LINE_COUNTER="$PROJECT_ROOT/tools/line-counter.sh"
fi

# STEP 3: Run with current branch (splits need explicit branch)
if [ -f "$LINE_COUNTER" ]; then
    if [ "$IS_SPLIT" = true ]; then
        # Splits need explicit branch parameter
        "$LINE_COUNTER" -c "$(git branch --show-current)"
    else
        # Normal efforts can auto-detect
        "$LINE_COUNTER"
    fi
else
    echo "❌ Line counter not found!"
    exit 1
fi
```

## Common Pitfalls to Avoid

❌ **DON'T assume git root = project root in splits**
```bash
# WRONG for splits
git rev-parse --show-toplevel  # Returns split directory, not main project!
```

❌ **DON'T use relative paths from split directory**
```bash
# WRONG - tools/ doesn't exist in split directory
./tools/line-counter.sh
```

❌ **DON'T hardcode the main project path**
```bash
# WRONG - Path might change
/home/vscode/workspaces/idpbuilder-oci-mgmt/tools/line-counter.sh
```

## The Foolproof Approach

```bash
# This works for BOTH splits and normal efforts
measure_lines() {
    echo "📊 Starting line measurement..."
    
    # Try multiple strategies in order
    LINE_COUNTER=""
    
    # Strategy 1: Look for tools/ going up from current directory
    SEARCH_DIR=$(pwd)
    MAX_LEVELS=10  # Prevent infinite loop
    LEVEL=0
    
    while [ "$SEARCH_DIR" != "/" ] && [ $LEVEL -lt $MAX_LEVELS ]; do
        if [ -f "$SEARCH_DIR/tools/line-counter.sh" ]; then
            LINE_COUNTER="$SEARCH_DIR/tools/line-counter.sh"
            echo "✅ Found via directory search: $LINE_COUNTER"
            break
        fi
        SEARCH_DIR=$(dirname "$SEARCH_DIR")
        LEVEL=$((LEVEL + 1))
    done
    
    # Strategy 2: If not found, look for orchestrator-state.json
    if [ -z "$LINE_COUNTER" ]; then
        SEARCH_DIR=$(pwd)
        while [ "$SEARCH_DIR" != "/" ]; do
            if [ -f "$SEARCH_DIR/orchestrator-state.json" ]; then
                if [ -f "$SEARCH_DIR/tools/line-counter.sh" ]; then
                    LINE_COUNTER="$SEARCH_DIR/tools/line-counter.sh"
                    echo "✅ Found via orchestrator-state.json: $LINE_COUNTER"
                    break
                fi
            fi
            SEARCH_DIR=$(dirname "$SEARCH_DIR")
        done
    fi
    
    # Strategy 3: Check known relative locations
    if [ -z "$LINE_COUNTER" ]; then
        for relative_path in \
            "../../../../tools/line-counter.sh" \
            "../../../tools/line-counter.sh" \
            "../../tools/line-counter.sh" \
            "../tools/line-counter.sh"; do
            if [ -f "$relative_path" ]; then
                LINE_COUNTER=$(realpath "$relative_path")
                echo "✅ Found via relative path: $LINE_COUNTER"
                break
            fi
        done
    fi
    
    # Execute if found
    if [ -n "$LINE_COUNTER" ] && [ -f "$LINE_COUNTER" ]; then
        CURRENT_BRANCH=$(git branch --show-current)
        echo "📈 Measuring branch: $CURRENT_BRANCH"
        "$LINE_COUNTER" -c "$CURRENT_BRANCH"
    else
        echo "❌ ERROR: Could not find line-counter.sh"
        echo "Tried multiple strategies from: $(pwd)"
        exit 1
    fi
}

# Run it
measure_lines
```

## Remember for Splits

1. **You're in a separate git repo** - Don't rely on git root
2. **Search UP the directory tree** - The main project is above you
3. **Use explicit branch parameter** - Auto-detection may fail in separate checkout
4. **Check multiple locations** - Be thorough in your search