# FIX INSTRUCTIONS - DEMO IMPLEMENTATION RECOVERY
Generated: 2025-09-10T21:22:55Z
Trigger: R291_DEMO_GATE_FAILURE

## CRITICAL FAILURE DETECTED

### Issue: Demo Implementations Not Completed
The SW Engineers spawned at 20:14:00Z for demo implementation have not completed their tasks.
No demo-features.sh files were created in any of the 4 efforts requiring demos.

## REQUIRED EFFORTS NEEDING DEMOS

### 1. image-builder
- Location: efforts/phase2/wave1/image-builder/
- Branch: idpbuilder-oci-build-push/phase2/wave1/image-builder
- Required Files:
  - DEMO-PLAN.md (missing)
  - demo-features.sh (missing)

### 2. gitea-client (original)
- Location: efforts/phase2/wave1/gitea-client/
- Branch: idpbuilder-oci-build-push/phase2/wave1/gitea-client
- Required Files:
  - DEMO-PLAN.md (missing)
  - demo-features.sh (missing)

### 3. gitea-client-split-001
- Location: efforts/phase2/wave1/gitea-client-split-001/
- Branch: idpbuilder-oci-build-push/phase2/wave1/gitea-client-split-001
- Required Files:
  - DEMO-PLAN.md (missing)
  - demo-features.sh (missing)

### 4. gitea-client-split-002
- Location: efforts/phase2/wave1/gitea-client-split-002/
- Branch: idpbuilder-oci-build-push/phase2/wave1/gitea-client-split-002
- Required Files:
  - DEMO-PLAN.md (missing)
  - demo-features.sh (missing)

## RECOVERY PLAN

### Step 1: Re-spawn SW Engineers for Demo Implementation
Spawn 4 SW Engineers in parallel (R151 compliant) with state DEMO_IMPLEMENTATION.

### Step 2: Engineer Instructions Template

For each engineer, provide these specific instructions:

#### Image Builder Engineer:
```
cd efforts/phase2/wave1/image-builder
git checkout idpbuilder-oci-build-push/phase2/wave1/image-builder

CREATE DEMO-PLAN.md with:
- Feature: Image building and registry operations
- Demo scenarios:
  1. Build simple image from source
  2. Push to local registry
  3. Verify image metadata
  4. List available images

CREATE demo-features.sh:
#!/bin/bash
set -e
echo "🎬 DEMO: Image Builder Features"
echo "================================"

# Demo 1: Build image
echo "📦 Building sample image..."
# Add actual demo commands based on implementation

# Demo 2: Registry operations
echo "🚀 Pushing to registry..."
# Add registry demo

# Demo 3: Verify operations
echo "✅ Verifying image..."
# Add verification

echo "✅ Demo completed successfully!"
exit 0
```

#### Gitea Client Engineer (Original):
```
cd efforts/phase2/wave1/gitea-client
git checkout idpbuilder-oci-build-push/phase2/wave1/gitea-client

CREATE DEMO-PLAN.md with:
- Feature: Gitea API client operations
- Demo scenarios:
  1. Connect to Gitea instance
  2. List repositories
  3. Create new repository
  4. Clone and push content

CREATE demo-features.sh with appropriate demo logic
```

#### Gitea Client Split-001 Engineer:
```
cd efforts/phase2/wave1/gitea-client-split-001
git checkout idpbuilder-oci-build-push/phase2/wave1/gitea-client-split-001

CREATE DEMO-PLAN.md with:
- Feature: Core Gitea authentication and connection
- Demo scenarios specific to split-001 functionality

CREATE demo-features.sh with split-specific demos
```

#### Gitea Client Split-002 Engineer:
```
cd efforts/phase2/wave1/gitea-client-split-002
git checkout idpbuilder-oci-build-push/phase2/wave1/gitea-client-split-002

CREATE DEMO-PLAN.md with:
- Feature: Advanced Gitea operations
- Demo scenarios specific to split-002 functionality

CREATE demo-features.sh with split-specific demos
```

## SUCCESS CRITERIA

Each effort MUST have:
1. ✅ DEMO-PLAN.md describing what the demo will show
2. ✅ demo-features.sh that:
   - Is executable (chmod +x)
   - Exits with 0 on success
   - Shows actual feature functionality
   - Has clear output messages
   - Handles basic error cases
3. ✅ All demos committed and pushed to their branches
4. ✅ Demos can run repeatedly without failure

## VALIDATION REQUIREMENTS (R291)

After implementation, ALL demos must:
- Build successfully
- Pass all tests
- Run demo-features.sh without errors
- Exit with status 0

Any failure requires immediate return to ERROR_RECOVERY.

## TIMELINE
- Recovery Start: 2025-09-10T21:22:55Z
- Expected Completion: Within 30 minutes
- Validation Gate: R291 mandatory verification

## NEXT STEPS AFTER FIX
1. Monitor all 4 engineers for completion
2. Verify each demo-features.sh exists and runs
3. Transition to INTEGRATION_DEMO_TESTING
4. Spawn Code Reviewer for R291 gate verification
5. Only proceed to RETROFIT_COMPLETE if all gates pass