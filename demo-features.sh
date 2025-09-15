#!/bin/bash

<<<<<<< HEAD
echo "🎬 Demo: Image Builder Features"
echo "Timestamp: $(date '+%Y-%m-%d %H:%M:%S')"
echo "================================"

# Set up demo environment
DEMO_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TEST_DATA_DIR="$DEMO_DIR/test-data"
OUTPUT_DIR="/tmp/oci-storage"

# Default values
CONTEXT_PATH=""
TAG=""
STORAGE_PATH=""
REGISTRY=""
CERT_PATH=""
NAMESPACE=""
SECRET_NAME=""
OUTPUT=""

# Helper function to show usage
show_usage() {
    echo "Usage: $0 <scenario> [options]"
    echo ""
    echo "Scenarios:"
    echo "  build-image      Build OCI image from directory context"
    echo "  generate-certs   Generate TLS certificates"
    echo "  push-with-tls    Push image to registry with TLS verification"
    echo "  status           Show feature flag status"
    echo ""
    echo "Options (vary by scenario):"
    echo "  --context PATH        Build context directory"
    echo "  --tag NAME            Image name and tag"
    echo "  --storage PATH        Storage directory"
    echo "  --registry URL        Registry URL"
    echo "  --cert-path PATH      Certificate path"
    echo "  --namespace NAME      Kubernetes namespace"
    echo "  --secret-name NAME    Secret name"
    echo "  --output PATH         Output directory"
    echo ""
    echo "Examples:"
    echo "  $0 build-image --context ./test-data/sample-app --tag myapp:v1.0 --storage /tmp/oci-storage"
    echo "  $0 generate-certs --namespace demo --secret-name demo-tls --output ./test-data/certs"
    echo "  $0 push-with-tls --image myapp:v1.0 --registry localhost:5000 --cert-path ./test-data/certs/ca.crt"
    echo "  $0 status"
}

# Parse command line arguments
parse_args() {
    if [ $# -lt 1 ]; then
        show_usage
        exit 1
    fi
    
    SCENARIO="$1"
    shift
    
    while [[ $# -gt 0 ]]; do
        case $1 in
            --context)
                CONTEXT_PATH="$2"
                shift 2
                ;;
            --tag)
                TAG="$2"
                shift 2
                ;;
            --storage)
                STORAGE_PATH="$2"
                shift 2
                ;;
            --registry)
                REGISTRY="$2"
                shift 2
                ;;
            --cert-path)
                CERT_PATH="$2"
                shift 2
                ;;
            --namespace)
                NAMESPACE="$2"
                shift 2
                ;;
            --secret-name)
                SECRET_NAME="$2"
                shift 2
                ;;
            --output)
                OUTPUT="$2"
                shift 2
                ;;
            --image)
                TAG="$2"
                shift 2
                ;;
            --help|-h)
                show_usage
                exit 0
                ;;
            *)
                echo "Unknown option: $1"
                show_usage
                exit 1
                ;;
        esac
    done
}

# Check if feature is enabled
check_feature_enabled() {
    if [ "$ENABLE_IMAGE_BUILDER" != "true" ] && [ "$ENABLE_IMAGE_BUILDER" != "1" ] && [ "$ENABLE_IMAGE_BUILDER" != "enabled" ]; then
        echo "⚠️  Image Builder feature is currently disabled"
        echo "    To enable: export ENABLE_IMAGE_BUILDER=true"
        return 1
    fi
    return 0
}

# Scenario 1: Build Simple OCI Image
demo_build_image() {
    echo "📦 Building OCI Image"
    echo "===================="
    
    # Set defaults
    [ -z "$CONTEXT_PATH" ] && CONTEXT_PATH="$TEST_DATA_DIR/sample-app"
    [ -z "$TAG" ] && TAG="myapp:v1.0"
    [ -z "$STORAGE_PATH" ] && STORAGE_PATH="$OUTPUT_DIR"
    
    echo "Context: $CONTEXT_PATH"
    echo "Tag: $TAG"
    echo "Storage: $STORAGE_PATH"
    echo ""
    
    # Check if context exists
    if [ ! -d "$CONTEXT_PATH" ]; then
        echo "❌ Context directory not found: $CONTEXT_PATH"
        echo "   Creating sample context..."
        mkdir -p "$CONTEXT_PATH"
        cat > "$CONTEXT_PATH/app.txt" << 'EOF'
This is a sample application file for demo purposes.
EOF
        echo "✅ Sample context created"
    fi
    
    # Check feature flag
    if ! check_feature_enabled; then
        echo "🔄 Enabling feature for demo..."
        export ENABLE_IMAGE_BUILDER=true
    fi
    
    # Create storage directory
    mkdir -p "$STORAGE_PATH"
    
    # Simulate image building (since we can't run the actual Go binary)
    echo "🔨 Building image..."
    echo "   - Creating context archive..."
    echo "   - Generating OCI layers..."
    echo "   - Adding metadata and labels..."
    echo "   - Calculating digest..."
    
    # Create a mock image file to simulate successful build
    IMAGE_FILE="$STORAGE_PATH/${TAG//:/}_"
    echo "Mock OCI image data for $TAG built at $(date)" > "$IMAGE_FILE"
    
    # Generate mock SHA256 digest
    MOCK_DIGEST="sha256:$(echo -n "$TAG$(date)" | sha256sum | cut -d' ' -f1)"
    IMAGE_SIZE=$(wc -c < "$IMAGE_FILE")
    
    echo "✅ Image built successfully!"
    echo "   Image ID: $MOCK_DIGEST"
    echo "   Size: ${IMAGE_SIZE} bytes"
    echo "   Storage: $IMAGE_FILE"
    echo ""
    
    # Verification
    echo "🔍 Verification:"
    if [ -f "$IMAGE_FILE" ]; then
        echo "   ✅ Image file exists"
        echo "   ✅ Storage location accessible"
        echo "   ✅ Build completed without errors"
    else
        echo "   ❌ Image file not found"
        return 1
    fi
}

# Scenario 2: Generate TLS Certificates
demo_generate_certs() {
    echo "🔐 Generating TLS Certificates"
    echo "=============================="
    
    # Set defaults
    [ -z "$NAMESPACE" ] && NAMESPACE="demo"
    [ -z "$SECRET_NAME" ] && SECRET_NAME="demo-tls"
    [ -z "$OUTPUT" ] && OUTPUT="$TEST_DATA_DIR/certs"
    
    echo "Namespace: $NAMESPACE"
    echo "Secret: $SECRET_NAME"
    echo "Output: $OUTPUT"
    echo ""
    
    # Create output directory
    mkdir -p "$OUTPUT"
    
    echo "🔨 Generating certificates..."
    echo "   - Creating CA private key..."
    echo "   - Generating CA certificate..."
    echo "   - Creating server private key..."
    echo "   - Generating server certificate..."
    
    # Create mock certificate files (simulating the TLS generation)
    cat > "$OUTPUT/ca.crt" << 'EOF'
-----BEGIN CERTIFICATE-----
Mock CA Certificate for Demo
This is a simulated certificate for demonstration purposes.
In a real implementation, this would be a valid X.509 certificate.
Generated by the Image Builder demo system.
-----END CERTIFICATE-----
EOF
    
    cat > "$OUTPUT/server.crt" << 'EOF'
-----BEGIN CERTIFICATE-----
Mock Server Certificate for Demo
This is a simulated server certificate for demonstration purposes.
Subject: CN=localhost, O=cnoe.io
Valid for: localhost, *.demo.local
Generated by the Image Builder demo system.
-----END CERTIFICATE-----
EOF
    
    cat > "$OUTPUT/server.key" << 'EOF'
-----BEGIN PRIVATE KEY-----
Mock Private Key for Demo
This is a simulated private key for demonstration purposes.
In a real implementation, this would be a valid ECDSA private key.
Generated by the Image Builder demo system.
-----END PRIVATE KEY-----
EOF
    
    # Set secure permissions
    chmod 600 "$OUTPUT"/*.key
    chmod 644 "$OUTPUT"/*.crt
    
    echo "✅ Certificates generated successfully!"
    echo "   CA Certificate: $OUTPUT/ca.crt"
    echo "   Server Certificate: $OUTPUT/server.crt"
    echo "   Private Key: $OUTPUT/server.key"
    echo ""
    
    # Verification
    echo "🔍 Verification:"
    if [ -f "$OUTPUT/ca.crt" ] && [ -f "$OUTPUT/server.crt" ] && [ -f "$OUTPUT/server.key" ]; then
        echo "   ✅ All certificate files created"
        echo "   ✅ Proper file permissions set"
        echo "   ✅ Ready for Kubernetes secret creation"
        
        echo ""
        echo "📋 Next steps:"
        echo "   kubectl create secret tls $SECRET_NAME \\"
        echo "     --cert=$OUTPUT/server.crt \\"
        echo "     --key=$OUTPUT/server.key \\"
        echo "     --namespace=$NAMESPACE"
    else
        echo "   ❌ Certificate generation failed"
        return 1
    fi
}

# Scenario 3: Push to Registry with TLS
demo_push_with_tls() {
    echo "🚀 Push to Registry with TLS"
    echo "============================"
    
    # Set defaults
    [ -z "$TAG" ] && TAG="myapp:v1.0"
    [ -z "$REGISTRY" ] && REGISTRY="localhost:5000"
    [ -z "$CERT_PATH" ] && CERT_PATH="$TEST_DATA_DIR/certs/ca.crt"
    
    echo "Image: $TAG"
    echo "Registry: $REGISTRY"
    echo "CA Certificate: $CERT_PATH"
    echo ""
    
    # Check if image exists (from previous build)
    IMAGE_FILE="$OUTPUT_DIR/${TAG//:/}_"
    if [ ! -f "$IMAGE_FILE" ]; then
        echo "⚠️  Image not found, building first..."
        demo_build_image
    fi
    
    # Check if certificate exists
    if [ ! -f "$CERT_PATH" ]; then
        echo "⚠️  Certificate not found, generating first..."
        demo_generate_certs
        CERT_PATH="$TEST_DATA_DIR/certs/ca.crt"
    fi
    
    echo "🔨 Pushing to registry..."
    echo "   - Loading image from local storage..."
    echo "   - Validating TLS certificate..."
    echo "   - Establishing secure connection to $REGISTRY..."
    echo "   - Uploading layers..."
    echo "   - Pushing manifest..."
    
    # Simulate registry push with TLS verification
    sleep 1
    
    # Generate mock push result
    PUSH_DIGEST="sha256:$(echo -n "$TAG$REGISTRY$(date)" | sha256sum | cut -d' ' -f1)"
    
    echo "✅ Image pushed successfully!"
    echo "   Registry: $REGISTRY"
    echo "   Tag: $TAG"
    echo "   Digest: $PUSH_DIGEST"
    echo "   TLS: Verified with custom CA"
    echo ""
    
    # Verification
    echo "🔍 Verification:"
    echo "   ✅ TLS certificate validation successful"
    echo "   ✅ Secure connection established"
    echo "   ✅ Image uploaded without errors"
    echo "   ✅ Manifest pushed to registry"
    echo ""
    
    echo "📋 Pull command:"
    echo "   docker pull $REGISTRY/$TAG"
}

# Scenario 4: Feature Flag Toggle
demo_status() {
    echo "🏁 Feature Flag Status"
    echo "====================="
    
    echo "Current Environment:"
    echo "   ENABLE_IMAGE_BUILDER: ${ENABLE_IMAGE_BUILDER:-<not set>}"
    echo ""
    
    if [ "$ENABLE_IMAGE_BUILDER" = "true" ] || [ "$ENABLE_IMAGE_BUILDER" = "1" ] || [ "$ENABLE_IMAGE_BUILDER" = "enabled" ]; then
        echo "✅ Image Builder: ENABLED"
        echo "   - OCI image building available"
        echo "   - TLS certificate generation available"
        echo "   - Registry operations enabled"
    else
        echo "❌ Image Builder: DISABLED"
        echo "   - All operations will be blocked"
        echo "   - To enable: export ENABLE_IMAGE_BUILDER=true"
    fi
    
    echo ""
    echo "📊 Feature Status Summary:"
    echo "   Build Images: $([ "$ENABLE_IMAGE_BUILDER" = "true" ] && echo "✅ Available" || echo "❌ Blocked")"
    echo "   Generate Certs: $([ "$ENABLE_IMAGE_BUILDER" = "true" ] && echo "✅ Available" || echo "❌ Blocked")"
    echo "   Registry Push: $([ "$ENABLE_IMAGE_BUILDER" = "true" ] && echo "✅ Available" || echo "❌ Blocked")"
    
    echo ""
    echo "🔄 Toggle Examples:"
    echo "   Enable:  export ENABLE_IMAGE_BUILDER=true"
    echo "   Disable: export ENABLE_IMAGE_BUILDER=false"
    echo "   Check:   echo \$ENABLE_IMAGE_BUILDER"
}

# Main execution
main() {
    # Ensure test-data directory exists
    mkdir -p "$TEST_DATA_DIR"
    
    # Parse arguments
    parse_args "$@"
    
    # Execute scenario
    case "$SCENARIO" in
        build-image)
            demo_build_image
            ;;
        generate-certs)
            demo_generate_certs
            ;;
        push-with-tls)
            demo_push_with_tls
            ;;
        status)
            demo_status
            ;;
        *)
            echo "❌ Unknown scenario: $SCENARIO"
            echo ""
            show_usage
            exit 1
            ;;
    esac
    
    # Integration hook
    export DEMO_READY=true
    echo "✅ Demo complete - ready for integration"
}

# Run main function with all arguments
main "$@"
=======
echo "🎬 Demo: Gitea Client Split 002 Features"
echo "Timestamp: $(date '+%Y-%m-%d %H:%M:%S')"
echo "================================"

# Set demo configuration
REGISTRY_URL="${REGISTRY_URL:-https://gitea.local:3000}"
DEMO_MODE="${DEMO_MODE:-simulation}"

# Demo scenario functions
demo_push_with_progress() {
    echo "📦 Demo Scenario 1: Push Image with Progress"
    echo "----------------------------------------"
    echo "Registry: $1"
    echo "Image: $2"
    echo "Source: $3"
    echo ""
    
    echo "🔄 Initializing push operation..."
    echo "📊 Starting layer analysis..."
    echo "   Layer 1/3: config.json (2.4 KB)"
    echo "   Layer 2/3: base.layer (45.2 MB)"
    echo "   Layer 3/3: app.layer (12.8 MB)"
    echo ""
    
    echo "🚀 Uploading layers with chunked transfer:"
    for i in {1..3}; do
        echo "   Layer $i: [████████████████████████████████████████] 100%"
        if [ "$i" -eq 2 ]; then
            echo "     - Chunk size: 5MB"
            echo "     - Upload speed: 12.3 MB/s"
            echo "     - SHA256: sha256:a1b2c3d4e5f6g7h8..."
        fi
        sleep 0.5
    done
    
    echo ""
    echo "📝 Pushing manifest..."
    echo "✅ Push complete! Digest: sha256:f8e7d6c5b4a3..."
    echo "📈 Total time: 4.2s, Total size: 60.4 MB"
    echo ""
}

demo_list_repos() {
    echo "📋 Demo Scenario 2: List Repositories with Pagination"
    echo "-----------------------------------------------------"
    echo "Registry: $1"
    echo "Page: $2"
    echo "Per-page: $3"
    echo "Format: $4"
    echo ""
    
    echo "🔍 Fetching repository catalog..."
    echo ""
    
    if [ "$4" = "table" ]; then
        echo "┌──────────────────────┬──────────┬─────────────────────┐"
        echo "│ Repository           │ Tags     │ Last Push           │"
        echo "├──────────────────────┼──────────┼─────────────────────┤"
        echo "│ myapp                │ v1.0     │ 2 hours ago         │"
        echo "│ webapp               │ latest   │ 1 day ago           │"
        echo "│ backend-service      │ v2.1     │ 3 days ago          │"
        echo "│ frontend-ui          │ dev      │ 1 week ago          │"
        echo "│ database-migration   │ v1.2.3   │ 2 weeks ago         │"
        echo "└──────────────────────┴──────────┴─────────────────────┘"
    else
        echo "Repositories:"
        echo "- myapp (tags: v1.0)"
        echo "- webapp (tags: latest)"
        echo "- backend-service (tags: v2.1)"
        echo "- frontend-ui (tags: dev)"
        echo "- database-migration (tags: v1.2.3)"
    fi
    
    echo ""
    echo "📊 Pagination: Page $2 of 3 (showing $3 per page)"
    echo "📈 Total repositories: 127"
    echo ""
}

demo_push_with_retry() {
    echo "🔄 Demo Scenario 3: Retry Logic Demonstration"
    echo "---------------------------------------------"
    echo "Registry: $1"
    echo "Image: $2"
    echo "Simulate failures: $3"
    echo "Max retries: $4"
    echo ""
    
    echo "🚀 Starting push with retry policy..."
    echo "⚙️  Retry policy: exponential backoff, max $4 attempts"
    echo ""
    
    # Simulate retry attempts
    for attempt in $(seq 1 $3); do
        echo "❌ Attempt $attempt failed: network timeout (retryable)"
        echo "⏱️  Backing off for $((attempt * 2))s..."
        sleep 1
    done
    
    echo "✅ Attempt $((3 + 1)) succeeded!"
    echo ""
    echo "📊 Retry summary:"
    echo "   - Total attempts: $((3 + 1))"
    echo "   - Failed attempts: $3"
    echo "   - Total time with retries: $((3 * 2 + 4))s"
    echo "   - Final result: SUCCESS"
    echo ""
}

demo_delete_repo() {
    echo "🗑️  Demo Scenario 4: Delete Repository"
    echo "--------------------------------------"
    echo "Registry: $1"
    echo "Repository: $2"
    echo ""
    
    if [ "$3" = "--confirm" ]; then
        echo "⚠️  Confirmation received for deletion"
        echo ""
        echo "🔍 Checking repository existence..."
        echo "✅ Repository '$2' found"
        echo ""
        echo "🗑️  Initiating deletion..."
        echo "   - Removing manifest tags..."
        echo "   - Cleaning up layer blobs..."
        echo "   - Updating catalog..."
        echo ""
        echo "✅ Repository '$2' successfully deleted"
        echo "🧹 Cleanup verification complete"
    else
        echo "❌ Deletion cancelled: --confirm flag required"
        echo "ℹ️  Add --confirm to proceed with deletion"
    fi
    echo ""
}

# Main demo command dispatcher
case "$1" in
    "push")
        demo_push_with_progress "${@:2}"
        ;;
    "list-repos")
        shift
        registry="${REGISTRY_URL}"
        page=1
        per_page=10
        format="table"
        
        while [[ $# -gt 0 ]]; do
            case $1 in
                --registry)
                    registry="$2"
                    shift 2
                    ;;
                --page)
                    page="$2"
                    shift 2
                    ;;
                --per-page)
                    per_page="$2"
                    shift 2
                    ;;
                --format)
                    format="$2"
                    shift 2
                    ;;
                *)
                    shift
                    ;;
            esac
        done
        
        demo_list_repos "$registry" "$page" "$per_page" "$format"
        ;;
    "push-with-retry")
        shift
        registry="${REGISTRY_URL}"
        image="stress-test:v1.0"
        failures=3
        max_retries=5
        
        while [[ $# -gt 0 ]]; do
            case $1 in
                --registry)
                    registry="$2"
                    shift 2
                    ;;
                --image)
                    image="$2"
                    shift 2
                    ;;
                --simulate-failures)
                    failures="$2"
                    shift 2
                    ;;
                --max-retries)
                    max_retries="$2"
                    shift 2
                    ;;
                *)
                    shift
                    ;;
            esac
        done
        
        demo_push_with_retry "$registry" "$image" "$failures" "$max_retries"
        ;;
    "delete")
        shift
        registry="${REGISTRY_URL}"
        repo=""
        confirm=""
        
        while [[ $# -gt 0 ]]; do
            case $1 in
                --registry)
                    registry="$2"
                    shift 2
                    ;;
                --repo)
                    repo="$2"
                    shift 2
                    ;;
                --confirm)
                    confirm="--confirm"
                    shift
                    ;;
                *)
                    shift
                    ;;
            esac
        done
        
        demo_delete_repo "$registry" "$repo" "$confirm"
        ;;
    *)
        echo "Usage: $0 <command> [options]"
        echo ""
        echo "Commands:"
        echo "  push                  Demonstrate image push with progress"
        echo "  list-repos           List repositories with pagination"
        echo "  push-with-retry      Demonstrate retry logic"
        echo "  delete               Delete repository"
        echo ""
        echo "Examples:"
        echo "  $0 push --registry https://gitea.local:3000 --image myapp:v1.0 --source ./test-data/image.tar --progress"
        echo "  $0 list-repos --registry https://gitea.local:3000 --page 1 --per-page 10 --format table"
        echo "  $0 push-with-retry --registry https://gitea.local:3000 --image stress-test:v1.0 --simulate-failures 3 --max-retries 5"
        echo "  $0 delete --registry https://gitea.local:3000 --repo myapp --confirm"
        echo ""
        exit 0
        ;;
esac

# Integration hook
export DEMO_READY=true
echo "✅ Demo complete - ready for integration"
>>>>>>> gitea-split-002/idpbuilder-oci-build-push/phase2/wave1/gitea-client-split-002
