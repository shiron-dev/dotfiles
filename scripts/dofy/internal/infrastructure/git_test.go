package infrastructure_test

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"testing"

	"github.com/shiron-dev/dotfiles/scripts/dofy/internal/di"
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

func TestGitInfrastructureImpl_SetGitDir(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
	}{
		{"test", args{"test"}},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			infra, err := di.InitializeTestInfrastructureSet(os.Stdout, os.Stderr)
			if err != nil {
				t.Fatal(err)
			}

			g := infra.GitInfrastructure
			g.SetGitDir(tt.args.path)
		})
	}
}

func TestGitInfrastructureImpl_GitDifftool(t *testing.T) {
	gitRepo, filePath := makeTestFile(t)

	type args struct {
		ctx  context.Context
		path []string
	}
	tests := []struct {
		name    string
		args    args
		gitRepo string
		wantErr bool
	}{
		{"no error", args{context.Background(), []string{filePath}}, gitRepo, false},
		{"error", args{context.Background(), []string{"not_exist"}}, t.TempDir(), true},
		{"not set git dir", args{context.Background(), []string{"."}}, "", true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			infra, err := di.InitializeTestInfrastructureSet(os.Stdout, os.Stderr)
			if err != nil {
				t.Fatal(err)
			}

			g := infra.GitInfrastructure

			if tt.gitRepo != "" {
				g.SetGitDir(tt.gitRepo)
			}

			sout := &bytes.Buffer{}
			serror := &bytes.Buffer{}
			if err := g.GitDifftool(tt.args.ctx, sout, serror, tt.args.path...); (err != nil) != tt.wantErr {
				t.Errorf("GitInfrastructureImpl.GitDifftool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGitInfrastructureImpl_CheckoutFile(t *testing.T) {
	gitRepo, filePath := makeTestFile(t)
	gitRepo2, filePath2 := makeTestFile(t)
	gitRepo3, filePath3 := makeTestFile(t)

	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		gitRepo string
		wantErr bool
	}{
		{"no error", args{filePath}, gitRepo, false},
		{"no update", args{filePath2}, gitRepo2, false},
		{"error", args{filePath3}, gitRepo3, true},
		{"not exist", args{"not_exist"}, t.TempDir(), true},
		{"not set git dir", args{"."}, "", true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			infra, err := di.InitializeTestInfrastructureSet(os.Stdout, os.Stderr)
			if err != nil {
				t.Fatal(err)
			}

			g := infra.GitInfrastructure

			if tt.gitRepo != "" {
				g.SetGitDir(tt.gitRepo)
			}

			if !tt.wantErr {
				cmd := exec.Command("git", "add", tt.args.path)
				cmd.Dir = tt.gitRepo

				if err := cmd.Run(); err != nil {
					t.Fatal(err)
				}

				cmd = exec.Command("git", "commit", "-m", "test")
				cmd.Dir = tt.gitRepo
				cmd.Stderr = os.Stderr

				if err := cmd.Run(); err != nil {
					t.Fatal(err)
				}

				if tt.name != "no update" {
					file, err := os.OpenFile(tt.args.path, os.O_WRONLY, 0o666)
					if err != nil {
						t.Fatal(err)
					}
					defer file.Close()

					if _, err = file.WriteString("test"); err != nil {
						t.Fatal(err)
					}
				}
			}

			if err := g.CheckoutFile(tt.args.path); (err != nil) != tt.wantErr {
				t.Errorf("GitInfrastructureImpl.CheckoutFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
