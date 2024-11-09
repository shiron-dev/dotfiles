package usecase

import (
	"bufio"
	"context"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/pkg/errors"
	"github.com/shiron-dev/dotfiles/scripts/dofy/internal/domain"
	"github.com/shiron-dev/dotfiles/scripts/dofy/internal/infrastructure"
)

type DepsUsecase interface {
	InstallHomebrew(ctx context.Context) error
	InstallGit() error
	CloneDotfiles() error
	InstallBrewBundle() error

	showBrewDiff(diffBundles []domain.BrewBundle, diffTmpBundles []domain.BrewBundle) error
	updateBrewfile() error
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

func (d *DepsUsecaseImpl) InstallHomebrew(ctx context.Context) error {
	d.printOutUC.PrintMdf(`
## Installing Homebrew
`)

	if d.depsInfrastructure.CheckInstalled("brew") {
		d.printOutUC.Println("Homebrew is already installed")
	} else {
		err := d.brewUC.InstallHomebrew(ctx)
		if err != nil {
			return errors.Wrap(err, "deps usecase: failed to install Homebrew")
		}
	}

	return nil
}

func (d *DepsUsecaseImpl) InstallGit() error {
	d.printOutUC.PrintMdf(`
## Installing required packages with Homebrew

- git
`)

	if d.depsInfrastructure.CheckInstalled("git") {
		d.printOutUC.Println("git is already installed")
	} else {
		err := d.brewUC.InstallFormula("git")
		if err != nil {
			return errors.Wrap(err, "deps usecase: failed to install git")
		}
	}

	return nil
}

func (d *DepsUsecaseImpl) CloneDotfiles() error {
	d.printOutUC.PrintMdf(`
## Git clone dotfiles repository

https://github.com/shiron-dev/dotfiles.git
`)

	usr, err := user.Current()
	if err != nil {
		return errors.Wrap(err, "deps usecase: failed to get current user")
	}

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
			return errors.Wrap(err, "deps usecase: failed to clone dotfiles repository")
		}
	}

	return nil
}

func (d *DepsUsecaseImpl) InstallBrewBundle() error {
	usr, err := user.Current()
	if err != nil {
		return errors.Wrap(err, "deps usecase: failed to get current user")
	}

	d.printOutUC.PrintMdf(`
## Installing brew packages

Install the packages using Homebrew Bundle.
`)

	err = d.brewUC.DumpTmpBrewBundle()
	if err != nil {
		return errors.Wrap(err, "deps usecase: failed to dump tmp Brewfile")
	}

	diffBundles, diffTmpBundles, err := d.brewUC.CheckDiffBrewBundle(
		usr.HomeDir+"/projects/dotfiles/data/Brewfile",
		usr.HomeDir+"/projects/dotfiles/data/Brewfile.tmp",
	)
	if err != nil {
		return errors.Wrap(err, "deps usecase: failed to check diff Brewfile")
	}

	if len(diffBundles) > 0 {
		d.printOutUC.PrintMdf(`
### Install brew packages with Brewfile
`)

		if err := d.brewUC.InstallBrewBundle(); err != nil {
			return errors.Wrap(err, "deps usecase: failed to install brew packages")
		}
	} else {
		d.printOutUC.Println("No new packages to install")
	}

	if len(diffTmpBundles) > 0 {
		err := d.showBrewDiff(diffBundles, diffTmpBundles)
		if err != nil {
			return errors.Wrap(err, "deps usecase: failed to update Brewfile")
		}
	}

	return nil
}

func (d *DepsUsecaseImpl) showBrewDiff(diffBundles []domain.BrewBundle, diffTmpBundles []domain.BrewBundle) error {
	var diffNames string
	for _, diff := range diffTmpBundles {
		diffNames += color.GreenString("+ " + diff.Name + "\n")
	}

	for _, diff := range diffBundles {
		diffNames += color.RedString("- " + diff.Name + "\n")
	}

	d.printOutUC.Println(color.RedString("The dotfiles Brewfile and the currently installed package are different."))
	d.printOutUC.PrintMdf(`
### Update Brewfile

diff:
` + diffNames + `

What will you do to resolve the diff?

1. run ` + "`brew bundlecleanup`" + `
2. update the Brewfile with the currently installed packages
3. do nothing
4. exit
`)
	d.printOutUC.Print("What do you run? [1-4]: ")

	return d.updateBrewfile()
}

func (d *DepsUsecaseImpl) updateBrewfile() error {
	usr, err := user.Current()
	if err != nil {
		return errors.Wrap(err, "deps usecase: failed to get current user")
	}

	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		switch strings.TrimSpace(scanner.Text()) {
		case "1":
			d.printOutUC.Println("Running `brew bundle cleanup`")

			err := d.brewUC.CleanupBrewBundle(true)
			if err != nil {
				return errors.Wrap(err, "deps usecase: failed to run brew bundle cleanup")
			}
		case "2":
			d.printOutUC.Println("Open Brewfile with code")

			if err := d.depsInfrastructure.OpenWithCode(
				usr.HomeDir+"/projects/dotfiles/data/Brewfile",
				usr.HomeDir+"/projects/dotfiles/data/Brewfile.tmp",
			); err != nil {
				return errors.Wrap(err, "deps usecase: failed to open with code")
			}
		case "3":
			d.printOutUC.Println("Do nothing")
		default:
			d.printOutUC.Println("Exit")
			os.Exit(0)
		}
	}

	return nil
}