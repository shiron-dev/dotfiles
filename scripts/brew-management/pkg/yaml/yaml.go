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
	if !utils.FileExists(filePath) {
		return nil, fmt.Errorf("YAML file not found: %s", filePath)
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read YAML file: %w", err)
	}

	// Remove yaml-language-server comment if present
	content := string(data)
	lines := strings.Split(content, "\n")
	var filteredLines []string
	for _, line := range lines {
		if !strings.HasPrefix(strings.TrimSpace(line), "# yaml-language-server:") {
			filteredLines = append(filteredLines, line)
		}
	}
	cleanContent := strings.Join(filteredLines, "\n")

	var config types.PackageGrouped
	if err := yaml.Unmarshal([]byte(cleanContent), &config); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	return &config, nil
}

// LoadSimpleConfig loads a simple YAML configuration file
func LoadSimpleConfig(filePath string) (*types.PackageSimple, error) {
	if !utils.FileExists(filePath) {
		return nil, fmt.Errorf("YAML file not found: %s", filePath)
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read YAML file: %w", err)
	}

	// Remove yaml-language-server comment if present
	content := string(data)
	lines := strings.Split(content, "\n")
	var filteredLines []string
	for _, line := range lines {
		if !strings.HasPrefix(strings.TrimSpace(line), "# yaml-language-server:") {
			filteredLines = append(filteredLines, line)
		}
	}
	cleanContent := strings.Join(filteredLines, "\n")

	var config types.PackageSimple
	if err := yaml.Unmarshal([]byte(cleanContent), &config); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	return &config, nil
}

// SaveGroupedConfig saves a grouped configuration to YAML file
func SaveGroupedConfig(config *types.PackageGrouped, filePath string) error {
	// Sort groups by priority
	type groupEntry struct {
		name  string
		group types.Group
	}
	var sortedGroups []groupEntry
	for name, group := range config.Groups {
		sortedGroups = append(sortedGroups, groupEntry{name, group})
	}
	sort.Slice(sortedGroups, func(i, j int) bool {
		return sortedGroups[i].group.Priority < sortedGroups[j].group.Priority
	})

	// Recreate the groups map in sorted order
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
	content := "# yaml-language-server: $schema=schemas/packages-grouped.schema.json\n"
	content += "# YAML-based Brew packages configuration with groups and tags\n"
	content += "# This file defines packages with group and tag classification\n\n"
	content += string(data)

	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write YAML file: %w", err)
	}

	return nil
}

// SaveSimpleConfig saves a simple configuration to YAML file
func SaveSimpleConfig(config *types.PackageSimple, filePath string) error {
	// Sort all arrays
	sort.Strings(config.Taps)
	sort.Strings(config.Brews)
	sort.Strings(config.Casks)
	sort.Slice(config.MasApps, func(i, j int) bool {
		return config.MasApps[i].Name < config.MasApps[j].Name
	})

	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal YAML: %w", err)
	}

	// Add yaml-language-server comment
	content := "# yaml-language-server: $schema=schemas/packages-simple.schema.json\n"
	content += "# YAML-based Brew packages configuration\n"
	content += "# This file defines packages to be installed via Homebrew\n\n"
	content += string(data)

	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write YAML file: %w", err)
	}

	return nil
}

// GetFilteredPackages returns packages filtered by groups, tags, and exclusions
func GetFilteredPackages(config *types.PackageGrouped, options *types.InstallOptions) []types.Package {
	var allPackages []types.Package

	// Apply profile first if specified
	if options.Profile != "" {
		profile, exists := config.Profiles[options.Profile]
		if exists {
			options.Groups = append(options.Groups, profile.Groups...)
			options.Tags = append(options.Tags, profile.Tags...)
			options.ExcludeTags = append(options.ExcludeTags, profile.ExcludeTags...)
		}
	}

	// Collect packages from specified groups or all groups
	groupsToProcess := options.Groups
	if len(groupsToProcess) == 0 {
		for groupName := range config.Groups {
			groupsToProcess = append(groupsToProcess, groupName)
		}
	}

	for _, groupName := range groupsToProcess {
		if utils.ContainsString(options.ExcludeGroups, groupName) {
			continue
		}
		
		group, exists := config.Groups[groupName]
		if !exists {
			continue
		}

		for _, pkg := range group.Packages {
			// Apply tag filters
			if len(options.Tags) > 0 && !utils.HasIntersection(pkg.Tags, options.Tags) {
				continue
			}
			if len(options.ExcludeTags) > 0 && utils.HasIntersection(pkg.Tags, options.ExcludeTags) {
				continue
			}

			allPackages = append(allPackages, pkg)
		}
	}

	return allPackages
}

// GetInstalledPackages retrieves currently installed brew packages
func GetInstalledPackages() (*types.PackageSimple, error) {
	result := &types.PackageSimple{}

	// Get taps
	if tapsOutput, err := utils.RunCommand("brew", "tap"); err == nil {
		result.Taps = strings.Fields(strings.TrimSpace(tapsOutput))
	}

	// Get formulae
	if brewsOutput, err := utils.RunCommand("brew", "list", "--formula"); err == nil {
		result.Brews = strings.Fields(strings.TrimSpace(brewsOutput))
	}

	// Get casks
	if casksOutput, err := utils.RunCommand("brew", "list", "--cask"); err == nil {
		result.Casks = strings.Fields(strings.TrimSpace(casksOutput))
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
						result.MasApps = append(result.MasApps, types.MasApp{
							Name: strings.TrimSpace(parts[1]),
							ID:   id,
						})
					}
				}
			}
		}
	}

	return result, nil
} 
