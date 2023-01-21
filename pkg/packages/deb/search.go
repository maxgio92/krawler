package deb

import (
	"github.com/maxgio92/krawler/pkg/packages"
)

type SearchOptions struct {
	*packages.SearchOptions
}

// NewSearchOptions returns a pointer to a SearchOptions object from a pointer to a packages.SearchOptions, and
// overriding architectures and seedURLs.
func NewSearchOptions(options *packages.SearchOptions, architectures []packages.Architecture, seedURLs []string) *SearchOptions {
	return &SearchOptions{
		packages.NewSearchOptions(
			options.PackageName(),
			architectures,
			seedURLs,
			options.Verbosity(),
			options.ProgressMessage(),
			options.PackageFileNames()...,
		)}
}
