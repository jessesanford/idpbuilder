// Package stack provides types for defining software stacks composed of OCI images.
package stack

import (
	"time"
	
	"github.com/cnoe-io/idpbuilder/pkg/oci"
)

// StackConfiguration represents a complete software stack definition.
type StackConfiguration struct {
	Name         string              `json:"name" yaml:"name"`
	Version      string              `json:"version" yaml:"version"`
	Description  string              `json:"description,omitempty" yaml:"description,omitempty"`
	Components   []StackComponent    `json:"components" yaml:"components"`
	Dependencies []StackDependency   `json:"dependencies,omitempty" yaml:"dependencies,omitempty"`
	Metadata     map[string]string   `json:"metadata,omitempty" yaml:"metadata,omitempty"`
	CreatedAt    time.Time           `json:"created_at" yaml:"created_at"`
	UpdatedAt    time.Time           `json:"updated_at" yaml:"updated_at"`
	Status       StackStatus         `json:"status" yaml:"status"`
}

// StackComponent represents a single component within a software stack.
type StackComponent struct {
	Name          string                    `json:"name" yaml:"name"`
	Type          string                    `json:"type" yaml:"type"`
	Version       string                    `json:"version" yaml:"version"`
	OCIReference  oci.OCIReference          `json:"oci_reference" yaml:"oci_reference"`
	Configuration map[string]interface{}    `json:"configuration,omitempty" yaml:"configuration,omitempty"`
	Dependencies  []string                  `json:"dependencies,omitempty" yaml:"dependencies,omitempty"`
	Ports         []StackPort               `json:"ports,omitempty" yaml:"ports,omitempty"`
	Volumes       []StackVolume             `json:"volumes,omitempty" yaml:"volumes,omitempty"`
	Environment   map[string]string         `json:"environment,omitempty" yaml:"environment,omitempty"`
	HealthCheck   *StackHealthCheck         `json:"health_check,omitempty" yaml:"health_check,omitempty"`
	Resources     *StackResources           `json:"resources,omitempty" yaml:"resources,omitempty"`
	Status        ComponentStatus           `json:"status" yaml:"status"`
}

// StackPort defines a network port exposed by a stack component.
type StackPort struct {
	Name       string `json:"name,omitempty" yaml:"name,omitempty"`
	Port       int    `json:"port" yaml:"port"`
	Protocol   string `json:"protocol,omitempty" yaml:"protocol,omitempty"`
	TargetPort int    `json:"target_port,omitempty" yaml:"target_port,omitempty"`
	External   bool   `json:"external,omitempty" yaml:"external,omitempty"`
}

// StackVolume defines a storage volume used by a stack component.
type StackVolume struct {
	Name     string `json:"name" yaml:"name"`
	Path     string `json:"path" yaml:"path"`
	Type     string `json:"type" yaml:"type"`
	Size     string `json:"size,omitempty" yaml:"size,omitempty"`
	ReadOnly bool   `json:"read_only,omitempty" yaml:"read_only,omitempty"`
}

// StackHealthCheck defines health check configuration for a component.
type StackHealthCheck struct {
	Type                string   `json:"type" yaml:"type"`
	Endpoint            string   `json:"endpoint,omitempty" yaml:"endpoint,omitempty"`
	Command             []string `json:"command,omitempty" yaml:"command,omitempty"`
	IntervalSeconds     int      `json:"interval_seconds,omitempty" yaml:"interval_seconds,omitempty"`
	TimeoutSeconds      int      `json:"timeout_seconds,omitempty" yaml:"timeout_seconds,omitempty"`
	InitialDelaySeconds int      `json:"initial_delay_seconds,omitempty" yaml:"initial_delay_seconds,omitempty"`
}

// StackResources defines resource requirements and limits for a component.
type StackResources struct {
	CPU         string `json:"cpu,omitempty" yaml:"cpu,omitempty"`
	Memory      string `json:"memory,omitempty" yaml:"memory,omitempty"`
	CPULimit    string `json:"cpu_limit,omitempty" yaml:"cpu_limit,omitempty"`
	MemoryLimit string `json:"memory_limit,omitempty" yaml:"memory_limit,omitempty"`
}

// StackDependency represents an external dependency required by the stack.
type StackDependency struct {
	Name      string `json:"name" yaml:"name"`
	Type      string `json:"type" yaml:"type"`
	Version   string `json:"version" yaml:"version"`
	Reference string `json:"reference,omitempty" yaml:"reference,omitempty"`
	Optional  bool   `json:"optional,omitempty" yaml:"optional,omitempty"`
}

// StackStatus represents the overall status of a stack deployment.
type StackStatus string

// Stack status constants.
const (
	StackStatusPending   StackStatus = "pending"
	StackStatusDeploying StackStatus = "deploying"
	StackStatusRunning   StackStatus = "running"
	StackStatusFailed    StackStatus = "failed"
	StackStatusStopped   StackStatus = "stopped"
	StackStatusUpdating  StackStatus = "updating"
)

// ComponentStatus represents the status of an individual stack component.
type ComponentStatus string

// Component status constants.
const (
	ComponentStatusPending   ComponentStatus = "pending"
	ComponentStatusStarting  ComponentStatus = "starting"
	ComponentStatusRunning   ComponentStatus = "running"
	ComponentStatusFailed    ComponentStatus = "failed"
	ComponentStatusStopped   ComponentStatus = "stopped"
	ComponentStatusUnhealthy ComponentStatus = "unhealthy"
)