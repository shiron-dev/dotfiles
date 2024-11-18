package infrastructure_test

import (
	"os"
	"os/exec"
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

	const notInstalledCommand = "not_installed_command"

	_, err = exec.LookPath(notInstalledCommand)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	installed := deps.CheckInstalled(notInstalledCommand)

	if installed {
		t.Fatal("not_installed is installed")
	}

	const installedCommand = "git"

	_, err = exec.LookPath(installedCommand)
	if err != nil {
		t.Fatal(err)
	}

	installed = deps.CheckInstalled(installedCommand)

	if !installed {
		t.Fatal("git is not installed")
	}
}
