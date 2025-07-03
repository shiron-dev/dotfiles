package brew

import (
	"fmt"
	"strconv"
	"strings"

	"brew-manager/pkg/types"
	"brew-manager/pkg/utils"
)

// InstallPackages installs packages based on configuration and options
func InstallPackages(filteredPackages []types.FilteredPackage, options *types.InstallOptions) error {
	if err := utils.CheckPrerequisites(); err != nil {
		return err
	}

	// Group packages by type
	packagesByType := make(map[string][]types.PackageInfo)
	for _, filteredPkg := range filteredPackages {
		packagesByType[filteredPkg.Type] = append(packagesByType[filteredPkg.Type], filteredPkg.PackageInfo)
	}

	// Install in order: taps, brews, casks, mas
	order := []string{"tap", "brew", "cask", "mas"}

	for _, pkgType := range order {
		pkgInfos := packagesByType[pkgType]
		if len(pkgInfos) == 0 {
			continue
		}

		if shouldSkipType(pkgType, options) {
			utils.PrintStatus(utils.Yellow, fmt.Sprintf("Skipping %s packages as requested", pkgType))
			continue
		}

		if err := installPackagesByType(pkgType, pkgInfos, options); err != nil {
			return fmt.Errorf("failed to install %s packages: %w", pkgType, err)
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

// installPackagesByType installs packages of a specific type
func installPackagesByType(pkgType string, pkgInfos []types.PackageInfo, options *types.InstallOptions) error {
	utils.PrintStatus(utils.Blue, fmt.Sprintf("Installing %s packages...", pkgType))

	for _, pkgInfo := range pkgInfos {
		if options.Verbose {
			utils.PrintStatus(utils.Cyan, fmt.Sprintf("Processing %s: %s", pkgType, pkgInfo.Name))
		}

		if options.DryRun {
			utils.PrintStatus(utils.Yellow, fmt.Sprintf("[DRY RUN] Would install %s: %s", pkgType, pkgInfo.Name))
			continue
		}

		if err := installSinglePackage(pkgType, pkgInfo, options.Verbose); err != nil {
			utils.PrintStatus(utils.Red, fmt.Sprintf("Failed to install %s: %s - %v", pkgType, pkgInfo.Name, err))
			// Decide if we should continue or stop on error. For now, continue.
			continue
		}

		utils.PrintStatus(utils.Green, fmt.Sprintf("Installed %s: %s", pkgType, pkgInfo.Name))
	}

	return nil
}

// installSinglePackage installs a single package
func installSinglePackage(pkgType string, pkgInfo types.PackageInfo, verbose bool) error {
	switch pkgType {
	case "tap":
		return installTap(pkgInfo.Name, verbose)
	case "brew":
		return installBrew(pkgInfo.Name, verbose)
	case "cask":
		return installCask(pkgInfo.Name, verbose)
	case "mas":
		// 'mas' type requires ID. Ensure it's present.
		if pkgInfo.ID == 0 {
			return fmt.Errorf("missing ID for mas package: %s", pkgInfo.Name)
		}
		return installMas(pkgInfo.Name, pkgInfo.ID, verbose)
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
