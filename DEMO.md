# Auth Flow Demo Documentation

## Overview
This demo showcases the authentication flow implementation that integrates flag overrides with secret-based credentials. The flow provides proper precedence where command-line flags take priority over Kubernetes secrets.

## Demo Features
- **Flag Override**: Command-line credentials override Kubernetes secrets
- **Secret Fallback**: Kubernetes secrets used when no flags provided
- **Error Handling**: Graceful error when no credentials available
- **Debug Logging**: Clear logs showing credential source

## Running the Demo

### Prerequisites
- Kubernetes cluster access (for secret scenarios)
- Authentication flow implementation in pkg/oci/flow.go

### Demo Scenarios

#### Scenario 1: Flag Override
```bash
./demo-auth-flow.sh --with-flags
```
Demonstrates that flag credentials take precedence over secrets.

#### Scenario 2: Secret Fallback
```bash
./demo-auth-flow.sh --with-secrets
```
Shows fallback to Kubernetes secrets when no flags provided.

#### Scenario 3: No Credentials Error
```bash
./demo-auth-flow.sh --no-creds
```
Demonstrates proper error handling when no credentials are available.

## Integration with Other Systems
- Export `DEMO_READY=true` when complete
- Entry point: `./demo-auth-flow.sh`
- Batch mode: `DEMO_BATCH=true ./demo-auth-flow.sh`

## Expected Outputs
Each scenario produces clear output showing the credential source and validation that the correct precedence is followed.