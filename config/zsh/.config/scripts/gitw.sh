#!/bin/bash

# 使用法を表示する関数
usage() {
    echo "使用法: $0 <リンク元ディレクトリ(A)> <リンク先ディレクトリ(B)>"
    exit 1
}

if [ "$#" -ne 2 ]; then
    usage
fi

SRC_DIR="$1"
DEST_DIR="$2"

if [ ! -d "$SRC_DIR" ]; then
    echo "エラー: リンク元ディレクトリ '$SRC_DIR' が存在しません。"
    exit 1
fi
if [ ! -d "$DEST_DIR" ]; then
    echo "エラー: リンク先ディレクトリ '$DEST_DIR' が存在しません。"
    exit 1
fi

ABS_SRC_DIR=$(cd "$SRC_DIR" && pwd)
ABS_DEST_DIR=$(cd "$DEST_DIR" && pwd)

echo "--- 処理を開始します ---"
echo "元(Repo):  $ABS_SRC_DIR"
echo "先(Home):  $ABS_DEST_DIR"
echo "------------------------"

for file_path in "$ABS_SRC_DIR"/{.,}*; do
    filename=$(basename "$file_path")

    # 1. 基本的な除外
    if [ "$filename" == "." ] || [ "$filename" == ".." ] || [ "$filename" == "*" ]; then
        continue
    fi

    # 2. .git の除外 (第一の防壁)
    # ここで除外すれば、下の処理には一切進みません
    if [ "$filename" == ".git" ]; then
        # echo "無視: .git ディレクトリは触りません"
        continue
    fi

    target_path="$ABS_DEST_DIR/$filename"

    # --- 判定ロジック ---

    # ケースA: リンク先に「実体」があり、かつ「リンクではない」場合
    if [ -e "$target_path" ] && [ ! -L "$target_path" ]; then
        
        # ★★★ 【重要】 二重の安全装置 (Fail-Safe) ★★★
        # 万が一、上の除外漏れがあってもここで .git の破壊を阻止します
        if [ "$filename" == ".git" ]; then
            echo "警告: .git の上書きを検知しました。処理を緊急スキップします。"
            continue
        fi
        
        echo "衝突検知: $filename は先に実体があります。"
        
        # Aにある元ファイルを削除
        rm -rf "$file_path"
        
        # Bの実体をAに移動
        mv "$target_path" "$file_path"
        echo "  -> 先(B)の実体を 元(A)へ移動(上書き)しました。"
        
        # リンクを作成
        ln -s "$file_path" "$target_path"
        echo "  -> リンク作成: $filename -> $target_path"

    # ケースB: リンク先にすでに「シンボリックリンク」が存在する場合
    elif [ -L "$target_path" ]; then
        echo "スキップ (リンク済): $filename"

    # ケースC: リンク先に何もない場合
    else
        ln -s "$file_path" "$target_path"
        echo "リンク作成: $filename -> $target_path"
    fi
done

# --- 先(B)にだけ存在するファイルを元(A)に回収してリンクを張る ---
echo "------------------------"
echo "--- 先(B)にのみ存在するファイルを確認中 ---"

for target_path in "$ABS_DEST_DIR"/{.,}*; do
    filename=$(basename "$target_path")

    # 1. 基本的な除外
    if [ "$filename" == "." ] || [ "$filename" == ".." ] || [ "$filename" == "*" ]; then
        continue
    fi

    # 2. .git の除外
    if [ "$filename" == ".git" ]; then
        continue
    fi

    file_path="$ABS_SRC_DIR/$filename"

    # 元(A)に存在しない & 先(B)に実体がある(リンクではない)場合のみ処理
    if [ ! -e "$file_path" ] && [ ! -L "$file_path" ] && [ -e "$target_path" ] && [ ! -L "$target_path" ]; then
        echo "回収検知: $filename は先(B)にのみ存在します。"

        # Bの実体をAに移動
        mv "$target_path" "$file_path"
        echo "  -> 先(B)の実体を 元(A)へ移動しました。"

        # リンクを作成
        ln -s "$file_path" "$target_path"
        echo "  -> リンク作成: $filename -> $target_path"
    fi
done

echo "--- 完了しました ---"
