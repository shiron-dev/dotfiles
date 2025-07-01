HISTFILE=$ZDOTDIR/.zsh-history
HISTSIZE=100000
SAVEHIST=10000000
HISTORY_IGNORE="nv"
setopt inc_append_history
setopt hist_ignore_dups
setopt share_history
setopt AUTO_CD
setopt AUTO_PARAM_KEYS

is_cursor() {
  [[ "$PAGER" == "head -n 10000 | cat" ]]
}

zshaddhistory() {
    if is_cursor; then
        return 1
    fi
    return 0
}

export LC_ALL="en_US.UTF-8"
