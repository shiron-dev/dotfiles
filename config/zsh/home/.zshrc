HISTFILE=$HOME/.zsh-history
HISTSIZE=100000
SAVEHIST=10000000
HISTORY_IGNORE="nv"
setopt inc_append_history
setopt hist_ignore_dups
setopt share_history
setopt AUTO_CD
setopt AUTO_PARAM_KEYS

source ~/.config/zsh/tools.zsh

# Startup time
ZSH_TIME=$(/opt/homebrew/bin/gdate +%s%3N)
DIFF=$(echo "$ZSH_TIME - $ZSH_STARTUP_TIME" | bc)
echo zsh startup time "$DIFF" ms

if (which zprof > /dev/null 2>&1) ;then
  zprof
fi
