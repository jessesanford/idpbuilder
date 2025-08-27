package multistage

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
)

// DockerfileParser handles parsing of multi-stage Dockerfiles
type DockerfileParser struct {
	stagePattern   *regexp.Regexp
	copyFromPattern *regexp.Regexp
}

// NewDockerfileParser creates a new Dockerfile parser
func NewDockerfileParser() *DockerfileParser {
	return &DockerfileParser{
		stagePattern:   regexp.MustCompile(`(?i)^FROM\s+([^\s]+)(?:\s+AS\s+([^\s]+))?`),
		copyFromPattern: regexp.MustCompile(`(?i)^COPY\s+--from=([^\s]+)`),
	}
}

// Parse parses a multi-stage Dockerfile and returns a StageGraph
func (p *DockerfileParser) Parse(dockerfile io.Reader) (*StageGraph, error) {
	scanner := bufio.NewScanner(dockerfile)
	stages := []BuildStage{}
	currentStage := BuildStage{
		Commands:  []Command{},
		BuildArgs: make(map[string]string),
	}
	
	dependencies := make(map[string][]string)
	stageIndex := 0
	
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		
		// Check for FROM instruction (new stage)
		if matches := p.stagePattern.FindStringSubmatch(line); matches != nil {
			// Save previous stage if it exists
			if currentStage.BaseImage != "" {
				stages = append(stages, currentStage)
			}
			
			// Start new stage
			baseImage := matches[1]
			stageName := matches[2]
			
			if stageName == "" {
				stageName = fmt.Sprintf("stage-%d", stageIndex)
			}
			
			currentStage = BuildStage{
				Name:      stageName,
				BaseImage: baseImage,
				Commands:  []Command{},
				BuildArgs: make(map[string]string),
			}
			
			dependencies[stageName] = []string{}
			stageIndex++
			continue
		}
		
		// Parse other commands
		command, err := p.parseCommand(line)
		if err != nil {
			return nil, fmt.Errorf("error parsing command '%s': %w", line, err)
		}
		
		// Track COPY --from dependencies
		if command.Type == "COPY" && command.From != "" {
			if deps, exists := dependencies[currentStage.Name]; exists {
				// Add dependency if not already present
				found := false
				for _, dep := range deps {
					if dep == command.From {
						found = true
						break
					}
				}
				if !found {
					dependencies[currentStage.Name] = append(deps, command.From)
				}
			}
		}
		
		currentStage.Commands = append(currentStage.Commands, command)
	}
	
	// Add the last stage
	if currentStage.BaseImage != "" {
		stages = append(stages, currentStage)
	}
	
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading dockerfile: %w", err)
	}
	
	// Calculate execution order
	executionOrder, err := p.calculateExecutionOrder(stages, dependencies)
	if err != nil {
		return nil, err
	}
	
	return &StageGraph{
		Stages:         stages,
		Dependencies:   dependencies,
		ExecutionOrder: executionOrder,
	}, nil
}

// parseCommand parses a single Dockerfile command
func (p *DockerfileParser) parseCommand(line string) (Command, error) {
	parts := strings.Fields(line)
	if len(parts) == 0 {
		return Command{}, fmt.Errorf("empty command")
	}
	
	cmdType := strings.ToUpper(parts[0])
	args := parts[1:]
	
	command := Command{
		Type: cmdType,
		Args: args,
	}
	
	// Special handling for COPY --from
	if cmdType == "COPY" && len(args) > 0 && strings.HasPrefix(args[0], "--from=") {
		fromArg := args[0]
		command.From = strings.TrimPrefix(fromArg, "--from=")
		command.Args = args[1:] // Remove --from argument
	}
	
	return command, nil
}

// calculateExecutionOrder determines the order in which stages should be built
func (p *DockerfileParser) calculateExecutionOrder(stages []BuildStage, dependencies map[string][]string) ([]string, error) {
	stageMap := make(map[string]bool)
	for _, stage := range stages {
		stageMap[stage.Name] = true
	}
	
	// Topological sort using Kahn's algorithm
	inDegree := make(map[string]int)
	adjList := make(map[string][]string)
	
	// Initialize in-degree and adjacency list
	for _, stage := range stages {
		inDegree[stage.Name] = 0
		adjList[stage.Name] = []string{}
	}
	
	// Build adjacency list and calculate in-degrees
	for stageName, deps := range dependencies {
		for _, dep := range deps {
			if !stageMap[dep] {
				return nil, fmt.Errorf("stage '%s' depends on undefined stage '%s'", stageName, dep)
			}
			adjList[dep] = append(adjList[dep], stageName)
			inDegree[stageName]++
		}
	}
	
	// Find all nodes with no incoming edges
	queue := []string{}
	for stageName, degree := range inDegree {
		if degree == 0 {
			queue = append(queue, stageName)
		}
	}
	
	result := []string{}
	
	for len(queue) > 0 {
		// Remove a node from queue
		current := queue[0]
		queue = queue[1:]
		result = append(result, current)
		
		// For each neighbor of the current node
		for _, neighbor := range adjList[current] {
			inDegree[neighbor]--
			if inDegree[neighbor] == 0 {
				queue = append(queue, neighbor)
			}
		}
	}
	
	// Check for circular dependencies
	if len(result) != len(stages) {
		return nil, fmt.Errorf("circular dependency detected in stage graph")
	}
	
	return result, nil
}