---
name: init-software-factory-noninteractive
description: Initialize Software Factory 3.0 project non-interactively (no questions, all parameters provided)
---

# /init-software-factory-noninteractive

╔═══════════════════════════════════════════════════════════════════════════════╗
║                  SOFTWARE FACTORY 3.0 - NON-INTERACTIVE                      ║
║                      PROJECT INITIALIZATION COMMAND                           ║
║                                                                               ║
║ Purpose: Fully automated initialization for CI/CD, experiment harnesses,     ║
║          and batch project creation. NO interactive questions.               ║
║ Flow: Same as /init-software-factory but architect skips Q&A state          ║
╚═══════════════════════════════════════════════════════════════════════════════╝

## 📖 USAGE

### Command Signature

```bash
/init-software-factory-noninteractive \
  --project-name <name> \
  --project-description "<description>" \
  --language <language> \
  --project-type <type> \
  --build-system <system> \
  --test-framework <framework> \
  --codebase-type <new|existing> \
  [--repo-url <url>] \
  [--frameworks <comma-separated-list>] \
  [--architecture <pattern>] \
  [--deployment <environment>] \
  [--coverage-target <percentage>] \
  [--additional-context <text>]
```

### Example: Minimal Required Parameters

```bash
/init-software-factory-noninteractive \
  --project-name "simple-api" \
  --project-description "REST API for task management with CRUD operations" \
  --language "python" \
  --project-type "api-service" \
  --build-system "poetry" \
  --test-framework "pytest" \
  --codebase-type "new"

# Defaults applied automatically:
# - frameworks: ["fastapi"] (auto-detected for Python API)
# - architecture: "three-tier" (auto-detected for API service)
# - deployment: "cloud-agnostic"
# - containers: "docker"
# - coverage-target: 70
```

### Example: Full Parameters for Existing Codebase

```bash
/init-software-factory-noninteractive \
  --project-name "idpbuilder-push" \
  --project-description "$(cat <<'EOF'
Implement idpbuilder push command to upload OCI images to Gitea registry at
https://gitea.cnoe.localtest.me:8443/ with authentication via -username and
-password flags, disable default certificate checks.

Context: idpbuilder is a Go CLI tool. It currently can build images but cannot
push them. We need to integrate go-containerregistry library for OCI operations.
EOF
)" \
  --language "go" \
  --project-type "cli" \
  --build-system "make" \
  --test-framework "go-test" \
  --codebase-type "existing" \
  --repo-url "https://github.com/cnoe-io/idpbuilder" \
  --frameworks "cobra,go-containerregistry" \
  --architecture "cli-extension" \
  --coverage-target "70"
```

## 📋 PARAMETER SPECIFICATION

### Required Parameters (Cannot Proceed Without)

| Parameter | Type | Description | Example |
|-----------|------|-------------|---------|
| `--project-name` | string | Alphanumeric project name/prefix | `idpbuilder-push` |
| `--project-description` | string | Full PRD or project description | `"Implement push command..."` |
| `--language` | string | Primary programming language | `go`, `python`, `typescript`, `rust` |
| `--project-type` | string | Type of project | `cli`, `library`, `api-service`, `web-app` |
| `--build-system` | string | Build tool | `make`, `npm`, `cargo`, `gradle`, `go` |
| `--test-framework` | string | Testing framework | `pytest`, `jest`, `go-test`, `junit` |
| `--codebase-type` | enum | New or existing codebase | `new`, `existing` |

### Optional Parameters with Auto-Detection

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `--repo-url` | string | (none) | Repository URL if existing codebase |
| `--fork-url` | string | (none) | Fork URL if applicable |
| `--main-branch` | string | `main` | Main branch name |
| `--frameworks` | list | (language defaults) | Key libraries/frameworks (comma-separated) |
| `--architecture` | string | (auto-detect from type) | Architecture pattern |
| `--deployment` | string | `cloud-agnostic` | Target deployment environment |
| `--containers` | string | `docker` | Container technology |
| `--ci-cd` | string | `github-actions` | CI/CD system |
| `--coverage-target` | number | `70` | Code coverage percentage |
| `--design-patterns` | list | (language idioms) | Design patterns to follow |
| `--performance` | string | `standard` | Performance requirements |
| `--security` | string | `basic-auth` | Security requirements |
| `--compliance` | string | `none` | Compliance requirements |
| `--additional-context` | string | (none) | Extra context for architect |

## 🤖 AUTO-DETECTION LOGIC

### Language-Specific Defaults

When optional parameters not provided, these defaults apply:

#### Go Projects
```yaml
frameworks: ["standard library"]
build_system: "make"
test_framework: "go test"
architecture: "cli" if cli, "microservices" if api-service
containers: "docker"
deployment: "k8s" if api-service, "binary" if cli
```

#### Python Projects
```yaml
frameworks: ["fastapi" if api-service, "click" if cli, "pytest" always]
build_system: "poetry"
test_framework: "pytest"
architecture: "three-tier" if web-app, "cli" if cli
containers: "docker"
deployment: "cloud-agnostic"
```

#### TypeScript Projects
```yaml
frameworks: ["next.js" if web-app, "express" if api-service]
build_system: "npm"
test_framework: "jest"
architecture: "next-app-router" if web-app, "rest-api" if api-service
containers: "docker"
deployment: "vercel" if web-app, "cloud-agnostic" if api-service
```

#### Rust Projects
```yaml
frameworks: ["tokio" if async, "clap" if cli]
build_system: "cargo"
test_framework: "cargo test"
architecture: "cli" if cli, "async-runtime" if service
containers: "docker-scratch"
deployment: "binary"
```

### Architecture Pattern Detection

Based on `--project-type`:
- `cli` → `cli-tool` or `cli-extension`
- `library` → `library` or `plugin-architecture`
- `api-service` → `rest-api`, `graphql-api`, or `microservices`
- `web-app` → `spa`, `ssr`, or `jamstack`

## 🔄 NON-INTERACTIVE FLOW

### Key Difference from Interactive Mode

**Interactive** (`/init-software-factory`):
```
Orchestrator → SPAWN_ARCHITECT_MASTER_PLANNING
  → Architect: INIT_REQUIREMENTS_GATHERING (asks ~20 questions)
  → Architect: INIT_DECOMPOSE_PRD
  → Architect: INIT_SYNTHESIZE_PLAN
```

**Non-Interactive** (`/init-software-factory-noninteractive`):
```
Orchestrator creates requirements bundle → SPAWN_ARCHITECT_MASTER_PLANNING
  → Architect detects non_interactive_mode: true
  → Architect: SKIPS INIT_REQUIREMENTS_GATHERING (NO QUESTIONS!)
  → Architect: INIT_DECOMPOSE_PRD (reads from bundle)
  → Architect: INIT_SYNTHESIZE_PLAN
```

### Implementation Steps

#### Step 1: Parse Command Parameters

```bash
# Parse all parameters
PROJECT_NAME=""
PROJECT_DESCRIPTION=""
LANGUAGE=""
PROJECT_TYPE=""
BUILD_SYSTEM=""
TEST_FRAMEWORK=""
CODEBASE_TYPE=""
# ... parse all --param values

# Validate required parameters
if [ -z "$PROJECT_NAME" ] || [ -z "$PROJECT_DESCRIPTION" ] || \
   [ -z "$LANGUAGE" ] || [ -z "$PROJECT_TYPE" ] || \
   [ -z "$BUILD_SYSTEM" ] || [ -z "$TEST_FRAMEWORK" ] || \
   [ -z "$CODEBASE_TYPE" ]; then
    echo "❌ ERROR: Missing required parameters"
    echo "Required: --project-name, --project-description, --language, --project-type, --build-system, --test-framework, --codebase-type"
    exit 1
fi
```

#### Step 2: Apply Auto-Detection for Optional Parameters

```bash
# Apply language-specific defaults
case "$LANGUAGE" in
    go)
        FRAMEWORKS="${FRAMEWORKS:-standard library}"
        ARCHITECTURE="${ARCHITECTURE:-cli}"
        [ "$PROJECT_TYPE" == "api-service" ] && ARCHITECTURE="microservices"
        ;;
    python)
        if [ "$PROJECT_TYPE" == "api-service" ]; then
            FRAMEWORKS="${FRAMEWORKS:-fastapi}"
            ARCHITECTURE="${ARCHITECTURE:-three-tier}"
        elif [ "$PROJECT_TYPE" == "cli" ]; then
            FRAMEWORKS="${FRAMEWORKS:-click}"
            ARCHITECTURE="${ARCHITECTURE:-cli-tool}"
        fi
        ;;
    typescript)
        if [ "$PROJECT_TYPE" == "web-app" ]; then
            FRAMEWORKS="${FRAMEWORKS:-next.js}"
            ARCHITECTURE="${ARCHITECTURE:-next-app-router}"
        elif [ "$PROJECT_TYPE" == "api-service" ]; then
            FRAMEWORKS="${FRAMEWORKS:-express}"
            ARCHITECTURE="${ARCHITECTURE:-rest-api}"
        fi
        ;;
    rust)
        FRAMEWORKS="${FRAMEWORKS:-clap}"
        ARCHITECTURE="${ARCHITECTURE:-cli-tool}"
        [ "$PROJECT_TYPE" == "api-service" ] && FRAMEWORKS="tokio,axum"
        ;;
esac

# Apply universal defaults
DEPLOYMENT="${DEPLOYMENT:-cloud-agnostic}"
CONTAINERS="${CONTAINERS:-docker}"
CI_CD="${CI_CD:-github-actions}"
COVERAGE_TARGET="${COVERAGE_TARGET:-70}"
MAIN_BRANCH="${MAIN_BRANCH:-main}"
PERFORMANCE="${PERFORMANCE:-standard}"
SECURITY="${SECURITY:-basic-auth}"
COMPLIANCE="${COMPLIANCE:-none}"
```

#### Step 3: Create Requirements Bundle

```bash
# Create init-state-temp.json with ALL requirements pre-populated
cat > init-state-temp.json <<EOF
{
  "project_name": "$PROJECT_NAME",
  "initial_description": "$PROJECT_DESCRIPTION",
  "non_interactive_mode": true,
  "requirements": {
    "codebase": {
      "type": "$CODEBASE_TYPE",
      "upstream_url": "${REPO_URL:-}",
      "fork_url": "${FORK_URL:-}",
      "main_branch": "$MAIN_BRANCH"
    },
    "technology": {
      "languages": ["$LANGUAGE"],
      "frameworks": [$(echo "$FRAMEWORKS" | sed 's/,/","/g' | sed 's/^/"/' | sed 's/$/"/') ],
      "build_system": "$BUILD_SYSTEM",
      "test_framework": "$TEST_FRAMEWORK",
      "code_generation_tools": []
    },
    "architecture": {
      "type": "$PROJECT_TYPE",
      "pattern": "$ARCHITECTURE",
      "design_patterns": [$(echo "${DESIGN_PATTERNS:-}" | sed 's/,/","/g' | sed 's/^/"/' | sed 's/$/"/')],
      "existing_patterns": []
    },
    "deployment": {
      "environment": "$DEPLOYMENT",
      "containerization": "$CONTAINERS",
      "ci_cd": "$CI_CD",
      "infrastructure_as_code": []
    },
    "quality": {
      "performance": "$PERFORMANCE",
      "security": "$SECURITY",
      "compliance": "$COMPLIANCE",
      "coverage_target": $COVERAGE_TARGET
    },
    "project_specifics": {
      "end_users": "developers",
      "problem": "$PROJECT_DESCRIPTION",
      "competitors": [],
      "documentation": "README + inline comments"
    }
  },
  "additional_context": "$ADDITIONAL_CONTEXT",
  "timestamp": "$(date -Iseconds)"
}
EOF

echo "✅ Requirements bundle created in init-state-temp.json"
echo "✅ Non-interactive mode: Architect will NOT ask questions"
```

#### Step 4: Check for Existing PRD (PRD Skip Logic)

```bash
# Check if PRD already exists for this project
PRD_FILE="prd/${PROJECT_NAME}-prd.md"

if [ -f "$PRD_FILE" ]; then
    echo "✅ Existing PRD found: $PRD_FILE"
    echo "✅ Skipping Product Manager PRD generation states"
    echo "✅ Will proceed directly to SPAWN_ARCHITECT_MASTER_PLANNING"

    # Set flag in requirements bundle to skip PM states
    jq '.prd_pre_exists = true' init-state-temp.json > tmp.json && mv tmp.json init-state-temp.json

    PRD_SKIP=true
else
    echo "📋 No existing PRD found"
    echo "📋 Product Manager will generate PRD from project description"
    echo "📋 Flow: SPAWN_PRODUCT_MANAGER_PRD_CREATION → ..."

    # Ensure flag is false
    jq '.prd_pre_exists = false' init-state-temp.json > tmp.json && mv tmp.json init-state-temp.json

    PRD_SKIP=false
fi
```

**PRD Skip Decision Logic:**
```
IF prd/${PROJECT_NAME}-prd.md exists:
  → Set prd_pre_exists = true in init-state-temp.json
  → Orchestrator will skip SPAWN_PRODUCT_MANAGER_PRD_CREATION
  → Orchestrator transitions directly to SPAWN_ARCHITECT_MASTER_PLANNING
ELSE:
  → Set prd_pre_exists = false
  → Orchestrator follows normal PRD generation flow
  → SPAWN_PRODUCT_MANAGER_PRD_CREATION → WAITING_FOR_PRD_CREATION → ...
```

#### Step 5: Initialize Orchestrator State

```bash
# Create orchestrator-state-v3.json with initial state
cat > orchestrator-state-v3.json <<EOF
{
  "current_state": "INIT",
  "project_name": "$PROJECT_NAME",
  "project_description": "$PROJECT_DESCRIPTION",
  "project_type": "$PROJECT_TYPE",
  "non_interactive_mode": true,
  "prd_pre_exists": $PRD_SKIP,
  "timestamp": "$(date -Iseconds)",
  "state_history": []
}
EOF

echo "✅ Orchestrator state initialized"
```

#### Step 6: Follow Same SF 3.0 Mandatory Sequence

From this point, the flow depends on PRD existence:

**Path 1: PRD Already Exists (PRD_SKIP=true)**
- INIT
- SPAWN_ARCHITECT_MASTER_PLANNING (architect detects non_interactive_mode and skips Q&A)
- WAITING_FOR_MASTER_ARCHITECTURE
- SPAWN_CODE_REVIEWER_PROJECT_TEST_PLANNING
- ... (rest of SF 3.0 mandatory sequence)

**Path 2: No PRD - Complete Description (PRD_SKIP=false, PM returns CONTINUE=TRUE)**
- INIT
- SPAWN_PRODUCT_MANAGER_PRD_CREATION
- WAITING_FOR_PRD_CREATION (PM analyzes, generates complete PRD, sets CONTINUE=TRUE)
- SPAWN_ARCHITECT_MASTER_PLANNING (architect detects non_interactive_mode and skips Q&A)
- WAITING_FOR_MASTER_ARCHITECTURE
- ... (rest of SF 3.0 mandatory sequence)

**Path 3: No PRD - Incomplete Description (PRD_SKIP=false, PM returns CONTINUE=FALSE)**
- INIT
- SPAWN_PRODUCT_MANAGER_PRD_CREATION
- WAITING_FOR_PRD_CREATION (PM analyzes, generates partial PRD with gaps, sets CONTINUE=FALSE)
- WAITING_FOR_PRD_VALIDATION (human edits PRD to fill [NEEDS INPUT] markers)
- [HUMAN: Run /continue-orchestrating]
- SPAWN_PRODUCT_MANAGER_PRD_VALIDATION
- (Validation passes) SPAWN_ARCHITECT_MASTER_PLANNING
- WAITING_FOR_MASTER_ARCHITECTURE
- ... (rest of SF 3.0 mandatory sequence)

## 🔴🔴🔴 ARCHITECT NON-INTERACTIVE DETECTION 🔴🔴🔴

### How Architect Detects Non-Interactive Mode

In `agent-states/initialization/architect/INIT_REQUIREMENTS_GATHERING/rules.md`, add:

```bash
# STEP 1: Check for non-interactive mode
if [ -f "$CLAUDE_PROJECT_DIR/init-state-temp.json" ]; then
    NON_INTERACTIVE=$(jq -r '.non_interactive_mode // false' init-state-temp.json)

    if [ "$NON_INTERACTIVE" == "true" ]; then
        echo "✅ Non-interactive mode detected"
        echo "✅ Requirements bundle found - skipping Q&A"
        echo "✅ Transitioning directly to INIT_DECOMPOSE_PRD"

        # Load requirements from bundle instead of asking
        # Transition immediately to INIT_DECOMPOSE_PRD
        exit 0
    fi
fi

# STEP 2: If not non-interactive, proceed with normal Q&A
echo "Interactive mode - gathering requirements..."
# Ask ~20 questions...
```

### Architect INIT_DECOMPOSE_PRD for Non-Interactive

```bash
# Check for requirements bundle
if [ -f "$CLAUDE_PROJECT_DIR/init-state-temp.json" ]; then
    # Load ALL requirements from bundle
    LANGUAGE=$(jq -r '.requirements.technology.languages[0]' init-state-temp.json)
    FRAMEWORKS=$(jq -r '.requirements.technology.frameworks | join(",")' init-state-temp.json)
    BUILD_SYSTEM=$(jq -r '.requirements.technology.build_system' init-state-temp.json)
    PROJECT_TYPE=$(jq -r '.requirements.architecture.type' init-state-temp.json)
    # ... load all fields

    echo "✅ Requirements loaded from bundle (non-interactive mode)"
else
    # Load from Q&A results (interactive mode)
    echo "✅ Requirements loaded from Q&A session"
fi

# Proceed with decomposition using loaded requirements
```

## 🚦 STATE TRANSITIONS

**IDENTICAL to /init-software-factory** except:
- Architect skips INIT_REQUIREMENTS_GATHERING state
- All requirements pre-loaded from init-state-temp.json

Same SF 3.0 mandatory sequence applies with all State Manager bookends.

## 🔴 CRITICAL RULES

**All SF 3.0 rules apply exactly as in /init-software-factory:**
- State Manager bookend pattern at every transition
- R287 TODO persistence
- R288 state file updates
- R322 mandatory stops
- R405 continuation flags
- R341 TDD requirements
- R342 early integration branches
- R506 pre-commit enforcement

## 📊 SUCCESS CRITERIA

Initialization succeeds when:
- ✅ All parameters parsed and validated
- ✅ Auto-detection applied for optional parameters
- ✅ Requirements bundle created in init-state-temp.json
- ✅ Architect NEVER asked interactive questions
- ✅ All SF 3.0 planning artifacts created
- ✅ orchestrator-state-v3.json shows current_state: "WAVE_START"
- ✅ Ready for /continue-software-factory

## 🎯 USE CASES

### 1. Experiment Harness

```bash
for config in config_a config_b config_c; do
  instance_dir="/tmp/sf3-experiment-${config}"

  # Setup instance
  cp -r software-factory-template "$instance_dir"
  cd "$instance_dir"

  # Initialize non-interactively
  /init-software-factory-noninteractive \
    --project-name "test-project-${config}" \
    --project-description "$(cat test-prd.txt)" \
    --language "go" \
    --project-type "cli" \
    --build-system "make" \
    --test-framework "go-test" \
    --codebase-type "new"

  # Continue to development
  /continue-software-factory

  # Measure results
  measure_metrics "$instance_dir"
done
```

### 2. CI/CD Pipeline

```yaml
# .github/workflows/test-sf3.yml
jobs:
  test-initialization:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Initialize SF3 Project
        run: |
          /init-software-factory-noninteractive \
            --project-name "ci-test-project" \
            --project-description "${{ github.event.pull_request.body }}" \
            --language "python" \
            --project-type "api-service" \
            --build-system "poetry" \
            --test-framework "pytest" \
            --codebase-type "new"
```

### 3. Batch Project Creation

```bash
# Create multiple projects from specification file
while IFS=, read -r name desc lang type; do
  /init-software-factory-noninteractive \
    --project-name "$name" \
    --project-description "$desc" \
    --language "$lang" \
    --project-type "$type" \
    --build-system "auto" \
    --test-framework "auto" \
    --codebase-type "new"
done < projects.csv
```

## ⚠️ LIMITATIONS

1. **No Custom Q&A**: Cannot ask clarifying questions - all info must be in parameters
2. **Auto-Detection May Be Wrong**: Defaults may not match your specific needs
3. **Less Flexible**: Cannot adapt to nuanced requirements like interactive mode
4. **Requires Complete PRD**: Project description must be comprehensive

## 💡 BEST PRACTICES

1. **Detailed Descriptions**: Provide comprehensive project descriptions with all context
2. **Explicit Parameters**: Don't rely on auto-detection for critical choices
3. **Test First**: Run with minimal projects before production use
4. **Version Control**: Track init-state-temp.json for reproducibility
5. **Error Handling**: Have fallback to interactive mode if automation fails

---

**For interactive initialization with Q&A, use `/init-software-factory` instead.**
