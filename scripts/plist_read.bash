#!/usr/bin/env bash
set -euo pipefail

# --- 事前チェックと初期設定 ---
script_dir=$(cd "$(dirname "$0")" && pwd)
yaml_file="${script_dir}/../data/plist.yaml"
base_out_dir=$(cd "${script_dir}/../config" && pwd)

# 依存関係 yq のチェック
yq_bin="yq"
if ! command -v "$yq_bin" &>/dev/null; then
  echo "❌ エラー: yq が必要ですが、インストールされていません。" >&2
  exit 1
fi

# サブコマンドのチェック
if [ $# -eq 0 ]; then
  echo "❌ エラー: サブコマンド (export, import, check のいずれか) を指定してください。" >&2
  echo "使用法: $0 {export|import|check}" >&2
  exit 1
fi
subcommand="$1"

# --- ヘルパー関数 ---

# plistファイルのMD5ハッシュ値を計算する
# @param $1: plistファイルのパス
get_plist_hash() {
  local file_path="$1"
  if [ -f "$file_path" ]; then
    # plutil でXML形式に変換し、ハッシュ値の一貫性を保つ
    (plutil -convert xml1 -o - "$file_path" 2>/dev/null | md5) || true
  else
    echo ""
  fi
}

# --- サブコマンドの実装 ---

# macOSのdefaultsから設定を .plist と .txt ファイルにエクスポートする
do_export() {
  echo "🚀 設定のエクスポート処理を開始します..."
  echo "📂 出力先のベースディレクトリ: $base_out_dir"
  echo ""

  successful_domains=()

  while IFS=$'\t' read -r name path domain; do
    echo "--- 処理開始: $name ($domain) ---"

    local final_out_dir="${base_out_dir}/${path}"
    mkdir -p "$final_out_dir"

    # --- 修正箇所 ---
    local plist_out_file
    if [[ "${domain}" == *.plist ]]; then
      plist_out_file="${domain}"
    else
      plist_out_file="${domain}.plist"
    fi
    # --- 修正ここまで ---
    
    local final_plist_path="${final_out_dir}/${plist_out_file}"
    local txt_out_file="${domain}.txt"
    local final_txt_path="${final_out_dir}/${txt_out_file}"

    # 表示用に相対パスを生成
    local display_plist_path
    display_plist_path=$(echo "$final_plist_path" | sed -e "s,^${base_out_dir}/,," -e "s,//,/,g")
    local display_txt_path
    display_txt_path=$(echo "$final_txt_path" | sed -e "s,^${base_out_dir}/,," -e "s,//,/,g")

    # [1/2] .plist をエクスポート
    echo "  [1/2] .plist をエクスポート中..."
    echo "        -> ${display_plist_path}"
    
    local temp_plist_path
    temp_plist_path=$(mktemp 2>/dev/null || mktemp -t 'plist-temp')
    if [ -z "$temp_plist_path" ]; then
        echo "        ⚠️  一時ファイルの作成に失敗しました。" >&2
        continue
    fi

    if defaults export "$domain" "$temp_plist_path"; then
      local hash_before
      hash_before=$(get_plist_hash "$final_plist_path")
      local hash_after
      hash_after=$(get_plist_hash "$temp_plist_path")

      if [ "$hash_before" != "$hash_after" ]; then
        echo "        ✅ 更新を検知しました。ファイルを保存します。"
        mv "$temp_plist_path" "$final_plist_path"
        successful_domains+=("$domain")
      else
        echo "        ℹ️  内容は変更ありませんでした。"
        rm "$temp_plist_path"
      fi
    else
      echo "        ⚠️  .plistのエクスポートに失敗しました (ドメインが存在しない可能性があります)。" >&2
      rm "$temp_plist_path"
    fi

    # [2/2] .txt を生成
    echo "  [2/2] .txt を生成中..."
    echo "        -> ${display_txt_path}"
    if defaults read "$domain" > "$final_txt_path"; then
      echo "        ✅ .txt ファイルの生成が完了しました。"
    else
      echo "        ⚠️  .txt ファイルの生成に失敗しました。" >&2
      rm -f "$final_txt_path"
    fi

    echo "--- 処理完了: $name ---"
    echo ""
  done < <(yq eval '.applications[] | [.name, .path, .domain] | @tsv' "$yaml_file")

  echo "🎉 すべてのエクスポート処理が完了しました。"
  echo ""

  if [ ${#successful_domains[@]} -gt 0 ]; then
    echo "---"
    echo "✅ 正常に更新されたplistドメイン一覧:"
    for d in "${successful_domains[@]}"; do
      echo "  - $d"
    done
  else
    echo "---"
    echo "ℹ️ 今回の実行で内容が更新されたplistはありませんでした。"
  fi
}

# .plist ファイルからmacOSのdefaultsに設定をインポートする
do_import() {
  echo "🚀 設定のインポート処理を開始します..."
  echo "📂 設定ファイルのベースディレクトリ: $base_out_dir"
  echo ""

  while IFS=$'\t' read -r name path domain; do
    echo "--- 処理開始: $name ($domain) ---"

    local final_out_dir="${base_out_dir}/${path}"

    # --- 修正箇所 ---
    local plist_in_file
    if [[ "${domain}" == *.plist ]]; then
      plist_in_file="${domain}"
    else
      plist_in_file="${domain}.plist"
    fi
    # --- 修正ここまで ---

    local final_plist_path="${final_out_dir}/${plist_in_file}"

    if [ -f "$final_plist_path" ]; then
      echo "      ${final_plist_path} からインポートしています..."
      if defaults import "$domain" "$final_plist_path"; then
        echo "      ✅ '$domain' の設定を正常にインポートしました。"
      else
        echo "      ❌ '$domain' の設定のインポート中にエラーが発生しました。" >&2
      fi
    else
      echo "      ⚠️  スキップ: plistファイルが見つかりません (${final_plist_path})。"
    fi
    echo ""
  done < <(yq eval '.applications[] | [.name, .path, .domain] | @tsv' "$yaml_file")

  echo "🎉 すべてのインポート処理が完了しました。"
  echo "ℹ️  注意: 設定を反映させるには、一部のアプリケーションの再起動が必要な場合があります。"
}

# 現在の設定と保存された .txt ファイルとの差分をチェックする
do_check() {
  echo "🚀 設定の差分チェック処理を開始します..."
  echo "📂 設定ファイルのベースディレクトリ: $base_out_dir"
  echo ""

  local changes_found=0

  while IFS=$'\t' read -r name path domain; do
    echo "--- チェック中: $name ($domain) ---"

    local txt_file_path="${base_out_dir}/${path}/${domain}.txt"
    
    if [ ! -f "$txt_file_path" ]; then
      echo "  ⚠️  スキップ: 保存された.txtファイルが見つかりません (${txt_file_path})。"
      echo ""
      continue
    fi

    # 比較のために現在設定を一時ファイルに書き出す
    local temp_txt_path
    temp_txt_path=$(mktemp 2>/dev/null || mktemp -t 'check-temp')
    if ! defaults read "$domain" > "$temp_txt_path" 2>/dev/null; then
      echo "  ℹ️  このドメインの現在の設定が見つかりません。差分なしと見なします。"
      rm "$temp_txt_path"
      echo ""
      continue
    fi

    echo "  現在の設定と ${txt_file_path} を比較しています..."
    # diffコマンドで差分を表示 (-u: unified形式, --color=always: 常に色付け)
    if diff --color=always -u "$txt_file_path" "$temp_txt_path"; then
      echo "  ✅ 差分はありませんでした。"
    else
      changes_found=1
      # diffが差分を出力するので、ここでは補足メッセージのみ表示
      echo "  👆 '$domain' に差分が見つかりました。"
    fi
    
    rm "$temp_txt_path"
    echo ""
  done < <(yq eval '.applications[] | [.name, .path, .domain] | @tsv' "$yaml_file")
  
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
    do_import
    ;;
  check)
    do_check
    ;;
  *)
    echo "❌ エラー: 無効なサブコマンド '$subcommand' です。" >&2
    echo "使用法: $0 {export|import|check}" >&2
    exit 1
    ;;
esac
