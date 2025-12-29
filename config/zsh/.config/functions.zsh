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

function gitbc() { git checkout -q main && git for-each-ref refs/heads/ "--format=%(refname:short)" | while read branch; do mergeBase=$(git merge-base main $branch) && [[ $(git cherry main $(git commit-tree $(git rev-parse "$branch^{tree}") -p $mergeBase -m _)) == "-"* ]] && git branch -D $branch; done; }

alias gc="ghq get"

function _ghq-fzf() {
  local src=$(ghq list | fzf --preview "bat --color=always --style=header,grid --line-range :80 $(ghq root)/{}/README.*")
  if [ -n "$src" ]; then
    BUFFER="cursor $(ghq root)/$src"
    zle accept-line
  fi
  zle -R -c
}

function ghq-fzf() {
  local src=$(ghq list | fzf --preview "bat --color=always --style=header,grid --line-range :80 $(ghq root)/{}/README.*")
  if [ -n "$src" ]; then
    cd $(ghq root)/$src || exit
  fi
}

zle -N _ghq-fzf
bindkey '^]' _ghq-fzf
alias pj='ghq-fzf'

alias o.='open .'
alias c.='code .'
alias ci.='code-insiders .'
alias cu.='cursor .'
alias i.='idea .'
alias g.='goland .'

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

asn() {
  curl "ipinfo.io/$1/org"
}

zle -N _navi_widget
bindkey '^g' _navi_widget

git-todo() {
  local tool="${1:-code}"
  shift
  git log --author="$(git config user.name)" --name-only --pretty=format:"" | sort -u | xargs git grep -l "TODO" | xargs "$tool" "$@"
}

docker-ssh() {
  local target="${1:-$(
    docker ps --format "{{.ID}}\t{{.Names}}\t{{.Image}}" | 
    fzf --height 40% --reverse | 
    awk '{print $1}'
  )}"

  [[ -z "$target" ]] && return

  docker exec -it "$target" sh -c "
    if command -v bash >/dev/null 2>&1; then
      exec bash
    elif command -v sh >/dev/null 2>&1; then
      exec sh
    else
      echo 'Error: No shell found.' && exit 1
    fi
  "
}

docker-up() {
  local image="${1:-$(
    docker images --format "{{.Repository}}:{{.Tag}}\t{{.ID}}" |
    fzf --height 40% --reverse --prompt="Select an image > " |
    awk '{print $2}'
  )}"
  [[ -z "$image" ]] && return

  docker run -it "$image" sh -c "
    if command -v bash >/dev/null 2>&1; then exec bash;
    elif command -v sh >/dev/null 2>&1; then exec sh;
    else echo 'Error: No shell found.' && exit 1; fi
  "

  local running_containers=$(docker ps --filter "ancestor=$image" --format "{{.ID}}")
  if [[ -n "$running_containers" ]]; then
    echo -e "\033[31mNotice: Containers from $image are still running/exist:\033[0m"
    echo "$running_containers" | xargs -I {} echo "  - {} (Use 'docker rm --force {}' to remove)"
  fi
}

docker-stop() {
  local ids=$(
    docker ps --format "{{.ID}}\t{{.Names}}\t{{.Status}}" |
    fzf -m --height 40% --reverse --prompt="STOP container(s) > " |
    awk '{print $1}'
  )

  [[ -n "$ids" ]] && echo "$ids" | xargs docker stop
}

docker-rm() {
  local ids=$(
    docker ps -a --format "{{.ID}}\t{{.Names}}\t{{.Status}}" |
    fzf -m --height 40% --reverse --prompt="REMOVE container(s) > " |
    awk '{print $1}'
  )

  [[ -n "$ids" ]] && echo "$ids" | xargs docker rm --force
}

docker-copy() {
  local mode="container"
  local docker_cmd="ps -a"
  local format="{{.ID}}\t{{.Names}}\t{{.Image}}"
  local prompt="Copy Container ID > "

  if [[ "$1" == "-i" || "$1" == "--image" ]]; then
    mode="image"
    docker_cmd="images"
    format="{{.ID}}\t{{.Repository}}:{{.Tag}}"
    prompt="Copy Image ID > "
  fi

  local id=$(
    docker $docker_cmd --format "$format" |
    fzf --height 40% --reverse --prompt="$prompt" |
    awk '{print $1}'
  )

  if [[ -n "$id" ]]; then
    echo -n "$id" | pbcopy
    echo "Copied ${mode} ID: $id"
  fi
}

docker-logs() {
  local follow=""
  [[ "$1" == "-f" ]] && follow="-f"

  local target=$(
    docker ps -a --format "{{.ID}}\t{{.Names}}\t{{.Status}}" |
    fzf --height 40% --reverse --prompt="Select Container for logs ${follow} > " |
    awk '{print $1}'
  )

  [[ -z "$target" ]] && return

  docker logs $follow "$target"
}

docker-f() {
  local choices=(
    "docker-ssh      : Login to container"
    "docker-up       : Start container from image"
    "docker-logs     : View logs"
    "docker-logs -f  : Stream logs"
    "docker-stop     : Stop container"
    "docker-rm       : Remove container"
    "docker-copy     : Copy container ID"
    "docker-copy -i  : Copy image ID"
  )

  local selected=$(
    printf "%s\n" "${choices[@]}" | 
    fzf --height 40% --reverse --prompt="Docker Utils > "
  )

  [[ -z "$selected" ]] && return

  local cmd=$(echo "$selected" | cut -d ':' -f1 | xargs)
  
  echo "Running: $cmd"
  eval "$cmd"
}

function git-trim-eof-newlines() {
  local all_files=false
  if [[ "$1" == "-a" || "$1" == "--all" ]]; then
    all_files=true
  fi

  local files
  if [[ "$all_files" == true ]]; then
    files=("${(@f)$(git ls-files)}")
  else
    files=("${(@f)$(git diff --name-only && git diff --cached --name-only | sort -u)}")
  fi

  if [[ -z "$files" ]]; then
    if [[ "$all_files" == true ]]; then
      echo "No files managed by git found."
    else
      echo "No changed files found."
    fi
    return 0
  fi

  for file in $files; do
    if [[ -f "$file" ]]; then
      if file "$file" | grep -q "text"; then
        perl -i -0777 -pe 's/\n+\z/\n/' "$file"
        echo "Trimmed EOF: $file"
      fi
    fi
  done

  echo "Done."
}
alias cai='git-trim-eof-newlines'
