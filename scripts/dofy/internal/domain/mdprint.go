package domain

import (
	"regexp"

	"github.com/fatih/color"
)

type MdPrinterType struct {
	Name string
	Reg  *regexp.Regexp
	Col  *color.Color
}

func GetMdPrinter() []MdPrinterType {
	return []MdPrinterType{
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
}
