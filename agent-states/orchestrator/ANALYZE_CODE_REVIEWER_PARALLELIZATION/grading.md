# ANALYZE_CODE_REVIEWER_PARALLELIZATION - Grading Rubric

## Grade Composition
- **Parallelization Analysis**: 40%
- **State File Updates**: 30%
- **Acknowledgment Protocol**: 20%
- **Validation Completeness**: 10%

## Detailed Scoring

### Parallelization Analysis (40%)
| Criterion | Points | Requirement |
|-----------|--------|-------------|
| Wave plan read with Read tool | 10% | MUST use Read tool on wave implementation plan |
| Blocking efforts identified | 10% | ALL "Can Parallelize: No" efforts found |
| Parallel groups created | 10% | Correct grouping of parallel efforts |
| Spawn sequence defined | 10% | Clear step-by-step execution plan |

### State File Updates (30%)
| Criterion | Points | Requirement |
|-----------|--------|-------------|
| Parallelization plan saved | 15% | Complete plan in orchestrator-state.json |
| Timestamp recorded | 5% | Analysis timestamp included |
| Spawn sequence documented | 10% | All steps clearly defined |

### Acknowledgment Protocol (20%)
| Criterion | Points | Requirement |
|-----------|--------|-------------|
| Plan location cited | 5% | Full path to wave plan stated |
| Decision output shown | 10% | Blocking vs parallel efforts listed |
| Strategy commitment stated | 5% | Clear statement of spawn strategy |

### Validation Completeness (10%)
| Criterion | Points | Requirement |
|-----------|--------|-------------|
| All checks passed | 5% | Validation function executed |
| Ready for transition | 5% | Explicit readiness confirmation |

## Critical Failures (Automatic F)

1. **Skipping this state entirely** → 0% (MISSION CRITICAL FAILURE)
2. **Not using Read tool on wave plan** → Maximum 50% (R218 violation)
3. **Spawning without analysis** → 0% (Protocol violation)
4. **Wrong parallelization groups** → Maximum 40% (Analysis failure)

## Grade Calculation Example

```yaml
parallelization_analysis_grade:
  wave_plan_read: 10/10  # Used Read tool
  blocking_identified: 10/10  # Found E3.1.1
  parallel_groups: 10/10  # Correctly grouped E3.1.2-E3.1.5
  spawn_sequence: 10/10  # Clear 2-step plan
  
state_file_updates:
  plan_saved: 15/15  # Complete plan in yaml
  timestamp: 5/5  # Included
  sequence_documented: 10/10  # All steps defined
  
acknowledgment:
  location_cited: 5/5  # Full path shown
  decision_output: 10/10  # All efforts listed
  strategy_commitment: 5/5  # Clear statement
  
validation:
  checks_passed: 5/5  # All validations run
  ready_confirmation: 5/5  # Explicit statement
  
total: 100/100  # Grade: A
```

## Performance Metrics

- **Time in state**: Should complete in <2 minutes
- **Read operations**: Minimum 1 (wave plan)
- **State file updates**: Minimum 1
- **Output lines**: 20-50 (showing analysis)