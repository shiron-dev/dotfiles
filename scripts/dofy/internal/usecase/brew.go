package usecase

// import (
// 	"io"
// 	"net/http"
// 	"os"
// 	"os/exec"
// 	"runtime"
// 	"strings"

// 	"github.com/shiron-dev/dotfiles/scripts/dofy/internal/infrastructure"
// )

// type BrewUsecase interface {
// 	InstallHomebrew()
// }

// type BrewUsecaseImpl struct {
// 	brewInfrastructure infrastructure.BrewInfrastructure
// 	printOutUC         PrintOutUsecase
// }

// func NewBrewUsecase(brewInfrastructure infrastructure.BrewInfrastructure, printOutUC PrintOutUsecase) BrewUsecase {
// 	return &BrewUsecaseImpl{
// 		brewInfrastructure: brewInfrastructure,
// 		printOutUC:         printOutUC,
// 	}
// }

// func (b *BrewUsecaseImpl) InstallHomebrew() {
// 	b.printOutUC.PrintMd(`
// ### Installing Homebrew
// `)

// 	url := "https://raw.githubusercontent.com/Homebrew/install/master/install.sh"

// 	resp, err := http.Get(url)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer resp.Body.Close()
// 	bytes, _ := io.ReadAll(resp.Body)

// 	cmd := exec.Command("/bin/bash", "-c", string(bytes))
// 	cmd.Stdout = *b.printOutUC.GetOut()
// 	cmd.Stderr = *b.printOutUC.GetError()
// 	err = cmd.Run()
// 	if err != nil {
// 		panic(err)
// 	}

// 	b.printOutUC.PrintMd(`
// ### Set Homebrew environment
// `)

// 	var brewPath string
// 	switch runtime.GOOS {
// 	case "darwin":
// 		brewPath = "/opt/homebrew/bin/brew"
// 	case "linux":
// 		brewPath = "/home/linuxbrew/.linuxbrew/bin/brew"
// 	}
// 	cmd = exec.Command("/bin/bash", "-c", `(echo; echo 'eval "$(/home/linuxbrew/.linuxbrew/bin/brew shellenv)"') >> ~/.bashrc`)
// 	cmd.Stdout = *b.printOutUC.GetOut()
// 	cmd.Stderr = *b.printOutUC.GetError()
// 	err = cmd.Run()
// 	if err != nil {
// 		panic(err)
// 	}

// 	cmd = exec.Command(brewPath, "shellenv")
// 	shellenv, _ := cmd.Output()
// 	for _, line := range strings.Split(string(shellenv), "\n") {
// 		if strings.HasPrefix(line, "export PATH=") {
// 			cmd = exec.Command("sh", "-c", "echo "+strings.Replace(line, "export PATH=", "", 1))
// 			out, _ := cmd.Output()
// 			os.Setenv("PATH", strings.Trim(string(out), "\""))
// 		}
// 	}
// }
