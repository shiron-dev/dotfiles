package usecase

import "github.com/shiron-dev/dotfiles/scripts/dofy/internal/infrastructure"

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

type configError struct {
	err error
}

func (e *configError) Error() string {
	return "ConfigUC: " + e.err.Error()
}

func (c *ConfigUsecaseImpl) ScanEnvInfo() (*EnvInfo, error) {
	gos, err := c.configInfrastructure.GetOS()
	if err != nil {
		return nil, &configError{err}
	}

	osVersion, err := c.configInfrastructure.GetOSVersion()
	if err != nil {
		return nil, &configError{err}
	}

	arch, err := c.configInfrastructure.GetArch()
	if err != nil {
		return nil, &configError{err}
	}

	return &EnvInfo{
		os:        gos,
		osVersion: osVersion,
		arch:      arch,
	}, nil
}
