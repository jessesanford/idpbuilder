# PreCompact Hook Limitations and SF 2.0 Solution

## The Reality of PreCompact Hooks

### What PreCompact Hooks CANNOT Do:
1. **Cannot communicate with agents** - Hooks are isolated scripts with no access to Claude's internal state
2. **Cannot force agents to save TODOs** - No ability to trigger agent actions
3. **Cannot access in-memory data** - Only works with files already on disk
4. **Cannot delay compaction** - Compaction proceeds regardless of hook success
5. **Cannot guarantee preservation** - If TODOs aren't saved to disk, they're lost

### What PreCompact Hooks CAN Do:
1. **Read existing files** from the filesystem
2. **Create marker files** for post-compaction detection
3. **Copy/backup files** that already exist
4. **Log information** about what was preserved

## The Software Factory 2.0 Solution

Since PreCompact hooks cannot technically enforce TODO saving, we use **behavioral enforcement through rules**:

### 1. Mandatory Save Rules (R187-R190)

#### R187 - TODO Save Triggers
- **BLOCKING criticality** - Agents MUST save TODOs:
  - Within 30 seconds of using TodoWrite tool
  - Before any state transition
  - Before spawning other agents
  - After completing milestones
  - When encountering errors

#### R188 - TODO Save Frequency
- **Every 10 messages** exchanged
- **Every 15 minutes** of active work
- **After every 200 lines** of code
- **Before high-memory operations**

#### R189 - TODO Commit Protocol
- **Must commit within 60 seconds** of save
- **Must push to remote** immediately
- **Must verify push success**

#### R190 - TODO Recovery Verification
- **Must verify recovery** after compaction
- **Must use TodoWrite tool** to load (not just read)
- **Must deduplicate** after loading

### 2. Grading Enforcement

Agents face severe penalties for non-compliance:
- Missing save after TodoWrite: **-20%**
- Missing 15-minute save: **-15%**
- No saves before compaction: **-50%**
- Lost TODOs during compaction: **-75%**
- Complete TODO loss: **-100% (IMMEDIATE FAILURE)**

### 3. Agent Acknowledgment

All agents must acknowledge these rules at startup:
```bash
R187: TODO save triggers [BLOCKING]
R188: TODO save frequency [BLOCKING]
R189: TODO commit protocol [BLOCKING]
R190: TODO recovery verification [BLOCKING]
```

## How It Works in Practice

### Before Compaction:
1. **Agents save TODOs regularly** (compelled by rules, not hooks)
2. **TODOs are committed and pushed** to git repository
3. **Files exist on disk** in `todos/` directory

### During Compaction:
1. **PreCompact hook runs** (`utilities/pre-compact.sh`)
2. **Finds latest TODO file** from disk (line 223)
3. **Copies to backup location** (`/tmp/todos-precompact.txt`)
4. **Creates marker file** (`/tmp/compaction_marker.txt`)

### After Compaction:
1. **Agents detect marker** (automatic check in pre-flight)
2. **Find their TODO files** in `todos/` directory
3. **Load into TodoWrite tool** (not just read!)
4. **Continue work** with preserved state

## The Critical Understanding

**PreCompact hooks are a last-resort safety net, NOT the primary solution.**

The real solution is:
1. **Behavioral enforcement** through strict rules
2. **Grading penalties** for non-compliance
3. **Regular saves** as standard practice
4. **Git persistence** for true durability

## Why This Works

1. **Agents are motivated** by grading penalties
2. **Rules are BLOCKING** - cannot proceed without compliance
3. **Multiple save triggers** ensure frequent persistence
4. **Git provides backup** beyond local filesystem
5. **Recovery is verifiable** through R190

## The Bottom Line

- **PreCompact hooks CANNOT save agent memory**
- **Agents MUST save their own TODOs regularly**
- **Rules R187-R190 compel this behavior**
- **Grading enforcement ensures compliance**
- **The hook only preserves what's already saved**

This is why R187 ends with:
> **Remember:** PreCompact CANNOT save your TODOs for you. Only YOU can prevent TODO loss!

And why all agents must acknowledge these rules as BLOCKING requirements.