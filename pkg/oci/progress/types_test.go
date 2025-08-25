package progress

import (
	"testing"
	"time"
)

func TestProgressStatus_String(t *testing.T) {
	tests := []struct {
		status   ProgressStatus
		expected string
	}{
		{ProgressStatusStarted, "Started"},
		{ProgressStatusInProgress, "InProgress"},
		{ProgressStatusCompleted, "Completed"},
		{ProgressStatusFailed, "Failed"},
		{ProgressStatusCanceled, "Canceled"},
		{ProgressStatus(999), "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := tt.status.String(); got != tt.expected {
				t.Errorf("ProgressStatus.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestNewProgressEvent(t *testing.T) {
	id := "test-id"
	operation := "test-operation"
	phase := "test-phase"
	status := ProgressStatusInProgress

	event := NewProgressEvent(id, operation, phase, status)

	if event.ID != id {
		t.Errorf("Expected ID %s, got %s", id, event.ID)
	}
	if event.Operation != operation {
		t.Errorf("Expected Operation %s, got %s", operation, event.Operation)
	}
	if event.Phase != phase {
		t.Errorf("Expected Phase %s, got %s", phase, event.Phase)
	}
	if event.Status != status {
		t.Errorf("Expected Status %v, got %v", status, event.Status)
	}
	if event.Timestamp.IsZero() {
		t.Errorf("Expected timestamp to be set")
	}
}

func TestProgressEvent_WithProgress(t *testing.T) {
	event := NewProgressEvent("id", "operation", "phase", ProgressStatusInProgress)
	event.WithProgress(50, 100)

	if event.Current != 50 {
		t.Errorf("Expected Current 50, got %d", event.Current)
	}
	if event.Total != 100 {
		t.Errorf("Expected Total 100, got %d", event.Total)
	}
	if event.Percent != 50.0 {
		t.Errorf("Expected Percent 50.0, got %f", event.Percent)
	}
}

func TestProgressEvent_WithProgressZeroTotal(t *testing.T) {
	event := NewProgressEvent("id", "operation", "phase", ProgressStatusInProgress)
	event.WithProgress(50, 0)

	if event.Percent != 0.0 {
		t.Errorf("Expected Percent 0.0 when total is 0, got %f", event.Percent)
	}
}

func TestProgressEvent_WithMessage(t *testing.T) {
	message := "test message"
	event := NewProgressEvent("id", "operation", "phase", ProgressStatusInProgress)
	event.WithMessage(message)

	if event.Message != message {
		t.Errorf("Expected Message %s, got %s", message, event.Message)
	}
}

func TestProgressEvent_WithDuration(t *testing.T) {
	duration := 5 * time.Second
	event := NewProgressEvent("id", "operation", "phase", ProgressStatusInProgress)
	event.WithDuration(duration)

	if event.Duration != duration {
		t.Errorf("Expected Duration %v, got %v", duration, event.Duration)
	}
}

func TestProgressEvent_WithETA(t *testing.T) {
	eta := 30 * time.Second
	event := NewProgressEvent("id", "operation", "phase", ProgressStatusInProgress)
	event.WithETA(eta)

	if event.ETA == nil {
		t.Errorf("Expected ETA to be set")
	}
	if *event.ETA != eta {
		t.Errorf("Expected ETA %v, got %v", eta, *event.ETA)
	}
}

func TestProgressEvent_Chaining(t *testing.T) {
	event := NewProgressEvent("id", "operation", "phase", ProgressStatusInProgress).
		WithProgress(25, 100).
		WithMessage("test message").
		WithDuration(10 * time.Second).
		WithETA(30 * time.Second)

	if event.Current != 25 {
		t.Errorf("Expected Current 25, got %d", event.Current)
	}
	if event.Total != 100 {
		t.Errorf("Expected Total 100, got %d", event.Total)
	}
	if event.Percent != 25.0 {
		t.Errorf("Expected Percent 25.0, got %f", event.Percent)
	}
	if event.Message != "test message" {
		t.Errorf("Expected Message 'test message', got %s", event.Message)
	}
	if event.Duration != 10*time.Second {
		t.Errorf("Expected Duration 10s, got %v", event.Duration)
	}
	if event.ETA == nil || *event.ETA != 30*time.Second {
		t.Errorf("Expected ETA 30s, got %v", event.ETA)
	}
}

func TestNewBuildProgress(t *testing.T) {
	imageName := "test-image:latest"
	buildID := "build-123"
	totalSteps := 5

	bp := NewBuildProgress(imageName, buildID, totalSteps)

	if bp.ImageName != imageName {
		t.Errorf("Expected ImageName %s, got %s", imageName, bp.ImageName)
	}
	if bp.BuildID != buildID {
		t.Errorf("Expected BuildID %s, got %s", buildID, bp.BuildID)
	}
	if bp.TotalSteps != totalSteps {
		t.Errorf("Expected TotalSteps %d, got %d", totalSteps, bp.TotalSteps)
	}
	if bp.StartTime.IsZero() {
		t.Errorf("Expected StartTime to be set")
	}
	if bp.LastUpdate.IsZero() {
		t.Errorf("Expected LastUpdate to be set")
	}
}

func TestNewPushProgress(t *testing.T) {
	imageName := "test-image:latest"
	registry := "registry.example.com"
	pushID := "push-456"
	totalLayers := 3

	pp := NewPushProgress(imageName, registry, pushID, totalLayers)

	if pp.ImageName != imageName {
		t.Errorf("Expected ImageName %s, got %s", imageName, pp.ImageName)
	}
	if pp.Registry != registry {
		t.Errorf("Expected Registry %s, got %s", registry, pp.Registry)
	}
	if pp.PushID != pushID {
		t.Errorf("Expected PushID %s, got %s", pushID, pp.PushID)
	}
	if pp.TotalLayers != totalLayers {
		t.Errorf("Expected TotalLayers %d, got %d", totalLayers, pp.TotalLayers)
	}
	if pp.StartTime.IsZero() {
		t.Errorf("Expected StartTime to be set")
	}
	if pp.LastUpdate.IsZero() {
		t.Errorf("Expected LastUpdate to be set")
	}
}

func TestBuildProgressFields(t *testing.T) {
	bp := NewBuildProgress("image", "build-id", 10)

	// Test that all fields are properly initialized to zero values
	if bp.Step != "" {
		t.Errorf("Expected Step to be empty, got %s", bp.Step)
	}
	if bp.StepNumber != 0 {
		t.Errorf("Expected StepNumber to be 0, got %d", bp.StepNumber)
	}
	if bp.LayersBuilt != 0 {
		t.Errorf("Expected LayersBuilt to be 0, got %d", bp.LayersBuilt)
	}
	if bp.TotalLayers != 0 {
		t.Errorf("Expected TotalLayers to be 0, got %d", bp.TotalLayers)
	}
	if bp.BytesProcessed != 0 {
		t.Errorf("Expected BytesProcessed to be 0, got %d", bp.BytesProcessed)
	}
	if bp.TotalBytes != 0 {
		t.Errorf("Expected TotalBytes to be 0, got %d", bp.TotalBytes)
	}
}

func TestPushProgressFields(t *testing.T) {
	pp := NewPushProgress("image", "registry", "push-id", 5)

	// Test that all fields are properly initialized to zero values
	if pp.LayersPushed != 0 {
		t.Errorf("Expected LayersPushed to be 0, got %d", pp.LayersPushed)
	}
	if pp.LayersSkipped != 0 {
		t.Errorf("Expected LayersSkipped to be 0, got %d", pp.LayersSkipped)
	}
	if pp.BytesUploaded != 0 {
		t.Errorf("Expected BytesUploaded to be 0, got %d", pp.BytesUploaded)
	}
	if pp.TotalBytes != 0 {
		t.Errorf("Expected TotalBytes to be 0, got %d", pp.TotalBytes)
	}
	if pp.CurrentLayer != "" {
		t.Errorf("Expected CurrentLayer to be empty, got %s", pp.CurrentLayer)
	}
}