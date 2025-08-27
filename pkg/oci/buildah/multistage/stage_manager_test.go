package multistage

import (
	"testing"
)

func TestNewStageManager(t *testing.T) {
	tests := []struct {
		name  string
		graph *StageGraph
	}{
		{
			name: "creates stage manager with simple graph",
			graph: &StageGraph{
				Stages: []BuildStage{
					{Name: "builder", BaseImage: "golang:1.19"},
					{Name: "runtime", BaseImage: "alpine:latest"},
				},
				Dependencies: map[string][]string{
					"runtime": {"builder"},
				},
				ExecutionOrder: []string{"builder", "runtime"},
			},
		},
		{
			name: "creates stage manager with empty graph",
			graph: &StageGraph{
				Stages:         []BuildStage{},
				Dependencies:   make(map[string][]string),
				ExecutionOrder: []string{},
			},
		},
		{
			name: "creates stage manager with complex dependency graph",
			graph: &StageGraph{
				Stages: []BuildStage{
					{Name: "base", BaseImage: "ubuntu:20.04"},
					{Name: "builder", BaseImage: "golang:1.19"},
					{Name: "tester", BaseImage: "golang:1.19"},
					{Name: "runtime", BaseImage: "alpine:latest"},
				},
				Dependencies: map[string][]string{
					"builder": {"base"},
					"tester":  {"base"},
					"runtime": {"builder", "tester"},
				},
				ExecutionOrder: []string{"base", "builder", "tester", "runtime"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sm := NewStageManager(tt.graph)
			if sm == nil {
				t.Error("NewStageManager() returned nil")
			}
			if sm.graph != tt.graph {
				t.Error("StageManager.graph not set correctly")
			}
			if sm.target != "" {
				t.Errorf("Expected empty target, got %q", sm.target)
			}
		})
	}
}

func TestStageManager_SetTarget(t *testing.T) {
	graph := &StageGraph{
		Stages: []BuildStage{
			{Name: "builder", BaseImage: "golang:1.19"},
			{Name: "tester", BaseImage: "golang:1.19"},
			{Name: "runtime", BaseImage: "alpine:latest"},
		},
		Dependencies: map[string][]string{
			"runtime": {"builder"},
			"tester":  {"builder"},
		},
		ExecutionOrder: []string{"builder", "tester", "runtime"},
	}

	tests := []struct {
		name        string
		target      string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid target stage",
			target:      "runtime",
			expectError: false,
		},
		{
			name:        "valid target stage - builder",
			target:      "builder",
			expectError: false,
		},
		{
			name:        "empty target should be allowed",
			target:      "",
			expectError: false,
		},
		{
			name:        "invalid target stage",
			target:      "nonexistent",
			expectError: true,
			errorMsg:    "target stage 'nonexistent' not found",
		},
		{
			name:        "case sensitive target",
			target:      "Runtime",
			expectError: true,
			errorMsg:    "target stage 'Runtime' not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create fresh stage manager for each test to avoid state pollution
			sm := NewStageManager(graph)
			err := sm.SetTarget(tt.target)
			
			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				} else if err.Error() != tt.errorMsg {
					t.Errorf("Expected error message %q, got %q", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if sm.target != tt.target {
					t.Errorf("Expected target %q, got %q", tt.target, sm.target)
				}
			}
		})
	}
}

func TestStageManager_GetExecutionStages(t *testing.T) {
	graph := &StageGraph{
		Stages: []BuildStage{
			{Name: "base", BaseImage: "ubuntu:20.04"},
			{Name: "builder", BaseImage: "golang:1.19"},
			{Name: "tester", BaseImage: "golang:1.19"},
			{Name: "runtime", BaseImage: "alpine:latest"},
		},
		Dependencies: map[string][]string{
			"builder": {"base"},
			"tester":  {"base"},
			"runtime": {"builder"},
		},
		ExecutionOrder: []string{"base", "builder", "tester", "runtime"},
	}

	tests := []struct {
		name     string
		target   string
		expected []string
	}{
		{
			name:     "no target returns all stages",
			target:   "",
			expected: []string{"base", "builder", "tester", "runtime"},
		},
		{
			name:     "target runtime returns base, builder, runtime",
			target:   "runtime",
			expected: []string{"base", "builder", "runtime"},
		},
		{
			name:     "target builder returns base, builder",
			target:   "builder",
			expected: []string{"base", "builder"},
		},
		{
			name:     "target tester returns base, tester",
			target:   "tester",
			expected: []string{"base", "tester"},
		},
		{
			name:     "target base returns only base",
			target:   "base",
			expected: []string{"base"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sm := NewStageManager(graph)
			if tt.target != "" {
				err := sm.SetTarget(tt.target)
				if err != nil {
					t.Fatalf("SetTarget failed: %v", err)
				}
			}
			
			result := sm.GetExecutionStages()
			
			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d stages, got %d: %v", len(tt.expected), len(result), result)
				return
			}
			
			for i, stage := range result {
				if stage != tt.expected[i] {
					t.Errorf("At index %d: expected %q, got %q", i, tt.expected[i], stage)
				}
			}
		})
	}
}

func TestStageManager_markNeededStages(t *testing.T) {
	// Test complex dependency graph with multiple levels
	graph := &StageGraph{
		Stages: []BuildStage{
			{Name: "base", BaseImage: "ubuntu:20.04"},
			{Name: "tools", BaseImage: "ubuntu:20.04"},
			{Name: "builder", BaseImage: "golang:1.19"},
			{Name: "tester", BaseImage: "golang:1.19"},
			{Name: "runtime", BaseImage: "alpine:latest"},
		},
		Dependencies: map[string][]string{
			"tools":   {"base"},
			"builder": {"base", "tools"},
			"tester":  {"base", "tools"},
			"runtime": {"builder"},
		},
		ExecutionOrder: []string{"base", "tools", "builder", "tester", "runtime"},
	}

	tests := []struct {
		name     string
		target   string
		expected map[string]bool
	}{
		{
			name:   "runtime stage marks all its dependencies",
			target: "runtime",
			expected: map[string]bool{
				"base":    true,
				"tools":   true,
				"builder": true,
				"runtime": true,
			},
		},
		{
			name:   "builder stage marks its dependencies",
			target: "builder",
			expected: map[string]bool{
				"base":    true,
				"tools":   true,
				"builder": true,
			},
		},
		{
			name:   "base stage only marks itself",
			target: "base",
			expected: map[string]bool{
				"base": true,
			},
		},
		{
			name:   "tools stage marks base and itself",
			target: "tools",
			expected: map[string]bool{
				"base":  true,
				"tools": true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sm := NewStageManager(graph)
			needed := make(map[string]bool)
			
			sm.markNeededStages(tt.target, needed)
			
			// Check that all expected stages are marked
			for stage, shouldBe := range tt.expected {
				if needed[stage] != shouldBe {
					t.Errorf("Stage %q: expected %v, got %v", stage, shouldBe, needed[stage])
				}
			}
			
			// Check that no unexpected stages are marked
			for stage := range needed {
				if !tt.expected[stage] {
					t.Errorf("Unexpected stage %q was marked as needed", stage)
				}
			}
		})
	}
}

// TestStageManager_GetExecutionStages_EmptyGraph tests behavior with empty graph
func TestStageManager_GetExecutionStages_EmptyGraph(t *testing.T) {
	graph := &StageGraph{
		Stages:         []BuildStage{},
		Dependencies:   make(map[string][]string),
		ExecutionOrder: []string{},
	}
	
	sm := NewStageManager(graph)
	result := sm.GetExecutionStages()
	
	if len(result) != 0 {
		t.Errorf("Expected empty result for empty graph, got %v", result)
	}
}

// TestStageManager_CircularDependencyHandling tests behavior when graph has circular dependencies
func TestStageManager_CircularDependencyHandling(t *testing.T) {
	// Note: This test verifies that markNeededStages doesn't infinite loop
	// The circular dependency should have been caught by the parser
	graph := &StageGraph{
		Stages: []BuildStage{
			{Name: "stage1", BaseImage: "ubuntu:20.04"},
			{Name: "stage2", BaseImage: "alpine:latest"},
		},
		Dependencies: map[string][]string{
			"stage1": {"stage2"},
			"stage2": {"stage1"}, // Circular dependency
		},
		ExecutionOrder: []string{"stage1", "stage2"}, // Invalid order, but test shouldn't hang
	}
	
	sm := NewStageManager(graph)
	err := sm.SetTarget("stage1")
	if err != nil {
		t.Fatalf("SetTarget failed: %v", err)
	}
	
	// This should not hang due to the circular dependency guard
	result := sm.GetExecutionStages()
	
	// The exact result depends on implementation, but it shouldn't hang
	if len(result) > 2 {
		t.Errorf("Expected at most 2 stages, got %d: %v", len(result), result)
	}
}