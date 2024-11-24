package usecase_test

import (
	"bytes"
	"io"
	"testing"

	mock_infrastructure "github.com/shiron-dev/dotfiles/scripts/dofy/gen/mock/infrastructure"
	"github.com/shiron-dev/dotfiles/scripts/dofy/internal/di"
	"go.uber.org/mock/gomock"
)

func TestAnsibleUsecaseImpl_CheckPlaybook(t *testing.T) {
	t.Parallel()

	type args struct {
		invPath      string
		playbookPath string
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"normal", args{"test", "test"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mAnsible := mock_infrastructure.NewMockAnsibleInfrastructure(ctrl)
			mPrintOut := mock_infrastructure.NewMockPrintOutInfrastructure(ctrl)

			sout := io.Writer(&bytes.Buffer{})
			serror := io.Writer(&bytes.Buffer{})

			mPrintOut.EXPECT().Print(gomock.Any()).AnyTimes()
			mPrintOut.EXPECT().GetOut().Return(&sout)
			mPrintOut.EXPECT().GetError().Return(&serror)
			mAnsible.EXPECT().CheckPlaybook(
				gomock.Eq(tt.args.invPath),
				gomock.Eq(tt.args.playbookPath),
				gomock.Eq(sout),
				gomock.Eq(serror),
			).Return(nil)

			uc, err := di.InitializeTestUsecaseSet(mAnsible, nil, nil, nil, nil, nil, mPrintOut)
			if err != nil {
				t.Fatal(err)
			}

			a := uc.AnsibleUsecase

			if err := a.CheckPlaybook(tt.args.invPath, tt.args.playbookPath); (err != nil) != tt.wantErr {
				t.Errorf("AnsibleUsecaseImpl.CheckPlaybook() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAnsibleUsecaseImpl_RunPlaybook(t *testing.T) {
	t.Parallel()

	type args struct {
		invPath      string
		playbookPath string
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"normal", args{"test", "test"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mAnsible := mock_infrastructure.NewMockAnsibleInfrastructure(ctrl)
			mPrintOut := mock_infrastructure.NewMockPrintOutInfrastructure(ctrl)

			sout := io.Writer(&bytes.Buffer{})
			serror := io.Writer(&bytes.Buffer{})

			mPrintOut.EXPECT().Print(gomock.Any()).AnyTimes()
			mPrintOut.EXPECT().GetOut().Return(&sout)
			mPrintOut.EXPECT().GetError().Return(&serror)
			mAnsible.EXPECT().RunPlaybook(
				gomock.Eq(tt.args.invPath),
				gomock.Eq(tt.args.playbookPath),
				gomock.Eq(sout),
				gomock.Eq(serror),
			).Return(nil)

			uc, err := di.InitializeTestUsecaseSet(mAnsible, nil, nil, nil, nil, nil, mPrintOut)
			if err != nil {
				t.Fatal(err)
			}

			a := uc.AnsibleUsecase

			if err := a.RunPlaybook(tt.args.invPath, tt.args.playbookPath); (err != nil) != tt.wantErr {
				t.Errorf("AnsibleUsecaseImpl.RunPlaybook() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
