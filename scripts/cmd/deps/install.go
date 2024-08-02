package deps

import (
	"os"
	"os/exec"
	"os/user"

	"github.com/shiron-dev/dotfiles/scripts/cmd/printout"
)

func installHomebrew() {
	cmd := exec.Command("/bin/bash", "-c", "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install.sh)")
	out, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	printout.Println(string(out))
}

func installWithBrew(pkg string) {
	cmd := exec.Command("brew", "install", pkg)
	cmd.Run()
}

func installBrewBundle() {
	usr, _ := user.Current()
	cmd := exec.Command("brew", "bundle", "--file", usr.HomeDir+"/projects/dotfiles/data/Brewfile")
	cmd.Stdout = os.Stdout
	cmd.Run()
}
