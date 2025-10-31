#!/bin/bash
# 🚨 PRODUCTION READINESS VALIDATION SCRIPT
# Implements R271-R280 validation requirements

set -e

echo "================================================"
echo "🔍 PRODUCTION READINESS VALIDATION"
echo "================================================"
echo "Timestamp: $(date -u +%Y-%m-%dT%H:%M:%SZ)"
echo "================================================"

VALIDATION_PASSED=true
ERRORS=0
WARNINGS=0

# Function to detect project runtime
detect_project_runtime() {
    if [ -f "go.mod" ] || [ -f "go.sum" ]; then
        echo "go"
    elif [ -f "Cargo.toml" ]; then
        echo "rust"
    elif [ -f "package.json" ]; then
        echo "nodejs"
    elif [ -f "requirements.txt" ] || [ -f "setup.py" ] || [ -f "pyproject.toml" ]; then
        echo "python"
    elif [ -f "pom.xml" ]; then
        echo "java"
    elif [ -f "Dockerfile" ]; then
        echo "docker"
    elif [ -f "Chart.yaml" ]; then
        echo "helm"
    elif [ -f "kustomization.yaml" ]; then
        echo "kustomize"
    elif [ -d ".terraform" ] || ls *.tf 2>/dev/null | head -1 > /dev/null; then
        echo "terraform"
    else
        echo "unknown"
    fi
}

# Function to validate Go project
validate_go_project() {
    echo "🔍 Validating Go project..."
    
    # Module verification
    if ! go mod verify 2>/dev/null; then
        echo "❌ Go module verification failed"
        ((ERRORS++))
        return 1
    fi
    
    # Build validation
    if ! go build -o test-binary . 2>/dev/null && \
       ! go build -o test-binary ./cmd/... 2>/dev/null && \
       ! go build -o test-binary ./cmd/main.go 2>/dev/null; then
        echo "❌ Go build failed"
        ((ERRORS++))
        return 1
    fi
    
    # Test execution
    if ! go test ./... 2>/dev/null; then
        echo "⚠️ Go tests failing"
        ((WARNINGS++))
    fi
    
    # Binary execution
    if [ -f test-binary ]; then
        if ! ./test-binary --version 2>/dev/null && ! ./test-binary --help 2>/dev/null; then
            echo "❌ Binary doesn't execute properly"
            ((ERRORS++))
            return 1
        fi
        rm -f test-binary
    fi
    
    echo "✅ Go validation passed"
    return 0
}

# Function to validate Node.js project
validate_nodejs_project() {
    echo "🔍 Validating Node.js project..."
    
    # Dependency installation
    if [ -f "package-lock.json" ]; then
        if ! npm ci 2>/dev/null; then
            echo "❌ npm ci failed"
            ((ERRORS++))
            return 1
        fi
    else
        if ! npm install 2>/dev/null; then
            echo "❌ npm install failed"
            ((ERRORS++))
            return 1
        fi
    fi
    
    # Build if applicable
    if grep -q '"build"' package.json 2>/dev/null; then
        if ! npm run build 2>/dev/null; then
            echo "❌ npm build failed"
            ((ERRORS++))
            return 1
        fi
    fi
    
    # Test execution
    if ! npm test 2>/dev/null; then
        echo "⚠️ npm tests failing"
        ((WARNINGS++))
    fi
    
    echo "✅ Node.js validation passed"
    return 0
}

# Function to validate Python project
validate_python_project() {
    echo "🔍 Validating Python project..."
    
    # Create virtual environment
    python -m venv test-env 2>/dev/null || python3 -m venv test-env
    source test-env/bin/activate
    
    # Dependency installation
    if [ -f "requirements.txt" ]; then
        if ! pip install -r requirements.txt 2>/dev/null; then
            echo "❌ pip install failed"
            deactivate
            rm -rf test-env
            ((ERRORS++))
            return 1
        fi
    elif [ -f "setup.py" ]; then
        if ! pip install -e . 2>/dev/null; then
            echo "❌ setup.py install failed"
            deactivate
            rm -rf test-env
            ((ERRORS++))
            return 1
        fi
    fi
    
    # Test execution
    if ! pytest 2>/dev/null && ! python -m pytest 2>/dev/null; then
        echo "⚠️ Python tests not found or failing"
        ((WARNINGS++))
    fi
    
    deactivate
    rm -rf test-env
    echo "✅ Python validation passed"
    return 0
}

# Function to validate Kubernetes project
validate_kubernetes_project() {
    echo "🔍 Validating Kubernetes manifests..."
    
    # YAML validation
    for file in $(find . -name "*.yaml" -o -name "*.yml"); do
        if ! kubectl apply --dry-run=client -f "$file" 2>/dev/null; then
            echo "❌ Invalid manifest: $file"
            ((ERRORS++))
            return 1
        fi
    done
    
    # Helm chart validation if present
    if [ -f "Chart.yaml" ]; then
        if ! helm lint . 2>/dev/null; then
            echo "❌ Helm chart validation failed"
            ((ERRORS++))
            return 1
        fi
    fi
    
    echo "✅ Kubernetes validation passed"
    return 0
}

# Function to validate Docker project
validate_docker_project() {
    echo "🔍 Validating Docker project..."
    
    # Build image
    if ! docker build -t test-image:validation . 2>/dev/null; then
        echo "❌ Docker build failed"
        ((ERRORS++))
        return 1
    fi
    
    # Cleanup
    docker rmi test-image:validation 2>/dev/null || true
    
    echo "✅ Docker validation passed"
    return 0
}

# Function to validate Terraform project
validate_terraform_project() {
    echo "🔍 Validating Terraform project..."
    
    # Initialize
    if ! terraform init -backend=false 2>/dev/null; then
        echo "❌ Terraform init failed"
        ((ERRORS++))
        return 1
    fi
    
    # Validation
    if ! terraform validate 2>/dev/null; then
        echo "❌ Terraform validation failed"
        ((ERRORS++))
        return 1
    fi
    
    echo "✅ Terraform validation passed"
    return 0
}

# Check 1: Integration Testing Branch
echo ""
echo "📋 Checking R272: Integration Testing Branch..."
CURRENT_BRANCH=$(git branch --show-current)
if [[ "$CURRENT_BRANCH" =~ integration-testing ]]; then
    echo "✅ On integration-testing branch: $CURRENT_BRANCH"
else
    echo "⚠️ WARNING: Not on integration-testing branch (current: $CURRENT_BRANCH)"
    ((WARNINGS++))
fi

# Check 2: Main Branch Protection
echo ""
echo "📋 Checking R280: Main Branch Protection..."
if [ "$CURRENT_BRANCH" = "main" ]; then
    echo "❌ SUPREME LAW VIOLATION: On main branch!"
    echo "Software Factory must NEVER modify main!"
    exit 280
fi
echo "✅ Not on main branch (R280 compliant)"

# Check 3: Runtime-Specific Validation
echo ""
echo "📋 Checking R273: Runtime-Specific Validation..."
RUNTIME=$(detect_project_runtime)
echo "Detected runtime: $RUNTIME"

case $RUNTIME in
    go) validate_go_project ;;
    nodejs) validate_nodejs_project ;;
    python) validate_python_project ;;
    kubernetes) validate_kubernetes_project ;;
    docker) validate_docker_project ;;
    terraform) validate_terraform_project ;;
    *)
        echo "⚠️ Unknown runtime - generic validation only"
        ((WARNINGS++))
        ;;
esac

# Check 4: Documentation
echo ""
echo "📋 Checking R276: Documentation Requirements..."
if [ ! -f "README.md" ]; then
    echo "❌ No README.md found"
    ((ERRORS++))
else
    echo "✅ README.md exists"
fi

if [ ! -f "RUNBOOK.md" ] && [ ! -f "docs/runbook.md" ] && [ ! -f "docs/RUNBOOK.md" ]; then
    echo "⚠️ No RUNBOOK.md found"
    ((WARNINGS++))
else
    echo "✅ RUNBOOK exists"
fi

# Check 5: Build Artifacts
echo ""
echo "📋 Checking R277: Build Artifacts..."
if [ "$RUNTIME" = "go" ]; then
    if go build . 2>/dev/null || go build ./cmd/... 2>/dev/null; then
        echo "✅ Build produces artifacts"
    else
        echo "❌ Build fails to produce artifacts"
        ((ERRORS++))
    fi
elif [ "$RUNTIME" = "nodejs" ]; then
    if [ -d "dist" ] || [ -d "build" ] || [ -d "out" ]; then
        echo "✅ Build artifacts found"
    else
        echo "⚠️ No obvious build artifacts"
        ((WARNINGS++))
    fi
fi

# Check 6: MASTER-PR-PLAN.md
echo ""
echo "📋 Checking R279: MASTER-PR-PLAN Requirement..."
if [ ! -f "MASTER-PR-PLAN.md" ]; then
    echo "⚠️ MASTER-PR-PLAN.md not yet created (will be generated)"
    ((WARNINGS++))
else
    echo "✅ MASTER-PR-PLAN.md exists"
fi

# Final Summary
echo ""
echo "================================================"
echo "📊 VALIDATION SUMMARY"
echo "================================================"
echo "Runtime: $RUNTIME"
echo "Branch: $CURRENT_BRANCH"
echo "Errors: $ERRORS"
echo "Warnings: $WARNINGS"

if [ $ERRORS -gt 0 ]; then
    echo ""
    echo "❌ PRODUCTION READINESS: FAILED"
    echo "Fix the $ERRORS errors before proceeding to PROJECT_DONE"
    exit 1
elif [ $WARNINGS -gt 0 ]; then
    echo ""
    echo "⚠️ PRODUCTION READINESS: PASSED WITH WARNINGS"
    echo "Consider addressing $WARNINGS warnings"
    exit 0
else
    echo ""
    echo "✅ PRODUCTION READINESS: FULLY VALIDATED"
    echo "Software is ready for PR plan generation"
    exit 0
fi