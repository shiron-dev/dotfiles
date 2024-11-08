package usecase

import (
	"context"

	"github.com/shiron-dev/dotfiles/scripts/dofy/internal/infrastructure"
)

type BrewUsecase interface {
	InstallHomebrew(ctx context.Context) error
	InstallFormula(formula string) error
}

type BrewUsecaseImpl struct {
	brewInfrastructure infrastructure.BrewInfrastructure
	printOutUC         PrintOutUsecase
	configUC           ConfigUsecase
}

func NewBrewUsecase(
	brewInfrastructure infrastructure.BrewInfrastructure,
	printOutUC PrintOutUsecase,
	configUC ConfigUsecase,
) *BrewUsecaseImpl {
	return &BrewUsecaseImpl{
		brewInfrastructure: brewInfrastructure,
		printOutUC:         printOutUC,
		configUC:           configUC,
	}
}

type BrewError struct {
	err error
}

func (e *BrewError) Error() string {
	return "BrewUC: " + e.err.Error()
}

func (b *BrewUsecaseImpl) InstallHomebrew(ctx context.Context) error {
	b.printOutUC.PrintMdf(`
### Installing Homebrew
`)

	err := b.brewInfrastructure.InstallHomebrew(ctx, *b.printOutUC.GetOut(), *b.printOutUC.GetError())
	if err != nil {
		return &BrewError{err}
	}

	b.printOutUC.PrintMdf(`
### Set Homebrew environment
`)

	var brewPath string

	cfg, err := b.configUC.ScanEnvInfo()
	if err != nil {
		return &BrewError{err}
	}

	switch cfg.os {
	case "darwin":
		brewPath = "/opt/homebrew/bin/brew"
	case "linux":
		brewPath = "/home/linuxbrew/.linuxbrew/bin/brew"
	}

	err = b.brewInfrastructure.SetHomebrewEnv(brewPath)
	if err != nil {
		return &BrewError{err}
	}

	return nil
}

func (b *BrewUsecaseImpl) InstallFormula(formula string) error {
	b.printOutUC.PrintMdf(`
### Installing %s (with Homebrew)
`, formula)

	err := b.brewInfrastructure.InstallFormula(formula)
	if err != nil {
		return &BrewError{err}
	}

	return nil
}
