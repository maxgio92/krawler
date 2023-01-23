package packages

import (
	log "github.com/sirupsen/logrus"

	"github.com/maxgio92/krawler/pkg/output"
)

type SearchOptions struct {
	packageName      string
	architectures    []Architecture
	packageFileNames []string
	seedURLs         []string
	*output.ProgressOptions
	progressMessage string
	*MPSCQueue
	verbosity output.Verbosity
	logger    *output.Logger
}

func NewSearchOptions(packageName string, architectures []Architecture, seedURLs []string, verbosity output.Verbosity, progressMessage string, packageFileNames ...string) *SearchOptions {
	logger := output.NewLogger()
	logger.SetLevel(log.Level(verbosity))
	logger.SetFormatter(&log.TextFormatter{
		ForceColors:      true,
		DisableTimestamp: true,
	})

	progressOptions := output.NewProgressOptions(len(seedURLs), progressMessage)

	queue := NewMPSCQueue(len(seedURLs))

	return &SearchOptions{
		packageName:      packageName,
		architectures:    architectures,
		packageFileNames: packageFileNames,
		seedURLs:         seedURLs,
		ProgressOptions:  progressOptions,
		progressMessage:  progressMessage,
		MPSCQueue:        queue,
		verbosity:        verbosity,
		logger:           logger,
	}
}

func (o *SearchOptions) PackageName() string {
	return o.packageName
}

func (o *SearchOptions) PackageFileNames() []string {
	return o.packageFileNames
}

func (o *SearchOptions) SeedURLs() []string {
	return o.seedURLs
}

func (o *SearchOptions) Log() *output.Logger {
	return o.logger
}

func (o *SearchOptions) Verbosity() output.Verbosity {
	return o.verbosity
}

func (o *SearchOptions) Architectures() []Architecture {
	return o.architectures
}

func (o *SearchOptions) ProgressMessage() string {
	return o.progressMessage
}
