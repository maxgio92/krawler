package distro

var DistroByType = map[Type]Distro{}

func NewDistro(distroType Type) (Distro, error) {
	distro, ok := DistroByType[distroType]

	if !ok {
		return nil, errDistroNotFound
	}

	return distro, nil
}
