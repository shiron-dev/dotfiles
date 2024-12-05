eval "$(sheldon source)"

export DYLD_LIBRARY_PATH="/opt/homebrew/lib/"

export PATH="/opt/homebrew/opt/ruby/bin:$PATH"

eval "$(/opt/homebrew/bin/mise activate zsh)"

export VOLTA_HOME="$HOME/.volta"
export PATH="$VOLTA_HOME/bin:$PATH"
export PATH="/usr/local/opt/curl/bin:$PATH"

# pnpm
export PNPM_HOME="/Users/shiron/Library/pnpm"
case ":$PATH:" in
  *":$PNPM_HOME:"*) ;;
  *) export PATH="$PNPM_HOME:$PATH" ;;
esac
# pnpm end

function gi() { curl -sLw "\n" https://www.toptal.com/developers/gitignore/api/$@ ;}

alias sed='gsed'

# Amazon Q post block. Keep at the bottom of this file.
[[ -f "${HOME}/Library/Application Support/amazon-q/shell/zshrc.post.zsh" ]] && builtin source "${HOME}/Library/Application Support/amazon-q/shell/zshrc.post.zsh"

export PATH="/opt/homebrew/opt/mysql-client@8.0/bin:$PATH"
export PATH="/opt/homebrew/opt/gawk/libexec/gnubin:$PATH"

ZSH_TIME=$(/opt/homebrew/bin/gdate +%s%3N)
DIFF=$(echo "$ZSH_TIME - $ZSH_STARTUP_TIME" | bc)
echo zsh startup time $DIFF ms

if (which zprof > /dev/null 2>&1) ;then
  zprof
fi
