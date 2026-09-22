#!/bin/bash
PROJECT_NAME=$(basename "$PWD")
say -v Victoria "${PROJECT_NAME}、完了"
osascript -e "display notification with title \"${PROJECT_NAME}: 完了\" sound name \"Submarine\""
