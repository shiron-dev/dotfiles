#!/usr/bin/env bash
set -euo pipefail

# --- 事前チェックと初期設定 ---
script_dir=$(cd "$(dirname "$0")" && pwd)
base_yaml_file="${script_dir}/../data/plist.yaml"
config_dir="${script_dir}/../config"
key_config_dir="${script_dir}/../data/plist"

# 依存関係 yq のチェック
yq_bin="yq"
if ! command -v "$yq_bin" &>/dev/null; then
  echo "❌ エラー: yq が必要ですが、インストールされていません。" >&2
  exit 1
fi

# サブコマンドのチェック
if [ $# -eq 0 ]; then
  echo "❌ エラー: サブコマンド (export, import, check のいずれか) を指定してください。" >&2
  echo "使用法: $0 {export|import [-y]|check}" >&2
  exit 1
fi
subcommand="$1"
shift # サブコマンドを引数リストから削除

# --- ヘルパー関数 ---

# 指定されたドメインの管理対象となるキーのリストを取得する
# @param $1: ドメイン名 (例: com.apple.dock)
# @param $2: サニタイズされたドメイン名 (例: com.apple.dock)
# @return: 管理対象となるキーのリスト(スペース区切り)
get_filtered_keys() {
  local domain="$1"
  local sanitized_domain="$2"
  local key_yaml_file="${key_config_dir}/${sanitized_domain}.yaml"

  # 現在のシステムに設定されている全キーリストを取得 (堅牢な方法)
  local all_keys
  all_keys=$(defaults export "$domain" - | plutil -p - | grep '=>' | sed -E 's/^[[:space:]]*"([^"]+)"[[:space:]]*=>.+$/\1/' || echo "")

  if [ ! -f "$key_yaml_file" ]; then
    # キー設定ファイルがなければ、すべてのキーを対象とする
    echo "$all_keys"
    return
  fi

  local include_all
  include_all=$($yq_bin eval '.include_all // false' "$key_yaml_file")
  
  if [ "$include_all" = "true" ]; then
    local excluded_keys
    excluded_keys=$($yq_bin eval '.exclude[]?' "$key_yaml_file")
    
    # all_keysからexcluded_keysを除外する
    echo "$all_keys" | grep -vFf <(echo "$excluded_keys" | tr ' ' '\n') | tr '\n' ' '
  else
    # includeが指定されている場合 (今回は include_all: true のユースケースを優先)
    # local included_keys=$($yq_bin eval '.include[]?' "$key_yaml_file")
    # echo "$included_keys"
    echo "" # include_all: false の場合はキーを返さない
  fi
}


# --- サブコマンドの実装 ---

# 設定を .txt ファイルにエクスポートする
do_export() {
  echo "🚀 設定のエクスポート処理を開始します..."
  echo "📂 出力先のベースディレクトリ: $config_dir"
  echo ""

  while IFS=$'\t' read -r name path domain; do
    local sanitized_domain=${domain//\//-}
    echo "--- 処理開始: $name ($domain) ---"

    local final_out_dir="${config_dir}/${path}"
    mkdir -p "$final_out_dir"
    
    local final_txt_path="${final_out_dir}/${sanitized_domain}.txt"
    local display_txt_path
    display_txt_path=$(echo "$final_txt_path" | sed -e "s,^${config_dir}/,," -e "s,//,/,g")

    echo "  [1/1] .txt を生成中 (キーフィルタ適用)..."
    echo "        -> ${display_txt_path}"

    # 管理対象のキーリストを取得
    local filtered_keys
    filtered_keys=$(get_filtered_keys "$domain" "$sanitized_domain")

    if [ -z "$filtered_keys" ]; then
      echo "        ℹ️  管理対象のキーがないため、ファイルを空にします。"
      # ファイルを空にする
      >"$final_txt_path"
      echo "--- 処理完了: $name ---"
      echo ""
      continue
    fi
    
    # 一時ファイルにエクスポート
    local temp_txt_path
    temp_txt_path=$(mktemp 2>/dev/null || mktemp -t 'export-temp')
    
    for key in $filtered_keys; do
      # readの出力は不安定なことがあるため、plist形式でexportしたものを変換して値を取得
      local value
      value=$(defaults export "$domain" - | plutil -extract "$key" xml1 - -o - | sed -e '1d;$d' -e 's/^[[:space:]]*//' -e 's/[[:space:]]*$//')
      echo "\"$key\" = $value;" >> "$temp_txt_path"
    done
    
    # 既存ファイルと比較して更新があれば置き換え
    if ! diff -q "$final_txt_path" "$temp_txt_path" >/dev/null 2>&1; then
      echo "        ✅ 更新を検知しました。ファイルを保存します。"
      mv "$temp_txt_path" "$final_txt_path"
    else
      echo "        ℹ️  内容は変更ありませんでした。"
      rm "$temp_txt_path"
    fi
    
    echo "--- 処理完了: $name ---"
    echo ""
  done < <(yq eval '.applications[] | [.name, .path, .domain] | @tsv' "$base_yaml_file")

  echo "🎉 すべてのエクスポート処理が完了しました。"
}


# .txt ファイルから設定をインポートする
do_import() {
  local force_import=false
  if [[ "${1:-}" == "-y" ]]; then
    force_import=true
  fi
  
  if ! $force_import; then
    echo "🔄 インポート前に現在の設定との差分を確認します..."
    echo ""
    local changes_found
    changes_found=$(do_check --quiet)
    
    if [ "$changes_found" -eq 1 ]; then
        echo "---"
        do_check
        echo "---"
        
        read -p "☝️ 設定に差分が見つかりました。インポートを実行しますか？ (y/N): " -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            echo "🚫 インポートをキャンセルしました。"
            exit 0
        fi
    else
        echo "✅ 差分はありません。インポート処理は不要です。"
        exit 0
    fi
  fi

  echo "🚀 設定のインポート処理を開始します..."
  if $force_import; then
      echo "ℹ️  -y オプションが指定されたため、確認をスキップして実行します。"
  fi
  echo "📂 設定ファイルのベースディレクトリ: $config_dir"
  echo ""

  while IFS=$'\t' read -r name path domain; do
    local sanitized_domain=${domain//\//-}
    echo "--- 処理開始: $name ($domain) ---"

    local txt_file_path="${config_dir}/${path}/${sanitized_domain}.txt"
    if [ ! -f "$txt_file_path" ]; then
      echo "      ⚠️  スキップ: .txtファイルが見つかりません (${txt_file_path})。"
      continue
    fi

    # 管理対象のキーリストを取得
    mapfile -t manageable_keys < <(get_filtered_keys "$domain" "$sanitized_domain" | tr ' ' '\n')

    # .txtファイルを一行ずつ読み込む
    while IFS= read -r line || [[ -n "$line" ]]; do
      # 空行はスキップ
      if [[ -z "$line" ]]; then continue; fi

      # "key" = value; の形式からキーと値をパース
      local key
      key=$(echo "$line" | sed -n 's/^[[:space:]]*"\([^"]*\)".*/\1/p')
      local value
      value=$(echo "$line" | sed -n 's/.*= \(.*\);/\1/p')

      if [ -z "$key" ] || [ -z "$value" ]; then
        echo "      ⚠️  不正な行をスキップ: $line"
        continue
      fi

      # 管理対象のキーかチェック
      local should_write=false
      for mkey in "${manageable_keys[@]}"; do
        if [[ "$mkey" == "$key" ]]; then
          should_write=true
          break
        fi
      done

      if $should_write; then
        echo "      書き込み中: $key"
        # ⚠️ 型情報が失われる可能性があります
        defaults write "$domain" "$key" "$value"
      fi

    done < "$txt_file_path"

    echo "--- 処理完了: $name ---"
    echo ""
  done < <(yq eval '.applications[] | [.name, .path, .domain] | @tsv' "$base_yaml_file")

  echo "🎉 すべてのインポート処理が完了しました。"
  echo "ℹ️  注意: 設定を反映させるには、一部のアプリケーションの再起動が必要な場合があります。"
}

# 現在の設定と保存されたファイルとの差分をチェックする
do_check() {
  local quiet_mode=false
  if [[ "${1:-}" == "--quiet" ]]; then
    quiet_mode=true
  fi

  if ! $quiet_mode; then
    echo "🚀 設定の差分チェック処理を開始します..."
    echo "📂 設定ファイルのベースディレクトリ: $config_dir"
    echo ""
  fi

  local changes_found=0

  while IFS=$'\t' read -r name path domain; do
    local sanitized_domain=${domain//\//-}

    if ! $quiet_mode; then
      echo "--- チェック中: $name ($domain) ---"
    fi

    local txt_file_path="${config_dir}/${path}/${sanitized_domain}.txt"
    if [ ! -f "$txt_file_path" ]; then
      if ! $quiet_mode; then
        echo "  ⚠️  スキップ: 保存された.txtファイルが見つかりません (${txt_file_path})。"
        echo ""
      fi
      continue
    fi

    # 現在の設定から一時ファイルを作成
    local temp_txt_path
    temp_txt_path=$(mktemp 2>/dev/null || mktemp -t 'check-temp')
    local filtered_keys
    filtered_keys=$(get_filtered_keys "$domain" "$sanitized_domain")

    if [ -n "$filtered_keys" ]; then
       for key in $filtered_keys; do
          local value
          value=$(defaults export "$domain" - | plutil -extract "$key" xml1 - -o - | sed -e '1d;$d' -e 's/^[[:space:]]*//' -e 's/[[:space:]]*$//')
          echo "\"$key\" = $value;" >> "$temp_txt_path"
       done
    fi

    if ! diff -q "$txt_file_path" "$temp_txt_path" >/dev/null; then
        changes_found=1
        if ! $quiet_mode; then
            echo "  現在の設定と ${txt_file_path} を比較しています..."
            diff --color=always -u "$txt_file_path" "$temp_txt_path" || true
            echo "  👆 '$domain' に差分が見つかりました。"
        fi
    else
        if ! $quiet_mode; then
            echo "  ✅ 差分はありませんでした。"
        fi
    fi
    
    rm "$temp_txt_path"
    if ! $quiet_mode; then
        echo ""
    fi
  done < <(yq eval '.applications[] | [.name, .path, .domain] | @tsv' "$base_yaml_file")
  
  if $quiet_mode; then
    echo $changes_found
    return
  fi

  echo "🎉 すべてのチェック処理が完了しました。"
  if [ $changes_found -eq 0 ]; then
    echo "✅ 現在の設定と保存されている設定の間に差分はありませんでした。"
  else
    echo "⚠️  差分が見つかりました。'./script.sh export' を実行して、保存されている設定を更新してください。"
  fi
}


# --- メインロジック ---

case "$subcommand" in
  export)
    do_export
    ;;
  import)
    do_import "$@"
    ;;
  check)
    do_check
    ;;
  *)
    echo "❌ エラー: 無効なサブコマンド '$subcommand' です。" >&2
    echo "使用法: $0 {export|import [-y]|check}" >&2
    exit 1
    ;;
esac
