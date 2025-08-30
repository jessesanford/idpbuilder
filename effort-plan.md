# CLI Commands Implementation Plan

## Effort Infrastructure Metadata (R209)
**EFFORT_NAME**: cli-commands
**BRANCH**: idpbuilder-oci-mvp/phase2/wave2/cli-commands
**WORKING_DIRECTORY**: /home/vscode/workspaces/idpbuilder-oci-mvp/efforts/phase2/wave2/cli-commands
**PARENT_WAVE**: Phase 2 Wave 2
**SEQUENTIAL_POSITION**: 1 of 2 (MUST complete before integration-testing)

## Critical Effort Metadata (FROM WAVE PLAN - R211)
**Branch**: `idpbuilder-oci-mvp/phase2/wave2/cli-commands`
**Can Parallelize**: No
**Parallel With**: None
**Size Estimate**: 500 lines (MUST be <800)
**Dependencies**: 
  - Phase 1: Certificate infrastructure (all efforts completed)
  - Phase 2 Wave 1: Build wrapper (2.1.1), Registry client (2.1.2)

## Overview
- **Effort**: Implement CLI commands for OCI build and push operations
- **Phase**: 2, Wave: 2
- **Estimated Size**: 500 lines (450 target, 500 soft limit)
- **Implementation Time**: 1.5 days

## Dependency Analysis (R219 Compliance)

### Phase 1 Dependencies (Certificate Infrastructure)
**What We Import**:
```go
// From Phase 1 efforts (assumed package structure based on patterns)
import (
    "github.com/idpbuilder/idpbuilder-oci-mvp/pkg/certs/extractor"
    "github.com/idpbuilder/idpbuilder-oci-mvp/pkg/certs/trust"
    "github.com/idpbuilder/idpbuilder-oci-mvp/pkg/certs/validator"
    "github.com/idpbuilder/idpbuilder-oci-mvp/pkg/certs/fallback"
)
```

**How It Influences Implementation**:
- CLI initialization will auto-extract certificates using Phase 1 components
- Build/Push commands will configure trust before operations
- Error messages will leverage validator's detailed failure modes
- --insecure flag will trigger fallback handler

### Phase 2 Wave 1 Dependencies (Build & Registry)
**What We Import**:
```go
// From Wave 1 efforts (verified to exist)
import (
    "github.com/idpbuilder/idpbuilder-oci-mvp/efforts/phase2/wave1/buildah-build-wrapper/pkg/build"
    "github.com/idpbuilder/idpbuilder-oci-mvp/efforts/phase2/wave1/gitea-registry-client/pkg/registry"
)
```

**How It Influences Implementation**:
- CLI commands are thin wrappers around Builder and GiteaClient
- Pass trust manager from Phase 1 to Wave 1 components
- Forward progress callbacks to CLI output
- Handle authentication using registry client patterns

## File Structure
```
cli-commands/
├── pkg/
│   ├── cmd/
│   │   ├── build.go          # Build command (~50 lines)
│   │   ├── push.go           # Push command (~50 lines)
│   │   └── root.go           # Root command setup (~30 lines)
│   ├── oci/
│   │   ├── certs/
│   │   │   ├── auto_config.go    # Auto certificate config (~80 lines)
│   │   │   └── initializer.go    # Certificate init logic (~70 lines)
│   │   ├── commands/
│   │   │   ├── build_handler.go  # Build command logic (~100 lines)
│   │   │   └── push_handler.go   # Push command logic (~100 lines)
│   │   └── config/
│   │       └── settings.go       # OCI configuration (~50 lines)
│   └── main.go                    # Entry point (~20 lines)
└── tests/
    └── unit/
        ├── cmd_test.go            # Command tests (included in count)
        └── config_test.go         # Config tests (included in count)
```

## Implementation Steps

### Step 1: Command Structure Setup (130 lines)
**Files**: `pkg/main.go`, `pkg/cmd/root.go`, `pkg/cmd/build.go`, `pkg/cmd/push.go`

1. Create main.go entry point:
   - Initialize root command
   - Add version info
   - Handle panic recovery

2. Implement root command:
   - Set up Cobra command structure
   - Configure global flags (--verbose, --config)
   - Add subcommand registration

3. Create build command:
   - Define flags: --file, --context, --tag, --platform
   - Add validation for required arguments
   - Set up command help and examples

4. Create push command:
   - Define flags: --insecure, --username, --password
   - Parse image reference argument
   - Add usage documentation

**Size checkpoint**: ~130 lines

### Step 2: Certificate Auto-Configuration (150 lines)
**Files**: `pkg/oci/certs/auto_config.go`, `pkg/oci/certs/initializer.go`

1. Implement auto_config.go:
   ```go
   func AutoConfigureCertificates() error {
       // 1. Check if Kind cluster exists
       // 2. Extract Gitea certificates using Phase 1 extractor
       // 3. Configure trust store using Phase 1 trust manager
       // 4. Validate certificate chain using Phase 1 validator
       // 5. Cache result for subsequent commands
   }
   ```

2. Implement initializer.go:
   ```go
   func InitializeTrustEnvironment() (*trust.Manager, error) {
       // 1. Auto-configure if not already done
       // 2. Load existing trust configuration
       // 3. Set up environment variables
       // 4. Return configured trust manager
   }
   ```

**Integration points**:
- Use Phase 1 KindCertExtractor
- Use Phase 1 TrustStoreManager
- Use Phase 1 CertValidator

**Size checkpoint**: ~280 lines total

### Step 3: Build Command Implementation (100 lines)
**Files**: `pkg/oci/commands/build_handler.go`

1. Implement build handler:
   ```go
   func ExecuteBuild(dockerfile, context, tag string, opts BuildOptions) error {
       // 1. Initialize trust environment
       trustMgr, err := certs.InitializeTrustEnvironment()
       
       // 2. Create Wave 1 Builder with trust config
       builder := build.NewBuilder(build.WithTrustManager(trustMgr))
       
       // 3. Set up progress reporting
       progressChan := make(chan build.Progress)
       
       // 4. Execute build
       return builder.BuildImage(ctx, dockerfile, context, tag, opts)
   }
   ```

2. Add progress reporting to terminal
3. Handle build errors with clear messages
4. Support platform selection

**Size checkpoint**: ~380 lines total

### Step 4: Push Command Implementation (100 lines)
**Files**: `pkg/oci/commands/push_handler.go`

1. Implement push handler:
   ```go
   func ExecutePush(imageRef string, opts PushOptions) error {
       // 1. Initialize trust environment (unless --insecure)
       var trustMgr *trust.Manager
       if !opts.Insecure {
           trustMgr, err = certs.InitializeTrustEnvironment()
       }
       
       // 2. Create Wave 1 GiteaClient
       client := registry.NewGiteaClient(
           registry.WithTrustManager(trustMgr),
           registry.WithInsecure(opts.Insecure),
       )
       
       // 3. Authenticate if credentials provided
       if opts.Username != "" {
           client.Authenticate(opts.Username, opts.Password)
       }
       
       // 4. Execute push with progress
       return client.Push(ctx, imageRef)
   }
   ```

2. Add push progress display
3. Handle authentication errors
4. Support --insecure bypass mode

**Size checkpoint**: ~480 lines total

### Step 5: Configuration Management (20 lines)
**Files**: `pkg/oci/config/settings.go`

1. Define configuration structure:
   ```go
   type OCIConfig struct {
       RegistryURL      string // Default: gitea.local:443
       CacheDir         string // Default: ~/.idpbuilder/cache
       AutoExtractCerts bool   // Default: true
       InsecureMode     bool   // Default: false
   }
   ```

2. Support environment variables:
   - IDPBUILDER_OCI_REGISTRY
   - IDPBUILDER_OCI_CACHE_DIR
   - IDPBUILDER_OCI_INSECURE

**Final size**: ~500 lines

## Size Management
- **Estimated Lines**: 500
- **Measurement Tool**: /home/vscode/workspaces/idpbuilder-oci-mvp/tools/line-counter.sh
- **Check Frequency**: After each step completion
- **Split Threshold**: 700 lines (warning), 800 lines (stop)

### Size Monitoring Points
1. After Step 1 (Command Structure): Check = 130 lines
2. After Step 2 (Certificates): Check = 280 lines  
3. After Step 3 (Build Command): Check = 380 lines
4. After Step 4 (Push Command): Check = 480 lines
5. After Step 5 (Config): Final Check = 500 lines

### If Approaching Limit
- Defer advanced configuration features to post-MVP
- Simplify error messages and help text
- Move shared utilities to common package
- Reduce test verbosity

## Test Requirements
- **Unit Tests**: 80% coverage minimum
- **Test Focus**:
  - Command flag parsing and validation
  - Configuration loading from environment
  - Certificate detection logic
  - Error handling paths
- **Test Files**: Included in line count (~50 lines for tests)

## Pattern Compliance
- **Command Structure**: Follow Cobra patterns if adding to existing CLI, else create simple structure
- **Error Handling**: Return errors with context, don't panic
- **Configuration**: Environment variables override defaults
- **Progress Display**: Use standard output for user feedback
- **Security**: Never log sensitive credentials

## Success Criteria
1. ✅ `idpbuilder build` command successfully builds test images
2. ✅ `idpbuilder push` command pushes to Gitea registry
3. ✅ Certificates auto-configure on first use
4. ✅ --insecure flag provides bypass option
5. ✅ Clear error messages for all failure modes
6. ✅ Size under 500 lines (measured with line-counter.sh)
7. ✅ 80% unit test coverage achieved

## Risk Mitigation
1. **Size Risk**: Pre-identified features to defer, clear split points
2. **Integration Risk**: Test with Wave 1 components early
3. **Certificate Risk**: Provide manual fallback with --insecure

## Integration Points Summary
```
Phase 1 Components:
├── KindCertExtractor    → Used in auto_config.go
├── TrustStoreManager    → Used in initializer.go
├── CertValidator        → Used in auto_config.go
└── FallbackHandler      → Used when --insecure flag set

Wave 1 Components:
├── Builder              → Used in build_handler.go
└── GiteaClient          → Used in push_handler.go
```

## Next Steps After Implementation
1. Run line-counter.sh to verify size compliance
2. Execute unit tests to verify 80% coverage
3. Code review by Code Reviewer agent
4. Fix any identified issues
5. Commit and push to branch
6. Signal completion for integration-testing effort to begin

---

**Document Version**: 1.0
**Created**: 2025-08-30
**Created By**: Code Reviewer Agent
**State**: EFFORT_PLAN_CREATION
**Status**: READY_FOR_IMPLEMENTATION