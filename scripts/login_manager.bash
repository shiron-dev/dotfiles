#!/bin/bash
set -euo pipefail

script_dir=$(cd "$(dirname "$0")" && pwd)
data_dir="${script_dir}/../data"
yaml_file="${data_dir}/login.yaml"


_get_current_items_as_yaml() {
  local osascript_cmd='
    tell application "System Events"
      set output to ""
      repeat with li in login items
        if path of li is not missing value then
          set output to output & (path of li) & "\t" & (hidden of li) & "\n"
        end if
      end repeat
      return output
    end tell'

  while IFS=$'\t' read -r path_val hidden_val; do
    if [[ -n "$path_val" ]]; then
      printf -- "- path: %s\n  hidden: %s\n" "$path_val" "$hidden_val"
    fi
  done < <(osascript -e "$osascript_cmd" 2>/dev/null)
}

export_login_items() {
  echo "🚀 Exporting current login items to $yaml_file..."
  mkdir -p "$data_dir"
  _get_current_items_as_yaml > "$yaml_file"
  echo "✅ Export complete."
}

check_login_items() {
  echo "🔎 Checking for differences against $yaml_file..."
  if [[ ! -f "$yaml_file" ]]; then
    echo "🚨 Error: YAML file not found. Please run 'export' first." >&2
    exit 1
  fi

  local temp_yaml
  temp_yaml=$(mktemp)
  _get_current_items_as_yaml > "$temp_yaml"

  echo "---"
  echo "■ Difference (login.yaml <-> current state):"
  diff -u "$yaml_file" "$temp_yaml" || true
  echo "---"

  rm "$temp_yaml"
  echo "✅ Check complete."
}

import_login_items() {
  echo "📥 Importing login items from $yaml_file..."
  if [[ ! -f "$yaml_file" ]]; then
    echo "🚨 Error: YAML file not found. Please run 'export' first." >&2
    exit 1
  fi

  echo "  - Clearing all current login items..."
  osascript -e 'tell application "System Events" to delete every login item'

  echo "  - Adding items defined in YAML file..."
  local path=""
  local hidden=""
  while IFS= read -r line || [[ -n "$line" ]]; do
    if [[ "$line" == *"- path: "* ]]; then
      path="${line#*- path: }"
    elif [[ "$line" == *"hidden: "* ]]; then
      hidden="${line#*hidden: }"
      if [[ -n "$path" && -n "$hidden" ]]; then
        echo "    - Adding: $path (hidden: $hidden)"
        osascript -e "tell application \"System Events\" to make new login item at end with properties {path:\"$path\", hidden:$hidden}"
        path=""
        hidden=""
      fi
    fi
  done < "$yaml_file"

  echo "✅ Import complete. System login items have been synced."
}

command="${1:-}"

case "$command" in
  export)
    export_login_items
    ;;
  check)
    check_login_items
    ;;
  import)
    import_login_items
    ;;
  *)
    echo "Usage: $0 {export|check|import}"
    echo
    echo "  export: Save current login items to ${yaml_file}"
    echo "  check:  Compare current login items with the YAML file"
    echo "  import: Sync login items based on the YAML file (deletes all current items first)"
    exit 1
    ;;
esac
