package scraper

import (
	"fmt"
	"net/url"
)

func sliceContains(s []string, e string) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}

	return false
}

func getHostnamesFromURLs(urls []string) ([]string, error) {
	hostnames := []string{}

	for _, v := range urls {
		url, err := url.Parse(v)
		if err != nil {
			return nil, fmt.Errorf("No hostnames found in URLs: %s", urls)
		}

		hostnames = append(hostnames, url.Host)
	}

	return hostnames, nil
}
