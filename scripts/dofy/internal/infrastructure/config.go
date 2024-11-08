package infrastructure

import (
	"os/exec"
	"runtime"
	"strings"
)

type ConfigInfrastructure interface {
	GetOS() (string, error)
	GetOSVersion() (string, error)
	GetArch() (string, error)
}

type ConfigInfrastructureImpl struct {
}

func NewConfigInfrastructure() ConfigInfrastructure {
	return &ConfigInfrastructureImpl{}
}

func (c *ConfigInfrastructureImpl) GetOS() (string, error) {
	return runtime.GOOS, nil
}

func (c *ConfigInfrastructureImpl) GetOSVersion() (string, error) {
	osVersion, err := exec.Command("uname", "-r").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(osVersion)), nil
}

func (c *ConfigInfrastructureImpl) GetArch() (string, error) {
	arch, err := exec.Command("uname", "-p").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(arch)), nil
}
