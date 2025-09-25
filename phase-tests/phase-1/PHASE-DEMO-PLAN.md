# PHASE 1 DEMO PLAN - COMMAND SKELETON & FOUNDATION
**Project**: idpbuilder-push
**Phase**: 1 - Command Structure and Validation
**Generated**: 2025-09-22
**Demo Author**: Code Reviewer Agent

## Overview

This demo plan provides interactive scenarios to showcase Phase 1 functionality of the idpbuilder push command. These demos are designed to validate the command skeleton, flag handling, and validation logic through user-facing demonstrations.

## Demo Environment Setup

### Prerequisites
```bash
# Ensure idpbuilder is installed
idpbuilder --version

# Create test directory
mkdir -p ~/demo/phase1
cd ~/demo/phase1

# Create sample image files for demos
touch sample-image.tar
touch test-app.tar
echo "mock image content" > mock-image.tar
```

### Environment Variables (Optional)
```bash
# Set default credentials via environment
export IDPBUILDER_REGISTRY_USERNAME=demo-user
export IDPBUILDER_REGISTRY_PASSWORD=demo-pass
export IDPBUILDER_INSECURE=false
```

## Demo Scenarios

### Demo 1: Command Discovery and Help
**Objective**: Demonstrate that the push command is properly integrated and discoverable.

#### Steps:
```bash
# 1. Show idpbuilder help to find push command
$ idpbuilder --help

Expected Output:
idpbuilder CLI for managing developer platforms

Usage:
  idpbuilder [command]

Available Commands:
  create      Create resources
  get         Get resources
  push        Push OCI image to Gitea registry  <-- NEW!
  help        Help about any command

# 2. Get specific help for push command
$ idpbuilder push --help

Expected Output:
Push an OCI image to the configured Gitea registry

Usage:
  idpbuilder push IMAGE REGISTRY [flags]

Flags:
  -u, --username string   Registry username
  -p, --password string   Registry password
  -k, --insecure         Skip TLS verification (use for self-signed certificates)
  -h, --help             help for push

# 3. Try help with just 'help push'
$ idpbuilder help push

[Same output as above]
```

**Success Criteria**:
- ✅ Push command appears in main help
- ✅ Detailed help available for push
- ✅ All flags documented

### Demo 2: Argument Validation
**Objective**: Demonstrate proper validation of command arguments.

#### Steps:
```bash
# 1. Try with no arguments (should fail)
$ idpbuilder push

Expected Error:
Error: requires exactly 2 arguments: IMAGE and REGISTRY
Usage:
  idpbuilder push IMAGE REGISTRY [flags]

# 2. Try with only one argument (should fail)
$ idpbuilder push sample-image.tar

Expected Error:
Error: requires exactly 2 arguments: IMAGE and REGISTRY
Usage:
  idpbuilder push IMAGE REGISTRY [flags]

# 3. Try with too many arguments (should fail)
$ idpbuilder push sample-image.tar registry.com extra-arg

Expected Error:
Error: too many arguments provided
Usage:
  idpbuilder push IMAGE REGISTRY [flags]

# 4. Try with correct number of arguments (validates but doesn't push yet)
$ idpbuilder push sample-image.tar https://gitea.cnoe.localtest.me

Expected Output:
Validating inputs...
✓ Image path: sample-image.tar
✓ Registry URL: https://gitea.cnoe.localtest.me
Error: Push operation not yet implemented (Phase 1 validates only)
```

**Success Criteria**:
- ✅ Rejects incorrect argument counts
- ✅ Clear error messages
- ✅ Shows usage on error

### Demo 3: Flag Handling and Parsing
**Objective**: Demonstrate all flag variations work correctly.

#### Steps:
```bash
# 1. Test with long flags
$ idpbuilder push sample-image.tar https://gitea.cnoe.localtest.me \
    --username demo-user \
    --password demo-pass \
    --insecure

Expected Output:
Configuration parsed:
  Image: sample-image.tar
  Registry: https://gitea.cnoe.localtest.me
  Username: demo-user
  Password: [hidden]
  Insecure: true
Error: Push operation not yet implemented (Phase 1 validates only)

# 2. Test with short flags
$ idpbuilder push sample-image.tar https://gitea.cnoe.localtest.me \
    -u demo-user \
    -p demo-pass \
    -k

[Same configuration output as above]

# 3. Test with mixed flag styles
$ idpbuilder push sample-image.tar https://gitea.cnoe.localtest.me \
    --username demo-user \
    -p demo-pass \
    --insecure

[Same configuration output as above]

# 4. Test with equals syntax
$ idpbuilder push sample-image.tar https://gitea.cnoe.localtest.me \
    --username=demo-user \
    --password=demo-pass

Expected Output:
Configuration parsed:
  Image: sample-image.tar
  Registry: https://gitea.cnoe.localtest.me
  Username: demo-user
  Password: [hidden]
  Insecure: false  <-- Default value
```

**Success Criteria**:
- ✅ All flag formats work
- ✅ Short and long forms equivalent
- ✅ Defaults applied correctly

### Demo 4: Input Validation - Security
**Objective**: Demonstrate security validations prevent dangerous inputs.

#### Steps:
```bash
# 1. Test path traversal prevention
$ idpbuilder push ../../../etc/passwd https://gitea.cnoe.localtest.me

Expected Error:
Error: Invalid image path - path traversal not allowed
  Path: ../../../etc/passwd
  Issue: Path cannot contain '..' sequences

# 2. Test invalid registry URL
$ idpbuilder push sample-image.tar not-a-valid-url

Expected Error:
Error: Invalid registry URL
  URL: not-a-valid-url
  Issue: URL must start with http:// or https://

# 3. Test missing protocol
$ idpbuilder push sample-image.tar gitea.cnoe.localtest.me

Expected Error:
Error: Invalid registry URL - protocol required
  URL: gitea.cnoe.localtest.me
  Hint: Use https://gitea.cnoe.localtest.me

# 4. Test unsupported protocol
$ idpbuilder push sample-image.tar ftp://gitea.cnoe.localtest.me

Expected Error:
Error: Unsupported protocol 'ftp'
  Supported protocols: http, https
```

**Success Criteria**:
- ✅ Path traversal blocked
- ✅ URL validation enforced
- ✅ Clear security messages

### Demo 5: Credential Handling
**Objective**: Demonstrate credential validation and precedence.

#### Steps:
```bash
# 1. Test incomplete credentials (username without password)
$ idpbuilder push sample-image.tar https://gitea.cnoe.localtest.me \
    --username demo-user

Expected Error:
Error: Incomplete credentials
  Password required when username is provided
  Use both --username and --password, or neither (to use defaults)

# 2. Test incomplete credentials (password without username)
$ idpbuilder push sample-image.tar https://gitea.cnoe.localtest.me \
    --password demo-pass

Expected Error:
Error: Incomplete credentials
  Username required when password is provided
  Use both --username and --password, or neither (to use defaults)

# 3. Test with no credentials (should indicate default will be used)
$ idpbuilder push sample-image.tar https://gitea.cnoe.localtest.me

Expected Output:
Configuration parsed:
  Image: sample-image.tar
  Registry: https://gitea.cnoe.localtest.me
  Credentials: Will use defaults from 'idpbuilder get secrets gitea'
Error: Push operation not yet implemented (Phase 1 validates only)

# 4. Test environment variable usage
$ export IDPBUILDER_REGISTRY_USERNAME=env-user
$ export IDPBUILDER_REGISTRY_PASSWORD=env-pass
$ idpbuilder push sample-image.tar https://gitea.cnoe.localtest.me

Expected Output:
Configuration parsed:
  Image: sample-image.tar
  Registry: https://gitea.cnoe.localtest.me
  Username: env-user (from environment)
  Password: [hidden] (from environment)
```

**Success Criteria**:
- ✅ Validates credential pairs
- ✅ Environment variables work
- ✅ Clear precedence rules

### Demo 6: Error Message Quality
**Objective**: Demonstrate helpful, actionable error messages.

#### Steps:
```bash
# 1. Multiple validation errors
$ idpbuilder push "" ""

Expected Error:
Error: Multiple validation failures:
  1. Image path cannot be empty
  2. Registry URL cannot be empty

Usage:
  idpbuilder push IMAGE REGISTRY [flags]

# 2. File not found with suggestion
$ idpbuilder push non-existent.tar https://gitea.cnoe.localtest.me

Expected Warning:
Warning: Image file not found: non-existent.tar
  Current directory: ~/demo/phase1
  Available files: sample-image.tar, test-app.tar, mock-image.tar
  (Validation continues for Phase 1)

# 3. HTTP registry warning
$ idpbuilder push sample-image.tar http://localhost:5000

Expected Output:
Configuration parsed:
  Image: sample-image.tar
  Registry: http://localhost:5000
⚠️  Warning: Using insecure HTTP protocol
  Consider using HTTPS or add --insecure flag to acknowledge
```

**Success Criteria**:
- ✅ Errors are specific
- ✅ Suggestions provided
- ✅ Warnings for security

### Demo 7: Integration with Existing Commands
**Objective**: Show push command works alongside existing idpbuilder commands.

#### Steps:
```bash
# 1. Use get secrets to see default credentials
$ idpbuilder get secrets gitea

Output:
gitea-admin-password: admin123
gitea-registry-username: gitea
gitea-registry-password: gitea123

# 2. Show push will use these defaults
$ idpbuilder push sample-image.tar https://gitea.cnoe.localtest.me

Expected Output:
Configuration parsed:
  Image: sample-image.tar
  Registry: https://gitea.cnoe.localtest.me
  Credentials: Using defaults from secrets (gitea/[hidden])
Error: Push operation not yet implemented (Phase 1 validates only)

# 3. Override defaults with flags
$ idpbuilder push sample-image.tar https://gitea.cnoe.localtest.me \
    --username override-user \
    --password override-pass

Expected Output:
Configuration parsed:
  Image: sample-image.tar
  Registry: https://gitea.cnoe.localtest.me
  Username: override-user (from flags, overrides defaults)
  Password: [hidden]
```

**Success Criteria**:
- ✅ Integrates with get secrets
- ✅ Defaults work correctly
- ✅ Overrides work as expected

### Demo 8: Insecure Mode for Self-Signed Certificates
**Objective**: Demonstrate the --insecure flag for development environments.

#### Steps:
```bash
# 1. Normal HTTPS (would verify certificates)
$ idpbuilder push sample-image.tar https://gitea.cnoe.localtest.me

Expected Output:
Configuration parsed:
  Image: sample-image.tar
  Registry: https://gitea.cnoe.localtest.me
  TLS Verification: Enabled (default)

# 2. With --insecure flag
$ idpbuilder push sample-image.tar https://gitea.cnoe.localtest.me --insecure

Expected Output:
Configuration parsed:
  Image: sample-image.tar
  Registry: https://gitea.cnoe.localtest.me
  TLS Verification: DISABLED
⚠️  Warning: TLS certificate verification is disabled
  This is insecure and should only be used for development

# 3. Short form -k
$ idpbuilder push sample-image.tar https://gitea.cnoe.localtest.me -k

[Same output as above with warning]
```

**Success Criteria**:
- ✅ Insecure flag recognized
- ✅ Clear security warning
- ✅ Both -k and --insecure work

## Interactive Demo Script

### Live Demo Flow (5-10 minutes)
```markdown
1. Introduction (30 seconds)
   - "Today we're demonstrating Phase 1 of the new push command"
   - "This phase implements the command structure and validation"

2. Discovery (1 minute)
   - Run: idpbuilder --help
   - Point out new push command
   - Run: idpbuilder push --help
   - Explain the command structure

3. Basic Usage (2 minutes)
   - Try with wrong arguments (show validation)
   - Try with correct arguments
   - Explain validation vs actual push (Phase 1 scope)

4. Flag Demonstration (2 minutes)
   - Show long flags
   - Show short flags
   - Show mixed usage
   - Demonstrate defaults

5. Security Features (2 minutes)
   - Attempt path traversal (blocked)
   - Try invalid URLs (rejected)
   - Show --insecure flag with warning

6. Integration (1 minute)
   - Show get secrets command
   - Explain credential defaults
   - Show override behavior

7. Q&A (1-2 minutes)
   - Address questions
   - Mention upcoming phases
```

## Demo Verification Checklist

### Pre-Demo Checks
- [ ] idpbuilder installed and working
- [ ] Test files created in demo directory
- [ ] Environment variables cleared or set as needed
- [ ] Terminal recording software ready (if recording)

### During Demo
- [ ] Commands typed clearly (not too fast)
- [ ] Errors shown are intentional and explained
- [ ] Security features emphasized
- [ ] Integration points highlighted

### Post-Demo
- [ ] All Phase 1 features demonstrated
- [ ] Questions answered
- [ ] Next steps explained (Phase 2-4)

## Troubleshooting Common Demo Issues

### Issue: Command not found
```bash
# Solution: Ensure idpbuilder is in PATH
which idpbuilder
export PATH=$PATH:/path/to/idpbuilder
```

### Issue: Test files missing
```bash
# Solution: Recreate test files
touch sample-image.tar test-app.tar
```

### Issue: Environment variables interfering
```bash
# Solution: Clear environment
unset IDPBUILDER_REGISTRY_USERNAME
unset IDPBUILDER_REGISTRY_PASSWORD
unset IDPBUILDER_INSECURE
```

## Demo Success Metrics

### Technical Success
- All commands execute without unexpected errors
- Validation catches all invalid inputs
- Help text is complete and accurate
- Flags work in all formats

### User Experience Success
- Commands are intuitive
- Error messages are helpful
- Security is transparent
- Integration is seamless

## Notes for Demo Presenters

### Key Messages
1. **Phase 1 Focus**: This phase establishes the foundation
2. **TDD Approach**: Tests written first, implementation follows
3. **Security First**: Validation prevents common attacks
4. **User Friendly**: Clear errors and helpful messages
5. **Integration**: Works with existing idpbuilder ecosystem

### What NOT to Demo (Yet)
- Actual image pushing (Phase 3)
- Registry authentication (Phase 2)
- Progress indicators (Phase 4)
- Advanced OCI operations (Future)

### Addressing Questions
- "When will it actually push?": Phase 3 implements the push operation
- "What registries are supported?": Any OCI-compliant registry, starting with Gitea
- "Is Docker required?": No, uses OCI standards directly
- "Can it build images?": No, it pushes existing images only

---

*This demo plan ensures Phase 1 functionality is thoroughly demonstrated and validated through interactive scenarios that stakeholders can follow.*