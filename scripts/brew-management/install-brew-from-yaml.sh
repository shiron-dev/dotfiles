#!/bin/bash

# YAML-based Brew installer script
# This script reads a YAML file and installs packages using Homebrew

set -euo pipefail

# Colors for output
readonly RED='\033[0;31m'
readonly GREEN='\033[0;32m'
readonly YELLOW='\033[1;33m'
readonly BLUE='\033[0;34m'
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

Install Homebrew packages from a YAML configuration file.

Options:
    -h, --help          Show this help message
    -v, --verbose       Enable verbose output
    -d, --dry-run       Show what would be installed without actually installing
    -t, --taps-only     Install only taps
    -b, --brews-only    Install only brew formulae
    -c, --casks-only    Install only casks
    -m, --mas-only      Install only Mac App Store apps
    --skip-taps         Skip installing taps
    --skip-brews        Skip installing brew formulae
    --skip-casks        Skip installing casks
    --skip-mas          Skip installing Mac App Store apps

Arguments:
    YAML_FILE          Path to YAML configuration file (default: ${DEFAULT_YAML_FILE})

Examples:
    $0                                      # Install all packages from default YAML file
    $0 my-packages.yml                      # Install from custom YAML file
    $0 --dry-run                           # Show what would be installed
    $0 --brews-only                        # Install only brew formulae
    $0 --skip-mas my-packages.yml          # Install everything except Mac App Store apps

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
        print_status "$RED" "Error: Homebrew is not installed. Please install Homebrew first."
        exit 1
    fi
    
    # Check if yq is installed for YAML parsing
    if ! command_exists yq; then
        print_status "$YELLOW" "Warning: yq is not installed. Installing yq for YAML parsing..."
        brew install yq
    fi
    
    # Check if mas is available for Mac App Store installations
    if ! command_exists mas; then
        print_status "$YELLOW" "Warning: mas is not installed. Mac App Store apps will be skipped."
        SKIP_MAS=true
    fi
    
    print_status "$GREEN" "Prerequisites check completed."
}

# Function to parse YAML and extract array elements
parse_yaml_array() {
    local yaml_file=$1
    local key_path=$2
    
    if [[ -f "$yaml_file" ]]; then
        yq eval ".${key_path}[]" "$yaml_file" 2>/dev/null | grep -v "^null$" || true
    fi
}

# Function to parse YAML and extract MAS apps with id
parse_mas_apps() {
    local yaml_file=$1
    
    if [[ -f "$yaml_file" ]]; then
        yq eval '.mas_apps[] | .name + " (ID: " + (.id | tostring) + ")"' "$yaml_file" 2>/dev/null | grep -v "^null" || true
    fi
}

# Function to install taps
install_taps() {
    local yaml_file=$1
    local taps uncategorized_taps
    
    print_status "$BLUE" "Installing Homebrew taps..."
    
    taps=$(parse_yaml_array "$yaml_file" "taps")
    uncategorized_taps=$(parse_yaml_array "$yaml_file" "uncategorized_taps")
    
    # Combine both regular and uncategorized taps
    local all_taps
    all_taps=$(echo -e "$taps\n$uncategorized_taps" | grep -v '^$' | sort -u)
    
    if [[ -z "$all_taps" ]]; then
        print_status "$YELLOW" "No taps found in YAML file."
        return 0
    fi
    
    while IFS= read -r tap; do
        if [[ -n "$tap" ]]; then
            if [[ "$DRY_RUN" == true ]]; then
                print_status "$YELLOW" "DRY RUN: Would install tap: $tap"
            else
                print_status "$GREEN" "Installing tap: $tap"
                if [[ "$VERBOSE" == true ]]; then
                    brew tap "$tap"
                else
                    brew tap "$tap" >/dev/null 2>&1 || print_status "$RED" "Failed to install tap: $tap"
                fi
            fi
        fi
    done <<< "$all_taps"
    
    print_status "$GREEN" "Taps installation completed."
}

# Function to install brew formulae
install_brews() {
    local yaml_file=$1
    local brews uncategorized_brews
    
    print_status "$BLUE" "Installing Homebrew formulae..."
    
    brews=$(parse_yaml_array "$yaml_file" "brews")
    uncategorized_brews=$(parse_yaml_array "$yaml_file" "uncategorized_brews")
    
    # Combine both regular and uncategorized brews
    local all_brews
    all_brews=$(echo -e "$brews\n$uncategorized_brews" | grep -v '^$' | sort -u)
    
    if [[ -z "$all_brews" ]]; then
        print_status "$YELLOW" "No brew formulae found in YAML file."
        return 0
    fi
    
    while IFS= read -r brew_package; do
        if [[ -n "$brew_package" ]]; then
            if [[ "$DRY_RUN" == true ]]; then
                print_status "$YELLOW" "DRY RUN: Would install brew: $brew_package"
            else
                print_status "$GREEN" "Installing brew: $brew_package"
                if [[ "$VERBOSE" == true ]]; then
                    brew install "$brew_package"
                else
                    brew install "$brew_package" >/dev/null 2>&1 || print_status "$RED" "Failed to install brew: $brew_package"
                fi
            fi
        fi
    done <<< "$all_brews"
    
    print_status "$GREEN" "Brew formulae installation completed."
}

# Function to install casks
install_casks() {
    local yaml_file=$1
    local casks uncategorized_casks
    
    print_status "$BLUE" "Installing Homebrew casks..."
    
    casks=$(parse_yaml_array "$yaml_file" "casks")
    uncategorized_casks=$(parse_yaml_array "$yaml_file" "uncategorized_casks")
    
    # Combine both regular and uncategorized casks
    local all_casks
    all_casks=$(echo -e "$casks\n$uncategorized_casks" | grep -v '^$' | sort -u)
    
    if [[ -z "$all_casks" ]]; then
        print_status "$YELLOW" "No casks found in YAML file."
        return 0
    fi
    
    while IFS= read -r cask; do
        if [[ -n "$cask" ]]; then
            if [[ "$DRY_RUN" == true ]]; then
                print_status "$YELLOW" "DRY RUN: Would install cask: $cask"
            else
                print_status "$GREEN" "Installing cask: $cask"
                if [[ "$VERBOSE" == true ]]; then
                    brew install --cask "$cask"
                else
                    brew install --cask "$cask" >/dev/null 2>&1 || print_status "$RED" "Failed to install cask: $cask"
                fi
            fi
        fi
    done <<< "$all_casks"
    
    print_status "$GREEN" "Casks installation completed."
}

# Function to install Mac App Store apps
install_mas_apps() {
    local yaml_file=$1
    
    if [[ "$SKIP_MAS" == true ]]; then
        print_status "$YELLOW" "Skipping Mac App Store apps (mas not available)."
        return 0
    fi
    
    print_status "$BLUE" "Installing Mac App Store apps..."
    
    # Parse MAS apps with yq from both sections
    local mas_ids uncategorized_mas_ids
    mas_ids=$(yq eval '.mas_apps[].id' "$yaml_file" 2>/dev/null | grep -v "^null$" || true)
    uncategorized_mas_ids=$(yq eval '.uncategorized_mas_apps[].id' "$yaml_file" 2>/dev/null | grep -v "^null$" || true)
    
    # Combine both regular and uncategorized MAS app IDs
    local all_mas_ids
    all_mas_ids=$(echo -e "$mas_ids\n$uncategorized_mas_ids" | grep -v '^$' | sort -u)
    
    if [[ -z "$all_mas_ids" ]]; then
        print_status "$YELLOW" "No Mac App Store apps found in YAML file."
        return 0
    fi
    
    while IFS= read -r app_id; do
        if [[ -n "$app_id" ]]; then
            local app_name
            # Try to get name from both sections
            app_name=$(yq eval ".mas_apps[] | select(.id == $app_id) | .name" "$yaml_file" 2>/dev/null || true)
            if [[ -z "$app_name" ]]; then
                app_name=$(yq eval ".uncategorized_mas_apps[] | select(.id == $app_id) | .name" "$yaml_file" 2>/dev/null || true)
            fi
            
            if [[ "$DRY_RUN" == true ]]; then
                print_status "$YELLOW" "DRY RUN: Would install MAS app: $app_name (ID: $app_id)"
            else
                print_status "$GREEN" "Installing MAS app: $app_name (ID: $app_id)"
                if [[ "$VERBOSE" == true ]]; then
                    mas install "$app_id"
                else
                    mas install "$app_id" >/dev/null 2>&1 || print_status "$RED" "Failed to install MAS app: $app_name (ID: $app_id)"
                fi
            fi
        fi
    done <<< "$all_mas_ids"
    
    print_status "$GREEN" "Mac App Store apps installation completed."
}

# Function to main installation process
main() {
    local yaml_file="$DEFAULT_YAML_FILE"
    
    # Parse command line arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            -h|--help)
                usage
                exit 0
                ;;
            -v|--verbose)
                VERBOSE=true
                shift
                ;;
            -d|--dry-run)
                DRY_RUN=true
                shift
                ;;
            -t|--taps-only)
                TAPS_ONLY=true
                shift
                ;;
            -b|--brews-only)
                BREWS_ONLY=true
                shift
                ;;
            -c|--casks-only)
                CASKS_ONLY=true
                shift
                ;;
            -m|--mas-only)
                MAS_ONLY=true
                shift
                ;;
            --skip-taps)
                SKIP_TAPS=true
                shift
                ;;
            --skip-brews)
                SKIP_BREWS=true
                shift
                ;;
            --skip-casks)
                SKIP_CASKS=true
                shift
                ;;
            --skip-mas)
                SKIP_MAS=true
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
    
    # Initialize variables
    VERBOSE=${VERBOSE:-false}
    DRY_RUN=${DRY_RUN:-false}
    TAPS_ONLY=${TAPS_ONLY:-false}
    BREWS_ONLY=${BREWS_ONLY:-false}
    CASKS_ONLY=${CASKS_ONLY:-false}
    MAS_ONLY=${MAS_ONLY:-false}
    SKIP_TAPS=${SKIP_TAPS:-false}
    SKIP_BREWS=${SKIP_BREWS:-false}
    SKIP_CASKS=${SKIP_CASKS:-false}
    SKIP_MAS=${SKIP_MAS:-false}
    
    # Check prerequisites
    check_prerequisites
    
    if [[ "$DRY_RUN" == true ]]; then
        print_status "$YELLOW" "=== DRY RUN MODE - No packages will be actually installed ==="
    fi
    
    # Execute based on options
    if [[ "$TAPS_ONLY" == true ]]; then
        install_taps "$yaml_file"
    elif [[ "$BREWS_ONLY" == true ]]; then
        install_brews "$yaml_file"
    elif [[ "$CASKS_ONLY" == true ]]; then
        install_casks "$yaml_file"
    elif [[ "$MAS_ONLY" == true ]]; then
        install_mas_apps "$yaml_file"
    else
        # Install all categories unless specifically skipped
        [[ "$SKIP_TAPS" != true ]] && install_taps "$yaml_file"
        [[ "$SKIP_BREWS" != true ]] && install_brews "$yaml_file"
        [[ "$SKIP_CASKS" != true ]] && install_casks "$yaml_file"
        [[ "$SKIP_MAS" != true ]] && install_mas_apps "$yaml_file"
    fi
    
    print_status "$GREEN" "All installations completed successfully!"
}

# Run main function with all arguments
main "$@" 
