package usecase

import (
	"dofy/internal/domain"
	"dofy/internal/infrastructure"
	"fmt"
	"os"
	"reflect"
)

type PrintOutUsecase interface {
	PrintMd(format string, a ...interface{})
	Println(str string)
	Print(str string)
	PrintObj(obj interface{})
	SetLogOutput() *os.File
}

type PrintOutUsecaseImpl struct {
	printOutInfrastructure infrastructure.PrintOutInfrastructure
}

func NewPrintOutUsecase(printOutInfrastructure infrastructure.PrintOutInfrastructure) PrintOutUsecase {
	return &PrintOutUsecaseImpl{
		printOutInfrastructure: printOutInfrastructure,
	}
}

func (p *PrintOutUsecaseImpl) PrintMd(format string, a ...interface{}) {
	str := fmt.Sprintf(format, a...)

	for _, p := range domain.MdPrinter {
		if p.Name == "underline" {
			str = p.Reg.ReplaceAllStringFunc(str, func(s string) string {
				return p.Reg.ReplaceAllString(s, p.Col.Sprint("$1"))
			})
		} else {
			str = p.Reg.ReplaceAllStringFunc(str, func(s string) string {
				return p.Col.SprintFunc()(s)
			})
		}
	}

	p.Println(str)
}

func (p *PrintOutUsecaseImpl) PrintObj(obj interface{}) {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		p.Println(field.Name + ": " + v.Field(i).String())
	}
}

func (p *PrintOutUsecaseImpl) Println(str string) {
	p.printOutInfrastructure.Print(str + "\n")
}

func (p *PrintOutUsecaseImpl) Print(str string) {
	p.printOutInfrastructure.Print(str)
}

func (p *PrintOutUsecaseImpl) SetLogOutput() *os.File {
	return p.printOutInfrastructure.SetLogOutput()
}
