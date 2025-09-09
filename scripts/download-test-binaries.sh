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