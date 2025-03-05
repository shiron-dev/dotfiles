package domain

import "strings"

type BrewBundleType uint

const (
	BrewBundleTypeTap BrewBundleType = iota
	BrewBundleTypeFormula
	BrewBundleTypeCask
	BrewBundleTypeMas
)

type BrewBundle struct {
	Name       string
	Others     []string
	BundleType BrewBundleType
	Categories []string
	Mode       string
}

func BrewBundleTypeFromString(str string) BrewBundleType {
	switch str {
	case "tap":
		return BrewBundleTypeTap
	case "brew":
		return BrewBundleTypeFormula
	case "cask":
		return BrewBundleTypeCask
	case "mas":
		return BrewBundleTypeMas
	default:
		return BrewBundleTypeFormula
	}
}

func (b BrewBundle) BundleFormat() string {
	str := ""

	switch b.BundleType {
	case BrewBundleTypeTap:
		str = "tap"
	case BrewBundleTypeFormula:
		str = "brew"
	case BrewBundleTypeCask:
		str = "cask"
	case BrewBundleTypeMas:
		str = "mas"
	}

	str += " \"" + b.Name + "\""
	if len(b.Others) > 0 {
		str += ", " + strings.Join(b.Others, ", ")
	}

	return str
}

func (b BrewBundle) String() string {
	str := "{ "

	str += b.BundleFormat()

	str += " { " + strings.Join(b.Categories, ", ") + " }"
	str += " :" + b.Mode

	str += " }"

	return str
}
