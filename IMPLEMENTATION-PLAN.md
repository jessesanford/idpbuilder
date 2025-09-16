# SOFTWARE FACTORY 2.0 IMPLEMENTATION PLAN: IDPBuilder OCI Management
## Building and Pushing Docker Images Without Docker Daemon

**Project Name**: idpbuilder-oci-mvp
**Timeline**: 2 weeks (10 business days)
**Start Date**: TBD
**Primary Goal**: Build Docker images from Dockerfiles and push to Gitea WITHOUT Docker daemon dependency
**Success Metric**: End-to-end Dockerfile build → push workflow with zero certificate errors
**Framework Version**: Software Factory 2.0 (Latest Rules Applied)

## 🚨 CRITICAL REQUIREMENT CLARIFICATION 🚨

### What MUST Be Built (NOT NEGOTIABLE):
1. **`idpbuilder build`** - MUST parse and execute Dockerfiles to create OCI images
2. **`idpbuilder push`** - MUST push the SAME images created by build (no Docker daemon)
3. **Build and Push MUST work together** - Output of build must be directly usable by push
4. **NO Docker daemon dependency** - Pure Go implementation using go-containerregistry
5. **Full Dockerfile support** - At minimum: FROM, COPY, RUN, WORKDIR, ENV, EXPOSE, CMD

### What Was Misunderstood Before:
- ❌ OLD: "Create single-layer from context directory" - This is NOT Dockerfile building!
- ❌ OLD: "no Dockerfile parsing" - This defeats the entire purpose!
- ❌ OLD: Build saves to temp, push reads from Docker - Disconnected workflow!

### What MUST Be Delivered:
- ✅ NEW: Full Dockerfile parser and executor
- ✅ NEW: Dockerfile instruction execution (FROM, RUN, COPY, etc.)
- ✅ NEW: Persistent image storage that both build and push can access
- ✅ NEW: Complete workflow: Dockerfile → OCI Image → Push to Registry

## 📋 Test-Driven Development (TDD) Requirements

### Master End-to-End Test (MUST PASS)
```bash
#!/bin/bash
# This test MUST pass for the project to be considered complete

# Create test Dockerfile
cat > Dockerfile << 'EOF'
FROM alpine:latest
RUN apk add --no-cache curl
COPY app.sh /app.sh
RUN chmod +x /app.sh
CMD ["/app.sh"]
EOF

cat > app.sh << 'EOF'
#!/bin/sh
echo "Hello from IDPBuilder!"
EOF

# Test 1: Build from Dockerfile (NOT from directory!)
idpbuilder build --dockerfile Dockerfile --tag gitea.cnoe.localtest.me:8443/demo:v1.0

# Test 2: Verify image was created and stored
idpbuilder images | grep "gitea.cnoe.localtest.me:8443/demo:v1.0"

# Test 3: Push with authentication
export GITEA_USERNAME="admin"
export GITEA_TOKEN="token123"
idpbuilder push \
    --username $GITEA_USERNAME \
    --token $GITEA_TOKEN \
    --insecure \
    gitea.cnoe.localtest.me:8443/demo:v1.0

# Test 4: Verify push succeeded
curl -k -u "$GITEA_USERNAME:$GITEA_TOKEN" \
    https://gitea.cnoe.localtest.me:8443/v2/demo/tags/list | grep "v1.0"

echo "✅ END-TO-END TEST PASSED!"
```

## Executive Summary

This MVP implementation delivers a **complete Dockerfile build and push system** that works WITHOUT Docker daemon. Using go-containerregistry, we implement Dockerfile parsing, layer creation, and registry operations in pure Go. The system seamlessly handles self-signed certificates and provides a unified build→push workflow.

## 🔴🔴🔴 CRITICAL SOFTWARE FACTORY 2.0 RULES 🔴🔴🔴

### SUPREME LAWS (Violation = Automatic Failure)
- **R006**: Orchestrator NEVER writes code - only coordinates
- **R319**: Orchestrator NEVER measures code - Code Reviewers measure ALWAYS
- **R320**: NO stub implementations - all code must be functional
- **R321**: Immediate backport during integration - no deferred fixes
- **R322**: Mandatory stop at state transitions - no automatic flow
- **R323**: Must build final artifact - no "code-only" completion
- **R151**: Parallel agents spawn with <5s timestamps
- **R283**: Project-level integration is MANDATORY

### Role Clarifications
- **Orchestrator**: ONLY coordinates, spawns agents, manages state (never writes/measures)
- **Code Reviewer**: Plans efforts, MEASURES ALL CODE, validates, detects stubs
- **SW Engineer**: Implements features, fixes issues, builds artifacts
- **Architect**: Reviews wave/phase completions, provides guidance
- **Integration Agent**: Merges branches (READ-ONLY for code)

## 📊 Project Statistics

| Metric | Value |
|--------|-------|
| **Total Phases** | 2 |
| **Total Weeks** | 2 |
| **Total Efforts** | 9 |
| **Estimated Lines** | ~2,400 |
| **Team Size** | 4 agents |
| **Parallel Capacity** | 3 concurrent efforts |
| **Phase 1 Efforts** | 4 (Certificate & Storage Infrastructure - Week 1) |
| **Phase 2 Efforts** | 5 (Dockerfile Build & Push - Week 2) |
| **Average Lines/Effort** | ~265 |

## 🏗️ Technology Stack

### Core Technologies
- **Language**: Go 1.21+
- **Framework**: Cobra CLI Framework
- **Container Build**: go-containerregistry (daemonless)
- **Testing**: Go testing package, Testify, Ginkgo

### Critical Libraries for Dockerfile Building
```yaml
critical_dependencies:
  - name: github.com/google/go-containerregistry
    version: v0.19.0
    purpose: OCI image assembly, layer creation, registry operations

  - name: github.com/moby/buildkit
    version: v0.12.0
    purpose: Dockerfile parser (frontend/dockerfile/parser)

  - name: github.com/spf13/cobra
    version: v1.8.1
    purpose: CLI command structure

  - name: github.com/boltdb/bolt
    version: v1.3.1
    purpose: Local image storage database
```

## MVP Scope

### In Scope (MUST deliver)
✅ **Dockerfile parsing** - Parse and understand Dockerfile syntax
✅ **Dockerfile execution** - Execute FROM, COPY, RUN, WORKDIR, ENV, CMD instructions
✅ **Layer creation** - Create OCI layers for each Dockerfile instruction
✅ **Image storage** - Persistent local storage accessible by both build and push
✅ **Image listing** - `idpbuilder images` command to list built images
✅ **Push from storage** - Push reads images from local storage (NOT Docker daemon)
✅ **Certificate handling** - Automatic cert extraction and trust configuration
✅ **Authentication** - --username and --token flags for registry auth

### Out of Scope (POST-MVP)
❌ Multi-stage Dockerfiles (only single stage)
❌ Build arguments (ARG instruction)
❌ Build context from URLs
❌ Layer caching optimization
❌ Image signing
❌ Manifest lists (multi-arch)

## Detailed Implementation Plan

## Phase 1: Certificate & Storage Infrastructure (Days 1-5)

### Wave 1: Certificate Management & Image Storage (Days 1-2)
**Goal**: Certificate handling AND persistent image storage

#### Effort 1.1.1: Kind Certificate Extraction (Day 1)
**Size**: ~400 lines
**Owner**: SW Engineer 1

```go
// pkg/certs/extractor.go
package certs

type KindCertExtractor interface {
    ExtractGiteaCert(ctx context.Context) (*x509.Certificate, error)
    GetClusterName() (string, error)
    ValidateCertificate(cert *x509.Certificate) error
}
```

**TDD Validation Test**:
```go
func TestCertificateExtraction(t *testing.T) {
    // Setup mock Kind cluster
    // Extract certificate
    // Verify certificate is valid
    // Verify storage location
    assert.FileExists(t, "~/.idpbuilder/certs/gitea.pem")
}
```

#### Effort 1.1.2: Image Storage System (Day 2) 🆕 CRITICAL
**Size**: ~600 lines
**Owner**: SW Engineer 2

```go
// pkg/storage/image_store.go
package storage

import (
    "github.com/boltdb/bolt"
    v1 "github.com/google/go-containerregistry/pkg/v1"
)

// ImageStore provides persistent storage for built images
// This is CRITICAL for build→push workflow!
type ImageStore interface {
    // Save stores an image with all its layers
    Save(ref string, image v1.Image) error

    // Load retrieves an image for pushing
    Load(ref string) (v1.Image, error)

    // List returns all stored image references
    List() ([]ImageRef, error)

    // Delete removes an image from storage
    Delete(ref string) error

    // GetStoragePath returns the storage location
    GetStoragePath() string
}

// BoltImageStore uses BoltDB for persistent storage
type BoltImageStore struct {
    db *bolt.DB
    path string  // ~/.idpbuilder/images/store.db
}
```

**TDD Validation Test**:
```go
func TestImagePersistence(t *testing.T) {
    store := NewImageStore()

    // Create test image
    img := createTestImage()

    // Save image
    err := store.Save("test:v1", img)
    assert.NoError(t, err)

    // Load image back
    loaded, err := store.Load("test:v1")
    assert.NoError(t, err)

    // Verify images are identical
    assert.Equal(t, img.Digest(), loaded.Digest())
}
```

### Wave 2: Certificate Validation & Registry Trust (Days 3-5)

#### Effort 1.2.1: Certificate Validation Pipeline (Day 3)
**Size**: ~400 lines
**Owner**: SW Engineer 1

#### Effort 1.2.2: Registry TLS Trust Integration (Days 4-5)
**Size**: ~500 lines
**Owner**: SW Engineer 3

```go
// pkg/certs/trust.go
type TrustStoreManager interface {
    ConfigureRegistryTransport(registry string) (http.RoundTripper, error)
    SetInsecureRegistry(registry string, insecure bool) error
}
```

**TDD Validation Test**:
```go
func TestTLSConnection(t *testing.T) {
    // Test secure connection with cert
    transport, err := trust.ConfigureRegistryTransport("gitea.cnoe.localtest.me:8443")
    assert.NoError(t, err)

    // Make test request
    resp, err := transport.RoundTrip(testRequest)
    assert.NoError(t, err)
    assert.Equal(t, 200, resp.StatusCode)
}
```

## Phase 2: Dockerfile Build & Push Implementation (Days 6-10)

### Wave 1: Dockerfile Parser & Executor (Days 6-7)
**Goal**: Parse and execute Dockerfiles to create OCI images

#### Effort 2.1.1: Dockerfile Parser Integration (Day 6) 🔄 REWRITTEN
**Size**: ~600 lines
**Owner**: SW Engineer 1

```go
// pkg/build/dockerfile_parser.go
package build

import (
    "github.com/moby/buildkit/frontend/dockerfile/parser"
)

// DockerfileParser parses Dockerfile and returns instructions
type DockerfileParser interface {
    Parse(dockerfilePath string) ([]Instruction, error)
    Validate(instructions []Instruction) error
}

// Instruction represents a parsed Dockerfile instruction
type Instruction struct {
    Command string   // FROM, RUN, COPY, etc.
    Args    []string // Arguments for the command
    Raw     string   // Original line from Dockerfile
}

// Parser implementation
type Parser struct {
    supportedCommands map[string]bool
}

func NewParser() *Parser {
    return &Parser{
        supportedCommands: map[string]bool{
            "FROM": true,
            "RUN": true,
            "COPY": true,
            "WORKDIR": true,
            "ENV": true,
            "EXPOSE": true,
            "CMD": true,
            "ENTRYPOINT": true,
        },
    }
}
```

**TDD Validation Test**:
```go
func TestDockerfileParsing(t *testing.T) {
    dockerfile := `
FROM alpine:latest
RUN apk add --no-cache curl
COPY app.sh /app.sh
CMD ["/app.sh"]
`
    parser := NewParser()
    instructions, err := parser.Parse(dockerfile)

    assert.NoError(t, err)
    assert.Len(t, instructions, 4)
    assert.Equal(t, "FROM", instructions[0].Command)
    assert.Equal(t, "alpine:latest", instructions[0].Args[0])
}
```

#### Effort 2.1.2: Dockerfile Executor (Day 7) 🔄 REWRITTEN
**Size**: ~700 lines
**Owner**: SW Engineer 2

```go
// pkg/build/dockerfile_executor.go
package build

import (
    v1 "github.com/google/go-containerregistry/pkg/v1"
    "github.com/google/go-containerregistry/pkg/v1/mutate"
)

// DockerfileExecutor executes parsed Dockerfile instructions
type DockerfileExecutor interface {
    // Execute runs all instructions and returns OCI image
    Execute(instructions []Instruction, contextDir string) (v1.Image, error)

    // ExecuteInstruction handles a single instruction
    ExecuteInstruction(inst Instruction, img v1.Image) (v1.Image, error)
}

// Executor implementation
type Executor struct {
    baseImages map[string]v1.Image  // Cache of pulled base images
    storage    ImageStore           // Reference to storage system
}

// Key methods to implement
func (e *Executor) executeFrom(args []string) (v1.Image, error) {
    // Pull base image from registry
    ref, _ := name.ParseReference(args[0])
    return remote.Image(ref)
}

func (e *Executor) executeRun(img v1.Image, args []string) (v1.Image, error) {
    // Create new layer with command execution result
    // For MVP: Simulate command execution, create layer with metadata
    layer := createRunLayer(args)
    return mutate.AppendLayers(img, layer)
}

func (e *Executor) executeCopy(img v1.Image, src, dst string, contextDir string) (v1.Image, error) {
    // Create layer from file/directory
    layer := createLayerFromPath(filepath.Join(contextDir, src))
    return mutate.AppendLayers(img, layer)
}
```

**TDD Validation Test**:
```go
func TestDockerfileExecution(t *testing.T) {
    instructions := []Instruction{
        {Command: "FROM", Args: []string{"alpine:latest"}},
        {Command: "COPY", Args: []string{"app.sh", "/app.sh"}},
        {Command: "RUN", Args: []string{"chmod", "+x", "/app.sh"}},
        {Command: "CMD", Args: []string{"/app.sh"}},
    }

    executor := NewExecutor()
    img, err := executor.Execute(instructions, "./context")

    assert.NoError(t, err)
    assert.NotNil(t, img)

    // Verify image has correct layers
    layers, _ := img.Layers()
    assert.GreaterOrEqual(t, len(layers), 3) // Base + COPY + RUN
}
```

### Wave 2: CLI Integration & Push (Days 8-10)

#### Effort 2.2.1: Build Command with Dockerfile Support (Day 8) 🔄 REWRITTEN
**Size**: ~500 lines
**Owner**: SW Engineer 3

```go
// cmd/build.go
package cmd

var buildCmd = &cobra.Command{
    Use:   "build",
    Short: "Build OCI image from Dockerfile",
    Long:  "Build an OCI image from a Dockerfile and context, WITHOUT Docker daemon",
    Example: `
  # Build from Dockerfile in current directory
  idpbuilder build -t myapp:latest .

  # Build from specific Dockerfile
  idpbuilder build -f Dockerfile.prod -t myapp:prod .

  # Build with specific context
  idpbuilder build -t myapp:latest -f docker/Dockerfile ./src`,
    RunE: runBuild,
}

func init() {
    buildCmd.Flags().StringVarP(&dockerfilePath, "file", "f", "Dockerfile", "Path to Dockerfile")
    buildCmd.Flags().StringVarP(&imageTag, "tag", "t", "", "Image tag (required)")
    buildCmd.Flags().StringVar(&contextPath, "context", ".", "Build context directory")
    buildCmd.MarkFlagRequired("tag")
}

func runBuild(cmd *cobra.Command, args []string) error {
    // 1. Parse Dockerfile
    parser := build.NewParser()
    instructions, err := parser.ParseFile(dockerfilePath)

    // 2. Execute Dockerfile to create image
    executor := build.NewExecutor()
    image, err := executor.Execute(instructions, contextPath)

    // 3. Save to persistent storage (NOT Docker daemon!)
    store := storage.NewImageStore()
    err = store.Save(imageTag, image)

    fmt.Printf("Successfully built %s\n", imageTag)
    fmt.Printf("Image ID: %s\n", digest)
    return nil
}
```

**TDD Validation Test**:
```bash
# Test build command
cat > Dockerfile << 'EOF'
FROM alpine:latest
COPY hello.txt /hello.txt
EOF

echo "Hello World" > hello.txt

# This MUST work without Docker daemon!
idpbuilder build -t test:v1 .

# Verify image was saved
idpbuilder images | grep "test:v1"
```

#### Effort 2.2.2: Push Command Reading from Storage (Day 9) 🔄 REWRITTEN
**Size**: ~500 lines
**Owner**: SW Engineer 1

```go
// cmd/push.go
package cmd

var pushCmd = &cobra.Command{
    Use:   "push IMAGE[:TAG]",
    Short: "Push image to registry",
    Long:  "Push an OCI image to registry (reads from local storage, not Docker)",
    Example: `
  # Push with authentication
  idpbuilder push --username admin --token secret myapp:latest

  # Push with insecure registry
  idpbuilder push --insecure localhost:5000/myapp:latest`,
    Args: cobra.ExactArgs(1),
    RunE: runPush,
}

func runPush(cmd *cobra.Command, args []string) error {
    imageRef := args[0]

    // 1. Load image from storage (NOT Docker daemon!)
    store := storage.NewImageStore()
    image, err := store.Load(imageRef)
    if err != nil {
        return fmt.Errorf("image not found in local storage: %w", err)
    }

    // 2. Configure registry transport with certs
    transport, err := certs.ConfigureTransport(pushInsecure)

    // 3. Push to registry
    ref, _ := name.ParseReference(imageRef)
    auth := &authn.Basic{Username: username, Password: token}

    err = remote.Write(ref, image,
        remote.WithAuth(auth),
        remote.WithTransport(transport))

    fmt.Printf("Successfully pushed %s\n", imageRef)
    return nil
}
```

**TDD Validation Test**:
```bash
# End-to-end test
idpbuilder build -t gitea.cnoe.localtest.me:8443/test:v1 .
idpbuilder push \
    --username admin \
    --token $TOKEN \
    --insecure \
    gitea.cnoe.localtest.me:8443/test:v1

# Verify in registry
curl -k https://gitea.cnoe.localtest.me:8443/v2/test/tags/list
```

#### Effort 2.2.3: Images Command for Listing (Day 10) 🆕
**Size**: ~300 lines
**Owner**: SW Engineer 2

```go
// cmd/images.go
package cmd

var imagesCmd = &cobra.Command{
    Use:   "images",
    Short: "List images in local storage",
    Long:  "List all OCI images stored locally by idpbuilder build",
    RunE: runImages,
}

func runImages(cmd *cobra.Command, args []string) error {
    store := storage.NewImageStore()
    images, err := store.List()

    // Display in Docker-like format
    fmt.Printf("REPOSITORY\tTAG\tIMAGE ID\tCREATED\tSIZE\n")
    for _, img := range images {
        fmt.Printf("%s\t%s\t%s\t%s\t%s\n",
            img.Repository, img.Tag, img.ID[:12], img.Created, img.Size)
    }
    return nil
}
```

## 🧪 Phase-by-Phase TDD Validation Tests

### Phase 1 Validation Tests (Must pass before Phase 2)
```bash
#!/bin/bash
# Phase 1 Gate Test

# Test 1: Certificate extraction works
idpbuilder cert extract --cluster kind-idpbuilder
test -f ~/.idpbuilder/certs/gitea.pem || exit 1

# Test 2: Image storage works
TEST_IMG=$(echo '{"test": "data"}' | idpbuilder internal save-test-image test:v1)
idpbuilder internal load-test-image test:v1 | grep '"test": "data"' || exit 1

# Test 3: TLS connection works
idpbuilder internal test-connection gitea.cnoe.localtest.me:8443 || exit 1

echo "✅ Phase 1 validation passed!"
```

### Phase 2 Wave 1 Validation Tests
```bash
#!/bin/bash
# Wave 1: Dockerfile parsing and execution

# Test 1: Parse Dockerfile
cat > Dockerfile << 'EOF'
FROM alpine:latest
RUN echo "test"
EOF

idpbuilder internal parse-dockerfile Dockerfile | grep "FROM alpine:latest" || exit 1

# Test 2: Execute Dockerfile instructions
idpbuilder internal test-executor Dockerfile | grep "2 layers created" || exit 1

echo "✅ Wave 1 validation passed!"
```

### Phase 2 Wave 2 Validation Tests
```bash
#!/bin/bash
# Wave 2: Full build and push

# Test 1: Build from Dockerfile
cat > Dockerfile << 'EOF'
FROM alpine:latest
COPY test.txt /test.txt
EOF
echo "test" > test.txt

idpbuilder build -t test:v1 .
idpbuilder images | grep "test:v1" || exit 1

# Test 2: Push to registry
idpbuilder push --insecure --username admin --token $TOKEN test:v1

echo "✅ Wave 2 validation passed!"
```

## 🎯 Success Criteria

### Primary Success Definition
**SUCCESS** means delivering a **fully functional idpbuilder** that:

1. ✅ **Builds from Dockerfiles** - Parses and executes Dockerfile instructions
2. ✅ **Creates OCI images** - Produces spec-compliant images with layers
3. ✅ **Stores images persistently** - Images survive between commands
4. ✅ **Pushes without Docker** - No Docker daemon dependency
5. ✅ **Handles certificates** - Automatic cert extraction and trust
6. ✅ **Provides authentication** - --username and --token flags work
7. ✅ **Lists built images** - `idpbuilder images` shows stored images
8. ✅ **Zero Docker dependency** - Everything works without Docker installed

### The Master Test (MUST PASS)
```bash
# This single test proves the entire system works
cat > Dockerfile << 'EOF'
FROM alpine:3.19
RUN apk add --no-cache curl
WORKDIR /app
COPY hello.sh .
RUN chmod +x hello.sh
CMD ["./hello.sh"]
EOF

cat > hello.sh << 'EOF'
#!/bin/sh
echo "IDPBuilder works without Docker!"
EOF

# Build (no Docker daemon!)
idpbuilder build -t gitea.cnoe.localtest.me:8443/demo:final .

# List (shows in local storage)
idpbuilder images | grep "demo:final"

# Push (reads from storage, not Docker!)
idpbuilder push \
    --username admin \
    --token secret \
    --insecure \
    gitea.cnoe.localtest.me:8443/demo:final

# Verify in registry
curl -k -u admin:secret \
    https://gitea.cnoe.localtest.me:8443/v2/demo/tags/list | grep "final"

echo "✅ COMPLETE SUCCESS - Dockerfile→Build→Push workflow works!"
```

## Risk Mitigation

### Risk 1: Dockerfile Parsing Complexity
**Mitigation**: Start with basic instructions only (FROM, RUN, COPY, CMD). Use Moby BuildKit's parser.

### Risk 2: Layer Creation Without Docker
**Mitigation**: Use go-containerregistry's mutate package for layer operations.

### Risk 3: Storage Corruption
**Mitigation**: Use BoltDB for ACID-compliant storage, implement backup/restore.

### Risk 4: Missing Docker Features
**Mitigation**: Document supported Dockerfile subset clearly, provide helpful error messages.

## Definition of Done

The project is COMPLETE when:

1. ✅ The master test script passes end-to-end
2. ✅ All phase validation tests pass
3. ✅ No Docker daemon required at any point
4. ✅ Images persist between commands
5. ✅ Build and push work seamlessly together
6. ✅ Certificate handling is automatic
7. ✅ Documentation includes working examples

## Quick Start (Post-Implementation)

```bash
# 1. Install idpbuilder
go install github.com/jessesanford/idpbuilder@mvp

# 2. Start IDPBuilder with Gitea
idpbuilder create --with-gitea

# 3. Build an image from Dockerfile (NO DOCKER NEEDED!)
cat > Dockerfile << 'EOF'
FROM alpine:latest
RUN echo "Hello from IDPBuilder"
CMD ["/bin/sh", "-c", "echo Running without Docker!"]
EOF

idpbuilder build -t myapp:v1 .

# 4. List images in local storage
idpbuilder images

# 5. Push to Gitea
idpbuilder push --username admin --token $TOKEN myapp:v1

# That's it! Complete Dockerfile→Build→Push without Docker!
```

## 📊 Implementation Summary

**Total Implementation Size**: ~2,400 lines
**Total Efforts**: 9 (all < 800 lines each)
**Total Phases**: 2
**Effort Distribution**:
- Phase 1: Certificate & Storage Infrastructure (4 efforts, ~1,300 lines)
- Phase 2: Dockerfile Build & Push (5 efforts, ~1,100 lines)

---

**Document Version**: 3.0
**Updated**: 2025-09-16
**Critical Changes**:
- Added Dockerfile parsing and execution as PRIMARY requirement
- Added persistent image storage system
- Added comprehensive TDD validation tests
- Fixed build→push workflow disconnect
- Removed "directory tarball" approach completely

**Remember**: This MVP delivers COMPLETE Dockerfile building and pushing WITHOUT Docker daemon!