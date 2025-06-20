package infrastructure_test

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/shiron-dev/dotfiles/scripts/dofy/internal/di"
	"github.com/shiron-dev/dotfiles/scripts/dofy/internal/test/util"
)

func createFile(t *testing.T, path string, content string) {
	t.Helper()

	//nolint:gosec
	file, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	_, err = file.WriteString(content)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFileInfrastructureImpl_ReadFile(t *testing.T) {
	t.Parallel()

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
	t.Parallel()

	path, err := util.MakeUnOpenableFile(t)
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		path string
		data []byte
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"test", args{filepath.Join(t.TempDir(), "test.txt"), []byte("test")}, false},
		{"test 2 line", args{filepath.Join(t.TempDir(), "test.txt"), []byte("test\nt")}, false},
		{"unopenable", args{path, []byte("abc")}, true},
	}

	for _, tt := range tests {
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
