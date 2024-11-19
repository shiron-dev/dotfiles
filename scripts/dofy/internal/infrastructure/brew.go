package infrastructure

import (
	"bufio"
	"context"
	"io"
	"net/http"
	"os"
	"os/exec"
	"sort"
	"strings"

	"github.com/pkg/errors"
	"github.com/shiron-dev/dotfiles/scripts/dofy/internal/domain"
)

type BrewInfrastructure interface {
	InstallHomebrew(ctx context.Context, sout io.Writer, serror io.Writer) error
	SetHomebrewEnv(goos string) error
	InstallFormula(pkg string, sout io.Writer, serror io.Writer) error
	DumpTmpBrewBundle(path string, isMac bool, sout io.Writer, serror io.Writer) error
	InstallBrewBundle(path string, sout io.Writer, serror io.Writer) error
	CleanupBrewBundle(path string, isForce bool, sout io.Writer, serror io.Writer) error
	ReadBrewBundle(path string) ([]domain.BrewBundle, error)
	WriteBrewBundle(path string, bundles []domain.BrewBundle) error
}

type BrewInfrastructureImpl struct{}

func NewBrewInfrastructure() *BrewInfrastructureImpl {
	return &BrewInfrastructureImpl{}
}

func (b *BrewInfrastructureImpl) InstallHomebrew(ctx context.Context, sout io.Writer, serror io.Writer) error {
	url := "https://raw.githubusercontent.com/Homebrew/install/master/install.sh"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return errors.Wrap(err, "brew infrastructure: failed to create request")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "brew infrastructure: failed to send request")
	}
	defer resp.Body.Close()

	bytes, _ := io.ReadAll(resp.Body)

	//nolint:gosec
	cmd := exec.Command("/bin/bash", "-c", string(bytes))
	cmd.Stdout = sout
	cmd.Stderr = serror

	if err = cmd.Run(); err != nil {
		return errors.Wrap(err, "brew infrastructure: failed to run command")
	}

	return nil
}

func (b *BrewInfrastructureImpl) SetHomebrewEnv(goos string) error {
	var brewPath string

	switch goos {
	case "darwin":
		brewPath = "/opt/homebrew/bin/brew"
	case "linux":
		brewPath = "/home/linuxbrew/.linuxbrew/bin/brew"
	}

	cmd := exec.Command(brewPath, "shellenv")

	shellenv, err := cmd.Output()
	if err != nil {
		return errors.Wrap(err, "brew infrastructure: failed to get shellenv")
	}

	for _, line := range strings.Split(string(shellenv), "\n") {
		if strings.HasPrefix(line, "export PATH=") {
			//nolint:gosec
			cmd = exec.Command("sh", "-c", "echo "+strings.Replace(line, "export PATH=", "", 1))
			out, _ := cmd.Output()
			os.Setenv("PATH", strings.Trim(string(out), "\""))
		}
	}

	return nil
}

func (b *BrewInfrastructureImpl) InstallFormula(formula string, sout io.Writer, serror io.Writer) error {
	cmd := exec.Command("brew", "install", formula)
	cmd.Stdout = sout
	cmd.Stderr = serror

	if err := cmd.Run(); err != nil {
		return errors.Wrap(err, "brew infrastructure: failed to run brew install command")
	}

	return nil
}

func (b *BrewInfrastructureImpl) DumpTmpBrewBundle(path string, isMac bool, sout io.Writer, serror io.Writer) error {
	if _, err := os.Stat(path); err == nil {
		os.Remove(path)
	}

	args := []string{"bundle", "dump", "--tap", "--formula"}

	if isMac {
		args = append(args, "--cask", "--mas")
	}

	args = append(args, "--file", path)

	cmd := exec.Command("brew", args...)
	cmd.Stdout = sout
	cmd.Stderr = serror

	if err := cmd.Run(); err != nil {
		return errors.Wrap(err, "brew infrastructure: failed to run brew bundle dump command")
	}

	return nil
}

func (b *BrewInfrastructureImpl) InstallBrewBundle(path string, sout io.Writer, serror io.Writer) error {
	cmd := exec.Command("brew", "bundle", "--no-lock", "--file", path)
	cmd.Stdout = sout
	cmd.Stderr = serror

	if err := cmd.Run(); err != nil {
		return errors.Wrap(err, "brew infrastructure: failed to run brew bundle command")
	}

	return nil
}

func (b *BrewInfrastructureImpl) CleanupBrewBundle(path string, isForce bool, sout io.Writer, serror io.Writer) error {
	forceFlag := ""

	if isForce {
		forceFlag = "--force"
	}

	cmd := exec.Command("brew", "bundle", "cleanup", forceFlag, "--file", path)
	cmd.Stdout = sout
	cmd.Stderr = serror

	if err := cmd.Run(); err != nil {
		return errors.Wrap(err, "brew infrastructure: failed to run brew bundle cleanup command")
	}

	return nil
}

func (b *BrewInfrastructureImpl) ReadBrewBundle(path string) ([]domain.BrewBundle, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "deps infrastructure: failed to open file")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var bundles []domain.BrewBundle

	lastCategories := []string{}

	for scanner.Scan() {
		line := scanner.Text()

		if len(line) == 0 {
			continue
		}

		if line[0] == '#' {
			lastCategories = getCategory(line, lastCategories)

			continue
		}

		spInd := strings.Index(line, " ")
		if spInd == -1 {
			continue
		}

		prefix := line[:spInd]
		formula := line[spInd+1:]
		cFormula := strings.Split(formula, ",")

		others := []string{}
		for _, c := range cFormula[1:] {
			others = append(others, strings.TrimSpace(c))
		}

		bundles = append(bundles, domain.BrewBundle{
			Name:       strings.TrimSpace(strings.Trim(strings.ReplaceAll(cFormula[0], ",", ""), "\"")),
			Others:     others,
			BundleType: domain.BrewBundleTypeFromString(prefix),
			Categories: append([]string{}, lastCategories...),
		})
	}

	return bundles, nil
}

func getCategory(line string, lastCategories []string) []string {
	count := 0

	for _, c := range line {
		if c == '#' {
			count++
		} else {
			break
		}
	}

	size := len(lastCategories)
	for i := count - 1; i < size; i++ {
		lastCategories = lastCategories[:len(lastCategories)-1]
	}

	lastCategories = append(lastCategories, strings.TrimSpace(line[count:]))

	return lastCategories
}

func (b *BrewInfrastructureImpl) WriteBrewBundle(path string, bundles []domain.BrewBundle) error {
	file, err := os.Create(path)
	if err != nil {
		return errors.Wrap(err, "deps infrastructure: failed to create file")
	}
	defer file.Close()

	if _, err := file.WriteString("# Brewfile made by dofy\n"); err != nil {
		return errors.Wrap(err, "deps infrastructure: failed to write file")
	}

	bundleMap := sortByCategories(bundles)

	lastCategories := []string{}

	writeString := ""

	for _, bundle := range bundleMap {
		for i, cate := range bundle.Categories {
			if len(lastCategories) > i && lastCategories[i] == cate {
				continue
			}

			for j := 0; j <= i; j++ {
				if j == 0 {
					writeString += "\n"
				}

				writeString += "#"
			}

			writeString += " " + cate + "\n"
		}

		lastCategories = bundle.Categories

		writeString += bundle.String() + "\n"
	}

	if _, err := file.WriteString(writeString); err != nil {
		return errors.Wrap(err, "deps infrastructure: failed to write file")
	}

	return nil
}

type cateKey string

func toCateKey(cate []string) cateKey {
	return cateKey(strings.Join(cate, ","))
}

func sortByCategories(bundles []domain.BrewBundle) []domain.BrewBundle {
	categoriesOrder := [][]string{}

	bundleMap := make(map[cateKey][]domain.BrewBundle)

	for _, bundle := range bundles {
		if _, ok := bundleMap[toCateKey(bundle.Categories)]; !ok {
			categoriesOrder = append(categoriesOrder, bundle.Categories)
		}

		bundleMap[toCateKey(bundle.Categories)] = append(bundleMap[toCateKey(bundle.Categories)], bundle)
	}

	for _, bundles := range bundleMap {
		sort.Slice(bundles, func(i, j int) bool {
			if bundles[i].BundleType != bundles[j].BundleType {
				return bundles[i].BundleType < bundles[j].BundleType
			}

			return bundles[i].Name < bundles[j].Name
		})
	}

	var sortedBundles []domain.BrewBundle

	for _, categories := range categoriesOrder {
		sortedBundles = append(sortedBundles, bundleMap[toCateKey(categories)]...)
	}

	return sortedBundles
}
