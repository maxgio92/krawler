package deb

import (
	"io"
)

type Package struct {
	Name        string
	Arch        string
	Version     string
	Release     string
	Location    string
	Url         string
	fileReaders []io.Reader
}

type PackageLocation struct {
	Href string
}

func (p *Package) GetName() string {
	return p.Name
}

func (p *Package) GetVersion() string {
	return p.Version
}

func (p *Package) GetRelease() string {
	return p.Release
}

func (p *Package) GetArch() string {
	return p.Arch
}

func (p *Package) GetLocation() string {
	return p.Location
}

func (p *Package) URL() string {
	return p.Url
}

func (p *Package) FileReaders() []io.Reader {
	return p.fileReaders
}
