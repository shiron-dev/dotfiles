package infrastructure_test

import (
	"os"
	"os/exec"
	"runtime"
	"strings"
	"testing"

	"github.com/shiron-dev/dotfiles/scripts/dofy/internal/di"
)

func TestGetOS(t *testing.T) {
	t.Parallel()

	infra, err := di.InitializeTestInfrastructureSet(os.Stdout, os.Stderr)
	if err != nil {
		t.Fatal(err)
	}

	config := infra.ConfigInfrastructure

	osStr, err := config.GetOS()
	if err != nil {
		t.Fatal(err)
	}

	if osStr == "" {
		t.Fatal("os is empty")
	}

	t.Log(osStr)
}

func TestGetOsVersion(t *testing.T) {
	t.Parallel()

	infra, err := di.InitializeTestInfrastructureSet(os.Stdout, os.Stderr)
	if err != nil {
		t.Fatal(err)
	}

	config := infra.ConfigInfrastructure

	osVersion, err := config.GetOSVersion()
	if err != nil {
		t.Fatal(err)
	}

	if osVersion == "" {
		t.Fatal("osVersion is empty")
	}

	t.Log(osVersion)
}

func TestGetArch(t *testing.T) {
	t.Parallel()

	infra, err := di.InitializeTestInfrastructureSet(os.Stdout, os.Stderr)
	if err != nil {
		t.Fatal(err)
	}

	config := infra.ConfigInfrastructure

	arch, err := config.GetArch()
	if err != nil {
		t.Fatal(err)
	}

	if arch == "" {
		t.Fatal("arch is empty")
	}

	t.Log(arch)
}

func TestConfigInfrastructureImpl_GetOS(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		{"no error", runtime.GOOS, false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			infra, err := di.InitializeTestInfrastructureSet(os.Stdout, os.Stderr)
			if err != nil {
				t.Fatal(err)
			}

			c := infra.ConfigInfrastructure

			got, err := c.GetOS()
			if (err != nil) != tt.wantErr {
				t.Errorf("ConfigInfrastructureImpl.GetOS() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ConfigInfrastructureImpl.GetOS() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfigInfrastructureImpl_GetOSVersion(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		{"no error", func() string {
			osVersion, err := exec.Command("uname", "-r").Output()
			if err != nil {
				t.Fatal(err)
			}

			return strings.TrimSpace(string(osVersion))
		}(), false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			infra, err := di.InitializeTestInfrastructureSet(os.Stdout, os.Stderr)
			if err != nil {
				t.Fatal(err)
			}

			c := infra.ConfigInfrastructure

			got, err := c.GetOSVersion()
			if (err != nil) != tt.wantErr {
				t.Errorf("ConfigInfrastructureImpl.GetOSVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ConfigInfrastructureImpl.GetOSVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfigInfrastructureImpl_GetArch(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		{"no error", func() string {
			osVersion, err := exec.Command("uname", "-p").Output()
			if err != nil {
				t.Fatal(err)
			}

			return strings.TrimSpace(string(osVersion))
		}(), false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			infra, err := di.InitializeTestInfrastructureSet(os.Stdout, os.Stderr)
			if err != nil {
				t.Fatal(err)
			}

			c := infra.ConfigInfrastructure

			got, err := c.GetArch()
			if (err != nil) != tt.wantErr {
				t.Errorf("ConfigInfrastructureImpl.GetArch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ConfigInfrastructureImpl.GetArch() = %v, want %v", got, tt.want)
			}
		})
	}
}
