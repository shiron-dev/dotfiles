package cmd

import (
	"fmt"
	"path/filepath"

	"brew-manager/pkg/types"
	"brew-manager/pkg/utils"
	"brew-manager/pkg/validate"

	"github.com/spf13/cobra"
)

var (
	all        bool
	schemaFile string
)

// validateCmd represents the validate command
var validateCmd = &cobra.Command{
	Use:   "validate [yaml_file]",
	Short: "Validate YAML configuration files against their schemas",
	Long: `Validate YAML configuration files against their JSON schemas.

Examples:
  brew-manager validate                                          # Validate all YAML files
  brew-manager validate packages.yml                             # Validate specific file
  brew-manager validate --schema packages-grouped.schema.json packages-grouped.yml
  brew-manager validate --all --verbose                          # Validate all with verbose output`,
	Run: func(cmd *cobra.Command, args []string) {
		// Build validate options
		options := &types.ValidateOptions{
			Verbose:    verbose,
			All:        all,
			SchemaFile: schemaFile,
		}

		if all {
			// Validate all YAML files in data directory
			dataDir := filepath.Dir(getDefaultYAMLPath("packages.yaml"))
			
			if err := validate.ValidateAllYAMLFiles(dataDir, options); err != nil {
				utils.PrintStatus(utils.Red, fmt.Sprintf("Validation failed: %v", err))
				return
			}
		} else {
			// Validate specific file
			var yamlFile string
			if len(args) > 0 {
				yamlFile = args[0]
			} else {
				// Default to grouped config file
				yamlFile = getDefaultYAMLPath("packages.yaml")
			}

			if err := validate.ValidateYAMLFile(yamlFile, options); err != nil {
				utils.PrintStatus(utils.Red, fmt.Sprintf("Validation failed: %v", err))
				return
			}
		}

		utils.PrintStatus(utils.Green, "Validation completed successfully!")
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)

	// Validate options
	validateCmd.Flags().BoolVarP(&all, "all", "a", false, "Validate all YAML files")
	validateCmd.Flags().StringVar(&schemaFile, "schema", "", "Use specific schema file")
} 
