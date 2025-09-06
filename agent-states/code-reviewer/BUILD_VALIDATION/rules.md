# Code Reviewer - BUILD_VALIDATION State Rules

## 🔴🔴🔴 STOP! STATE RULE READING IS ABSOLUTELY FIRST! 🔴🔴🔴

**YOU HAVE ENTERED BUILD_VALIDATION STATE - YOU MUST READ AND ACKNOWLEDGE ALL STATE RULES BEFORE DOING ANY STATE WORK!**

## State Context
You are performing build validation on integrated code. This is a TECHNICAL ASSESSMENT that only Code Reviewers can perform. Orchestrators are FORBIDDEN from running builds or validating artifacts - that's YOUR job.

## 🎯 STATE OBJECTIVES

In the BUILD_VALIDATION state, you are responsible for:

1. **Complete Build Execution**
   - Run full build process
   - Compile all code
   - Generate binaries/artifacts
   - Package for deployment
   - **🚨 R323: BUILD FINAL DELIVERABLE ARTIFACT (MANDATORY)**

2. **Build Verification**
   - Verify all modules compile
   - Check for build warnings
   - Validate output artifacts
   - Ensure reproducible builds
   - **🚨 R323: VERIFY FINAL ARTIFACT EXISTS AND WORKS**

3. **Issue Documentation**
   - Document any build failures
   - Track compilation warnings
   - Note missing dependencies
   - Create comprehensive report
   - **🚨 R323: DOCUMENT ARTIFACT PATH, SIZE, TYPE**

4. **Backport Requirement Tracking**
   - Track ALL changes needed to fix build issues
   - Document which effort branches need fixes
   - Note specific fixes required for backporting

5. **🔴🔴🔴 MANDATORY FINAL ARTIFACT (R323) 🔴🔴🔴**
   - **MUST build the final deliverable binary/package**
   - **MUST verify artifact exists and is executable**
   - **MUST document artifact location and metadata**
   - **CANNOT complete without built artifact**

## 📝 REQUIRED ACTIONS

### Step 1: Clean Build Environment
```bash
cd $INTEGRATION_DIR  # Or workspace provided by orchestrator

# Clean any previous build artifacts
if [ -f "Makefile" ]; then
    make clean || true
elif [ -f "package.json" ]; then
    rm -rf node_modules dist build
elif [ -f "go.mod" ]; then
    go clean -cache
elif [ -f "Cargo.toml" ]; then
    cargo clean
elif [ -f "pom.xml" ]; then
    mvn clean
elif [ -f "build.gradle" ] || [ -f "build.gradle.kts" ]; then
    ./gradlew clean
fi

echo "✅ Build environment cleaned"
```

### Step 2: Execute Full Build (R323 MANDATORY)
```bash
# 🚨🚨🚨 R323: MUST BUILD FINAL DELIVERABLE ARTIFACT 🚨🚨🚨
# Detect and run appropriate build command
BUILD_SUCCESS=false
BUILD_OUTPUT="BUILD-OUTPUT.log"
FINAL_ARTIFACT=""

if [ -f "Makefile" ]; then
    echo "📦 Building with Make..."
    # Try multiple make targets to ensure final artifact is built
    if make clean && (make || make build || make all) > "$BUILD_OUTPUT" 2>&1; then
        BUILD_SUCCESS=true
        # Find the built artifact
        FINAL_ARTIFACT=$(find . -type f -executable -o -name "*.exe" | head -1)
    fi
elif [ -f "package.json" ]; then
    echo "📦 Building Node.js project..."
    npm install > "$BUILD_OUTPUT" 2>&1
    if npm run build >> "$BUILD_OUTPUT" 2>&1; then
        BUILD_SUCCESS=true
        FINAL_ARTIFACT=$(find dist build -type f 2>/dev/null | head -1)
    fi
elif [ -f "go.mod" ]; then
    echo "📦 Building Go project..."
    # Build with output name if possible
    PROJECT_NAME=$(basename $(pwd))
    if go build -o "$PROJECT_NAME" ./... > "$BUILD_OUTPUT" 2>&1; then
        BUILD_SUCCESS=true
        FINAL_ARTIFACT="$PROJECT_NAME"
    fi
elif [ -f "Cargo.toml" ]; then
    echo "📦 Building Rust project..."
    if cargo build --release > "$BUILD_OUTPUT" 2>&1; then
        BUILD_SUCCESS=true
    fi
elif [ -f "pom.xml" ]; then
    echo "📦 Building Maven project..."
    if mvn package > "$BUILD_OUTPUT" 2>&1; then
        BUILD_SUCCESS=true
    fi
elif [ -f "build.gradle" ] || [ -f "build.gradle.kts" ]; then
    echo "📦 Building Gradle project..."
    if ./gradlew build > "$BUILD_OUTPUT" 2>&1; then
        BUILD_SUCCESS=true
    fi
elif [ -f "CMakeLists.txt" ]; then
    echo "📦 Building CMake project..."
    mkdir -p build && cd build
    if cmake .. && make > "../$BUILD_OUTPUT" 2>&1; then
        BUILD_SUCCESS=true
    fi
    cd ..
else
    echo "⚠️ No standard build system detected"
    echo "Attempting generic compilation..."
    # Try to find and compile source files based on language
fi

if [ "$BUILD_SUCCESS" = true ]; then
    echo "✅ Build completed successfully"
else
    echo "❌ Build failed - see $BUILD_OUTPUT for details"
fi
```

### Step 3: Analyze Build Output
```bash
# Check for warnings and errors
echo "Analyzing build output..."

# Extract warnings
grep -i "warning" "$BUILD_OUTPUT" > BUILD-WARNINGS.txt || echo "No warnings found"

# Extract errors
grep -i "error" "$BUILD_OUTPUT" > BUILD-ERRORS.txt || echo "No errors found"

# Count statistics
WARNING_COUNT=$(wc -l < BUILD-WARNINGS.txt)
ERROR_COUNT=$(wc -l < BUILD-ERRORS.txt)

echo "Build Analysis:"
echo "- Warnings: $WARNING_COUNT"
echo "- Errors: $ERROR_COUNT"
```

### Step 4: Verify Artifacts (R323 CRITICAL)
```bash
# 🚨🚨🚨 R323: MANDATORY FINAL ARTIFACT VERIFICATION 🚨🚨🚨
echo "🔍 Verifying build artifacts per R323..."

ARTIFACTS_FOUND=""
ARTIFACT_DETAILS=""

# Check common artifact locations
for artifact_dir in "dist" "build" "target" "out" "bin"; do
    if [ -d "$artifact_dir" ]; then
        echo "✅ Found artifact directory: $artifact_dir"
        ls -la "$artifact_dir" | head -10
        ARTIFACTS_FOUND="yes"
    fi
done

# Check for executables (primary artifacts)
EXECUTABLES=$(find . -type f -executable -not -path "./.git/*" 2>/dev/null | grep -v "\.sh$")
if [ -n "$EXECUTABLES" ]; then
    echo "✅ Found executable artifacts:"
    echo "$EXECUTABLES"
    ARTIFACTS_FOUND="yes"
    
    # Document artifact details per R323
    for exe in $EXECUTABLES; do
        echo "📦 Artifact: $exe"
        echo "   Size: $(du -h "$exe" | cut -f1)"
        echo "   Type: $(file "$exe" | cut -d: -f2)"
        echo "   Permissions: $(ls -l "$exe" | cut -d' ' -f1)"
    done
fi

# Check for package artifacts (JARs, WARs, etc)
PACKAGES=$(find . -name "*.jar" -o -name "*.war" -o -name "*.tar.gz" -o -name "*.zip" 2>/dev/null)
if [ -n "$PACKAGES" ]; then
    echo "✅ Found package artifacts:"
    echo "$PACKAGES"
    ARTIFACTS_FOUND="yes"
fi

# R323 ENFORCEMENT: NO ARTIFACT = FAILURE
if [ -z "$ARTIFACTS_FOUND" ]; then
    echo "🚨🚨🚨 R323 VIOLATION: NO FINAL ARTIFACT BUILT! 🚨🚨🚨"
    echo "This is a BLOCKING failure - project cannot be marked complete without deliverable!"
    BUILD_SUCCESS=false
fi
```

### Step 5: Create Build Validation Report
```bash
cat > BUILD-VALIDATION-REPORT.md << EOF
# Build Validation Report
Date: $(date)
State: BUILD_VALIDATION
Agent: Code Reviewer
Workspace: $(pwd)

## Build Status
- Status: $(if [ "$BUILD_SUCCESS" = true ]; then echo "SUCCESS"; else echo "FAILED"; fi)
- Build System: [Make/npm/go/cargo/maven/gradle/other]
- Duration: [Time taken]
- Backport Required: $(if [ "$BUILD_SUCCESS" = false ]; then echo "Yes"; else echo "No"; fi)

## Build Output Analysis
- Total Warnings: $WARNING_COUNT
- Total Errors: $ERROR_COUNT
- Build Log: See BUILD-OUTPUT.log

## Artifacts Generated (R323 MANDATORY)
- **Final Artifact Path**: $FINAL_ARTIFACT
- **Artifact Exists**: $(if [ -n "$FINAL_ARTIFACT" ] && [ -f "$FINAL_ARTIFACT" ]; then echo "✅ YES"; else echo "❌ NO - R323 VIOLATION!"; fi)
- **Artifact Size**: $(if [ -f "$FINAL_ARTIFACT" ]; then du -h "$FINAL_ARTIFACT" | cut -f1; else echo "N/A"; fi)
- **Artifact Type**: $(if [ -f "$FINAL_ARTIFACT" ]; then file "$FINAL_ARTIFACT" | cut -d: -f2; else echo "N/A"; fi)
- **Build Command Used**: [make/npm run build/go build/etc]
- Location: [dist/build/target/etc]
- Files:
  $(ls -la dist/ build/ target/ out/ bin/ 2>/dev/null | head -10)

## Issues Encountered
### Compilation Errors
$(head -20 BUILD-ERRORS.txt)

### Build Warnings
$(head -20 BUILD-WARNINGS.txt)

### Missing Dependencies
[List any missing dependencies detected]

### Configuration Issues
[List any configuration problems]

## Backport Requirements
$(if [ "$BUILD_SUCCESS" = false ]; then
    echo "### CRITICAL: The following fixes need backporting to source branches:"
    echo "1. [Describe fix needed] → [Which effort branch]"
    echo "2. [Describe fix needed] → [Which effort branch]"
    echo ""
    echo "Per R321, these MUST be fixed in source branches before integration can continue."
else
    echo "No backporting required - build successful."
fi)

## Recommendation
$(if [ "$BUILD_SUCCESS" = true ]; then
    echo "PROCEED_TO_PR_PLAN - Build successful, artifacts generated"
else
    echo "IMMEDIATE_BACKPORT_REQUIRED - Build failures must be fixed in source branches"
fi)

## Next Steps
$(if [ "$BUILD_SUCCESS" = true ]; then
    echo "1. Orchestrator should proceed with PR plan creation"
    echo "2. All artifacts are ready for deployment"
else
    echo "1. Orchestrator must spawn SW Engineers to fix issues in source branches"
    echo "2. After fixes are backported, re-run integration from clean state"
    echo "3. Then re-validate build"
fi)
EOF

echo "✅ Build validation report created: BUILD-VALIDATION-REPORT.md"
```

## ⚠️ CRITICAL REQUIREMENTS

### Build Must Be Reproducible
- Same commands must produce same output
- No timestamps in artifacts (unless versioned)
- Dependencies must be locked/pinned
- Build environment must be documented

### Track Everything for Backporting
**Per R321, integration branches are READ-ONLY for code:**
- Document EVERY issue that prevents build success
- Identify WHICH effort branch has the problem
- Specify EXACTLY what fix is needed
- Orchestrator will spawn SW Engineers to apply fixes to source branches

### Zero Tolerance for Critical Errors
- Compilation errors = MUST FIX IN SOURCE
- Linking errors = MUST FIX IN SOURCE  
- Missing dependencies = MUST FIX IN SOURCE
- Runtime errors in build = MUST FIX IN SOURCE

## 🚫 FORBIDDEN ACTIONS

1. **NEVER fix code directly in integration branch** (R321 violation)
2. **NEVER skip build failures** - document them
3. **NEVER fake build success** - be honest
4. **NEVER proceed with partial builds** - all or nothing
5. **NEVER apply patches directly** - fixes go to source branches

## ✅ SUCCESS CRITERIA

Before reporting completion:
- [ ] Build process executed completely
- [ ] All build output captured in logs
- [ ] Warnings and errors analyzed
- [ ] Artifacts verified (if build succeeded)
- [ ] Comprehensive report created
- [ ] Backport requirements documented (if needed)
- [ ] All findings committed and pushed

## 🔄 STATE TRANSITIONS

Your role is to VALIDATE and REPORT. The orchestrator will decide next state based on your findings:

- If build succeeds → Orchestrator moves to PR_PLAN_CREATION
- If build fails → Orchestrator moves to IMMEDIATE_BACKPORT_REQUIRED
- If issues found → Orchestrator spawns SW Engineers to fix source branches

## 📊 VALIDATION STANDARDS

You are the technical expert. Your assessment must be:
1. **Accurate** - No false positives or negatives
2. **Complete** - Check everything thoroughly
3. **Actionable** - Clear about what needs fixing
4. **Traceable** - Document which effort has issues

## 💡 TIPS FOR SUCCESS

1. **Clean Builds**: Always start with clean environment
2. **Capture Everything**: Verbose logging helps debugging
3. **Be Specific**: "Module X fails with error Y" not "build failed"
4. **Think Backporting**: Always identify source branch for fixes
5. **Document Clearly**: Orchestrator relies on your report

## 🎓 GRADING CRITERIA

You will be evaluated on:
1. **Build Execution** (25%)
   - Correct build commands used
   - All modules built
   
2. **Issue Identification** (25%)
   - All errors caught
   - Root causes identified
   
3. **Documentation** (25%)
   - Clear, comprehensive report
   - Actionable recommendations
   
4. **Backport Tracking** (25%)
   - Source branches identified
   - Fix requirements specified

Remember: You are the ONLY agent allowed to run builds and validate artifacts. The orchestrator depends on your technical assessment to make coordination decisions!