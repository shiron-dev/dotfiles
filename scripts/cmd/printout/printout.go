package printout

import (
	"fmt"
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

	fmt.Println(str)
}

func Println(str string) {
	fmt.Println(str)
}

func PrintObj(obj interface{}) {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		Println(field.Name + ": " + v.Field(i).String())
	}
}
