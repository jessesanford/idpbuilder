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

// Constants for build time estimation and caching
const (
	defaultBuildTime     = 30 * time.Second
	minCacheableSize     = 10
	maxParallelStages    = 10
	baseImagePullTime    = 5 * time.Second
	packageInstallTime   = 30 * time.Second
	regularCommandTime   = 10 * time.Second
	fileOperationTime    = 2 * time.Second
	defaultInstruction   = 1 * time.Second
	maxInstructions      = 20
	maxBuildTime         = 5 * time.Minute
)

// Analyzer examines Dockerfiles for multi-stage optimization opportunities.
type Analyzer struct {
	stageRegex    *regexp.Regexp
	copyFromRegex *regexp.Regexp
	argRegex      *regexp.Regexp
}

// wrapErr creates formatted errors with context
func wrapErr(err error, msg string) error {
	if err != nil {
		return fmt.Errorf("%s: %w", msg, err)
	}
	return nil
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
		return nil, wrapErr(err, "failed to parse stages")
	}
	if len(stages) == 0 {
		return nil, fmt.Errorf("no stages found in dockerfile")
	}

	dependencies := a.analyzeDependencies(stages)
	parallelGroups := a.identifyParallelGroups(stages, dependencies)

	return &api.StageAnalysis{
		Stages:            stages,
		Dependencies:      dependencies,
		ParallelGroups:    parallelGroups,
		CacheableStages:   a.identifyCacheableStages(stages),
		EstimatedTime:     a.estimateTotalBuildTime(stages, parallelGroups),
		CriticalPath:      a.calculateCriticalPath(stages, dependencies),
		OptimizationScore: a.calculateOptimizationScore(stages, parallelGroups),
		Warnings:          a.generateWarnings(stages, dependencies),
	}, nil
}

// parseStages extracts individual build stages from the Dockerfile content.
func (a *Analyzer) parseStages(dockerfile []byte) ([]*api.Stage, error) {
	var stages []*api.Stage
	scanner := bufio.NewScanner(bytes.NewReader(dockerfile))
	var currentStage *api.Stage
	stageIndex := 0

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if matches := a.stageRegex.FindStringSubmatch(line); matches != nil {
			if currentStage != nil {
				stages = append(stages, currentStage)
			}
			stageName := matches[2]
			if stageName == "" {
				stageName = fmt.Sprintf("stage%d", stageIndex)
			}
			currentStage = &api.Stage{
				Name: stageName, BaseImage: matches[1], Instructions: []string{}, Dependencies: []string{}, Cacheable: true, BuildArgs: make(map[string]string),
			}
			stageIndex++
		}

		if currentStage != nil {
			currentStage.Instructions = append(currentStage.Instructions, line)
			a.analyzeInstruction(currentStage, line)
		}
	}

	if currentStage != nil {
		stages = append(stages, currentStage)
	}
	if err := scanner.Err(); err != nil {
		return nil, wrapErr(err, "error reading dockerfile")
	}
	a.refineStagAnalysis(stages)
	return stages, nil
}

// analyzeInstruction examines individual Dockerfile instructions for optimization insights.
func (a *Analyzer) analyzeInstruction(stage *api.Stage, instruction string) {
	upper := strings.ToUpper(instruction)

	if matches := a.copyFromRegex.FindStringSubmatch(instruction); matches != nil {
		if _, err := strconv.Atoi(matches[1]); err != nil {
			stage.Dependencies = append(stage.Dependencies, matches[1])
		}
	}

	if matches := a.argRegex.FindStringSubmatch(instruction); matches != nil {
		argValue := ""
		if len(matches) > 2 {
			argValue = matches[2]
		}
		stage.BuildArgs[matches[1]] = argValue
	}

	if strings.HasPrefix(upper, "RUN ") && (strings.Contains(upper, "UPDATE") || strings.Contains(upper, "APKUPDATE")) {
		stage.Cacheable = false
	}

	stage.EstimatedBuildTime += a.estimateInstructionTime(instruction)
}

// estimateInstructionTime provides rough estimates for different instruction types.
func (a *Analyzer) estimateInstructionTime(instruction string) time.Duration {
	upper := strings.ToUpper(instruction)
	switch {
	case strings.HasPrefix(upper, "FROM"):
		return baseImagePullTime
	case strings.HasPrefix(upper, "RUN") && (strings.Contains(upper, "INSTALL") || strings.Contains(upper, "UPDATE")):
		return packageInstallTime
	case strings.HasPrefix(upper, "RUN"):
		return regularCommandTime
	case strings.HasPrefix(upper, "COPY") || strings.HasPrefix(upper, "ADD"):
		return fileOperationTime
	default:
		return defaultInstruction
	}
}

// refineStagAnalysis performs post-processing to refine stage analysis.
func (a *Analyzer) refineStagAnalysis(stages []*api.Stage) {
	stageMap := make(map[string]*api.Stage)
	for _, stage := range stages {
		stageMap[stage.Name] = stage
	}
	for _, stage := range stages {
		if baseStage, exists := stageMap[stage.BaseImage]; exists {
			stage.Dependencies = append(stage.Dependencies, baseStage.Name)
		}
		stage.Dependencies = a.removeDuplicates(stage.Dependencies)
	}
}

// analyzeDependencies creates a comprehensive dependency map between stages.
func (a *Analyzer) analyzeDependencies(stages []*api.Stage) map[string][]string {
	deps := make(map[string][]string)
	for _, stage := range stages {
		deps[stage.Name] = stage.Dependencies
	}
	return deps
}

// identifyParallelGroups finds groups of stages that can be built in parallel.
func (a *Analyzer) identifyParallelGroups(stages []*api.Stage, dependencies map[string][]string) [][]string {
	levels := a.calculateDependencyLevels(stages, dependencies)
	levelGroups := make(map[int][]string)
	for stageName, level := range levels {
		levelGroups[level] = append(levelGroups[level], stageName)
	}

	var groups [][]string
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

	var calcLevel func(string) int
	calcLevel = func(name string) int {
		if level, exists := levels[name]; exists {
			return level
		}
		if visited[name] {
			return 0
		}
		visited[name] = true
		maxLevel := -1
		for _, dep := range dependencies[name] {
			if depLevel := calcLevel(dep); depLevel > maxLevel {
				maxLevel = depLevel
			}
		}
		levels[name] = maxLevel + 1
		visited[name] = false
		return levels[name]
	}

	for _, stage := range stages {
		calcLevel(stage.Name)
	}
	return levels
}

// identifyCacheableStages determines which stages are good candidates for caching.
func (a *Analyzer) identifyCacheableStages(stages []*api.Stage) []string {
	var cacheable []string
	for _, stage := range stages {
		if stage.Cacheable && a.hasStableInstructions(stage) {
			cacheable = append(cacheable, stage.Name)
		}
	}
	return cacheable
}

// hasStableInstructions checks if a stage has stable, cacheable instructions.
func (a *Analyzer) hasStableInstructions(stage *api.Stage) bool {
	for _, inst := range stage.Instructions {
		if strings.Contains(strings.ToUpper(inst), "TIMESTAMP") ||
			strings.Contains(strings.ToUpper(inst), "DATE") ||
			strings.Contains(strings.ToUpper(inst), "RANDOM") {
			return false
		}
	}
	return true
}

// calculateCriticalPath identifies the longest dependency path through the stages.
func (a *Analyzer) calculateCriticalPath(stages []*api.Stage, dependencies map[string][]string) []string {
	levels := a.calculateDependencyLevels(stages, dependencies)
	var maxLevel int
	var criticalStage string
	for name, level := range levels {
		if level > maxLevel {
			maxLevel, criticalStage = level, name
		}
	}
	var path []string
	for current := criticalStage; current != ""; {
		path = append([]string{current}, path...)
		var next string
		var nextLevel = -1
		for _, dep := range dependencies[current] {
			if levels[dep] > nextLevel {
				nextLevel, next = levels[dep], dep
			}
		}
		current = next
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
	for _, group := range parallelGroups {
		var maxTime time.Duration
		for _, name := range group {
			if stage := stageMap[name]; stage != nil && stage.EstimatedBuildTime > maxTime {
				maxTime = stage.EstimatedBuildTime
			}
		}
		totalTime += maxTime
	}
	return totalTime
}

// calculateOptimizationScore assigns a score indicating the optimization potential.
func (a *Analyzer) calculateOptimizationScore(stages []*api.Stage, parallelGroups [][]string) int {
	if len(stages) <= 1 {
		return 0
	}
	score := int(float64(len(stages)) / float64(len(parallelGroups)) * 20)
	for _, group := range parallelGroups {
		if len(group) > 1 {
			score += len(group) * 5
		}
	}
	return score
}

// generateWarnings identifies potential issues and optimization opportunities.
func (a *Analyzer) generateWarnings(stages []*api.Stage, dependencies map[string][]string) []string {
	var warnings []string
	for _, stage := range stages {
		if len(stage.Instructions) > maxInstructions {
			warnings = append(warnings, fmt.Sprintf("Stage '%s' has %d instructions", stage.Name, len(stage.Instructions)))
		}
		if stage.EstimatedBuildTime > maxBuildTime {
			warnings = append(warnings, fmt.Sprintf("Stage '%s' has long build time", stage.Name))
		}
	}
	if len(stages) > 3 && len(a.identifyParallelGroups(stages, dependencies)) < 2 {
		warnings = append(warnings, "Limited parallelization opportunities")
	}
	return warnings
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