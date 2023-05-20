package alpm

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"github.com/Jguer/go-alpm/v2"
	"github.com/maxgio92/krawler/pkg/packages"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
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

// SearchPackages looks for the package of which the specified package names, parsing the remote
// repository DB, and returns a slice of packages.Package.
// It possibly returns an error.
func SearchPackages(dbURL string, packageNames []string) ([]packages.Package, error) {
	tmpdir, err := os.MkdirTemp(os.TempDir(), "*")
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
	err = untar(trr, local)
	if err != nil {
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

	packageList := localDb.Search(packageNames).Slice()

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

func untar(source *tar.Reader, target string) error {
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
