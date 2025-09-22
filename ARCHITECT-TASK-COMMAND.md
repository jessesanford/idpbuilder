# Task Command for Architect Agent

Copy and paste this command to spawn the architect agent:

```
architect

Your task is to create a comprehensive Master Implementation Plan for adding OCI image build and push capabilities to IDPBuilder's stacks feature.

MANDATORY REQUIREMENTS:
1. Use the SF 2.0 template at templates/MASTER-IMPLEMENTATION-PLAN.md
2. Copy it to ./IMPLEMENTATION-PLAN.md and fill ALL placeholders
3. Read the detailed requirements at: ARCHITECT-PROMPT-IDPBUILDER-OCI.md

PROJECT CONTEXT:
- IDPBuilder is a kubebuilder-like system with K8s controllers on Kind cluster
- Has "stacks" feature - ArgoCD apps pulled into local Gitea
- Gitea at gitea.cnoe.localtest.me supports OCI but has self-signed certs
- Current Docker daemon needs ugly config changes to work with Gitea

SOLUTION REQUIREMENTS:
- Build containers from LOCAL Dockerfiles in stacks using Buildah Go libraries
- Push ONLY to local Gitea registry (no other registries)
- Handle self-signed certificates (InsecureSkipVerify or cert bundle)
- Manual CLI triggers only: `idpbuilder stack build [stack-name]`
- Enable ArgoCD to pull images without internet access
- 5 phases, 38-45 efforts, 6-7 weeks, ~8,000-10,000 lines

CRITICAL ARCHITECTURE FOCUS:
- Phase 1 Wave 1 MUST define ALL interfaces for stacks integration
- Study existing codebase: https://github.com/jessesanford/idpbuilder
- Focus on extending stacks feature, not generic OCI support
- Prioritize working solution with InsecureSkipVerify first

Read the full requirements document first, then create the IMPLEMENTATION-PLAN.md with all architectural decisions.
```

## Alternative Shorter Version:

```
architect

Create Master Implementation Plan for IDPBuilder stacks container build feature.

Context: IDPBuilder stacks need local container builds without Docker daemon config.

Requirements:
- Build from stacks' Dockerfiles with Buildah libs, push ONLY to gitea.cnoe.localtest.me
- Handle self-signed certs, manual CLI triggers only
- Template: cp templates/MASTER-IMPLEMENTATION-PLAN.md ./IMPLEMENTATION-PLAN.md
- Full details: ARCHITECT-PROMPT-IDPBUILDER-OCI.md
- Study codebase: https://github.com/jessesanford/idpbuilder

Critical: Focus on stacks integration, not generic OCI. Front-load ALL interfaces to Phase 1.
```