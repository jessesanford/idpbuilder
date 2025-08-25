package stack

// Stack Status Constants
const (
	// StackStatusPending indicates the stack is being initialized
	StackStatusPending StackStatus = "pending"

	// StackStatusReady indicates the stack is fully deployed and operational
	StackStatusReady StackStatus = "ready"

	// StackStatusFailed indicates the stack deployment failed
	StackStatusFailed StackStatus = "failed"

	// StackStatusUpdating indicates the stack is being updated
	StackStatusUpdating StackStatus = "updating"

	// StackStatusDestroying indicates the stack is being destroyed
	StackStatusDestroying StackStatus = "destroying"
)

// Component Status Constants
const (
	// ComponentStatusPending indicates the component is being initialized
	ComponentStatusPending ComponentStatus = "pending"

	// ComponentStatusRunning indicates the component is running successfully
	ComponentStatusRunning ComponentStatus = "running"

	// ComponentStatusFailed indicates the component failed to start or run
	ComponentStatusFailed ComponentStatus = "failed"

	// ComponentStatusStopped indicates the component is stopped
	ComponentStatusStopped ComponentStatus = "stopped"

	// ComponentStatusUpdating indicates the component is being updated
	ComponentStatusUpdating ComponentStatus = "updating"
)

// Component Type Constants
const (
	// ComponentTypeApplication represents an application component
	ComponentTypeApplication ComponentType = "application"

	// ComponentTypeDatabase represents a database component
	ComponentTypeDatabase ComponentType = "database"

	// ComponentTypeService represents a service component
	ComponentTypeService ComponentType = "service"

	// ComponentTypeMiddleware represents middleware components (cache, message queue, etc.)
	ComponentTypeMiddleware ComponentType = "middleware"

	// ComponentTypeProxy represents proxy or load balancer components
	ComponentTypeProxy ComponentType = "proxy"

	// ComponentTypeMonitoring represents monitoring and observability components
	ComponentTypeMonitoring ComponentType = "monitoring"
)

// Dependency Type Constants
const (
	// DependencyTypeRuntime indicates a runtime dependency (required at runtime)
	DependencyTypeRuntime DependencyType = "runtime"

	// DependencyTypeBuild indicates a build-time dependency
	DependencyTypeBuild DependencyType = "build"

	// DependencyTypeOptional indicates an optional dependency
	DependencyTypeOptional DependencyType = "optional"

	// DependencyTypeNetwork indicates a network dependency
	DependencyTypeNetwork DependencyType = "network"
)

// Default Configuration Keys
const (
	// ConfigKeyReplicas specifies the number of replicas for a component
	ConfigKeyReplicas = "replicas"

	// ConfigKeyMemoryLimit specifies the memory limit for a component
	ConfigKeyMemoryLimit = "memoryLimit"

	// ConfigKeyCPULimit specifies the CPU limit for a component
	ConfigKeyCPULimit = "cpuLimit"

	// ConfigKeyEnvironment specifies environment variables
	ConfigKeyEnvironment = "environment"

	// ConfigKeyPorts specifies exposed ports
	ConfigKeyPorts = "ports"

	// ConfigKeyVolumes specifies volume mounts
	ConfigKeyVolumes = "volumes"
)

// Version Constraint Operators
const (
	// ExactVersionConstraint requires an exact version match
	ExactVersionConstraint = "="

	// MinimumVersionConstraint requires a minimum version
	MinimumVersionConstraint = ">="

	// CompatibleVersionConstraint requires a compatible version (tilde constraint)
	CompatibleVersionConstraint = "~"

	// PessimisticVersionConstraint requires a compatible version (caret constraint)
	PessimisticVersionConstraint = "^"
)