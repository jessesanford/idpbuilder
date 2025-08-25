// Package stack provides types and utilities for managing multi-component application stacks.
//
// This package defines the core stack configuration types used throughout the
// idpbuilder-oci-mgmt system for defining, managing, and deploying complex
// application stacks composed of multiple interconnected components.
//
// # Stack Components
//
// The package supports various component types:
//   - Applications (web services, APIs, workers)
//   - Databases (SQL, NoSQL, time-series)
//   - Services (message queues, caches, proxies)
//   - Middleware (authentication, monitoring, logging)
//   - Infrastructure components (load balancers, storage)
//
// # Basic Usage
//
// Creating a simple stack configuration:
//
//	stack := &stack.StackConfiguration{
//		Name:    "web-stack",
//		Version: "1.0.0",
//		Status:  stack.StackStatusPending,
//		Components: []stack.StackComponent{
//			{
//				Name:    "web-app",
//				Type:    stack.ComponentTypeApplication,
//				Version: "2.1.0",
//				Status:  stack.ComponentStatusPending,
//				OCIReference: &stack.OCIReference{
//					Registry:   "docker.io",
//					Repository: "myorg/webapp",
//					Tag:        "v2.1.0",
//				},
//			},
//			{
//				Name:    "postgres",
//				Type:    stack.ComponentTypeDatabase,
//				Version: "14.0",
//				Status:  stack.ComponentStatusPending,
//				Configuration: map[string]interface{}{
//					stack.ConfigKeyReplicas: 1,
//					stack.ConfigKeyMemoryLimit: "2Gi",
//				},
//			},
//		},
//		Dependencies: []stack.StackDependency{
//			{
//				Component: "web-app",
//				DependsOn: "postgres",
//				Type:      stack.DependencyTypeRuntime,
//			},
//		},
//	}
//
// # Component Management
//
// Retrieving and managing stack components:
//
//	component, err := stack.GetComponentByName("web-app")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	dependencies := stack.GetDependenciesFor("web-app")
//	for _, dep := range dependencies {
//		fmt.Printf("Component %s depends on %s\n", dep.Component, dep.DependsOn)
//	}
//
// # Validation
//
// Stack configurations include comprehensive validation:
//
//	if err := stack.IsValid(); err != nil {
//		log.Fatal("Invalid stack configuration:", err)
//	}
//
// # Status Tracking
//
// Components and stacks track their lifecycle status:
//   - Pending: Initial state, not yet deployed
//   - Ready/Running: Successfully deployed and operational
//   - Failed: Deployment or runtime failure
//   - Updating: In process of being updated
//   - Stopped/Destroying: Shutdown states
//
// # Configuration Management
//
// Components support flexible configuration via key-value maps:
//   - Resource limits (CPU, memory, storage)
//   - Scaling parameters (replicas, autoscaling)
//   - Environment variables
//   - Network and volume configurations
//
// The package defines standard configuration keys for common settings,
// but supports arbitrary key-value pairs for component-specific needs.
//
// # OCI Integration
//
// Stack components integrate with OCI images and artifacts:
//   - Reference OCI images for containerized components
//   - Support for registries, tags, and digest-based references
//   - Validation of OCI reference format and completeness
//
// This enables stacks to define precise versions of all components
// and ensure reproducible deployments across environments.
package stack