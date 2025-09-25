# 🚨🚨🚨 BLOCKING RULE R323: MANDATORY FINAL ARTIFACT BUILD 🚨🚨🚨

## Rule Definition
**ID**: R323.0.0  
**Category**: Artifact Management / Project Completion  
**Criticality**: 🚨🚨🚨 BLOCKING  
**Enforcement**: MANDATORY - NO EXCEPTIONS

## Overview
NO project, phase, or major effort can be marked as complete without building and verifying the final deliverable artifact. Projects without built artifacts are fundamentally incomplete and represent a CRITICAL FAILURE of the Software Factory.

## 🔴🔴🔴 ABSOLUTE REQUIREMENTS 🔴🔴🔴

### 1. BUILD_VALIDATION State Requirements
The Code Reviewer MUST during BUILD_VALIDATION:
```bash
# 1. Check for build system
if [ -f Makefile ]; then
    make clean
    make || make build || make all
elif [ -f package.json ]; then
    npm run build || npm run compile
elif [ -f build.sh ]; then
    ./build.sh
elif [ -f CMakeLists.txt ]; then
    mkdir -p build && cd build && cmake .. && make
fi

# 2. Verify artifact exists
# Check for common artifact patterns
ls -la bin/ dist/ build/ target/ out/ *.exe *.app *.jar *.war

# 3. Document artifact details
echo "ARTIFACT_PATH: $(find . -type f -executable -o -name '*.jar' -o -name '*.exe' | head -1)"
echo "ARTIFACT_SIZE: $(du -h $ARTIFACT_PATH)"
echo "ARTIFACT_TYPE: $(file $ARTIFACT_PATH)"
```

### 2. Orchestrator Verification
The Orchestrator MUST during BUILD_VALIDATION:
- Verify Code Reviewer built the artifact
- Confirm artifact path in validation report
- BLOCK transition to SUCCESS if no artifact exists
- Document artifact location in state file

### 3. PRODUCTION_READY_VALIDATION Requirements
MUST include:
- Test execution using the ACTUAL built artifact
- Verification that artifact runs correctly
- Performance testing of the built artifact
- Documentation of artifact behavior

### 4. SUCCESS State Prerequisites
Cannot transition to SUCCESS without:
- `final_artifact_path` documented in state file
- `artifact_size` recorded
- `artifact_type` identified
- `build_command` that created it documented

## 🚨 Grading Penalties

### SEVERE VIOLATIONS (-50% to -100%)
- **No artifact built**: -50% immediate penalty
- **SUCCESS without artifact**: -75% penalty
- **Claiming completion without deliverable**: -100% FAIL

### MODERATE VIOLATIONS (-25%)
- Artifact built but not documented
- Artifact location not in state file
- Build command not recorded

### MINOR VIOLATIONS (-10%)
- Artifact built late in process
- Missing artifact metadata

## Implementation Checklist

### For Code Reviewer Agent
```yaml
BUILD_VALIDATION:
  mandatory_tasks:
    - identify_build_system
    - execute_build_command
    - verify_artifact_exists
    - document_artifact_location
    - test_artifact_execution
    - report_artifact_details
```

### For Orchestrator Agent
```yaml
BUILD_VALIDATION:
  verification:
    - artifact_built: true
    - artifact_path: /path/to/artifact
    - artifact_tested: true
    - can_transition: only_if_artifact_exists

SUCCESS:
  required_fields:
    - final_artifact_path
    - artifact_size
    - artifact_type
    - build_timestamp
```

## Common Artifact Patterns

### By Language/Framework
- **Go**: Binary in current dir or `bin/`
- **Java**: `*.jar` in `target/` or `build/`
- **Node.js**: Bundle in `dist/` or `build/`
- **Python**: Package in `dist/` or executable
- **C/C++**: Binary in `build/` or `bin/`
- **Rust**: Binary in `target/release/`

### Build Commands Priority
1. `make` or `make build` (if Makefile exists)
2. `npm run build` (if package.json exists)
3. `go build` (if go.mod exists)
4. `mvn package` (if pom.xml exists)
5. `gradle build` (if build.gradle exists)
6. `cargo build --release` (if Cargo.toml exists)

## Validation Report Format

```markdown
## BUILD VALIDATION REPORT

### Artifact Build Status
- Build System: [Makefile/npm/gradle/etc]
- Build Command: [exact command used]
- Build Duration: [time taken]
- Build Success: ✅/❌

### Final Artifact Details
- **Path**: `/absolute/path/to/artifact`
- **Size**: [size in MB]
- **Type**: [executable/jar/bundle/etc]
- **Permissions**: [file permissions]
- **Dependencies**: [if applicable]

### Artifact Verification
- Execution Test: ✅/❌
- Basic Functionality: ✅/❌
- Performance Check: ✅/❌
```

## State File Requirements

```yaml
final_artifact:
  path: /path/to/artifact
  size: "15.2MB"
  type: "executable"
  build_command: "make build"
  build_timestamp: "2024-01-20T10:30:00Z"
  verified: true
  test_status: "passed"
```

## Enforcement Message

When violation detected:
```
🚨🚨🚨 R323 VIOLATION: NO FINAL ARTIFACT BUILT! 🚨🚨🚨

A Software Factory project WITHOUT a built artifact is like a car factory 
that never builds cars - it's a FUNDAMENTAL FAILURE!

Required Actions:
1. STOP all other work immediately
2. Identify the build system
3. Execute the build command
4. Verify artifact exists
5. Test the artifact
6. Document in state file

This is BLOCKING - cannot proceed without artifact!
```

## Related Rules
- R007: Size limit compliance (applies to source, not artifacts)
- R206: State machine validation (BUILD_VALIDATION state)
- R288: State file updates (must include artifact info)

## Implementation Notes

1. **Build Early**: Start building artifacts during early validation
2. **Test Actual Artifact**: Don't just test source - test the built artifact
3. **Document Everything**: Path, size, type, build command
4. **Verify Repeatedly**: Check artifact exists at multiple gates

## Why This Rule Exists

The Software Factory exists to BUILD SOFTWARE. A project marked as "SUCCESS" without a built artifact is like:
- A bakery that never bakes bread
- A factory that never produces products
- A construction project with no building

This represents a FUNDAMENTAL FAILURE of purpose and must be prevented at all costs.

---
*Rule Created: 2024-01-20*
*Severity: 🚨🚨🚨 BLOCKING*
*Penalty: -50% to -100%*