package deps

import (
	"bufio"
	"os"
	"os/exec"
	"os/user"
	"strings"

	"github.com/shiron-dev/dotfiles/scripts/cmd/printout"
)

func installHomebrew() {
	cmd := exec.Command("/bin/bash", "-c", "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install.sh)")
	cmd.Stdout = printout.Out
	cmd.Run()
}

func installWithBrew(pkg string) {
	cmd := exec.Command("brew", "install", pkg)
	cmd.Run()
}

func dumpTmpBrewBundle() {
	usr, _ := user.Current()
	cmd := exec.Command("brew", "bundle", "dump", "--tap", "--formula", "--cask", "--mas", "--file", usr.HomeDir+"/projects/dotfiles/data/Brewfile.tmp")
	cmd.Stdout = printout.Out
	cmd.Run()
}

func installBrewBundle() {
	usr, _ := user.Current()
	cmd := exec.Command("brew", "bundle", "--no-lock", "--file", usr.HomeDir+"/projects/dotfiles/data/Brewfile")
	cmd.Stdout = printout.Out
	cmd.Run()
}

func checkBrewBundle() {
	usr, _ := user.Current()
	cmd := exec.Command("brew", "bundle", "check", "--file", usr.HomeDir+"/projects/dotfiles/data/Brewfile")
	cmd.Stdout = printout.Out
	cmd.Run()
}

func cleanupBrewBundle(isForce bool) {
	usr, _ := user.Current()
	forceFlag := ""
	if isForce {
		forceFlag = "--force"
	}
	cmd := exec.Command("brew", "bundle", "cleanup", forceFlag, "--file", usr.HomeDir+"/projects/dotfiles/data/Brewfile")
	cmd.Stdout = printout.Out
	cmd.Run()
}

type BrewBundleType uint

const (
	BrewBundleTypeTap BrewBundleType = iota
	BrewBundleTypeFormula
	BrewBundleTypeCask
	BrewBundleTypeMas
)

type BrewBundle struct {
	name       string
	bundleType BrewBundleType
}

func checkDiffBrewBundle(bundlePath string, tmpPath string) ([]BrewBundle, []BrewBundle) {
	bundles := readBrewBundle(bundlePath)
	tmpBundles := readBrewBundle(tmpPath)
	tmpBundlesMap := make(map[string]bool)
	var diffBundles []BrewBundle
	for _, bundle := range bundles {
		isFound := false
		for _, tmpBundle := range tmpBundles {
			if bundle.name == tmpBundle.name && bundle.bundleType == tmpBundle.bundleType {
				isFound = true
				tmpBundlesMap[bundle.name] = true
				break
			}
		}
		if !isFound {
			diffBundles = append(diffBundles, bundle)
		}
	}
	var diffTmpBundles []BrewBundle
	for _, tmpBundle := range tmpBundles {
		if _, ok := tmpBundlesMap[tmpBundle.name]; !ok {
			diffTmpBundles = append(diffTmpBundles, tmpBundle)
		}
	}

	return diffBundles, diffTmpBundles
}

func readBrewBundle(path string) []BrewBundle {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var bundles []BrewBundle
	for scanner.Scan() {
		line := scanner.Text()
		sp := strings.Split(line, " ")
		if len(sp) < 2 || sp[0] == "#" {
			continue
		}
		var bundleType BrewBundleType
		switch sp[0] {
		case "tap":
			bundleType = BrewBundleTypeTap
		case "brew":
			bundleType = BrewBundleTypeFormula
		case "cask":
			bundleType = BrewBundleTypeCask
		case "mas":
			bundleType = BrewBundleTypeMas
		default:
			continue
		}
		bundles = append(bundles, BrewBundle{
			name:       strings.Trim(sp[1], "\""),
			bundleType: bundleType,
		})
	}
	return bundles
}
