# Phase 1 Integration Instructions for Integration Agent

## Mission

Integrate ALL Phase 1 wave branches into a single phase integration branch following Sequential Rebuild Model (R009/R282/R283).

## Integration Context

- **Phase**: 1
- **Waves to Integrate**: 2 waves (Wave 1 and Wave 2)
- **Integration Type**: Phase integration (combines multiple waves)
- **Target Repository**: https://github.com/jessesanford/idpbuilder.git
- **Base Branch**: main (as specified in target-repo-config.yaml)

## Wave Branches to Integrate

### Wave 1 (CONVERGED)
- **Branch**: `idpbuilder-oci-push/phase1/wave1/integration`
- **Status**: CONVERGED (0 bugs remaining, build passing)
- **Efforts**: 1.1.1, 1.1.2, 1.1.3, 1.1.4

### Wave 2 (CONVERGED)
- **Branch**: `idpbuilder-oci-push/phase1/wave2/integration`
- **Status**: CONVERGED (0 bugs remaining, build passing)
- **Base**: Built on wave1 integration (cascade branching per R308)
- **Efforts**: 1.2.1, 1.2.2, 1.2.3

## Integration Strategy (Sequential Rebuild Model - R009/R282/R283)

**CRITICAL**: This is a PHASE integration following Sequential Rebuild Model:

### Base Branch Selection (R282/R283)
- **Base**: First effort of PHASE = effort from Wave 1 (`idpbuilder-oci-push/phase1/wave1/integration`)
- **NOT**: Last wave's integration branch (that's CASCADE, different model!)
- **Reason**: Phase integration tests SEQUENTIAL MERGEABILITY from phase start

### Sequential Merge Pattern
1. Clone from `idpbuilder-oci-push/phase1/wave1/integration` (first wave = phase base)
2. Create new branch: `idpbuilder-oci-push/phase1/integration`
3. Merge wave2 integration into the phase branch
4. Resolve any conflicts
5. Build and validate
6. Push to remote

## Workspace Setup

1. **Create Phase Integration Workspace**:
   ```bash
   mkdir -p efforts/phase1/integration
   cd efforts/phase1/integration
   ```

2. **Clone from Base Branch**:
   ```bash
   git clone --single-branch --branch idpbuilder-oci-push/phase1/wave1/integration \
     https://github.com/jessesanford/idpbuilder.git .
   ```

3. **Create Phase Integration Branch**:
   ```bash
   git checkout -b idpbuilder-oci-push/phase1/integration
   ```

## Merge Sequence

### Merge Wave 2
```bash
git merge origin/idpbuilder-oci-push/phase1/wave2/integration --no-ff \
  -m "Integrate Wave 2 into Phase 1 integration

Wave 2 adds TLS/certificate management capabilities:
- Certificate validation and management
- TLS configuration
- Security protocols

Per R308 incremental branching and R009 sequential rebuild model.
"
```

**Expected**: Clean merge (Wave 2 was built on Wave 1, so should merge cleanly)

## Conflict Resolution Protocol

If conflicts occur (unlikely given cascade branching):

1. **Document Conflicts**: Record which files have conflicts
2. **Analyze Root Cause**: Determine why cascade branching didn't prevent this
3. **Resolve Conflicts**: Maintain functionality from both waves
4. **Validate Resolution**: Ensure both wave features work after resolution
5. **Document in Report**: Include conflict details in integration report

## Build Validation (R265/R323)

After merging all waves:

1. **Run Build**:
   ```bash
   make build
   ```

2. **Verify Success**: Check that binary is created

3. **Record Metrics**:
   - Build duration
   - Binary size
   - Any warnings

4. **Create Build Report**: Document validation results per R323

## Testing Requirements

Per R265, execute comprehensive testing:

1. **Unit Tests**: Run all unit tests for merged code
2. **Integration Tests**: Verify wave interactions work correctly
3. **Build Validation**: Confirm artifact generation
4. **Functionality Check**: Basic smoke tests

## Success Criteria

Phase integration is successful when:

- ✅ Wave 2 branch merged cleanly into phase branch
- ✅ No unresolved conflicts
- ✅ Build completes successfully
- ✅ All tests pass
- ✅ Binary artifact generated
- ✅ Phase integration branch pushed to remote
- ✅ Integration report created

## Deliverables

1. **Integration Report**: Comprehensive markdown report with:
   - Merge details (commits, conflicts, resolutions)
   - Build validation results
   - Test execution summary
   - Artifacts generated
   - Recommendations for next steps

2. **Updated Branch**: `idpbuilder-oci-push/phase1/integration` pushed to remote

3. **Build Artifact**: Binary generated and validated per R323

## Integration Rules Compliance

This integration MUST follow:

- **R009**: Sequential Rebuild Model for phase integration
- **R282**: Phase integration base = first effort of phase
- **R283**: Validation of sequential mergeability
- **R307**: Integration iteration protocol
- **R308**: Incremental branching strategy (why waves cascade)
- **R329**: Only Integration Agent performs merges
- **R265**: Comprehensive testing requirements
- **R323**: Artifact generation and verification

## Notes for Integration Agent

- **Cascade Context**: Wave 2 was built ON TOP of Wave 1 (R308 cascade branching)
- **Expected Outcome**: Clean merge since Wave 2 already contains Wave 1's changes
- **If Conflicts Occur**: This indicates a problem - Wave 2 should have had Wave 1's code
- **Sequential Rebuild**: We're testing that phase can be rebuilt from scratch sequentially
- **Next Phase Integration**: This creates the base for Phase 2 integration later

## Orchestrator Will Monitor For

- Integration report completion
- Branch pushed to remote
- Build validation results
- Any error conditions requiring intervention

**Good luck with the integration!**
