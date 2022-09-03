package rpm

import (
	"encoding/xml"
	"fmt"
)

type Package struct {
	XMLName     xml.Name        `xml:"package"`
	Name        string          `xml:"name"`
	Arch        string          `xml:"arch"`
	Version     PackageVersion  `xml:"version"`
	Summary     string          `xml:"summary"`
	Description string          `xml:"description"`
	Packager    string          `xml:"packager"`
	URL         string          `xml:"url"`
	Time        PackageTime     `xml:"time"`
	Size        PackageSize     `xml:"size"`
	Location    PackageLocation `xml:"location"`
	Format      PackageFormat   `xml:"format"`
}

func (p *Package) GetName() string {
	return p.Name
}

func (p *Package) GetVersion() string {
	return p.Version.Ver
}

func (p *Package) GetRelease() string {
	return p.Version.Rel
}

func (p *Package) GetArch() string {
	return p.Arch
}

func (p *Package) String() string {
	return fmt.Sprintf("%s-%s-%s.%s", p.GetName(), p.GetVersion(), p.GetRelease(), p.GetArch())
}

func (p *Package) GetLocation() string {
	return p.Location.Href
}

func (p *Package) Unpack() ([]byte, error) {
	return []byte{}, nil
}
