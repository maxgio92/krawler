package distro

type Config struct {
	// A list of Mirrors to scrape.
	Mirrors []Mirror

	// The mirrored repositories.
	Repositories []Repository

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
	// the provided Filter-type filter.
	GetPackages(Filter) ([]Package, error)
}

type Version string

type Type string

type Mirror struct {
	// The base URL of the package mirror
	// (e.g. https://mirrors.kernel.org/<distribution>)
	URL string
}

type Repository struct {
	Name string
	URI  URITemplate
}

type URITemplate string

type Package string

// A package filter string prefix.
type Filter Prefix

type Prefix string
