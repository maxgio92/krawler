package distro

const (
	X86_64Arch Arch = "x86_64"

	// Default architecture for which scrape for packages
	DefaultArch = X86_64Arch

	// Default regex to base the distro version detection on
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
	CentosDefaultConfig Config = Config{
		Mirrors: []Mirror{
			{
				Url: "https://mirrors.edge.kernel.org/centos/",
				Repositories: []Repository{
					{Name: "BaseOS", PackagesURITemplate: PackagesURITemplate("/BaseOS/" + DefaultArch + "/os/Packages/")},
					{Name: "AppStream", PackagesURITemplate: PackagesURITemplate("/AppStream/" + DefaultArch + "/os/Packages/")},
					{Name: "Devel", PackagesURITemplate: PackagesURITemplate("/Devel/" + DefaultArch + "/os/Packages/")},
				},
			},
			{
				Url: "https://archive.kernel.org/centos-vault/",
				Repositories: []Repository{
					{Name: "base", PackagesURITemplate: PackagesURITemplate("/os/" + DefaultArch + "/Packages/")},
					{Name: "updates", PackagesURITemplate: PackagesURITemplate("/updates/" + DefaultArch + "/Packages/")},
				},
			},
		},
		Archs:    []Arch{DefaultArch},
		Versions: nil,
	}

	debugScrape = true
)
