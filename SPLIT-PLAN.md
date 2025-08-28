# Split Plan for E4.1.3-custom-contexts

## Size Analysis
- **Current Size**: 1147 lines (verified by line-counter.sh)
- **Limit**: 800 lines
- **Excess**: 347 lines
- **Required splits**: 2

## Split Architecture Overview

The current implementation has 5 main files totaling 1142 lines:
- `types.go`: 68 lines (foundational types)
- `resolver.go`: 196 lines (main resolver logic)
- `url_context.go`: 273 lines (URL fetching)
- `git_context.go`: 278 lines (Git cloning)
- `archive_context.go`: 327 lines (Archive extraction)

The split strategy divides the effort into two logical, cohesive parts:
1. **Core Context Framework** - Foundation and resolver
2. **Context Implementations** - Specific context handlers

## Split 001: Core Context Framework (Target: ~470 lines)

### Files Included
- `pkg/oci/build/contexts/types.go` (68 lines) - Interfaces and types
- `pkg/oci/build/contexts/resolver.go` (196 lines) - Main resolver logic
- `pkg/oci/build/contexts/url_context.go` (273 lines) - URL fetching implementation

**Total**: 537 lines

### Functionality
- Core Context interface definition
- ContextType enumeration and configuration
- Main ContextResolver that detects and delegates to implementations
- URL context implementation for HTTP/HTTPS sources
- LocalContextImpl (included in resolver.go)
- Cleanup coordination and error handling

### Why These Files Together
- `types.go` and `resolver.go` form the core framework
- `url_context.go` is included as it's the most commonly used context type
- URL context serves as a reference implementation for other context types
- This split provides a fully functional base that can handle local and URL contexts

### Dependencies
- None (foundational split)

### Implementation Notes
- Implement complete Context interface
- Ensure URLContextImpl is fully functional
- LocalContextImpl is already in resolver.go
- All cleanup mechanisms must work
- Include caching support for URL contexts

## Split 002: Additional Context Implementations (Target: ~605 lines)

### Files Included
- `pkg/oci/build/contexts/git_context.go` (278 lines) - Git repository support
- `pkg/oci/build/contexts/archive_context.go` (327 lines) - Archive extraction

**Total**: 605 lines

### Functionality
- GitContextImpl for cloning git repositories
- Support for branches, tags, and specific commits
- ArchiveContextImpl for extracting compressed archives
- Support for tar, tar.gz, tar.bz2, and zip formats
- Path validation and security checks
- Temporary directory management for both types

### Why These Files Together
- Both are specialized context implementations
- Both require external operations (git commands, archive extraction)
- Similar cleanup patterns (temporary directories)
- Can be tested independently from core framework

### Dependencies
- **Requires Split 001** (needs Context interface, ContextType, and ContextConfig from types.go)
- The resolver in Split 001 already has stub references to GitContextImpl and ArchiveContextImpl

### Implementation Notes
- Must implement the Context interface from Split 001
- GitContextImpl handles authentication and shallow cloning
- ArchiveContextImpl includes security validations
- Both must properly clean up temporary resources
- Integration points with resolver are already defined

## Implementation Order

### Phase 1: Split 001 - Core Context Framework
1. SW Engineer implements types, resolver, and URL context
2. Code Reviewer validates <800 lines and functionality
3. Tests for core resolver and URL fetching
4. Merge to parent branch

### Phase 2: Split 002 - Additional Context Implementations  
1. SW Engineer implements git and archive contexts
2. Imports types from Split 001
3. Code Reviewer validates <800 lines and integration
4. Tests for git cloning and archive extraction
5. Merge to parent branch

## Integration Strategy

### Split Boundaries
- **Clean Interface**: Split 002 only depends on the Context interface from Split 001
- **No Circular Dependencies**: Split 001 is self-contained
- **Stub References**: Resolver in Split 001 already references the types that Split 002 will implement

### Testing Strategy
- Split 001: Test resolver with local and URL contexts
- Split 002: Test git and archive contexts independently
- Integration: After both splits, test complete resolver with all context types

## Verification Checklist

### Split 001 Verification
- [ ] Size under 800 lines (current: 537 lines)
- [ ] Types and interfaces complete
- [ ] Resolver can detect all context types
- [ ] URL context fully functional
- [ ] Local context works
- [ ] Cleanup mechanisms operational

### Split 002 Verification  
- [ ] Size under 800 lines (current: 605 lines)
- [ ] Git context implements Context interface
- [ ] Archive context implements Context interface
- [ ] Both integrate with resolver from Split 001
- [ ] Proper cleanup of temporary resources
- [ ] Security validations in place

## Risk Mitigation

### Potential Issues
1. **Interface Changes**: If Context interface needs modification, both splits affected
   - Mitigation: Interface is well-defined and stable
   
2. **Resolver Integration**: Split 002 implementations must match resolver expectations
   - Mitigation: Resolver already has typed references (GitContextImpl, ArchiveContextImpl)

3. **Testing Dependencies**: Split 002 tests need Split 001 types
   - Mitigation: Can mock interfaces for unit tests

## Success Criteria

### Split 001 Success
- ✅ Compiles independently
- ✅ URL contexts can be fetched
- ✅ Local contexts work
- ✅ Under 800 lines
- ✅ All tests pass

### Split 002 Success
- ✅ Compiles with Split 001 imports
- ✅ Git repositories can be cloned
- ✅ Archives can be extracted
- ✅ Under 800 lines
- ✅ Integration tests pass

## Final Integration

After both splits are complete and merged:
1. Full context resolution works for all types
2. Seamless switching between context types
3. Proper cleanup on all code paths
4. Complete test coverage across all context types
5. Ready for use in OCI build system

## Notes for SW Engineers

### Split 001 Implementation
- Start with types.go as the foundation
- Implement resolver with all detection logic
- URLContextImpl needs full fetch, cache, and cleanup
- LocalContextImpl is already in resolver.go - keep it there
- Focus on making URL context robust as reference implementation

### Split 002 Implementation
- Import types from Split 001 package
- GitContextImpl needs git command execution
- Consider using go-git library or os/exec for git operations
- ArchiveContextImpl must validate paths to prevent directory traversal
- Both implementations must handle cleanup properly

This split plan ensures each part is under the 800-line limit while maintaining logical cohesion and proper dependencies.