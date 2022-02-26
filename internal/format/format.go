package format

import (
	"bufio"
	"encoding/json"

	"gopkg.in/yaml.v2"

	"github.com/olekukonko/tablewriter"
)

type FormatType string

const (
	Text FormatType = "text"
	Json FormatType = "json"
	Yaml FormatType = "yaml"
)

func Encode(output *bufio.Writer,
	objects interface{},
	format FormatType) (*bufio.Writer, error) {
	switch format {
	case Json:
		return encodeJson(output, objects)
	case Text:
		return encodeText(output, objects)
	case Yaml:
		return encodeYaml(output, objects)
	default:
		return encodeText(output, objects)
	}
}

func encodeJson(output *bufio.Writer, objects interface{}) (*bufio.Writer, error) {
	json, err := json.Marshal(objects)
	if err != nil {
		return nil, err
	}

	_, err = output.Write(json)
	if err != nil {
		return nil, err
	}

	return output, nil
}

func encodeYaml(output *bufio.Writer, objects interface{}) (*bufio.Writer, error) {
	yaml, err := yaml.Marshal(objects)
	if err != nil {
		return nil, err
	}

	_, err = output.Write(yaml)
	if err != nil {
		return nil, err
	}

	return output, nil
}

func encodeText(output *bufio.Writer, objects interface{}) (*bufio.Writer, error) {
	return encodeTableFromStructs(output, objects)
}

func encodeTableFromStructs(output *bufio.Writer, objects interface{}) (*bufio.Writer, error) {
	printer := tablewriter.NewWriter(output)

	printer.SetStructs(objects)
	printer.Render()

	return output, nil
}
