# 🔴🔴🔴 RULE R374: Pre-Planning Research Protocol (SUPREME LAW) 🔴🔴🔴

## Category
SUPREME LAW - Planning Phase Requirement

## Criticality
BLOCKING - Violation = -50% grading penalty

## Description
Code reviewers and architects MUST research ALL existing implementations across the entire effort repository BEFORE creating any plans. This research MUST be documented in effort plans and used to prevent duplicate implementations.

## Rationale
Without researching existing code:
1. Duplicate implementations proliferate
2. Incompatible APIs emerge
3. Integration becomes impossible
4. Architecture degrades into chaos
5. Massive rework becomes necessary

## Requirements

### Phase 1: Comprehensive Research (BEFORE ANY PLANNING)

#### 1.1 Current Wave Analysis
```bash
# For current wave branches
for branch in $(git branch -r | grep "phase${PHASE}-wave${WAVE}"); do
    echo "Analyzing $branch..."
    git checkout $branch

    # Find all interfaces
    grep -r "type.*interface" --include="*.go" . > interfaces.txt

    # Find all major functions
    grep -r "^func " --include="*.go" . > functions.txt

    # Find all structs
    grep -r "^type.*struct" --include="*.go" . > structs.txt
done
```

#### 1.2 Previous Wave Analysis
```bash
# Check ALL previous waves
for wave in $(seq 1 $((CURRENT_WAVE - 1))); do
    for branch in $(git branch -r | grep "phase${PHASE}-wave${wave}"); do
        git checkout $branch
        # Same analysis as above
    done
done
```

#### 1.3 Previous Phase Analysis
```bash
# Check ALL previous phases
for phase in $(seq 1 $((CURRENT_PHASE - 1))); do
    for branch in $(git branch -r | grep "phase${phase}"); do
        git checkout $branch
        # Same analysis as above
    done
done
```

### Phase 2: Documentation Requirements

#### 2.1 Effort Plan MUST Include
```markdown
## Pre-Planning Research Results (MANDATORY - R374)

### Existing Interfaces Found
| Interface | Location | Signature | Must Implement |
|-----------|----------|-----------|----------------|
| Registry | branch4/registry/types.go | Push(ctx, image string, content io.Reader) error | YES |
| Storage | branch2/storage/interface.go | Store(key string, data []byte) error | YES |

### Existing Implementations to Reuse
| Component | Location | Purpose | How to Use |
|-----------|----------|---------|------------|
| ImageParser | branch3/parser/image.go | Parses OCI images | Import and use directly |
| ConfigLoader | branch1/config/loader.go | Loads YAML configs | Import, don't recreate |

### APIs Already Defined
| API | Method | Signature | Notes |
|-----|--------|-----------|-------|
| RegistryAPI | Push | (ctx, image, content) | MUST use this exact signature |
| StorageAPI | Get | (ctx, key) ([]byte, error) | DO NOT create alternative |

### Forbidden Duplications
- DO NOT create new Registry interface (exists in branch4)
- DO NOT create new image parsing logic (exists in branch3)
- DO NOT implement alternative Push method signatures
- DO NOT create competing storage abstractions

### Required Integrations
- MUST implement Registry interface from branch4
- MUST use ImageParser from branch3
- MUST import ConfigLoader from branch1
```

### Phase 3: Validation Requirements

#### 3.1 Planning Validation Checklist
- [ ] Searched current wave for existing code
- [ ] Searched previous waves for existing code
- [ ] Searched previous phases for existing code
- [ ] Documented all found interfaces
- [ ] Documented all reusable implementations
- [ ] Listed forbidden duplications
- [ ] Specified required integrations

#### 3.2 Implementation Instructions
```markdown
## Implementation Requirements (from R374 Research)

### MUST USE These Existing Components:
1. Registry interface from effort-4 (branch: phase2-wave1-effort-4)
   - Import: `import "github.com/project/effort-repo/branch4/registry"`
   - Implement: All Registry interface methods with EXACT signatures

2. ImageParser from effort-3
   - Import: `import "github.com/project/effort-repo/branch3/parser"`
   - Use: `parser.ParseImage(imageData)`

### FORBIDDEN - Do Not Create:
1. New Registry interface (R373 violation)
2. New image parsing functions (use ImageParser)
3. Alternative Push/Pull method signatures
```

## Enforcement Mechanisms

### Code Reviewer Responsibilities
1. **BEFORE creating effort plan**:
   - Execute research protocol
   - Document ALL findings
   - Include in effort plan

2. **IN effort plan**:
   - Mandatory "Pre-Planning Research Results" section
   - List all interfaces to implement
   - List all code to reuse
   - List forbidden duplications

3. **DURING review**:
   - Verify no duplicates created
   - Confirm interfaces implemented correctly
   - Check that existing code was reused

### SW Engineer Responsibilities
1. **READ research results in plan**
2. **IMPORT specified existing packages**
3. **IMPLEMENT specified interfaces exactly**
4. **REUSE specified components**
5. **NEVER create alternatives**

### Architect Responsibilities
1. **VERIFY research was conducted**
2. **VALIDATE no duplicates exist**
3. **CONFIRM interfaces properly implemented**
4. **REJECT if existing code not reused**

## Validation Scripts

### Pre-Planning Research Script
```bash
#!/bin/bash
# R374-pre-planning-research.sh

EFFORT_REPO="/path/to/effort-repo"
OUTPUT_FILE="pre-planning-research.md"

echo "# Pre-Planning Research Results" > $OUTPUT_FILE
echo "Generated: $(date)" >> $OUTPUT_FILE
echo "" >> $OUTPUT_FILE

cd $EFFORT_REPO

echo "## Existing Interfaces" >> $OUTPUT_FILE
find . -name "*.go" -exec grep -l "type.*interface" {} \; | while read file; do
    echo "### $file" >> $OUTPUT_FILE
    grep "type.*interface" $file -A 10 >> $OUTPUT_FILE
    echo "" >> $OUTPUT_FILE
done

echo "## Existing Major Types" >> $OUTPUT_FILE
find . -name "*.go" -exec grep "^type.*struct" {} \; | head -20 >> $OUTPUT_FILE

echo "## Key Functions Found" >> $OUTPUT_FILE
for keyword in Push Pull Upload Download Store Retrieve Create Delete; do
    echo "### ${keyword} functions:" >> $OUTPUT_FILE
    grep -r "func.*${keyword}(" --include="*.go" . | head -5 >> $OUTPUT_FILE
    echo "" >> $OUTPUT_FILE
done

echo "Research complete. Results in $OUTPUT_FILE"
```

### Validation During Planning
```bash
#!/bin/bash
# R374-validate-plan.sh

PLAN_FILE=$1

echo "Validating R374 compliance in plan..."

# Check for mandatory section
if ! grep -q "Pre-Planning Research Results" $PLAN_FILE; then
    echo "ERROR: Missing mandatory 'Pre-Planning Research Results' section"
    exit 1
fi

# Check for interface documentation
if ! grep -q "Existing Interfaces Found" $PLAN_FILE; then
    echo "ERROR: Missing interface research documentation"
    exit 1
fi

# Check for reuse documentation
if ! grep -q "Existing Implementations to Reuse" $PLAN_FILE; then
    echo "ERROR: Missing reuse documentation"
    exit 1
fi

# Check for forbidden duplications
if ! grep -q "Forbidden Duplications" $PLAN_FILE; then
    echo "ERROR: Missing forbidden duplications list"
    exit 1
fi

echo "Plan passes R374 validation"
```

## Example: Correct Research Documentation

```markdown
# Effort Plan: Implement Container Upload Feature

## Pre-Planning Research Results (R374)

### Research Conducted
- Searched phase2-wave1-* branches: Found Registry interface
- Searched phase2-wave2-* branches: Found client implementations
- Searched phase1-* branches: Found authentication utilities

### Existing Interfaces Found
1. **Registry Interface** (phase2-wave1-effort-4/pkg/registry/interface.go)
   ```go
   type Registry interface {
       Push(ctx context.Context, image string, content io.Reader) error
       Pull(ctx context.Context, image string) (io.Reader, error)
   }
   ```
   **ACTION**: MUST implement this interface exactly

2. **Authenticator Interface** (phase1-wave3-effort-2/pkg/auth/types.go)
   ```go
   type Authenticator interface {
       Authenticate(ctx context.Context, creds Credentials) error
   }
   ```
   **ACTION**: MUST use for authentication

### Existing Implementations to Reuse
1. **BasicAuthenticator** (phase1-wave3-effort-2/pkg/auth/basic.go)
   - Provides basic auth implementation
   - **ACTION**: Import and use, don't recreate

2. **ImageValidator** (phase2-wave1-effort-3/pkg/validate/image.go)
   - Validates OCI image format
   - **ACTION**: Import and use for validation

### Forbidden Duplications
- DO NOT create new Registry interface
- DO NOT create alternative Push/Pull signatures
- DO NOT implement new authentication (use BasicAuthenticator)
- DO NOT write image validation (use ImageValidator)

### Required Implementation
The new GiteaClient MUST:
1. Implement Registry interface with exact signatures
2. Import and use BasicAuthenticator
3. Import and use ImageValidator
4. NOT create any alternative interfaces
```

## Grading Impact
- No research conducted: **-50% penalty**
- Research not documented: **-30% penalty**
- Existing code not identified: **-25% penalty**
- Duplicates not prevented: **-100% (via R373)**

## Integration with Other Rules
- **R373**: Enforces the reuse identified by this research
- **R362**: Prevents architectural changes to found interfaces
- **R251**: Effort plan creation incorporates research
- **R255**: Code review validates research was done

## Tags
#supreme-law #planning #research #code-reuse #r374