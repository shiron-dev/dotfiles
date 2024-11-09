package infrastructure

import (
	"context"
	"io"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"strings"

	"github.com/pkg/errors"
)

type BrewInfrastructure interface {
	InstallHomebrew(ctx context.Context, sout io.Writer, serror io.Writer) error
	SetHomebrewEnv(brewPath string) error
	InstallFormula(pkg string) error
	DumpTmpBrewBundle(sout io.Writer, serror io.Writer) error
	InstallBrewBundle(sout io.Writer, serror io.Writer) error
	CleanupBrewBundle(isForce bool, sout io.Writer, serror io.Writer) error
}

type BrewInfrastructureImpl struct{}

func NewBrewInfrastructure() *BrewInfrastructureImpl {
	return &BrewInfrastructureImpl{}
}

func (b *BrewInfrastructureImpl) InstallHomebrew(ctx context.Context, sout io.Writer, serror io.Writer) error {
	url := "https://raw.githubusercontent.com/Homebrew/install/master/install.sh"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return errors.Wrap(err, "brew infrastructure: failed to create request")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "brew infrastructure: failed to send request")
	}
	defer resp.Body.Close()

	bytes, _ := io.ReadAll(resp.Body)

	//nolint:gosec
	cmd := exec.Command("/bin/bash", "-c", string(bytes))
	cmd.Stdout = sout
	cmd.Stderr = serror

	if err = cmd.Run(); err != nil {
		return errors.Wrap(err, "brew infrastructure: failed to run command")
	}

	return nil
}

func (b *BrewInfrastructureImpl) SetHomebrewEnv(brewPath string) error {
	cmd := exec.Command(brewPath, "shellenv")

	shellenv, err := cmd.Output()
	if err != nil {
		return errors.Wrap(err, "brew infrastructure: failed to get shellenv")
	}

	for _, line := range strings.Split(string(shellenv), "\n") {
		if strings.HasPrefix(line, "export PATH=") {
			//nolint:gosec
			cmd = exec.Command("sh", "-c", "echo "+strings.Replace(line, "export PATH=", "", 1))
			out, _ := cmd.Output()
			os.Setenv("PATH", strings.Trim(string(out), "\""))
		}
	}

	return nil
}

func (b *BrewInfrastructureImpl) InstallFormula(formula string) error {
	cmd := exec.Command("brew", "install", formula)
	if err := cmd.Run(); err != nil {
		return errors.Wrap(err, "brew infrastructure: failed to run brew install command")
	}

	return nil
}

func (b *BrewInfrastructureImpl) DumpTmpBrewBundle(sout io.Writer, serror io.Writer) error {
	usr, _ := user.Current()
	path := usr.HomeDir + "/projects/dotfiles/data/Brewfile.tmp"

	if _, err := os.Stat(path); err == nil {
		os.Remove(path)
	}

	cmd := exec.Command("brew", "bundle", "dump", "--tap", "--formula", "--cask", "--mas", "--file", path)
	cmd.Stdout = sout
	cmd.Stderr = serror

	if err := cmd.Run(); err != nil {
		return errors.Wrap(err, "brew infrastructure: failed to run brew bundle dump command")
	}

	return nil
}

func (b *BrewInfrastructureImpl) InstallBrewBundle(sout io.Writer, serror io.Writer) error {
	usr, _ := user.Current()
	//nolint:gosec
	cmd := exec.Command("brew", "bundle", "--no-lock", "--file", usr.HomeDir+"/projects/dotfiles/data/Brewfile")
	cmd.Stdout = sout
	cmd.Stderr = serror

	if err := cmd.Run(); err != nil {
		return errors.Wrap(err, "brew infrastructure: failed to run brew bundle command")
	}

	return nil
}

// func checkBrewBundle() {
// 	usr, _ := user.Current()
// 	cmd := exec.Command("brew", "bundle", "check", "--file", usr.HomeDir+"/projects/dotfiles/data/Brewfile")
// 	cmd.Stdout = infrastructure.Out
// 	cmd.Stderr = infrastructure.Error
// 	err := cmd.Run()
// 	if err != nil {
// 		panic(err)
// 	}
// }

func (b *BrewInfrastructureImpl) CleanupBrewBundle(isForce bool, sout io.Writer, serror io.Writer) error {
	usr, _ := user.Current()
	forceFlag := ""

	if isForce {
		forceFlag = "--force"
	}

	//nolint:gosec
	cmd := exec.Command("brew", "bundle", "cleanup", forceFlag, "--file", usr.HomeDir+"/projects/dotfiles/data/Brewfile")
	cmd.Stdout = sout
	cmd.Stderr = serror

	if err := cmd.Run(); err != nil {
		return errors.Wrap(err, "brew infrastructure: failed to run brew bundle cleanup command")
	}

	return nil
}
