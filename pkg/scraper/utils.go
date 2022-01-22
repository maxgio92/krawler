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

func getHostnamesFromURLs(URLs []string) ([]string, error) {
	hostnames := []string{}

	for _, v := range URLs {
		url, err := url.Parse(v)
		if err != nil {
			return nil, fmt.Errorf("No hostnames found in URLs: %s", URLs)
		}
		hostnames = append(hostnames, url.Host)
	}
	return hostnames, nil
}
