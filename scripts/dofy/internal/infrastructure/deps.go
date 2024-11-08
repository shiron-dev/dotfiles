package infrastructure

import (
	"bufio"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
	"github.com/shiron-dev/dotfiles/scripts/dofy/internal/domain"
)

type DepsInfrastructure interface {
	CheckInstalled(name string) bool
	ReadBrewBundle(path string) ([]domain.BrewBundle, error)
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

func (d *DepsInfrastructureImpl) ReadBrewBundle(path string) ([]domain.BrewBundle, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "deps infrastructure: failed to open file")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var bundles []domain.BrewBundle

	for scanner.Scan() {
		line := scanner.Text()
		spBlank := strings.Split(line, " ")

		if len(spBlank) < 2 || spBlank[0] == "#" {
			continue
		}

		var bundleType domain.BrewBundleType

		switch spBlank[0] {
		case "tap":
			bundleType = domain.BrewBundleTypeTap
		case "brew":
			bundleType = domain.BrewBundleTypeFormula
		case "cask":
			bundleType = domain.BrewBundleTypeCask
		case "mas":
			bundleType = domain.BrewBundleTypeMas
		default:
			continue
		}

		bundles = append(bundles, domain.BrewBundle{
			Name:       strings.Trim(spBlank[1], "\""),
			BundleType: bundleType,
		})
	}

	return bundles, nil
}

func (d *DepsInfrastructureImpl) OpenWithCode(path ...string) error {
	args := []string{"-n", "-w"}
	args = append(args, path...)

	if err := exec.Command("code", args...).Run(); err != nil {
		return errors.Wrap(err, "deps infrastructure: failed to open with code")
	}

	return nil
}
