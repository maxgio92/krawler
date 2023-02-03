package kernelrelease

import p "github.com/maxgio92/krawler/pkg/packages"

func GetKernelReleasesFromPackages(packages []p.Package, prefix string) ([]KernelRelease, error) {
	releases := []KernelRelease{}

	for _, pkg := range packages {
		kr := &KernelRelease{}

		err := kr.BuildFromPackage(pkg)
		if err != nil {
			return []KernelRelease{}, err
		}

		if kr.Fullversion != "" {
			releases = append(releases, *kr)
		}
	}

	return unique(releases), nil
}

func unique(kernelReleases []KernelRelease) []KernelRelease {
	krs := make([]KernelRelease, 0, len(kernelReleases))
	m := make(map[string]bool)

	for _, v := range kernelReleases {
		if _, ok := m[v.SHA256Sum()]; !ok {
			m[v.SHA256Sum()] = true

			krs = append(krs, v)
		}
	}

	return krs
}
