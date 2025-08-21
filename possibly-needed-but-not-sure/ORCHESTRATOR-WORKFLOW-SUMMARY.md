# Orchestrator Workflow Summary

## Key Concepts

### Working Copies (NOT Worktrees)
Each effort gets an **independent sparse-checkout working copy**:
- Location: `/workspaces/efforts/phase{X}/wave{Y}/effort{Z}-{name}/`
- Method: Sparse clone with only needed directories (pkg, apis, cmd, test, hack)
- Purpose: Isolation between efforts, minimal disk usage

### Branch Strategy
```
Standard Effort:  /phase{X}/wave{Y}/effort{Z}-{name}
Split Original:   /phase{X}/wave{Y}/effort{Z}-{name}-to-be-split
Split Parts:      /phase{X}/wave{Y}/effort{Z}-{name}-part{N}
Integration:      phase{X}-integration
```

### State Management
- **File**: `/workspaces/agent-configs/tmc-orchestrator-impl-8-20-2025/orchestrator-state.yaml`
- **Version Controlled**: Committed to repo after every change
- **Purpose**: Recovery from interruption, progress tracking

## Workflow Steps

### 1. Create Working Copy
```bash
# For effort E1.1.1 (api-types-core)
mkdir -p /workspaces/efforts/phase1/wave1/effort1-api-types-core
cd /workspaces/efforts/phase1/wave1/effort1-api-types-core
git clone --no-checkout https://github.com/jessesanford/kcp.git .
git sparse-checkout init --cone
git sparse-checkout set pkg apis cmd test hack
git checkout main
git checkout -b /phase1/wave1/effort1-api-types-core
git push -u origin /phase1/wave1/effort1-api-types-core
```

### 2. Task Agent
Provide paths, NOT content:
```markdown
Working directory: /workspaces/efforts/phase1/wave1/effort1-api-types-core
Branch: /phase1/wave1/effort1-api-types-core
Instructions: READ /workspaces/agent-configs/tmc-orchestrator-impl-8-20-2025/PHASE1-SPECIFIC-IMPL-PLAN-8-20-25.md#E1.1.1
```

### 3. Review Code
```markdown
Branch to review: /phase1/wave1/effort1-api-types-core
Guide: READ /workspaces/agent-configs/tmc-orchestrator-impl-8-20-2025/KCP-CODE-REVIEWER-COMPREHENSIVE-GUIDE.md
Measure with: /workspaces/kcp-shared-tools/tmc-pr-line-counter.sh
```

### 4. Handle Splits (if needed)
```bash
# Rename original
cd /workspaces/efforts/phase1/wave1/effort1-api-types-core
git branch -m /phase1/wave1/effort1-api-types-core-to-be-split

# Create split working copies
for N in 1 2 3; do
    mkdir -p /workspaces/efforts/phase1/wave1/effort1-api-types-core-split${N}
    cd /workspaces/efforts/phase1/wave1/effort1-api-types-core-split${N}
    git clone --no-checkout https://github.com/jessesanford/kcp.git .
    git sparse-checkout init --cone
    git sparse-checkout set pkg apis cmd test hack
    git checkout main
    git checkout -b /phase1/wave1/effort1-api-types-core-part${N}
done
```

### 5. Update State
```bash
cd /workspaces/agent-configs
# Edit tmc-orchestrator-impl-8-20-2025/orchestrator-state.yaml
git add tmc-orchestrator-impl-8-20-2025/orchestrator-state.yaml
git commit -m "state: completed E1.1.1"
git push
```

## Directory Structure

```
/workspaces/
├── agent-configs/
│   └── tmc-orchestrator-impl-8-20-2025/
│       ├── orchestrator-state.yaml          # Version controlled state
│       ├── ORCHESTRATOR-*.md                # Instructions
│       └── phase{X}/wave{Y}/effort{Z}/
│           └── code-review/                 # Review outputs
└── efforts/                                 # Working copies
    ├── phase1/
    │   ├── wave1/
    │   │   ├── effort1-api-types-core/     # Sparse checkout
    │   │   ├── effort2-synctarget-types/   # Sparse checkout
    │   │   └── effort1-api-types-core-split1/  # If split needed
    │   └── wave2/
    └── phase2/
```

## Critical Rules

1. **NO WORKTREES** - Use sparse checkouts
2. **State in repo** - `/workspaces/agent-configs/tmc-orchestrator-impl-8-20-2025/orchestrator-state.yaml`
3. **Unique branches** - Each effort gets its own branch
4. **Isolated working copies** - Each effort/split gets its own directory
5. **Provide paths, not content** - Keep orchestrator context minimal

## Common Commands

### Check State
```bash
cat /workspaces/agent-configs/tmc-orchestrator-impl-8-20-2025/orchestrator-state.yaml
```

### List Working Copies
```bash
ls -la /workspaces/efforts/phase*/wave*/effort*/
```

### Check Branch Status
```bash
cd /workspaces/efforts/phase1/wave1/effort1-api-types-core
git branch -vv
git status
```

### Measure Size
```bash
/workspaces/kcp-shared-tools/tmc-pr-line-counter.sh -c /phase1/wave1/effort1-api-types-core
```

## Recovery

If orchestrator is interrupted:
1. Read state from `/workspaces/agent-configs/tmc-orchestrator-impl-8-20-2025/orchestrator-state.yaml`
2. Check working copies in `/workspaces/efforts/`
3. Resume from last completed effort
4. Continue orchestration

This workflow ensures clean separation, version control, and recovery capability.