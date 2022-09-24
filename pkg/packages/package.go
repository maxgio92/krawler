package packages

import "io"

type Package interface {
	GetName() string
	GetVersion() string
	GetRelease() string
	GetArch() string
	GetLocation() string
	String() string
	URL() string
	FileReaders() []io.Reader
}

type Filter struct {
	packageName      string
	packageFileNames []string
}

func NewFilter(packageName string, packageFileNames ...string) *Filter {
	return &Filter{
		packageName:      packageName,
		packageFileNames: packageFileNames,
	}
}

func (f *Filter) String() string {
	return f.packageName
}
