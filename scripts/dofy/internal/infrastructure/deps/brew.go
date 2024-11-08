package deps

// import (
// 	"bufio"
// 	"github.com/shiron-dev/dotfiles/scripts/dofy/internal/infrastructure"
// 	"io"
// 	"net/http"
// 	"os"
// 	"os/exec"
// 	"os/user"
// 	"runtime"
// 	"strings"
// )

// func installHomebrew() {
// 	infrastructure.PrintMd(`
// ### Installing Homebrew
// `)

// 	url := "https://raw.githubusercontent.com/Homebrew/install/master/install.sh"

// 	resp, err := http.Get(url)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer resp.Body.Close()
// 	bytes, _ := io.ReadAll(resp.Body)

// 	cmd := exec.Command("/bin/bash", "-c", string(bytes))
// 	cmd.Stdout = infrastructure.Out
// 	cmd.Stderr = infrastructure.Error
// 	err = cmd.Run()
// 	if err != nil {
// 		panic(err)
// 	}

// 	infrastructure.PrintMd(`
// ### Set Homebrew environment
// `)

// 	var brewPath string
// 	switch runtime.GOOS {
// 	case "darwin":
// 		brewPath = "/opt/homebrew/bin/brew"
// 	case "linux":
// 		brewPath = "/home/linuxbrew/.linuxbrew/bin/brew"
// 	}
// 	cmd = exec.Command("/bin/bash", "-c", `(echo; echo 'eval "$(/home/linuxbrew/.linuxbrew/bin/brew shellenv)"') >> ~/.bashrc`)
// 	cmd.Stdout = infrastructure.Out
// 	cmd.Stderr = infrastructure.Error
// 	err = cmd.Run()
// 	if err != nil {
// 		panic(err)
// 	}

// 	cmd = exec.Command(brewPath, "shellenv")
// 	shellenv, _ := cmd.Output()
// 	for _, line := range strings.Split(string(shellenv), "\n") {
// 		if strings.HasPrefix(line, "export PATH=") {
// 			cmd = exec.Command("sh", "-c", "echo "+strings.Replace(line, "export PATH=", "", 1))
// 			out, _ := cmd.Output()
// 			os.Setenv("PATH", strings.Trim(string(out), "\""))
// 		}
// 	}
// }

// func installWithBrew(pkg string) {
// 	cmd := exec.Command("brew", "install", pkg)
// 	err := cmd.Run()
// 	if err != nil {
// 		panic(err)
// 	}
// }

// func dumpTmpBrewBundle() {
// 	usr, _ := user.Current()
// 	path := usr.HomeDir + "/projects/dotfiles/data/Brewfile.tmp"

// 	if _, err := os.Stat(path); err == nil {
// 		os.Remove(path)
// 	}

// 	cmd := exec.Command("brew", "bundle", "dump", "--tap", "--formula", "--cask", "--mas", "--file", path)
// 	cmd.Stdout = infrastructure.Out
// 	cmd.Stderr = infrastructure.Error
// 	err := cmd.Run()
// 	if err != nil {
// 		panic(err)
// 	}
// }

// func installBrewBundle() {
// 	usr, _ := user.Current()
// 	cmd := exec.Command("brew", "bundle", "--no-lock", "--file", usr.HomeDir+"/projects/dotfiles/data/Brewfile")
// 	cmd.Stdout = infrastructure.Out
// 	cmd.Stderr = infrastructure.Error
// 	err := cmd.Run()
// 	if err != nil {
// 		panic(err)
// 	}
// }

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
