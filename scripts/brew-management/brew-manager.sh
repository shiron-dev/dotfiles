#!/bin/bash

# Brew Management Unified Script
# This script provides a unified interface for all brew management operations

set -euo pipefail

# Colors for output
readonly RED='\033[0;31m'
readonly GREEN='\033[0;32m'
readonly YELLOW='\033[1;33m'
readonly BLUE='\033[0;34m'
readonly CYAN='\033[0;36m'
readonly NC='\033[0m' # No Color

# Script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Function to print colored output
print_status() {
    local color=$1
    local message=$2
    echo -e "${color}${message}${NC}"
}

# Function to print usage
usage() {
    cat << EOF
Usage: $0 <command> [options]

Unified brew package management tool with YAML configuration support.

Commands:
    install             Install packages from YAML configuration (with groups/tags support)
    install-simple      Install packages from simple YAML configuration  
    sync                Sync installed packages to YAML configuration (with groups/tags)
    sync-simple         Sync installed packages to simple YAML configuration
    convert             Convert Brewfile to YAML format
    validate            Validate YAML configuration files against their schemas
    list-groups         List available groups in grouped YAML
    list-tags           List available tags in grouped YAML
    list-profiles       List available profiles in grouped YAML

Global Options:
    -h, --help          Show this help message
    -v, --verbose       Enable verbose output
    -d, --dry-run       Show what would be done without actually doing it

Examples:
    $0 install --groups core,development           # Install core and development groups
    $0 install --profile developer                 # Install using developer profile
    $0 install-simple                              # Install from simple YAML format
    $0 sync --auto-detect                          # Sync with auto group/tag detection
    $0 sync-simple                                 # Sync to simple YAML format
    $0 convert Brewfile packages.yml               # Convert Brewfile to YAML
    $0 list-groups                                 # List all available groups

Use '$0 <command> --help' for more information on a specific command.

EOF
}

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to execute subcommands
execute_subcommand() {
    local command=$1
    shift
    
    case $command in
        install)
            exec "$SCRIPT_DIR/install-brew-grouped.sh" "$@"
            ;;
        install-simple)
            exec "$SCRIPT_DIR/install-brew-from-yaml.sh" "$@"
            ;;
        sync)
            exec "$SCRIPT_DIR/sync-brew-grouped.sh" "$@"
            ;;
        sync-simple)
            exec "$SCRIPT_DIR/sync-brew-to-yaml.sh" "$@"
            ;;
        convert)
            exec "$SCRIPT_DIR/convert-brewfile-to-yaml.sh" "$@"
            ;;
        validate)
            exec "$SCRIPT_DIR/validate-yaml.sh" "$@"
            ;;
        list-groups)
            exec "$SCRIPT_DIR/install-brew-grouped.sh" --list-groups "$@"
            ;;
        list-tags)
            exec "$SCRIPT_DIR/install-brew-grouped.sh" --list-tags "$@"
            ;;
        list-profiles)
            exec "$SCRIPT_DIR/install-brew-grouped.sh" --list-profiles "$@"
            ;;
        *)
            print_status "$RED" "Error: Unknown command '$command'"
            echo
            usage
            exit 1
            ;;
    esac
}

# Main function
main() {
    # Check for no arguments
    if [[ $# -eq 0 ]]; then
        usage
        exit 1
    fi
    
    # Handle global help
    case $1 in
        -h|--help|help)
            usage
            exit 0
            ;;
    esac
    
    # Check for required tools
    if ! command_exists brew; then
        print_status "$RED" "Error: Homebrew is not installed. Please install Homebrew first."
        exit 1
    fi
    
    if ! command_exists yq; then
        print_status "$YELLOW" "Warning: yq is not installed. Installing yq for YAML parsing..."
        brew install yq
    fi
    
    # Get command and execute
    local command=$1
    shift
    
    execute_subcommand "$command" "$@"
}

# Run main function with all arguments
main "$@" 
