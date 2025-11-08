// Package registry provides interfaces and implementations for pushing
// OCI images to container registries.
//
// This package supports:
//   - Pushing images to OCI-compatible registries
//   - Progress reporting during layer uploads
//   - Registry validation and connectivity checks
//   - Building fully qualified image references
//
// The primary interface is Client, which abstracts registry operations.
// Implementations use go-containerregistry for OCI compatibility.
package registry
