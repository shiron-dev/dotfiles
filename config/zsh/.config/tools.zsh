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

# My tools
export PATH="$PATH":"/Users/shiron/projects/tools/bin"

# My Aliases
alias docker-compose-rm="docker compose down --rmi all --volumes --remove-orphans"
alias lsusb="system_profiler SPUSBDataType"

# My Functions
_notify() {
  afplay /System/Library/Sounds/Hero.aiff
  osascript -e 'display notification with title "Terminal"'
}
alias notify=_notify

function moveToTrash() {
  local p
  for p in "$@"; do
    if [[ "$p" == -* ]]; then
      continue
    fi

    if [ -e "$p" ]; then
      date=$(/usr/bin/env date "+%Y-%m-%d_%H-%M-%S")
      /usr/bin/env mkdir -p ~/.Trash/$p_$date/..
      /usr/bin/env mv "$p" ~/.Trash/$p_$date
    else
      /usr/bin/env echo "Error: '$p' does not exist."
    fi
  done
}
alias rm='moveToTrash'

function gi() { curl -sLw "\n" https://www.toptal.com/developers/gitignore/api/$@ ;}
