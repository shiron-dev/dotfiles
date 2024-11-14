#!/usr/bin/env bash

script_dir=$(dirname "$(realpath "${BASH_SOURCE[0]}")")

cd "$script_dir/.." && docker build -t shiron-dev/dotfiles:latest -f "$script_dir/Dockerfile" .

docker run -it shiron-dev/dotfiles:latest
