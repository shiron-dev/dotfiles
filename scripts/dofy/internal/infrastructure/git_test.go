package infrastructure_test

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/shiron-dev/dotfiles/scripts/dofy/internal/di"
	"github.com/shiron-dev/dotfiles/scripts/dofy/internal/infrastructure"
	"github.com/shiron-dev/dotfiles/scripts/dofy/internal/test/util"
)

func makeTestFile(t *testing.T) (string, string) {
	t.Helper()

	gitRepo := util.MakeGitRepo(t)
	path := gitRepo + "/test"

	file, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	_, err = file.WriteString("test")
	if err != nil {
		t.Fatal(err)
	}

	return gitRepo, path
}

func TestGitDifftool(t *testing.T) {
	t.Parallel()

	infra, err := di.InitializeTestInfrastructureSet(os.Stdout, os.Stderr)
	if err != nil {
		t.Fatal(err)
	}

	git := infra.GitInfrastructure

	gitRepo, filePath := makeTestFile(t)

	err = git.GitDifftool(context.Background(), os.Stdout, os.Stderr, filePath)
	if !errors.Is(err, infrastructure.ErrGitDirNotSet) {
		t.Fatal(err)
	}

	git.SetGitDir(gitRepo)

	err = git.GitDifftool(context.Background(), os.Stdout, os.Stderr, filePath)
	if err != nil {
		t.Fatal(err)
	}

	err = git.GitDifftool(context.Background(), os.Stdout, os.Stderr, filepath.Join(gitRepo, "not_exist"))
	if err == nil {
		t.Fatal("difftool should fail")
	}
}

func TestCheckoutFile(t *testing.T) {
	t.Parallel()

	infra, err := di.InitializeTestInfrastructureSet(os.Stdout, os.Stderr)
	if err != nil {
		t.Fatal(err)
	}

	git := infra.GitInfrastructure

	gitRepo, filePath := makeTestFile(t)

	if err := git.CheckoutFile(filePath); !errors.Is(err, infrastructure.ErrGitDirNotSet) {
		t.Fatal(err)
	}

	git.SetGitDir(gitRepo)

	if err := git.CheckoutFile(filePath); err == nil {
		t.Fatal("checkout file should fail")
	}

	cmd := exec.Command("git", "add", filePath)
	cmd.Dir = gitRepo

	if err := cmd.Run(); err != nil {
		t.Fatal(err)
	}

	cmd = exec.Command("git", "commit", "-m", "test")
	cmd.Dir = gitRepo
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		t.Fatal(err)
	}

	if err := git.CheckoutFile(filePath); err != nil {
		t.Fatal(err)
	}

	file, err := os.OpenFile(filePath, os.O_WRONLY, 0o666)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	if _, err = file.WriteString("test"); err != nil {
		t.Fatal(err)
	}

	if err = git.CheckoutFile(filePath); err != nil {
		t.Fatal(err)
	}
}
