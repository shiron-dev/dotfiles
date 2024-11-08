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
