/*
Copyright 2024 The idpbuilder Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package root

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	config *Config
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "idpbuilder-push",
	Short: "Push and manage container images to registries",
	Long: `idpbuilder-push is a CLI tool for building and pushing container images to registries.
It supports multiple build strategies including Dockerfile, buildpacks, and Kaniko.

This tool is designed to work with Kubernetes environments and can integrate with
GitRepository and LocalBuild custom resources.`,
	Version: "v0.1.0",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return config.Validate()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&config.KubeConfig, "kubeconfig", "", "path to kubeconfig file (default: $KUBECONFIG or ~/.kube/config)")
	rootCmd.PersistentFlags().StringVar(&config.Context, "context", "", "kubeconfig context to use")
	rootCmd.PersistentFlags().StringVarP(&config.Namespace, "namespace", "n", "default", "kubernetes namespace")
	rootCmd.PersistentFlags().BoolVarP(&config.Verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolVar(&config.DryRun, "dry-run", false, "dry-run mode (don't actually push)")

	// Registry flags
	rootCmd.PersistentFlags().StringVar(&config.Registry.URL, "registry", "localhost:5000", "registry URL")

	// Build flags
	rootCmd.PersistentFlags().StringVar(&config.Build.Strategy, "build-strategy", "dockerfile", "build strategy")
	rootCmd.PersistentFlags().StringVar(&config.Build.Context, "build-context", ".", "build context directory")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	config = NewConfig()
}