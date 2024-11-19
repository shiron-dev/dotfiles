package infrastructure_test

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/shiron-dev/dotfiles/scripts/dofy/internal/di"
	"github.com/shiron-dev/dotfiles/scripts/dofy/internal/domain"
)

// Related to `../test/data/brew_test.brewfile`.
//
//nolint:gochecknoglobals
var testBundles = []domain.BrewBundle{
	{
		Name:       "git",
		Others:     []string{},
		BundleType: domain.BrewBundleTypeFormula,
		Categories: []string{"cat 1", "cat 1.1", "cat 1.1.1"},
	},
	{
		Name:       "go",
		Others:     []string{},
		BundleType: domain.BrewBundleTypeFormula,
		Categories: []string{"cat 1", "cat 1.2"},
	},
	{
		Name:       "homebrew/cask",
		Others:     []string{},
		BundleType: domain.BrewBundleTypeTap,
		Categories: []string{"cat 2"},
	},
	{
		Name:       "mas",
		Others:     []string{"id: 1234567890"},
		BundleType: domain.BrewBundleTypeMas,
		Categories: []string{"cat 2"},
	},
}

func TestSetHomebrewEnv(t *testing.T) {
	t.Parallel()

	infra, err := di.InitializeTestInfrastructureSet(os.Stdout, os.Stderr)
	if err != nil {
		t.Fatal(err)
	}

	brew := infra.BrewInfrastructure

	err = brew.SetHomebrewEnv("")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	err = brew.SetHomebrewEnv(runtime.GOOS)
	if err != nil {
		t.Fatal(err)
	}
}

//nolint:paralleltest
func TestInstallFormula(t *testing.T) {
	infra, err := di.InitializeTestInfrastructureSet(os.Stdout, os.Stderr)
	if err != nil {
		t.Fatal(err)
	}

	brew := infra.BrewInfrastructure

	const notExistFormula = "not_exist_formula"

	cmd := exec.Command("brew", "info", notExistFormula)
	err = cmd.Run()

	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	outBuffer := &bytes.Buffer{}

	errBuffer := &bytes.Buffer{}

	err = brew.InstallFormula(notExistFormula, outBuffer, errBuffer)
	if err == nil {
		t.Fatal("expected error, got nil", outBuffer.String(), errBuffer.String())
	}

	const existFormula = "go"

	cmd = exec.Command("brew", "info", existFormula)
	if err = cmd.Run(); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	outBuffer.Reset()

	errBuffer.Reset()

	err = brew.InstallFormula(existFormula, outBuffer, errBuffer)
	if err != nil {
		t.Fatal(err, outBuffer.String(), errBuffer.String())
	}
}

func TestInstallTap(t *testing.T) {
	infra, err := di.InitializeTestInfrastructureSet(os.Stdout, os.Stderr)
	if err != nil {
		t.Fatal(err)
	}

	brew := infra.BrewInfrastructure

	const notExistFormula = "not_exist_formula"

	outBuffer := &bytes.Buffer{}

	errBuffer := &bytes.Buffer{}

	err = brew.InstallTap(notExistFormula, outBuffer, errBuffer)
	if err == nil {
		t.Fatal("expected error, got nil", outBuffer.String(), errBuffer.String())
	}

	const existFormula = "Homebrew/bundle"

	outBuffer.Reset()

	errBuffer.Reset()

	err = brew.InstallTap(existFormula, outBuffer, errBuffer)
	if err != nil {
		t.Fatal(err, outBuffer.String(), errBuffer.String())
	}
}

func TestInstallByMas(t *testing.T) {
	infra, err := di.InitializeTestInfrastructureSet(os.Stdout, os.Stderr)
	if err != nil {
		t.Fatal(err)
	}

	brew := infra.BrewInfrastructure

	const notExistFormula = "123"

	cmd := exec.Command("mas", "info", notExistFormula)
	err = cmd.Run()

	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	outBuffer := &bytes.Buffer{}

	errBuffer := &bytes.Buffer{}

	err = brew.InstallByMas(notExistFormula, outBuffer, errBuffer)
	if err == nil {
		t.Fatal("expected error, got nil", outBuffer.String(), errBuffer.String())
	}

	const existFormula = "497799835"

	cmd = exec.Command("mas", "info", existFormula)
	if err = cmd.Run(); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	outBuffer.Reset()

	errBuffer.Reset()

	err = brew.InstallByMas(existFormula, outBuffer, errBuffer)
	if err != nil {
		t.Fatal(err, outBuffer.String(), errBuffer.String())
	}
}

func TestDumpTmpBrewBundle(t *testing.T) {
	t.Parallel()

	infra, err := di.InitializeTestInfrastructureSet(os.Stdout, os.Stderr)
	if err != nil {
		t.Fatal(err)
	}

	brew := infra.BrewInfrastructure

	path := filepath.Join(t.TempDir(), "/Brewfile.tmp")

	outBuffer := &bytes.Buffer{}

	errBuffer := &bytes.Buffer{}

	err = brew.DumpTmpBrewBundle(path, false, outBuffer, errBuffer)
	if err != nil {
		t.Fatal(err, outBuffer.String(), errBuffer.String())
	}

	if _, err := os.Stat(path); err != nil {
		t.Fatal(err)
	}
}

func TestInstallBrewBundle(t *testing.T) {
	t.Parallel()

	infra, err := di.InitializeTestInfrastructureSet(os.Stdout, os.Stderr)
	if err != nil {
		t.Fatal(err)
	}

	brew := infra.BrewInfrastructure

	path := filepath.Join(t.TempDir(), "/Brewfile")

	if file, err := os.Create(path); err != nil {
		t.Fatal(err)
	} else {
		_, err = file.WriteString("brew \"go\"\n")
		if err != nil {
			t.Fatal(err)
		}

		err = file.Close()
		if err != nil {
			t.Fatal(err)
		}
	}

	outBuffer := &bytes.Buffer{}

	errBuffer := &bytes.Buffer{}

	err = brew.InstallBrewBundle(path, outBuffer, errBuffer)
	if err != nil {
		t.Fatal(err, outBuffer.String(), errBuffer.String())
	}
}

//nolint:cyclop
func TestReadBrewBundle(t *testing.T) {
	t.Parallel()

	infra, err := di.InitializeTestInfrastructureSet(os.Stdout, os.Stderr)
	if err != nil {
		t.Fatal(err)
	}

	brew := infra.BrewInfrastructure

	path, err := filepath.Abs("../test/data/brew_test.brewfile")
	if err != nil {
		t.Fatal(err)
	}

	bundles, err := brew.ReadBrewBundle(path)
	if err != nil {
		t.Fatal(err)
	}

	if len(bundles) != len(testBundles) {
		t.Fatalf("expected %d, got %d", len(testBundles), len(bundles))
	}

	for ind, bundle := range bundles {
		if bundle.Name != testBundles[ind].Name {
			t.Fatalf("expected %s, got %s", testBundles[ind].Name, bundle.Name)
		}

		if bundle.BundleType != testBundles[ind].BundleType {
			t.Fatalf("expected %d, got %d", testBundles[ind].BundleType, bundle.BundleType)
		}

		if len(bundle.Categories) != len(testBundles[ind].Categories) {
			t.Fatalf("expected %d, got %d", len(testBundles[ind].Categories), len(bundle.Categories))
		}

		for j, cat := range bundle.Categories {
			if cat != testBundles[ind].Categories[j] {
				t.Fatalf("expected %s, got %s", testBundles[ind].Categories[j], cat)
			}
		}
	}
}

func TestWriteBrewBundle(t *testing.T) {
	t.Parallel()

	infra, err := di.InitializeTestInfrastructureSet(os.Stdout, os.Stderr)
	if err != nil {
		t.Fatal(err)
	}

	brew := infra.BrewInfrastructure

	path := t.TempDir() + "/test.brewfile"

	if _, err := os.Stat(path); err == nil {
		t.Fatal("file already exists")
	}

	err = brew.WriteBrewBundle(path, testBundles)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(path); err != nil {
		t.Fatal(err)
	}

	file, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	correctFile, err := os.ReadFile("../test/data/brew_test.brewfile")
	if err != nil {
		t.Fatal(err)
	}

	if string(file) != string(correctFile) {
		t.Fatalf("expected %s, got %s", string(correctFile), string(file))
	}
}
