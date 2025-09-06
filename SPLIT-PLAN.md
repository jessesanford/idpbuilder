# Split Plan for CLI Commands Implementation

## Overview
**Effort**: cli-commands (Phase 2, Wave 2, E2.2.1)
**Current Size**: 800 lines (AT HARD LIMIT)
**Splits Required**: 3
**Created**: 2025-09-05 02:44:00 UTC
**Planner**: Code Reviewer Agent

## ⚠️ SPLIT INTEGRITY NOTICE ⚠️
ALL splits below belong to THIS effort ONLY: phase2/wave2/cli-commands
NO splits should reference efforts outside this path!

## Split Inventory

| Split | Description | Size Estimate | Files | Dependencies |
|-------|-------------|---------------|-------|--------------|
| 001 | Core CLI Infrastructure | ~350 lines | config, progress, flags | None |
| 002 | Build Command | ~220 lines | build cmd + tests | Split 001 |
| 003 | Push Command | ~230 lines | push cmd + tests | Split 001 |

## Detailed Split Plans

### SPLIT-001: Core CLI Infrastructure
**Branch**: `idpbuilder-oci-go-cr/phase2/wave2/cli-commands-split-001`
**Size**: ~350 lines
**Purpose**: Foundational CLI utilities and shared components

#### Files to Include:
- `pkg/cli/config.go` (173 lines)
- `pkg/cli/config_test.go` (183 lines - reduced from original)
- Total: ~356 lines (within limit)

#### Functionality:
- Configuration management with viper
- Environment variable expansion
- Config file loading/saving
- Default configuration paths

#### Dependencies:
- External: viper, testify
- Internal: None (foundational split)

#### Implementation Instructions:
1. Create split-001 directory structure
2. Implement configuration management
3. Write comprehensive tests
4. Ensure standalone compilation
5. Measure with line-counter.sh

---

### SPLIT-002: Build Command and Progress
**Branch**: `idpbuilder-oci-go-cr/phase2/wave2/cli-commands-split-002`
**Size**: ~400 lines
**Purpose**: Build command with progress reporting

#### Files to Include:
- `pkg/cmd/build/build.go` (104 lines)
- `pkg/cmd/build/build_test.go` (107 lines)
- `pkg/cli/progress.go` (270 lines - will need trimming)
- `pkg/cli/progress_test.go` (180 lines - reduced)
- Total target: ~400 lines

#### Functionality:
- `idpbuilder build` command implementation
- Progress bar/spinner functionality
- Build context validation
- Platform parsing
- Integration with Phase 1 builder

#### Dependencies:
- Split 001: Configuration utilities
- Phase 2 Wave 1: builder package
- External: cobra, go-containerregistry

#### Implementation Instructions:
1. Import Split 001 for config
2. Implement build command with cobra
3. Add progress reporting
4. Write command tests
5. Ensure size under 400 lines

---

### SPLIT-003: Push Command and Flags
**Branch**: `idpbuilder-oci-go-cr/phase2/wave2/cli-commands-split-003`
**Size**: ~250 lines
**Purpose**: Push command and shared flags

#### Files to Include:
- `pkg/cmd/push/push.go` (141 lines)
- `pkg/cmd/push/push_test.go` (188 lines - will need reduction)
- `pkg/cmd/flags.go` (52 lines)
- `pkg/cmd/root.go` (modified to ~20 lines)
- Total target: ~250 lines

#### Functionality:
- `idpbuilder push` command implementation
- Registry authentication
- Certificate trust setup
- Command integration into root
- Shared flag definitions

#### Dependencies:
- Split 001: Configuration
- Split 002: Progress reporting (imported)
- Phase 2 Wave 1: registry, certs packages
- External: cobra

#### Implementation Instructions:
1. Import Split 001 and 002
2. Implement push command
3. Add certificate handling
4. Integrate commands into root
5. Ensure proper error handling

## Split Execution Strategy

### Phase 1: Foundation (Split 001)
- Start immediately
- No dependencies
- Focus on config management
- Must pass all tests independently

### Phase 2: Build Functionality (Split 002)
- Start after Split 001 review
- Import Split 001 as dependency
- Focus on build command and progress
- Trim progress tests if needed for size

### Phase 3: Push and Integration (Split 003)
- Start after Split 002 review
- Import both previous splits
- Complete CLI integration
- Final root command assembly

## Deduplication Matrix

| Component | Split 001 | Split 002 | Split 003 |
|-----------|-----------|-----------|-----------|
| config.go | ✅ | ❌ | ❌ |
| config_test.go | ✅ | ❌ | ❌ |
| progress.go | ❌ | ✅ | ❌ |
| progress_test.go | ❌ | ✅ | ❌ |
| build.go | ❌ | ✅ | ❌ |
| build_test.go | ❌ | ✅ | ❌ |
| push.go | ❌ | ❌ | ✅ |
| push_test.go | ❌ | ❌ | ✅ |
| flags.go | ❌ | ❌ | ✅ |
| root.go | ❌ | ❌ | ✅ |

## Size Management Strategy

### Test Reduction Techniques:
1. Combine related test cases
2. Use table-driven tests more efficiently
3. Remove redundant assertions
4. Focus on critical path testing

### Code Optimization:
1. Extract common error messages
2. Consolidate similar functions
3. Use more concise variable names where appropriate
4. Remove excessive comments (keep essential ones)

## Verification Checklist
- [ ] No file appears in multiple splits
- [ ] All files from original effort covered
- [ ] Each split compiles independently
- [ ] Dependencies properly ordered
- [ ] Each split <400 lines (well under 800 limit)
- [ ] Tests maintained for critical functionality
- [ ] Integration points clearly defined

## Risk Mitigation
1. **Size Overrun**: Each split targets ~50% of limit for safety
2. **Dependency Issues**: Clear import hierarchy defined
3. **Test Coverage**: Maintain >80% coverage despite size constraints
4. **Integration**: Final split handles all command wiring

## Success Criteria
- All splits under 400 lines (50% of limit)
- Each split passes tests independently
- Final integration preserves all functionality
- No code duplication between splits
- Clean dependency chain maintained

## Next Actions
1. Orchestrator creates split-001 workspace
2. SW Engineer implements Split 001 (config)
3. Code Review of Split 001
4. Proceed sequentially through splits
5. Final integration testing after all splits complete