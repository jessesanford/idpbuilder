package version

import (
	"fmt"
	"runtime"

	"github.com/cnoe-io/idpbuilder/pkg/cmd/helpers"
	"github.com/spf13/cobra"
)

var (
	// Version information (set during build)
	Version   = "dev"
	GitCommit = "unknown"
	GitBranch = "unknown"
	BuildDate = "unknown"
	GoVersion = runtime.Version()

	// Command flags
	short  bool
	output string
)

// VersionCmd represents the version command
var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display version information",
	Long: `Display version information for idpbuilder including:
- Version number
- Git commit and branch
- Build date and Go version
- Platform information

Examples:
  # Show full version information
  idpbuilder version

  # Show short version only
  idpbuilder version --short

  # Show version in JSON format
  idpbuilder version -o json
`,
	RunE: runVersion,
}

// VersionInfo holds version information
type VersionInfo struct {
	Version    string `json:"version" yaml:"version"`
	GitCommit  string `json:"gitCommit" yaml:"gitCommit"`
	GitBranch  string `json:"gitBranch" yaml:"gitBranch"`
	BuildDate  string `json:"buildDate" yaml:"buildDate"`
	GoVersion  string `json:"goVersion" yaml:"goVersion"`
	Platform   string `json:"platform" yaml:"platform"`
	Arch       string `json:"arch" yaml:"arch"`
}

func init() {
	VersionCmd.Flags().BoolVar(&short, "short", false, "Show only the version number")
	VersionCmd.Flags().StringVarP(&output, "output", "o", "", "Output format (json, yaml)")

	// Register flag completion
	VersionCmd.RegisterFlagCompletionFunc("output", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"json", "yaml"}, cobra.ShellCompDirectiveNoFileComp
	})
}

// runVersion executes the version command
func runVersion(cmd *cobra.Command, args []string) error {
	// Show short version if requested
	if short {
		fmt.Println(Version)
		return nil
	}

	// Build version information
	versionInfo := VersionInfo{
		Version:   Version,
		GitCommit: GitCommit,
		GitBranch: GitBranch,
		BuildDate: BuildDate,
		GoVersion: GoVersion,
		Platform:  runtime.GOOS,
		Arch:      runtime.GOARCH,
	}

	// Display based on output format
	if output != "" {
		return displayVersionFormatted(versionInfo)
	}

	return displayVersionDefault(versionInfo)
}

// displayVersionDefault displays version information in default format
func displayVersionDefault(info VersionInfo) error {
	fmt.Printf("idpbuilder version: %s\n", info.Version)
	fmt.Printf("Git commit: %s\n", info.GitCommit)
	fmt.Printf("Git branch: %s\n", info.GitBranch)
	fmt.Printf("Build date: %s\n", info.BuildDate)
	fmt.Printf("Go version: %s\n", info.GoVersion)
	fmt.Printf("Platform: %s/%s\n", info.Platform, info.Arch)
	return nil
}

// displayVersionFormatted displays version information in specified format
func displayVersionFormatted(info VersionInfo) error {
	// Validate output format
	format, err := helpers.ValidateOutputFormat(output)
	if err != nil {
		return fmt.Errorf("invalid output format: %w", err)
	}

	// Create printer and display
	printer := helpers.NewPrinter(format)
	return printer.Print(info)
}