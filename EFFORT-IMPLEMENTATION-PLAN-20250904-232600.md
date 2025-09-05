# Effort Implementation Plan: CLI Commands Implementation

## 🚨 CRITICAL EFFORT METADATA (FROM WAVE PLAN)
**Effort**: E2.2.1 - CLI Commands Implementation  
**Branch**: `idpbuilder-oci-go-cr/phase2/wave2/cli-commands`  
**Can Parallelize**: No (depends on Wave 1 completion)  
**Parallel With**: None (single effort in Wave 2)  
**Size Estimate**: 600 lines (MUST be <800)  
**Dependencies**: Phase 2 Wave 1 (go-containerregistry builder and gitea client)  

## 📋 Source Information
**Wave Plan**: PHASE-2-WAVE-2-IMPLEMENTATION-PLAN.md  
**Effort Section**: E2.2.1  
**Created By**: Code Reviewer Agent  
**Date**: 2025-09-04  
**Extracted**: 2025-09-04 23:26:00 UTC  

## 🚀 Parallelization Context
**Can Parallelize**: No  
**Parallel With**: N/A  
**Blocking Status**: N/A - Single effort in wave  
**Parallel Group**: N/A  
**Orchestrator Guidance**: Spawn after Wave 1 integration complete  

## 📁 Files to Create

### Primary Implementation Files
```yaml
new_files:
  - path: cmd/build.go
    lines: ~150
    purpose: Implement 'idpbuilder build' command using cobra
    exports:
      - buildCmd (cobra.Command)
      - runBuild (function)
      
  - path: cmd/push.go
    lines: ~150
    purpose: Implement 'idpbuilder push' command using cobra
    exports:
      - pushCmd (cobra.Command)
      - runPush (function)
      
  - path: cmd/flags.go
    lines: ~80
    purpose: Common flag definitions shared across commands
    exports:
      - insecureFlag (bool)
      - contextFlag (string)
      - tagFlag (string)
      - configFlag (string)
      
  - path: pkg/cli/config.go
    lines: ~120
    purpose: Configuration management using viper
    exports:
      - Config (struct)
      - LoadConfig (function)
      - SaveConfig (function)
      
  - path: pkg/cli/progress.go
    lines: ~100
    purpose: Progress reporting for build and push operations
    exports:
      - ProgressReporter (interface)
      - NewProgressBar (function)
      - UpdateProgress (function)
```

### Test Files
```yaml
test_files:
  - path: cmd/build_test.go
    coverage_target: 80%
    test_types:
      - unit
      - command parsing
      
  - path: cmd/push_test.go
    coverage_target: 80%
    test_types:
      - unit
      - command parsing
      
  - path: pkg/cli/config_test.go
    coverage_target: 80%
    test_types:
      - unit
      - configuration loading
      
  - path: pkg/cli/progress_test.go
    coverage_target: 80%
    test_types:
      - unit
      - progress updates
```

## 📦 Files to Import/Reuse

### From Phase 2 Wave 1 Integration
```yaml
wave1_imports:
  - source: pkg/builder/builder.go
    from_effort: E2.1.1 (go-containerregistry-image-builder)
    usage: Use Builder interface for image assembly
    key_types:
      - Builder interface
      - BuildOptions struct
    
  - source: pkg/registry/gitea_client.go
    from_effort: E2.1.2 (gitea-registry-client)
    usage: Use GiteaClient for registry push operations
    key_types:
      - GiteaClient interface
      - PushOptions struct
      
  - source: pkg/registry/auth.go
    from_effort: E2.1.2
    usage: Authentication for registry operations
    key_types:
      - Authenticator interface
```

### From Phase 1 (Certificate Infrastructure)
```yaml
phase1_imports:
  - source: pkg/certs/trust.go
    usage: Certificate trust configuration
    key_types:
      - TrustManager interface
      
  - source: pkg/certs/extractor.go
    usage: Certificate extraction from Kind
    key_types:
      - CertExtractor interface
```

### External Dependencies (from go.mod)
```yaml
external_imports:
  - package: github.com/spf13/cobra
    version: v1.8.1
    usage: CLI command framework
    
  - package: github.com/spf13/viper
    version: v1.19.0
    usage: Configuration management
    
  - package: github.com/sirupsen/logrus
    version: v1.9.3
    usage: Structured logging
```

## 🔗 Dependencies

### Effort Dependencies
- **Must Complete First**: Phase 2 Wave 1 Integration (contains builder and registry client)
- **Can Run in Parallel With**: None (single effort in wave)
- **Blocks**: Phase 2 Integration

### Technical Dependencies
- Builder from E2.1.1 (image assembly)
- GiteaClient from E2.1.2 (registry operations)
- Certificate infrastructure from Phase 1

## 📝 Implementation Instructions

### Step-by-Step Guide

1. **Setup and Verification**
   ```bash
   # Verify correct directory and branch
   cd efforts/phase2/wave2/cli-commands
   git branch --show-current  # Should be: idpbuilder-oci-go-cr/phase2/wave2/cli-commands
   
   # Verify Wave 1 integration is available
   git log --oneline -5 idpbuilder-oci-go-cr/phase2/wave1-integration-20250904-212505
   ```

2. **Create Command Structure**
   ```bash
   # Create cmd directory
   mkdir -p cmd
   mkdir -p pkg/cli
   ```

3. **Implement Build Command (cmd/build.go)**
   ```go
   package cmd
   
   import (
       "context"
       "fmt"
       "os"
       
       "github.com/spf13/cobra"
       "github.com/jessesanford/idpbuilder/pkg/builder"
       "github.com/jessesanford/idpbuilder/pkg/cli"
   )
   
   var buildCmd = &cobra.Command{
       Use:   "build",
       Short: "Assemble OCI image from context directory",
       Long:  `Assemble a single-layer OCI image from a directory using go-containerregistry.
               The image is stored locally as an OCI tarball.`,
       Example: `  idpbuilder build --context ./app --tag myapp:v1
     idpbuilder build --context . --tag myimage:latest`,
       RunE: runBuild,
   }
   
   func init() {
       buildCmd.Flags().String("context", ".", "Build context directory")
       buildCmd.Flags().String("tag", "", "Image tag (required)")
       buildCmd.MarkFlagRequired("tag")
   }
   
   func runBuild(cmd *cobra.Command, args []string) error {
       contextPath, _ := cmd.Flags().GetString("context")
       tag, _ := cmd.Flags().GetString("tag")
       
       // Create progress reporter
       progress := cli.NewProgressBar("Building image")
       defer progress.Finish()
       
       // Create builder with options
       b := builder.New(builder.WithProgress(progress))
       
       // Build the image
       progress.UpdateMessage("Assembling layers from " + contextPath)
       if err := b.BuildImage(context.Background(), contextPath, tag); err != nil {
           return fmt.Errorf("build failed: %w", err)
       }
       
       progress.UpdateMessage("Image built successfully: " + tag)
       return nil
   }
   ```

4. **Implement Push Command (cmd/push.go)**
   ```go
   package cmd
   
   import (
       "context"
       "fmt"
       
       "github.com/spf13/cobra"
       "github.com/jessesanford/idpbuilder/pkg/registry"
       "github.com/jessesanford/idpbuilder/pkg/cli"
       "github.com/jessesanford/idpbuilder/pkg/certs"
   )
   
   var pushCmd = &cobra.Command{
       Use:   "push IMAGE[:TAG]",
       Short: "Push image to Gitea registry",
       Long:  `Push a container image to the builtin Gitea registry with certificate support.
               Automatically handles certificate trust configuration.`,
       Example: `  idpbuilder push myapp:v1
     idpbuilder push --insecure myapp:latest`,
       Args: cobra.ExactArgs(1),
       RunE: runPush,
   }
   
   func init() {
       pushCmd.Flags().Bool("insecure", false, "Skip TLS certificate verification")
   }
   
   func runPush(cmd *cobra.Command, args []string) error {
       image := args[0]
       insecure, _ := cmd.Flags().GetBool("insecure")
       
       // Create progress reporter
       progress := cli.NewProgressBar("Pushing " + image)
       defer progress.Finish()
       
       // Setup certificate trust (unless --insecure)
       var trustManager certs.TrustManager
       if !insecure {
           progress.UpdateMessage("Configuring certificate trust")
           trustManager = certs.NewTrustManager()
           if err := trustManager.SetupGiteaTrust(); err != nil {
               return fmt.Errorf("certificate setup failed: %w", err)
           }
       }
       
       // Create registry client
       client := registry.NewGiteaClient(
           registry.WithProgress(progress),
           registry.WithTrustManager(trustManager),
           registry.WithInsecure(insecure),
       )
       
       // Push the image
       progress.UpdateMessage("Pushing to registry")
       if err := client.Push(context.Background(), image); err != nil {
           return fmt.Errorf("push failed: %w", err)
       }
       
       progress.UpdateMessage("Successfully pushed " + image)
       return nil
   }
   ```

5. **Implement Common Flags (cmd/flags.go)**
   ```go
   package cmd
   
   import "github.com/spf13/cobra"
   
   var (
       // Global flags
       configFile string
       verbose    bool
       quiet      bool
   )
   
   // AddGlobalFlags adds global flags to the root command
   func AddGlobalFlags(cmd *cobra.Command) {
       cmd.PersistentFlags().StringVar(&configFile, "config", "", "Config file path")
       cmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
       cmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "Quiet mode")
   }
   
   // AddBuildFlags adds build-specific flags
   func AddBuildFlags(cmd *cobra.Command) {
       cmd.Flags().String("exclude", "", "Exclusion patterns (comma-separated)")
       cmd.Flags().String("platform", "linux/amd64", "Target platform")
   }
   
   // AddPushFlags adds push-specific flags
   func AddPushFlags(cmd *cobra.Command) {
       cmd.Flags().String("registry", "", "Registry URL (default: auto-detect)")
       cmd.Flags().Int("retry", 3, "Number of retry attempts")
   }
   ```

6. **Configuration Management (pkg/cli/config.go)**
   ```go
   package cli
   
   import (
       "fmt"
       "os"
       "path/filepath"
       
       "github.com/spf13/viper"
   )
   
   type Config struct {
       Registry   RegistryConfig   `mapstructure:"registry"`
       Build      BuildConfig      `mapstructure:"build"`
       Logging    LoggingConfig    `mapstructure:"logging"`
   }
   
   type RegistryConfig struct {
       URL      string `mapstructure:"url"`
       Username string `mapstructure:"username"`
       Insecure bool   `mapstructure:"insecure"`
   }
   
   type BuildConfig struct {
       Context   string   `mapstructure:"context"`
       Exclude   []string `mapstructure:"exclude"`
       CacheDir  string   `mapstructure:"cache_dir"`
   }
   
   type LoggingConfig struct {
       Level  string `mapstructure:"level"`
       Format string `mapstructure:"format"`
   }
   
   // LoadConfig loads configuration from file and environment
   func LoadConfig(configPath string) (*Config, error) {
       viper.SetConfigType("yaml")
       viper.SetEnvPrefix("IDPBUILDER")
       viper.AutomaticEnv()
       
       // Set defaults
       viper.SetDefault("registry.url", "auto")
       viper.SetDefault("build.cache_dir", filepath.Join(os.Getenv("HOME"), ".idpbuilder/cache"))
       viper.SetDefault("logging.level", "info")
       viper.SetDefault("logging.format", "text")
       
       // Load config file if specified
       if configPath != "" {
           viper.SetConfigFile(configPath)
           if err := viper.ReadInConfig(); err != nil {
               return nil, fmt.Errorf("failed to read config: %w", err)
           }
       }
       
       var config Config
       if err := viper.Unmarshal(&config); err != nil {
           return nil, fmt.Errorf("failed to unmarshal config: %w", err)
       }
       
       return &config, nil
   }
   ```

7. **Progress Reporting (pkg/cli/progress.go)**
   ```go
   package cli
   
   import (
       "fmt"
       "io"
       "os"
       "sync"
   )
   
   // ProgressReporter provides progress updates
   type ProgressReporter interface {
       UpdateMessage(message string)
       UpdateProgress(current, total int64)
       Finish()
   }
   
   // ProgressBar implements a simple progress bar
   type ProgressBar struct {
       mu      sync.Mutex
       writer  io.Writer
       message string
       current int64
       total   int64
   }
   
   // NewProgressBar creates a new progress reporter
   func NewProgressBar(initialMessage string) ProgressReporter {
       pb := &ProgressBar{
           writer:  os.Stdout,
           message: initialMessage,
       }
       pb.render()
       return pb
   }
   
   func (pb *ProgressBar) UpdateMessage(message string) {
       pb.mu.Lock()
       defer pb.mu.Unlock()
       pb.message = message
       pb.render()
   }
   
   func (pb *ProgressBar) UpdateProgress(current, total int64) {
       pb.mu.Lock()
       defer pb.mu.Unlock()
       pb.current = current
       pb.total = total
       pb.render()
   }
   
   func (pb *ProgressBar) render() {
       // Clear line and render progress
       fmt.Fprintf(pb.writer, "\r%-50s", pb.message)
       if pb.total > 0 {
           percent := int((pb.current * 100) / pb.total)
           fmt.Fprintf(pb.writer, " [%3d%%]", percent)
       }
   }
   
   func (pb *ProgressBar) Finish() {
       fmt.Fprintln(pb.writer) // New line after progress
   }
   ```

## ✅ Test Requirements

### Coverage Requirements
- **Minimum Coverage**: 80%
- **Critical Paths**: 100% coverage required
- **Error Handling**: All error cases must be tested

### Test Categories
```yaml
required_tests:
  unit_tests:
    - Command parsing and validation
    - Configuration loading
    - Progress reporting
    - Flag handling
    
  integration_tests:
    - Full build command workflow
    - Full push command workflow
    - Certificate handling integration
    - Error scenario handling
    
  cli_tests:
    - Help text validation
    - Flag combination testing
    - Exit code verification
```

## 📏 Size Constraints
**Target Size**: 600 lines  
**Maximum Size**: 800 lines (HARD LIMIT)  
**Current Estimate**: ~600 lines  

### Size Monitoring Protocol
```bash
# Check size regularly during implementation
cd /home/vscode/workspaces/idpbuilder-oci-go-cr/efforts/phase2/wave2/cli-commands

# Find project root and use line counter
PROJECT_ROOT=/home/vscode/workspaces/idpbuilder-oci-go-cr
$PROJECT_ROOT/tools/line-counter.sh

# If approaching 700 lines:
# 1. Alert Code Reviewer
# 2. Focus on core functionality only
# 3. Defer nice-to-have features
```

## 🏁 Completion Criteria

### Implementation Checklist
- [ ] All command files created (build.go, push.go, flags.go)
- [ ] Configuration management implemented
- [ ] Progress reporting functional
- [ ] Size verified under 800 lines
- [ ] Integration with Wave 1 components complete

### Quality Checklist
- [ ] Test coverage ≥80%
- [ ] All tests passing
- [ ] No linting errors
- [ ] Error messages clear and helpful
- [ ] Commands have proper help text

### Integration Checklist
- [ ] Build command works end-to-end
- [ ] Push command works with Gitea
- [ ] Certificate handling integrated
- [ ] --insecure flag functional
- [ ] Configuration loading works

### Review Checklist
- [ ] Self-review completed
- [ ] Code committed and pushed
- [ ] Ready for Code Reviewer assessment
- [ ] Work log updated

## 📊 Progress Tracking

### Work Log Updates
The SW Engineer should update work-log.md with:
- Files created and their line counts
- Integration points verified
- Tests written and coverage achieved
- Any issues encountered

## ⚠️ Important Notes

### Key Integration Points
1. **Builder Interface**: Use the Builder from pkg/builder (Wave 1)
2. **Registry Client**: Use GiteaClient from pkg/registry (Wave 1)
3. **Certificate Trust**: Integrate with Phase 1 certificate infrastructure
4. **Error Handling**: Ensure clear error messages for certificate issues

### Command Structure
The commands should integrate into the existing idpbuilder CLI:
```
idpbuilder
├── create    (existing)
├── build     (new - this effort)
└── push      (new - this effort)
```

### Common Pitfalls to Avoid
1. **Size Limit**: Keep implementation focused on core functionality
2. **Dependencies**: Use existing Wave 1 components, don't reimplement
3. **Testing**: Write tests alongside implementation
4. **Error Messages**: Make certificate errors clear and actionable
5. **Progress**: Provide meaningful progress updates during operations

## 📚 References

### Source Documents
- [Master Implementation Plan](/home/vscode/workspaces/idpbuilder-oci-go-cr/IMPLEMENTATION-PLAN.md)
- [Phase 2 Implementation Plan](/home/vscode/workspaces/idpbuilder-oci-go-cr/phase-plans/phase2/PHASE-2-IMPLEMENTATION-PLAN.md)
- [Wave 1 Integration](/home/vscode/workspaces/idpbuilder-oci-go-cr/efforts/phase2/wave1/integration-workspace/)

### Code References
- [Wave 1 Builder](/home/vscode/workspaces/idpbuilder-oci-go-cr/pkg/builder/)
- [Wave 1 Registry](/home/vscode/workspaces/idpbuilder-oci-go-cr/pkg/registry/)
- [Phase 1 Certs](/home/vscode/workspaces/idpbuilder-oci-go-cr/pkg/certs/)

### External Documentation
- [Cobra Documentation](https://cobra.dev/)
- [Viper Documentation](https://github.com/spf13/viper)
- [go-containerregistry](https://github.com/google/go-containerregistry)

---

**Remember**: This is the ONLY effort in Wave 2. Focus on creating clean, user-friendly CLI commands that leverage the existing Wave 1 infrastructure. The goal is to provide a seamless experience for building and pushing OCI images to Gitea with automatic certificate handling.