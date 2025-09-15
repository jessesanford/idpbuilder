# Upstream Fixes Required for E2.2.1 and Related Efforts

**Date**: 2025-09-15
**Priority**: HIGH - Must fix before E2.2.2 implementation
**Affects**: E2.1.1 (image-builder), E2.2.1 (cli-commands)

## Issues Found During Testing

### 1. 🔴 Feature Flag Blocking Build Command

**Location**: `pkg/build/image_builder.go:42-44`
**Effort**: E2.1.1 (image-builder)

**Problem**:
```go
// BuildImage blocks unless environment variable is set
if !IsImageBuilderEnabled() {
    return nil, ErrFeatureDisabled
}
```

**Impact**:
- Build command fails with "image builder feature is disabled"
- Requires `export ENABLE_IMAGE_BUILDER=true` to work
- Poor user experience - should be enabled by default

**Fix Required**:
```go
// Option 1: Remove feature flag entirely (RECOMMENDED)
// Delete the check - feature should always be enabled in production

// Option 2: Default to enabled
func IsImageBuilderEnabled() bool {
    value := os.Getenv(EnableImageBuilderFlag)
    // Default to true if not explicitly disabled
    if value == "" {
        return true
    }
    return value != "false" && value != "0" && value != "disabled"
}
```

### 2. 🔴 Hardcoded Test Credentials

**Location**: `pkg/gitea/client.go:145-152`
**Effort**: E2.2.1 (cli-commands)

**Problem**:
```go
func getRegistryUsername() string {
    // Hardcoded during testing
    return "giteaAdmin"
}

func getRegistryPassword() string {
    // Hardcoded token from testing
    return "fd5e24fb3666d66621b12e8f0545757ea1ec74cb"
}
```

**Impact**:
- Credentials are hardcoded from our test environment
- Will fail for any other Gitea instance
- Security risk if committed to upstream

**Fix Required**:
- Will be completely replaced in E2.2.2 with proper credential management
- For now, should return empty strings with proper error handling

### 3. 🔴 Placeholder Manifest (Already Documented)

**Location**: `pkg/gitea/client.go:118-141`
**Effort**: E2.2.1 (cli-commands)

**Problem**: Returns placeholder JSON instead of real image
**Fix**: Complete rewrite in E2.2.2

## Fix Protocol

### Phase 1: Immediate Fixes (Before E2.2.2)

1. **Remove Feature Flag Requirement**
   - Edit `pkg/build/image_builder.go`
   - Remove lines 42-44 (the feature flag check)
   - OR change default to enabled

2. **Remove Hardcoded Credentials**
   - Edit `pkg/gitea/client.go`
   - Return empty strings from credential functions
   - Let E2.2.2 implement proper credential management

### Phase 2: Production Implementation (E2.2.2)

1. **Replace ALL placeholder code**
2. **Implement proper credential management**
3. **Add real image persistence and loading**
4. **Complete production push functionality**

## Testing After Fixes

```bash
# Should work WITHOUT setting ENABLE_IMAGE_BUILDER
./idpbuilder build --context ./test-app --tag myapp:v1.0.0

# Should fail gracefully with "no credentials" error
./idpbuilder push myapp:v1.0.0
# Error: No credentials provided (use --username/--token or set env vars)
```

## Backport Requirements (R321)

Per R321, these fixes must be backported to the original effort branches:
- E2.1.1 (image-builder) - Remove feature flag
- E2.2.1 (cli-commands) - Remove hardcoded credentials

The fixes should be applied in the original effort directories before integration.