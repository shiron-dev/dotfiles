package usecase

import (
	"context"

	"github.com/shiron-dev/dotfiles/scripts/dofy/internal/infrastructure"
)

type DepsUsecase interface {
	InstallDeps(ctx context.Context) error
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

	return nil
}

// func InstallDeps() {

// 	infrastructure.PrintMd(`
// ## Installing required packages with Homebrew

// - git
// `)

// 	infrastructure.PrintMd("### Installing git")
// 	if checkInstalled("git") {
// 		infrastructure.Println("git is already installed")
// 	} else {
// 		infrastructure.Println("Installing git")
// 		installWithBrew("git")
// 	}

// 	infrastructure.PrintMd(`
// ## Git clone dotfiles repository

// https://github.com/shiron-dev/dotfiles.git
// `)

// 	usr, _ := user.Current()
// 	if _, err := os.Stat(usr.HomeDir + "/projects/dotfiles"); err == nil {
// 		infrastructure.Println("dotfiles directory already exists")
// 	} else {
// 		infrastructure.Println("Cloning dotfiles repository")
// 		cmd := exec.Command("git", "clone", "https://github.com/shiron-dev/dotfiles.git", usr.HomeDir+"/projects/dotfiles")
// 		cmd.Run()
// 	}

// 	infrastructure.PrintMd(`
// ## Installing brew packages

// Install the packages using Homebrew Bundle.
// `)

// 	dumpTmpBrewBundle()
// 	diffBundle, diffTmpBundles := checkDiffBrewBundle(
// 		usr.HomeDir+"/projects/dotfiles/data/Brewfile",
// 		usr.HomeDir+"/projects/dotfiles/data/Brewfile.tmp",
// 	)
// 	if len(diffBundle) > 0 {
// 		infrastructure.Println("Installing brew packages")
// 		installBrewBundle()
// 	} else {
// 		infrastructure.Println("No new packages to install")
// 	}
// 	if len(diffTmpBundles) > 0 {
// 		var diffNames string
// 		for _, diff := range diffTmpBundles {
// 			diffNames += "- " + diff.name + "\n"
// 		}
// 		infrastructure.Println(color.RedString("The dotfiles Brewfile and the currently installed package are different."))
// 		infrastructure.PrintMd(`
// ### Update Brewfile

// diff:
// ` + diffNames + `

// What will you do to resolve the diff?

// 1. run ` + "`brew bundlecleanup`" + `
// 2. update the Brewfile with the currently installed packages
// 3. do nothing
// 4. exit
// `)
// 		infrastructure.Print("What do you run? [1-4]: ")
// 		scanner := bufio.NewScanner(os.Stdin)
// 		if scanner.Scan() {
// 			switch strings.TrimSpace(scanner.Text()) {
// 			case "1":
// 				infrastructure.Println("Running `brew bundle cleanup`")
// 				cleanupBrewBundle(true)
// 			case "2":
// 				infrastructure.Println("Open Brewfile with code")
// 				utils.OpenWithCode(
// 					usr.HomeDir+"/projects/dotfiles/data/Brewfile",
// 					usr.HomeDir+"/projects/dotfiles/data/Brewfile.tmp",
// 				)
// 			case "3":
// 				infrastructure.Println("Do nothing")
// 			default:
// 				infrastructure.Println("Exit")
// 				os.Exit(0)
// 			}
// 		}
// 	}

// 	infrastructure.PrintMd(`
// ### Install brew packages with Brewfile
// 	`)
// 	installBrewBundle()
// }
