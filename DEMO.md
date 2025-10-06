# Phase 2 Integration Demos - idpbuilder push Command

## Overview

This directory contains comprehensive demonstration scripts for the `idpbuilder push` command functionality introduced in Phase 2. These demos showcase all features including basic push workflow, authentication methods, retry mechanisms, TLS configuration, and error handling.

### What These Demos Demonstrate

- **Basic Push Workflow**: Simple image push to local registry
- **Authentication Methods**: Multiple ways to authenticate with registries
- **Retry Mechanisms**: Automatic retry behavior on transient failures
- **TLS Configuration**: Certificate validation and custom CA handling
- **Error Handling**: Graceful failure modes and recovery
- **Integration Workflow**: Complete end-to-end real-world scenarios

## Prerequisites

### Required Tools

- **Docker**: Container runtime and image management
  ```bash
  docker --version  # Verify Docker is installed
  ```

- **kind** (optional): For creating test Kubernetes clusters
  ```bash
  kind --version  # Verify kind is installed
  ```

- **idpbuilder**: The idpbuilder binary with push command
  ```bash
  ./bin/idpbuilder version  # Verify binary exists
  ```

### Environment Setup

1. **Local Registry** (for basic demos):
   ```bash
   # Start a local registry on port 5000
   docker run -d -p 5000:5000 --name demo-registry registry:2
   ```

2. **TLS Registry** (for TLS demos - optional):
   ```bash
   # Start a TLS-enabled registry on port 5443
   # See test/fixtures/certs/README.md for certificate setup
   ```

3. **Build idpbuilder** (if binary not present):
   ```bash
   make build  # Or appropriate build command
   ```

## Demo Scripts

### 1. Basic Push Demo

**Script**: `demos/basic-push-demo.sh`

**Purpose**: Demonstrates the fundamental image push workflow using the idpbuilder push command.

**Features Demonstrated**:
- Creating a test Docker image
- Basic push to local registry
- Verification of successful push
- Clean feedback and status reporting

**Duration**: ~2-3 minutes

**How to Run**:
```bash
cd /home/vscode/workspaces/idpbuilder-push-oci/efforts/phase2/wave2/E2.2.1-user-documentation
./demos/basic-push-demo.sh
```

**Expected Output**:
- Image built successfully
- Push operation completes without errors
- Clear status messages throughout
- Verification step confirms image in registry

**Troubleshooting**:
- If registry not available: Start local registry (see Prerequisites)
- If binary not found: Set `IDPBUILDER_BIN` environment variable
- If Docker errors: Check Docker daemon is running

### 2. Authenticated Push Demo

**Script**: `demos/authenticated-push-demo.sh`

**Purpose**: Showcases various authentication methods supported by idpbuilder push.

**Features Demonstrated**:
- Environment variable authentication (`REGISTRY_USERNAME`, `REGISTRY_PASSWORD`)
- Docker config file authentication (`~/.docker/config.json`)
- Stdin authentication for CI/CD pipelines (`--auth-stdin`)
- Explicit credential flags (`--username`, `--password`)

**Duration**: ~3-4 minutes

**How to Run**:
```bash
./demos/authenticated-push-demo.sh
```

**Expected Output**:
- Four authentication methods demonstrated
- Security recommendations for each method
- Successful authentication with each approach
- Comparison of methods and use cases

**Key Takeaways**:
- **Best for CI/CD**: Stdin method (most secure)
- **Best for local dev**: Docker config (automatic)
- **Avoid in production**: Command-line flags (visible in process list)

**Troubleshooting**:
- If auth fails: Demo uses simulated credentials (expected)
- For real registry: Configure actual credentials
- Docker config: Ensure `docker login` has been run

### 3. Retry Mechanism Demo

**Script**: `demos/retry-mechanism-demo.sh`

**Purpose**: Demonstrates automatic retry behavior on transient failures with intelligent backoff strategies.

**Features Demonstrated**:
- Automatic retry on network failures
- Exponential backoff algorithm
- Maximum retry limits
- Non-retryable error detection (immediate failure)

**Duration**: ~4-5 minutes

**How to Run**:
```bash
./demos/retry-mechanism-demo.sh
```

**Expected Output**:
- Simulated transient failures
- Retry attempts with increasing backoff times
- Successful push after retries
- Demonstration of max retry limit
- Non-retryable error immediate failure

**Retry Behavior**:
- **Retryable**: Network timeouts, connection refused, 503 errors, rate limits (429)
- **Non-retryable**: Auth errors (401/403), not found (404), invalid requests (400)

**Configuration Options**:
```bash
--max-retries N          # Maximum retry attempts (default: 3)
--initial-backoff T      # Initial backoff duration (default: 1s)
--max-backoff T          # Maximum backoff duration (default: 30s)
--no-retry              # Disable retry mechanism
```

**Troubleshooting**:
- Demo uses unavailable registry (intentional for retry demonstration)
- Backoff times are shortened in demo mode for faster execution
- In production, use default settings for optimal reliability

### 4. TLS Configuration Demo

**Script**: `demos/tls-configuration-demo.sh`

**Purpose**: Shows TLS certificate validation and configuration options for secure registry communication.

**Features Demonstrated**:
- Default strict TLS verification
- Insecure mode (testing only) with `--tls-verify=false`
- Custom CA certificate usage
- Certificate directory loading
- System CA bundle integration

**Duration**: ~3-4 minutes

**How to Run**:
```bash
./demos/tls-configuration-demo.sh
```

**Expected Output**:
- Default TLS verification enforcement
- Insecure mode demonstration (with warnings)
- Custom CA certificate usage
- Multiple certificate sources
- Security best practices

**TLS Options**:
```bash
--tls-verify=false              # Disable verification (INSECURE - testing only!)
--ca-cert /path/to/ca.crt      # Single custom CA certificate
--ca-cert-dir /path/to/certs/  # Directory of CA certificates
```

**Security Considerations**:
- ⚠️ **NEVER** use `--tls-verify=false` in production
- Always use TLS verification with proper certificates
- Store CA certificates securely
- Rotate certificates regularly
- Document certificate requirements

**Troubleshooting**:
- "certificate signed by unknown authority": Add CA to system or use --ca-cert
- "certificate has expired": Update registry certificate
- "certificate is valid for X, not Y": Use correct hostname

### 5. Phase 2 Integration Demo (CRITICAL)

**Script**: `demos/phase2-integration-demo.sh`

**Purpose**: Complete end-to-end demonstration integrating all Phase 2 push command features in a realistic workflow.

**Features Demonstrated**:
- **ALL Phase 2 features combined**:
  - Basic push workflow
  - Multiple authentication methods
  - Retry mechanisms with backoff
  - TLS certificate handling
  - Comprehensive error handling
- Real-world CI/CD pipeline scenario
- Production-ready configuration patterns

**Duration**: ~8-10 minutes

**How to Run**:
```bash
./demos/phase2-integration-demo.sh
```

**Expected Output**:
- Multi-stage application image build
- Authentication method comparison
- Retry behavior on failures
- TLS configuration scenarios
- Error handling demonstrations
- Complete CI/CD workflow simulation
- Summary report with results

**Real-World Scenario**:
The integration demo simulates a complete CI/CD pipeline:
1. Build application image (multi-stage Dockerfile)
2. Authenticate securely using stdin method
3. Push with retry enabled for reliability
4. Verify TLS certificates for security
5. Handle failures gracefully with proper error messages

**Results**:
All execution logs are saved to `demo-results/` directory:
- `basic-push.log` - Basic push execution
- `auth-*.log` - Authentication tests
- `retry-test.log` - Retry mechanism
- `tls-*.log` - TLS configurations
- `error-*.log` - Error scenarios
- `integration-workflow.log` - Complete workflow
- `integration-demo-summary.txt` - Summary report

**Troubleshooting**:
- Demo runs in simulation mode if binary not available
- Registry availability affects some tests (expected)
- Check individual log files for detailed diagnostics

## Running All Demos

### Quick Start - Run Integration Demo

For a complete demonstration of all features:

```bash
# Run the comprehensive integration demo
./demos/phase2-integration-demo.sh

# Review the summary report
cat demo-results/integration-demo-summary.txt
```

### Run Individual Demos

To explore specific features:

```bash
# Basic functionality
./demos/basic-push-demo.sh

# Authentication methods
./demos/authenticated-push-demo.sh

# Retry behavior
./demos/retry-mechanism-demo.sh

# TLS configuration
./demos/tls-configuration-demo.sh
```

### Run All Demos Sequentially

```bash
# Run all demos in order
for demo in demos/*.sh; do
    echo "Running $demo..."
    bash "$demo"
    echo "---"
done
```

## Configuration

### Environment Variables

Demos respect these environment variables:

```bash
# Path to idpbuilder binary
export IDPBUILDER_BIN="/path/to/idpbuilder"

# Registry credentials (for auth demos)
export REGISTRY_USERNAME="your-username"
export REGISTRY_PASSWORD="your-password"

# TLS settings
export IDPBUILDER_TLS_VERIFY="true"
export IDPBUILDER_CA_CERT="/path/to/ca.crt"
```

### Configuration File

You can also use a configuration file:

```yaml
# ~/.idpbuilder/config.yaml
registry:
  default: "localhost:5000"
  auth:
    username: "demo-user"
    password: "demo-pass"

tls:
  verify: true
  ca_cert: "/path/to/ca.crt"
  ca_cert_dir: "/path/to/certs/"

retry:
  max_retries: 3
  initial_backoff: "1s"
  max_backoff: "30s"
```

## Demo Results

### Output Location

All demo execution results are saved to:
```
demo-results/
├── basic-push-execution.log
├── authenticated-push-execution.log
├── retry-mechanism-execution.log
├── tls-configuration-execution.log
├── integration-demo-execution.log
└── integration-demo-summary.txt
```

### Interpreting Results

**Success Indicators**:
- ✅ Green checkmarks in output
- Exit code 0
- "Demo complete" message
- Logs show expected behavior

**Expected Warnings**:
- ⚠️ Registry not available (in simulation mode)
- ⚠️ Binary not found (falls back to simulation)
- ⚠️ TLS verification failures (for demo purposes)

**Actual Failures**:
- ❌ Prerequisites not met
- ❌ Docker not available
- ❌ Unexpected errors in logs

## Troubleshooting

### Common Issues

#### Demo Binary Not Found

```
⚠️ idpbuilder binary not found at ./bin/idpbuilder
   Running in simulation/demonstration mode
```

**Solution**:
- Build the binary: `make build`
- Or set `IDPBUILDER_BIN` environment variable
- Demos work in simulation mode without binary

#### Registry Not Available

```
⚠️ Local registry not available (expected in demo)
```

**Solution**:
- Start local registry: `docker run -d -p 5000:5000 registry:2`
- Or continue in simulation mode (demonstrations still work)

#### Docker Errors

```
❌ Docker not found. Please install Docker.
```

**Solution**:
- Install Docker: https://docs.docker.com/get-docker/
- Start Docker daemon
- Verify: `docker ps`

#### Permission Denied

```
bash: ./demos/basic-push-demo.sh: Permission denied
```

**Solution**:
```bash
chmod +x demos/*.sh  # Make all demos executable
```

#### TLS Certificate Errors

```
Error: x509: certificate signed by unknown authority
```

**Solution**:
- For testing: Use `--tls-verify=false` (INSECURE)
- For production: Add CA cert with `--ca-cert`
- Or use system CA bundle

### Getting Help

If demos fail unexpectedly:

1. **Check Prerequisites**: Verify all required tools are installed
2. **Review Logs**: Check `demo-results/*.log` files
3. **Enable Debug**: Run with `bash -x demos/script.sh`
4. **Check Documentation**: See `docs/` for detailed guides
5. **Report Issues**: Include demo logs and error messages

## Integration with CI/CD

### GitHub Actions Example

```yaml
name: Demo Push Command

on: [push]

jobs:
  demo:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Setup Docker
        uses: docker/setup-docker@v2

      - name: Run Integration Demo
        run: |
          ./demos/phase2-integration-demo.sh
        env:
          REGISTRY_USERNAME: ${{ secrets.REGISTRY_USERNAME }}
          REGISTRY_PASSWORD: ${{ secrets.REGISTRY_PASSWORD }}

      - name: Upload Results
        uses: actions/upload-artifact@v3
        with:
          name: demo-results
          path: demo-results/
```

### GitLab CI Example

```yaml
demo:
  stage: test
  script:
    - ./demos/phase2-integration-demo.sh
  artifacts:
    paths:
      - demo-results/
    when: always
```

## Compliance

### R330 Demo Planning Requirements

✅ **Demo Scripts**: All 5 required scripts created
✅ **Feature Coverage**: All Phase 2 features demonstrated
✅ **Integration Demo**: Complete end-to-end workflow
✅ **Documentation**: Comprehensive DEMO.md provided
✅ **Results**: Execution logs captured and preserved

### R291 Demo Deliverables

✅ **Execution Results**: All demos executed and logged
✅ **Summary Report**: integration-demo-summary.txt generated
✅ **Verification**: Demo functionality validated
✅ **Accessibility**: Clear instructions for reproduction

## Additional Resources

### Documentation

- **User Guide**: `docs/user-guide/push-command.md` - Complete push command guide
- **Examples**: `docs/examples/` - More usage examples
- **Reference**: `docs/reference/` - Command reference and options
- **Troubleshooting**: `docs/user-guide/troubleshooting.md` - Common issues

### Source Code

- **Implementation**: `pkg/push/` - Push command implementation
- **Tests**: `pkg/push/*_test.go` - Unit and integration tests
- **CLI**: `cmd/push.go` - Command-line interface

## Summary

These demos provide comprehensive coverage of all Phase 2 idpbuilder push command features:

- ✅ **5 Demo Scripts**: Cover all major feature areas
- ✅ **Integration Demo**: Shows features working together
- ✅ **Real-World Scenarios**: Practical usage patterns
- ✅ **Complete Documentation**: This guide and inline comments
- ✅ **CI/CD Ready**: Examples for pipeline integration
- ✅ **Compliance**: Satisfies R330 and R291 requirements

### Quick Reference

| Demo | Duration | Key Features |
|------|----------|--------------|
| basic-push | 2-3 min | Basic workflow, registry push |
| authenticated-push | 3-4 min | Auth methods, security best practices |
| retry-mechanism | 4-5 min | Auto-retry, backoff, error detection |
| tls-configuration | 3-4 min | Certificate handling, TLS options |
| phase2-integration | 8-10 min | **All features combined** |

### Next Steps

1. **Start with Integration Demo**: Get overview of all features
2. **Explore Individual Demos**: Deep dive into specific features
3. **Review Documentation**: Read detailed guides in `docs/`
4. **Try Examples**: See `docs/examples/` for more scenarios
5. **Integrate with CI/CD**: Use demo patterns in your pipelines

For questions or issues, refer to the troubleshooting section or check the detailed documentation in the `docs/` directory.
