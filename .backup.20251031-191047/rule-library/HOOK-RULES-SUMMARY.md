# Hook and Utility Rules Summary

## New Rules Added to Registry (R171-R175)

### R171.0.0 - PreCompact Hook Configuration
- **Purpose**: Define the ONLY automatic hook used in SF 2.0
- **Key Point**: PreCompact hook in `.claude/settings.json` with inline commands
- **File**: `R171-precompact-hook.md`

### R172.0.0 - Utility Script Execution  
- **Purpose**: Clarify that scripts in `utilities/` are manual, not hooks
- **Key Point**: Must be executed manually, never automatic
- **File**: `R172-utility-scripts.md`

### R173.0.0 - State Preservation Protocol
- **Purpose**: Define dual approach - automatic (minimal) + manual (comprehensive)
- **Key Point**: PreCompact creates marker, utilities provide full preservation
- **File**: `R173-state-preservation.md`

### R174.0.0 - Context Recovery Detection
- **Purpose**: Mandate checking for `/tmp/compaction_marker.txt` on startup
- **Key Point**: All agents must detect and respond to compaction
- **File**: `R174-recovery-detection.md`

### R175.0.0 - Manual Utility Usage
- **Purpose**: Document that utilities require explicit manual execution
- **Key Point**: Clear documentation prevents automation expectations
- **File**: `R175-manual-utilities.md`

## Rule Compliance Checklist

### For Setup Scripts
- [ ] Only configure PreCompact hook (R171)
- [ ] Place scripts in `utilities/` not `hooks/` (R172)
- [ ] Make utilities executable (R175)
- [ ] Document manual execution requirement (R175)

### For Agents
- [ ] Check for compaction marker on startup (R174)
- [ ] Include recovery procedures (R173)
- [ ] Document when to run utilities (R175)
- [ ] Clean up marker after reading (R174)

### For Documentation
- [ ] Call them "utilities" not "hooks" (R172)
- [ ] Explain automatic vs manual clearly (R173)
- [ ] Provide execution examples (R175)
- [ ] Set correct expectations (R171)

## Validation Commands

```bash
# Verify correct structure
ls -la utilities/           # Should exist
ls -la hooks/ 2>&1          # Should NOT exist

# Check hook configuration
cat .claude/settings.json | grep -c "PreCompact"  # Should be >0
cat .claude/settings.json | grep -c "post_compact" # Should be 0

# Test marker detection
echo "TEST" > /tmp/compaction_marker.txt
# Run agent - should detect and clean up
```

## Common Violations to Avoid

1. ❌ Creating a `hooks/` directory (violates R172)
2. ❌ Configuring fictional hooks in settings.json (violates R171)
3. ❌ Expecting automatic utility execution (violates R172, R175)
4. ❌ Not checking for marker file (violates R174)
5. ❌ Only using automatic preservation (violates R173)

## Correct Implementation

✅ PreCompact hook in settings.json (automatic)
✅ Utility scripts in utilities/ (manual)
✅ Marker detection in all agents
✅ Clear documentation about manual execution
✅ Dual preservation strategy