#!/usr/bin/env bash
set -euo pipefail

# plist.yamlのパス
yaml_file="$(dirname "$0")/../data/plist.yaml"

# yqコマンドが必要
yq_bin="yq"
if ! command -v "$yq_bin" &>/dev/null; then
  echo "Error: yq is required but not installed." >&2
  exit 1
fi

# 各applicationごとに処理
yq eval '.applications[]' "$yaml_file" | \
while read -r line; do
  # name, path, domainを抽出
  if [[ $line == "- name:"* ]]; then
    name=$(echo "$line" | awk '{print $3}')
    read -r path_line
    path=$(echo "$path_line" | awk '{print $2}')
    read -r domain_line
    domain=$(echo "$domain_line" | awk '{print $2}')
    # defaults readでエクスポート
    mkdir -p "$path"
    defaults read "$domain" > "$path/$domain.txt"
    echo "Exported $domain to $path/$domain.txt"
  fi
done 
