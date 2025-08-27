package multistage

import (
	"fmt"
)

// StageManager coordinates the execution of multi-stage builds
type StageManager struct {
	graph  *StageGraph
	target string
}

// NewStageManager creates a new stage manager
func NewStageManager(graph *StageGraph) *StageManager {
	return &StageManager{graph: graph}
}

// SetTarget sets the target stage for builds
func (sm *StageManager) SetTarget(target string) error {
	if target == "" {
		return nil
	}
	
	for _, stage := range sm.graph.Stages {
		if stage.Name == target {
			sm.target = target
			return nil
		}
	}
	return fmt.Errorf("target stage '%s' not found", target)
}

// GetExecutionStages returns stages that need to be built for the target
func (sm *StageManager) GetExecutionStages() []string {
	if sm.target == "" {
		return sm.graph.ExecutionOrder
	}
	
	// Find stages needed for target (including dependencies)
	needed := make(map[string]bool)
	sm.markNeededStages(sm.target, needed)
	
	result := []string{}
	for _, stage := range sm.graph.ExecutionOrder {
		if needed[stage] {
			result = append(result, stage)
		}
	}
	return result
}

// markNeededStages recursively marks stages needed for target
func (sm *StageManager) markNeededStages(stage string, needed map[string]bool) {
	if needed[stage] {
		return
	}
	needed[stage] = true
	
	for _, dep := range sm.graph.Dependencies[stage] {
		sm.markNeededStages(dep, needed)
	}
}