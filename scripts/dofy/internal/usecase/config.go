package usecase

import (
	"github.com/pkg/errors"
	"github.com/shiron-dev/dotfiles/scripts/dofy/internal/infrastructure"
)

type ConfigUsecase interface {
	ScanEnvInfo() (*EnvInfo, error)
}

type ConfigUsecaseImpl struct {
	configInfrastructure infrastructure.ConfigInfrastructure
}

func NewConfigUsecase(configInfrastructure infrastructure.ConfigInfrastructure) *ConfigUsecaseImpl {
	return &ConfigUsecaseImpl{
		configInfrastructure: configInfrastructure,
	}
}

type EnvInfo struct {
	os        string
	osVersion string
	arch      string
}

func (c *ConfigUsecaseImpl) ScanEnvInfo() (*EnvInfo, error) {
	gos, err := c.configInfrastructure.GetOS()
	if err != nil {
		return nil, errors.Wrap(err, "config usecase: failed to get os")
	}

	osVersion, err := c.configInfrastructure.GetOSVersion()
	if err != nil {
		return nil, errors.Wrap(err, "config usecase: failed to get os version")
	}

	arch, err := c.configInfrastructure.GetArch()
	if err != nil {
		return nil, errors.Wrap(err, "config usecase: failed to get arch")
	}

	return &EnvInfo{
		os:        gos,
		osVersion: osVersion,
		arch:      arch,
	}, nil
}
