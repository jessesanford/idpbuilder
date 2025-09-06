# Rule R274: Production Readiness Checklist

## Rule Statement
Every Software Factory project MUST complete a comprehensive production readiness checklist before declaring SUCCESS. This ensures the software is not just functional but ready for real-world deployment.

## Criticality Level
**BLOCKING** - Cannot proceed to SUCCESS without checklist completion
Violation = -40% grade penalty

## Core Principle
**"Production-ready means ready for real users, not just demos"**

## The Production Readiness Checklist

### 1. Build & Compilation ✓
```yaml
build_requirements:
  - clean_checkout_builds: "Fresh clone must build successfully"
  - dependency_resolution: "All dependencies resolve correctly"
  - build_reproducibility: "Build produces consistent output"
  - build_documentation: "Clear build instructions exist"
  - build_time: "Reasonable build time (<10 minutes)"
```

### 2. Testing Coverage ✓
```yaml
test_requirements:
  - unit_tests_pass: "All unit tests passing"
  - integration_tests_pass: "Integration tests passing"
  - test_coverage: "Reasonable coverage (>60% for new code)"
  - test_documentation: "How to run tests documented"
  - ci_ready: "Tests can run in CI environment"
```

### 3. Documentation ✓
```yaml
documentation_requirements:
  - readme_complete: "README with purpose, setup, usage"
  - api_documented: "API/Interface documentation"
  - deployment_guide: "How to deploy to production"
  - troubleshooting: "Common issues and solutions"
  - changelog: "Version history or changelog"
```

### 4. Deployment Readiness ✓
```yaml
deployment_requirements:
  - configuration_management: "Externalized configuration"
  - secrets_management: "No hardcoded secrets"
  - health_checks: "Health/liveness endpoints"
  - logging_configured: "Structured logging in place"
  - monitoring_ready: "Metrics/monitoring hooks"
```

### 5. Security & Compliance ✓
```yaml
security_requirements:
  - dependency_scanning: "No critical vulnerabilities"
  - secret_scanning: "No exposed credentials"
  - security_headers: "Proper security headers (if web)"
  - authentication: "Auth mechanism in place (if needed)"
  - authorization: "Proper access controls"
```

### 6. Performance ✓
```yaml
performance_requirements:
  - load_tested: "Basic load testing completed"
  - resource_limits: "Memory/CPU limits defined"
  - response_times: "Acceptable response times"
  - scalability_plan: "Can scale horizontally/vertically"
  - optimization_done: "No obvious performance issues"
```

### 7. Operations ✓
```yaml
operations_requirements:
  - runbook_exists: "Operational runbook created"
  - backup_strategy: "Data backup plan (if applicable)"
  - disaster_recovery: "DR plan documented"
  - rollback_procedure: "How to rollback deployment"
  - support_contacts: "Who to contact for issues"
```

## Validation Script

```bash
#!/bin/bash
# production-readiness-check.sh

CHECKLIST_PASS=0
CHECKLIST_FAIL=0
CHECKLIST_WARN=0

check_item() {
    local category="$1"
    local item="$2"
    local command="$3"
    local required="$4"  # true/false
    
    echo -n "  [$category] $item: "
    
    if eval "$command" > /dev/null 2>&1; then
        echo "✅ PASS"
        ((CHECKLIST_PASS++))
        return 0
    else
        if [ "$required" = "true" ]; then
            echo "❌ FAIL (REQUIRED)"
            ((CHECKLIST_FAIL++))
            return 1
        else
            echo "⚠️  WARN (Recommended)"
            ((CHECKLIST_WARN++))
            return 0
        fi
    fi
}

echo "======================================"
echo "🔍 PRODUCTION READINESS CHECKLIST"
echo "======================================"

# Build Checks
echo "📦 Build & Compilation:"
check_item "BUILD" "Project builds" "make build || go build . || npm run build" "true"
check_item "BUILD" "Dependencies resolve" "go mod download || npm ci || pip install -r requirements.txt" "true"
check_item "BUILD" "Build instructions exist" "test -f README.md && grep -qi 'build' README.md" "true"

# Test Checks
echo "\n🧪 Testing:"
check_item "TEST" "Tests exist" "find . -name '*_test.go' -o -name '*.test.js' -o -name 'test_*.py' | head -1" "true"
check_item "TEST" "Tests pass" "make test || go test ./... || npm test" "true"
check_item "TEST" "Test documentation" "grep -qi 'test' README.md" "false"

# Documentation Checks
echo "\n📚 Documentation:"
check_item "DOCS" "README exists" "test -f README.md" "true"
check_item "DOCS" "Runbook exists" "test -f RUNBOOK.md -o -f docs/runbook.md" "true"
check_item "DOCS" "API docs" "test -f API.md -o -f docs/api.md -o -f swagger.json" "false"

# Deployment Checks
echo "\n🚀 Deployment:"
check_item "DEPLOY" "No hardcoded secrets" "! grep -r 'password.*=.*[\"\']' --include='*.go' --include='*.js' --include='*.py' ." "true"
check_item "DEPLOY" "Config externalized" "test -f .env.example -o -f config.yaml.example" "true"
check_item "DEPLOY" "Dockerfile exists" "test -f Dockerfile" "false"

# Security Checks
echo "\n🔒 Security:"
check_item "SEC" "No vulnerable deps" "go list -m all | nancy || npm audit --audit-level=high" "false"
check_item "SEC" "Secret scanning clean" "! git grep -E '(api_key|secret|token|password)\s*=\s*[\"\'][^\"\']+[\"\']'" "true"

# Operations Checks
echo "\n⚙️ Operations:"
check_item "OPS" "Health check endpoint" "grep -r 'health' --include='*.go' --include='*.js' --include='*.py' ." "false"
check_item "OPS" "Logging configured" "grep -r 'log\.' --include='*.go' --include='*.js' --include='*.py' . | head -1" "true"
check_item "OPS" "Metrics/monitoring" "grep -r 'metrics\|prometheus' --include='*.go' --include='*.js' ." "false"

echo "\n======================================"
echo "📊 RESULTS:"
echo "  ✅ Passed: $CHECKLIST_PASS"
echo "  ❌ Failed: $CHECKLIST_FAIL"
echo "  ⚠️  Warnings: $CHECKLIST_WARN"

if [ $CHECKLIST_FAIL -gt 0 ]; then
    echo "\n🚫 PRODUCTION READINESS: FAILED"
    echo "Fix the failed required items before proceeding."
    exit 1
else
    echo "\n✅ PRODUCTION READINESS: PASSED"
    if [ $CHECKLIST_WARN -gt 0 ]; then
        echo "Consider addressing warnings for production deployment."
    fi
    exit 0
fi
```

## Checklist by Project Type

### Web Services
- Health/readiness endpoints
- Rate limiting configured
- CORS properly set
- TLS/HTTPS ready
- Session management

### CLI Tools
- Help text complete
- Version flag works
- Exit codes documented
- Shell completion available
- Man page or help docs

### Libraries
- API stability guaranteed
- Semantic versioning
- Breaking changes documented
- Examples provided
- Benchmarks available

### Infrastructure (Terraform/K8s)
- State management planned
- Resource tagging
- Cost estimation done
- Destroy plan tested
- RBAC configured

## Integration with Other Rules

### Prerequisites
- R273: Runtime validation passed
- All tests passing

### Enables
- R275: Deployment verification
- R279: PR plan generation

## Grading Impact

**Required Items (MUST PASS):**
- Builds successfully: -20% if fails
- Tests pass: -15% if fails
- Documentation exists: -10% if missing
- No hardcoded secrets: -15% if found

**Recommended Items (SHOULD PASS):**
- Each warning: -2% (max -10%)

## Common Issues and Solutions

### Issue: "Tests exist but don't run"
**Solution**: Ensure test command in package.json, Makefile, or documented

### Issue: "Builds locally but not in clean environment"
**Solution**: Test with fresh clone, document all prerequisites

### Issue: "No deployment documentation"
**Solution**: Create RUNBOOK.md with step-by-step deployment

### Issue: "Hardcoded configuration"
**Solution**: Externalize to env vars, config files, or flags

## Summary

R274 ensures that Software Factory output isn't just "code that works" but "software ready for production". This comprehensive checklist prevents common deployment failures and ensures operational readiness.