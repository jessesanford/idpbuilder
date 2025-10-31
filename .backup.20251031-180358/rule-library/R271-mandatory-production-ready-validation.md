# Rule R271: Mandatory Production-Ready Validation (SUPREME LAW)

## Rule Statement
Every Software Factory project MUST produce production-ready, deployable software validated in an integration-testing branch WITHOUT touching main. The software MUST be immediately runnable/deployable based on its runtime type.

## Criticality Level
**SUPREME LAW** - Violation = -100% AUTOMATIC FAILURE
No project can succeed without proven deployable software.

## Core Principle
**"Software Factory produces VALIDATED SOFTWARE, not just code"**

## Detailed Requirements

### 1. Production-Ready Definition by Runtime

#### Go Projects (Binary)
```bash
# MUST succeed:
go build -o app . || go build -o app ./cmd/...
./app --version || ./app --help
go test ./...
```

#### Python Projects (Package/Service)
```bash
# MUST succeed:
pip install -e . || python setup.py develop
python -c "import package_name"
python -m package_name --help || python app.py
pytest || python -m pytest
```

#### Node.js Projects (Service/CLI)
```bash
# MUST succeed:
npm install || yarn install
npm start || npm run dev
npm test
npm run build  # if applicable
```

#### Kubernetes Projects (Manifests/Operators)
```bash
# MUST succeed:
kubectl apply --dry-run=client -f manifests/
helm lint charts/  # if using Helm
kustomize build overlays/production  # if using Kustomize
```

#### Terraform Projects (Infrastructure)
```bash
# MUST succeed:
terraform init
terraform validate
terraform plan
```

#### Docker Projects (Containers)
```bash
# MUST succeed:
docker build -t test-image .
docker run --rm test-image
# Health check must pass
```

### 2. Integration Testing Branch Protocol

**CRITICAL**: Main branch must NEVER be modified directly!

```bash
validate_production_ready() {
    # Create integration-testing branch from main's HEAD
    git checkout main
    git pull origin main
    INTEGRATE_WAVE_EFFORTS_BRANCH="integration-testing-$(date +%Y%m%d-%H%M%S)"
    git checkout -b "$INTEGRATE_WAVE_EFFORTS_BRANCH"
    
    # Merge all effort branches in order
    while read -r effort_branch; do
        echo "Integrating $effort_branch..."
        git merge --no-ff "origin/$effort_branch" || {
            echo "Resolving conflicts for $effort_branch"
            # Document resolution in MASTER-PR-PLAN.md
        }
        
        # Validate after each merge
        validate_build || return 1
    done < effort-order.txt
    
    # Final validation
    full_production_validation || return 1
    
    echo "✅ Software validated in $INTEGRATE_WAVE_EFFORTS_BRANCH"
    echo "✅ Main branch untouched"
}
```

### 3. Validation Checkpoints

- **After Each Wave**: Partial functionality test
- **After Each Phase**: Integration test
- **After Final Integration**: Complete production validation
- **Before PROJECT_DONE**: Full deployment test

### 4. External User Validation

A new user MUST be able to:
```bash
git clone <repo>
git checkout <integration-testing-branch>
# Follow README/RUNBOOK
# Software must work as documented
```

## Enforcement Mechanism

### Technical Enforcement
```python
def validate_before_success(state_data):
    """Cannot transition to PROJECT_DONE without validation"""
    
    # Check integration-testing branch exists
    if not integration_testing_branch_exists():
        return ValidationResult(
            passed=False,
            reason="No integration-testing branch found"
        )
    
    # Verify main untouched
    if main_branch_modified():
        return ValidationResult(
            passed=False,
            reason="Main branch was modified! SUPREME LAW VIOLATION!"
        )
    
    # Validate software builds/runs
    if not software_is_deployable():
        return ValidationResult(
            passed=False,
            reason="Software not production-ready"
        )
    
    return ValidationResult(passed=True)
```

## Integration with Other Rules

### Prerequisites
- R272: Integration Testing Branch created
- R273: Runtime-specific validation passed
- R274: Production checklist complete

### Enables
- R279: MASTER-PR-PLAN generation
- PROJECT_DONE state transition

## Grading Impact

**-100% AUTOMATIC FAILURE if:**
- No integration-testing branch created
- Software doesn't build/run
- Main branch was modified
- External user cannot run software

**Grade Modifiers:**
- Perfect integration (no conflicts): +10%
- All tests passing: +10%
- Complete documentation: +10%

## Common Violations to Avoid

### ❌ FORBIDDEN
```bash
# NEVER merge to main directly
git checkout main
git merge phase-integration  # VIOLATION!

# NEVER push to main
git push origin main  # VIOLATION!
```

### ✅ REQUIRED
```bash
# Always use integration-testing branch
git checkout -b integration-testing-$(date +%Y%m%d-%H%M%S)
# Test everything here
# Generate PR plan for humans
```

## State Machine Integration

This rule BLOCKS transition to PROJECT_DONE state.
The `INTEGRATE_WAVE_EFFORTS_TESTING` state must complete successfully before `PR_PLAN_CREATION`.

## Summary

R271 ensures every Software Factory project produces real, working, deployable software while respecting that main branch integration requires human review. The Factory proves everything works in an integration-testing branch and provides a plan for human-executed PRs.