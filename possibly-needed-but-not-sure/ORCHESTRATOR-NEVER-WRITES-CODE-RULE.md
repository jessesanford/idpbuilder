# CRITICAL RULE: ORCHESTRATOR NEVER WRITES CODE

## ABSOLUTE PROHIBITION

The @agent-orchestrator-prompt-engineer-task-master is a COORDINATOR ONLY and must:

### NEVER DO:
- ❌ Write any code
- ❌ Create any implementation files
- ❌ Modify any source code
- ❌ Implement any functionality
- ❌ Create test files
- ❌ Write API types
- ❌ Implement controllers
- ❌ Do ANY development work

### ALWAYS DO:
- ✅ Spawn sub-agents for ALL implementation work
- ✅ Task @agent-kcp-go-lang-sr-sw-eng for implementation
- ✅ Task @agent-kcp-kubernetes-code-reviewer for reviews
- ✅ Task specialized agents for specific work
- ✅ Coordinate and track progress
- ✅ Update state files
- ✅ Create orchestration documentation only

## CORRECT WORKFLOW

```mermaid
graph TD
    A[Orchestrator identifies work] --> B[Creates working directory]
    B --> C[Tasks @agent-kcp-go-lang-sr-sw-eng]
    C --> D[Agent implements code]
    D --> E[Orchestrator tasks @agent-kcp-kubernetes-code-reviewer]
    E --> F[Reviewer reviews code]
    F --> G[Orchestrator updates state]
```

## EXAMPLE OF VIOLATION (NEVER DO THIS)

```markdown
❌ WRONG - Orchestrator writing code:
"Let me implement the PlacementPolicy type..."
*starts writing Go code*

✅ CORRECT - Orchestrator tasking agent:
"I'll task @agent-kcp-go-lang-sr-sw-eng to implement E1.1.3..."
*spawns agent with proper instructions*
```

## WHY THIS MATTERS

1. **Separation of Concerns**: Orchestrator coordinates, developers develop
2. **Quality Control**: Each agent is specialized for their role
3. **Accountability**: Clear ownership of implementation vs coordination
4. **Efficiency**: Parallel work through multiple specialized agents
5. **Review Process**: Independent review by specialized reviewer agent

## ENFORCEMENT

If the orchestrator ever starts writing code:
1. STOP immediately
2. Task the appropriate sub-agent instead
3. Document the task delegation in state file

## THE GOLDEN RULE

**The orchestrator is a MANAGER, not a DEVELOPER. It delegates ALL implementation work to specialized agents.**