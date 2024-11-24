package usecase_test

import (
	"bytes"
	"io"
	"reflect"
	"testing"

	"github.com/fatih/color"
	mock_infrastructure "github.com/shiron-dev/dotfiles/scripts/dofy/gen/mock/infrastructure"
	"github.com/shiron-dev/dotfiles/scripts/dofy/internal/di"
	"github.com/shiron-dev/dotfiles/scripts/dofy/internal/domain"
	"go.uber.org/mock/gomock"
)

func TestPrintOutUsecaseImpl_PrintMdf(t *testing.T) {
	t.Parallel()

	type args struct {
		format string
		a      []interface{}
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{"print", args{"test PrintMdf", nil}, "test PrintMdf\n"},
		{"print h1", args{"# test PrintMdf", nil}, func() string {
			for _, printer := range domain.GetMdPrinter() {
				if printer.Name != "h1" {
					continue
				}

				return printer.Col.SprintFunc()("# test PrintMdf") + "\n"
			}

			return ""
		}()},
		{"print underline", args{"test __PrintMdf__", nil}, func() string {
			for _, printer := range domain.GetMdPrinter() {
				if printer.Name != "underline" {
					continue
				}

				return "test " + color.New(color.Underline).Sprint("PrintMdf") + "\n"
			}

			return ""
		}()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mPrintOut := mock_infrastructure.NewMockPrintOutInfrastructure(ctrl)

			mPrintOut.EXPECT().Print(gomock.Eq(tt.want)).Return()

			uc, err := di.InitializeTestUsecaseSet(nil, nil, nil, nil, nil, nil, mPrintOut)
			if err != nil {
				t.Fatal(err)
			}

			p := uc.PrintOutUsecase

			p.PrintMdf(tt.args.format, tt.args.a...)
		})
	}
}

func TestPrintOutUsecaseImpl_PrintObj(t *testing.T) {
	t.Parallel()

	type args struct {
		obj interface{}
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"print",
			args{struct {
				a string
				b int
				c bool
				d bool
			}{"test", 1, true, false}},
			"a: test\nb: 1\nc: true\nd: false\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mPrintOut := mock_infrastructure.NewMockPrintOutInfrastructure(ctrl)

			uc, err := di.InitializeTestUsecaseSet(nil, nil, nil, nil, nil, nil, mPrintOut)
			if err != nil {
				t.Fatal(err)
			}

			mPrintOut.EXPECT().Print(gomock.Eq(tt.want)).Return()

			p := uc.PrintOutUsecase

			p.PrintObj(tt.args.obj)
		})
	}
}

func TestPrintOutUsecaseImpl_Println(t *testing.T) {
	t.Parallel()

	type args struct {
		str string
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{"print", args{"test Print"}, "test Print\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mPrintOut := mock_infrastructure.NewMockPrintOutInfrastructure(ctrl)

			uc, err := di.InitializeTestUsecaseSet(nil, nil, nil, nil, nil, nil, mPrintOut)
			if err != nil {
				t.Fatal(err)
			}

			mPrintOut.EXPECT().Print(gomock.Eq(tt.want)).Return()

			p := uc.PrintOutUsecase

			p.Println(tt.args.str)
		})
	}
}

func TestPrintOutUsecaseImpl_Print(t *testing.T) {
	t.Parallel()

	type args struct {
		str string
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{"print", args{"test Print"}, "test Print"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mPrintOut := mock_infrastructure.NewMockPrintOutInfrastructure(ctrl)

			uc, err := di.InitializeTestUsecaseSet(nil, nil, nil, nil, nil, nil, mPrintOut)
			if err != nil {
				t.Fatal(err)
			}

			mPrintOut.EXPECT().Print(gomock.Eq(tt.want)).Return()

			p := uc.PrintOutUsecase

			p.Print(tt.args.str)
		})
	}
}

//nolint:paralleltest
func TestPrintOutUsecaseImpl_SetLogOutput(t *testing.T) {
	t.Skip("skipping test; not running")
}

func TestPrintOutUsecaseImpl_GetOut(t *testing.T) {
	t.Parallel()

	sout := &bytes.Buffer{}
	soutWriter := io.Writer(sout)

	tests := []struct {
		name string
		want *io.Writer
	}{
		{"test1", &soutWriter},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mPrintOut := mock_infrastructure.NewMockPrintOutInfrastructure(ctrl)

			uc, err := di.InitializeTestUsecaseSet(nil, nil, nil, nil, nil, nil, mPrintOut)
			if err != nil {
				t.Fatal(err)
			}

			mPrintOut.EXPECT().GetOut().Return(&soutWriter)

			p := uc.PrintOutUsecase

			if got := p.GetOut(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PrintOutUsecaseImpl.GetOut() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrintOutUsecaseImpl_GetError(t *testing.T) {
	t.Parallel()

	serror := &bytes.Buffer{}
	serrorWriter := io.Writer(serror)

	tests := []struct {
		name string
		want *io.Writer
	}{
		{"test1", &serrorWriter},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mPrintOut := mock_infrastructure.NewMockPrintOutInfrastructure(ctrl)

			uc, err := di.InitializeTestUsecaseSet(nil, nil, nil, nil, nil, nil, mPrintOut)
			if err != nil {
				t.Fatal(err)
			}

			mPrintOut.EXPECT().GetError().Return(&serrorWriter)

			p := uc.PrintOutUsecase

			if got := p.GetError(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PrintOutUsecaseImpl.GetError() = %v, want %v", got, tt.want)
			}
		})
	}
}
