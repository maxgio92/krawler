package debian

import (
	"fmt"

	"github.com/maxgio92/krawler/internal/utils"

	"pault.ag/go/archive"
)

func GetReleasesFromPackages(packages []archive.Package) ([]string, error) {
	releases := []string{}
	if len(packages) > 0 {
		for _, v := range packages {
			releases = append(releases, fmt.Sprintf("%s-%s-%s", v.Version.Version, v.Version.Revision, v.Architecture.CPU))
		}
	}

	return utils.Unique(releases), nil
}
