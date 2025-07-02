# Brew Manager

A command-line tool for managing Homebrew packages with support for groups, tags, and profiles.

## Features

- **Sync**: Automatically sync installed packages to YAML configuration
- **Install**: Install packages from YAML configuration with filtering
- **Convert**: Convert Brewfile to YAML format
- **Validate**: Validate YAML configuration files
- **Generate**: Generate JSON schema from Go structs

## Schema Generation

The tool automatically generates JSON schema files from Go struct definitions to ensure they stay in sync with the code.

### Manual Schema Generation

```bash
# Generate schema file
./brew-manager generate schema

# Generate to custom location
./brew-manager generate schema --output /path/to/schema.json
```

### Automatic Schema Generation

```bash
# Use Go generate to regenerate all schemas
go generate
```

### Schema Validation (CI/CD)

```bash
# Check if schema is in sync with Go structs
./check-schema.sh
```

## Commands

### Sync

Synchronize currently installed packages with YAML configuration:

```bash
# Basic sync
./brew-manager sync

# Dry run to see what would be added
./brew-manager sync --dry-run

# Interactive mode with group/tag assignment
./brew-manager sync --interactive

# Auto-detect groups and tags
./brew-manager sync --auto-detect --sort
```

### Install

Install packages from YAML configuration:

```bash
# Install all packages
./brew-manager install

# Install specific groups
./brew-manager install --groups development,productivity

# Install with tags
./brew-manager install --tags essential,cli

# Use profile
./brew-manager install --profile developer
```

### Convert

Convert Brewfile to YAML format:

```bash
# Convert Brewfile to grouped YAML
./brew-manager convert /path/to/Brewfile output.yml

# With verbose output
./brew-manager convert /path/to/Brewfile output.yml --verbose
```

### Validate

Validate YAML configuration files:

```bash
# Validate single file
./brew-manager validate packages.yaml

# Validate all YAML files in directory
./brew-manager validate --all /path/to/data

# Verbose validation
./brew-manager validate packages.yaml --verbose
```

## Configuration Structure

The YAML configuration follows this structure:

```yaml
# yaml-language-server: $schema=schemas/packages-grouped.schema.json

groups:
  development:
    description: Development tools and environments
    priority: 2
    packages:
      - name: git
        type: brew
        tags: [essential, version-control]
      - name: visual-studio-code
        type: cask
        tags: [editor, development]

  productivity:
    description: Productivity applications
    priority: 3
    packages:
      - name: notion
        type: cask
        tags: [productivity, notes]

profiles:
  minimal:
    description: Minimal development setup
    groups: [development]
    tags: [essential]
  
  full:
    description: Complete setup
    groups: [development, productivity]
```

## Development

### Building

```bash
go build -o brew-manager
```

### Running Tests

```bash
go test ./...
```

### Generating Schemas

```bash
# Generate schema from Go structs
go generate

# Check schema sync
./check-schema.sh
```

### Schema Customization

To modify the JSON schema, update the struct tags in `pkg/types/types.go`:

```go
type Package struct {
    Name string `json:"name" jsonschema:"title=Package Name,description=Package name,required,minLength=1"`
    Type string `json:"type" jsonschema:"title=Package Type,description=Package type,required,enum=tap,enum=brew,enum=cask,enum=mas"`
    // ...
}
```

Then run `go generate` to update the schema file.

## CI/CD Integration

Add this to your CI pipeline to ensure schemas stay in sync:

```yaml
- name: Check Schema Sync
  run: |
    cd scripts/brew-management
    ./check-schema.sh
``` 
