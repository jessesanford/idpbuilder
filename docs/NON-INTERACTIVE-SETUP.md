# 🚀 Non-Interactive Setup Guide

## Overview

Software Factory 2.0 now supports non-interactive setup using a configuration file. This allows you to:
- Automate project creation in CI/CD pipelines
- Quickly spin up multiple projects with similar configurations
- Version control your project setup configuration
- Avoid repetitive manual input

## Quick Start

1. **Copy the example configuration:**
   ```bash
   cp setup-config-example.yaml my-project-config.yaml
   ```

2. **Edit the configuration file:**
   ```bash
   nano my-project-config.yaml
   ```

3. **Run setup with config:**
   ```bash
   ./setup-noninteractive.sh --config my-project-config.yaml
   ```

## Configuration File Format

The configuration file uses YAML format with the following structure:

### Project Information
```yaml
project:
  name: "my-awesome-project"           # Required: Project name
  description: "A cool project"        # Optional: Project description
  target_dir: "/workspaces/my-project" # Required: Installation directory
  github_url: "https://github.com/..."  # Optional: Git remote URL
```

### Technology Stack
```yaml
technology:
  primary_language: "Go"                # Required: Main language
  frameworks:                          # List of frameworks/libraries
    - "Kubernetes/KCP"
    - "gRPC"
    - "Prometheus"
```

Available languages and their frameworks:
- **Go**: Kubernetes/KCP, Gin Web Framework, GORM, Cobra CLI, gRPC, Prometheus
- **Python**: Django, FastAPI, Flask, SQLAlchemy, Celery, PyTorch, NumPy/Pandas  
- **TypeScript**: React, Next.js, Node.js/Express, NestJS, Vue.js, Angular
- **Java**, **Rust**, **C++**: Use `frameworks_custom` for custom list

### Agent Configuration
```yaml
agents:
  selected:                            # Agents to include
    - "Software Engineer"
    - "Code Reviewer"
    - "Architect"
  expertise:                           # Areas of expertise
    - "Cloud Architecture"
    - "API Design"
    - "Testing Strategies"
```

### Implementation Planning
```yaml
implementation:
  plan_type: "generate"                # Options: generate, existing, idpbuilder-example
  
  # If plan_type is "generate":
  project_type: "Kubernetes Controller/Operator"
  estimated_loc: 5000
  num_phases: 3
  test_coverage: 80
  
  # If plan_type is "existing":
  # existing_plan_path: "/path/to/plan.md"
  
  # If plan_type is "idpbuilder-example":
  # No additional config needed
```

### Development Constraints
```yaml
constraints:
  max_lines_per_effort: 800           # Max lines per effort
  max_parallel_agents: 3              # Max concurrent agents
  code_review: "mandatory"            # Review requirement
  security_level: 1                   # 0=Standard, 1=Enhanced, 2=Maximum
```

### Directory Handling
```yaml
directory_handling:
  if_exists: "backup"                 # Options: backup, delete, ask
  create_parent: "create"             # Options: create, fail, ask
```

### Additional Options
```yaml
options:
  skip_git_init: false                # Skip git initialization
  skip_remote: false                  # Skip adding git remote
  verbose: true                       # Show detailed output
  create_line_counter: true           # Create line counter for Go/K8s
```

## Example Configurations

### Minimal Configuration
```yaml
project:
  name: "minimal-project"
  target_dir: "/workspaces/minimal"

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
  estimated_loc: 2000
  num_phases: 2
  test_coverage: 75

constraints:
  max_lines_per_effort: 600
```

### Kubernetes Controller Project
```yaml
project:
  name: "kcp-controller"
  description: "KCP multi-tenant controller"
  target_dir: "/workspaces/kcp-controller"
  github_url: "https://github.com/myorg/kcp-controller"

technology:
  primary_language: "Go"
  frameworks:
    - "Kubernetes/KCP"
    - "gRPC"
    - "Prometheus"

agents:
  selected:
    - "Software Engineer"
    - "Code Reviewer"
    - "Architect"
  expertise:
    - "Cloud Architecture"
    - "API Design"
    - "Testing Strategies"
    - "Security"

implementation:
  plan_type: "generate"
  project_type: "Kubernetes Controller/Operator"
  estimated_loc: 8000
  num_phases: 4
  test_coverage: 85

constraints:
  max_lines_per_effort: 800
  max_parallel_agents: 3
  code_review: "mandatory"
  security_level: 2

directory_handling:
  if_exists: "backup"
  create_parent: "create"

options:
  verbose: true
  create_line_counter: true
```

### Using IDPBuilder Example
```yaml
project:
  name: "idpbuilder-oci"
  description: "Add OCI build capabilities to IDPBuilder"
  target_dir: "/workspaces/idpbuilder-oci"

technology:
  primary_language: "Go"
  frameworks:
    - "Kubernetes/KCP"

agents:
  selected:
    - "Software Engineer"
    - "Code Reviewer"
    - "Architect"

implementation:
  plan_type: "idpbuilder-example"  # Uses built-in example

constraints:
  max_lines_per_effort: 800
  max_parallel_agents: 3
```

### Using Existing Plan
```yaml
project:
  name: "existing-plan-project"
  target_dir: "/workspaces/my-project"

technology:
  primary_language: "TypeScript"
  frameworks:
    - "React"
    - "Next.js"

agents:
  selected:
    - "Software Engineer"
    - "Code Reviewer"

implementation:
  plan_type: "existing"
  existing_plan_path: "/home/user/my-plan.md"

constraints:
  max_lines_per_effort: 700
```

## Command Line Options

```bash
./setup-noninteractive.sh [options]

Options:
  --config, -c <file>  Configuration file for non-interactive setup
  --verbose, -v        Show detailed output
  --help, -h          Show help message

Examples:
  ./setup-noninteractive.sh --config setup.yaml
  ./setup-noninteractive.sh -c setup.yaml -v
  ./setup-noninteractive.sh --help
```

## Automation Examples

### CI/CD Pipeline
```yaml
# GitHub Actions example
- name: Setup Software Factory 2.0
  run: |
    git clone https://github.com/yourorg/software-factory-2.0-template
    cd software-factory-2.0-template
    ./setup-noninteractive.sh --config ../project-config.yaml
```

### Batch Project Creation
```bash
#!/bin/bash
# Create multiple projects from configs

for config in configs/*.yaml; do
    echo "Setting up project from $config..."
    ./setup-noninteractive.sh --config "$config"
done
```

### Docker Container Setup
```dockerfile
FROM ubuntu:latest
COPY setup-config.yaml /tmp/
COPY software-factory-2.0-template /opt/sf2/
WORKDIR /opt/sf2
RUN ./setup-noninteractive.sh --config /tmp/setup-config.yaml
```

## Validation Rules

The setup script validates:

1. **Required Fields**:
   - `project.name` must be provided
   - `technology.primary_language` must be specified
   - If `plan_type` is "generate", `project_type` is required

2. **Directory Handling**:
   - Target directory path is converted to absolute if relative
   - Parent directory permissions are checked
   - Existing directories handled per configuration

3. **Agent Configuration**:
   - Orchestrator is always included (automatic)
   - Invalid agent names are ignored

4. **Constraints**:
   - Numeric values must be positive integers
   - Security level must be 0, 1, or 2

## Troubleshooting

### Common Issues

1. **"Configuration file not found"**
   - Check file path is correct
   - Ensure file has .yaml or .yml extension

2. **"Required field missing"**
   - Verify `project.name` and `technology.primary_language` are set
   - Check YAML syntax is valid

3. **"Permission denied"**
   - Ensure write permissions for target directory
   - Check parent directory permissions

4. **"Invalid plan_type"**
   - Must be: generate, existing, or idpbuilder-example
   - Check spelling and case

### Debug Mode

Run with verbose flag for detailed output:
```bash
./setup-noninteractive.sh --config my-config.yaml --verbose
```

### Validation Script

Test your configuration without running setup:
```bash
# Dry run (not implemented yet, but planned)
./setup-noninteractive.sh --config my-config.yaml --dry-run
```

## Migration from Interactive

To create a config from an existing interactive session:

1. Run interactive setup once
2. Note your answers
3. Create config file with those values
4. Use config for future similar projects

## Best Practices

1. **Version Control Configs**: Keep configuration files in git
2. **Template Configs**: Create templates for common project types
3. **Environment Variables**: Use for sensitive data (planned feature)
4. **Validation**: Test configs in development before production
5. **Documentation**: Comment complex configurations

## Security Considerations

- Configuration files may contain sensitive information
- Use appropriate file permissions (600 or 640)
- Don't commit files with real credentials
- Use environment variables for secrets (planned)

## Future Enhancements

Planned features:
- Environment variable substitution
- Config file validation command
- Dry-run mode
- Config generation from existing projects
- Remote config file support (URLs)
- Config file encryption support

## Support

For issues or questions:
- Check example configuration: `setup-config-example.yaml`
- Review this documentation
- Run with `--verbose` flag for debugging
- Submit issues to project repository