# Split Plan Summary: registry-auth-types

## Overview
**Effort**: registry-auth-types  
**Total Size**: 965 lines (VIOLATION - exceeds 800 line limit)  
**Splits Required**: 2  
**Planner**: @agent-code-reviewer (code-reviewer-1756082516)  
**Planning Date**: 2025-08-25 00:40:28 UTC  

## Size Violation Details
- **Limit**: 800 lines
- **Actual**: 965 lines  
- **Overage**: 165 lines
- **Measurement Tool**: ${PROJECT_ROOT}/tools/line-counter.sh

## Split Strategy
Divided by logical package boundaries to maintain functional cohesion:

### Split 001: Authentication Components
- **Files**: pkg/auth/* (3 files) + pkg/doc.go
- **Size**: 649 lines (COMPLIANT)
- **Functionality**: All authentication types, credentials, and documentation
- **Branch**: phase1/wave1/registry-auth-types-split-001

### Split 002: Certificate Components  
- **Files**: pkg/certs/* (2 files)
- **Size**: 310 lines (COMPLIANT)
- **Functionality**: Certificate and TLS configuration types
- **Branch**: phase1/wave1/registry-auth-types-split-002

## Implementation Order
1. **Split 001 First** - Core authentication types
2. **Split 002 Second** - Certificate types (can reference auth types if needed)

## Key Benefits of This Split
- Clean package boundary separation
- Each split is independently compilable
- Both splits well under limit (headroom for fixes)
- Logical functional grouping maintained
- No file duplication between splits

## Files Created
- `SPLIT-PLAN-001.md` - Detailed plan for authentication split
- `SPLIT-PLAN-002.md` - Detailed plan for certificate split
- `work-log.md` - Updated with split planning details

## Next Steps for Orchestrator
1. Create split-001 and split-002 directories
2. Set up branches for each split
3. Assign SW Engineers to implement splits sequentially
4. Each split gets independent review
5. Merge splits back to parent branch after all complete

## Compliance Status
✅ Both splits under 800 line limit  
✅ No file appears in multiple splits  
✅ Clear boundaries defined  
✅ Implementation instructions provided  
✅ Single reviewer handled all planning (R199)