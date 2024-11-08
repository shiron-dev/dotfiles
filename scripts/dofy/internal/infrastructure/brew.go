package infrastructure

import (
	"context"
	"io"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"strings"
)

type BrewInfrastructure interface {
	InstallHomebrew(ctx context.Context, sout io.Writer, serror io.Writer) error
	SetHomebrewEnv(brewPath string) error
	InstallFormula(pkg string) error
	InstallBrewBundle(sout io.Writer, serror io.Writer) error
}

type BrewInfrastructureImpl struct{}

func NewBrewInfrastructure() *BrewInfrastructureImpl {
	return &BrewInfrastructureImpl{}
}

type BrewError struct {
	err error
}

func (e *BrewError) Error() string {
	return "BrewUC: " + e.err.Error()
}

func (b *BrewInfrastructureImpl) InstallHomebrew(ctx context.Context, sout io.Writer, serror io.Writer) error {
	url := "https://raw.githubusercontent.com/Homebrew/install/master/install.sh"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return &BrewError{err}
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return &BrewError{err}
	}
	defer resp.Body.Close()

	bytes, _ := io.ReadAll(resp.Body)

	//nolint:gosec
	cmd := exec.Command("/bin/bash", "-c", string(bytes))
	cmd.Stdout = sout
	cmd.Stderr = serror

	if err = cmd.Run(); err != nil {
		return &BrewError{err}
	}

	return nil
}

func (b *BrewInfrastructureImpl) SetHomebrewEnv(brewPath string) error {
	cmd := exec.Command(brewPath, "shellenv")

	shellenv, err := cmd.Output()
	if err != nil {
		return &BrewError{err}
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

func (b *BrewInfrastructureImpl) InstallFormula(formula string) error {
	cmd := exec.Command("brew", "install", formula)
	if err := cmd.Run(); err != nil {
		return &BrewError{err}
	}

	return nil
}

// func (b *BrewInfrastructureImpl) DumpTmpBrewBundle(sout io.Writer, serror io.Writer) error {
// 	usr, _ := user.Current()
// 	path := usr.HomeDir + "/projects/dotfiles/data/Brewfile.tmp"

// 	if _, err := os.Stat(path); err == nil {
// 		os.Remove(path)
// 	}

// 	cmd := exec.Command("brew", "bundle", "dump", "--tap", "--formula", "--cask", "--mas", "--file", path)
// 	cmd.Stdout = sout
// 	cmd.Stderr = serror

// 	if err := cmd.Run(); err != nil {
// 		return &BrewError{err}
// 	}

// 	return nil
// }

func (b *BrewInfrastructureImpl) InstallBrewBundle(sout io.Writer, serror io.Writer) error {
	usr, _ := user.Current()
	//nolint:gosec
	cmd := exec.Command("brew", "bundle", "--no-lock", "--file", usr.HomeDir+"/projects/dotfiles/data/Brewfile")
	cmd.Stdout = sout
	cmd.Stderr = serror

	if err := cmd.Run(); err != nil {
		return &BrewError{err}
	}

	return nil
}

// func checkBrewBundle() {
// 	usr, _ := user.Current()
// 	cmd := exec.Command("brew", "bundle", "check", "--file", usr.HomeDir+"/projects/dotfiles/data/Brewfile")
// 	cmd.Stdout = infrastructure.Out
// 	cmd.Stderr = infrastructure.Error
// 	err := cmd.Run()
// 	if err != nil {
// 		panic(err)
// 	}
// }

// func cleanupBrewBundle(isForce bool) {
// 	usr, _ := user.Current()
// 	forceFlag := ""
// 	if isForce {
// 		forceFlag = "--force"
// 	}
// 	cmd := exec.Command("brew", "bundle", "cleanup", forceFlag, "--file", usr.HomeDir+"/projects/dotfiles/data/Brewfile")
// 	cmd.Stdout = infrastructure.Out
// 	cmd.Stderr = infrastructure.Error
// 	err := cmd.Run()
// 	if err != nil {
// 		panic(err)
// 	}
// }

// type BrewBundleType uint

// const (
// 	BrewBundleTypeTap BrewBundleType = iota
// 	BrewBundleTypeFormula
// 	BrewBundleTypeCask
// 	BrewBundleTypeMas
// )

// type BrewBundle struct {
// 	name       string
// 	bundleType BrewBundleType
// }

// func checkDiffBrewBundle(bundlePath string, tmpPath string) ([]BrewBundle, []BrewBundle) {
// 	bundles := readBrewBundle(bundlePath)
// 	tmpBundles := readBrewBundle(tmpPath)
// 	tmpBundlesMap := make(map[string]bool)
// 	var diffBundles []BrewBundle
// 	for _, bundle := range bundles {
// 		isFound := false
// 		for _, tmpBundle := range tmpBundles {
// 			if bundle.name == tmpBundle.name && bundle.bundleType == tmpBundle.bundleType {
// 				isFound = true
// 				tmpBundlesMap[bundle.name] = true
// 				break
// 			}
// 		}
// 		if !isFound {
// 			diffBundles = append(diffBundles, bundle)
// 		}
// 	}
// 	var diffTmpBundles []BrewBundle
// 	for _, tmpBundle := range tmpBundles {
// 		if _, ok := tmpBundlesMap[tmpBundle.name]; !ok {
// 			diffTmpBundles = append(diffTmpBundles, tmpBundle)
// 		}
// 	}

// 	return diffBundles, diffTmpBundles
// }

// func readBrewBundle(path string) []BrewBundle {
// 	file, err := os.Open(path)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer file.Close()
// 	scanner := bufio.NewScanner(file)
// 	var bundles []BrewBundle
// 	for scanner.Scan() {
// 		line := scanner.Text()
// 		sp := strings.Split(line, " ")
// 		if len(sp) < 2 || sp[0] == "#" {
// 			continue
// 		}
// 		var bundleType BrewBundleType
// 		switch sp[0] {
// 		case "tap":
// 			bundleType = BrewBundleTypeTap
// 		case "brew":
// 			bundleType = BrewBundleTypeFormula
// 		case "cask":
// 			bundleType = BrewBundleTypeCask
// 		case "mas":
// 			bundleType = BrewBundleTypeMas
// 		default:
// 			continue
// 		}
// 		bundles = append(bundles, BrewBundle{
// 			name:       strings.Trim(sp[1], "\""),
// 			bundleType: bundleType,
// 		})
// 	}
// 	return bundles
// }
