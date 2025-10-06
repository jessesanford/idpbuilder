# idpbuilder push

## Synopsis

Push OCI-compliant container images to any OCI registry.

```bash
idpbuilder push [flags]
```

## Description

The `push` command uploads locally built OCI container images to a remote OCI-compliant registry. It supports various authentication methods, handles network retries automatically, and provides detailed progress reporting.

## Flags

### Authentication Flags

- `--username <string>`: Registry username for authentication (highest precedence)
- `--password <string>`: Registry password or token for authentication (highest precedence)

### Registry Flags

- `--insecure`: Skip TLS certificate verification (use for self-signed certificates)
- `--platform <string>`: Target platform for multi-architecture images (e.g., linux/amd64)

### Operation Flags

- `--timeout <duration>`: Maximum time to wait for push operation (default: 5m)
- `--retry-attempts <int>`: Maximum number of retry attempts on failure (default: 10)
- `--quiet`, `-q`: Suppress progress output
- `--verbose`, `-v`: Enable detailed logging

### Global Flags

- `--help`, `-h`: Display help for push command
- `--version`: Display version information

## Authentication

Authentication credentials are evaluated in the following precedence order:

1. **Command-line flags** (`--username` and `--password`) - highest priority
2. **Environment variables** (`IDPBUILDER_REGISTRY_USER` and `IDPBUILDER_REGISTRY_PASSWORD`)
3. **Docker config file** (`~/.docker/config.json`) - fallback

See [Authentication Guide](../user-guide/authentication.md) for detailed information.

## Environment Variables

The following environment variables are supported:

- `IDPBUILDER_REGISTRY_USER`: Default registry username
- `IDPBUILDER_REGISTRY_PASSWORD`: Default registry password or token
- `IDPBUILDER_REGISTRY_INSECURE`: Set to "true" to skip TLS verification
- `IDPBUILDER_PUSH_TIMEOUT`: Default push operation timeout (e.g., "10m")

See [Environment Variables Reference](../reference/environment-vars.md) for complete list.

## Examples

### Basic Push

Push an image to a public registry:

```bash
idpbuilder push
```

### Push with Authentication

Push using command-line credentials:

```bash
idpbuilder push --username myuser --password mytoken
```

Push using environment variables:

```bash
export IDPBUILDER_REGISTRY_USER=myuser
export IDPBUILDER_REGISTRY_PASSWORD=mytoken
idpbuilder push
```

### Push to Insecure Registry

Push to a registry with self-signed certificate:

```bash
idpbuilder push --insecure
```

### Push with Retry Configuration

Customize retry behavior:

```bash
idpbuilder push --retry-attempts 5 --timeout 10m
```

## Return Codes

- `0`: Success - image pushed successfully
- `1`: General error
- `2`: Authentication failure
- `3`: Network/connectivity error
- `4`: Image not found or invalid
- `5`: Registry error or incompatibility
- `6`: Operation timeout

See [Error Codes Reference](../reference/error-codes.md) for detailed information.

## Output

The push command provides real-time progress feedback:

```
Pushing image myimage:latest to registry.example.com/repo...
✓ Layer 1/3 pushed (2.5 MB)
✓ Layer 2/3 pushed (15.3 MB)
✓ Layer 3/3 pushed (1.2 MB)
✓ Manifest pushed successfully
Image pushed: registry.example.com/repo/myimage:latest
```

Use `--quiet` to suppress progress output or `--verbose` for detailed logging.

## See Also

- [Getting Started Guide](../user-guide/getting-started.md)
- [Push Command Detailed Usage](../user-guide/push-command.md)
- [Authentication Configuration](../user-guide/authentication.md)
- [Troubleshooting Guide](../user-guide/troubleshooting.md)
- [Basic Examples](../examples/basic-push.md)
- [Advanced Examples](../examples/advanced-push.md)

## Notes

- The push command automatically detects and pushes all layers of multi-architecture images
- Network failures trigger automatic retry with exponential backoff (up to 10 attempts by default)
- Progress reporting is disabled in non-TTY environments (e.g., CI/CD pipelines)
- Large images may take several minutes to push depending on network speed
