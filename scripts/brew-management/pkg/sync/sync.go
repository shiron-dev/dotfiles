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
	// Check if file exists before backup
	fileExists := utils.FileExists(filePath)
	
	if options.Backup && fileExists {
		if err := utils.CreateBackup(filePath); err != nil {
			return fmt.Errorf("failed to create backup: %w", err)
		}
	}

	// Load existing config
	config, err := yamlPkg.LoadGroupedConfig(filePath)
	if err != nil {
		return fmt.Errorf("failed to load grouped config: %w", err)
	}
	
	// Notify if we're starting with an empty configuration
	if !fileExists || len(config.Groups) == 0 {
		utils.PrintStatus(utils.Yellow, "Starting with empty configuration. All installed packages will be added.")
	}

	// Get currently installed packages
	installedPackagesMap, installedMasApps, err := yamlPkg.GetInstalledPackages()
	if err != nil {
		return fmt.Errorf("failed to get installed packages: %w", err)
	}

	// Find missing packages
	missingPackages := findMissingPackages(config, installedPackagesMap, installedMasApps)
	
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



// MissingPackage represents a package that is installed but not in config
type MissingPackage struct {
	Name string
	Type string
	ID   int64 // For mas apps
}

// findMissingPackages finds packages that are installed but not in the config
func findMissingPackages(config *types.PackageGrouped, installedPackagesMap map[string][]string, installedMasApps []types.MasApp) []MissingPackage {
	var missing []MissingPackage

	// Get all packages from config
	configPackages := make(map[string]bool) // Stores "type:name" or "mas:id"
	for _, group := range config.Groups {
		for pkgType, pkgInfos := range group.Packages {
			for _, pkgInfo := range pkgInfos {
				var key string
				if pkgType == "mas" {
					key = fmt.Sprintf("mas:%d", pkgInfo.ID)
				} else {
					key = fmt.Sprintf("%s:%s", pkgType, pkgInfo.Name)
				}
				configPackages[key] = true
			}
		}
	}

	// Check taps
	for _, tap := range installedPackagesMap["taps"] {
		key := fmt.Sprintf("tap:%s", tap)
		if !configPackages[key] {
			missing = append(missing, MissingPackage{Name: tap, Type: "tap"})
		}
	}

	// Check brews
	for _, brew := range installedPackagesMap["brews"] {
		key := fmt.Sprintf("brew:%s", brew)
		if !configPackages[key] {
			missing = append(missing, MissingPackage{Name: brew, Type: "brew"})
		}
	}

	// Check casks
	for _, cask := range installedPackagesMap["casks"] {
		key := fmt.Sprintf("cask:%s", cask)
		if !configPackages[key] {
			missing = append(missing, MissingPackage{Name: cask, Type: "cask"})
		}
	}

	// Check mas apps
	for _, app := range installedMasApps {
		key := fmt.Sprintf("mas:%d", app.ID)
		if !configPackages[key] {
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

	// Sort package types for consistent output order
	sortedTypes := make([]string, 0, len(packagesByType))
	for pkgType := range packagesByType {
		sortedTypes = append(sortedTypes, pkgType)
	}
	sort.Strings(sortedTypes)

	for _, pkgType := range sortedTypes {
		pkgs := packagesByType[pkgType]
		utils.PrintStatus(utils.Cyan, fmt.Sprintf("%s packages:", strings.Title(pkgType)))
		// Sort packages within each type by name for consistent output
		sort.Slice(pkgs, func(i, j int) bool {
			return pkgs[i].Name < pkgs[j].Name
		})
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
		defaultGroup = "uncategorized"
	}

	// Ensure default group exists and has Packages map initialized
	if _, exists := config.Groups[defaultGroup]; !exists {
		config.Groups[defaultGroup] = types.Group{
			Description: "Uncategorized packages",
			Priority:    10,
			Packages:    make(map[string][]types.PackageInfo),
		}
	} else {
		// Ensure Packages map is initialized even if group exists
		grp := config.Groups[defaultGroup]
		if grp.Packages == nil {
			grp.Packages = make(map[string][]types.PackageInfo)
			config.Groups[defaultGroup] = grp
		}
	}


	for _, pkg := range missing {
		targetGroup := defaultGroup
		tags := options.DefaultTags // Use default tags from options

		if options.Interactive {
			// Auto-detect group and tags for suggestion if AutoDetect is true
			suggestedGroup := targetGroup
			suggestedTags := tags
			if options.AutoDetect {
				suggestedGroup = utils.AutoDetectGroup(pkg.Name, pkg.Type)
				suggestedTags = utils.AutoDetectTags(pkg.Name, pkg.Type)
			}

			response, err := promptForPackageAssignment(pkg, suggestedGroup, suggestedTags)
			if err != nil {
				return fmt.Errorf("interactive prompt failed: %w", err)
			}
			targetGroup = response.Group
			tags = response.Tags
		}

		// Ensure target group exists and has Packages map initialized
		if _, exists := config.Groups[targetGroup]; !exists {
			config.Groups[targetGroup] = types.Group{
				Description: fmt.Sprintf("Auto-created group for %s", targetGroup),
				Priority:    5, // Default priority for auto-created groups
				Packages:    make(map[string][]types.PackageInfo),
			}
		} else {
			grp := config.Groups[targetGroup]
			if grp.Packages == nil {
				grp.Packages = make(map[string][]types.PackageInfo)
				config.Groups[targetGroup] = grp
			}
		}


		// Create PackageInfo
		newPackageInfo := types.PackageInfo{
			Name: pkg.Name,
			Tags: tags,
		}

		if pkg.Type == "mas" {
			newPackageInfo.ID = pkg.ID
		}

		// Add to group
		group := config.Groups[targetGroup]
		group.Packages[pkg.Type] = append(group.Packages[pkg.Type], newPackageInfo)
		config.Groups[targetGroup] = group // Update the map with the modified group

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
	if tagsString == "" { // Handle case where user enters empty string for no tags
		tags = []string{}
	}


	return &PackageAssignment{
		Group: group,
		Tags:  tags,
	}, nil
}

// sortGroupedPackages sorts packages within each group by type and then by name
func sortGroupedPackages(config *types.PackageGrouped) {
	for groupName, group := range config.Groups {
		for pkgType, pkgInfos := range group.Packages {
			sort.Slice(pkgInfos, func(i, j int) bool {
				return pkgInfos[i].Name < pkgInfos[j].Name
			})
			group.Packages[pkgType] = pkgInfos // Update the slice in the map
		}
		// Note: Map keys (package types) are not sorted here, but yaml.v3 usually sorts them.
		config.Groups[groupName] = group
	}
}
