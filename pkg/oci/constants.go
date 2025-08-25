package oci

// OCI Media Types as defined in the OCI Image Specification
const (
	// MediaTypeManifest specifies the media type for a manifest
	MediaTypeManifest = "application/vnd.oci.image.manifest.v1+json"
	
	// MediaTypeImageIndex specifies the media type for an image index
	MediaTypeImageIndex = "application/vnd.oci.image.index.v1+json"
	
	// MediaTypeImageConfig specifies the media type for image config
	MediaTypeImageConfig = "application/vnd.oci.image.config.v1+json"
	
	// MediaTypeImageLayer specifies the media type for image layers
	MediaTypeImageLayer = "application/vnd.oci.image.layer.v1.tar+gzip"
	
	// MediaTypeImageLayerNonDistributable specifies non-distributable layers
	MediaTypeImageLayerNonDistributable = "application/vnd.oci.image.layer.nondistributable.v1.tar+gzip"
	
	// MediaTypeDescriptor specifies the media type for descriptors
	MediaTypeDescriptor = "application/vnd.oci.descriptor.v1+json"
)

// Docker compatibility media types for legacy registries
const (
	// DockerMediaTypeManifest for Docker v2 Schema 2 manifests
	DockerMediaTypeManifest = "application/vnd.docker.distribution.manifest.v2+json"
	
	// DockerMediaTypeManifestList for Docker v2 manifest lists
	DockerMediaTypeManifestList = "application/vnd.docker.distribution.manifest.list.v2+json"
	
	// DockerMediaTypeConfig for Docker v2 image config
	DockerMediaTypeConfig = "application/vnd.docker.container.image.v1+json"
	
	// DockerMediaTypeLayer for Docker v2 layers
	DockerMediaTypeLayer = "application/vnd.docker.image.rootfs.diff.tar.gzip"
)

// Registry defaults and configuration constants
const (
	// DefaultRegistry is the default container registry
	DefaultRegistry = "docker.io"
	
	// DefaultNamespace is the default namespace for single-name images
	DefaultNamespace = "library"
	
	// DefaultTag is the default tag when none is specified
	DefaultTag = "latest"
	
	// DefaultDigestAlgorithm is the default algorithm for content digests
	DefaultDigestAlgorithm = "sha256"
	
	// MaxManifestSize is the maximum allowed manifest size in bytes (4MB)
	MaxManifestSize = 4 * 1024 * 1024
	
	// MaxLayerSize is the maximum allowed layer size in bytes (10GB)
	MaxLayerSize = 10 * 1024 * 1024 * 1024
)

// Registry port constants
const (
	// DefaultHTTPPort for insecure registry connections
	DefaultHTTPPort = "80"
	
	// DefaultHTTPSPort for secure registry connections  
	DefaultHTTPSPort = "443"
	
	// DefaultRegistryPort for Docker registry protocol
	DefaultRegistryPort = "5000"
)

// OCI specification version constants
const (
	// OCIImageSpecVersion is the supported OCI image spec version
	OCIImageSpecVersion = "1.0.1"
	
	// OCIRuntimeSpecVersion is the supported OCI runtime spec version
	OCIRuntimeSpecVersion = "1.0.2"
	
	// OCIDistributionSpecVersion is the supported OCI distribution spec version
	OCIDistributionSpecVersion = "1.0.1"
)