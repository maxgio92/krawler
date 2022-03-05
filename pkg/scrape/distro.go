package scrape

var DistroByName = map[DistroName]Distro{}

type Distro interface {
	GetPackages(MirrorsConfig, Filter) ([]Package, error)
}

type DistroName int16

const (
	Unknown Distro = iota
	Centos
)

func Factory(distro DistroName) (Distro, error) {
	distro, ok := DistroByName[distro]

	if !ok {
		return nil, errDistroNotFound
	}

	return distro, nil
}
