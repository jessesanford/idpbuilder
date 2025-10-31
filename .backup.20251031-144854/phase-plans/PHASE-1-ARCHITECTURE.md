# Phase 1 Architecture Plan
# IDPBuilder OCI Push Command - Foundation & Interfaces

**Phase**: Phase 1 - Foundation & Interfaces
**Created**: 2025-10-29
**Architect**: @agent-architect
**Fidelity Level**: **PSEUDOCODE** (high-level patterns, library choices)

---

## Executive Summary

Phase 1 establishes the complete foundation for the OCI push command through interface-first development. This phase enables **maximum parallelization** by defining all contracts upfront in Wave 1, allowing four independent implementation teams to work simultaneously in Wave 2.

**Phase Goals:**
- Define all package interfaces (Docker, Registry, Auth, TLS)
- Implement core packages using go-containerregistry library
- Enable parallel development through frozen contracts
- Establish testing patterns for subsequent phases

**Phase Outcomes:**
- All interfaces documented and compiled
- Four core packages fully implemented
- Unit test coverage ≥85%
- Foundation ready for Phase 2 command integration

---

## Adaptation Notes

### Context from Master Architecture

This is **Phase 1 of 3** in the overall implementation plan:
- **Phase 1** (This Phase): Foundation & Interfaces - 2 waves, 8 efforts
- **Phase 2**: Core Push Functionality - 3 waves, 6 efforts
- **Phase 3**: Testing & Integration - 2 waves, 4 efforts

**Master Architecture Decisions Applied:**
- Interface-first design for R307 compliance (independent mergeability)
- Incremental branching per R308 (Wave 2 builds on Wave 1 integration)
- Conservative effort sizing (400-600 lines, safety margin below 800 limit)
- go-containerregistry library mandated by user requirements
- No certificate management (only --insecure flag support)

### First Phase - No Previous Lessons

Since this is Phase 1, there are no previous phase lessons to adapt from. Future phase architecture documents will include:
- Patterns that worked well in previous phases
- Integration challenges encountered
- Performance characteristics observed
- Testing approaches that proved effective

---

## High-Level Patterns

### Core Architecture Pattern: Interface-First Modular Design

**Pattern**: Dependency Injection with Interface Contracts

**Pseudocode Structure**:
```
ARCHITECTURE InterfaceFirstModular:
  WAVE_1:
    - Define ALL interfaces upfront
    - Establish contracts between packages
    - Document expected behaviors
    - No implementations yet

  WAVE_2:
    - Implement interfaces independently
    - Four parallel development tracks
    - Each team codes against frozen interfaces
    - Integration happens at wave boundary
```

**Why This Pattern:**
- Enables parallel development (R307 compliance)
- Prevents interface changes mid-implementation
- Clear separation of concerns
- Testability through mocking
- Independent branch mergeability

### Package Structure Pattern

**Pattern**: Hexagonal Architecture (Ports and Adapters)

**Pseudocode**:
```
PACKAGE_STRUCTURE:
  Core Domain (pkg/):
    - docker/     # Port: Docker daemon operations
    - registry/   # Port: OCI registry operations
    - auth/       # Port: Authentication provision
    - tls/        # Port: TLS configuration

  Command Layer (cmd/):
    - push.go     # Adapter: CLI interface to core domain

  External Dependencies:
    - Docker daemon API (adapter)
    - go-containerregistry library (adapter)
    - OS TLS/cert pools (adapter)

FLOW:
  User Command → CLI Adapter → Core Domain (via interfaces) → External Adapters
```

**Benefits**:
- Core logic independent of external dependencies
- Easy to mock for testing
- Can swap implementations without changing core
- Clear boundaries between layers

### Error Handling Pattern

**Pattern**: Typed Errors with Context Wrapping

**Pseudocode**:
```
ERROR_HIERARCHY:
  BaseError
    ├─ DockerError
    │   ├─ ImageNotFoundError
    │   ├─ DaemonConnectionError
    │   └─ ImageInvalidError
    ├─ RegistryError
    │   ├─ AuthenticationError
    │   ├─ NetworkError
    │   └─ PushFailedError
    ├─ ValidationError
    │   └─ InvalidInputError
    └─ TLSError
        └─ CertificateError

ERROR_WRAPPING:
  FUNCTION doOperation():
    result, err = externalCall()
    IF err != nil:
      RETURN wrapError(err, "context about what failed", additionalMetadata)
```

**Implementation Strategy**:
- Use Go's error wrapping (errors.Wrap)
- Include actionable context in error messages
- Preserve error chain for debugging
- Map errors to exit codes in command layer

---

## Library Choices

### Primary OCI Library: go-containerregistry

**Choice**: `github.com/google/go-containerregistry` v0.19.0+
**Version**: Latest stable (v0.19.0 as of architecture date)

**Justification**:
- **User-mandated**: Explicitly requested in project requirements
- **Industry standard**: Used by major projects (k8s, skaffold, ko)
- **Complete OCI support**: Full implementation of OCI Image and Distribution specs
- **Well-maintained**: Active development by Google
- **Excellent documentation**: Clear API and examples

**Key Features Used**:
```
LIBRARY_USAGE go-containerregistry:
  - v1.Image interface: Represents OCI images
  - remote.Push(): Push images to registries
  - authn.Authenticator: Authentication abstraction
  - transport.Wrapper: Custom HTTP transport (for TLS config)
  - name.ParseReference(): Image reference parsing
```

**Alternatives Considered**:
- **docker/distribution**: Rejected (server-side focus, less client support)
- **Custom OCI implementation**: Rejected (too complex, reinventing wheel)

### Docker Daemon Integration

**Choice**: `github.com/docker/docker` (Docker Engine API client) v24.0.0+
**Version**: Latest stable Docker client library

**Justification**:
- **Official Docker client**: Maintained by Docker Inc.
- **Complete API coverage**: All Docker daemon operations
- **Stable API**: Well-established, backward compatible
- **Required for requirement**: Project mandates reading from Docker daemon

**Key Features Used**:
```
LIBRARY_USAGE docker-client:
  - ImageList(): List images in daemon
  - ImageInspect(): Get image details
  - ImageSave(): Export image as tar (for conversion)
  - NewClientWithOpts(): Client initialization with options
```

**Integration Point with go-containerregistry**:
```
INTEGRATION Docker-to-OCI:
  1. Use Docker client to export image as tar stream
  2. Use go-containerregistry's tarball.ImageFromPath() to convert
  3. Result is v1.Image compatible with go-containerregistry push
```

### Command Line Interface

**Choice**: `github.com/spf13/cobra` (already in IDPBuilder)
**Version**: Existing version in IDPBuilder (v1.x)

**Justification**:
- **Already integrated**: IDPBuilder uses cobra for all commands
- **Consistency**: Matches existing CLI patterns
- **No new dependency**: Zero additional overhead

**Companion**: `github.com/spf13/viper` for configuration management

### Testing Frameworks

**Unit Testing**:
- **Choice**: Go standard `testing` package
- **Mocking**: `github.com/stretchr/testify/mock`
- **Assertions**: `github.com/stretchr/testify/assert`

**Integration Testing**:
- **Choice**: Go `testing` with Docker containers
- **Registry**: Gitea container via Docker Compose
- **Docker daemon**: Host Docker socket

**Justification**:
- Standard Go tools (no exotic frameworks)
- testify widely used in Go community
- Easy CI/CD integration

---

## Conceptual Interfaces (Wave 1 Focus)

### Docker Client Interface (pkg/docker/interface.go)

**Purpose**: Abstract Docker daemon interactions for image operations

**Pseudocode Contract**:
```
INTERFACE DockerClient:
  METHOD ImageExists(ctx, imageName) -> (exists bool, error):
    PURPOSE: Check if image exists in local Docker daemon
    INPUT: context for cancellation, image name (e.g., "myapp:latest")
    OUTPUT: true if exists, false otherwise, error on daemon issues
    ERRORS: DaemonConnectionError, InvalidImageNameError

  METHOD GetImage(ctx, imageName) -> (v1.Image, error):
    PURPOSE: Retrieve image from Docker as OCI v1.Image
    INPUT: context, image name
    OUTPUT: OCI-compatible image object, error if not found
    ERRORS: ImageNotFoundError, ConversionError
    NOTES: Internally converts Docker image to OCI format

  METHOD ValidateImageName(imageName) -> error:
    PURPOSE: Validate image name follows OCI naming spec
    INPUT: image name string
    OUTPUT: nil if valid, error with details if invalid
    ERRORS: ValidationError
    VALIDATION: Check format, no special chars, valid tag syntax

  METHOD Close() -> error:
    PURPOSE: Clean up Docker client resources
    OUTPUT: error if cleanup fails
```

**Usage Example (Pseudocode)**:
```
USAGE DockerClient:
  client = NewDockerClient()
  DEFER client.Close()

  err = client.ValidateImageName("myapp:latest")
  IF err != nil: RETURN err

  exists, err = client.ImageExists(ctx, "myapp:latest")
  IF err != nil: RETURN wrapError(err, "checking image existence")
  IF NOT exists: RETURN ImageNotFoundError("myapp:latest not in daemon")

  image, err = client.GetImage(ctx, "myapp:latest")
  IF err != nil: RETURN err

  RETURN image
```

### Registry Client Interface (pkg/registry/interface.go)

**Purpose**: Abstract OCI registry push operations

**Pseudocode Contract**:
```
INTERFACE RegistryClient:
  METHOD Push(ctx, image, targetRef, progressCallback) -> error:
    PURPOSE: Push OCI image to registry with progress reporting
    INPUT:
      - ctx: context for cancellation/timeout
      - image: v1.Image (from Docker or elsewhere)
      - targetRef: fully qualified reference (registry/namespace/image:tag)
      - progressCallback: optional function for progress updates
    OUTPUT: error if push fails
    ERRORS: AuthenticationError, NetworkError, RegistryUnavailableError
    SIDE_EFFECTS: Uploads layers and manifest to remote registry

  METHOD BuildImageReference(registryURL, imageName) -> (string, error):
    PURPOSE: Construct fully qualified image reference
    INPUT:
      - registryURL: base registry URL (e.g., "https://gitea.cnoe.localtest.me:8443")
      - imageName: image with tag (e.g., "myapp:latest")
    OUTPUT: full reference (e.g., "gitea.cnoe.localtest.me:8443/giteaadmin/myapp:latest")
    ERRORS: ValidationError
    LOGIC: Parse registry, apply defaults (username=giteaadmin), combine

  METHOD ValidateRegistry(ctx, registryURL) -> error:
    PURPOSE: Verify registry is reachable (ping check)
    INPUT: registry URL
    OUTPUT: error if unreachable or invalid
    ERRORS: NetworkError, InvalidURLError

TYPE ProgressCallback = FUNCTION(ProgressUpdate):
  PURPOSE: Called during layer uploads for progress reporting

TYPE ProgressUpdate = STRUCT:
  LayerDigest:  string   # SHA256 digest of layer
  LayerSize:    int64    # Total layer size in bytes
  BytesPushed:  int64    # Bytes uploaded so far
  Status:       string   # "uploading", "complete", "exists" (layer already on registry)
```

**Usage Example (Pseudocode)**:
```
USAGE RegistryClient:
  authProvider = NewBasicAuthProvider(username, password)
  tlsConfig = NewTLSConfigProvider(insecure=true)
  client = NewRegistryClient(authProvider, tlsConfig)

  targetRef, err = client.BuildImageReference(
    "https://gitea.cnoe.localtest.me:8443",
    "myapp:latest"
  )
  IF err != nil: RETURN err

  err = client.ValidateRegistry(ctx, registryURL)
  IF err != nil: RETURN wrapError(err, "registry unreachable")

  progressFunc = FUNCTION(update):
    PRINT "Layer", update.LayerDigest, ":", update.BytesPushed, "/", update.LayerSize

  err = client.Push(ctx, image, targetRef, progressFunc)
  IF err != nil: RETURN err
```

### Authentication Provider Interface (pkg/auth/interface.go)

**Purpose**: Abstract credential management and authentication

**Pseudocode Contract**:
```
INTERFACE AuthProvider:
  METHOD GetAuthenticator() -> (authn.Authenticator, error):
    PURPOSE: Return go-containerregistry compatible authenticator
    OUTPUT: authn.Authenticator (from go-containerregistry/pkg/authn)
    ERRORS: ValidationError if credentials malformed
    NOTES: Converts internal credentials to library-specific format

  METHOD ValidateCredentials() -> error:
    PURPOSE: Pre-flight check that credentials are well-formed
    OUTPUT: error with details if invalid, nil if valid
    CHECKS:
      - Username not empty
      - Password not empty
      - No control characters in username
      - Password supports special chars (no validation, just warning)

TYPE Credentials = STRUCT:
  Username: string
  Password: string  # Supports special characters, unicode, 256+ chars

FUNCTION NewBasicAuthProvider(username, password) -> AuthProvider:
  PURPOSE: Factory for basic auth (username/password)
  RETURNS: AuthProvider instance with basic auth credentials
```

**Usage Example (Pseudocode)**:
```
USAGE AuthProvider:
  provider = NewBasicAuthProvider("giteaadmin", "myP@ss!123")

  err = provider.ValidateCredentials()
  IF err != nil: RETURN wrapError(err, "invalid credentials")

  authenticator, err = provider.GetAuthenticator()
  IF err != nil: RETURN err

  # authenticator is passed to go-containerregistry's remote.Push()
```

### TLS Configuration Provider Interface (pkg/tls/interface.go)

**Purpose**: Abstract TLS/SSL configuration for registry connections

**Pseudocode Contract**:
```
INTERFACE TLSConfigProvider:
  METHOD GetTLSConfig() -> *tls.Config:
    PURPOSE: Return tls.Config for HTTP transport
    OUTPUT: tls.Config with InsecureSkipVerify set based on mode
    NOTES:
      - Insecure mode: InsecureSkipVerify = true (skip cert validation)
      - Secure mode: Use system cert pool

  METHOD IsInsecure() -> bool:
    PURPOSE: Query whether insecure mode is enabled
    OUTPUT: true if --insecure flag set, false otherwise

TYPE TLSConfig = STRUCT:
  InsecureSkipVerify: bool  # true = bypass cert checking (curl -k equivalent)

FUNCTION NewTLSConfigProvider(insecure bool) -> TLSConfigProvider:
  PURPOSE: Factory for TLS config
  INPUT: insecure flag from CLI
  RETURNS: TLSConfigProvider instance
```

**Usage Example (Pseudocode)**:
```
USAGE TLSConfigProvider:
  tlsProvider = NewTLSConfigProvider(insecure=true)  # from --insecure flag

  IF tlsProvider.IsInsecure():
    WARN "TLS certificate verification disabled (insecure mode)"

  tlsConfig = tlsProvider.GetTLSConfig()
  # tlsConfig is passed to HTTP transport used by go-containerregistry
```

---

## Data Flow Patterns

### Complete Push Workflow (High-Level)

**Flow**: End-to-End Image Push

**Pseudocode**:
```
WORKFLOW PushImage:
  INPUT: imageName (e.g., "myapp:latest"), flags (username, password, registry, insecure)

  PHASE 1: Initialization
    dockerClient = NewDockerClient()
    DEFER dockerClient.Close()

    authProvider = NewBasicAuthProvider(flags.username, flags.password)
    tlsProvider = NewTLSConfigProvider(flags.insecure)
    registryClient = NewRegistryClient(authProvider, tlsProvider)

  PHASE 2: Input Validation
    err = dockerClient.ValidateImageName(imageName)
    IF err != nil: EXIT with ValidationError

    err = authProvider.ValidateCredentials()
    IF err != nil: EXIT with AuthenticationError

  PHASE 3: Docker Operations
    exists, err = dockerClient.ImageExists(ctx, imageName)
    IF err != nil: EXIT with DockerError
    IF NOT exists: EXIT with ImageNotFoundError

    image, err = dockerClient.GetImage(ctx, imageName)
    IF err != nil: EXIT with DockerError

  PHASE 4: Registry Preparation
    registryURL = flags.registry OR DEFAULT_REGISTRY
    targetRef, err = registryClient.BuildImageReference(registryURL, imageName)
    IF err != nil: EXIT with ValidationError

    err = registryClient.ValidateRegistry(ctx, registryURL)
    IF err != nil: EXIT with NetworkError

  PHASE 5: Push Execution
    progressFunc = FUNCTION(update):
      DisplayProgressBar(update.LayerDigest, update.BytesPushed, update.LayerSize, update.Status)

    err = registryClient.Push(ctx, image, targetRef, progressFunc)
    IF err != nil: EXIT with PushFailedError

  PHASE 6: Success
    PRINT "Successfully pushed", imageName, "to", targetRef
    EXIT 0
```

### Wave 1 Development Flow (Interface Definition Only)

**Flow**: How Wave 1 efforts are developed

**Pseudocode**:
```
WAVE_1_PROCESS:
  FOR EACH package IN [docker, registry, auth, tls]:
    EFFORT = NewEffort("Define " + package + " interface")

    STEP 1: Create interface.go file
      DEFINE all method signatures
      DOCUMENT each method with GoDoc
      DEFINE error types
      DEFINE data structures (e.g., Credentials, TLSConfig, ProgressUpdate)
      NO IMPLEMENTATIONS

    STEP 2: Create constructor signatures
      DEFINE NewClient/NewProvider functions (signatures only)
      DOCUMENT factory function purposes
      NO IMPLEMENTATIONS

    STEP 3: Add integration tests
      VERIFY interface compiles
      VERIFY method signatures are correct
      NO UNIT TESTS (nothing to test yet)

    STEP 4: Commit and push
      BRANCH: phase1-wave1-effort-{package}-interface
      COMMIT: "feat: define {package} interface for OCI push"
      SIZE CHECK: Should be ~120-200 lines

  INTEGRATION:
    BRANCH: phase1-wave1-integration
    MERGE all 4 efforts sequentially
    VERIFY: go build succeeds (all interfaces compile together)
    TAG: phase1-wave1-complete
```

### Wave 2 Development Flow (Parallel Implementation)

**Flow**: How Wave 2 efforts are developed in parallel

**Pseudocode**:
```
WAVE_2_PROCESS:
  # All efforts branch from phase1-wave1-integration
  BASE_BRANCH = "phase1-wave1-integration"

  PARALLEL_EFFORT DockerImplementation:
    BRANCH FROM: phase1-wave1-integration
    BRANCH NAME: phase1-wave2-effort-docker-impl

    IMPLEMENT pkg/docker/client.go:
      TYPE dockerClient STRUCT:
        apiClient *docker.Client  # Docker Engine API client

      IMPLEMENT ImageExists(ctx, imageName):
        images, err = apiClient.ImageList(ctx, filters)
        IF err: RETURN false, wrapError(err)
        FOR image IN images:
          IF image.Name == imageName: RETURN true, nil
        RETURN false, nil

      IMPLEMENT GetImage(ctx, imageName):
        # Export Docker image as tar
        tarStream, err = apiClient.ImageSave(ctx, [imageName])
        IF err: RETURN nil, wrapError(err)

        # Convert tar to v1.Image using go-containerregistry
        image, err = tarball.ImageFromReader(tarStream)
        IF err: RETURN nil, wrapError(err)

        RETURN image, nil

      IMPLEMENT ValidateImageName(imageName):
        # Use go-containerregistry's name.ParseReference
        _, err = name.ParseReference(imageName)
        RETURN err

      IMPLEMENT Close():
        RETURN apiClient.Close()

    ADD UNIT TESTS:
      TEST ImageExists with mocked Docker client
      TEST GetImage conversion logic
      TEST ValidateImageName with valid/invalid names
      TEST Close cleanup
      COVERAGE TARGET: 85%+

    SIZE CHECK: ~500 lines total
    COMMIT and PUSH

  PARALLEL_EFFORT RegistryImplementation:
    BRANCH FROM: phase1-wave1-integration
    BRANCH NAME: phase1-wave2-effort-registry-impl

    IMPLEMENT pkg/registry/client.go:
      TYPE registryClient STRUCT:
        authProvider AuthProvider
        tlsProvider  TLSConfigProvider
        transport    http.RoundTripper

      IMPLEMENT Push(ctx, image, targetRef, progressCallback):
        # Parse target reference
        ref, err = name.ParseReference(targetRef)
        IF err: RETURN wrapError(err)

        # Get authenticator from auth provider
        auth, err = authProvider.GetAuthenticator()
        IF err: RETURN wrapError(err)

        # Configure transport with TLS config
        transport = &http.Transport{
          TLSClientConfig: tlsProvider.GetTLSConfig(),
        }

        # Use go-containerregistry's remote.Push with progress
        err = remote.Push(ref, image,
          remote.WithAuth(auth),
          remote.WithTransport(transport),
          remote.WithProgress(progressCallback)  # Wrapper for callback
        )
        RETURN err

      IMPLEMENT BuildImageReference(registryURL, imageName):
        # Parse registry URL, apply defaults
        registry = parseRegistryHost(registryURL)  # Extract host:port
        namespace = "giteaadmin"  # Default username

        # Parse image name to extract repository and tag
        parts = parseImageName(imageName)  # e.g., "myapp:latest" -> repo="myapp", tag="latest"

        fullRef = registry + "/" + namespace + "/" + parts.repo + ":" + parts.tag
        RETURN fullRef, nil

      IMPLEMENT ValidateRegistry(ctx, registryURL):
        # Ping registry (e.g., GET /v2/)
        ref = registryURL + "/v2/"
        resp, err = http.Get(ref)
        IF err: RETURN NetworkError
        IF resp.StatusCode != 200: RETURN RegistryUnavailableError
        RETURN nil

    ADD UNIT TESTS:
      TEST Push with mocked image and registry
      TEST BuildImageReference with various inputs
      TEST ValidateRegistry with reachable/unreachable registries
      COVERAGE TARGET: 85%+

    SIZE CHECK: ~550 lines total
    COMMIT and PUSH

  PARALLEL_EFFORT AuthImplementation:
    BRANCH FROM: phase1-wave1-integration
    BRANCH NAME: phase1-wave2-effort-auth-impl

    IMPLEMENT pkg/auth/basic.go:
      TYPE basicAuthProvider STRUCT:
        username string
        password string

      IMPLEMENT GetAuthenticator():
        # Convert to go-containerregistry's authn.Authenticator
        auth = authn.FromConfig(authn.AuthConfig{
          Username: username,
          Password: password,
        })
        RETURN auth, nil

      IMPLEMENT ValidateCredentials():
        IF username == "": RETURN ValidationError("username required")
        IF password == "": RETURN ValidationError("password required")
        IF containsControlChars(username): RETURN ValidationError("username contains invalid characters")
        # Password: no validation, supports all special characters
        RETURN nil

    ADD UNIT TESTS:
      TEST GetAuthenticator returns valid authenticator
      TEST ValidateCredentials with valid/invalid inputs
      TEST special characters in password (unicode, quotes, spaces)
      TEST long passwords (256+ characters)
      COVERAGE TARGET: 90%+

    SIZE CHECK: ~300 lines total
    COMMIT and PUSH

  PARALLEL_EFFORT TLSImplementation:
    BRANCH FROM: phase1-wave1-integration
    BRANCH NAME: phase1-wave2-effort-tls-impl

    IMPLEMENT pkg/tls/config.go:
      TYPE tlsConfigProvider STRUCT:
        insecure bool

      IMPLEMENT GetTLSConfig():
        IF insecure:
          RETURN &tls.Config{
            InsecureSkipVerify: true,  # Skip certificate verification
          }
        ELSE:
          # Use system cert pool
          certPool, err = x509.SystemCertPool()
          IF err:
            WARN "Failed to load system certs, using empty pool"
            certPool = x509.NewCertPool()

          RETURN &tls.Config{
            RootCAs: certPool,
          }

      IMPLEMENT IsInsecure():
        RETURN insecure

    ADD UNIT TESTS:
      TEST GetTLSConfig in insecure mode returns InsecureSkipVerify=true
      TEST GetTLSConfig in secure mode returns system cert pool
      TEST IsInsecure returns correct value
      COVERAGE TARGET: 90%+

    SIZE CHECK: ~200 lines total
    COMMIT and PUSH

  INTEGRATION:
    BRANCH: phase1-wave2-integration
    MERGE all 4 parallel efforts
    RUN all unit tests (from all 4 packages)
    VERIFY go build succeeds
    VERIFY 85%+ test coverage achieved
    TAG: phase1-wave2-complete
```

---

## Error Handling Strategy (Detailed)

### Error Taxonomy

**Error Categories**:
```
ERROR_CATEGORIES:
  1. ValidationError: Input doesn't meet requirements
     - InvalidImageNameError
     - InvalidRegistryURLError
     - InvalidCredentialsError

  2. DockerError: Docker daemon issues
     - DaemonConnectionError (daemon not running)
     - ImageNotFoundError (image doesn't exist locally)
     - ImageConversionError (tar export failed)

  3. AuthenticationError: Credential problems
     - InvalidUsernamePasswordError
     - RegistryAuthFailedError (401/403 from registry)

  4. NetworkError: Connectivity issues
     - RegistryUnreachableError
     - TimeoutError
     - DNSResolutionError

  5. RegistryError: Registry-specific failures
     - PushFailedError (layer upload failed)
     - ManifestUploadError
     - StorageQuotaError (registry out of space)

  6. TLSError: Certificate/SSL issues
     - CertificateVerificationError (secure mode, bad cert)
     - TLSHandshakeError
```

### Error Wrapping Pattern

**Pseudocode**:
```
PATTERN ErrorWrapping:
  # Go's error wrapping: errors.Wrap(err, "context")

  FUNCTION operation():
    result, err = lowLevelCall()
    IF err != nil:
      RETURN errors.Wrap(err, "high-level context about what failed")
    RETURN result, nil

  EXAMPLE:
    image, err = dockerClient.GetImage(ctx, "myapp:latest")
    IF err != nil:
      RETURN errors.Wrap(err, "failed to retrieve image from Docker daemon")

    # Error chain: "failed to retrieve image from Docker daemon: connection refused"
```

### Error Response Format

**Pseudocode**:
```
ERROR_RESPONSE_FORMAT:
  FUNCTION formatError(err error, exitCode int):
    PRINT "Error:", getMainMessage(err)

    IF hasDetails(err):
      PRINT "Details:", getDetails(err)

    IF hasActionableSuggestion(err):
      PRINT "Suggestion:", getSuggestion(err)

    IF verboseMode:
      PRINT "Stack trace:", getStackTrace(err)

    EXIT exitCode

EXAMPLE:
  Error: Image 'myapp:latest' not found in Docker daemon
  Details: Searched local Docker images, no match found
  Suggestion: Run 'docker images' to list available images, or build the image with 'docker build -t myapp:latest .'
```

### Exit Code Mapping

**Pseudocode**:
```
EXIT_CODE_MAP:
  SUCCESS                  -> 0
  GENERAL_ERROR            -> 1  (ValidationError, unknown errors)
  AUTHENTICATION_ERROR     -> 2  (InvalidCredentialsError, RegistryAuthFailedError)
  NETWORK_ERROR            -> 3  (RegistryUnreachableError, TimeoutError, TLSError)
  IMAGE_NOT_FOUND_ERROR    -> 4  (ImageNotFoundError, DaemonConnectionError)

FUNCTION mapErrorToExitCode(err error):
  SWITCH typeOf(err):
    CASE AuthenticationError: RETURN 2
    CASE NetworkError, RegistryError, TLSError: RETURN 3
    CASE DockerError: RETURN 4
    DEFAULT: RETURN 1
```

---

## Security Considerations

### Authentication Security

**Pattern**: Secure Credential Handling

**Pseudocode**:
```
SECURITY AuthenticationHandling:
  PRINCIPLE: Never log passwords or auth tokens

  LOGGING:
    FUNCTION logAuthAttempt(username, registry):
      LOG "Authenticating user", username, "to registry", registry
      # NEVER log password or derived tokens

  ERROR_MESSAGES:
    FUNCTION sanitizeError(err error):
      message = err.String()
      IF containsPassword(message):
        message = redactPassword(message)  # Replace password with "***"
      RETURN message

  STORAGE:
    # Phase 1: No credential storage (flags only)
    # Future: OS keychain integration (out of scope)
```

### TLS Security Strategy

**Pattern**: Secure by Default, Insecure by Explicit Flag

**Pseudocode**:
```
SECURITY TLSConfiguration:
  DEFAULT_MODE: Secure (verify certificates)

  INSECURE_MODE: Only when --insecure flag explicitly set
    WARNING: Display prominent warning to user
      "WARNING: TLS certificate verification disabled. Connection is not secure."

  IMPLEMENTATION:
    IF insecureFlag:
      WARN "Insecure mode enabled - bypassing certificate verification"
      tlsConfig.InsecureSkipVerify = true
    ELSE:
      tlsConfig.InsecureSkipVerify = false
      tlsConfig.RootCAs = systemCertPool()
```

### Input Validation Security

**Pattern**: Validate All Inputs Against Injection Attacks

**Pseudocode**:
```
SECURITY InputValidation:
  THREAT: Command injection via image names or registry URLs

  VALIDATION ValidateImageName(imageName):
    # Use OCI spec parser (go-containerregistry's name.ParseReference)
    # Rejects invalid characters, shell metacharacters

    DISALLOWED_CHARS = [";", "&", "|", "`", "$", "(", ")", "<", ">", "\n"]

    FOR char IN DISALLOWED_CHARS:
      IF imageName.contains(char):
        RETURN ValidationError("Image name contains invalid character: " + char)

    # Further validation by OCI parser
    _, err = name.ParseReference(imageName)
    RETURN err

  VALIDATION ValidateRegistryURL(url):
    # Parse as URL, ensure http/https scheme
    parsed, err = url.Parse(url)
    IF err: RETURN ValidationError("Invalid URL format")

    IF parsed.Scheme NOT IN ["http", "https"]:
      RETURN ValidationError("Registry URL must use http or https")

    # Prevent SSRF by disallowing internal IPs (future enhancement)
    # For now, allow all (user controls registry URL)
```

---

## Performance Strategy (Conceptual)

### Memory Efficiency

**Strategy**: Stream Image Layers, No Full Buffering

**Pseudocode**:
```
PERFORMANCE MemoryEfficiency:
  PROBLEM: Large images (GB scale) would exhaust memory if fully buffered

  SOLUTION: Streaming approach
    WHEN pushing image:
      FOR EACH layer IN image:
        layerStream = image.GetLayerStream(layer.digest)  # Stream, not full read

        # go-containerregistry handles streaming internally
        uploadLayerStream(layerStream, registry)

        layerStream.Close()  # Release immediately

  RESULT: Memory footprint ~200MB (constant, regardless of image size)
```

### Progress Reporting Performance

**Strategy**: Throttle Progress Updates to Avoid UI Spam

**Pseudocode**:
```
PERFORMANCE ProgressReporting:
  PROBLEM: Layer uploads trigger progress callbacks very frequently (per chunk)

  SOLUTION: Throttle updates to user
    THROTTLE_INTERVAL = 100ms  # Update UI at most every 100ms

    lastUpdateTime = now()

    CALLBACK progressUpdate(update):
      currentTime = now()

      IF currentTime - lastUpdateTime > THROTTLE_INTERVAL:
        DisplayProgressToUser(update)
        lastUpdateTime = currentTime

      # Always process final "complete" updates
      IF update.Status == "complete":
        DisplayProgressToUser(update)
```

### Parallel Layer Uploads (Future Enhancement)

**Note**: go-containerregistry handles concurrency internally

**Conceptual Approach**:
```
PERFORMANCE ParallelUploads:
  # go-containerregistry already supports concurrent layer uploads
  # No custom implementation needed in Phase 1

  FUTURE: If needed, can tune concurrency:
    remote.Push(ref, image,
      remote.WithJobs(4)  # Max 4 concurrent layer uploads
    )
```

---

## Testing Strategy (High-Level)

### Unit Testing Approach

**Pattern**: Test Each Package Independently with Mocks

**Pseudocode**:
```
TESTING UnitTests:
  PACKAGE docker:
    TEST ImageExists:
      MOCK Docker API client
      SETUP: Mock returns image list with target image
      CALL: ImageExists(ctx, "myapp:latest")
      VERIFY: Returns true, no error

      SETUP: Mock returns empty image list
      CALL: ImageExists(ctx, "missing:latest")
      VERIFY: Returns false, no error

      SETUP: Mock returns connection error
      CALL: ImageExists(ctx, "myapp:latest")
      VERIFY: Returns false, DaemonConnectionError

    TEST GetImage:
      MOCK Docker API client
      SETUP: Mock ImageSave returns tar stream
      SETUP: Mock tarball conversion succeeds
      CALL: GetImage(ctx, "myapp:latest")
      VERIFY: Returns v1.Image, no error

      SETUP: Mock ImageSave fails
      CALL: GetImage(ctx, "myapp:latest")
      VERIFY: Returns nil, ImageConversionError

  PACKAGE registry:
    TEST Push:
      MOCK go-containerregistry's remote.Push
      MOCK auth provider, TLS provider
      SETUP: Mock returns success
      CALL: Push(ctx, image, "registry/repo:tag", nil)
      VERIFY: No error

      SETUP: Mock returns 401 Unauthorized
      CALL: Push(ctx, image, "registry/repo:tag", nil)
      VERIFY: Returns AuthenticationError

    TEST BuildImageReference:
      CALL: BuildImageReference("https://gitea.cnoe.localtest.me:8443", "myapp:latest")
      VERIFY: Returns "gitea.cnoe.localtest.me:8443/giteaadmin/myapp:latest"

      CALL: BuildImageReference("https://custom.com", "app:v1.0")
      VERIFY: Returns "custom.com/giteaadmin/app:v1.0"

  PACKAGE auth:
    TEST ValidateCredentials:
      provider = NewBasicAuthProvider("admin", "pass123")
      CALL: ValidateCredentials()
      VERIFY: No error

      provider = NewBasicAuthProvider("", "pass")
      CALL: ValidateCredentials()
      VERIFY: Returns ValidationError

    TEST GetAuthenticator:
      provider = NewBasicAuthProvider("user", "P@ss!$pecial")
      CALL: GetAuthenticator()
      VERIFY: Returns authn.Authenticator, no error

  PACKAGE tls:
    TEST GetTLSConfig insecure mode:
      provider = NewTLSConfigProvider(insecure=true)
      config = provider.GetTLSConfig()
      VERIFY: config.InsecureSkipVerify == true

    TEST GetTLSConfig secure mode:
      provider = NewTLSConfigProvider(insecure=false)
      config = provider.GetTLSConfig()
      VERIFY: config.InsecureSkipVerify == false
      VERIFY: config.RootCAs != nil (system cert pool loaded)
```

### Integration Testing Approach (Phase 3, Preview)

**Pattern**: Test Components Together with Real Dependencies

**Pseudocode**:
```
TESTING IntegrationTests:
  # Phase 3 Wave 1 will implement these

  TEST E2E Push Workflow:
    SETUP: Start Gitea registry in Docker container
    SETUP: Build test image and load into Docker daemon

    CALL: idpbuilder push testimage:latest --insecure --password testpass

    VERIFY: Image exists in Gitea registry
    VERIFY: Manifest uploaded correctly
    VERIFY: All layers present

    CLEANUP: Stop Gitea container, remove test image

  TEST Authentication Failure:
    SETUP: Gitea registry running
    CALL: idpbuilder push testimage:latest --password wrongpass
    VERIFY: Exit code 2 (authentication error)
    VERIFY: Error message mentions "authentication failed"

  TEST Insecure TLS Mode:
    SETUP: Gitea with self-signed cert
    CALL: idpbuilder push testimage:latest --insecure --password testpass
    VERIFY: Push succeeds despite invalid cert

    CALL: idpbuilder push testimage:latest --password testpass  # No --insecure
    VERIFY: Push fails with TLS error
```

### Test Coverage Targets

**Targets by Package**:
```
COVERAGE_TARGETS:
  pkg/docker:    85% minimum
  pkg/registry:  85% minimum
  pkg/auth:      90% minimum (critical security component)
  pkg/tls:       90% minimum (critical security component)
  cmd/push:      80% minimum (Phase 2)

  OVERALL:       85% minimum across all packages

MEASUREMENT:
  COMMAND: go test -cover ./...
  CI/CD: Fail build if coverage drops below target
```

---

## Wave Architecture Details

### Wave 1: Interface & Contract Definitions

**Objective**: Define all interfaces to freeze contracts for Wave 2 parallelization

**Efforts**:

#### Effort 1.1.1: Docker Client Interface Definition
- **File**: `pkg/docker/interface.go`
- **Size**: ~150 lines
- **Scope**:
  - DockerClient interface with 4 methods (ImageExists, GetImage, ValidateImageName, Close)
  - NewClient constructor signature
  - Error types: DaemonConnectionError, ImageNotFoundError, ImageConversionError, ValidationError
  - Comprehensive GoDoc for all methods
- **Acceptance Criteria**:
  - Compiles successfully (`go build`)
  - All methods have clear documentation
  - Error types defined
  - No implementation code (interface only)

**Pseudocode Structure**:
```
FILE pkg/docker/interface.go:
  PACKAGE docker

  IMPORT context, v1 from go-containerregistry

  COMMENT: Client interface for Docker daemon operations
  INTERFACE Client:
    METHOD ImageExists(ctx, imageName) -> (bool, error)
    METHOD GetImage(ctx, imageName) -> (v1.Image, error)
    METHOD ValidateImageName(imageName) -> error
    METHOD Close() -> error

  COMMENT: Error types for Docker operations
  TYPE DaemonConnectionError STRUCT { wraps error }
  TYPE ImageNotFoundError STRUCT { imageName string }
  TYPE ImageConversionError STRUCT { wraps error }
  TYPE ValidationError STRUCT { message string }

  COMMENT: Factory function signature
  FUNCTION NewClient() -> (Client, error)
```

#### Effort 1.1.2: Registry Client Interface Definition
- **File**: `pkg/registry/interface.go`
- **Size**: ~180 lines
- **Scope**:
  - RegistryClient interface with 3 methods (Push, BuildImageReference, ValidateRegistry)
  - ProgressCallback function type
  - ProgressUpdate struct (LayerDigest, LayerSize, BytesPushed, Status)
  - NewClient constructor signature
  - Error types: AuthenticationError, NetworkError, RegistryUnavailableError, PushFailedError
  - GoDoc comments
- **Acceptance Criteria**:
  - Compiles successfully
  - Progress reporting types well-defined
  - Integration points with auth.Provider and tls.ConfigProvider specified
  - All methods documented

**Pseudocode Structure**:
```
FILE pkg/registry/interface.go:
  PACKAGE registry

  IMPORT context, v1 from go-containerregistry

  INTERFACE Client:
    METHOD Push(ctx, image, targetRef, progressCallback) -> error
    METHOD BuildImageReference(registryURL, imageName) -> (string, error)
    METHOD ValidateRegistry(ctx, registryURL) -> error

  TYPE ProgressCallback = FUNCTION(ProgressUpdate)

  TYPE ProgressUpdate STRUCT:
    LayerDigest  string
    LayerSize    int64
    BytesPushed  int64
    Status       string  # "uploading" | "complete" | "exists"

  TYPE AuthenticationError STRUCT { wraps error }
  TYPE NetworkError STRUCT { wraps error }
  TYPE RegistryUnavailableError STRUCT { registryURL string }
  TYPE PushFailedError STRUCT { wraps error, layer string }

  FUNCTION NewClient(authProvider, tlsConfig) -> (Client, error)
```

#### Effort 1.1.3: Auth & TLS Interface Definitions
- **Files**: `pkg/auth/interface.go`, `pkg/tls/interface.go`
- **Size**: ~120 lines (combined)
- **Scope**:
  - auth.Provider interface (GetAuthenticator, ValidateCredentials)
  - auth.Credentials struct
  - tls.ConfigProvider interface (GetTLSConfig, IsInsecure)
  - tls.Config struct
  - Factory function signatures
  - Error types: ValidationError
  - GoDoc comments
- **Acceptance Criteria**:
  - Both interfaces compile
  - Compatible with go-containerregistry's authn package
  - Clear separation of concerns (auth vs TLS)
  - Integration points defined

**Pseudocode Structure**:
```
FILE pkg/auth/interface.go:
  PACKAGE auth

  IMPORT authn from go-containerregistry

  INTERFACE Provider:
    METHOD GetAuthenticator() -> (authn.Authenticator, error)
    METHOD ValidateCredentials() -> error

  TYPE Credentials STRUCT:
    Username string
    Password string

  TYPE ValidationError STRUCT { field, message string }

  FUNCTION NewBasicAuthProvider(username, password) -> Provider

FILE pkg/tls/interface.go:
  PACKAGE tls

  IMPORT crypto/tls

  INTERFACE ConfigProvider:
    METHOD GetTLSConfig() -> *tls.Config
    METHOD IsInsecure() -> bool

  TYPE Config STRUCT:
    InsecureSkipVerify bool

  FUNCTION NewConfigProvider(insecure bool) -> ConfigProvider
```

#### Effort 1.1.4: Command Structure & Flag Definitions
- **File**: `cmd/push.go` (skeleton)
- **Size**: ~200 lines
- **Scope**:
  - Cobra command definition for `idpbuilder push`
  - All flags: --registry, --username, --password, -k/--insecure, --verbose
  - Default values (registry, username)
  - Flag validation function signatures (no implementations)
  - Help text and usage examples
  - Command execution function signature (RunE)
  - No actual push logic
- **Acceptance Criteria**:
  - Command registers with cobra
  - All flags defined with correct types and defaults
  - Help text complete and accurate
  - Compiles and can be invoked (though does nothing yet)

**Pseudocode Structure**:
```
FILE cmd/push.go:
  PACKAGE cmd

  IMPORT cobra, viper, docker, registry, auth, tls

  VARIABLE pushCmd = COBRA_COMMAND:
    Use:   "push IMAGE"
    Short: "Push Docker image to OCI registry"
    Long:  "Push a Docker image from local daemon to OCI registry (default: Gitea)"
    Example: "idpbuilder push myapp:latest --password mypassword"
    Args:  cobra.ExactArgs(1)  # Require image name
    RunE:  runPushCommand  # Function signature only, implementation in Phase 2

  FLAGS:
    pushCmd.Flags().String("registry", "https://gitea.cnoe.localtest.me:8443", "Registry URL")
    pushCmd.Flags().String("username", "giteaadmin", "Registry username")
    pushCmd.Flags().String("password", "", "Registry password (REQUIRED)")
    pushCmd.Flags().BoolP("insecure", "k", false, "Skip TLS certificate verification")
    pushCmd.Flags().Bool("verbose", false, "Enable verbose output")

    pushCmd.MarkFlagRequired("password")  # Password is mandatory

  FUNCTION runPushCommand(cmd, args) -> error:
    # Signature only - implementation in Phase 2
    RETURN errors.New("not implemented yet")

  FUNCTION validateFlags() -> error:
    # Signature only

  FUNCTION init():
    rootCmd.AddCommand(pushCmd)  # Register with IDPBuilder root command
```

**Wave 1 Integration**:
```
INTEGRATION phase1-wave1-integration:
  1. Create integration branch from main
  2. Merge effort 1.1.1 (Docker interface)
  3. Merge effort 1.1.2 (Registry interface)
  4. Merge effort 1.1.3 (Auth & TLS interfaces)
  5. Merge effort 1.1.4 (Command structure)
  6. Run: go build ./...
  7. Verify: All packages compile together, no conflicts
  8. Document: Interface contracts in architecture notes
  9. Tag: phase1-wave1-complete
```

**Wave 1 Verification**:
- Total new lines: ~650 (well under 800 limit)
- All efforts independently mergeable to main
- Build succeeds (green build guaranteed - interfaces only)
- All interfaces frozen for Wave 2

---

### Wave 2: Core Package Implementations (Parallel)

**Objective**: Implement all packages in parallel using frozen Wave 1 interfaces

**Branch Strategy**: All efforts branch from `phase1-wave1-integration`

**Efforts**:

#### Effort 1.2.1: Docker Client Implementation
- **Files**: `pkg/docker/client.go`, `pkg/docker/image_validator.go`, `pkg/docker/client_test.go`
- **Size**: ~500 lines
- **Scope**:
  - Implement Client interface using Docker Engine API
  - ImageExists: List images, check for match
  - GetImage: Export as tar, convert to v1.Image using go-containerregistry
  - ValidateImageName: Use go-containerregistry's name.ParseReference
  - Close: Clean up Docker client
  - Unit tests with mocked Docker client (testify/mock)
  - Coverage: 85%+
- **Acceptance Criteria**:
  - All interface methods implemented
  - Handles Docker daemon connection errors gracefully
  - Image name validation per OCI spec
  - Unit tests pass with 85%+ coverage

**Implementation Pseudocode**:
```
FILE pkg/docker/client.go:
  PACKAGE docker

  IMPORT docker API client, go-containerregistry tarball, name

  TYPE dockerClient STRUCT:
    apiClient *docker.Client

  FUNCTION NewClient() -> (Client, error):
    client, err = docker.NewClientWithOpts(
      docker.FromEnv,                # Use DOCKER_HOST env var
      docker.WithAPIVersionNegotiation(),  # Auto-negotiate API version
    )
    IF err != nil:
      RETURN nil, DaemonConnectionError{err}

    RETURN &dockerClient{apiClient: client}, nil

  METHOD (c *dockerClient) ImageExists(ctx, imageName) -> (bool, error):
    filters = filters.NewArgs()
    filters.Add("reference", imageName)

    images, err = c.apiClient.ImageList(ctx, docker.ImageListOptions{Filters: filters})
    IF err != nil:
      RETURN false, DaemonConnectionError{err}

    RETURN len(images) > 0, nil

  METHOD (c *dockerClient) GetImage(ctx, imageName) -> (v1.Image, error):
    # Export image as tar stream
    tarStream, err = c.apiClient.ImageSave(ctx, []string{imageName})
    IF err != nil:
      IF errors.Is(err, docker.ErrImageNotFound):
        RETURN nil, ImageNotFoundError{imageName}
      RETURN nil, DaemonConnectionError{err}
    DEFER tarStream.Close()

    # Convert tar to v1.Image
    image, err = tarball.ImageFromReader(tarStream)
    IF err != nil:
      RETURN nil, ImageConversionError{err}

    RETURN image, nil

  METHOD (c *dockerClient) ValidateImageName(imageName) -> error:
    _, err = name.ParseReference(imageName)
    IF err != nil:
      RETURN ValidationError{"Invalid image name format: " + err.Error()}
    RETURN nil

  METHOD (c *dockerClient) Close() -> error:
    IF c.apiClient != nil:
      RETURN c.apiClient.Close()
    RETURN nil

FILE pkg/docker/client_test.go:
  PACKAGE docker

  IMPORT testing, testify/mock, testify/assert

  TEST TestImageExists_Found:
    mockClient = NewMockDockerClient()
    mockClient.On("ImageList", mock.Anything, mock.Anything).Return(
      []docker.ImageSummary{{ID: "abc123", RepoTags: []string{"myapp:latest"}}},
      nil,
    )

    client = &dockerClient{apiClient: mockClient}
    exists, err = client.ImageExists(context.Background(), "myapp:latest")

    ASSERT err == nil
    ASSERT exists == true

  TEST TestImageExists_NotFound:
    mockClient = NewMockDockerClient()
    mockClient.On("ImageList", mock.Anything, mock.Anything).Return(
      []docker.ImageSummary{},
      nil,
    )

    client = &dockerClient{apiClient: mockClient}
    exists, err = client.ImageExists(context.Background(), "missing:latest")

    ASSERT err == nil
    ASSERT exists == false

  TEST TestGetImage_Success:
    mockClient = NewMockDockerClient()
    mockTarStream = createMockTarStream()  # Helper to create valid tar stream
    mockClient.On("ImageSave", mock.Anything, mock.Anything).Return(mockTarStream, nil)

    client = &dockerClient{apiClient: mockClient}
    image, err = client.GetImage(context.Background(), "myapp:latest")

    ASSERT err == nil
    ASSERT image != nil

  TEST TestValidateImageName_Valid:
    client = &dockerClient{}

    err = client.ValidateImageName("myapp:latest")
    ASSERT err == nil

    err = client.ValidateImageName("registry.io/namespace/app:v1.0")
    ASSERT err == nil

  TEST TestValidateImageName_Invalid:
    client = &dockerClient{}

    err = client.ValidateImageName("invalid name with spaces")
    ASSERT err != nil

    err = client.ValidateImageName("app:tag;rm -rf /")  # Command injection attempt
    ASSERT err != nil
```

#### Effort 1.2.2: Registry Client Implementation
- **Files**: `pkg/registry/client.go`, `pkg/registry/pusher.go`, `pkg/registry/client_test.go`
- **Size**: ~550 lines
- **Scope**:
  - Implement Client interface using go-containerregistry
  - Push: Use remote.Push with auth and TLS config
  - BuildImageReference: Parse registry URL, apply defaults, construct full ref
  - ValidateRegistry: Ping /v2/ endpoint
  - Integrate with auth.Provider and tls.ConfigProvider
  - Progress callback wrapper for go-containerregistry
  - Unit tests with mocked go-containerregistry
  - Coverage: 85%+
- **Acceptance Criteria**:
  - Push workflow implemented end-to-end
  - Authentication integrated via auth.Provider
  - TLS configuration integrated via tls.ConfigProvider
  - Progress reporting functional
  - Unit tests pass with 85%+ coverage

**Implementation Pseudocode**:
```
FILE pkg/registry/client.go:
  PACKAGE registry

  IMPORT go-containerregistry remote, name, authn, v1, http

  TYPE registryClient STRUCT:
    authProvider auth.Provider
    tlsProvider  tls.ConfigProvider

  FUNCTION NewClient(authProv, tlsProv) -> (Client, error):
    RETURN &registryClient{
      authProvider: authProv,
      tlsProvider:  tlsProv,
    }, nil

  METHOD (c *registryClient) Push(ctx, image, targetRef, progressCallback) -> error:
    # Parse target reference
    ref, err = name.ParseReference(targetRef)
    IF err != nil:
      RETURN ValidationError{"Invalid target reference: " + err.Error()}

    # Get authenticator
    authenticator, err = c.authProvider.GetAuthenticator()
    IF err != nil:
      RETURN AuthenticationError{err}

    # Configure HTTP transport with TLS
    transport = &http.Transport{
      TLSClientConfig: c.tlsProvider.GetTLSConfig(),
    }

    # Prepare push options
    pushOptions = []remote.Option{
      remote.WithAuth(authenticator),
      remote.WithTransport(transport),
      remote.WithContext(ctx),
    }

    # Add progress callback if provided
    IF progressCallback != nil:
      progressHandler = wrapProgressCallback(progressCallback)
      pushOptions = append(pushOptions, remote.WithProgress(progressHandler))

    # Execute push
    err = remote.Push(ref, image, pushOptions...)
    IF err != nil:
      # Classify error
      IF isAuthError(err):
        RETURN AuthenticationError{err}
      ELSE IF isNetworkError(err):
        RETURN NetworkError{err}
      ELSE:
        RETURN PushFailedError{err, "unknown"}

    RETURN nil

  METHOD (c *registryClient) BuildImageReference(registryURL, imageName) -> (string, error):
    # Parse registry URL to extract host:port
    registryHost, err = parseRegistryHost(registryURL)
    IF err != nil:
      RETURN "", ValidationError{"Invalid registry URL: " + err.Error()}

    # Parse image name to extract repository and tag
    imageParts, err = parseImageNameAndTag(imageName)
    IF err != nil:
      RETURN "", ValidationError{"Invalid image name: " + err.Error()}

    # Apply defaults
    namespace = "giteaadmin"  # Default username for Gitea
    repository = imageParts.repo
    tag = imageParts.tag OR "latest"

    # Construct full reference
    fullRef = registryHost + "/" + namespace + "/" + repository + ":" + tag

    RETURN fullRef, nil

  METHOD (c *registryClient) ValidateRegistry(ctx, registryURL) -> error:
    # Ping registry /v2/ endpoint
    endpoint = registryURL + "/v2/"

    req, err = http.NewRequestWithContext(ctx, "GET", endpoint, nil)
    IF err != nil:
      RETURN ValidationError{"Invalid registry URL: " + err.Error()}

    # Use configured transport
    transport = &http.Transport{
      TLSClientConfig: c.tlsProvider.GetTLSConfig(),
    }
    client = &http.Client{Transport: transport}

    resp, err = client.Do(req)
    IF err != nil:
      RETURN NetworkError{err}
    DEFER resp.Body.Close()

    IF resp.StatusCode != 200 AND resp.StatusCode != 401:  # 401 is OK (auth required)
      RETURN RegistryUnavailableError{registryURL}

    RETURN nil

  FUNCTION wrapProgressCallback(callback ProgressCallback) -> chan v1.Update:
    # Convert go-containerregistry's progress format to our ProgressUpdate
    updates = make(chan v1.Update)

    GO ROUTINE:
      FOR update IN updates:
        progressUpdate = ProgressUpdate{
          LayerDigest: update.Digest.String(),
          LayerSize:   update.Total,
          BytesPushed: update.Complete,
          Status:      mapStatus(update.Status),
        }
        callback(progressUpdate)

    RETURN updates

  FUNCTION parseRegistryHost(url string) -> (string, error):
    parsed, err = url.Parse(url)
    IF err: RETURN "", err
    RETURN parsed.Host, nil  # e.g., "gitea.cnoe.localtest.me:8443"

FILE pkg/registry/client_test.go:
  PACKAGE registry

  TEST TestPush_Success:
    mockAuthProvider = NewMockAuthProvider()
    mockAuthProvider.On("GetAuthenticator").Return(authn.Anonymous, nil)

    mockTLSProvider = NewMockTLSProvider()
    mockTLSProvider.On("GetTLSConfig").Return(&tls.Config{})

    client = NewClient(mockAuthProvider, mockTLSProvider)

    mockImage = createMockImage()  # Helper to create test v1.Image

    # Note: This test requires mocking go-containerregistry's remote.Push
    # Or use integration test with real registry
    err = client.Push(context.Background(), mockImage, "registry.io/repo:tag", nil)

    ASSERT err == nil

  TEST TestBuildImageReference:
    client = &registryClient{}

    ref, err = client.BuildImageReference("https://gitea.cnoe.localtest.me:8443", "myapp:latest")
    ASSERT err == nil
    ASSERT ref == "gitea.cnoe.localtest.me:8443/giteaadmin/myapp:latest"

    ref, err = client.BuildImageReference("https://custom.io", "app:v1.0")
    ASSERT err == nil
    ASSERT ref == "custom.io/giteaadmin/app:v1.0"

    ref, err = client.BuildImageReference("https://registry.io", "notagimage")
    ASSERT err == nil
    ASSERT ref == "registry.io/giteaadmin/notagimage:latest"  # Default tag

  TEST TestValidateRegistry_Reachable:
    mockTLSProvider = NewMockTLSProvider()
    mockTLSProvider.On("GetTLSConfig").Return(&tls.Config{})

    client = &registryClient{tlsProvider: mockTLSProvider}

    # Use httptest to mock registry response
    mockServer = httptest.NewServer(http.HandlerFunc(func(w, r) {
      w.WriteHeader(200)
    }))
    DEFER mockServer.Close()

    err = client.ValidateRegistry(context.Background(), mockServer.URL)
    ASSERT err == nil
```

#### Effort 1.2.3: Authentication Implementation
- **Files**: `pkg/auth/basic.go`, `pkg/auth/validator.go`, `pkg/auth/basic_test.go`
- **Size**: ~300 lines
- **Scope**:
  - Implement Provider interface for basic auth
  - GetAuthenticator: Convert credentials to authn.Authenticator
  - ValidateCredentials: Check username/password well-formed
  - Support special characters in passwords (including quotes, unicode)
  - Password length supports 256+ characters
  - Unit tests for credential validation, special characters
  - Coverage: 90%+
- **Acceptance Criteria**:
  - Basic authentication implemented
  - Handles all special characters correctly
  - Password validation permissive (no restrictions on special chars)
  - Unit tests pass with 90%+ coverage

**Implementation Pseudocode**:
```
FILE pkg/auth/basic.go:
  PACKAGE auth

  IMPORT go-containerregistry authn, strings, unicode

  TYPE basicAuthProvider STRUCT:
    username string
    password string

  FUNCTION NewBasicAuthProvider(username, password) -> Provider:
    RETURN &basicAuthProvider{
      username: username,
      password: password,
    }

  METHOD (p *basicAuthProvider) GetAuthenticator() -> (authn.Authenticator, error):
    # Validate first
    err = p.ValidateCredentials()
    IF err != nil:
      RETURN nil, err

    # Use go-containerregistry's authn.FromConfig
    authenticator = authn.FromConfig(authn.AuthConfig{
      Username: p.username,
      Password: p.password,
    })

    RETURN authenticator, nil

  METHOD (p *basicAuthProvider) ValidateCredentials() -> error:
    IF p.username == "":
      RETURN ValidationError{"username", "Username is required"}

    IF p.password == "":
      RETURN ValidationError{"password", "Password is required"}

    # Check for control characters in username (not allowed)
    FOR _, char IN p.username:
      IF unicode.IsControl(char):
        RETURN ValidationError{"username", "Username contains invalid control characters"}

    # Password: No validation! All special characters allowed
    # Including: quotes, spaces, unicode, emoji, etc.
    # This is intentional - passwords can be anything

    RETURN nil

FILE pkg/auth/validator.go:
  PACKAGE auth

  FUNCTION ValidatePassword(password string) -> error:
    # Deliberately permissive - no restrictions
    # Just check not empty (already done in ValidateCredentials)
    RETURN nil

FILE pkg/auth/basic_test.go:
  PACKAGE auth

  TEST TestGetAuthenticator_Success:
    provider = NewBasicAuthProvider("admin", "password123")

    auth, err = provider.GetAuthenticator()

    ASSERT err == nil
    ASSERT auth != nil

  TEST TestValidateCredentials_Valid:
    provider = NewBasicAuthProvider("user", "P@ssw0rd!")
    err = provider.ValidateCredentials()
    ASSERT err == nil

    provider = NewBasicAuthProvider("admin", "simple")
    err = provider.ValidateCredentials()
    ASSERT err == nil

  TEST TestValidateCredentials_Invalid:
    provider = NewBasicAuthProvider("", "password")
    err = provider.ValidateCredentials()
    ASSERT err != nil
    ASSERT contains(err.Error(), "username")

    provider = NewBasicAuthProvider("admin", "")
    err = provider.ValidateCredentials()
    ASSERT err != nil
    ASSERT contains(err.Error(), "password")

    provider = NewBasicAuthProvider("admin\x00user", "pass")  # Control char
    err = provider.ValidateCredentials()
    ASSERT err != nil

  TEST TestSpecialCharactersInPassword:
    testCases = [
      "P@ssw0rd!#$%",
      "pass with spaces",
      "quote\"password\"here",
      "single'quotes",
      "unicode密码",
      "emoji🔒password",
      strings.Repeat("a", 256),  # Long password
      strings.Repeat("a", 1024), # Very long password
    ]

    FOR password IN testCases:
      provider = NewBasicAuthProvider("admin", password)
      err = provider.ValidateCredentials()
      ASSERT err == nil, "Failed for password: " + password

      auth, err = provider.GetAuthenticator()
      ASSERT err == nil
      ASSERT auth != nil
```

#### Effort 1.2.4: TLS Configuration Implementation
- **Files**: `pkg/tls/config.go`, `pkg/tls/config_test.go`
- **Size**: ~200 lines
- **Scope**:
  - Implement ConfigProvider interface
  - GetTLSConfig: Return tls.Config with InsecureSkipVerify or system certs
  - IsInsecure: Return insecure flag status
  - System cert pool loading with fallback
  - Unit tests for secure/insecure modes
  - Coverage: 90%+
- **Acceptance Criteria**:
  - TLS configuration provider implemented
  - Insecure mode bypasses certificate verification
  - Secure mode uses system cert pool
  - Handles system cert pool loading failures gracefully
  - Unit tests pass with 90%+ coverage

**Implementation Pseudocode**:
```
FILE pkg/tls/config.go:
  PACKAGE tls

  IMPORT crypto/tls, crypto/x509

  TYPE tlsConfigProvider STRUCT:
    insecure bool

  FUNCTION NewConfigProvider(insecure bool) -> ConfigProvider:
    RETURN &tlsConfigProvider{insecure: insecure}

  METHOD (p *tlsConfigProvider) GetTLSConfig() -> *tls.Config:
    IF p.insecure:
      # Insecure mode: skip certificate verification
      RETURN &tls.Config{
        InsecureSkipVerify: true,
      }

    # Secure mode: use system cert pool
    certPool, err = x509.SystemCertPool()
    IF err != nil:
      # Fallback: empty cert pool
      # Log warning but continue
      WARN "Failed to load system cert pool: " + err.Error()
      certPool = x509.NewCertPool()

    RETURN &tls.Config{
      RootCAs:            certPool,
      InsecureSkipVerify: false,
    }

  METHOD (p *tlsConfigProvider) IsInsecure() -> bool:
    RETURN p.insecure

FILE pkg/tls/config_test.go:
  PACKAGE tls

  TEST TestGetTLSConfig_Insecure:
    provider = NewConfigProvider(insecure=true)

    config = provider.GetTLSConfig()

    ASSERT config != nil
    ASSERT config.InsecureSkipVerify == true

  TEST TestGetTLSConfig_Secure:
    provider = NewConfigProvider(insecure=false)

    config = provider.GetTLSConfig()

    ASSERT config != nil
    ASSERT config.InsecureSkipVerify == false
    ASSERT config.RootCAs != nil  # System cert pool loaded

  TEST TestIsInsecure:
    provider = NewConfigProvider(insecure=true)
    ASSERT provider.IsInsecure() == true

    provider = NewConfigProvider(insecure=false)
    ASSERT provider.IsInsecure() == false
```

**Wave 2 Integration**:
```
INTEGRATION phase1-wave2-integration:
  1. Create integration branch from phase1-wave1-integration
  2. Merge effort 1.2.1 (Docker implementation)
  3. Merge effort 1.2.2 (Registry implementation)
  4. Merge effort 1.2.3 (Auth implementation)
  5. Merge effort 1.2.4 (TLS implementation)
  6. Run: go test ./... -cover
  7. Verify: All tests pass, coverage ≥85%
  8. Run: go build ./...
  9. Verify: All packages compile
  10. Tag: phase1-wave2-complete
```

**Wave 2 Verification**:
- Total new lines: ~1550 across 4 efforts
- All efforts independently mergeable
- Unit test coverage ≥85% achieved
- All interfaces from Wave 1 implemented
- Foundation ready for Phase 2 command integration

---

## Phase 1 Integration & Transition

### Phase 1 Complete Integration

**Integration Branch**: `phase1-integration`

**Process**:
```
PHASE_1_INTEGRATION:
  1. Create phase1-integration branch
  2. Merge phase1-wave1-integration
  3. Merge phase1-wave2-integration
  4. Run full test suite: go test ./...
  5. Verify coverage: go test -cover ./...
  6. Verify build: go build ./...
  7. Generate coverage report
  8. Tag: phase1-complete
```

**Success Criteria**:
- ✅ All interfaces defined and documented
- ✅ All packages implemented (Docker, Registry, Auth, TLS)
- ✅ Unit test coverage ≥85% overall
- ✅ All tests passing
- ✅ Build succeeds
- ✅ No uncommitted code
- ✅ Ready for Phase 2 command integration

### Handoff to Phase 2

**Phase 2 Will Build On**:
- All interfaces from pkg/ packages
- Docker client for image retrieval
- Registry client for push operations
- Authentication provider for credentials
- TLS configuration for secure/insecure mode

**Phase 2 Responsibilities** (Out of Scope for Phase 1):
- Command implementation (cmd/push.go RunE function)
- Flag processing and validation
- Orchestrating push workflow
- Progress reporting to user
- Error handling and exit codes
- Environment variable support
- Registry URL override logic

---

## Compliance Verification

### R307: Independent Branch Mergeability

**Verification**:
```
R307_COMPLIANCE:
  WAVE 1: All 4 efforts define interfaces only
    - No breaking changes possible (interfaces frozen)
    - Each effort can merge to main independently
    - Build always green (interfaces compile)

  WAVE 2: All 4 efforts implement frozen interfaces
    - No interface changes allowed
    - Each implementation independent of others
    - All efforts can merge to main independently
    - Build green (implementations satisfy interfaces)

  RESULT: ✅ COMPLIANT - Maximum parallelization achieved
```

### R308: Incremental Branching Strategy

**Verification**:
```
R308_COMPLIANCE:
  BRANCH_CHAIN:
    main
      └─ phase1-wave1-integration
           └─ phase1-wave2-integration (branches from wave1, NOT main)
                └─ phase1-integration
                     └─ phase2-wave1-integration (Phase 2 branches from Phase 1)

  INCREMENTAL_BUILD:
    - Wave 2 builds on Wave 1 (interfaces → implementations)
    - Phase 2 builds on Phase 1 (foundation → features)
    - Each wave adds functionality incrementally

  RESULT: ✅ COMPLIANT - Incremental branching maintained
```

### R359: No Code Deletion

**Verification**:
```
R359_COMPLIANCE:
  CHANGES:
    - Phase 1: ONLY adds new code (pkg/ packages, cmd/push.go skeleton)
    - No existing IDPBuilder code deleted
    - Pure additive enhancement

  RESULT: ✅ COMPLIANT - No code deletion
```

### R383: Metadata File Organization

**Verification**:
```
R383_COMPLIANCE:
  METADATA_STRUCTURE:
    .software-factory/
      └─ phase1/
           ├─ wave1/
           │    ├─ effort-docker-interface/
           │    │    └─ IMPLEMENTATION-PLAN--20251029-HHMMSS.md
           │    └─ ... (other efforts)
           └─ wave2/
                └─ ... (similar structure)

  WORKING_TREES: Clean (only code visible in pkg/, cmd/ roots)

  RESULT: ✅ COMPLIANT - All metadata in .software-factory/ with timestamps
```

### R220/R221: Size Limits

**Verification**:
```
SIZE_COMPLIANCE:
  WAVE_1:
    Effort 1.1.1: ~150 lines ✅ (under 800)
    Effort 1.1.2: ~180 lines ✅
    Effort 1.1.3: ~120 lines ✅
    Effort 1.1.4: ~200 lines ✅

  WAVE_2:
    Effort 1.2.1: ~500 lines ✅
    Effort 1.2.2: ~550 lines ✅
    Effort 1.2.3: ~300 lines ✅
    Effort 1.2.4: ~200 lines ✅

  SAFETY_MARGIN: All efforts well under 800 hard limit (200-350 line buffer)

  RESULT: ✅ COMPLIANT - Conservative effort sizing
```

### R340: Phase Architecture Fidelity

**Verification**:
```
R340_FIDELITY_COMPLIANCE:
  REQUIRED_LEVEL: PSEUDOCODE (high-level patterns, library choices)

  THIS_DOCUMENT:
    ✅ Design patterns identified (Interface-First, Hexagonal Architecture)
    ✅ Library choices justified (go-containerregistry, Docker client)
    ✅ Pseudocode examples provided (workflows, interfaces, implementations)
    ✅ NO concrete function signatures (those come in wave plans)
    ✅ NO actual code implementation (implementation in Wave 2)
    ✅ NO detailed interface definitions beyond pseudocode

  RESULT: ✅ COMPLIANT - PSEUDOCODE fidelity maintained
```

---

## Open Questions / Decisions Needed

### Decisions Made (Closed)

1. ✅ **Library Choice**: go-containerregistry (user-mandated)
2. ✅ **Interface-First**: Phase 1 Wave 1 defines all interfaces
3. ✅ **Parallelization**: 4 parallel implementations in Wave 2
4. ✅ **TLS Strategy**: Insecure flag only, no cert management
5. ✅ **Authentication**: Basic auth only (username/password)

### Pending Decisions (For Phase 2 or Later)

1. ❓ **Progress Bar Library**: Which Go library for progress display?
   - Options: cheggaaa/pb, schollz/progressbar, custom implementation
   - Decision needed in Phase 2 Wave 1 (Progress Reporting effort)

2. ❓ **Environment Variable Precedence**: Exact priority order?
   - Proposed: CLI flags > env vars > defaults
   - Confirm in Phase 2 Wave 2

3. ❓ **Verbose Mode Detail Level**: How verbose is --verbose?
   - Options: Debug logs, layer digests, HTTP traces
   - Decision needed in Phase 2 Wave 1

4. ❓ **Integration Test Infrastructure**: Docker Compose vs manual Gitea setup?
   - Options: Compose file, manual instructions, CI-specific setup
   - Decision needed in Phase 3 Wave 1

### Technical Debt / Future Enhancements (Out of Scope)

1. 🔮 **Credential Storage**: OS keychain integration (deferred to v2)
2. 🔮 **Certificate Management**: Custom cert bundles (deferred)
3. 🔮 **Multi-Registry Push**: Push to multiple registries (future)
4. 🔮 **Resume Partial Uploads**: Resume interrupted pushes (future)
5. 🔮 **Layer Deduplication**: Optimize repeated layer uploads (future)

---

## Next Steps

### Immediate Next Step: Wave 1 Implementation

**Action**: Orchestrator spawns 4 SW Engineers for Wave 1 efforts

**Branches**:
- `phase1-wave1-effort-docker-interface`
- `phase1-wave1-effort-registry-interface`
- `phase1-wave1-effort-auth-tls-interfaces`
- `phase1-wave1-effort-command-structure`

**Timeline Estimate**: 1-2 days (sequential, small efforts)

### After Wave 1: Wave 2 Implementation

**Action**: Orchestrator spawns 4 SW Engineers in parallel for Wave 2 efforts

**Branches** (all from `phase1-wave1-integration`):
- `phase1-wave2-effort-docker-impl`
- `phase1-wave2-effort-registry-impl`
- `phase1-wave2-effort-auth-impl`
- `phase1-wave2-effort-tls-impl`

**Timeline Estimate**: 2-3 days (parallel implementation)

### After Wave 2: Phase 1 Integration & Review

**Action**:
1. Orchestrator integrates all Wave 2 efforts
2. Architect reviews Phase 1 integration
3. Verify all success criteria met
4. Tag phase1-complete
5. Transition to Phase 2

**Timeline Estimate**: 1 day (integration + review)

---

## Appendix: Command Examples (User Perspective)

**Note**: These commands won't work until Phase 2 (command implementation)

### Basic Push (Default Registry)
```bash
idpbuilder push myapp:latest --password 'mypassword'
# Expected: Push to gitea.cnoe.localtest.me:8443/giteaadmin/myapp:latest
```

### Push with Custom Username
```bash
idpbuilder push myapp:latest --username developer --password 'myP@ss'
# Expected: Authenticate as 'developer' instead of default 'giteaadmin'
```

### Insecure Mode (Bypass TLS Verification)
```bash
idpbuilder push myapp:latest -k --password 'mypassword'
# Expected: Warning displayed, push proceeds without cert verification
```

### Custom Registry
```bash
idpbuilder push myapp:v1.0 --registry https://custom-registry.io --password 'pass'
# Expected: Push to custom-registry.io/giteaadmin/myapp:v1.0
```

### Verbose Mode
```bash
idpbuilder push myapp:latest --verbose --password 'pass'
# Expected: Detailed logs including layer digests, upload progress, HTTP interactions
```

### Using Environment Variables
```bash
export IDPBUILDER_REGISTRY_USERNAME=developer
export IDPBUILDER_REGISTRY_PASSWORD='complex!P@ss#123'
export IDPBUILDER_INSECURE=true
idpbuilder push myapp:latest
# Expected: Uses env vars for auth and insecure mode
```

---

## Glossary

- **OCI**: Open Container Initiative - standardizes container image formats
- **go-containerregistry**: Google's Go library for OCI registry operations
- **v1.Image**: go-containerregistry's interface representing an OCI image
- **authn.Authenticator**: go-containerregistry's authentication interface
- **Docker Engine API**: API for interacting with Docker daemon
- **Gitea**: Git service with built-in container registry
- **TLS**: Transport Layer Security (SSL)
- **InsecureSkipVerify**: TLS config option to bypass certificate verification
- **Basic Auth**: Authentication using username/password
- **Interface-First Development**: Define interfaces before implementations
- **Hexagonal Architecture**: Ports and Adapters pattern
- **Wave**: Group of related efforts within a phase
- **Effort**: Single focused implementation unit (<800 lines)
- **Integration Branch**: Branch merging all efforts in a wave
- **Pseudocode**: High-level code-like notation (not compilable)

---

## Document Status

**Status**: ✅ READY FOR IMPLEMENTATION
**Phase**: Phase 1 of 3
**Fidelity**: PSEUDOCODE (R340 compliant)
**Compliance**:
- ✅ R307: Independent Branch Mergeability (interface-first enables parallelization)
- ✅ R308: Incremental Branching Strategy (Wave 2 builds on Wave 1)
- ✅ R359: No Code Deletion (pure additive enhancement)
- ✅ R383: Metadata File Organization (all in .software-factory/)
- ✅ R220/R221: Size Limits (all efforts 150-550 lines)
- ✅ R340: Phase Architecture Fidelity (PSEUDOCODE level maintained)

**Next Action**: Orchestrator spawns SW Engineers for Phase 1 Wave 1
**Created By**: @agent-architect
**Date**: 2025-10-29

---

**END OF PHASE 1 ARCHITECTURE PLAN**
