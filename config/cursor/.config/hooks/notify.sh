#!/bin/bash

PROJECT_NAME=$(basename "$PWD")
STATUS=$1
TASK_DESC=$2

case $STATUS in
  "success")
    osascript -e "display notification \"✅ ${TASK_DESC}が完了しました\" with title \"${PROJECT_NAME}\" sound name \"Submarine\""
    ;;
  "error")
    osascript -e "display notification \"❌ ${TASK_DESC}でエラーが発生しました\" with title \"${PROJECT_NAME}\" sound name \"Basso\""
    ;;
esac
