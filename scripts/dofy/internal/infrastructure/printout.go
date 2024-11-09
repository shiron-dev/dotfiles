package infrastructure

import (
	"fmt"
	"io"
	"log"
	"os"
)

type PrintOutInfrastructure interface {
	Print(str string)
	SetLogOutput() *os.File

	GetOut() *io.Writer
	GetError() *io.Writer
}

type PrintOutInfrastructureImpl struct {
	out   io.Writer
	error io.Writer
}

func NewPrintOutInfrastructure() *PrintOutInfrastructureImpl {
	return &PrintOutInfrastructureImpl{
		out:   os.Stdout,
		error: os.Stderr,
	}
}

func (p *PrintOutInfrastructureImpl) Print(str string) {
	log.Print(str)
	fmt.Fprint(os.Stdout, str)
}

func (p *PrintOutInfrastructureImpl) SetLogOutput() *os.File {
	logfile, err := os.OpenFile("./dotfiles.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, filePermission)
	if err != nil {
		panic("cannot open ./dotfiles.log:" + err.Error())
	}

	log.SetOutput(logfile)

	log.SetFlags(log.Ldate | log.Ltime)

	p.out = io.MultiWriter(os.Stdout, logfile)
	p.error = io.MultiWriter(os.Stderr, logfile)

	return logfile
}

func (p *PrintOutInfrastructureImpl) GetOut() *io.Writer {
	return &p.out
}

func (p *PrintOutInfrastructureImpl) GetError() *io.Writer {
	return &p.error
}
