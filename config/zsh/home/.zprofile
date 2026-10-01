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

# Added by OrbStack: command-line tools and integration
# This won't be added again if you remove it.
source ~/.orbstack/shell/init.zsh 2>/dev/null || :
