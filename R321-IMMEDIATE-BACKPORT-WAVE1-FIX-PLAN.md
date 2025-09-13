# R321 IMMEDIATE BACKPORT FIX PLAN - Wave 1 Duplicate Test Helpers

## 🔴 CRITICAL R321 VIOLATION DETECTED
**Issue**: Wave 1 integration has duplicate test helpers
**Detection Time**: 2025-09-12T00:56:55Z
**Action Required**: Immediate backport to source branches per R321

## Problem Description
The Wave 1 integration branch contains duplicate test helper functions that were created independently in different effort branches. When merged together, these duplicates cause:
- Build failures due to redeclared functions
- Test execution errors
- Integration branch cannot build successfully

## Affected Wave 1 Efforts
1. **kind-cert-extraction** - Branch: idpbuilder-oci-build-push/phase1/wave1/kind-cert-extraction
2. **registry-tls-trust** - Branch: idpbuilder-oci-build-push/phase1/wave1/registry-tls-trust  
3. **registry-auth-types-split-001** - Branch: idpbuilder-oci-build-push/phase1/wave1/registry-auth-types-split-001
4. **registry-auth-types-split-002** - Branch: idpbuilder-oci-build-push/phase1/wave1/registry-auth-types-split-002

## Fix Strategy

### Step 1: Identify Duplicate Test Helpers
```bash
# In each effort directory, check for test helper functions
cd efforts/phase1/wave1/kind-cert-extraction
grep -r "func Test" pkg/
grep -r "func setup" pkg/
grep -r "func teardown" pkg/

# Look for common test utilities
find . -name "*_test.go" -exec grep -l "TestMain\|testHelper\|mockData" {} \;
```

### Step 2: Consolidate Test Helpers
Create a shared test utilities package that all efforts can use:

1. Create `pkg/testutil/helpers.go` in ONE effort (registry-tls-trust as the base)
2. Move all common test helpers to this package
3. Update other efforts to import and use the shared helpers
4. Remove duplicate definitions from other efforts

### Step 3: Fix Each Effort Branch

#### For kind-cert-extraction:
```bash
cd efforts/phase1/wave1/kind-cert-extraction
git checkout idpbuilder-oci-build-push/phase1/wave1/kind-cert-extraction
git pull origin idpbuilder-oci-build-push/phase1/wave1/kind-cert-extraction

# Remove duplicate test helpers
# Update imports to use shared testutil package
# Ensure tests still pass
go test ./...

git add -A
git commit -m "fix(R321): remove duplicate test helpers, use shared testutil package"
git push origin idpbuilder-oci-build-push/phase1/wave1/kind-cert-extraction
```

#### For registry-tls-trust (base with shared helpers):
```bash
cd efforts/phase1/wave1/registry-tls-trust
git checkout idpbuilder-oci-build-push/phase1/wave1/registry-tls-trust
git pull origin idpbuilder-oci-build-push/phase1/wave1/registry-tls-trust

# Create the shared testutil package
mkdir -p pkg/testutil
# Move common helpers here
# Ensure all tests pass
go test ./...

git add -A
git commit -m "fix(R321): create shared testutil package for Wave 1"
git push origin idpbuilder-oci-build-push/phase1/wave1/registry-tls-trust
```

#### For registry-auth-types splits:
```bash
# For split-001
cd efforts/phase1/wave1/registry-auth-types-split-001
git checkout idpbuilder-oci-build-push/phase1/wave1/registry-auth-types-split-001
# Remove duplicates, update imports
go test ./...
git commit -m "fix(R321): remove duplicate test helpers from split-001"
git push

# For split-002
cd efforts/phase1/wave1/registry-auth-types-split-002
git checkout idpbuilder-oci-build-push/phase1/wave1/registry-auth-types-split-002
# Remove duplicates, update imports
go test ./...
git commit -m "fix(R321): remove duplicate test helpers from split-002"
git push
```

### Step 4: Verify Each Branch Independently
```bash
# Each effort must build and test independently
for effort in kind-cert-extraction registry-tls-trust registry-auth-types-split-001 registry-auth-types-split-002; do
    cd efforts/phase1/wave1/$effort
    echo "Testing $effort..."
    go build ./...
    go test ./...
done
```

### Step 5: Create Verification Marker
```bash
# After all fixes are complete
touch efforts/phase1/wave1/R321_FIXES_COMPLETE.flag
echo "All Wave 1 duplicate test helpers fixed at $(date)" > efforts/phase1/wave1/R321_FIXES_COMPLETE.flag
```

## Success Criteria
- ✅ Each Wave 1 effort branch builds independently
- ✅ Each Wave 1 effort branch tests pass independently  
- ✅ No duplicate function definitions across efforts
- ✅ Shared test utilities in testutil package
- ✅ All branches pushed to remote
- ✅ R321_FIXES_COMPLETE.flag created

## Important Notes
- **DO NOT** fix in integration branch - R321 MANDATES source branch fixes
- **DO NOT** proceed with phase integration until these fixes are complete
- **MUST** verify each branch works independently before re-integration
- After fixes complete, orchestrator will transition to retry phase integration

## Timeline
- Expected completion: 30-45 minutes
- Must complete before phase integration can continue
- R321 is SUPREME LAW - no exceptions