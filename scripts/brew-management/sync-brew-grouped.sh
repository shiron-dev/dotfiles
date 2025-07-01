#!/bin/bash

# Enhanced Brew to YAML sync script with groups and tags support
# This script compares currently installed brew packages with YAML file
# and adds missing packages to "uncategorized" sections with optional group/tag assignment

set -euo pipefail

# Colors for output
readonly RED='\033[0;31m'
readonly GREEN='\033[0;32m'
readonly YELLOW='\033[1;33m'
readonly BLUE='\033[0;34m'
readonly CYAN='\033[0;36m'
readonly MAGENTA='\033[0;35m'
readonly NC='\033[0m' # No Color

# Default YAML file path
DEFAULT_YAML_FILE="$(dirname "$0")/../../data/brew/packages-grouped.yml"

# Function to print colored output
print_status() {
    local color=$1
    local message=$2
    echo -e "${color}${message}${NC}"
}

# Function to print usage
usage() {
    cat << EOF
Usage: $0 [OPTIONS] [YAML_FILE]

Sync currently installed Homebrew packages with YAML configuration file (with groups/tags support).
Adds missing packages to appropriate sections with optional group/tag assignment.

Options:
    -h, --help              Show this help message
    -v, --verbose           Enable verbose output
    -d, --dry-run           Show what would be added without actually modifying the file
    -b, --backup            Create backup of YAML file before modification
    -s, --sort              Sort packages alphabetically within categories
    --show-only             Only show missing packages without modifying the file
    
    # Group/Tag assignment for new packages
    --default-group GROUP   Assign this group to newly added packages (default: optional)
    --default-tags TAGS     Assign these tags to newly added packages (comma-separated)
    --interactive           Prompt for group/tag assignment for each new package
    --auto-detect           Try to auto-detect appropriate groups/tags based on package names

Arguments:
    YAML_FILE              Path to YAML configuration file (default: ${DEFAULT_YAML_FILE})

Examples:
    $0                                          # Sync with default settings
    $0 --dry-run                               # Show what would be added
    $0 --backup --default-group system         # Add missing packages to 'system' group
    $0 --interactive                           # Prompt for group/tag for each package
    $0 --auto-detect --sort                    # Auto-detect groups/tags and sort

EOF
}

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to check prerequisites
check_prerequisites() {
    print_status "$BLUE" "Checking prerequisites..."
    
    # Check if Homebrew is installed
    if ! command_exists brew; then
        print_status "$RED" "Error: Homebrew is not installed."
        exit 1
    fi
    
    # Check if yq is installed for YAML parsing
    if ! command_exists yq; then
        print_status "$YELLOW" "Warning: yq is not installed. Installing yq for YAML parsing..."
        brew install yq
    fi
    
    print_status "$GREEN" "Prerequisites check completed."
}

# Function to auto-detect group based on package name
auto_detect_group() {
    local package_name=$1
    local package_type=$2
    
    case $package_name in
        # Development tools
        git*|node*|python*|go|rust|java*|kotlin*|docker*|kubernetes*|terraform*|ansible*)
            echo "development"
            ;;
        # System utilities
        htop|tree|watch|stats|battery|raycast|alfred|1password*)
            echo "system"
            ;;
        # Creative tools
        figma*|obs|vlc|audacity|gimp|inkscape*)
            echo "creative"
            ;;
        # Productivity tools
        notion|slack|zoom|chrome*|firefox*|safari*|arc|brave*)
            echo "productivity"
            ;;
        # Core tools
        mas|brew|yq|jq)
            echo "core"
            ;;
        *)
            echo "optional"
            ;;
    esac
}

# Function to auto-detect tags based on package name
auto_detect_tags() {
    local package_name=$1
    local package_type=$2
    local tags=()
    
    case $package_name in
        # Programming languages
        *python*) tags+=("language" "python") ;;
        *node*|*npm*|*yarn*) tags+=("language" "javascript" "nodejs") ;;
        *go*) tags+=("language" "golang") ;;
        *rust*) tags+=("language" "rust") ;;
        *java*) tags+=("language" "java") ;;
        
        # Development tools
        git*) tags+=("version-control" "essential") ;;
        docker*) tags+=("container" "development") ;;
        *kubernetes*|*k8s*) tags+=("container" "orchestration") ;;
        terraform*) tags+=("infrastructure" "cloud") ;;
        ansible*) tags+=("automation" "infrastructure") ;;
        
        # CLI tools
        *cli*|bat|fd|fzf|ripgrep|htop|tree) tags+=("cli" "productivity") ;;
        
        # Browsers
        *chrome*) tags+=("browser" "google") ;;
        *firefox*) tags+=("browser" "mozilla") ;;
        *safari*) tags+=("browser" "apple") ;;
        arc) tags+=("browser" "modern") ;;
        
        # Communication
        slack) tags+=("communication" "team") ;;
        zoom) tags+=("video-call" "meeting") ;;
        
        # Creative
        figma*) tags+=("design" "ui-ux") ;;
        obs) tags+=("streaming" "recording") ;;
        vlc) tags+=("media-player" "video") ;;
        
        # System
        stats|battery) tags+=("monitoring" "system") ;;
        raycast|alfred) tags+=("launcher" "productivity") ;;
        1password*) tags+=("security" "password") ;;
    esac
    
    # Add package type as tag
    case $package_type in
        "brews") tags+=("formula") ;;
        "casks") tags+=("application") ;;
        "taps") tags+=("tap") ;;
        "mas") tags+=("app-store") ;;
    esac
    
    # Convert array to comma-separated string
    IFS=',' && echo "${tags[*]}"
}

# Function to prompt for group/tag assignment
prompt_for_assignment() {
    local package_name=$1
    local package_type=$2
    local suggested_group=$3
    local suggested_tags=$4
    
    echo ""
    print_status "$CYAN" "Package: $package_name ($package_type)"
    print_status "$YELLOW" "Suggested group: $suggested_group"
    print_status "$YELLOW" "Suggested tags: $suggested_tags"
    
    read -p "Enter group (or press Enter for suggested): " input_group
    if [[ -z "$input_group" ]]; then
        input_group="$suggested_group"
    fi
    
    read -p "Enter tags (comma-separated, or press Enter for suggested): " input_tags
    if [[ -z "$input_tags" ]]; then
        input_tags="$suggested_tags"
    fi
    
    echo "${input_group}|${input_tags}"
}

# Function to get currently installed packages (same as sync-brew-to-yaml.sh)
get_installed_packages() {
    local type=$1
    
    case $type in
        "taps")
            brew tap | sort
            ;;
        "brews")
            brew list --formula | sort
            ;;
        "casks")
            brew list --cask | sort
            ;;
        "mas")
            if command_exists mas; then
                mas list | sort -k2 | while IFS= read -r line; do
                    if [[ $line =~ ^([0-9]+)\ (.+)\ \([0-9.]+\)$ ]]; then
                        local id="${BASH_REMATCH[1]}"
                        local name="${BASH_REMATCH[2]}"
                        # Trim leading and trailing whitespace from name
                        name=$(echo "$name" | sed 's/^[[:space:]]*//;s/[[:space:]]*$//')
                        echo "$name|$id"
                    fi
                done
            fi
            ;;
    esac
}

# Function to get existing packages from YAML groups
get_existing_packages() {
    local yaml_file=$1
    local package_type=$2
    
    case "$package_type" in
        "taps")
            yq eval '.groups[].packages[] | select(.type == "tap") | .name' "$yaml_file" 2>/dev/null | grep -v "^null$" | sort || true
            ;;
        "brews")
            yq eval '.groups[].packages[] | select(.type == "brew") | .name' "$yaml_file" 2>/dev/null | grep -v "^null$" | sort || true
            ;;
        "casks")
            yq eval '.groups[].packages[] | select(.type == "cask") | .name' "$yaml_file" 2>/dev/null | grep -v "^null$" | sort || true
            ;;
        "mas")
            if yq eval '.groups[].packages[] | select(.type == "mas")' "$yaml_file" 2>/dev/null | grep -v "^null$" >/dev/null; then
                yq eval '.groups[].packages[] | select(.type == "mas") | .name + "|" + (.id | tostring)' "$yaml_file" 2>/dev/null | sort || true
            fi
            ;;
    esac
}

# Function to add package to uncategorized group
add_uncategorized_package() {
    local temp_yaml=$1
    local package_type=$2
    local package_name=$3
    local package_id=$4
    local group="uncategorized"
    local tags="uncategorized"
    
    # Check if uncategorized group exists, if not create it
    if ! yq eval '.groups.uncategorized' "$temp_yaml" >/dev/null 2>&1; then
        yq eval '.groups.uncategorized = {"description": "Uncategorized packages", "priority": 99, "packages": []}' -i "$temp_yaml"
    fi
    
    case "$package_type" in
        "taps")
            yq eval ".groups.uncategorized.packages += [{\"name\": \"$package_name\", \"type\": \"tap\", \"tags\": [\"$tags\"]}]" -i "$temp_yaml"
            ;;
        "brews")
            yq eval ".groups.uncategorized.packages += [{\"name\": \"$package_name\", \"type\": \"brew\", \"tags\": [\"$tags\"]}]" -i "$temp_yaml"
            ;;
        "casks")
            yq eval ".groups.uncategorized.packages += [{\"name\": \"$package_name\", \"type\": \"cask\", \"tags\": [\"$tags\"]}]" -i "$temp_yaml"
            ;;
        "mas")
            yq eval ".groups.uncategorized.packages += [{\"name\": \"$package_name\", \"type\": \"mas\", \"id\": $package_id, \"tags\": [\"$tags\"]}]" -i "$temp_yaml"
            ;;
    esac
}

# Function to sort packages within each group
sort_yaml_packages() {
    local temp_yaml=$1
    
    # Sort packages within each group by name
    for group in $(yq eval '.groups | keys | .[]' "$temp_yaml" 2>/dev/null); do
        if yq eval ".groups.${group}.packages" "$temp_yaml" >/dev/null 2>&1; then
            yq eval ".groups.${group}.packages |= sort_by(.name)" -i "$temp_yaml"
        fi
    done
}

# Function to add package to a specific group
add_package_to_group() {
    local temp_yaml=$1
    local group=$2
    local package_name=$3
    local package_type=$4
    local package_id=$5
    local tags=$6
    
    # Check if group exists, if not create it
    if ! yq eval ".groups.${group}" "$temp_yaml" >/dev/null 2>&1; then
        local priority=99
        case "$group" in
            "core") priority=1 ;;
            "development") priority=2 ;;
            "productivity") priority=3 ;;
            "creative") priority=4 ;;
            "system") priority=5 ;;
            *) priority=99 ;;
        esac
        
        yq eval ".groups.${group} = {\"description\": \"Auto-created group\", \"priority\": $priority, \"packages\": []}" -i "$temp_yaml"
    fi
    
    # Format tags
    local tags_formatted=""
    if [[ -n "$tags" ]]; then
        tags_formatted=$(echo "$tags" | sed 's/,/", "/g' | sed 's/^/"/;s/$/"/')
    fi
    
    # Add package to group based on type
    if [[ "$package_type" == "mas" && -n "$package_id" ]]; then
        yq eval ".groups.${group}.packages += [{\"name\": \"$package_name\", \"type\": \"$package_type\", \"id\": $package_id, \"tags\": [$tags_formatted]}]" -i "$temp_yaml"
    else
        yq eval ".groups.${group}.packages += [{\"name\": \"$package_name\", \"type\": \"$package_type\", \"tags\": [$tags_formatted]}]" -i "$temp_yaml"
    fi
}

# Function to find missing packages (same as sync-brew-to-yaml.sh)
find_missing_packages() {
    local yaml_file=$1
    local type=$2
    local installed_file="/tmp/installed_${type}.txt"
    local yaml_file_packages="/tmp/yaml_${type}.txt"
    local missing_file="/tmp/missing_${type}.txt"
    
    # Get installed packages
    get_installed_packages "$type" > "$installed_file"
    
    # Get YAML packages
    get_existing_packages "$yaml_file" "$type" > "$yaml_file_packages"
    
    # Find missing packages (in installed but not in YAML)
    if [[ -s "$installed_file" ]]; then
        comm -23 "$installed_file" "$yaml_file_packages" > "$missing_file"
        cat "$missing_file"
    fi
    
    # Cleanup
    rm -f "$installed_file" "$yaml_file_packages" "$missing_file"
}

# Function to add missing packages with group/tag support
add_missing_packages_grouped() {
    local yaml_file=$1
    local backup=$2
    local sort_packages=$3
    local default_group=$4
    local default_tags=$5
    local interactive=$6
    local auto_detect=$7
    
    # Create backup if requested
    if [[ "$backup" == true ]]; then
        local backup_file="${yaml_file}.backup.$(date +%Y%m%d_%H%M%S)"
        cp "$yaml_file" "$backup_file"
        print_status "$GREEN" "Backup created: $backup_file"
    fi
    
    # Create temporary file for new YAML content
    local temp_yaml="/tmp/new_packages_grouped.yml"
    cp "$yaml_file" "$temp_yaml"
    
    local has_changes=false
    
    # Process each package type
    for package_type in taps brews casks mas; do
        local missing_packages
        missing_packages=$(find_missing_packages "$yaml_file" "$package_type")
        
        if [[ -n "$missing_packages" ]]; then
            print_status "$BLUE" "Adding missing $package_type..."
            
            while IFS= read -r package_info; do
                if [[ -n "$package_info" ]]; then
                    local package_name group tags
                    
                    if [[ "$package_type" == "mas" ]]; then
                        if [[ "$package_info" =~ ^(.+)\|([0-9]+)$ ]]; then
                            package_name="${BASH_REMATCH[1]}"
                            local package_id="${BASH_REMATCH[2]}"
                        fi
                    else
                        package_name="$package_info"
                    fi
                    
                    # Determine group and tags
                    if [[ "$auto_detect" == true ]]; then
                        group=$(auto_detect_group "$package_name" "$package_type")
                        tags=$(auto_detect_tags "$package_name" "$package_type")
                    elif [[ "$interactive" == true ]]; then
                        local suggested_group suggested_tags assignment
                        suggested_group=$(auto_detect_group "$package_name" "$package_type")
                        suggested_tags=$(auto_detect_tags "$package_name" "$package_type")
                        assignment=$(prompt_for_assignment "$package_name" "$package_type" "$suggested_group" "$suggested_tags")
                        IFS='|' read -r group tags <<< "$assignment"
                    else
                        group="$default_group"
                        tags="$default_tags"
                    fi
                    
                    # Add package to YAML using new group structure
                    case $package_type in
                        "taps")
                            add_package_to_group "$temp_yaml" "$group" "$package_name" "tap" "" "$tags"
                            ;;
                        "brews")
                            add_package_to_group "$temp_yaml" "$group" "$package_name" "brew" "" "$tags"
                            ;;
                        "casks")
                            add_package_to_group "$temp_yaml" "$group" "$package_name" "cask" "" "$tags"
                            ;;
                        "mas")
                            add_package_to_group "$temp_yaml" "$group" "$package_name" "mas" "$package_id" "$tags"
                            ;;
                    esac
                    
                    print_status "$GREEN" "  Added $package_type: $package_name [group: $group, tags: $tags]"
                    has_changes=true
                fi
            done <<< "$missing_packages"
        fi
    done
    
    # Sort packages if requested
    if [[ "$sort_packages" == true && "$has_changes" == true ]]; then
        print_status "$BLUE" "Sorting packages alphabetically..."
        sort_yaml_packages "$temp_yaml"
    fi
    
    # Replace original file if changes were made
    if [[ "$has_changes" == true ]]; then
        mv "$temp_yaml" "$yaml_file"
        print_status "$GREEN" "✅ YAML file updated successfully with group/tag information!"
    else
        rm -f "$temp_yaml"
        print_status "$GREEN" "✅ No missing packages found. YAML file is up to date!"
    fi
}

# Function to display missing packages
display_missing_packages() {
    local yaml_file=$1
    
    print_status "$CYAN" "=== Checking for missing packages ==="
    
    local has_missing=false
    
    # Check each package type
    for package_type in taps brews casks mas; do
        local missing_packages
        missing_packages=$(find_missing_packages "$yaml_file" "$package_type")
        
        if [[ -n "$missing_packages" ]]; then
            print_status "$YELLOW" "\nMissing $package_type:"
            while IFS= read -r package_info; do
                if [[ -n "$package_info" ]]; then
                    if [[ "$package_type" == "mas" && "$package_info" =~ ^(.+)\|([0-9]+)$ ]]; then
                        local name="${BASH_REMATCH[1]}"
                        local id="${BASH_REMATCH[2]}"
                        echo "  - name: \"$name\""
                        echo "    id: $id"
                    else
                        echo "  - $package_info"
                    fi
                fi
            done <<< "$missing_packages"
            has_missing=true
        fi
    done
    
    if [[ "$has_missing" == false ]]; then
        print_status "$GREEN" "✅ All installed packages are already in the YAML file!"
    fi
}

# Main function
main() {
    local yaml_file="$DEFAULT_YAML_FILE"
    local dry_run=false
    local verbose=false
    local backup=false
    local sort_packages=false
    local show_only=false
    local default_group="optional"
    local default_tags=""
    local interactive=false
    local auto_detect=false
    
    # Parse command line arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            -h|--help)
                usage
                exit 0
                ;;
            -v|--verbose)
                verbose=true
                shift
                ;;
            -d|--dry-run)
                dry_run=true
                shift
                ;;
            -b|--backup)
                backup=true
                shift
                ;;
            -s|--sort)
                sort_packages=true
                shift
                ;;
            --show-only)
                show_only=true
                shift
                ;;
            --default-group)
                default_group="$2"
                shift 2
                ;;
            --default-tags)
                default_tags="$2"
                shift 2
                ;;
            --interactive)
                interactive=true
                shift
                ;;
            --auto-detect)
                auto_detect=true
                shift
                ;;
            -*)
                print_status "$RED" "Unknown option: $1"
                usage
                exit 1
                ;;
            *)
                yaml_file="$1"
                shift
                ;;
        esac
    done
    
    # Check if YAML file exists
    if [[ ! -f "$yaml_file" ]]; then
        print_status "$RED" "Error: YAML file not found: $yaml_file"
        exit 1
    fi
    
    print_status "$BLUE" "Using YAML file: $yaml_file"
    
    # Check prerequisites
    check_prerequisites
    
    if [[ "$show_only" == true || "$dry_run" == true ]]; then
        display_missing_packages "$yaml_file"
        if [[ "$dry_run" == true ]]; then
            print_status "$YELLOW" "\n=== DRY RUN MODE - No changes were made ==="
        fi
    else
        # Show missing packages first
        display_missing_packages "$yaml_file"
        
        # Add missing packages with group/tag support
        add_missing_packages_grouped "$yaml_file" "$backup" "$sort_packages" "$default_group" "$default_tags" "$interactive" "$auto_detect"
    fi
}

# Run main function with all arguments
main "$@" 
