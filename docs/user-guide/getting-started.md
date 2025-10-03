# Getting Started with idpbuilder push

This guide will help you quickly get started with pushing OCI images using the idpbuilder push command.

## Prerequisites

Before using the push command, ensure you have:

1. **idpbuilder installed**: The push command is part of idpbuilder
2. **Docker or compatible runtime**: For building OCI images locally
3. **Registry access**: Credentials for your target OCI registry
4. **Network connectivity**: Access to the target registry

## Your First Push

### Step 1: Build an Image

First, ensure you have a local image to push:

```bash
docker build -t myapp:latest .
```

### Step 2: Push to Registry

Push the image using basic authentication:

```bash
idpbuilder push --username myuser --password mytoken
```

That's it! Your image is now in the registry.

## Common Scenarios

### Push to Public Registry

For registries that don't require authentication:

```bash
idpbuilder push
```

### Push to Private Registry

Use environment variables for credentials:

```bash
export IDPBUILDER_REGISTRY_USER=myuser
export IDPBUILDER_REGISTRY_PASSWORD=mytoken
idpbuilder push
```

### Push to Self-Signed Registry

For development environments with self-signed certificates:

```bash
idpbuilder push --insecure
```

## Quick Troubleshooting

### Authentication Failed

If you see "authentication failed" errors:

1. Verify credentials are correct
2. Check if token has expired
3. Ensure user has push permissions

### TLS Certificate Error

If you see certificate verification errors:

1. Use `--insecure` flag for self-signed certificates
2. Or add the CA certificate to your system trust store

### Image Not Found

If the push fails to find the image:

1. Verify image exists locally: `docker images`
2. Check image name and tag are correct
3. Ensure Docker daemon is running

## Next Steps

- **Detailed Usage**: See [Push Command Guide](push-command.md) for comprehensive usage
- **Authentication**: Learn about [Authentication Methods](authentication.md)
- **Examples**: Explore [Basic Examples](../examples/basic-push.md) and [Advanced Examples](../examples/advanced-push.md)
- **Troubleshooting**: Review [Common Issues](troubleshooting.md) for problem resolution
