package distro

import (
	"github.com/maxgio92/krawler/pkg/output"
	"github.com/maxgio92/krawler/pkg/packages"
)

type Config struct {

	// A list of Mirrors to scrape.
	Mirrors []packages.Mirror

	// The mirrored repositories.
	Repositories []packages.Repository

	// A list of architecture for to which scrape packages.
	Architectures []packages.Architecture `json:"archs"`

	// A list of Distro versions.
	Versions []Version

	// Options for visual output.
	Output output.Options `json:"output"`
}

type Distro interface {

	// Configure expects distro.Config and arbitrary variables
	// for config fields that support templating.
	//Configure(Config, map[string]interface{}) error
	Configure(Config) error

	// GetPackages should return a slice of Package based on
	// the provided SearchOptions-type filter.
	SearchPackages(packages.SearchOptions) ([]packages.Package, error)
}

type Version string

type Type string
