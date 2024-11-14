#!/usr/bin/env bash

cat <<EOM

# shiron-dev dotfiles

Start setup.
For more information check the following link.

https://github.com/shiron-dev/dotfiles

EOM

if [ -d ~/projects/dotfiles ]; then
  cat <<EOM

dotfiles repository already exists.
Run setup from '~/projects/dotfiles'.

EOM
  cd ~/projects/dotfiles/scripts/ && go run ~/projects/dotfiles/scripts/main.go
else
  cat <<EOM

dotfiles repository does not exist.

EOM
  # TODO: Run from github release
  exit 1
fi
