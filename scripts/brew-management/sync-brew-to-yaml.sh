#!/bin/bash

# Brew to YAML sync script
# This script compares currently installed brew packages with YAML file
# and adds missing packages to an "uncategorized" section

set -euo pipefail

# Colors for output
readonly RED='\033[0;31m'
readonly GREEN='\033[0;32m'
readonly YELLOW='\033[1;33m'
readonly BLUE='\033[0;34m'
readonly CYAN='\033[0;36m'
readonly NC='\033[0m' # No Color

# Default YAML file path
DEFAULT_YAML_FILE="$(dirname "$0")/../../data/brew/packages.yml"

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

Sync currently installed Homebrew packages with YAML configuration file.
Adds missing packages to "uncategorized" sections.

Options:
    -h, --help          Show this help message
    -v, --verbose       Enable verbose output
    -d, --dry-run       Show what would be added without actually modifying the file
    -b, --backup        Create backup of YAML file before modification
    -s, --sort          Sort packages alphabetically within categories
    --show-only         Only show missing packages without modifying the file

Arguments:
    YAML_FILE          Path to YAML configuration file (default: ${DEFAULT_YAML_FILE})

Examples:
    $0                                      # Sync with default YAML file
    $0 --dry-run                           # Show what would be added
    $0 --backup my-packages.yml            # Sync with backup
    $0 --show-only                         # Only display missing packages

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

# Function to get currently installed packages
get_installed_packages() {
    local type=$1
    
    case $type in
        "taps")
            brew tap | sort
            ;;
        "brews")
            # Get all installed formulae including those from taps
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

# Function to get packages from YAML file
get_yaml_packages() {
    local yaml_file=$1
    local type=$2
    
    if [[ ! -f "$yaml_file" ]]; then
        return 0
    fi
    
    case $type in
        "taps")
            yq eval '.taps[]?' "$yaml_file" 2>/dev/null | grep -v "^null$" | sort || true
            ;;
        "brews")
            yq eval '.brews[]?' "$yaml_file" 2>/dev/null | grep -v "^null$" | sort || true
            ;;
        "casks")
            yq eval '.casks[]?' "$yaml_file" 2>/dev/null | grep -v "^null$" | sort || true
            ;;
        "mas")
            if yq eval '.mas_apps[]?' "$yaml_file" 2>/dev/null | grep -v "^null$" >/dev/null; then
                yq eval '.mas_apps[] | .name + "|" + (.id | tostring)' "$yaml_file" 2>/dev/null | sort || true
            fi
            ;;
    esac
}

# Function to find missing packages
find_missing_packages() {
    local yaml_file=$1
    local type=$2
    local installed_file="/tmp/installed_${type}.txt"
    local yaml_file_packages="/tmp/yaml_${type}.txt"
    local missing_file="/tmp/missing_${type}.txt"
    
    # Get installed packages
    get_installed_packages "$type" > "$installed_file"
    
    # Get YAML packages
    get_yaml_packages "$yaml_file" "$type" > "$yaml_file_packages"
    
    # Find missing packages (in installed but not in YAML)
    if [[ -s "$installed_file" ]]; then
        comm -23 "$installed_file" "$yaml_file_packages" > "$missing_file"
        cat "$missing_file"
    fi
    
    # Cleanup
    rm -f "$installed_file" "$yaml_file_packages" "$missing_file"
}

# Function to display missing packages
display_missing_packages() {
    local yaml_file=$1
    
    print_status "$CYAN" "=== Checking for missing packages ==="
    
    local has_missing=false
    
    # Check taps
    local missing_taps
    missing_taps=$(find_missing_packages "$yaml_file" "taps")
    if [[ -n "$missing_taps" ]]; then
        print_status "$YELLOW" "\nMissing taps:"
        echo "$missing_taps" | while IFS= read -r tap; do
            [[ -n "$tap" ]] && echo "  - $tap"
        done
        has_missing=true
    fi
    
    # Check brews
    local missing_brews
    missing_brews=$(find_missing_packages "$yaml_file" "brews")
    if [[ -n "$missing_brews" ]]; then
        print_status "$YELLOW" "\nMissing brew packages:"
        echo "$missing_brews" | while IFS= read -r brew_pkg; do
            [[ -n "$brew_pkg" ]] && echo "  - $brew_pkg"
        done
        has_missing=true
    fi
    
    # Check casks
    local missing_casks
    missing_casks=$(find_missing_packages "$yaml_file" "casks")
    if [[ -n "$missing_casks" ]]; then
        print_status "$YELLOW" "\nMissing casks:"
        echo "$missing_casks" | while IFS= read -r cask; do
            [[ -n "$cask" ]] && echo "  - $cask"
        done
        has_missing=true
    fi
    
    # Check MAS apps
    if command_exists mas; then
        local missing_mas
        missing_mas=$(find_missing_packages "$yaml_file" "mas")
        if [[ -n "$missing_mas" ]]; then
            print_status "$YELLOW" "\nMissing Mac App Store apps:"
            echo "$missing_mas" | while IFS= read -r mas_app; do
                if [[ -n "$mas_app" && "$mas_app" =~ ^(.+)\|([0-9]+)$ ]]; then
                    local name="${BASH_REMATCH[1]}"
                    local id="${BASH_REMATCH[2]}"
                    echo "  - name: \"$name\""
                    echo "    id: $id"
                fi
            done
            has_missing=true
        fi
    fi
    
    if [[ "$has_missing" == false ]]; then
        print_status "$GREEN" "✅ All installed packages are already in the YAML file!"
    fi
}

# Function to add missing packages to YAML
add_missing_packages() {
    local yaml_file=$1
    local backup=$2
    local sort_packages=$3
    
    # Create backup if requested
    if [[ "$backup" == true ]]; then
        local backup_file="${yaml_file}.backup.$(date +%Y%m%d_%H%M%S)"
        cp "$yaml_file" "$backup_file"
        print_status "$GREEN" "Backup created: $backup_file"
    fi
    
    # Create temporary file for new YAML content
    local temp_yaml="/tmp/new_packages.yml"
    cp "$yaml_file" "$temp_yaml"
    
    local has_changes=false
    
    # Add missing taps
    local missing_taps
    missing_taps=$(find_missing_packages "$yaml_file" "taps")
    if [[ -n "$missing_taps" ]]; then
        print_status "$BLUE" "Adding missing taps..."
        
        # Check if uncategorized_taps section exists
        if ! yq eval '.uncategorized_taps' "$temp_yaml" >/dev/null 2>&1; then
            # Add uncategorized_taps section
            yq eval '.uncategorized_taps = []' -i "$temp_yaml"
        fi
        
        while IFS= read -r tap; do
            if [[ -n "$tap" ]]; then
                print_status "$GREEN" "  Adding tap: $tap"
                yq eval ".uncategorized_taps += [\"$tap\"]" -i "$temp_yaml"
                has_changes=true
            fi
        done <<< "$missing_taps"
    fi
    
    # Add missing brews
    local missing_brews
    missing_brews=$(find_missing_packages "$yaml_file" "brews")
    if [[ -n "$missing_brews" ]]; then
        print_status "$BLUE" "Adding missing brew packages..."
        
        # Check if uncategorized_brews section exists
        if ! yq eval '.uncategorized_brews' "$temp_yaml" >/dev/null 2>&1; then
            # Add uncategorized_brews section
            yq eval '.uncategorized_brews = []' -i "$temp_yaml"
        fi
        
        while IFS= read -r brew_pkg; do
            if [[ -n "$brew_pkg" ]]; then
                print_status "$GREEN" "  Adding brew: $brew_pkg"
                yq eval ".uncategorized_brews += [\"$brew_pkg\"]" -i "$temp_yaml"
                has_changes=true
            fi
        done <<< "$missing_brews"
    fi
    
    # Add missing casks
    local missing_casks
    missing_casks=$(find_missing_packages "$yaml_file" "casks")
    if [[ -n "$missing_casks" ]]; then
        print_status "$BLUE" "Adding missing casks..."
        
        # Check if uncategorized_casks section exists
        if ! yq eval '.uncategorized_casks' "$temp_yaml" >/dev/null 2>&1; then
            # Add uncategorized_casks section
            yq eval '.uncategorized_casks = []' -i "$temp_yaml"
        fi
        
        while IFS= read -r cask; do
            if [[ -n "$cask" ]]; then
                print_status "$GREEN" "  Adding cask: $cask"
                yq eval ".uncategorized_casks += [\"$cask\"]" -i "$temp_yaml"
                has_changes=true
            fi
        done <<< "$missing_casks"
    fi
    
    # Add missing MAS apps
    if command_exists mas; then
        local missing_mas
        missing_mas=$(find_missing_packages "$yaml_file" "mas")
        if [[ -n "$missing_mas" ]]; then
            print_status "$BLUE" "Adding missing Mac App Store apps..."
            
            # Check if uncategorized_mas_apps section exists
            if ! yq eval '.uncategorized_mas_apps' "$temp_yaml" >/dev/null 2>&1; then
                # Add uncategorized_mas_apps section
                yq eval '.uncategorized_mas_apps = []' -i "$temp_yaml"
            fi
            
            while IFS= read -r mas_app; do
                if [[ -n "$mas_app" && "$mas_app" =~ ^(.+)\|([0-9]+)$ ]]; then
                    local name="${BASH_REMATCH[1]}"
                    local id="${BASH_REMATCH[2]}"
                    print_status "$GREEN" "  Adding MAS app: $name (ID: $id)"
                    yq eval ".uncategorized_mas_apps += [{\"name\": \"$name\", \"id\": $id}]" -i "$temp_yaml"
                    has_changes=true
                fi
            done <<< "$missing_mas"
        fi
    fi
    
    # Sort packages if requested
    if [[ "$sort_packages" == true && "$has_changes" == true ]]; then
        print_status "$BLUE" "Sorting packages alphabetically..."
        
        # Sort each array in the YAML file
        for section in taps brews casks uncategorized_taps uncategorized_brews uncategorized_casks; do
            if yq eval ".$section" "$temp_yaml" >/dev/null 2>&1; then
                yq eval ".$section |= sort" -i "$temp_yaml"
            fi
        done
        
        # Sort MAS apps by name
        for section in mas_apps uncategorized_mas_apps; do
            if yq eval ".$section" "$temp_yaml" >/dev/null 2>&1; then
                yq eval ".$section |= sort_by(.name)" -i "$temp_yaml"
            fi
        done
    fi
    
    # Replace original file if changes were made
    if [[ "$has_changes" == true ]]; then
        mv "$temp_yaml" "$yaml_file"
        print_status "$GREEN" "✅ YAML file updated successfully!"
    else
        rm -f "$temp_yaml"
        print_status "$GREEN" "✅ No missing packages found. YAML file is up to date!"
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
        
        # Add missing packages
        add_missing_packages "$yaml_file" "$backup" "$sort_packages"
    fi
}

# Run main function with all arguments
main "$@" 
