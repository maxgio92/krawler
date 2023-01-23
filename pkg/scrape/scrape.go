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

// CrawlFiles returns a list of file names found from the seed URL, filtered by file name regex.
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
		if !(strings.Contains(link, "../")) && link != "/" {
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

			// If the URL is of a file.
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
		err := co.Visit(seedURL.String())
		if err != nil {
			return nil, err
		}
	}

	return files, nil
}

// CrawlFolders returns a list of folder names found from each seed URL, filtered by folder name regex.
//
//nolint:funlen
func CrawlFolders(seedURLs []*url.URL, exactFolderRegex string, recursive bool, debug bool) ([]string, error) {
	var folders []string

	folderPattern := regexp.MustCompile(folderRegex)

	exactFolderPattern := regexp.MustCompile(exactFolderRegex)

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

	// Visit each specific folder.
	co.OnHTML("a[href]", func(e *colly.HTMLElement) {
		href := e.Attr("href")

		folderMatch := folderPattern.FindStringSubmatch(href)

		// if the URL is of a folder.
		//nolint:nestif
		if len(folderMatch) > 0 {

			// Do not traverse the hierarchy in reverse order.
			if strings.Contains(href, "../") || href == "/" {
				return
			}

			exactFolderMatch := exactFolderPattern.FindStringSubmatch(href)
			if len(exactFolderMatch) > 0 {

				hrefAbsURL, _ := url.Parse(e.Request.AbsoluteURL(href))
				if !urlSliceContains(seedURLs, hrefAbsURL) {

					folders = append(folders, path.Base(hrefAbsURL.Path))
				}
			}
			if recursive {
				//nolint:errcheck
				co.Visit(e.Request.AbsoluteURL(href))
			}
		}
	})

	co.OnRequest(func(r *colly.Request) {
		folderMatch := folderPattern.FindStringSubmatch(r.URL.String())

		// if the URL is not of a folder.
		if len(folderMatch) == 0 {
			r.Abort()
		}
	})

	// Visit each mirror root folder.
	for _, seedURL := range seedURLs {
		err := co.Visit(seedURL.String())
		if err != nil {
			return nil, err
		}
	}

	return folders, nil
}
