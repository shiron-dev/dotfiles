package infrastructure_test

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/shiron-dev/dotfiles/scripts/dofy/internal/di"
	"github.com/shiron-dev/dotfiles/scripts/dofy/internal/domain"
	"github.com/shiron-dev/dotfiles/scripts/dofy/internal/test/util"
)

// Related to `../test/data/brew_test.brewfile`.
//
//nolint:gochecknoglobals
var testBundles = []domain.BrewBundle{
	{
		Name:       "gh",
		Others:     []string{},
		BundleType: domain.BrewBundleTypeFormula,
		Categories: []string{"cat 1", "cat 1.1", "cat 1.1.1"},
	},
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

func TestBrewInfrastructureImpl_InstallHomebrew(t *testing.T) {
	t.Parallel()

	if !util.IsCI() {
		t.Skip("skipping test; not running on CI")
	}

	type args struct {
		ctx context.Context
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"no error", args{t.Context()}, false},
	}

	//nolint:paralleltest
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			infra, err := di.InitializeTestInfrastructureSet(os.Stdout, os.Stderr)
			if err != nil {
				t.Fatal(err)
			}

			b := infra.BrewInfrastructure
			sout := &bytes.Buffer{}
			serror := &bytes.Buffer{}

			if err := b.InstallHomebrew(tt.args.ctx, sout, serror); (err != nil) != tt.wantErr {
				t.Errorf("BrewInfrastructureImpl.InstallHomebrew() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
		})
	}
}

func TestBrewInfrastructureImpl_SetHomebrewEnv(t *testing.T) {
	t.Parallel()

	type args struct {
		goos string
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"error", args{""}, true},
		{"linux", args{"linux"}, runtime.GOOS != "linux"},
		{"darwin", args{"darwin"}, runtime.GOOS != "darwin"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			infra, err := di.InitializeTestInfrastructureSet(os.Stdout, os.Stderr)
			if err != nil {
				t.Fatal(err)
			}

			b := infra.BrewInfrastructure
			if err := b.SetHomebrewEnv(tt.args.goos); (err != nil) != tt.wantErr {
				t.Errorf("BrewInfrastructureImpl.SetHomebrewEnv() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

//nolint:paralleltest
func TestBrewInfrastructureImpl_InstallFormula(t *testing.T) {
	type args struct {
		formula string
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"error", args{"not_exist_formula"}, true},
		{"go", args{"go"}, false},
	}

	//nolint:paralleltest
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			infra, err := di.InitializeTestInfrastructureSet(os.Stdout, os.Stderr)
			if err != nil {
				t.Fatal(err)
			}

			if tt.wantErr {
				//nolint:gosec
				cmd := exec.Command("brew", "info", tt.args.formula)
				if err = cmd.Run(); err == nil {
					t.Fatalf("expected error, got nil")
				}
			}

			b := infra.BrewInfrastructure
			sout := &bytes.Buffer{}
			serror := &bytes.Buffer{}

			if err := b.InstallFormula(tt.args.formula, sout, serror); (err != nil) != tt.wantErr {
				t.Errorf("BrewInfrastructureImpl.InstallFormula() error = %v, wantErr %v, serror %s",
					err, tt.wantErr, serror.String())

				return
			}
		})
	}
}

//nolint:paralleltest
func TestBrewInfrastructureImpl_InstallTap(t *testing.T) {
	type args struct {
		formula string
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"error", args{"not_exist_formula"}, true},
		{"shiron-dev/tap", args{"shiron-dev/tap"}, false},
	}

	//nolint:paralleltest
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			infra, err := di.InitializeTestInfrastructureSet(os.Stdout, os.Stderr)
			if err != nil {
				t.Fatal(err)
			}

			if tt.wantErr {
				//nolint:gosec
				cmd := exec.Command("brew", "info", tt.args.formula)
				if err = cmd.Run(); err == nil {
					t.Fatalf("expected error, got nil")
				}
			}

			b := infra.BrewInfrastructure
			sout := &bytes.Buffer{}
			serror := &bytes.Buffer{}

			if err := b.InstallTap(tt.args.formula, sout, serror); (err != nil) != tt.wantErr {
				t.Errorf("BrewInfrastructureImpl.InstallTap() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
		})
	}
}

func TestBrewInfrastructureImpl_InstallByMas(t *testing.T) {
	t.Parallel()

	t.Skip("skipping test; not running on linux")
}

func TestBrewInfrastructureImpl_DumpTmpBrewBundle(t *testing.T) {
	t.Parallel()

	type args struct {
		path  string
		isMac bool
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"no error", args{filepath.Join(t.TempDir(), "/Brewfile.tmp"), false}, false},
		// {"mac mode", args{filepath.Join(t.TempDir(), "/Brewfile.tmp"), true}, runtime.GOOS != "darwin"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			infra, err := di.InitializeTestInfrastructureSet(os.Stdout, os.Stderr)
			if err != nil {
				t.Fatal(err)
			}

			b := infra.BrewInfrastructure
			sout := &bytes.Buffer{}
			serror := &bytes.Buffer{}

			if err := b.DumpTmpBrewBundle(tt.args.path, tt.args.isMac, sout, serror); (err != nil) != tt.wantErr {
				t.Errorf("BrewInfrastructureImpl.DumpTmpBrewBundle() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
		})
	}
}

//nolint:paralleltest
func TestBrewInfrastructureImpl_InstallBrewBundle(t *testing.T) {
	path := filepath.Join(t.TempDir(), "/Brewfile")

	//nolint:gosec
	if file, err := os.Create(path); err != nil {
		t.Fatal(err)
	} else {
		_, err = file.WriteString("brew \"go\"\n")
		if err != nil {
			t.Fatal(err)
		}
	}

	type args struct {
		path string
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"no error", args{path}, false},
		{"error", args{filepath.Join(t.TempDir(), "/not_exist")}, true},
	}

	//nolint:paralleltest
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			infra, err := di.InitializeTestInfrastructureSet(os.Stdout, os.Stderr)
			if err != nil {
				t.Fatal(err)
			}

			b := infra.BrewInfrastructure
			sout := &bytes.Buffer{}
			serror := &bytes.Buffer{}

			if err := b.InstallBrewBundle(tt.args.path, sout, serror); (err != nil) != tt.wantErr {
				t.Errorf("BrewInfrastructureImpl.InstallBrewBundle() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
		})
	}
}

func TestBrewInfrastructureImpl_CleanupBrewBundle(t *testing.T) {
	t.Parallel()

	t.Skip("skipping test; not running")
}

func TestBrewInfrastructureImpl_ReadBrewBundle(t *testing.T) {
	t.Parallel()

	path, err := filepath.Abs("../test/data/brew_test.brewfile")
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		path string
	}

	tests := []struct {
		name    string
		args    args
		want    []domain.BrewBundle
		wantErr bool
	}{
		{"no error", args{path}, testBundles, false},
		{"error", args{filepath.Join(t.TempDir(), "/not_exist")}, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			infra, err := di.InitializeTestInfrastructureSet(os.Stdout, os.Stderr)
			if err != nil {
				t.Fatal(err)
			}

			b := infra.BrewInfrastructure
			got, err := b.ReadBrewBundle(tt.args.path)

			if (err != nil) != tt.wantErr {
				t.Errorf("BrewInfrastructureImpl.ReadBrewBundle() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BrewInfrastructureImpl.ReadBrewBundle() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBrewInfrastructureImpl_WriteBrewBundle(t *testing.T) {
	t.Parallel()

	type args struct {
		path    string
		bundles []domain.BrewBundle
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"no error", args{filepath.Join(t.TempDir(), "/Brewfile"), testBundles}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			infra, err := di.InitializeTestInfrastructureSet(os.Stdout, os.Stderr)
			if err != nil {
				t.Fatal(err)
			}

			b := infra.BrewInfrastructure
			if err := b.WriteBrewBundle(tt.args.path, tt.args.bundles); (err != nil) != tt.wantErr {
				t.Errorf("BrewInfrastructureImpl.WriteBrewBundle() error = %v, wantErr %v", err, tt.wantErr)
			}

			file, err := os.ReadFile(tt.args.path)
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
		})
	}
}
