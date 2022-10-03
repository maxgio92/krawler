package rpm

import (
	"bytes"
	"encoding/xml"
	"io"
	"net/http"
	u "net/url"
	"path/filepath"

	"github.com/antchfx/xmlquery"
	"github.com/sassoftware/go-rpmutils"
	log "github.com/sirupsen/logrus"
)

const (
	metadataURI = "repodata/repomd.xml"
)

var (
	logger = log.New()
)

func init() {
	logger.SetLevel(log.DebugLevel)
	logger.SetFormatter(&log.TextFormatter{
		ForceColors:      true,
		DisableTimestamp: true,
	})
}

// GetPackagesFromRepositories returns a list of package of type Package with specified name,
// searching in the specified repositories.
func GetPackagesFromRepositories(repositoryURLs []*u.URL, packageName string, packageFileNames ...string) ([]Package, error) {
	var packages []Package

	for _, repoURL := range repositoryURLs {
		metadataURL, err := u.JoinPath(repoURL.String(), metadataURI)
		if err != nil {
			return nil, err
		}

		logger.WithField("url", repoURL).Info("Analysing repository")

		dbs, _ := getDBsFromMetadataURL(metadataURL)

		for _, db := range dbs {
			logger.WithField("type", db.Type).Info("Analysing DB")

			p, _ := getPackagesFromDB(repoURL.String(), db.GetLocation(), packageName, packageFileNames...)
			packages = append(packages, p...)
		}
	}

	return packages, nil
}

func getDBsFromMetadataURL(url string) (dbs []Data, err error) {
	u, err := u.Parse(url)
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusNotFound {
		logger.Debug("skipping")
		return
	}

	logger.Debug("Parsing repository metadata")

	doc, err := xmlquery.Parse(resp.Body)
	if err != nil {
		return
	}

	logger.Debug("Getting repository DBs")

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
			dbs = append(dbs, *data)
		default:
		}
	}

	//nolint:nakedret
	return
}

func getPackagesFromDB(repoURL string, dbURI string, packageName string, fileNames ...string) (packages []Package, err error) {
	return getPackagesFromXMLDB(repoURL, dbURI, packageName, fileNames...)
}

func getPackagesFromXMLDB(repoURL string, dbURI string, packageName string, fileNames ...string) (packages []Package, err error) {
	dbURL, err := u.JoinPath(repoURL, dbURI)
	if err != nil {
		return nil, err
	}

	u, err := u.Parse(dbURL)
	if err != nil {
		return nil, err
	}

	logger.WithField("url", u.String()).Debug("Downloading DB")

	gr, err := getGzipReaderFromURL(u.String())
	if err != nil {
		return nil, err
	}

	logger.WithField("uri", filepath.Base(dbURI)).Debug("Parsing DB")

	doc, err := xmlquery.Parse(gr)
	if err != nil {
		return nil, err
	}

	logger.WithField("package", packageName).Info("Querying DB")

	packagesXML, err := xmlquery.QueryAll(doc, "//package[name='"+packageName+"']")
	if err != nil {
		return nil, err
	}

	packages, err = buildPackagesFromXML(packagesXML, repoURL, fileNames...)
	if err != nil {
		return nil, err
	}
	//nolint:nakedret
	return
}

func buildPackagesFromXML(nodes []*xmlquery.Node, repositoryURL string, fileNames ...string) ([]Package, error) {
	packages := []Package{}

	for _, v := range nodes {
		p := &Package{}

		err := xml.Unmarshal([]byte(v.OutputXML(true)), p)
		if err != nil {
			return nil, err
		}

		p.url = repositoryURL + p.GetLocation()

		logger.WithField("fullname", filepath.Base(p.GetLocation())).Info("Opening package")

		fileReaders, err := getFileReadersFromPackageURL(p.url, fileNames...)
		if err != nil {
			return packages, err
		}
		p.fileReaders = fileReaders

		packages = append(packages, *p)
	}

	return packages, nil
}

func getFileReadersFromPackageURL(packageURL string, fileNames ...string) ([]io.Reader, error) {
	util, err := getRPMUtilFromPackageURL(packageURL)
	if err != nil {
		return nil, err
	}

	fileReaders, err := getFileReadersFromRPMUtil(util, fileNames...)
	if err != nil {
		return nil, err
	}

	return fileReaders, nil
}

func getRPMUtilFromPackageURL(packageURL string) (*rpmutils.Rpm, error) {
	u, err := u.Parse(packageURL)
	if err != nil {
		return nil, err
	}

	logger.WithField("url", u.String()).Debug("Downloading package")

	resp, err := http.Get(u.String())
	if err != nil {
		logger.WithError(err).Debug("Error downloading package")

		return nil, err
	}

	logger.Debug("Parsing package")

	rpm, err := rpmutils.ReadRpm(resp.Body)
	if err != nil {
		logger.WithError(err).Debug("Error parsing package")

		return nil, err
	}

	return rpm, nil
}

func getFileReadersFromRPMUtil(util *rpmutils.Rpm, names ...string) ([]io.Reader, error) {
	payload, err := util.PayloadReaderExtended()
	if err != nil {
		return nil, err
	}

	readers := []io.Reader{}

	fileName := ""
	limited := len(names) > 0

	for {
		fileInfo, err := payload.Next()
		if err == io.EOF {
			break
		}

		if limited {
			fileName = names[len(names)-1]
		}

		if fileName == "" || filepath.Base(fileInfo.Name()) == fileName {
			logger.WithField("name", fileName).Debug("Found file")

			buf, err := io.ReadAll(payload)
			if err != nil {
				return nil, err
			}
			r := bytes.NewReader(buf)
			r.Seek(0, io.SeekStart)

			readers = append(readers, r)

			if limited {
				names = names[:len(names)-1]
				if len(names) == 0 {
					return readers, nil
				}
			}
		}
	}

	return readers, nil
}
