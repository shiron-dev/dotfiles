# Amazon Q pre block. Keep at the top of this file.
[[ -f "${HOME}/Library/Application Support/amazon-q/shell/zprofile.pre.zsh" ]] && builtin source "${HOME}/Library/Application Support/amazon-q/shell/zprofile.pre.zsh"
# Q pre block. Keep at the top of this file.
HISTFILE=$ZDOTDIR/.zsh-history
HISTSIZE=100000
SAVEHIST=10000000
setopt inc_append_history
setopt share_history
setopt AUTO_CD
setopt AUTO_PARAM_KEYS

# iTerm2
bindkey "^[[H" beginning-of-line
bindkey "^[[F" end-of-line

# Homebrew
eval "$(/opt/homebrew/bin/brew shellenv)"

# Added by Toolbox App
export PATH="$PATH:/Users/shiron/Library/Application Support/JetBrains/Toolbox/scripts"

# gcc
alias gcc="/opt/homebrew/Cellar/gcc/13.1.0/bin/gcc-13"
alias g++="/opt/homebrew/Cellar/gcc/13.1.0/bin/g++-13"

# python
alias python=python3
alias pip=pip3

export LC_ALL=en_US.UTF-8
  
alias lsusb=system_profiler SPUSBDataType

# terraform -install-autocomplete
# autoload -U +X bashcompinit && bashcompinit
# complete -o nospace -C /opt/homebrew/bin/terraform terraform

export CPPFLAGS="-I/opt/homebrew/opt/openjdk@17/include"

export PATH="$PATH:/Users/shiron/Library/Android/sdk/platform-tools"

source ~/.config/peco.sh
export PATH="$PATH":"$HOME/.pub-cache/bin"

export SDKMAN_DIR=$(brew --prefix sdkman-cli)/libexec
[[ -s "${SDKMAN_DIR}/bin/sdkman-init.sh" ]] && source "${SDKMAN_DIR}/bin/sdkman-init.sh"


export PATH="$PATH":"/Users/shiron/projects/tools/bin"

alias docker-compose-rm="docker compose down --rmi all --volumes --remove-orphans"
alias tailscale='/Applications/Tailscale.app/Contents/MacOS/Tailscale'

# Q post block. Keep at the bottom of this file.
export GOPATH=$HOME/go
export GOBIN=$GOPATH/bin
export PATH=$PATH:$GOBIN


_notify-done() {
  afplay /System/Library/Sounds/Hero.aiff
  osascript -e 'display notification "Process is done!" with title "Terminal"'
}
alias notify=_notify-done

# Amazon Q post block. Keep at the bottom of this file.
[[ -f "${HOME}/Library/Application Support/amazon-q/shell/zprofile.post.zsh" ]] && builtin source "${HOME}/Library/Application Support/amazon-q/shell/zprofile.post.zsh"

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

eval $(thefuck --alias)
