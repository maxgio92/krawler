package distro

import "errors"

var (
	errDistroNotFound         = errors.New("no distribution found with the specified name")
	errDistroVersionsNotFound = errors.New("no distribution versions found with specified mirrors config")
	errPackagesNotFound       = errors.New("no packages found")
	errMirrorsNotFound        = errors.New("no mirrors found")
	errDomainsFromMirrorUrls  = errors.New("error while retrieving DNS names from mirrors URLs")
)
