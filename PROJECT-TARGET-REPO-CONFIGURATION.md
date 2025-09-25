# PROJECT-TARGET-REPO Configuration Guide

## Overview

The `PROJECT-TARGET-REPO` variable represents the URL of your actual project repository that Software Factory 2.0 will work on. This is **NOT** the Software Factory instance repository itself, but rather the target project you want to enhance or build.

## Configuration Locations

The `PROJECT-TARGET-REPO` value should be configured in the following locations:

### 1. setup-config.yaml (Initial Setup)
When running the initial setup with `./setup.sh --config setup-config.yaml`, specify your target repository:

```yaml
target_repository:
  url: "https://github.com/your-org/your-project.git"  # Your PROJECT-TARGET-REPO
  base_branch: "main"
  clone_depth: 100
  auth_method: "https"
```

### 2. target-repo-config.yaml (Runtime Configuration)
After setup, the target repository is configured in this file for runtime operations:

```yaml
target_repository:
  url: "https://github.com/your-org/your-project.git"  # Your PROJECT-TARGET-REPO
  base_branch: "main"
  fork_url: ""  # Optional: if working with a fork
  auth_method: "https"
```

### 3. Environment Variable (Optional)
For CI/CD or automation scenarios, you can export the variable:

```bash
export PROJECT_TARGET_REPO="https://github.com/your-org/your-project.git"
```

## Important Notes

### What NOT to Use
- **NEVER** use the Software Factory repository URL as your target
- **NEVER** use placeholder values like "idpbuilder" unless that's your actual project
- **NEVER** leave this unconfigured - it will cause complete failure

### Common Examples
Instead of hard-coded values like:
- ❌ `https://github.com/cnoe-io/idpbuilder.git`
- ❌ `https://github.com/user/software-factory-instance.git`

Use your actual project:
- ✅ `https://github.com/your-org/your-actual-project.git`
- ✅ `git@github.com:your-org/your-project.git` (for SSH)

## Validation

The orchestrator will validate the target repository configuration on startup:

1. **Existence Check**: Verifies target-repo-config.yaml exists
2. **URL Validation**: Ensures a valid repository URL is configured
3. **Separation Check**: Confirms the target is NOT the SF instance itself
4. **Access Check**: Verifies the repository can be cloned

## Migration from Hard-coded Values

If you're updating from an older template that had hard-coded repository URLs:

1. Update your `setup-config.yaml` with your actual target repository URL
2. Ensure `target-repo-config.yaml` has the correct URL
3. Run validation: `./verify-installation.sh`

## Related Rules

- **R191**: Target Repository Configuration enforcement
- **R192**: Repository Separation (SF vs Target)
- **R309**: Never Create Efforts in SF Repository

## Troubleshooting

### Error: "CRITICAL: target-repo-config.yaml not found!"
**Solution**: Create the file with your target repository URL

### Error: "Target is Software Factory!"
**Solution**: You've configured the SF repository as your target. Use your actual project repository URL instead.

### Error: "No target repository URL configured!"
**Solution**: Add the `target_repository.url` field to your configuration file

---

**Remember**: The Software Factory is the TOOL. Your project repository is the TARGET. Keep them separate!