package utils

import (
	"strings"

	kr "github.com/falcosecurity/driverkit/pkg/kernelrelease"

	p "github.com/maxgio92/krawler/pkg/packages"
)

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

func GetKernelReleaseListFromPackageList(packages []p.Package, prefix string) ([]kr.KernelRelease, error) {
	releases := []kr.KernelRelease{}

	for _, pkg := range packages {
		version := strings.TrimPrefix(pkg.String(), prefix)

		r := kr.FromString(version)
		releases = append(releases, r)
	}

	return UniqueKernelReleases(releases), nil
}
