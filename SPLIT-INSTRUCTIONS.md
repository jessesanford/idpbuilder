# Split Implementation Instructions for SW Engineers

## Overview
The registry-auth-types effort (965 lines) has been split into 2 independent parts to comply with the 800-line limit. These splits can be implemented in PARALLEL by different engineers or SEQUENTIALLY by one engineer.

## Split Summary
| Split | Description | Size | Branch Name |
|-------|-------------|------|-------------|
| 001 | OCI Types & Documentation | 661 lines | phase1/wave1/registry-auth-types-split-001 |
| 002 | Stack Types | 313 lines | phase1/wave1/registry-auth-types-split-002 |

## CRITICAL RULES FOR SW ENGINEERS

### 1. Size Compliance (MANDATORY)
```bash
# ALWAYS measure from WITHIN your split directory:
cd /path/to/your/split/directory
/workspaces/idpbuilder-oci-mgmt/tools/line-counter.sh  # NO PARAMETERS!

# STOP if size > 700 lines
# MUST stay under 800 lines
```

### 2. File Isolation (STRICT)
- Work ONLY on files assigned to your split
- DO NOT modify files from other splits
- DO NOT create shared dependencies

### 3. Branch Management
```bash
# For Split 001:
git checkout -b phase1/wave1/registry-auth-types-split-001

# For Split 002:
git checkout -b phase1/wave1/registry-auth-types-split-002
```

## Implementation Process

### For Split 001 Engineer (OCI Types)
1. **Setup**:
   ```bash
   git checkout -b phase1/wave1/registry-auth-types-split-001
   cd efforts/phase1/wave1/registry-auth-types
   ```

2. **Files to Work On**:
   - pkg/doc.go (39 lines)
   - pkg/oci/types.go (121 lines)
   - pkg/oci/manifest.go (124 lines)
   - pkg/oci/constants.go (56 lines)
   - pkg/oci/types_test.go (130 lines)
   - pkg/oci/manifest_test.go (191 lines)

3. **Testing**:
   ```bash
   go test ./pkg/oci/...
   ```

4. **Size Check**:
   ```bash
   /workspaces/idpbuilder-oci-mgmt/tools/line-counter.sh
   # Must show ≤661 lines
   ```

5. **Commit**:
   ```bash
   git add pkg/doc.go pkg/oci/
   git commit -m "feat: implement OCI types and manifest handling (split 001 of 2)"
   git push origin phase1/wave1/registry-auth-types-split-001
   ```

### For Split 002 Engineer (Stack Types)
1. **Setup**:
   ```bash
   git checkout -b phase1/wave1/registry-auth-types-split-002
   cd efforts/phase1/wave1/registry-auth-types
   ```

2. **Files to Work On**:
   - pkg/stack/types.go (107 lines)
   - pkg/stack/constants.go (42 lines)
   - pkg/stack/types_test.go (164 lines)

3. **Testing**:
   ```bash
   go test ./pkg/stack/...
   ```

4. **Size Check**:
   ```bash
   /workspaces/idpbuilder-oci-mgmt/tools/line-counter.sh
   # Must show ≤313 lines
   ```

5. **Commit**:
   ```bash
   git add pkg/stack/
   git commit -m "feat: implement Stack types and validation (split 002 of 2)"
   git push origin phase1/wave1/registry-auth-types-split-002
   ```

## Coordination Requirements

### For Parallel Implementation:
- **NO coordination needed** - splits are completely independent
- Each engineer works on their assigned split only
- Both can start immediately
- Both can be reviewed independently

### For Sequential Implementation:
- Start with either split (order doesn't matter)
- Complete and push first split
- Switch to second split branch
- Complete and push second split

## Quality Requirements
1. **Test Coverage**: Minimum 80% for each split
2. **Documentation**: Godoc comments on all public types
3. **No Lint Errors**: Run `golangci-lint run`
4. **Clean Commits**: One commit per split
5. **Size Compliance**: Stay under limits

## Common Issues and Solutions

### Issue: Size exceeds limit
**Solution**: Stop immediately, notify Code Reviewer for re-split

### Issue: Tests failing
**Solution**: Ensure you're only testing your split's packages

### Issue: Import conflicts
**Solution**: Splits should have no cross-dependencies

### Issue: Line counter shows wrong size
**Solution**: Ensure you're in the correct directory and using NO parameters

## Completion Checklist
- [ ] Branch created with correct name
- [ ] Only assigned files modified
- [ ] Tests passing (>80% coverage)
- [ ] Size verified (<800 lines, ideally <700)
- [ ] Code committed with clear message
- [ ] Branch pushed to origin
- [ ] Ready for Code Review

## After Implementation
1. Notify orchestrator that split is complete
2. Wait for Code Reviewer to validate
3. Do not merge until all splits are reviewed
4. Splits will be integrated back to main effort branch

## Important Notes
- These splits were created from an existing 965-line implementation
- The code already exists and has been tested
- Your job is to properly separate it into compliant chunks
- Maintain all existing functionality and tests
- Do not add new features during splitting