package packages

import "io"

type Package interface {
	GetName() string
	GetVersion() string
	GetRelease() string
	GetArch() string
	GetLocation() string
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

// TODO: rename String to PackageName.
func (f *Filter) String() string {
	return f.packageName
}

func (f *Filter) PackageFileNames() []string {
	return f.packageFileNames
}
