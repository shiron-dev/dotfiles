package infrastructure_test

import (
	"os"
	"strconv"
	"testing"

	"github.com/shiron-dev/dotfiles/scripts/dofy/internal/di"
	"github.com/shiron-dev/dotfiles/scripts/dofy/internal/test/util"
)

func createFile(t *testing.T, path string, content string) {
	t.Helper()

	file, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		t.Fatal(err)
	}
}

func TestReadFile(t *testing.T) {
	t.Parallel()

	infra, err := di.InitializeTestInfrastructureSet(os.Stdout, os.Stderr)
	if err != nil {
		t.Fatal(err)
	}

	file := infra.FileInfrastructure

	testStrs := []string{
		"test1",
		"test2",
		"test3\nt3",
	}

	dir := t.TempDir()

	for i, testStr := range testStrs {
		path := dir + "/test" + strconv.Itoa(i) + ".txt"

		createFile(t, path, testStr)

		content, err := file.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}

		if string(content) != testStr {
			t.Fatalf("expected %s, got %s", testStr, content)
		}
	}

	_, err = file.ReadFile(dir + "/not_exist.txt")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	path, err := util.MakeUnOpenableFile(t)
	if err != nil {
		t.Fatal(err)
	}

	_, err = file.ReadFile(path)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestWriteFile(t *testing.T) {
	t.Parallel()

	infra, err := di.InitializeTestInfrastructureSet(os.Stdout, os.Stderr)
	if err != nil {
		t.Fatal(err)
	}

	file := infra.FileInfrastructure

	testStrs := []string{
		"test1",
		"test2",
		"test3\nt3",
	}

	dir := t.TempDir()

	for i, testStr := range testStrs {
		path := dir + "/test" + strconv.Itoa(i) + ".txt"

		err := file.WriteFile(path, []byte(testStr))
		if err != nil {
			t.Fatal(err)
		}

		content, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}

		if string(content) != testStr {
			t.Fatalf("expected %s, got %s", testStr, content)
		}
	}

	path, err := util.MakeUnOpenableFile(t)
	if err != nil {
		t.Fatal(err)
	}

	err = file.WriteFile(path, []byte("test"))
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
