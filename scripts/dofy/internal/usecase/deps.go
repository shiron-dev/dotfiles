package usecase

import (
	"context"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"

	"github.com/shiron-dev/dotfiles/scripts/dofy/internal/infrastructure"
)

type DepsUsecase interface {
	InstallDeps(ctx context.Context) error
	brewBundle() error
}

type DepsUsecaseImpl struct {
	depsInfrastructure infrastructure.DepsInfrastructure
	printOutUC         PrintOutUsecase
	brewUC             BrewUsecase
}

func NewDepsUsecase(
	depsInfrastructure infrastructure.DepsInfrastructure,
	printOutUC PrintOutUsecase,
	brewUC BrewUsecase,
) *DepsUsecaseImpl {
	return &DepsUsecaseImpl{
		depsInfrastructure: depsInfrastructure,
		printOutUC:         printOutUC,
		brewUC:             brewUC,
	}
}

type DepsError struct {
	err error
}

func (e *DepsError) Error() string {
	return "DepsUC: " + e.err.Error()
}

func (d *DepsUsecaseImpl) InstallDeps(ctx context.Context) error {
	d.printOutUC.PrintMdf(`
## Installing Homebrew
`)

	if d.depsInfrastructure.CheckInstalled("brew") {
		d.printOutUC.Println("Homebrew is already installed")
	} else {
		err := d.brewUC.InstallHomebrew(ctx)
		if err != nil {
			return &DepsError{err}
		}
	}

	d.printOutUC.PrintMdf(`
## Installing required packages with Homebrew

- git
`)

	if d.depsInfrastructure.CheckInstalled("git") {
		d.printOutUC.Println("git is already installed")
	} else {
		err := d.brewUC.InstallFormula("git")
		if err != nil {
			return &DepsError{err}
		}
	}

	d.printOutUC.PrintMdf(`
## Git clone dotfiles repository

https://github.com/shiron-dev/dotfiles.git
`)

	usr, _ := user.Current()
	if _, err := os.Stat(usr.HomeDir + "/projects/dotfiles"); err == nil {
		d.printOutUC.Println("dotfiles directory already exists")
	} else {
		d.printOutUC.Println("Cloning dotfiles repository")

		//nolint:gosec
		cmd := exec.Command(
			"git",
			"clone",
			"https://github.com/shiron-dev/dotfiles.git",
			filepath.Join(usr.HomeDir, "/projects/dotfiles"),
		)
		if err := cmd.Run(); err != nil {
			return &DepsError{err}
		}
	}

	return d.brewBundle()
}

func (d *DepsUsecaseImpl) brewBundle() error {
	d.printOutUC.PrintMdf(`
## Installing brew packages

Install the packages using Homebrew Bundle.
`)

	d.printOutUC.PrintMdf(`
### Install brew packages with Brewfile
	`)

	if err := d.brewUC.InstallBrewBundle(); err != nil {
		return &DepsError{err}
	}

	return nil
}

// 	dumpTmpBrewBundle()
// 	diffBundle, diffTmpBundles := checkDiffBrewBundle(
// 		usr.HomeDir+"/projects/dotfiles/data/Brewfile",
// 		usr.HomeDir+"/projects/dotfiles/data/Brewfile.tmp",
// 	)

// 	if len(diffBundle) > 0 {
// 		d.printOutUC.Println("Installing brew packages")
// 		installBrewBundle()
// 	} else {
// 		d.printOutUC.Println("No new packages to install")
// 	}

// 	if len(diffTmpBundles) > 0 {
// 		var diffNames string
// 		for _, diff := range diffTmpBundles {
// 			diffNames += "- " + diff.name + "\n"
// 		}
// 		d.printOutUC.Println(color.RedString("The dotfiles Brewfile and the currently installed package are different."))
// 		d.printOutUC.PrintMdf(`
// ### Update Brewfile

// diff:
// ` + diffNames + `

// What will you do to resolve the diff?

// 1. run ` + "`brew bundlecleanup`" + `
// 2. update the Brewfile with the currently installed packages
// 3. do nothing
// 4. exit
// `)
// 		d.printOutUC.Print("What do you run? [1-4]: ")
// 		scanner := bufio.NewScanner(os.Stdin)
// 		if scanner.Scan() {
// 			switch strings.TrimSpace(scanner.Text()) {
// 			case "1":
// 				d.printOutUC.Println("Running `brew bundle cleanup`")
// 				cleanupBrewBundle(true)
// 			case "2":
// 				d.printOutUC.Println("Open Brewfile with code")
// 				utils.OpenWithCode(
// 					usr.HomeDir+"/projects/dotfiles/data/Brewfile",
// 					usr.HomeDir+"/projects/dotfiles/data/Brewfile.tmp",
// 				)
// 			case "3":
// 				d.printOutUC.Println("Do nothing")
// 			default:
// 				d.printOutUC.Println("Exit")
// 				os.Exit(0)
// 			}
// 		}
// 	}
