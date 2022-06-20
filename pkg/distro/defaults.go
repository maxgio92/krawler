package distro

const (
	X86_64Arch Arch = "x86_64"

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
			{
				Url: "https://mirrors.edge.kernel.org/centos/",
				Repositories: []Repository{
					{Name: "BaseOS", PackagesURITemplate: "/BaseOS/{{ Arch }}/os/Packages/"},
					{Name: "AppStream", PackagesURITemplate: "/AppStream/{{ Arch }}/os/Packages/"},
					{Name: "Devel", PackagesURITemplate: "/Devel/{{ Arch }}/os/Packages/"},
				},
			},
			{
				Url: "https://archive.kernel.org/centos-vault/",
				Repositories: []Repository{
					{Name: "base", PackagesURITemplate: "/os/{{ Arch }}/Packages/"},
					{Name: "updates", PackagesURITemplate: "/updates/{{ Arch }}/Packages/"},
				},
			},
		},
		Archs:    []Arch{DefaultArch},
		Versions: nil,
	}

	debugScrape = true
)
