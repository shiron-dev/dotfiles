package brew

import (
	"fmt"
	"strconv"
	"strings"

	"brew-manager/pkg/types"
	"brew-manager/pkg/utils"
)

// InstallPackages installs packages based on configuration and options
func InstallPackages(packages []types.Package, options *types.InstallOptions) error {
	if err := utils.CheckPrerequisites(); err != nil {
		return err
	}

	// Group packages by type
	packagesByType := make(map[string][]types.Package)
	for _, pkg := range packages {
		packagesByType[pkg.Type] = append(packagesByType[pkg.Type], pkg)
	}

	// Install in order: taps, brews, casks, mas
	order := []string{"tap", "brew", "cask", "mas"}
	
	for _, pkgType := range order {
		packages := packagesByType[pkgType]
		if len(packages) == 0 {
			continue
		}

		if shouldSkipType(pkgType, options) {
			utils.PrintStatus(utils.Yellow, fmt.Sprintf("Skipping %s packages as requested", pkgType))
			continue
		}

		if shouldInstallOnlyType(pkgType, options) || (!options.TapsOnly && !options.BrewsOnly && !options.CasksOnly && !options.MasOnly) {
			if err := installPackagesByType(pkgType, packages, options); err != nil {
				return fmt.Errorf("failed to install %s packages: %w", pkgType, err)
			}
		}
	}

	return nil
}

// shouldSkipType checks if a package type should be skipped
func shouldSkipType(pkgType string, options *types.InstallOptions) bool {
	switch pkgType {
	case "tap":
		return options.SkipTaps
	case "brew":
		return options.SkipBrews
	case "cask":
		return options.SkipCasks
	case "mas":
		return options.SkipMas
	}
	return false
}

// shouldInstallOnlyType checks if only this package type should be installed
func shouldInstallOnlyType(pkgType string, options *types.InstallOptions) bool {
	switch pkgType {
	case "tap":
		return options.TapsOnly
	case "brew":
		return options.BrewsOnly
	case "cask":
		return options.CasksOnly
	case "mas":
		return options.MasOnly
	}
	return false
}

// installPackagesByType installs packages of a specific type
func installPackagesByType(pkgType string, packages []types.Package, options *types.InstallOptions) error {
	utils.PrintStatus(utils.Blue, fmt.Sprintf("Installing %s packages...", pkgType))

	for _, pkg := range packages {
		if options.Verbose {
			utils.PrintStatus(utils.Cyan, fmt.Sprintf("Processing %s: %s", pkgType, pkg.Name))
		}

		if options.DryRun {
			utils.PrintStatus(utils.Yellow, fmt.Sprintf("[DRY RUN] Would install %s: %s", pkgType, pkg.Name))
			continue
		}

		if err := installSinglePackage(pkgType, pkg, options.Verbose); err != nil {
			utils.PrintStatus(utils.Red, fmt.Sprintf("Failed to install %s: %s - %v", pkgType, pkg.Name, err))
			continue
		}

		utils.PrintStatus(utils.Green, fmt.Sprintf("Installed %s: %s", pkgType, pkg.Name))
	}

	return nil
}

// installSinglePackage installs a single package
func installSinglePackage(pkgType string, pkg types.Package, verbose bool) error {
	switch pkgType {
	case "tap":
		return installTap(pkg.Name, verbose)
	case "brew":
		return installBrew(pkg.Name, verbose)
	case "cask":
		return installCask(pkg.Name, verbose)
	case "mas":
		return installMas(pkg.Name, pkg.ID, verbose)
	default:
		return fmt.Errorf("unknown package type: %s", pkgType)
	}
}

// installTap installs a tap
func installTap(name string, verbose bool) error {
	// Check if already installed
	if output, err := utils.RunCommand("brew", "tap"); err == nil {
		installedTaps := strings.Fields(output)
		for _, installed := range installedTaps {
			if installed == name {
				if verbose {
					utils.PrintStatus(utils.Yellow, fmt.Sprintf("Tap already installed: %s", name))
				}
				return nil
			}
		}
	}

	return utils.RunCommandSilent("brew", "tap", name)
}

// installBrew installs a brew formula
func installBrew(name string, verbose bool) error {
	// Check if already installed
	if output, err := utils.RunCommand("brew", "list", "--formula"); err == nil {
		installedBrews := strings.Fields(output)
		for _, installed := range installedBrews {
			if installed == name || strings.HasPrefix(name, installed+"/") {
				if verbose {
					utils.PrintStatus(utils.Yellow, fmt.Sprintf("Formula already installed: %s", name))
				}
				return nil
			}
		}
	}

	return utils.RunCommandSilent("brew", "install", name)
}

// installCask installs a cask
func installCask(name string, verbose bool) error {
	// Check if already installed
	if output, err := utils.RunCommand("brew", "list", "--cask"); err == nil {
		installedCasks := strings.Fields(output)
		for _, installed := range installedCasks {
			if installed == name {
				if verbose {
					utils.PrintStatus(utils.Yellow, fmt.Sprintf("Cask already installed: %s", name))
				}
				return nil
			}
		}
	}

	return utils.RunCommandSilent("brew", "install", "--cask", name)
}

// installMas installs a Mac App Store app
func installMas(name string, id int64, verbose bool) error {
	if !utils.CommandExists("mas") {
		return fmt.Errorf("mas is not installed, cannot install Mac App Store apps")
	}

	// Check if already installed
	if output, err := utils.RunCommand("mas", "list"); err == nil {
		lines := strings.Split(strings.TrimSpace(output), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, strconv.FormatInt(id, 10)) {
				if verbose {
					utils.PrintStatus(utils.Yellow, fmt.Sprintf("App already installed: %s", name))
				}
				return nil
			}
		}
	}

	return utils.RunCommandSilent("mas", "install", strconv.FormatInt(id, 10))
}

 
