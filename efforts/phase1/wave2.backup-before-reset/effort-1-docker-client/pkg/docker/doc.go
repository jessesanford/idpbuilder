// Package docker provides interfaces and implementations for interacting
// with the Docker daemon to retrieve OCI images.
//
// This package enables:
//   - Checking if images exist in the local Docker daemon
//   - Retrieving images in OCI v1.Image format
//   - Validating image names per OCI specification
//   - Managing Docker client lifecycle
//
// The primary interface is Client, which abstracts all Docker operations.
// Implementations use the Docker Engine API client library.
//
// Example usage:
//
//	client, err := docker.NewClient()
//	if err != nil {
//	    return err
//	}
//	defer client.Close()
//
//	exists, err := client.ImageExists(ctx, "myapp:latest")
//	if err != nil {
//	    return err
//	}
//
//	if exists {
//	    image, err := client.GetImage(ctx, "myapp:latest")
//	    // use image...
//	}
package docker
