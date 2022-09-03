package packages

import "io"

type Package interface {
	GetName() string
	GetVersion() string
	GetRelease() string
	GetArch() string
	GetLocation() string
	String() string

	// Unpack returns a list of readers.
	Unpack() ([]*io.Reader, error)

	URL() string
}

type Filter string
