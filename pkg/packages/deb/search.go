package deb

import (
	"github.com/maxgio92/krawler/pkg/output"
)

// TODO: filter by architecture
type SearchOptions struct {
	packageName string
	seedURLs    []string
	*output.ProgressOptions
	*MPSCQueue
}

func NewSearchOptions(packageName string, seedURLs []string, message ...string) *SearchOptions {

	progressOptions := output.NewProgressOptions(len(seedURLs), message...)

	syncOptions := NewSyncOptions(len(seedURLs))

	return &SearchOptions{
		packageName,
		seedURLs,
		progressOptions,
		syncOptions,
	}
}

func (o *SearchOptions) PackageName() string {
	return o.packageName
}

func (o *SearchOptions) SeedURLs() []string {
	return o.seedURLs
}
