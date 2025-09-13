# Wave 1 Re-Integration Plan (R327 Cascade)

## Context
Wave 1 integration is STALE and must be recreated with R321 fixes applied to all effort branches.

## Integration Requirements

### R327 Cascade Requirements
- Wave 1 must be re-integrated first (cascade starting point)
- All R321 fixes must be included from effort branches
- Fresh integration branch from main
- Proper merge order to avoid conflicts

### R321 Fixes Applied
- **2025-09-12 01:21:45 UTC**: All duplicate test helpers eliminated
- Shared testutil package created
- Each effort branch updated and pushed

## Integration Order

1. **Base Effort First**: E1.1.2-registry-tls-trust
   - Contains shared pkg/testutil/helpers.go
   - Other efforts depend on this

2. **Dependent Efforts**: 
   - E1.1.1-kind-cert-extraction (uses shared testutil)
   - E1.1.3-registry-auth-types-split-001 (simplified contains())
   - E1.1.3-registry-auth-types-split-002 (uses shared testutil)

## Integration Steps

1. Create fresh integration branch from main
2. Create integration workspace
3. Merge efforts in dependency order:
   - registry-tls-trust (base with shared utils)
   - kind-cert-extraction
   - registry-auth-types-split-001
   - registry-auth-types-split-002
4. Verify build and tests pass
5. Create R291 demo scripts
6. Push integration branch
7. Update orchestrator-state.json

## Success Criteria
- [ ] All efforts merged without conflicts
- [ ] Build passes
- [ ] Tests pass
- [ ] R291 demos functional
- [ ] Integration branch pushed
- [ ] State updated: stale_integrations.wave1.recreation_completed = true

## Next Steps After Completion
- Wave 2 re-integration (based on fixed Wave 1)
- Phase 1 re-integration (based on fixed Wave 2)
- Architect review of fixed integrations