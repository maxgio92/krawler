package rpm

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/xml"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"sync"

	"github.com/antchfx/xmlquery"
	"github.com/pkg/errors"
	rpmutils "github.com/sassoftware/go-rpmutils"
	log "github.com/sirupsen/logrus"
)

func init() {
	logger.SetLevel(log.InfoLevel)
	logger.SetFormatter(&log.TextFormatter{
		ForceColors:      true,
		DisableTimestamp: true,
	})
}

// SearchPackages crawls packages from the specified repositories,
// and returns a list of package of type Package with specified name.
func SearchPackages(repositoryURLs []*url.URL, packageName string, packageFileNames ...string) ([]Package, error) {
	var packages []Package

	packagesCh := make(chan Package)

	producersWG := sync.WaitGroup{}
	producersWG.Add(len(repositoryURLs))

	errCh := make(chan error)

	consumersDoneCh := make(chan bool, 1)

	// Run parallel packages producer workers.
	for _, r := range repositoryURLs {
		repoURL := r

		go producePackagesFromRepository(&producersWG, packagesCh, errCh, repoURL, packageName, packageFileNames...)
	}

	// Consume packages.
	go consumePackages(consumersDoneCh, &packages, packagesCh, errCh)

	// Wait for producers.
	producersWG.Wait()
	close(packagesCh)
	close(errCh)

	// Wait for consumers.
	<-consumersDoneCh

	return packages, nil
}

// producePackagesFromRepository crawls packages from specified repository as repositoryURL *URL,
// sends them over a packagesCh Package channel and signals completion through a WaitGroup.
// The waitGroup counter needs to be greater than zero.
func producePackagesFromRepository(waitGroup *sync.WaitGroup, packagesCh chan Package, errCh chan error, repositoryURL *url.URL, packageName string, packageFileNames ...string) {
	defer waitGroup.Done()

	repoURL := repositoryURL.String()

	metadataURL, err := url.JoinPath(repoURL, metadataPath)
	if err != nil {
		errCh <- err
	}

	logger.WithField("url", repoURL).Info("Analysing repository")

	dbs, err := getDBsFromMetadataURL(metadataURL)
	if err != nil {
		errCh <- err
	}

	for _, db := range dbs {
		logger.WithField("type", db.Type).Info("Analysing DB")

		err = getPackagesFromDB(packagesCh, repoURL, db.GetLocation(), packageName, packageFileNames...)
		if err != nil {
			errCh <- err
		}
	}
}

// consumePackages receives Package over a packagesCh Package channel,
// and signals completion through a bool channel done.
func consumePackages(done chan bool, packages *[]Package, packagesCh chan Package, errCh chan error) {
	for errCh != nil || packagesCh != nil {
		select {
		case p, ok := <-packagesCh:
			if ok {
				*packages = append(*packages, p)
				logger.WithField("name", p.Name).WithField("version", p.Version.Ver).WithField("release", p.Version.Rel).Info("New package found")

				continue
			}

			packagesCh = nil
		case e, ok := <-errCh:
			if ok {
				logger.WithError(e).Error()

				continue
			}

			errCh = nil
		}
	}
	done <- true
}

//nolint:cyclop
func getDBsFromMetadataURL(metadataURL string) ([]Data, error) {
	var dbs []Data

	u, err := url.Parse(metadataURL)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.Wrap(errMetadataURLNotValid, metadataURL)
	}

	if resp.Body == nil {
		return nil, errMetadataInvalidResponse
	}
	defer resp.Body.Close()

	logger.Debug("Parsing repository metadata")

	doc, err := xmlquery.Parse(resp.Body)
	if err != nil {
		return nil, err
	}

	logger.Debug("Getting repository DBs")

	datasXML, err := xmlquery.QueryAll(doc, metadataDataXPath)
	if err != nil {
		return nil, err
	}

	for _, v := range datasXML {
		data := &Data{}

		err = xml.Unmarshal([]byte(v.OutputXML(true)), data)
		if err != nil {
			return nil, err
		}

		switch data.Type {
		case primary:
			dbs = append(dbs, *data)
		default:
		}
	}

	return dbs, nil
}

func getPackagesFromDB(packagesCh chan Package, repoURL string, dbURI string, packageName string, fileNames ...string) error {
	return getPackagesFromXMLDB(packagesCh, repoURL, dbURI, packageName, fileNames...)
}

//nolint:cyclop
func getPackagesFromXMLDB(packagesCh chan Package, repoURL string, dbURI string, packageName string, fileNames ...string) (err error) {
	dbURL, err := url.JoinPath(repoURL, dbURI)
	if err != nil {
		return err
	}

	u, err := url.Parse(dbURL)
	if err != nil {
		return err
	}

	logger.WithField("url", u.String()).Debug("Downloading DB")

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, u.String(), nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.Wrap(errRepositoryURLNotValid, u.String())
	}

	if resp.Body == nil {
		return errRepositoryInvalidResponse
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

	done := make(chan bool, 1)

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
		done <- true
	}()

	pkgWorkers.Wait()
	close(errCh)

	<-done
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

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.WithError(err).Debug("Error downloading package")

		return nil, err
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, errPackageURLNotFound
	}

	if resp.Body == nil {
		return nil, errPackageURLInvalidResponse
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

//nolint:cyclop
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
		if errors.Is(err, io.EOF) {
			break
		}

		if limited {
			fileName = names[len(names)-1]
		}

		//nolint:nestif
		if fileName == "" || filepath.Base(fileInfo.Name()) == fileName {
			logger.WithField("name", fileName).Debug("Found file")

			var buf bytes.Buffer

			_, err = io.Copy(&buf, payload)
			if err != nil {
				return nil, err
			}

			r := bytes.NewReader(buf.Bytes())

			_, err := r.Seek(0, io.SeekStart)
			if err != nil {
				return nil, err
			}

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
