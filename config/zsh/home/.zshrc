source ~/.config/zsh/tools.zsh

# Startup time
ZSH_TIME=$(/opt/homebrew/bin/gdate +%s%3N)
DIFF=$(echo "$ZSH_TIME - $ZSH_STARTUP_TIME" | bc)
echo zsh startup time "$DIFF" ms

if (which zprof > /dev/null 2>&1) ;then
  zprof
fi
