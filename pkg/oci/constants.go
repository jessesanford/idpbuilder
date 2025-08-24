package oci

// Media Type constants from OCI Image Specification.
const (
	MediaTypeManifest     = "application/vnd.oci.image.manifest.v1+json"
	MediaTypeManifestList = "application/vnd.oci.image.index.v1+json"
	MediaTypeConfig       = "application/vnd.oci.image.config.v1+json"
	MediaTypeLayer        = "application/vnd.oci.image.layer.v1.tar+gzip"
	MediaTypeEmptyJSON    = "application/vnd.oci.empty.v1+json"
)

// Standard annotation keys from OCI Image Specification.
const (
	AnnotationCreated       = "org.opencontainers.image.created"
	AnnotationAuthors       = "org.opencontainers.image.authors"
	AnnotationURL           = "org.opencontainers.image.url"
	AnnotationSource        = "org.opencontainers.image.source"
	AnnotationVersion       = "org.opencontainers.image.version"
	AnnotationRevision      = "org.opencontainers.image.revision"
	AnnotationVendor        = "org.opencontainers.image.vendor"
	AnnotationLicenses      = "org.opencontainers.image.licenses"
	AnnotationTitle         = "org.opencontainers.image.title"
	AnnotationDescription   = "org.opencontainers.image.description"
)

// Default registry and platform constants.
const (
	DefaultRegistry   = "docker.io"
	DefaultTag        = "latest"
	DefaultNamespace  = "library"
	SchemaVersion     = 2
)

// Platform constants for common OS/arch combinations.
const (
	PlatformLinuxAMD64   = "linux/amd64"
	PlatformLinuxARM64   = "linux/arm64"
	PlatformLinuxARM     = "linux/arm"
	PlatformWindowsAMD64 = "windows/amd64"
	PlatformDarwinAMD64  = "darwin/amd64"
	PlatformDarwinARM64  = "darwin/arm64"
)

// Operating system constants.
const (
	OSLinux   = "linux"
	OSWindows = "windows"
	OSDarwin  = "darwin"
)

// Architecture constants.
const (
	ArchAMD64  = "amd64"
	ArchARM64  = "arm64"
	ArchARM    = "arm"
	Arch386    = "386"
)