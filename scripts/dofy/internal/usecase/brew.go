package usecase

import (
	"context"

	"github.com/pkg/errors"
	"github.com/shiron-dev/dotfiles/scripts/dofy/internal/domain"
	"github.com/shiron-dev/dotfiles/scripts/dofy/internal/infrastructure"
)

type BrewUsecase interface {
	InstallHomebrew(ctx context.Context) error
	InstallFormula(formula string) error
	InstallBrewBundle() error
	DumpTmpBrewBundle() error
	CheckDiffBrewBundle(bundlePath string, tmpPath string) ([]domain.BrewBundle, []domain.BrewBundle, error)
	CleanupBrewBundle(isForce bool) error
}

type BrewUsecaseImpl struct {
	brewInfrastructure infrastructure.BrewInfrastructure
	depsInfrastructure infrastructure.DepsInfrastructure
	printOutUC         PrintOutUsecase
	configUC           ConfigUsecase
}

func NewBrewUsecase(
	brewInfrastructure infrastructure.BrewInfrastructure,
	depsInfrastructure infrastructure.DepsInfrastructure,
	printOutUC PrintOutUsecase,
	configUC ConfigUsecase,
) *BrewUsecaseImpl {
	return &BrewUsecaseImpl{
		brewInfrastructure: brewInfrastructure,
		depsInfrastructure: depsInfrastructure,
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

func (b *BrewUsecaseImpl) DumpTmpBrewBundle() error {
	err := b.brewInfrastructure.DumpTmpBrewBundle(*b.printOutUC.GetOut(), *b.printOutUC.GetError())
	if err != nil {
		return errors.Wrap(err, "brew usecase: failed to dump Brewfile.tmp")
	}

	return nil
}

func (b *BrewUsecaseImpl) CheckDiffBrewBundle(
	bundlePath string,
	tmpPath string,
) ([]domain.BrewBundle, []domain.BrewBundle, error) {
	bundles, err := b.brewInfrastructure.ReadBrewBundle(bundlePath)
	if err != nil {
		return nil, nil, errors.Wrap(err, "brew usecase: failed to read Brewfile")
	}

	tmpBundles, err := b.brewInfrastructure.ReadBrewBundle(tmpPath)
	if err != nil {
		return nil, nil, errors.Wrap(err, "brew usecase: failed to read Brewfile.tmp")
	}

	tmpBundlesMap := make(map[string]bool)

	var diffBundles []domain.BrewBundle

	for _, bundle := range bundles {
		isFound := false

		for _, tmpBundle := range tmpBundles {
			if bundle.Name == tmpBundle.Name && bundle.BundleType == tmpBundle.BundleType {
				isFound = true
				tmpBundlesMap[bundle.Name] = true
			}
		}

		if !isFound {
			diffBundles = append(diffBundles, bundle)
		}
	}

	var diffTmpBundles []domain.BrewBundle

	for _, tmpBundle := range tmpBundles {
		if _, ok := tmpBundlesMap[tmpBundle.Name]; !ok {
			diffTmpBundles = append(diffTmpBundles, tmpBundle)
		}
	}

	return diffBundles, diffTmpBundles, nil
}

func (b *BrewUsecaseImpl) CleanupBrewBundle(isForce bool) error {
	err := b.brewInfrastructure.CleanupBrewBundle(isForce, *b.printOutUC.GetOut(), *b.printOutUC.GetError())
	if err != nil {
		return errors.Wrap(err, "brew usecase: failed to cleanup Brewfile")
	}

	return nil
}
