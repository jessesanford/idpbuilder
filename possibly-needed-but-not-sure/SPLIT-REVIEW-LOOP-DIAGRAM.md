# Branch Split Review Loop Process

## Visual Flow Diagram

```
┌─────────────────────┐
│   Original Effort   │
│   (>800 lines)      │
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│  Code Review by     │
│  @kcp-k8s-reviewer  │
└──────────┬──────────┘
           │
           ▼
      ┌────────┐
      │ >800?  │──NO──→ [ACCEPTED]
      └───┬────┘
          │YES
          ▼
    ┌─────────────┐
    │  Exception  │
    │  Warranted? │
    └──────┬──────┘
         /   \
       YES    NO
        │      │
        ▼      ▼
 [ACCEPTED  ┌─────────────────────┐
  WITH      │  Create Split Plan  │
  EXCEPTION]│  (APIs→Impl→Tests)  │
            └──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│  Implement Splits   │
│  (Sequential Only)  │
└──────────┬──────────┘
           │
           ▼
┌─────────────────────────┐
│  Review Each Split      │
│  @kcp-k8s-reviewer      │
└──────────┬──────────────┘
           │
    ┌──────▼──────┐
    │             │
    ▼             ▼
[All ≤800]    [Any >800]
    │             │
    │             ▼
    │      ┌──────────────┐
    │      │ Create       │
    │      │ Sub-Split    │
    │      │ Plan         │
    │      └──────┬───────┘
    │             │
    │             └────────┐
    │                      │
    ▼                      │
┌─────────────┐            │
│  ACCEPTED   │            │
│  All Splits │            │
│  ≤800 lines │            │
└─────────────┘            │
                           │
                   LOOP BACK
```

## Iterative Split Algorithm (With Exception Handling)

```python
def handle_effort_with_splits(effort):
    """
    Main orchestrator function for handling efforts that may need splitting
    """
    # Initial implementation
    effort_branch = implement_effort(effort)
    
    # Start review loop
    branches_to_review = [effort_branch]
    accepted_branches = []
    iteration = 0
    
    while branches_to_review:
        iteration += 1
        print(f"Review iteration {iteration}")
        
        for branch in branches_to_review[:]:  # Copy list to modify during iteration
            review = review_branch(branch)
            
            if review.line_count <= 800 and review.quality_pass:
                # Branch is acceptable
                accepted_branches.append(branch)
                branches_to_review.remove(branch)
                
            elif review.line_count > 800:
                # Check if reviewer grants exception
                if review.exception_granted:
                    # Reviewer determined split would break logic
                    print(f"Exception granted for {branch}: {review.exception_reason}")
                    accepted_branches.append({
                        'branch': branch,
                        'exception': review.exception_details,
                        'extra_validation': True
                    })
                    branches_to_review.remove(branch)
                else:
                    # Need to split - reviewer did not grant exception
                    split_plan = create_split_plan(branch, review)
                    split_branches = implement_splits(branch, split_plan)
                    
                    # Remove original, add splits for review
                    branches_to_review.remove(branch)
                    branches_to_review.extend(split_branches)
                
            else:
                # Quality issues - fix and re-review
                fix_quality_issues(branch, review.issues)
                # Keep in branches_to_review for next iteration
    
    return accepted_branches
```

## Split Priority Rules

### First Split (Always in this order)
1. **APIs/Types/Interfaces** - Public contracts first
2. **Core Implementation** - Business logic
3. **Tests** - Validation code
4. **Documentation** - Supporting docs

### Recursive Splits (If needed)
```
Original (2000 lines)
├── Part 1: APIs (300 lines) ✓
├── Part 2: Implementation (1200 lines) ✗ Too large
│   ├── Part 2.1: Core Controllers (600 lines) ✓
│   └── Part 2.2: Helper Functions (600 lines) ✓
└── Part 3: Tests (500 lines) ✓
```

## Review Loop State Machine

```yaml
states:
  IMPLEMENTING:
    next: REVIEWING
    
  REVIEWING:
    transitions:
      - condition: "lines <= 800 AND quality_pass"
        next: ACCEPTED
      - condition: "lines > 800"
        next: CREATING_SPLIT_PLAN
      - condition: "quality_fail"
        next: FIXING_ISSUES
        
  CREATING_SPLIT_PLAN:
    next: IMPLEMENTING_SPLITS
    
  IMPLEMENTING_SPLITS:
    next: REVIEWING  # Loop back
    
  FIXING_ISSUES:
    next: REVIEWING  # Loop back
    
  ACCEPTED:
    final: true
```

## Orchestrator Commands for Split Loop

```bash
# Initial review
review_result=$(review_effort_branch $EFFORT_BRANCH)

# Check if split needed
if [ $(echo $review_result | jq '.line_count') -gt 800 ]; then
    # Create split plan
    split_plan=$(create_split_plan $EFFORT_BRANCH)
    
    # Implement splits sequentially
    for split in $(echo $split_plan | jq -r '.splits[]'); do
        implement_split $split
        
        # Review the split immediately
        split_review=$(review_effort_branch $split)
        
        # If still too large, recurse
        if [ $(echo $split_review | jq '.line_count') -gt 800 ]; then
            # Create sub-split plan
            sub_split_plan=$(create_split_plan $split)
            # Continue recursion...
        fi
    done
fi
```

## Critical Rules for Split Reviews

1. **NEVER PARALLELIZE SPLITS** - Always sequential
2. **ALWAYS RE-REVIEW** - Every split gets full review
3. **MAINTAIN QUALITY** - Each split must build/test independently
4. **PRESERVE LOGIC** - Splits must be logically cohesive
5. **TRACK HIERARCHY** - Document parent->child relationships
6. **VERIFY INTEGRATION** - Merged splits = original functionality
7. **EXCEPTION AUTHORITY** - Only @agent-kcp-kubernetes-code-reviewer can grant >800 line exceptions

## Exception Granting Guidelines for @agent-kcp-kubernetes-code-reviewer

When reviewing a branch >800 lines, evaluate these factors:

```yaml
exception_evaluation:
  must_stay_together:
    - "Would splitting break atomic transactions?"
    - "Are interfaces tightly coupled to implementations?"
    - "Is this a complex state machine?"
    - "Would split duplicate significant code?"
    - "Is this mostly generated code?"
    
  if_exception_granted:
    document:
      - Specific reason why split would harm code quality
      - Risk assessment of large PR
      - Mitigation strategies (extra tests, multiple reviewers)
    
  if_exception_denied:
    provide:
      - Clear split boundaries
      - Logical groupings for splits
      - Dependencies between splits
```

### Example Exception Grant

```yaml
# In CODE-REVIEW-2025-08-20-1430.md
review_result:
  branch: /phase3/wave1/effort1-sync-engine
  line_count: 1350
  
  exception_granted: true
  exception_details:
    reason: "Complex bidirectional sync state machine"
    justification: |
      The sync engine implements a 7-state bidirectional synchronization
      protocol. States share transition functions and error recovery logic.
      Splitting would require either:
      1. Duplicating 400+ lines of shared transition code
      2. Creating artificial interfaces that reduce readability
      
      The current implementation is the minimal cohesive unit.
    
    risk_mitigation:
      - Required: 2 additional reviewers familiar with sync protocols
      - Required: Extended integration test suite
      - Required: Performance benchmarks before/after
      - Recommended: Split into 5+ logical commits within branch
    
    commit_structure:
      - "Commit 1: Sync state definitions and interfaces"
      - "Commit 2: State transition functions"  
      - "Commit 3: Error recovery mechanisms"
      - "Commit 4: Bidirectional sync implementation"
      - "Commit 5: Tests and validation"
```

## Example Split Tracking

```yaml
effort: E3.2.1-resource-sync
original_lines: 2500
status: split_complete
splits:
  - branch: E3.2.1-resource-sync-part1
    description: "API types and interfaces"
    lines: 400
    status: accepted
    
  - branch: E3.2.1-resource-sync-part2
    description: "Core sync logic"
    lines: 1500
    status: split_required
    sub_splits:
      - branch: E3.2.1-resource-sync-part2-subpart1
        description: "Reconciliation logic"
        lines: 700
        status: accepted
        
      - branch: E3.2.1-resource-sync-part2-subpart2
        description: "Status update logic"
        lines: 800
        status: accepted
        
  - branch: E3.2.1-resource-sync-part3
    description: "Tests and validation"
    lines: 600
    status: accepted
```

## Success Metrics

- All branches ≤ 800 lines
- All branches pass quality checks
- All branches integrate correctly
- Total functionality preserved
- Clean PR history maintained