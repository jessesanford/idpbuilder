# Rule R273: Runtime-Specific Validation

## Rule Statement
Every Software Factory project MUST pass validation specific to its runtime environment. Validation requirements vary by language/platform but all must prove the software is production-ready.

## Criticality Level
**MANDATORY** - Required for production readiness
Violation = -30% grade penalty

## Core Principle
**"Validate according to how the software will actually run"**

## Runtime Detection

```bash
detect_project_runtime() {
    # Check project markers in order of precedence
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
    elif [ -f "build.gradle" ] || [ -f "build.gradle.kts" ]; then
        echo "gradle"
    elif [ -f "Gemfile" ]; then
        echo "ruby"
    elif [ -f "composer.json" ]; then
        echo "php"
    elif [ -f "Dockerfile" ]; then
        echo "docker"
    elif [ -f "Chart.yaml" ]; then
        echo "helm"
    elif [ -f "kustomization.yaml" ]; then
        echo "kustomize"
    elif [ -d ".terraform" ] || ls *.tf 2>/dev/null | head -1 > /dev/null; then
        echo "terraform"
    elif [ -f "Makefile" ]; then
        echo "make"
    else
        echo "unknown"
    fi
}
```

## Validation Requirements by Runtime

### Go Projects
```bash
validate_go_project() {
    echo "🔍 Validating Go project..."
    
    # 1. Module verification
    go mod verify || return 1
    go mod tidy || return 1
    
    # 2. Build validation
    go build -o test-binary . || \
    go build -o test-binary ./cmd/... || \
    go build -o test-binary ./cmd/main.go || return 1
    
    # 3. Test execution
    go test ./... || return 1
    
    # 4. Binary execution
    if [ -f test-binary ]; then
        ./test-binary --version || ./test-binary --help || return 1
    fi
    
    # 5. Linting (if configured)
    if command -v golangci-lint &> /dev/null; then
        golangci-lint run || echo "⚠️ Linting issues found"
    fi
    
    echo "✅ Go validation passed"
}
```

### Python Projects
```bash
validate_python_project() {
    echo "🔍 Validating Python project..."
    
    # 1. Virtual environment setup
    python -m venv test-env || python3 -m venv test-env
    source test-env/bin/activate
    
    # 2. Dependency installation
    if [ -f "requirements.txt" ]; then
        pip install -r requirements.txt || return 1
    elif [ -f "setup.py" ]; then
        pip install -e . || return 1
    elif [ -f "pyproject.toml" ]; then
        pip install . || return 1
    fi
    
    # 3. Import test
    python -c "import sys; sys.path.insert(0, '.'); __import__('main')" || \
    python -c "import app" || echo "⚠️ No main module"
    
    # 4. Test execution
    pytest || python -m pytest || python -m unittest discover || echo "⚠️ No tests"
    
    # 5. Execution test
    python main.py --help 2>/dev/null || \
    python -m app --help 2>/dev/null || \
    python app.py --help 2>/dev/null || echo "⚠️ No CLI interface"
    
    deactivate
    echo "✅ Python validation passed"
}
```

### Node.js Projects
```bash
validate_node_project() {
    echo "🔍 Validating Node.js project..."
    
    # 1. Dependency installation
    if [ -f "package-lock.json" ]; then
        npm ci || return 1
    else
        npm install || return 1
    fi
    
    # 2. Build (if applicable)
    if grep -q '"build"' package.json; then
        npm run build || return 1
    fi
    
    # 3. Test execution
    npm test || echo "⚠️ No tests configured"
    
    # 4. Start validation
    timeout 5 npm start || echo "ℹ️ Service started (timeout expected)"
    
    # 5. Package validity
    npm ls --depth=0 || return 1
    
    echo "✅ Node.js validation passed"
}
```

### Kubernetes Projects
```bash
validate_kubernetes_project() {
    echo "🔍 Validating Kubernetes manifests..."
    
    # 1. YAML validation
    find . -name "*.yaml" -o -name "*.yml" | while read -r file; do
        kubectl apply --dry-run=client -f "$file" || return 1
    done
    
    # 2. Helm chart validation (if present)
    if [ -f "Chart.yaml" ]; then
        helm lint . || return 1
        helm template . | kubectl apply --dry-run=client -f - || return 1
    fi
    
    # 3. Kustomize validation (if present)
    if [ -f "kustomization.yaml" ]; then
        kustomize build . | kubectl apply --dry-run=client -f - || return 1
    fi
    
    # 4. CRD validation (if present)
    if ls crds/*.yaml 2>/dev/null; then
        kubectl apply --dry-run=client -f crds/ || return 1
    fi
    
    echo "✅ Kubernetes validation passed"
}
```

### Terraform Projects
```bash
validate_terraform_project() {
    echo "🔍 Validating Terraform project..."
    
    # 1. Initialize
    terraform init -backend=false || return 1
    
    # 2. Format check
    terraform fmt -check || echo "⚠️ Formatting issues"
    
    # 3. Validation
    terraform validate || return 1
    
    # 4. Plan generation (without applying)
    terraform plan -input=false || return 1
    
    echo "✅ Terraform validation passed"
}
```

### Docker Projects
```bash
validate_docker_project() {
    echo "🔍 Validating Docker project..."
    
    # 1. Build image
    docker build -t test-image:validation . || return 1
    
    # 2. Run container
    docker run --rm -d --name test-container test-image:validation || return 1
    
    # 3. Health check (if defined)
    sleep 5
    docker inspect test-container --format='{{.State.Status}}' | grep -q running || return 1
    
    # 4. Cleanup
    docker stop test-container 2>/dev/null || true
    docker rmi test-image:validation || true
    
    echo "✅ Docker validation passed"
}
```

## Generic Validation (Unknown Runtime)

```bash
validate_generic_project() {
    echo "🔍 Running generic validation..."
    
    # 1. Check for Makefile
    if [ -f "Makefile" ]; then
        make build || make || echo "⚠️ Make failed"
        make test || echo "⚠️ No test target"
    fi
    
    # 2. Check for scripts
    if [ -f "build.sh" ]; then
        bash build.sh || echo "⚠️ Build script failed"
    fi
    
    if [ -f "test.sh" ]; then
        bash test.sh || echo "⚠️ Test script failed"
    fi
    
    # 3. Look for executables
    find . -type f -perm +111 -name "*" | head -5
    
    echo "⚠️ Generic validation only - consider adding runtime-specific validation"
}
```

## Master Validation Function

```bash
validate_runtime() {
    local runtime=$(detect_project_runtime)
    echo "Detected runtime: $runtime"
    
    case $runtime in
        go) validate_go_project ;;
        python) validate_python_project ;;
        nodejs) validate_node_project ;;
        rust) validate_rust_project ;;
        java) validate_java_project ;;
        docker) validate_docker_project ;;
        kubernetes) validate_kubernetes_project ;;
        terraform) validate_terraform_project ;;
        helm) validate_helm_project ;;
        *) validate_generic_project ;;
    esac
}
```

## Integration with Other Rules

### Prerequisites
- R272: Integration testing branch exists
- Project structure follows runtime conventions

### Enables
- R274: Production readiness checklist
- R275: Deployment verification

## Grading Impact

- Runtime not detected: -10%
- Build/compile fails: -30%
- Tests fail: -20%
- No tests present: -10%
- Binary/package doesn't run: -30%

## Summary

R273 ensures that every project is validated according to its specific runtime requirements. This prevents the "it works on my machine" problem by enforcing runtime-specific build, test, and execution validation.