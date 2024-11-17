package infrastructure_test

import (
	"os"
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
