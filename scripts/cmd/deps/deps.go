package deps

import (
	"github.com/shiron-dev/dotfiles/scripts/cmd/printout"
)

func InstallDeps() {
	printout.PrintMd(`
## Installing dependencies

- Homebrew
`)

	if !checkInstalled("brew") {
		printout.PrintMd("### Installing Homebrew")
		installHomebrew()
	}
}
