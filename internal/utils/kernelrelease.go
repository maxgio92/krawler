package utils

import (
	"strings"

	"github.com/falcosecurity/driverkit/pkg/kernelrelease"
)

func KernelReleaseFromPackageName(packageName string, packagePrefix string) string {
	ss := strings.Split(packageName, ".")
	ss = strings.Split(strings.Split(packageName, "."+ss[len(ss)-1])[0], packagePrefix+"-")

	return ss[1]
}

func UniqueKernelReleases(kernelReleases []kernelrelease.KernelRelease) []kernelrelease.KernelRelease {
	krs := make([]kernelrelease.KernelRelease, 0, len(kernelReleases))
	m := make(map[kernelrelease.KernelRelease]bool)

	for _, v := range kernelReleases {
		if _, ok := m[v]; !ok {
			m[v] = true
			krs = append(krs, v)
		}
	}

	return krs
}
