package cmd

import (
	"brew-manager/pkg/sync"
	"brew-manager/pkg/types"
	"brew-manager/pkg/utils"
	"fmt"

	"github.com/spf13/cobra"
)

var (
	backup       bool
	sortPackages bool
	showOnly     bool
	defaultGroup string
	defaultTags  string
	interactive  bool
	autoDetect   bool
)

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync [yaml_file]",
	Short: "Sync installed packages to YAML configuration (with groups/tags)",
	Long: `Sync currently installed Homebrew packages with YAML configuration file (with groups/tags support).
Adds missing packages to appropriate sections with optional group/tag assignment.

Examples:
  brew-manager sync                                     # Sync with default settings
  brew-manager sync --dry-run                          # Show what would be added
  brew-manager sync --backup --default-group system    # Add missing packages to 'system' group
  brew-manager sync --interactive                      # Prompt for group/tag for each package
  brew-manager sync --auto-detect --sort               # Auto-detect groups/tags and sort`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get YAML file path
		yamlFile := getDefaultYAMLPath("packages-grouped.yml")
		if len(args) > 0 {
			yamlFile = args[0]
		}

		// Build sync options
		options := &types.SyncOptions{
			DryRun:       dryRun,
			Verbose:      verbose,
			Backup:       backup,
			Sort:         sortPackages,
			ShowOnly:     showOnly,
			DefaultGroup: defaultGroup,
			DefaultTags:  utils.SplitCommaSeparated(defaultTags),
			Interactive:  interactive,
			AutoDetect:   autoDetect,
		}

		// Perform sync
		if err := sync.SyncGroupedPackages(yamlFile, options); err != nil {
			utils.PrintStatus(utils.Red, fmt.Sprintf("Sync failed: %v", err))
			return
		}
	},
}



func init() {
	rootCmd.AddCommand(syncCmd)


	// Sync options for grouped format
	syncCmd.Flags().BoolVarP(&backup, "backup", "b", false, "Create backup of YAML file before modification")
	syncCmd.Flags().BoolVarP(&sortPackages, "sort", "s", false, "Sort packages alphabetically within categories")
	syncCmd.Flags().BoolVar(&showOnly, "show-only", false, "Only show missing packages without modifying the file")
	syncCmd.Flags().StringVar(&defaultGroup, "default-group", "", "Assign this group to newly added packages (default: optional)")
	syncCmd.Flags().StringVar(&defaultTags, "default-tags", "", "Assign these tags to newly added packages (comma-separated)")
	syncCmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "Prompt for group/tag assignment for each new package")
	syncCmd.Flags().BoolVarP(&autoDetect, "auto-detect", "a", false, "Try to auto-detect appropriate groups/tags based on package names")


} 
