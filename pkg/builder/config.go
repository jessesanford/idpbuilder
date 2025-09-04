package builder

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	v1 "github.com/google/go-containerregistry/pkg/v1"
)

// DockerfileInstruction represents a parsed Dockerfile instruction
type DockerfileInstruction struct {
	Command string
	Args    string
	Line    int
}

// ImageConfig represents basic container image configuration
type ImageConfig struct {
	Architecture string            `json:"architecture"`
	OS           string            `json:"os"`
	User         string            `json:"user,omitempty"`
	WorkingDir   string            `json:"workingDir,omitempty"`
	Env          []string          `json:"env,omitempty"`
	Cmd          []string          `json:"cmd,omitempty"`
	Entrypoint   []string          `json:"entrypoint,omitempty"`
	Labels       map[string]string `json:"labels,omitempty"`
	Created      time.Time         `json:"created"`
}

// ConfigManager handles image configuration management
type ConfigManager struct {
	config *ImageConfig
}

// NewConfigManager creates a new configuration manager
func NewConfigManager() *ConfigManager {
	return &ConfigManager{
		config: &ImageConfig{
			Architecture: "amd64",
			OS:           "linux",
			Labels:       make(map[string]string),
			Created:      time.Now(),
		},
	}
}

// SetUser sets the default user
func (cm *ConfigManager) SetUser(user string) {
	cm.config.User = user
}

// SetWorkingDir sets the working directory
func (cm *ConfigManager) SetWorkingDir(dir string) {
	cm.config.WorkingDir = dir
}

// AddEnv adds an environment variable
func (cm *ConfigManager) AddEnv(key, value string) {
	envVar := fmt.Sprintf("%s=%s", key, value)
	cm.config.Env = append(cm.config.Env, envVar)
}

// SetCmd sets the default command
func (cm *ConfigManager) SetCmd(cmd []string) {
	cm.config.Cmd = cmd
}

// SetEntrypoint sets the entrypoint command
func (cm *ConfigManager) SetEntrypoint(entrypoint []string) {
	cm.config.Entrypoint = entrypoint
}

// AddLabel adds a metadata label
func (cm *ConfigManager) AddLabel(key, value string) {
	cm.config.Labels[key] = value
}

// ToV1Config converts the configuration to a v1.Config
func (cm *ConfigManager) ToV1Config() (*v1.Config, error) {
	return &v1.Config{
		User:       cm.config.User,
		Env:        cm.config.Env,
		Cmd:        cm.config.Cmd,
		Entrypoint: cm.config.Entrypoint,
		WorkingDir: cm.config.WorkingDir,
		Labels:     cm.config.Labels,
	}, nil
}

// GetConfig returns the current image configuration
func (cm *ConfigManager) GetConfig() *ImageConfig {
	return cm.config
}

// parseDockerfileInstruction parses a single Dockerfile instruction line
func parseDockerfileInstruction(line string, lineNum int) (DockerfileInstruction, error) {
	parts := strings.SplitN(line, " ", 2)
	if len(parts) < 2 {
		return DockerfileInstruction{}, fmt.Errorf("invalid instruction format")
	}

	command := strings.ToUpper(strings.TrimSpace(parts[0]))
	args := strings.TrimSpace(parts[1])

	return DockerfileInstruction{
		Command: command,
		Args:    args,
		Line:    lineNum,
	}, nil
}

// ApplyInstructionToConfig applies a Dockerfile instruction to the configuration
func (cm *ConfigManager) ApplyInstructionToConfig(instruction DockerfileInstruction) error {
	switch instruction.Command {
	case "CMD":
		args, err := parseShellArgs(instruction.Args)
		if err != nil {
			return fmt.Errorf("failed to parse CMD: %w", err)
		}
		cm.SetCmd(args)
	case "ENTRYPOINT":
		args, err := parseShellArgs(instruction.Args)
		if err != nil {
			return fmt.Errorf("failed to parse ENTRYPOINT: %w", err)
		}
		cm.SetEntrypoint(args)
	case "WORKDIR":
		cm.SetWorkingDir(instruction.Args)
	case "USER":
		cm.SetUser(instruction.Args)
	case "ENV":
		parts := strings.SplitN(instruction.Args, " ", 2)
		if len(parts) == 2 {
			cm.AddEnv(parts[0], parts[1])
		}
	case "LABEL":
		parts := strings.SplitN(instruction.Args, " ", 2)
		if len(parts) == 2 {
			cm.AddLabel(parts[0], parts[1])
		}
	}
	return nil
}

// parseShellArgs parses shell-style arguments
func parseShellArgs(args string) ([]string, error) {
	args = strings.TrimSpace(args)
	if strings.HasPrefix(args, "[") && strings.HasSuffix(args, "]") {
		var result []string
		if err := json.Unmarshal([]byte(args), &result); err != nil {
			return nil, fmt.Errorf("invalid JSON array: %w", err)
		}
		return result, nil
	}
	return strings.Fields(args), nil
}