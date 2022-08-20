package distro

const (
	X86_64Arch Arch = "x86_64"

	// Default architecture for which scrape for packages.
	DefaultArch = X86_64Arch
	CentosType  = "centos"
)

var (
	debugScrape = true
)
