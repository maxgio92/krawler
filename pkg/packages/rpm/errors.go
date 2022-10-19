package rpm

import "errors"

var (
	errMetadataURLNotValid       = errors.New("metadata url is not valid")
	errMetadataInvalidResponse   = errors.New("metadata url returned an invalid response")
	errRepositoryURLNotValid     = errors.New("repository url is not valid")
	errRepositoryInvalidResponse = errors.New("repository url returned an invalid response")
	errPackageURLNotFound        = errors.New("package url not found")
	errPackageURLInvalidResponse = errors.New("package url returned an invalid response")
)
