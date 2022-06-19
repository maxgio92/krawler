package scrape

var DistroByType = map[DistroType]Distro{}

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
