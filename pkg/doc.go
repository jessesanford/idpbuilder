/*
Package pkg provides core types for the idpbuilder OCI management system.

Contains two sub-packages:

- oci: Types for OCI (Open Container Initiative) images, manifests, and repositories
- stack: Types for defining software stacks composed of OCI images

OCI Package provides:
  - OCIImage: Container image representation with metadata
  - OCIManifest: Image manifest following OCI spec
  - OCIDescriptor: Content descriptor for registry artifacts  
  - OCIRepository: Interface for repository operations

Stack Package provides:
  - StackConfiguration: Complete stack definition
  - StackComponent: Individual component within a stack
  - StackDependency: External dependencies

Example:

	ref := oci.OCIReference{
		Registry:   "docker.io",
		Repository: "library/nginx", 
		Tag:        "latest",
	}
	
	stack := &stack.StackConfiguration{
		Name:    "web-app",
		Version: "1.0.0",
		Components: []stack.StackComponent{
			{
				Name:         "frontend",
				Type:         stack.ComponentTypeService,
				OCIReference: ref,
			},
		},
	}
*/
package pkg