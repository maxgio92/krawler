package distro

type Config struct {
	// A list of Mirrors to scrape
	Mirrors []Mirror

	// A list of architecture for to which scrape packages
	Archs []Arch

	// A list of Distro versions
	Versions []Version
}

type Arch string

type Distro interface {
	Configure(Config) error
	GetPackages(Filter, map[string]interface{}) ([]Package, error)
}

type Version string

type Type string

type Mirror struct {
	// The base URL of the package mirror
	// (e.g. https://mirrors.kernel.org/<distribution>)
	URL string

	// The mirrored repositories
	//
	// TODO: actually we scrape all repositories for all mirrors,
	// independently.
	// Evaluate to decouple Repository from Mirror.
	Repositories []Repository
}

type Repository struct {
	Name string
	URI  URITemplate
}

type URITemplate string

type Package string

// A package filter string prefix.
type Filter string
