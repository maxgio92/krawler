package scrape

type Config struct {

	// A list of Mirrors to scrape
	Mirrors  []Mirror

	// A list of architecture for to which scrape packages
	Archs    []Arch

	// A list of Distro versions
	Versions []DistroVersion
}
