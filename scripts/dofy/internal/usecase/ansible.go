package usecase

import (
	"github.com/shiron-dev/dotfiles/scripts/dofy/internal/infrastructure"
)

type AnsibleUsecase interface {
	CheckPlaybook(invPath string, playbookPath string) error
	RunPlaybook(invPath string, playbookPath string) error
}

type AnsibleUsecaseImpl struct {
	ansibleInfrastructure  infrastructure.AnsibleInfrastructure
	PrintOutInfrastructure infrastructure.PrintOutInfrastructure
}

func NewAnsibleUsecase(
	ansibleInfrastructure infrastructure.AnsibleInfrastructure,
	PrintOutInfrastructure infrastructure.PrintOutInfrastructure,
) *AnsibleUsecaseImpl {
	return &AnsibleUsecaseImpl{}
}

func (a *AnsibleUsecaseImpl) CheckPlaybook(invPath string, playbookPath string) error {
	return a.ansibleInfrastructure.CheckPlaybook(
		invPath, playbookPath,
		*a.PrintOutInfrastructure.GetOut(),
		*a.PrintOutInfrastructure.GetError(),
	)
}

func (a *AnsibleUsecaseImpl) RunPlaybook(invPath string, playbookPath string) error {
	return a.ansibleInfrastructure.RunPlaybook(
		invPath, playbookPath,
		*a.PrintOutInfrastructure.GetOut(),
		*a.PrintOutInfrastructure.GetError(),
	)
}
