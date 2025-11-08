# Rule R013: Git Operations Protocol

## Rule Statement
All agents MUST follow standardized git operations protocols. Agents must verify repository state, use proper branch naming, commit atomically with descriptive messages, and push changes immediately after commits to prevent data loss.

## Criticality Level
**BLOCKING** - Improper git operations can corrupt repositories and lose work

## Enforcement Mechanism
- **Technical**: Verify git status before and after operations
- **Behavioral**: Atomic commits with immediate push
- **Grading**: -30% for uncommitted work, -50% for repository corruption

## Core Principle

```
Git Operations = Verify State → Make Changes → Commit Atomically → Push Immediately
NEVER leave uncommitted changes
NEVER work in wrong branch
ALWAYS push after commit
```

## Detailed Requirements

### Git Operation Standards

1. **Pre-Operation Verification**
   ```bash
   git status          # Check current state
   git branch          # Verify correct branch
   git remote -v       # Verify correct remote
   ```

2. **Commit Protocol**
   ```bash
   git add <specific-files>     # Add only necessary files
   git commit -m "type: clear, descriptive message"
   git push                      # Push immediately
   ```

3. **Branch Operations**
   ```bash
   # Always create from correct base
   git checkout <base-branch>
   git pull
   git checkout -b <new-branch>
   git push -u origin <new-branch>
   ```

### Forbidden Practices
- ❌ NEVER use `git add .` without verification
- ❌ NEVER force push to shared branches
- ❌ NEVER work with uncommitted changes
- ❌ NEVER modify git config without permission

### Recovery Protocol
If git operations fail:
1. Save work to temporary location
2. Verify repository state
3. Fix issues
4. Reapply changes
5. Commit and push immediately

## Relationship to Other Rules
- **R182**: Verify git repository
- **R184**: Verify git branch  
- **R312**: Git config immutability
- **R316**: Orchestrator commit restrictions

## Implementation Notes
- Agents must check git status every 5 operations
- All commits must include descriptive messages
- Push must occur within 60 seconds of commit