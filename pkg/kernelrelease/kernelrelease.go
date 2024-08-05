package kernelrelease

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
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

//nolint:cyclop
func (k *KernelRelease) BuildFromPackage(pkg p.Package) error {
	k.PackageName = pkg.GetName()
	k.PackageURL = pkg.URL()
	k.Architecture = Arch(pkg.GetArch())

	kernelVersion := versionStringFromPackage(pkg)
	match := kernelVersionPattern.FindStringSubmatch(kernelVersion)

	identifiers := make(map[string]string)

	for i, name := range kernelVersionPattern.SubexpNames() {
		if i > 0 && i <= len(match) {
			identifiers[name] = match[i]

			switch name {
			case "fullversion":
				k.Fullversion = match[i]
			case "version":
				k.Version, _ = strconv.Atoi(match[i])
			case "patchlevel":
				k.PatchLevel, _ = strconv.Atoi(match[i])
			case "sublevel":
				k.Sublevel, _ = strconv.Atoi(match[i])
			case "extraversion":
				k.Extraversion = match[i]
			case "fullextraversion":
				k.FullExtraversion = match[i]
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

func (k *KernelRelease) SHA256Sum() string {
	sha256.New()

	s := fmt.Sprintf("%s%s%s%s",
		k.Fullversion,
		k.FullExtraversion,
		k.PackageName,
		string(k.Architecture),
	)

	hash := sha256.Sum256([]byte(s))

	return hex.EncodeToString(hash[:])
}
