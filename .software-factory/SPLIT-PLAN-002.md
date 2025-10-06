# SPLIT-PLAN-002.md
## Split 002 of 3: Image Discovery & Pusher Implementation
**Planner**: Code Reviewer Agent
**Parent Effort**: E1.2.3-image-push-operations

<!-- ORCHESTRATOR METADATA PLACEHOLDER - DO NOT REMOVE -->
<!-- The orchestrator will add infrastructure metadata below: -->
<!-- WORKING_DIRECTORY, BRANCH, REMOTE, BASE_BRANCH, etc. -->
<!-- SW Engineers MUST read this metadata to navigate to the correct directory -->
<!-- END PLACEHOLDER -->

## Boundaries (CRITICAL: All splits MUST reference SAME effort!)
- **Previous Split**: Split 001 of phase1/wave2/E1.2.3-image-push-operations
  - Path: efforts/phase1/wave2/E1.2.3-image-push-operations/split-001/
  - Branch: phase1/wave2/image-push-operations-split-001
  - Summary: Implemented logging and progress reporting infrastructure
- **This Split**: Split 002 of phase1/wave2/E1.2.3-image-push-operations
  - Path: efforts/phase1/wave2/E1.2.3-image-push-operations/split-002/
  - Branch: phase1/wave2/image-push-operations-split-002
- **Next Split**: Split 003 of phase1/wave2/E1.2.3-image-push-operations
  - Path: efforts/phase1/wave2/E1.2.3-image-push-operations/split-003/
  - Branch: phase1/wave2/image-push-operations-split-003

## Files in This Split (EXCLUSIVE - no overlap with other splits)
- **pkg/push/discovery.go** (326 lines)
  - LocalImage struct and image discovery
  - Support for tarball and OCI layout formats
  - Filesystem scanning with filters
- **pkg/push/pusher.go** (363 lines)
  - ImagePusher implementation
  - Retry logic with exponential backoff
  - Registry interaction using go-containerregistry

**Total Lines**: 689 lines (within 800 line limit)

## Functionality
### Image Discovery (discovery.go)
- **LocalImage struct**: Represents discovered images
  - Name, Path, Format fields
  - v1.Image reference
- **DiscoveryOptions**: Configuration for discovery
  - BuildPath, Extensions, MaxSizeMB
  - FollowLinks support
- **Discovery Functions**:
  - DiscoverLocalImages: Main discovery entry point
  - DiscoverLocalImagesWithOptions: Configurable discovery
  - Support for .tar, .tar.gz, .tgz formats
  - OCI layout directory support
  - Image validation and loading

### Image Pusher (pusher.go)
- **ImagePusher struct**: Core pushing logic
  - Authentication handling
  - Transport configuration
  - Progress reporting integration
- **PusherOptions**: Configuration
  - MaxRetries with exponential backoff
  - Insecure registry support
  - Custom user agent
- **Push Methods**:
  - PushImage: Main push operation
  - PushImageWithContext: Context-aware pushing
  - Retry logic implementation
  - Error handling and recovery

## Dependencies
- **From Split 001**:
  - PushLogger (for logging operations)
  - ProgressReporter interface (for progress tracking)
- **External Libraries**:
  - github.com/google/go-containerregistry/pkg/v1
  - github.com/google/go-containerregistry/pkg/v1/layout
  - github.com/google/go-containerregistry/pkg/v1/tarball
  - github.com/google/go-containerregistry/pkg/v1/remote
  - github.com/google/go-containerregistry/pkg/authn
  - github.com/go-logr/logr

## Implementation Instructions
1. **Setup**:
   ```bash
   cd efforts/phase1/wave2/E1.2.3-image-push-operations
   git checkout phase1/wave2/image-push-operations-split-001
   git checkout -b phase1/wave2/image-push-operations-split-002
   ```

2. **Implement discovery.go**:
   - Define LocalImage struct with all fields
   - Implement DiscoveryOptions and defaults
   - Create discovery functions for filesystem scanning
   - Add support for tarball formats
   - Add support for OCI layout directories
   - Implement image validation logic
   - Add filtering by size and extension

3. **Implement pusher.go**:
   - Define ImagePusher struct
   - Implement PusherOptions with defaults
   - Create NewImagePusher constructor
   - Implement PushImage with retry logic
   - Add exponential backoff calculation
   - Integrate with ProgressReporter from Split 001
   - Add authentication handling
   - Implement transport configuration

4. **Integration Points**:
   - Import and use PushLogger from Split 001
   - Import and use ProgressReporter interface from Split 001
   - Wire up progress reporting in push operations
   - Add structured logging at key points

5. **Testing**:
   - Unit tests for discovery functions
   - Unit tests for pusher with mocked registry
   - Integration tests with local test images

6. **Validation**:
   ```bash
   # Ensure compilation with Split 001
   go build ./pkg/push/...

   # Run tests
   go test ./pkg/push/...

   # Measure size
   ${PROJECT_ROOT}/tools/line-counter.sh
   ```

## Split Branch Strategy
- Branch: `phase1/wave2/image-push-operations-split-002`
- Base: `phase1/wave2/image-push-operations-split-001`
- Must merge to: `phase1/wave2/image-push-operations` after review

## Success Criteria
- [ ] Both files compile without errors
- [ ] Proper integration with Split 001 interfaces
- [ ] Image discovery works for all supported formats
- [ ] Pusher implements retry with backoff correctly
- [ ] Progress reporting integrated throughout
- [ ] Implementation size <800 lines (target: ~690 lines)
- [ ] No references to operations orchestration code

## Notes for SW Engineer
- Build on top of Split 001 - don't reimplement logging/progress
- Use go-containerregistry library as specified in architecture
- Discovery should be robust to malformed images
- Pusher must handle network failures gracefully
- Exponential backoff should have jitter to avoid thundering herd
- Consider concurrent discovery for performance
- Validate images before attempting to push
- Handle both compressed and uncompressed tarballs