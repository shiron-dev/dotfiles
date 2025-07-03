package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"brew-manager/pkg/utils"

	"github.com/spf13/cobra"
)

var (
	verbose bool
	dryRun  bool
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "brew-manager",
	Short: "Unified brew package management tool with YAML configuration support",
	Long: `Unified brew package management tool with YAML configuration support.

This tool provides a unified interface for all brew management operations including:
- Installing packages from YAML configuration (with groups/tags support)
- Synchronizing installed packages to YAML configuration
- Converting Brewfile to YAML format
- Validating YAML configuration files
- Removing packages not defined in YAML configuration (prune)

Examples:
  brew-manager install --groups core,development
  brew-manager install --profile developer
  brew-manager sync --auto-detect
  brew-manager prune --dry-run
  brew-manager validate`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Check prerequisites for most commands, skip for validate, help, and completion
		commandName := cmd.Name()
		if commandName != "validate" && commandName != "help" && commandName != "completion" {
			if err := utils.CheckPrerequisites(); err != nil {
				utils.PrintStatus(utils.Red, fmt.Sprintf("Error: %v", err))
				os.Exit(1)
			}
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		utils.PrintStatus(utils.Red, fmt.Sprintf("Error: %v", err))
		os.Exit(1)
	}
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	rootCmd.PersistentFlags().BoolVarP(&dryRun, "dry-run", "d", false, "Show what would be done without actually doing it")
}

// getDefaultYAMLPath returns the default path for YAML configuration files
func getDefaultYAMLPath(filename string) string {
	// デフォルトのYAMLパス
	defaultPath := filepath.Join(os.Getenv("HOME"), "projects/github.com/shiron-dev/dotfiles/data/brew/packages.yaml")
	if filename == "" || filename == "packages-grouped.yml" || filename == "packages.yml" || filename == "packages.yaml" {
		return defaultPath
	}

	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return filepath.Join(".", "data", "brew", filename)
	}
	// Look for data directory relative to current location
	dataPath := filepath.Join(cwd, "../../data/brew", filename)
	if utils.FileExists(dataPath) {
		return dataPath
	}
	// Fall back to relative path
	return filepath.Join("data", "brew", filename)
}
