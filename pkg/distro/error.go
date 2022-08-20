package distro

import "errors"

var (
	ErrDistroNotConfigured      = errors.New("the distro has not been configured")
	ErrDistroNotFound           = errors.New("no distribution found with the specified name")
	ErrNoDistroVersionSpecified = errors.New("no versions specified")
	ErrDomainsFromMirrorUrls    = errors.New("error while retrieving DNS names from mirrors URLs")
)
