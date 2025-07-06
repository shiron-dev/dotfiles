#!/usr/bin/env bash
set -euo pipefail

script_dir=$(cd "$(dirname "$0")" && pwd)

yaml_file="${script_dir}/../data/plist.yaml"
base_out_dir=$(cd "${script_dir}/../config" && pwd)

yq_bin="yq"
if ! command -v "$yq_bin" &>/dev/null;
then
  echo "❌ Error: yq is required but not installed." >&2
  exit 1
fi

get_plist_hash() {
  local file_path="$1"
  if [ -f "$file_path" ]; then
    (plutil -convert xml1 -o - "$file_path" 2>/dev/null | md5) || true
  else
    echo ""
  fi
}

echo "🚀 設定ファイルのエクスポートおよび読み込み処理を開始します"
echo "📂 出力先のベースディレクトリ: $base_out_dir"
echo ""

successful_domains=()

while IFS=$'\t' read -r name path domain; do
  echo "--- 処理開始: $name ($domain) ---"

  final_out_dir="${base_out_dir}/${path}"
  mkdir -p "$final_out_dir"

  if [[ "${domain}" == *.plist ]]; then
    plist_out_file="${domain}"
  else
    plist_out_file="${domain}.plist"
  fi
  final_plist_path="${final_out_dir}/${plist_out_file}"

  txt_out_file="${domain}.txt"
  final_txt_path="${final_out_dir}/${txt_out_file}"

  display_plist_path=$(echo "$final_plist_path" | sed -e "s,^${base_out_dir}/,," -e "s,//,/,g")
  display_txt_path=$(echo "$final_txt_path" | sed -e "s,^${base_out_dir}/,," -e "s,//,/,g")

  echo "  [1/2] .plist をエクスポート中..."
  echo "        -> ${display_plist_path}"
  
  temp_plist_path=""
  temp_plist_path=$(mktemp 2>/dev/null || mktemp -t 'plist-temp')
  if [ -z "$temp_plist_path" ]; then
      echo "        ⚠️  一時ファイルの作成に失敗しました" >&2
      continue
  fi

  if defaults export "$domain" "$temp_plist_path"; then
    hash_before=$(get_plist_hash "$final_plist_path")
    hash_after=$(get_plist_hash "$temp_plist_path")

    if [ "$hash_before" != "$hash_after" ]; then
      echo "        ✅ 更新を検知しました"
      mv "$temp_plist_path" "$final_plist_path"
      successful_domains+=("$domain")
    else
      echo "        ℹ️  内容は変更ありませんでした"
      rm "$temp_plist_path"
    fi
  else
    echo "        ⚠️  設定が見つからないか、.plistのエクスポートに失敗しました" >&2
    rm "$temp_plist_path"
  fi

  echo "  [2/2] .txt を生成中..."
  echo "        -> ${display_txt_path}"
  if defaults read "$domain" > "$final_txt_path"; then
    echo "        ✅ .txt 生成完了"
  else
    echo "        ⚠️  設定が見つからないか、.txtの生成に失敗しました" >&2
    rm -f "$final_txt_path"
  fi

  echo "--- 処理完了: $name ---"
  echo ""
done < <(yq eval '.applications[] | [.name, .path, .domain] | @tsv' "$yaml_file")

echo "🎉 すべての処理が完了しました"
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
