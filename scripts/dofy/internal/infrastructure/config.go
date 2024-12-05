package infrastructure

import (
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/pkg/errors"
)

type ConfigInfrastructure interface {
	GetOS() (string, error)
	GetOSVersion() (string, error)
	GetArch() (string, error)
	GetDotfilesDir() (string, error)
}

type ConfigInfrastructureImpl struct{}

func NewConfigInfrastructure() *ConfigInfrastructureImpl {
	return &ConfigInfrastructureImpl{}
}

func (c *ConfigInfrastructureImpl) GetOS() (string, error) {
	return runtime.GOOS, nil
}

func (c *ConfigInfrastructureImpl) GetOSVersion() (string, error) {
	osVersion, err := exec.Command("uname", "-r").Output()
	if err != nil {
		return "", errors.Wrap(err, "config infrastructure: failed to get os version")
	}

	return strings.TrimSpace(string(osVersion)), nil
}

func (c *ConfigInfrastructureImpl) GetArch() (string, error) {
	arch, err := exec.Command("uname", "-p").Output()
	if err != nil {
		return "", errors.Wrap(err, "config infrastructure: failed to get arch")
	}

	return strings.TrimSpace(string(arch)), nil
}

func (c *ConfigInfrastructureImpl) GetDotfilesDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", errors.Wrap(err, "deps usecase: failed to get current user")
	}

	return filepath.Join(usr.HomeDir, "/projects/github.com/shiron-dev/dotfiles/"), nil
}
