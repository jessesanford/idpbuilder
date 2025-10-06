# Troubleshooting Guide

This guide helps you diagnose and resolve common issues with the idpbuilder push command.

## Common Issues

### Authentication Failures

#### Error: 401 Unauthorized

**Symptoms**:
```
Error: authentication failed: 401 Unauthorized
Failed to push image to registry.example.com/repo
```

**Causes**:
- Incorrect username or password
- Expired access token
- Token lacks required permissions

**Solutions**:

1. **Verify credentials**:
```bash
# Test with docker login first
docker login registry.example.com -u myuser -p mytoken
```

2. **Check token expiration**:
```bash
# Regenerate token if expired
# Update environment variable
export IDPBUILDER_REGISTRY_PASSWORD=new-token
```

3. **Validate token scope**:
- Ensure token has `write` or `push` permissions
- Check repository-specific access rights

4. **Clear cached credentials**:
```bash
# Remove Docker config
rm ~/.docker/config.json
# Try fresh login
docker login registry.example.com
```

#### Error: 403 Forbidden

**Symptoms**:
```
Error: access denied: 403 Forbidden
User does not have push access to repository
```

**Solutions**:
1. Verify user has write/push permissions in registry UI
2. Check if repository exists (create if needed)
3. Ensure correct repository path/namespace

### TLS Certificate Errors

#### Error: x509 Certificate Verification

**Symptoms**:
```
Error: x509: certificate signed by unknown authority
TLS handshake failed for registry.example.com
```

**Causes**:
- Self-signed certificate
- Missing CA certificate in system trust store
- Certificate chain incomplete

**Solutions**:

1. **Use --insecure flag** (development only):
```bash
idpbuilder push --insecure
```

2. **Add CA certificate to system** (production):
```bash
# Linux
sudo cp ca.crt /usr/local/share/ca-certificates/
sudo update-ca-certificates

# macOS
sudo security add-trusted-cert -d -r trustRoot -k /Library/Keychains/System.keychain ca.crt
```

3. **Set environment variable**:
```bash
export IDPBUILDER_REGISTRY_INSECURE=true
idpbuilder push
```

### Network and Connectivity Issues

#### Error: Connection Timeout

**Symptoms**:
```
Error: context deadline exceeded
Failed to connect to registry after 5m0s
```

**Solutions**:

1. **Increase timeout**:
```bash
idpbuilder push --timeout 15m
```

2. **Check network connectivity**:
```bash
# Test registry access
curl -I https://registry.example.com/v2/

# Check DNS resolution
nslookup registry.example.com
```

3. **Verify firewall rules**:
- Ensure outbound HTTPS (443) allowed
- Check if registry port (e.g., 5000) accessible

#### Error: Connection Refused

**Symptoms**:
```
Error: dial tcp: connection refused
Cannot connect to registry.example.com:5000
```

**Solutions**:
1. Verify registry is running: `curl http://registry.example.com:5000/v2/`
2. Check correct port in URL
3. Ensure registry is accessible from your network

### Image and Repository Issues

#### Error: Image Not Found

**Symptoms**:
```
Error: image not found locally
No image matching 'myapp:latest' found
```

**Solutions**:

1. **List local images**:
```bash
docker images | grep myapp
```

2. **Build image if missing**:
```bash
docker build -t myapp:latest .
```

3. **Check image name/tag**:
```bash
# Ensure exact match (case-sensitive)
docker images myapp:latest
```

#### Error: Repository Does Not Exist

**Symptoms**:
```
Error: repository not found
registry.example.com/myorg/myapp does not exist
```

**Solutions**:
1. Create repository in registry UI first
2. Check repository path is correct
3. Verify namespace/organization exists

### Retry and Performance Issues

#### Error: Max Retries Exceeded

**Symptoms**:
```
Error: max retries (10) exceeded
Last error: connection reset by peer
```

**Solutions**:

1. **Check network stability**:
```bash
# Test connectivity
ping -c 10 registry.example.com
```

2. **Reduce retry attempts for faster failure**:
```bash
idpbuilder push --retry-attempts 3
```

3. **Increase retry attempts for flaky networks**:
```bash
idpbuilder push --retry-attempts 20
```

#### Slow Push Performance

**Symptoms**:
- Push taking > 10 minutes for small images
- Progress stalls on specific layers

**Solutions**:

1. **Check network bandwidth**:
```bash
# Test upload speed
curl -T testfile.bin https://registry.example.com/upload-test
```

2. **Verify registry performance**:
- Check registry logs for errors
- Monitor registry resource usage

3. **Use verbose logging to identify bottleneck**:
```bash
idpbuilder push --verbose 2>&1 | tee push-debug.log
```

## Advanced Debugging

### Enable Debug Logging

Get detailed operation logs:

```bash
idpbuilder push --verbose 2>&1 | tee -a push.log
```

Look for:
- Authentication method used
- Retry attempts and backoff delays
- Layer upload progress
- Network errors

### Check Registry Compatibility

Verify OCI registry compatibility:

```bash
# Check registry supports OCI spec
curl -I https://registry.example.com/v2/

# Look for:
# HTTP/2 200
# docker-distribution-api-version: registry/2.0
```

### Test with Docker

Compare with Docker push:

```bash
# Tag for registry
docker tag myapp:latest registry.example.com/repo/myapp:latest

# Push with Docker
docker push registry.example.com/repo/myapp:latest

# If Docker works but idpbuilder doesn't, report issue
```

## Error Code Reference

Quick reference for return codes:

- `0`: Success
- `1`: General error (check logs)
- `2`: Authentication failure (verify credentials)
- `3`: Network error (check connectivity)
- `4`: Image not found (verify local image exists)
- `5`: Registry error (check registry status)
- `6`: Timeout (increase --timeout)

See [Error Codes Reference](../reference/error-codes.md) for detailed information.

## Getting Help

If issues persist:

1. **Collect debug information**:
```bash
idpbuilder push --verbose 2>&1 | tee debug.log
idpbuilder version > version.txt
docker version >> version.txt
```

2. **Check common issues**: Review this guide thoroughly

3. **Search existing issues**: Check project issue tracker

4. **Report new issue**: Provide:
   - Complete error message
   - Debug logs (`push --verbose`)
   - Environment details (OS, Docker version)
   - Steps to reproduce

## See Also

- [Getting Started Guide](getting-started.md)
- [Push Command Guide](push-command.md)
- [Authentication Configuration](authentication.md)
- [Error Codes Reference](../reference/error-codes.md)
