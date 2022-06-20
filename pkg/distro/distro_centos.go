package distro

import (
	"fmt"
	"net/url"
	"path"
	"regexp"
	"strings"

	"github.com/maxgio92/krawler/pkg/template"

	"github.com/gocolly/colly"
	d "github.com/gocolly/colly/debug"
)

func init() {
	DistroByType[CentosType] = &Centos{}
}

type Centos struct{}

// For each mirror, for each distro version, for each repository,
// for each architecture, scrape.
func (c *Centos) GetPackages(userConfig Config, filter Filter, allSettings map[string]interface{}) ([]Package, error) {

	var packages []Package

	// Merge custom config with default config.
	config, err := c.buildConfig(CentosDefaultConfig, userConfig)
	if err != nil {
		return nil, err
	}

	// Get distro versions for which to search packages.
	versions, err := c.buildVersions(config.Mirrors, config.Versions)
	if err != nil {
		return nil, err
	}

	versionsUrls, err := c.buildVersionsUrls(config.Mirrors, versions)
	if err != nil {
		return nil, err
	}

	// Apply repository packages URI for each provided architecture.
	//repositoriesUris, _ := c.buildRepositoriesUris(config.Mirrors, config.Archs, viper)
	repositoriesUris, err := c.buildRepositoriesUris(config.Mirrors, config.Archs, allSettings)
	if err != nil {
		return nil, err
	}

	for _, root := range versionsUrls {
		for _, repositoryUri := range repositoriesUris {

			// Get repository URL from URI.
			repositoryUrl, err := url.Parse(string(root.String() + repositoryUri))
			if err != nil {
				return nil, err
			}

			// Crawl packages based on filter.
			p, err := c.crawlPackages(repositoryUrl, filter, debugScrape)
			if err != nil {
				return nil, err
			}

			packages = append(packages, p...)
		}
	}

	return packages, nil
}

// Returns the final configuration, considering a default and a user-provided
// configuration.
func (c *Centos) buildConfig(def Config, user Config) (Config, error) {
	config, err := c.mergeConfig(def, user)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}

// Returns the final configuration by merging the default with the user provided.
func (c *Centos) mergeConfig(def Config, config Config) (Config, error) {
	if len(config.Archs) < 1 {
		config.Archs = def.Archs
	} else {

		for _, arch := range config.Archs {
			if arch == "" {
				config.Archs = def.Archs
				break
			}
		}
	}

	if len(config.Mirrors) < 1 {
		config.Mirrors = def.Mirrors
	} else {

		for i, mirror := range config.Mirrors {
			if mirror.Url == "" {
				config.Mirrors = def.Mirrors
				break
			}

			if len(mirror.Repositories) < 1 {
				config.Mirrors[i].Repositories = c.getDefaultRepositories()
			} else {
				for _, repository := range mirror.Repositories {
					if repository.PackagesURITemplate == "" {
						config.Mirrors[i].Repositories = c.getDefaultRepositories()
						break
					}
				}
			}
		}
	}

	return config, nil
}

// Returns the final configuration by overriding the default.
func (c *Centos) overrideConfig(def Config, override Config) (Config, error) {
	if len(override.Mirrors) > 0 {
		if override.Mirrors[0].Url != "" {
			return override, nil
		}
	}

	return def, nil
}

// Returns a list of distro versions, considering the user-provided configuration,
// and if not, the ones available on configured mirrors.
func (c *Centos) buildVersions(mirrors []Mirror, staticVersions []DistroVersion) ([]DistroVersion, error) {
	if staticVersions != nil {
		return staticVersions, nil
	}

	var dynamicVersions []DistroVersion

	dynamicVersions, err := c.crawlVersions(mirrors, debugScrape)
	if err != nil {
		return nil, err
	}

	return dynamicVersions, nil
}

// Returns the list of version-specific mirror URLs.
func (c *Centos) buildVersionsUrls(mirrors []Mirror, versions []DistroVersion) ([]*url.URL, error) {
	if (len(versions) > 0) && (len(mirrors) > 0) {
		var versionRoots []*url.URL

		for _, mirror := range mirrors {
			for _, version := range versions {

				versionRoot, err := url.Parse(mirror.Url + string(version))
				if err != nil {
					return nil, err
				}

				versionRoots = append(versionRoots, versionRoot)
			}
		}

		return versionRoots, nil
	}

	return nil, fmt.Errorf("no versions specified")
}

// Returns the list of the current available distro versions, by scraping
// the specified mirrors, dynamically.
func (c *Centos) crawlVersions(mirrors []Mirror, debug bool) ([]DistroVersion, error) {

	var versions []DistroVersion

	versionPattern := regexp.MustCompile(CentosMirrorsDistroVersionRegex)

	var roots []string
	for _, mirror := range mirrors {
		roots = append(roots, mirror.Url)
	}

	allowedDomains, err := getHostnamesFromURLs(roots)
	if err != nil || len(allowedDomains) < 1 {
		return nil, err
	}

	// Create the collector settings
	coOptions := []func(*colly.Collector){
		colly.AllowedDomains(allowedDomains...),
		colly.Async(false),
	}

	if debug {
		coOptions = append(coOptions, colly.Debugger(&d.LogDebugger{}))
	}

	// Create the collector
	co := colly.NewCollector(coOptions...)

	// Visit each distro version-specific folder
	co.OnHTML("a[href]", func(e *colly.HTMLElement) {
		distroVersionFolderMatch := versionPattern.FindStringSubmatch(e.Attr("href"))

		if len(distroVersionFolderMatch) > 0 {
			//nolint:errcheck
			co.Visit(e.Request.AbsoluteURL(e.Attr("href")))
		}
	})

	// Collect all the version folder names
	co.OnRequest(func(r *colly.Request) {
		if !stringSliceContains(roots, r.URL.String()) {
			versions = append(versions, DistroVersion(path.Base(r.URL.Path)))
		}
	})

	// Visit each mirror root folder
	for _, root := range roots {
		err := co.Visit(root)
		if err != nil {
			return nil, err
		}
	}

	return versions, nil
}

// Returns the list of repositories URLs, for each specified architecture.
func (c *Centos) buildRepositoriesUris(mirrors []Mirror, archs []Arch, allSettings map[string]interface{}) ([]string, error) {
	uris := []string{}

	for _, mirror := range mirrors {
		for _, repository := range mirror.Repositories {

			keys, err := template.GetSupportedVariables(string(repository.PackagesURITemplate))
			if err != nil {
				return nil, err
			}

			inventory := map[string][]interface{}{}
			for _, key := range keys {
				centosSettingsI := allSettings[string(CentosType)]

				centosSettings, ok := centosSettingsI.(map[string]interface{})
				if !ok {
					return nil, fmt.Errorf("CentOS configuration does not contain a valid structure")
				}

				centosSettingI := centosSettings[key]
				centosSetting, ok := centosSettingI.([]interface{})
				if !ok {
					return nil, fmt.Errorf("variable '%s' in CentOS repository configuration is not a valid slice", key)
				}

				for _, v := range centosSetting {
					if v != "" {
						inventory[key] = append(inventory[key], v)
					}
				}
			}

			if repository.PackagesURITemplate != "" {

				result, err := template.MultiplexAndExecute(string(repository.PackagesURITemplate), inventory)
				if err != nil {
					panic(err)
				}

				uris = append(uris, result...)
				fmt.Println(uris)
			}
		}
	}

	// Scrape for all possible repositories
	if len(uris) < 1 {
		uris = append(uris, "/")
	}

	return uris, nil
}

// Returns the list of default repositories from the default config
func (c *Centos) getDefaultRepositories() []Repository {
	var repositories []Repository

	for _, mirror := range CentosDefaultConfig.Mirrors {
		for _, repository := range mirror.Repositories {
			if !repositorySliceContains(repositories, repository) {
				repositories = append(repositories, repository)
			}
		}
	}

	return repositories
}

// Returns a list of packages found on each page URL,
// filtered by filter.
func (c *Centos) crawlPackages(seedUrl *url.URL, filter Filter, debug bool) ([]Package, error) {
	var packages []Package

	folderRegex := `.+\/$`
	folderPattern := regexp.MustCompile(folderRegex)

	packageRegex := fmt.Sprintf(`.+\.%s$`, CentosPackageFileExtension)
	packagePattern := regexp.MustCompile(packageRegex)

	filteredPackageRegex := `^` + string(filter) + `.+.` + CentosPackageFileExtension
	filteredPackagePattern := regexp.MustCompile(filteredPackageRegex)

	allowedDomains, err := getHostnamesFromURLs([]string{seedUrl.String()})
	if err != nil {
		return nil, err
	}

	// Create the collector settings
	coOptions := []func(*colly.Collector){
		colly.AllowedDomains(allowedDomains...),
		colly.Async(false),
	}

	if debug {
		coOptions = append(coOptions, colly.Debugger(&d.LogDebugger{}))
	}

	// Create the collector.
	co := colly.NewCollector(coOptions...)

	// Add the callback to Visit the linked resource, for each HTML element found
	co.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")

		// Do not traverse the hierarchy in reverse order.
		if !(strings.Contains(link, "../")) {
			//nolint:errcheck
			co.Visit(e.Request.AbsoluteURL(link))
		}
	})

	// Add the analysis callback to find package URLs, for each Visit call
	co.OnRequest(func(r *colly.Request) {
		folderMatch := folderPattern.FindStringSubmatch(r.URL.String())

		// If the URL is not of a folder.
		if len(folderMatch) == 0 {
			packageMatch := packagePattern.FindStringSubmatch(r.URL.String())

			// If the URL is of a package file.
			if len(packageMatch) > 0 {
				packageName := path.Base(r.URL.String())
				packageNameMatch := filteredPackagePattern.FindStringSubmatch(packageName)

				// If the URL matches the package filter regex.
				if len(packageNameMatch) > 0 {
					packages = append(packages, Package(packageName))
				}
			}

			// Otherwise abort the request.
			r.Abort()
		}
	})

	co.Visit(seedUrl.String())

	return packages, nil
}
