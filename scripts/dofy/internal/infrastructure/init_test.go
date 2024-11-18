package infrastructure_test

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"runtime"
	"testing"

	"github.com/pkg/errors"
	"github.com/shiron-dev/dotfiles/scripts/dofy/internal/di"
	"github.com/shiron-dev/dotfiles/scripts/dofy/internal/infrastructure"
)

func TestMain(m *testing.M) {
	infra, err := di.InitializeTestInfrastructureSet(os.Stdout, os.Stderr)
	if err != nil {
		panic(err)
	}

	err = setupBrew(infra.BrewInfrastructure)
	if err != nil {
		panic(err)
	}

	code := m.Run()
	os.Exit(code)
}

func setupBrew(brew infrastructure.BrewInfrastructure) error {
	_, err := exec.LookPath("brew")
	if err == nil {
		return nil
	}

	outBuffer := &bytes.Buffer{}

	errBuffer := &bytes.Buffer{}

	err = brew.InstallHomebrew(context.Background(), outBuffer, errBuffer)
	if err != nil {
		return errors.Wrap(err, "failed to install homebrew\n"+outBuffer.String()+"\n"+errBuffer.String())
	}

	brewPath := ""

	switch runtime.GOOS {
	case "darwin":
		brewPath = "/opt/homebrew/bin/brew"
	case "linux":
		brewPath = "/home/linuxbrew/.linuxbrew/bin/brew"
	}

	err = brew.SetHomebrewEnv(brewPath)
	if err != nil {
		return errors.Wrap(err, "failed to set homebrew env")
	}

	return nil
}
