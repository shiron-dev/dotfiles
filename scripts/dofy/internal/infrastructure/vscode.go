package infrastructure

import (
	"os/exec"

	"github.com/pkg/errors"
)

type VSCodeInfrastructure interface {
	ListExtensions() ([]string, error)
}

type VSCodeInfrastructureImpl struct{}

func NewVSCodeInfrastructure() *VSCodeInfrastructureImpl {
	return &VSCodeInfrastructureImpl{}
}

func (v *VSCodeInfrastructureImpl) ListExtensions() ([]string, error) {
	cmd := exec.Command("code", "--list-extensions")
	if out, err := cmd.Output(); err != nil {
		return nil, errors.Wrap(err, "vscode infrastructure: failed to list extensions")
	} else {
		return []string{string(out)}, nil
	}
}
