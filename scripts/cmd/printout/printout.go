package printout

import (
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"regexp"

	"github.com/fatih/color"
)

type mdPrinter struct {
	name string
	reg  *regexp.Regexp
	col  *color.Color
}

var printer = []mdPrinter{
	{
		"h1",
		regexp.MustCompile(`(?m)^# (.+)`),
		color.New(color.FgCyan),
	},
	{
		"h2",
		regexp.MustCompile(`(?m)^## (.*)`),
		color.New(color.FgCyan),
	},
	{
		"h3",
		regexp.MustCompile(`(?m)^### (.*)`),
		color.New(color.FgCyan),
	},
	{
		"h4",
		regexp.MustCompile(`(?m)^#### (.*)`),
		color.New(color.FgCyan),
	},
	{
		"h5",
		regexp.MustCompile(`(?m)^##### (.*)`),
		color.New(color.FgCyan),
	},
	{
		"h6",
		regexp.MustCompile(`(?m)^###### (.*)`),
		color.New(color.FgCyan),
	},
	{
		"bold",
		regexp.MustCompile(`\*\*(.*)\*\*`),
		color.New(color.Bold),
	},
	{
		"italic",
		regexp.MustCompile(`\*(.*)\*`),
		color.New(color.Italic),
	},
	{
		"code",
		regexp.MustCompile("`(.*)`"),
		color.New(color.FgHiWhite),
	},
	{
		"underline",
		regexp.MustCompile(`__(.*)__`),
		color.New(color.Underline),
	},
}

func PrintMd(format string, a ...interface{}) {
	str := fmt.Sprintf(format, a...)

	for _, p := range printer {
		if p.name == "underline" {
			str = p.reg.ReplaceAllStringFunc(str, func(s string) string {
				return p.reg.ReplaceAllString(s, p.col.Sprint("$1"))
			})
		} else {
			str = p.reg.ReplaceAllStringFunc(str, func(s string) string {
				return p.col.SprintFunc()(s)
			})
		}
	}

	Println(str)
}

func Println(str string) {
	log.Println(str)
	fmt.Println(str)
}

func Print(str string) {
	log.Print(str)
	fmt.Print(str)
}

func PrintObj(obj interface{}) {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		Println(field.Name + ": " + v.Field(i).String())
	}
}

func SetLogOutput() *os.File {
	logfile, err := os.OpenFile("./dotfiles.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic("cannnot open ./dotfiles.log:" + err.Error())
	}
	log.SetOutput(logfile)

	log.SetFlags(log.Ldate | log.Ltime)

	Out = io.MultiWriter(os.Stdout, logfile)
	return logfile
}

var Out io.Writer
