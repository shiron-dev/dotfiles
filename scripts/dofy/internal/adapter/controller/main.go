package controller

import (
	"bufio"
	"context"
	"os"
	"strings"

	"github.com/shiron-dev/dotfiles/scripts/dofy/internal/usecase"
)

type DofyController interface {
	Start()
	getMode() string
}

type DofyControllerImpl struct {
	printoutUC usecase.PrintOutUsecase
	configUC   usecase.ConfigUsecase
	depsUC     usecase.DepsUsecase
}

func NewDofyController(
	printoutUC usecase.PrintOutUsecase,
	configUC usecase.ConfigUsecase,
	depsUC usecase.DepsUsecase,
) *DofyControllerImpl {
	return &DofyControllerImpl{
		printoutUC: printoutUC,
		configUC:   configUC,
		depsUC:     depsUC,
	}
}

func (c *DofyControllerImpl) Start() {
	logfile := c.printoutUC.SetLogOutput()
	defer logfile.Close()

	c.printoutUC.PrintMdf(`

# shiron-dev dotfiles setup script

This script will install dependencies and setup dotfiles.

`)

	c.printoutUC.PrintMdf(`
## Load environment information

### Environment information
`)

	if envInfo, err := c.configUC.ScanEnvInfo(); err == nil {
		c.printoutUC.PrintObj(*envInfo)
	} else {
		panic(err)
	}

	c.printoutUC.PrintMdf(`
### Setup mode
`)

	mode := c.getMode()

	c.printoutUC.PrintMdf("Start setup in `" + mode + "` mode.")

	initialSetupFlag := c.depsUC.CheckInstalled("brew")

	err := c.depsUC.InstallHomebrew(context.Background())
	if err != nil {
		panic(err)
	}

	err = c.depsUC.InstallGit()
	if err != nil {
		panic(err)
	}

	err = c.depsUC.CloneDotfiles()
	if err != nil {
		panic(err)
	}

	err = c.depsUC.InstallBrewBundle(initialSetupFlag)
	if err != nil {
		panic(err)
	}
}

func (c *DofyControllerImpl) getMode() string {
	var mode string
	if len(os.Args) > 1 {
		mode = strings.ToLower(os.Args[1])

		c.printoutUC.Println("The mode is set by command line arguments.")
	} else {
		c.printoutUC.Print("What mode do you use? [standard]: ")

		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			mode = strings.ToLower(strings.TrimSpace(scanner.Text()))
			if mode == "" {
				mode = "standard"
			}
		}
	}

	return mode
}
