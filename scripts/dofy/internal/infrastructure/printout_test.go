package infrastructure_test

import (
	"bufio"
	"bytes"
	"io"
	"os"
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

func TestGetOut(t *testing.T) {
	t.Parallel()

	outBuffer := &bytes.Buffer{}

	errBuffer := &bytes.Buffer{}

	infra, err := di.InitializeTestInfrastructureSet(outBuffer, errBuffer)
	if err != nil {
		t.Fatal(err)
	}

	if infra.PrintOutInfrastructure.GetOut() == nil {
		t.Fatal("GetOut() is nil")
	}

	if *infra.PrintOutInfrastructure.GetOut() != io.Writer(outBuffer) {
		t.Errorf("GetOut() = %v, want %v", infra.PrintOutInfrastructure.GetOut(), outBuffer)
	}
}

func TestGetError(t *testing.T) {
	t.Parallel()

	outBuffer := &bytes.Buffer{}

	errBuffer := &bytes.Buffer{}

	infra, err := di.InitializeTestInfrastructureSet(outBuffer, errBuffer)
	if err != nil {
		t.Fatal(err)
	}

	if infra.PrintOutInfrastructure.GetError() == nil {
		t.Fatal("GetError() is nil")
	}

	if *infra.PrintOutInfrastructure.GetError() != io.Writer(errBuffer) {
		t.Errorf("GetError() = %v, want %v", infra.PrintOutInfrastructure.GetError(), errBuffer)
	}
}
