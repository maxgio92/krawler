package format

import (
	"strings"
	"encoding/json"

	"github.com/falcosecurity/driverkit/pkg/kernelrelease"
	"gopkg.in/yaml.v2"
)

type FormatType string

const (
	Text FormatType = "text"
	Json FormatType = "json"
	Yaml FormatType = "yaml"
)

func Encode(kernelReleases []kernelrelease.KernelRelease, format FormatType) (string, error) {
	switch format {
	case Json:
		return encodeJson(kernelReleases)
	case Text:
		return encodeText(kernelReleases), nil
	case Yaml:
		return encodeYaml(kernelReleases)
	default:
		return encodeText(kernelReleases), nil
	}
}

func encodeText(kernelReleases []kernelrelease.KernelRelease) string {
	var ss []string
	for _, v := range kernelReleases {
		ss = append(ss, v.Fullversion+v.FullExtraversion)
	}
	raw := strings.Join(ss, " ")
	return raw
}

func encodeJson(kernelReleases []kernelrelease.KernelRelease) (string, error) {
	json, err := json.Marshal(kernelReleases)
	if err != nil {
		return "", err
	}
	return string(json), nil
}

func encodeYaml(kernelReleases []kernelrelease.KernelRelease) (string, error) {
	yaml, err := yaml.Marshal(kernelReleases)
	if err != nil {
		return "", err
	}
	return string(yaml), nil
}
