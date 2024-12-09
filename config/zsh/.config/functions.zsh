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
      /usr/bin/env mkdir -p ~/.Trash/"$p_""$date"/..
      /usr/bin/env mv "$p" ~/.Trash/"$p_""$date"
    else
      /usr/bin/env echo "Error: '$p' does not exist."
    fi
  done
}
alias rm='moveToTrash'

function gi() { curl -sLw "\n" https://www.toptal.com/developers/gitignore/api/$@; }

alias gc="ghq get"

function _ghq-fzf() {
  local src=$(ghq list | fzf --preview "bat --color=always --style=header,grid --line-range :80 $(ghq root)/{}/README.*")
  if [ -n "$src" ]; then
    BUFFER="cd $(ghq root)/$src"
    zle accept-line
  fi
  zle -R -c
}
function ghq-fzf() {
  local src=$(ghq list | fzf --preview "bat --color=always --style=header,grid --line-range :80 $(ghq root)/{}/README.*")
  if [ -n "$src" ]; then
    cd $(ghq root)/"$src" || exit
  fi
}
zle -N _ghq-fzf
bindkey '^]' _ghq-fzf
alias pj='ghq-fzf'

alias o.='open .'
alias c.='code .'

function cg() {
  cd "$(git rev-parse --show-toplevel)" || exit
}

function y() {
  if [ "$1" != "" ]; then
    if [ -d "$1" ]; then
      yazi "$1"
    else
      yazi "$(zoxide query "$1")"
    fi
  else
    yazi
  fi
  return $?
}

function yazi-cd() {
  local tmp="$(mktemp -t "yazi-cwd.XXXXXX")" cwd
  yazi "$@" --cwd-file="$tmp"
  if cwd="$(command cat -- "$tmp")" && [ -n "$cwd" ] && [ "$cwd" != "$PWD" ]; then
    builtin cd -- "$cwd" || exit
  fi
  rm -f -- "$tmp"
}

alias yazi='yazi-cd'

alias ghb='gh browse'

_navi_call() {
  local result="$(navi "$@" </dev/tty)"
  printf "%s" "$result"
}

_navi_widget() {
  local -r input="${LBUFFER}"
  local -r last_command="$(echo "${input}" | navi fn widget::last_command)"
  local replacement="$last_command"

  if [ -z "$last_command" ]; then
    replacement="$(_navi_call --print)"
  elif [ "$LASTWIDGET" = "_navi_widget" ] && [ "$input" = "$previous_output" ]; then
    replacement="$(_navi_call --print --query "$last_command")"
  else
    replacement="$(_navi_call --print --best-match --query "$last_command")"
  fi

  if [ -n "$replacement" ]; then
    local -r find="${last_command}_NAVIEND"
    previous_output="${input}_NAVIEND"
    previous_output="${previous_output//$find/$replacement}"
  else
    previous_output="$input"
  fi

  zle kill-whole-line
  LBUFFER="${previous_output}"
  region_highlight=("P0 100 bold")
  zle redisplay
}
nv() {
  local -r result="$(_navi_call --print)"
  if [ -n "$result" ]; then
    print -z "$result"

    builtin fc -RW
  fi
}

zle -N _navi_widget
bindkey '^g' _navi_widget
