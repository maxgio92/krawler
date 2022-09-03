package packages

type Package interface {
	GetName() string
	GetVersion() string
	GetRelease() string
	GetArch() string
	String() string
	Unpack() ([]byte, error)
}

type Filter string
