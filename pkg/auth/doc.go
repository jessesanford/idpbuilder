// Package auth provides interfaces for registry authentication.
//
// This package supports:
//   - Basic username/password authentication
//   - Credential validation
//   - Integration with go-containerregistry authn
//
// The primary interface is Provider, which supplies authentication
// to registry clients.
package auth
