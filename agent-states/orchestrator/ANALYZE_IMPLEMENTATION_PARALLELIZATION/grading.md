# ANALYZE_IMPLEMENTATION_PARALLELIZATION - Grading Rubric

## Grade Composition
- **Plan Analysis Completeness**: 35%
- **Consistency Verification**: 25%
- **State File Updates**: 20%
- **Acknowledgment Protocol**: 15%
- **Directory Preparation**: 5%

## Detailed Scoring

### Plan Analysis Completeness (35%)
| Criterion | Points | Requirement |
|-----------|--------|-------------|
| All plans read with Read tool | 15% | MUST read every IMPLEMENTATION-PLAN.md |
| Metadata extracted correctly | 10% | Parallelization info from each plan |
| Dependencies identified | 10% | Inter-effort dependencies mapped |

### Consistency Verification (25%)
| Criterion | Points | Requirement |
|-----------|--------|-------------|
| Wave plan comparison | 10% | Verify against original wave plan |
| Metadata preservation check | 10% | R211 compliance verified |
| Conflict resolution | 5% | Handle any inconsistencies |

### State File Updates (20%)
| Criterion | Points | Requirement |
|-----------|--------|-------------|
| SW Engineer plan saved | 10% | Complete plan in orchestrator-state.yaml |
| Spawn sequence defined | 5% | Clear step-by-step execution |
| Timestamp recorded | 5% | Analysis timestamp included |

### Acknowledgment Protocol (15%)
| Criterion | Points | Requirement |
|-----------|--------|-------------|
| All plans listed | 5% | Every effort plan path shown |
| Decision output | 5% | Blocking vs parallel groups |
| Strategy commitment | 5% | Clear spawn strategy statement |

### Directory Preparation (5%)
| Criterion | Points | Requirement |
|-----------|--------|-------------|
| R208 compliance shown | 5% | Directory paths verified |

## Critical Failures (Automatic F)

1. **Skipping this state entirely** → 0% (MISSION CRITICAL FAILURE)
2. **Not reading any implementation plans** → Maximum 30% (Analysis failure)
3. **Spawning without parallelization analysis** → 0% (Protocol violation)
4. **Major inconsistency with wave plan** → Maximum 50% (Coordination failure)

## Grade Calculation Example

```yaml
implementation_parallelization_grade:
  plan_analysis:
    all_plans_read: 15/15  # Read 5 implementation plans
    metadata_extracted: 10/10  # Got parallelization info
    dependencies_mapped: 10/10  # Found blocking relationship
    
  consistency_verification:
    wave_plan_comparison: 10/10  # Matched wave metadata
    metadata_preservation: 10/10  # R211 compliance good
    conflict_resolution: 5/5  # No conflicts found
    
  state_file_updates:
    sw_engineer_plan: 10/10  # Complete plan saved
    spawn_sequence: 5/5  # 2-step sequence defined
    timestamp: 5/5  # Included
    
  acknowledgment:
    plans_listed: 5/5  # All 5 paths shown
    decision_output: 5/5  # Groups clearly stated
    strategy_commitment: 5/5  # Commitment stated
    
  directory_prep:
    r208_compliance: 5/5  # Directories verified
    
  total: 100/100  # Grade: A
```

## Performance Metrics

- **Time in state**: Should complete in 2-3 minutes
- **Read operations**: Minimum = number of efforts
- **State file updates**: Minimum 1
- **Consistency checks**: All efforts verified