# Phase 1 Wave 2 Integration Demo

## Overview
This demo showcases the integrated Wave 2 features for the idpbuilder push command with OCI registry support.

## Features Demonstrated

### 1. Push Command Structure (E1.2.1)
- Comprehensive CLI interface
- Help documentation
- Flag parsing and validation

### 2. Registry Authentication (E1.2.2)
- Username/password authentication
- Environment variable support
- Insecure mode for self-signed certificates

### 3. Retry Logic (E1.2.2)
- Exponential backoff strategy
- Constant backoff strategy
- Configurable retry parameters

### 4. Image Push Operations (E1.2.3)
- Image discovery from tarballs
- Registry operations
- Push implementation

## How to Run

```bash
# Make the demo script executable (if not already)
chmod +x demos/demo-wave2.sh

# Run the demo
./demos/demo-wave2.sh
```

## Expected Output
The demo will:
1. Display push command help
2. Show authentication options
3. Demonstrate retry configurations
4. Execute a dry-run push operation

## Integration Success Criteria
✅ All 6 Wave 2 efforts merged  
✅ BUG-007 (PushCmd redeclared) resolved  
✅ Build passes without errors  
✅ Demo script executes successfully  
