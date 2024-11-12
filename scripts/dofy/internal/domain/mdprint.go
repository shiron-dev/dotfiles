package domain

import (
	"regexp"
	"strings"

	"github.com/fatih/color"
)

type MdPrinterType struct {
	Name    string
	Reg     *regexp.Regexp
	Col     *color.Color
	Printer func(string) string
}

func GetMdPrinter() []MdPrinterType {
	return []MdPrinterType{
		{
			"h1",
			regexp.MustCompile(`(?m)^# (.+)`),
			color.New(color.FgCyan),
			nil,
		},
		{
			"h2",
			regexp.MustCompile(`(?m)^## (.*)`),
			color.New(color.FgCyan),
			nil,
		},
		{
			"h3",
			regexp.MustCompile(`(?m)^### (.*)`),
			color.New(color.FgCyan),
			nil,
		},
		{
			"h4",
			regexp.MustCompile(`(?m)^#### (.*)`),
			color.New(color.FgCyan),
			nil,
		},
		{
			"h5",
			regexp.MustCompile(`(?m)^##### (.*)`),
			color.New(color.FgCyan),
			nil,
		},
		{
			"h6",
			regexp.MustCompile(`(?m)^###### (.*)`),
			color.New(color.FgCyan),
			nil,
		},
		{
			"bold",
			regexp.MustCompile(`\*\*(.*)\*\*`),
			color.New(color.Bold),
			nil,
		},
		{
			"italic",
			regexp.MustCompile(`\*(.*)\*`),
			color.New(color.Italic),
			nil,
		},
		{
			"code",
			regexp.MustCompile("`(.*)`"),
			color.New(color.FgHiWhite),
			nil,
		},
		{
			"underline",
			nil,
			nil,
			func(str string) string {
				re := regexp.MustCompile(`__(.*)__`)

				return re.ReplaceAllStringFunc(str, func(s string) string {
					return re.ReplaceAllString(s, color.New(color.Underline).Sprint("$1"))
				})
			},
		},
		{
			"alert",
			nil,
			nil,
			func(str string) string {
				ret := ""

				re := regexp.MustCompile(`^> \[!(.+?)\]$`)

				var fg *color.Color = nil
				for _, line := range strings.Split(str, "\n") {
					if re.MatchString(line) {
						emoji := ""
						switch re.FindStringSubmatch(line)[1] {
						case "NOTE":
							fg = color.New(color.FgBlue)
							emoji = "📝"
						case "TIP":
							fg = color.New(color.FgGreen)
							emoji = "💡"
						case "IMPORTANT":
							fg = color.New(color.FgMagenta)
							emoji = "❗"
						case "WARNING":
							//nolint:mnd
							fg = color.RGB(255, 128, 0)
							emoji = "⚠️"
						case "CAUTION":
							fg = color.New(color.FgRed)
							emoji = "🚨"
						}
						ret += fg.SprintFunc()("|"+strings.ReplaceAll(line[1:], "!", emoji)) + "\n"

						continue
					}

					if !strings.HasPrefix(line, ">") {
						fg = nil
					}

					if fg != nil {
						ret += fg.SprintFunc()("|") + line[1:] + "\n"
					} else {
						ret += line + "\n"
					}
				}

				return ret
			},
		},
	}
}
