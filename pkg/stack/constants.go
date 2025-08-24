package stack

// Component type constants.
const (
	ComponentTypeService     = "service"
	ComponentTypeDatabase    = "database"
	ComponentTypeCache       = "cache"
	ComponentTypeProxy       = "proxy"
	ComponentTypeMonitoring  = "monitoring"
	ComponentTypeStorage     = "storage"
)

// Volume type constants.
const (
	VolumeTypePersistent = "persistent"
	VolumeTypeEphemeral  = "ephemeral"
	VolumeTypeHostPath   = "host-path"
	VolumeTypeConfigMap  = "config-map"
	VolumeTypeSecret     = "secret"
)

// Health check type constants.
const (
	HealthCheckTypeHTTP = "http"
	HealthCheckTypeTCP  = "tcp"
	HealthCheckTypeExec = "exec"
)

// Network protocol constants.
const (
	ProtocolTCP = "TCP"
	ProtocolUDP = "UDP"
)

// Default configuration values.
const (
	DefaultHealthCheckInterval = 30
	DefaultHealthCheckTimeout  = 5
	DefaultCPURequest          = "100m"
	DefaultMemoryRequest       = "128Mi"
	DefaultCPULimit            = "500m"
	DefaultMemoryLimit         = "512Mi"
)