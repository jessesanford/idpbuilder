package stack

import (
	"testing"
	"time"
	
	"github.com/cnoe-io/idpbuilder/pkg/oci"
)

func TestStackConfiguration_Creation(t *testing.T) {
	now := time.Now()
	stack := StackConfiguration{
		Name:      "test-stack",
		Version:   "1.0.0",
		CreatedAt: now,
		Status:    StackStatusPending,
	}

	if stack.Name != "test-stack" {
		t.Errorf("expected name test-stack, got %s", stack.Name)
	}
	if stack.Status != StackStatusPending {
		t.Errorf("expected status pending, got %s", stack.Status)
	}
}

func TestStackComponent_Creation(t *testing.T) {
	component := StackComponent{
		Name:    "web-server",
		Type:    ComponentTypeService,
		Version: "1.2.3",
		OCIReference: oci.OCIReference{
			Registry:   "docker.io",
			Repository: "nginx",
			Tag:        "latest",
		},
		Status: ComponentStatusRunning,
	}

	if component.Type != ComponentTypeService {
		t.Errorf("expected type %s, got %s", ComponentTypeService, component.Type)
	}
	if component.OCIReference.Repository != "nginx" {
		t.Errorf("expected repository nginx, got %s", component.OCIReference.Repository)
	}
}

func TestStackPort_Creation(t *testing.T) {
	port := StackPort{
		Name:     "http",
		Port:     80,
		Protocol: ProtocolTCP,
		External: true,
	}

	if port.Port != 80 {
		t.Errorf("expected port 80, got %d", port.Port)
	}
	if port.Protocol != ProtocolTCP {
		t.Errorf("expected protocol %s, got %s", ProtocolTCP, port.Protocol)
	}
	if !port.External {
		t.Error("expected external to be true")
	}
}

func TestStackVolume_Creation(t *testing.T) {
	volume := StackVolume{
		Name:     "data",
		Path:     "/var/lib/data",
		Type:     VolumeTypePersistent,
		Size:     "10Gi",
		ReadOnly: false,
	}

	if volume.Type != VolumeTypePersistent {
		t.Errorf("expected type %s, got %s", VolumeTypePersistent, volume.Type)
	}
	if volume.Size != "10Gi" {
		t.Errorf("expected size 10Gi, got %s", volume.Size)
	}
	if volume.ReadOnly {
		t.Error("expected read-only to be false")
	}
}

func TestStackHealthCheck_Creation(t *testing.T) {
	healthCheck := StackHealthCheck{
		Type:                HealthCheckTypeHTTP,
		Endpoint:            "http://localhost:8080/health",
		IntervalSeconds:     30,
		TimeoutSeconds:      5,
		InitialDelaySeconds: 10,
	}

	if healthCheck.Type != HealthCheckTypeHTTP {
		t.Errorf("expected type %s, got %s", HealthCheckTypeHTTP, healthCheck.Type)
	}
	if healthCheck.IntervalSeconds != 30 {
		t.Errorf("expected interval 30, got %d", healthCheck.IntervalSeconds)
	}
}

func TestStackResources_Creation(t *testing.T) {
	resources := StackResources{
		CPU:         "100m",
		Memory:      "128Mi",
		CPULimit:    "500m",
		MemoryLimit: "512Mi",
	}

	if resources.CPU != "100m" {
		t.Errorf("expected CPU 100m, got %s", resources.CPU)
	}
	if resources.MemoryLimit != "512Mi" {
		t.Errorf("expected memory limit 512Mi, got %s", resources.MemoryLimit)
	}
}

func TestStackDependency_Creation(t *testing.T) {
	dependency := StackDependency{
		Name:      "database",
		Type:      "postgresql",
		Version:   "13.0",
		Reference: "postgres:13",
		Optional:  false,
	}

	if dependency.Name != "database" {
		t.Errorf("expected name database, got %s", dependency.Name)
	}
	if dependency.Optional {
		t.Error("expected optional to be false")
	}
}

func TestStackStatus_Constants(t *testing.T) {
	statuses := []StackStatus{
		StackStatusPending,
		StackStatusDeploying,
		StackStatusRunning,
		StackStatusFailed,
		StackStatusStopped,
		StackStatusUpdating,
	}

	if len(statuses) != 6 {
		t.Errorf("expected 6 status constants, got %d", len(statuses))
	}
}

func TestComponentStatus_Constants(t *testing.T) {
	statuses := []ComponentStatus{
		ComponentStatusPending,
		ComponentStatusStarting,
		ComponentStatusRunning,
		ComponentStatusFailed,
		ComponentStatusStopped,
		ComponentStatusUnhealthy,
	}

	if len(statuses) != 6 {
		t.Errorf("expected 6 component status constants, got %d", len(statuses))
	}
}