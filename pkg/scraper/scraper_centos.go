package scraper

import (
	"errors"
	"fmt"
	"path"
	"regexp"
	"strings"

	"github.com/gocolly/colly"
)

var centosMirrorsDistroVersionRegex = `^(0|[1-9]\d*)(\.(0|[1-9]\d*)?)?(\.(0|[1-9]\d*)?)?(-[a-zA-Z\d][-a-zA-Z.\d]*)?(\+[a-zA-Z\d][-a-zA-Z.\d]*)?\/$`
var centosPackageFileExtension = "rpm"

func init() {
	ScraperByDistro[Centos] = &centosScraper{}
}

type centosScraper struct {
}

// Scrape to find packages
func (c centosScraper) Scrape(mirrorsConfig MirrorsConfig, packagePrefix string) ([]string, error) {
	mirrorSpecificVersionRootURLs, err := seekDistroVersionsURLs(mirrorsConfig)
	if err != nil {
		return nil, errors.New("No distribution versions found with specified mirrors config.")
	}

	if len(mirrorSpecificVersionRootURLs) > 0 {

		packages, err := scrape(mirrorsConfig, mirrorSpecificVersionRootURLs, packagePrefix)
		if err != nil {
			return nil, err
		}

		if len(packages) > 0 {
			return packages, nil
		}

		return nil, errors.New("No packages found.")
	}

	return nil, errors.New("No mirrors found.")
}

// Seek Distro version folders to cycle only on those packages folders directly!
func seekDistroVersionsURLs(mirrorsConfig MirrorsConfig) ([]string, error) {
	centosMirrorsDistroVersionPattern := regexp.MustCompile(centosMirrorsDistroVersionRegex)
	distroVersionsURLs := []string{}

	allowedDomains, err := getHostnamesFromURLs(mirrorsConfig.URLs)
	if err != nil {
		return nil, fmt.Errorf("Error while getting domains from mirrors root URLs: %s", mirrorsConfig.URLs)
	}

	co := colly.NewCollector(
		colly.AllowedDomains(allowedDomains...),
	)

	co.OnHTML("a[href]", func(e *colly.HTMLElement) {
		distroVersionFolderMatch := centosMirrorsDistroVersionPattern.FindStringSubmatch(e.Attr("href"))

		if len(distroVersionFolderMatch) > 0 {
			//nolint:errcheck
			co.Visit(e.Request.AbsoluteURL(e.Attr("href")))
		}
	})

	co.OnRequest(func(r *colly.Request) {

		if !sliceContains(mirrorsConfig.URLs, r.URL.String()) {
			distroVersionsURLs = append(distroVersionsURLs, r.URL.String())
		}
	})

	for _, mirrorRootURL := range mirrorsConfig.URLs {
		err := co.Visit(mirrorRootURL)
		if err != nil {
			return nil, err
		}
	}

	return distroVersionsURLs, nil
}

// Seek packages for each Distro version
func scrape(mirrorsConfig MirrorsConfig, versionRootURLs []string, packagePrefix string) ([]string, error) {
	var packages []string
	packageFilenameRegex := `^` + packagePrefix + `.+.` + centosPackageFileExtension

	allowedDomains, err := getHostnamesFromURLs(mirrorsConfig.URLs)
	if err != nil {
		return nil, fmt.Errorf("Error while getting domains from mirrors root URLs: %s", mirrorsConfig.URLs)
	}

	co := colly.NewCollector(
		colly.AllowedDomains(allowedDomains...),
	)

	co.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")

		// Do not traverse the hierarchy in inverse order
		if !(strings.Contains(link, "../")) {
			// TODO: implement a retry logic
			//nolint:errcheck
			co.Visit(e.Request.AbsoluteURL(link))
		}
	})

	co.OnRequest(func(r *colly.Request) {
		folderPattern := regexp.MustCompile(`.+\/$`)
		folderMatch := folderPattern.FindStringSubmatch(r.URL.String())

		// If the URL is of a folder
		if len(folderMatch) <= 0 {
			packagePattern := regexp.MustCompile(fmt.Sprintf(`.+\.%s$`, centosPackageFileExtension))
			packageMatch := packagePattern.FindStringSubmatch(r.URL.String())

			// If the URL is of a package file
			if len(packageMatch) > 0 {
				packageName := path.Base(r.URL.String())
				packageNamePattern := regexp.MustCompile(packageFilenameRegex)
				packageNameMatch := packageNamePattern.FindStringSubmatch(packageName)

				// If the URL matches the package filter regex
				if len(packageNameMatch) > 0 {
					packages = append(packages, packageName)
				}
			}

			r.Abort()
		}
	})

	packagesURIs := []string{}

	for _, arch := range mirrorsConfig.Archs {
		for _, uriFormat := range mirrorsConfig.PackagesURIFormats {
			packagesURIs = append(packagesURIs, fmt.Sprintf(uriFormat, arch))
		}
	}

	for _, versionRootURL := range versionRootURLs {
		for _, packagesURI := range packagesURIs {
			packagesURL := versionRootURL + packagesURI

			//nolint:errcheck
			co.Visit(packagesURL)
		}
	}

	return packages, nil
}
