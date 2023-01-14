package deb

import (
	"github.com/maxgio92/krawler/pkg/output"
	"pault.ag/go/archive"
	"sync"
)

// TODO: filter by architecture
type SearchOptions struct {
	packageName string
	seedURLs    []string
	*output.ProgressOptions
	*SyncOptions
}

func NewSearchOptions(packageName string, seedURLs []string, message ...string) *SearchOptions {

	progressOptions := output.NewProgressOptions(len(seedURLs), message...)

	waitGroup := sync.WaitGroup{}
	waitGroup.Add(len(seedURLs))

	resultCh := make(chan []archive.Package)

	errCh := make(chan error)

	doneCh := make(chan bool)

	syncOptions := NewSyncOptions(&waitGroup, resultCh, errCh, doneCh)

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
