# Push Command - Detailed Usage

This guide provides comprehensive information about using the `idpbuilder push` command.

## Command Syntax

```bash
idpbuilder push [flags]
```

## Image Detection

The push command automatically detects locally built images from your Docker daemon or compatible OCI runtime. It identifies:

- Images built with `docker build`
- Images created by `buildah` or `podman`
- Multi-architecture image manifests
- All image layers and metadata

## Flag Reference

### Authentication Flags

**`--username <string>`**
- Registry username for authentication
- Highest precedence (overrides environment variables)
- Used with `--password` for basic authentication

**`--password <string>`**
- Registry password or access token
- Required when `--username` is specified
- Security: avoid using in command history, prefer environment variables

### Registry Configuration

**`--insecure`**
- Skip TLS certificate verification
- Use for registries with self-signed certificates
- **Warning**: Only use in trusted development environments

**`--platform <string>`**
- Specify target platform for multi-arch images
- Format: `os/architecture` (e.g., `linux/amd64`, `linux/arm64`)
- Default: pushes all available platforms

### Operation Control

**`--timeout <duration>`**
- Maximum time to wait for push completion
- Format: duration string (e.g., `5m`, `30s`, `1h`)
- Default: `5m` (5 minutes)
- Recommendation: increase for large images

**`--retry-attempts <int>`**
- Maximum number of retry attempts on transient failures
- Default: `10`
- Range: `1-20`
- Exponential backoff applied between retries

**`--quiet`, `-q`**
- Suppress progress output
- Useful for scripting and CI/CD pipelines
- Errors still displayed to stderr

**`--verbose`, `-v`**
- Enable detailed logging
- Shows layer upload progress
- Displays retry attempts and backoff delays
- Useful for debugging

## Image Reference Formats

The push command accepts images in the standard OCI format:

### Basic Format
```
[registry/]repository[:tag]
```

### Examples

**Local repository**:
```bash
myapp:latest
```

**With registry**:
```bash
registry.example.com/myapp:v1.0
```

**With organization/namespace**:
```bash
registry.example.com/myorg/myapp:latest
```

**With port**:
```bash
registry.example.com:5000/myapp:dev
```

## Registry URL Formats

Supported registry URL formats:

- `registry.example.com` - HTTPS (default port 443)
- `registry.example.com:5000` - HTTPS with custom port
- `localhost:5000` - Local registry

## Flag Combinations

### Recommended Combinations

**Production push with authentication**:
```bash
idpbuilder push \
  --username $REGISTRY_USER \
  --password $REGISTRY_TOKEN \
  --timeout 10m
```

**Development push to local registry**:
```bash
idpbuilder push \
  --insecure \
  --timeout 2m
```

**CI/CD pipeline push (quiet mode)**:
```bash
idpbuilder push \
  --quiet \
  --retry-attempts 5
```

**Debug mode for troubleshooting**:
```bash
idpbuilder push \
  --verbose \
  --retry-attempts 3
```

### Invalid Combinations

These flag combinations will produce errors:

- `--quiet` and `--verbose` (mutually exclusive)
- `--timeout` with invalid duration format
- `--retry-attempts` outside valid range (1-20)

## Best Practices

### 1. Use Environment Variables for Credentials

Avoid exposing credentials in command line:

```bash
export IDPBUILDER_REGISTRY_USER=myuser
export IDPBUILDER_REGISTRY_PASSWORD=mytoken
idpbuilder push
```

### 2. Adjust Timeout for Image Size

For large images, increase the timeout:

```bash
idpbuilder push --timeout 15m  # for images > 1GB
```

### 3. Use Quiet Mode in Scripts

In automated scripts, suppress progress output:

```bash
if idpbuilder push --quiet; then
  echo "Push successful"
else
  echo "Push failed"
  exit 1
fi
```

### 4. Enable Verbose for Debugging

When troubleshooting issues:

```bash
idpbuilder push --verbose 2>&1 | tee push-debug.log
```

## See Also

- [Getting Started Guide](getting-started.md)
- [Authentication Configuration](authentication.md)
- [Troubleshooting Guide](troubleshooting.md)
- [Command Reference](../commands/push.md)
