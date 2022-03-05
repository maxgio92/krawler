package scrape

func init() {
	DistroByName[Centos] = &Centos{}
}

var (
	CentosMirrorsDistroVersionRegex = `^(0|[1-9]\d*)(\.(0|[1-9]\d*)?)?(\.(0|[1-9]\d*)?)?(-[a-zA-Z\d][-a-zA-Z.\d]*)?(\+[a-zA-Z\d][-a-zA-Z.\d]*)?\/$`

	CentosDefaultConfig = Config{
		Mirrors: []Mirror{
			Mirror{baseUrl: "https://mirrors.edge.kernel.org/centos/"},
			Mirror{baseUrl: "https://archive.kernel.org/centos-vault/"},
			Mirror{baseUrl: "https://mirror.nsc.liu.se/centos-store/"},
		},
		Archs: DefaultArch,
		// TODO: evaluate to be not configurable
		// Why do we need to configure this?
		PackagesURIFormats: []string{
			// Repo / Arch / Standard path/to/packages
			"/BaseOS/%s/os/Packages/",
			"/os/%s/Packages/",
			"/updates/%s/Packages/",
		}
	}
)

type Centos struct {
	mirrors []Mirror
	archs   []Arch
	filter  Filter
}

func (c *Centos) GetPackages() ([]Package, error) {

	c.mirrors = CentosDefaultConfig.Mirrors

	// TODO: scrape

	return []Package{}, nil
}
