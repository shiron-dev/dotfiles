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

# volta
export VOLTA_HOME="$HOME/.volta"
export PATH="$VOLTA_HOME/bin:$PATH"

# pnpm
export PNPM_HOME="/Users/shiron/Library/pnpm"
case ":$PATH:" in
  *":$PNPM_HOME:"*) ;;
  *) export PATH="$PNPM_HOME:$PATH" ;;
esac

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

# Golang
export GOPATH=$HOME/go
export GOBIN=$GOPATH/bin
export PATH=$PATH:$GOBIN

# lazy
alias lg='lazygit'
alias ld='lazydocker'

# fzf
source <(fzf --zsh)

# My tools
export PATH="$PATH":"/Users/shiron/projects/tools/bin"

if [ ! -f "/Users/shiron/projects/github.com/shiron-dev/dotfiles/scripts/dofy/dofy" ]; then
  echo "Building dofy..."
  (cd ~/projects/github.com/shiron-dev/dotfiles/scripts/dofy && go build -o dofy cmd/main.go)
fi
alias dofy="/Users/shiron/projects/github.com/shiron-dev/dotfiles/scripts/dofy/dofy"

# My Aliases
alias docker-compose-rm="docker compose down --rmi all --volumes --remove-orphans"
alias lsusb="system_profiler SPUSBDataType"

# My functions
source ~/.config/zsh/functions.zsh
