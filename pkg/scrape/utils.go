package scrape

import "net/url"

func getHostnamesFromURLs(urls []string) ([]string, error) {
	hostnames := []string{}

	for _, v := range urls {
		url, err := url.Parse(v)
		if err != nil {
			return nil, errDomainsFromSeedUrls
		}

		hostnames = append(hostnames, url.Host)
	}

	return hostnames, nil
}
