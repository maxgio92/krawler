package kernelrelease

import (
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/sassoftware/go-rpmutils"
)

const (
	CONFIG_GCC_VERSION = "CONFIG_GCC_VERSION"
	CONFIG_DIR         = ".config"
)

// TODO: Provide abstraction like GetCompilerVersionFromPackageArchive
func GetCompilerVersionFromRPMPackageURL(packageURL string) (gccVersion string, err error) {
	bufferSize := 64

	u, err := url.Parse(packageURL)
	if err != nil {
		return
	}

	resp, err := http.Get(u.String())
	if err != nil {
		return
	}

	rpm, err := rpmutils.ReadRpm(resp.Body)
	if err != nil {
		panic(err)
	}

	payload, err := rpm.PayloadReaderExtended()

	for {
		fileInfo, err := payload.Next()
		if err == io.EOF {
			break
		}

		if filepath.Base(fileInfo.Name()) == CONFIG_DIR {

			for {
				bytes := make([]byte, bufferSize)
				_, err := payload.Read(bytes)

				if strings.Contains(string(bytes), CONFIG_GCC_VERSION) {
					lines := strings.Split(string(bytes), "\n")

					for _, l := range lines {
						if strings.Contains(l, CONFIG_GCC_VERSION) {
							gccVersion = strings.Split(l, "=")[1]
						}
					}
				}
				if err == io.EOF {
					break
				}
			}

			goto exit
		}
	}
exit:
	return
}
