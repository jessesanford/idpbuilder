# Certificate Validation Test Fixtures

This directory contains test certificate files used by the certificate validation tests.

## Files

- `valid_cert.pem` - A valid test certificate for basic validation tests
- `root_ca.pem` - A test root CA certificate for chain validation tests
- `intermediate_ca.pem` - A test intermediate CA certificate for multi-level chains
- `chain.pem` - A complete certificate chain (leaf -> intermediate -> root)
- `invalid_cert.pem` - An invalid certificate for testing error handling

## Usage

These fixtures are used by the tests in the parent directory to ensure robust certificate validation functionality without needing to generate certificates at runtime.

## Security Note

These are test certificates only and should never be used in production environments.