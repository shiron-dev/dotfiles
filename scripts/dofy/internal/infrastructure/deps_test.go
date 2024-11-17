package infrastructure_test

import (
	"os"
	"testing"

	"github.com/shiron-dev/dotfiles/scripts/dofy/internal/di"
)

func TestCheckInstalled(t *testing.T) {
	t.Parallel()

	infra, err := di.InitializeTestInfrastructureSet(os.Stdout, os.Stderr)
	if err != nil {
		t.Fatal(err)
	}

	deps := infra.DepsInfrastructure

	installed := deps.CheckInstalled("git")

	if !installed {
		t.Fatal("git is not installed")
	}

	installed = deps.CheckInstalled("not_installed")

	if installed {
		t.Fatal("not_installed is installed")
	}
}
