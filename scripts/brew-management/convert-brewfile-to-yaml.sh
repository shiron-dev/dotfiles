#!/bin/bash

# Brewfile to YAML converter script
# This script converts a traditional Brewfile to YAML format

set -euo pipefail

# Colors for output
readonly RED='\033[0;31m'
readonly GREEN='\033[0;32m'
readonly YELLOW='\033[1;33m'
readonly BLUE='\033[0;34m'
readonly NC='\033[0m' # No Color

# Default paths
DEFAULT_BREWFILE="$(dirname "$0")/../data/brew/Brewfile"
DEFAULT_OUTPUT="$(dirname "$0")/../data/brew/packages.yml"

# Function to print colored output
print_status() {
    local color=$1
    local message=$2
    echo -e "${color}${message}${NC}"
}

# Function to print usage
usage() {
    cat << EOF
Usage: $0 [OPTIONS] [BREWFILE] [OUTPUT_FILE]

Convert a Brewfile to YAML format.

Options:
    -h, --help          Show this help message
    -v, --verbose       Enable verbose output
    -f, --force         Overwrite output file if it exists

Arguments:
    BREWFILE           Path to input Brewfile (default: ${DEFAULT_BREWFILE})
    OUTPUT_FILE        Path to output YAML file (default: ${DEFAULT_OUTPUT})

Examples:
    $0                                      # Convert default Brewfile to default YAML file
    $0 my-Brewfile packages.yml             # Convert custom Brewfile to custom YAML file
    $0 --force                              # Overwrite existing output file

EOF
}

# Function to extract taps from Brewfile
extract_taps() {
    local brewfile=$1
    grep '^tap ' "$brewfile" | sed 's/^tap "\(.*\)"/\1/' | sort -u
}

# Function to extract brews from Brewfile
extract_brews() {
    local brewfile=$1
    grep '^brew ' "$brewfile" | sed 's/^brew "\(.*\)"/\1/' | sort -u
}

# Function to extract casks from Brewfile
extract_casks() {
    local brewfile=$1
    grep '^cask ' "$brewfile" | sed 's/^cask "\(.*\)"/\1/' | sort -u
}

# Function to extract MAS apps from Brewfile
extract_mas_apps() {
    local brewfile=$1
    grep '^mas ' "$brewfile" | while IFS= read -r line; do
        # Extract name and id from lines like: mas "App Name", id: 123456789
        if [[ $line =~ mas\ \"([^\"]+)\",\ id:\ ([0-9]+) ]]; then
            local name="${BASH_REMATCH[1]}"
            local id="${BASH_REMATCH[2]}"
            echo "  - name: \"$name\""
            echo "    id: $id"
        fi
    done
}

# Function to convert Brewfile to YAML
convert_brewfile_to_yaml() {
    local brewfile=$1
    local output_file=$2
    
    print_status "$BLUE" "Converting Brewfile to YAML format..."
    print_status "$BLUE" "Input: $brewfile"
    print_status "$BLUE" "Output: $output_file"
    
    # Check if input file exists
    if [[ ! -f "$brewfile" ]]; then
        print_status "$RED" "Error: Brewfile not found: $brewfile"
        exit 1
    fi
    
    # Create output directory if it doesn't exist
    local output_dir
    output_dir=$(dirname "$output_file")
    mkdir -p "$output_dir"
    
    # Start writing YAML file
    cat > "$output_file" << 'EOF'
# YAML-based Brew packages configuration
# This file defines packages to be installed via Homebrew
# Generated from Brewfile

EOF
    
    # Extract and write taps
    local taps
    taps=$(extract_taps "$brewfile")
    if [[ -n "$taps" ]]; then
        echo "taps:" >> "$output_file"
        while IFS= read -r tap; do
            [[ -n "$tap" ]] && echo "  - $tap" >> "$output_file"
        done <<< "$taps"
        echo "" >> "$output_file"
    fi
    
    # Extract and write brews
    local brews
    brews=$(extract_brews "$brewfile")
    if [[ -n "$brews" ]]; then
        echo "# Homebrew formulae (brew install)" >> "$output_file"
        echo "brews:" >> "$output_file"
        while IFS= read -r brew_package; do
            [[ -n "$brew_package" ]] && echo "  - $brew_package" >> "$output_file"
        done <<< "$brews"
        echo "" >> "$output_file"
    fi
    
    # Extract and write casks
    local casks
    casks=$(extract_casks "$brewfile")
    if [[ -n "$casks" ]]; then
        echo "# Homebrew casks (brew install --cask)" >> "$output_file"
        echo "casks:" >> "$output_file"
        while IFS= read -r cask; do
            [[ -n "$cask" ]] && echo "  - $cask" >> "$output_file"
        done <<< "$casks"
        echo "" >> "$output_file"
    fi
    
    # Extract and write MAS apps
    local mas_apps
    mas_apps=$(extract_mas_apps "$brewfile")
    if [[ -n "$mas_apps" ]]; then
        echo "# Mac App Store apps (mas install)" >> "$output_file"
        echo "mas_apps:" >> "$output_file"
        echo "$mas_apps" >> "$output_file"
    fi
    
    print_status "$GREEN" "Conversion completed successfully!"
    print_status "$GREEN" "YAML file created: $output_file"
}

# Main function
main() {
    local brewfile="$DEFAULT_BREWFILE"
    local output_file="$DEFAULT_OUTPUT"
    local force=false
    local verbose=false
    
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
            -f|--force)
                force=true
                shift
                ;;
            -*)
                print_status "$RED" "Unknown option: $1"
                usage
                exit 1
                ;;
            *)
                if [[ "$brewfile" == "$DEFAULT_BREWFILE" ]]; then
                    brewfile="$1"
                elif [[ "$output_file" == "$DEFAULT_OUTPUT" ]]; then
                    output_file="$1"
                else
                    print_status "$RED" "Too many arguments"
                    usage
                    exit 1
                fi
                shift
                ;;
        esac
    done
    
    # Check if output file exists and force is not set
    if [[ -f "$output_file" && "$force" != true ]]; then
        print_status "$RED" "Error: Output file already exists: $output_file"
        print_status "$YELLOW" "Use --force to overwrite"
        exit 1
    fi
    
    # Convert Brewfile to YAML
    convert_brewfile_to_yaml "$brewfile" "$output_file"
    
    if [[ "$verbose" == true ]]; then
        print_status "$BLUE" "Generated YAML content:"
        cat "$output_file"
    fi
}

# Run main function with all arguments
main "$@" 
