// Package auth provides types and interfaces for OCI registry authentication.
//
// This package defines the core authentication types used throughout the
// idpbuilder-oci-mgmt system for connecting to and authenticating with
// OCI-compliant container registries.
//
// # Authentication Types
//
// The package supports multiple authentication methods:
//   - Basic Authentication (username/password)
//   - Bearer Token Authentication
//   - OAuth 2.0 Authentication
//   - Registry Identity Tokens
//
// # Basic Usage
//
// Creating basic authentication configuration:
//
//	config := &auth.AuthConfig{
//		Registry:  "registry.example.com",
//		Username:  "myuser",
//		Password:  "mypass",
//		AuthType:  auth.AuthTypeBasic,
//	}
//
//	creds, err := config.GetCredentials("registry.example.com")
//	if err != nil {
//		log.Fatal(err)
//	}
//
// # Bearer Token Authentication
//
// Using bearer tokens for registry access:
//
//	config := &auth.AuthConfig{
//		Registry: "gcr.io",
//		Token:    "bearer-token-here",
//		AuthType: auth.AuthTypeBearer,
//	}
//
// # Credential Storage
//
// The package provides thread-safe credential storage:
//
//	store := auth.NewCredentialStore()
//	err := store.Set("registry.example.com", creds)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	retrievedCreds, err := store.Get("registry.example.com")
//	if err != nil {
//		log.Fatal(err)
//	}
//
// # Certificate Management
//
// For secure connections and mutual TLS:
//
//	bundle := &certs.CertificateBundle{
//		CACert:     caCert,
//		ClientCert: clientCert,
//		ClientKey:  clientKey,
//		Registry:   "secure-registry.example.com",
//	}
//
//	tlsConfig := &certs.TLSConfig{
//		InsecureSkipVerify: false,
//		ServerName:         "secure-registry.example.com",
//		MinVersion:         certs.MinTLSVersion,
//	}
//
// # Security Considerations
//
// This package is designed with security best practices:
//   - Sensitive data (passwords, tokens) should be cleared from memory after use
//   - Certificate validation is enabled by default
//   - Secure string comparison is used for credential verification
//   - Credentials are never logged or exposed in error messages
//
// # Error Handling
//
// The package defines specific error types for different failure scenarios:
//   - ErrInvalidCredentials: Invalid or incomplete credential data
//   - ErrCredentialsNotFound: No credentials found for registry
//   - ErrAuthenticationFailed: Authentication with registry failed
//   - ErrTokenExpired: Authentication token has expired
//
// Always check for these specific errors when handling authentication failures.
package auth