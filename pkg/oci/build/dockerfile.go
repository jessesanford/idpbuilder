// Package build provides core Dockerfile parsing for OCI image builds.
package build

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// DockerfileParser handles parsing and validation of Dockerfile content.
type DockerfileParser struct {
	dockerfile string
	context    string
}

// ParsedDockerfile represents the parsed Dockerfile structure.
type ParsedDockerfile struct {
	Stages     []*BuildStage     `json:"stages"`
	BaseImages []string          `json:"base_images"`
	BuildArgs  map[string]string `json:"build_args"`
	Labels     map[string]string `json:"labels"`
}

// BuildStage represents a build stage in a multi-stage Dockerfile.
type BuildStage struct {
	Name         string         `json:"name"`
	BaseImage    string         `json:"base_image"`
	Instructions []*Instruction `json:"instructions"`
	Index        int            `json:"index"`
}

// Instruction represents a Dockerfile instruction.
type Instruction struct {
	Command string   `json:"command"`
	Args    []string `json:"args"`
	Line    int      `json:"line"`
}

// Valid Dockerfile commands
var validCommands = map[string]bool{
	"FROM": true, "RUN": true, "CMD": true, "LABEL": true, "EXPOSE": true,
	"ENV": true, "ADD": true, "COPY": true, "ENTRYPOINT": true, "VOLUME": true,
	"USER": true, "WORKDIR": true, "ARG": true, "HEALTHCHECK": true,
}

// NewDockerfileParser creates a new parser.
func NewDockerfileParser(dockerfile, context string) *DockerfileParser {
	return &DockerfileParser{dockerfile: dockerfile, context: context}
}

// Parse parses the Dockerfile and returns structured representation.
func (p *DockerfileParser) Parse() (*ParsedDockerfile, error) {
	lines, err := p.loadFile()
	if err != nil {
		return nil, err
	}

	if err := p.validateSyntax(lines); err != nil {
		return nil, err
	}

	stages, err := p.extractStages(lines)
	if err != nil {
		return nil, err
	}

	parsed := &ParsedDockerfile{
		Stages:     stages,
		BaseImages: make([]string, 0),
		BuildArgs:  make(map[string]string),
		Labels:     make(map[string]string),
	}

	// Extract global info
	p.extractGlobalInfo(parsed)
	return parsed, nil
}

// ValidateSyntax validates basic Dockerfile syntax.
func (p *DockerfileParser) ValidateSyntax() error {
	lines, err := p.loadFile()
	if err != nil {
		return err
	}
	return p.validateSyntax(lines)
}

// ExtractStages extracts build stages from Dockerfile.
func (p *DockerfileParser) ExtractStages() ([]*BuildStage, error) {
	lines, err := p.loadFile()
	if err != nil {
		return nil, err
	}
	return p.extractStages(lines)
}

// loadFile loads the Dockerfile content.
func (p *DockerfileParser) loadFile() ([]string, error) {
	dockerfilePath := p.dockerfile
	if !filepath.IsAbs(dockerfilePath) {
		dockerfilePath = filepath.Join(p.context, p.dockerfile)
	}

	file, err := os.Open(dockerfilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open Dockerfile: %w", err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		// Handle line continuations
		for strings.HasSuffix(line, "\\") && scanner.Scan() {
			line = strings.TrimSuffix(line, "\\") + " " + strings.TrimLeft(scanner.Text(), " \t")
		}
		lines = append(lines, line)
	}

	return lines, scanner.Err()
}

// validateSyntax performs basic syntax validation.
func (p *DockerfileParser) validateSyntax(lines []string) error {
	if len(lines) == 0 {
		return fmt.Errorf("dockerfile is empty")
	}

	hasFrom := false
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		parts := strings.Fields(trimmed)
		if len(parts) == 0 {
			continue
		}

		command := strings.ToUpper(parts[0])
		if !validCommands[command] {
			return fmt.Errorf("invalid command '%s' at line %d", command, i+1)
		}

		if command == "FROM" {
			hasFrom = true
			if len(parts) < 2 {
				return fmt.Errorf("FROM requires image at line %d", i+1)
			}
		} else if !hasFrom && command != "ARG" {
			return fmt.Errorf("FROM must be first instruction at line %d", i+1)
		}
	}

	if !hasFrom {
		return fmt.Errorf("dockerfile must contain FROM instruction")
	}
	return nil
}

// extractStages extracts build stages from lines.
func (p *DockerfileParser) extractStages(lines []string) ([]*BuildStage, error) {
	var stages []*BuildStage
	var currentStage *BuildStage
	stageIndex := 0

	for lineNum, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		parts := strings.Fields(trimmed)
		if len(parts) == 0 {
			continue
		}

		command := strings.ToUpper(parts[0])
		args := parts[1:]

		if command == "FROM" {
			// Start new stage
			if currentStage != nil {
				stages = append(stages, currentStage)
				stageIndex++
			}

			baseImage, stageName := p.parseFROM(strings.Join(args, " "))
			currentStage = &BuildStage{
				Name:         stageName,
				BaseImage:    baseImage,
				Instructions: make([]*Instruction, 0),
				Index:        stageIndex,
			}
		}

		// Add instruction to current stage
		if currentStage != nil {
			instruction := &Instruction{
				Command: command,
				Args:    args,
				Line:    lineNum + 1,
			}
			currentStage.Instructions = append(currentStage.Instructions, instruction)
		}
	}

	if currentStage != nil {
		stages = append(stages, currentStage)
	}

	return stages, nil
}

// parseFROM parses FROM instruction.
func (p *DockerfileParser) parseFROM(args string) (baseImage, stageName string) {
	parts := strings.Fields(args)
	if len(parts) == 0 {
		return "", ""
	}

	baseImage = parts[0]
	if strings.HasPrefix(baseImage, "--platform=") && len(parts) > 1 {
		baseImage = parts[1]
		parts = parts[1:]
	}

	if len(parts) >= 3 && strings.ToUpper(parts[1]) == "AS" {
		stageName = parts[2]
	} else {
		stageName = fmt.Sprintf("stage-%d", time.Now().UnixNano()%1000)
	}

	return baseImage, stageName
}

func (p *DockerfileParser) extractGlobalInfo(parsed *ParsedDockerfile) {
	seenImages := make(map[string]bool)
	for _, stage := range parsed.Stages {
		if !seenImages[stage.BaseImage] {
			parsed.BaseImages = append(parsed.BaseImages, stage.BaseImage)
			seenImages[stage.BaseImage] = true
		}
		for _, instruction := range stage.Instructions {
			switch instruction.Command {
			case "ARG":
				if len(instruction.Args) > 0 {
					arg := instruction.Args[0]
					if parts := strings.SplitN(arg, "=", 2); len(parts) == 2 {
						parsed.BuildArgs[parts[0]] = parts[1]
					} else {
						parsed.BuildArgs[parts[0]] = ""
					}
				}
			case "LABEL":
				for _, arg := range instruction.Args {
					if parts := strings.SplitN(arg, "=", 2); len(parts) == 2 {
						parsed.Labels[parts[0]] = parts[1]
					}
				}
			}
		}
	}
}

func (p *ParsedDockerfile) GetStageByName(name string) *BuildStage {
	for _, stage := range p.Stages {
		if stage.Name == name {
			return stage
		}
	}
	return nil
}

func (p *ParsedDockerfile) GetFinalStage() *BuildStage {
	if len(p.Stages) == 0 {
		return nil
	}
	return p.Stages[len(p.Stages)-1]
}

func (p *ParsedDockerfile) HasMultiStage() bool {
	return len(p.Stages) > 1
}
