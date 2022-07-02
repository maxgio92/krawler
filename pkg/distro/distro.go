package distro

var DistroByType = map[DistroType]Distro{}

func NewDistro(distroType DistroType) (Distro, error) {
	distro, ok := DistroByType[distroType]

	if !ok {
		return nil, errDistroNotFound
	}

	return distro, nil
}
