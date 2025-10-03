# Error Codes Reference

This document provides a complete reference of exit codes and error messages from the idpbuilder push command.

## Exit Codes

### 0 - Success

**Description**: Image pushed successfully

**Example Output**:
```
✓ Image pushed successfully
registry.example.com/repo/myapp:latest
```

**Action**: None required

---

### 1 - General Error

**Description**: Unspecified or general error occurred

**Common Causes**:
- Invalid command syntax
- Missing required arguments
- Internal error

**Example Output**:
```
Error: invalid argument
Failed to push image
```

**Resolution**:
1. Check command syntax
2. Review error message details
3. Enable verbose mode: `--verbose`

---

### 2 - Authentication Failure

**Description**: Failed to authenticate with registry

**Common Causes**:
- Incorrect username or password
- Expired access token
- Insufficient permissions
- Token lacks required scope

**Example Output**:
```
Error: authentication failed: 401 Unauthorized
Registry rejected credentials
```

**Resolution**:
1. Verify credentials are correct
2. Check token expiration
3. Validate user has push permissions
4. Regenerate access token if needed

**See**: [Authentication Guide](../user-guide/authentication.md)

---

### 3 - Network/Connectivity Error

**Description**: Network communication failed

**Common Causes**:
- Registry unreachable
- DNS resolution failure
- Network timeout
- Firewall blocking connection
- Connection reset

**Example Output**:
```
Error: dial tcp: connection refused
Cannot reach registry.example.com:5000
```

**Resolution**:
1. Verify registry URL is correct
2. Test network connectivity: `curl https://registry.example.com/v2/`
3. Check firewall rules
4. Verify DNS resolution
5. Increase timeout: `--timeout 15m`

**See**: [Troubleshooting Guide](../user-guide/troubleshooting.md#network-and-connectivity-issues)

---

### 4 - Image Not Found

**Description**: Specified image does not exist locally

**Common Causes**:
- Image name or tag incorrect
- Image not built yet
- Docker daemon not running

**Example Output**:
```
Error: image not found: myapp:latest
No matching image in local registry
```

**Resolution**:
1. List local images: `docker images`
2. Build image if missing: `docker build -t myapp:latest .`
3. Verify image name/tag spelling
4. Ensure Docker daemon is running

---

### 5 - Registry Error

**Description**: Registry returned an error

**Common Causes**:
- Registry unavailable or down
- Repository does not exist
- Registry storage full
- Registry configuration issue
- OCI spec incompatibility

**Example Output**:
```
Error: registry error: 500 Internal Server Error
Registry failed to process request
```

**Resolution**:
1. Check registry status/health
2. Verify repository exists
3. Check registry logs for details
4. Contact registry administrator
5. Test with Docker: `docker push`

---

### 6 - Operation Timeout

**Description**: Operation exceeded timeout limit

**Common Causes**:
- Large image size
- Slow network connection
- Registry performance issues
- Timeout value too low

**Example Output**:
```
Error: context deadline exceeded
Operation timed out after 5m0s
```

**Resolution**:
1. Increase timeout: `--timeout 15m`
2. Check network speed
3. Verify registry performance
4. For large images, use: `--timeout 30m`

---

## Common Error Messages

### Authentication Errors

| Error Message | Code | Cause | Resolution |
|--------------|------|-------|------------|
| `401 Unauthorized` | 2 | Invalid credentials | Verify username/password |
| `403 Forbidden` | 2 | Insufficient permissions | Check user access rights |
| `Token expired` | 2 | Expired access token | Generate new token |

### Network Errors

| Error Message | Code | Cause | Resolution |
|--------------|------|-------|------------|
| `connection refused` | 3 | Registry not reachable | Verify registry is running |
| `connection reset` | 3 | Network interruption | Check network stability |
| `DNS resolution failed` | 3 | Invalid hostname | Verify registry URL |
| `TLS handshake failed` | 3 | Certificate issue | Use `--insecure` or add CA cert |

### Image Errors

| Error Message | Code | Cause | Resolution |
|--------------|------|-------|------------|
| `image not found` | 4 | Image doesn't exist | Build image first |
| `invalid image format` | 4 | Corrupted image | Rebuild image |
| `layer not found` | 4 | Incomplete image | Rebuild from scratch |

### Registry Errors

| Error Message | Code | Cause | Resolution |
|--------------|------|-------|------------|
| `repository not found` | 5 | Repository doesn't exist | Create repository first |
| `500 Internal Server Error` | 5 | Registry malfunction | Check registry logs |
| `insufficient storage` | 5 | Registry disk full | Contact administrator |

## Error Handling Examples

### Check Exit Code in Scripts

```bash
idpbuilder push --quiet

case $? in
  0)
    echo "✓ Success"
    ;;
  2)
    echo "✗ Authentication failed"
    exit 1
    ;;
  3)
    echo "✗ Network error"
    exit 1
    ;;
  4)
    echo "✗ Image not found"
    exit 1
    ;;
  5)
    echo "✗ Registry error"
    exit 1
    ;;
  6)
    echo "✗ Operation timeout"
    exit 1
    ;;
  *)
    echo "✗ Unknown error"
    exit 1
    ;;
esac
```

### Retry on Specific Errors

```bash
#!/bin/bash

idpbuilder push
EXIT_CODE=$?

# Retry only on network errors
if [ $EXIT_CODE -eq 3 ]; then
  echo "Network error - retrying..."
  sleep 5
  idpbuilder push
fi
```

## Debugging

### Enable Verbose Output

For detailed error information:

```bash
idpbuilder push --verbose 2>&1 | tee error.log
```

### Capture Full Error Context

```bash
#!/bin/bash
set -x  # Enable command tracing

idpbuilder push --verbose 2>&1 | tee -a debug.log
EXIT_CODE=${PIPESTATUS[0]}

echo "Exit code: $EXIT_CODE" >> debug.log
echo "Timestamp: $(date)" >> debug.log
```

## See Also

- [Troubleshooting Guide](../user-guide/troubleshooting.md)
- [Push Command Reference](../commands/push.md)
- [Authentication Guide](../user-guide/authentication.md)
