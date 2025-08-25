// Package pkg provides OCI (Open Container Initiative) types and utilities
// for idpbuilder-oci-mgmt.
//
// This package includes comprehensive support for OCI image references,
// manifests, descriptors, and registry operations. It provides foundational
// types for working with container images and registries in the context
// of infrastructure development platform management.
//
// Key components:
//   - OCIReference: Container image reference parsing and validation
//   - OCIImage: Container image metadata and operations
//   - OCIManifest: OCI manifest handling and validation
//   - OCIDescriptor: OCI descriptor types for content addressing
//
// The types support standard OCI specifications and provide safe abstractions
// for common container registry operations including authentication, manifest
// parsing, and image reference handling.
//
// This package is designed to integrate with Kubernetes and container runtime
// environments while maintaining compatibility with standard OCI tooling.
//
// Example usage:
//   ref, err := oci.ParseOCIReference("registry.example.com/myapp:v1.0.0")
//   if err != nil {
//       log.Fatal(err)
//   }
//   
//   image := &oci.OCIImage{
//       Reference: ref,
//       Registry: "registry.example.com",
//       Repository: "myapp",
//       Tag: "v1.0.0",
//   }
//
// For stack and build integration, see the stack package which extends
// these foundation types with deployment-specific functionality.
package pkg