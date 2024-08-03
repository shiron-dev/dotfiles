package deps

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strings"

	"github.com/fatih/color"
	"github.com/shiron-dev/dotfiles/scripts/cmd/printout"
	"github.com/shiron-dev/dotfiles/scripts/cmd/utils"
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

	dumpTmpBrewBundle()
	diffBundle, diffTmpBundles := checkDiffBrewBundle(
		usr.HomeDir+"/projects/dotfiles/data/Brewfile",
		usr.HomeDir+"/projects/dotfiles/data/Brewfile.tmp",
	)
	if len(diffBundle) > 0 {
		printout.Println("Installing brew packages")
		installBrewBundle()
	} else {
		printout.Println("No new packages to install")
	}
	if len(diffTmpBundles) > 0 {
		var diffNames string
		for _, diff := range diffTmpBundles {
			diffNames += "- " + diff.name + "\n"
		}
		fmt.Println(color.RedString("The dotfiles Brewfile and the currently installed package are different."))
		printout.PrintMd(`
### Update Brewfile

diff:
` + diffNames + `

What will you do to resolve the diff?

1. run ` + "`brew bundlecleanup`" + `
2. update the Brewfile with the currently installed packages
3. do nothing
4. exit
`)
		fmt.Print("What do you run? [1-4]: ")
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			switch strings.TrimSpace(scanner.Text()) {
			case "1":
				printout.Println("Running `brew bundle cleanup`")
				cleanupBrewBundle(true)
			case "2":
				printout.Println("Open Brewfile with code")
				utils.OpenWithCode(
					usr.HomeDir+"/projects/dotfiles/data/Brewfile",
					usr.HomeDir+"/projects/dotfiles/data/Brewfile.tmp",
				)
			case "3":
				printout.Println("Do nothing")
			default:
				printout.Println("Exit")
				os.Exit(0)
			}
		}
	}

	printout.PrintMd(`
### Install brew packages with Brewfile
	`)
	installBrewBundle()
}
