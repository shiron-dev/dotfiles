package util

import (
	"os"
	"testing"

	"github.com/pkg/errors"
)

func IsCI() bool {
	return os.Getenv("GITHUB_ACTIONS") == "true"
}

func MakeUnOpenableFile(t *testing.T) (string, error) {
	t.Helper()

	path := t.TempDir() + "/unopenable"

	//nolint:gosec
	file, err := os.Create(path)
	if err != nil {
		return path, errors.Wrap(err, "failed to create file")
	}

	defer func() {
		if err := file.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	err = os.Chmod(path, 0)

	return path, errors.Wrap(err, "failed to chmod")
}
