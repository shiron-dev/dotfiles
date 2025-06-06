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

	_, err = file.WriteString("test")
	if err != nil {
		t.Fatal(err)
	}

	return gitRepo, path
}

func TestGitInfrastructureImpl_SetGitDir(t *testing.T) {
	t.Parallel()

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
	t.Parallel()

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
		{"no error", args{t.Context(), []string{filePath}}, gitRepo, false},
		{"error", args{t.Context(), []string{"not_exist"}}, t.TempDir(), true},
		{"not set git dir", args{t.Context(), []string{"."}}, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			infra, err := di.InitializeTestInfrastructureSet(os.Stdout, os.Stderr)
			if err != nil {
				t.Fatal(err)
			}

			git := infra.GitInfrastructure

			if tt.gitRepo != "" {
				git.SetGitDir(tt.gitRepo)
			}

			sout := &bytes.Buffer{}
			serror := &bytes.Buffer{}

			if err := git.GitDifftool(tt.args.ctx, sout, serror, tt.args.path...); (err != nil) != tt.wantErr {
				t.Errorf("GitInfrastructureImpl.GitDifftool() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
		})
	}
}

func TestGitInfrastructureImpl_CheckoutFile(t *testing.T) {
	t.Parallel()

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
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			infra, err := di.InitializeTestInfrastructureSet(os.Stdout, os.Stderr)
			if err != nil {
				t.Fatal(err)
			}

			git := infra.GitInfrastructure

			if tt.gitRepo != "" {
				git.SetGitDir(tt.gitRepo)
			}

			if !tt.wantErr {
				//nolint:gosec
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
					//nolint:gosec
					file, err := os.OpenFile(tt.args.path, os.O_WRONLY, 0o666)
					if err != nil {
						t.Fatal(err)
					}

					defer func() {
						if err := file.Close(); err != nil {
							t.Fatal(err)
						}
					}()

					if _, err = file.WriteString("test"); err != nil {
						t.Fatal(err)
					}
				}
			}

			if err := git.CheckoutFile(tt.args.path); (err != nil) != tt.wantErr {
				t.Errorf("GitInfrastructureImpl.CheckoutFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGitInfrastructureImpl_IsGitDiff(t *testing.T) {
	t.Parallel()

	gitRepo, filePath := makeTestFile(t)

	type args struct {
		path []string
	}

	tests := []struct {
		name    string
		args    args
		gitRepo string
		want    bool
		wantErr bool
	}{
		{"no error", args{[]string{filePath}}, gitRepo, false, false},
		{"no git repo", args{[]string{filePath}}, "", false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			infra, err := di.InitializeTestInfrastructureSet(os.Stdout, os.Stderr)
			if err != nil {
				t.Fatal(err)
			}

			git := infra.GitInfrastructure

			if tt.gitRepo != "" {
				git.SetGitDir(tt.gitRepo)
			}

			got, err := git.IsGitDiff(tt.args.path...)
			if (err != nil) != tt.wantErr {
				t.Errorf("GitInfrastructureImpl.IsGitDiff() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if got != tt.want {
				t.Errorf("GitInfrastructureImpl.IsGitDiff() = %v, want %v", got, tt.want)
			}
		})
	}
}
