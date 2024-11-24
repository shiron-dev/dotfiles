package infrastructure

import (
	"io"
	"os/exec"
)

type AnsibleInfrastructure interface {
	CheckPlaybook(invPath string, playbookPath string, sout io.Writer, serror io.Writer) error
	RunPlaybook(invPath string, playbookPath string, sout io.Writer, serror io.Writer) error
}

type AnsibleInfrastructureImpl struct{}

func NewAnsibleInfrastructure() *AnsibleInfrastructureImpl {
	return &AnsibleInfrastructureImpl{}
}

func (a *AnsibleInfrastructureImpl) CheckPlaybook(
	invPath string,
	playbookPath string,
	sout io.Writer,
	serror io.Writer,
) error {
	cmd := exec.Command("ansible-playbook", "-i", invPath, playbookPath, "-C")
	cmd.Stdout = sout
	cmd.Stderr = serror

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func (a *AnsibleInfrastructureImpl) RunPlaybook(
	invPath string,
	playbookPath string,
	sout io.Writer,
	serror io.Writer,
) error {
	cmd := exec.Command("ansible-playbook", "-i", invPath, playbookPath)
	cmd.Stdout = sout
	cmd.Stderr = serror

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
