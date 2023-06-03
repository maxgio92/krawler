package alpm

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/Jguer/go-alpm/v2"
	"github.com/pkg/errors"
	"github.com/spf13/afero"

	"github.com/maxgio92/krawler/pkg/packages"
)

type Package struct {
	Name         string
	Version      string
	Release      string
	Architecture string
	Location     string
	url          string
	fileReaders  []io.Reader
	Files        []string
}

func (p *Package) GetName() string          { return p.Name }
func (p *Package) GetVersion() string       { return p.Version }
func (p *Package) GetRelease() string       { return p.Release }
func (p *Package) GetArch() string          { return p.Architecture }
func (p *Package) GetLocation() string      { return p.Location }
func (p *Package) URL() string              { return p.url }
func (p *Package) FileReaders() []io.Reader { return p.fileReaders }

const (
	root              = "/"
	ALPMDBVersionFile = "ALPM_DB_VERSION"
	ALPMDBVersion     = 9
)

func SearchPackages(so *SearchOptions) ([]packages.Package, error) {
	var result []packages.Package

	search := func(dbURL string) {
		searchPackagesFromDB(
			func() {
				so.Progress(1)
				so.SigProducerCompletion()
			},
			so, dbURL)
	}

	collect := func() {
		so.Consume(
			func(p ...packages.Package) {
				so.Log().Info("scanned db")
				if len(p) > 0 {
					result = append(result, p...)
					so.Log().Infof("new %d packages found", len(p))
				}
			},
			func(e error) {
				so.Log().Error(e)
			},
		)
	}

	// Run search producers.
	for _, v := range so.SeedURLs() {
		dbURL := v
		go search(dbURL)
	}

	// Run collect consumer.
	go collect()

	// Wait for producers and consumers to complete and cleanup.
	so.WaitAndClose()

	return result, nil
}

func searchPackagesFromDB(doneFunc func(), so *SearchOptions, dbURL string) {
	defer doneFunc()

	p, err := doSearchPackagesFromDB(dbURL, so.PackageNames())
	if err != nil {
		so.SendError(errors.Wrap(err, "searching packages from db"))
	}
	so.SendMessage(p...)
}

// doSearchPackagesFromDB looks for the package of which the specified package names, parsing the remote
// repository DB, and returns a slice of packages.Package.
// It possibly returns an error.
func doSearchPackagesFromDB(dbURL string, packageNames []string) ([]packages.Package, error) {
	fs := afero.NewOsFs()

	tmpdir, err := afero.TempDir(fs, os.TempDir(), "krawler")
	if err != nil {
		return nil, errors.Wrap(err, "error creating local DB temporeary directory")
	}

	req, err := http.NewRequest(http.MethodGet, dbURL, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error creating HTTP request")
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "error doing HTTP request")
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, nil
	}

	gzr, err := gzip.NewReader(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "error reading gzip response")
	}
	defer gzr.Close()

	trr := tar.NewReader(gzr)
	local := path.Join(tmpdir, "local")
	err = untar(trr, fs, local)
	if err != nil {
		if errors.Is(err, ErrDirEmpty) {
			return []packages.Package{}, nil
		}
		return nil, err
	}

	if err = createALPMDBVersionFile(path.Join(local, ALPMDBVersionFile), strconv.Itoa(ALPMDBVersion)); err != nil {
		return nil, errors.Wrap(err, "error creating ALPM DB version file")
	}

	h, err := alpm.Initialize(root, tmpdir)
	if err != nil {
		return nil, errors.Wrap(err, "error initializing Arch Linux Package Manager handler")
	}
	defer h.Release()

	localDb, err := h.LocalDB()
	if err != nil {
		return nil, err
	}

	var packageList []alpm.IPackage
	for _, v := range packageNames {
		list := localDb.Search([]string{v}).Slice()
		packageList = append(packageList, list...)
	}
	os.Remove(tmpdir)

	ps := []packages.Package{}
	for _, p := range packageList {
		ps = append(ps, &Package{
			Name:         p.Name(),
			Version:      p.Version(),
			Release:      "",
			Architecture: p.Architecture(),
			Location:     p.FileName(),
			url:          p.URL(),
			fileReaders:  nil,
		})
	}

	return ps, nil
}

func untar(source *tar.Reader, fs afero.Fs, target string) error {
	err := fs.MkdirAll(target, 0755)
	if err != nil {
		return errors.Wrap(err, "creating the target directory")
	}

	for {
		header, err := source.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		if header == nil {
			continue
		}

		target := path.Join(target, header.Name)

		switch header.Typeflag {

		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return errors.Wrap(err, "error on untar creating directory")
				}
			}

		case tar.TypeReg:
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return errors.Wrap(err, "error on untar creating file")
			}
			if _, err := io.Copy(f, source); err != nil {
				return errors.Wrap(err, "error on untar copying file content")
			}
			if err = f.Close(); err != nil {
				return errors.Wrap(err, "error on untar closing file")
			}

		default:
			return fmt.Errorf("error on untar: file %s type not supported", header.Name)
		}
	}

	d, err := fs.Open(target)
	if _, err = d.Readdirnames(1); err != nil {
		if err == io.EOF {
			return ErrDirEmpty
		}
	}
	defer d.Close()

	return nil
}

func createALPMDBVersionFile(filename, version string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	if _, err = io.Copy(f, bytes.NewBuffer([]byte(version))); err != nil {
		return err
	}

	return nil
}
