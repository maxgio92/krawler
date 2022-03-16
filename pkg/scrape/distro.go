package scrape

var DistroByType = map[DistroType]Distro{}

type Distro interface {
	GetPackages(Config, Filter) ([]Package, error)
}

type DistroVersion string

type DistroType string

const (
	CentosType = "centos"
)

func Factory(distroType DistroType) (Distro, error) {
	distro, ok := DistroByType[distroType]

	if !ok {
		return nil, errDistroNotFound
	}

	return distro, nil
}
