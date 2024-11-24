package infrastructure

import (
	"io"
	"os/exec"

	"github.com/pkg/errors"
)

type AnsibleInfrastructure interface {
	SetWorkingDir(workingDir string)
	CheckPlaybook(invPath string, playbookPath string, sout io.Writer, serror io.Writer) error
	RunPlaybook(invPath string, playbookPath string, sout io.Writer, serror io.Writer) error
}

type AnsibleInfrastructureImpl struct {
	workingDir string
}

func NewAnsibleInfrastructure() *AnsibleInfrastructureImpl {
	return &AnsibleInfrastructureImpl{
		workingDir: "",
	}
}

func (a *AnsibleInfrastructureImpl) SetWorkingDir(workingDir string) {
	a.workingDir = workingDir
}

func (a *AnsibleInfrastructureImpl) CheckPlaybook(
	invPath string,
	playbookPath string,
	sout io.Writer,
	serror io.Writer,
) error {
	if a.workingDir == "" {
		return errors.New("ansible infrastructure: working directory is not set")
	}

	cmd := exec.Command("ansible-playbook", "-i", invPath, playbookPath, "-C")
	cmd.Stdout = sout
	cmd.Stderr = serror

	if err := cmd.Run(); err != nil {
		return errors.Wrap(err, "ansible infrastructure: failed to check playbook")
	}

	return nil
}

func (a *AnsibleInfrastructureImpl) RunPlaybook(
	invPath string,
	playbookPath string,
	sout io.Writer,
	serror io.Writer,
) error {
	if a.workingDir == "" {
		return errors.New("ansible infrastructure: working directory is not set")
	}

	cmd := exec.Command("ansible-playbook", "-i", invPath, playbookPath)
	cmd.Stdout = sout
	cmd.Stderr = serror

	if err := cmd.Run(); err != nil {
		return errors.Wrap(err, "ansible infrastructure: failed to run playbook")
	}

	return nil
}
