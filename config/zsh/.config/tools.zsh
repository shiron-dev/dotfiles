# Homebrew
eval "$(/opt/homebrew/bin/brew shellenv)"
export DYLD_LIBRARY_PATH="$DYLD_LIBRARY_PATH:/opt/homebrew/lib/"

# sheldon
eval "$(sheldon source)"

# mise
eval "$(/opt/homebrew/bin/mise activate zsh)"

# iTerm2
bindkey "^[[H" beginning-of-line
bindkey "^[[F" end-of-line

# pnpm
export PNPM_HOME="/Users/shiron/Library/pnpm"
case ":$PATH:" in
  *":$PNPM_HOME:"*) ;;
  *) export PATH="$PNPM_HOME:$PATH" ;;
esac

# volta
export VOLTA_HOME="$HOME/.volta"
export PATH="$VOLTA_HOME/bin:$PATH"

# gnu
export PATH="$PATH:/opt/homebrew/opt/gawk/libexec/gnubin"
alias gcc="gcc-14"
alias g++="g++-14"
alias sed='gsed'

# python
alias python="python3"
alias pip="pip3"

# adb
export PATH="$PATH:/Users/shiron/Library/Android/sdk/platform-tools"

# Jetbrains Toolbox
export PATH="$PATH:/Users/shiron/Library/Application Support/JetBrains/Toolbox/scripts"

# Flutter
export PATH="$PATH":"$HOME/.pub-cache/bin"

# MySQL
export PATH="$PATH:/opt/homebrew/opt/mysql-client@8.0/bin"

# Tailscale
alias tailscale='/Applications/Tailscale.app/Contents/MacOS/Tailscale'

# OrbStack
source ~/.orbstack/shell/init.zsh 2>/dev/null || :

# Golang
export GOPATH=$HOME/go
export GOBIN=$GOPATH/bin
export PATH=$PATH:$GOBIN

# lazy
alias lg='lazygit'
alias ld='lazydocker'

# fzf
source <(fzf --zsh)

# zoxide
eval "$(zoxide init zsh)"

# Android

export ANDROID_HOME="/Users/$USER/Library/Android/sdk"
export PATH="$PATH":"$ANDROID_HOME/tools":"$ANDROID_HOME/build-tools/35.0.0"

# My tools
export PATH="$PATH":"/Users/shiron/projects/tools/bin"
export PATH="$PATH":"/Users/shiron/projects/github.com/shiron-dev/arcanum-hue/bin"

if [ ! -f "/Users/shiron/projects/github.com/shiron-dev/dotfiles/scripts/dofy/dofy" ]; then
  echo "Building dofy..."
  (cd ~/projects/github.com/shiron-dev/dotfiles/scripts/dofy && go build -o dofy cmd/main.go)
fi
alias dofy="/Users/shiron/projects/github.com/shiron-dev/dotfiles/scripts/dofy/dofy"

# My Aliases
alias grep="ggrep"

alias docker-compose-rm="docker compose down --rmi all --volumes --remove-orphans"
alias lsusb="system_profiler SPUSBDataType"
alias gic="git clean -Xdf"

# My functions
source ~/.config/zsh/functions.zsh
