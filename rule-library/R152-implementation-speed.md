# 🚨 RULE R152 - Implementation Speed Requirements

**Criticality:** CRITICAL - Performance metric  
**Grading Impact:** 30% of SW Engineer grade  
**Enforcement:** CONTINUOUS - Tracked per session

## Rule Statement

SW Engineers MUST maintain an implementation speed of >50 lines per hour to remain employed.

## Speed Metrics

### Minimum Thresholds
| Metric | Minimum | Target | Excellent |
|--------|---------|--------|-----------|
| Lines/Hour | 50 | 75 | 100+ |
| Functions/Hour | 3 | 5 | 8+ |
| Tests/Hour | 2 | 4 | 6+ |
| Commits/Hour | 2 | 3 | 4+ |

### Calculation Method
```python
def calculate_implementation_speed(session):
    # Exclude breaks, planning, reviews
    actual_coding_time = session.total_time - session.non_coding_time
    
    # Count only NEW lines (not modifications)
    net_lines_added = session.lines_added - session.lines_removed
    
    # Calculate hourly rate
    hours = actual_coding_time / 3600
    lines_per_hour = net_lines_added / hours if hours > 0 else 0
    
    return {
        'lines_per_hour': lines_per_hour,
        'grade': min(lines_per_hour / 50, 1.0) * 30  # 30% of total grade
    }
```

## What Counts Toward Speed

### Included ✅
- New source code lines
- New test code lines
- Documentation in code files
- Configuration as code
- Generated code you modified

### Excluded ❌
- Pure generated code (unmodified)
- Copied/pasted code from examples
- Comments added to existing code
- Whitespace and formatting
- Deleted lines

## Speed Tracking

### Real-time Monitoring
```bash
track_implementation_speed() {
    local start_time=$1
    local current_time=$(date +%s)
    local elapsed=$((current_time - start_time))
    local hours=$(echo "scale=2; $elapsed / 3600" | bc)
    
    # Get line count
    local lines=$(git diff --stat | tail -1 | awk '{print $4}')
    local speed=$(echo "scale=1; $lines / $hours" | bc)
    
    echo "📊 SPEED METRICS"
    echo "Time elapsed: ${hours}h"
    echo "Lines written: $lines"
    echo "Speed: ${speed} lines/hour"
    
    if (( $(echo "$speed < 50" | bc -l) )); then
        echo "⚠️ WARNING: Below minimum speed!"
        echo "🎯 Need $(echo "50 * $hours - $lines" | bc) more lines"
    fi
}
```

### Work Log Updates
```bash
# Required in work-log.md
## Speed Metrics
- Session start: 2025-08-26T14:00:00Z
- Current time: 2025-08-26T15:30:00Z
- Lines written: 85
- Current speed: 56.7 lines/hour ✅
- Status: ON TRACK
```

## Speed Optimization Strategies

### 1. Batch Similar Tasks
```bash
# GOOD: Implement all similar functions together
implement_crud_operations() {
    # Create all CRUD operations in sequence
    write_create_function
    write_read_function
    write_update_function
    write_delete_function
    # Higher speed through pattern repetition
}
```

### 2. Use Templates and Snippets
```bash
# GOOD: Leverage templates for speed
generate_from_template() {
    local entity=$1
    sed "s/{{ENTITY}}/$entity/g" template.go > "${entity}.go"
    # Quickly customize
    edit_for_specifics
}
```

### 3. Write Tests Alongside Code
```bash
# GOOD: Parallel implementation
write_function_and_test() {
    # Write function (25 lines)
    implement_function
    # Immediately write test (15 lines)
    implement_test
    # 40 lines in quick succession
}
```

## Common Speed Killers

### VIOLATION: Over-Planning
```bash
# ❌ WRONG: 45 minutes planning, 15 minutes coding
# Result: 20 lines in 1 hour = 20 lines/hour FAIL
```

### VIOLATION: Perfectionism
```bash
# ❌ WRONG: Rewriting same function 5 times
# Result: 100 lines written, 80 deleted = 20 net lines/hour FAIL
```

### VIOLATION: Debugging Instead of Writing
```bash
# ❌ WRONG: Spending hour fixing one bug
# Result: 5 lines changed = 5 lines/hour FAIL
```

## Recovery Strategies

### When Behind Schedule
```bash
catch_up_strategy() {
    echo "⚠️ Currently behind speed target"
    echo "📋 RECOVERY PLAN:"
    echo "1. Switch to simpler tasks"
    echo "2. Use more templates"
    echo "3. Batch similar implementations"
    echo "4. Defer complex debugging"
    echo "5. Focus on line count for next 30 min"
}
```

### Speed Burst Technique
```bash
# When you need to boost speed
speed_burst() {
    echo "🚀 SPEED BURST MODE"
    # 1. Pick repetitive tasks
    # 2. Disable all distractions
    # 3. Use snippets/templates
    # 4. Target: 100+ lines in 30 minutes
    # 5. Clean up after burst
}
```

## Grading Formula

```python
def grade_implementation_speed(sw_engineer):
    sessions = sw_engineer.get_coding_sessions()
    
    total_weighted_speed = 0
    total_weight = 0
    
    for session in sessions:
        # Recent sessions weighted more
        recency_weight = calculate_recency_weight(session)
        
        # Calculate session speed
        speed = session.lines_written / session.hours_coding
        
        # Apply weight
        total_weighted_speed += speed * recency_weight
        total_weight += recency_weight
    
    average_speed = total_weighted_speed / total_weight
    
    # Grade calculation (30% of total)
    if average_speed >= 100:
        return 30  # Full marks
    elif average_speed >= 75:
        return 27  # 90%
    elif average_speed >= 50:
        return 21  # 70%
    else:
        return average_speed / 50 * 21  # Proportional
```

## Speed vs Quality Balance

### Acceptable Trade-offs
```
Speed Priority (First Pass):
- Basic functionality ✅
- Core happy path ✅
- Minimal error handling ✅
- Basic tests ✅

Quality Priority (Second Pass):
- Edge cases
- Comprehensive error handling
- Full test coverage
- Performance optimization
```

### Unacceptable Shortcuts
```
NEVER sacrifice:
- Compilation ❌
- Basic correctness ❌
- Security basics ❌
- Critical tests ❌
```

## Daily Speed Report

```yaml
daily_speed_report:
  date: "2025-08-26"
  sw_engineer: "sw-engineer-1"
  
  sessions:
    - start: "09:00"
      end: "11:00"
      lines: 120
      speed: 60
      status: "PASS"
    
    - start: "13:00"
      end: "15:00"
      lines: 95
      speed: 47.5
      status: "WARNING"
    
    - start: "15:30"
      end: "17:00"
      lines: 85
      speed: 56.7
      status: "PASS"
  
  daily_total:
    hours: 5.5
    lines: 300
    average_speed: 54.5
    grade: "PASS"
    
  trend: "improving"
  recommendation: "Maintain current pace"
```

---
**Remember:** Speed matters. Below 50 lines/hour = Below employment threshold.