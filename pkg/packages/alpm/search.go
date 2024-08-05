//go:build archlinux

package alpm

import (
	"github.com/maxgio92/krawler/pkg/packages"
)

type SearchOptions struct {
	*packages.SearchOptions
	packageNames []string
}

// NewSearchOptions returns a pointer to a SearchOptions object from a pointer to a packages.SearchOptions, and
// overriding architectures and seedURLs.
func NewSearchOptions(options *packages.SearchOptions, seedURLs []string, packageNames []string) *SearchOptions {
	return &SearchOptions{
		packages.NewSearchOptions(
			options.PackageName(),
			nil,
			seedURLs,
			options.Verbosity(),
			options.ProgressMessage(),
			options.PackageFileNames()...,
		),
		packageNames,
	}
}

func (o *SearchOptions) PackageNames() []string {
	return o.packageNames
}
