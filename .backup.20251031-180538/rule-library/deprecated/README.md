# Deprecated Rules Directory

**Purpose**: Preserve deprecated rules for historical reference and migration support

## Directory Structure

```
deprecated/
├── sf2/          # SF 2.0-specific rules deprecated in SF 3.0
└── README.md     # This file
```

## SF 2.0 Deprecated Rules

The `sf2/` directory contains rules that were deprecated during the SF 2.0 → SF 3.0 migration:

| Rule | Title | Reason | SF 3.0 Replacement |
|------|-------|--------|-------------------|
| R187 | TODO Save Triggers | Consolidated into comprehensive TODO persistence rule | R287 |
| R188 | TODO Save Frequency | Consolidated into comprehensive TODO persistence rule | R287 |
| R189 | TODO Commit Protocol | Consolidated into comprehensive TODO persistence rule | R287 |
| R190 | TODO Recovery Verification | Consolidated into comprehensive TODO persistence rule | R287 |
| R236 | Mandatory State Rule Reading | Consolidated into unified state rule enforcement | R290 |
| R237 | State Rule Verification Enforcement | Consolidated into unified state rule enforcement | R290 |
| R252 | Mandatory State File Updates | Consolidated into atomic update protocol | R288 |
| R253 | Mandatory State File Commit/Push | Consolidated into atomic update protocol | R288 |
| R296 | Branch Marking Protocol | Split workflow obsolete in SF 3.0 | N/A (workflow change) |

**Total Deprecated**: 9 rules

## Why Preserve Deprecated Rules?

1. **Historical Reference**: Understanding rule evolution and decision-making
2. **Migration Support**: Projects mid-transition may still reference old rules
3. **Audit Trail**: Complete documentation of system changes over time
4. **Educational Value**: Learning from past approaches and improvements

## Using Deprecated Rules

### For SF 3.0 Projects
- **DO NOT** reference deprecated rules in new work
- Use SF 3.0 replacement rules instead (see table above)
- Deprecated rules are for reference only

### For SF 2.0 Migration Projects
- Complete current work with SF 2.0 rules if mid-wave
- Transition to SF 3.0 rules at wave/phase boundaries
- See `docs/DEPRECATED-RULES-SF2.md` for migration guidance

## Related Documentation

- **Deprecated Rules Details**: `/docs/DEPRECATED-RULES-SF2.md`
- **New SF 3.0 Rules**: `/docs/NEW-RULES-SF3.md`
- **Migration Report**: `/sf3-implementation/SF3-MIGRATION-COMPLETE-REPORT.md`
- **SF 3.0 Architecture**: `/docs/SOFTWARE-FACTORY-3.0-ARCHITECTURE.md`

---

**Created**: 2025-10-22 (Week 13 - SF 2.0 → SF 3.0 Cut-Over)
**Status**: Complete
