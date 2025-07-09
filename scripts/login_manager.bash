#!/bin/bash
set -euo pipefail

script_dir=$(cd "$(dirname "$0")" && pwd)
data_dir="${script_dir}/../data"
yaml_file="${data_dir}/login.yaml"

# --- 内部関数 ---
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

  if ! osascript -e "$osascript_cmd" 2> /dev/null | while IFS=$'\t' read -r path_val hidden_val; do
      if [[ -n "$path_val" ]]; then
        path_val="${path_val%$'\r'}"
        hidden_val="${hidden_val%$'\r'}"
        printf -- "- path: %s\n  hidden: %s\n" "$path_val" "$hidden_val"
      fi
    done; then
    echo "🚨 Error: osascript failed to get login items." >&2
    return 1
  fi
}

# --- コマンド関数 ---

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

  ( 
    local temp_yaml
    temp_yaml=$(mktemp)
    trap 'rm -f "$temp_yaml"' EXIT
    _get_current_items_as_yaml > "$temp_yaml"

    echo "---"
    echo "■ Difference (login.yaml <-> current state):"
    diff -u "$yaml_file" "$temp_yaml" || true
    echo "---"

    echo "✅ Check complete."
  )
}

import_login_items() {
  echo "📥 Importing login items from $yaml_file..."
  if [[ ! -f "$yaml_file" ]]; then
    echo "🚨 Error: YAML file not found. Please run 'export' first." >&2
    exit 1
  fi

  echo "  - Clearing all current login items..."
  if ! osascript -e 'tell application "System Events" to delete every login item'; then
    echo "🚨 Error: Failed to clear login items." >&2
    exit 1
  fi

  echo "  - Adding items defined in YAML file..."
  local path=""
  local hidden=""

  while IFS= read -r line || [[ -n "$line" ]]; do
    if [[ "$line" == "- path: "* ]]; then
      path="${line#*- path: }"
      path="${path%"${path##*[![:space:]]}"}"
      path="${path%$'\r'}"
    
    elif [[ "$line" == *"hidden: "* ]]; then
      hidden="${line#*hidden: }"
      hidden="${hidden#"${hidden%%[![:space:]]*}"}"
      hidden="${hidden%$'\r'}"
      
      if [[ -n "$path" && -n "$hidden" ]]; then
        echo "    - Adding: $path (hidden: $hidden)"
        local path_escaped="${path//\"/\\\"}"
        if ! osascript -e "tell application \"System Events\" to make new login item at end with properties {path:\"$path_escaped\", hidden:$hidden}"; then
          echo "🚨 Error: Failed to add login item: $path" >&2
          exit 1
        fi
        path=""
        hidden=""
      fi
    fi
  done < "$yaml_file"

  echo "✅ Import complete. System login items have been synced."
}

# ✨ 新しく追加した機能 ✨
open_login_items() {
  echo "📂 This will open all applications defined in $yaml_file."
  if [[ ! -f "$yaml_file" ]]; then
    echo "🚨 Error: YAML file not found. Please run 'export' first." >&2
    exit 1
  fi

  # ユーザーに確認を求める (-n 1 で1文字だけ読み込む)
  read -p "Are you sure you want to open all apps? (y/N) " -n 1 -r
  echo # 改行して表示を整える

  # 入力が 'y' または 'Y' の場合のみ処理を続行
  if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo "  - Opening applications..."
    local path=""
    while IFS= read -r line || [[ -n "$line" ]]; do
      # 'path:' を含む行からパスを抽出
      if [[ "$line" == "- path: "* ]]; then
        path="${line#*- path: }"
        path="${path%"${path##*[![:space:]]}"}" # 末尾の空白を削除
        path="${path%$'\r'}"                  # 末尾の改行コードを削除

        if [[ -n "$path" ]]; then
          echo "    - Opening: $path"
          # 'open' コマンドでアプリケーションを起動
          if ! open "$path"; then
            echo "⚠️  Warning: Failed to open application: $path" >&2
          fi
          path="" # 次のアイテムのためにリセット
        fi
      fi
    done < "$yaml_file"
    echo "✅ All specified applications have been launched."
  else
    echo "🚫 Operation cancelled."
  fi
}


# --- メイン処理 ---
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
  open) # ✨ openコマンドを追加
    open_login_items
    ;;
  *)
    echo "Usage: $0 {export|check|import|open}"
    echo
    echo "  export: Save current login items to ${yaml_file}"
    echo "  check:  Compare current login items with the YAML file"
    echo "  import: Sync login items based on the YAML file (deletes all current items first)"
    echo "  open:   Open all applications listed in the YAML file" # ✨ 説明を追加
    exit 1
    ;;
esac
