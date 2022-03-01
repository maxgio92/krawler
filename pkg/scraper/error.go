package scraper

import "errors"

var (
	errScraperNotFound              = errors.New("no scraper found for the specified Linux distribution")
	errDistributionVersionsNotFound = errors.New("no distribution versions found with specified mirrors config")
	errPackagesNotFound             = errors.New("no packages found")
	errMirrorsNotFound              = errors.New("no mirrors found")
	errDomainsFromMirrorUrls        = errors.New("error while retrieving DNS names from mirrors URLs")
)
