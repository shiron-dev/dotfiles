package infrastructure_test

import (
	"os"
	"path/filepath"
	"reflect"
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

func TestFileInfrastructureImpl_ReadFile(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{"test", args{filepath.Join(t.TempDir(), "test.txt")}, []byte("test"), false},
		{"test 2 line", args{filepath.Join(t.TempDir(), "test.txt")}, []byte("test\nt"), false},
		{"not exist", args{filepath.Join(t.TempDir(), "not_exist.txt")}, nil, true},
		{"unopenable", args{filepath.Join(t.TempDir(), "unopenable")}, nil, true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			infra, err := di.InitializeTestInfrastructureSet(os.Stdout, os.Stderr)
			if err != nil {
				t.Fatal(err)
			}

			if tt.want != nil {
				createFile(t, tt.args.path, string(tt.want))
			}

			if tt.name == "unopenable" {
				path, err := util.MakeUnOpenableFile(t)
				if err != nil {
					t.Fatal(err)
				}
				tt.args.path = path
			}

			f := infra.FileInfrastructure
			got, err := f.ReadFile(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("FileInfrastructureImpl.ReadFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FileInfrastructureImpl.ReadFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFileInfrastructureImpl_WriteFile(t *testing.T) {
	type args struct {
		path string
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			infra, err := di.InitializeTestInfrastructureSet(os.Stdout, os.Stderr)
			if err != nil {
				t.Fatal(err)
			}

			f := infra.FileInfrastructure
			if err := f.WriteFile(tt.args.path, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("FileInfrastructureImpl.WriteFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
