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

	"github.com/maxgio92/krawler/pkg/packages"

	"github.com/antchfx/xmlquery"
	"github.com/pkg/errors"
	rpmutils "github.com/sassoftware/go-rpmutils"
)

// SearchPackages crawls packages from the specified repositories,
// and returns a list of package of type Package with specified name.
// func SearchPackages(repositoryURLs []*url.URL, packageName string, packageFileNames ...string) ([]Package, error) {
func SearchPackages(so *SearchOptions) ([]packages.Package, error) {
	var result []packages.Package

	search := func(repoURL string) {
		searchPackagesFromRepository(
			func() {
				so.Progress(1)
				so.SigProducerCompletion()
			},
			so, repoURL)
	}

	collect := func() {
		so.Consume(
			func(p ...packages.Package) {
				so.Log().Debug("Scanned DB")
				if len(p) > 0 {
					result = append(result, p...)
					so.Log().Infof("New %d packages found", len(p))
				}
			},
			func(e error) {
				so.Log().Error(e)
			},
		)
	}

	// Run search producers.
	for _, v := range so.SeedURLs() {
		repoURL := v
		go search(repoURL)
	}

	// Run collect consumer.
	go collect()

	// Wait for producers and consumers to complete and cleanup.
	so.WaitAndClose()

	return result, nil
}

// searchPackagesFromRepository crawls packages from specified repository as repositoryURL *URL,
// sends them over a packagesCh Package channel and signals completion through a WaitGroup.
// The waitGroup counter needs to be greater than zero.
func searchPackagesFromRepository(doneFunc func(), so *SearchOptions, repoURL string) {
	defer doneFunc()

	metadataURL, err := url.JoinPath(repoURL, metadataPath)
	if err != nil {
		so.SendError(err)

		return
	}

	so.Log().WithField("url", repoURL).Info("Analysing repository")

	dbs, err := getPrimaryDBsFromMetadataURL(metadataURL)
	if err != nil {
		so.SendError(err)

		return
	}

	dbURLs := []string{}
	for _, db := range dbs {
		dbURL, _ := url.JoinPath(repoURL, db.GetLocation())
		dbURLs = append(dbURLs, dbURL)
	}

	for _, dbURL := range dbURLs {
		so.Log().WithField("url", dbURL).Info("Analysing DB")
		searchPackagesFromDB(so, repoURL, dbURL)
	}
}

//nolint:cyclop
func getPrimaryDBsFromMetadataURL(metadataURL string) ([]Data, error) {
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

// func searchPackagesFromDB(doneFunc func(), so *SearchOptions, repoURL, dbURL string) {
func searchPackagesFromDB(so *SearchOptions, repoURL, dbURL string) {
	xmlDB, err := getPackagesXMLDBFromURL(so, dbURL)
	if err != nil {
		so.SendError(err)
	}

	queue := packages.NewMPSCQueue(len(xmlDB))

	for _, v := range xmlDB {
		node := v

		go func() {
			defer queue.SigProducerCompletion()

			p := &Package{}

			err := xml.Unmarshal([]byte(node.OutputXML(true)), p)
			if err != nil {
				queue.SendError(err)

				return
			}

			p.url, err = url.JoinPath(repoURL, p.GetLocation())
			if err != nil {
				queue.SendError(err)

				return
			}

			so.Log().WithField("fullname", filepath.Base(p.GetLocation())).Debug("Opening package")

			fileReaders, err := getFileReadersFromPackageURL(p.url, so.PackageFileNames()...)
			if err != nil {
				queue.SendError(err)

				return
			}

			p.fileReaders = fileReaders

			so.Log().WithField("version", p.Version.Ver).WithField("release", p.Version.Rel).WithField("name", p.Name).Debug("found package")
			queue.SendMessage(p)
		}()
	}

	go func() {
		queue.Consume(
			func(p ...packages.Package) {
				so.SendMessage(p...)
			},
			func(e error) {
				so.Log().Error(e)
			},
		)
	}()

	queue.WaitAndClose()
}

func getPackagesXMLDBFromURL(so *SearchOptions, dbURL string) ([]*xmlquery.Node, error) {
	u, err := url.Parse(dbURL)
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
		return nil, errors.Wrap(errRepositoryURLNotValid, u.String())
	}

	if resp.Body == nil {
		return nil, errRepositoryInvalidResponse
	}
	defer resp.Body.Close()

	gr, err := gzip.NewReader(resp.Body)
	if err != nil {
		return nil, err
	}
	defer gr.Close()

	var packagesXML []*xmlquery.Node

	sp, err := xmlquery.CreateStreamParser(
		gr,
		dataPackageXPath,
		dataPackageXPath+"[name='"+so.PackageName()+"']")
	if err != nil {
		return nil, err
	}

	for {
		n, err := sp.Read()
		if err != nil {
			break
		}

		packagesXML = append(packagesXML, n)
	}

	return packagesXML, nil
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
