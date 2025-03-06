#!/usr/bin/env bash

cat <<EOM

# shiron-dev dotfiles

Start setup.
For more information check the following link.

https://github.com/shiron-dev/dotfiles

EOM

if [ -d ~/projects/github.com/shiron-dev/dotfiles ]; then
  cat <<EOM

dotfiles repository already exists.
Run setup from '~/projects/github.com/shiron-dev/dotfiles'.

EOM
  cd ~/projects/github.com/shiron-dev/dotfiles/scripts/dofy/ && go run ~/projects/github.com/shiron-dev/dotfiles/scripts/dofy/cmd/main.go
else
  cat <<EOM

dotfiles repository does not exist.

EOM
  # TODO: Run from github release
  exit 1
fi
