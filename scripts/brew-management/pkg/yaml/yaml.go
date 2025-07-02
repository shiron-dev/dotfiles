package yaml

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"brew-manager/pkg/types"
	"brew-manager/pkg/utils"

	"gopkg.in/yaml.v3"
)

// LoadGroupedConfig loads a grouped YAML configuration file
func LoadGroupedConfig(filePath string) (*types.PackageGrouped, error) {
	// If file doesn't exist, return a default empty configuration
	if !utils.FileExists(filePath) {
		return createDefaultGroupedConfig(), nil
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read YAML file: %w", err)
	}

	// If file is empty or contains only whitespace/comments, return default config
	content := string(data)
	lines := strings.Split(content, "\n")
	var filteredLines []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		// Skip comments and empty lines
		if !strings.HasPrefix(trimmed, "#") && trimmed != "" {
			filteredLines = append(filteredLines, line)
		}
	}

	// If no meaningful content found, return default config
	if len(filteredLines) == 0 || strings.TrimSpace(strings.Join(filteredLines, "\n")) == "" {
		return createDefaultGroupedConfig(), nil
	}

	// Remove yaml-language-server comment if present for parsing
	var cleanLines []string
	for _, line := range lines {
		if !strings.HasPrefix(strings.TrimSpace(line), "# yaml-language-server:") {
			cleanLines = append(cleanLines, line)
		}
	}
	cleanContent := strings.Join(cleanLines, "\n")

	var config types.PackageGrouped
	if err := yaml.Unmarshal([]byte(cleanContent), &config); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	// Ensure Groups is initialized
	if config.Groups == nil {
		config.Groups = make(map[string]types.Group)
	}

	// Ensure Profiles is initialized
	if config.Profiles == nil {
		config.Profiles = make(map[string]types.Profile)
	}

	return &config, nil
}

// createDefaultGroupedConfig creates a default empty configuration
func createDefaultGroupedConfig() *types.PackageGrouped {
	return &types.PackageGrouped{
		Groups:   make(map[string]types.Group),
		Profiles: make(map[string]types.Profile),
	}
}

// SaveGroupedConfig saves a grouped configuration to YAML file
func SaveGroupedConfig(config *types.PackageGrouped, filePath string) error {
	// Create directory if it doesn't exist
	if err := utils.EnsureDir(filePath); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Sort groups by priority
	type groupEntry struct {
		name  string
		group types.Group
	}
	var sortedGroups []groupEntry
	for name, group := range config.Groups {
		// Sort packages within each group by type and then by name
		for pkgType, pkgInfos := range group.Packages {
			sort.Slice(pkgInfos, func(i, j int) bool {
				return pkgInfos[i].Name < pkgInfos[j].Name
			})
			group.Packages[pkgType] = pkgInfos
		}
		// Sort package types (keys of the map)
		// This part is tricky as map keys are not ordered in Go serialization directly.
		// For consistent YAML output, we might need a custom Marshal or a structure that preserves order.
		// However, yaml.v3 generally tries to sort map keys alphabetically by default.
		sortedGroups = append(sortedGroups, groupEntry{name, group})
	}
	sort.Slice(sortedGroups, func(i, j int) bool {
		return sortedGroups[i].group.Priority < sortedGroups[j].group.Priority
	})

	// Recreate the groups map in sorted order for top-level groups
	orderedGroups := make(map[string]types.Group)
	for _, entry := range sortedGroups {
		orderedGroups[entry.name] = entry.group
	}
	config.Groups = orderedGroups

	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal YAML: %w", err)
	}

	// Add yaml-language-server comment
	// TODO: Update the schema path if it's hosted or standardized
	content := "# yaml-language-server: $schema=~/github.com/shiron-dev/dotfiles/scripts/brew-management/packages.schema.json\n\n"
	content += string(data)

	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write YAML file: %w", err)
	}

	return nil
}

// GetFilteredPackages returns packages filtered by groups, tags, and exclusions
func GetFilteredPackages(config *types.PackageGrouped, options *types.InstallOptions) []types.FilteredPackage {
	var allPackages []types.FilteredPackage

	// Apply profile first if specified
	if options.Profile != "" {
		profile, exists := config.Profiles[options.Profile]
		if exists {
			// Deduplicate groups and tags
			options.Groups = utils.UniqueStrings(append(options.Groups, profile.Groups...))
			options.Tags = utils.UniqueStrings(append(options.Tags, profile.Tags...))
			// TODO: Handle profile.ExcludeTags
		}
	}

	groupsToProcess := options.Groups
	if len(groupsToProcess) == 0 {
		for groupName := range config.Groups {
			groupsToProcess = append(groupsToProcess, groupName)
		}
	}
	groupsToProcess = utils.UniqueStrings(groupsToProcess) // Ensure unique group processing

	for _, groupName := range groupsToProcess {
		group, exists := config.Groups[groupName]
		if !exists {
			continue
		}

		for pkgType, pkgInfos := range group.Packages {
			for _, pkgInfo := range pkgInfos {
				// Apply tag filters
				if len(options.Tags) > 0 && !utils.HasIntersection(pkgInfo.Tags, options.Tags) {
					continue
				}
				// TODO: Apply options.ExcludeTags if implemented

				allPackages = append(allPackages, types.FilteredPackage{
					PackageInfo: pkgInfo,
					Type:        pkgType,
				})
			}
		}
	}

	// Sort final list for consistent processing order, by Type then by Name
	sort.Slice(allPackages, func(i, j int) bool {
		if allPackages[i].Type != allPackages[j].Type {
			return allPackages[i].Type < allPackages[j].Type
		}
		return allPackages[i].Name < allPackages[j].Name
	})

	return allPackages
}

// GetInstalledPackages retrieves currently installed brew packages
func GetInstalledPackages() (map[string][]string, []types.MasApp, error) {
	result := make(map[string][]string)
	var masApps []types.MasApp

	// Get taps
	if tapsOutput, err := utils.RunCommand("brew", "tap"); err == nil {
		result["taps"] = strings.Fields(strings.TrimSpace(tapsOutput))
	}

	// Get formulae
	if brewsOutput, err := utils.RunCommand("brew", "list", "--formula"); err == nil {
		result["brews"] = strings.Fields(strings.TrimSpace(brewsOutput))
	}

	// Get casks
	if casksOutput, err := utils.RunCommand("brew", "list", "--cask"); err == nil {
		result["casks"] = strings.Fields(strings.TrimSpace(casksOutput))
	}

	// Get mas apps if mas is available
	if utils.CommandExists("mas") {
		if masOutput, err := utils.RunCommand("mas", "list"); err == nil {
			lines := strings.Split(strings.TrimSpace(masOutput), "\n")
			for _, line := range lines {
				if line == "" {
					continue
				}
				parts := strings.SplitN(line, " ", 2)
				if len(parts) >= 2 {
					var id int64
					if _, err := fmt.Sscanf(parts[0], "%d", &id); err == nil {
						masApps = append(masApps, types.MasApp{
							Name: strings.TrimSpace(parts[1]),
							ID:   id,
						})
					}
				}
			}
		}
	}

	return result, masApps, nil
}
