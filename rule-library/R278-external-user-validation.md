# Rule R278: External User Validation

## Rule Statement
Every Software Factory project MUST be validated from the perspective of an external user who has no prior knowledge of the project. The software must be usable by following only the provided documentation.

## Criticality Level
**BLOCKING** - External users must be able to use the software
Violation = -35% grade penalty

## Core Principle
**"If external users can't use it, it's not ready"**

## External User Simulation Test

```bash
#!/bin/bash
# external-user-test.sh

echo "🧑‍💻 Starting External User Validation Test"
echo "========================================="
echo "Simulating a new user with no prior knowledge..."

# Create clean test environment
TEST_DIR="/tmp/external-user-test-$(date +%s)"
mkdir -p "$TEST_DIR"
cd "$TEST_DIR"

# Test Step 1: Can user clone the repository?
echo "📦 Step 1: Cloning repository..."
if ! git clone "$REPO_URL" project; then
    echo "❌ FAIL: Cannot clone repository"
    exit 1
fi
cd project

# Test Step 2: Can user find documentation?
echo "📚 Step 2: Looking for documentation..."
if [ ! -f "README.md" ]; then
    echo "❌ FAIL: No README.md found"
    exit 1
fi

if [ ! -f "RUNBOOK.md" ] && [ ! -f "docs/runbook.md" ]; then
    echo "⚠️ WARNING: No RUNBOOK found"
fi

# Test Step 3: Can user understand what the software does?
echo "🤔 Step 3: Understanding the software..."
if ! grep -q "## Overview\|## What is\|## Description\|## Purpose" README.md; then
    echo "❌ FAIL: No clear description of what software does"
    exit 1
fi

# Test Step 4: Are prerequisites clear?
echo "📋 Step 4: Checking prerequisites..."
if ! grep -q "## Prerequisites\|## Requirements\|## Dependencies" README.md RUNBOOK.md 2>/dev/null; then
    echo "⚠️ WARNING: Prerequisites not clearly documented"
fi

# Test Step 5: Can user build the software?
echo "🔨 Step 5: Building the software..."
BUILD_SUCCESS=false

# Try common build commands
if [ -f "Makefile" ] && grep -q "^build:" Makefile; then
    echo "Found Makefile, trying 'make build'..."
    if make build; then
        BUILD_SUCCESS=true
    fi
elif [ -f "go.mod" ]; then
    echo "Found go.mod, trying 'go build'..."
    if go build .; then
        BUILD_SUCCESS=true
    fi
elif [ -f "package.json" ]; then
    echo "Found package.json, trying 'npm install && npm run build'..."
    if npm install && npm run build; then
        BUILD_SUCCESS=true
    fi
elif [ -f "setup.py" ] || [ -f "pyproject.toml" ]; then
    echo "Found Python project, trying 'pip install -e .'..."
    if pip install -e .; then
        BUILD_SUCCESS=true
    fi
fi

if [ "$BUILD_SUCCESS" = "false" ]; then
    echo "❌ FAIL: Cannot build software following documentation"
    exit 1
fi

# Test Step 6: Can user run the software?
echo "🚀 Step 6: Running the software..."
RUN_SUCCESS=false

# Try to find and run the software
if [ -f "Makefile" ] && grep -q "^run:" Makefile; then
    timeout 5 make run && RUN_SUCCESS=true
elif [ -x "./app" ] || [ -x "./binary" ] || [ -x "./main" ]; then
    BINARY=$(find . -maxdepth 1 -type f -executable | head -1)
    timeout 5 "$BINARY" --help && RUN_SUCCESS=true
elif [ -f "package.json" ] && grep -q '"start"' package.json; then
    timeout 5 npm start && RUN_SUCCESS=true
elif [ -f "main.py" ] || [ -f "app.py" ]; then
    timeout 5 python main.py --help && RUN_SUCCESS=true
fi

if [ "$RUN_SUCCESS" = "false" ]; then
    echo "⚠️ WARNING: Cannot easily run the software"
fi

# Test Step 7: Can user run tests?
echo "🧪 Step 7: Running tests..."
TEST_SUCCESS=false

if [ -f "Makefile" ] && grep -q "^test:" Makefile; then
    make test && TEST_SUCCESS=true
elif [ -f "go.mod" ]; then
    go test ./... && TEST_SUCCESS=true
elif [ -f "package.json" ] && grep -q '"test"' package.json; then
    npm test && TEST_SUCCESS=true
elif command -v pytest &> /dev/null; then
    pytest && TEST_SUCCESS=true
fi

if [ "$TEST_SUCCESS" = "false" ]; then
    echo "⚠️ WARNING: No obvious way to run tests"
fi

# Test Step 8: Is help available?
echo "❓ Step 8: Checking for help..."
HELP_AVAILABLE=false

# Check for help command
if [ -x "./app" ]; then
    ./app --help &> /dev/null && HELP_AVAILABLE=true
    ./app -h &> /dev/null && HELP_AVAILABLE=true
fi

# Check for man page or help docs
if [ -f "man/app.1" ] || [ -f "docs/help.md" ]; then
    HELP_AVAILABLE=true
fi

if [ "$HELP_AVAILABLE" = "false" ]; then
    echo "⚠️ WARNING: No help system found"
fi

# Cleanup
cd /
rm -rf "$TEST_DIR"

echo "========================================="
echo "✅ External User Validation Complete"
```

## User Journey Validation

### Persona 1: Developer User
```yaml
developer_journey:
  goals:
    - Clone and build project
    - Run tests
    - Make modifications
    - Deploy locally
  
  validation_points:
    - Clear build instructions
    - Development setup guide
    - Contribution guidelines
    - Local deployment steps
```

### Persona 2: Operations User
```yaml
operations_journey:
  goals:
    - Deploy to production
    - Monitor application
    - Troubleshoot issues
    - Perform maintenance
  
  validation_points:
    - Deployment runbook exists
    - Monitoring setup documented
    - Troubleshooting guide present
    - Maintenance procedures clear
```

### Persona 3: End User
```yaml
end_user_journey:
  goals:
    - Install software
    - Use basic features
    - Get help when needed
    - Report issues
  
  validation_points:
    - Installation guide exists
    - Feature documentation clear
    - Help system available
    - Issue reporting process defined
```

## Documentation Quality Metrics

```python
def assess_documentation_quality():
    """Measure how well documented the project is"""
    
    score = 100
    findings = []
    
    # Check README
    if not os.path.exists("README.md"):
        score -= 25
        findings.append("No README.md found")
    else:
        readme_content = open("README.md").read()
        
        # Check for essential sections
        required_sections = [
            ("overview", r"#{1,2}\s*(Overview|Description|What is|Purpose)"),
            ("installation", r"#{1,2}\s*(Installation|Install|Getting Started)"),
            ("usage", r"#{1,2}\s*(Usage|How to|Examples)"),
            ("prerequisites", r"#{1,2}\s*(Prerequisites|Requirements|Dependencies)"),
        ]
        
        for section_name, pattern in required_sections:
            if not re.search(pattern, readme_content, re.IGNORECASE):
                score -= 10
                findings.append(f"Missing {section_name} section")
    
    # Check for RUNBOOK
    if not os.path.exists("RUNBOOK.md"):
        score -= 15
        findings.append("No RUNBOOK.md found")
    
    # Check for examples
    if not os.path.exists("examples/") and not "example" in readme_content.lower():
        score -= 10
        findings.append("No examples provided")
    
    # Check for API documentation
    if not any(os.path.exists(f) for f in ["API.md", "docs/api.md", "swagger.json"]):
        score -= 5
        findings.append("No API documentation")
    
    return score, findings
```

## Self-Contained Test

```bash
# The project must be self-contained
test_self_contained() {
    # Start with clean environment
    docker run --rm -v $(pwd):/workspace ubuntu:latest bash -c "
        apt-get update
        apt-get install -y git curl build-essential
        
        cd /workspace
        
        # Project should specify its requirements
        if [ -f .tool-versions ]; then
            # asdf tool versions specified
            echo '✓ Tool versions specified'
        elif [ -f .nvmrc ]; then
            # Node version specified
            echo '✓ Node version specified'
        elif [ -f .python-version ]; then
            # Python version specified
            echo '✓ Python version specified'
        fi
        
        # Should be able to build without external deps
        make build || go build . || npm install && npm run build
    "
}
```

## Common External User Issues

### Issue: "I don't know where to start"
**Solution**: README must have "Quick Start" section

### Issue: "Build fails with missing dependency"
**Solution**: All dependencies must be documented

### Issue: "I can't find how to configure it"
**Solution**: Configuration guide with examples required

### Issue: "It works but I don't know what it does"
**Solution**: Clear purpose/overview in documentation

### Issue: "I need help but don't know where to ask"
**Solution**: Support section with contact information

## Validation Checklist

### Minimum Requirements (Must Pass)
- [ ] Repository clones successfully
- [ ] README.md exists and is helpful
- [ ] Build instructions work
- [ ] Software runs after building
- [ ] Purpose is clearly explained

### Recommended Requirements (Should Pass)
- [ ] RUNBOOK.md for operations
- [ ] Examples provided
- [ ] Tests can be run
- [ ] Help system available
- [ ] Troubleshooting guide exists

### Excellence Indicators
- [ ] Video tutorials
- [ ] Interactive documentation
- [ ] Multiple examples
- [ ] Architecture diagrams
- [ ] Performance benchmarks

## Integration with Other Rules

### Prerequisites
- R274: Production readiness
- R276: Runbook exists

### Enables
- R279: PR plan can reference user validation
- SUCCESS state

## Grading Impact

**Failure Points:**
- No README: -15%
- Cannot build following docs: -20%
- Cannot run after building: -15%
- No clear purpose: -10%
- Missing prerequisites: -10%

**Warning Points:**
- No examples: -5%
- No troubleshooting guide: -5%
- No help system: -5%

## Summary

R278 ensures that Software Factory projects are usable by real external users who have no inside knowledge. This is the ultimate test of whether the documentation and user experience are complete. If someone outside the team can't clone, build, and run the software, it's not ready for the world.