package distro

import "errors"

var (
	errDistroNotConfigured      = errors.New("the distro has not been configured")
	errDistroNotFound           = errors.New("no distribution found with the specified name")
	errNoDistroVersionSpecified = errors.New("no versions specified")
	errDomainsFromMirrorUrls    = errors.New("error while retrieving DNS names from mirrors URLs")
)
