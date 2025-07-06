#!/usr/bin/env bash

set -e

REPO_URL="https://github.com/shiron-dev/dotfiles.git"
REPO_PATH="$HOME/projects/github.com/shiron-dev/dotfiles"

cat <<EOM

# shiron-dev dotfiles

Start setup.
For more information check the following link.

https://github.com/shiron-dev/dotfiles

EOM

# Check Homebrew
if ! command -v brew >/dev/null 2>&1; then
  echo "[INFO] Homebrew not found. Installing..."
  /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
else
  echo "[INFO] Homebrew found."
fi

# Check git
if ! command -v git >/dev/null 2>&1; then
  echo "[INFO] git not found. Installing via Homebrew..."
  brew install git
else
  echo "[INFO] git found."
fi

# Clone dotfiles if not exists
if [ ! -d "$REPO_PATH" ]; then
  echo "[INFO] Cloning dotfiles repository to $REPO_PATH ..."
  git clone "$REPO_URL" "$REPO_PATH"
else
  echo "[INFO] dotfiles repository already exists at $REPO_PATH."
fi

# Check Go
if ! command -v go >/dev/null 2>&1; then
  echo "[INFO] Go not found. Installing via Homebrew..."
  brew install go
else
  echo "[INFO] Go found."
fi

# Install brew-management
echo "[INFO] Installing brew-management..."
cd "$REPO_PATH/scripts/brew-management" && go install

cat <<EOM

✅ 初期セットアップが完了しました。

次のコマンドを順に実行してください：

brew-management install
cd $REPO_PATH/scripts/ansible
ansible-playbook -i hosts.yml site.yml
$REPO_PATH/scripts/login_manager.bash check
$REPO_PATH/scripts/login_manager.bash import

EOM
