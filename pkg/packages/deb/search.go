package deb

import (
	"github.com/maxgio92/krawler/pkg/output"
	log "github.com/sirupsen/logrus"
)

// TODO: filter by architecture
type SearchOptions struct {
	packageName string
	seedURLs    []string
	*output.ProgressOptions
	*MPSCQueue
	logger *output.Logger
}

func NewSearchOptions(packageName string, seedURLs []string, verbosity output.Verbosity, message ...string) *SearchOptions {
	logger := output.NewLogger()
	logger.SetLevel(log.Level(verbosity))
	logger.SetFormatter(&log.TextFormatter{
		ForceColors:      true,
		DisableTimestamp: true,
	})

	progressOptions := output.NewProgressOptions(len(seedURLs), message...)

	syncOptions := NewSyncOptions(len(seedURLs))

	return &SearchOptions{
		packageName,
		seedURLs,
		progressOptions,
		syncOptions,
		logger,
	}
}

func (o *SearchOptions) PackageName() string {
	return o.packageName
}

func (o *SearchOptions) SeedURLs() []string {
	return o.seedURLs
}

func (o *SearchOptions) Log() *output.Logger {
	return o.logger
}

func (o *SearchOptions) Verbosity() output.Verbosity {
	return output.Verbosity(o.logger.Level)
}
