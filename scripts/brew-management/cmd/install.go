package cmd

import (
	"fmt"
	"sort"

	"brew-manager/pkg/brew"
	"brew-manager/pkg/types"
	"brew-manager/pkg/utils"
	yamlPkg "brew-manager/pkg/yaml"

	"github.com/spf13/cobra"
)

var (
	groups        string
	tags          string
	profile       string
	skipTaps      bool
	skipBrews     bool
	skipCasks     bool
	skipMas       bool
	listGroups    bool
	listTags      bool
	listProfiles  bool
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install [yaml_file]",
	Short: "Install packages from YAML configuration (with groups/tags support)",
	Long: `Install Homebrew packages from a YAML configuration file with group/tag support.

Examples:
  brew-manager install                                    # Install all packages
  brew-manager install --groups core,development         # Install only core and development groups
  brew-manager install --tags essential,productivity     # Install packages with essential or productivity tags
  brew-manager install --profile developer               # Install using developer profile
  brew-manager install --exclude-tags experimental       # Install all except experimental packages
  brew-manager install --groups development --brews-only # Install only brew packages from development group`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get YAML file path
		yamlFile := getDefaultYAMLPath("packages.yaml")
		if len(args) > 0 {
			yamlFile = args[0]
		}

		// Handle list commands
		if listGroups || listTags || listProfiles {
			if err := handleListCommands(yamlFile); err != nil {
				utils.PrintStatus(utils.Red, fmt.Sprintf("Error: %v", err))
				return
			}
			return
		}

		// Build install options
		options := &types.InstallOptions{
			DryRun:         dryRun,
			Verbose:        verbose,
			Groups:         utils.SplitCommaSeparated(groups),
			Tags:           utils.SplitCommaSeparated(tags),
			Profile:        profile,
			SkipTaps:       skipTaps,
			SkipBrews:      skipBrews,
			SkipCasks:      skipCasks,
			SkipMas:        skipMas,
		}

		// Load configuration
		config, err := yamlPkg.LoadGroupedConfig(yamlFile)
		if err != nil {
			utils.PrintStatus(utils.Red, fmt.Sprintf("Error loading config: %v", err))
			return
		}

		// Get filtered packages
		packages := yamlPkg.GetFilteredPackages(config, options)

		if len(packages) == 0 {
			utils.PrintStatus(utils.Yellow, "No packages found matching the specified criteria.")
			return
		}

		utils.PrintStatus(utils.Blue, fmt.Sprintf("Found %d packages to install", len(packages)))

		// Install packages
		if err := brew.InstallPackages(packages, options); err != nil {
			utils.PrintStatus(utils.Red, fmt.Sprintf("Installation failed: %v", err))
			return
		}

		utils.PrintStatus(utils.Green, "Installation completed successfully!")
	},
}

func handleListCommands(yamlFile string) error {
	config, err := yamlPkg.LoadGroupedConfig(yamlFile)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if listGroups {
		utils.PrintStatus(utils.Cyan, "Available Groups:")
		
		// Create slice for sorting
		type groupInfo struct {
			name     string
			priority int
			desc     string
		}
		var groups []groupInfo
		
		for name, group := range config.Groups {
			groups = append(groups, groupInfo{
				name:     name,
				priority: group.Priority,
				desc:     group.Description,
			})
		}
		
		// Sort by priority
		sort.Slice(groups, func(i, j int) bool {
			return groups[i].priority < groups[j].priority
		})
		
		for _, group := range groups {
			fmt.Printf("  %s: %s (priority: %d)\n", group.name, group.desc, group.priority)
		}
	}

	if listTags {
		utils.PrintStatus(utils.Cyan, "Available Tags:")
		tagSet := make(map[string]bool)
		
		for _, group := range config.Groups {
			for _, pkg := range group.Packages {
				for _, tag := range pkg.Tags {
					tagSet[tag] = true
				}
			}
		}
		
		var tags []string
		for tag := range tagSet {
			tags = append(tags, tag)
		}
		sort.Strings(tags)
		
		for _, tag := range tags {
			fmt.Printf("  - %s\n", tag)
		}
	}

	if listProfiles {
		utils.PrintStatus(utils.Cyan, "Available Profiles:")
		for name, profile := range config.Profiles {
			fmt.Printf("  %s: %s\n", name, profile.Description)
			if len(profile.Groups) > 0 {
				fmt.Printf("    Groups: %v\n", profile.Groups)
			}
			if len(profile.Tags) > 0 {
				fmt.Printf("    Tags: %v\n", profile.Tags)
			}
			if len(profile.ExcludeTags) > 0 {
				fmt.Printf("    Exclude Tags: %v\n", profile.ExcludeTags)
			}
		}
	}

	return nil
}

func init() {
	rootCmd.AddCommand(installCmd)

	// Installation filters
	installCmd.Flags().StringVarP(&groups, "groups", "g", "", "Install only specified groups (comma-separated)")
	installCmd.Flags().StringVarP(&tags, "tags", "t", "", "Install only packages with specified tags (comma-separated)")
	installCmd.Flags().StringVarP(&profile, "profile", "p", "", "Install using predefined profile")

	// Package type skip flags
	installCmd.Flags().BoolVar(&skipTaps, "skip-taps", false, "Skip installing taps")
	installCmd.Flags().BoolVar(&skipBrews, "skip-brews", false, "Skip installing brew formulae")
	installCmd.Flags().BoolVar(&skipCasks, "skip-casks", false, "Skip installing casks")
	installCmd.Flags().BoolVar(&skipMas, "skip-mas", false, "Skip installing Mac App Store apps")

	// List commands
	installCmd.Flags().BoolVar(&listGroups, "list-groups", false, "List available groups")
	installCmd.Flags().BoolVar(&listTags, "list-tags", false, "List available tags")
	installCmd.Flags().BoolVar(&listProfiles, "list-profiles", false, "List available profiles")
} 
