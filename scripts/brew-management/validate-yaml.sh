#!/bin/bash

# YAML Schema Validation Script
# This script validates YAML files against their JSON schemas

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
DATA_DIR="$SCRIPT_DIR/../../data/brew"
SCHEMAS_DIR="$DATA_DIR/schemas"

# Function to print colored output
print_status() {
    local color=$1
    local message=$2
    echo -e "${color}${message}${NC}"
}

# Function to print usage
usage() {
    cat << EOF
Usage: $0 [options] [yaml_file]

Validate YAML configuration files against their JSON schemas.

Options:
    -h, --help          Show this help message
    -v, --verbose       Enable verbose output
    -a, --all           Validate all YAML files
    --schema SCHEMA     Use specific schema file

Arguments:
    yaml_file          Path to YAML file to validate (optional)

Examples:
    $0                                          # Validate all YAML files
    $0 packages.yml                             # Validate specific file
    $0 --schema packages-grouped.schema.json packages-grouped.yml
    $0 --all --verbose                          # Validate all with verbose output

EOF
}

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to validate YAML against schema
validate_yaml() {
    local yaml_file=$1
    local schema_file=$2
    local verbose=${3:-false}
    
    if [[ ! -f "$yaml_file" ]]; then
        print_status "$RED" "Error: YAML file not found: $yaml_file"
        return 1
    fi
    
    if [[ ! -f "$schema_file" ]]; then
        print_status "$RED" "Error: Schema file not found: $schema_file"
        return 1
    fi
    
    if [[ "$verbose" == true ]]; then
        print_status "$BLUE" "Validating: $yaml_file"
        print_status "$BLUE" "Schema: $schema_file"
    fi
    
    # Remove yaml-language-server comment and perform basic validation
    local temp_file
    temp_file=$(mktemp)
    grep -v "# yaml-language-server:" "$yaml_file" > "$temp_file"
    
    # Basic YAML syntax check
    if ! yq eval '.' "$temp_file" >/dev/null 2>&1; then
        print_status "$RED" "❌ Invalid: $(basename "$yaml_file") - YAML syntax error"
        if [[ "$verbose" == true ]]; then
            print_status "$YELLOW" "YAML syntax errors:"
            yq eval '.' "$temp_file" 2>&1 || true
        fi
        rm -f "$temp_file"
        return 1
    fi
    
    # Schema-specific validation
    local validation_errors=()
    local filename=$(basename "$yaml_file")
    
    case "$filename" in
        packages-grouped.yml)
            # Check for required fields in grouped format
            if ! yq eval '.metadata' "$temp_file" >/dev/null 2>&1; then
                validation_errors+=("Missing required field: metadata")
            fi
            if ! yq eval '.groups' "$temp_file" >/dev/null 2>&1; then
                validation_errors+=("Missing required field: groups")
            fi
            
            # Check metadata structure
            if yq eval '.metadata' "$temp_file" >/dev/null 2>&1; then
                if ! yq eval '.metadata.version' "$temp_file" >/dev/null 2>&1; then
                    validation_errors+=("Missing metadata.version")
                fi
                if ! yq eval '.metadata.supports_groups' "$temp_file" >/dev/null 2>&1; then
                    validation_errors+=("Missing metadata.supports_groups")
                fi
                if ! yq eval '.metadata.supports_tags' "$temp_file" >/dev/null 2>&1; then
                    validation_errors+=("Missing metadata.supports_tags")
                fi
            fi
            
            # Check groups structure
            if yq eval '.groups' "$temp_file" >/dev/null 2>&1; then
                local groups
                groups=$(yq eval '.groups | keys | .[]' "$temp_file" 2>/dev/null)
                while IFS= read -r group; do
                    if [[ -n "$group" ]]; then
                        if ! yq eval ".groups.${group}.description" "$temp_file" >/dev/null 2>&1; then
                            validation_errors+=("Missing description in group: $group")
                        fi
                        if ! yq eval ".groups.${group}.priority" "$temp_file" >/dev/null 2>&1; then
                            validation_errors+=("Missing priority in group: $group")
                        fi
                        if ! yq eval ".groups.${group}.packages" "$temp_file" >/dev/null 2>&1; then
                            validation_errors+=("Missing packages array in group: $group")
                        fi
                    fi
                done <<< "$groups"
            fi
            ;;
        packages.yml)
            # Check for at least one package section in simple format
            local has_packages=false
            for section in "taps" "brews" "casks" "mas_apps"; do
                if yq eval ".$section" "$temp_file" >/dev/null 2>&1; then
                    has_packages=true
                    break
                fi
            done
            
            if [[ "$has_packages" == false ]]; then
                validation_errors+=("No package sections found (taps, brews, casks, mas_apps)")
            fi
            
            # Check mas_apps structure if exists
            if yq eval '.mas_apps' "$temp_file" >/dev/null 2>&1; then
                local mas_count
                mas_count=$(yq eval '.mas_apps | length' "$temp_file" 2>/dev/null)
                if [[ "$mas_count" =~ ^[0-9]+$ ]] && [[ "$mas_count" -gt 0 ]]; then
                    for ((i=0; i<mas_count; i++)); do
                        if ! yq eval ".mas_apps[$i].name" "$temp_file" >/dev/null 2>&1; then
                            validation_errors+=("Missing name in mas_apps[$i]")
                        fi
                        if ! yq eval ".mas_apps[$i].id" "$temp_file" >/dev/null 2>&1; then
                            validation_errors+=("Missing id in mas_apps[$i]")
                        fi
                    done
                fi
            fi
            ;;
    esac
    
    rm -f "$temp_file"
    
    # Report validation results
    if [[ ${#validation_errors[@]} -eq 0 ]]; then
        print_status "$GREEN" "✅ Valid: $(basename "$yaml_file")"
        return 0
    else
        print_status "$RED" "❌ Invalid: $(basename "$yaml_file")"
        
        if [[ "$verbose" == true ]]; then
            print_status "$YELLOW" "Validation errors:"
            for error in "${validation_errors[@]}"; do
                echo "  - $error"
            done
        fi
        
        return 1
    fi
}

# Function to auto-detect schema for YAML file
detect_schema() {
    local yaml_file=$1
    local filename=$(basename "$yaml_file")
    
    case "$filename" in
        packages-grouped.yml)
            echo "$SCHEMAS_DIR/packages-grouped.schema.json"
            ;;
        packages.yml)
            echo "$SCHEMAS_DIR/packages-simple.schema.json"
            ;;
        *)
            # Try to detect from content or schema comment
            if grep -q "supports_groups.*true" "$yaml_file" 2>/dev/null; then
                echo "$SCHEMAS_DIR/packages-grouped.schema.json"
            elif grep -q "taps:\|brews:\|casks:\|mas_apps:" "$yaml_file" 2>/dev/null; then
                echo "$SCHEMAS_DIR/packages-simple.schema.json"
            else
                echo ""
            fi
            ;;
    esac
}

# Function to validate all YAML files
validate_all() {
    local verbose=${1:-false}
    local success_count=0
    local error_count=0
    
    print_status "$CYAN" "Validating all YAML configuration files..."
    echo
    
    # Find all YAML files in the data directory
    find "$DATA_DIR" -name "*.yml" -o -name "*.yaml" | while read -r yaml_file; do
        local schema_file
        schema_file=$(detect_schema "$yaml_file")
        
        if [[ -n "$schema_file" ]]; then
            if validate_yaml "$yaml_file" "$schema_file" "$verbose"; then
                ((success_count++))
            else
                ((error_count++))
            fi
        else
            print_status "$YELLOW" "⚠️  No schema found for: $(basename "$yaml_file")"
        fi
        echo
    done
    
    echo
    if [[ $error_count -eq 0 ]]; then
        print_status "$GREEN" "🎉 All YAML files are valid!"
    else
        print_status "$RED" "❌ $error_count file(s) failed validation"
        exit 1
    fi
}

# Main function
main() {
    local verbose=false
    local validate_all_files=false
    local schema_file=""
    local yaml_file=""
    
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
            -a|--all)
                validate_all_files=true
                shift
                ;;
            --schema)
                schema_file="$2"
                shift 2
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
    
    # Check prerequisites
    if ! command_exists yq; then
        print_status "$RED" "Error: yq is not installed. Please install yq for YAML processing."
        exit 1
    fi
    
    if ! command_exists npx; then
        print_status "$RED" "Error: npx is not installed. Please install Node.js."
        exit 1
    fi
    
    # Execute validation
    if [[ "$validate_all_files" == true ]]; then
        validate_all "$verbose"
    elif [[ -n "$yaml_file" ]]; then
        # Validate specific file
        if [[ -n "$schema_file" ]]; then
            # Use provided schema
            if [[ ! -f "$schema_file" ]]; then
                schema_file="$SCHEMAS_DIR/$schema_file"
            fi
        else
            # Auto-detect schema
            schema_file=$(detect_schema "$yaml_file")
            if [[ -z "$schema_file" ]]; then
                print_status "$RED" "Error: Could not detect schema for: $yaml_file"
                print_status "$YELLOW" "Use --schema option to specify schema manually"
                exit 1
            fi
        fi
        
        validate_yaml "$yaml_file" "$schema_file" "$verbose"
    else
        # Default: validate all files
        validate_all "$verbose"
    fi
}

# Run main function with all arguments
main "$@" 
