package distro

import (
	"net/url"

	"github.com/maxgio92/krawler/pkg/packages"
)

func RepositorySliceContains(s []packages.Repository, e packages.Repository) bool {
	for _, v := range s {
		if v.URI == e.URI {
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
			return nil, ErrDomainsFromMirrorUrls
		}

		hostnames = append(hostnames, url.Host)
	}

	return hostnames, nil
}
