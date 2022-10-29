package distro

//nolint:nosnakecase
const (
	X86_64Arch Arch = "x86_64"

	// Default architecture for which scrape for packages.
	DefaultArch          = X86_64Arch
	CentosType           = "centos"
	AmazonLinuxV1Type    = "amazonlinux1"
	AmazonLinuxV2Type    = "amazonlinux2"
	AmazonLinuxV2022Type = "amazonlinux2022"
)
