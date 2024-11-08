package usecase

import (
	"context"

	"github.com/pkg/errors"
	"github.com/shiron-dev/dotfiles/scripts/dofy/internal/infrastructure"
)

type BrewUsecase interface {
	InstallHomebrew(ctx context.Context) error
	InstallFormula(formula string) error

	InstallBrewBundle() error
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

func (b *BrewUsecaseImpl) InstallHomebrew(ctx context.Context) error {
	b.printOutUC.PrintMdf(`
### Installing Homebrew
`)

	err := b.brewInfrastructure.InstallHomebrew(ctx, *b.printOutUC.GetOut(), *b.printOutUC.GetError())
	if err != nil {
		return errors.Wrap(err, "brew usecase: failed to install Homebrew")
	}

	b.printOutUC.PrintMdf(`
### Set Homebrew environment
`)

	var brewPath string

	cfg, err := b.configUC.ScanEnvInfo()
	if err != nil {
		return errors.Wrap(err, "brew usecase: failed to get environment info")
	}

	switch cfg.os {
	case "darwin":
		brewPath = "/opt/homebrew/bin/brew"
	case "linux":
		brewPath = "/home/linuxbrew/.linuxbrew/bin/brew"
	}

	err = b.brewInfrastructure.SetHomebrewEnv(brewPath)
	if err != nil {
		return errors.Wrap(err, "brew usecase: failed to set Homebrew environment")
	}

	return nil
}

func (b *BrewUsecaseImpl) InstallFormula(formula string) error {
	b.printOutUC.PrintMdf(`
### Installing %s (with Homebrew)
`, formula)

	err := b.brewInfrastructure.InstallFormula(formula)
	if err != nil {
		return errors.Wrap(err, "brew usecase: failed to install formula")
	}

	return nil
}

func (b *BrewUsecaseImpl) InstallBrewBundle() error {
	err := b.brewInfrastructure.InstallBrewBundle(*b.printOutUC.GetOut(), *b.printOutUC.GetError())
	if err != nil {
		return errors.Wrap(err, "brew usecase: failed to install Brewfile")
	}

	return nil
}
