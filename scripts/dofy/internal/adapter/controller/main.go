package controller

import (
	"bufio"
	"context"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/shiron-dev/dotfiles/scripts/dofy/internal/usecase"
)

type DofyController interface {
	Start()
	getMode() string
	getYN(msg string, def bool) bool
}

type DofyControllerImpl struct {
	ansibleUC  usecase.AnsibleUsecase
	printoutUC usecase.PrintOutUsecase
	configUC   usecase.ConfigUsecase
	depsUC     usecase.DepsUsecase
	vsCodeUC   usecase.VSCodeUsecase
}

func NewDofyController(
	ansibleUC usecase.AnsibleUsecase,
	printoutUC usecase.PrintOutUsecase,
	configUC usecase.ConfigUsecase,
	depsUC usecase.DepsUsecase,
	vsCodeUC usecase.VSCodeUsecase,
) *DofyControllerImpl {
	return &DofyControllerImpl{
		ansibleUC:  ansibleUC,
		printoutUC: printoutUC,
		configUC:   configUC,
		depsUC:     depsUC,
		vsCodeUC:   vsCodeUC,
	}
}

//nolint:funlen,cyclop
func (c *DofyControllerImpl) Start() {
	logfile := c.printoutUC.SetLogOutput()
	defer func() {
		if err := logfile.Close(); err != nil {
			log.Fatal(err)
		}
	}()

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

	initialSetupFlag := !c.depsUC.CheckInstalled("brew")

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

	dotPath, err := c.configUC.GetDotfilesDir()
	if err != nil {
		panic(err)
	}

	err = c.vsCodeUC.SaveExtensions()
	if err != nil {
		panic(err)
	}

	c.ansibleUC.SetWorkingDir(filepath.Join(dotPath, "scripts/ansible"))

	c.printoutUC.PrintMdf(`
## Run Ansible playbook

`)

	ok := c.getYN("Do you want to run Ansible?", true)
	if ok {
		err = c.ansibleUC.RunPlaybook("hosts.yml", "site.yml")
		if err != nil {
			panic(err)
		}
	} else {
		err = c.ansibleUC.CheckPlaybook("hosts.yml", "site.yml")
		if err != nil {
			panic(err)
		}
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

func (c *DofyControllerImpl) getYN(msg string, def bool) bool {
	var ynStr string
	if def {
		ynStr = "Y/n"
	} else {
		ynStr = "y/N"
	}

	c.printoutUC.Print(msg + " [" + ynStr + "]: ")

	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		ynStr = strings.ToLower(strings.TrimSpace(scanner.Text()))
	}

	if ynStr == "" {
		return def
	}

	return ynStr == "y"
}
