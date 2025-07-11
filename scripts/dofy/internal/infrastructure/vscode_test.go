package infrastructure_test

import (
	"os"
	"testing"

	"github.com/shiron-dev/dotfiles/scripts/dofy/internal/di"
	"github.com/shiron-dev/dotfiles/scripts/dofy/internal/test/util"
)

func TestVSCodeInfrastructureImpl_ListExtensions(t *testing.T) {
	t.Parallel()

	if util.IsCI() {
		t.Skip("skipping test; running on CI")
	}

	tests := []struct {
		name    string
		wantErr bool
	}{
		{"no error", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			infra, err := di.InitializeTestInfrastructureSet(os.Stdout, os.Stderr)
			if err != nil {
				t.Fatal(err)
			}

			v := infra.VSCodeInfrastructure

			_, err = v.ListExtensions()

			if (err != nil) != tt.wantErr {
				t.Errorf("VSCodeInfrastructureImpl.ListExtensions() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
		})
	}
}
