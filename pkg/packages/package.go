package packages

import (
	"github.com/maxgio92/krawler/pkg/output"
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

type PackageOptions struct {
	packageName      string
	packageFileNames []string
	verbosity        output.Verbosity
}

func NewFilter(verbosity output.Verbosity, packageName string, packageFileNames ...string) *PackageOptions {
	return &PackageOptions{
		packageName:      packageName,
		packageFileNames: packageFileNames,
	}
}

func (o *PackageOptions) PackageName() string {
	return o.packageName
}

func (o *PackageOptions) PackageFileNames() []string {
	return o.packageFileNames
}

func (o *PackageOptions) Verbosity() output.Verbosity {
	return o.verbosity
}
