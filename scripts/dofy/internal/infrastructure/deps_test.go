package infrastructure_test

import (
	"os"
	"os/exec"
	"testing"

	"github.com/shiron-dev/dotfiles/scripts/dofy/internal/di"
)

func TestDepsInfrastructureImpl_CheckInstalled(t *testing.T) {
	t.Parallel()

	type args struct {
		name string
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{"not_installed", args{"not_installed"}, false},
		{"installed", args{"git"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			infra, err := di.InitializeTestInfrastructureSet(os.Stdout, os.Stderr)
			if err != nil {
				t.Fatal(err)
			}

			_, err = exec.LookPath(tt.args.name)
			installed := err == nil

			if installed != tt.want {
				t.Fatalf("expected %v, got %v", tt.want, installed)
			}

			d := infra.DepsInfrastructure
			if got := d.CheckInstalled(tt.args.name); got != tt.want {
				t.Errorf("DepsInfrastructureImpl.CheckInstalled() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDepsInfrastructureImpl_OpenWithCode(t *testing.T) {
	t.Parallel()

	_, err := exec.LookPath("code")
	hasCode := err == nil

	type args struct {
		path []string
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"no error", args{[]string{"-h"}}, !hasCode},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			infra, err := di.InitializeTestInfrastructureSet(os.Stdout, os.Stderr)
			if err != nil {
				t.Fatal(err)
			}

			d := infra.DepsInfrastructure
			if err := d.OpenWithCode(tt.args.path...); (err != nil) != tt.wantErr {
				t.Errorf("DepsInfrastructureImpl.OpenWithCode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
