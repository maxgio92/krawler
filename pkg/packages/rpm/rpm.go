package rpm

import (
	"bytes"
	"compress/gzip"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"sync"

	"github.com/antchfx/xmlquery"
	rpmutils "github.com/sassoftware/go-rpmutils"
	log "github.com/sirupsen/logrus"
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
func GetPackagesFromRepositories(repositoryURLs []*url.URL, packageName string, packageFileNames ...string) ([]Package, error) {
	var packages []Package
	packagesCh := make(chan Package)
	defer close(packagesCh)

	repoWorkers := sync.WaitGroup{}
	repoWorkers.Add(len(repositoryURLs))

	errCh := make(chan error)
	defer close(errCh)

	for _, r := range repositoryURLs {
		repoURL := r
		go func() {
			defer repoWorkers.Done()

			err := GetPackagesFromRepository(packagesCh, repoURL, packageName, packageFileNames...)
			if err != nil {
				errCh <- err
				logger.WithError(err).WithField("repository", repoURL).Error("error found getting packages from repository")
				return
			}
		}()
	}

	// Acquire packages.
	go func() {
		for errCh != nil {
			select {
			case p, ok := <-packagesCh:
				if ok {
					packages = append(packages, p)
					logger.WithField("name", p.Name).WithField("version", p.Version.Ver).WithField("release", p.Version.Rel).Info("New package found")
				}
			case e, ok := <-errCh:
				if !ok {
					errCh = nil
					continue
				}
				logger.WithError(e).Debug()
			}
		}
	}()
	repoWorkers.Wait()

	return packages, nil
}

// GetPackagesFromRepository returns a list of package of type Package with specified name,
// searching in the specified repository.
func GetPackagesFromRepository(packagesCh chan Package, repositoryURL *url.URL, packageName string, packageFileNames ...string) error {
	repoURL := repositoryURL.String()

	metadataURL, err := url.JoinPath(repoURL, metadataPath)
	if err != nil {
		return err
	}
	logger.WithField("url", repoURL).Info("Analysing repository")

	dbs, err := getDBsFromMetadataURL(metadataURL)
	if err != nil {
		return err
	}
	for _, db := range dbs {
		logger.WithField("type", db.Type).Info("Analysing DB")

		err = getPackagesFromDB(packagesCh, repoURL, db.GetLocation(), packageName, packageFileNames...)
		if err != nil {
			return err
		}
	}
	return nil
}

func getDBsFromMetadataURL(metadataURL string) (dbs []Data, err error) {
	u, err := url.Parse(metadataURL)
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("metadata url %s not valid", metadataURL)
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

	datasXML, err := xmlquery.QueryAll(doc, metadataDataXPath)
	if err != nil {
		return
	}

	for _, v := range datasXML {
		data := &Data{}

		err = xml.Unmarshal([]byte(v.OutputXML(true)), data)
		if err != nil {
			return
		}

		switch data.Type {
		case primary:
			dbs = append(dbs, *data)
		default:
		}
	}

	//nolint:nakedret
	return
}

func getPackagesFromDB(packagesCh chan Package, repoURL string, dbURI string, packageName string, fileNames ...string) error {
	return getPackagesFromXMLDB(packagesCh, repoURL, dbURI, packageName, fileNames...)
}

func getPackagesFromXMLDB(packagesCh chan Package, repoURL string, dbURI string, packageName string, fileNames ...string) (err error) {
	//nolint:typecheck
	dbURL, err := url.JoinPath(repoURL, dbURI)
	if err != nil {
		return err
	}

	u, err := url.Parse(dbURL)
	if err != nil {
		return err
	}
	logger.WithField("url", u.String()).Debug("Downloading DB")

	resp, err := http.Get(u.String())
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("repository url %s not valid", u.String())
	}
	if resp.Body == nil {
		return fmt.Errorf("repository url returned an invalid response")
	}
	defer resp.Body.Close()

	gr, err := gzip.NewReader(resp.Body)
	if err != nil {
		return err
	}
	defer gr.Close()
	logger.WithField("uri", filepath.Base(dbURI)).Debug("Parsing DB")

	var packagesXML []*xmlquery.Node

	sp, err := xmlquery.CreateStreamParser(gr, dataPackageXPath, dataPackageXPath+"[name='"+packageName+"']")
	if err != nil {
		return err
	}
	logger.WithField("package", packageName).Debug("Querying DB")
	for {
		n, err := sp.Read()
		if err != nil {
			break
		}
		packagesXML = append(packagesXML, n)
	}

	buildPackagesFromXMLNodes(packagesCh, packagesXML, repoURL, fileNames...)

	return nil
}

func buildPackagesFromXMLNodes(packages chan Package, nodes []*xmlquery.Node, repositoryURL string, fileNames ...string) {
	pkgWorkers := sync.WaitGroup{}
	pkgWorkers.Add(len(nodes))

	errCh := make(chan error)
	defer close(errCh)

	for _, v := range nodes {
		node := v

		go func() {
			defer pkgWorkers.Done()

			p, err := buildPackageFromXML(node, repositoryURL, fileNames...)
			if err != nil {
				errCh <- err
				return
			}
			packages <- *p
		}()
	}

	go func() {
		for err := range errCh {
			logger.WithError(err).Error("Error found building package from XML")
		}
	}()

	pkgWorkers.Wait()
}

func buildPackageFromXML(node *xmlquery.Node, repositoryURL string, fileNames ...string) (*Package, error) {
	p := &Package{}

	err := xml.Unmarshal([]byte(node.OutputXML(true)), p)
	if err != nil {
		return nil, err
	}

	p.url, err = url.JoinPath(repositoryURL, p.GetLocation())
	if err != nil {
		return nil, err
	}

	logger.WithField("fullname", filepath.Base(p.GetLocation())).Debug("Opening package")

	fileReaders, err := getFileReadersFromPackageURL(p.url, fileNames...)
	if err != nil {
		return nil, err
	}
	p.fileReaders = fileReaders

	return p, nil
}

func getFileReadersFromPackageURL(packageURL string, fileNames ...string) ([]io.Reader, error) {
	u, err := url.Parse(packageURL)
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

	rpm, err := rpmutils.ReadRpm(resp.Body)
	if err != nil {
		logger.WithError(err).Debug("Error parsing package")
		return nil, err
	}

	fileReaders, err := getFileReadersFromRPMUtil(rpm, fileNames...)
	if err != nil {
		return nil, err
	}

	return fileReaders, nil
}

func getFileReadersFromRPMUtil(util *rpmutils.Rpm, names ...string) ([]io.Reader, error) {
	payload, err := util.PayloadReaderExtended()
	if err != nil {
		return nil, err
	}

	var readers []io.Reader

	fileName := ""
	limited := len(names) > 0

	logger.WithField("files", names).Debug("Looking for files")

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
			if err != nil {
				return nil, err
			}
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
