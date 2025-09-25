# Archived Fix States

This directory contains archived fix state files from completed fix cascades, hotfixes, and error recovery operations.

## Directory Structure

```
archived-fixes/
├── README.md (this file)
├── 2025/
│   ├── 01/
│   │   ├── auth-hotfix-20250121-143000.json
│   │   ├── gitea-api-fix-20250122-091500.json
│   │   └── ...
│   ├── 02/
│   │   └── ...
│   └── ...
└── ...
```

## Archival Process (Per R375)

### When to Archive
Fix state files are archived when:
- Fix cascade completes successfully
- All validation passes
- All target branches updated
- No pending recovery actions

### How to Archive

1. **Complete the fix**
   ```bash
   # Verify fix is complete
   cat orchestrator-[fix-name]-state.json | jq '.status'
   # Should show: "COMPLETED"
   ```

2. **Archive the state file**
   ```bash
   # Set variables
   TIMESTAMP=$(date +%Y%m%d-%H%M%S)
   FIX_NAME="gitea-api-fix"  # Your fix identifier
   YEAR=$(date +%Y)
   MONTH=$(date +%m)

   # Create archive directory
   mkdir -p archived-fixes/${YEAR}/${MONTH}

   # Move and rename with timestamp
   mv orchestrator-${FIX_NAME}-state.json \
      archived-fixes/${YEAR}/${MONTH}/${FIX_NAME}-${TIMESTAMP}.json
   ```

3. **Commit the archive**
   ```bash
   git add archived-fixes/
   git commit -m "archive: ${FIX_NAME} completed at ${TIMESTAMP}"
   git push
   ```

## Archive Naming Convention

```
[fix-identifier]-[YYYYMMDD]-[HHMMSS].json
```

Examples:
- `auth-hotfix-20250121-143000.json`
- `pr367-backport-20250122-091500.json`
- `critical-security-fix-20250123-180000.json`

## Searching Archives

### Find all fixes for a specific month
```bash
ls archived-fixes/2025/01/*.json
```

### Search for a specific fix
```bash
find archived-fixes -name "*auth-hotfix*.json"
```

### View fix summary
```bash
cat archived-fixes/2025/01/auth-hotfix-*.json | jq '{
  fix: .fix_identifier,
  type: .fix_type,
  status: .status,
  created: .created_at,
  completed: .archival_info.archived_at
}'
```

### Check validation results
```bash
cat archived-fixes/2025/01/fix-*.json | jq '.validation_results'
```

## Archive Retention

- Archives are kept indefinitely for audit trail
- Never delete archived fix states
- Can be compressed yearly if needed:
  ```bash
  tar -czf archived-fixes-2025.tar.gz archived-fixes/2025/
  ```

## Recovery from Archives

If you need to reference a previous fix:

```bash
# Copy archive to temp for reference (don't modify archive!)
cp archived-fixes/2025/01/auth-hotfix-*.json /tmp/reference-fix.json

# View the fix pattern
cat /tmp/reference-fix.json | jq '.'
```

## Statistics

Generate monthly fix statistics:
```bash
# Count fixes per month
ls archived-fixes/2025/01/*.json 2>/dev/null | wc -l

# List fix types
for f in archived-fixes/2025/01/*.json; do
  jq -r '.fix_type' "$f"
done | sort | uniq -c

# Average completion time
for f in archived-fixes/2025/01/*.json; do
  jq '.metrics.total_time_minutes' "$f"
done | awk '{sum+=$1; count++} END {print sum/count " minutes average"}'
```

## Important Notes

1. **Never modify archived files** - They are historical records
2. **Always archive, never delete** - Maintains complete audit trail
3. **Include timestamp in filename** - Prevents overwrites
4. **Organize by year/month** - Keeps directory manageable
5. **Commit immediately after archiving** - Ensures persistence

---

Per R375 - Fix State File Management Protocol