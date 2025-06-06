package usecase_test

import (
	"bytes"
	"context"
	"io"
	"testing"

	mock_infrastructure "github.com/shiron-dev/dotfiles/scripts/dofy/gen/mock/infrastructure"
	"github.com/shiron-dev/dotfiles/scripts/dofy/internal/di"
	"github.com/shiron-dev/dotfiles/scripts/dofy/internal/domain"
	"go.uber.org/mock/gomock"
)

func TestDepsUsecaseImpl_CheckInstalled(t *testing.T) {
	t.Parallel()

	type args struct {
		name string
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{"test", args{"test"}, true},
		{"test2", args{"test2"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mDeps := mock_infrastructure.NewMockDepsInfrastructure(ctrl)
			mCfg := mock_infrastructure.NewMockConfigInfrastructure(ctrl)
			mBrew := mock_infrastructure.NewMockBrewInfrastructure(ctrl)
			mFile := mock_infrastructure.NewMockFileInfrastructure(ctrl)
			mGit := mock_infrastructure.NewMockGitInfrastructure(ctrl)
			mPrintOut := mock_infrastructure.NewMockPrintOutInfrastructure(ctrl)

			mDeps.EXPECT().CheckInstalled(tt.args.name).Return(tt.want)

			uc, err := di.InitializeTestUsecaseSet(nil, mBrew, mCfg, mDeps, mFile, mGit, mPrintOut)
			if err != nil {
				t.Fatal(err)
			}

			d := uc.DepsUsecase

			if got := d.CheckInstalled(tt.args.name); got != tt.want {
				t.Errorf("DepsUsecaseImpl.CheckInstalled() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDepsUsecaseImpl_InstallHomebrew(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx context.Context
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"no error", args{t.Context()}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mDeps := mock_infrastructure.NewMockDepsInfrastructure(ctrl)
			mCfg := mock_infrastructure.NewMockConfigInfrastructure(ctrl)
			mBrew := mock_infrastructure.NewMockBrewInfrastructure(ctrl)
			mFile := mock_infrastructure.NewMockFileInfrastructure(ctrl)
			mGit := mock_infrastructure.NewMockGitInfrastructure(ctrl)
			mPrintOut := mock_infrastructure.NewMockPrintOutInfrastructure(ctrl)

			sout := io.Writer(&bytes.Buffer{})
			serror := io.Writer(&bytes.Buffer{})

			mPrintOut.EXPECT().Print(gomock.Any()).AnyTimes()
			mPrintOut.EXPECT().GetOut().Return(&sout)
			mPrintOut.EXPECT().GetError().Return(&serror)
			mDeps.EXPECT().CheckInstalled(gomock.Eq("brew")).Return(false)
			mCfg.EXPECT().GetOS().Return("testOS", nil)
			mCfg.EXPECT().GetOSVersion().Return("testOSVersion", nil)
			mCfg.EXPECT().GetArch().Return("testArch", nil)
			mBrew.EXPECT().InstallHomebrew(gomock.Eq(tt.args.ctx), gomock.Eq(sout), gomock.Eq(serror)).Return(nil)
			mBrew.EXPECT().SetHomebrewEnv(gomock.Eq("testOS")).Return(nil)

			uc, err := di.InitializeTestUsecaseSet(nil, mBrew, mCfg, mDeps, mFile, mGit, mPrintOut)
			if err != nil {
				t.Fatal(err)
			}

			d := uc.DepsUsecase

			if err := d.InstallHomebrew(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("DepsUsecaseImpl.InstallHomebrew() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDepsUsecaseImpl_InstallGit(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{"no error", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mDeps := mock_infrastructure.NewMockDepsInfrastructure(ctrl)
			mCfg := mock_infrastructure.NewMockConfigInfrastructure(ctrl)
			mBrew := mock_infrastructure.NewMockBrewInfrastructure(ctrl)
			mFile := mock_infrastructure.NewMockFileInfrastructure(ctrl)
			mGit := mock_infrastructure.NewMockGitInfrastructure(ctrl)
			mPrintOut := mock_infrastructure.NewMockPrintOutInfrastructure(ctrl)

			mDeps.EXPECT().CheckInstalled(gomock.Eq("git")).Return(false)
			mPrintOut.EXPECT().Print(gomock.Any()).AnyTimes()

			uc, err := di.InitializeTestUsecaseSet(nil, mBrew, mCfg, mDeps, mFile, mGit, mPrintOut)
			if err != nil {
				t.Fatal(err)
			}

			d := uc.DepsUsecase

			if err := d.InstallGit(); (err != nil) != tt.wantErr {
				t.Errorf("DepsUsecaseImpl.InstallGit() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDepsUsecaseImpl_CloneDotfiles(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{"no error", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mDeps := mock_infrastructure.NewMockDepsInfrastructure(ctrl)
			mCfg := mock_infrastructure.NewMockConfigInfrastructure(ctrl)
			mBrew := mock_infrastructure.NewMockBrewInfrastructure(ctrl)
			mFile := mock_infrastructure.NewMockFileInfrastructure(ctrl)
			mGit := mock_infrastructure.NewMockGitInfrastructure(ctrl)
			mPrintOut := mock_infrastructure.NewMockPrintOutInfrastructure(ctrl)

			uc, err := di.InitializeTestUsecaseSet(nil, mBrew, mCfg, mDeps, mFile, mGit, mPrintOut)
			if err != nil {
				t.Fatal(err)
			}

			mPrintOut.EXPECT().Print(gomock.Any()).AnyTimes()

			d := uc.DepsUsecase

			if err := d.CloneDotfiles(); (err != nil) != tt.wantErr {
				t.Errorf("DepsUsecaseImpl.CloneDotfiles() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDepsUsecaseImpl_InstallBrewBundle(t *testing.T) {
	t.Parallel()

	type args struct {
		forceInstall bool
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"no error", args{true}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mDeps := mock_infrastructure.NewMockDepsInfrastructure(ctrl)
			mCfg := mock_infrastructure.NewMockConfigInfrastructure(ctrl)
			mBrew := mock_infrastructure.NewMockBrewInfrastructure(ctrl)
			mFile := mock_infrastructure.NewMockFileInfrastructure(ctrl)
			mGit := mock_infrastructure.NewMockGitInfrastructure(ctrl)
			mPrintOut := mock_infrastructure.NewMockPrintOutInfrastructure(ctrl)

			sout := io.Writer(&bytes.Buffer{})
			serror := io.Writer(&bytes.Buffer{})

			mPrintOut.EXPECT().Print(gomock.Any()).AnyTimes()
			mPrintOut.EXPECT().GetOut().Return(&sout).AnyTimes()
			mPrintOut.EXPECT().GetError().Return(&serror).AnyTimes()
			mCfg.EXPECT().GetOS().Return("testOS", nil)
			mCfg.EXPECT().GetOSVersion().Return("testOSVersion", nil)
			mCfg.EXPECT().GetArch().Return("testArch", nil)
			mBrew.EXPECT().InstallTap(gomock.Any(), sout, serror).Return(nil)
			mBrew.EXPECT().DumpTmpBrewBundle(gomock.Any(), false, sout, serror).Return(nil)
			mBrew.EXPECT().InstallBrewBundle(gomock.Any(), sout, serror).Return(nil)
			mBrew.EXPECT().ReadBrewBundle(gomock.Any()).Return([]domain.BrewBundle{
				{
					Name: "git", Others: []string{}, BundleType: domain.BrewBundleTypeFormula, Categories: []string{},
				},
			},
				nil).AnyTimes()
			mBrew.EXPECT().WriteBrewBundle(gomock.Any(), gomock.Eq([]domain.BrewBundle{
				{
					Name: "git", Others: []string{}, BundleType: domain.BrewBundleTypeFormula, Categories: []string{},
				},
			})).Return(nil).AnyTimes()

			uc, err := di.InitializeTestUsecaseSet(nil, mBrew, mCfg, mDeps, mFile, mGit, mPrintOut)
			if err != nil {
				t.Fatal(err)
			}

			d := uc.DepsUsecase

			if err := d.InstallBrewBundle(tt.args.forceInstall); (err != nil) != tt.wantErr {
				t.Errorf("DepsUsecaseImpl.InstallBrewBundle() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDepsUsecaseImpl_Finish(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{"no error", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mDeps := mock_infrastructure.NewMockDepsInfrastructure(ctrl)
			mCfg := mock_infrastructure.NewMockConfigInfrastructure(ctrl)
			mBrew := mock_infrastructure.NewMockBrewInfrastructure(ctrl)
			mFile := mock_infrastructure.NewMockFileInfrastructure(ctrl)
			mGit := mock_infrastructure.NewMockGitInfrastructure(ctrl)
			mPrintOut := mock_infrastructure.NewMockPrintOutInfrastructure(ctrl)

			mPrintOut.EXPECT().Print(gomock.Any()).AnyTimes()

			uc, err := di.InitializeTestUsecaseSet(nil, mBrew, mCfg, mDeps, mFile, mGit, mPrintOut)
			if err != nil {
				t.Fatal(err)
			}

			d := uc.DepsUsecase

			if err := d.Finish(); (err != nil) != tt.wantErr {
				t.Errorf("DepsUsecaseImpl.Finish() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
