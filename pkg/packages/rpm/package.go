package rpm

import "fmt"

func (p *Package) GetName() string {
	return p.Name
}

func (p *Package) GetVersion() string {
	return p.Version.Ver
}

func (p *Package) GetRelease() string {
	return p.Version.Rel
}

func (p *Package) GetArch() string {
	return p.Arch
}

func (p *Package) String() string {
	return fmt.Sprintf("%s-%s-%s.%s", p.GetName(), p.GetVersion(), p.GetRelease(), p.GetArch())
}

func (p *Package) GetLocation() string {
	return p.Location.Href
}

func (p *Package) Unpack() ([]byte, error) {
	return []byte{}, nil
}
