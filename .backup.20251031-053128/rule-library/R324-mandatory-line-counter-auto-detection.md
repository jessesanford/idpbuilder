# 🔴🔴🔴 RULE R324 - MANDATORY LINE COUNTER AUTO-DETECTION 🔴🔴🔴

**Criticality**: BLOCKING - Failure to use auto-detection = -100% IMMEDIATE FAILURE
**Affected Agents**: code-reviewer, sw-engineer
**Created**: 2025-09-06
**Updated**: 2025-09-06
**Purpose**: Enforce ONLY auto-detection mode for line-counter.sh tool

## 🚨🚨🚨 CRITICAL UPDATE: NO MANUAL BASE SPECIFICATION ALLOWED! 🚨🚨🚨

### THE TOOL IS NOW SMART - IT AUTO-DETECTS EVERYTHING!

**MANDATORY REQUIREMENTS:**
1. **ALWAYS use ${PROJECT_ROOT}/tools/line-counter.sh** - NO EXCEPTIONS
2. **NO PARAMETERS for base detection** - Tool auto-detects EVERYTHING!
3. **NEVER use -b parameter** - That's DEPRECATED/WRONG syntax!
4. **NEVER do manual counting** - AUTOMATIC FAILURE (-100%)
5. **Tool shows detected base** - Look for "🎯 Detected base:" in output

## HOW THE AUTO-DETECTION WORKS

The line-counter.sh tool automatically determines the correct base branch by:

1. **Pattern Detection**: Analyzes branch name for phase/wave/effort patterns
2. **Split Detection**: Recognizes split-XXX patterns and chains them correctly
3. **Integration Detection**: Identifies integration branches
4. **Project Prefix Support**: Handles optional project prefixes automatically

### For Regular Efforts:
- Branch: `phase1/wave1/my-effort`
- Auto-detected base: `main` (if phase1/wave1)
- Auto-detected base: `phase1/integration` (if later wave)

### For Splits:
- Branch: `phase1/wave1/my-effort--split-001`
- Auto-detected base: `phase1/wave1/my-effort` (original)
- Branch: `phase1/wave1/my-effort--split-002`
- Auto-detected base: `phase1/wave1/my-effort--split-001` (previous split)

### For Integration Branches:
- Branch: `phase1/wave1/integration`
- Auto-detected base: `phase1/integration` or `main`

## CORRECT USAGE (MANDATORY PATTERN)

```bash
# STEP 1: Navigate to effort directory
cd /path/to/effort/directory
pwd  # Confirm location

# STEP 2: Ensure code is committed
git status  # Must be clean
git add -A && git commit -m "feat: ready for measurement" && git push

# STEP 3: Find project root
PROJECT_ROOT=$(pwd)
while [ "$PROJECT_ROOT" != "/" ]; do 
    [ -f "$PROJECT_ROOT/orchestrator-state-v3.json" ] && break
    PROJECT_ROOT=$(dirname "$PROJECT_ROOT")
done
echo "Project root: $PROJECT_ROOT"

# STEP 4: RUN THE TOOL - NO BASE PARAMETERS!
$PROJECT_ROOT/tools/line-counter.sh

# Output will show:
# 🎯 Detected base: [automatically determined base]
# 📦 Analyzing branch: [current branch]
# ✅ Total non-generated lines: [count]
```

## AUTOMATIC FAILURES (-100% GRADE)

### ❌❌❌ THESE WILL FAIL YOU IMMEDIATELY:

```bash
# WRONG - Manual base specification
$PROJECT_ROOT/tools/line-counter.sh -b main  # -100% FAILURE!

# WRONG - Using git diff with manual base
git diff main --stat  # -100% FAILURE!

# WRONG - Manual counting
wc -l *.go  # -100% FAILURE!
find . -name "*.go" | xargs wc -l  # -100% FAILURE!

# WRONG - Old tool locations
/workspaces/kcp-shared-tools/line-counter.sh  # -100% FAILURE!
```

## TOOL OUTPUT INTERPRETATION

The tool provides clear output showing its auto-detection:

```
🎯 Detected base: phase1/integration
   ↳ Source: Pattern matching for phase1/wave2/api-client
📦 Analyzing branch: phase1/wave2/api-client
📊 Generating diff against base branch...
✅ Total non-generated lines: 487
```

### Key Output Elements:
- **🎯 Detected base**: Shows what base the tool determined
- **Source**: Explains how it made the decision
- **📦 Analyzing branch**: Confirms what's being measured
- **✅ Total non-generated lines**: The ONLY number that matters

## ENFORCEMENT IN CODE REVIEWER

The code reviewer MUST:
1. Always use line-counter.sh for measurements
2. Never specify base branches manually
3. Document the exact command used
4. Include the full tool output showing auto-detection
5. Report the detected base in review documentation

## ENFORCEMENT IN SW ENGINEER

The SW engineer MUST:
1. Measure regularly during development
2. Use ONLY line-counter.sh 
3. Stop immediately if approaching 700 lines
4. Never attempt manual counting

## GRADING PENALTIES

- **-100%**: Using -b parameter with line-counter.sh
- **-100%**: Manual counting (wc, find, etc.)
- **-100%**: Using git diff with wrong base
- **-100%**: Not using line-counter.sh at all
- **-50%**: Not documenting auto-detected base
- **-30%**: Misinterpreting tool output

## RATIONALE

Manual base specification leads to:
- Incorrect measurements (measuring against wrong base)
- Inflated line counts (counting all code since main)
- Failed size compliance (11,876 lines when should be ~500)
- Wasted effort on unnecessary splits

Auto-detection ensures:
- Correct base every time
- Accurate measurements
- Proper size compliance
- Efficient development

## VERIFICATION CHECKLIST

Before measuring:
- [ ] In correct effort directory
- [ ] Code committed and pushed
- [ ] Found project root with orchestrator-state-v3.json
- [ ] Using line-counter.sh from project root
- [ ] NO -b parameter specified
- [ ] Ready to document auto-detected base

## REFERENCES

- R198: Line counter usage (general)
- R200: Measure only changeset
- R304: Mandatory line counter enforcement
- R319: Orchestrator never measures (doesn't apply to code-reviewer!)

---

**REMEMBER**: The tool is SMART. Let it do its job. Don't try to outsmart it!