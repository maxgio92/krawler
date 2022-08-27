package utils

import (
	"strings"

	kr "github.com/falcosecurity/driverkit/pkg/kernelrelease"
	p "github.com/maxgio92/krawler/pkg/packages"
)

func KernelReleaseFromPackageName(packageName string, packagePrefix string) string {
	ss := strings.Split(packageName, ".")
	ss = strings.Split(strings.Split(packageName, "."+ss[len(ss)-1])[0], packagePrefix+"-")

	return ss[1]
}

func UniqueKernelReleases(kernelReleases []kr.KernelRelease) []kr.KernelRelease {
	krs := make([]kr.KernelRelease, 0, len(kernelReleases))
	m := make(map[kr.KernelRelease]bool)

	for _, v := range kernelReleases {
		if _, ok := m[v]; !ok {
			m[v] = true

			krs = append(krs, v)
		}
	}

	return krs
}

func GetKernelReleaseListFromPackageList(packages []p.Package, packagePrefix string) ([]kr.KernelRelease, error) {
	kernelReleases := []kr.KernelRelease{}

	for _, v := range packages {
		s := KernelReleaseFromPackageName(v.(string), packagePrefix)
		r := kr.FromString(s)
		kernelReleases = append(kernelReleases, r)
	}

	return UniqueKernelReleases(kernelReleases), nil
}
