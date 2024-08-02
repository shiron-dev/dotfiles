package deps

import (
	"os"
	"os/exec"
	"os/user"

	"github.com/shiron-dev/dotfiles/scripts/cmd/printout"
)

func InstallDeps() {
	printout.PrintMd(`
## Installing dependencies

- Homebrew
`)

	printout.PrintMd("### Installing Homebrew")
	if checkInstalled("brew") {
		printout.Println("Homebrew is already installed")
	} else {
		printout.Println("Installing Homebrew")
		installHomebrew()
	}

	printout.PrintMd(`
## Installing required packages with Homebrew

- git
`)

	printout.PrintMd("### Installing git")
	if checkInstalled("git") {
		printout.Println("git is already installed")
	} else {
		printout.Println("Installing git")
		installWithBrew("git")
	}

	printout.PrintMd(`
## Git clone dotfiles repository

https://github.com/shiron-dev/dotfiles.git
`)

	usr, _ := user.Current()
	if _, err := os.Stat(usr.HomeDir + "/projects/dotfiles"); err == nil {
		printout.Println("dotfiles directory already exists")
	} else {
		printout.Println("Cloning dotfiles repository")
		cmd := exec.Command("git", "clone", "https://github.com/shiron-dev/dotfiles.git", usr.HomeDir+"/projects/dotfiles")
		cmd.Run()
	}

	printout.PrintMd(`
## Installing brew packages

Install the packages using Homebrew Bundle.
`)

	installBrewBundle()
}
