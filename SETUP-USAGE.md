# 🚀 Software Factory 2.0 Setup Usage Guide

## Quick Start

### Interactive Setup (Default)
```bash
./setup.sh
```
This will prompt you for all configuration options interactively.

### Non-Interactive Setup (With Config File)
```bash
./setup.sh --config my-config.yaml
```
This uses a configuration file to set up the project without any prompts.

## Setup Methods

### Method 1: Using setup.sh (Recommended)

The main `setup.sh` script handles both interactive and non-interactive modes:

```bash
# Interactive mode (prompts for all options)
./setup.sh

# Non-interactive mode (uses config file)
./setup.sh --config setup-config.yaml

# Show help
./setup.sh --help
```

**How it works:**
- If no arguments: Runs interactive setup
- If `--config` flag: Automatically redirects to `setup-noninteractive.sh`
- Handles both modes seamlessly

### Method 2: Direct Non-Interactive Script

You can also call the non-interactive script directly:

```bash
./setup-noninteractive.sh --config setup-config.yaml
```

This is what `setup.sh` calls internally when you use `--config`.

## Configuration File

For non-interactive setup, create a YAML configuration file:

### 1. Copy the Example
```bash
cp setup-config-example.yaml my-project-config.yaml
```

### 2. Edit the Configuration
```yaml
project:
  name: "my-awesome-project"
  description: "My project description"
  target_dir: "/workspaces/my-project"
  github_url: "https://github.com/myorg/myrepo"

technology:
  primary_language: "Go"
  frameworks:
    - "Kubernetes/KCP"
    - "gRPC"

agents:
  selected:
    - "Software Engineer"
    - "Code Reviewer"
    - "Architect"
  expertise:
    - "Cloud Architecture"
    - "API Design"

implementation:
  plan_type: "generate"
  project_type: "Kubernetes Controller/Operator"
  estimated_loc: 5000
  num_phases: 3
  test_coverage: 80

constraints:
  max_lines_per_effort: 800
  max_parallel_agents: 3
  code_review: "mandatory"
  security_level: 1

directory_handling:
  if_exists: "backup"
  create_parent: "create"

options:
  verbose: true
  create_line_counter: true
```

### 3. Run Setup
```bash
./setup.sh --config my-project-config.yaml
```

## Command Line Options

### setup.sh Options
```
./setup.sh [options]

Options:
  (no options)         Interactive setup mode
  --config, -c <file>  Non-interactive setup with config file
  --help, -h          Show help message

Examples:
  ./setup.sh                              # Interactive
  ./setup.sh --config project.yaml        # Non-interactive
  ./setup.sh -c project.yaml              # Short form
```

### setup-noninteractive.sh Options
```
./setup-noninteractive.sh [options]

Options:
  --config, -c <file>  Configuration file (required)
  --verbose, -v        Show detailed output
  --help, -h          Show help message

Examples:
  ./setup-noninteractive.sh --config project.yaml
  ./setup-noninteractive.sh -c project.yaml -v
```

## Interactive vs Non-Interactive

### When to Use Interactive Mode
- First-time setup
- Exploring available options
- Small, one-off projects
- Learning the system

### When to Use Non-Interactive Mode
- CI/CD pipelines
- Automated deployments
- Multiple similar projects
- Reproducible setups
- Docker/container builds

## Troubleshooting

### "Config file not found"
```bash
# Check file exists
ls -la setup-config.yaml

# Check you're in the right directory
pwd

# Use absolute path if needed
./setup.sh --config /full/path/to/config.yaml
```

### "setup-noninteractive.sh not found"
```bash
# Ensure both scripts are present
ls -la setup*.sh

# Make them executable
chmod +x setup.sh setup-noninteractive.sh
```

### "Permission denied"
```bash
# Make scripts executable
chmod +x setup.sh setup-noninteractive.sh

# Check target directory permissions
ls -la /workspaces/
```

### Interactive Mode Starts Instead of Non-Interactive
```bash
# Make sure you're using --config flag
./setup.sh --config my-config.yaml  # Correct
./setup.sh my-config.yaml           # Wrong - will start interactive
```

## Examples

### Minimal Setup
```bash
# Create minimal config
cat > minimal.yaml << 'EOF'
project:
  name: "test-project"
  target_dir: "/tmp/test-project"
technology:
  primary_language: "Python"
  frameworks:
    - "FastAPI"
agents:
  selected:
    - "Software Engineer"
implementation:
  plan_type: "generate"
  project_type: "API Service"
  estimated_loc: 1000
  num_phases: 2
  test_coverage: 75
constraints:
  max_lines_per_effort: 500
EOF

# Run setup
./setup.sh --config minimal.yaml
```

### Full Kubernetes Project
```bash
# Use the example config
./setup.sh --config setup-config-example.yaml
```

### Batch Setup
```bash
#!/bin/bash
# Setup multiple projects

for config in configs/*.yaml; do
    echo "Setting up from $config..."
    ./setup.sh --config "$config"
done
```

## Best Practices

1. **Test Config First**: Run with a test directory first
2. **Version Control Configs**: Keep config files in git
3. **Use Absolute Paths**: Avoid confusion with relative paths
4. **Backup Existing**: Set `if_exists: "backup"` for safety
5. **Verbose Mode**: Use `-v` for debugging

## Summary

- **One Script, Two Modes**: `setup.sh` handles both interactive and non-interactive
- **Config Flag**: Use `--config` for non-interactive mode
- **Automatic Redirect**: `setup.sh --config` calls `setup-noninteractive.sh`
- **Help Available**: Both scripts have `--help` flags
- **Example Provided**: `setup-config-example.yaml` shows all options