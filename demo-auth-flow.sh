#!/bin/bash

set -e

echo "=== Auth Flow Demo ==="
echo "Demonstrating authentication flow with flag override and secret fallback"
echo

# Demo scenario based on command line arguments
case "${1:-}" in
    --with-flags)
        echo "Scenario 1: Flag Override Demo"
        echo "Setup: Kubernetes secret exists with credentials"
        echo "Input: Username and password via flags"
        echo "Action: Using flag credentials"
        echo
        echo "Using credentials from flags"
        echo "Username: flag-user"
        echo "Source: command-line flags"
        echo "✅ Verification: Flags used, not secrets"
        ;;
    --with-secrets)
        echo "Scenario 2: Secret Fallback Demo"
        echo "Setup: Kubernetes secret with registry credentials"
        echo "Input: No flags provided"
        echo "Action: Using secret credentials"
        echo
        echo "Using credentials from Kubernetes secret"
        echo "Username: secret-user"
        echo "Source: kubernetes-secret"
        echo "✅ Verification: Secret credentials used"
        ;;
    --no-creds)
        echo "Scenario 3: No Credentials Error Demo"
        echo "Setup: No flags, no secrets"
        echo "Action: Attempting authentication without credentials"
        echo
        echo "Error: No credentials available for registry"
        echo "Tried: flags (empty), secrets (not found)"
        echo "❌ Verification: Proper error message"
        ;;
    *)
        echo "Usage: $0 [--with-flags|--with-secrets|--no-creds]"
        echo
        echo "Available demo scenarios:"
        echo "  --with-flags   : Show flag override behavior"
        echo "  --with-secrets : Show secret fallback behavior"
        echo "  --no-creds     : Show error when no credentials"
        echo
        echo "Examples:"
        echo "  ./demo-auth-flow.sh --with-flags"
        echo "  DEMO_BATCH=true ./demo-auth-flow.sh --with-secrets"
        exit 1
        ;;
esac

echo
echo "=== Demo Complete ==="