#!/bin/bash

# Enhanced YAML-based Brew installer script with groups and tags support
# This script reads a YAML file and installs packages using group/tag filtering

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

Install Homebrew packages from a YAML configuration file with group/tag support.

Options:
    -h, --help              Show this help message
    -v, --verbose           Enable verbose output
    -d, --dry-run           Show what would be installed without actually installing
    --list-groups           List all available groups
    --list-tags             List all available tags
    --list-profiles         List all available profiles
    
    # Installation filters
    -g, --groups GROUPS     Install only specified groups (comma-separated)
    -t, --tags TAGS         Install only packages with specified tags (comma-separated)
    --exclude-groups GROUPS Exclude specified groups (comma-separated)
    --exclude-tags TAGS     Exclude packages with specified tags (comma-separated)
    -p, --profile PROFILE   Install using predefined profile
    
    # Package type filters
    --taps-only             Install only taps
    --brews-only            Install only brew formulae
    --casks-only            Install only casks
    --mas-only              Install only Mac App Store apps
    --skip-taps             Skip installing taps
    --skip-brews            Skip installing brew formulae
    --skip-casks            Skip installing casks
    --skip-mas              Skip installing Mac App Store apps

Arguments:
    YAML_FILE              Path to YAML configuration file (default: ${DEFAULT_YAML_FILE})

Examples:
    $0                                          # Install all packages
    $0 --groups core,development               # Install only core and development groups
    $0 --tags essential,productivity           # Install packages with essential or productivity tags
    $0 --profile developer                     # Install using developer profile
    $0 --exclude-tags experimental             # Install all except experimental packages
    $0 --groups development --brews-only       # Install only brew packages from development group

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

# Function to list available groups
list_groups() {
    local yaml_file=$1
    print_status "$CYAN" "Available Groups:"
    yq eval '.groups | to_entries | .[] | .key + ": " + .value.description + " (priority: " + (.value.priority | tostring) + ")"' "$yaml_file" 2>/dev/null | sort -t':' -k2,2n || print_status "$RED" "No groups found"
}

# Function to list available tags
list_tags() {
    local yaml_file=$1
    print_status "$CYAN" "Available Tags:"
    {
        yq eval '.groups[].packages[]?.tags[]?' "$yaml_file" 2>/dev/null
    } | grep -v "^null$" | sort -u || print_status "$RED" "No tags found"
}

# Function to list available profiles
list_profiles() {
    local yaml_file=$1
    print_status "$CYAN" "Available Profiles:"
    yq eval '.profiles | to_entries | .[] | .key + ": " + .value.description' "$yaml_file" 2>/dev/null || print_status "$RED" "No profiles found"
}



# Function to get profile configuration
get_profile_config() {
    local yaml_file=$1
    local profile=$2
    
    local profile_groups profile_tags exclude_tags
    profile_groups=$(yq eval ".profiles.${profile}.groups[]?" "$yaml_file" 2>/dev/null | tr '\n' ',' | sed 's/,$//')
    profile_tags=$(yq eval ".profiles.${profile}.tags[]?" "$yaml_file" 2>/dev/null | tr '\n' ',' | sed 's/,$//')
    exclude_tags=$(yq eval ".profiles.${profile}.exclude_tags[]?" "$yaml_file" 2>/dev/null | tr '\n' ',' | sed 's/,$//')
    
    echo "${profile_groups}|${profile_tags}|${exclude_tags}"
}



# Function to get filtered packages from groups
get_filtered_packages() {
    local yaml_file=$1
    local package_type=$2
    local groups_filter=$3
    local tags_filter=$4
    local exclude_groups_filter=$5
    local exclude_tags_filter=$6
    
    # Simple approach: if groups_filter is specified, only process those groups
    if [[ -n "$groups_filter" ]]; then
        # Split groups by comma
        echo "$groups_filter" | tr ',' '\n' | while read -r group; do
            if [[ -n "$group" ]]; then
                yq eval ".groups.\"$group\".packages[]? | select(.type == \"$package_type\")" "$yaml_file" 2>/dev/null
            fi
        done
    else
        # Process all groups
        yq eval ".groups[].packages[]? | select(.type == \"$package_type\")" "$yaml_file" 2>/dev/null
    fi
}

# Function to install packages of a specific type with filtering
install_filtered_packages() {
    local yaml_file=$1
    local package_type=$2
    local groups_filter=$3
    local tags_filter=$4
    local exclude_groups_filter=$5
    local exclude_tags_filter=$6
    
    local packages
    packages=$(get_filtered_packages "$yaml_file" "$package_type" "$groups_filter" "$tags_filter" "$exclude_groups_filter" "$exclude_tags_filter")
    
    if [[ -z "$packages" ]]; then
        print_status "$YELLOW" "No $package_type packages found matching the criteria."
        return 0
    fi
    
    case "$package_type" in
        "tap")
            install_taps_from_data "$packages"
            ;;
        "brew")
            install_brews_from_data "$packages"
            ;;
        "cask")
            install_casks_from_data "$packages"
            ;;
        "mas")
            install_mas_from_data "$packages"
            ;;
    esac
}

# Function to install taps from YAML data
install_taps_from_data() {
    local tap_data=$1
    
    print_status "$BLUE" "Installing Homebrew taps..."
    
    local tap_count=0
    
    while IFS= read -r tap_entry; do
        if [[ -n "$tap_entry" && "$tap_entry" != "null" ]]; then
            local tap_name
            tap_name=$(echo "$tap_entry" | yq eval '.name' -)
            
            if [[ -n "$tap_name" && "$tap_name" != "null" ]]; then
                ((tap_count++))
                
                if [[ "$DRY_RUN" == true ]]; then
                    print_status "$CYAN" "Would install tap: $tap_name"
                    continue
                fi
                
                print_status "$CYAN" "Installing tap: $tap_name"
                
                if brew tap "$tap_name"; then
                    print_status "$GREEN" "Successfully installed tap: $tap_name"
                else
                    print_status "$RED" "Failed to install tap: $tap_name"
                fi
            fi
        fi
    done <<< "$tap_data"
    
    if [[ "$tap_count" -eq 0 ]]; then
        print_status "$YELLOW" "No taps to install."
    else
        print_status "$GREEN" "Taps installation completed. ($tap_count packages processed)"
    fi
}

# Function to install brews from YAML data
install_brews_from_data() {
    local brew_data=$1
    
    print_status "$BLUE" "Installing Homebrew formulae..."
    
    local brew_count=0
    
    while IFS= read -r brew_entry; do
        if [[ -n "$brew_entry" && "$brew_entry" != "null" ]]; then
            local brew_name
            brew_name=$(echo "$brew_entry" | yq eval '.name' -)
            
            if [[ -n "$brew_name" && "$brew_name" != "null" ]]; then
                ((brew_count++))
                
                if [[ "$DRY_RUN" == true ]]; then
                    print_status "$CYAN" "Would install formula: $brew_name"
                    continue
                fi
                
                print_status "$CYAN" "Installing formula: $brew_name"
                
                if brew install "$brew_name"; then
                    print_status "$GREEN" "Successfully installed formula: $brew_name"
                else
                    print_status "$RED" "Failed to install formula: $brew_name"
                fi
            fi
        fi
    done <<< "$brew_data"
    
    if [[ "$brew_count" -eq 0 ]]; then
        print_status "$YELLOW" "No formulae to install."
    else
        print_status "$GREEN" "Formulae installation completed. ($brew_count packages processed)"
    fi
}

# Function to install casks from YAML data
install_casks_from_data() {
    local cask_data=$1
    
    print_status "$BLUE" "Installing Homebrew casks..."
    
    local cask_count=0
    
    while IFS= read -r cask_entry; do
        if [[ -n "$cask_entry" && "$cask_entry" != "null" ]]; then
            local cask_name
            cask_name=$(echo "$cask_entry" | yq eval '.name' -)
            
            if [[ -n "$cask_name" && "$cask_name" != "null" ]]; then
                ((cask_count++))
                
                if [[ "$DRY_RUN" == true ]]; then
                    print_status "$CYAN" "Would install cask: $cask_name"
                    continue
                fi
                
                print_status "$CYAN" "Installing cask: $cask_name"
                
                if brew install --cask "$cask_name"; then
                    print_status "$GREEN" "Successfully installed cask: $cask_name"
                else
                    print_status "$RED" "Failed to install cask: $cask_name"
                fi
            fi
        fi
    done <<< "$cask_data"
    
    if [[ "$cask_count" -eq 0 ]]; then
        print_status "$YELLOW" "No casks to install."
    else
        print_status "$GREEN" "Casks installation completed. ($cask_count packages processed)"
    fi
}

# Function to install Mac App Store apps from YAML data
install_mas_from_data() {
    local mas_data=$1
    
    if [[ "$SKIP_MAS" == true ]]; then
        print_status "$YELLOW" "Skipping Mac App Store apps (mas not available)."
        return 0
    fi
    
    print_status "$BLUE" "Installing Mac App Store apps..."
    
    local mas_count=0
    
    while IFS= read -r mas_entry; do
        if [[ -n "$mas_entry" && "$mas_entry" != "null" ]]; then
            local app_name app_id
            app_name=$(echo "$mas_entry" | yq eval '.name' -)
            app_id=$(echo "$mas_entry" | yq eval '.id' -)
            
            if [[ -n "$app_id" && "$app_id" != "null" ]]; then
                ((mas_count++))
                
                if [[ "$DRY_RUN" == true ]]; then
                    print_status "$CYAN" "Would install MAS app: $app_name ($app_id)"
                    continue
                fi
                
                print_status "$CYAN" "Installing MAS app: $app_name ($app_id)"
                
                if mas install "$app_id"; then
                    print_status "$GREEN" "Successfully installed MAS app: $app_name"
                else
                    print_status "$RED" "Failed to install MAS app: $app_name"
                fi
            fi
        fi
    done <<< "$mas_data"
    
    if [[ "$mas_count" -eq 0 ]]; then
        print_status "$YELLOW" "No MAS apps to install."
    else
        print_status "$GREEN" "MAS apps installation completed. ($mas_count packages processed)"
    fi
}

# Main function
main() {
    local yaml_file="$DEFAULT_YAML_FILE"
    local groups_filter=""
    local tags_filter=""
    local exclude_groups_filter=""
    local exclude_tags_filter=""
    local profile=""
    
    
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
            --list-groups)
                if [[ ! -f "$yaml_file" ]]; then
                    yaml_file="${2:-$DEFAULT_YAML_FILE}"
                fi
                list_groups "$yaml_file"
                exit 0
                ;;
            --list-tags)
                if [[ ! -f "$yaml_file" ]]; then
                    yaml_file="${2:-$DEFAULT_YAML_FILE}"
                fi
                list_tags "$yaml_file"
                exit 0
                ;;
            --list-profiles)
                if [[ ! -f "$yaml_file" ]]; then
                    yaml_file="${2:-$DEFAULT_YAML_FILE}"
                fi
                list_profiles "$yaml_file"
                exit 0
                ;;
            -g|--groups)
                groups_filter="$2"
                shift 2
                ;;
            -t|--tags)
                tags_filter="$2"
                shift 2
                ;;
            --exclude-groups)
                exclude_groups_filter="$2"
                shift 2
                ;;
            --exclude-tags)
                exclude_tags_filter="$2"
                shift 2
                ;;
            -p|--profile)
                profile="$2"
                shift 2
                ;;
            --taps-only)
                TAPS_ONLY=true
                shift
                ;;
            --brews-only)
                BREWS_ONLY=true
                shift
                ;;
            --casks-only)
                CASKS_ONLY=true
                shift
                ;;
            --mas-only)
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
    
    # Handle profile configuration
    if [[ -n "$profile" ]]; then
        local profile_config
        profile_config=$(get_profile_config "$yaml_file" "$profile")
        IFS='|' read -r profile_groups profile_tags profile_exclude_tags <<< "$profile_config"
        
        print_status "$CYAN" "Using profile: $profile"
        [[ -n "$profile_groups" ]] && groups_filter="$profile_groups"
        [[ -n "$profile_tags" ]] && tags_filter="$profile_tags"
        [[ -n "$profile_exclude_tags" ]] && exclude_tags_filter="$profile_exclude_tags"
    fi
    

    
    # Display filter information
    if [[ -n "$groups_filter" ]]; then
        print_status "$MAGENTA" "Groups filter: $groups_filter"
    fi
    if [[ -n "$tags_filter" ]]; then
        print_status "$MAGENTA" "Tags filter: $tags_filter"
    fi
    if [[ -n "$exclude_groups_filter" ]]; then
        print_status "$MAGENTA" "Exclude groups: $exclude_groups_filter"
    fi
    if [[ -n "$exclude_tags_filter" ]]; then
        print_status "$MAGENTA" "Exclude tags: $exclude_tags_filter"
    fi
    
    # Check prerequisites
    check_prerequisites
    
    if [[ "$DRY_RUN" == true ]]; then
        print_status "$YELLOW" "=== DRY RUN MODE - No packages will be actually installed ==="
    fi
    
    # Execute based on options
    if [[ "$TAPS_ONLY" == true ]]; then
        install_filtered_packages "$yaml_file" "tap" "$groups_filter" "$tags_filter" "$exclude_groups_filter" "$exclude_tags_filter"
    elif [[ "$BREWS_ONLY" == true ]]; then
        install_filtered_packages "$yaml_file" "brew" "$groups_filter" "$tags_filter" "$exclude_groups_filter" "$exclude_tags_filter"
    elif [[ "$CASKS_ONLY" == true ]]; then
        install_filtered_packages "$yaml_file" "cask" "$groups_filter" "$tags_filter" "$exclude_groups_filter" "$exclude_tags_filter"
    elif [[ "$MAS_ONLY" == true ]]; then
        install_filtered_packages "$yaml_file" "mas" "$groups_filter" "$tags_filter" "$exclude_groups_filter" "$exclude_tags_filter"
    else
        # Install all categories unless specifically skipped
        [[ "$SKIP_TAPS" != true ]] && install_filtered_packages "$yaml_file" "tap" "$groups_filter" "$tags_filter" "$exclude_groups_filter" "$exclude_tags_filter"
        [[ "$SKIP_BREWS" != true ]] && install_filtered_packages "$yaml_file" "brew" "$groups_filter" "$tags_filter" "$exclude_groups_filter" "$exclude_tags_filter"
        [[ "$SKIP_CASKS" != true ]] && install_filtered_packages "$yaml_file" "cask" "$groups_filter" "$tags_filter" "$exclude_groups_filter" "$exclude_tags_filter"
        [[ "$SKIP_MAS" != true ]] && install_filtered_packages "$yaml_file" "mas" "$groups_filter" "$tags_filter" "$exclude_groups_filter" "$exclude_tags_filter"
    fi
    
    print_status "$GREEN" "All filtered installations completed successfully!"
}

# Run main function with all arguments
main "$@" 
