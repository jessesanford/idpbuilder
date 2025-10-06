# Authentication Configuration

This guide explains how to configure authentication for pushing images to OCI registries.

## Authentication Methods

The push command supports three authentication methods with a specific precedence order:

### 1. Command-Line Flags (Highest Priority)

Explicitly provide credentials via flags:

```bash
idpbuilder push --username myuser --password mytoken
```

**Pros**:
- Explicit and clear
- Overrides all other methods

**Cons**:
- Credentials visible in shell history
- Not suitable for production

### 2. Environment Variables (Medium Priority)

Set credentials via environment variables:

```bash
export IDPBUILDER_REGISTRY_USER=myuser
export IDPBUILDER_REGISTRY_PASSWORD=mytoken
idpbuilder push
```

**Pros**:
- No credentials in command line
- Easy CI/CD integration
- Can be set in shell profiles

**Cons**:
- Visible to all processes
- Need to manage in environment

### 3. Docker Config File (Fallback)

Use existing Docker credentials from `~/.docker/config.json`:

```bash
docker login registry.example.com
idpbuilder push
```

**Pros**:
- Reuses existing Docker credentials
- Most secure (encrypted storage)
- No manual configuration needed

**Cons**:
- Requires Docker CLI
- May not work in all environments

## Credential Precedence

When multiple methods are configured, precedence is:

```
Command-line flags > Environment variables > Docker config
```

Example:

```bash
# Docker config has credentials for registry.example.com
# Environment variables also set
export IDPBUILDER_REGISTRY_USER=envuser

# This uses "flaguser" (highest precedence)
idpbuilder push --username flaguser --password flagpass

# This uses "envuser" (env var precedence)
idpbuilder push

# This uses Docker config credentials (fallback)
unset IDPBUILDER_REGISTRY_USER IDPBUILDER_REGISTRY_PASSWORD
idpbuilder push
```

## Security Best Practices

### 1. Never Hardcode Credentials

❌ **Don't do this**:
```bash
idpbuilder push --username admin --password secret123
```

✅ **Do this**:
```bash
idpbuilder push --username "$REGISTRY_USER" --password "$REGISTRY_TOKEN"
```

### 2. Use Access Tokens

Prefer registry access tokens over passwords:

```bash
# Generate token in registry UI
export IDPBUILDER_REGISTRY_PASSWORD=ghp_xxxxxxxxxxxxxxxxxxxx
```

### 3. Limit Token Scope

Create tokens with minimal required permissions:
- Read/Write access to specific repositories only
- Time-limited tokens for temporary access
- Separate tokens for different environments

### 4. Rotate Credentials Regularly

Implement credential rotation:

```bash
# Rotate every 90 days
export TOKEN_EXPIRY="2024-03-01"
export IDPBUILDER_REGISTRY_PASSWORD=$(get-fresh-token)
```

### 5. Use Secrets Management

In production, use secrets managers:

**Kubernetes Secret**:
```bash
export IDPBUILDER_REGISTRY_USER=$(kubectl get secret registry-creds -o jsonpath='{.data.username}' | base64 -d)
export IDPBUILDER_REGISTRY_PASSWORD=$(kubectl get secret registry-creds -o jsonpath='{.data.password}' | base64 -d)
```

**AWS Secrets Manager**:
```bash
export IDPBUILDER_REGISTRY_PASSWORD=$(aws secretsmanager get-secret-value --secret-id registry-token --query SecretString --output text)
```

## Token Management

### Creating Access Tokens

Different registries have different token creation processes:

**GitHub Container Registry**:
1. Go to Settings > Developer settings > Personal access tokens
2. Generate new token with `write:packages` scope
3. Use token as password

**GitLab Container Registry**:
1. Go to Settings > Access Tokens
2. Create token with `write_registry` scope
3. Use token as password

**Harbor**:
1. Go to User Profile > User Settings
2. Create Robot Account or CLI secret
3. Use as username/password pair

### Token Validation

Test token before using in production:

```bash
# Test with explicit credentials
idpbuilder push --username testuser --password testtoken --verbose

# Check for authentication errors
if [ $? -eq 2 ]; then
  echo "Token validation failed"
  exit 1
fi
```

## Common Authentication Issues

### Issue: 401 Unauthorized

**Causes**:
- Incorrect username or password
- Expired token
- Insufficient permissions

**Solutions**:
1. Verify credentials are correct
2. Check token expiration date
3. Ensure token has write permissions
4. Test with `docker login` first

### Issue: 403 Forbidden

**Causes**:
- Token lacks required scope
- Repository access denied
- Rate limiting

**Solutions**:
1. Verify token has `write` or `push` scope
2. Check repository permissions
3. Wait and retry if rate-limited

### Issue: Credential Precedence Confusion

**Problem**: Using wrong credentials source

**Solution**: Check precedence order:
```bash
# Debug which credentials are used
idpbuilder push --verbose 2>&1 | grep -i "auth"
```

## Environment-Specific Configuration

### Development

Use local Docker config:

```bash
docker login localhost:5000 --username dev --password dev
idpbuilder push --insecure
```

### CI/CD Pipelines

Use environment variables from secrets:

```yaml
# GitHub Actions
- name: Push Image
  env:
    IDPBUILDER_REGISTRY_USER: ${{ secrets.REGISTRY_USER }}
    IDPBUILDER_REGISTRY_PASSWORD: ${{ secrets.REGISTRY_TOKEN }}
  run: idpbuilder push
```

### Production

Use service accounts with restricted tokens:

```bash
export IDPBUILDER_REGISTRY_USER=service-account
export IDPBUILDER_REGISTRY_PASSWORD=$(vault read -field=token secret/registry/prod-token)
idpbuilder push
```

## See Also

- [Getting Started Guide](getting-started.md)
- [Push Command Guide](push-command.md)
- [Environment Variables Reference](../reference/environment-vars.md)
- [Troubleshooting Guide](troubleshooting.md)
