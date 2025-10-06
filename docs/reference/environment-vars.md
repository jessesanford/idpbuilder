# Environment Variables Reference

This document lists all environment variables supported by the idpbuilder push command.

## Authentication Variables

### IDPBUILDER_REGISTRY_USER

**Description**: Default registry username for authentication

**Type**: String

**Default**: None (empty)

**Example**:
```bash
export IDPBUILDER_REGISTRY_USER=myuser
```

**Precedence**: Lower than command-line flags, higher than Docker config

---

### IDPBUILDER_REGISTRY_PASSWORD

**Description**: Default registry password or access token

**Type**: String (sensitive)

**Default**: None (empty)

**Example**:
```bash
export IDPBUILDER_REGISTRY_PASSWORD=mytoken
```

**Security Note**: Use access tokens instead of passwords when possible

**Precedence**: Lower than command-line flags, higher than Docker config

---

## Registry Configuration Variables

### IDPBUILDER_REGISTRY_INSECURE

**Description**: Skip TLS certificate verification

**Type**: Boolean string ("true" or "false")

**Default**: "false"

**Example**:
```bash
export IDPBUILDER_REGISTRY_INSECURE=true
```

**Warning**: Only use in trusted development environments

---

### IDPBUILDER_REGISTRY_URL

**Description**: Default registry URL

**Type**: URL string

**Default**: Detected from image reference

**Example**:
```bash
export IDPBUILDER_REGISTRY_URL=registry.example.com:5000
```

---

## Operation Control Variables

### IDPBUILDER_PUSH_TIMEOUT

**Description**: Default timeout for push operations

**Type**: Duration string (e.g., "5m", "30s", "1h")

**Default**: "5m" (5 minutes)

**Example**:
```bash
export IDPBUILDER_PUSH_TIMEOUT=10m
```

---

### IDPBUILDER_PUSH_RETRY_ATTEMPTS

**Description**: Maximum number of retry attempts

**Type**: Integer (1-20)

**Default**: 10

**Example**:
```bash
export IDPBUILDER_PUSH_RETRY_ATTEMPTS=15
```

---

### IDPBUILDER_PUSH_QUIET

**Description**: Suppress progress output

**Type**: Boolean string ("true" or "false")

**Default**: "false"

**Example**:
```bash
export IDPBUILDER_PUSH_QUIET=true
```

---

### IDPBUILDER_PUSH_VERBOSE

**Description**: Enable detailed logging

**Type**: Boolean string ("true" or "false")

**Default**: "false"

**Example**:
```bash
export IDPBUILDER_PUSH_VERBOSE=true
```

**Note**: Mutually exclusive with IDPBUILDER_PUSH_QUIET

---

## Proxy and Network Variables

### HTTP_PROXY / HTTPS_PROXY

**Description**: HTTP/HTTPS proxy server URL

**Type**: URL string

**Default**: None

**Example**:
```bash
export HTTPS_PROXY=http://proxy.example.com:8080
```

**Note**: Standard environment variables, not specific to idpbuilder

---

### NO_PROXY

**Description**: Comma-separated list of hosts to exclude from proxy

**Type**: Comma-separated string

**Default**: None

**Example**:
```bash
export NO_PROXY=localhost,127.0.0.1,registry.local
```

---

## Precedence Rules

When multiple configuration sources are present:

### For Authentication:
1. Command-line flags (`--username`, `--password`) - highest
2. Environment variables (`IDPBUILDER_REGISTRY_USER/PASSWORD`)
3. Docker config file (`~/.docker/config.json`) - lowest

### For Other Settings:
1. Command-line flags (e.g., `--timeout`) - highest
2. Environment variables (e.g., `IDPBUILDER_PUSH_TIMEOUT`)
3. Default values - lowest

## Usage Examples

### Basic Configuration

```bash
# Set credentials
export IDPBUILDER_REGISTRY_USER=myuser
export IDPBUILDER_REGISTRY_PASSWORD=mytoken

# Configure operation
export IDPBUILDER_PUSH_TIMEOUT=10m
export IDPBUILDER_PUSH_RETRY_ATTEMPTS=15

# Run push
idpbuilder push
```

### CI/CD Configuration

```bash
# Use secrets from CI/CD platform
export IDPBUILDER_REGISTRY_USER=$CI_REGISTRY_USER
export IDPBUILDER_REGISTRY_PASSWORD=$CI_REGISTRY_TOKEN

# Quiet mode for cleaner logs
export IDPBUILDER_PUSH_QUIET=true

# Execute push
idpbuilder push
```

### Development Configuration

```bash
# Local registry with self-signed cert
export IDPBUILDER_REGISTRY_INSECURE=true
export IDPBUILDER_PUSH_VERBOSE=true

# Push with debug output
idpbuilder push
```

## See Also

- [Authentication Guide](../user-guide/authentication.md)
- [Push Command Reference](../commands/push.md)
- [CI/CD Integration Examples](../examples/ci-integration.md)
