package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/shiron-dev/dotfiles/scripts/cmd/conf"
	"github.com/shiron-dev/dotfiles/scripts/cmd/deps"
	"github.com/shiron-dev/dotfiles/scripts/cmd/printout"
)

func main() {
	printout.PrintMd(`

# shiron-dev dotfiles setup script

This script will install dependencies and setup dotfiles.

`)

	printout.PrintMd(`
## Load environment information

### Environment information
`)
	envInfo := conf.ScanEnvInfo()
	printout.PrintObj(*envInfo)

	printout.PrintMd(`
### Setup mode
`)

	var mode string
	if len(os.Args) > 1 {
		mode = strings.ToLower(os.Args[1])
		fmt.Println("The mode is set by command line arguments.")
	} else {
		fmt.Print("What mode do you use? [standard]: ")
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			mode = strings.ToLower(strings.TrimSpace(scanner.Text()))
			if mode == "" {
				mode = "standard"
			}
		}
	}

	printout.PrintMd("Start setup in `" + mode + "` mode.")

	deps.InstallDeps()
}
