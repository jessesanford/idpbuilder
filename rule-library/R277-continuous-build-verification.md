# Rule R277: Continuous Build Verification

## Rule Statement
Software Factory projects MUST verify that the software builds and passes tests continuously throughout development, not just at the end. Each wave and phase completion requires build verification.

## Criticality Level
**MANDATORY** - Prevents accumulation of broken code
Violation = -20% grade penalty per failed checkpoint

## Core Principle
**"Continuous integration prevents integration hell"**

## Build Verification Checkpoints

### 1. After Each Effort Completion
```bash
verify_effort_build() {
    local effort_branch="$1"
    echo "🔨 Verifying build for effort: $effort_branch"
    
    git checkout "$effort_branch"
    
    # Build must succeed
    if ! make build 2>/dev/null || ! go build . 2>/dev/null; then
        echo "❌ Build failed for $effort_branch"
        return 1
    fi
    
    # Tests must pass
    if ! make test 2>/dev/null || ! go test ./... 2>/dev/null; then
        echo "⚠️ Tests failing for $effort_branch"
        # Warning only at effort level
    fi
    
    echo "✅ Effort $effort_branch builds successfully"
}
```

### 2. After Each Wave Integration
```bash
verify_wave_integration_build() {
    local wave_branch="$1"
    echo "🌊 Verifying wave integration: $wave_branch"
    
    # Create temp integration branch
    git checkout -b temp-wave-verify
    
    # Merge all wave efforts
    for effort in $(get_wave_efforts "$wave_branch"); do
        git merge --no-ff "$effort" || {
            echo "❌ Merge conflict in wave integration"
            return 1
        }
    done
    
    # Build must succeed
    if ! make build || ! go build .; then
        echo "❌ Wave integration build failed"
        return 1
    fi
    
    # Tests must pass
    if ! make test || ! go test ./...; then
        echo "❌ Wave integration tests failed"
        return 1
    fi
    
    # Check for performance regression
    if command -v make benchmark &> /dev/null; then
        make benchmark || echo "⚠️ Performance regression detected"
    fi
    
    echo "✅ Wave $wave_branch integration verified"
}
```

### 3. After Each Phase Integration
```bash
verify_phase_integration_build() {
    local phase_num="$1"
    echo "📦 Verifying phase $phase_num integration"
    
    # Create temp phase integration
    git checkout -b temp-phase-verify
    
    # Merge all phase waves
    for wave in $(get_phase_waves "$phase_num"); do
        git merge --no-ff "$wave" || {
            echo "❌ Merge conflict in phase integration"
            return 1
        }
    done
    
    # Full build validation
    if ! make all || ! go build -v ./...; then
        echo "❌ Phase build failed"
        return 1
    fi
    
    # All tests must pass
    if ! make test-all || ! go test -v ./...; then
        echo "❌ Phase tests failed"
        return 1
    fi
    
    # Integration tests
    if ! make test-integration; then
        echo "❌ Integration tests failed"
        return 1
    fi
    
    # Check binary/package works
    if [ -f Makefile ] && grep -q "^run:" Makefile; then
        timeout 10 make run || echo "⚠️ Runtime check failed"
    fi
    
    echo "✅ Phase $phase_num fully integrated and verified"
}
```

## Continuous Verification Schedule

```yaml
verification_schedule:
  effort_complete:
    - build: required
    - unit_tests: required
    - lint: recommended
    
  wave_complete:
    - build: required
    - all_tests: required
    - integration_test: required
    - performance_check: recommended
    
  phase_complete:
    - full_build: required
    - all_tests: required
    - integration_tests: required
    - deployment_test: required
    - security_scan: recommended
    
  project_complete:
    - production_build: required
    - full_test_suite: required
    - deployment_verification: required
    - load_testing: recommended
    - security_audit: required
```

## Build Health Metrics

```bash
track_build_health() {
    cat > build-health.md << EOF
# Build Health Report
Generated: $(date)

## Build Success Rate
| Level | Total | Passed | Failed | Success Rate |
|-------|-------|--------|--------|--------------|
| Efforts | 45 | 43 | 2 | 95.6% |
| Waves | 12 | 11 | 1 | 91.7% |
| Phases | 3 | 3 | 0 | 100% |

## Test Coverage Trend
| Phase | Coverage | Change |
|-------|----------|--------|
| Phase 1 | 65% | - |
| Phase 2 | 72% | +7% |
| Phase 3 | 78% | +6% |

## Build Time Trend
| Checkpoint | Duration | Status |
|------------|----------|--------|
| Effort avg | 45s | ✅ |
| Wave avg | 2m 30s | ✅ |
| Phase avg | 5m 15s | ✅ |
| Full build | 8m 42s | ✅ |

## Failed Build Analysis
| Branch | Failure Type | Resolution |
|--------|--------------|------------|
| effort-23 | Import cycle | Fixed in next commit |
| wave-7 | Test timeout | Increased timeout |
EOF
}
```

## Early Warning System

```bash
continuous_build_monitor() {
    while true; do
        # Check latest commit
        LATEST_COMMIT=$(git rev-parse HEAD)
        
        # Try to build
        if ! make build 2>/dev/null; then
            alert "🚨 Build broken at commit $LATEST_COMMIT"
            
            # Bisect to find breaking commit
            git bisect start
            git bisect bad HEAD
            git bisect good HEAD~10
            git bisect run make build
            
            BREAKING_COMMIT=$(git bisect view --oneline | head -1)
            alert "💔 Breaking commit: $BREAKING_COMMIT"
            
            git bisect reset
        fi
        
        sleep 300  # Check every 5 minutes
    done
}
```

## Build Verification Matrix

| Stage | Build | Unit Tests | Integration | Deploy Test | Required |
|-------|-------|------------|-------------|-------------|----------|
| Effort Done | ✅ | ✅ | - | - | Yes |
| Wave Done | ✅ | ✅ | ✅ | - | Yes |
| Phase Done | ✅ | ✅ | ✅ | ✅ | Yes |
| Project Done | ✅ | ✅ | ✅ | ✅ | Yes |

## Failure Recovery Protocol

### Build Failure at Effort Level
1. Fix in effort branch
2. Rerun verification
3. Continue if passes

### Build Failure at Wave Level
1. Identify conflicting efforts
2. Fix integration issues
3. Rerun wave verification
4. May need to split efforts

### Build Failure at Phase Level
1. Major integration issue
2. May need architectural review
3. Consider phase restructuring
4. Must fix before next phase

## Integration with Other Rules

### Prerequisites
- R273: Runtime validation
- Clean code practices

### Enables
- R271: Production readiness
- R279: PR plan generation

### Related Rules
- R007: Size limits (smaller = easier to verify)
- R031: Code review (catch issues early)

## Common Issues and Prevention

### Issue: "Works on my machine"
**Prevention**: Always verify in clean environment

### Issue: "Tests pass individually but fail together"
**Prevention**: Run full test suite at each checkpoint

### Issue: "Build time increases exponentially"
**Prevention**: Monitor build time trends, optimize regularly

### Issue: "Integration reveals fundamental conflicts"
**Prevention**: Continuous verification catches conflicts early

## Grading Impact

- No effort-level verification: -5%
- No wave-level verification: -10%
- No phase-level verification: -15%
- Build broken for >1 day: -20%
- No build health tracking: -5%

## Automation Support

```bash
# Add to CI/CD pipeline
cat > .github/workflows/continuous-build.yml << 'EOF'
name: Continuous Build Verification

on:
  push:
    branches: ['phase*/wave*/effort*']
  pull_request:
    branches: ['main', 'integration-testing-*']

jobs:
  verify-build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Detect Runtime
        run: |
          source utilities/detect-runtime.sh
          echo "RUNTIME=$(detect_project_runtime)" >> $GITHUB_ENV
      
      - name: Setup Environment
        run: |
          case $RUNTIME in
            go) setup-go@v4 ;;
            node) setup-node@v3 ;;
            python) setup-python@v4 ;;
          esac
      
      - name: Build
        run: make build || go build . || npm run build
      
      - name: Test
        run: make test || go test ./... || npm test
      
      - name: Integration Test
        if: contains(github.ref, 'wave') || contains(github.ref, 'phase')
        run: make test-integration
      
      - name: Report Status
        if: always()
        run: |
          echo "Build Status: ${{ job.status }}" >> build-health.log
          echo "Branch: ${{ github.ref }}" >> build-health.log
          echo "Time: $(date)" >> build-health.log
EOF
```

## Summary

R277 ensures that build health is maintained continuously throughout development, not just verified at the end. This prevents the accumulation of technical debt and integration problems that can derail projects at the finish line.