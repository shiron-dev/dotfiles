package main

import (
	"bufio"
	"dofy/internal/infrastructure"
	"dofy/internal/infrastructure/conf"
	"dofy/internal/infrastructure/deps"
	"os"
	"strings"
)

func main() {
	logfile := infrastructure.SetLogOutput()
	defer logfile.Close()

	infrastructure.PrintMd(`

# shiron-dev dotfiles setup script

This script will install dependencies and setup dotfiles.

`)

	infrastructure.PrintMd(`
## Load environment information

### Environment information
`)
	envInfo := conf.ScanEnvInfo()
	infrastructure.PrintObj(*envInfo)

	infrastructure.PrintMd(`
### Setup mode
`)

	var mode string
	if len(os.Args) > 1 {
		mode = strings.ToLower(os.Args[1])
		infrastructure.Println("The mode is set by command line arguments.")
	} else {
		infrastructure.Print("What mode do you use? [standard]: ")
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			mode = strings.ToLower(strings.TrimSpace(scanner.Text()))
			if mode == "" {
				mode = "standard"
			}
		}
	}

	infrastructure.PrintMd("Start setup in `" + mode + "` mode.")

	deps.InstallDeps()
}
