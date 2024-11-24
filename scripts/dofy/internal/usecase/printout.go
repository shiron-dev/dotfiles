package usecase

import (
	"fmt"
	"io"
	"os"
	"reflect"

	"github.com/shiron-dev/dotfiles/scripts/dofy/internal/domain"
	"github.com/shiron-dev/dotfiles/scripts/dofy/internal/infrastructure"
)

type PrintOutUsecase interface {
	PrintMdf(format string, a ...interface{})
	Println(str string)
	Print(str string)
	PrintObj(obj interface{})
	SetLogOutput() *os.File

	GetOut() *io.Writer
	GetError() *io.Writer
}

type PrintOutUsecaseImpl struct {
	printOutInfrastructure infrastructure.PrintOutInfrastructure
}

func NewPrintOutUsecase(printOutInfrastructure infrastructure.PrintOutInfrastructure) *PrintOutUsecaseImpl {
	return &PrintOutUsecaseImpl{
		printOutInfrastructure: printOutInfrastructure,
	}
}

func (p *PrintOutUsecaseImpl) PrintMdf(format string, a ...interface{}) {
	str := fmt.Sprintf(format, a...)

	for _, printer := range domain.GetMdPrinter() {
		if printer.Printer != nil {
			str = printer.Printer(str)
		} else {
			str = printer.Reg.ReplaceAllStringFunc(str, func(s string) string {
				return printer.Col.SprintFunc()(s)
			})
		}
	}

	p.Print(str)
}

func (p *PrintOutUsecaseImpl) PrintObj(obj interface{}) {
	refT := reflect.TypeOf(obj)
	refV := reflect.ValueOf(obj)

	str := ""

	for i := range refT.NumField() {
		field := refT.Field(i)
		value := refV.Field(i)

		str += field.Name + ": " + fmt.Sprintf("%v", value) + "\n"
	}

	p.Print(str)
}

func (p *PrintOutUsecaseImpl) Println(str string) {
	p.printOutInfrastructure.Print(str + "\n")
}

func (p *PrintOutUsecaseImpl) Print(str string) {
	p.printOutInfrastructure.Print(str)
}

const filePermission = 0o666

func (p *PrintOutUsecaseImpl) SetLogOutput() *os.File {
	logFile, err := os.OpenFile("./dotfiles.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, filePermission)
	if err != nil {
		panic("cannot open ./dotfiles.log:" + err.Error())
	}

	p.printOutInfrastructure.SetLogOutput(logFile)

	return logFile
}

func (p *PrintOutUsecaseImpl) GetOut() *io.Writer {
	return p.printOutInfrastructure.GetOut()
}

func (p *PrintOutUsecaseImpl) GetError() *io.Writer {
	return p.printOutInfrastructure.GetError()
}
