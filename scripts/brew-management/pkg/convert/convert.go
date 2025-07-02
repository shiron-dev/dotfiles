package convert

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"brew-manager/pkg/types"
	"brew-manager/pkg/utils"
	yamlPkg "brew-manager/pkg/yaml"
)

// ConvertBrewfileToYAML converts a Brewfile to YAML format
func ConvertBrewfileToYAML(brewfilePath, yamlPath string, grouped bool, verbose bool) error {
	if !utils.FileExists(brewfilePath) {
		return fmt.Errorf("Brewfile not found: %s", brewfilePath)
	}

	if verbose {
		utils.PrintStatus(utils.Blue, fmt.Sprintf("Converting Brewfile: %s", brewfilePath))
	}

	// Parse Brewfile
	brewfileData, err := parseBrewfile(brewfilePath, verbose)
	if err != nil {
		return fmt.Errorf("failed to parse Brewfile: %w", err)
	}

	// Convert to grouped format (only supported format now)
	groupedConfig := convertToGroupedFormat(brewfileData, verbose)
	if err := yamlPkg.SaveGroupedConfig(groupedConfig, yamlPath); err != nil {
		return fmt.Errorf("failed to save grouped YAML: %w", err)
	}

	utils.PrintStatus(utils.Green, fmt.Sprintf("Successfully converted Brewfile to: %s", yamlPath))
	return nil
}

// BrewfileData represents parsed Brewfile content
type BrewfileData struct {
	Taps    []string
	Brews   []string
	Casks   []string
	MasApps []types.MasApp
}

// parseBrewfile parses a Brewfile and extracts package information
func parseBrewfile(filePath string, verbose bool) (*BrewfileData, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open Brewfile: %w", err)
	}
	defer file.Close()

	data := &BrewfileData{
		Taps:    []string{},
		Brews:   []string{},
		Casks:   []string{},
		MasApps: []types.MasApp{},
	}

	scanner := bufio.NewScanner(file)
	lineNum := 0

	// Regex patterns for different types
	tapRegex := regexp.MustCompile(`^tap\s+["']?([^"'\s]+)["']?`)
	brewRegex := regexp.MustCompile(`^brew\s+["']?([^"'\s]+)["']?`)
	caskRegex := regexp.MustCompile(`^cask\s+["']?([^"'\s]+)["']?`)
	masRegex := regexp.MustCompile(`^mas\s+["']?([^"']+?)["']?\s*,\s*id:\s*(\d+)`)

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if verbose {
			utils.PrintStatus(utils.Cyan, fmt.Sprintf("Processing line %d: %s", lineNum, line))
		}

		// Parse tap
		if matches := tapRegex.FindStringSubmatch(line); len(matches) > 1 {
			data.Taps = append(data.Taps, matches[1])
			continue
		}

		// Parse brew
		if matches := brewRegex.FindStringSubmatch(line); len(matches) > 1 {
			data.Brews = append(data.Brews, matches[1])
			continue
		}

		// Parse cask
		if matches := caskRegex.FindStringSubmatch(line); len(matches) > 1 {
			data.Casks = append(data.Casks, matches[1])
			continue
		}

		// Parse mas
		if matches := masRegex.FindStringSubmatch(line); len(matches) > 2 {
			id, err := strconv.ParseInt(matches[2], 10, 64)
			if err != nil {
				utils.PrintStatus(utils.Yellow, fmt.Sprintf("Warning: Invalid mas ID on line %d: %s", lineNum, matches[2]))
				continue
			}
			data.MasApps = append(data.MasApps, types.MasApp{
				Name: matches[1],
				ID:   id,
			})
			continue
		}

		// Unknown line format
		if verbose {
			utils.PrintStatus(utils.Yellow, fmt.Sprintf("Warning: Unrecognized line format on line %d: %s", lineNum, line))
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading Brewfile: %w", err)
	}

	if verbose {
		utils.PrintStatus(utils.Green, fmt.Sprintf("Parsed %d taps, %d brews, %d casks, %d mas apps", 
			len(data.Taps), len(data.Brews), len(data.Casks), len(data.MasApps)))
	}

	return data, nil
}



// convertToGroupedFormat converts BrewfileData to grouped YAML format
func convertToGroupedFormat(data *BrewfileData, verbose bool) *types.PackageGrouped {
	if verbose {
		utils.PrintStatus(utils.Blue, "Converting to grouped YAML format with auto-detection")
	}

	config := &types.PackageGrouped{
		Groups:   make(map[string]types.Group),
		Profiles: make(map[string]types.Profile),
	}

	// Initialize groups with the new Packages structure
	initializedGroups := make(map[string]types.Group)
	groupDefinitions := map[string]struct {
		description string
		priority    int
	}{
		"core":         {"Essential development tools", 1},
		"development":  {"Development tools and environments", 2},
		"productivity": {"Productivity and office applications", 3},
		"creative":     {"Creative and multimedia tools", 4},
		"system":       {"System utilities and tools", 5},
		"optional":     {"Optional and uncategorized tools", 10},
	}

	for name, def := range groupDefinitions {
		initializedGroups[name] = types.Group{
			Description: def.description,
			Priority:    def.priority,
			Packages:    make(map[string][]types.PackageInfo), // Initialize the map
		}
	}

	// Add taps
	for _, tapName := range data.Taps {
		groupName := utils.AutoDetectGroup(tapName, "tap")
		tags := utils.AutoDetectTags(tapName, "tap")
		
		pkgInfo := types.PackageInfo{
			Name: tapName,
			Tags: tags,
		}
		
		group := initializedGroups[groupName] // Get the group
		group.Packages["tap"] = append(group.Packages["tap"], pkgInfo) // Append to the 'tap' type slice
		initializedGroups[groupName] = group // Update the map with the modified group
		
		if verbose {
			utils.PrintStatus(utils.Cyan, fmt.Sprintf("Added tap '%s' to group '%s' with tags: %v", tapName, groupName, tags))
		}
	}

	// Add brews
	for _, brewName := range data.Brews {
		groupName := utils.AutoDetectGroup(brewName, "brew")
		tags := utils.AutoDetectTags(brewName, "brew")
		
		pkgInfo := types.PackageInfo{
			Name: brewName,
			Tags: tags,
		}
		
		group := initializedGroups[groupName]
		group.Packages["brew"] = append(group.Packages["brew"], pkgInfo)
		initializedGroups[groupName] = group
		
		if verbose {
			utils.PrintStatus(utils.Cyan, fmt.Sprintf("Added brew '%s' to group '%s' with tags: %v", brewName, groupName, tags))
		}
	}

	// Add casks
	for _, caskName := range data.Casks {
		groupName := utils.AutoDetectGroup(caskName, "cask")
		tags := utils.AutoDetectTags(caskName, "cask")
		
		pkgInfo := types.PackageInfo{
			Name: caskName,
			Tags: tags,
		}
		
		group := initializedGroups[groupName]
		group.Packages["cask"] = append(group.Packages["cask"], pkgInfo)
		initializedGroups[groupName] = group
		
		if verbose {
			utils.PrintStatus(utils.Cyan, fmt.Sprintf("Added cask '%s' to group '%s' with tags: %v", caskName, groupName, tags))
		}
	}

	// Add mas apps
	for _, app := range data.MasApps {
		groupName := utils.AutoDetectGroup(app.Name, "mas")
		tags := utils.AutoDetectTags(app.Name, "mas")
		
		pkgInfo := types.PackageInfo{
			Name: app.Name,
			Tags: tags,
			ID:   app.ID,
		}
		
		group := initializedGroups[groupName]
		group.Packages["mas"] = append(group.Packages["mas"], pkgInfo)
		initializedGroups[groupName] = group
		
		if verbose {
			utils.PrintStatus(utils.Cyan, fmt.Sprintf("Added mas app '%s' to group '%s' with tags: %v", app.Name, groupName, tags))
		}
	}

	config.Groups = initializedGroups // Assign the fully populated groups to the config

	// Add default profiles
	config.Profiles = map[string]types.Profile{
		"minimal": {
			Description: "Minimal development setup",
			Groups:      []string{"core"},
			Tags:        []string{"essential"},
		},
		"developer": {
			Description: "Full development environment",
			Groups:      []string{"core", "development"},
			ExcludeTags: []string{"experimental"},
		},
		"full": {
			Description: "Complete setup with all tools",
			Groups:      []string{"core", "development", "productivity", "creative", "system"},
		},
	}

	return config
}

// ValidateBrewfile validates the syntax of a Brewfile
func ValidateBrewfile(filePath string, verbose bool) error {
	if !utils.FileExists(filePath) {
		return fmt.Errorf("Brewfile not found: %s", filePath)
	}

	if verbose {
		utils.PrintStatus(utils.Blue, fmt.Sprintf("Validating Brewfile: %s", filePath))
	}

	_, err := parseBrewfile(filePath, verbose)
	if err != nil {
		utils.PrintStatus(utils.Red, fmt.Sprintf("❌ Invalid Brewfile: %v", err))
		return err
	}

	utils.PrintStatus(utils.Green, "✅ Brewfile is valid")
	return nil
} 
