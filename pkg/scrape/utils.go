package scrape

import "net/url"

//nolint:unparam
func getHostnamesFromURLs(urls []*url.URL) ([]string, error) {
	hostnames := []string{}

	for _, v := range urls {
		hostnames = append(hostnames, v.Host)
	}

	return hostnames, nil
}

func urlSliceContains(us []*url.URL, u *url.URL) bool {
	for _, v := range us {
		if v == u {
			return true
		}
	}

	return false
}
