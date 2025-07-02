package cmd

import (
	"fmt"
	"strings"

	"brew-manager/pkg/types"
	"brew-manager/pkg/utils"
	"brew-manager/pkg/yaml"

	"github.com/spf13/cobra"
)

var (
	skipTapsInPrune  bool
	skipBrewsInPrune bool
	skipCasksInPrune bool
	skipMasInPrune   bool
	confirmAll       bool
)

// pruneCmd represents the prune command
var pruneCmd = &cobra.Command{
	Use:   "prune [yaml_file]",
	Short: "Remove packages not defined in YAML configuration",
	Long: `Remove Homebrew packages that are currently installed but not defined in the YAML configuration file.

This command will:
1. Read the YAML configuration file
2. Compare with currently installed packages
3. Remove packages that are not in the configuration

Examples:
  brew-manager prune                        # Use default packages.yaml
  brew-manager prune packages.yaml          # Use specific YAML file
  brew-manager prune --dry-run              # Show what would be removed
  brew-manager prune --skip-brews           # Only remove casks, taps, and mas apps
  brew-manager prune --confirm-all          # Remove all without individual confirmation`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get YAML file path
		yamlFile := getDefaultYAMLPath("packages.yaml")
		if len(args) > 0 {
			yamlFile = args[0]
		}

		// Build prune options
		options := &types.PruneOptions{
			DryRun:     dryRun,
			Verbose:    verbose,
			SkipTaps:   skipTapsInPrune,
			SkipBrews:  skipBrewsInPrune,
			SkipCasks:  skipCasksInPrune,
			SkipMas:    skipMasInPrune,
			ConfirmAll: confirmAll,
		}

		// Perform prune
		if err := prunePackages(yamlFile, options); err != nil {
			utils.PrintStatus(utils.Red, fmt.Sprintf("Prune failed: %v", err))
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(pruneCmd)

	// Prune-specific flags
	pruneCmd.Flags().BoolVar(&skipTapsInPrune, "skip-taps", false, "Skip removing taps")
	pruneCmd.Flags().BoolVar(&skipBrewsInPrune, "skip-brews", false, "Skip removing brew formulae")
	pruneCmd.Flags().BoolVar(&skipCasksInPrune, "skip-casks", false, "Skip removing casks")
	pruneCmd.Flags().BoolVar(&skipMasInPrune, "skip-mas", false, "Skip removing Mac App Store apps")
	pruneCmd.Flags().BoolVar(&confirmAll, "confirm-all", false, "Remove all packages without individual confirmation")
}

// prunePackages removes packages not defined in the YAML configuration
func prunePackages(yamlFile string, options *types.PruneOptions) error {
	// Load YAML configuration
	config, err := yaml.LoadGroupedConfig(yamlFile)
	if err != nil {
		return fmt.Errorf("failed to load YAML configuration: %w", err)
	}

	// Get all packages from YAML
	yamlPackages := getAllPackagesFromConfig(config)

	// Get currently installed packages
	installedPackages, installedMasApps, err := yaml.GetInstalledPackages()
	if err != nil {
		return fmt.Errorf("failed to get installed packages: %w", err)
	}

	// Find packages to remove
	packagesToRemove := findPackagesToRemove(yamlPackages, installedPackages, installedMasApps, options)

	if len(packagesToRemove["taps"]) == 0 && len(packagesToRemove["brews"]) == 0 &&
		len(packagesToRemove["casks"]) == 0 && len(packagesToRemove["mas"]) == 0 {
		utils.PrintStatus(utils.Green, "No packages to remove. All installed packages are defined in YAML configuration.")
		return nil
	}

	// Show what will be removed
	showRemovalSummary(packagesToRemove)

	if options.DryRun {
		utils.PrintStatus(utils.Yellow, "[DRY RUN] No packages were actually removed.")
		return nil
	}

	// Confirm removal unless --confirm-all is used
	if !options.ConfirmAll {
		if !confirmRemoval() {
			utils.PrintStatus(utils.Yellow, "Prune operation cancelled.")
			return nil
		}
	}

	// Remove packages in reverse order: mas, casks, brews, taps
	order := []string{"mas", "cask", "brew", "tap"}
	for _, pkgType := range order {
		if err := removePackagesByType(pkgType, packagesToRemove[pkgType], options); err != nil {
			return fmt.Errorf("failed to remove %s packages: %w", pkgType, err)
		}
	}

	utils.PrintStatus(utils.Green, "Prune operation completed successfully.")
	return nil
}

// getAllPackagesFromConfig extracts all packages from the configuration
func getAllPackagesFromConfig(config *types.PackageGrouped) map[string]map[string]bool {
	result := map[string]map[string]bool{
		"tap":  make(map[string]bool),
		"brew": make(map[string]bool),
		"cask": make(map[string]bool),
		"mas":  make(map[string]bool),
	}

	for _, group := range config.Groups {
		for _, pkg := range group.Packages {
			switch pkg.Type {
			case "tap":
				result["tap"][pkg.Name] = true
			case "brew":
				result["brew"][pkg.Name] = true
			case "cask":
				result["cask"][pkg.Name] = true
			case "mas":
				result["mas"][fmt.Sprintf("%d", pkg.ID)] = true
			}
		}
	}

	return result
}

// findPackagesToRemove identifies packages to remove
func findPackagesToRemove(yamlPackages map[string]map[string]bool,
	installedPackages map[string][]string, installedMasApps []types.MasApp,
	options *types.PruneOptions) map[string][]string {

	result := map[string][]string{
		"taps":  []string{},
		"brews": []string{},
		"casks": []string{},
		"mas":   []string{},
	}

	// Check taps
	if !options.SkipTaps {
		for _, installed := range installedPackages["taps"] {
			if !yamlPackages["tap"][installed] {
				result["taps"] = append(result["taps"], installed)
			}
		}
	}

	// Check brews
	if !options.SkipBrews {
		for _, installed := range installedPackages["brews"] {
			if !yamlPackages["brew"][installed] {
				result["brews"] = append(result["brews"], installed)
			}
		}
	}

	// Check casks
	if !options.SkipCasks {
		for _, installed := range installedPackages["casks"] {
			if !yamlPackages["cask"][installed] {
				result["casks"] = append(result["casks"], installed)
			}
		}
	}

	// Check mas apps
	if !options.SkipMas {
		for _, installed := range installedMasApps {
			idStr := fmt.Sprintf("%d", installed.ID)
			if !yamlPackages["mas"][idStr] {
				result["mas"] = append(result["mas"], fmt.Sprintf("%d (%s)", installed.ID, installed.Name))
			}
		}
	}

	return result
}

// showRemovalSummary displays what will be removed
func showRemovalSummary(packagesToRemove map[string][]string) {
	utils.PrintStatus(utils.Blue, "Packages to be removed:")

	if len(packagesToRemove["taps"]) > 0 {
		utils.PrintStatus(utils.Yellow, fmt.Sprintf("Taps (%d):", len(packagesToRemove["taps"])))
		for _, tap := range packagesToRemove["taps"] {
			fmt.Printf("  - %s\n", tap)
		}
	}

	if len(packagesToRemove["brews"]) > 0 {
		utils.PrintStatus(utils.Yellow, fmt.Sprintf("Brew formulae (%d):", len(packagesToRemove["brews"])))
		for _, brew := range packagesToRemove["brews"] {
			fmt.Printf("  - %s\n", brew)
		}
	}

	if len(packagesToRemove["casks"]) > 0 {
		utils.PrintStatus(utils.Yellow, fmt.Sprintf("Casks (%d):", len(packagesToRemove["casks"])))
		for _, cask := range packagesToRemove["casks"] {
			fmt.Printf("  - %s\n", cask)
		}
	}

	if len(packagesToRemove["mas"]) > 0 {
		utils.PrintStatus(utils.Yellow, fmt.Sprintf("Mac App Store apps (%d):", len(packagesToRemove["mas"])))
		for _, mas := range packagesToRemove["mas"] {
			fmt.Printf("  - %s\n", mas)
		}
	}
}

// confirmRemoval asks for user confirmation
func confirmRemoval() bool {
	fmt.Print("\nAre you sure you want to remove these packages? [y/N]: ")
	var response string
	fmt.Scanln(&response)
	response = strings.ToLower(strings.TrimSpace(response))
	return response == "y" || response == "yes"
}

// removePackagesByType removes packages of a specific type
func removePackagesByType(pkgType string, packages []string, options *types.PruneOptions) error {
	if len(packages) == 0 {
		return nil
	}

	utils.PrintStatus(utils.Blue, fmt.Sprintf("Removing %s packages...", pkgType))

	for _, pkg := range packages {
		if options.Verbose {
			utils.PrintStatus(utils.Cyan, fmt.Sprintf("Processing %s: %s", pkgType, pkg))
		}

		if err := removeSinglePackage(pkgType, pkg, options.Verbose); err != nil {
			utils.PrintStatus(utils.Red, fmt.Sprintf("Failed to remove %s: %s - %v", pkgType, pkg, err))
			continue
		}

		utils.PrintStatus(utils.Green, fmt.Sprintf("Removed %s: %s", pkgType, pkg))
	}

	return nil
}

// removeSinglePackage removes a single package
func removeSinglePackage(pkgType string, pkg string, verbose bool) error {
	switch pkgType {
	case "tap":
		return removeTap(pkg, verbose)
	case "brew":
		return removeBrew(pkg, verbose)
	case "cask":
		return removeCask(pkg, verbose)
	case "mas":
		// Extract ID from "ID (Name)" format
		parts := strings.SplitN(pkg, " ", 2)
		if len(parts) > 0 {
			return removeMas(parts[0], verbose)
		}
		return fmt.Errorf("invalid mas package format: %s", pkg)
	default:
		return fmt.Errorf("unknown package type: %s", pkgType)
	}
}

// removeTap removes a tap
func removeTap(name string, verbose bool) error {
	return utils.RunCommandSilent("brew", "untap", name)
}

// removeBrew removes a brew formula
func removeBrew(name string, verbose bool) error {
	return utils.RunCommandSilent("brew", "uninstall", name)
}

// removeCask removes a cask
func removeCask(name string, verbose bool) error {
	return utils.RunCommandSilent("brew", "uninstall", "--cask", name)
}

// removeMas removes a Mac App Store app
func removeMas(id string, verbose bool) error {
	if !utils.CommandExists("mas") {
		return fmt.Errorf("mas is not installed, cannot remove Mac App Store apps")
	}
	return utils.RunCommandSilent("mas", "uninstall", id)
}
