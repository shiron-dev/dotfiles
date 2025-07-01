package types

// PackageGrouped represents the grouped YAML configuration format
type PackageGrouped struct {
	Metadata Metadata          `yaml:"metadata"`
	Groups   map[string]Group  `yaml:"groups"`
	Profiles map[string]Profile `yaml:"profiles"`
}

// Metadata contains version and feature information
type Metadata struct {
	Version        string `yaml:"version"`
	SupportsGroups bool   `yaml:"supports_groups"`
	SupportsTags   bool   `yaml:"supports_tags"`
}

// Group represents a package group with description and priority
type Group struct {
	Description string    `yaml:"description"`
	Priority    int       `yaml:"priority"`
	Packages    []Package `yaml:"packages"`
}

// Package represents a single package with type and metadata
type Package struct {
	Name        string   `yaml:"name"`
	Type        string   `yaml:"type"`
	Tags        []string `yaml:"tags,omitempty"`
	Description string   `yaml:"description,omitempty"`
	ID          int64    `yaml:"id,omitempty"` // For mas apps
}

// Profile represents an installation profile
type Profile struct {
	Description  string   `yaml:"description"`
	Groups       []string `yaml:"groups,omitempty"`
	Tags         []string `yaml:"tags,omitempty"`
	ExcludeTags  []string `yaml:"exclude_tags,omitempty"`
}

// PackageSimple represents the simple YAML configuration format
type PackageSimple struct {
	Taps    []string `yaml:"taps,omitempty"`
	Brews   []string `yaml:"brews,omitempty"`
	Casks   []string `yaml:"casks,omitempty"`
	MasApps []MasApp `yaml:"mas_apps,omitempty"`
}

// MasApp represents a Mac App Store application
type MasApp struct {
	Name string `yaml:"name"`
	ID   int64  `yaml:"id"`
}

// InstallOptions represents installation configuration
type InstallOptions struct {
	DryRun         bool
	Verbose        bool
	Groups         []string
	Tags           []string
	ExcludeGroups  []string
	ExcludeTags    []string
	Profile        string
	TapsOnly       bool
	BrewsOnly      bool
	CasksOnly      bool
	MasOnly        bool
	SkipTaps       bool
	SkipBrews      bool
	SkipCasks      bool
	SkipMas        bool
}

// SyncOptions represents synchronization configuration
type SyncOptions struct {
	DryRun       bool
	Verbose      bool
	Backup       bool
	Sort         bool
	ShowOnly     bool
	DefaultGroup string
	DefaultTags  []string
	Interactive  bool
	AutoDetect   bool
}

// ValidateOptions represents validation configuration
type ValidateOptions struct {
	Verbose    bool
	All        bool
	SchemaFile string
} 
