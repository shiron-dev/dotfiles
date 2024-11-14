package infrastructure

import (
	"os/exec"

	"github.com/pkg/errors"
)

type DepsInfrastructure interface {
	CheckInstalled(name string) bool
	OpenWithCode(path ...string) error
}

type DepsInfrastructureImpl struct{}

func NewDepsInfrastructure() *DepsInfrastructureImpl {
	return &DepsInfrastructureImpl{}
}

func (d *DepsInfrastructureImpl) CheckInstalled(name string) bool {
	_, err := exec.LookPath(name)

	return err == nil
}

func (d *DepsInfrastructureImpl) OpenWithCode(path ...string) error {
	args := []string{"-n", "-w"}
	args = append(args, path...)

	if err := exec.Command("code", args...).Run(); err != nil {
		return errors.Wrap(err, "deps infrastructure: failed to open with code")
	}

	return nil
}
