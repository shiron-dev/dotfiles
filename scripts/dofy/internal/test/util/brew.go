package util

import (
	"bytes"
	"context"
	"os/exec"
	"runtime"

	"github.com/pkg/errors"
	"github.com/shiron-dev/dotfiles/scripts/dofy/internal/infrastructure"
)

func SetupBrew(brew infrastructure.BrewInfrastructure) error {
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

	err = brew.SetHomebrewEnv(runtime.GOOS)
	if err != nil {
		return errors.Wrap(err, "failed to set homebrew env")
	}

	outBuffer.Reset()

	errBuffer.Reset()

	err = brew.InstallTap("Homebrew/bundle", outBuffer, errBuffer)
	if err != nil {
		return errors.Wrap(err, "failed to install tap\n"+outBuffer.String()+"\n"+errBuffer.String())
	}

	return nil
}
