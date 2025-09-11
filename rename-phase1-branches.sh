#!/bin/bash

# Script to rename Phase 1 Wave 1 branches to include project prefix
set -e

echo "🔄 Renaming Phase 1 Wave 1 branches on remote..."

# Clone the target repository
cd /tmp
rm -rf idpbuilder-rename
git clone https://github.com/jessesanford/idpbuilder.git idpbuilder-rename
cd idpbuilder-rename

# Function to rename a branch on remote
rename_branch() {
    local old_name="$1"
    local new_name="$2"
    
    echo "Renaming: $old_name -> $new_name"
    
    # Check if old branch exists
    if git show-ref --verify --quiet "refs/remotes/origin/$old_name"; then
        # Checkout the old branch
        git checkout -b "$old_name" "origin/$old_name" 2>/dev/null || git checkout "$old_name"
        
        # Create new branch with new name
        git checkout -b "$new_name"
        
        # Push new branch to remote
        git push origin "$new_name"
        
        # Delete old branch from remote
        git push origin --delete "$old_name"
        
        echo "✅ Renamed: $old_name -> $new_name"
    else
        echo "⚠️  Branch not found: $old_name"
    fi
}

# Rename main effort branches
rename_branch "phase1/wave1/effort-kind-cert-extraction" "idpbuilder-oci-build-push/phase1/wave1/kind-cert-extraction"
rename_branch "phase1/wave1/effort-registry-tls-trust" "idpbuilder-oci-build-push/phase1/wave1/registry-tls-trust"
rename_branch "phase1/wave1/registry-auth-types" "idpbuilder-oci-build-push/phase1/wave1/registry-auth-types"

# Rename split branches
rename_branch "phase1/wave1/registry-auth-types-split-001" "idpbuilder-oci-build-push/phase1/wave1/registry-auth-types-split-001"
rename_branch "phase1/wave1/registry-auth-types--split-002" "idpbuilder-oci-build-push/phase1/wave1/registry-auth-types-split-002"

# Also check for alternate split naming
rename_branch "phase1/wave1/registry-auth-types--split-001" "idpbuilder-oci-build-push/phase1/wave1/registry-auth-types-split-001"

echo "✅ Branch renaming complete!"