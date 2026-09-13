#!/bin/bash

# ツールのパス (同じディレクトリ内のgetmicvolを使用)
TOOL_PATH="$(dirname "$0")/../tools/getmicvol"

# ツールを実行
# マイク未使用時は "OFF"、使用中は "0"〜"100" が返ってくる
OUTPUT=$($TOOL_PATH)

# マイクが使われていない場合
if [ "$OUTPUT" == "OFF" ]; then
  # ここで「使われていない時の表示」を設定
  # 例: マイクアイコンに斜線など (SF Symbolsが使える場合は :mic.slash: など)
  echo ":mic.slash:" 
  echo "---"
  echo "Status: Microphone is inactive"
  exit 0
fi

# 以下はマイク使用中のレベル表示ロジック（数値として扱う）
VOL=$OUTPUT

if [ "$VOL" -eq 0 ]; then
  BAR="--"
  COLOR="gray"
elif [ "$VOL" -lt 20 ]; then
  BAR="▂ "
  COLOR="white"
elif [ "$VOL" -lt 40 ]; then
  BAR="▂▃"
  COLOR="white"
elif [ "$VOL" -lt 60 ]; then
  BAR="▂▃▅"
  COLOR="white"
elif [ "$VOL" -lt 80 ]; then
  BAR="▂▃▅▆"
  COLOR="orange"
else
  BAR="▂▃▅▆▇"
  COLOR="red"
fi

echo "$BAR | color=$COLOR size=12 font=Menlo"
echo "---"
echo "Status: Live Input"
echo "Level: $VOL%"
