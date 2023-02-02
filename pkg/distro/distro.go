package distro

import (
	"github.com/maxgio92/krawler/pkg/output"
	"github.com/maxgio92/krawler/pkg/packages"
	"github.com/maxgio92/krawler/pkg/utils/template"
)

type Config struct {
	// A list of Mirrors to scrape.
	Mirrors []packages.Mirror

	// The mirrored repositories.
	Repositories []packages.Repository

	// A list of architecture for to which scrape packages.
	Archs []packages.Architecture

	// A list of Distro versions.
	Versions []Version

	// Options for visual output.
	Output output.Options `json:"output,omitempty"`
}

type Distro interface {
	// Configure expects distro.Config and arbitrary variables
	// for config fields that support templating.
	Configure(Config) error

	// GetPackages should return a slice of Package based on
	// the provided SearchOptions-type filter.
	SearchPackages(packages.SearchOptions) ([]packages.Package, error)
}

type Version string

type Type string

// BuildTemplates computes templated Config fields by evaluating the template against a set of variables,
// expected as a map of string to interface argument.
// As of now, only the URI field of Config.Repositories is a supported field to be templated.
func (c *Config) BuildTemplates(vars map[string]interface{}) error {
	uris := []string{}

	for _, repository := range c.Repositories {
		if repository.URI != "" {
			result, err := template.MultiplexAndExecute(string(repository.URI), vars)
			if err != nil {
				return err
			}

			uris = append(uris, result...)
		}
	}

	r := []packages.Repository{}
	for _, v := range uris {
		r = append(r, packages.Repository{Name: "", URI: packages.URITemplate(v)})
	}

	c.Repositories = r

	return nil
}
