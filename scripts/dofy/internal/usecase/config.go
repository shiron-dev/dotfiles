package usecase

import "dofy/internal/infrastructure"

type ConfigUsecase interface {
	ScanEnvInfo() (*EnvInfo, error)
}

type ConfigUsecaseImpl struct {
	configInfrastructure infrastructure.ConfigInfrastructure
}

func NewConfigUsecase(configInfrastructure infrastructure.ConfigInfrastructure) ConfigUsecase {
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
	os, err := c.configInfrastructure.GetOS()
	if err != nil {
		return nil, err
	}

	osVersion, err := c.configInfrastructure.GetOSVersion()
	if err != nil {
		return nil, err
	}

	arch, err := c.configInfrastructure.GetArch()
	if err != nil {
		return nil, err
	}

	return &EnvInfo{
		os:        os,
		osVersion: osVersion,
		arch:      arch,
	}, nil
}
