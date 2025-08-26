package optimizer

import (
	"bufio"
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/cnoe-io/idpbuilder/pkg/oci/api"
)

// Analyzer examines Dockerfiles for multi-stage optimization opportunities.
// It parses stage structure, identifies dependencies, and calculates parallelization groups.
type Analyzer struct {
	stageRegex    *regexp.Regexp
	copyFromRegex *regexp.Regexp
	argRegex      *regexp.Regexp
}

// NewAnalyzer creates a new Dockerfile analyzer with initialized patterns.
func NewAnalyzer() *Analyzer {
	return &Analyzer{
		stageRegex:    regexp.MustCompile(`(?i)^FROM\s+([^\s]+)(?:\s+AS\s+([^\s]+))?`),
		copyFromRegex: regexp.MustCompile(`(?i)COPY\s+--from=([^\s]+)`),
		argRegex:      regexp.MustCompile(`(?i)^ARG\s+([^=\s]+)(?:=([^\s]+))?`),
	}
}

// Analyze parses a Dockerfile and performs comprehensive stage analysis.
func (a *Analyzer) Analyze(dockerfile []byte) (*api.StageAnalysis, error) {
	if len(dockerfile) == 0 {
		return nil, fmt.Errorf("dockerfile content cannot be empty")
	}

	stages, err := a.parseStages(dockerfile)
	if err != nil {
		return nil, fmt.Errorf("failed to parse stages: %w", err)
	}

	if len(stages) == 0 {
		return nil, fmt.Errorf("no stages found in dockerfile")
	}

	// Analyze dependencies between stages
	dependencies := a.analyzeDependencies(stages)

	// Identify parallel execution groups
	parallelGroups := a.identifyParallelGroups(stages, dependencies)

	// Identify cacheable stages
	cacheableStages := a.identifyCacheableStages(stages)

	// Calculate critical path
	criticalPath := a.calculateCriticalPath(stages, dependencies)

	// Estimate total build time
	estimatedTime := a.estimateTotalBuildTime(stages, parallelGroups)

	// Calculate optimization score
	optimizationScore := a.calculateOptimizationScore(stages, parallelGroups)

	// Generate warnings for potential issues
	warnings := a.generateWarnings(stages, dependencies)

	return &api.StageAnalysis{
		Stages:            stages,
		Dependencies:      dependencies,
		ParallelGroups:    parallelGroups,
		CacheableStages:   cacheableStages,
		EstimatedTime:     estimatedTime,
		CriticalPath:      criticalPath,
		OptimizationScore: optimizationScore,
		Warnings:          warnings,
	}, nil
}

// parseStages extracts individual build stages from the Dockerfile content.
func (a *Analyzer) parseStages(dockerfile []byte) ([]*api.Stage, error) {
	var stages []*api.Stage
	scanner := bufio.NewScanner(bytes.NewReader(dockerfile))
	
	var currentStage *api.Stage
	stageIndex := 0
	lineNumber := 0

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		lineNumber++

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Check for FROM instruction (start of new stage)
		if matches := a.stageRegex.FindStringSubmatch(line); matches != nil {
			// Save the previous stage if it exists
			if currentStage != nil {
				stages = append(stages, currentStage)
			}

			// Start new stage
			baseImage := matches[1]
			stageName := matches[2]
			
			// If no stage name specified, generate one
			if stageName == "" {
				stageName = fmt.Sprintf("stage%d", stageIndex)
			}

			currentStage = &api.Stage{
				Name:         stageName,
				BaseImage:    baseImage,
				Instructions: []string{},
				Dependencies: []string{},
				Cacheable:    true, // Default to cacheable, will be refined
				BuildArgs:    make(map[string]string),
			}
			
			stageIndex++
		}

		// Add instruction to current stage
		if currentStage != nil {
			currentStage.Instructions = append(currentStage.Instructions, line)
			
			// Analyze specific instructions
			a.analyzeInstruction(currentStage, line)
		}
	}

	// Don't forget the last stage
	if currentStage != nil {
		stages = append(stages, currentStage)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading dockerfile: %w", err)
	}

	// Post-process stages to refine analysis
	a.refineStagAnalysis(stages)

	return stages, nil
}

// analyzeInstruction examines individual Dockerfile instructions for optimization insights.
func (a *Analyzer) analyzeInstruction(stage *api.Stage, instruction string) {
	upperInstruction := strings.ToUpper(instruction)

	// Check for COPY --from instructions (stage dependencies)
	if matches := a.copyFromRegex.FindStringSubmatch(instruction); matches != nil {
		fromStage := matches[1]
		// Check if it's a stage name (not a stage index)
		if _, err := strconv.Atoi(fromStage); err != nil {
			stage.Dependencies = append(stage.Dependencies, fromStage)
		}
	}

	// Check for ARG instructions
	if matches := a.argRegex.FindStringSubmatch(instruction); matches != nil {
		argName := matches[1]
		argValue := ""
		if len(matches) > 2 {
			argValue = matches[2]
		}
		stage.BuildArgs[argName] = argValue
	}

	// Instructions that make stages less cacheable
	if strings.HasPrefix(upperInstruction, "RUN ") && 
		(strings.Contains(upperInstruction, "APT-GET UPDATE") ||
		 strings.Contains(upperInstruction, "YUM UPDATE") ||
		 strings.Contains(upperInstruction, "APKUPDATE")) {
		stage.Cacheable = false
	}

	// Estimate build time based on instruction complexity
	estimatedTime := a.estimateInstructionTime(instruction)
	stage.EstimatedBuildTime += estimatedTime
}

// estimateInstructionTime provides rough estimates for different instruction types.
func (a *Analyzer) estimateInstructionTime(instruction string) time.Duration {
	upperInstruction := strings.ToUpper(instruction)
	
	switch {
	case strings.HasPrefix(upperInstruction, "FROM"):
		return 5 * time.Second // Base image pull/check
	case strings.HasPrefix(upperInstruction, "RUN"):
		// Estimate based on command complexity
		if strings.Contains(upperInstruction, "INSTALL") ||
		   strings.Contains(upperInstruction, "UPDATE") ||
		   strings.Contains(upperInstruction, "UPGRADE") {
			return 30 * time.Second // Package operations
		}
		return 10 * time.Second // Regular RUN commands
	case strings.HasPrefix(upperInstruction, "COPY") ||
		 strings.HasPrefix(upperInstruction, "ADD"):
		return 2 * time.Second // File operations
	default:
		return 1 * time.Second // Other instructions
	}
}

// refineStagAnalysis performs post-processing to refine stage analysis.
func (a *Analyzer) refineStagAnalysis(stages []*api.Stage) {
	stageMap := make(map[string]*api.Stage)
	for _, stage := range stages {
		stageMap[stage.Name] = stage
	}

	// Resolve dependencies and add implicit dependencies
	for _, stage := range stages {
		// Check base image dependencies
		if baseStage, exists := stageMap[stage.BaseImage]; exists {
			stage.Dependencies = append(stage.Dependencies, baseStage.Name)
		}

		// Remove duplicates from dependencies
		stage.Dependencies = a.removeDuplicates(stage.Dependencies)
	}
}

// analyzeDependencies creates a comprehensive dependency map between stages.
func (a *Analyzer) analyzeDependencies(stages []*api.Stage) map[string][]string {
	dependencies := make(map[string][]string)
	
	for _, stage := range stages {
		dependencies[stage.Name] = stage.Dependencies
	}

	return dependencies
}

// identifyParallelGroups finds groups of stages that can be built in parallel.
func (a *Analyzer) identifyParallelGroups(stages []*api.Stage, dependencies map[string][]string) [][]string {
	var groups [][]string
	processed := make(map[string]bool)
	
	// Create levels based on dependency depth
	levels := a.calculateDependencyLevels(stages, dependencies)
	
	// Group stages by level (stages at same level can run in parallel)
	levelGroups := make(map[int][]string)
	for stageName, level := range levels {
		levelGroups[level] = append(levelGroups[level], stageName)
	}

	// Convert level groups to parallel groups
	for level := 0; level < len(levelGroups); level++ {
		if group, exists := levelGroups[level]; exists && len(group) > 0 {
			groups = append(groups, group)
		}
	}

	return groups
}

// calculateDependencyLevels determines the dependency depth for each stage.
func (a *Analyzer) calculateDependencyLevels(stages []*api.Stage, dependencies map[string][]string) map[string]int {
	levels := make(map[string]int)
	visited := make(map[string]bool)

	var calculateLevel func(string) int
	calculateLevel = func(stageName string) int {
		if level, exists := levels[stageName]; exists {
			return level
		}

		if visited[stageName] {
			// Circular dependency detected
			return 0
		}

		visited[stageName] = true
		maxDepLevel := -1

		for _, dep := range dependencies[stageName] {
			depLevel := calculateLevel(dep)
			if depLevel > maxDepLevel {
				maxDepLevel = depLevel
			}
		}

		level := maxDepLevel + 1
		levels[stageName] = level
		visited[stageName] = false

		return level
	}

	// Calculate levels for all stages
	for _, stage := range stages {
		calculateLevel(stage.Name)
	}

	return levels
}

// identifyCacheableStages determines which stages are good candidates for caching.
func (a *Analyzer) identifyCacheableStages(stages []*api.Stage) []string {
	var cacheable []string
	
	for _, stage := range stages {
		if stage.Cacheable {
			// Additional heuristics for cacheability
			if a.hasStableInstructions(stage) {
				cacheable = append(cacheable, stage.Name)
			}
		}
	}

	return cacheable
}

// hasStableInstructions checks if a stage has stable, cacheable instructions.
func (a *Analyzer) hasStableInstructions(stage *api.Stage) bool {
	for _, instruction := range stage.Instructions {
		upperInstruction := strings.ToUpper(instruction)
		
		// Instructions that often change
		if strings.Contains(upperInstruction, "TIMESTAMP") ||
		   strings.Contains(upperInstruction, "DATE") ||
		   strings.Contains(upperInstruction, "RANDOM") {
			return false
		}
	}
	
	return true
}

// calculateCriticalPath identifies the longest dependency path through the stages.
func (a *Analyzer) calculateCriticalPath(stages []*api.Stage, dependencies map[string][]string) []string {
	levels := a.calculateDependencyLevels(stages, dependencies)
	
	// Find the stage with the highest level
	var maxLevel int
	var criticalStage string
	
	for stageName, level := range levels {
		if level > maxLevel {
			maxLevel = level
			criticalStage = stageName
		}
	}

	// Build critical path by following dependencies backwards
	var path []string
	current := criticalStage
	
	for current != "" {
		path = append([]string{current}, path...) // Prepend to maintain order
		
		// Find the dependency with the highest level
		var nextStage string
		var nextLevel = -1
		
		for _, dep := range dependencies[current] {
			if levels[dep] > nextLevel {
				nextLevel = levels[dep]
				nextStage = dep
			}
		}
		
		current = nextStage
	}

	return path
}

// estimateTotalBuildTime calculates the expected total build time considering parallelization.
func (a *Analyzer) estimateTotalBuildTime(stages []*api.Stage, parallelGroups [][]string) time.Duration {
	stageMap := make(map[string]*api.Stage)
	for _, stage := range stages {
		stageMap[stage.Name] = stage
	}

	var totalTime time.Duration

	// Sum the maximum time for each parallel group
	for _, group := range parallelGroups {
		var maxGroupTime time.Duration
		for _, stageName := range group {
			if stage, exists := stageMap[stageName]; exists {
				if stage.EstimatedBuildTime > maxGroupTime {
					maxGroupTime = stage.EstimatedBuildTime
				}
			}
		}
		totalTime += maxGroupTime
	}

	return totalTime
}

// calculateOptimizationScore assigns a score indicating the optimization potential.
func (a *Analyzer) calculateOptimizationScore(stages []*api.Stage, parallelGroups [][]string) int {
	if len(stages) <= 1 {
		return 0 // No optimization possible
	}

	// Base score calculation
	stageCount := len(stages)
	groupCount := len(parallelGroups)
	
	// Higher score for more parallelization opportunities
	parallelizationScore := int(float64(stageCount) / float64(groupCount) * 20)

	// Bonus for multiple parallel stages in groups
	parallelStageBonus := 0
	for _, group := range parallelGroups {
		if len(group) > 1 {
			parallelStageBonus += len(group) * 5
		}
	}

	return parallelizationScore + parallelStageBonus
}

// generateWarnings identifies potential issues and optimization opportunities.
func (a *Analyzer) generateWarnings(stages []*api.Stage, dependencies map[string][]string) []string {
	var warnings []string

	// Check for circular dependencies
	if a.hasCircularDependencies(dependencies) {
		warnings = append(warnings, "Circular dependencies detected between stages")
	}

	// Check for overly complex stages
	for _, stage := range stages {
		if len(stage.Instructions) > 20 {
			warnings = append(warnings, fmt.Sprintf("Stage '%s' has many instructions (%d), consider splitting", stage.Name, len(stage.Instructions)))
		}
		
		if stage.EstimatedBuildTime > 5*time.Minute {
			warnings = append(warnings, fmt.Sprintf("Stage '%s' has long estimated build time (%v)", stage.Name, stage.EstimatedBuildTime))
		}
	}

	// Check for missed parallelization opportunities
	if len(stages) > 3 && len(a.identifyParallelGroups(stages, dependencies)) < 2 {
		warnings = append(warnings, "Limited parallelization opportunities found, consider restructuring stages")
	}

	return warnings
}

// hasCircularDependencies detects circular dependencies in the stage graph.
func (a *Analyzer) hasCircularDependencies(dependencies map[string][]string) bool {
	visited := make(map[string]bool)
	recursionStack := make(map[string]bool)

	var hasCycle func(string) bool
	hasCycle = func(node string) bool {
		visited[node] = true
		recursionStack[node] = true

		for _, dep := range dependencies[node] {
			if !visited[dep] && hasCycle(dep) {
				return true
			} else if recursionStack[dep] {
				return true
			}
		}

		recursionStack[node] = false
		return false
	}

	for node := range dependencies {
		if !visited[node] && hasCycle(node) {
			return true
		}
	}

	return false
}

// removeDuplicates removes duplicate entries from a string slice.
func (a *Analyzer) removeDuplicates(slice []string) []string {
	seen := make(map[string]bool)
	var result []string

	for _, item := range slice {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}

	return result
}