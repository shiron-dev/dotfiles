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
}

type PrintOutInfrastructureImpl struct {
	out   io.Writer
	error io.Writer
}

func NewPrintOutInfrastructure() PrintOutInfrastructure {
	return &PrintOutInfrastructureImpl{}
}

func (p *PrintOutInfrastructureImpl) Print(str string) {
	log.Print(str)
	fmt.Print(str)
}

func (p *PrintOutInfrastructureImpl) SetLogOutput() *os.File {
	logfile, err := os.OpenFile("./dotfiles.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic("cannnot open ./dotfiles.log:" + err.Error())
	}
	log.SetOutput(logfile)

	log.SetFlags(log.Ldate | log.Ltime)

	p.out = io.MultiWriter(os.Stdout, logfile)
	p.error = io.MultiWriter(os.Stderr, logfile)
	return logfile
}
