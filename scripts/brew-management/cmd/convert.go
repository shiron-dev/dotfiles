package cmd

import (
	"fmt"

	"brew-manager/pkg/convert"
	"brew-manager/pkg/utils"

	"github.com/spf13/cobra"
)

var (
	grouped bool
)

// convertCmd represents the convert command
var convertCmd = &cobra.Command{
	Use:   "convert <brewfile> <yaml_file>",
	Short: "Convert Brewfile to YAML format",
	Long: `Convert a Brewfile to YAML format.

Examples:
  brew-manager convert Brewfile packages.yaml      # Convert to grouped YAML format
  brew-manager convert --grouped Brewfile packages.yaml  # Convert to grouped YAML format`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		brewfilePath := args[0]
		yamlPath := args[1]

		if err := convert.ConvertBrewfileToYAML(brewfilePath, yamlPath, grouped, verbose); err != nil {
			utils.PrintStatus(utils.Red, fmt.Sprintf("Conversion failed: %v", err))
			return
		}

		utils.PrintStatus(utils.Green, fmt.Sprintf("Successfully converted %s to %s", brewfilePath, yamlPath))
	},
}

func init() {
	rootCmd.AddCommand(convertCmd)

	// Convert options
	convertCmd.Flags().BoolVarP(&grouped, "grouped", "g", false, "Convert to grouped YAML format with auto-detection")
} 
