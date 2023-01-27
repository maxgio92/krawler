package deb

import (
	"github.com/maxgio92/krawler/pkg/packages"
)

type SearchOptions struct {
	// Deb components are needed to filter indexed packages.
	// More on this here: https://wiki.debian.org/DebianRepository.
	components []string
	*packages.SearchOptions
}

// NewSearchOptions returns a pointer to a SearchOptions object from a pointer to a packages.SearchOptions, and
// overriding architectures and seedURLs.
func NewSearchOptions(options *packages.SearchOptions, architectures []packages.Architecture, seedURLs []string, components []string) *SearchOptions {
	return &SearchOptions{
		components,
		packages.NewSearchOptions(
			options.PackageName(),
			architectures,
			seedURLs,
			options.Verbosity(),
			options.ProgressMessage(),
			options.PackageFileNames()...,
		),
	}
}

func (s *SearchOptions) Components() []string {
	return s.components
}
