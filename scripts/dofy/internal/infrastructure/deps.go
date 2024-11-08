package infrastructure

import "os/exec"

type DepsInfrastructure interface {
	CheckInstalled(name string) bool
}

type DepsInfrastructureImpl struct{}

func NewDepsInfrastructure() *DepsInfrastructureImpl {
	return &DepsInfrastructureImpl{}
}

func (d *DepsInfrastructureImpl) CheckInstalled(name string) bool {
	_, err := exec.LookPath(name)

	return err == nil
}
