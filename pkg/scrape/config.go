package scrape

type Config struct {

	// A list of Mirrors to scrape
	Mirrors            []Mirror

	// A list of architecture for which scrape packages
	Archs              []Arch

	// TODO: split URI into: Distro version, Repository, Architecture, Package URL / Packages DB URL
	PackagesURIFormats []PackagesURIFormat
}
