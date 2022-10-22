package amazonlinux2

import (
	"context"
	"fmt"
	"github.com/maxgio92/krawler/pkg/distro"
	p "github.com/maxgio92/krawler/pkg/packages"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func (a *AmazonLinux2) dereferenceRepositoryURL(src *url.URL, arch distro.Arch) (*url.URL, error) {
	var dest *url.URL

	mirrorListURL, err := url.JoinPath(src.String(), string(arch), "mirror.list")
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, mirrorListURL, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Amazon Linux 2 repository URL not valid to be dereferenced")
	}

	if resp.Body == nil {
		return nil, fmt.Errorf("empty response from Amazon Linux 2 repository reference URL")
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Get first repository URL available, no matter what the geolocation.
	s := strings.Split(string(b), "\n")[0]

	dest, err = url.Parse(s)
	if err != nil {
		return nil, err
	}

	return dest, nil
}

// Returns the list of default repositories from the default config.
func (a *AmazonLinux2) getDefaultRepositories() []p.Repository {
	var repositories []p.Repository

	for _, repository := range DefaultConfig.Repositories {
		if !distro.RepositorySliceContains(repositories, repository) {
			repositories = append(repositories, repository)
		}
	}

	return repositories
}
