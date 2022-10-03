package kernelrelease

import (
	"bufio"
	"io"
	"strings"
	"unicode"

	p "github.com/maxgio92/krawler/pkg/packages"
)

const (
	ConfigCompilerVersion = "CONFIG_GCC_VERSION"
)

func GetCompilerVersionFromKernelPackage(pkg p.Package) (string, error) {
	return getCompilerVersionFromFileReaders(pkg.FileReaders())
}

func getCompilerVersionFromFileReaders(files []io.Reader) (string, error) {
	for _, r := range files {
		fileScanner := bufio.NewScanner(r)
		fileScanner.Split(bufio.ScanLines)

		for fileScanner.Scan() {
			line := fileScanner.Text()
			if strings.Contains(line, ConfigCompilerVersion) {

				compilerVersion, err := parseConfig(line, ConfigCompilerVersion)
				if err == nil {
					return compilerVersion, nil
				}
				return "", err
			}
		}
		err := fileScanner.Err()
		if err != nil {
			return "", err
		}
	}

	return "", ErrKernelCompilerVersionNotFound
}

func parseConfig(line string, key string) (string, error) {
	tokens := strings.FieldsFunc(line, func(c rune) bool {
		return unicode.Is(unicode.Space, c) || unicode.Is(unicode.Sm, c)
	})
	if len(tokens) > 1 {
		return tokens[len(tokens)-1], nil
	}

	return "", ErrKernelConfigValueNotFound
}
