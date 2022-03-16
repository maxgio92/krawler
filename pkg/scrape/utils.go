package scrape

import (
	"net/url"
)

func repositorySliceContains(s []Repository, e Repository) bool {
	for _, v := range s {
		if v.PackagesURIFormat == e.PackagesURIFormat {
			return true
		}
	}

	return false
}

func stringSliceContains(s []string, e string) bool {
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
			return nil, errDomainsFromMirrorUrls
		}

		hostnames = append(hostnames, url.Host)
	}

	return hostnames, nil
}
