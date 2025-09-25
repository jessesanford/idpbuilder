# Rule R275: Deployment Verification

## Rule Statement
Every Software Factory project MUST verify that the software can actually be deployed and run in a production-like environment. Building is not enough - the software must be deployable.

## Criticality Level
**MANDATORY** - Required for SUCCESS
Violation = -30% grade penalty

## Core Principle
**"If it can't be deployed, it's not done"**

## Deployment Verification by Type

### 1. Web Services/APIs
```bash
verify_web_service_deployment() {
    echo "🚀 Verifying web service deployment..."
    
    # Start the service
    PORT=8080
    ./app &
    APP_PID=$!
    sleep 5
    
    # Verify it's running
    if ! ps -p $APP_PID > /dev/null; then
        echo "❌ Service failed to start"
        return 1
    fi
    
    # Test health endpoint
    curl -f "http://localhost:$PORT/health" || \
    curl -f "http://localhost:$PORT/healthz" || \
    curl -f "http://localhost:$PORT/" || {
        echo "❌ Service not responding"
        kill $APP_PID 2>/dev/null
        return 1
    }
    
    # Cleanup
    kill $APP_PID
    echo "✅ Web service deployment verified"
}
```

### 2. CLI Applications
```bash
verify_cli_deployment() {
    echo "🚀 Verifying CLI deployment..."
    
    # Build/Install
    make install || go install . || npm install -g . || {
        echo "❌ Installation failed"
        return 1
    }
    
    # Verify command available
    COMMAND_NAME=$(basename $(pwd))
    which $COMMAND_NAME || {
        echo "❌ Command not in PATH"
        return 1
    }
    
    # Test basic commands
    $COMMAND_NAME --version || $COMMAND_NAME version || {
        echo "❌ Version command failed"
        return 1
    }
    
    $COMMAND_NAME --help || $COMMAND_NAME help || {
        echo "❌ Help command failed"
        return 1
    }
    
    echo "✅ CLI deployment verified"
}
```

### 3. Kubernetes Applications
```bash
verify_kubernetes_deployment() {
    echo "🚀 Verifying Kubernetes deployment..."
    
    # Create test namespace
    TEST_NS="test-deploy-$(date +%s)"
    kubectl create namespace $TEST_NS
    
    # Deploy application
    kubectl apply -f manifests/ -n $TEST_NS || \
    helm install test-release . -n $TEST_NS || {
        echo "❌ Deployment failed"
        kubectl delete namespace $TEST_NS
        return 1
    }
    
    # Wait for rollout
    kubectl rollout status deployment -n $TEST_NS --timeout=5m || {
        echo "❌ Rollout failed"
        kubectl delete namespace $TEST_NS
        return 1
    }
    
    # Verify pods running
    kubectl get pods -n $TEST_NS | grep -q Running || {
        echo "❌ No pods running"
        kubectl delete namespace $TEST_NS
        return 1
    }
    
    # Cleanup
    kubectl delete namespace $TEST_NS
    echo "✅ Kubernetes deployment verified"
}
```

### 4. Docker Containers
```bash
verify_docker_deployment() {
    echo "🚀 Verifying Docker deployment..."
    
    # Build image
    docker build -t deploy-test:latest . || {
        echo "❌ Docker build failed"
        return 1
    }
    
    # Run container
    docker run -d --name deploy-test -p 8080:8080 deploy-test:latest || {
        echo "❌ Container run failed"
        return 1
    }
    
    sleep 5
    
    # Check if running
    docker ps | grep -q deploy-test || {
        echo "❌ Container not running"
        docker rm -f deploy-test 2>/dev/null
        return 1
    }
    
    # Test connectivity (if applicable)
    curl -f http://localhost:8080/ || echo "ℹ️ No web interface"
    
    # Check logs for errors
    docker logs deploy-test 2>&1 | grep -i error && {
        echo "⚠️ Errors found in logs"
    }
    
    # Cleanup
    docker rm -f deploy-test
    docker rmi deploy-test:latest
    echo "✅ Docker deployment verified"
}
```

### 5. Terraform Infrastructure
```bash
verify_terraform_deployment() {
    echo "🚀 Verifying Terraform deployment..."
    
    # Initialize
    terraform init || {
        echo "❌ Terraform init failed"
        return 1
    }
    
    # Plan
    terraform plan -out=tfplan || {
        echo "❌ Terraform plan failed"
        return 1
    }
    
    # Validate plan creates resources
    terraform show -json tfplan | jq '.resource_changes | length' | {
        read count
        if [ "$count" -eq 0 ]; then
            echo "⚠️ No resources to create"
        else
            echo "✅ Plan creates $count resources"
        fi
    }
    
    # Would apply in real deployment
    echo "ℹ️ Skipping actual apply (would create real resources)"
    
    echo "✅ Terraform deployment verified"
}
```

### 6. Python Packages
```bash
verify_python_deployment() {
    echo "🚀 Verifying Python package deployment..."
    
    # Create virtual env
    python -m venv test-deploy
    source test-deploy/bin/activate
    
    # Install package
    pip install . || {
        echo "❌ Package installation failed"
        deactivate
        return 1
    }
    
    # Test import
    python -c "import $(basename $(pwd))" || {
        echo "❌ Package import failed"
        deactivate
        return 1
    }
    
    # Test CLI if exists
    if grep -q "console_scripts" setup.py 2>/dev/null; then
        SCRIPT_NAME=$(grep console_scripts setup.py | grep -oP '\w+(?=\s*=)')
        $SCRIPT_NAME --help || echo "⚠️ CLI not working"
    fi
    
    deactivate
    rm -rf test-deploy
    echo "✅ Python deployment verified"
}
```

## Master Deployment Verification

```bash
verify_deployment() {
    local project_type=$(detect_project_type)
    
    echo "📦 Detected project type: $project_type"
    echo "Starting deployment verification..."
    
    case $project_type in
        web_service)
            verify_web_service_deployment ;;
        cli)
            verify_cli_deployment ;;
        kubernetes)
            verify_kubernetes_deployment ;;
        docker)
            verify_docker_deployment ;;
        terraform)
            verify_terraform_deployment ;;
        python_package)
            verify_python_deployment ;;
        *)
            echo "⚠️ Unknown project type - skipping deployment verification"
            return 0 ;;
    esac
}
```

## Deployment Readiness Indicators

### Must Have
- Software starts without errors
- Basic functionality accessible
- No crash on startup
- Logs show normal operation

### Should Have
- Health checks pass
- Metrics exported
- Graceful shutdown works
- Configuration validated

### Nice to Have
- Auto-scaling tested
- Load balancing works
- Failover tested
- Monitoring integrated

## Integration with Other Rules

### Prerequisites
- R273: Runtime validation passed
- R274: Production checklist complete

### Enables
- R278: External user validation
- R279: PR plan generation

## Common Deployment Issues

### Issue: "Works locally but not when deployed"
**Solution**: Test in clean environment, check all dependencies

### Issue: "Container runs but service unreachable"
**Solution**: Verify port bindings, check network configuration

### Issue: "Deployment succeeds but app crashes"
**Solution**: Check logs, verify environment variables, test health checks

### Issue: "Terraform plan works but apply would fail"
**Solution**: Verify credentials, check resource quotas, test in sandbox

## Grading Impact

- Deployment fails completely: -30%
- Deployment succeeds but app doesn't work: -20%
- Missing deployment configuration: -15%
- No deployment verification attempted: -10%

## Summary

R275 ensures that Software Factory projects don't just compile/build but can actually be deployed and run in real environments. This catches the critical gap between "code complete" and "production ready".