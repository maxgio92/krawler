package rpm

import (
	"compress/gzip"
	"encoding/xml"
	"net/http"
	u "net/url"

	"github.com/antchfx/xmlquery"
	"github.com/maxgio92/krawler/pkg/distro"
)

const (
	repoMetadataURI = "repodata/repomd.xml"
)

func GetPackagesFromRepositories(repositoryURLs []*u.URL, name string, debug bool) ([]distro.Package, error) {
	var packages []distro.Package

	for _, repoURL := range repositoryURLs {
		DBs, err := getDBsFromRepoMetaDataURL(repoURL.String() + "/" + repoMetadataURI)
		if err != nil {
			return nil, err
		}

		for _, v := range DBs {
			DBURL := repoURL.String() + v.Location.Href
			rpmPackages, err := getPackagesFromRepoXMLDBURL(DBURL, name)
			if err != nil {
				return nil, err
			}

			for _, v := range rpmPackages {
				packages = append(packages, distro.Package(v.Name+"-"+v.Version.Ver+"-"+v.Version.Rel+"."+v.Arch+".rpm"))
			}
		}
	}

	return packages, nil
}

func getDBsFromRepoMetaDataURL(url string) (DBs []Data, err error) {
	u, err := u.Parse(url)
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}

	doc, err := xmlquery.Parse(resp.Body)
	if err != nil {
		return
	}

	DatasXML, err := xmlquery.QueryAll(doc, "//repomd/data")
	if err != nil {
		return
	}

	for _, v := range DatasXML {
		data := &Data{}
		err = xml.Unmarshal([]byte(v.OutputXML(true)), data)
		if err != nil {
			return
		}
		switch data.Type {
		case "primary":
			DBs = append(DBs, *data)
		default:
			break
		}
	}

	return
}

func getPackagesFromRepoXMLDBURL(us string, packageName string) (packages []Package, err error) {
	u, err := u.Parse(us)
	if err != nil {
		return nil, err
	}

	gr, err := getGzipReaderFromURL(u.String())
	if err != nil {
		return
	}

	doc, err := xmlquery.Parse(gr)
	if err != nil {
		return
	}

	packagesXML, err := xmlquery.QueryAll(doc, "//package[name='"+packageName+"']")
	if err != nil {
		return
	}

	for _, v := range packagesXML {
		p := &Package{}
		err = xml.Unmarshal([]byte(v.OutputXML(true)), p)
		if err != nil {
			return
		}
		packages = append(packages, *p)

	}
	return
}

func getGzipReaderFromURL(us string) (*gzip.Reader, error) {
	u, err := u.Parse(us)
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
