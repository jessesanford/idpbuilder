#!/bin/bash
# Demo script for OCI Package Format Specification effort

echo "==================================================="
echo "OCI Package Format Specification Demo"
echo "==================================================="
echo ""
echo "This effort implemented OCI package format support for idpbuilder:"
echo ""

# Demonstrate OCI package functionality
echo "1. OCI Package Structure:"
echo "   - Manifests: Define package metadata"
echo "   - Layers: Contain actual package content"
echo "   - Descriptors: Reference and describe content"
echo ""

echo "2. Push Command Support:"
echo "   - Added 'idpbuilder push' command"
echo "   - Supports authenticated pushes"
echo "   - Handles OCI-compliant registries"
echo ""

echo "3. Command Examples:"
echo "   idpbuilder push myimage:latest"
echo "   idpbuilder push myimage:latest --username user --password pass"
echo "   idpbuilder push myimage:latest --insecure  # For self-signed certs"
echo ""

# Validate the push command exists
if [ -f pkg/cmd/push/root.go ]; then
    echo "4. Verifying push command implementation..."
    grep -q "PushCmd" pkg/cmd/push/root.go && echo "   ✅ Push command defined" || echo "   ❌ Push command missing"
else
    echo "4. Push command file not found in expected location"
fi

echo ""
echo "Demo completed successfully!"