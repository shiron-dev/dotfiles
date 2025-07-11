package types

// PackageGrouped represents the grouped YAML configuration format
type PackageGrouped struct {
	Groups   map[string]Group  `yaml:"groups" json:"groups" jsonschema:"title=Package Groups,description=Package groups definition,required"`
	Profiles map[string]Profile `yaml:"profiles" json:"profiles" jsonschema:"title=Installation Profiles,description=Installation profiles - predefined combinations"`
}

// Group represents a package group with description and priority
type Group struct {
	Description string                       `yaml:"description" json:"description" jsonschema:"title=Description,description=Human-readable description of the group,required,minLength=1"`
	Priority    int                          `yaml:"priority" json:"priority" jsonschema:"title=Priority,description=Installation priority (lower numbers install first),required,minimum=1,maximum=99"`
	Packages    map[string][]PackageInfo `yaml:"packages" json:"packages" jsonschema:"title=Packages,description=Packages in this group,required"`
}

// PackageInfo represents a single package with metadata, where the type is determined by the map key in Group.Packages
type PackageInfo struct {
	Name        string   `yaml:"name" json:"name" jsonschema:"title=Package Name,description=Package name,required,minLength=1"`
	Tags        []string `yaml:"tags,omitempty" json:"tags,omitempty" jsonschema:"title=Tags,description=Tags for categorization and filtering,uniqueItems"`
	Description string   `yaml:"description,omitempty" json:"description,omitempty" jsonschema:"title=Description,description=Optional description of the package,minLength=1"`
	ID          int64    `yaml:"id,omitempty" json:"id,omitempty" jsonschema:"title=App Store ID,description=Mac App Store ID (required for mas type),minimum=1"` // For mas apps
}

// Profile represents an installation profile
type Profile struct {
	Description  string   `yaml:"description" json:"description" jsonschema:"title=Description,description=Human-readable description of the profile,required,minLength=1"`
	Groups       []string `yaml:"groups,omitempty" json:"groups,omitempty" jsonschema:"title=Groups,description=Groups to include in this profile,uniqueItems"`
	Tags         []string `yaml:"tags,omitempty" json:"tags,omitempty" jsonschema:"title=Tags,description=Tags to include in this profile,uniqueItems"`
	ExcludeTags  []string `yaml:"exclude_tags,omitempty" json:"exclude_tags,omitempty" jsonschema:"title=Exclude Tags,description=Tags to exclude from this profile,uniqueItems"`
}

// FilteredPackage represents a package with its type, used for filtering results
type FilteredPackage struct {
	PackageInfo
	Type string
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
	Profile        string
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

// PruneOptions represents prune configuration
type PruneOptions struct {
	DryRun     bool
	Verbose    bool
	SkipTaps   bool
	SkipBrews  bool
	SkipCasks  bool
	SkipMas    bool
	ConfirmAll bool
} 
