package main

import (
	"github.com/shiron-dev/dotfiles/scripts/cmd/deps"
	"github.com/shiron-dev/dotfiles/scripts/cmd/printout"
)

func main() {
	printout.PrintMd(`

# shiron-dev dotfiles setup script

This script will install dependencies and setup dotfiles.

`)
	deps.InstallDeps()
}
