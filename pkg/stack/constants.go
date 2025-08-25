package stack

// Stack component types define the different types of components in a stack
const (
	// ComponentTypeApplication represents application components
	ComponentTypeApplication = "application"
	
	// ComponentTypeInfrastructure represents infrastructure components
	ComponentTypeInfrastructure = "infrastructure"
	
	// ComponentTypeService represents service components
	ComponentTypeService = "service"
	
	// ComponentTypeConfig represents configuration components
	ComponentTypeConfig = "config"
	
	// ComponentTypeSecret represents secret/credential components
	ComponentTypeSecret = "secret"
)

// Stack component states define the lifecycle states of components
const (
	// ComponentStatePending indicates component is waiting to be processed
	ComponentStatePending = "pending"
	
	// ComponentStateProcessing indicates component is being deployed
	ComponentStateProcessing = "processing"
	
	// ComponentStateReady indicates component is successfully deployed
	ComponentStateReady = "ready"
	
	// ComponentStateFailed indicates component deployment failed
	ComponentStateFailed = "failed"
	
	// ComponentStateTerminating indicates component is being removed
	ComponentStateTerminating = "terminating"
)

// Stack validation constants
const (
	// MaxStackNameLength is the maximum allowed length for stack names
	MaxStackNameLength = 63
	
	// MaxComponentsPerStack is the maximum number of components per stack
	MaxComponentsPerStack = 100
	
	// MaxDependencyDepth is the maximum depth of component dependencies
	MaxDependencyDepth = 10
)

// Default stack configuration values
const (
	// DefaultStackNamespace is the default Kubernetes namespace for stacks
	DefaultStackNamespace = "idp-system"
	
	// DefaultStackTimeout is the default timeout for stack operations (in seconds)
	DefaultStackTimeout = 600
	
	// DefaultRetryAttempts is the default number of retry attempts
	DefaultRetryAttempts = 3
)