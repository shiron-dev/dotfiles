package usecase

import (
	"github.com/pkg/errors"
	"github.com/shiron-dev/dotfiles/scripts/dofy/internal/infrastructure"
)

type AnsibleUsecase interface {
	CheckPlaybook(invPath string, playbookPath string) error
	RunPlaybook(invPath string, playbookPath string) error
}

type AnsibleUsecaseImpl struct {
	ansibleInfrastructure infrastructure.AnsibleInfrastructure
	printOutUC            PrintOutUsecase
}

func NewAnsibleUsecase(
	ansibleInfrastructure infrastructure.AnsibleInfrastructure,
	printOutUC PrintOutUsecase,
) *AnsibleUsecaseImpl {
	return &AnsibleUsecaseImpl{
		ansibleInfrastructure: ansibleInfrastructure,
		printOutUC:            printOutUC,
	}
}

func (a *AnsibleUsecaseImpl) CheckPlaybook(invPath string, playbookPath string) error {
	a.printOutUC.PrintMdf(`
## Check Ansible playbook
`)

	err := a.ansibleInfrastructure.CheckPlaybook(
		invPath, playbookPath,
		*a.printOutUC.GetOut(),
		*a.printOutUC.GetError(),
	)
	if err != nil {
		return errors.Wrap(err, "ansible usecase: failed to check playbook")
	}

	return nil
}

func (a *AnsibleUsecaseImpl) RunPlaybook(invPath string, playbookPath string) error {
	a.printOutUC.PrintMdf(`
## Run Ansible playbook
`)

	err := a.ansibleInfrastructure.RunPlaybook(
		invPath, playbookPath,
		*a.printOutUC.GetOut(),
		*a.printOutUC.GetError(),
	)
	if err != nil {
		return errors.Wrap(err, "ansible usecase: failed to run playbook")
	}

	return nil
}
