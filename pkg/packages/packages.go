package packages

type Package interface{}

type Repository struct {
	Name string
	URI  URITemplate
}

type URITemplate string

type Mirror struct {
	Name string
	// The base URL of the package mirror
	// (e.g. https://mirrors.kernel.org/<distribution>)
	URL string
}

type Filter string
