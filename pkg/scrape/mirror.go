package scrape

type Mirror struct {

	// The base URL of the package mirror
	// (e.g. https://mirrors.kernel.org/<distribution>)
	Url          string

	// The mirrored repositories
	Repositories []Repository
}
