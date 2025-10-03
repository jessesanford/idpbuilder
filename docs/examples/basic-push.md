# Basic Push Examples

This document provides simple, copy-paste ready examples for common push scenarios.

## Push to Public Registry

Push an image without authentication:

```bash
# Build image
docker build -t myapp:latest .

# Push to public registry
idpbuilder push
```

## Push with Basic Authentication

### Using Command-Line Flags

```bash
idpbuilder push \
  --username myuser \
  --password mytoken
```

### Using Environment Variables

```bash
export IDPBUILDER_REGISTRY_USER=myuser
export IDPBUILDER_REGISTRY_PASSWORD=mytoken
idpbuilder push
```

### Using Docker Config

```bash
# Login with Docker first
docker login registry.example.com -u myuser -p mytoken

# Push using saved credentials
idpbuilder push
```

## Push to Development Registry

Push to local or development registry with self-signed certificate:

```bash
# Local registry (insecure)
idpbuilder push --insecure

# With custom port
idpbuilder push --insecure
```

## Push with Multiple Tags

Push the same image with different tags:

```bash
# Tag image multiple times
docker tag myapp:latest myapp:v1.0.0
docker tag myapp:latest myapp:stable

# Push each tag
idpbuilder push  # pushes :latest
idpbuilder push  # pushes :v1.0.0
idpbuilder push  # pushes :stable
```

## Quick Script Examples

### Simple Push Script

```bash
#!/bin/bash
set -e

# Configuration
IMAGE="myapp:latest"

# Push with error handling
if idpbuilder push; then
  echo "✓ Image pushed successfully"
else
  echo "✗ Push failed"
  exit 1
fi
```

### Push with Retry

```bash
#!/bin/bash

# Push with custom retry configuration
idpbuilder push \
  --retry-attempts 5 \
  --timeout 10m
```

## See Also

- [Advanced Examples](advanced-push.md)
- [CI/CD Integration](ci-integration.md)
- [Push Command Guide](../user-guide/push-command.md)
