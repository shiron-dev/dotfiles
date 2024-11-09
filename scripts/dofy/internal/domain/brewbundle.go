package domain

type BrewBundleType uint

const (
	BrewBundleTypeTap BrewBundleType = iota
	BrewBundleTypeFormula
	BrewBundleTypeCask
	BrewBundleTypeMas
)

type BrewBundle struct {
	Name       string
	BundleType BrewBundleType
}

func (b BrewBundle) String() string {
	var str string

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

	str += " " + b.Name

	return str
}
