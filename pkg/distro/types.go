package distro

type Config struct {

	// A list of Mirrors to scrape
	Mirrors []Mirror

	// A list of architecture for to which scrape packages
	Archs []Arch

	// A list of Distro versions
	Versions []DistroVersion
}

type Arch string

type Distro interface {
	GetPackages(Config, Filter) ([]Package, error)
}

type DistroVersion string

type DistroType string

type Mirror struct {

	// The base URL of the package mirror
	// (e.g. https://mirrors.kernel.org/<distribution>)
	Url string

	// The mirrored repositories
	//
	// TODO: actually we scrape all repositories for all mirrors,
	// independently.
	// Evaluate to decouple Repository from Mirror.
	Repositories []Repository
}

type Repository struct {
	Name                string
	PackagesURITemplate PackagesURITemplate
}

type PackagesURITemplate string

type Package string

// A package filter string prefix
type Filter string
