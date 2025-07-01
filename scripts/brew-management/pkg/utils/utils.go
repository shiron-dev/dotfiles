package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

var (
	Red     = color.New(color.FgRed)
	Green   = color.New(color.FgGreen)
	Yellow  = color.New(color.FgYellow)
	Blue    = color.New(color.FgBlue)
	Cyan    = color.New(color.FgCyan)
	Magenta = color.New(color.FgMagenta)
)

// PrintStatus prints colored status messages
func PrintStatus(colorFunc *color.Color, message string) {
	colorFunc.Println(message)
}

// CommandExists checks if a command is available in the system PATH
func CommandExists(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}

// RunCommand executes a command and returns the output
func RunCommand(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

// RunCommandSilent executes a command without capturing output
func RunCommandSilent(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	return cmd.Run()
}

// CheckPrerequisites verifies that required tools are installed
func CheckPrerequisites() error {
	PrintStatus(Blue, "Checking prerequisites...")

	// Check if Homebrew is installed
	if !CommandExists("brew") {
		return fmt.Errorf("homebrew is not installed. Please install Homebrew first")
	}

	// Check if yq is available
	if !CommandExists("yq") {
		PrintStatus(Yellow, "Warning: yq is not installed. Installing yq for YAML parsing...")
		if err := RunCommandSilent("brew", "install", "yq"); err != nil {
			return fmt.Errorf("failed to install yq: %w", err)
		}
	}

	PrintStatus(Green, "Prerequisites check completed.")
	return nil
}

// GetDefaultYAMLPath returns the default path for YAML configuration files
func GetDefaultYAMLPath(filename string) string {
	defaultPath := filepath.Join(os.Getenv("HOME"), "projects/github.com/shiron-dev/dotfiles/data/brew/packages.yaml")
	if filename == "" || filename == "packages-grouped.yml" || filename == "packages.yml" || filename == "packages.yaml" {
		return defaultPath
	}
	dir, _ := os.Getwd()
	return filepath.Join(dir, "../../data/brew", filename)
}

// FileExists checks if a file exists
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// ContainsString checks if a slice contains a string
func ContainsString(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

// SplitCommaSeparated splits a comma-separated string into a slice
func SplitCommaSeparated(str string) []string {
	if str == "" {
		return nil
	}
	parts := strings.Split(str, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// HasIntersection checks if two string slices have any common elements
func HasIntersection(slice1, slice2 []string) bool {
	for _, s1 := range slice1 {
		for _, s2 := range slice2 {
			if s1 == s2 {
				return true
			}
		}
	}
	return false
}

// AutoDetectGroup attempts to determine the appropriate group for a package
func AutoDetectGroup(packageName, packageType string) string {
	name := strings.ToLower(packageName)

	switch {
	case strings.Contains(name, "git") || strings.Contains(name, "node") ||
		strings.Contains(name, "python") || strings.Contains(name, "go") ||
		strings.Contains(name, "rust") || strings.Contains(name, "java") ||
		strings.Contains(name, "docker") || strings.Contains(name, "terraform") ||
		strings.Contains(name, "ansible"):
		return "development"
	case name == "htop" || name == "tree" || name == "watch" ||
		name == "stats" || name == "battery" || name == "raycast" ||
		strings.Contains(name, "1password"):
		return "system"
	case strings.Contains(name, "figma") || name == "obs" || name == "vlc" ||
		name == "audacity" || name == "gimp" || name == "inkscape":
		return "creative"
	case name == "notion" || name == "slack" || name == "zoom" ||
		strings.Contains(name, "chrome") || strings.Contains(name, "firefox") ||
		name == "arc" || name == "brave":
		return "productivity"
	case name == "mas" || name == "brew" || name == "yq" || name == "jq":
		return "core"
	default:
		return "optional"
	}
}

// AutoDetectTags attempts to determine appropriate tags for a package
func AutoDetectTags(packageName, packageType string) []string {
	name := strings.ToLower(packageName)
	var tags []string

	switch {
	case strings.Contains(name, "python"):
		tags = append(tags, "language", "python")
	case strings.Contains(name, "node") || strings.Contains(name, "npm") || strings.Contains(name, "yarn"):
		tags = append(tags, "language", "javascript", "nodejs")
	case strings.Contains(name, "go"):
		tags = append(tags, "language", "golang")
	case strings.Contains(name, "rust"):
		tags = append(tags, "language", "rust")
	case strings.Contains(name, "java"):
		tags = append(tags, "language", "java")
	case strings.Contains(name, "git"):
		tags = append(tags, "version-control", "essential")
	case strings.Contains(name, "docker"):
		tags = append(tags, "container", "development")
	case strings.Contains(name, "terraform"):
		tags = append(tags, "infrastructure", "cloud")
	case strings.Contains(name, "ansible"):
		tags = append(tags, "automation", "infrastructure")
	case name == "bat" || name == "fd" || name == "fzf" || name == "ripgrep" || name == "htop" || name == "tree":
		tags = append(tags, "cli", "productivity")
	case strings.Contains(name, "chrome"):
		tags = append(tags, "browser", "google")
	case strings.Contains(name, "firefox"):
		tags = append(tags, "browser", "mozilla")
	case name == "arc":
		tags = append(tags, "browser", "modern")
	case name == "slack":
		tags = append(tags, "communication", "team")
	case name == "zoom":
		tags = append(tags, "video-call", "meeting")
	case strings.Contains(name, "figma"):
		tags = append(tags, "design", "ui-ux")
	case name == "obs":
		tags = append(tags, "streaming", "recording")
	case name == "vlc":
		tags = append(tags, "media-player", "video")
	case name == "stats" || name == "battery":
		tags = append(tags, "monitoring", "system")
	case name == "raycast":
		tags = append(tags, "launcher", "productivity")
	case strings.Contains(name, "1password"):
		tags = append(tags, "security", "password")
	}

	// Add package type as tag
	switch packageType {
	case "brew":
		tags = append(tags, "formula")
	case "cask":
		tags = append(tags, "application")
	case "tap":
		tags = append(tags, "tap")
	case "mas":
		tags = append(tags, "app-store")
	}

	return tags
}

// EnsureDir creates the directory for a file path if it doesn't exist
func EnsureDir(filePath string) error {
	dir := filepath.Dir(filePath)
	if dir == "." || dir == "/" {
		return nil
	}

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
		PrintStatus(Green, fmt.Sprintf("Created directory: %s", dir))
	}
	return nil
}

// CreateBackup creates a backup file with timestamp
func CreateBackup(filePath string) error {
	if !FileExists(filePath) {
		return nil
	}

	backupPath := filePath + ".backup"

	input, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file for backup: %w", err)
	}

	err = os.WriteFile(backupPath, input, 0644)
	if err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}

	PrintStatus(Green, fmt.Sprintf("Backup created: %s", backupPath))
	return nil
}
