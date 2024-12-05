package infrastructure_test

import (
	"os"
	"testing"

	"github.com/shiron-dev/dotfiles/scripts/dofy/internal/di"
)

func TestAnsibleInfrastructureImpl_CheckPlaybook(t *testing.T) {
	t.Parallel()

	t.Skip("Skip this test because it requires ansible-playbook command")
}

func TestAnsibleInfrastructureImpl_RunPlaybook(t *testing.T) {
	t.Parallel()

	t.Skip("Skip this test because it requires ansible-playbook command")
}

func TestAnsibleInfrastructureImpl_SetWorkingDir(t *testing.T) {
	t.Parallel()

	type args struct {
		workingDir string
	}

	tests := []struct {
		name string
		args args
	}{
		{"normal", args{"test"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			infra, err := di.InitializeTestInfrastructureSet(os.Stdout, os.Stderr)
			if err != nil {
				t.Fatal(err)
			}

			a := infra.AnsibleInfrastructure

			a.SetWorkingDir(tt.args.workingDir)
		})
	}
}
