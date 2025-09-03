# Integration Agent Core Rules

## Identity Rules
- R260 - Integration Agent Core Requirements
- R267 - Integration Agent Grading Criteria

## Planning Rules
- R261 - Integration Planning Requirements

## Operation Rules
- R262 - Merge Operation Protocols (SUPREME LAW)
- R264 - Work Log Tracking Requirements
- R265 - Integration Testing Requirements
- R266 - Upstream Bug Documentation Protocol (SUPREME LAW)

## Documentation Rules
- R263 - Integration Documentation Requirements

## General Rules Also Applied
- R002 - Agent Acknowledgment Protocol
- R013 - Git Operations
- R014 - Branch Naming Convention
- R015 - Commit Message Format

## State Machine States
1. **INIT** - Startup and acknowledgment
2. **PLANNING** - Analyze branches, create plan
3. **MERGING** - Execute integration per plan
4. **TESTING** - Build and test (document only)
5. **REPORTING** - Complete documentation
6. **COMPLETED** - Terminal state

## Grading Breakdown
- **50% Completeness**
  - 20% Branch merging
  - 15% Conflict resolution
  - 10% Branch integrity
  - 5% Final validation
- **50% Documentation**
  - 25% Work log quality
  - 25% Integration report quality