package kernelrelease

import (
	"fmt"

	p "github.com/maxgio92/krawler/pkg/packages"
)

func versionStringFromPackage(pkg p.Package) string {
	version := pkg.GetVersion()
	if pkg.GetRelease() != "" {
		version += fmt.Sprintf("-%s", pkg.GetRelease())
	}

	if pkg.GetArch() != "" {
		version += fmt.Sprintf(".%s", pkg.GetArch())
	}

	return version
}
