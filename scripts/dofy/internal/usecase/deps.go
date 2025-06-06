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

//nolint:interfacebloat
type DepsUsecase interface {
	CheckInstalled(name string) bool
	InstallHomebrew(ctx context.Context) error
	InstallGit() error
	CloneDotfiles() error
	InstallBrewBundle(forceInstall bool) error
	Finish() error

	showBrewDiff(
		diffBundles []domain.BrewBundle,
		diffTmpBundles []domain.BrewBundle,
		brewPath string, brewTmpPath string,
	) error
	updateBrewfile(brewPath string, brewTmpPath string) error
	resolveBrewDiff(brewPath string, brewTmpPath string) error
	resolveBrewDiffWithEditor(ctx context.Context, brewPath string) error
	fmtBrewfile(brewPath string) error
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

func (d *DepsUsecaseImpl) CheckInstalled(name string) bool {
	return d.depsInfrastructure.CheckInstalled(name)
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
		err := d.brewUC.InstallFormula("git", domain.BrewBundleTypeFormula)
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

	if _, err := os.Stat(usr.HomeDir + "/projects/github.com/shiron-dev/dotfiles"); err == nil {
		d.printOutUC.Println("dotfiles directory already exists")
	} else {
		d.printOutUC.Println("Cloning dotfiles repository")

		//nolint:gosec
		cmd := exec.Command(
			"git",
			"clone",
			"https://github.com/shiron-dev/dotfiles.git",
			filepath.Join(usr.HomeDir, "/projects/github.com/shiron-dev/dotfiles"),
		)
		if err := cmd.Run(); err != nil {
			return errors.Wrap(err, "deps usecase: failed to clone dotfiles repository")
		}
	}

	return nil
}

func (d *DepsUsecaseImpl) InstallBrewBundle(forceInstall bool) error {
	usr, err := user.Current()
	if err != nil {
		return errors.Wrap(err, "deps usecase: failed to get current user")
	}

	brewPath := usr.HomeDir + "/projects/github.com/shiron-dev/dotfiles/data/brew/Brewfile"
	brewTmpPath := usr.HomeDir + "/projects/github.com/shiron-dev/dotfiles/data/brew/Brewfile.tmp"

	d.printOutUC.PrintMdf(`
## Installing brew packages

Install the packages using Homebrew Bundle.
`)

// 	d.printOutUC.PrintMdf(`
// ### Install brew bundle

// ` + "`brew tap Homebrew/bundle`\n")

// 	if err := d.brewUC.InstallFormula("Homebrew/bundle", domain.BrewBundleTypeTap); err != nil {
// 		return errors.Wrap(err, "deps usecase: failed to install Homebrew/bundle")
// 	}

	err = d.brewUC.DumpTmpBrewBundle(brewTmpPath)
	if err != nil {
		return errors.Wrap(err, "deps usecase: failed to dump tmp Brewfile")
	}

	diffBundles, diffTmpBundles, err := d.brewUC.CheckDiffBrewBundle(brewPath, brewTmpPath)
	if err != nil {
		return errors.Wrap(err, "deps usecase: failed to check diff Brewfile")
	}

	if !forceInstall && len(diffTmpBundles)+len(diffTmpBundles) > 0 {
		err := d.showBrewDiff(diffBundles, diffTmpBundles, brewPath, brewTmpPath)
		if err != nil {
			return errors.Wrap(err, "deps usecase: failed to update Brewfile")
		}
	}

	d.printOutUC.PrintMdf(`
### Format Brewfile
`)

	err = d.fmtBrewfile(brewPath)
	if err != nil {
		return errors.Wrap(err, "deps usecase: failed to format Brewfile")
	}

	d.printOutUC.PrintMdf(`
### Install brew packages with Brewfile
`)

	if err := d.brewUC.InstallBrewBundle(brewPath); err != nil {
		return errors.Wrap(err, "deps usecase: failed to install brew packages")
	}

	return nil
}

func (d *DepsUsecaseImpl) Finish() error {
	d.printOutUC.PrintMdf(`
## Finish
`)

	return nil
}

func (d *DepsUsecaseImpl) showBrewDiff(
	diffBundles []domain.BrewBundle,
	diffTmpBundles []domain.BrewBundle,
	brewPath string, brewTmpPath string,
) error {
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

	return d.updateBrewfile(brewPath, brewTmpPath)
}

func (d *DepsUsecaseImpl) updateBrewfile(brewPath string, brewTmpPath string) error {
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		switch strings.TrimSpace(scanner.Text()) {
		case "1":
			d.printOutUC.PrintMdf("#### Open Brewfile with code\n")

			if err := d.resolveBrewDiff(brewPath, brewTmpPath); err != nil {
				return errors.Wrap(err, "deps usecase: failed to resolve Brewfile diff")
			}

			d.printOutUC.Println("Running `brew bundle cleanup`")

			if err := d.brewUC.CleanupBrewBundle(brewPath, true); err != nil {
				return errors.Wrap(err, "deps usecase: failed to run brew bundle cleanup")
			}

		case "2":
			d.printOutUC.Println("Running `brew bundle cleanup`")

			if err := d.brewUC.CleanupBrewBundle(brewPath, true); err != nil {
				return errors.Wrap(err, "deps usecase: failed to run brew bundle cleanup")
			}
		case "3":
			d.printOutUC.Println("Do nothing")
		default:
			d.printOutUC.Println("Exit")

			if err := d.Finish(); err != nil {
				return errors.Wrap(err, "deps usecase: failed to finish")
			}
		}
	}

	return nil
}

func (d *DepsUsecaseImpl) resolveBrewDiff(brewPath string, brewTmpPath string) error {
	diffBundles, diffTmpBundles, err := d.brewUC.CheckDiffBrewBundle(brewPath, brewTmpPath)
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

		d.gitInfrastructure.SetGitDir(filepath.Dir(brewPath))

		if err := d.gitInfrastructure.CheckoutFile(brewPath); err != nil {
			panic(err)
		}
	}()

	bundles, err := d.brewInfrastructure.ReadBrewBundle(brewPath)
	if err != nil {
		return errors.Wrap(err, "deps usecase: failed to read Brewfile")
	}

	if err = d.brewInfrastructure.WriteBrewBundle(brewPath, mergeDiff(bundles, diffTmpBundles, diffBundles)); err != nil {
		return errors.Wrap(err, "deps usecase: failed to write Brewfile")
	}

	d.printOutUC.PrintMdf(`
> [!NOTE]
> If you do not want to change it, do a process kill (ctrl + c)
`)

	d.resolveBrewDiffWithEditorCount = 0
	if err := d.resolveBrewDiffWithEditor(ctx, brewPath); err != nil {
		return errors.Wrap(err, "deps usecase: failed to resolve Brewfile diff with editor")
	}

	return nil
}

func (d *DepsUsecaseImpl) resolveBrewDiffWithEditor(ctx context.Context, brewPath string) error {
	d.resolveBrewDiffWithEditorCount++

	if d.resolveBrewDiffWithEditorCount > resolveBrewDiffWithEditorMaxCount {
		d.printOutUC.PrintMdf(`
> [!CAUTION]
> Abort because brewfile was not updated
`)

		return errResolveBrewDiffWithEditorMaxCount
	}

	d.gitInfrastructure.SetGitDir(filepath.Dir(brewPath))

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

		err := d.resolveBrewDiffWithEditor(ctx, brewPath)
		if err != nil {
			return errors.Wrap(err, "deps usecase: failed to resolve Brewfile diff with editor")
		}
	}

	return nil
}

func (d *DepsUsecaseImpl) fmtBrewfile(brewPath string) error {
	bundles, err := d.brewInfrastructure.ReadBrewBundle(brewPath)
	if err != nil {
		return errors.Wrap(err, "deps usecase: failed to read Brewfile")
	}

	if err := d.brewInfrastructure.WriteBrewBundle(brewPath, bundles); err != nil {
		return errors.Wrap(err, "deps usecase: failed to write Brewfile")
	}

	return nil
}

func mergeDiff(
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
