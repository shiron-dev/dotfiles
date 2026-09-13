#!/bin/bash

usage() {
    echo "Usage: $0 {clone|pull} [LOCAL_PATH] [REMOTE_PATH] [-f]"
    echo "  clone: 手元(LOCAL) --> コピー先(REMOTE) ※マーカーは先に作成"
    echo "  pull : コピー先(REMOTE) --> 手元(LOCAL)"
    exit 1
}

# --- 引数解析 ---
FORCE=false
PARAMS=()
for arg in "$@"; do
    if [[ "$arg" == "-f" ]]; then FORCE=true; else PARAMS+=("$arg"); fi
done

MODE=${PARAMS[0]}
LOCAL_DIR=${PARAMS[1]}   # 手元（保護対象）
REMOTE_DIR=${PARAMS[2]}  # コピー先（マーカー設置場所）
MARKER=".last_sync"

if [[ -z "$MODE" || -z "$LOCAL_DIR" || -z "$REMOTE_DIR" ]]; then
    usage
fi

LOCAL_PATH="${LOCAL_DIR%/}/"
REMOTE_PATH="${REMOTE_DIR%/}/"
# マーカーのパスをコピー先に設定
MARKER_PATH="$REMOTE_DIR/$MARKER"

# --- 確認プロンプト ---
confirm_proceed() {
    if [ "$FORCE" = true ]; then
        echo -n "上記の競合がありますが、コピーを強行しますか？ (y/N): "
        read -r CONFIRM < /dev/tty
        if [[ "$CONFIRM" != "y" && "$CONFIRM" != "Y" ]]; then
            echo "ユーザーにより中断されました。"
            exit 1
        fi
    else
        echo "エラー: 手元に更新があるため中断しました。解決するか -f を使用してください。"
        exit 1
    fi
}

# --- 安全チェック ---
check_local_safety() {
    # コピー先にあるマーカーを確認
    if [ ! -f "$MARKER_PATH" ]; then
        echo "【警告】コピー先($REMOTE_DIR)に同期マーカーが見つかりません。"
        confirm_proceed
        return
    fi

    echo "--- 手元の安全チェック開始 (基準: コピー先のマーカー) ---"

    # 1. コピー先のマーカー時刻以降に「手元」で編集されたファイルがあるか
    LOCAL_MODIFIED=$(find "$LOCAL_PATH" -type f -newer "$MARKER_PATH" \
        -not -path "*/.git/*")

    if [ -n "$LOCAL_MODIFIED" ]; then
        echo "【競合】前回の同期以降に「手元」で以下のファイルが更新されています:"
        echo "------------------------------------------------"
        echo "$LOCAL_MODIFIED"
        echo "------------------------------------------------"
        confirm_proceed
    fi

    # 2. 「手元」の方が「コピー先」よりタイムスタンプが新しいファイルがあるか
    REVERSE_NEWER=$(rsync -n -avu --exclude='.git' --exclude="$MARKER" --out-format="%n" "$LOCAL_PATH" "$REMOTE_PATH" \
        | grep -v 'incremental file list' | grep -v 'sending list' | grep -v 'total size' | grep -v '^$')

    if [ -n "$REVERSE_NEWER" ]; then
        echo "【競合】コピー先より新しいタイムスタンプのファイルが「手元」にあります:"
        echo "------------------------------------------------"
        echo "$REVERSE_NEWER"
        echo "------------------------------------------------"
        confirm_proceed
    fi
    echo "チェック完了: 手元に競合はありません。"
}

# --- メイン処理 ---
case "$MODE" in
    "clone")
        if [ -d "$REMOTE_DIR" ] && [ "$(ls -A "$REMOTE_DIR" 2>/dev/null)" ]; then
            echo "【エラー】コピー先ディレクトリが空ではありません。"
            confirm_proceed
        fi
        
        echo "Clone: 手元($LOCAL_PATH) --> コピー先($REMOTE_PATH)"
        mkdir -p "$REMOTE_DIR"
        rsync -av --exclude='.git' "$LOCAL_PATH" "$REMOTE_PATH"
        
        # 同期完了マーカーを【コピー先】に作成
        touch "$MARKER_PATH"
        echo "マーカーをコピー先に作成しました: $MARKER_PATH"
        ;;

    "pull")
        if [ ! -d "$REMOTE_DIR" ]; then
            echo "【エラー】コピー先($REMOTE_DIR)が存在しません。"
            exit 1
        fi
        
        # 手元の状態を、コピー先のマーカーと比較してチェック
        check_local_safety
        
        echo "Pull: コピー先($REMOTE_PATH) --> 手元($LOCAL_PATH)"
        # マーカー自体は手元にコピーしないよう exclude する
        rsync -av --exclude='.git' --exclude="$MARKER" "$REMOTE_PATH" "$LOCAL_PATH"
        
        # Pull成功後、コピー先のマーカーを更新
        touch "$MARKER_PATH"
        echo "コピー先のマーカーを更新しました。"
        ;;

    *)
        usage
        ;;
esac

echo "完了しました。"
