package optimizer

import (
	"fmt"

	"github.com/jessesanford/idpbuilder/pkg/oci/api"
)

// GraphBuilder builds dependency graphs from build stages
type GraphBuilder struct {
	nodes map[string]*Node
	edges map[string][]string
}

// Node represents a stage in the dependency graph
type Node struct {
	Stage    *api.Stage
	Level    int
	Visited  bool
	Children []string
}

// NewGraphBuilder creates a new graph builder
func NewGraphBuilder() *GraphBuilder {
	return &GraphBuilder{
		nodes: make(map[string]*Node),
		edges: make(map[string][]string),
	}
}

// Build constructs a dependency graph from stages
func (g *GraphBuilder) Build(stages []*api.Stage) (*api.DependencyGraph, error) {
	if len(stages) == 0 {
		return nil, fmt.Errorf("no stages provided")
	}

	// Reset builder state
	g.nodes = make(map[string]*Node)
	g.edges = make(map[string][]string)

	// Build nodes
	for _, stage := range stages {
		if stage.Name == "" {
			return nil, fmt.Errorf("stage name cannot be empty")
		}
		g.nodes[stage.Name] = &Node{
			Stage:    stage,
			Children: make([]string, 0),
		}
		g.edges[stage.Name] = make([]string, 0)
	}

	// Build edges from dependencies
	for _, stage := range stages {
		for _, dep := range stage.Dependencies {
			if _, exists := g.nodes[dep]; !exists {
				return nil, fmt.Errorf("dependency %s not found for stage %s", dep, stage.Name)
			}
			g.edges[dep] = append(g.edges[dep], stage.Name)
			g.nodes[dep].Children = append(g.nodes[dep].Children, stage.Name)
		}
	}

	// Perform topological sort to validate and calculate levels
	sorted, err := g.topologicalSort()
	if err != nil {
		return nil, err
	}

	// Calculate levels for parallel execution
	levels := g.calculateLevels()

	// Convert to API format
	apiNodes := make(map[string]*api.GraphNode)
	for name, node := range g.nodes {
		apiNodes[name] = &api.GraphNode{
			Stage:        node.Stage,
			Dependencies: node.Stage.Dependencies,
			Level:        node.Level,
		}
	}

	graph := &api.DependencyGraph{
		Nodes: apiNodes,
		Edges: g.edges,
	}

	// Store additional metadata in the first node for debugging
	if len(sorted) > 0 && len(levels) > 0 {
		// Graph is valid, levels calculated
	}

	return graph, nil
}

// topologicalSort performs Kahn's algorithm for topological sorting
func (g *GraphBuilder) topologicalSort() ([]string, error) {
	// Calculate in-degrees
	inDegree := make(map[string]int)
	for name := range g.nodes {
		inDegree[name] = len(g.nodes[name].Stage.Dependencies)
	}

	// Initialize queue with nodes having no dependencies
	queue := make([]string, 0)
	for name, degree := range inDegree {
		if degree == 0 {
			queue = append(queue, name)
		}
	}

	result := make([]string, 0, len(g.nodes))

	// Process queue
	for len(queue) > 0 {
		// Dequeue
		current := queue[0]
		queue = queue[1:]
		result = append(result, current)

		// Process children
		for _, child := range g.nodes[current].Children {
			inDegree[child]--
			if inDegree[child] == 0 {
				queue = append(queue, child)
			}
		}
	}

	// Check for cycles
	if len(result) != len(g.nodes) {
		return nil, fmt.Errorf("circular dependency detected in stage graph")
	}

	return result, nil
}

// calculateLevels groups nodes by dependency level for parallel execution
func (g *GraphBuilder) calculateLevels() [][]string {
	levels := make([][]string, 0)
	processed := make(map[string]bool)

	for len(processed) < len(g.nodes) {
		currentLevel := make([]string, 0)

		// Find nodes with all dependencies processed
		for name, node := range g.nodes {
			if processed[name] {
				continue
			}

			allDepsProcessed := true
			for _, dep := range node.Stage.Dependencies {
				if !processed[dep] {
					allDepsProcessed = false
					break
				}
			}

			if allDepsProcessed {
				currentLevel = append(currentLevel, name)
				node.Level = len(levels)
			}
		}

		if len(currentLevel) == 0 {
			break // Safety check
		}

		levels = append(levels, currentLevel)

		// Mark current level as processed
		for _, name := range currentLevel {
			processed[name] = true
		}
	}

	return levels
}