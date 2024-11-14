package usecase

import (
	"bufio"
	"context"
	"os"
	"os/exec"
	"os/signal"
	"os/user"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/fatih/color"
	"github.com/pkg/errors"
	"github.com/shiron-dev/dotfiles/scripts/dofy/internal/domain"
	"github.com/shiron-dev/dotfiles/scripts/dofy/internal/infrastructure"
)

const dofyBrewCategory = "Add by dofy"

const resolveBrewDiffWithEditorMaxCount = 3

var errResolveBrewDiffWithEditorMaxCount = errors.New("resolve brew diff with editor max count error")

type DepsUsecase interface {
	InstallHomebrew(ctx context.Context) error
	InstallGit() error
	CloneDotfiles() error
	InstallBrewBundle() error

	showBrewDiff(diffBundles []domain.BrewBundle, diffTmpBundles []domain.BrewBundle) error
	updateBrewfile() error
	resolveBrewDiff() error
	resolveBrewDiffWithEditor(ctx context.Context) error
	mergeDiff(base []domain.BrewBundle, add []domain.BrewBundle, sub []domain.BrewBundle) []domain.BrewBundle
}

type DepsUsecaseImpl struct {
	depsInfrastructure infrastructure.DepsInfrastructure
	brewInfrastructure infrastructure.BrewInfrastructure
	fileInfrastructure infrastructure.FileInfrastructure
	gitInfrastructure  infrastructure.GitInfrastructure
	printOutUC         PrintOutUsecase
	brewUC             BrewUsecase

	resolveBrewDiffWithEditorCount int
}

func NewDepsUsecase(
	depsInfrastructure infrastructure.DepsInfrastructure,
	brewInfrastructure infrastructure.BrewInfrastructure,
	fileInfrastructure infrastructure.FileInfrastructure,
	gitInfrastructure infrastructure.GitInfrastructure,
	printOutUC PrintOutUsecase,
	brewUC BrewUsecase,
) *DepsUsecaseImpl {
	return &DepsUsecaseImpl{
		depsInfrastructure:             depsInfrastructure,
		brewInfrastructure:             brewInfrastructure,
		fileInfrastructure:             fileInfrastructure,
		gitInfrastructure:              gitInfrastructure,
		printOutUC:                     printOutUC,
		brewUC:                         brewUC,
		resolveBrewDiffWithEditorCount: 0,
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

	if len(diffTmpBundles) > 0 {
		err := d.showBrewDiff(diffBundles, diffTmpBundles)
		if err != nil {
			return errors.Wrap(err, "deps usecase: failed to update Brewfile")
		}
	}

	d.printOutUC.PrintMdf(`
### Install brew packages with Brewfile
`)

	if err := d.brewUC.InstallBrewBundle(); err != nil {
		return errors.Wrap(err, "deps usecase: failed to install brew packages")
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

1. update the Brewfile with the currently installed packages
2. run ` + "`brew bundlecleanup`" + `
3. do nothing
4. exit
`)
	d.printOutUC.Print("What do you run? [1-4]: ")

	return d.updateBrewfile()
}

func (d *DepsUsecaseImpl) updateBrewfile() error {
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		switch strings.TrimSpace(scanner.Text()) {
		case "1":
			d.printOutUC.PrintMdf("#### Open Brewfile with code\n")

			if err := d.resolveBrewDiff(); err != nil {
				return errors.Wrap(err, "deps usecase: failed to resolve Brewfile diff")
			}

			d.printOutUC.PrintMdf("#### Running `brew bundle install`\n")

			if err := d.brewUC.InstallBrewBundle(); err != nil {
				return errors.Wrap(err, "deps usecase: failed to install brew packages")
			}

			d.printOutUC.PrintMdf("#### Running `brew bundle cleanup`\n")

			if err := d.brewUC.CleanupBrewBundle(true); err != nil {
				return errors.Wrap(err, "deps usecase: failed to run brew bundle cleanup")
			}

		case "2":
			d.printOutUC.Println("Running `brew bundle cleanup`")

			if err := d.brewUC.CleanupBrewBundle(true); err != nil {
				return errors.Wrap(err, "deps usecase: failed to run brew bundle cleanup")
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

func (d *DepsUsecaseImpl) resolveBrewDiff() error {
	usr, err := user.Current()
	if err != nil {
		return errors.Wrap(err, "deps usecase: failed to get current user")
	}

	brewPath := usr.HomeDir + "/projects/dotfiles/data/Brewfile"

	diffBundles, diffTmpBundles, err := d.brewUC.CheckDiffBrewBundle(brewPath,
		usr.HomeDir+"/projects/dotfiles/data/Brewfile.tmp")
	if err != nil {
		return errors.Wrap(err, "deps usecase: failed to check diff Brewfile")
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	endFlag := false

	defer func() {
		endFlag = true

		stop()
	}()

	go func() {
		<-ctx.Done()

		if endFlag {
			return
		}

		d.printOutUC.PrintMdf(`
> [!WARNING]
> The Brewfile changes have been discarded.
`)

		if err := d.gitInfrastructure.CheckoutFile(brewPath); err != nil {
			panic(err)
		}
	}()

	bundles, err := d.brewInfrastructure.ReadBrewBundle(brewPath)
	if err != nil {
		return errors.Wrap(err, "deps usecase: failed to read Brewfile")
	}

	err = d.brewInfrastructure.WriteBrewBundle(d.mergeDiff(bundles, diffTmpBundles, diffBundles), brewPath)
	if err != nil {
		return errors.Wrap(err, "deps usecase: failed to write Brewfile")
	}

	d.printOutUC.PrintMdf(`
> [!NOTE]
> If you do not want to change it, do a process kill (ctrl + c)
`)

	d.resolveBrewDiffWithEditorCount = 0
	if err := d.resolveBrewDiffWithEditor(ctx); err != nil {
		return errors.Wrap(err, "deps usecase: failed to resolve Brewfile diff with editor")
	}

	return nil
}

func (d *DepsUsecaseImpl) resolveBrewDiffWithEditor(ctx context.Context) error {
	d.resolveBrewDiffWithEditorCount++

	if d.resolveBrewDiffWithEditorCount > resolveBrewDiffWithEditorMaxCount {
		d.printOutUC.PrintMdf(`
> [!CAUTION]
> Abort because brewfile was not updated
`)

		return errResolveBrewDiffWithEditorMaxCount
	}

	usr, err := user.Current()
	if err != nil {
		return errors.Wrap(err, "deps usecase: failed to get current user")
	}

	brewPath := usr.HomeDir + "/projects/dotfiles/data/Brewfile"

	if err := d.gitInfrastructure.GitDifftool(
		ctx,
		*d.printOutUC.GetOut(),
		*d.printOutUC.GetError(),
		brewPath,
	); err != nil {
		return errors.Wrap(err, "deps usecase: failed to open with code")
	}

	if data, err := d.fileInfrastructure.ReadFile(brewPath); err != nil {
		return errors.Wrap(err, "deps usecase: failed to read Brewfile")
	} else if strings.Contains(string(data), "# "+dofyBrewCategory) {
		d.printOutUC.PrintMdf(`
> [!CAUTION]
> Update your brewfile
`)

		err := d.resolveBrewDiffWithEditor(ctx)
		if err != nil {
			return errors.Wrap(err, "deps usecase: failed to resolve Brewfile diff with editor")
		}
	}

	return nil
}

func (d *DepsUsecaseImpl) mergeDiff(
	base []domain.BrewBundle,
	add []domain.BrewBundle,
	sub []domain.BrewBundle,
) []domain.BrewBundle {
	for _, diff := range add {
		diff.Categories = []string{dofyBrewCategory}
		base = append(base, diff)
	}

	for _, diff := range sub {
		for i, bundle := range base {
			if bundle.Name == diff.Name {
				base = append(base[:i], base[i+1:]...)
			}
		}
	}

	return base
}
