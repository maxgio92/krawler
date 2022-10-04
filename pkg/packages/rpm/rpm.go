package rpm

import (
	"bytes"
	"compress/gzip"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	u "net/url"
	"path/filepath"
	"sync"

	"github.com/antchfx/xmlquery"
	rpmutils "github.com/sassoftware/go-rpmutils"
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
	packagesCh := make(chan Package)
	defer close(packagesCh)

	repoWorkers := sync.WaitGroup{}
	repoWorkers.Add(len(repositoryURLs))

	errCh := make(chan error)
	defer close(errCh)

	for _, r := range repositoryURLs {

		repoURL := r.String()
		go func() {
			metadataURL, err := u.JoinPath(repoURL, metadataURI)
			if err != nil {
				errCh <- err
			} else {
				logger.WithField("url", repoURL).Info("Analysing repository")

				dbs, _ := getDBsFromMetadataURL(metadataURL)

				for _, db := range dbs {
					logger.WithField("type", db.Type).Info("Analysing DB")

					getPackagesFromDB(packagesCh, repoURL, db.GetLocation(), packageName, packageFileNames...)
				}
			}
			repoWorkers.Done()
		}()
	}

	// Acquire packages.
	go func() {
		for p := range packagesCh {
			packages = append(packages, p)
		}
	}()

	go func() {
		for e := range errCh {
			logger.WithError(e).Debug()
		}
	}()

	repoWorkers.Wait()

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
		return nil, fmt.Errorf("metadata url not found")
	}
	if resp.Body == nil {
		return nil, fmt.Errorf("metadata url returned an invalid response")
	}
	defer resp.Body.Close()

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

func getPackagesFromDB(packages chan<- Package, repoURL string, dbURI string, packageName string, fileNames ...string) error {
	return getPackagesFromXMLDB(packages, repoURL, dbURI, packageName, fileNames...)
}

func getPackagesFromXMLDB(packages chan<- Package, repoURL string, dbURI string, packageName string, fileNames ...string) (err error) {
	//nolint:typecheck
	dbURL, err := u.JoinPath(repoURL, dbURI)
	if err != nil {
		return err
	}

	u, err := u.Parse(dbURL)
	if err != nil {
		return err
	}

	logger.WithField("url", u.String()).Debug("Downloading DB")

	resp, err := http.Get(u.String())
	if err != nil {
		return err
	}
	if resp.Body == nil {
		return fmt.Errorf("repository url not found")
	}
	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("repository url returned an invalid response")
	}
	defer resp.Body.Close()

	gr, err := gzip.NewReader(resp.Body)
	defer gr.Close()
	if err != nil {
		return err
	}

	logger.WithField("uri", filepath.Base(dbURI)).Debug("Parsing DB")

	doc, err := xmlquery.Parse(gr)
	if err != nil {
		return err
	}

	logger.WithField("package", packageName).Info("Querying DB")

	packagesXML, err := xmlquery.QueryAll(doc, "//package[name='"+packageName+"']")
	if err != nil {
		return err
	}

	err = buildPackagesFromXML(packages, packagesXML, repoURL, fileNames...)
	if err != nil {
		return err
	}

	return nil
}

func buildPackagesFromXML(packages chan<- Package, nodes []*xmlquery.Node, repositoryURL string, fileNames ...string) error {
	pkgWorkers := sync.WaitGroup{}
	pkgWorkers.Add(len(nodes))

	errCh := make(chan error)
	defer close(errCh)

	for _, node := range nodes {

		v := node

		go func() {
			p := &Package{}

			err := xml.Unmarshal([]byte(v.OutputXML(true)), p)
			if err != nil {
				errCh <- err
			}

			p.url = repositoryURL + p.GetLocation()

			logger.WithField("fullname", filepath.Base(p.GetLocation())).Info("Opening package")

			fileReaders, err := getFileReadersFromPackageURL(p.url, fileNames...)
			if err != nil {
				errCh <- err
			}
			p.fileReaders = fileReaders

			packages <- *p

			pkgWorkers.Done()
		}()
	}

	pkgWorkers.Wait()

	return nil
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
	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("package url not found")
	}
	if resp.Body == nil {
		return nil, fmt.Errorf("packge url returned an invalid response")
	}
	defer resp.Body.Close()

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

			var buf bytes.Buffer
			_, err = io.Copy(&buf, payload)
			//buf, err := io.ReadAll(payload)
			if err != nil {
				return nil, err
			}
			//r := bytes.NewReader(buf)
			r := bytes.NewReader(buf.Bytes())
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
