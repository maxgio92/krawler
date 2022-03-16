package scrape

const (
	// Default architecture for which scrape for packages
	DefaultArch = X86_64Arch

	// Default regex to base the distro version detection on
	CentosMirrorsDistroVersionRegex = `^(0|[1-9]\d*)(\.(0|[1-9]\d*)?)?(\.(0|[1-9]\d*)?)?(-[a-zA-Z\d][-a-zA-Z.\d]*)?(\+[a-zA-Z\d][-a-zA-Z.\d]*)?\/$`
	CentosPackageFileExtension      = "rpm"
)

var (
	// Default configuration for scrape Centos (RPM) packages
	CentosDefaultConfig Config = Config{
		Mirrors: []Mirror{
			Mirror{
				Url: "https://mirrors.edge.kernel.org/centos/",
				Repositories: []Repository{
					Repository{Name: "BaseOS", PackagesURIFormat: "/BaseOS/%s/os/Packages/"},
					Repository{Name: "AppStream", PackagesURIFormat: "/AppStream/%s/os/Packages/"},
					Repository{Name: "Devel", PackagesURIFormat: "/Devel/%s/os/Packages/"},
				},
			},
			Mirror{
				Url: "https://archive.kernel.org/centos-vault/",
				Repositories: []Repository{
					Repository{Name: "base", PackagesURIFormat: "/os/%s/Packages/"},
					Repository{Name: "updates", PackagesURIFormat: "/updates/%s/Packages/"},
				},
			},
		},
		Archs:    []Arch{DefaultArch},
		Versions: nil,
	}

	debugScrape = true
)
