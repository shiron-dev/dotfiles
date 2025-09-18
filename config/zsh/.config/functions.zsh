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
    BUFFER="code $(ghq root)/$src"
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

# ~/.zshrc に以下を追記または修正
ssh-docker() {
  # 実行するシェルのリスト
  local shells_to_try=("/bin/bash" "/bin/sh")
  local shell_to_exec=""

  # 実行対象のコンテナIDを格納する変数
  local target_container=""

  if [ -n "$1" ]; then
    # --- 引数が指定されている場合 ---
    target_container=$1
  else
    # --- 引数が指定されていない場合 (fzfで選択) ---
    target_container=$(docker ps -a --format "table {{.ID}}\t{{.Names}}\t{{.Image}}\t{{.Status}}" | sed 1d | fzf --height 40% --layout=reverse --prompt="Select a container > " | awk '{print $1}')
  fi

  # コンテナが選択されている場合のみ実行
  if [ -n "$target_container" ]; then
    # コンテナ内で利用可能なシェルを探す
    for shell in "${shells_to_try[@]}"; do
      if docker exec -i "$target_container" test -x "$shell" 2>/dev/null; then
        shell_to_exec=$shell
        break
      fi
    done

    # 利用可能なシェルがあれば実行、なければエラーメッセージ
    if [ -n "$shell_to_exec" ]; then
      docker exec -it "$target_container" "$shell_to_exec"
    else
      echo "Error: Could not find a valid shell (/bin/bash or /bin/sh) in the container."
      return 1
    fi
  fi
}
# ~/.zshrc に以下を追記
docker-up() {
  local image
  # docker images の結果から fzf でイメージを選択し、IMAGE ID を取得
  image=$(docker images --format "table {{.Repository}}:{{.Tag}}\t{{.ID}}" | sed 1d | fzf --height 40% --layout=reverse --prompt="Select an image to run > " | awk '{print $2}')

  # イメージが選択された場合のみコンテナを起動
  if [ -n "$image" ]; then
    echo "Starting a new container from image: $image"
    docker run -it "$image" /bin/sh -c "([ -x /bin/bash ] && exec /bin/bash) || exec /bin/sh"

    # Show in red if the container is not stopped, and guide how to stop it
    docker ps -a --filter "id=$image" --format "{{.ID}}\t{{.Status}}" | while read id status; do
      if [[ "$status" != *"Exited"* ]]; then
      echo -e "\033[31mContainer is not stopped. You can stop it with: docker rm --force $id\033[0m"
      fi
    done
  fi
}
# ~/.zshrc に以下を追記
# docker-stop: Select and stop one or more running containers
docker-stop() {
  local containers
  containers=$(docker ps --format "table {{.ID}}\t{{.Names}}\t{{.Image}}\t{{.Status}}" | sed 1d | fzf -m --height 40% --layout=reverse --prompt="Select container(s) to STOP > " | awk '{print $1}')

  if [ -n "$containers" ]; then
    echo "$containers" | xargs docker stop
  fi
}

# docker-rm: Select and forcibly remove one or more containers (stops them if running)
docker-rm() {
  local containers
  containers=$(docker ps -a --format "table {{.ID}}\t{{.Names}}\t{{.Image}}\t{{.Status}}" | sed 1d | fzf -m --height 40% --layout=reverse --prompt="Select container(s) to REMOVE > " | awk '{print $1}')

  if [ -n "$containers" ]; then
    echo "$containers" | xargs docker rm --force
  fi
}
# docker-c: fzfでコンテナIDをクリップボードにコピー
function docker-c() {
  local container_id
  container_id=$(docker ps -a --format "{{.ID}}\t{{.Names}}\t{{.Image}}\t{{.Status}}" | fzf --height 40% --reverse)

  if [[ -n "$container_id" ]]; then
    # fzfで選択した行からID部分（先頭の単語）のみを抽出
    echo -n "$container_id" | awk '{print $1}' | pbcopy
    echo "Copied container ID to clipboard: $(echo -n "$container_id" | awk '{print $1}')"
  else
    echo "No container selected."
  fi
}
# docker-i: fzfでイメージIDをクリップボードにコピー
function docker-i() {
  local image_id
  image_id=$(docker images --format "{{.ID}}\t{{.Repository}}:{{.Tag}}\t{{.Size}}" | fzf --height 40% --reverse)

  if [[ -n "$image_id" ]]; then
    # fzfで選択した行からID部分（先頭の単語）のみを抽出
    echo -n "$image_id" | awk '{print $1}' | pbcopy
    echo "Copied image ID to clipboard: $(echo -n "$image_id" | awk '{print $1}')"
  else
    echo "No image selected."
  fi
}

# ##########################
# sekai-ssh
autoload -Uz compinit ; compinit # なかったら追記

function sekai-ssh() {
    ip=`aws ec2 describe-instances --output=text --filters "Name=tag:Name,Values=$1" --query "Reservations[].Instances[].PublicIpAddress"`
    ssh s28628@$ip
}

_refresh-sekai-ssh-hosts() {
    rm -f ~/.ssh/_sekai-ssh.hosts
    aws ec2 describe-tags --output text --filters 'Name=resource-type, Values=instance' 'Name=key, Values=Name' --query 'Tags[*].Value' | tr '\t' '\n' | sort > ~/.ssh/_sekai-ssh.hosts
}

 _sekai-ssh() {
    local -a hosts
    hosts=( ${(f)"$(<~/.ssh/_sekai-ssh.hosts)"} )
    _values 'hosts' "${hosts[@]}"
}

compdef _sekai-ssh sekai-ssh
