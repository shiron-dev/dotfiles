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

function gi() { curl -sLw "\n" https://www.toptal.com/developers/gitignore/api/$@; }

alias gc="ghq get"

function ghq-fzf() {
  local src=$(ghq list | fzf --preview "bat --color=always --style=header,grid --line-range :80 $(ghq root)/{}/README.*")
  if [ -n "$src" ]; then
    BUFFER="cd $(ghq root)/$src"
    zle accept-line
  fi
  zle -R -c
}
zle -N ghq-fzf
bindkey '^]' ghq-fzf
alias pj='ghq-fzf'

alias o.='open .'
alias c.='code .'

function cg(){
  cd "$(git rev-parse --show-toplevel)"
}