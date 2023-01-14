package deb

// TODO: filter by architecture
type SearchOptions struct {
	packageName string
	distURLs    []string
}

func NewSearchOptions(packageName string, distURLs []string) *SearchOptions {
	return &SearchOptions{
		packageName,
		distURLs,
	}
}

func (o *SearchOptions) PackageName() string {
	return o.packageName
}

func (o *SearchOptions) DistURLs() []string {
	return o.distURLs
}
