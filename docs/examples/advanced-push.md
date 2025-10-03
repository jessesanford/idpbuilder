# Advanced Push Examples

This document provides examples for complex push scenarios and advanced use cases.

## Multi-Architecture Images

Push multi-platform images:

```bash
# Build for multiple platforms
docker buildx build --platform linux/amd64,linux/arm64 -t myapp:latest .

# Push specific platform
idpbuilder push --platform linux/amd64

# Push all platforms (default behavior)
idpbuilder push
```

## Complex Authentication Scenarios

### Using Tokens from Secrets Manager

**AWS Secrets Manager**:
```bash
#!/bin/bash
TOKEN=$(aws secretsmanager get-secret-value \
  --secret-id registry-token \
  --query SecretString \
  --output text)

idpbuilder push --username robot --password "$TOKEN"
```

**HashiCorp Vault**:
```bash
#!/bin/bash
export IDPBUILDER_REGISTRY_PASSWORD=$(vault kv get -field=token secret/registry/prod)
idpbuilder push
```

### Dynamic Credential Rotation

```bash
#!/bin/bash
get_fresh_token() {
  # Generate new token
  curl -X POST https://registry.example.com/api/token \
    -d "grant_type=client_credentials" \
    | jq -r '.access_token'
}

# Use fresh token
export IDPBUILDER_REGISTRY_PASSWORD=$(get_fresh_token)
idpbuilder push
```

## Custom Registry Configurations

### Push to Registry with Custom Port

```bash
idpbuilder push
```

### Push Through HTTP Proxy

```bash
export HTTP_PROXY=http://proxy.example.com:8080
export HTTPS_PROXY=http://proxy.example.com:8080
idpbuilder push
```

## Batch Operations

### Push Multiple Images

```bash
#!/bin/bash
IMAGES=(
  "app-frontend:latest"
  "app-backend:latest"
  "app-worker:latest"
)

for image in "${IMAGES[@]}"; do
  echo "Pushing $image..."
  idpbuilder push || {
    echo "Failed to push $image"
    exit 1
  }
done
```

### Parallel Push with Different Registries

```bash
#!/bin/bash

# Function to push to specific registry
push_to_registry() {
  local registry=$1
  idpbuilder push
}

# Push to multiple registries in parallel
push_to_registry "registry1.example.com" &
push_to_registry "registry2.example.com" &
wait

echo "All pushes complete"
```

## Error Handling and Logging

### Comprehensive Error Handling

```bash
#!/bin/bash
set -euo pipefail

# Enable verbose logging
LOGFILE="push-$(date +%Y%m%d-%H%M%S).log"

idpbuilder push --verbose 2>&1 | tee "$LOGFILE"

EXIT_CODE=${PIPESTATUS[0]}

case $EXIT_CODE in
  0)
    echo "✓ Success"
    ;;
  2)
    echo "✗ Authentication failed - check credentials"
    exit 1
    ;;
  3)
    echo "✗ Network error - check connectivity"
    exit 1
    ;;
  *)
    echo "✗ Unknown error (code: $EXIT_CODE)"
    exit 1
    ;;
esac
```

### Retry with Exponential Backoff

```bash
#!/bin/bash

MAX_ATTEMPTS=5
ATTEMPT=1

while [ $ATTEMPT -le $MAX_ATTEMPTS ]; do
  echo "Attempt $ATTEMPT of $MAX_ATTEMPTS..."

  if idpbuilder push --quiet; then
    echo "✓ Push successful"
    exit 0
  fi

  if [ $ATTEMPT -lt $MAX_ATTEMPTS ]; then
    DELAY=$((2 ** ATTEMPT))
    echo "Waiting ${DELAY}s before retry..."
    sleep $DELAY
  fi

  ATTEMPT=$((ATTEMPT + 1))
done

echo "✗ Failed after $MAX_ATTEMPTS attempts"
exit 1
```

## Performance Optimization

### Large Image Push

```bash
# For images > 1GB
idpbuilder push \
  --timeout 30m \
  --retry-attempts 15
```

### Quiet Mode for Performance

```bash
# Disable progress output for faster operation
idpbuilder push --quiet
```

## Conditional Push Logic

### Push Only on Main Branch

```bash
#!/bin/bash

if [ "$GIT_BRANCH" = "main" ]; then
  echo "Main branch detected - pushing to production"
  export IDPBUILDER_REGISTRY_USER=$PROD_USER
  export IDPBUILDER_REGISTRY_PASSWORD=$PROD_TOKEN
  idpbuilder push
else
  echo "Not main branch - skipping push"
fi
```

### Push with Tag Validation

```bash
#!/bin/bash

# Only push if tag is semantic version
if [[ "$IMAGE_TAG" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
  echo "Valid semantic version: $IMAGE_TAG"
  idpbuilder push
else
  echo "Invalid tag format: $IMAGE_TAG"
  exit 1
fi
```

## See Also

- [Basic Examples](basic-push.md)
- [CI/CD Integration](ci-integration.md)
- [Push Command Guide](../user-guide/push-command.md)
- [Troubleshooting Guide](../user-guide/troubleshooting.md)
