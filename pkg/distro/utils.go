package distro

import (
	"github.com/maxgio92/krawler/pkg/packages"
)

func RepositorySliceContains(s []packages.Repository, e packages.Repository) bool {
	for _, v := range s {
		if v.URI == e.URI {
			return true
		}
	}

	return false
}
