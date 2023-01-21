package packages

import (
	"io"
)

type Package interface {
	GetName() string
	GetVersion() string
	GetRelease() string
	GetArch() string
	GetLocation() string
	URL() string
	FileReaders() []io.Reader
}

type Architecture string
