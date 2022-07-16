package scrape

import (
	"net/url"
	"path"
	"regexp"
	"strings"

	"github.com/gocolly/colly"
	d "github.com/gocolly/colly/debug"
)

// Returns a list of files found on each page URL, filtered by filename regex.
//nolint:funlen,revive,stylecheck
func CrawlFiles(seedUrl *url.URL, exactFileRegex string, debug bool) ([]string, error) {
	var files []string

	folderPattern := regexp.MustCompile(folderRegex)

	exactFilePattern := regexp.MustCompile(exactFileRegex)

	fileRegex := strings.TrimPrefix(exactFileRegex, "^")
	filePattern := regexp.MustCompile(fileRegex)

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

	//nolint:errcheck
	co.Visit(seedUrl.String())

	return files, nil
}
