package scrape

import (
	"fmt"
	"net/url"
	"path"
	"regexp"
	"strings"

	"github.com/gocolly/colly"
	d "github.com/gocolly/colly/debug"
)

// Returns a list of file names found from the seed URL, filtered by file name regex.
//
//nolint:funlen
func CrawlFiles(seedURLs []*url.URL, exactFileRegex string, debug bool) ([]string, error) {
	var files []string

	folderPattern := regexp.MustCompile(folderRegex)

	exactFilePattern := regexp.MustCompile(exactFileRegex)

	fileRegex := strings.TrimPrefix(exactFileRegex, "^")
	filePattern := regexp.MustCompile(fileRegex)

	allowedDomains := getHostnamesFromURLs(seedURLs)

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

	// Add the analysis callback to find file URLs, for each Visit call
	co.OnRequest(func(r *colly.Request) {
		folderMatch := folderPattern.FindStringSubmatch(r.URL.String())

		// If the URL is not of a folder.
		if len(folderMatch) == 0 {
			fileMatch := filePattern.FindStringSubmatch(r.URL.String())

			// If the URL is of a file file.
			if len(fileMatch) > 0 {
				fileName := path.Base(r.URL.String())
				fileNameMatch := exactFilePattern.FindStringSubmatch(fileName)

				// If the URL matches the file filter regex.
				if len(fileNameMatch) > 0 {
					files = append(files, fileName)
				}
			}

			// Otherwise abort the request.
			r.Abort()
		}
	})

	// Visit each mirror root folder.
	for _, seedURL := range seedURLs {
		//nolint:errcheck
		co.Visit(seedURL.String())
	}

	return files, nil
}

// Returns a list of folder names found from each seed URL, filtered by folder name regex.
func CrawlFolders(seedURLs []*url.URL, regex string, debug bool) ([]string, error) {
	var versions []string

	versionPattern := regexp.MustCompile(regex)

	allowedDomains := getHostnamesFromURLs(seedURLs)
	if len(allowedDomains) < 1 {
		//nolint:goerr113
		return nil, fmt.Errorf("invalid seed urls")
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

	// Visit each distro version-specific folder.
	co.OnHTML("a[href]", func(e *colly.HTMLElement) {
		distroVersionFolderMatch := versionPattern.FindStringSubmatch(e.Attr("href"))

		if len(distroVersionFolderMatch) > 0 {
			//nolint:errcheck
			co.Visit(e.Request.AbsoluteURL(e.Attr("href")))
		}
	})

	// Collect all the version folder names.
	co.OnRequest(func(r *colly.Request) {
		if !urlSliceContains(seedURLs, r.URL) {
			versions = append(versions, path.Base(r.URL.Path))
		}
	})

	// Visit each mirror root folder.
	for _, seedURL := range seedURLs {
		err := co.Visit(seedURL.String())
		if err != nil {
			return nil, err
		}
	}

	return versions, nil
}
