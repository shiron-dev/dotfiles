package sync

import (
	"fmt"
	"sort"
	"strings"

	"brew-manager/pkg/types"
	"brew-manager/pkg/utils"
	yamlPkg "brew-manager/pkg/yaml"

	"github.com/AlecAivazis/survey/v2"
)

// SyncGroupedPackages synchronizes installed packages with grouped YAML config
func SyncGroupedPackages(filePath string, options *types.SyncOptions) error {
	if options.Backup {
		if err := utils.CreateBackup(filePath); err != nil {
			return fmt.Errorf("failed to create backup: %w", err)
		}
	}

	// Load existing config
	config, err := yamlPkg.LoadGroupedConfig(filePath)
	if err != nil {
		return fmt.Errorf("failed to load grouped config: %w", err)
	}

	// Get currently installed packages
	installedPackages, err := yamlPkg.GetInstalledPackages()
	if err != nil {
		return fmt.Errorf("failed to get installed packages: %w", err)
	}

	// Find missing packages
	missingPackages := findMissingPackages(config, installedPackages)
	
	if len(missingPackages) == 0 {
		utils.PrintStatus(utils.Green, "No new packages found. Configuration is up to date.")
		return nil
	}

	utils.PrintStatus(utils.Cyan, fmt.Sprintf("Found %d new packages", len(missingPackages)))

	if options.ShowOnly {
		showMissingPackages(missingPackages)
		return nil
	}

	if options.DryRun {
		utils.PrintStatus(utils.Yellow, "[DRY RUN] Would add the following packages:")
		showMissingPackages(missingPackages)
		return nil
	}

	// Add missing packages to config
	if err := addMissingPackagesToGrouped(config, missingPackages, options); err != nil {
		return fmt.Errorf("failed to add missing packages: %w", err)
	}

	// Sort packages if requested
	if options.Sort {
		sortGroupedPackages(config)
	}

	// Save updated config
	if err := yamlPkg.SaveGroupedConfig(config, filePath); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	utils.PrintStatus(utils.Green, fmt.Sprintf("Successfully synchronized %d new packages", len(missingPackages)))
	return nil
}

// SyncSimplePackages synchronizes installed packages with simple YAML config
func SyncSimplePackages(filePath string, options *types.SyncOptions) error {
	if options.Backup {
		if err := utils.CreateBackup(filePath); err != nil {
			return fmt.Errorf("failed to create backup: %w", err)
		}
	}

	// Load existing config
	config, err := yamlPkg.LoadSimpleConfig(filePath)
	if err != nil {
		return fmt.Errorf("failed to load simple config: %w", err)
	}

	// Get currently installed packages
	installedPackages, err := yamlPkg.GetInstalledPackages()
	if err != nil {
		return fmt.Errorf("failed to get installed packages: %w", err)
	}

	// Find and add missing packages
	var totalAdded int

	// Add missing taps
	for _, tap := range installedPackages.Taps {
		if !utils.ContainsString(config.Taps, tap) {
			if options.DryRun {
				utils.PrintStatus(utils.Yellow, fmt.Sprintf("[DRY RUN] Would add tap: %s", tap))
			} else {
				config.Taps = append(config.Taps, tap)
				utils.PrintStatus(utils.Green, fmt.Sprintf("Added tap: %s", tap))
			}
			totalAdded++
		}
	}

	// Add missing brews
	for _, brew := range installedPackages.Brews {
		if !utils.ContainsString(config.Brews, brew) {
			if options.DryRun {
				utils.PrintStatus(utils.Yellow, fmt.Sprintf("[DRY RUN] Would add brew: %s", brew))
			} else {
				config.Brews = append(config.Brews, brew)
				utils.PrintStatus(utils.Green, fmt.Sprintf("Added brew: %s", brew))
			}
			totalAdded++
		}
	}

	// Add missing casks
	for _, cask := range installedPackages.Casks {
		if !utils.ContainsString(config.Casks, cask) {
			if options.DryRun {
				utils.PrintStatus(utils.Yellow, fmt.Sprintf("[DRY RUN] Would add cask: %s", cask))
			} else {
				config.Casks = append(config.Casks, cask)
				utils.PrintStatus(utils.Green, fmt.Sprintf("Added cask: %s", cask))
			}
			totalAdded++
		}
	}

	// Add missing mas apps
	for _, app := range installedPackages.MasApps {
		found := false
		for _, existing := range config.MasApps {
			if existing.ID == app.ID {
				found = true
				break
			}
		}
		if !found {
			if options.DryRun {
				utils.PrintStatus(utils.Yellow, fmt.Sprintf("[DRY RUN] Would add mas app: %s", app.Name))
			} else {
				config.MasApps = append(config.MasApps, app)
				utils.PrintStatus(utils.Green, fmt.Sprintf("Added mas app: %s", app.Name))
			}
			totalAdded++
		}
	}

	if totalAdded == 0 {
		utils.PrintStatus(utils.Green, "No new packages found. Configuration is up to date.")
		return nil
	}

	if options.DryRun {
		return nil
	}

	// Sort if requested
	if options.Sort {
		sort.Strings(config.Taps)
		sort.Strings(config.Brews)
		sort.Strings(config.Casks)
		sort.Slice(config.MasApps, func(i, j int) bool {
			return config.MasApps[i].Name < config.MasApps[j].Name
		})
	}

	// Save updated config
	if err := yamlPkg.SaveSimpleConfig(config, filePath); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	utils.PrintStatus(utils.Green, fmt.Sprintf("Successfully synchronized %d new packages", totalAdded))
	return nil
}

// MissingPackage represents a package that is installed but not in config
type MissingPackage struct {
	Name string
	Type string
	ID   int64 // For mas apps
}

// findMissingPackages finds packages that are installed but not in the config
func findMissingPackages(config *types.PackageGrouped, installed *types.PackageSimple) []MissingPackage {
	var missing []MissingPackage

	// Get all packages from config
	configPackages := make(map[string]bool)
	for _, group := range config.Groups {
		for _, pkg := range group.Packages {
			key := fmt.Sprintf("%s:%s", pkg.Type, pkg.Name)
			configPackages[key] = true
		}
	}

	// Check taps
	for _, tap := range installed.Taps {
		key := fmt.Sprintf("tap:%s", tap)
		if !configPackages[key] {
			missing = append(missing, MissingPackage{Name: tap, Type: "tap"})
		}
	}

	// Check brews
	for _, brew := range installed.Brews {
		key := fmt.Sprintf("brew:%s", brew)
		if !configPackages[key] {
			missing = append(missing, MissingPackage{Name: brew, Type: "brew"})
		}
	}

	// Check casks
	for _, cask := range installed.Casks {
		key := fmt.Sprintf("cask:%s", cask)
		if !configPackages[key] {
			missing = append(missing, MissingPackage{Name: cask, Type: "cask"})
		}
	}

	// Check mas apps
	for _, app := range installed.MasApps {
		found := false
		for _, group := range config.Groups {
			for _, pkg := range group.Packages {
				if pkg.Type == "mas" && pkg.ID == app.ID {
					found = true
					break
				}
			}
			if found {
				break
			}
		}
		if !found {
			missing = append(missing, MissingPackage{Name: app.Name, Type: "mas", ID: app.ID})
		}
	}

	return missing
}

// showMissingPackages displays missing packages grouped by type
func showMissingPackages(packages []MissingPackage) {
	packagesByType := make(map[string][]MissingPackage)
	for _, pkg := range packages {
		packagesByType[pkg.Type] = append(packagesByType[pkg.Type], pkg)
	}

	for pkgType, pkgs := range packagesByType {
		utils.PrintStatus(utils.Cyan, fmt.Sprintf("%s packages:", strings.Title(pkgType)))
		for _, pkg := range pkgs {
			if pkg.Type == "mas" {
				fmt.Printf("  - %s (ID: %d)\n", pkg.Name, pkg.ID)
			} else {
				fmt.Printf("  - %s\n", pkg.Name)
			}
		}
	}
}

// addMissingPackagesToGrouped adds missing packages to grouped config
func addMissingPackagesToGrouped(config *types.PackageGrouped, missing []MissingPackage, options *types.SyncOptions) error {
	defaultGroup := options.DefaultGroup
	if defaultGroup == "" {
		defaultGroup = "optional"
	}

	// Ensure default group exists
	if _, exists := config.Groups[defaultGroup]; !exists {
		config.Groups[defaultGroup] = types.Group{
			Description: "Uncategorized packages",
			Priority:    10,
			Packages:    []types.Package{},
		}
	}

	for _, pkg := range missing {
		targetGroup := defaultGroup
		tags := options.DefaultTags

		// Auto-detect or interactive assignment
		if options.AutoDetect {
			targetGroup = utils.AutoDetectGroup(pkg.Name, pkg.Type)
			detectedTags := utils.AutoDetectTags(pkg.Name, pkg.Type)
			tags = append(tags, detectedTags...)
		}

		if options.Interactive {
			response, err := promptForPackageAssignment(pkg, targetGroup, tags)
			if err != nil {
				return fmt.Errorf("interactive prompt failed: %w", err)
			}
			targetGroup = response.Group
			tags = response.Tags
		}

		// Ensure target group exists
		if _, exists := config.Groups[targetGroup]; !exists {
			config.Groups[targetGroup] = types.Group{
				Description: fmt.Sprintf("Auto-created group for %s", targetGroup),
				Priority:    5,
				Packages:    []types.Package{},
			}
		}

		// Create package
		newPackage := types.Package{
			Name: pkg.Name,
			Type: pkg.Type,
			Tags: tags,
		}

		if pkg.Type == "mas" {
			newPackage.ID = pkg.ID
		}

		// Add to group
		group := config.Groups[targetGroup]
		group.Packages = append(group.Packages, newPackage)
		config.Groups[targetGroup] = group

		utils.PrintStatus(utils.Green, fmt.Sprintf("Added %s '%s' to group '%s'", pkg.Type, pkg.Name, targetGroup))
	}

	return nil
}

// PackageAssignment represents user input for package assignment
type PackageAssignment struct {
	Group string
	Tags  []string
}

// promptForPackageAssignment prompts user for group and tag assignment
func promptForPackageAssignment(pkg MissingPackage, suggestedGroup string, suggestedTags []string) (*PackageAssignment, error) {
	utils.PrintStatus(utils.Cyan, fmt.Sprintf("Package: %s (%s)", pkg.Name, pkg.Type))

	var group string
	prompt := &survey.Input{
		Message: "Enter group:",
		Default: suggestedGroup,
	}
	if err := survey.AskOne(prompt, &group); err != nil {
		return nil, err
	}

	var tagsString string
	tagsPrompt := &survey.Input{
		Message: "Enter tags (comma-separated):",
		Default: strings.Join(suggestedTags, ","),
	}
	if err := survey.AskOne(tagsPrompt, &tagsString); err != nil {
		return nil, err
	}

	tags := utils.SplitCommaSeparated(tagsString)

	return &PackageAssignment{
		Group: group,
		Tags:  tags,
	}, nil
}

// sortGroupedPackages sorts packages within each group
func sortGroupedPackages(config *types.PackageGrouped) {
	for groupName, group := range config.Groups {
		sort.Slice(group.Packages, func(i, j int) bool {
			// Sort by type first, then by name
			if group.Packages[i].Type != group.Packages[j].Type {
				return group.Packages[i].Type < group.Packages[j].Type
			}
			return group.Packages[i].Name < group.Packages[j].Name
		})
		config.Groups[groupName] = group
	}
} 
