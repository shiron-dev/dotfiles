# Work machine specific shell setup.
# Sourced at the very end of ~/.zshrc, so PATH tweaks here win.

export DOTFILES_PROFILE=work

# gcloud
if [ -f "$(brew --prefix)/share/google-cloud-sdk/path.zsh.inc" ]; then
  source "$(brew --prefix)/share/google-cloud-sdk/path.zsh.inc"
fi

### MANAGED BY RANCHER DESKTOP START (DO NOT EDIT)
export PATH="$HOME/.rd/bin:$PATH"
### MANAGED BY RANCHER DESKTOP END (DO NOT EDIT)

# Unity CLI
[ -f "$HOME/.unity/env" ] && . "$HOME/.unity/env"

# >>> uloop PATH >>>
export PATH="$HOME/.local/bin:$PATH"
# <<< uloop PATH <<<
