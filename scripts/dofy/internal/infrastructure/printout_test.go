package infrastructure_test

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/shiron-dev/dotfiles/scripts/dofy/internal/di"
)

const filePermission = 0o666

//nolint:funlen,cyclop
func TestPrint(t *testing.T) {
	t.Parallel()

	outBuffer := &bytes.Buffer{}

	errBuffer := &bytes.Buffer{}

	infra, err := di.InitializeTestInfrastructureSet(outBuffer, errBuffer)
	if err != nil {
		t.Fatal(err)
	}

	testStrs := []string{
		"test1",
		"test2",
		"test3\nt3",
	}

	logFile, err := os.OpenFile("./dotfiles.log", os.O_CREATE|os.O_RDWR, filePermission)
	if err != nil {
		panic("cannot open ./dotfiles.log:" + err.Error())
	}
	defer logFile.Close()

	if err := logFile.Truncate(0); err != nil {
		t.Fatal(err)
	}

	if _, err := logFile.Seek(0, 0); err != nil {
		t.Fatal(err)
	}

	infra.PrintOutInfrastructure.SetLogOutput(logFile)

	if logFile == nil {
		t.Fatal("logFile is nil")
	}

	logLastLine := 0

	for _, testStr := range testStrs {
		outBuffer.Reset()
		infra.PrintOutInfrastructure.Print(testStr)

		if outBuffer.String() != testStr {
			t.Errorf("Print() = %q, want %q", outBuffer.String(), testStr)
		}

		if _, err := logFile.Seek(0, 0); err != nil {
			t.Fatal(err)
		}

		scanner := bufio.NewScanner(logFile)
		str := ""

		lineCount := 0

		for scanner.Scan() {
			if str != "" {
				str += "\n"
			}

			str += scanner.Text()
			lineCount++

			if lineCount <= logLastLine {
				str = ""
			}
		}

		logLastLine = lineCount

		if !strings.HasSuffix(str, testStr) {
			t.Errorf("log file = %q, want %q", str, testStr)
		}
	}
}

func TestPrintOutInfrastructureImpl_Print(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name       string
		wantSout   string
		wantSerror string
		args       args
	}{
		{"test1", "test1", "", args{"test1"}},
		{"test2", "test2\nt2", "", args{"test2\nt2"}},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			sout := &bytes.Buffer{}
			serror := &bytes.Buffer{}
			infra, err := di.InitializeTestInfrastructureSet(sout, serror)
			if err != nil {
				t.Fatal(err)
			}

			p := infra.PrintOutInfrastructure

			logPath := filepath.Join(t.TempDir(), tt.name+".log")

			logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_RDWR, filePermission)
			if err != nil {
				t.Fatal(err)
			}
			defer logFile.Close()

			p.SetLogOutput(logFile)

			p.Print(tt.args.str)

			if gotSout := sout.String(); gotSout != tt.wantSout {
				t.Errorf("GitInfrastructureImpl.GitDifftool() = %v, want %v", gotSout, tt.wantSout)
			}
			if gotSerror := serror.String(); gotSerror != tt.wantSerror {
				t.Errorf("GitInfrastructureImpl.GitDifftool() = %v, want %v", gotSerror, tt.wantSerror)
			}

			if _, err := logFile.Seek(0, 0); err != nil {
				t.Fatal(err)
			}

			str := ""
			scanner := bufio.NewScanner(logFile)

			for scanner.Scan() {
				if str != "" {
					str += "\n"
				}

				str += scanner.Text()
			}

			if !strings.HasSuffix(str, tt.wantSout) {
				t.Errorf("log file = %q, want %q", str, tt.wantSout)
			}
		})
	}
}

func TestPrintOutInfrastructureImpl_SetLogOutput(t *testing.T) {
	file, err := os.Create(filepath.Join(t.TempDir(), "test.log"))
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		logFile *os.File
	}
	tests := []struct {
		name string
		args args
	}{
		{"test1", args{file}},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			sout := &bytes.Buffer{}
			serror := &bytes.Buffer{}
			infra, err := di.InitializeTestInfrastructureSet(sout, serror)
			if err != nil {
				t.Fatal(err)
			}

			p := infra.PrintOutInfrastructure

			p.SetLogOutput(tt.args.logFile)
		})
	}
}

func TestPrintOutInfrastructureImpl_GetOut(t *testing.T) {
	sout := &bytes.Buffer{}
	serror := &bytes.Buffer{}

	soutWriter := io.Writer(sout)

	tests := []struct {
		name string
		want *io.Writer
	}{
		{"test1", &soutWriter},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			infra, err := di.InitializeTestInfrastructureSet(sout, serror)
			if err != nil {
				t.Fatal(err)
			}

			p := infra.PrintOutInfrastructure

			if got := p.GetOut(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PrintOutInfrastructureImpl.GetOut() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrintOutInfrastructureImpl_GetError(t *testing.T) {
	sout := &bytes.Buffer{}
	serror := &bytes.Buffer{}

	sserrorWriter := io.Writer(serror)

	tests := []struct {
		name string
		want *io.Writer
	}{
		{"test1", &sserrorWriter},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			infra, err := di.InitializeTestInfrastructureSet(sout, serror)
			if err != nil {
				t.Fatal(err)
			}

			p := infra.PrintOutInfrastructure

			if got := p.GetError(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PrintOutInfrastructureImpl.GetError() = %v, want %v", got, tt.want)
			}
		})
	}
}
