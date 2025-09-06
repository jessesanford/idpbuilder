# Architect Agent Task: Create Master Implementation Plan for IDPBuilder OCI Build & Push Feature

## Context and Mission

You are the Software Architect for the IDPBuilder project. Your task is to create a comprehensive MASTER-IMPLEMENTATION-PLAN.md for adding OCI (Open Container Initiative) image build and push capabilities to IDPBuilder. This feature will enable IDPBuilder to build container images from local Dockerfiles within stacks and push them to the local Gitea registry, eliminating the need for Docker daemon configuration changes and enabling fully air-gapped development workflows.

## Critical Instructions

### 1. MANDATORY: Use the Software Factory 2.0 Template

You MUST use the template located at `templates/MASTER-IMPLEMENTATION-PLAN.md` as your starting point:

```bash
cp templates/MASTER-IMPLEMENTATION-PLAN.md ./IMPLEMENTATION-PLAN.md
```

Fill in ALL placeholders marked with [BRACKETS]. Do not skip any sections.

### 2. Project Requirements

**Core Functionality Required:**
- Build container images from LOCAL Dockerfiles only using Buildah (rootless, daemonless)
- Push images ONLY to the local IDPBuilder Gitea instance (gitea.cnoe.localtest.me)
- Handle Gitea's self-signed certificate properly (skip TLS verify or load cert bundle)
- Integrate with IDPBuilder's "stacks" concept for container builds
- Manual trigger for builds and pushes (user-initiated, not automatic)
- Enable ArgoCD apps to pull images from local Gitea without internet access
- Progress reporting and logging

**Technical Context:**
- **IDPBuilder Architecture**: 
  - Kubebuilder-like system with K8s controllers running on Kind cluster
  - Codebase: https://github.com/jessesanford/idpbuilder
  - Go-based tool for creating Internal Developer Platforms
  
- **Stacks System**: 
  - ArgoCD applications that get pulled into local Gitea SCM and loaded on local ArgoCD
  - Reference: https://cnoe.io/docs/tutorials/idpbuilder/idpbuilder-stack
  - Stacks can contain Dockerfiles that need building
  - Convenient for testing K8s manifests/Helm charts locally
  
- **Gitea Integration**: 
  - Always available at gitea.cnoe.localtest.me
  - Supports OCI registry but uses self-signed certificates
  - Current Docker daemon requires ugly JSON config changes and restart to work
  - This feature eliminates need for Docker daemon configuration
  
- **Buildah Requirements**:
  - MUST use Buildah Go libraries (not CLI) for direct integration into idpbuilder
  - MUST handle insecure TLS/skip verification for self-signed certs
  - OR load Gitea's cert into Buildah's cert bundle programmatically
  - Build only from local Dockerfiles (no remote Git repos)
  
- **Workflow Vision**:
  - User points to a stack containing a Dockerfile
  - IDPBuilder builds the container using embedded Buildah
  - IDPBuilder pushes to local Gitea registry (handling certs properly)
  - ArgoCD apps pull from Gitea without reaching external internet

### 3. Software Factory 2.0 Structure Requirements

You must structure the implementation using the Phase/Wave/Effort system:

- **Phases**: 5 total phases (high-level milestones)
- **Waves**: 3-5 waves per phase (groups of related work)
- **Efforts**: Individual work units, each MUST be ≤800 lines of code

**Target Metrics:**
- Total implementation: ~8,000-10,000 lines
- 38-45 total efforts across all phases
- Timeline: 6-7 weeks
- Test coverage: Phase 1: 70%, Phase 2: 80%, Phase 3: 90%, Phase 4: 85%, Phase 5: 80%

### 4. Architecture Principles for Parallelization

**CRITICAL**: Design for maximum parallelization by front-loading decisions:

#### Phase 1 MUST Define ALL:
- **Wave 1**: API contracts, interfaces, and data types
  - Stack build configuration types (linking stacks to Dockerfiles)
  - Gitea registry authentication types
  - Certificate handling types (skip TLS, cert bundle loading)
  - Image manifest types for OCI
  - Error types and codes
  - All public interfaces

- **Wave 2**: Core abstractions and patterns
  - Builder interface (wrapping Buildah Go libraries)
  - Gitea registry client interface
  - Certificate management pattern (insecure skip vs cert loading)
  - Stack-to-container mapping pattern
  - Progress reporting interface

- **Wave 3**: Integration contracts
  - CLI command structure (e.g., `idpbuilder stack build`, `idpbuilder stack push`)
  - Stack configuration schema extensions for container builds
  - Integration with existing IDPBuilder controllers/reconcilers

#### Subsequent Phases Can Then Parallelize:
- Phase 2-5 efforts can work independently once contracts are defined
- Multiple engineers can implement different registry types simultaneously
- CLI, API, and backend work can proceed in parallel

### 5. Specific Architectural Decisions to Make

In your master plan, make these HIGH-LEVEL decisions (leave implementation details for phase plans):

1. **Stacks Integration Strategy**:
   - How to extend stack configuration to include Dockerfile paths
   - Stack discovery mechanism for buildable containers
   - Relationship between stack manifests and container builds
   - Build trigger mechanism (manual CLI commands)

2. **Buildah Integration**:
   - How to embed Buildah Go libraries into IDPBuilder
   - Build context management for local Dockerfiles only
   - Rootless build execution within IDPBuilder process
   - Build cache management strategy

3. **Gitea Registry Architecture**:
   - Gitea OCI registry client implementation
   - Certificate handling approach (prefer skip TLS initially, cert bundle as enhancement)
   - Authentication with Gitea (use existing IDPBuilder credentials if available)
   - Registry URL construction (always gitea.cnoe.localtest.me)

4. **Certificate Management**:
   - Primary approach: InsecureSkipVerify for Gitea's self-signed cert
   - Alternative: Programmatic cert bundle loading from IDPBuilder's Kind cluster
   - How to avoid Docker daemon configuration changes
   - Buildah cert configuration approach

5. **CLI Design**:
   - Command structure aligned with stacks (e.g., `idpbuilder stack build [stack-name]`)
   - Flag conventions for build options
   - Manual trigger approach (no automatic builds)
   - Progress output and error reporting

6. **Testing Strategy**:
   - Unit tests with mock Buildah interfaces
   - Integration tests with test Gitea instance
   - End-to-end tests with real stacks and Dockerfiles
   - Certificate handling test scenarios

### 6. Phase Structure Guidelines

**Phase 1: Foundation & Contracts** (Week 1)
- ALL interfaces, types, and contracts for stacks integration
- Stack configuration schema extensions
- Basic build and push working with Gitea (InsecureSkipVerify)
- Buildah Go library integration foundation
- ~2000 lines total

**Phase 2: MVP Core - Stack Build & Push** (Week 2-3)
- Complete Buildah integration with local Dockerfile support
- Gitea OCI registry client implementation
- Stack discovery and build mapping
- CLI commands: `idpbuilder stack build` and `idpbuilder stack push`
- ~2500 lines total

**Phase 3: Production Readiness** (Week 4)
- Proper certificate handling (cert bundle loading option)
- Advanced error handling and retry logic
- Build optimization and caching
- Progress reporting and logging
- ~2000 lines total

**Phase 4: Enhanced Stack Features** (Week 5-6)
- Multi-stage Dockerfile support
- Build args and secrets management
- Batch stack operations
- Integration with IDPBuilder controllers
- ~2000 lines total

**Phase 5: Polish & Documentation** (Week 7)
- CLI UX improvements
- Comprehensive stacks build documentation
- Integration examples with ArgoCD apps
- Performance tuning for large images
- ~1500 lines total

### 7. Risk Factors to Address

Consider and plan mitigation for:
- Buildah Go library integration complexity and API stability
- Gitea's self-signed certificate handling across different environments
- Stack configuration breaking changes
- Build context size limitations for large Dockerfiles
- Kind cluster resource constraints when building large images
- Potential conflicts with existing IDPBuilder controllers
- Network issues between IDPBuilder and Gitea in Kind cluster

### 8. Success Criteria

Define measurable success criteria:
- Build time for typical stack application: <2 minutes
- Successful push to Gitea with self-signed cert: 100% reliability
- No Docker daemon configuration required: Zero manual setup
- ArgoCD can pull built images from Gitea: Full air-gapped workflow
- CLI integration seamless with existing `idpbuilder` commands
- Resource usage: <500MB RAM for typical builds
- Support for common Dockerfile patterns in stacks

## Deliverable Requirements

Your IMPLEMENTATION-PLAN.md must include:

1. **Complete project overview** with all sections from template
2. **Technology stack decisions** (Buildah, go-containerregistry, etc.)
3. **Detailed phase breakdown** with estimated lines per phase
4. **Dependency graph** showing parallelization opportunities
5. **Risk management** section with specific mitigations
6. **Integration strategy** for merging work
7. **Success metrics** that are measurable
8. **Resource allocation** showing how many parallel efforts possible

## Important Reminders

- **DO NOT** include low-level implementation details - save those for phase plans
- **DO NOT** plan for multiple registries - focus ONLY on Gitea at gitea.cnoe.localtest.me
- **DO NOT** include automatic build triggers - all builds are manual via CLI
- **DO** make architectural decisions that enable parallelization
- **DO** ensure every effort can be ≤800 lines by proper decomposition
- **DO** front-load ALL interface and contract definitions to Phase 1
- **DO** study the existing IDPBuilder codebase at https://github.com/jessesanford/idpbuilder
- **DO** focus on stacks integration as the primary use case
- **DO** prioritize certificate handling solution that works immediately (InsecureSkipVerify)
- **DO** plan for comprehensive testing with real stacks and Dockerfiles

## Output Format

Generate the complete IMPLEMENTATION-PLAN.md file with:
- All [PLACEHOLDERS] replaced with actual content
- Specific technology choices justified
- Clear phase/wave/effort structure
- Realistic timeline and estimates
- Comprehensive risk analysis

Remember: This master plan sets the architecture and enables parallel development. Make decisions that maximize parallelization while maintaining quality and coherence.

---

**Note**: After creating the master plan, phase plans will be created separately with detailed implementation specifications for each wave and effort. Your master plan should provide enough architectural guidance to enable this without requiring constant coordination.