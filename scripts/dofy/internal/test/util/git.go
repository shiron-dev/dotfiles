package util

import (
	"os"
	"os/exec"
	"testing"
)

const gitRepoPermission = 0o755

func MakeGitRepo(t *testing.T) string {
	t.Helper()

	workingDir := t.TempDir()
	path := workingDir + "/git_temp"

	err := os.MkdirAll(path, gitRepoPermission)
	if err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command("git", "init")
	cmd.Dir = path

	if err = cmd.Run(); err != nil {
		t.Fatal(err)
	}

	return path
}
