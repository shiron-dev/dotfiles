#!/usr/bin/env bash
set -euo pipefail

# update-sops-secrets.sh
# Wrapper script for sops-encrypt that manages file permissions

usage() {
  echo "Usage: $0 <filepath>" >&2
  echo "  Encrypts a secrets file using the project's makefile sops-encrypt target." >&2
  echo "  File must match *.secrets.* pattern (but not *.secrets.*.*)." >&2
  exit 1
}

# Check argument
if [ $# -ne 1 ]; then
  usage
fi

filepath="$1"

# Validate file pattern: *.secrets.* but NOT *.secrets.*.*
filename=$(basename "$filepath")
if [[ ! "$filename" =~ ^.*\.secrets\.[^.]+$ ]]; then
  echo "Error: File '$filename' does not match *.secrets.* pattern (or matches *.secrets.*.* which is not allowed)" >&2
  exit 1
fi

# Check if file exists
if [ ! -f "$filepath" ]; then
  echo "Error: File '$filepath' not found" >&2
  exit 1
fi

# Get absolute path
filepath=$(cd "$(dirname "$filepath")" && pwd)/$(basename "$filepath")

# Check if we're in a git repository and get project root
if ! project_root=$(git -C "$(dirname "$filepath")" rev-parse --show-toplevel 2>/dev/null); then
  echo "Error: '$filepath' is not in a git repository" >&2
  exit 1
fi

# Check for makefile with sops-encrypt target
makefile_path=""
if [ -f "$project_root/makefile" ]; then
  makefile_path="$project_root/makefile"
elif [ -f "$project_root/Makefile" ]; then
  makefile_path="$project_root/Makefile"
else
  echo "Error: No makefile or Makefile found in project root '$project_root'" >&2
  exit 1
fi

if ! grep -q '^sops-encrypt:' "$makefile_path"; then
  echo "Error: sops-encrypt target not found in '$makefile_path'" >&2
  exit 1
fi

# Check if sops-encrypt target supports FILE parameter
supports_file_param=false
if grep -q '\$(FILE)' "$makefile_path"; then
  supports_file_param=true
fi

echo "[INFO] Starting sops-encrypt for: $filepath"
echo "[INFO] FILE parameter support: $supports_file_param"

# Define the sops file path
sops_filepath="${filepath}.sops"

# Function to get file permissions (macOS compatible)
get_permissions() {
  local file="$1"
  if [ -f "$file" ]; then
    stat -f '%p' "$file"
  else
    echo ""
  fi
}

# Function to restore permissions
restore_permissions() {
  local file="$1"
  local perms="$2"
  if [ -n "$perms" ] && [ -f "$file" ]; then
    # Extract the last 4 digits (octal permissions)
    local octal_perms="${perms: -4}"
    chmod "$octal_perms" "$file"
  fi
}

# Save original permissions
original_perms=$(get_permissions "$filepath")
original_sops_perms=$(get_permissions "$sops_filepath")

echo "[INFO] Saved permissions - secrets file: $original_perms, sops file: ${original_sops_perms:-N/A}"

# Make files writable
chmod +w "$filepath"
if [ -f "$sops_filepath" ]; then
  chmod +w "$sops_filepath"
fi

echo "[INFO] Made files writable"

# Run sops-encrypt
exit_code=0
if [ "$supports_file_param" = true ]; then
  if ! make -C "$project_root" sops-encrypt FILE="$filepath"; then
    exit_code=$?
    echo "[ERROR] sops-encrypt failed with exit code $exit_code" >&2
  fi
else
  if ! make -C "$project_root" sops-encrypt; then
    exit_code=$?
    echo "[ERROR] sops-encrypt failed with exit code $exit_code" >&2
  fi
fi

# Restore original permissions
echo "[INFO] Restoring original permissions..."
restore_permissions "$filepath" "$original_perms"
restore_permissions "$sops_filepath" "$original_sops_perms"

# If sops file was newly created, set it to read-only
if [ -z "$original_sops_perms" ] && [ -f "$sops_filepath" ]; then
  chmod -w "$sops_filepath"
  echo "[INFO] New sops file created, set to read-only"
fi

echo "[INFO] Done"
exit $exit_code
