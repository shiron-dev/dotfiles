package usecase_test

import (
	"reflect"
	"testing"

	mock_infrastructure "github.com/shiron-dev/dotfiles/scripts/dofy/gen/mock/infrastructure"
	"github.com/shiron-dev/dotfiles/scripts/dofy/internal/di"
	"github.com/shiron-dev/dotfiles/scripts/dofy/internal/usecase"
	"go.uber.org/mock/gomock"
)

func TestConfigUsecaseImpl_ScanEnvInfo(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		want    *usecase.EnvInfo
		wantErr bool
	}{
		{"test", &usecase.EnvInfo{"testOS", "testOSVersion", "testArch", false}, false},
		{"darwin", &usecase.EnvInfo{"darwin", "testOSVersion", "testArch", true}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mCfg := mock_infrastructure.NewMockConfigInfrastructure(ctrl)

			mCfg.EXPECT().GetOS().Return(tt.want.OS, nil)
			mCfg.EXPECT().GetOSVersion().Return(tt.want.OSVersion, nil)
			mCfg.EXPECT().GetArch().Return(tt.want.Arch, nil)

			uc, err := di.InitializeTestUsecaseSet(nil, nil, mCfg, nil, nil, nil, nil)
			if err != nil {
				t.Fatal(err)
			}

			c := uc.ConfigUsecase

			got, err := c.ScanEnvInfo()
			if (err != nil) != tt.wantErr {
				t.Errorf("ConfigUsecaseImpl.ScanEnvInfo() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConfigUsecaseImpl.ScanEnvInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}
