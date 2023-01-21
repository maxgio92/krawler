package distro

import (
	"github.com/maxgio92/krawler/pkg/packages"
)

type Config struct {
	// A list of Mirrors to scrape.
	Mirrors []packages.Mirror

	// The mirrored repositories.
	Repositories []packages.Repository

	// A list of architecture for to which scrape packages.
	Archs []Arch

	// A list of Distro versions.
	Versions []Version
}

type Arch string

type Distro interface {
	// Configure expects distro.Config and arbitrary variables
	// for config fields that support templating.
	Configure(Config, map[string]interface{}) error

	// GetPackages should return a slice of Package based on
	// the provided PackageOptions-type filter.
	GetPackages(packages.PackageOptions) ([]packages.Package, error)
}

type Version string

type Type string
