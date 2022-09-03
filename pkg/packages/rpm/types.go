package rpm

import (
	"encoding/xml"
)

type RepositoryMetadata struct {
	XMLName  xml.Name `xml:"repomd"`
	Revision string   `xml:"revision"`
	Data     []Data   `xml:"data"`
}

type Data struct {
	Type     string   `xml:"type,attr"`
	Location Location `xml:"location"`
}

type Location struct {
	Href string `xml:"href,attr"`
}

type PrimaryRepositoryMetadata struct {
	XMLName  xml.Name  `xml:"metadata"`
	Packages []Package `xml:"package"`
}

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

type PackageVersion struct {
	XMLName xml.Name `xml:"version"`
	Epoch   string   `xml:"epoch,attr"`
	Ver     string   `xml:"ver,attr"`
	Rel     string   `xml:"rel,attr"`
}

type PackageTime struct {
	File  string `xml:"file,attr"`
	Build string `xml:"build,attr"`
}

type PackageSize struct {
	Package   string `xml:"package,attr"`
	Installed string `xml:"installed,attr"`
	Archive   string `xml:"archive,attr"`
}

type PackageLocation struct {
	XMLName xml.Name `xml:"location"`
	Href    string   `xml:"href,attr"`
}

type PackageFormat struct {
	XMLName     xml.Name           `xml:"format"`
	License     string             `xml:"license"`
	Vendor      string             `xml:"vendor"`
	Group       string             `xml:"group"`
	Buildhost   string             `xml:"buildhost"`
	HeaderRange PackageHeaderRange `xml:"header-range"`
	Requires    PackageRequires    `xml:"requires"`
	Provides    PackageProvides    `xml:"provides"`
}

type PackageHeaderRange struct {
	Start string `xml:"start,attr"`
	End   string `xml:"end,attr"`
}

type PackageProvides struct {
	XMLName xml.Name `xml:"provides"`
	Entries []Entry  `xml:"entry"`
}

type PackageRequires struct {
	XMLName xml.Name `xml:"requires"`
	Entries []Entry  `xml:"entry"`
}

type Entry struct {
	XMLName xml.Name `xml:"entry"`
	Name    string   `xml:"name,attr"`
}
