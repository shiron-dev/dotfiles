package infrastructure_test

import (
	"os"
	"os/exec"
	"runtime"
	"strings"
	"testing"

	"github.com/shiron-dev/dotfiles/scripts/dofy/internal/di"
)

func TestConfigInfrastructureImpl_GetOS(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		{"no error", runtime.GOOS, false},
	}

	for _, tt := range tests {
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
	t.Parallel()

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
	t.Parallel()

	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		{"no error", func() string {
			arch, err := exec.Command("uname", "-p").Output()
			if err != nil {
				t.Fatal(err)
			}

			return strings.TrimSpace(string(arch))
		}(), false},
	}

	for _, tt := range tests {
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
