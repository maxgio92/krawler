package distro

const (
	X86_64Arch Arch = "x86_64"

	// Default architecture for which scrape for packages.
	DefaultArch = X86_64Arch

	// Default regex to base the distro version detection on.
	CentosMirrorsDistroVersionRegex = `^(0|[1-9]\d*)(\.(0|[1-9]\d*)?)?(\.(0|[1-9]\d*)?)?(-[a-zA-Z\d][-a-zA-Z.\d]*)?(\+[a-zA-Z\d][-a-zA-Z.\d]*)?\/$`
	CentosPackageFileExtension      = "rpm"
	CentosType                      = "centos"
)

var (
	// Default configuration for scrape Centos (RPM) packages.
	//
	// TODO: support for templating default PackagesURI.
	// As of now URI templating depends on distro's viper.AllSettings()
	// data.
	CentosDefaultConfig = Config{
		Mirrors: []Mirror{
			{
				URL: "https://mirrors.edge.kernel.org/centos/",
				Repositories: []Repository{
					{Name: "BaseOS", URI: URITemplate("/BaseOS/" + DefaultArch + "/os/Packages/")},
					{Name: "AppStream", URI: URITemplate("/AppStream/" + DefaultArch + "/os/Packages/")},
					{Name: "Devel", URI: URITemplate("/Devel/" + DefaultArch + "/os/Packages/")},
				},
			},
			{
				URL: "https://archive.kernel.org/centos-vault/",
				Repositories: []Repository{
					{Name: "base", URI: URITemplate("/os/" + DefaultArch + "/Packages/")},
					{Name: "updates", URI: URITemplate("/updates/" + DefaultArch + "/Packages/")},
				},
			},
		},
		Archs:    []Arch{DefaultArch},
		Versions: nil,
	}

	debugScrape = true
)
