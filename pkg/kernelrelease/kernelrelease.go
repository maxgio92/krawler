package kernelrelease

import (
	"fmt"
	"regexp"
	"strconv"

	p "github.com/maxgio92/krawler/pkg/packages"
)

type Arch string

type Archs map[Arch]string

type KernelRelease struct {
	Fullversion      string `json:"full_version"`
	Version          int    `json:"version"`
	PatchLevel       int    `json:"patch_level"`
	Sublevel         int    `json:"sublevel"`
	Extraversion     string `json:"extra_version"`
	FullExtraversion string `json:"full_extra_version"`
	Architecture     Arch   `json:"architecture"`
	PackageName      string `json:"package_name"`
	PackageURL       string `json:"package_url"`
	CompilerVersion  string `json:"compiler_version"`
}

var kernelVersionPattern = regexp.MustCompile(`(?P<fullversion>^(?P<version>0|[1-9]\d*)\.(?P<patchlevel>0|[1-9]\d*)[.+]?(?P<sublevel>0|[1-9]\d*)?)(?P<fullextraversion>[-.+](?P<extraversion>0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)([\.+~](0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-_]*))*)?(\+[0-9a-zA-Z-]+(\.[0-9a-zA-Z-]+)*)?$`)

//nolint:cyclop
func (k *KernelRelease) BuildFromPackage(pkg p.Package) error {
	k.PackageName = pkg.GetName()
	k.PackageURL = pkg.URL()
	k.Architecture = Arch(pkg.GetArch())

	kernelVersion := VersionStringFromPackage(pkg)
	match := kernelVersionPattern.FindStringSubmatch(kernelVersion)

	identifiers := make(map[string]string)

	for i, name := range kernelVersionPattern.SubexpNames() {
		if i > 0 && i <= len(match) {
			var err error

			identifiers[name] = match[i]

			switch name {
			case "fullversion":
				k.Fullversion = match[i]
			case "version":
				k.Version, err = strconv.Atoi(match[i])
			case "patchlevel":
				k.PatchLevel, err = strconv.Atoi(match[i])
			case "sublevel":
				k.Sublevel, err = strconv.Atoi(match[i])
			case "extraversion":
				k.Extraversion = match[i]
			case "fullextraversion":
				k.FullExtraversion = match[i]
			}

			if err != nil {
				return err
			}
		}
	}

	compilerVersion, err := GetCompilerVersionFromKernelPackage(pkg)
	if err != nil {
		k.CompilerVersion = ""
	}

	k.CompilerVersion = compilerVersion

	return nil
}

func VersionStringFromPackage(pkg p.Package) string {
	version := pkg.GetVersion()
	if pkg.GetRelease() != "" {
		version += fmt.Sprintf("-%s", pkg.GetRelease())
	}

	if pkg.GetArch() != "" {
		version += fmt.Sprintf(".%s", pkg.GetArch())
	}

	return version
}

func UniqueKernelReleases(kernelReleases []KernelRelease) []KernelRelease {
	krs := make([]KernelRelease, 0, len(kernelReleases))
	m := make(map[KernelRelease]bool)

	for _, v := range kernelReleases {
		if _, ok := m[v]; !ok {
			m[v] = true

			krs = append(krs, v)
		}
	}

	return krs
}

func GetKernelReleaseListFromPackageList(packages []p.Package, prefix string) ([]KernelRelease, error) {
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

	return UniqueKernelReleases(releases), nil
}
