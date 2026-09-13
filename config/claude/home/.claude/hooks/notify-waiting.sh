#!/bin/bash
PROJECT_NAME=$(basename "$PWD")
say -v Victoria "${PROJECT_NAME}、入力待ち"
osascript -e "display notification with title \"${PROJECT_NAME}: 入力待ち\" sound name \"Submarine\""
