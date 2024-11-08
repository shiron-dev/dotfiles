package controller

import "dofy/internal/usecase"

type Controller interface {
	Start()
}

type ControllerImpl struct {
	printoutUC usecase.PrintOutUsecase
	configUC   usecase.ConfigUsecase
}

func NewController(printoutUC usecase.PrintOutUsecase, configUC usecase.ConfigUsecase) Controller {
	return &ControllerImpl{
		printoutUC: printoutUC,
		configUC:   configUC,
	}
}

func (c *ControllerImpl) Start() {
	logfile := c.printoutUC.SetLogOutput()
	defer logfile.Close()

	c.printoutUC.PrintMd(`

# shiron-dev dotfiles setup script

This script will install dependencies and setup dotfiles.

`)

	c.printoutUC.PrintMd(`
## Load environment information

### Environment information
`)
	envInfo, err := c.configUC.ScanEnvInfo()
	if err != nil {
		panic(err)
	}
	c.printoutUC.PrintObj(*envInfo)

	// 	infrastructure.PrintMd(`
	// ### Setup mode
	// `)

	// 	var mode string
	// 	if len(os.Args) > 1 {
	// 		mode = strings.ToLower(os.Args[1])
	// 		infrastructure.Println("The mode is set by command line arguments.")
	// 	} else {
	// 		infrastructure.Print("What mode do you use? [standard]: ")
	// 		scanner := bufio.NewScanner(os.Stdin)
	// 		if scanner.Scan() {
	// 			mode = strings.ToLower(strings.TrimSpace(scanner.Text()))
	// 			if mode == "" {
	// 				mode = "standard"
	// 			}
	// 		}
	// 	}

	// 	infrastructure.PrintMd("Start setup in `" + mode + "` mode.")

	// deps.InstallDeps()
}
