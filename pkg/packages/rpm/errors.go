package rpm

import "errors"

var (
	metadataUrlNotValidErr       = errors.New("metadata url is not valid")
	metadataInvalidResponseErr   = errors.New("metadata url returned an invalid response")
	repositoryUrlNotValidErr     = errors.New("repository url is not valid")
	repositoryInvalidResponseErr = errors.New("repository url returned an invalid response")
	packageUrlNotFoundErr        = errors.New("package url not found")
	packageUrlInvalidResponseErr = errors.New("package url returned an invalid response")
)
