package scraper

import "fmt"

var ScraperByDistro = map[Distro]Scraper{}

type Scraper interface {
	Scrape(MirrorsConfig, string) ([]string, error)
}

type Distro int16

const (
	Unknown Distro = iota
	Centos
)

func Factory(distro Distro) (Scraper, error) {
	scraper, ok := ScraperByDistro[distro]

	if !ok {
		return nil, fmt.Errorf("No scraper found for Linux distribution: %s", distro)
	}

	return scraper, nil
}
