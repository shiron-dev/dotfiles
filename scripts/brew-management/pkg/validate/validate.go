package validate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"brew-manager/pkg/types"
	"brew-manager/pkg/utils"
	yamlPkg "brew-manager/pkg/yaml"

	"gopkg.in/yaml.v3"
)

// ValidateYAMLFile validates a YAML file against its schema
func ValidateYAMLFile(filePath string, options *types.ValidateOptions) error {
	if !utils.FileExists(filePath) {
		return fmt.Errorf("YAML file not found: %s", filePath)
	}

	if options.Verbose {
		utils.PrintStatus(utils.Blue, fmt.Sprintf("Validating: %s", filePath))
	}

	// Read and clean YAML content
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read YAML file: %w", err)
	}

	// Remove yaml-language-server comment
	content := string(data)
	lines := strings.Split(content, "\n")
	var filteredLines []string
	for _, line := range lines {
		if !strings.HasPrefix(strings.TrimSpace(line), "# yaml-language-server:") {
			filteredLines = append(filteredLines, line)
		}
	}
	cleanContent := strings.Join(filteredLines, "\n")

	// Basic YAML syntax check
	var tempYAML interface{}
	if err := yaml.Unmarshal([]byte(cleanContent), &tempYAML); err != nil {
		utils.PrintStatus(utils.Red, fmt.Sprintf("❌ Invalid: %s - YAML syntax error", filepath.Base(filePath)))
		if options.Verbose {
			utils.PrintStatus(utils.Yellow, fmt.Sprintf("YAML syntax errors: %v", err))
		}
		return err
	}

	// Determine file type and validate structure
	filename := filepath.Base(filePath)
	var validationErrors []string

	switch {
	case strings.Contains(filename, "grouped"):
		validationErrors = validateGroupedYAML(cleanContent, options.Verbose)
	case strings.Contains(filename, "packages"):
		validationErrors = append(validationErrors, "Simple YAML format is no longer supported")
	default:
		// Try to detect format by content
		if strings.Contains(cleanContent, "groups:") && strings.Contains(cleanContent, "metadata:") {
			validationErrors = validateGroupedYAML(cleanContent, options.Verbose)
		} else {
			validationErrors = append(validationErrors, "Simple YAML format is no longer supported")
		}
	}

	// Report validation results
	if len(validationErrors) == 0 {
		utils.PrintStatus(utils.Green, fmt.Sprintf("✅ Valid: %s", filepath.Base(filePath)))
		return nil
	} else {
		utils.PrintStatus(utils.Red, fmt.Sprintf("❌ Invalid: %s", filepath.Base(filePath)))
		if options.Verbose {
			utils.PrintStatus(utils.Yellow, "Validation errors:")
			for _, error := range validationErrors {
				fmt.Printf("  - %s\n", error)
			}
		}
		return fmt.Errorf("validation failed with %d errors", len(validationErrors))
	}
}

// validateGroupedYAML validates grouped YAML format
func validateGroupedYAML(content string, verbose bool) []string {
	var errors []string

	// Try to parse as grouped config
	var config types.PackageGrouped
	if err := yaml.Unmarshal([]byte(content), &config); err != nil {
		errors = append(errors, fmt.Sprintf("Failed to parse as grouped format: %v", err))
		return errors
	}

	// Check required fields
	if config.Metadata.Version == "" {
		errors = append(errors, "Missing metadata.version")
	}

	if len(config.Groups) == 0 {
		errors = append(errors, "Missing groups section")
	}

	// Validate groups
	for groupName, group := range config.Groups {
		if group.Description == "" {
			errors = append(errors, fmt.Sprintf("Missing description in group: %s", groupName))
		}
		if group.Priority == 0 {
			errors = append(errors, fmt.Sprintf("Missing or zero priority in group: %s", groupName))
		}

		// Validate packages in group
		for i, pkg := range group.Packages {
			if pkg.Name == "" {
				errors = append(errors, fmt.Sprintf("Missing name in group %s, package %d", groupName, i))
			}
			if pkg.Type == "" {
				errors = append(errors, fmt.Sprintf("Missing type in group %s, package %s", groupName, pkg.Name))
			}
			if pkg.Type == "mas" && pkg.ID == 0 {
				errors = append(errors, fmt.Sprintf("Missing ID for mas app in group %s, package %s", groupName, pkg.Name))
			}
		}
	}

	// Validate profiles if present
	for profileName, profile := range config.Profiles {
		if profile.Description == "" {
			errors = append(errors, fmt.Sprintf("Missing description in profile: %s", profileName))
		}
	}

	return errors
}



// ValidateAllYAMLFiles validates all YAML files in the data directory
func ValidateAllYAMLFiles(dataDir string, options *types.ValidateOptions) error {
	if !utils.FileExists(dataDir) {
		return fmt.Errorf("data directory not found: %s", dataDir)
	}

	utils.PrintStatus(utils.Blue, fmt.Sprintf("Validating all YAML files in: %s", dataDir))

	var hasErrors bool

	// Walk through directory and find YAML files
	err := filepath.Walk(dataDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if file is a YAML file
		if !info.IsDir() && (strings.HasSuffix(path, ".yml") || strings.HasSuffix(path, ".yaml")) {
			// Skip schema files
			if strings.Contains(path, "schema") {
				return nil
			}

			if err := ValidateYAMLFile(path, options); err != nil {
				hasErrors = true
			}
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to walk directory: %w", err)
	}

	if hasErrors {
		return fmt.Errorf("validation failed for one or more files")
	}

	utils.PrintStatus(utils.Green, "All YAML files are valid!")
	return nil
}

// TestYAMLLoad tests loading YAML files
func TestYAMLLoad(filePath string, verbose bool) error {
	if verbose {
		utils.PrintStatus(utils.Blue, fmt.Sprintf("Testing YAML load: %s", filePath))
	}

	filename := filepath.Base(filePath)
	
	if strings.Contains(filename, "grouped") {
		_, err := yamlPkg.LoadGroupedConfig(filePath)
		if err != nil {
			utils.PrintStatus(utils.Red, fmt.Sprintf("Failed to load grouped config: %v", err))
			return err
		}
		utils.PrintStatus(utils.Green, "Successfully loaded grouped config")
	} else {
		utils.PrintStatus(utils.Red, "Simple YAML format is no longer supported")
	}

	return nil
} 
