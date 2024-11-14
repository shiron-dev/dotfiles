package infrastructure

import (
	"fmt"
	"io"
	"log"
	"os"
)

type PrintOutInfrastructure interface {
	Print(str string)
	SetLogOutput(logFile *os.File)

	GetOut() *io.Writer
	GetError() *io.Writer
}

type PrintOutInfrastructureImpl struct {
	out    io.Writer
	err    io.Writer
	stdout io.Writer
	stderr io.Writer
}

func NewPrintOutInfrastructure(stdout io.Writer, stderr io.Writer) *PrintOutInfrastructureImpl {
	return &PrintOutInfrastructureImpl{
		out:    stdout,
		err:    stderr,
		stdout: stdout,
		stderr: stderr,
	}
}

func (p *PrintOutInfrastructureImpl) Print(str string) {
	log.Print(str)
	fmt.Fprint(p.stdout, str)
}

func (p *PrintOutInfrastructureImpl) SetLogOutput(logFile *os.File) {
	log.SetOutput(logFile)

	log.SetFlags(log.Ldate | log.Ltime)

	p.out = io.MultiWriter(p.stdout, logFile)
	p.err = io.MultiWriter(p.stderr, logFile)
}

func (p *PrintOutInfrastructureImpl) GetOut() *io.Writer {
	return &p.out
}

func (p *PrintOutInfrastructureImpl) GetError() *io.Writer {
	return &p.err
}
