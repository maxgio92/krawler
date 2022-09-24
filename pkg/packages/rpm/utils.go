package rpm

import (
	"compress/gzip"
	"net/http"
	u "net/url"
)

func getGzipReaderFromURL(url string) (*gzip.Reader, error) {
	u, err := u.Parse(url)
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}

	r, err := gzip.NewReader(resp.Body)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	return r, nil
}
