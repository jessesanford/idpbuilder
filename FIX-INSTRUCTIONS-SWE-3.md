# Fix Instructions for SW Engineer 3 - Test Infrastructure Setup
Date: 2025-09-09
Assigned by: orchestrator/COORDINATE_BUILD_FIXES
Priority: HIGH

## Your Assignment
Fix missing etcd binary issue for controller tests

## Working Directory
```bash
cd /home/vscode/workspaces/idpbuilder-oci-build-push/efforts/project/integration-workspace
```

## Fix Plan Reference
See: FIX-PLAN-BUILD-FAILURES.md Error 3 (lines 147-228)

## Problem Description
Tests fail with: `fork/exec ../../../bin/k8s/1.29.1-linux-amd64/etcd: no such file or directory`
The test environment (envtest) needs Kubernetes binaries including etcd.

## Solution Approach
Create a script to download test binaries and ensure they're available for tests.

## Specific Tasks

### Task 1: Create Binary Download Script
**Create file**: scripts/download-test-binaries.sh
```bash
#!/bin/bash
# Download test binaries for envtest

# Set up directories
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJECT_ROOT="$( cd "$SCRIPT_DIR/.." && pwd )"
BIN_DIR="$PROJECT_ROOT/bin"

# Create bin directory structure
mkdir -p "$BIN_DIR"

# Install setup-envtest if not present
if ! command -v setup-envtest &> /dev/null; then
    echo "Installing setup-envtest..."
    go install sigs.k8s.io/controller-runtime/tools/setup-envtest@latest
fi

# Download binaries for version 1.29.1
echo "Downloading Kubernetes test binaries v1.29.1..."
KUBEBUILDER_ASSETS=$(setup-envtest use 1.29.1 --bin-dir "$BIN_DIR" -p path)

if [ $? -eq 0 ]; then
    echo "✅ Test binaries successfully installed"
    echo "Location: $KUBEBUILDER_ASSETS"
    
    # Export for immediate use
    export KUBEBUILDER_ASSETS
    
    # Create a marker file for other scripts
    echo "$KUBEBUILDER_ASSETS" > "$BIN_DIR/.kubebuilder-assets-path"
    
    # Verify etcd exists
    if [ -f "$KUBEBUILDER_ASSETS/etcd" ]; then
        echo "✅ etcd binary verified at: $KUBEBUILDER_ASSETS/etcd"
    else
        echo "⚠️ Warning: etcd binary not found in expected location"
    fi
else
    echo "❌ Failed to download test binaries"
    exit 1
fi
```

### Task 2: Make Script Executable
```bash
chmod +x scripts/download-test-binaries.sh
```

### Task 3: Run the Script
```bash
./scripts/download-test-binaries.sh
```

### Task 4: Update Test to Use Downloaded Binaries (Optional)
If the test still can't find binaries, update pkg/controllers/custompackage/controller_test.go:

Look for the BinaryAssetsDirectory configuration (around lines 47-48) and ensure it points to the correct location, or remove it to let envtest use the downloaded binaries:

**Option A: Remove hardcoded path (Recommended)**
Comment out or remove the BinaryAssetsDirectory line to let envtest find the binaries automatically.

**Option B: Update path**
Update to use the actual downloaded path from KUBEBUILDER_ASSETS environment variable.

## Testing Requirements
```bash
# Step 1: Run the download script
./scripts/download-test-binaries.sh

# Step 2: Verify binaries exist
ls -la bin/k8s/*/etcd
# Should show etcd binary

# Step 3: Run controller tests
export KUBEBUILDER_ASSETS=$(cat bin/.kubebuilder-assets-path)
go test ./pkg/controllers/custompackage -v
# Should NOT fail with "missing etcd" error

# Step 4: Create completion marker
touch FIX-COMPLETE-SWE-3.marker
echo "Test infrastructure setup completed at $(date)" > FIX-COMPLETE-SWE-3.marker
```

## Success Indicators
- ✅ Test binaries downloaded successfully
- ✅ etcd binary exists in bin directory
- ✅ Controller tests can start test environment
- ✅ No "fork/exec etcd: no such file" errors
- ✅ FIX-COMPLETE-SWE-3.marker created

## Important Notes
- This fix enables the test infrastructure for Error 4 fix
- May take a few minutes to download binaries (about 50MB)
- Independent fix - can be done in parallel with Errors 1 & 2
- Once complete, SWE-4 can proceed with Error 4 fix