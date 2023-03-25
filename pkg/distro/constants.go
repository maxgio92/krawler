package distro

import "github.com/maxgio92/krawler/pkg/packages"

const (
	X8664Arch packages.Architecture = "x86_64"

	// Default architecture for which scrape for packages.
	DefaultArch          = X8664Arch
	CentosType           = "centos"
	AmazonLinuxV1Type    = "amazonlinux"
	AmazonLinuxV2Type    = "amazonlinux2"
	AmazonLinuxV2022Type = "amazonlinux2022"
	DebianType           = "debian"
	UbuntuType           = "ubuntu"
	FedoraType           = "fedora"
)
