package usecase_test

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	mock_infrastructure "github.com/shiron-dev/dotfiles/scripts/dofy/gen/mock/infrastructure"
	"github.com/shiron-dev/dotfiles/scripts/dofy/internal/di"
	"github.com/shiron-dev/dotfiles/scripts/dofy/internal/domain"
	"go.uber.org/mock/gomock"
)

func TestBrewUsecaseImpl_InstallHomebrew(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx context.Context
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"normal", args{t.Context()}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mBrew := mock_infrastructure.NewMockBrewInfrastructure(ctrl)
			mDeps := mock_infrastructure.NewMockDepsInfrastructure(ctrl)
			mPrintOut := mock_infrastructure.NewMockPrintOutInfrastructure(ctrl)
			mCfg := mock_infrastructure.NewMockConfigInfrastructure(ctrl)

			sout := io.Writer(&bytes.Buffer{})
			serror := io.Writer(&bytes.Buffer{})

			mPrintOut.EXPECT().Print(gomock.Any()).AnyTimes()
			mPrintOut.EXPECT().GetOut().Return(&sout)
			mPrintOut.EXPECT().GetError().Return(&serror)
			mCfg.EXPECT().GetOS().Return("testOS", nil)
			mCfg.EXPECT().GetOSVersion().Return("testOSVersion", nil)
			mCfg.EXPECT().GetArch().Return("testArch", nil)
			mBrew.EXPECT().InstallHomebrew(gomock.Eq(tt.args.ctx), gomock.Eq(sout), gomock.Eq(serror)).Return(nil)
			mBrew.EXPECT().SetHomebrewEnv(gomock.Eq("testOS")).Return(nil)

			uc, err := di.InitializeTestUsecaseSet(nil, mBrew, mCfg, mDeps, nil, nil, mPrintOut)
			if err != nil {
				t.Fatal(err)
			}

			b := uc.BrewUsecase

			if err := b.InstallHomebrew(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("BrewUsecaseImpl.InstallHomebrew() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBrewUsecaseImpl_InstallFormula(t *testing.T) {
	t.Parallel()

	type args struct {
		formula string
		bType   domain.BrewBundleType
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"normal", args{"formula", domain.BrewBundleTypeFormula}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mBrew := mock_infrastructure.NewMockBrewInfrastructure(ctrl)
			mDeps := mock_infrastructure.NewMockDepsInfrastructure(ctrl)
			mPrintOut := mock_infrastructure.NewMockPrintOutInfrastructure(ctrl)
			mCfg := mock_infrastructure.NewMockConfigInfrastructure(ctrl)

			uc, err := di.InitializeTestUsecaseSet(nil, mBrew, mCfg, mDeps, nil, nil, mPrintOut)
			if err != nil {
				t.Fatal(err)
			}

			mPrintOut.EXPECT().Print(gomock.Any())

			b := uc.BrewUsecase

			if err := b.InstallFormula(tt.args.formula, tt.args.bType); (err != nil) != tt.wantErr {
				t.Errorf("BrewUsecaseImpl.InstallFormula() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBrewUsecaseImpl_InstallBrewBundle(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "/Brewfile")

	//nolint:gosec
	if file, err := os.Create(path); err != nil {
		t.Fatal(err)
	} else {
		_, err = file.WriteString("brew \"go\"\n")
		if err != nil {
			t.Fatal(err)
		}
	}

	type args struct {
		path string
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"no error", args{path}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mBrew := mock_infrastructure.NewMockBrewInfrastructure(ctrl)
			mDeps := mock_infrastructure.NewMockDepsInfrastructure(ctrl)
			mPrintOut := mock_infrastructure.NewMockPrintOutInfrastructure(ctrl)
			mCfg := mock_infrastructure.NewMockConfigInfrastructure(ctrl)

			sout := io.Writer(&bytes.Buffer{})
			serror := io.Writer(&bytes.Buffer{})

			mPrintOut.EXPECT().Print(gomock.Any()).AnyTimes()
			mPrintOut.EXPECT().GetOut().Return(&sout)
			mPrintOut.EXPECT().GetError().Return(&serror)
			mBrew.EXPECT().InstallBrewBundle(gomock.Eq(tt.args.path), gomock.Eq(sout), gomock.Eq(serror)).Return(nil)

			uc, err := di.InitializeTestUsecaseSet(nil, mBrew, mCfg, mDeps, nil, nil, mPrintOut)
			if err != nil {
				t.Fatal(err)
			}

			b := uc.BrewUsecase

			if err := b.InstallBrewBundle(tt.args.path); (err != nil) != tt.wantErr {
				t.Errorf("BrewUsecaseImpl.InstallBrewBundle() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBrewUsecaseImpl_DumpTmpBrewBundle(t *testing.T) {
	t.Parallel()

	type args struct {
		path string
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"no error", args{filepath.Join(t.TempDir(), "/Brewfile")}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mBrew := mock_infrastructure.NewMockBrewInfrastructure(ctrl)
			mDeps := mock_infrastructure.NewMockDepsInfrastructure(ctrl)
			mPrintOut := mock_infrastructure.NewMockPrintOutInfrastructure(ctrl)
			mCfg := mock_infrastructure.NewMockConfigInfrastructure(ctrl)

			mCfg.EXPECT().GetOS().Return("testOS", nil)
			mCfg.EXPECT().GetOSVersion().Return("testOSVersion", nil)
			mCfg.EXPECT().GetArch().Return("testArch", nil)

			sout := io.Writer(&bytes.Buffer{})
			serror := io.Writer(&bytes.Buffer{})

			mPrintOut.EXPECT().GetOut().Return(&sout)
			mPrintOut.EXPECT().GetError().Return(&serror)
			mBrew.EXPECT().DumpTmpBrewBundle(
				gomock.Eq(tt.args.path),
				gomock.Eq(false),
				gomock.Eq(sout),
				gomock.Eq(serror),
			).Return(nil)

			uc, err := di.InitializeTestUsecaseSet(nil, mBrew, mCfg, mDeps, nil, nil, mPrintOut)
			if err != nil {
				t.Fatal(err)
			}

			b := uc.BrewUsecase

			if err := b.DumpTmpBrewBundle(tt.args.path); (err != nil) != tt.wantErr {
				t.Errorf("BrewUsecaseImpl.DumpTmpBrewBundle() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBrewUsecaseImpl_CheckDiffBrewBundle(t *testing.T) {
	t.Parallel()

	type args struct {
		bundlePath string
		tmpPath    string
	}

	tests := []struct {
		name    string
		args    args
		want    []domain.BrewBundle
		want1   []domain.BrewBundle
		wantErr bool
	}{
		{
			"no error",
			args{filepath.Join(t.TempDir(), "/Brewfile"), filepath.Join(t.TempDir(), "/Brewfile")},
			[]domain.BrewBundle{{Name: "bundle", Others: nil, BundleType: domain.BrewBundleTypeFormula, Categories: nil}},
			[]domain.BrewBundle{{Name: "tmp", Others: nil, BundleType: domain.BrewBundleTypeTap, Categories: nil}},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mBrew := mock_infrastructure.NewMockBrewInfrastructure(ctrl)
			mDeps := mock_infrastructure.NewMockDepsInfrastructure(ctrl)
			mPrintOut := mock_infrastructure.NewMockPrintOutInfrastructure(ctrl)
			mCfg := mock_infrastructure.NewMockConfigInfrastructure(ctrl)

			mBrew.EXPECT().ReadBrewBundle(tt.args.bundlePath).Return(tt.want, nil)
			mBrew.EXPECT().ReadBrewBundle(tt.args.tmpPath).Return(tt.want1, nil)

			uc, err := di.InitializeTestUsecaseSet(nil, mBrew, mCfg, mDeps, nil, nil, mPrintOut)
			if err != nil {
				t.Fatal(err)
			}

			b := uc.BrewUsecase

			got, got1, err := b.CheckDiffBrewBundle(tt.args.bundlePath, tt.args.tmpPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("BrewUsecaseImpl.CheckDiffBrewBundle() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BrewUsecaseImpl.CheckDiffBrewBundle() got = %v, want %v", got, tt.want)
			}

			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("BrewUsecaseImpl.CheckDiffBrewBundle() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestBrewUsecaseImpl_CleanupBrewBundle(t *testing.T) {
	t.Parallel()

	t.Skip("skipping test; not running")
}
